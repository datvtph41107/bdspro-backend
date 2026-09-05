package client

import "testing"

func TestAssistantClientReadyFailsWhenProcessDidNotInjectClient(t *testing.T) {
	client := &AssistantClient{}
	if err := client.ready(); err == nil {
		t.Fatal("ready() error = nil, want unavailable error")
	}
}

func TestAssistantClientCloseDoesNotOwnProcessConnection(t *testing.T) {
	client := &AssistantClient{}
	if err := client.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}
