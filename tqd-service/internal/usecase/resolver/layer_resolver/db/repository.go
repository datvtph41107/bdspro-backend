package db

import (
	"context"
	"encoding/json"
)

type IConfigRepository interface {
	GetByKey(ctx context.Context, key string) (json.RawMessage, error)
	GetAll(ctx context.Context) (map[string]json.RawMessage, error)
	UpdateByKey(ctx context.Context, key string, value json.RawMessage, updatedBy string) error
}
