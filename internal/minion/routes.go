package minion

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

type minionHTTPWrapper struct {
	minion *Minion
}

type lsReponseFormat int

const (
	formatPlain = 0
	formatJSON  = 1
)

func NewMinionRoute(minion *Minion) http.Handler {
	wrapper := &minionHTTPWrapper{
		minion,
	}

	handler := mux.NewRouter()
	handler.HandleFunc("/api/ping", pingHandler)
	handler.HandleFunc("/api/ls", wrapper.lsHandler)
	handler.HandleFunc("/api/file/get", wrapper.getHandler)
	handler.HandleFunc("/api/file/put", wrapper.putHandler)
	handler.HandleFunc("/api/proc/exec", wrapper.procExecHandler)
	handler.HandleFunc("/api/proc/{id}/info", wrapper.procInfoHandler)
	handler.HandleFunc("/api/proc/{id}/kill", wrapper.procKillHandler)
	handler.HandleFunc("/api/proc/ps", wrapper.procPsHandler)

	return handler
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}

func (wrapper *minionHTTPWrapper) lsHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	if paths := q["path"]; len(paths) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("parameter path not specified"))
		return
	}

	path := q["path"][0]

	stat, err := os.Stat(path)
	if os.IsNotExist(err) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("path not found"))
		return
	}

	if !stat.IsDir() {
		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte("path not directory"))
		return
	}

	format := formatJSON
	if formats := q["fmt"]; len(formats) > 0 {
		switch formats[0] {
		case "json":
			format = formatJSON
		case "plain":
			format = formatPlain
		default:
			w.WriteHeader(http.StatusNotAcceptable)
			w.Write([]byte("illegal parameter fmt"))
			return
		}
	}

	files, err := wrapper.minion.ReadDir(path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error while read path. " + err.Error()))
		return
	}

	if format == formatJSON {
		resp, err := json.Marshal(files)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/json; encoding=UTF-8")
		w.Write(resp)
	} else if format == formatPlain {
		resp := formatFiles(files)
		w.Header().Set("Content-Type", "text/plain; encoding=UTF-8")
		w.Write(resp)
	}

}

func formatFiles(files []DirItem) []byte {
	b := bytes.NewBuffer([]byte{})
	b.WriteString(fmt.Sprintf("total %d\n", len(files)))
	for _, f := range files {
		var t string = "f"
		var s = f.Size
		if f.IsDir {
			t = "d"
			s = 0
		} else if f.Error != "" {
			t = "e"
			s = 0
		}
		if f.Error != "" {
			b.WriteString(fmt.Sprintf("%s %15d %30s %s(%s)\n", t, s, "", f.Name, f.Error))
		} else {
			b.WriteString(fmt.Sprintf("%s %15d %30s %s\n", t, s, f.Modified, f.Name))
		}
	}
	return b.Bytes()
}

func (wrapper *minionHTTPWrapper) getHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	if paths := q["path"]; len(paths) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("parameter path not specified"))
		return
	}

	path := q["path"][0]

	f, err := wrapper.minion.GetFile(path)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error while read path. " + err.Error()))
		return
	}
	defer f.Close()

	if _, err := io.Copy(w, f); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error while read path. " + err.Error()))
	}
}

func (wrapper *minionHTTPWrapper) putHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Only PUT supported"))
		return
	}

	q := r.URL.Query()

	if paths := q["path"]; len(paths) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("parameter path not specified"))
		return
	}

	var mode os.FileMode = 0644
	if modes := q["mode"]; len(modes) != 0 {
		if m, err := strconv.Atoi(modes[0]); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalide parameter mode"))
		} else {
			mode = os.FileMode(m)
		}

		return
	}

	path := q["path"][0]

	err := wrapper.minion.PutFile(path, mode, r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("error while read path. " + err.Error()))
		return
	}
}

func (wrapper *minionHTTPWrapper) procExecHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Only POST supported"))
		return
	}

	rq := struct {
		Cmd  string   `json:"cmd"`
		Args []string `json:"args"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&rq); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error while parse request body. " + err.Error()))
		return
	}

	if len(rq.Cmd) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Illegal json format. 'cmd' not found"))
		return
	}

	id, err := wrapper.minion.Exec(rq.Cmd, rq.Args...)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error while execute command. " + err.Error()))
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"procId": id,
	})
	w.Header().Set("Content-Type", "application/json")
}

func (wrapper *minionHTTPWrapper) procInfoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	info, err := wrapper.minion.ProcInfo(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	json.NewEncoder(w).Encode(info)
	w.Header().Set("Content-Type", "application/json")
}

func (wrapper *minionHTTPWrapper) procKillHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	err := wrapper.minion.ProcKill(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (wrapper *minionHTTPWrapper) procPsHandler(w http.ResponseWriter, r *http.Request) {
	ps := wrapper.minion.ProcPs()
	json.NewEncoder(w).Encode(ps)
	w.Header().Set("Content-Type", "application/json")
}