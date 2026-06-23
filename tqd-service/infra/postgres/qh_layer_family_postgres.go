package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/enums"
	"tqd/internal/interface/repo"

	"gorm.io/gorm"
)

type qhLayerFamilyRepo struct {
	db *gorm.DB
}

// NewQHLayerFamilyRepository @bind: internal/interface/repo.QHLayerFamilyRepository
func NewQHLayerFamilyRepository(db *gorm.DB) repo.QHLayerFamilyRepository {
	return &qhLayerFamilyRepo{db: db}
}

func (r *qhLayerFamilyRepo) Create(ctx context.Context, row *qh_domain.QHLayerFamily) error {
	return r.db.WithContext(ctx).Omit("Layers").Create(row).Error
}

func (r *qhLayerFamilyRepo) Update(ctx context.Context, row *qh_domain.QHLayerFamily) error {
	return r.db.WithContext(ctx).Model(row).
		Omit("created_at", "created_by", "id", "Layers").
		Select("Name", "SortNumber").
		Updates(row).Error
}

func (r *qhLayerFamilyRepo) Delete(ctx context.Context, id uint64) error {
	if id == 0 {
		return nil
	}
	now := time.Now()
	return r.db.WithContext(ctx).Model(&qh_domain.QHLayerFamily{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", now).Error
}

func (r *qhLayerFamilyRepo) GetByID(ctx context.Context, id uint64) (*qh_domain.QHLayerFamily, error) {
	var row qh_domain.QHLayerFamily
	err := r.db.WithContext(ctx).
		Omit("Layers").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("get qh_layer_family by id: %w", err)
	}
	return &row, nil
}

func (r *qhLayerFamilyRepo) List(ctx context.Context, offset, limit int, search, sort string) ([]qh_domain.QHLayerFamily, int64, error) {
	q := r.db.WithContext(ctx).Model(&qh_domain.QHLayerFamily{}).Where("deleted_at IS NULL")
	search = strings.TrimSpace(search)
	if search != "" {
		like := "%" + search + "%"
		q = q.Where("name ILIKE ?", like)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count qh_layer_families: %w", err)
	}

	if sort == "" {
		sort = "sort_number DESC, id DESC"
	}

	var rows []qh_domain.QHLayerFamily
	err := q.Omit("Layers").Order(sort).Offset(offset).Limit(limit).Find(&rows).Error
	if err != nil {
		return nil, 0, fmt.Errorf("list qh_layer_families: %w", err)
	}
	return rows, total, nil
}

func (r *qhLayerFamilyRepo) ListClient(ctx context.Context, offset, limit int, search string) ([]qh_domain.QHLayerFamily, int64, error) {
	q := r.db.WithContext(ctx).Model(&qh_domain.QHLayerFamily{}).
		Where("qh_layer_families.deleted_at IS NULL").
		Where(`EXISTS (
			SELECT 1 FROM qh_layers l
			WHERE l.family_id = qh_layer_families.id
			  AND l.deleted_at IS NULL
			  AND l.status = ?
			  AND l.display_order > 0
		)`, enums.LayerStatusActive)

	search = strings.TrimSpace(search)
	if search != "" {
		like := "%" + search + "%"
		q = q.Where("qh_layer_families.name ILIKE ?", like)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count client qh_layer_families: %w", err)
	}

	var rows []qh_domain.QHLayerFamily
	err := q.Omit("Layers").
		Order("qh_layer_families.sort_number ASC, qh_layer_families.id ASC").
		Offset(offset).
		Limit(limit).
		Find(&rows).Error
	if err != nil {
		return nil, 0, fmt.Errorf("list client qh_layer_families: %w", err)
	}
	return rows, total, nil
}
