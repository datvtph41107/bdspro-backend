package httpmiddleware

import (
	"errors"
	"net/http"
	"syscall"

	"github.com/gin-gonic/gin"
)

type HTTPTransportOutcome string

const (
	HTTPTransportOutcomeCompleted   HTTPTransportOutcome = "completed"
	HTTPTransportOutcomeInterrupted HTTPTransportOutcome = "interrupted"
)

type HTTPTransportReason string

const (
	HTTPTransportReasonNone            HTTPTransportReason = "none"
	HTTPTransportReasonBrokenPipe      HTTPTransportReason = "broken_pipe"
	HTTPTransportReasonConnectionReset HTTPTransportReason = "connection_reset"
	HTTPTransportReasonAbortHandler    HTTPTransportReason = "abort_handler"
)

func classifyHTTPTransport(
	ctx *gin.Context,
) (
	HTTPTransportOutcome,
	HTTPTransportReason,
) {
	for _, contextError := range ctx.Errors {
		if contextError == nil ||
			contextError.Err == nil {
			continue
		}

		switch {
		case errors.Is(
			contextError.Err,
			syscall.EPIPE,
		):
			return HTTPTransportOutcomeInterrupted,
				HTTPTransportReasonBrokenPipe

		case errors.Is(
			contextError.Err,
			syscall.ECONNRESET,
		):
			return HTTPTransportOutcomeInterrupted,
				HTTPTransportReasonConnectionReset

		case errors.Is(
			contextError.Err,
			http.ErrAbortHandler,
		):
			return HTTPTransportOutcomeInterrupted,
				HTTPTransportReasonAbortHandler
		}
	}

	return HTTPTransportOutcomeCompleted,
		HTTPTransportReasonNone
}
