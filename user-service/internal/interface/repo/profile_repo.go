package repo

import (
	_dto "common/domain/dto"
	_enum "common/domain/enum"
	"context"
	"user/enums"
	"user/internal/dto"
	"user/internal/models"
)

type IProfileRepo interface {
	GlobalSearchProfile(c context.Context, text string, pageable _dto.Pagable) (*[]models.UserProfileSearch, int64, error)
	FriendSearchProfile(c context.Context, profileId uint64, dto dto.ProfileSearch) (*[]models.UserProfileSearch, int64, error)
	UpdateProfile(c context.Context, profileId uint64, dto *dto.UserInfoRequest) error
	GetByProfileID(profileId uint64) (*models.UserProfileEntity, error)
	GetProfileByIds(ids []uint64) ([]*models.UserProfileEntity, error)
	GetProfileByPhones(phones []string) ([]*models.UserProfileEntity, error)
	ListAllUsers(ctx context.Context, req *dto.ListUsersRequest) ([]models.UserProfileEntity, uint32, error)
	GetUserDetailByID(ctx context.Context, profileID uint64) (*models.UserProfileEntity, error)
	IsSystemRootProfile(ctx context.Context, profileID uint64) (bool, error)
	CreateUserProfile(ctx context.Context, profile *models.UserProfileEntity) (*models.UserProfileEntity, error)
	UpdateUserProfile(ctx context.Context, profile *models.UserProfileEntity) (*models.UserProfileEntity, error)
	DeleteUserProfile(ctx context.Context, profileID uint64, deletedBy *uint64, reason string) error
	UpdateUserProfileStatus(ctx context.Context, profileID uint64, status _enum.EUserStatus) error
	GetByPhone(phone string) (*models.UserProfileEntity, error)
	SearchByPhoneOrName(ctx context.Context, text string) ([]models.UserProfileEntity, error)
	UpdatePerson(c context.Context, profileId uint64, req *_dto.UserV3DTO) error

	PatchProfileV3(c context.Context, profileId uint64, req *_dto.UserV3DTO) error
	UpdatePrivacySetting(ctx context.Context, profileID uint64, visibility uint8, profileVisibility uint8, statusOnline enums.EOnlineStatus) error
	UpdateLastSeen(ctx context.Context, profileID uint64) error
}
