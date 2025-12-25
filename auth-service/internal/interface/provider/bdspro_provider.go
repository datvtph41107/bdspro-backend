package provider

import (
	"auth/internal/dto"
	"context"
)

type BdsproProvider interface {
	GetDealMember(ctx context.Context, dealId uint64, userId uint64) (*dto.DealMember, error)
}
