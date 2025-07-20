package main

import (
	. "GolangAPI/src"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create a new Gin router instance
	router := gin.Default()

	// 1. Simple Get/Post
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

	// 如果要給 Group path, 就是 http://localhost:8080/v1/Simple
	// 可以不給 Group path 這樣路由就變成 http://localhost:8080/Simple
	//v1 := router.Group("")
	// 2. User Router
	v1 := router.Group("/v1")
	AddUserRouter(v1)

	// 3. Simple Router
	AddSimpleRouter(v1)

	router.Run(":8080") // Start the server on port 8080
}
