package services

import (
	model "GolangAPI/models"
	repository "GolangAPI/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

	// model改用binding tag, 這邊就不用再用 validator 驗證了
	// 用套件 github.com/go-playground/validator/v10
	// 來驗證User的validate:"required"屬性
	// validate := validator.New()
	// err = validate.Struct(user)
	// if err != nil {
	// 	for _, err := range err.(validator.ValidationErrors) {
	// 		c.JSON(http.StatusBadRequest, "error : "+err.Namespace()+" "+err.Tag()+" "+err.Param())
	// 		return
	// 	}
	// }

	insertedUser, err := repository.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "error : "+err.Error())
		return
	}
	c.JSON(http.StatusOK, insertedUser)
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
