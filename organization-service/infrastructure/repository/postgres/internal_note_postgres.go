package postgres

import (
	"context"
	"organization/internal/domain/entity"
	"organization/internal/domain/repository"
	"organization/internal/enums"
	"time"

	"gorm.io/gorm"
)

type InternalNoteModel struct {
	ID         uint64                `gorm:"primaryKey;autoIncrement"`
	DealID     uint64                `gorm:"not null;index"`
	ActorID    uint64                `gorm:"not null"`
	Content    string                `gorm:"type:text;not null"`
	ActionType enums.ActionType      `gorm:"not null;index"`
	CreatedAt  time.Time             `gorm:"autoCreateTime"`
	UpdatedAt  time.Time             `gorm:"autoUpdateTime"`
}

func (InternalNoteModel) TableName() string {
	return "internal_notes"
}

type InternalNotePostgresRepository struct {
	db *gorm.DB
}

func NewInternalNotePostgresRepository(db *gorm.DB) repository.InternalNoteRepository {
	return &InternalNotePostgresRepository{db: db}
}

func (r *InternalNotePostgresRepository) Create(ctx context.Context, note *entity.InternalNote) (*entity.InternalNote, error) {
	model := &InternalNoteModel{
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

func (r *InternalNotePostgresRepository) GetByDealID(ctx context.Context, dealID uint64, actionTypes []enums.ActionType, page, size uint32) ([]*entity.InternalNote, uint32, error) {
	var models []InternalNoteModel
	var total int64

	query := r.db.WithContext(ctx).Where("deal_id = ?", dealID)
	
	// Lọc theo action types nếu có
	if len(actionTypes) > 0 {
		query = query.Where("action_type IN ?", actionTypes)
	}
	
	// Đếm tổng số
	if err := query.Model(&InternalNoteModel{}).Count(&total).Error; err != nil {
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
	notes := make([]*entity.InternalNote, len(models))
	for i, model := range models {
		notes[i] = &entity.InternalNote{
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

func (r *InternalNotePostgresRepository) GetByID(ctx context.Context, id uint64) (*entity.InternalNote, error) {
	var model InternalNoteModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		return nil, err
	}
	
	return &entity.InternalNote{
		ID:         model.ID,
		DealID:     model.DealID,
		ActorID:    model.ActorID,
		Content:    model.Content,
		ActionType: model.ActionType,
		CreatedAt:  model.CreatedAt,
		UpdatedAt:  model.UpdatedAt,
	}, nil
} 