package httpmiddleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"

	_requestcontext "common/requestcontext"

	"github.com/gin-gonic/gin"
)

// Recovery panic -> http 500
func Recovery(
	logger *slog.Logger,
) gin.HandlerFunc {
	if logger == nil {
		logger = slog.Default()
	}

	return gin.CustomRecovery(
		func(
			ctx *gin.Context,
			recovered any,
		) {
			requestID, _ := _requestcontext.
				FromContext(
					ctx.Request.Context(),
				)

			logger.ErrorContext(
				ctx.Request.Context(),
				"gateway request panic",
				slog.String(
					"event_name",
					"http.request.panic",
				),
				slog.String(
					"request_id",
					requestID,
				),
				slog.String(
					"http_method",
					ctx.Request.Method,
				),
				slog.String(
					"http_route",
					normalizedRoute(ctx),
				),
				slog.Any(

					"panic",
					string(debug.Stack()),
				),
			)

			// Header request ID đã được set
			// trc khi chain chạy
			if ctx.Writer.Written() {
				ctx.Abort()
				return
			}

			ctx.AbortWithStatus(
				http.StatusInternalServerError,
			)
		},
	)
}
