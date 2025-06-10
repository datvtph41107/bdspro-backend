package iusecase

import (
	"context"
	"social/internal/dto"
)

type UserClient interface {
	IsFriends(ctx context.Context, profileID uint64, friendIDs []uint64) (bool, error)
	GetUser(ctx context.Context, profileID uint64) (*dto.UserProfile, error)
	GetUsers(ctx context.Context, profileIDs []uint64) ([]*dto.UserProfile, error)
	// GetProfileByIds(ctx context.Context, ids []uint64) ([]*dto.UserProfile, error)
}
