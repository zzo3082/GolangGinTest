package Repository

import (
	"GolangAPI/database"
	. "GolangAPI/models"
	model "GolangAPI/models"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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
func CreateUser(user User) (User, error) {
	result := database.DBConn.Create(&user)
	if result.Error != nil {
		return User{}, result.Error // 如果有錯誤，回傳空的 User 結構
	}
	return user, nil
}

// Insert Multiple Users
func CreateUsers(users []User) error {
	result := database.DBConn.Create(&users)
	if result.Error != nil {
		return result.Error // 如果有錯誤，回傳錯誤
	}
	return nil
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

// 批量加入大量資料 2種方法
// 1. 分批加入
func CreateUsersBatch(users []User) error {
	// 每次處理 1000 筆資料
	batchSize := 1
	return database.DBConn.CreateInBatches(users, batchSize).Error
}

// 2. 使用 SQL 指令
func CreateUsersBulk(users []User) error {
	// 設定每批插入的筆數
	batchSize := 1
	valueStrings := make([]string, 0, batchSize)
	valueArgs := make([]interface{}, 0, batchSize*4) // 每個 User 有 4 個欄位

	for i, user := range users {
		// 每個 User 的值格式為 "(?, ?, ?, ?)"
		valueStrings = append(valueStrings, "(?, ?, ?, ?)")
		valueArgs = append(valueArgs, user.ID, user.Name, user.Password, user.Email)

		// 每 batchSize 筆或最後一批執行插入
		if len(valueStrings) == batchSize || i == len(users)-1 {
			stmt := `INSERT INTO users (ID, Name, Password, ZzoEmail) VALUES ` + strings.Join(valueStrings, ",")
			sqlDb, _ := database.DBConn.DB()
			_, err := sqlDb.Exec(stmt, valueArgs...)
			if err != nil {
				return err
			}
			// 清空批次資料
			valueStrings = valueStrings[:0]
			valueArgs = valueArgs[:0]
		}
	}
	return nil
}

// Transaction
func CreateUsersTransaction(users []User) error {
	tx := database.DBConn.Begin()
	if err := tx.Error; err != nil {
		return err
	}

	for _, user := range users {
		if err := tx.Create(&user).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// CheckUserPassword
func CheckUserPassword(username, password string) (User, error) {
	var user User
	err := database.DBConn.Where("name = ? AND password = ?", username, password).First(&user).Error
	if err != nil {
		return User{}, err // 如果找不到使用者，回傳錯誤
	}
	return user, nil // 找到使用者，回傳使用者資料
}

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
