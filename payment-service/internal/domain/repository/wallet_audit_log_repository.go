package repository

import (
	"context"

	"payment/internal/domain/entity"
)

type WalletAuditLogRepository interface {
	CreateAuditLog(ctx context.Context, log *entity.WalletAuditLog) error
}
