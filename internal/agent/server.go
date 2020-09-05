package agent

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	httpServer *http.Server
}

// NewServer create new Agent server
func NewServer(listen, keyfile, certfile, cafile string) (*Server, error) {
	r := mux.NewRouter()
	r.HandleFunc("/api/ping", pingHandler)

	tls, err := tlsConfig(certfile, keyfile, cafile)
	if err != nil {
		return nil, err
	}

	s := &http.Server{
		Handler:   r,
		Addr:      listen,
		TLSConfig: tls,
	}

	return &Server{
		httpServer: s,
	}, nil
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}

// Run start http server
func (srv *Server) Run() error {
	return srv.httpServer.ListenAndServeTLS("", "")
}

// tlsConfig create tls config for http server
func tlsConfig(certfile, keyfile, cafile string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certfile, keyfile)
	if err != nil {
		return nil, err
	}

	caCert, err := ioutil.ReadFile(cafile)
	if err != nil {
		return nil, err
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}, nil

}
