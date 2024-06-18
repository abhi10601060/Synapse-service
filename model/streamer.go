package model

import (
	"log"
	"strconv"
	"github.com/gorilla/websocket"
)

type Streamer struct {
	UserId string
	WsConn *websocket.Conn
	Stream *Stream
}

func (s *Streamer) ListenToWs() {
	defer func() {
		s.WsConn.Close()
	}()
	for {
		msgType, msg, err := s.WsConn.ReadMessage()
		if e, ok := err.(*websocket.CloseError); ok &&
			(e.Code == websocket.CloseNormalClosure || e.Code == websocket.CloseNoStatusReceived) {
			log.Println("Error in reading message for streamer : ", err)
			break
		}
		log.Println("From Streamer : " + s.UserId + ", Message type: " + strconv.Itoa(msgType) + ", this is msg : " + string(msg))

		WsMessageHandler(msg, s.Stream, false)
	}
}