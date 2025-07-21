package middlewares

import (
	model "GolangAPI/models"
	"regexp"
	"unicode"

	"github.com/go-playground/validator/v10"
)

// 用 regex 驗證密碼格式
func UserPasd(field validator.FieldLevel) bool {
	password := field.Field().String()

	// Step 1: 檢查長度與只包含英數
	match, _ := regexp.MatchString(`^[A-Za-z0-9]{6,}$`, password)
	if !match {
		return false
	}

	// Step 2: 確認同時包含「至少一個英文字母」與「至少一個數字」
	hasLetter := false
	hasDigit := false

	for _, c := range password {
		if unicode.IsLetter(c) {
			hasLetter = true
		} else if unicode.IsDigit(c) {
			hasDigit = true
		}
	}

	return hasLetter && hasDigit
}

func CkeckUserList(field validator.StructLevel) {
	users := field.Current().Interface().(model.Users)
	if users.Count != len(users.UserList) {
		field.ReportError(users.Count,
			"Count of user list",
			"Count",
			"UsersCount",
			"CkeckUserListCountFailed")
	}
}
