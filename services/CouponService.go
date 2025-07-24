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
		c.JSON(http.StatusBadRequest, gin.H{"error": "傳入值無法轉成coupon, 請確認欄位."})
		return
	}

	coupon.CreatedAt = time.Now()

	coupon, err = repository.CreateCoupon(coupon)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}

	// 其他service檔案的方法, 可以直接使用
	err = CreateRedisCoupon(coupon)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "錯誤的輸入."})
		return
	}
	// 1. 去 redis 判斷這個 coupon 的 current_uses/max_uses 是否達到上限了
	err = CheckAddCache(claimCouponReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		// 這邊不能rollback
		return
	}
	// 2. 去db找看看coupon (可能可以把coupon都改到redis)
	coupon, err := repository.GetCoupon(claimCouponReq.CouponCode)
	if err != nil {
		c.JSON(http.StatusNotFound, err.Error())
		RollbackCouponUsesCache(claimCouponReq.CouponCode)
		return
	}

	// 確認日期
	now := time.Now()
	if now.Before(coupon.StartDate) || now.After(coupon.EndDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "現在不是優惠券的使用期間."})
		RollbackCouponUsesCache(claimCouponReq.CouponCode)
		return
	}

	// 確認 currentUses
	if coupon.CurrentUses >= coupon.MaxUses {
		c.JSON(http.StatusBadRequest, gin.H{"error": "優惠券發放數量已達上限."})
		RollbackCouponUsesCache(claimCouponReq.CouponCode)
		return
	}

	// update coupon.current_uses 跟 insert user_coupon 紀錄
	// 用 transaction 綁在一起, 一個失敗 交易就失敗
	userId := middlewares.GetSessionUserId(c)
	err = repository.ClaimCouponTransaction(userId, coupon)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":        "ClaimCouponTransaction 失敗",
			"errorMessage": err.Error(),
		})
		RollbackCouponUsesCache(claimCouponReq.CouponCode)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ClaimCoupon 成功."})

}
