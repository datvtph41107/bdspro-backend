package providers

import (
	_dto "common/domain/dto"
	"context"
	"time"
	"user/internal/dto"
)

// ProfileProvider interface cho việc xử lý logic liên quan đến profile
type ProfileProvider interface {
	GetByProfileID(ctx context.Context, profileID uint64) (*dto.ProfileDTO, error)
	GetByProfileIDV3(ctx context.Context, profileID uint64) (*_dto.UserV3DTO, error)
	SaveProfile(ctx context.Context, profile *dto.ProfileDTO) error
	CreateProfile(ctx context.Context, authID uint64, planID *uint64, name, email, phone, avatar string) (*dto.ProfileDTO, error)
	UpdateProfile(ctx context.Context, profile *dto.ProfileDTO) error
	SoftDeleteProfile(ctx context.Context, profileID uint64) error
	HardDeleteProfile(ctx context.Context, profileID uint64) error
	CheckAccountStatus(ctx context.Context, profileID uint64) (isLocked bool, canLogin bool, err error)
	LockAccount(ctx context.Context, profileID uint64) error

	// Profile utilities
	MakeUserProfileEntity(authID uint64, planID *uint64, name, email, phone, avatar string) *dto.ProfileDTO
	GetPlanInfo(profileID uint64) (*uint64, *time.Time)
	GetMapProfileByIds(ctx context.Context, ids []uint64) (map[uint64]*dto.ProfileDTO, error)
}
