package postgres

import (
	"context"
	"fmt"
	"strings"

	"notification/internal/domain"
	"notification/internal/dto"
	usecase "notification/internal/usecase"

	"gorm.io/gorm"
)

type propertyHistoryPostgres struct {
	db *gorm.DB
}

func NewActivityHistoryPostgres(db *gorm.DB) usecase.PropertyHistoryStore {
	return &propertyHistoryPostgres{db: db}
}

func (r *propertyHistoryPostgres) Create(
	ctx context.Context,
	entity *domain.PropertyHistory,
) (*domain.PropertyHistory, error) {
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		return nil, fmt.Errorf("activity_history.create failed: %w", err)
	}
	return entity, nil
}

func (r *propertyHistoryPostgres) Search(
	ctx context.Context,
	q *dto.PropertyHistoryQueryDTO,
) ([]*domain.PropertyHistory, int64, error) {
	var (
		histories []*domain.PropertyHistory
		total     int64
	)

	db := r.db.WithContext(ctx).Model(&domain.PropertyHistory{})
	if q.SubjectID != nil {
		db = db.Where("subject_id = ?", *q.SubjectID)
	}
	if q.Action != nil {
		db = db.Where("action = ?", *q.Action)
	}

	if q.ActorID != nil {
		db = db.Where("actor_id = ?", *q.ActorID)
	}

	if q.Keyword != nil {
		keyword := strings.TrimSpace(*q.Keyword)
		if keyword != "" {
			kw := "%" + keyword + "%"
			db = db.Where(
				"(description ILIKE ? OR metadata::text ILIKE ?)",
				kw, kw,
			)
		}
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	order := "occurred_at DESC"
	if q.Sort != "" {
		order = q.Sort
	}

	if q.Limit > 0 {
		db = db.Limit(q.Limit)
	}

	if q.Offset >= 0 {
		db = db.Offset(q.Offset)
	}
	if err := db.
		Order(order).
		Find(&histories).Error; err != nil {
		return nil, 0, err
	}

	return histories, total, nil
}

func (r *propertyHistoryPostgres) GetByID(
	ctx context.Context,
	id uint64,
) (*domain.PropertyHistory, error) {

	var entity domain.PropertyHistory

	if err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&entity).Error; err != nil {

		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &entity, nil
}

func (r *propertyHistoryPostgres) Delete(
	ctx context.Context,
	id uint64,
) error {
	if err := r.db.WithContext(ctx).
		Delete(&domain.PropertyHistory{}, id).Error; err != nil {
	}

	return nil
}
