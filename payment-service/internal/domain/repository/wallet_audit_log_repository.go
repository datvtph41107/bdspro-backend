package repository

import (
	"context"

	"github.com/tuandq2112/payment-service/internal/domain/entity"
)

type WalletAuditLogRepository interface {
	CreateAuditLog(ctx context.Context, log *entity.WalletAuditLog) error
}