package facade

import (
	"context"
	"core/kernel/contract"

	"github.com/google/uuid"
)

func WithUUID(ctx context.Context, g contract.UUIDGenerator) context.Context {
	return with(ctx, uuidKey, g)
}

func UUID(ctx context.Context) string {
	if v := get(ctx, uuidKey); v != nil {
		return v.(contract.UUIDGenerator).New()
	}
	return uuid.NewString()
}
