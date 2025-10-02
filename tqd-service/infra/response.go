package infra

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response represents a standard API response
type Response struct {
	Code         int         `json:"code"`
	Message      string      `json:"message"`
	Data         interface{} `json:"data,omitempty"`
	Total        *int64      `json:"total,omitempty"`
	TotalElement *int64      `json:"totalElement,omitempty"`
}

// SuccessResponse creates a success response
func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    http.StatusOK,
		Message: "Success",
		Data:    data,
	})
}

// SuccessResponseWithPagination creates a success response with pagination
func SuccessResponseWithPagination(c *gin.Context, data interface{}, total int64) {
	c.JSON(http.StatusOK, Response{
		Code:         0,
		Message:      "Success",
		Data:         data,
		Total:        &total,
		TotalElement: &total,
	})
}

// ErrorResponse creates an error response
func ErrorResponse(c *gin.Context, code int, message string) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
	})
}

// BadRequestResponse creates a bad request response
func BadRequestResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusBadRequest, message)
}

// NotFoundResponse creates a not found response
func NotFoundResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusNotFound, message)
}

// InternalServerErrorResponse creates an internal server error response
func InternalServerErrorResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusInternalServerError, message)
}

// UnauthorizedResponse creates an unauthorized response
func UnauthorizedResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusUnauthorized, message)
}

// ForbiddenResponse creates a forbidden response
func ForbiddenResponse(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusForbidden, message)
}
