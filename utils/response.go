package utils

import "github.com/gin-gonic/gin"

// Response General response structure
// @Description General API response structure
type Response struct {
	Code    int    `json:"code" example:"200"`
	Message string `json:"message" example:"success"`
	Data    any    `json:"data"`
}

// Success Success response
func Success(c *gin.Context, data any) {
	c.JSON(200, Response{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

// Error Error response
func Error(c *gin.Context, code int, message string) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// BadRequest returns 400 error
func BadRequest(c *gin.Context, message string) {
	Error(c, 400, message)
}

// Unauthorized returns 401 error
func Unauthorized(c *gin.Context, message string) {
	Error(c, 401, message)
}

// NotFound returns 404 error
func NotFound(c *gin.Context, message string) {
	Error(c, 404, message)
}

// InternalError returns 500 error
func InternalError(c *gin.Context, message string) {
	Error(c, 500, message)
}

