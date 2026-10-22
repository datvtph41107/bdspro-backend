package repository

import (
	"context"
	"strings"

	"mdn4/modules/bdspro/entity"

	"gorm.io/gorm"
)

type PropertyHistoryFilter struct {
	SubjectID uint64
	Action    *uint32
	ActorID   *uint64
	Keyword   string
	Page      uint32
	Size      uint32
}

type PropertyHistoryRepository interface {
	ListBySubject(ctx context.Context, filter PropertyHistoryFilter) ([]*entity.PropertyHistory, int64, error)
}

type propertyHistoryPostgres struct{ db *gorm.DB }

func NewPropertyHistoryPostgres(db *gorm.DB) PropertyHistoryRepository {
	return &propertyHistoryPostgres{db: db}
}

func (r *propertyHistoryPostgres) ListBySubject(ctx context.Context, filter PropertyHistoryFilter) ([]*entity.PropertyHistory, int64, error) {
	if r.db == nil || filter.SubjectID == 0 {
		return nil, 0, nil
	}
	q := r.db.WithContext(ctx).Model(&entity.PropertyHistory{}).
		Where("subject_id = ? AND deleted_at IS NULL", filter.SubjectID)
	if filter.Action != nil {
		q = q.Where("action = ?", *filter.Action)
	}
	if filter.ActorID != nil && *filter.ActorID > 0 {
		q = q.Where("actor_id = ?", *filter.ActorID)
	}
	if kw := strings.TrimSpace(filter.Keyword); kw != "" {
		like := "%" + kw + "%"
		q = q.Where("(description ILIKE ? OR metadata::text ILIKE ?)", like, like)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	limit := int(filter.Size)
	if limit <= 0 {
		limit = 20
	}
	if limit > 999 {
		limit = 999
	}
	offset := int(filter.Page) * limit

	var list []*entity.PropertyHistory
	err := q.Order("occurred_at DESC").Limit(limit).Offset(offset).Find(&list).Error
	return list, total, err
}
