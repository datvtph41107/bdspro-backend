package _routes

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ResponseDTO là cấu trúc JSON chuẩn
type ResponseDTO struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"` // Data có thể là bất kỳ kiểu dữ liệu nào
}

// ResponseDTO là cấu trúc JSON chuẩn
type Except struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *Except) Error() string {
	return fmt.Sprintf("%d: %s", e.Code, e.Message)
}

func RouteResult(c *gin.Context, result interface{}, err error) {
	if err != nil {
		var specErr *Except
		if errors.As(err, &specErr) { // Kiểm tra lỗi có phải AuthError không
			c.JSON(http.StatusOK, ResponseDTO{
				Code:    specErr.Code,
				Message: specErr.Message,
			})
		} else {
			fmt.Println("Lỗi khác:", result)
			c.JSON(http.StatusOK, ResponseDTO{
				Code:    500,
				Message: err.Error(),
			})
		}
		return
	}

	c.JSON(http.StatusOK, ResponseDTO{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    result,
	})
}
