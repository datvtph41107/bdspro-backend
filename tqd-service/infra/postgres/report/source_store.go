package reportpostgres

import (
	"context"
	legacy "tqd/infra/postgres"
	"tqd/internal/report"
)

/**
 * SourceStore dùng lại các query parcel/region hiện có trong giai đoạn target source.
 *
 * Việc bọc lại ở đây giúp internal/report không phụ thuộc package postgres cũ.
 * Khi migrate thật, các query có thể chuyển dần vào package này mà không đổi report service.
 */
type SourceStore struct {
	workspace *legacy.MapWorkspacePostgres
}

func NewSourceStore(workspace *legacy.MapWorkspacePostgres) *SourceStore {
	return &SourceStore{workspace: workspace}
}

func (s *SourceStore) FindParcel(
	ctx context.Context,
	parcelID uint64,
) (report.Parcel, bool, error) {
	row, err := s.workspace.GetParcelWorkspacePreview(ctx, parcelID)
	if err != nil {
		return report.Parcel{}, false, err
	}
	if row == nil || row.ParcelID == 0 {
		return report.Parcel{}, false, nil
	}

	return report.Parcel{
		ID:         row.ParcelID,
		MapNumber:  row.MapNumber,
		LandNumber: row.LandNumber,
		AreaSqm:    row.AreaSqm,
		Location: report.Location{
			Address:      row.Address,
			Province:     row.Province,
			ProvinceCode: row.ProvinceCode,
			WardCode:     row.WardCode,
		},
		Spatial: report.Spatial{
			Bounds: report.Bounds{
				MinLon: row.MinLon,
				MinLat: row.MinLat,
				MaxLon: row.MaxLon,
				MaxLat: row.MaxLat,
			},
			Centroid: report.Point{
				Lat: row.CentroidLat,
				Lon: row.CentroidLon,
			},
		},
	}, true, nil
}

func (s *SourceStore) FindRegion(
	ctx context.Context,
	regionID uint64,
) (report.Region, bool, error) {
	row, err := s.workspace.GetRegionWorkspacePreview(ctx, regionID)
	if err != nil {
		return report.Region{}, false, err
	}
	if row == nil || row.RegionID == 0 {
		return report.Region{}, false, nil
	}

	return report.Region{
		ID:               row.RegionID,
		Name:             row.Name,
		DisplayName:      row.DisplayName,
		PlanningName:     row.PlanningName,
		LayerName:        row.LayerName,
		LayerDisplayName: row.LayerDisplayName,
		LandUseName:      row.LandUseName,
		AreaSqm:          row.AreaSqm,
		Location: report.Location{
			Address:      row.PlanningName,
			Province:     row.Province,
			ProvinceCode: row.ProvinceCode,
			WardCode:     row.WardCode,
		},
		Spatial: report.Spatial{
			Bounds: report.Bounds{
				MinLon: row.MinLon,
				MinLat: row.MinLat,
				MaxLon: row.MaxLon,
				MaxLat: row.MaxLat,
			},
			Centroid: report.Point{
				Lat: row.CenterLat,
				Lon: row.CenterLon,
			},
		},
	}, true, nil
}
