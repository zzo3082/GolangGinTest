package middlewares

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// 指定輸出格式的 Logger 中間件
func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(params gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] %s %s %d \n",
			params.ClientIP,   // 客戶端IP
			params.TimeStamp,  // 日誌時間戳
			params.Method,     // HTTP方法
			params.Path,       // 請求路由
			params.StatusCode) // HTTP狀態碼
	})

}
