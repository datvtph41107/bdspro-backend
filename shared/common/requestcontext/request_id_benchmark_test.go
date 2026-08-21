package requestcontext

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkPrepareHTTPRequest(benchmark *testing.B) {
	benchmark.ReportAllocs()
	benchmark.ResetTimer()

	for index := 0; index < benchmark.N; index++ {
		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()
		PrepareHTTPRequest(response, request)
	}
}
