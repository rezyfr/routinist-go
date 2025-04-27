package middleware

import "github.com/gin-gonic/gin"

func ContentTypeApplicationJson() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Accept", "application/json")
		c.Next()
	}
}
