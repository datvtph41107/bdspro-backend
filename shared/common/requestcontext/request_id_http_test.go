package requestcontext

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newRequestIDTestServer(
	t *testing.T,
) *httptest.Server {
	t.Helper()

	handler := http.HandlerFunc(
		func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			descisionID :=
				PrepareHTTPRequest(
					w,
					r,
				)

			contextReqID, ok :=
				FromContext(
					r.Context(),
				)

			if !ok {
				http.Error(
					w,
					"request ID missing from context",
					http.StatusInternalServerError,
				)
				return
			}

			if contextReqID != descisionID {
				http.Error(
					w,
					"request ID context drift",
					http.StatusInternalServerError,
				)
				return
			}

			if r.URL.Path == "/fail" {
				http.Error(
					w,
					"simulated failure",
					http.StatusServiceUnavailable,
				)
				return
			}

			_, _ = fmt.Fprint(
				w,
				descisionID,
			)
		},
	)

	return httptest.NewServer(handler)
}

func TestRequestIDHTTPRoundTrip(t *testing.T) {
	server := newRequestIDTestServer(t)
	defer server.Close()

	testCases := []struct {
		name          string
		path          string
		values        []string
		wantStatus    int
		wantExactID   string
		wantGenerated bool
	}{
		{
			name:          "missing generates",
			path:          "/echo",
			wantStatus:    http.StatusOK,
			wantGenerated: true,
		},
		{
			name:        "valid preserves",
			path:        "/echo",
			values:      []string{"client-request-123"},
			wantStatus:  http.StatusOK,
			wantExactID: "client-request-123",
		},
		{
			name:          "invalid replaces",
			path:          "/echo",
			values:        []string{"unsafe request id"},
			wantStatus:    http.StatusOK,
			wantGenerated: true,
		},
		{
			name: "multiple replaces",
			path: "/echo",
			values: []string{
				"request-a",
				"request-b",
			},
			wantStatus:    http.StatusOK,
			wantGenerated: true,
		},
		{
			name:        "failure still echoes",
			path:        "/fail",
			values:      []string{"request-failure-123"},
			wantStatus:  http.StatusServiceUnavailable,
			wantExactID: "request-failure-123",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			request, err := http.NewRequest(
				http.MethodGet,
				server.URL+testCase.path,
				nil,
			)
			if err != nil {
				t.Fatalf("NewRequest() error = %v", err)
			}

			for _, value := range testCase.values {
				request.Header.Add(
					RequestIDHeader,
					value,
				)
			}

			response, err := server.Client().Do(request)
			if err != nil {
				t.Fatalf("client.Do() error = %v", err)
			}
			defer response.Body.Close()

			body, err := io.ReadAll(response.Body)
			if err != nil {
				t.Fatalf("ReadAll() error = %v", err)
			}

			if response.StatusCode != testCase.wantStatus {
				t.Fatalf(
					"status = %d, want %d",
					response.StatusCode,
					testCase.wantStatus,
				)
			}

			responseRequestID := response.Header.Get(RequestIDHeader)

			if !IsValidRequestID(responseRequestID) {
				t.Fatalf(
					"response request ID invalid: %q",
					responseRequestID,
				)
			}

			if testCase.wantExactID != "" &&
				responseRequestID != testCase.wantExactID {
				t.Fatalf(
					"request ID = %q, want %q",
					responseRequestID,
					testCase.wantExactID,
				)
			}

			if testCase.wantGenerated {
				for _, value := range testCase.values {
					if responseRequestID == value {
						t.Fatalf(
							"expected generated request ID, got original %q",
							responseRequestID,
						)
					}
				}
			}

			if testCase.path == "/echo" &&
				string(body) != responseRequestID {
				t.Fatalf(
					"body = %q, response ID = %q",
					string(body),
					responseRequestID,
				)
			}
		})
	}
}
