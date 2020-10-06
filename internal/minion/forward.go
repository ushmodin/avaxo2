package minion

import (
	"errors"

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
