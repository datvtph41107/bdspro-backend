package services

import (
	"context"

	sharepb "pb/types/shared"
)

// ProfileService owns caller profile-info reads migrated into modules/user.
type ProfileService interface {
	GetInfoMe(ctx context.Context) (*sharepb.UserV3Proto, error)
}
