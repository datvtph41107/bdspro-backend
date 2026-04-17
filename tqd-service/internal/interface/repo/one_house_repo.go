package repo

import (
	"context"
	"tqd/internal/domain"
)

type OneHouseRepo interface {
	GetByID(ctx context.Context, id string) (*domain.OneHouse, error)
	GetByPropertyUUID(ctx context.Context, uuid string) (*domain.OneHouse, error)
}
