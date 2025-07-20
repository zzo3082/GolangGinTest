package src

import (
	"GolangAPI/services"

	"github.com/gin-gonic/gin"
)

func AddUserRouter(router *gin.RouterGroup) {
	user := router.Group("/user")

	user.GET("/", services.FindAllUsers)
	user.POST("/", services.PostUser)
}
