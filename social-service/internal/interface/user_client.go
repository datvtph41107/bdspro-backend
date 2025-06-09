package iusecase

import (
	"context"
)

type UserClient interface {
	IsFriends(ctx context.Context, profileID uint64, friendIDs []uint64) (bool, error)
	// GetProfileByIds(ctx context.Context, ids []uint64) ([]*dto.UserProfile, error)
}
