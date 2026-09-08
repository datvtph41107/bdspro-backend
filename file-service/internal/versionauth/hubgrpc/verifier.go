package hubgrpc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"file/internal/versionauth"
	hubpb "pb/types/hub"

	"google.golang.org/grpc"
)

type apiKeyClient interface {
	VerifyApiKey(
		context.Context,
		*hubpb.VerifyApiKeyRequest,
		...grpc.CallOption,
	) (*hubpb.VerifyApiKeyResponse, error)
}

// Verifier adapts Hub's API-key RPC to File-service's narrow semantic
// capability. It owns call timeout/error normalization only; the gRPC
// connection itself is process-owned by File-service composition.
type Verifier struct {
	client  apiKeyClient
	timeout time.Duration
}

var _ versionauth.Verifier = (*Verifier)(nil)

func New(client apiKeyClient, timeout time.Duration) *Verifier {
	return &Verifier{
		client:  client,
		timeout: timeout,
	}
}

func (v *Verifier) Verify(ctx context.Context, apiKey string) error {
	if v == nil {
		return versionauth.ErrUnavailable
	}

	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		return versionauth.ErrInvalid
	}

	if v.client == nil || v.timeout <= 0 {
		return versionauth.ErrUnavailable
	}

	callCtx, cancel := context.WithTimeout(ctx, v.timeout)
	defer cancel()

	resp, err := v.client.VerifyApiKey(
		callCtx,
		&hubpb.VerifyApiKeyRequest{ApiKey: apiKey},
	)
	if err != nil {
		return fmt.Errorf(
			"%w: Hub VerifyApiKey: %v",
			versionauth.ErrUnavailable,
			err,
		)
	}

	// A successful transport call without a response is not evidence that the
	// caller supplied an invalid key. The authority failed to produce a
	// decision, so classify it as unavailable.
	if resp == nil {
		return versionauth.ErrUnavailable
	}

	if !resp.GetValid() {
		message := strings.TrimSpace(resp.GetMessage())
		if message == "" {
			return versionauth.ErrInvalid
		}

		return fmt.Errorf("%w: %s", versionauth.ErrInvalid, message)
	}

	return nil
}
