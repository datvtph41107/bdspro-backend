package repository

import (
	"context"
	"organization/internal/domain/entity"
)

type DealMemberRepository interface {
	CreateMember(ctx context.Context, dealMember *entity.DealMember) (*entity.DealMember, error)
}