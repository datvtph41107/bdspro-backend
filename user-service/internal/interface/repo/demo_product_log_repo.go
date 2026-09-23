package repo

import (
	"context"
	"user/internal/domain/user"
)

type IDemoProductLogRepo interface {
	CreateBatch(ctx context.Context, logs []*user.DemoProductLog) error
	ListProductIDsByProfileID(ctx context.Context, profileID uint64) ([]uint64, error)
	DeleteByProfileID(ctx context.Context, profileID uint64) error
	DeleteByProfileIDAndProductIDs(ctx context.Context, profileID uint64, productIDs []uint64) error

	CreateContactBatch(ctx context.Context, logs []*user.DemoContactLog) error
	ListContactIDsByProfileID(ctx context.Context, profileID uint64) ([]uint64, error)
	DeleteContactsByProfileID(ctx context.Context, profileID uint64) error
	DeleteContactsByProfileIDAndContactIDs(ctx context.Context, profileID uint64, contactIDs []uint64) error
}
