package property_repo

import (
	"bdspro/internal/domain"
	"context"
)

type AuditLogRepository interface {
	Create(ctx context.Context, log *domain.PropertyAuditLog) error
	GetByProperty(ctx context.Context, propertyID uint64, limit int) ([]*domain.PropertyAuditLog, error)
}
