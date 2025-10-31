package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/v3/jwk"
	"github.com/rahulSailesh-shah/converSense/pkg/auth"
)


func AuthMiddleware(authKeys jwk.Set) gin.HandlerFunc {
	return func(c *gin.Context){
		userId, err := auth.UserFromToken(c.Request, authKeys)
		if err != nil {
			c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return
		}
		c.Set("userId", userId)
		c.Next()
	}
}
