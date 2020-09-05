package agent

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	router   *mux.Router
	listen   string
	keyfile  string
	certfile string
}

// NewServer create new Agent server
func NewServer(listen, key, cert string) (*Server, error) {
	r := mux.NewRouter()
	r.HandleFunc("/api/ping", pingHandler)
	return &Server{
		router:   r,
		listen:   listen,
		keyfile:  key,
		certfile: cert,
	}, nil
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}

// Run start http server
func (srv *Server) Run() error {
	s := &http.Server{
		Handler: srv.router,
		Addr:    srv.listen,
	}
	return s.ListenAndServeTLS(srv.certfile, srv.keyfile)
}
