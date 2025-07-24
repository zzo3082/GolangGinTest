package Repository

import (
	"GolangAPI/database"
	. "GolangAPI/models"
	"GolangAPI/models/enums"
	"fmt"
	"time"

	"gorm.io/gorm"
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
func UpdateCouponAfterClaimed(tx *gorm.DB, coupon Coupon) error {
	coupon.UpdatedAt = time.Now()
	coupon.CurrentUses++

	result := tx.Model(&coupon).Where("id = ?", coupon.ID).Updates(coupon)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func CreateUserCoupon(tx *gorm.DB, userid int, couponid int) error {
	userCoupons := UserCoupon{
		UserID:    int64(userid),
		CouponID:  int64(couponid),
		Status:    enums.String(enums.UNUSED),
		ClaimedAt: time.Now(),
	}
	result := tx.Create(userCoupons)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func ClaimCouponTransaction(userid int, coupon Coupon) error {
	// 把 UpdateCouponAfterClaimed 跟 CreateUserCoupon 包在 Transaction 內
	// 兩個都成功可以更新到 db
	return database.DBConn.Transaction(func(tx *gorm.DB) error {
		// 這邊開始 Transaction
		if err := UpdateCouponAfterClaimed(tx, coupon); err != nil {
			return err
		}

		if err := CreateUserCoupon(tx, userid, int(coupon.ID)); err != nil {
			return err
		}
		return nil
		// Transaction 結束, 沒有 err 回傳 nil
	})
}
