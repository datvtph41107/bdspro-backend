package service

import (
	_jwt "common/jwt"
	"context"
	"errors"
	hubpb "pb/types/hub"
	"testing"
	"time"

	"google.golang.org/grpc"
)

type fakeHubAPIKeyClient struct {
	hubpb.ApiKeyServiceClient
	resp *hubpb.VerifyApiKeyResponse
	err  error
}

func (f fakeHubAPIKeyClient) VerifyApiKey(
	context.Context,
	*hubpb.VerifyApiKeyRequest,
	...grpc.CallOption,
) (*hubpb.VerifyApiKeyResponse, error) {
	return f.resp, f.err
}

func TestHubAPIKeyVerifierReturnsIdentity(t *testing.T) {
	t.Parallel()

	verifier := &HubAPIKeyVerifier{client: fakeHubAPIKeyClient{resp: &hubpb.VerifyApiKeyResponse{
		Valid: true,
		Data:  &hubpb.ApiKeyProto{Id: 7, AppName: "release-uploader"},
	}}}
	got, err := verifier.VerifyAPIKey(context.Background(), "0123456789abcdef")
	if err != nil {
		t.Fatalf("VerifyAPIKey() error = %v", err)
	}
	if got.ID != 7 || got.AppName != "release-uploader" {
		t.Fatalf("identity = %#v", got)
	}
}

func TestHubAPIKeyVerifierRejectsInvalidResponse(t *testing.T) {
	t.Parallel()

	verifier := &HubAPIKeyVerifier{client: fakeHubAPIKeyClient{resp: &hubpb.VerifyApiKeyResponse{Valid: false}}}
	_, err := verifier.VerifyAPIKey(context.Background(), "0123456789abcdef")
	if !errors.Is(err, _jwt.ErrAPIKeyInvalid) {
		t.Fatalf("error = %v, want ErrAPIKeyInvalid", err)
	}
}

func TestHubAPIKeyVerifierMapsTransportFailure(t *testing.T) {
	t.Parallel()

	verifier := &HubAPIKeyVerifier{client: fakeHubAPIKeyClient{err: errors.New("hub unavailable")}}
	_, err := verifier.VerifyAPIKey(context.Background(), "0123456789abcdef")
	if !errors.Is(err, _jwt.ErrAPIKeyUnavailable) {
		t.Fatalf("error = %v, want ErrAPIKeyUnavailable", err)
	}
}

type blockingHubAPIKeyClient struct {
	hubpb.ApiKeyServiceClient
}

func (blockingHubAPIKeyClient) VerifyApiKey(
	ctx context.Context,
	_ *hubpb.VerifyApiKeyRequest,
	_ ...grpc.CallOption,
) (*hubpb.VerifyApiKeyResponse, error) {
	<-ctx.Done()
	return nil, ctx.Err()
}

func TestHubAPIKeyVerifierPreservesParentCancellation(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	verifier := &HubAPIKeyVerifier{
		client:  blockingHubAPIKeyClient{},
		timeout: time.Second,
	}
	_, err := verifier.VerifyAPIKey(ctx, "0123456789abcdef")
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("error = %v, want context.Canceled", err)
	}
}

func TestHubAPIKeyVerifierPreservesDeadlineExceeded(t *testing.T) {
	t.Parallel()

	verifier := &HubAPIKeyVerifier{
		client:  blockingHubAPIKeyClient{},
		timeout: time.Millisecond,
	}
	_, err := verifier.VerifyAPIKey(context.Background(), "0123456789abcdef")
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("error = %v, want context.DeadlineExceeded", err)
	}
}
