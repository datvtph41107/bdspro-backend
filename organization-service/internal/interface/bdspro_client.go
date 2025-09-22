package iusecase

import (
	"context"
	bdspropb "pb/types/bdspro"
)

// @bind: internal/interface
type BdsproClient interface {
	GetCountByOwner(ctx context.Context, req *bdspropb.GetCountByOwnerRequest) (*bdspropb.GetCountByOwnerResponse, error)
}
