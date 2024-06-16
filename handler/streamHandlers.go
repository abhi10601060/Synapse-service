package handler

import "github.com/gin-gonic/gin"

func StartStream(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Stream started",
	})
}
