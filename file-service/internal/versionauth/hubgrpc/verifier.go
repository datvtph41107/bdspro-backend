package hubgrpc

import (
	"context"
	"errors"
	"strings"
	"time"

	hubpb "pb/types/hub"

	"google.golang.org/grpc"
)

var ErrUnavailable = errors.New("hub api key verifier unavailable")

type apiKeyClient interface {
	VerifyApiKey(context.Context, *hubpb.VerifyApiKeyRequest, ...grpc.CallOption) (*hubpb.VerifyApiKeyResponse, error)
}

// Verifier adapts Hub's API-key RPC to File-service's narrow semantic
// capability. It owns call timeout/error normalization only; the gRPC
// connection itself is process-owned by File-service composition.
type Verifier struct {
	client  apiKeyClient
	timeout time.Duration
}

func New(client apiKeyClient, timeout time.Duration) *Verifier {
	return &Verifier{client: client, timeout: timeout}
}

func (v *Verifier) Verify(ctx context.Context, apiKey string) error {
	if v == nil || v.client == nil || v.timeout <= 0 {
		return ErrUnavailable
	}
	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		return errors.New("API key không hợp lệ")
	}

	callCtx, cancel := context.WithTimeout(ctx, v.timeout)
	defer cancel()

	resp, err := v.client.VerifyApiKey(callCtx, &hubpb.VerifyApiKeyRequest{ApiKey: apiKey})
	if err != nil {
		return err
	}
	if resp == nil || !resp.GetValid() {
		message := "API key không hợp lệ"
		if resp != nil && strings.TrimSpace(resp.GetMessage()) != "" {
			message = strings.TrimSpace(resp.GetMessage())
		}
		return errors.New(message)
	}
	return nil
}
