package iusecase

import (
	"context"
	"crm/internal/domain"
)

type UserClient interface {
	GetProfileById(ctx context.Context, profileId uint64) (*domain.Profile, error)
}
