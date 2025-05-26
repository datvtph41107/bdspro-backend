package _routes

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

// ResponseDTO là cấu trúc JSON chuẩn
type ResponseDTO struct {
	Code          int         `json:"code"`
	Message       string      `json:"message"`
	Data          interface{} `json:"data,omitempty"` // Data có thể là bất kỳ kiểu dữ liệu nào
	TotalElements *int64      `json:"totalElements,omitempty"`
	Errors        interface{} `json:"errors,omitempty"`
}

// ResponseDTO là cấu trúc JSON chuẩn
type Except struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors,omitempty"`
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
				Errors:  specErr.Errors,
			})
		} else {
			fmt.Println("Lỗi khác:", err)
			c.JSON(http.StatusOK, ResponseDTO{
				Code:    500,
				Message: err.Error(),
			})
		}
		return
	}

	if reflect.TypeOf(result) == reflect.TypeOf(ResponseDTO{}) {
		c.JSON(http.StatusOK, result)
		return
	}

	// json := jsoniter.Con
	// rs, _ := json.MarshalIndent(result, "", "  ")
	c.JSON(http.StatusOK, ResponseDTO{
		Code:    0,
		Message: "Success",
		Data:    result,
	})
}

// func RouteResult(c *gin.Context, result ResponseDTO) {
// 	c.JSON(http.StatusOK, result)
// }
