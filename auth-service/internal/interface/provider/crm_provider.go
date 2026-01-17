package provider

import (
	"context"
)

type CrmProvider interface {
	GetOriginID(ctx context.Context, profileID uint64) (uint64, error)
}
