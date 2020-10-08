package util

import (
	"io"
	"log"
	"net"
	"time"

	"github.com/gorilla/websocket"
)

func ForwardLocalTraffic(wsConn *websocket.Conn, localConn net.Conn) {
	buffer := make([]byte, 1024)
	defer wsConn.Close()
	for {
		n, err := localConn.Read(buffer)
		if err == io.EOF || err == io.ErrClosedPipe {
			break
		} else if err != nil {
			log.Printf("Error while read from local connection: %v\n", err)
			break
		}
		if n == 0 {
			time.Sleep(100)
			continue
		}
		w, err := wsConn.NextWriter(websocket.BinaryMessage)
		if err != nil {
			log.Printf("Can't open websocket writer: %v\n", err)
			return
		}
		_, err = w.Write(buffer[0:n])
		if err != nil {
			log.Printf("Error while write to web socket: %v\n", err)
			w.Close()
			break
		}
		w.Close()
	}

}

func ForwardWebsocketTraffic(localConn net.Conn, wsConn *websocket.Conn) {
	defer localConn.Close()

	for {
		mt, reader, err := wsConn.NextReader()
		if err == io.EOF || err == io.ErrClosedPipe || err == io.ErrUnexpectedEOF {
			break
		}
		if mt == websocket.PingMessage {
			log.Println("Ping was received")
			continue
		}
		if err != nil {
			log.Printf("Can't open websocket reader: %v\n", err)
			break
		}
		if mt != websocket.BinaryMessage {
			log.Printf("incorrect websocket reader message type\n")
			break
		}
		// for n, err := io.Copy(localConn, reader); err != nil; {
		// 	if n == 0 {
		// 		time.Sleep(100)
		// 	}
		if n, err := io.Copy(localConn, reader); err != nil {
			break
		} else if n == 0 {
			time.Sleep(100)
		}
	}

}

func SendPings(wsConn *websocket.Conn) {
	for {
		if err := wsConn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
			break
		}
		log.Println("Ping was send")
		<-time.Tick(30 * time.Second)
	}
}
