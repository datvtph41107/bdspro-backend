package iusecase

import (
	"context"
	"crm/internal/domain"
)

type UserClient interface {
	GetProfileById(ctx context.Context, profileId uint64) (*domain.Profile, error)
	GetProfileByPhones(ctx context.Context, phones []string) (map[string]*domain.Profile, error)
	ValidateProfileIds(ctx context.Context, profileIds []uint64) (bool, error)
}
