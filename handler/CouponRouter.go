package src

import (
	"GolangAPI/services"

	"github.com/gin-gonic/gin"
)

func AddCouponRouter(router *gin.RouterGroup) {
	coupon := router.Group("/coupon")
	coupon.POST("/", services.CreateCoupon)
}
