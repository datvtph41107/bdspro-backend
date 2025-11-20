package iprovider

import "context"

// IAuthClient interface cho auth client
type IAuthClient interface {
	GetPushTokensByProfileId(ctx context.Context, profileID uint64) ([]string, error)
}
