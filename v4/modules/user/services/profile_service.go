package services

import (
	"context"

	sharepb "pb/types/shared"
)

// ProfileService owns caller profile-info reads migrated into modules/user.
type ProfileService interface {
	GetInfoMe(ctx context.Context) (*sharepb.UserV3Proto, error)
	// GetProfileInfoV3 is GET /v3/user/profile/{access}/{id} (public profile + related stats).
	GetProfileInfoV3(ctx context.Context, id uint64, access string, timestamp int64) (*sharepb.TotalResponseProto, error)
}
