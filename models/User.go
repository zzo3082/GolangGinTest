package model

// User model
type User struct {
	ID       int    `json:"UserId" validate:"required"`
	Name     string `json:"UserName" validate:"required"`
	Password string `json:"UserPassword" validate:"required"`
	Email    string `json:"UserEmail"`
}
