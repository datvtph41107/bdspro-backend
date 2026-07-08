package provider

import (
	"context"
	"crm/internal/domain"
	sharepb "pb/types/shared"
)

type UserClient interface {
	GetProfileById(ctx context.Context, profileId uint64) (*domain.Profile, error)
	GetProfileWithPhoneById(ctx context.Context, profileId uint64) (*domain.Profile, error)
	GetProfileByPhones(ctx context.Context, phones []string) (map[string]*domain.Profile, error)
	ValidateProfileIds(ctx context.Context, profileIds []uint64) (bool, error)
	GetProfileByIds(ctx context.Context, ids []uint64) ([]*sharepb.ProfileItem, error)
}