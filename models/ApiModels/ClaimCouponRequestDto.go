package apimodels

type ClaimCouponRequestDto struct {
	CouponCode string `json:"code" binding:"required"`
}
