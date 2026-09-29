package walletpostgres

import (
	"context"
	"time"

	paymentpb "pb/types/payment"

	"payment/internal/domain/wallet"
	walletuc "payment/internal/usecase/wallet"

	"gorm.io/gorm"
)

type WalletTransactionModel struct {
	ID                uint32 `gorm:"primaryKey"`
	WalletID          uint32
	Type              wallet.TransactionType
	Amount            int64
	Status            wallet.TransactionStatus
	RelatedService    string
	RelatedID         string
	ExternalPaymentID string
	PaymentMethod     string
	Description       string
	CreatedBy         uint32
	ApprovedBy        uint32
	TransactionCode   string `gorm:"unique"`
	CreatedAt         time.Time
	UpdatedAt         time.Time

	DeletedAt *time.Time `gorm:"index"`
	IsDeleted bool       `gorm:"default:false"`
}

func (m *WalletTransactionModel) TableName() string {
	return "wallet_transactions"
}

func WalletTransactionModelToEntity(m *WalletTransactionModel) *wallet.WalletTransaction {
	return &wallet.WalletTransaction{
		Id:                m.ID,
		WalletId:          m.WalletID,
		Type:              m.Type,
		Amount:            m.Amount,
		Status:            m.Status,
		RelatedService:    m.RelatedService,
		RelatedId:         m.RelatedID,
		ExternalPaymentId: m.ExternalPaymentID,
		PaymentMethod:     m.PaymentMethod,
		Description:       m.Description,
		CreatedBy:         m.CreatedBy,
		ApprovedBy:        m.ApprovedBy,
		TransactionCode:   m.TransactionCode,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
	}
}

func WalletTransactionEntityToModel(e *wallet.WalletTransaction) *WalletTransactionModel {
	return &WalletTransactionModel{
		ID:                e.Id,
		WalletID:          e.WalletId,
		Type:              e.Type,
		Amount:            e.Amount,
		Status:            e.Status,
		RelatedService:    e.RelatedService,
		RelatedID:         e.RelatedId,
		ExternalPaymentID: e.ExternalPaymentId,
		PaymentMethod:     e.PaymentMethod,
		Description:       e.Description,
		CreatedBy:         e.CreatedBy,
		ApprovedBy:        e.ApprovedBy,
		TransactionCode:   e.TransactionCode,
		CreatedAt:         e.CreatedAt,
		UpdatedAt:         e.UpdatedAt,
	}
}

type walletTransactionRepository struct {
	db *gorm.DB
}

func NewWalletTransactionRepository(db *gorm.DB) walletuc.WalletTransactionRepository {
	return &walletTransactionRepository{db: db}
}

func (r *walletTransactionRepository) GetTransactions(ctx context.Context, walletId uint32, req *paymentpb.GetWalletTransactionsRequest) ([]*paymentpb.WalletTransaction, int64, error) {
	var models []*WalletTransactionModel
	var total int64
	page := 0
	size := 10
	query := dbFromContext(ctx, r.db).Model(&WalletTransactionModel{}).Where("wallet_id = ?", walletId)

	if req != nil {
		if req.Type != "" {
			query = query.Where("type = ?", req.Type)
		}
		if req.Status != "" {
			query = query.Where("status = ?", req.Status)
		}
		if req.FromDate != "" {
			query = query.Where("created_at >= ?", req.FromDate)
		}
		if req.ToDate != "" {
			query = query.Where("created_at <= ?", req.ToDate)
		}
		if req.Page != 0 {
			page = int(req.Page)
		}
		if req.Size != 0 {
			size = int(req.Size)
		}
	}

	if err := query.Offset(page).Limit(size).Find(&models).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	transactions := make([]*paymentpb.WalletTransaction, len(models))
	for i, model := range models {
		entity := WalletTransactionModelToEntity(model)
		transactions[i] = &paymentpb.WalletTransaction{
			Id:       entity.Id,
			WalletId: entity.WalletId,
			Type:     string(entity.Type),
			Amount:   wallet.MinorToWire(entity.Amount),
			Status:   string(entity.Status),
		}
		transactions[i].RelatedService = entity.RelatedService
		transactions[i].RelatedId = entity.RelatedId
		transactions[i].ExternalPaymentId = entity.ExternalPaymentId
		transactions[i].CreatedAt = entity.CreatedAt.Format(time.RFC3339)
		transactions[i].UpdatedAt = entity.UpdatedAt.Format(time.RFC3339)
	}

	return transactions, total, nil
}

