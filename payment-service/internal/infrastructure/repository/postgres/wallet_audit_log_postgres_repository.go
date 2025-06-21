package postgres

import (
	"context"
	"time"

	"github.com/tuandq2112/payment-service/internal/domain/entity"
	"github.com/tuandq2112/payment-service/internal/domain/repository"
	"gorm.io/gorm"
)

type WalletAuditLogModel struct {
	ID          uint32 `gorm:"primaryKey"`
	WalletId    uint32
	OldBalance  float64
	NewBalance  float64
	Reason      string
	ChangedBy   uint32
	ChangedAt   time.Time
}

func (m *WalletAuditLogModel) TableName() string {
	return "wallet_audit_logs"
}

func WalletAuditLogModelToEntity(m *WalletAuditLogModel) *entity.WalletAuditLog {
	return &entity.WalletAuditLog{
		Id:          m.ID,
		WalletId:    m.WalletId,
		OldBalance:  m.OldBalance,
		NewBalance:  m.NewBalance,
		Reason:      m.Reason,
		ChangedBy:   m.ChangedBy,
		ChangedAt:   m.ChangedAt,
	}
}

func WalletAuditLogEntityToModel(e *entity.WalletAuditLog) *WalletAuditLogModel {
	return &WalletAuditLogModel{
		ID:          e.Id,
		WalletId:    e.WalletId,
		OldBalance:  e.OldBalance,
		NewBalance:  e.NewBalance,
		Reason:      e.Reason,
		ChangedBy:   e.ChangedBy,
		ChangedAt:   e.ChangedAt,
	}
}

type walletAuditLogRepository struct {
	db *gorm.DB
}

func NewWalletAuditLogRepository(db *gorm.DB) repository.WalletAuditLogRepository {
	return &walletAuditLogRepository{db: db}
}

func (r *walletAuditLogRepository) GetWalletAuditLog(ctx context.Context, req *entity.WalletAuditLog) (*entity.WalletAuditLog, error) {
	var model WalletAuditLogModel
	if err := r.db.WithContext(ctx).Where("wallet_id = ?", req.WalletId).First(&model).Error; err != nil {
		return nil, err
	}
	return WalletAuditLogModelToEntity(&model), nil
}

func (r *walletAuditLogRepository) CreateAuditLog(ctx context.Context, req *entity.WalletAuditLog) error {
	model := WalletAuditLogEntityToModel(req)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	return nil
}
func (r *walletAuditLogRepository) UpdateWalletAuditLog(ctx context.Context, req *entity.WalletAuditLog) (*entity.WalletAuditLog, error) {
	model := WalletAuditLogEntityToModel(req)
	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		return nil, err
	}
	return WalletAuditLogModelToEntity(model), nil
}

func (r *walletAuditLogRepository) DeleteWalletAuditLog(ctx context.Context, req *entity.WalletAuditLog) (*entity.WalletAuditLog, error) {
	if err := r.db.WithContext(ctx).Delete(&WalletAuditLogModel{}, req.Id).Error; err != nil {
		return nil, err
	}
	return req, nil
}

