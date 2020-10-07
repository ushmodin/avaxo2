package util

import (
	"log"
	"net"
	"time"

	"github.com/gorilla/websocket"
	"github.com/ushmodin/avaxo2/internal/model"
)

func ForwardLocalTraffic(wsConn *websocket.Conn, localConn net.Conn) {
	defer localConn.Close()
	defer wsConn.Close()

	buffer := make([]byte, 1024)
	for {
		buffSize, err := localConn.Read(buffer)
		if err != nil {
			log.Printf("Error while read from local connection: %v\n", err)
			return
		}
		if buffSize == 0 {
			time.Sleep(100)
			continue
		}

		err = wsConn.WriteJSON(model.ForwardPacket{
			Type: model.ForwardBytes,
			Body: model.ForwardPacketBody{
				Bytes: buffer[:buffSize],
			},
		})
		if err != nil {
			log.Printf("Error while writer to web socket: %v\n", err)
			return
		}
	}

}

func ForwardWebsocketTraffic(localConn net.Conn, wsConn *websocket.Conn) {
	var data model.ForwardPacket
	defer localConn.Close()
	defer wsConn.Close()

	for {
		if err := wsConn.ReadJSON(&data); err != nil {
			log.Printf("Error while read packet from websocket: %v\n", err)
		}
		if data.Type != model.ForwardBytes {
			log.Printf("Incorrect packet from websocket. Type: %d\n", data.Type)
			return
		}
		_, err := localConn.Write(data.Body.Bytes)
		if err != nil {
			log.Printf("Error while writer to local connection: %v\n", err)
			return
		}
	}
}
