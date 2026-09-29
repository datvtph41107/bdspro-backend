package repo

import (
	"context"
	qh_domain "tqd/internal/domain/qh"
)

type OneHouseRepo interface {
	GetByID(ctx context.Context, id uint64) (*qh_domain.OneHouse, error)
	GetByPropertyUUID(ctx context.Context, uuid string) (*qh_domain.OneHouse, error)
}
