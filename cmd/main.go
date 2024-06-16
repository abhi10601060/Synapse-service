package main

import (
	"fmt"
	"synapse/stream/handler"
	"synapse/stream/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("Welcome to Synapse Stream")

	r := gin.Default()

	stream := r.Group("/stream")
	{
		stream.GET("/start", middleware.Authorize, handler.StartStream)
	}

	r.Run(":8010")
}
