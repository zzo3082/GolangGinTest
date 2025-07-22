package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	// "gopkg.in/mgo.v2" // 版本偏舊
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDb *mongo.Collection

func MongoDBConnect() {
	// 載入 .env 檔案
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 從環境變數中讀取連線字串
	mongoDbConnString := os.Getenv("mongoDbConnString")
	if mongoDbConnString == "" {
		log.Fatal("mongoDbConnString not set in .env file")
	}
	// 設置 MongoDB 客戶端選項
	clientOptions := options.Client().ApplyURI(mongoDbConnString)

	// 設置連線超時
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 連接到 MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}

	// 選擇資料庫和集合
	MongoDb = client.Database("GolangMongoDb").Collection("User")
}

// 套件版本偏舊 我 mongodb 裝新版的 連不到
// func MongoDBConnect() {
// 	// 載入 .env 檔案
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}

// 	// 從環境變數中讀取連線字串
// 	mongoDbConnString := os.Getenv("mongoDbConnString")
// 	if mongoDbConnString == "" {
// 		log.Fatal("mongoDbConnString not set in .env file")
// 	}

// 	session, err := mgo.Dial("mongodb://localhost:27017")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	MongoDb = session.DB("GolangApi").C("Test")
// }
