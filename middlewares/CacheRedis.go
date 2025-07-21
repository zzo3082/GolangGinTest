package middlewares

import (
	redisDB "GolangAPI/database"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/pquerna/ffjson/ffjson"
)

// 這個裝飾器會先嘗試從 Redis 獲取資料，如果沒有找到，則調用原始處理函數(RedisOneUser)從 MySQL 獲取資料，並將結果存入 Redis
func CacheOneUserDecorator(h gin.HandlerFunc, param string, readKeyPattern string, empty interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		keyId := c.Param(param)                        // 取得路由參數 這邊範例是 id
		redisKey := fmt.Sprintf(readKeyPattern, keyId) // 生成 Redis key, 這邊範例給的是("user:%s" ,id) > user:1
		conn := redisDB.RedisDefaultPool.Get()         // 取得 Redis 連線
		defer conn.Close()                             // 確保連線在使用後被關閉

		data, err := redis.Bytes(conn.Do("GET", redisKey)) // 嘗試用剛剛組成的 redisKey 從 Redis 獲取資料
		if err != nil {                                    // 如果 Redis 中沒有資料, err 會有值
			h(c)                                  // 呼叫原始的處理函數, 這裡會去 MySQL 獲取資料, 這邊使用的是 UserRepository.RedisOneUser
			dbResult, exists := c.Get("dbResult") // Repository.RedisOneUser 裡面有把 db 撈出來的結果寫進 c.Set("dbResult", user)
			if !exists {
				dbResult = empty
			}
			redisData, _ := ffjson.Marshal(dbResult)  // db有撈到值的話 將資料轉換為 []byte
			conn.Do("SETEX", redisKey, 30, redisData) // 把 []byte 資料存入 Redis, 並設定 30 秒的過期時間跟剛剛組成的 redisKey
			c.JSON(http.StatusOK, gin.H{
				"message": "Get From MySql",
				"data":    dbResult,
			}) // 回傳從 MySQL 獲取的資料
			return
		}
		ffjson.Unmarshal(data, &empty) // 如果 Redis 中有資料, 將資料反序列化到 empty 變數中
		c.JSON(http.StatusOK, gin.H{
			"message": "Get From Cache",
			"data":    empty,
		}) // 回傳從 Redis 獲取的資料
	}
}
