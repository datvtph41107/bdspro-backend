package httpmiddleware

import (
	_request "common/request"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

func Recovery(
	logger *slog.Logger,
) gin.HandlerFunc {
	if logger == nil {
		logger = slog.Default()
	}

	return gin.CustomRecoveryWithWriter(
		nil,
		func(
			ctx *gin.Context,
			recovered any,
		) {
			requestID, _ :=
				_request.RequestIDFromContext(
					ctx.Request.Context(),
				)

			recoveredErr :=
				asRecoveredError(
					recovered,
				)

			_ = ctx.Error(
				recoveredErr,
			)

			attrs := []slog.Attr{
				slog.String(
					"event_name",
					"http.request.panic",
				),
				slog.String(
					"request_id",
					requestID,
				),
			}

			if operationID, ok :=
				_request.OperationIDFromContext(
					ctx.Request.Context(),
				); ok {

				attrs = append(
					attrs,
					slog.String(
						"operation_id",
						operationID,
					),
				)
			}

			attrs = append(
				attrs,
				slog.String(
					"http_method",
					ctx.Request.Method,
				),
				slog.String(
					"http_route",
					normalizedRoute(ctx),
				),
				slog.Bool(
					"response_committed",
					ctx.Writer.Written(),
				),
				slog.String(
					"panic_type",
					fmt.Sprintf(
						"%T",
						recovered,
					),
				),
				slog.String(
					"stack_trace",
					string(
						debug.Stack(),
					),
				),
			)

			logger.LogAttrs(
				ctx.Request.Context(),
				slog.LevelError,
				"gateway request panic",
				attrs...,
			)

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

func asRecoveredError(
	recovered any,
) error {
	if err, ok :=
		recovered.(error); ok {

		return err
	}

	return fmt.Errorf(
		"%v",
		recovered,
	)
}
