package _errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AppError định nghĩa lỗi chung của ứng dụng
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Error triển khai interface `error`
func (e *AppError) Error() string {
	return e.Message
}

// NewAppError tạo lỗi mới
func NewAppError(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Middleware xử lý lỗi và trả về JSON chuẩn
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // Xử lý request

		// Kiểm tra nếu có lỗi xảy ra
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			if appErr, ok := err.(*AppError); ok {
				// Nếu là AppError, trả về JSON theo format chuẩn
				c.JSON(appErr.Code, gin.H{"code": appErr.Code, "message": appErr.Message})
				return
			}
			// Lỗi không xác định
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "Internal Server Error"})
		}
	}
}
