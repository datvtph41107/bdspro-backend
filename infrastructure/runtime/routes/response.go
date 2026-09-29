package _routes

import (
	"log/slog"
	"net/http"
	"reflect"

	_errors "common/errors"
	"common/logging"

	"github.com/gin-gonic/gin"
)

// ResponseDTO là cấu trúc JSON chuẩn
type ResponseDTO struct {
	Code          int         `json:"code"`
	Message       string      `json:"message"`
	Data          interface{} `json:"data,omitempty"`
	TotalElements *int64      `json:"totalElements,omitempty"`
	Errors        interface{} `json:"errors,omitempty"`
}

func RouteResult(c *gin.Context, result interface{}, err error) {
	if err != nil {
		if application, ok := _errors.As(err); ok {
			code := application.Spec().LegacyCode()
			if occurrenceCode, exists := application.LegacyCode(); exists {
				code = occurrenceCode
			}
			if code == 0 {
				code = int32(application.Code())
			}

			var fieldErrors map[string]string
			if violations := application.Violations(); len(violations) > 0 {
				fieldErrors = make(map[string]string, len(violations))
				for _, violation := range violations {
					if violation.Field != "" {
						fieldErrors[violation.Field] = violation.Description
					}
				}
			}

			c.JSON(http.StatusOK, ResponseDTO{
				Code:    int(code),
				Message: application.PublicMessage(),
				Errors:  fieldErrors,
			})
			return
		}

		logging.WithComponent(c.Request.Context(), "http.route_result").Error(
			"direct HTTP request failed",
			slog.Any("error", err),
		)
		c.JSON(http.StatusOK, ResponseDTO{
			Code:    500,
			Message: "internal server error",
		})
		return
	}

	if reflect.TypeOf(result) == reflect.TypeOf(ResponseDTO{}) {
		c.JSON(http.StatusOK, result)
		return
	}

	c.JSON(http.StatusOK, ResponseDTO{
		Code:    0,
		Message: "Success",
		Data:    result,
	})
}
