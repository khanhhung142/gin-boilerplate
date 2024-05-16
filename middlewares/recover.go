package middlewares

import (
	"gin-boilerplate/consts"
	"gin-boilerplate/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RecoverMiddleware() gin.HandlerFunc {
	// recover from panic
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				errCast, ok := err.(error)
				if !ok {
					errCast = consts.CodeErrorUnknown
				}
				logger.Error(c, "panic", errCast)
				c.Header("Content-Type", "application/json")
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			}
		}()
		c.Next()
	}
}
