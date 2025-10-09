package repo

import (
	_dto "common/domain/dto"
	"context"
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
}
