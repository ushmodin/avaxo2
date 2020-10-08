package minion

import (
	"errors"
	"fmt"
	"net"

	"github.com/gorilla/websocket"
	"github.com/ushmodin/avaxo2/internal/model"
)

func forwardInit(conn *websocket.Conn) error {
	var rq model.ForwardPacket

	err := conn.ReadJSON(&rq)
	if err != nil {
		return err
	}

	if rq.Type != model.ForwardInit {
		return errors.New("Illegal init packet")
	}

	return conn.WriteJSON(model.ForwardPacket{
		Type: model.ForwardOK,
	})
}

func forwardGetTarget(conn *websocket.Conn) (string, error) {
	var rq model.ForwardPacket

	err := conn.ReadJSON(&rq)
	if err != nil {
		return "", err
	}

	if rq.Type != model.ForwardTarget {
		return "", errors.New("Illegal init packet")
	}

	err = conn.WriteJSON(model.ForwardPacket{
		Type: model.ForwardOK,
	})
	if err != nil {
		return "", err
	}
	return rq.Body.Str, nil
}

func forwardConnect(wsConn *websocket.Conn, target string) (net.Conn, error) {
	var rq model.ForwardPacket

	err := wsConn.ReadJSON(&rq)
	if err != nil {
		return nil, err
	}

	if rq.Type != model.ForwardConnect {
		return nil, errors.New("Illegal init packet")
	}

	tcpConn, err := net.Dial("tcp", target)
	if err != nil {
		wsConn.WriteJSON(model.ForwardPacket{
			Type: model.ForwardError,
			Body: model.ForwardPacketBody{
				Str: fmt.Sprintf("Connect error: %s", err),
			},
		})
		if err != nil {
			return nil, err
		}

	}
	err = wsConn.WriteJSON(model.ForwardPacket{
		Type: model.ForwardOK,
	})
	if err != nil {
		return nil, err
	}
	return tcpConn, nil
}
