package repo

import (
	_enum "common/domain/enum"
	"context"
	"mdn4/modules/user/dto"
	entity "mdn4/modules/user/entity"
)

type IAdminRepo interface {
	GetByID(ctx context.Context, req uint64) (*entity.AdminProfile, error)
	UpdateUserStatus(profileID uint64, status _enum.EUserStatus) error

	// Admin Management methods
	ListAdmin(ctx context.Context, req *dto.AdminListRequest) ([]entity.AdminProfile, int64, error)
	CreateAdmin(ctx context.Context, admin *entity.AdminProfile) (*entity.AdminProfile, error)
	UpdateAdmin(ctx context.Context, req *entity.AdminProfile) (*entity.AdminProfile, error)
	UpdateAuthID(ctx context.Context, adminID, authID uint64) error
	IsSystemRootAdmin(ctx context.Context, id uint64) (bool, error)
	IsSuperAdminByProfileID(ctx context.Context, profileID uint64) (bool, error)
	DeleteAdmin(ctx context.Context, authID uint64) error
	GetDetail(ctx context.Context, req uint64) (*entity.AdminProfile, error)
	GetMapByAuthIDs(ctx context.Context, authIDs []uint64) (map[uint64]*entity.AdminProfile, error)
	GetMapByIDs(ctx context.Context, ids []uint64) (map[uint64]*entity.AdminProfile, error)
}
