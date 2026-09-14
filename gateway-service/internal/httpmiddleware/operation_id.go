package httpmiddleware

import (
	_request "common/request"
	"errors"
	"net/http"

	_httpresponse "gateway/internal/httpresponse"

	"github.com/gin-gonic/gin"
)

type operationIDPreparer func(http.ResponseWriter, *http.Request) (_request.OperationIDDecision, error)

func OperationID() gin.HandlerFunc { return operationID(prepareHTTPOperationID) }

func prepareHTTPOperationID(w http.ResponseWriter, r *http.Request) (_request.OperationIDDecision, error) {
	decision, err := _request.ResolveOperationID(r.Header.Values(OperationIDHeader))
	if err != nil {
		return _request.OperationIDDecision{}, err
	}
	bound, err := _request.BindOperationID(r.Context(), decision.ID)
	if err != nil {
		return _request.OperationIDDecision{}, err
	}
	r.Header.Set(OperationIDHeader, decision.ID)
	w.Header().Set(OperationIDHeader, decision.ID)
	*r = *r.WithContext(bound)
	return decision, nil
}

func operationID(prepare operationIDPreparer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_, err := prepare(ctx.Writer, ctx.Request)
		if err == nil {
			ctx.Next()
			return
		}
		abortOperationIDFailure(ctx, err)
	}
}

func abortOperationIDFailure(ctx *gin.Context, err error) {
	problem := _httpresponse.NewProblem(
		http.StatusInternalServerError,
		"request.operation_id.unavailable",
		"internal server error",
	)
	if errors.Is(err, _request.ErrInvalidOperationID) || errors.Is(err, _request.ErrMultipleOperationIDs) {
		problem = _httpresponse.NewProblem(
			http.StatusBadRequest,
			"request.operation_id.invalid",
			"invalid X-Operation-ID header",
		)
	}
	_httpresponse.WriteProblem(ctx.Request.Context(), ctx.Writer, problem)
	ctx.Abort()
}
