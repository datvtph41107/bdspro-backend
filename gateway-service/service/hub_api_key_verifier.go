package service

import (
	"common"
	_jwt "common/jwt"
	"context"
	"fmt"
	hubpb "pb/types/hub"
	"strings"
	"time"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

/**
 * HubAPIKeyVerifier resolves API key identity through Hub service.
 */
type HubAPIKeyVerifier struct {
	conn    *grpc.ClientConn
	client  hubpb.ApiKeyServiceClient
	timeout time.Duration
}

/**
 * NewHubAPIKeyVerifier creates the Gateway API key authority client.
 */
func NewHubAPIKeyVerifier(env string, opts []grpc.DialOption) (*HubAPIKeyVerifier, error) {
	port := viper.GetString("service.port.hub")
	endpoint := fmt.Sprintf("%s:%s", common.GetPrefixProtobufUrl(env, "hub"), port)
	conn, err := grpc.Dial(endpoint, opts...)
	if err != nil {
		return nil, fmt.Errorf("dial hub api key service: %w", err)
	}
	return &HubAPIKeyVerifier{
		conn:    conn,
		client:  hubpb.NewApiKeyServiceClient(conn),
		timeout: 2 * time.Second,
	}, nil
}

/**
 * VerifyAPIKey returns non-secret identity for one valid key.
 */
func (v *HubAPIKeyVerifier) VerifyAPIKey(ctx context.Context, rawKey string) (_jwt.VerifiedAPIKey, error) {
	if v == nil || v.client == nil {
		return _jwt.VerifiedAPIKey{}, _jwt.ErrAPIKeyUnavailable
	}

	timeout := v.timeout
	if timeout <= 0 {
		timeout = 2 * time.Second
	}
	callCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	resp, err := v.client.VerifyApiKey(callCtx, &hubpb.VerifyApiKeyRequest{ApiKey: rawKey})
	if err != nil {
		return _jwt.VerifiedAPIKey{}, fmt.Errorf("%w: %v", _jwt.ErrAPIKeyUnavailable, err)
	}
	if resp == nil || !resp.GetValid() || resp.GetData() == nil {
		return _jwt.VerifiedAPIKey{}, _jwt.ErrAPIKeyInvalid
	}

	data := resp.GetData()
	appName := strings.TrimSpace(data.GetAppName())
	if data.GetId() == 0 || appName == "" {
		return _jwt.VerifiedAPIKey{}, _jwt.ErrAPIKeyInvalid
	}
	return _jwt.VerifiedAPIKey{ID: data.GetId(), AppName: appName}, nil
}

/**
 * Close releases the Hub connection.
 */
func (v *HubAPIKeyVerifier) Close() error {
	if v == nil || v.conn == nil {
		return nil
	}
	return v.conn.Close()
}
