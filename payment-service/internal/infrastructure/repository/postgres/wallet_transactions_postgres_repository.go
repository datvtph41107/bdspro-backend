package postgres

import (
	"context"
	"fmt"
	"time"

	paymentpb "pb/types/payment"

	"github.com/tuandq2112/payment-service/internal/domain/entity"
	"github.com/tuandq2112/payment-service/internal/domain/repository"
	"gorm.io/gorm"
)

type WalletTransactionModel struct {
	ID                uint32    `gorm:"primaryKey"`
	WalletID          uint32
	Type              entity.TransactionType
	Amount            float64
	Status            entity.TransactionStatus
	RelatedService    string
	RelatedID         string
	ExternalPaymentID string
	TransactionCode   string `gorm:"unique"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

func (m *WalletTransactionModel) TableName() string {
	return "wallet_transactions"
}

func WalletTransactionModelToEntity(m *WalletTransactionModel) *entity.WalletTransaction {
	return &entity.WalletTransaction{
		Id:               m.ID,
		WalletId:         m.WalletID,
		Type:             m.Type,
		Amount:           m.Amount,
		Status:           m.Status,
		RelatedService:   m.RelatedService,
		RelatedId:        m.RelatedID,
		ExternalPaymentId: m.ExternalPaymentID,
		TransactionCode:  m.TransactionCode,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
	}
}

func WalletTransactionEntityToModel(e *entity.WalletTransaction) *WalletTransactionModel {
	return &WalletTransactionModel{
		ID:                e.Id,
		WalletID:          e.WalletId,
		Type:              e.Type,
		Amount:            e.Amount,
		Status:            e.Status,
		RelatedService:    e.RelatedService,
		RelatedID:         e.RelatedId,
		ExternalPaymentID: e.ExternalPaymentId,
		TransactionCode:   e.TransactionCode,
		CreatedAt:         e.CreatedAt,
		UpdatedAt:         e.UpdatedAt,
	}
}

type walletTransactionRepository struct {
	db *gorm.DB
}

func NewWalletTransactionRepository(db *gorm.DB) repository.WalletTransactionsRepository {
	return &walletTransactionRepository{db: db}
}

func (r *walletTransactionRepository) GetTransactions(ctx context.Context, walletId uint32, filter *paymentpb.TransactionFilter, pagination *paymentpb.PaginationRequest) ([]*paymentpb.WalletTransaction, int64, error) {
	var models []*WalletTransactionModel
	var total int64

	query := r.db.WithContext(ctx).Model(&WalletTransactionModel{}).Where("wallet_id = ?", walletId)

	if filter != nil {
		if filter.Type != "" {
			query = query.Where("type = ?", filter.Type)
		}
		if filter.Status != "" {
			query = query.Where("status = ?", filter.Status)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(int(pagination.Page * pagination.Size)).Limit(int(pagination.Size)).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	transactions := make([]*paymentpb.WalletTransaction, len(models))
	for i, model := range models {
		entity := WalletTransactionModelToEntity(model)
		transactions[i] = &paymentpb.WalletTransaction{
			Id:               entity.Id,
			WalletId:         entity.WalletId,
			Type:             string(entity.Type),
			Amount:           entity.Amount,
			Status:           string(entity.Status),
		}
		transactions[i].RelatedService = entity.RelatedService
		transactions[i].RelatedId = entity.RelatedId
		transactions[i].ExternalPaymentId = entity.ExternalPaymentId
		transactions[i].CreatedAt = entity.CreatedAt.Format(time.RFC3339)
		transactions[i].UpdatedAt = entity.UpdatedAt.Format(time.RFC3339)
	}

	return transactions, total, nil
}

func (r *walletTransactionRepository) GetTotalAmountByType(ctx context.Context, walletId uint32, transactionType string) (float64, error) {
	var total float64
	if err := r.db.WithContext(ctx).Model(&WalletTransactionModel{}).
		Where("wallet_id = ? AND type = ?", walletId, transactionType).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (r *walletTransactionRepository) GetTransactionCountByStatus(ctx context.Context, walletId uint32, status string) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&WalletTransactionModel{}).
		Where("wallet_id = ? AND status = ?", walletId, status).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *walletTransactionRepository) CreateTransaction(ctx context.Context, transaction *entity.WalletTransaction) error {
	model := WalletTransactionEntityToModel(transaction)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *walletTransactionRepository) UpdateTransaction(ctx context.Context, transaction *entity.WalletTransaction) error {
	model := WalletTransactionEntityToModel(transaction)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *walletTransactionRepository) GetTransactionByDealId(ctx context.Context, dealId uint32) (*entity.WalletTransaction, error) {
	var model WalletTransactionModel
	if err := r.db.WithContext(ctx).Where("related_id = ?", dealId).First(&model).Error; err != nil {
		return nil, err
	}
	return WalletTransactionModelToEntity(&model), nil
}

func (r *walletTransactionRepository) GetTransactionsReport(ctx context.Context, req *paymentpb.ReportRequest) ([]*paymentpb.WalletTransaction, int64, float64, error) {
	var models []*WalletTransactionModel
	var total int64
	var totalAmount float64

	query := r.db.WithContext(ctx).Model(&WalletTransactionModel{})
	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, 0, err
	}

	if err := query.Select("COALESCE(SUM(amount), 0)").Scan(&totalAmount).Error; err != nil {
		return nil, 0, 0, err
	}

	if err := query.Find(&models).Error; err != nil {
		return nil, 0, 0, err
	}

	transactions := make([]*paymentpb.WalletTransaction, len(models))
	for i, model := range models {
		entity := WalletTransactionModelToEntity(model)
		transactions[i] = &paymentpb.WalletTransaction{
			Id:               entity.Id,
			WalletId:         entity.WalletId,
			Type:             string(entity.Type),
			Amount:           entity.Amount,
			Status:           string(entity.Status),
		}
		transactions[i].RelatedService = entity.RelatedService
		transactions[i].RelatedId = entity.RelatedId
		transactions[i].ExternalPaymentId = entity.ExternalPaymentId
		transactions[i].CreatedAt = entity.CreatedAt.Format(time.RFC3339)
		transactions[i].UpdatedAt = entity.UpdatedAt.Format(time.RFC3339)
	}

	return transactions, total, totalAmount, nil
}

func (r *walletTransactionRepository) ExportTransactions(ctx context.Context, req *paymentpb.ReportRequest) ([]byte, error) {
	// TODO: Implement export functionality
	return nil, nil
}

func (r *walletTransactionRepository) GetDealTransactions(ctx context.Context, filter *paymentpb.TransactionFilter, pagination *paymentpb.PaginationRequest) ([]*paymentpb.DealTransaction, int64, error) {
	var models []*WalletTransactionModel
	var total int64

	query := r.db.WithContext(ctx).Model(&WalletTransactionModel{})
	if filter != nil {
		if filter.Type != "" {
			query = query.Where("type = ?", filter.Type)
		}
		if filter.Status != "" {
			query = query.Where("status = ?", filter.Status)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(int(pagination.Page * pagination.Size)).Limit(int(pagination.Size)).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	transactions := make([]*paymentpb.DealTransaction, len(models))
	for i, model := range models {
		transactions[i] = &paymentpb.DealTransaction{
			Id:          model.ID,
			Type:        string(model.Type),
			Amount:      model.Amount,
			Status:      string(model.Status),
			CreatedAt:   model.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   model.UpdatedAt.Format(time.RFC3339),
		}
	}

	return transactions, total, nil
}

func (r *walletTransactionRepository) GetDealTransactionById(ctx context.Context, id uint32) (*paymentpb.DealTransaction, error) {
	var model WalletTransactionModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		return nil, err
	}
	return &paymentpb.DealTransaction{
		Id:          model.ID,
		Type:        string(model.Type),
		Amount:      model.Amount,
		Status:      string(model.Status),
		CreatedAt:   model.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   model.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (r *walletTransactionRepository) UpdateDealTransaction(ctx context.Context, transaction *paymentpb.DealTransaction) error {
	model := &WalletTransactionModel{
		ID:     transaction.Id,
		Type:   entity.TransactionType(transaction.Type)	,
		Amount: transaction.Amount,
		Status: entity.TransactionStatus(transaction.Status),
	}
	return r.db.WithContext(ctx).Save(model).Error
}

type TransactionTypeModel struct {
	ID          uint32 `gorm:"primaryKey"`
	Name        string
	Code        string
	Category    string
	IsActive    bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (m *TransactionTypeModel) TableName() string {
	return "transaction_types"
}

func (r *walletTransactionRepository) GetTransactionTypes(ctx context.Context) ([]*paymentpb.TransactionType, error) {
	var models []*TransactionTypeModel
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}

	types := make([]*paymentpb.TransactionType, len(models))
	for i, model := range models {
		types[i] = &paymentpb.TransactionType{
			Id:       model.ID,
			Name:     model.Name,
			Code:     model.Code,
			Category: model.Category,
			IsActive: model.IsActive,
		}
	}
	return types, nil
}

func (r *walletTransactionRepository) CreateTransactionType(ctx context.Context, transactionType *paymentpb.TransactionType) error {
	model := &TransactionTypeModel{
		ID:       transactionType.Id,
		Name:     transactionType.Name,
		Code:     transactionType.Code,
		Category: transactionType.Category,
		IsActive: transactionType.IsActive,
	}
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *walletTransactionRepository) UpdateTransactionType(ctx context.Context, transactionType *paymentpb.TransactionType) error {
	fmt.Printf("Transaction Type: %+v\n", transactionType)
	model := &TransactionTypeModel{
		ID:       transactionType.Id,
		Name:     transactionType.Name,
		Code:     transactionType.Code,
		Category: transactionType.Category,
		IsActive: transactionType.IsActive,
	}
	return r.db.WithContext(ctx).Model(&TransactionTypeModel{}).Where("id = ?", transactionType.Id).Updates(model).Error
}

func (r *walletTransactionRepository) DeleteTransactionType(ctx context.Context, id uint32) error {
	return r.db.WithContext(ctx).Delete(&TransactionTypeModel{}, id).Error
}

func (r *walletTransactionRepository) GetTransactionByCode(ctx context.Context, code string) (*entity.WalletTransaction, error) {
	var model WalletTransactionModel
	if err := r.db.WithContext(ctx).Where("transaction_code = ?", code).First(&model).Error; err != nil {
		return nil, err
	}
	return WalletTransactionModelToEntity(&model), nil
}

