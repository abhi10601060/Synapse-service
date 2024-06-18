package model

import (
	"encoding/json"
	"log"
	"strconv"
	"github.com/gorilla/websocket"
)

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

		WsMessageHandler(msg, v.Stream, true)
	}
}

func WsMessageHandler(msg []byte, stream *Stream, fromViewer bool) {
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

