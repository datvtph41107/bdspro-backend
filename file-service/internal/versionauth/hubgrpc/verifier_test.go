package hubgrpc

import (
	"context"
	"errors"
	"testing"
	"time"

	hubpb "pb/types/hub"

	"google.golang.org/grpc"
)

type fakeAPIKeyClient struct {
	response *hubpb.VerifyApiKeyResponse
	err      error
	request  *hubpb.VerifyApiKeyRequest
	wait     bool
}

func (f *fakeAPIKeyClient) VerifyApiKey(ctx context.Context, request *hubpb.VerifyApiKeyRequest, _ ...grpc.CallOption) (*hubpb.VerifyApiKeyResponse, error) {
	f.request = request
	if f.wait {
		<-ctx.Done()
		return nil, ctx.Err()
	}
	return f.response, f.err
}

func TestVerifyAcceptsValidHubEvidence(t *testing.T) {
	client := &fakeAPIKeyClient{response: &hubpb.VerifyApiKeyResponse{Valid: true}}
	verifier := New(client, time.Second)

	if err := verifier.Verify(context.Background(), " key-1 "); err != nil {
		t.Fatalf("Verify() error = %v", err)
	}
	if client.request == nil || client.request.GetApiKey() != "key-1" {
		t.Fatalf("request = %#v", client.request)
	}
}

func TestVerifyPreservesHubDenialMessage(t *testing.T) {
	client := &fakeAPIKeyClient{response: &hubpb.VerifyApiKeyResponse{Valid: false, Message: "revoked"}}
	verifier := New(client, time.Second)

	err := verifier.Verify(context.Background(), "key-1")
	if err == nil || err.Error() != "revoked" {
		t.Fatalf("Verify() error = %v", err)
	}
}

func TestVerifyFailsClosedWhenClientUnavailable(t *testing.T) {
	if err := New(nil, time.Second).Verify(context.Background(), "key-1"); !errors.Is(err, ErrUnavailable) {
		t.Fatalf("Verify() error = %v, want ErrUnavailable", err)
	}
}

func TestVerifyHonorsConfiguredTimeout(t *testing.T) {
	client := &fakeAPIKeyClient{wait: true}
	verifier := New(client, time.Millisecond)

	if err := verifier.Verify(context.Background(), "key-1"); err == nil {
		t.Fatal("Verify() error = nil, want timeout")
	}
}
