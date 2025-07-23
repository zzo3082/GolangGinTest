package Repository

import (
	"GolangAPI/database"
	. "GolangAPI/models"
	"fmt"
)

// Create Coupon
func CreateCoupon(coupon Coupon) (Coupon, error) {
	result := database.DBConn.Create(&coupon)
	if result.Error != nil {
		return Coupon{}, result.Error // 如果有錯誤，回傳空的 Coupon 結構
	}
	return coupon, nil
}

// 查 Coupon
func GetCoupon(code string) (Coupon, error) {
	var coupon Coupon
	result := database.DBConn.Where("code = ?", code).First(&coupon)
	if result.Error != nil {
		return Coupon{}, fmt.Errorf("coupon not found with code %s", code)
	}
	return coupon, nil
}
