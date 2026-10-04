package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// JSONRecovery catches panics and formats them into a clean JSON 500 error response
func JSONRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[PANIC RECOVERED] %v\n%s", r, debug.Stack())
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error": gin.H{
						"code":    "INTERNAL_SERVER_ERROR",
						"message": "An unexpected server error occurred",
					},
				})
			}
		}()
		c.Next()
	}
}
