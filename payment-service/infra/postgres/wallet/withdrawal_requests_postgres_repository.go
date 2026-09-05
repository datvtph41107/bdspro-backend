package walletpostgres

import (
	"context"

	"time"

	"payment/internal/domain/wallet"
	walletuc "payment/internal/usecase/wallet"

	"gorm.io/gorm"
)

type WithdrawalRequestModel struct {
	ID              uint32 `gorm:"primaryKey"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	WalletId        uint32
	Amount          int64
	Status          string
	PaymentMethod   string
	BankAccountInfo string
	RequestedBy     uint32
	ApprovedBy      uint32
}

func (m *WithdrawalRequestModel) TableName() string {
	return "withdrawal_requests"
}

func WithdrawalRequestModelToEntity(m *WithdrawalRequestModel) *wallet.WithdrawalRequest {
	return &wallet.WithdrawalRequest{
		Id:              m.ID,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
		WalletId:        m.WalletId,
		Amount:          m.Amount,
		Status:          m.Status,
		PaymentMethod:   m.PaymentMethod,
		BankAccountInfo: m.BankAccountInfo,
		RequestedBy:     m.RequestedBy,
		ApprovedBy:      m.ApprovedBy,
	}
}

func WithdrawalRequestEntityToModel(e *wallet.WithdrawalRequest) *WithdrawalRequestModel {
	return &WithdrawalRequestModel{
		ID:              e.Id,
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
		WalletId:        e.WalletId,
		Amount:          e.Amount,
		Status:          e.Status,
		PaymentMethod:   e.PaymentMethod,
		BankAccountInfo: e.BankAccountInfo,
		RequestedBy:     e.RequestedBy,
		ApprovedBy:      e.ApprovedBy,
	}
}

type withdrawalRequestsPostgresRepository struct {
	db *gorm.DB
}

func NewWithdrawalRequestsRepository(db *gorm.DB) walletuc.WithdrawalRepository {
	return &withdrawalRequestsPostgresRepository{
		db: db,
	}
}

func (r *withdrawalRequestsPostgresRepository) CreateWithdrawalRequest(ctx context.Context, req *wallet.WithdrawalRequest) error {
	model := WithdrawalRequestEntityToModel(req)
	return dbFromContext(ctx, r.db).Create(model).Error
}

func (r *withdrawalRequestsPostgresRepository) GetWithdrawalRequestById(ctx context.Context, id uint32) (*wallet.WithdrawalRequest, error) {
	var model WithdrawalRequestModel
	if err := dbFromContext(ctx, r.db).Where("id = ? AND is_deleted = ?", id, false).First(&model).Error; err != nil {
		return nil, err
	}
	return WithdrawalRequestModelToEntity(&model), nil
}

func (r *withdrawalRequestsPostgresRepository) UpdateWithdrawalRequest(ctx context.Context, request *wallet.WithdrawalRequest) error {
	model := WithdrawalRequestEntityToModel(request)
	return dbFromContext(ctx, r.db).Save(model).Error
}
