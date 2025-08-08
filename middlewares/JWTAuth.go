package middlewares

import (
	"GolangAPI/database"
	models "GolangAPI/models"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(c *gin.Context) {
	// Get Cookie
	tokenStr, err := c.Cookie("Authorization")
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// 驗證 JWT
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		return []byte(os.Getenv("JWTSecret")), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		// log.Default(err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// 上面 jwt.WithValidMethods 有驗證 exp 了
		// if float64(time.Now().Unix()) > claims["exp"].(float64) {
		// 	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
		// 	return
		// }

		// 查db user
		var user models.User
		database.DBConn.Where("id = ?", int(claims["user_id"].(float64))).First(&user)
		if user.ID == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		c.Set("user", user)
		// 在 Router "user.GET("/validate", middlewares.RequireAuth, services.ValidateJWT)"
		// 使用 c.Next() 下一步是到 services.ValidateJWT 回傳 api 結果
		c.Next()
	} else {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}
}
