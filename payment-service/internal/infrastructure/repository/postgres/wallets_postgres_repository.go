package postgres

import (
	"context"

	"time"

	"payment/internal/domain/entity"
	"payment/internal/domain/repository"

	"gorm.io/gorm"
)

type WalletModel struct {
	ID        uint32 `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	UserId    uint32
	Balance   float64
	Currency  string

	DeletedAt *time.Time `gorm:"index"`
	IsDeleted bool       `gorm:"default:false"`
}

func (m *WalletModel) TableName() string {
	return "wallets"
}

func WalletModelToEntity(m *WalletModel) *entity.Wallet {
	return &entity.Wallet{
		Id:        m.ID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		UserId:    m.UserId,
		Balance:   m.Balance,
		Currency:  m.Currency,
	}
}

func WalletEntityToModel(e *entity.Wallet) *WalletModel {
	return &WalletModel{
		ID:        e.Id,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
		UserId:    e.UserId,
		Balance:   e.Balance,
		Currency:  e.Currency,
	}
}

type walletsPostgresRepository struct {
	db *gorm.DB
}

type transaction struct {
	tx *gorm.DB
}

func (t *transaction) Commit() error {
	return t.tx.Commit().Error
}

func (t *transaction) Rollback() error {
	return t.tx.Rollback().Error
}

func NewWalletsRepository(db *gorm.DB) repository.WalletsRepository {
	return &walletsPostgresRepository{
		db: db,
	}
}

func (r *walletsPostgresRepository) GetWalletById(ctx context.Context, id uint32) (*entity.Wallet, error) {
	var model WalletModel
	if err := r.db.WithContext(ctx).First(&model, id).Error; err != nil {
		return nil, err
	}
	return WalletModelToEntity(&model), nil
}

func (r *walletsPostgresRepository) GetWalletByIdForUpdate(ctx context.Context, id uint32) (*entity.Wallet, error) {
	var model WalletModel
	if err := r.db.WithContext(ctx).Set("gorm:for_update", true).First(&model, id).Error; err != nil {
		return nil, err
	}
	return WalletModelToEntity(&model), nil
}

func (r *walletsPostgresRepository) UpdateWallet(ctx context.Context, wallet *entity.Wallet) error {
	model := WalletEntityToModel(wallet)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *walletsPostgresRepository) BeginTx(ctx context.Context) repository.Transaction {
	return &transaction{
		tx: r.db.WithContext(ctx).Begin(),
	}
}

func (r *walletsPostgresRepository) CreateWallet(ctx context.Context, wallet *entity.Wallet) (*entity.Wallet, error) {
	model := WalletEntityToModel(wallet)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return nil, err
	}
	return WalletModelToEntity(model), nil
}

func (r *walletsPostgresRepository) GetWalletByUserId(ctx context.Context, userId uint32) (*entity.Wallet, error) {
	var model WalletModel
	if err := r.db.WithContext(ctx).Where("user_id = ? AND is_deleted = ?", userId, false).First(&model).Error; err != nil {
		return nil, err
	}
	return WalletModelToEntity(&model), nil
}

func (r *walletsPostgresRepository) GetWalletBySubAccount(ctx context.Context, subAccount string) (*entity.Wallet, error) {
	var wallet entity.Wallet
	if err := r.db.WithContext(ctx).Where("sub_account = ? AND is_deleted = ?", subAccount, false).First(&wallet).Error; err != nil {
		return nil, err
	}
	return &wallet, nil
}
