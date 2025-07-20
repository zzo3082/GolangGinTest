package Repository

import (
	"GolangAPI/database"
	. "GolangAPI/models"
)

// 這個檔案處理 User 資料庫CRUD操作

func FindAllUsers() []User {
	var users []User
	database.DBConn.Find(&users)
	return users
}

func FindByUserId(id int) User {
	var user User
	database.DBConn.Where("id = ?", id).First(&user)
	return user
}

// Post User
func CreateUser(user User) User {
	database.DBConn.Create(&user)
	return user
}

// Delete User
func DeleteUser(id int) bool {
	var user User
	result := database.DBConn.Where("id = ?", id).Delete(&user).RowsAffected
	return result == 1
}

// Update User
func UpdateUser(id int, user User) User {
	database.DBConn.Model(&user).Where("id = ?", id).Updates(user)
	return user
}
