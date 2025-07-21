package model

// User model
type User struct {
	ID       int    `json:"UserId" validate:"required"`
	Name     string `json:"UserName" validate:"required"`
	Password string `json:"UserPassword" validate:"required"`
	Email    string `json:"UserEmail" gorm:"Column:ZzoEmail"` // db內欄位名稱是 ZzoEmail 但這邊是 Email 不是對應的話可以加 gorm:"Column:ZzoEmail"
}
