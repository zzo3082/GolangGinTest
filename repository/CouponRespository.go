package Repository

import (
	"GolangAPI/database"
	. "GolangAPI/models"
	"fmt"
	"time"
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
	result := database.DBConn.Where("code = ? and isdeleted = ?", code, 0).First(&coupon)
	if result.Error != nil {
		return Coupon{}, fmt.Errorf("coupon not found with code %s", code)
	}
	return coupon, nil
}

// 拿取優惠券後更新 coupon
func UpdateCouponAfterClaimed(coupon Coupon) error {
	coupon.UpdatedAt = time.Now()
	coupon.CurrentUses++

	result := database.DBConn.Model(&coupon).Where("id = ?", coupon.ID).Updates(coupon)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func CreateUserCoupon(userid int, couponid int) error {
	userCoupons := UserCoupon{
		UserID:    int64(userid),
		CouponID:  int64(couponid),
		Status:    "UNUSED",
		ClaimedAt: time.Now(),
	}
	result := database.DBConn.Create(userCoupons)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
