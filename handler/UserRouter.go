package src

import (
	sessions "GolangAPI/middlewares"
	"GolangAPI/services"

	"github.com/gin-gonic/gin"
)

// 管理 User 路由的檔案

func AddUserRouter(router *gin.RouterGroup) {
	// user := router.Group("/user")
	user := router.Group("/user", sessions.SetSession())

	user.GET("/", services.FindAllUsers)
	user.GET("/:id", services.FindByUserId)
	user.POST("/", services.PostUser)
	user.POST("/batch", services.PostUsers) // 批量新增用戶
	user.PUT("/:id", services.PutUser)

	user.POST("/login", services.Login) // 登入用戶

	user.Use(sessions.AuthSession()) // 使用者登入後的驗證中間件
	{
		user.DELETE("/:id", services.DeleteUser)
		user.GET("/logout", services.LogOut)          // 登出用戶
		user.GET("/check", services.CheckUserSession) // 檢查使用者登入狀態
	}
}
