package iusecase

import "context"

// SyncProvider dùng cho check cache và cập nhật timestamp key (cnv:c, ...).
type SyncProvider interface {
	GetKey(ctx context.Context, prefix string, resourceId uint64) string
	GetKeyWithoutMe(ctx context.Context, prefix string, resourceId uint64) string
	IsCacheChanged(ctx context.Context, key string, timestamp int64) (bool, bool)
	PutTimeRequest(ctx context.Context, key string, timestamp int64)
	GetTimestamp(ctx context.Context, key string) int64
	PutTimestamp(ctx context.Context, key string, timestamp int64) error
}
