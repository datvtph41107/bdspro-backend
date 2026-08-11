package httpmiddleware

import (
	_requestcontext "common/requestcontext"

	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	invalidIdempotencyKeyMessage     = "invalid Idempotency-Key header"
	idempotencyKeyUnavailableMessage = "idempotency identity unavailable"
)

// IdempotencyKey accepts an optional caller-owned command identity.
// It validates and binds a supplied key, but never originates one.
func IdempotencyKey() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, _, err := _requestcontext.PrepareHTTPIdempotencyKey(ctx.Request)
		if err == nil {
			ctx.Next()
			return
		}

		statusCode := http.StatusInternalServerError
		message := idempotencyKeyUnavailableMessage
		if errors.Is(err, _requestcontext.ErrInvalidIdempotencyKey) ||
			errors.Is(err, _requestcontext.ErrMultipleIdempotencyKeys) {
			statusCode = http.StatusBadRequest
			message = invalidIdempotencyKeyMessage
		}

		ctx.AbortWithStatusJSON(statusCode, gin.H{"error": message})
	}
}
