package main

import (
	"GolangAPI/database"
	. "GolangAPI/handler"
	"GolangAPI/middlewares"
	model "GolangAPI/models"
	"io"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// log設置
func setupLogger() {
	// os.Create 回回傳一個檔案
	f, _ := os.Create("gin.log")
	// gin.DefaultWriter 是 gin 的日誌輸出
	gin.DefaultWriter = io.MultiWriter(f, os.Stdout) // 將日誌輸出到文件和控制台
}

func main() {
	// 設置日誌輸出
	setupLogger()

	// Create a new Gin router instance
	router := gin.Default()

	// 註冊自定義驗證規則
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("ZzoUserPasd", middlewares.UserPasd) // 註冊自定義驗證規則
		// StructLevel 的驗證, 給方法要是validator.StructLevel, 後面給 model
		v.RegisterStructValidation(middlewares.CkeckUserList, model.Users{}) // 註冊結構體驗證
	}

	// 這邊去讀 middleware 的 log 格式, 也可以加入簡單的auth驗證
	//router.Use(gin.BasicAuth(gin.Accounts{"Tom": "123456"}), middlewares.Logger())
	router.Use(gin.Recovery(), middlewares.Logger())

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

	// 4. 連資料庫
	go func() {
		database.DBConnect()
	}()

	router.Run(":8080") // Start the server on port 8080
}
