package admin_postgres

import (
	"bdspro/internal/domain"
	admin_repo "bdspro/internal/repo/admin"
	_db "common/db"
	"context"
	"time"

	"gorm.io/gorm"
)

type AdminTransactionPostgres struct {
	DB *_db.TransactionRepo
}

func NewAdminTransactionPostgres(db *_db.TransactionRepo) admin_repo.IAdminTransactionRepo {
	return &AdminTransactionPostgres{
		DB: db,
	}
}

func (r *AdminTransactionPostgres) GetList(ctx context.Context, filter *admin_repo.TransactionFilter) ([]*domain.Transaction, int64, error) {
	var transactions []*domain.Transaction
	var total int64

	query := r.DB.GetDB(ctx).Model(&domain.Transaction{}).
		Where("deleted_at IS NULL")

	// Apply filters
	if filter.Keyword != nil && *filter.Keyword != "" {
		keyword := "%" + *filter.Keyword + "%"
		query = query.Where("transaction_name ILIKE ? OR description ILIKE ?", keyword, keyword)
	}

	if filter.TransactionType != nil {
		query = query.Where("transaction_type = ?", *filter.TransactionType)
	}

	if filter.TransactionStatus != nil {
		query = query.Where("transaction_status = ?", *filter.TransactionStatus)
	}

	if filter.ApprovalStatus != nil {
		query = query.Where("approval_status = ?", *filter.ApprovalStatus)
	}

	if filter.ProductId != nil {
		query = query.Where("product_id = ?", *filter.ProductId)
	}

	if filter.ContactId != nil {
		query = query.Where("contact_id = ?", *filter.ContactId)
	}

	if filter.OwnerId != nil {
		query = query.Where("owner_id = ?", *filter.OwnerId)
	}

	if filter.OwnerType != nil {
		query = query.Where("owner_type = ?", *filter.OwnerType)
	}

	if filter.StartDate != nil {
		startDate, err := time.Parse(time.RFC3339, *filter.StartDate)
		if err == nil {
			query = query.Where("transaction_date >= ?", startDate)
		}
	}

	if filter.EndDate != nil {
		endDate, err := time.Parse(time.RFC3339, *filter.EndDate)
		if err == nil {
			query = query.Where("transaction_date <= ?", endDate)
		}
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pagination
	offset := filter.Page * filter.Size
	query = query.Offset(offset).Limit(filter.Size)

	// Order by created_at desc
	query = query.Order("created_at DESC")

	// Preload product
	query = query.Preload("Product")

	if err := query.Find(&transactions).Error; err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

func (r *AdminTransactionPostgres) GetDetail(ctx context.Context, id uint64) (*domain.Transaction, error) {
	var transaction domain.Transaction

	err := r.DB.GetDB(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		Preload("Product").
		First(&transaction).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return &transaction, nil
}

func (r *AdminTransactionPostgres) UpdateApprovalStatus(ctx context.Context, id uint64, status string, approvedBy *uint32) error {
	updates := map[string]interface{}{
		"approval_status": status,
	}

	if approvedBy != nil {
		updates["approved_by"] = *approvedBy
		now := time.Now()
		updates["approval_date"] = now
	}

	return r.DB.GetDB(ctx).
		Model(&domain.Transaction{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *AdminTransactionPostgres) Delete(ctx context.Context, id uint64) error {
	return r.DB.GetDB(ctx).
		Model(&domain.Transaction{}).
		Where("id = ?", id).
		Update("deleted_at", time.Now()).Error
}
