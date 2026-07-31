package providers

import (
	"context"
	"user/internal/dto"
)

type BdsproProvider interface {
	GetDealMember(ctx context.Context, dealId uint64, userId uint64) (*dto.DealMember, error)
}
