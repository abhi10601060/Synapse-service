package main

import (
	"fmt"
	"log"
	"synapse/stream/auth"
	"synapse/stream/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Welcome to Synapse Stream")

	r := gin.Default()

	r.GET("/start", authorize, handler.StartStream)

	r.Run(":8010")
}

func authorize(c *gin.Context) {
	tokenStr := c.Request.Header.Get("Authentication-Token")

	if tokenStr == "" {
		c.JSON(405, gin.H{
			"message": "Header token missing",
		})
		c.Abort()
		return
	}

	isTokenValid, err := auth.IsAuthorizedToken(tokenStr)
	if err != nil {
		c.JSON(401, gin.H{
			"message": "In correct Header Token",
		})
		c.Abort()
		return
	}

	if !isTokenValid {
		c.JSON(401, gin.H{
			"message": "Unauthorized Header Token",
		})
		c.Abort()
		return
	}
	log.Println("Authorization Passed wit valid token...")
}
