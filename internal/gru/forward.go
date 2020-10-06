package gru

import (
	"errors"

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
