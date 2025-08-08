// Redis.go - 負責初始化與管理 Redis 連線池，提供全域 RedisDefaultPool 供全專案使用
package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/joho/godotenv"
)

var RedisDefaultPool *redis.Pool

// 容器化用這個
func newPool(host, port, password string) *redis.Pool {
	addr := host + ":" + port
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			opts := []redis.DialOption{}
			if password != "" {
				opts = append(opts, redis.DialPassword(password))
			}
			return redis.Dial("tcp", addr, opts...)
		},
	}
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file. If docker is fine.")
	}
	// 從環境變數讀取連線資訊
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD") // 可選：如果 Redis 無密碼，設為空

	// 驗證環境變數
	if redisHost == "" || redisPort == "" {
		log.Fatal("REDIS_HOST or REDIS_PORT not set")
	}

	// 初始化 Redis 連線池
	RedisDefaultPool = newPool(redisHost, redisPort, redisPassword)

	// 測試連線
	conn := RedisDefaultPool.Get()
	defer conn.Close()
	tmp, err := conn.Do("PING")
	fmt.Printf("Redis PING response: %v\n", tmp)
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	log.Println("Connected to Redis successfully")
}

// func newPool(addr string) *redis.Pool {
// 	return &redis.Pool{
// 		MaxIdle:     3,                 // Maximum number of idle connections
// 		IdleTimeout: 240 * time.Second, // Wait for connection to be available
// 		Dial: func() (redis.Conn, error) {
// 			return redis.Dial("tcp", addr)
// 		}, // Close connections after use
// 	}
// }

// func init() {
// 	// 載入 .env 檔案
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}

// 	// 從環境變數中讀取連線字串
// 	redisConnString := os.Getenv("redisConnString")
// 	if redisConnString == "" {
// 		log.Fatal("redisConnString not set in .env file")
// 	}
// 	RedisDefaultPool = newPool(redisConnString) // Default Redis address
// }
