package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/interface/repo"

	"gorm.io/gorm"
)

type qhAuthorityIssuringRepo struct {
	db *gorm.DB
}

// NewQHAuthorityIssuringRepository @bind: internal/interface/repo.QHAuthorityIssuringRepository
func NewQHAuthorityIssuringRepository(db *gorm.DB) repo.QHAuthorityIssuringRepository {
	return &qhAuthorityIssuringRepo{db: db}
}

func (r *qhAuthorityIssuringRepo) Create(ctx context.Context, row *qh_domain.QHAuthorityIssuring) error {
	return r.db.WithContext(ctx).Create(row).Error
}

func (r *qhAuthorityIssuringRepo) Update(ctx context.Context, row *qh_domain.QHAuthorityIssuring) error {
	return r.db.WithContext(ctx).Model(row).
		Omit("created_at", "created_by", "id").
		Select("Name", "Code", "Description").
		Updates(row).Error
}

func (r *qhAuthorityIssuringRepo) Delete(ctx context.Context, id uint64) error {
	if id == 0 {
		return nil
	}
	now := time.Now()
	return r.db.WithContext(ctx).Model(&qh_domain.QHAuthorityIssuring{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", now).Error
}

func (r *qhAuthorityIssuringRepo) GetByID(ctx context.Context, id uint64) (*qh_domain.QHAuthorityIssuring, error) {
	var row qh_domain.QHAuthorityIssuring
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("get qh_authority_issuring by id: %w", err)
	}
	return &row, nil
}

func (r *qhAuthorityIssuringRepo) GetByCode(ctx context.Context, code string) (*qh_domain.QHAuthorityIssuring, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return nil, nil
	}
	var row qh_domain.QHAuthorityIssuring
	err := r.db.WithContext(ctx).
		Where("code = ? AND deleted_at IS NULL", code).
		First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("get qh_authority_issuring by code: %w", err)
	}
	return &row, nil
}

func (r *qhAuthorityIssuringRepo) List(ctx context.Context, offset, limit int) ([]qh_domain.QHAuthorityIssuring, int64, error) {
	q := r.db.WithContext(ctx).Model(&qh_domain.QHAuthorityIssuring{}).Where("deleted_at IS NULL")

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count qh_authority_issuring: %w", err)
	}

	var rows []qh_domain.QHAuthorityIssuring
	err := q.Order("id ASC").Offset(offset).Limit(limit).Find(&rows).Error
	if err != nil {
		return nil, 0, fmt.Errorf("list qh_authority_issuring: %w", err)
	}
	return rows, total, nil
}
