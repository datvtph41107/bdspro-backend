package hubgrpc

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"file/internal/versionauth"
	hubpb "pb/types/hub"

	"google.golang.org/grpc"
)

type fakeAPIKeyClient struct {
	response *hubpb.VerifyApiKeyResponse
	err      error
	request  *hubpb.VerifyApiKeyRequest
	wait     bool
}

func (f *fakeAPIKeyClient) VerifyApiKey(
	ctx context.Context,
	request *hubpb.VerifyApiKeyRequest,
	_ ...grpc.CallOption,
) (*hubpb.VerifyApiKeyResponse, error) {
	f.request = request

	if f.wait {
		<-ctx.Done()
		return nil, ctx.Err()
	}

	return f.response, f.err
}

func TestVerifyAcceptsValidHubEvidence(t *testing.T) {
	client := &fakeAPIKeyClient{
		response: &hubpb.VerifyApiKeyResponse{Valid: true},
	}
	verifier := New(client, time.Second)

	if err := verifier.Verify(context.Background(), " key-1 "); err != nil {
		t.Fatalf("Verify() error = %v", err)
	}

	if client.request == nil || client.request.GetApiKey() != "key-1" {
		t.Fatalf("request = %#v", client.request)
	}
}

func TestVerifyRejectsEmptyAPIKeyAsInvalid(t *testing.T) {
	client := &fakeAPIKeyClient{
		response: &hubpb.VerifyApiKeyResponse{Valid: true},
	}
	verifier := New(client, time.Second)

	err := verifier.Verify(context.Background(), " \t ")

	if !errors.Is(err, versionauth.ErrInvalid) {
		t.Fatalf("Verify() error = %v, want ErrInvalid", err)
	}

	if client.request != nil {
		t.Fatalf("VerifyApiKey() called with request = %#v", client.request)
	}
}

func TestVerifyPreservesHubDenialMessage(t *testing.T) {
	client := &fakeAPIKeyClient{
		response: &hubpb.VerifyApiKeyResponse{
			Valid:   false,
			Message: "revoked",
		},
	}
	verifier := New(client, time.Second)

	err := verifier.Verify(context.Background(), "key-1")

	if !errors.Is(err, versionauth.ErrInvalid) {
		t.Fatalf("Verify() error = %v, want ErrInvalid", err)
	}

	if err == nil || !strings.Contains(err.Error(), "revoked") {
		t.Fatalf("Verify() error = %v, want Hub denial message", err)
	}
}

func TestVerifyFailsClosedWhenClientUnavailable(t *testing.T) {
	err := New(nil, time.Second).Verify(context.Background(), "key-1")

	if !errors.Is(err, versionauth.ErrUnavailable) {
		t.Fatalf("Verify() error = %v, want ErrUnavailable", err)
	}
}

func TestVerifyMapsHubTransportFailureToUnavailable(t *testing.T) {
	client := &fakeAPIKeyClient{
		err: errors.New("hub down"),
	}
	verifier := New(client, time.Second)

	err := verifier.Verify(context.Background(), "key-1")

	if !errors.Is(err, versionauth.ErrUnavailable) {
		t.Fatalf("Verify() error = %v, want ErrUnavailable", err)
	}

	if !strings.Contains(err.Error(), "hub down") {
		t.Fatalf("Verify() error = %v, want transport evidence", err)
	}
}

func TestVerifyMapsNilHubResponseToUnavailable(t *testing.T) {
	client := &fakeAPIKeyClient{}
	verifier := New(client, time.Second)

	err := verifier.Verify(context.Background(), "key-1")

	if !errors.Is(err, versionauth.ErrUnavailable) {
		t.Fatalf("Verify() error = %v, want ErrUnavailable", err)
	}
}

func TestVerifyHonorsConfiguredTimeout(t *testing.T) {
	client := &fakeAPIKeyClient{wait: true}
	verifier := New(client, time.Millisecond)

	err := verifier.Verify(context.Background(), "key-1")

	if !errors.Is(err, versionauth.ErrUnavailable) {
		t.Fatalf("Verify() error = %v, want ErrUnavailable", err)
	}
}
