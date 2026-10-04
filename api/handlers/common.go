package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIError represents the standard error response body
type APIError struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

// ErrorResponse wraps the APIError
type ErrorResponse struct {
	Error APIError `json:"error"`
}

// RespondError sends a standardized JSON error response
func RespondError(c *gin.Context, status int, code, message string, details ...string) {
	c.AbortWithStatusJSON(status, ErrorResponse{
		Error: APIError{
			Code:    code,
			Message: message,
			Details: details,
		},
	})
}

// RespondSuccess sends a standardized JSON response
func RespondSuccess(c *gin.Context, status int, data interface{}) {
	if data == nil {
		c.Status(status)
		return
	}
	c.JSON(status, data)
}

// CheckContentType ensures the request has application/json Content-Type if body is present
func ValidateJSONContentType(c *gin.Context) bool {
	if c.Request.ContentLength > 0 {
		ct := c.ContentType()
		if ct != "application/json" && ct != "" {
			RespondError(c, http.StatusUnsupportedMediaType, "UNSUPPORTED_MEDIA_TYPE",
				"Content-Type must be application/json")
			return false
		}
	}
	return true
}
