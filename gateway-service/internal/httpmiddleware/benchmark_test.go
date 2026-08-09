package httpmiddleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func BenchmarkGatewayRequestBoundary(
	benchmark *testing.B,
) {
	gin.SetMode(gin.TestMode)

	router := gin.New()

	router.Use(
		RequestID(nil),
		RequestLifecycle(nil),
	)

	router.GET(
		"/ok",
		func(ctx *gin.Context) {
			ctx.Status(http.StatusNoContent)
		},
	)

	benchmark.ReportAllocs()
	benchmark.ResetTimer()

	for index := 0; index < benchmark.N; index++ {
		request := httptest.NewRequest(
			http.MethodGet,
			"/ok",
			nil,
		)

		response := httptest.NewRecorder()

		router.ServeHTTP(
			response,
			request,
		)

		if response.Code != http.StatusNoContent {
			benchmark.Fatal(
				"unexpected status",
			)
		}
	}
}
