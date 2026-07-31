package repo

import (
	_enum "common/domain/enum"
	"context"
	"user/internal/dto"
	models "user/internal/models"
)

type IAdminRepo interface {
	GetByID(ctx context.Context, req uint64) (*models.AdminProfile, error)
	UpdateUserStatus(profileID uint64, status _enum.EUserStatus) error

	// Admin Management methods
	ListAdmin(ctx context.Context, req *dto.AdminListRequest) ([]models.AdminProfile, int64, error)
	CreateAdmin(ctx context.Context, admin *models.AdminProfile) (*models.AdminProfile, error)
	UpdateAdmin(ctx context.Context, req *models.AdminProfile) (*models.AdminProfile, error)
	UpdateAuthID(ctx context.Context, adminID, authID uint64) error
	DeleteAdmin(ctx context.Context, authID uint64) error
	AssignRole(ctx context.Context, authID uint64, roleKey string) error
	GetDetail(ctx context.Context, req uint64) (*models.AdminProfile, error)
	GetMapByAuthIDs(ctx context.Context, authIDs []uint64) (map[uint64]*models.AdminProfile, error)
	GetMapByIDs(ctx context.Context, ids []uint64) (map[uint64]*models.AdminProfile, error)
}
