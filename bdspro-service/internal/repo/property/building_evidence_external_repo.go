package property_repo

import (
	"bdspro/internal/domain"
	"context"
)

// PropertyBuildingInfoRepo — quản lý bảng property_building_info.
type PropertyBuildingInfoRepository interface {
	GetByID(ctx context.Context, id uint64) (*domain.PropertyBuildingInfo, error)
	GetByLineageID(ctx context.Context, lineageID uint64) (*domain.PropertyBuildingInfo, error)
	Create(ctx context.Context, entity *domain.PropertyBuildingInfo) error
	UpdateFields(ctx context.Context, id uint64, fields map[string]any) error
}

// PropertyEvidenceRepo — quản lý bảng property_edvidence.
type PropertyEvidenceRepository interface {
	GetByID(ctx context.Context, id uint64) (*domain.PropertyEdvidence, error)
	GetByLineageID(ctx context.Context, lineageID uint64) (*domain.PropertyEdvidence, error)
	Create(ctx context.Context, entity *domain.PropertyEdvidence) error
	UpdateFields(ctx context.Context, id uint64, fields map[string]any) error
}

type PropertyAmenityRepository interface {
	// ReplaceAmenities: xóa toàn bộ liên kết cũ → insert list mới trong cùng tx.
	ReplaceAmenities(ctx context.Context, lineageID uint64, amenityIDs []uint64) error
	GetAmenityIDsByLineageID(ctx context.Context, lineageID uint64) ([]uint64, error)
}

// ─────────────────────────────────────────────────────────────────────────────
// PropertyExternalRefRepo
// ─────────────────────────────────────────────────────────────────────────────

type PropertyExternalRefRepository interface {
	GetByID(ctx context.Context, id uint64) (*domain.PropertyExternalRef, error)
	GetByLineageID(ctx context.Context, lineageID uint64) ([]*domain.PropertyExternalRef, error)
	Create(ctx context.Context, entity *domain.PropertyExternalRef) error
	UpdateFields(ctx context.Context, id uint64, fields map[string]any) error
}

type AreaRegionRepository interface {
	ReplaceAreaRegions(ctx context.Context, lineageID uint64, regionIDs []uint64) error
	GetAreaRegionIDsByLineageID(ctx context.Context, lineageID uint64) ([]uint64, error)
}
