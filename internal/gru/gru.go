package gru

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/ushmodin/avaxo2/internal/model"
	"github.com/ushmodin/avaxo2/internal/settings"
	"github.com/ushmodin/avaxo2/internal/util"
)

type Gru struct {
	httpClient *http.Client
}

func NewGru(certfile, keyfile, cafile string) (*Gru, error) {
	tls, err := util.TLSConfig(certfile, keyfile, cafile)
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tls,
		},
	}

	return &Gru{
		httpClient: client,
	}, nil
}

func (gru *Gru) Ls(minion, path string) ([]model.DirItem, error) {
	host, err := getMinionHost(minion)
	if err != nil {
		return nil, err
	}

	u := &url.URL{
		Scheme: "https",
		Host:   host,
		Path:   "/api/ls",
		RawQuery: (&url.Values{
			"path": {path},
			"fmt":  {"json"},
		}).Encode(),
	}
	rsp, err := gru.httpClient.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()
	var res []model.DirItem
	err = json.NewDecoder(rsp.Body).Decode(&res)
	return res, err
}

// GetFile get file from minion by path
func (gru *Gru) GetFile(minion, path string) (io.ReadCloser, error) {
	host, err := getMinionHost(minion)
	if err != nil {
		return nil, err
	}

	u := &url.URL{
		Scheme: "https",
		Host:   host,
		Path:   "/api/file/get",
		RawQuery: (&url.Values{
			"path": {path},
		}).Encode(),
	}
	rsp, err := gru.httpClient.Get(u.String())
	if err != nil {
		return nil, err
	}
	return rsp.Body, nil
}

// PutFile send file to minion
func (gru *Gru) PutFile(minion, path string, r io.Reader) error {
	host, err := getMinionHost(minion)
	if err != nil {
		return err
	}

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
	_, err = gru.httpClient.Do(req)
	return err
}

func getMinionHost(val string) (string, error) {
	u, err := url.ParseRequestURI(val)
	if err == nil {
		return u.Host + ":" + u.Port(), nil
	}

	addr, err := settings.GetMinionAddress(val)

	if err == nil {
		return addr.Host, nil
	}
	return "", errors.New("Minion not found")
}

// Exec execute command cmd with args on minion
func (gru *Gru) Exec(cmd string, args []string, nowait bool, timeout int) error {
	return nil
}
