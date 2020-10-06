package gru

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/ushmodin/avaxo2/internal/model"
	"github.com/ushmodin/avaxo2/internal/util"
)

type Port struct {
	httpClient *http.Client
	dialer     *websocket.Dialer
}

func NewPort(certfile, keyfile, cafile string) (*Port, error) {
	tls, err := util.TLSConfig(certfile, keyfile, cafile)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tls,
		},
	}

	dialer := &websocket.Dialer{
		TLSClientConfig: tls,
	}

	return &Port{
		httpClient: client,
		dialer:     dialer,
	}, nil
}

func (port *Port) Ls(host, path string) ([]model.DirItem, error) {
	u := &url.URL{
		Scheme: "https",
		Host:   host,
		Path:   "/api/ls",
		RawQuery: (&url.Values{
			"path": {path},
			"fmt":  {"json"},
		}).Encode(),
	}
	rsp, err := port.httpClient.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	var res []model.DirItem
	err = json.NewDecoder(rsp.Body).Decode(&res)
	return res, err
}

// GetFile get file from minion by path
func (port *Port) GetFile(host, path string) (io.ReadCloser, error) {
	u := &url.URL{
		Scheme: "https",
		Host:   host,
		Path:   "/api/file/get",
		RawQuery: (&url.Values{
			"path": {path},
		}).Encode(),
	}
	rsp, err := port.httpClient.Get(u.String())
	if err != nil {
		return nil, err
	}
	return rsp.Body, nil
}

// PutFile send file to minion
func (port *Port) PutFile(host, path string, r io.Reader) error {

	u := &url.URL{
		Scheme: "https",
		Host:   host,
		Path:   "/api/file/put",
		RawQuery: (&url.Values{
			"path": {path},
		}).Encode(),
	}

	req, err := http.NewRequest("PUT", u.String(), r)
	if err != nil {
		return err
	}
	_, err = port.httpClient.Do(req)
	return err
}

// Exec execute command cmd with args on minion
func (port *Port) Exec(host string, cmd string, args []string) (string, error) {
	u := &url.URL{
		Scheme: "https",
		Host:   host,
		Path:   "/api/proc/exec",
	}
	rq := model.ExecRq{
		Cmd:  cmd,
		Args: args,
	}
	body, err := json.Marshal(rq)
	if err != nil {
		return "", err
	}

	rsp, err := port.httpClient.Post(u.String(), "application/json", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	var res model.ExecRs
	err = json.NewDecoder(rsp.Body).Decode(&res)
	return res.ProcID, nil
}

func (port *Port) ProcTail(host string, procID string) (io.ReadCloser, error) {
	u := &url.URL{
		Scheme: "https",
		Host:   host,
		Path:   "/api/proc/" + procID + "/tail",
	}
	rsp, err := port.httpClient.Get(u.String())
	if err != nil {
		return nil, err
	}
	return rsp.Body, nil
}

func (port *Port) WsForward(host string) (*websocket.Conn, error) {
	u := &url.URL{
		Scheme: "wss",
		Host:   host,
		Path:   "/api/forward",
	}
	c, _, err := port.dialer.Dial(u.String(), nil)
	return c, err
}
