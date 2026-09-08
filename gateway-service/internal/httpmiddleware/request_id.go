package httpmiddleware

import (
	_request "common/request"

	"github.com/gin-gonic/gin"
)

// RequestID canonicalizes one request ID at the public HTTP boundary.
func RequestID(recorder RequestIDDecisionRecorder) gin.HandlerFunc {
	if recorder == nil {
		recorder = noopRequestIDDecisionRecorder{}
	}
	return func(ctx *gin.Context) {
		decision := _request.ResolveRequestID(ctx.Request.Header.Values(RequestIDHeader))
		requestCtx := _request.WithRequestID(ctx.Request.Context(), decision.ID)
		ctx.Request = ctx.Request.WithContext(requestCtx)
		ctx.Request.Header.Set(RequestIDHeader, decision.ID)
		ctx.Writer.Header().Set(RequestIDHeader, decision.ID)
		recorder.RecordRequestIDDecision(requestCtx, decision)
		ctx.Next()
	}
}
