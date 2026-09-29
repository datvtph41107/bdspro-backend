package iusecase

import (
	"chat/internal/dto"
	"context"
)

type IUserClient interface {
	GetUserByIds(ctx context.Context, ids []uint64) ([]*dto.UserProfile, error)
}
