package httpmiddleware

import (
	"time"

	_request "common/request"

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

		requestID, _ := _request.
			RequestIDFromContext(
				ctx.Request.Context(),
			)

		operationID, _ :=
			_request.
				OperationIDFromContext(
					ctx.Request.Context(),
				)

		route := normalizedRoute(
			ctx,
		)

		recorder.RecordHTTPRequest(
			ctx.Request.Context(),
			HTTPRequestResult{
				RequestID:        requestID,
				OperationID:      operationID,
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
