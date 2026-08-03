package repo

import (
	_crud "common/domain/crud"
	"context"
	"user/internal/models"
)

type IKYCRepo interface {
	_crud.ICrudRepo[models.KYCEntity]
	GetByProfileID(ctx context.Context, profileID uint64) (*models.KYCEntity, error)
	GetMapByProfileIDs(ctx context.Context, profileIDs []uint64) (map[uint64]*models.KYCEntity, error)
	GetListWithFilter(ctx context.Context, search string, status uint32, page, size int, sortBy, sortOrder string) ([]models.KYCEntity, int64, error)
}
