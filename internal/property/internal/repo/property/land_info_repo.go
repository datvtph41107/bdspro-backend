package property_repo

import (
	"bdspro/internal/domain"
	"context"
)

// PropertyLandInfoRepo — quản lý bảng property_land_info.
type PropertyLandInfoRepository interface {
	GetByID(ctx context.Context, id uint64) (*domain.PropertyLandInfo, error)
	GetByLineageID(ctx context.Context, lineageID uint64) (*domain.PropertyLandInfo, error)
	Create(ctx context.Context, entity *domain.PropertyLandInfo) error
	UpdateFields(ctx context.Context, id uint64, fields map[string]any) error
}
