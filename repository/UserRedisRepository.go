package Repository

import (
	"GolangAPI/database"
	model "GolangAPI/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 從db撈出一個User, 並將結果寫入 gin.Context
// 要給 CacheOneUserDecorator 方法用的
// 第一次因為 Redis 沒有資料, 會呼叫這個方法
func RedisOneUser(c *gin.Context) {
	id := c.Param("id")
	if id == "0" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	user := model.User{}
	database.DBConn.Where("id = ?", id).First(&user)
	c.Set("dbResult", user)
}
