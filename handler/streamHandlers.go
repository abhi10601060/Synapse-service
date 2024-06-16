package handler

import (
	"log"
	"synapse/stream/auth"

	"github.com/gin-gonic/gin"
)

func StartStream(c *gin.Context) {

	id := auth.GetUserIdFromToken(c.Request.Header.Get("Authentication-Token"))
	log.Println("received id : ", id)

	c.JSON(200, gin.H{
		"message": "Stream started with id: " + id,
	})
}
