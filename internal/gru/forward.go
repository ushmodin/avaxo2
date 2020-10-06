package gru

import (
	"errors"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/ushmodin/avaxo2/internal/model"
)

func forwardInit(conn *websocket.Conn) error {
	rq := model.ForwardPacket{
		Type: model.ForwardInit,
		Body: model.ForwardPacketBody{
			Str: "1.0",
		},
	}

	if err := conn.WriteJSON(rq); err != nil {
		return err
	}

	var rs model.ForwardPacket
	if err := conn.ReadJSON(&rs); err != nil {
		return err
	}
	if rs.Type != model.ForwardOK {
		return errors.New("Illegal response for init packet")
	}
	return nil
}

func forwardSendTarget(conn *websocket.Conn, target string) error {
	rq := model.ForwardPacket{
		Type: model.ForwardTarget,
		Body: model.ForwardPacketBody{
			Str: target,
		},
	}

	if err := conn.WriteJSON(rq); err != nil {
		return err
	}

	var rs model.ForwardPacket
	if err := conn.ReadJSON(&rs); err != nil {
		return err
	}
	if rs.Type != model.ForwardOK {
		return errors.New("Illegal response for init packet")
	}
	return nil
}

func forwardConnect(conn *websocket.Conn) error {
	rq := model.ForwardPacket{
		Type: model.ForwardConnect,
	}

	if err := conn.WriteJSON(rq); err != nil {
		return err
	}

	var rs model.ForwardPacket
	if err := conn.ReadJSON(&rs); err != nil {
		return err
	}
	if rs.Type != model.ForwardOK {
		if rs.Type == model.ForwardError {
			return fmt.Errorf("Error while minion connect to target: %s", rs.Body.Str)
		}
		return errors.New("Illegal response for connect packet")
	}
	return nil
}
