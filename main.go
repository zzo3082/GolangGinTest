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

	// Define a POST endpoint with a path parameter id
	router.POST("/ping/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.JSON(200, gin.H{
			"PostId": id,
		})
	})

	router.Run(":8080") // Start the server on port 8080
}
