package apimodels

import . "GolangAPI/models"

type UserResponse struct {
	ID    int    `json:"UserId"`
	Name  string `json:"UserName"`
	Email string `json:"UserEmail"`
}

func NewUserResponse(u User) UserResponse {
	return UserResponse{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}
