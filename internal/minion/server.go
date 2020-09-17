package minion

import (
	"net/http"

	"github.com/ushmodin/avaxo2/internal/util"
)

type Server struct {
	httpServer *http.Server
}

// NewServer create new Agent server
func NewServer(listen, keyfile, certfile, cafile string) (*Server, error) {
	tls, err := util.TLSConfig(certfile, keyfile, cafile)
	if err != nil {
		return nil, err
	}

	s := &http.Server{
		Handler:   NewMinionRoute(NewMinion()),
		Addr:      listen,
		TLSConfig: tls,
	}

	return &Server{
		httpServer: s,
	}, nil
}

// Run start http server
func (srv *Server) Run() error {
	return srv.httpServer.ListenAndServeTLS("", "")
}
