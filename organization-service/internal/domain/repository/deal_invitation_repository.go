package repository

import (
	"context"
	"organization/internal/domain/entity"
)

type DealInvitationRepository interface {
	Create(ctx context.Context, invitation *entity.DealInvitation) (*entity.DealInvitation, error)
	Update(ctx context.Context, invitation *entity.DealInvitation) (*entity.DealInvitation, error)
	GetByID(ctx context.Context, id uint64) (*entity.DealInvitation, error)
	GetByDealIDAndInviteeID(ctx context.Context, dealID, inviteeID uint64) (*entity.DealInvitation, error)
	GetByDealID(ctx context.Context, dealID uint64, page, size int) ([]*entity.DealInvitation, uint32, error)
	GetPendingByInviteeID(ctx context.Context, inviteeID uint64, page, size int) ([]*entity.DealInvitation, uint32, error)
	GetAcceptedByDealID(ctx context.Context, dealID uint64) ([]*entity.DealInvitation, error)
	CountByDealID(ctx context.Context, dealID uint64) (uint32, error)
	CountPendingByInviteeID(ctx context.Context, inviteeID uint64) (uint32, error)
	Delete(ctx context.Context, id uint64) error
	ExistsByDealIDAndInviteeID(ctx context.Context, dealID, inviteeID uint64) (bool, error)
} 