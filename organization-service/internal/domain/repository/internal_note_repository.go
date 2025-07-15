package repository

import (
	"context"
	"organization/internal/domain/entity"
	"organization/internal/enums"
)

type InternalNoteRepository interface {
	Create(ctx context.Context, note *entity.InternalNote) (*entity.InternalNote, error)
	GetByDealID(ctx context.Context, dealID uint64, actionTypes []enums.ActionType, page, limit uint32) ([]*entity.InternalNote, uint32, error)
	GetByID(ctx context.Context, id uint64) (*entity.InternalNote, error)
} 