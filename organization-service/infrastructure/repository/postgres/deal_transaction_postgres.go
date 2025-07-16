package postgres

import (
	"context"
	"fmt"
	"organization/internal/domain/entity"
	"organization/internal/dto"
	"organization/internal/enums"
	"time"

	"gorm.io/gorm"
)

type DealTransactionPostgresRepository struct {
	db *gorm.DB
}

func NewDealTransactionPostgresRepository(db *gorm.DB) *DealTransactionPostgresRepository {
	return &DealTransactionPostgresRepository{db: db}
}

func (r *DealTransactionPostgresRepository) Create(ctx context.Context, transaction *entity.DealTransaction) error {
	return r.db.WithContext(ctx).Create(transaction).Error
}

func (r *DealTransactionPostgresRepository) Update(ctx context.Context, transaction *entity.DealTransaction) error {
	return r.db.WithContext(ctx).Save(transaction).Error
}

func (r *DealTransactionPostgresRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&entity.DealTransaction{}, id).Error
}

func (r *DealTransactionPostgresRepository) GetByID(ctx context.Context, id uint64) (*entity.DealTransaction, error) {
	var transaction entity.DealTransaction
	err := r.db.WithContext(ctx).
		Preload("AttachDocuments").
		Preload("Deal").
		Where("id = ?", id).
		First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *DealTransactionPostgresRepository) GetByDealID(ctx context.Context, dealID uint64, transactionType *enums.TransactionType, page, limit int) ([]*entity.DealTransaction, int64, error) {
	var transactions []*entity.DealTransaction
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.DealTransaction{}).Where("deal_id = ?", dealID)
	
	if transactionType != nil {
		query = query.Where("type = ?", *transactionType)
	}

	// Đếm tổng số records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Lấy data với pagination
	offset := (page - 1) * limit
	err := query.
		Preload("AttachDocuments").
		Preload("Deal").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&transactions).Error

	if err != nil {
		return nil, 0, err
	}

	return transactions, total, nil
}

func (r *DealTransactionPostgresRepository) Approve(ctx context.Context, id uint64, approvedBy uint64, status enums.ApprovedStatus, rejectReason string) error {
	updates := map[string]interface{}{
		"status":        status,
		"approved_by":   approvedBy,
		"approved_at":   time.Now(),
		"reject_reason": rejectReason,
	}

	return r.db.WithContext(ctx).
		Model(&entity.DealTransaction{}).
		Where("id = ?", id).
		Updates(updates).Error
}

func (r *DealTransactionPostgresRepository) GetFinancialSummary(ctx context.Context, dealID uint64) (*dto.DealFinancialSummary, error) {
	// Lấy thông tin deal để có target profit
	var deal entity.Deal
	if err := r.db.WithContext(ctx).Where("id = ?", dealID).First(&deal).Error; err != nil {
		return nil, fmt.Errorf("deal not found: %w", err)
	}

	// Tính tổng doanh thu (chỉ những transaction đã được approve)
	var totalRevenue float64
	err := r.db.WithContext(ctx).
		Model(&entity.DealTransaction{}).
		Where("deal_id = ? AND type = ? AND status = ?", 
			dealID, enums.TransactionTypeRevenue, enums.ApprovedStatusApproved).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalRevenue).Error
	if err != nil {
		return nil, fmt.Errorf("failed to calculate total revenue: %w", err)
	}

	// Tính tổng chi phí (chỉ những transaction đã được approve)
	var totalExpense float64
	err = r.db.WithContext(ctx).
		Model(&entity.DealTransaction{}).
		Where("deal_id = ? AND type = ? AND status = ?", 
			dealID, enums.TransactionTypeExpense, enums.ApprovedStatusApproved).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&totalExpense).Error
	if err != nil {
		return nil, fmt.Errorf("failed to calculate total expense: %w", err)
	}

	// Tính lợi nhuận
	netProfit := totalRevenue - totalExpense

	// Tính tỷ lệ lợi nhuận
	var profitMargin float64
	if totalRevenue > 0 {
		profitMargin = (netProfit / totalRevenue) * 100
	}

	// Xác định trạng thái lợi nhuận
	profitStatus := "on_track"
	if deal.TargetProfit > 0 {
		if netProfit < deal.TargetProfit*0.8 {
			profitStatus = "below_target"
		} else if netProfit > deal.TargetProfit*1.2 {
			profitStatus = "above_target"
		}
	}

	return &dto.DealFinancialSummary{
		DealID:       dealID,
		TotalRevenue: totalRevenue,
		TotalExpense: totalExpense,
		NetProfit:    netProfit,
		ProfitMargin: profitMargin,
		TargetProfit: deal.TargetProfit,
		ProfitStatus: profitStatus,
	}, nil
}

func (r *DealTransactionPostgresRepository) GetTotalAmountByType(ctx context.Context, dealID uint64, transactionType enums.TransactionType) (float64, error) {
	var total float64
	err := r.db.WithContext(ctx).
		Model(&entity.DealTransaction{}).
		Where("deal_id = ? AND type = ? AND status = ?", 
			dealID, transactionType, enums.ApprovedStatusApproved).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error
	if err != nil {
		return 0, fmt.Errorf("failed to calculate total amount: %w", err)
	}
	return total, nil
} 