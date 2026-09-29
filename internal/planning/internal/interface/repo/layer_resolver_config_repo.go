package repo

import (
	"context"
	"encoding/json"
)

// IQHLayerResolverConfigRepo - Repository cho QHLayerResolverConfig (thiết lập đánh giá tính toán)
type IQHLayerResolverConfigRepo interface {
	GetByKey(ctx context.Context, key string) (json.RawMessage, error)
	GetAll(ctx context.Context) (map[string]json.RawMessage, error)
	UpdateByKey(ctx context.Context, key string, value json.RawMessage, updatedBy string) error
}
