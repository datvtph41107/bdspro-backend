package _jwt

import (
	"common/identity"

	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClassifyCallerMatrix(t *testing.T) {
	t.Parallel()

	apiKey := &VerifiedAPIKey{ID: 7, AppName: " partner-app "}
	cases := []struct {
		name  string
		input CallerClassificationInput
		want  identity.Caller
	}{
		{
			name:  "anonymous public route",
			input: CallerClassificationInput{AllowAnonymous: true},
			want:  identity.Caller{Kind: identity.CallerAnonymous},
		},
		{
			name:  "api key only",
			input: CallerClassificationInput{APIKey: apiKey},
			want:  identity.Caller{Kind: identity.CallerAPIKey, APIKeyID: 7, APIKeyApp: "partner-app"},
		},
		{
			name:  "jwt user",
			input: CallerClassificationInput{Principal: &Principal{ProfileId: 42}},
			want:  identity.Caller{Kind: identity.CallerUser},
		},
		{
			name: "jwt user through verified api client",
			input: CallerClassificationInput{
				Principal: &Principal{ProfileId: 42},
				APIKey:    apiKey,
			},
			want: identity.Caller{Kind: identity.CallerUser, APIKeyID: 7, APIKeyApp: "partner-app"},
		},
	}

	for _, test := range cases {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got, err := ClassifyCaller(test.input)
			if err != nil {
				t.Fatalf("ClassifyCaller() error = %v", err)
			}
			if got != test.want {
				t.Fatalf("ClassifyCaller() = %#v, want %#v", got, test.want)
			}
		})
	}
}

func TestClassifyCallerRejectsMissingPrivateAuthentication(t *testing.T) {
	t.Parallel()

	_, err := ClassifyCaller(CallerClassificationInput{})
	if !errors.Is(err, ErrAuthenticationRequired) {
		t.Fatalf("error = %v, want ErrAuthenticationRequired", err)
	}
}

func TestClassifyCallerRejectsVerifiedTokenWithoutActorIdentity(t *testing.T) {
	t.Parallel()

	_, err := ClassifyCaller(CallerClassificationInput{
		Principal: &Principal{Role: "admin", Type: "ACCESS"},
	})
	if !errors.Is(err, ErrCallerClassification) {
		t.Fatalf("error = %v, want ErrCallerClassification", err)
	}
}

func TestClassifyCallerRejectsInvalidVerifiedAPIKeyIdentity(t *testing.T) {
	t.Parallel()

	for _, key := range []*VerifiedAPIKey{
		{ID: 0, AppName: "partner"},
		{ID: 7, AppName: ""},
		{ID: 7, AppName: "partner\nroot"},
	} {
		if _, err := ClassifyCaller(CallerClassificationInput{APIKey: key}); !errors.Is(err, ErrCallerClassification) {
			t.Fatalf("ClassifyCaller(%#v) error = %v", key, err)
		}
	}
}

type callerTestVerifier struct {
	result VerifiedAPIKey
	err    error
	calls  int
}

func (f *callerTestVerifier) VerifyAPIKey(_ context.Context, _ string) (VerifiedAPIKey, error) {
	f.calls++
	return f.result, f.err
}

func TestResolveHTTPCallerCanonicalizesLegacyAPIKey(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/private", nil)
	req.Header.Set(LegacyAPIKeyHeader, "0123456789abcdef")
	verifier := &callerTestVerifier{result: VerifiedAPIKey{ID: 9, AppName: "partner"}}

	resolution, err := ResolveHTTPCaller(req.Context(), req, nil, false, verifier)
	if err != nil {
		t.Fatalf("ResolveHTTPCaller() error = %v", err)
	}
	if resolution.Caller.Kind != identity.CallerAPIKey || resolution.Caller.APIKeyID != 9 {
		t.Fatalf("caller = %#v", resolution.Caller)
	}
	if got := req.Header.Get(APIKeyHeader); got != "0123456789abcdef" {
		t.Fatalf("canonical header = %q", got)
	}
	if got := req.Header.Get(LegacyAPIKeyHeader); got != "" {
		t.Fatalf("legacy header remained = %q", got)
	}
	if verifier.calls != 1 {
		t.Fatalf("verifier calls = %d, want 1", verifier.calls)
	}
}

func TestResolveHTTPCallerDoesNotCallVerifierWithoutAPIKey(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/public", nil)
	verifier := &callerTestVerifier{result: VerifiedAPIKey{ID: 9, AppName: "partner"}}
	resolution, err := ResolveHTTPCaller(req.Context(), req, nil, true, verifier)
	if err != nil {
		t.Fatalf("ResolveHTTPCaller() error = %v", err)
	}
	if resolution.Caller.Kind != identity.CallerAnonymous || verifier.calls != 0 {
		t.Fatalf("resolution = %#v, calls = %d", resolution, verifier.calls)
	}
}

func TestResolveHTTPCallerRequiresAuthorityForPresentedAPIKey(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/private", nil)
	req.Header.Set(APIKeyHeader, "0123456789abcdef")
	_, err := ResolveHTTPCaller(req.Context(), req, nil, false, nil)
	if !errors.Is(err, ErrAPIKeyUnavailable) {
		t.Fatalf("error = %v, want ErrAPIKeyUnavailable", err)
	}
}

func TestResolveHTTPCallerPreservesCancellation(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/private", nil)
	req.Header.Set(APIKeyHeader, "0123456789abcdef")
	verifier := &callerTestVerifier{err: context.Canceled}
	_, err := ResolveHTTPCaller(req.Context(), req, nil, false, verifier)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want context.Canceled", err)
	}
}
