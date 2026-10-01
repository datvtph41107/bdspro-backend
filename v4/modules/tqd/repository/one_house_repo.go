package repository

import (
	"context"
	qh "userv4/modules/tqd/entity/qh"
)

type OneHouseRepo interface {
	GetByID(ctx context.Context, id uint64) (*qh.OneHouse, error)
	GetByPropertyUUID(ctx context.Context, uuid string) (*qh.OneHouse, error)
}
