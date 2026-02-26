package util

import (
	"github.com/gin-gonic/gin"
)

// ErrorResponse represents a standard error structure
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

// JSONError sends a standard JSON error response
func JSONError(c *gin.Context, status int, message string, err error) {
	resp := ErrorResponse{
		Code:    status,
		Message: message,
	}
	if err != nil {
		resp.Error = err.Error()
	}
	c.AbortWithStatusJSON(status, resp)
}

// BadRequest sends a 400 Bad Request response
func BadRequest(c *gin.Context, message string, err error) {
	JSONError(c, 400, message, err)
}

// InternalError sends a 500 Internal Server Error response
func InternalError(c *gin.Context, message string, err error) {
	JSONError(c, 500, message, err)
}

// Unauthorized sends a 401 Unauthorized response
func Unauthorized(c *gin.Context, message string) {
	JSONError(c, 401, message, nil)
}

// Forbidden sends a 403 Forbidden response
func Forbidden(c *gin.Context, message string) {
	JSONError(c, 403, message, nil)
}

// SuccessResponse sends a standard JSON success response (optional wrapper)
func SuccessResponse(c *gin.Context, status int, data interface{}) {
	c.JSON(status, data)
}
