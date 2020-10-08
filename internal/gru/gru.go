package gru

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/url"
	"os"
	"time"

	"github.com/ushmodin/avaxo2/internal/model"
	"github.com/ushmodin/avaxo2/internal/settings"
	"github.com/ushmodin/avaxo2/internal/util"
)

type Gru struct {
	port *Port
}

func NewGru(certfile, keyfile, cafile string) (*Gru, error) {
	port, err := NewPort(certfile, keyfile, cafile)
	if err != nil {
		return nil, err
	}
	return &Gru{
		port: port,
	}, nil
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

func (gru *Gru) Ls(minion, path string, jsonFormat bool) error {
	host, err := getMinionHost(minion)
	if err != nil {
		return err
	}

	files, err := gru.port.Ls(host, path)
	if err != nil {
		return err
	}

	if jsonFormat {
		out, _ := json.Marshal(files)
		fmt.Println(string(out))
	} else {
		out := model.PrintFiles(files)
		fmt.Println(string(out))
	}
	return nil
}

func (gru *Gru) GetFile(minion, remote, local string) error {
	host, err := getMinionHost(minion)
	if err != nil {
		return err
	}

	var dest io.WriteCloser
	if local != "" {
		dest, err := os.OpenFile(local, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer dest.Close()
	} else {
		dest = os.Stdout
	}

	reader, err := gru.port.GetFile(host, remote)
	if err != nil {
		return err
	}
	defer reader.Close()

	_, err = io.Copy(dest, reader)
	return err
}

func (gru *Gru) PutFile(minion, path, remote string) error {
	host, err := getMinionHost(minion)
	if err != nil {
		return err
	}

	src, err := os.Open(path)
	if err != nil {
		return err
	}
	defer src.Close()

	return gru.port.PutFile(host, remote, src)
}

func (gru *Gru) Exec(minion, cmd string, args []string, nowait bool, timeout int) error {
	host, err := getMinionHost(minion)
	if err != nil {
		return err
	}

	procID, err := gru.port.Exec(host, cmd, args)
	if err != nil {
		return err
	}

	if nowait {
		fmt.Printf("ProcID: %s\n", procID)
		return nil
	}

	reader, err := gru.port.ProcTail(host, procID)
	if err != nil {
		return err
	}
	defer reader.Close()

	if timeout <= 0 {
		io.Copy(os.Stdout, reader)
		return nil
	}

	copyFinished := make(chan error, 1)
	go func() {
		_, err := io.Copy(os.Stdout, reader)
		copyFinished <- err
	}()

	select {
	case err := <-copyFinished:
		return err
	case <-time.Tick(time.Duration(timeout) * time.Second):
		return nil
	}
}

func (gru *Gru) Forward(minion string, port int, target string) error {
	host, err := getMinionHost(minion)
	if err != nil {
		return err
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		return err
	}
	defer listener.Close()
	log.Printf("Listing port %d with target %s on minion %s\n", port, target, minion)
	for {
		localConn, err := listener.Accept()
		if err != nil {
			return err
		}
		go gru.forwardConnection(localConn, host, target)
	}
}

func (gru *Gru) forwardConnection(localConn net.Conn, host string, target string) {
	defer localConn.Close()

	wsConn, err := gru.port.WsForward(host)
	if err != nil {
		log.Printf("Error while connect to minion: %v", err)
		return
	}
	defer wsConn.Close()

	if err := forwardInit(wsConn); err != nil {
		log.Printf("Forward protocol initialization error: %v", err)
		return
	}

	if err := forwardSendTarget(wsConn, target); err != nil {
		log.Printf("Forward protocol initialization error: %v", err)
		return
	}

	if err := forwardConnect(wsConn); err != nil {
		log.Printf("Forward protocol initialization error: %v", err)
		return
	}

	go util.SendPings(wsConn)
	go util.ForwardLocalTraffic(wsConn, localConn)
	util.ForwardWebsocketTraffic(localConn, wsConn)
}
