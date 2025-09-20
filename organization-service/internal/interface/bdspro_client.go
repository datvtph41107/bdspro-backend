package iusecase

import (
	"context"
	pb_bdspro "pb/types/bdspro"
)

// @bind: internal/interface
type BdsproClient interface {
	GetCountByOwner(ctx context.Context, req *pb_bdspro.GetCountByOwnerRequest) (*pb_bdspro.GetCountByOwnerResponse, error)
}
