package repo

import (
	_crud "common/domain/crud"
	"context"
	"user/internal/domain/user"
)

type IKYCRepo interface {
	_crud.ICrudRepo[user.KYCEntity]
	GetByProfileID(ctx context.Context, profileID uint64) (*user.KYCEntity, error)
	GetMapByProfileIDs(ctx context.Context, profileIDs []uint64) (map[uint64]*user.KYCEntity, error)
	GetListWithFilter(ctx context.Context, search string, status uint32, page, size int, sortBy, sortOrder string) ([]user.KYCEntity, int64, error)
}
