package middleware

import (
	"log"
	"synapse/stream/auth"
	"github.com/gin-gonic/gin"
)

func Authorize(c *gin.Context) {
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
