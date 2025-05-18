package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"routinist/internal/auth"
	"strings"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Missing Authorization header"})
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid Authorization header format"})
			return
		}

		claims, err := auth.ParseJWT(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid or expired token"})
			return
		}

		// Attach user ID or email to context
		c.Set("user_id", claims.ID)
		c.Set("email", claims.Email)
		c.Set("expired_at", claims.ExpiresAt.Time)

		c.Next()
	}
}
