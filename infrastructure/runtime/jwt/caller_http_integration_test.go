package _jwt

import (
	"common/identity"

	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
)

type integrationAPIKeyVerifier struct{}

func (integrationAPIKeyVerifier) VerifyAPIKey(ctx context.Context, rawKey string) (VerifiedAPIKey, error) {
	select {
	case <-ctx.Done():
		return VerifiedAPIKey{}, ctx.Err()
	default:
	}
	if rawKey != "0123456789abcdef" {
		return VerifiedAPIKey{}, ErrAPIKeyInvalid
	}
	return VerifiedAPIKey{ID: 7, AppName: "integration-client"}, nil
}

type callerHTTPResponse struct {
	Kind       string `json:"kind"`
	APIKeyID   uint64 `json:"apiKeyId"`
	APIKeyApp  string `json:"apiKeyApp"`
	Canonical  string `json:"canonicalApiKey"`
	HasActor   bool   `json:"hasActor"`
	StatusText string `json:"status"`
}

func TestCallerResolutionOverHTTP(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowAnonymous := r.URL.Query().Get("public") == "1"
		var principal *Principal
		if value := r.Header.Get("X-Test-Verified-Profile"); value != "" {
			profileID, err := strconv.ParseUint(value, 10, 64)
			if err != nil || profileID == 0 {
				http.Error(w, "invalid test principal", http.StatusBadRequest)
				return
			}
			principal = &Principal{ProfileId: profileID}
		}

		resolution, err := ResolveHTTPCaller(
			r.Context(),
			r,
			principal,
			allowAnonymous,
			integrationAPIKeyVerifier{},
		)
		if err != nil {
			statusCode := http.StatusUnauthorized
			if errors.Is(err, ErrAPIKeyHeaderInvalid) {
				statusCode = http.StatusBadRequest
			}
			http.Error(w, err.Error(), statusCode)
			return
		}
		_ = json.NewEncoder(w).Encode(callerHTTPResponse{
			Kind:       string(resolution.Caller.Kind),
			APIKeyID:   resolution.Caller.APIKeyID,
			APIKeyApp:  resolution.Caller.APIKeyApp,
			Canonical:  r.Header.Get(APIKeyHeader),
			HasActor:   principal != nil,
			StatusText: "ok",
		})
	}))
	defer server.Close()

	t.Run("anonymous", func(t *testing.T) {
		resp := requestCallerHTTP(t, server.URL+"?public=1", nil)
		if resp.Kind != string(identity.CallerAnonymous) || resp.APIKeyID != 0 || resp.HasActor {
			t.Fatalf("response = %#v", resp)
		}
	})

	t.Run("api key", func(t *testing.T) {
		headers := http.Header{LegacyAPIKeyHeader: []string{"0123456789abcdef"}}
		resp := requestCallerHTTP(t, server.URL, headers)
		if resp.Kind != string(identity.CallerAPIKey) || resp.APIKeyID != 7 || resp.Canonical == "" || resp.HasActor {
			t.Fatalf("response = %#v", resp)
		}
	})

	t.Run("user", func(t *testing.T) {
		headers := http.Header{"X-Test-Verified-Profile": []string{"42"}}
		resp := requestCallerHTTP(t, server.URL, headers)
		if resp.Kind != string(identity.CallerUser) || !resp.HasActor || resp.APIKeyID != 0 {
			t.Fatalf("response = %#v", resp)
		}
	})

	t.Run("user through api key client", func(t *testing.T) {
		headers := http.Header{
			"X-Test-Verified-Profile": []string{"42"},
			APIKeyHeader:              []string{"0123456789abcdef"},
		}
		resp := requestCallerHTTP(t, server.URL, headers)
		if resp.Kind != string(identity.CallerUser) || !resp.HasActor || resp.APIKeyID != 7 {
			t.Fatalf("response = %#v", resp)
		}
	})

	t.Run("concurrent anonymous", func(t *testing.T) {
		const total = 100
		var wait sync.WaitGroup
		errorsCh := make(chan error, total)
		for index := 0; index < total; index++ {
			wait.Add(1)
			go func() {
				defer wait.Done()
				req, err := http.NewRequest(http.MethodGet, server.URL+"?public=1", nil)
				if err != nil {
					errorsCh <- err
					return
				}
				response, err := http.DefaultClient.Do(req)
				if err != nil {
					errorsCh <- err
					return
				}
				defer response.Body.Close()
				if response.StatusCode != http.StatusOK {
					errorsCh <- errors.New("unexpected status")
				}
			}()
		}
		wait.Wait()
		close(errorsCh)
		for err := range errorsCh {
			if err != nil {
				t.Fatal(err)
			}
		}
	})
}

func requestCallerHTTP(t *testing.T, url string, headers http.Header) callerHTTPResponse {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal(err)
	}
	for key, values := range headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", response.StatusCode)
	}
	var result callerHTTPResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	return result
}
