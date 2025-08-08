package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DBConn *gorm.DB

var err error

func DBConnect() {
	// 載入 .env 檔案
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file. If docker is fine.!!!!!!!!!!!!!!!!!!!!!!!!")
	}
	// 讀取環境變數
	mysqlHost := os.Getenv("MYSQL_HOST")
	mysqlPort := os.Getenv("MYSQL_PORT")
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	mysqlDatabase := os.Getenv("MYSQL_DATABASE")
	var mysqlDSN string

	// 驗證環境變數
	if mysqlHost == "" || mysqlPort == "" || mysqlUser == "" || mysqlDatabase == "" {
		panic("Missing required MySQL environment variables")
	}
	if mysqlPassword == "" {
		mysqlDSN = fmt.Sprintf("root:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			mysqlUser, mysqlHost, mysqlPort, mysqlDatabase)
	} else {
		// 連接到 MySQL
		mysqlDSN = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			mysqlUser, mysqlPassword, mysqlHost, mysqlPort, mysqlDatabase)
	}

	// **新增重試邏輯**
	maxRetries := 10
	retryInterval := 5 * time.Second

	for i := 0; i < maxRetries; i++ {
		log.Printf("Connecting to MySQL atempt %d/%d...", i+1, maxRetries)
		DBConn, err = gorm.Open(mysql.Open(mysqlDSN), &gorm.Config{})
		if err == nil {
			log.Println("Connecting to MySQL successfully")
			return // 連線成功，結束函式
		}
		log.Printf("Failed to connect: %v", err)
		time.Sleep(retryInterval) // 等待一段時間後再重試
	}

	// 超過最大重試次數後仍失敗
	log.Fatal("Failed to connect to MySQL after multiple retries")
}
