package model

import (
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

func (s *Streamer) listenToWs(){
	defer func(){
		s.WsConn.Close()
	}()
	for{
		msgType, msg, err := s.WsConn.ReadMessage()
		if e, ok :=  err.(*websocket.CloseError); ok && 
		(e.Code == websocket.CloseNormalClosure || e.Code == websocket.CloseNoStatusReceived) {
			log.Println("Error in reading message for streamer : ", err)
			break
		}
		log.Println("From : " + s.UserId + ", Message type: " + strconv.Itoa(msgType) + ", this is msg : " + string(msg))
	}
}


type Viewer struct {
	UserId string
	WsConn *websocket.Conn
	Stream *Stream
}

func (v *Viewer) listenToWs(){
	defer func(){
		v.WsConn.Close()
	}()
	for{
		msgType, msg, err := v.WsConn.ReadMessage()
		if e, ok :=  err.(*websocket.CloseError); ok && 
		(e.Code == websocket.CloseNormalClosure || e.Code == websocket.CloseNoStatusReceived) {
			log.Println("Error in reading message for Viewer : ", err)
			break
		}
		log.Println("From : " + v.UserId + ", Message type: " + strconv.Itoa(msgType) + ", this is msg : " + string(msg))
	}
}


type Stream struct {
	Id     string
	Streamer *Streamer
	Viewers *map[string] *Viewer
	*sync.RWMutex
}
