package walletpostgres

import (
	"context"

	"time"

	"payment/internal/domain/wallet"
	walletuc "payment/internal/usecase/wallet"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type WalletModel struct {
	ID        uint32 `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	UserId    uint32
	Balance   int64
	Currency  string

	DeletedAt *time.Time `gorm:"index"`
	IsDeleted bool       `gorm:"default:false"`
}

func (m *WalletModel) TableName() string {
	return "wallets"
}

func WalletModelToEntity(m *WalletModel) *wallet.Wallet {
	return &wallet.Wallet{
		Id:        m.ID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		UserId:    m.UserId,
		Balance:   m.Balance,
		Currency:  m.Currency,
	}
}

func WalletEntityToModel(e *wallet.Wallet) *WalletModel {
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

func NewWalletsRepository(db *gorm.DB) walletuc.WalletRepository {
	return &walletsPostgresRepository{
		db: db,
	}
}

func (r *walletsPostgresRepository) GetWalletById(ctx context.Context, id uint32) (*wallet.Wallet, error) {
	var model WalletModel
	if err := dbFromContext(ctx, r.db).First(&model, id).Error; err != nil {
		return nil, err
	}
	return WalletModelToEntity(&model), nil
}

func (r *walletsPostgresRepository) GetWalletByIdForUpdate(ctx context.Context, id uint32) (*wallet.Wallet, error) {
	var model WalletModel
	if err := dbFromContext(ctx, r.db).Clauses(clause.Locking{Strength: "UPDATE"}).First(&model, id).Error; err != nil {
		return nil, err
	}
	return WalletModelToEntity(&model), nil
}

func (r *walletsPostgresRepository) UpdateWallet(ctx context.Context, wallet *wallet.Wallet) error {
	model := WalletEntityToModel(wallet)
	return dbFromContext(ctx, r.db).Save(model).Error
}

func (r *walletsPostgresRepository) BeginTx(ctx context.Context) walletuc.Transaction {
	return &transaction{tx: dbFromContext(ctx, r.db).Begin()}
}

func (r *walletsPostgresRepository) CreateWallet(ctx context.Context, wallet *wallet.Wallet) (*wallet.Wallet, error) {
	model := WalletEntityToModel(wallet)
	if err := dbFromContext(ctx, r.db).Create(model).Error; err != nil {
		return nil, err
	}
	return WalletModelToEntity(model), nil
}

func (r *walletsPostgresRepository) GetWalletByUserId(ctx context.Context, userId uint32) (*wallet.Wallet, error) {
	var model WalletModel
	if err := dbFromContext(ctx, r.db).Where("user_id = ? AND is_deleted = ?", userId, false).First(&model).Error; err != nil {
		return nil, err
	}
	return WalletModelToEntity(&model), nil
}

func (r *walletsPostgresRepository) GetWalletBySubAccount(ctx context.Context, subAccount string) (*wallet.Wallet, error) {
	var wallet wallet.Wallet
	if err := dbFromContext(ctx, r.db).Where("sub_account = ? AND is_deleted = ?", subAccount, false).First(&wallet).Error; err != nil {
		return nil, err
	}
	return &wallet, nil
}
