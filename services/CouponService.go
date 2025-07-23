package services

import (
	"GolangAPI/middlewares"
	model "GolangAPI/models"
	apiModel "GolangAPI/models/ApiModels"
	repository "GolangAPI/repository"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Create Coupon
func CreateCoupon(c *gin.Context) {
	coupon := model.Coupon{}
	err := c.Bind(&coupon)
	if err != nil {
		c.JSON(http.StatusBadRequest, "傳入值無法轉成coupon, 請確認欄位.")
		return
	}

	coupon.CreatedAt = time.Now()

	coupon, err = repository.CreateCoupon(coupon)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Error : "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Create Coupon Done!",
		"Coupon":  coupon,
	})

}

func ClaimCoupon(c *gin.Context) {
	claimCouponReq := apiModel.ClaimCouponRequestDto{}
	err := c.Bind(&claimCouponReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, "錯誤的輸入.")
		return
	}

	// 找看看coupon
	coupon, err := repository.GetCoupon(claimCouponReq.CouponCode)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
	}

	// 確認日期
	now := time.Now()
	if now.Before(coupon.StartDate) || now.After(coupon.EndDate) {
		c.JSON(http.StatusBadRequest, "error : 現在不是優惠券的使用期間.")
	}

	// 確認 currentUses
	if coupon.CurrentUses >= coupon.MaxUses {
		c.JSON(http.StatusBadRequest, "error : 優惠券發放數量已達上限.")
		return
	}

	// 更新 Coupon current_uses, updateAt
	err = repository.UpdateCouponAfterClaimed(coupon)
	if err != nil {
		c.JSON(http.StatusBadRequest, "error : UpdateCouponAfterClaimed 失敗.")
		return
	}

	userId := middlewares.GetSessionUserId(c)

	// 新增 user coupon 紀錄
	err = repository.CreateUserCoupon(userId, int(coupon.ID))
	if err != nil {
		c.JSON(http.StatusBadRequest, "error : CreateUserCoupon 失敗.")
		return
	}

	c.JSON(http.StatusOK, "message : ClaimCoupon 成功.")

}
