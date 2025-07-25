package services

import (
	model "GolangAPI/models"
	repository "GolangAPI/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// 寫 User 邏輯操作的檔案, 驗證輸入 回傳

// Get User
func FindAllUsers(c *gin.Context) {
	// c.JSON(http.StatusOK, userList)
	// 從 db 撈 users 出來
	users := repository.FindAllUsers()
	if len(users) == 0 {
		c.JSON(http.StatusNotFound, "error : No users found.")
		return
	}
	c.JSON(http.StatusOK, users)
}

// Get User by ID
func FindByUserId(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Param("id"))
	user := repository.FindByUserId(userId)
	if user.ID == 0 {
		c.JSON(http.StatusNotFound, "error : User not found.")
		return
	}
	c.JSON(http.StatusOK, user)
}

// Post User
func PostUser(c *gin.Context) {
	user := model.User{}

	// 用 binging 屬性來驗證輸入, 在main註冊後, 可加入字定義的 middleware
	// c.BINDJSON 只有判斷屬性是否有對應 , 沒有判斷屬性 required
	// ex : UserId是int, 如果輸入string, 會加入err
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, "error : "+err.Error())
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "error : "+err.Error())
		return
	}
	user.Password = string(hashedPassword)
	insertedUser, err := repository.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "error : "+err.Error())
		return
	}
	c.JSON(http.StatusOK, insertedUser)
}

// Post Multiple Users
func PostUsers(c *gin.Context) {
	users := model.Users{}

	// 用 binging 屬性來驗證輸入, 在main註冊後, 可加入字定義的 middleware
	err := c.BindJSON(&users)
	if err != nil {
		c.JSON(http.StatusBadRequest, "error : "+err.Error())
		return
	}

	// 把userList的密碼加密
	for _, user := range users.UserList {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, "error : "+err.Error())
			return
		}
		user.Password = string(hashedPassword)
	}

	err = repository.CreateUsers(users.UserList) // 1. 直接DB.Create
	//err = repository.CreateUsersBatch(users.UserList) // 2. 使用batch分批 Create
	// err = repository.CreateUsersBulk(users.UserList) // 3. 使用 SQL 指令批量插入
	if err != nil {
		c.JSON(http.StatusInternalServerError, "error : "+err.Error())
		return
	}
	c.JSON(http.StatusOK, "message : PostUsers Successed.")
}

// Delete User
func DeleteUser(c *gin.Context) {
	userId, _ := strconv.Atoi(c.Param("id"))
	isDeleted := repository.DeleteUser(userId)
	if !isDeleted {
		c.JSON(http.StatusNotFound, "error : User not found.")
		return
	}
	c.JSON(http.StatusOK, "message : DeleteUser Successed.")
}

// Update User
func PutUser(c *gin.Context) {
	updateUser := model.User{}
	err := c.BindJSON(&updateUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, "error : "+err.Error())
		return
	}
	updateUser = repository.UpdateUser(updateUser.ID, updateUser)
	if updateUser.ID == 0 {
		c.JSON(http.StatusNotFound, "error : User not found.")
		return
	}
	c.JSON(http.StatusOK, "message : UpdateUser Successed. UserId: "+strconv.Itoa(updateUser.ID))
}
