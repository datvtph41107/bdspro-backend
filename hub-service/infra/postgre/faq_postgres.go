package postgres

import (
	"context"
	"strings"

	_dto "common/domain/dto"
	_provider "common/provider"
	"hub/internal/domain"
	"hub/internal/repo"

	_db "common/db"
)

// FAQPostgres implementation của FAQ repository
type FAQPostgres struct {
	_provider.CrudRepo[domain.FAQEntity]
}

// NewFAQRepo tạo mới FAQRepository
func NewFAQRepo(db *_db.TransactionRepo) repo.IFAQRepo {
	repo := &FAQPostgres{}
	repo.Init(repo, db)
	return repo
}

// GetListWithFilter retrieves FAQs with filters and pagination
func (r *FAQPostgres) GetListWithFilter(ctx context.Context, question string, groupKey string, pagable _dto.IPagable) ([]*domain.FAQEntity, int64, error) {
	var faqs []*domain.FAQEntity
	var total int64

	db := r.GetDB(ctx).Model(&domain.FAQEntity{})

	// Apply filters
	question = strings.TrimSpace(question)
	groupKey = strings.ToLower(strings.TrimSpace(groupKey))

	if question != "" {
		db = db.Where("question ILIKE ?", "%"+question+"%")
	}
	if groupKey != "" {
		db = db.Where("group_key = ?", groupKey)
	}

	// Count total
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := pagable.GetOffset()
	limit := pagable.GetLimit()
	db = db.Offset(offset).Limit(limit)

	// Order by created_at desc
	db = db.Order("updated_at DESC, id DESC")

	// Execute query
	if err := db.Find(&faqs).Error; err != nil {
		return nil, 0, err
	}

	return faqs, total, nil
}

// GetSimpleListWithText retrieves FAQs with text search (question or answer) and pagination
func (r *FAQPostgres) GetSimpleListWithText(ctx context.Context, text string, groupKey string, pagable _dto.IPagable) ([]*domain.FAQEntity, int64, error) {
	var faqs []*domain.FAQEntity
	var total int64

	db := r.GetDB(ctx).Model(&domain.FAQEntity{})

	text = strings.TrimSpace(text)
	groupKey = strings.ToLower(strings.TrimSpace(groupKey))

	// Apply text search: lowercase và thay khoảng trống bằng %
	if text != "" {
		// Lowercase
		searchText := strings.ToLower(text)
		// Thay khoảng trống bằng %
		searchText = strings.ReplaceAll(searchText, " ", "%")
		// Tạo pattern LIKE
		likePattern := "%" + searchText + "%"

		// Search trong cả question và answer (lowercase)
		db = db.Where("LOWER(question) LIKE ? OR LOWER(answer) LIKE ?", likePattern, likePattern)
	}

	// Apply groupKey filter
	if groupKey != "" {
		db = db.Where("group_key = ?", groupKey)
	}

	// Count total
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := pagable.GetOffset()
	limit := pagable.GetLimit()
	db = db.Offset(offset).Limit(limit)

	// Order by created_at desc
	db = db.Order("updated_at DESC, id DESC")

	// Execute query
	if err := db.Find(&faqs).Error; err != nil {
		return nil, 0, err
	}

	return faqs, total, nil
}
