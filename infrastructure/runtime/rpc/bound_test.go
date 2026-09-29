package rpc

import (
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

func TestNewBoundClientReturnsConstructionError(t *testing.T) {
	type boundClient struct{}

	client, cleanup, err := NewBoundClient(
		ClientConfig{},
		func(grpc.ClientConnInterface) *boundClient {
			t.Fatal("bind must not run when channel construction fails")
			return nil
		},
	)

	if err == nil {
		t.Fatal("expected construction error")
	}
	if client != nil {
		t.Fatal("client must be nil on construction error")
	}
	if cleanup != nil {
		t.Fatal("cleanup must be nil on construction error")
	}
}

func TestNewBoundClientBindsAndOwnsCleanup(t *testing.T) {
	var conn *grpc.ClientConn

	client, cleanup, err := NewBoundClient(
		ClientConfig{
			Target:          "user:8201",
			BackoffMaxDelay: 5 * time.Second,
			Credentials:     insecure.NewCredentials(),
		},
		func(cc grpc.ClientConnInterface) grpc.ClientConnInterface {
			var ok bool
			conn, ok = cc.(*grpc.ClientConn)
			if !ok {
				t.Fatalf("connection type = %T, want *grpc.ClientConn", cc)
			}
			return cc
		},
	)
	if err != nil {
		t.Fatalf("NewBoundClient() error = %v", err)
	}
	if client == nil {
		t.Fatal("bound client is nil")
	}
	if cleanup == nil {
		t.Fatal("cleanup is nil")
	}
	if conn == nil {
		t.Fatal("bind was not called")
	}

	cleanup()

	if state := conn.GetState(); state != connectivity.Shutdown {
		t.Fatalf("connection state after cleanup = %v, want SHUTDOWN", state)
	}
}
