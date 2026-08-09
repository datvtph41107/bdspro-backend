package httpmiddleware

import (
	_requestcontext "common/requestcontext"

	"github.com/gin-gonic/gin"
)

// RequestID canonicalize request ID tại public boundary
func RequestID(
	recorder RequestIDDecisionRecorder,
) gin.HandlerFunc {
	if recorder == nil {
		// Do nothing
		recorder = noopRequestIDDecisionRecorder{}
	}

	return func(ctx *gin.Context) {
		decision := _requestcontext.
			PrepareHTTPRequestDecision(
				ctx.Writer,
				ctx.Request,
			)

		recorder.RecordRequestIDDecision(
			ctx.Request.Context(),
			decision,
		)

		ctx.Next()
	}
}
