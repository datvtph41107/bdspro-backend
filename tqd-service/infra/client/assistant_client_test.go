package client

import (
	"errors"
	"testing"
	"time"

	qhprorpc "common/rpc"

	assistantpb "pb/types/assistant"

	"google.golang.org/grpc/credentials/insecure"
)

func TestAssistantClientReadyReturnsConstructionError(t *testing.T) {
	want := errors.New("construction failed")
	client := &AssistantClient{initErr: want}
	if err := client.ready(); !errors.Is(err, want) {
		t.Fatalf("ready() error = %v, want construction error", err)
	}
}

func TestAssistantClientCloseReleasesOwnedConnection(t *testing.T) {
	conn, err := qhprorpc.NewClient(qhprorpc.ClientConfig{
		Target:          "passthrough:///assistant-client-test",
		BackoffMaxDelay: time.Second,
		Credentials:     insecure.NewCredentials(),
	})
	if err != nil {
		t.Fatalf("create test client connection: %v", err)
	}

	client := &AssistantClient{
		conn:   conn,
		client: assistantpb.NewAssistantServiceClient(conn),
	}
	if err := client.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	if client.conn != nil || client.client != nil {
		t.Fatal("Close() retained assistant connection/client references")
	}
	if err := client.ready(); err == nil {
		t.Fatal("ready() error = nil after Close")
	}
	if err := client.Close(); err != nil {
		t.Fatalf("second Close() error = %v", err)
	}
}
