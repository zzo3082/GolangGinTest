package middlewares

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

const userKey = "session_id"

// use cookie to store session id
func SetSession() gin.HandlerFunc {
	store := cookie.NewStore([]byte("userKey"))
	return sessions.Sessions("zzoSession", store)
}

// Authenticate user session
func AuthSession() gin.HandlerFunc {
	return func(c *gin.Context) {
		sessoin := sessions.Default(c)
		sessionId := sessoin.Get(userKey)
		if sessionId == nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "這邊要登入才能進"})
			return
		}
		c.Next()
	}
}

// Save Session
func SaveSession(c *gin.Context, userId int) {
	sessoin := sessions.Default(c)
	sessoin.Set(userKey, userId)
	if err := sessoin.Save(); err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": "無法儲存 session"})
		return
	}
}

// Clear Session
func ClearSession(c *gin.Context) {
	sessoin := sessions.Default(c)
	sessoin.Clear()
	if err := sessoin.Save(); err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": "無法清除 session"})
		return
	}
}

// Get Session User ID
func GetSessionUserId(c *gin.Context) int {
	sessoin := sessions.Default(c)
	sessionId := sessoin.Get(userKey)
	if sessionId == nil {
		return -1
	}
	userId, ok := sessionId.(int)
	if !ok {
		return -1
	}
	return userId
}
