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

		HandleWsMessage(msg, s.Stream, false)
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

		HandleWsMessage(msg, v.Stream, true)
	}
}

func HandleWsMessage(msg []byte, stream *Stream, fromViewer bool) {
	jsonMap := make(map[string]json.RawMessage)
	err := json.Unmarshal(msg, &jsonMap)
	log.Println("JSON MAP FROM HANDLE MESSAGE IS : ", jsonMap)
	if err != nil {
		log.Println("error in unmarshal : ", err)
		return
	}

	log.Println("Type received from message : ", string(jsonMap["type"]))
	receivedMsgType, _ := strconv.Atoi(string(jsonMap["type"]))

	if receivedMsgType == 1 {
		stream.SendChatMessage(msg)
	} else if receivedMsgType == 0 {

		var receivedConnectionMessage ConnectionMessage
		err := json.Unmarshal(jsonMap["data"], &receivedConnectionMessage)
		if err != nil {
			log.Println("error in json parse for connection message: ", err)
			return
		}

		if fromViewer {
			stream.SendSdpToStreamer(&receivedConnectionMessage)
		} else {
			stream.SendSdpToViewer(&receivedConnectionMessage)
		}
	}
}

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
