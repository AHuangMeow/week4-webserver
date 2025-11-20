package utils

import "github.com/gin-gonic/gin"

// Response General response structure
// @Description General API response structure
type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
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

func BadRequest(c *gin.Context, message string) {
	Error(c, 400, message)
}

func Unauthorized(c *gin.Context, message string) {
	Error(c, 401, message)
}

func NotFound(c *gin.Context, message string) {
	Error(c, 404, message)
}

func InternalError(c *gin.Context, message string) {
	Error(c, 500, message)
}
