package requestcontext

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func BenchmarkIsValidRequestID(benchmark *testing.B) {
	requestID := "req_0123456789abcdef0123456789abcdef"

	benchmark.ReportAllocs()
	benchmark.ResetTimer()

	for index := 0; index < benchmark.N; index++ {
		if !IsValidRequestID(requestID) {
			benchmark.Fatal("valid ID rejected")
		}
	}
}

func BenchmarkNewRequestID(benchmark *testing.B) {
	benchmark.ReportAllocs()
	benchmark.ResetTimer()

	for index := 0; index < benchmark.N; index++ {
		_ = NewRequestID()
	}
}

func BenchmarkPrepareHTTPRequest(benchmark *testing.B) {
	benchmark.ReportAllocs()
	benchmark.ResetTimer()

	for index := 0; index < benchmark.N; index++ {
		request := httptest.NewRequest(
			http.MethodGet,
			"/",
			nil,
		)

		response := httptest.NewRecorder()

		PrepareHTTPRequest(
			response,
			request,
		)
	}
}
