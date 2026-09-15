package middlewares

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	commonjwt "common/jwt"
)

func TestWriteErrorResponsePreservesHistoricalJWTContract(
	t *testing.T,
) {
	tests := []struct {
		name    string
		err     error
		message string
	}{
		{
			name: "authorization header missing",
			err: commonjwt.
				AuthorizationHeaderMissingFault(),
			message: "Authorization header is missing",
		},
		{
			name:    "token missing",
			err:     commonjwt.TokenMissingFault(),
			message: "Missing token",
		},
		{
			name: "invalid or expired token",
			err: commonjwt.
				InvalidOrExpiredTokenFault(nil),
			message: "Invalid or expired token",
		},
	}

	for _, test := range tests {
		t.Run(
			test.name,
			func(t *testing.T) {
				recorder :=
					httptest.NewRecorder()

				writeErrorResponse(
					context.Background(),
					recorder,
					test.err,
				)

				if recorder.Code !=
					http.StatusUnauthorized {
					t.Fatalf(
						"status = %d, want %d",
						recorder.Code,
						http.StatusUnauthorized,
					)
				}

				if got :=
					recorder.Header().Get(
						"Content-Type",
					); got != "application/json" {
					t.Fatalf(
						"content type = %q",
						got,
					)
				}

				want :=
					"{\"code\":401," +
						"\"message\":\"" +
						test.message +
						"\",\"data\":{}}\n"

				if got :=
					recorder.Body.String(); got != want {
					t.Fatalf(
						"body = %q, want %q",
						got,
						want,
					)
				}
			},
		)
	}
}
