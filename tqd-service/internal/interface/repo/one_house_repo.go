package repo

import (
	"context"
	"tqd/internal/domain"
)

type OneHouseRepo interface {
	GetByID(ctx context.Context, id uint64) (*domain.OneHouse, error)
	GetByPropertyUUID(ctx context.Context, uuid string) (*domain.OneHouse, error)
}
