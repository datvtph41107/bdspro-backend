package filehttp

import (
	"testing"
	"time"
)

func TestPublicURLOwnsFileDeliveryRoute(t *testing.T) {
	client, err := New(Config{
		BaseURL:        "http://127.0.0.1:8002",
		PublicBaseURL:  "https://files.qhpro.vn/",
		ServiceName:    "tqd-service",
		ServiceAuthKey: "test-secret",
		Timeout:        time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	got, err := client.PublicURL("pencoded/path")
	if err != nil {
		t.Fatal(err)
	}

	const want = "https://files.qhpro.vn/v1/file/load/pencoded%2Fpath/local"

	if got != want {
		t.Fatalf("PublicURL() = %q, want %q", got, want)
	}
}

func TestPublicURLRejectsEmptyPath(t *testing.T) {
	client, err := New(Config{
		BaseURL:        "http://127.0.0.1:8002",
		ServiceName:    "tqd-service",
		ServiceAuthKey: "test-secret",
		Timeout:        time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	if _, err := client.PublicURL(" "); err == nil {
		t.Fatal("PublicURL(empty) error = nil")
	}
}

func TestNewRejectsInvalidPublicBaseURL(t *testing.T) {
	if _, err := New(Config{
		BaseURL:        "http://127.0.0.1:8002",
		PublicBaseURL:  "not-a-url",
		ServiceName:    "tqd-service",
		ServiceAuthKey: "test-secret",
		Timeout:        time.Second,
	}); err == nil {
		t.Fatal("New(invalid PublicBaseURL) error = nil")
	}
}
