package repo

import (
	_dto "common/domain/dto"
	_enum "common/domain/enum"
	"context"
	"time"
	"user/enums"
	"user/internal/domain/user"
	"user/internal/dto"
)

type IProfileRepo interface {
	GlobalSearchProfile(c context.Context, text string, pageable _dto.Pagable) (*[]user.UserProfileSearch, int64, error)
	FriendSearchProfile(c context.Context, profileId uint64, dto dto.ProfileSearch) (*[]user.UserProfileSearch, int64, error)
	UpdateProfile(c context.Context, profileId uint64, dto *dto.UserInfoRequest) error
	GetByProfileID(profileId uint64) (*user.UserProfileEntity, error)
	GetProfileByIds(ids []uint64) ([]*user.UserProfileEntity, error)
	GetProfileByPhones(phones []string) ([]*user.UserProfileEntity, error)
	ListAllUsers(ctx context.Context, req *dto.ListUsersRequest) ([]user.UserProfileEntity, uint32, error)
	GetUserDetailByID(ctx context.Context, profileID uint64) (*user.UserProfileEntity, error)
	IsSystemRootProfile(ctx context.Context, profileID uint64) (bool, error)
	CreateUserProfile(ctx context.Context, profile *user.UserProfileEntity) (*user.UserProfileEntity, error)
	UpdateUserProfile(ctx context.Context, profile *user.UserProfileEntity) (*user.UserProfileEntity, error)
	DeleteUserProfile(ctx context.Context, profileID uint64, deletedBy *uint64, reason string) error
	UpdateUserProfileStatus(ctx context.Context, profileID uint64, status _enum.EUserStatus) error
	GetByPhone(phone string) (*user.UserProfileEntity, error)
	SearchByPhoneOrName(ctx context.Context, text string) ([]user.UserProfileEntity, error)
	UpdatePerson(c context.Context, profileId uint64, req *_dto.UserV3DTO) error

	PatchProfileV3(c context.Context, profileId uint64, req *_dto.UserV3DTO) error
	UpdatePrivacySetting(ctx context.Context, profileID uint64, visibility uint8, profileVisibility uint8, statusOnline enums.EOnlineStatus) error
	UpdateLastSeen(ctx context.Context, profileID uint64) error

	// Thống kê user mới theo ngày (created_at) từ fromTime
	GetNewUsersDailyStats(ctx context.Context, fromTime time.Time) ([]dto.UserGrowthDailyItem, error)
	// Thống kê user hoạt động theo ngày (login/refresh activity) từ fromTime
	GetActiveUsersDailyStats(ctx context.Context, fromTime time.Time) ([]dto.UserGrowthDailyItem, error)
	// Thống kê user mới theo giờ (created_at) từ fromTime
	GetNewUsersHourlyStats(ctx context.Context, fromTime time.Time) ([]dto.UserGrowthDailyItem, error)
	// Thống kê user hoạt động theo giờ (login/refresh activity) từ fromTime
	GetActiveUsersHourlyStats(ctx context.Context, fromTime time.Time) ([]dto.UserGrowthDailyItem, error)
}
