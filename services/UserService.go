package services

import (
	model "GolangAPI/models"
	"net/http"

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
