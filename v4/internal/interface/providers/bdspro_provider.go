package providers

import (
	"context"
	"userv4/modules/user/dto"
)

type BdsproProvider interface {
	GetDealMember(ctx context.Context, dealId uint64, userId uint64) (*dto.DealMember, error)
}
