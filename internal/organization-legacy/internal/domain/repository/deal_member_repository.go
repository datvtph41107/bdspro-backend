package repository

import (
	"context"
	"organization/internal/domain/entity"
	"organization/internal/dto"
	"organization/internal/enums"
)

type DealMemberRepository interface {
	CreateMember(ctx context.Context, dealMember *entity.DealMember) (*entity.DealMember, error)
	CreateMembers(ctx context.Context, dealMembers []entity.DealMember) ([]entity.DealMember, error)
	GetByDealIDAndMemberID(ctx context.Context, dealID, memberID uint64) (*entity.DealMember, error)
	DeleteMember(ctx context.Context, dealID, memberID uint64) error
	UpdateCommission(ctx context.Context, dealID, memberID uint64, commissionValue float64, commissionType enums.CommissionType, note string) error
	UpdateNote(ctx context.Context, dealID, memberID uint64, note string) error
	UpdateUnilateralStatus(ctx context.Context, dealID, memberID uint64, isUnilateral bool) error
	UpdateMembersCommission(ctx context.Context, dealID uint64, members []*entity.DealMember) error
	GetDealMembers(ctx context.Context, payload dto.SearchMembersRequest) ([]*entity.DealMember, error)
	CalculateTotalCommission(ctx context.Context, dealID uint64) (float64, error)
	GetCommissionStats(ctx context.Context, dealID uint64) (*entity.CommissionStats, error)
	GetByIds(ctx context.Context, ids []uint64) ([]*entity.DealMember, error)
	GetDealMembersByUserIds(ctx context.Context, dealID uint64, userIds []uint64) ([]*entity.DealMember, error)
	IsMemberDeal(ctx context.Context, dealID, memberID uint64) (bool, error)
	UpdateDoneInvestment(ctx context.Context, dealID, memberID uint64) error
	UpdateRole(ctx context.Context, dealID, memberID uint64, roleKey enums.RoleKey) error

	Create(ctx context.Context, invitation *entity.DealMember) (*entity.DealMember, error)
	Update(ctx context.Context, invitation *entity.DealMember) (*entity.DealMember, error)
	GetByID(ctx context.Context, id uint64) (*entity.DealMember, error)
	GetByDealIDAndInviteeID(ctx context.Context, dealID, inviteeID uint64) (*entity.DealMember, error)
	GetByDealID(ctx context.Context, dealID uint64, page, size int) ([]*entity.DealMember, uint32, error)
	GetPendingByInviteeID(ctx context.Context, inviteeID uint64, page, size int) ([]*entity.DealMember, uint32, error)
	GetAcceptedByDealID(ctx context.Context, dealID uint64) ([]*entity.DealMember, error)
	CountByDealID(ctx context.Context, dealID uint64) (uint32, error)
	CountPendingByInviteeID(ctx context.Context, inviteeID uint64) (uint32, error)
	Delete(ctx context.Context, id uint64) error
	ExistsByDealIDAndInviteeID(ctx context.Context, dealID, inviteeID uint64) (bool, error)
	// GetByIds(ctx context.Context, dealId uint64, ids []uint64) ([]*entity.DealMember, error)
}