func (r *walletTransactionRepository) GetTotalAmountByType(ctx context.Context, walletId uint32, transactionType string) (int64, error) {
	var total int64
	if err := dbFromContext(ctx, r.db).Model(&WalletTransactionModel{}).
		Where("wallet_id = ? AND type = ?", walletId, transactionType).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (r *walletTransactionRepository) GetTransactionCountByStatus(ctx context.Context, walletId uint32, status string) (int64, error) {
	var count int64
	if err := dbFromContext(ctx, r.db).Model(&WalletTransactionModel{}).
		Where("wallet_id = ? AND status = ?", walletId, status).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
func (r *walletTransactionRepository) CountTransactionsByStatus(ctx context.Context, status string) (int64, error) {
	var count int64
	if err := dbFromContext(ctx, r.db).Model(&WalletTransactionModel{}).
		Where("status = ?", status).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *walletTransactionRepository) CreateTransaction(ctx context.Context, transaction *wallet.WalletTransaction) error {
	model := WalletTransactionEntityToModel(transaction)
	return dbFromContext(ctx, r.db).Create(model).Error
}

func (r *walletTransactionRepository) UpdateTransaction(ctx context.Context, transaction *wallet.WalletTransaction) error {
	model := WalletTransactionEntityToModel(transaction)
	return dbFromContext(ctx, r.db).Save(model).Error
}

func (r *walletTransactionRepository) GetTransactionByDealId(ctx context.Context, dealId uint32) (*wallet.WalletTransaction, error) {
	var model WalletTransactionModel
	if err := dbFromContext(ctx, r.db).Where("related_id = ?", dealId).First(&model).Error; err != nil {
		return nil, err
	}
	return WalletTransactionModelToEntity(&model), nil
}

func (r *walletTransactionRepository) GetTransactionsReport(ctx context.Context, req *paymentpb.ReportRequest) ([]*paymentpb.WalletTransaction, int64, int64, error) {
	var models []*WalletTransactionModel
	var total int64
	var totalAmount int64

	query := dbFromContext(ctx, r.db).Model(&WalletTransactionModel{})
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
			Id:       entity.Id,
			WalletId: entity.WalletId,
			Type:     string(entity.Type),
			Amount:   wallet.MinorToWire(entity.Amount),
			Status:   string(entity.Status),
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

func (r *walletTransactionRepository) GetDealTransactions(ctx context.Context, req *paymentpb.GetDealTransactionsRequest) ([]*paymentpb.DealTransaction, int64, error) {
	var models []*WalletTransactionModel
	var total int64
	page := 0
	size := 10
	query := dbFromContext(ctx, r.db).Model(&WalletTransactionModel{})
	if req != nil {
		if req.Type != "" {
			query = query.Where("type = ?", req.Type)
		}
		if req.Status != "" {
			query = query.Where("status = ?", req.Status)
		}
		if req.FromDate != "" {
			query = query.Where("created_at >= ?", req.FromDate)
		}
		if req.ToDate != "" {
			query = query.Where("created_at <= ?", req.ToDate)
		}
		if req.Page != 0 {
			page = int(req.Page)
		}
		if req.Size != 0 {
			size = int(req.Size)
		}
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(page).Limit(size).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	transactions := make([]*paymentpb.DealTransaction, len(models))
	for i, model := range models {
		transactions[i] = &paymentpb.DealTransaction{
			Id:        model.ID,
			Type:      string(model.Type),
			Amount:    wallet.MinorToWire(model.Amount),
			Status:    string(model.Status),
			CreatedAt: model.CreatedAt.Format(time.RFC3339),
			UpdatedAt: model.UpdatedAt.Format(time.RFC3339),
		}
	}

	return transactions, total, nil
}

func (r *walletTransactionRepository) GetDealTransactionById(ctx context.Context, id uint32) (*paymentpb.DealTransaction, error) {
	var model WalletTransactionModel
	if err := dbFromContext(ctx, r.db).First(&model, id).Error; err != nil {
		return nil, err
	}
	return &paymentpb.DealTransaction{
		Id:        model.ID,
		Type:      string(model.Type),
		Amount:    wallet.MinorToWire(model.Amount),
		Status:    string(model.Status),
		CreatedAt: model.CreatedAt.Format(time.RFC3339),
		UpdatedAt: model.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (r *walletTransactionRepository) UpdateDealTransaction(ctx context.Context, transaction *paymentpb.DealTransaction) error {
	amountMinor, err := wallet.MinorFromWire(transaction.Amount)
	if err != nil {
		return err
	}
	model := &WalletTransactionModel{
		ID:     transaction.Id,
		Type:   wallet.TransactionType(transaction.Type),
		Amount: amountMinor,
		Status: wallet.TransactionStatus(transaction.Status),
	}
	return dbFromContext(ctx, r.db).Save(model).Error
}

type TransactionTypeModel struct {
	ID        uint32 `gorm:"primaryKey"`
	Name      string
	Code      string
	Category  string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (m *TransactionTypeModel) TableName() string {
	return "transaction_types"
}

func (r *walletTransactionRepository) GetTransactionTypes(ctx context.Context) ([]*paymentpb.TransactionType, error) {
	var models []*TransactionTypeModel
	if err := dbFromContext(ctx, r.db).Find(&models).Error; err != nil {
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
	return dbFromContext(ctx, r.db).Create(model).Error
}

func (r *walletTransactionRepository) UpdateTransactionType(ctx context.Context, transactionType *paymentpb.TransactionType) error {
	model := &TransactionTypeModel{
		ID:       transactionType.Id,
		Name:     transactionType.Name,
		Code:     transactionType.Code,
		Category: transactionType.Category,
		IsActive: transactionType.IsActive,
	}
	return dbFromContext(ctx, r.db).Model(&TransactionTypeModel{}).Where("id = ?", transactionType.Id).Updates(model).Error
}

func (r *walletTransactionRepository) DeleteTransactionType(ctx context.Context, id uint32) error {
	return dbFromContext(ctx, r.db).Delete(&TransactionTypeModel{}, id).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_at": time.Now(),
		}).Error
}

func (r *walletTransactionRepository) GetTransactionByCode(ctx context.Context, code string) (*wallet.WalletTransaction, error) {
	var model WalletTransactionModel
	if err := dbFromContext(ctx, r.db).Where("transaction_code = ? AND is_deleted = ?", code, false).First(&model).Error; err != nil {
		return nil, err
	}
	return WalletTransactionModelToEntity(&model), nil
}
