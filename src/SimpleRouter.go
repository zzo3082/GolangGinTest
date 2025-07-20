package src

import (
	"GolangAPI/services"

	"github.com/gin-gonic/gin"
)

func AddSimpleRouter(router *gin.RouterGroup) {
	// 注意這邊大小寫有差, 路由打成小寫, http://localhost:8080/v1/simple, 會404
	simple := router.Group("/Simple")

	simple.GET("/", services.Get)

	// 這邊的 :id 是路由參數
	// 跟 Post 方法內的 c.Param("id") 對應
	// 例如 http://localhost:8080/v1/Simple/123
	simple.POST("/:id", services.Post)
}
