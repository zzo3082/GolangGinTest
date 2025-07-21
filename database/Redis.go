// Redis.go - 負責初始化與管理 Redis 連線池，提供全域 RedisDefaultPool 供全專案使用
package database

import (
	"log"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/joho/godotenv"
)

var RedisDefaultPool *redis.Pool

func newPool(addr string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,                 // Maximum number of idle connections
		IdleTimeout: 240 * time.Second, // Wait for connection to be available
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", addr)
		}, // Close connections after use
	}
}

func init() {
	// 載入 .env 檔案
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 從環境變數中讀取連線字串
	redisConnString := os.Getenv("redisConnString")
	if redisConnString == "" {
		log.Fatal("redisConnString not set in .env file")
	}
	RedisDefaultPool = newPool(redisConnString) // Default Redis address
}
