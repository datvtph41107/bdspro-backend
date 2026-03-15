package property_repo

import (
	"bdspro/internal/domain"
	"context"
)

type PropertyMediaRepository interface {
	GetByID(ctx context.Context, id uint64) (*domain.PropertyMedia, error)
	GetByLineageID(ctx context.Context, lineageID uint64) ([]domain.PropertyMedia, error)
	Create(ctx context.Context, entity *domain.PropertyMedia) error
	UpdateFields(ctx context.Context, id uint64, fields map[string]any) error
	// DeleteByLineageIDExcept xóa mọi media của lineage không nằm trong keepIDs. keepIDs rỗng = xóa hết.
	DeleteByLineageIDExcept(ctx context.Context, lineageID uint64, keepIDs []uint64) error
}
