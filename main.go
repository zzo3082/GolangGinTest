package main

import "github.com/gin-gonic/gin"

func main() {
	// Create a new Gin router instance
	router := gin.Default()

	// Define a simple GET endpoint
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"MESSAGE": "pong",
		})
	})
	router.Run(":8080") // Start the server on port 8080
}
