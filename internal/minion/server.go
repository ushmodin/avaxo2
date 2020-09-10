package minion

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

// NewServer create new Agent server
func NewServer(listen, keyfile, certfile, cafile string) (*Server, error) {
	tls, err := tlsConfig(certfile, keyfile, cafile)
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
