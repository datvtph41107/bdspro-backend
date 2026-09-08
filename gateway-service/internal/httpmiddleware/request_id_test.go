package httpmiddleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	_request "common/request"
)

func TestGatewayResponsesRetainRequestID(
	t *testing.T,
) {
	testCases := []struct {
		name       string
		method     string
		path       string
		requestID  string
		wantStatus int
	}{
		{
			name:       "success",
			method:     http.MethodGet,
			path:       "/ok",
			requestID:  "request-success-123",
			wantStatus: http.StatusOK,
		},
		{
			name:       "authentication failure",
			method:     http.MethodGet,
			path:       "/unauthorized",
			requestID:  "request-auth-123",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "authorization failure",
			method:     http.MethodGet,
			path:       "/forbidden",
			requestID:  "request-forbidden-123",
			wantStatus: http.StatusForbidden,
		},
		{
			name:       "route not found",
			method:     http.MethodGet,
			path:       "/missing",
			requestID:  "request-not-found-123",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "method not allowed",
			method:     http.MethodPost,
			path:       "/ok",
			requestID:  "request-method-123",
			wantStatus: http.StatusMethodNotAllowed,
		},
		{
			name:       "downstream unavailable",
			method:     http.MethodGet,
			path:       "/unavailable",
			requestID:  "request-unavailable-123",
			wantStatus: http.StatusServiceUnavailable,
		},
		{
			name:       "panic",
			method:     http.MethodGet,
			path:       "/panic",
			requestID:  "request-panic-123",
			wantStatus: http.StatusInternalServerError,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			recorder := &testRecorder{}
			router := newTestRouter(recorder)

			request := httptest.NewRequest(
				testCase.method,
				testCase.path,
				nil,
			)

			request.Header.Set(
				RequestIDHeader,
				testCase.requestID,
			)

			response := httptest.NewRecorder()

			router.ServeHTTP(
				response,
				request,
			)

			if response.Code != testCase.wantStatus {
				t.Fatalf(
					"status = %d, want %d",
					response.Code,
					testCase.wantStatus,
				)
			}

			responseRequestID := response.Header().Get(
				RequestIDHeader,
			)

			if responseRequestID != testCase.requestID {
				t.Fatalf(
					"response request ID = %q, want %q",
					responseRequestID,
					testCase.requestID,
				)
			}

			if len(recorder.results) != 1 {
				t.Fatalf(
					"recorded results = %d, want 1",
					len(recorder.results),
				)
			}

			result := recorder.results[0]

			if result.RequestID != testCase.requestID {
				t.Fatalf(
					"recorded request ID = %q",
					result.RequestID,
				)
			}

			if result.StatusCode != testCase.wantStatus {
				t.Fatalf(
					"recorded status = %d, want %d",
					result.StatusCode,
					testCase.wantStatus,
				)
			}
		})
	}
}

func TestGatewayRecordsRequestIDDecision(
	t *testing.T,
) {
	testCases := []struct {
		name       string
		values     []string
		wantSource _request.RequestIDSource
		wantReason _request.RequestIDReason
	}{
		{
			name:       "missing",
			wantSource: _request.RequestIDSourceGenerated,
			wantReason: _request.RequestIDReasonMissing,
		},
		{
			name:       "valid",
			values:     []string{"client-request-123"},
			wantSource: _request.RequestIDSourceClient,
			wantReason: _request.RequestIDReasonValid,
		},
		{
			name:       "invalid",
			values:     []string{"unsafe request id"},
			wantSource: _request.RequestIDSourceGenerated,
			wantReason: _request.RequestIDReasonInvalid,
		},
		{
			name: "multiple",
			values: []string{
				"request-a",
				"request-b",
			},
			wantSource: _request.RequestIDSourceGenerated,
			wantReason: _request.RequestIDReasonMultiple,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			recorder := &testRecorder{}
			router := newTestRouter(recorder)

			request := httptest.NewRequest(
				http.MethodGet,
				"/ok",
				nil,
			)

			for _, value := range testCase.values {
				request.Header.Add(
					RequestIDHeader,
					value,
				)
			}

			response := httptest.NewRecorder()

			router.ServeHTTP(
				response,
				request,
			)

			if len(recorder.decisions) != 1 {
				t.Fatalf(
					"decisions = %d, want 1",
					len(recorder.decisions),
				)
			}

			decision := recorder.decisions[0]

			if decision.Source != testCase.wantSource {
				t.Fatalf(
					"source = %q, want %q",
					decision.Source,
					testCase.wantSource,
				)
			}

			if decision.Reason != testCase.wantReason {
				t.Fatalf(
					"reason = %q, want %q",
					decision.Reason,
					testCase.wantReason,
				)
			}
		})
	}
}
