package postgres

import (
	"context"
	"fmt"

	"notification/internal/domain"
	"notification/internal/dto"
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
	q dto.PropertyHistoryQueryDTO,
) ([]domain.PropertyHistoryEntity, int64, error) {

	var (
		histories []domain.PropertyHistoryEntity
		total     int64
	)

	baseQuery := r.db.WithContext(ctx).
		Model(&domain.PropertyHistoryEntity{})

	/* ---------- FILTER ---------- */

	if q.SubjectID != nil {
		baseQuery = baseQuery.Where("subject_id = ?", *q.SubjectID)
	}

	if q.SubjectType != nil {
		baseQuery = baseQuery.Where("subject_type = ?", *q.SubjectType)
	}

	if q.Action != nil {
		baseQuery = baseQuery.Where("action = ?", *q.Action)
	}

	if q.ActorID != nil {
		baseQuery = baseQuery.Where("actor_id = ?", *q.ActorID)
	}

	if q.Keyword != nil && *q.Keyword != "" {
		kw := "%" + *q.Keyword + "%"
		baseQuery = baseQuery.Where(`
			description ILIKE ?
			OR metadata::text ILIKE ?
		`, kw, kw)
	}

	if err := baseQuery.
		Session(&gorm.Session{}).
		Count(&total).Error; err != nil {

		return nil, 0, fmt.Errorf(
			"property_history.count failed: %w",
			err,
		)
	}

	order := "occurred_at DESC"
	if q.Sort != "" {
		order = q.Sort
	}

	query := baseQuery.Order(order)

	if q.Limit > 0 {
		query = query.Limit(q.Limit)
	}

	if q.Offset >= 0 {
		query = query.Offset(q.Offset)
	}

	if err := query.Find(&histories).Error; err != nil {
		return nil, 0, fmt.Errorf(
			"property_history.search failed: %w",
			err,
		)
	}

	return histories, total, nil
}

func (r *propertyHistoryPostgres) GetByID(
	ctx context.Context,
	id uint64,
) (*domain.PropertyHistoryEntity, error) {

	var entity domain.PropertyHistoryEntity

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
		Delete(&domain.PropertyHistoryEntity{}, id).Error; err != nil {
		return fmt.Errorf("property_history.delete failed: %w", err)
	}

	return nil
}
