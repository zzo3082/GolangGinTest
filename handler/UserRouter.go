package src

import (
	"GolangAPI/services"

	"github.com/gin-gonic/gin"
)

// 管理 User 路由的檔案

func AddUserRouter(router *gin.RouterGroup) {
	user := router.Group("/user")

	user.GET("/", services.FindAllUsers)
	user.GET("/:id", services.FindByUserId)
	user.POST("/", services.PostUser)
	user.POST("/batch", services.PostUsers) // 批量新增用戶
	user.DELETE("/:id", services.DeleteUser)
	user.PUT("/:id", services.PutUser)
}
