package agent

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type lsReponseFormat int

const (
	formatPlain = 0
	formatJSON  = 1
)

type lsData struct {
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	Modified string `json:"modified"`
	IsDir    bool   `json:"isDir"`
}

func lsHandler(w http.ResponseWriter, r *http.Request) {
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

	files, err := ls(path)
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

func ls(path string) ([]lsData, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	res := make([]lsData, len(files))
	for i, v := range files {
		res[i].Name = v.Name()
		res[i].Size = v.Size()
		res[i].Modified = v.ModTime().Format(time.RFC822)
		res[i].IsDir = v.IsDir()
	}
	return res, nil
}

func formatFiles(files []lsData) []byte {
	b := bytes.NewBuffer([]byte{})
	b.WriteString(fmt.Sprintf("total %d\n", len(files)))
	for _, f := range files {
		var t string = "f"
		var s = f.Size
		if f.IsDir {
			t = "d"
			s = 0
		}
		b.WriteString(fmt.Sprintf("%s %15d %s %s\n", t, s, f.Modified, f.Name))
	}
	return b.Bytes()
}
