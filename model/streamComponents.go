package model

import (
	"encoding/json"
	"log"
	"strconv"
	"sync"

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

		jsonMap := make(map[string] json.RawMessage)
		err = json.Unmarshal(msg, &jsonMap)
		if err != nil {
			log.Println("error in unmarshal : ", err)
		}
	}
}

type Viewer struct {
	UserId string
	WsConn *websocket.Conn
	Stream *Stream
}

func (v *Viewer) ListenToWs() {
	defer func() {
		v.WsConn.Close()
	}()
	for {
		msgType, msg, err := v.WsConn.ReadMessage()
		if e, ok := err.(*websocket.CloseError); ok &&
			(e.Code == websocket.CloseNormalClosure || e.Code == websocket.CloseNoStatusReceived) {
			log.Println("Error in reading message for Viewer : ", err)
			break
		}
		log.Println("From Viewer : " + v.UserId + ", Message type: " + strconv.Itoa(msgType) + ", this is msg : " + string(msg))

		jsonMap := make(map[string] json.RawMessage)
		err = json.Unmarshal(msg, &jsonMap)
		if err != nil {
			log.Println("error in unmarshal : ", err)
		}
	}
}

type Stream struct {
	Id       string
	Streamer *Streamer
	Viewers  map[string]*Viewer
	sync.RWMutex
}
