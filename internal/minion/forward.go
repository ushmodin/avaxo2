package minion

import (
	"errors"
	"fmt"
	"log"
	"net"
	"time"

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
		err = wsConn.WriteJSON(model.ForwardPacket{
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

func forwardLocalTraffic(wsConn *websocket.Conn, localConn net.Conn) {
	defer localConn.Close()
	defer wsConn.Close()

	for {
		buffer := make([]byte, 1024)
		buffSize, err := localConn.Read(buffer)
		if err != nil {
			log.Printf("Error while read from local connection: %v\n", err)
			return
		}
		if buffSize == 0 {
			time.Sleep(100)
			continue
		}

		log.Printf("%d bytes was read from local\n", buffSize)
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
		log.Printf("bytes packet was send to websocket\n")
	}

}

func forwardWebsocketTraffic(localConn net.Conn, wsConn *websocket.Conn) {
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
		log.Printf("bytes packet was read from websocket\n")
		n, err := localConn.Write(data.Body.Bytes)
		if err != nil {
			log.Printf("Error while writer to local connection: %v\n", err)
			return
		}
		log.Printf("%d bytes was write to local connection", n)
	}
}
