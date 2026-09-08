package walletpostgres

import (
	"context"
	"time"

	"payment/internal/domain/wallet"
	walletuc "payment/internal/usecase/wallet"

	"gorm.io/gorm"
)

type WalletAuditLogModel struct {
	ID         uint32 `gorm:"primaryKey"`
	WalletId   uint32
	OldBalance int64
	NewBalance int64
	Reason     string
	ChangedBy  uint32
	ChangedAt  time.Time

	DeletedAt *time.Time `gorm:"index"`
	IsDeleted bool       `gorm:"default:false"`
}

func (m *WalletAuditLogModel) TableName() string {
	return "wallet_audit_logs"
}

func WalletAuditLogModelToEntity(m *WalletAuditLogModel) *wallet.WalletAuditLog {
	return &wallet.WalletAuditLog{
		Id:         m.ID,
		WalletId:   m.WalletId,
		OldBalance: m.OldBalance,
		NewBalance: m.NewBalance,
		Reason:     m.Reason,
		ChangedBy:  m.ChangedBy,
		ChangedAt:  m.ChangedAt,
	}
}

func WalletAuditLogEntityToModel(e *wallet.WalletAuditLog) *WalletAuditLogModel {
	return &WalletAuditLogModel{
		ID:         e.Id,
		WalletId:   e.WalletId,
		OldBalance: e.OldBalance,
		NewBalance: e.NewBalance,
		Reason:     e.Reason,
		ChangedBy:  e.ChangedBy,
		ChangedAt:  e.ChangedAt,
	}
}

type walletAuditLogRepository struct {
	db *gorm.DB
}

func NewWalletAuditLogRepository(db *gorm.DB) walletuc.WalletAuditRepository {
	return &walletAuditLogRepository{db: db}
}

func (r *walletAuditLogRepository) GetWalletAuditLog(ctx context.Context, req *wallet.WalletAuditLog) (*wallet.WalletAuditLog, error) {
	var model WalletAuditLogModel
	if err := dbFromContext(ctx, r.db).Where("wallet_id = ? AND is_deleted = ?", req.WalletId, false).First(&model).Error; err != nil {
		return nil, err
	}
	return WalletAuditLogModelToEntity(&model), nil
}

func (r *walletAuditLogRepository) CreateAuditLog(ctx context.Context, req *wallet.WalletAuditLog) error {
	model := WalletAuditLogEntityToModel(req)
	if err := dbFromContext(ctx, r.db).Create(model).Error; err != nil {
		return err
	}
	return nil
}
func (r *walletAuditLogRepository) UpdateWalletAuditLog(ctx context.Context, req *wallet.WalletAuditLog) (*wallet.WalletAuditLog, error) {
	model := WalletAuditLogEntityToModel(req)
	if err := dbFromContext(ctx, r.db).Save(model).Error; err != nil {
		return nil, err
	}
	return WalletAuditLogModelToEntity(model), nil
}

func (r *walletAuditLogRepository) DeleteWalletAuditLog(ctx context.Context, req *wallet.WalletAuditLog) (*wallet.WalletAuditLog, error) {
	if err := dbFromContext(ctx, r.db).Model(&WalletAuditLogModel{}).
		Where("id = ?", req.Id).
		Updates(map[string]interface{}{
			"is_deleted": true,
			"deleted_at": time.Now(),
		}).Error; err != nil {
		return nil, err
	}
	return req, nil
}
