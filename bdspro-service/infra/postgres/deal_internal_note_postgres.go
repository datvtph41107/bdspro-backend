package postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/enums"
	"bdspro/internal/repo"
	"context"

	"gorm.io/gorm"
)

type InternalNotePostgresRepository struct {
	db *gorm.DB
}

func NewInternalNotePostgresRepository(db *gorm.DB) repo.DealInternalNoteRepo {
	return &InternalNotePostgresRepository{db: db}
}

func (r *InternalNotePostgresRepository) Create(ctx context.Context, note *domain.DealInternalNote) (*domain.DealInternalNote, error) {
	model := &domain.DealInternalNote{
		DealID:     note.DealID,
		ActorID:    note.ActorID,
		Content:    note.Content,
		ActionType: note.ActionType,
	}

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, err
	}

	note.ID = model.ID
	note.CreatedAt = model.CreatedAt
	note.UpdatedAt = model.UpdatedAt

	return note, nil
}

func (r *InternalNotePostgresRepository) GetByDealID(ctx context.Context, dealID uint64, actionTypes []enums.DealActionType, page, size uint32) ([]*domain.DealInternalNote, uint32, error) {
	var models []domain.DealInternalNote
	var total int64

	query := r.db.WithContext(ctx).Where("deal_id = ?", dealID)

	// Lọc theo action types nếu có
	if len(actionTypes) > 0 {
		query = query.Where("action_type IN ?", actionTypes)
	}

	// Đếm tổng số
	if err := query.Model(&domain.DealInternalNote{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Lấy dữ liệu có phân trang
	offset := page * size
	if err := query.Order("created_at DESC").
		Offset(int(offset)).
		Limit(int(size)).
		Find(&models).Error; err != nil {
		return nil, 0, err
	}

	// Convert to entities
	notes := make([]*domain.DealInternalNote, len(models))
	for i, model := range models {
		notes[i] = &domain.DealInternalNote{
			ID:         model.ID,
			DealID:     model.DealID,
			ActorID:    model.ActorID,
			Content:    model.Content,
			ActionType: model.ActionType,
			CreatedAt:  model.CreatedAt,
			UpdatedAt:  model.UpdatedAt,
		}
	}

	return notes, uint32(total), nil
}

func (r *InternalNotePostgresRepository) GetByID(ctx context.Context, id uint64) (*domain.DealInternalNote, error) {
	var model domain.DealInternalNote
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return nil, err
	}

	return &domain.DealInternalNote{
		ID:         model.ID,
		DealID:     model.DealID,
		ActorID:    model.ActorID,
		Content:    model.Content,
		ActionType: model.ActionType,
		CreatedAt:  model.CreatedAt,
		UpdatedAt:  model.UpdatedAt,
	}, nil
}
