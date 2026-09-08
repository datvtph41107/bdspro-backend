package iusecase

import (
	"chat/internal/dto"
	"context"
)

type IUserClient interface {
	// GetUserById(ctx context.Context, id uint64) (*dto.UserProfile, error)
	GetUserById(ctx context.Context, id uint64) (*dto.UserProfile, error)
	GetUserByIds(ctx context.Context, ids []uint64) ([]*dto.UserProfile, error)
}
