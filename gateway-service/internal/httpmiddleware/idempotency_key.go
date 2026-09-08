package httpmiddleware

import (
	_request "common/request"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	invalidIdempotencyKeyMessage     = "invalid Idempotency-Key header"
	idempotencyKeyUnavailableMessage = "idempotency identity unavailable"
)

func IdempotencyKey() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		key, ok, err := _request.ParseIdempotencyKey(ctx.Request.Header.Values(IdempotencyKeyHeader))
		if err == nil && !ok {
			ctx.Next()
			return
		}
		if err == nil {
			bound, bindErr := _request.BindIdempotencyKey(ctx.Request.Context(), key)
			if bindErr == nil {
				ctx.Request.Header.Set(IdempotencyKeyHeader, key)
				ctx.Request = ctx.Request.WithContext(bound)
				ctx.Next()
				return
			}
			err = bindErr
		}
		statusCode := http.StatusInternalServerError
		message := idempotencyKeyUnavailableMessage
		if errors.Is(err, _request.ErrInvalidIdempotencyKey) || errors.Is(err, _request.ErrMultipleIdempotencyKeys) {
			statusCode = http.StatusBadRequest
			message = invalidIdempotencyKeyMessage
		}
		ctx.AbortWithStatusJSON(statusCode, gin.H{"error": message})
	}
}
