package agent

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/ping", pingHandler)
	r.HandleFunc("/api/ls", lsHandler)
	return r
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}
