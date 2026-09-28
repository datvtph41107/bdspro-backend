package services

import (
	"context"

	bdspropb "pb/types/bdspro"
)

// PublicService owns public catalog/enum flows for modules/bdspro.
type PublicService interface {
	GetEnums(ctx context.Context, timestamp int64) (*bdspropb.GetEnumsResponse, error)
	GetEnumByName(ctx context.Context, name string) (*bdspropb.GetEnumResponse, error)
}
