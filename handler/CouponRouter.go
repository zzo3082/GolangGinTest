package src

import (
	sessions "GolangAPI/middlewares"
	"GolangAPI/services"

	"github.com/gin-gonic/gin"
)

func AddCouponRouter(router *gin.RouterGroup) {
	coupon := router.Group("/coupon", sessions.SetSession())
	coupon.POST("/", services.CreateCoupon)

	coupon.Use(sessions.AuthSession()) // 使用者登入後的驗證中間件
	{
		coupon.POST("/cliam", services.ClaimCoupon)
	}
}
