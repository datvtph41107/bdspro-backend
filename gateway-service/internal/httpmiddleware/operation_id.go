package httpmiddleware

import (
	"errors"
	"net/http"

	_requestcontext "common/requestcontext"

	"github.com/gin-gonic/gin"
)

type operationIDPreparer func(
	http.ResponseWriter,
	*http.Request,
) (
	_requestcontext.OperationIDDecision,
	error,
)

/**
 * OperationID establishes one logical-operation identity
 * after the request has crossed the CORS boundary.
 */
func OperationID() gin.HandlerFunc {
	return operationID(
		_requestcontext.
			PrepareHTTPOperationID,
	)
}

func operationID(
	prepare operationIDPreparer,
) gin.HandlerFunc {
	return func(
		ctx *gin.Context,
	) {
		_, err :=
			prepare(
				ctx.Writer,
				ctx.Request,
			)

		if err == nil {
			ctx.Next()

			return
		}

		abortOperationIDFailure(
			ctx,
			err,
		)
	}
}

func abortOperationIDFailure(
	ctx *gin.Context,
	err error,
) {
	switch {
	case errors.Is(
		err,
		_requestcontext.
			ErrInvalidOperationID,
	),
		errors.Is(
			err,
			_requestcontext.
				ErrMultipleOperationIDs,
		):

		ctx.AbortWithStatusJSON(
			http.StatusBadRequest,
			gin.H{
				"error": "invalid X-Operation-ID header",
			},
		)

	default:
		ctx.AbortWithStatusJSON(
			http.StatusInternalServerError,
			gin.H{
				"error": "internal server error",
			},
		)
	}
}
