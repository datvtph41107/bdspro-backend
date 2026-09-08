package httpmiddleware

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"

	_request "common/request"

	"github.com/gin-gonic/gin"
)

type testRecorder struct {
	decisions []_request.RequestIDDecision
	results   []HTTPRequestResult
}

func (recorder *testRecorder) RecordRequestIDDecision(
	_ context.Context,
	decision _request.RequestIDDecision,
) {
	recorder.decisions = append(recorder.decisions, decision)
}

func (recorder *testRecorder) RecordHTTPRequest(
	_ context.Context,
	result HTTPRequestResult,
) {
	recorder.results = append(recorder.results, result)
}

func newTestRouter(recorder *testRecorder) *gin.Engine {
	gin.SetMode(gin.TestMode)

	logger := slog.New(
		slog.NewJSONHandler(
			&bytes.Buffer{},
			nil,
		),
	)

	router := gin.New()
	router.HandleMethodNotAllowed = true

	router.Use(
		RequestID(recorder),
		RequestLifecycle(recorder),
		Recovery(logger),
	)

	router.GET("/ok", func(ctx *gin.Context) {
		requestID, ok := _request.RequestIDFromContext(
			ctx.Request.Context(),
		)

		if !ok {
			ctx.AbortWithStatus(
				http.StatusInternalServerError,
			)
			return
		}

		ctx.String(
			http.StatusOK,
			requestID,
		)
	})

	router.GET("/unauthorized", func(ctx *gin.Context) {
		ctx.AbortWithStatus(
			http.StatusUnauthorized,
		)
	})

	router.GET("/forbidden", func(ctx *gin.Context) {
		ctx.AbortWithStatus(
			http.StatusForbidden,
		)
	})

	router.GET("/unavailable", func(ctx *gin.Context) {
		ctx.AbortWithStatus(
			http.StatusServiceUnavailable,
		)
	})

	router.GET("/panic", func(*gin.Context) {
		panic("simulated panic")
	})

	return router
}
