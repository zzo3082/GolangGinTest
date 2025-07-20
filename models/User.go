package model

import "GolangAPI/database"

type User struct {
	ID       int    `json:"UserId" validate:"required"`
	Name     string `json:"UserName" validate:"required"`
	Password string `json:"UserPassword" validate:"required"`
	Email    string `json:"UserEmail"`
}

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
