package postgre

import (
	_utils "common/utils"
	"context"

	"gorm.io/gorm"
)

// @bind: social/internal/interface.ITransaction
type PostgreTransaction struct {
	db *gorm.DB
}

func NewPostgreTransaction(db *gorm.DB) *PostgreTransaction {
	return &PostgreTransaction{db: db}
}

func (t *PostgreTransaction) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return t.db.Transaction(func(tx *gorm.DB) error {
		// Truyền tx vào context để repo có thể dùng
		profileId := _utils.GetProfileIdWithContext(ctx)
		organizationId := _utils.GetOrganizationIdFromContext(ctx)

		newCtx := context.WithValue(ctx, "tx", tx)
		newCtx = context.WithValue(newCtx, "profileId", profileId)
		newCtx = context.WithValue(newCtx, "organizationId", organizationId)

		return fn(newCtx)
	})
}

func GetDB(ctx context.Context, defaultDB *gorm.DB) *gorm.DB {
	if tx, ok := ctx.Value("tx").(*gorm.DB); ok && tx != nil {
		return tx.WithContext(ctx)
	}
	return defaultDB
}
