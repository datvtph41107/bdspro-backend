package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"tqd/internal/domain"
	"tqd/internal/interface/repo"

	_dto "common/domain/dto"

	"gorm.io/gorm"
)

type contactLabelPostgres struct {
	db *gorm.DB
}

func NewContactLabelPostgresRepo(db *gorm.DB) repo.ContactLabelRepository {
	return &contactLabelPostgres{db: db}
}

func (r *contactLabelPostgres) Create(ctx context.Context, contactLabel *domain.ContactLabel) error {
	return r.db.WithContext(ctx).Create(contactLabel).Error
}

func (r *contactLabelPostgres) GetByID(ctx context.Context, id uint64) (*domain.ContactLabel, error) {
	var contactLabel domain.ContactLabel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&contactLabel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("contact label not found")
		}
		return nil, err
	}
	return &contactLabel, nil
}

func (r *contactLabelPostgres) GetByCode(ctx context.Context, code string) (*domain.ContactLabel, error) {
	var contactLabel domain.ContactLabel
	err := r.db.WithContext(ctx).Where("code = ?", code).First(&contactLabel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil, nil when not found to handle in usecase
		}
		return nil, err
	}
	return &contactLabel, nil
}

func (r *contactLabelPostgres) Update(ctx context.Context, contactLabel *domain.ContactLabel) error {
	return r.db.WithContext(ctx).Save(contactLabel).Error
}

func (r *contactLabelPostgres) Delete(ctx context.Context, id uint64) error {
	// Check if it's a system label first
	var contactLabel domain.ContactLabel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&contactLabel).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("contact label not found")
		}
		return err
	}

	if contactLabel.IsSystem {
		return fmt.Errorf("cannot delete system contact label")
	}

	return r.db.WithContext(ctx).Delete(&domain.ContactLabel{}, id).Error
}

func (r *contactLabelPostgres) List(ctx context.Context, pagable _dto.Pagable, search string, isActive *bool, isSystem *bool) ([]domain.ContactLabel, int64, error) {
	var contactLabels []domain.ContactLabel
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.ContactLabel{})

	// Apply filters
	if search != "" {
		searchTerm := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(code) LIKE ? OR LOWER(description) LIKE ?",
			searchTerm, searchTerm, searchTerm)
	}

	if isActive != nil {
		query = query.Where("is_active = ?", *isActive)
	}

	if isSystem != nil {
		query = query.Where("is_system = ?", *isSystem)
	}

	// Count total
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Apply pagination and sorting
	err = query.
		Order("sort_order ASC, created_at DESC").
		Offset(pagable.GetOffset()).
		Limit(pagable.GetLimit()).
		Find(&contactLabels).Error
	if err != nil {
		return nil, 0, err
	}

	return contactLabels, total, nil
}

func (r *contactLabelPostgres) UpdateContactCount(ctx context.Context, id uint64, count int) error {
	return r.db.WithContext(ctx).
		Model(&domain.ContactLabel{}).
		Where("id = ?", id).
		Update("contact_count", count).Error
}
