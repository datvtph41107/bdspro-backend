package repo

import (
	_dto "common/domain/dto"
	"context"
	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/dto"
)

type IParcelRepo interface {
	// Preview & Search
	GetParcelPreview(ctx context.Context, parcelID uint64) (*qh_domain.QHParcelInfo, error)
	SearchByAddress(ctx context.Context, keyword string, mapNumber, landNumber string, pagable _dto.Pagable) ([]qh_domain.QHParcelInfo, int64, error)
	SearchPublicParcels(ctx context.Context, pagable _dto.Pagable) ([]qh_domain.QHParcelInfo, int64, error)

	// Basic CRUD
	GetByID(ctx context.Context, id uint64) (*qh_domain.Parcel, error)
	GetByIDs(ctx context.Context, ids []uint64) ([]qh_domain.Parcel, error)
	CreateBatch(ctx context.Context, parcels []qh_domain.Parcel) error
	FindByLocation(ctx context.Context, lat, lng float64) (*qh_domain.Parcel, error)
	FindIntersectPolygon(ctx context.Context, points []Point, limit, offset int, filter *dto.PolygonTargetFilter) ([]qh_domain.Parcel, int64, error)
	FindWithinPolygon(ctx context.Context, points []Point, limit, offset int, filter *dto.PolygonTargetFilter) ([]qh_domain.Parcel, int64, error)
	FindRegionByLocation(ctx context.Context, lat, lng float64) (*dto.RegionInfoResponse, error)
	FindRegionsByPolygon(ctx context.Context, points []Point, intersect bool, zoom uint32, limit int, offset int, filter *dto.PolygonTargetFilter, skipCount bool) ([]dto.RegionInfoResponse, int64, error)
	SummarizePolygonAnalysis(ctx context.Context, points []Point, intersect bool, zoom uint32, filter *dto.PolygonTargetFilter, withQuickInfo bool) (*dto.PolygonAnalysisSummary, []dto.ParcelLayerInfo, error)
	CountParcelsByPolygon(ctx context.Context, points []Point, intersect bool, filter *dto.PolygonTargetFilter) (int64, error)
	GetPolygonQuickInfo(ctx context.Context, points []Point, intersect bool, zoom uint32, filter *dto.PolygonTargetFilter) ([]dto.ParcelLayerInfo, error)

	// Detail API
	GetParcelInfo(ctx context.Context, parcelID uint64) (*dto.ParcelInfoResponse, error)
	GetParcelLayerRows(ctx context.Context, parcelID uint64) ([]dto.ParcelLayerRow, error)
	GetParcelLayerRowsV2(ctx context.Context, parcelID uint64) ([]dto.ParcelLayerRowV2, error)

	GetParcelSeoSource(ctx context.Context, parcelID uint64) (*dto.ParcelSeoSource, error)
	ListParcelSeoSourcesForGenerate(ctx context.Context, limit uint32) ([]dto.ParcelSeoSource, error)
	UpdateParcelSeoID(ctx context.Context, parcelID uint64, seoID uint64) error
	GetParcelQuickInfo(ctx context.Context, parcelID uint64) ([]dto.ParcelLayerInfo, error)

	// GetZoneGeometry — Lấy geometry của một zone cụ thể
	GetZoneGeometry(ctx context.Context, parcelID, layerID, zoneID uint64) (string, *dto.RegionInfoResponse, error)
}

type Point struct {
	Lat float64
	Lng float64
}
