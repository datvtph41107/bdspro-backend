package client

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFileClientGetFileUsesCompleteFileReadContract(t *testing.T) {
	var gotPath string
	var gotServiceName string
	var gotServiceAuth string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.EscapedPath()
		gotServiceName = r.Header.Get("X-Service-Name")
		gotServiceAuth = r.Header.Get("X-Service-Auth")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "payload")
	}))
	defer server.Close()

	client := &FileClient{
		baseURL:        server.URL,
		serviceAuthKey: "test-service-key",
		httpClient:     server.Client(),
	}

	body, err := client.GetFile(context.Background(), "pOpaque_ref-42")
	if err != nil {
		t.Fatalf("GetFile() error = %v", err)
	}
	if string(body) != "payload" {
		t.Fatalf("body = %q, want payload", body)
	}
	if gotPath != "/v1/file/load/pOpaque_ref-42/local" {
		t.Fatalf("path = %q", gotPath)
	}
	if gotServiceName != fileClientServiceName {
		t.Fatalf("X-Service-Name = %q", gotServiceName)
	}
	if gotServiceAuth == "" {
		t.Fatal("X-Service-Auth is empty")
	}
}
