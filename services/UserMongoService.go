package services

import (
	model "GolangAPI/models"
	repository "GolangAPI/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func MongoCresteUser(c *gin.Context) {
	user := model.User{}
	err := c.Bind(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Error : "+err.Error())
		return
	}

	repository.MongoCreater(user)
	c.JSON(http.StatusOK, gin.H{
		"message": "Create Mongo User Done.",
		"User":    user,
	})
}

// find all users
func MongoFindUsers(c *gin.Context) {
	users := repository.MongoFindAllUsers()
	c.JSON(http.StatusOK, gin.H{
		"message": "Find Mongo User Done.",
		"User":    users,
	})
}

// find user by id
func MongoFindUserById(c *gin.Context) {
	id := c.Param("id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Error : id can not convert to int.")
		return
	}

	user := repository.MongoFindUserById(userId)
	c.JSON(http.StatusOK, gin.H{
		"message": "Find Mongo User Done.",
		"User":    user,
	})
}

// update user
func MongoUpdateUser(c *gin.Context) {
	user := model.User{}
	err := c.Bind(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Error : "+err.Error())
		return
	}

	user = repository.MongoUpdateUser(user.ID, user)
	if user.ID == 0 {
		c.JSON(http.StatusBadRequest, "Error : Can't find User with Id")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Update Mongo User Done.",
		"User":    user,
	})
}

// delete user
func MongoDeleteUser(c *gin.Context) {
	id := c.Param("id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Error : id can not convert to int.")
		return
	}

	repository.MongoDeleteUser(userId)
	c.JSON(http.StatusOK, gin.H{
		"message": "Delete Mongo User Done.",
	})
}
