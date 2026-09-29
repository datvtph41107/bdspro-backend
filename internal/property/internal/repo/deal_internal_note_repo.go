package repo

import (
	"bdspro/internal/domain"
	"bdspro/internal/enums"
	"context"
)

type DealInternalNoteRepo interface {
	Create(ctx context.Context, note *domain.DealInternalNote) (*domain.DealInternalNote, error)
	GetByDealID(ctx context.Context, dealID uint64, actionTypes []enums.DealActionType, page, limit uint32) ([]*domain.DealInternalNote, uint32, error)
	GetByID(ctx context.Context, id uint64) (*domain.DealInternalNote, error)
}
