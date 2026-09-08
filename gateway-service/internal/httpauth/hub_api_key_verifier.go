package httpauth

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	hubpb "pb/types/hub"

	"google.golang.org/grpc"
)

type HubAPIKeyVerifier struct {
	client  hubpb.ApiKeyServiceClient
	timeout time.Duration
}

func NewHubAPIKeyVerifier(conn *grpc.ClientConn) (*HubAPIKeyVerifier, error) {
	if conn == nil {
		return nil, ErrAPIKeyUnavailable
	}
	return &HubAPIKeyVerifier{client: hubpb.NewApiKeyServiceClient(conn), timeout: 2 * time.Second}, nil
}

func (v *HubAPIKeyVerifier) VerifyAPIKey(ctx context.Context, raw string) (VerifiedAPIKey, error) {
	if ctx == nil || v == nil || v.client == nil {
		return VerifiedAPIKey{}, ErrAPIKeyUnavailable
	}

	timeout := v.timeout
	if timeout <= 0 {
		timeout = 2 * time.Second
	}

	callCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	resp, err := v.client.VerifyApiKey(callCtx, &hubpb.VerifyApiKeyRequest{ApiKey: raw})
	if err != nil {
		switch {
		case errors.Is(err, context.Canceled), errors.Is(callCtx.Err(), context.Canceled):
			return VerifiedAPIKey{}, context.Canceled
		case errors.Is(err, context.DeadlineExceeded), errors.Is(callCtx.Err(), context.DeadlineExceeded):
			return VerifiedAPIKey{}, context.DeadlineExceeded
		default:
			return VerifiedAPIKey{}, fmt.Errorf("%w: %v", ErrAPIKeyUnavailable, err)
		}
	}
	if resp == nil || !resp.GetValid() || resp.GetData() == nil {
		return VerifiedAPIKey{}, ErrAPIKeyInvalid
	}

	data := resp.GetData()
	app := strings.TrimSpace(data.GetAppName())
	if data.GetId() == 0 || app == "" {
		return VerifiedAPIKey{}, ErrAPIKeyInvalid
	}

	return VerifiedAPIKey{ID: data.GetId(), AppName: app}, nil
}
