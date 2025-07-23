package services

import (
	"GolangAPI/middlewares"
	apimodel "GolangAPI/models/ApiModels"
	repository "GolangAPI/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Login
func Login(c *gin.Context) {
	// 用 form 表單傳入資訊
	// name := c.PostForm("name")
	// password := c.PostForm("password")
	var loginInfo apimodel.LoginInfoDto
	err := c.BindJSON(&loginInfo)
	if err != nil {
		c.JSON(http.StatusBadRequest, "error : "+err.Error())
		return
	}

	user, err := repository.CheckUserPassword(loginInfo.UserName, loginInfo.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "error : Invalid username or password.")
		return
	}
	// 登入成功後儲存 session
	middlewares.SaveSession(c, user.ID)

	c.JSON(http.StatusOK, gin.H{
		"Message": "Login successful",
		"User":    user,
		"Session": middlewares.GetSessionUserId(c),
	})
}

// LogOut
func LogOut(c *gin.Context) {
	middlewares.ClearSession(c)
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
}

// check if user is logged in
func CheckUserSession(c *gin.Context) {
	sessionId := middlewares.GetSessionUserId(c)
	if sessionId == -1 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not logged in"})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"Message": "Login successful",
		"Session": sessionId,
	})
}
