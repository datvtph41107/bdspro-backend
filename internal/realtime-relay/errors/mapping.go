package errors

import "net/http"

var CodeToHTTPStatus = map[int32]int{
	4001: http.StatusNotFound,
	4002: http.StatusForbidden,
	4003: http.StatusForbidden,
	4004: http.StatusRequestEntityTooLarge,
	4005: http.StatusForbidden,
	4006: http.StatusBadRequest,
	4007: http.StatusForbidden,
	4008: http.StatusConflict,
	4009: http.StatusConflict,
	4010: http.StatusUnprocessableEntity,
}