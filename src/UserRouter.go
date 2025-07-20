package src

import (
	"GolangAPI/services"

	"github.com/gin-gonic/gin"
)

func AddUserRouter(router *gin.RouterGroup) {
	user := router.Group("/user")

	user.GET("/", services.FindAllUsers)
	user.GET("/:id", services.FindByUserId)
	user.POST("/", services.PostUser)
	user.DELETE("/:id", services.DeleteUser)
	user.PUT("/:id", services.PutUser)
}
