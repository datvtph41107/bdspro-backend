package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/interface/repo"

	"gorm.io/gorm"
)

type landUseDictRepo struct {
	db *gorm.DB
}

func NewLandUseDictRepository(db *gorm.DB) repo.LandUseDictRepository {
	return &landUseDictRepo{db: db}
}

func (r *landUseDictRepo) GetByCode(ctx context.Context, code string) (*qh_domain.LandUseDict, error) {
	var dict qh_domain.LandUseDict
	err := r.db.WithContext(ctx).
		Where("code = ?", code).
		First(&dict).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("get land use dict by code failed: %w", err)
	}
	return &dict, nil
}

func (r *landUseDictRepo) List(ctx context.Context, filter *repo.LandUseDictFilter) ([]qh_domain.LandUseDict, int64, error) {
	var dicts []qh_domain.LandUseDict
	var total int64

	query := r.db.WithContext(ctx).Model(&qh_domain.LandUseDict{})

	if filter != nil {
		if filter.ActiveOnly != nil && *filter.ActiveOnly {
			query = query.Where("is_active = ?", true)
		}
		if filter.Search != nil && *filter.Search != "" {
			searchTerm := "%" + strings.ToLower(*filter.Search) + "%"
			query = query.Where(
				"LOWER(name) LIKE ? OR LOWER(code) LIKE ?",
				searchTerm, searchTerm,
			)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count land use dict failed: %w", err)
	}

	if filter != nil && filter.Pagable != nil {
		offset := filter.Pagable.GetOffset()
		limit := filter.Pagable.GetLimit()
		query = query.Offset(offset).Limit(limit)
	}

	err := query.Order("display_priority DESC, code ASC").Find(&dicts).Error
	if err != nil {
		return nil, 0, fmt.Errorf("list land use dict failed: %w", err)
	}

	return dicts, total, nil
}
