package postgres

import (
	"context"

	"time"

	"payment/internal/domain/entity"
	"payment/internal/domain/repository"

	"gorm.io/gorm"
)

type WithdrawalRequestModel struct {
	ID              uint32 `gorm:"primaryKey"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
	WalletId        uint32
	Amount          float64
	Status          string
	PaymentMethod   string
	BankAccountInfo string
	RequestedBy     uint32
	ApprovedBy      uint32
}

func (m *WithdrawalRequestModel) TableName() string {
	return "withdrawal_requests"
}

func WithdrawalRequestModelToEntity(m *WithdrawalRequestModel) *entity.WithdrawalRequest {
	return &entity.WithdrawalRequest{
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

func WithdrawalRequestEntityToModel(e *entity.WithdrawalRequest) *WithdrawalRequestModel {
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

func NewWithdrawalRequestsRepository(db *gorm.DB) repository.WithdrawalRequestsRepository {
	return &withdrawalRequestsPostgresRepository{
		db: db,
	}
}

func (r *withdrawalRequestsPostgresRepository) CreateWithdrawalRequest(ctx context.Context, req *entity.WithdrawalRequest) error {
	model := WithdrawalRequestEntityToModel(req)
	return r.db.WithContext(ctx).Create(model).Error
}

func (r *withdrawalRequestsPostgresRepository) GetWithdrawalRequestById(ctx context.Context, id uint32) (*entity.WithdrawalRequest, error) {
	var model WithdrawalRequestModel
	if err := r.db.WithContext(ctx).Where("id = ? AND is_deleted = ?", id, false).First(&model).Error; err != nil {
		return nil, err
	}
	return WithdrawalRequestModelToEntity(&model), nil
}

func (r *withdrawalRequestsPostgresRepository) UpdateWithdrawalRequest(ctx context.Context, request *entity.WithdrawalRequest) error {
	model := WithdrawalRequestEntityToModel(request)
	return r.db.WithContext(ctx).Save(model).Error
}
