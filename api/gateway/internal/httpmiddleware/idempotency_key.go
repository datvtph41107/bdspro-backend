package httpmiddleware

import (
	_request "common/request"
	"errors"
	"net/http"

	_httpresponse "gateway/internal/httpresponse"

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

		problem := _httpresponse.NewProblem(
			http.StatusInternalServerError,
			"request.idempotency_key.unavailable",
			idempotencyKeyUnavailableMessage,
		)
		if errors.Is(err, _request.ErrInvalidIdempotencyKey) || errors.Is(err, _request.ErrMultipleIdempotencyKeys) {
			problem = _httpresponse.NewProblem(
				http.StatusBadRequest,
				"request.idempotency_key.invalid",
				invalidIdempotencyKeyMessage,
			)
		}
		_httpresponse.WriteProblem(ctx.Request.Context(), ctx.Writer, problem)
		ctx.Abort()
	}
}
