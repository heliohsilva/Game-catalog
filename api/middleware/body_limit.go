package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// MaxBodySize restricts the maximum bytes allowed in a request body
func MaxBodySize(limitBytes int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Body != nil {
			c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, limitBytes)
		}
		c.Next()
	}
}
