package rpc

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func TestNewClientValidatesConstructionInputs(t *testing.T) {
	base := ClientConfig{
		Target:          "user:8201",
		BackoffMaxDelay: 5 * time.Second,
		Credentials:     insecure.NewCredentials(),
	}
	for name, mutate := range map[string]func(*ClientConfig){
		"target":      func(c *ClientConfig) { c.Target = " " },
		"backoff":     func(c *ClientConfig) { c.BackoffMaxDelay = 0 },
		"credentials": func(c *ClientConfig) { c.Credentials = nil },
	} {
		t.Run(name, func(t *testing.T) {
			cfg := base
			mutate(&cfg)
			if _, err := NewClient(cfg); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

func TestNewClientBuildsCallerOwnedChannel(t *testing.T) {
	conn, err := NewClient(ClientConfig{
		Target:          "user:8201",
		BackoffMaxDelay: 5 * time.Second,
		Credentials:     insecure.NewCredentials(),
	})
	if err != nil {
		t.Fatalf("NewClient() error = %v", err)
	}
	if conn == nil {
		t.Fatal("nil connection")
	}
	if err := conn.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}

func TestUnavailableReturnsUnavailable(t *testing.T) {
	conn := Unavailable(fmt.Errorf("dial config invalid"))
	err := conn.Invoke(context.Background(), "/test.Service/Call", nil, nil)
	if status.Code(err) != codes.Unavailable {
		t.Fatalf("code = %v", status.Code(err))
	}
	if !strings.Contains(err.Error(), "dial config invalid") {
		t.Fatalf("error = %v", err)
	}
}
