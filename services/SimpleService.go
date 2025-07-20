package services

import (
	"github.com/gin-gonic/gin"
)

// Simple Get
func Get(c *gin.Context) {
	c.JSON(200, gin.H{"MESSAGE": "pong"})
}

// Simple Post
func Post(c *gin.Context) {
	id := c.Param("id")
	c.JSON(200, gin.H{"PostId": id})
}
