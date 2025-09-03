package repository

import (
	"context"
	"organization/internal/domain/entity"
	"organization/internal/enums"
)

type DealMemberRepository interface {
	CreateMember(ctx context.Context, dealMember *entity.DealMember) (*entity.DealMember, error)
	GetByDealIDAndMemberID(ctx context.Context, dealID, memberID uint64) (*entity.DealMember, error)
	DeleteMember(ctx context.Context, dealID, memberID uint64) error
	UpdateCommission(ctx context.Context, dealID, memberID uint64, commissionValue float64, commissionType enums.CommissionType, note string) error
	UpdateNote(ctx context.Context, dealID, memberID uint64, note string) error
	UpdateUnilateralStatus(ctx context.Context, dealID, memberID uint64, isUnilateral bool) error
	UpdateMembersCommission(ctx context.Context, dealID uint64, members []*entity.DealMember) error
	GetDealMembers(ctx context.Context, dealID uint64) ([]*entity.DealMember, error)
	CalculateTotalCommission(ctx context.Context, dealID uint64) (float64, error)
	GetCommissionStats(ctx context.Context, dealID uint64) (*entity.CommissionStats, error)
	GetByIds(ctx context.Context, ids []uint64) ([]*entity.DealMember, error)
	GetDealMembersByUserIds(ctx context.Context, dealID uint64, userIds []uint64) ([]*entity.DealMember, error)
}
