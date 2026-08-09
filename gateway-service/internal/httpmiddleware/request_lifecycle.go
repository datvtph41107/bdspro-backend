package httpmiddleware

import (
	"time"

	_requestcontext "common/requestcontext"

	"github.com/gin-gonic/gin"
)

// RequestLifecycle ghi nhận kết quả cuối của HTTP attempt
func RequestLifecycle(
	recorder HTTPRequestRecorder,
) gin.HandlerFunc {
	if recorder == nil {
		recorder = noopHTTPRequestRecorder{}
	}

	return func(ctx *gin.Context) {
		startedAt := time.Now()

		ctx.Next()

		transportOutcome,
			transportReason :=
			classifyHTTPTransport(
				ctx,
			)

		requestID, _ := _requestcontext.
			FromContext(
				ctx.Request.Context(),
			)

		route := normalizedRoute(
			ctx,
		)

		recorder.RecordHTTPRequest(
			ctx.Request.Context(),
			HTTPRequestResult{
				RequestID:        requestID,
				Method:           ctx.Request.Method,
				Route:            route,
				StatusCode:       ctx.Writer.Status(),
				Duration:         time.Since(startedAt),
				Aborted:          ctx.IsAborted(),
				ErrorCount:       len(ctx.Errors),
				TransportOutcome: transportOutcome,
				TransportReason:  transportReason,
			},
		)
	}
}
