package model


import (
	"encoding/json"
	"log"
	"sync"
	"github.com/gorilla/websocket"
)

type Stream struct {
	Id       string
	Streamer *Streamer
	Viewers  map[string]*Viewer
	sync.RWMutex
}

func (s *Stream) SendChatMessage(msg []byte) {

	log.Println("Strated writting message to all viewers...")

	err := s.Streamer.WsConn.WriteMessage(1, msg)
	if e, ok := err.(*websocket.CloseError); ok &&
		(e.Code == websocket.CloseNormalClosure || e.Code == websocket.CloseNoStatusReceived) {
		log.Println("Error in sending message to streamer : ", err)
		return
	}

	for _, viewer := range s.Viewers {
		err := viewer.WsConn.WriteMessage(1, msg)
		if e, ok := err.(*websocket.CloseError); ok &&
			(e.Code == websocket.CloseNormalClosure || e.Code == websocket.CloseNoStatusReceived) {
			log.Println("Error in sending message to streamer : ", err)
		}
	}
}

func (s *Stream) SendSdpToViewer(connectionMessage *ConnectionMessage) {
	viewerId := connectionMessage.ViewerUserId

	viewer, isExist := s.Viewers[viewerId]
	if !isExist {
		log.Println("SendSdpToViewer viewer not exist...")
		return
	}
	msgJson, err := json.Marshal(connectionMessage)
	if err != nil {
		log.Println("SendSdpToViewer json marshal error for connectionMessage : ", err)
		return
	}
	err = viewer.WsConn.WriteMessage(1, msgJson)
	if err != nil {
		log.Println("SendSdpToViewer error in writeMessage...")
	}
	log.Println("SendSdpToViewer sdp sent to Viewer successfully...")
}

func (s *Stream) SendSdpToStreamer(connectionMessage *ConnectionMessage) {
	msgJson, err := json.Marshal(connectionMessage)
	if err != nil {
		log.Println("SendSdpToStreamer json marshal error for connectionMessage : ", err)
		return
	}
	err = s.Streamer.WsConn.WriteMessage(1, msgJson)
	if err != nil {
		log.Println("SendSdpToStreamer error in writeMessage...")
	}
	log.Println("SendSdpToStreamer sdp sent to Streamer successfully...")
}
