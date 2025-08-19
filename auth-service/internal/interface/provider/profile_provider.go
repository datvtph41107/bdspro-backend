package provider

import (
	"auth/internal/dto"
	"context"
	"time"
)

// ProfileProvider interface cho việc xử lý logic liên quan đến profile
type ProfileProvider interface {
	GetByProfileID(ctx context.Context, profileID uint64) (*dto.ProfileDTO, error)
	SaveProfile(ctx context.Context, profile *dto.ProfileDTO) error
	CreateProfile(ctx context.Context, authID uint64, planID *uint64, name, email, phone, avatar string) (*dto.ProfileDTO, error)
	UpdateProfile(ctx context.Context, profile *dto.ProfileDTO) error
	SoftDeleteProfile(ctx context.Context, profileID uint64) error

	// Profile utilities
	MakeUserProfileEntity(authID uint64, planID *uint64, name, email, phone, avatar string) *dto.ProfileDTO
	GetPlanInfo(profileID uint64) (*uint64, *time.Time)
}
