package httpauth

import (
	"net/http"

	_httpresponse "gateway/internal/httpresponse"

	"github.com/gin-gonic/gin"
)

func abortProblem(ctx *gin.Context, status int, code, detail string) {
	_httpresponse.WriteProblem(
		ctx.Request.Context(),
		ctx.Writer,
		_httpresponse.NewProblem(status, code, detail),
	)
	ctx.Abort()
}

func abortInternalProblem(ctx *gin.Context, code string) {
	abortProblem(ctx, http.StatusInternalServerError, code, "internal server error")
}
