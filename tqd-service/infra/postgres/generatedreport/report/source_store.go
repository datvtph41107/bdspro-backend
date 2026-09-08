package postgresreport

import (
	"context"

	"tqd/internal/dto"
	"tqd/internal/usecase/generatedreport/application"
)

// workspaceReader is the exact Workspace read capability consumed by Report.
//
// Workspace owns parcel/region preview queries. Report owns the translation
// from those read facts into its application language.
type workspaceReader interface {
	GetParcelWorkspacePreview(
		ctx context.Context,
		parcelID uint64,
	) (*dto.ParcelWorkspacePreviewRow, error)

	GetRegionWorkspacePreview(
		ctx context.Context,
		regionID uint64,
	) (*dto.RegionWorkspacePreviewRow, error)
}

type SourceStore struct {
	workspace workspaceReader
}

func NewSourceStore(workspace workspaceReader) *SourceStore {
	return &SourceStore{workspace: workspace}
}

func (s *SourceStore) FindParcel(
	ctx context.Context,
	parcelID uint64,
) (application.Parcel, bool, error) {
	row, err := s.workspace.GetParcelWorkspacePreview(ctx, parcelID)
	if err != nil {
		return application.Parcel{}, false, err
	}
	if row == nil || row.ParcelID == 0 {
		return application.Parcel{}, false, nil
	}

	return application.Parcel{
		ID:         row.ParcelID,
		MapNumber:  row.MapNumber,
		LandNumber: row.LandNumber,
		AreaSqm:    row.AreaSqm,
		Location: application.Location{
			Address:      row.Address,
			Province:     row.Province,
			ProvinceCode: row.ProvinceCode,
			WardCode:     row.WardCode,
		},
		Spatial: application.Spatial{
			Bounds: application.Bounds{
				MinLon: row.MinLon,
				MinLat: row.MinLat,
				MaxLon: row.MaxLon,
				MaxLat: row.MaxLat,
			},
			Centroid: application.Point{
				Lat: row.CentroidLat,
				Lon: row.CentroidLon,
			},
		},
	}, true, nil
}

func (s *SourceStore) FindRegion(
	ctx context.Context,
	regionID uint64,
) (application.Region, bool, error) {
	row, err := s.workspace.GetRegionWorkspacePreview(ctx, regionID)
	if err != nil {
		return application.Region{}, false, err
	}
	if row == nil || row.RegionID == 0 {
		return application.Region{}, false, nil
	}

	return application.Region{
		ID:               row.RegionID,
		Name:             row.Name,
		DisplayName:      row.DisplayName,
		PlanningName:     row.PlanningName,
		LayerName:        row.LayerName,
		LayerDisplayName: row.LayerDisplayName,
		LandUseName:      row.LandUseName,
		AreaSqm:          row.AreaSqm,
		Location: application.Location{
			Address:      row.PlanningName,
			Province:     row.Province,
			ProvinceCode: row.ProvinceCode,
			WardCode:     row.WardCode,
		},
		Spatial: application.Spatial{
			Bounds: application.Bounds{
				MinLon: row.MinLon,
				MinLat: row.MinLat,
				MaxLon: row.MaxLon,
				MaxLat: row.MaxLat,
			},
			Centroid: application.Point{
				Lat: row.CenterLat,
				Lon: row.CenterLon,
			},
		},
	}, true, nil
}
