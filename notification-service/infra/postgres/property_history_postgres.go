package postgres

import (
	"context"
	"fmt"

	"notification/internal/domain"
	"notification/internal/repo"

	"gorm.io/gorm"
)

type propertyHistoryPostgres struct {
	db *gorm.DB
}

func NewActivityHistoryPostgres(db *gorm.DB) repo.PropertyHistoryRepo {
	return &propertyHistoryPostgres{db: db}
}

func (r *propertyHistoryPostgres) Create(
	ctx context.Context,
	entity *domain.PropertyHistoryEntity,
) (*domain.PropertyHistoryEntity, error) {
	if err := r.db.WithContext(ctx).Create(entity).Error; err != nil {
		return nil, fmt.Errorf("activity_history.create failed: %w", err)
	}

	return entity, nil
}

func (r *propertyHistoryPostgres) Search(
	ctx context.Context,
	q repo.PropertyHistoryQuery,
) ([]domain.PropertyHistoryEntity, error) {

	var histories []domain.PropertyHistoryEntity

	db := r.db.WithContext(ctx).
		Model(&domain.PropertyHistoryEntity{})

	if q.SubjectID != nil {
		db = db.Where("subject_id = ?", *q.SubjectID)
	}

	if q.SubjectType != nil {
		db = db.Where("subject_type = ?", *q.SubjectType)
	}

	if q.Action != nil {
		db = db.Where("action = ?", *q.Action)
	}

	if q.ActorID != nil {
		db = db.Where("actor_id = ?", *q.ActorID)
	}

	if q.Keyword != nil && *q.Keyword != "" {
		kw := "%" + *q.Keyword + "%"
		db = db.Where(`
			description ILIKE ?
			OR metadata::text ILIKE ?
		`, kw, kw)
	}

	order := "occurred_at DESC"
	if q.Sort != "" {
		order = q.Sort
	}

	db = db.Order(order)

	if q.Limit > 0 {
		db = db.Limit(q.Limit)
	}

	if q.Offset >= 0 {
		db = db.Offset(q.Offset)
	}

	if err := db.Find(&histories).Error; err != nil {
		return nil, fmt.Errorf("activity_history.search failed: %w", err)
	}

	return histories, nil
}
