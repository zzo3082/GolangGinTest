package Repository

import (
	"GolangAPI/database"
	. "GolangAPI/models"
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

func MongoCreater(user User) User {
	ctx := context.Background()
	_, err := database.MongoDb.InsertOne(ctx, user)
	if err != nil {
		log.Printf("Failed to insert user: %v", err)
		return User{}
	}
	return user
}

func MongoFindAllUsers() []User {
	var users []User
	ctx := context.Background()
	cursor, err := database.MongoDb.Find(ctx, bson.M{}) // 這邊的 cusor 是查詢結果的游標(cusor)
	if err != nil {
		log.Printf("Failed to find users: %v", err)
		return nil
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &users); err != nil {
		log.Printf("Failed to decode users: %v", err)
		return nil
	}
	return users
}

func MongoFindUserById(id int) User {
	var user User
	ctx := context.Background()
	err := database.MongoDb.FindOne(ctx, bson.M{"id": id}).Decode(&user) // Decode 把結果 反序列化成 struct
	if err != nil {
		log.Printf("Failed to find user: %v", err)
		return User{}
	}
	return user
}

func MongoUpdateUser(id int, user User) User {
	ctx := context.Background()
	filter := bson.M{"id": id}     // 是 mapping 用的 類似 json { id = id}
	update := bson.M{"$set": user} // 更新物件的語法
	_, err := database.MongoDb.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Failed to update user: %v", err)
		return User{}
	}
	return user
}

func MongoDeleteUser(id int) bool {
	ctx := context.Background()
	filter := bson.M{"id": id}
	_, err := database.MongoDb.DeleteOne(ctx, filter)
	if err != nil {
		log.Printf("Failed to delete user: %v", err)
		return false
	}
	return true
}
