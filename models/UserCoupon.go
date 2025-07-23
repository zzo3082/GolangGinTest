package model

import "time"

type UserCoupon struct {
	UserID    int64     `json:"user_id"`
	CouponID  int64     `json:"coupon_id"`
	Status    string    `json:"status"`
	ClaimedAt time.Time `json:"claimed_at"`
}
