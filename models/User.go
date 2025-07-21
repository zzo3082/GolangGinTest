package model

// User model
// type User struct {
// 	ID       int    `json:"UserId" validate:"required"`
// 	Name     string `json:"UserName" validate:"required"`
// 	Password string `json:"UserPassword" validate:"required"`
// 	Email    string `json:"UserEmail" gorm:"Column:ZzoEmail"` // db內欄位名稱是 ZzoEmail 但這邊是 Email 不是對應的話可以加 gorm:"Column:ZzoEmail"
// }

// User model 改用 binding tag
type User struct {
	ID       int    `json:"UserId" binding:"required"`
	Name     string `json:"UserName" binding:"required,gt=5"`
	Password string `json:"UserPassword" binding:"required,gt=4,max=20,ZzoUserPasd"` // gt=4 是大於4個字元, ZzoUserPasd 是自定義驗證規則
	Email    string `json:"UserEmail" gorm:"Column:ZzoEmail"`                        // db內欄位名稱是 ZzoEmail 但這邊是 Email 不是對應的話可以加 gorm:"Column:ZzoEmail"
}

type Users struct {
	UserList []User `json:"Users" binding:"required,gt=0,lt=30"` // gt=0 是大於0個元素, lt=30 是小於30個元素
	Count    int    `json:"UsersCount"`
}
