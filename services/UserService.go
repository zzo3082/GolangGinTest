package services

import (
	model "GolangAPI/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var userList = []model.User{}

// Get User
func FindAllUsers(c *gin.Context) {
	c.JSON(http.StatusOK, userList)
}

// Post User
func PostUser(c *gin.Context) {
	user := model.User{}

	// c.BINDJSON 只有判斷屬性是否有對應 , 沒有判斷屬性 required
	// ex : UserId是int, 如果輸入string, 會加入err
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, "error : "+err.Error())
		return
	}

	// 用套件 github.com/go-playground/validator/v10
	// 來驗證User的validate:"required"屬性
	validate := validator.New()
	err = validate.Struct(user)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			c.JSON(http.StatusBadRequest, "error : "+err.Namespace()+" "+err.Tag()+" "+err.Param())
			return
		}
	}

	userList = append(userList, user)
	c.JSON(http.StatusOK, "message : PostUser Successed.")

}

// Delete User
func DeleteUser(c *gin.Context) {
	userIdStr := c.Param("id") // 這個撈出來是 string
	// 用 strconv.Atoi 轉成 int, 失敗會有error 所以前面要多一個 err 去接收, 不要的話可以用 _ 代替
	userId, _ := strconv.Atoi(userIdStr)
	for _, user := range userList {
		log.Println(user)
		if user.ID == userId {
			userList = append(userList[:userId-1], userList[userId:]...)
			c.JSON(http.StatusOK, "message : DeleteUser Successed.")
			return
		}
	}

	c.JSON(http.StatusNotFound, "error : User not found.")
}

// Put User
func PutUser(c *gin.Context) {
	updateUser := model.User{}
	err := c.BindJSON(&updateUser)
	if err != nil {
		c.JSON(http.StatusBadRequest, "error : "+err.Error())
		return
	}
	userId, _ := strconv.Atoi(c.Param("id"))
	for index, user := range userList {
		if user.ID == userId {
			// 更新用戶資料
			userList[index] = updateUser
			log.Println("Updated User:", userList[index])
			c.JSON(http.StatusOK, "message : PutUser Successed.")
			return
		}
	}
	c.JSON(http.StatusNotFound, "error : User not found.")
}
