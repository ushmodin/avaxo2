package util

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
)

// TLSConfig create tls config for http server
func TLSConfig(certfile, keyfile, cafile string) (*tls.Config, error) {
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
		Certificates:       []tls.Certificate{cert},
		ClientCAs:          caCertPool,
		ClientAuth:         tls.RequireAndVerifyClientCert,
		InsecureSkipVerify: true,
	}, nil

}
