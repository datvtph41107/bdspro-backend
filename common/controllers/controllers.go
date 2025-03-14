package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response là cấu trúc JSON chuẩn
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"` // Data có thể là bất kỳ kiểu dữ liệu nào
}

// BaseController chứa các phương thức dùng chung
type BaseController struct{}

// Success trả về response thành công
func (bc *BaseController) Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    data,
	})
}

// Error trả về response lỗi
func (bc *BaseController) Error(c *gin.Context, code int, message string) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}
