package handler

import (
	"log"
	"net/http"
	"synapse/stream/auth"
	"synapse/stream/model"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var (
	Manager = model.StreamManager{
		Streams: map[string]*model.Stream{},
	}

	websocketupgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
)

func StartStream(c *gin.Context) {
	userId := auth.GetUserIdFromToken(c.Request.Header.Get("Authentication-Token"))
	if userId == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Couldn't parse UserId from Token",
		})
		c.Abort()
		return
	}
	log.Println("received id : ", userId)

	wsConn, err := websocketupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("err during websocket upgrade in start Stream : ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Couldn't Upgrade websocket connection...",
		})
		c.Abort()
		return
	}

	streamer := model.Streamer{
		UserId: userId,
		WsConn: wsConn,
	}
	go streamer.ListenToWs()

	Manager.StartStream(&streamer)
}

func JoinStream(c *gin.Context) {
	userId := auth.GetUserIdFromToken(c.Request.Header.Get("Authentication-Token"))
	if userId == "" {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Couldn't parse UserId from Token",
		})
		c.Abort()
		return
	}
	log.Println("received id : ", userId)

	pid := c.Param("pid")
	log.Println("Received Pid from viewer is : ", pid)

	if isExist := Manager.IsStreamExist(pid); !isExist {
		c.JSON(404, gin.H{
			"message": "Required stream does not exist",
		})
		c.Abort()
		return
	}

	wsConn, err := websocketupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("err during websocket upgrade in start Stream : ", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Couldn't Upgrade websocket connection...",
		})
		c.Abort()
		return
	}

	viewer := model.Viewer{
		UserId: userId,
		WsConn: wsConn,
	}

	isAdded := Manager.AddViewerToStream(pid, &viewer)

	if !isAdded {
		wsConn.Close()
		c.Abort()
		return
	}

	go viewer.ListenToWs()
}
