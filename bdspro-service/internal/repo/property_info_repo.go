package repo

import (
	"bdspro/internal/domain"
	"context"
)

type PropertyInfoRepo interface {
	Save(ctx context.Context, info *domain.PropertyInfo) error
}
