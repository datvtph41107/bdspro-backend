package providers

import "context"

type UpdateDataProvider interface {
	GetRefreshIds(c context.Context, resource string, ownerId uint64) ([]uint64, error)
	GetUpdatedAtOfId(c context.Context, resource string, id uint64) (int64, error)
}
