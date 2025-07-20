package database

import (
	"log"
	"os"

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
		log.Fatal("Error loading .env file")
	}

	// 從環境變數中讀取連線字串
	dbConnString := os.Getenv("dbConnString")
	if dbConnString == "" {
		log.Fatal("dbConnString not set in .env file")
	}

	DBConn, err = gorm.Open(mysql.Open(dbConnString), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}
}
