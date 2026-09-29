package mapper

import (
	tqdpb "pb/types/tqd"
	"tqd/internal/dto"
)

type WorkspaceMapper struct{}

func NewWorkspaceMapper() *WorkspaceMapper {
	return &WorkspaceMapper{}
}

func (m *WorkspaceMapper) ToProtoPoint(v dto.SpatialPointDTO) *tqdpb.SpatialPoint {
	return &tqdpb.SpatialPoint{
		Lat: v.Lat,
		Lon: v.Lon,
	}
}

func (m *WorkspaceMapper) ToProtoBounds(v dto.SpatialBoundsDTO) *tqdpb.SpatialBounds {
	return &tqdpb.SpatialBounds{
		MinLon: v.MinLon,
		MinLat: v.MinLat,
		MaxLon: v.MaxLon,
		MaxLat: v.MaxLat,
	}
}

func (m *WorkspaceMapper) ToProtoGeometryPreview(v dto.GeometryPreviewDTO) *tqdpb.GeometryPreview {
	return &tqdpb.GeometryPreview{
		Type:     v.Type,
		GeoJson:  v.GeoJSON,
		Centroid: m.ToProtoPoint(v.Centroid),
		Bounds:   m.ToProtoBounds(v.Bounds),
		Srid:     v.SRID,
	}
}

func (m *WorkspaceMapper) ToProtoPreviewHint(v dto.PreviewHintDTO) *tqdpb.PreviewHint {
	return &tqdpb.PreviewHint{
		Mode:            v.Mode,
		Bounds:          m.ToProtoBounds(v.Bounds),
		Centroid:        m.ToProtoPoint(v.Centroid),
		FitPaddingRatio: v.FitPaddingRatio,
		ThumbnailUrl:    v.ThumbnailURL,
		HighlightLayer:  v.HighlightLayer,
	}
}

func (m *WorkspaceMapper) ToProtoFocusTarget(v dto.SpatialFocusTargetDTO) *tqdpb.SpatialFocusTarget {
	return &tqdpb.SpatialFocusTarget{
		Type:     v.Type,
		Bounds:   m.ToProtoBounds(v.Bounds),
		Centroid: m.ToProtoPoint(v.Centroid),
		GeoJson:  v.GeoJSON,
		Srid:     v.SRID,
	}
}

func (m *WorkspaceMapper) ToProtoLocation(v dto.WorkspaceLocationDTO) *tqdpb.WorkspaceLocation {
	return &tqdpb.WorkspaceLocation{
		Address:      v.Address,
		Province:     v.Province,
		ProvinceCode: v.ProvinceCode,
		WardCode:     v.WardCode,
	}
}

func (m *WorkspaceMapper) ToProtoParcelMeta(v dto.WorkspaceParcelMetaDTO) *tqdpb.WorkspaceParcelMeta {
	return &tqdpb.WorkspaceParcelMeta{
		MapNumber:  v.MapNumber,
		LandNumber: v.LandNumber,
		AreaSqm:    v.AreaSqm,
	}
}

func (m *WorkspaceMapper) ToProtoFollowedParcel(v dto.FollowedParcelPreviewDTO) *tqdpb.FollowedParcelPreview {
	return &tqdpb.FollowedParcelPreview{
		Id:              v.ID,
		Type:            v.Type,
		FollowId:        v.FollowID,
		UserId:          v.UserID,
		ParcelId:        v.ParcelID,
		Title:           v.Title,
		Subtitle:        v.Subtitle,
		Parcel:          m.ToProtoParcelMeta(v.Parcel),
		Location:        m.ToProtoLocation(v.Location),
		Lat:             v.Lat,
		Lon:             v.Lon,
		GeometryPreview: m.ToProtoGeometryPreview(v.GeometryPreview),
		Preview:         m.ToProtoPreviewHint(v.Preview),
		FocusTarget:     m.ToProtoFocusTarget(v.FocusTarget),
		FollowedAt:      v.FollowedAt,
		Actions: &tqdpb.FollowedParcelActions{
			Focus:        v.Actions.Focus,
			Remove:       v.Actions.Remove,
			Compare:      v.Actions.Compare,
			CreateReport: v.Actions.CreateReport,
		},
	}
}

func (m *WorkspaceMapper) ToProtoFollowedParcels(items []dto.FollowedParcelPreviewDTO) []*tqdpb.FollowedParcelPreview {
	result := make([]*tqdpb.FollowedParcelPreview, 0, len(items))
	for _, item := range items {
		result = append(result, m.ToProtoFollowedParcel(item))
	}
	return result
}

func (m *WorkspaceMapper) ToProtoHistory(v dto.ViewHistoryPreviewDTO) *tqdpb.ViewHistoryPreview {
	return &tqdpb.ViewHistoryPreview{
		Id:              v.ID,
		Type:            v.Type,
		HistoryId:       v.HistoryID,
		UserId:          v.UserID,
		EntityType:      v.EntityType,
		EntityId:        v.EntityID,
		ParcelId:        v.ParcelID,
		Title:           v.Title,
		Subtitle:        v.Subtitle,
		Parcel:          m.ToProtoParcelMeta(v.Parcel),
		Location:        m.ToProtoLocation(v.Location),
		Lat:             v.Lat,
		Lon:             v.Lon,
		GeometryPreview: m.ToProtoGeometryPreview(v.GeometryPreview),
		Preview:         m.ToProtoPreviewHint(v.Preview),
		FocusTarget:     m.ToProtoFocusTarget(v.FocusTarget),
		ViewedAt:        v.ViewedAt,
		ViewContext: &tqdpb.ViewHistoryContext{
			Source: v.ViewContext.Source,
			Zoom:   v.ViewContext.Zoom,
		},
		Actions: &tqdpb.ViewHistoryActions{
			Focus:  v.Actions.Focus,
			Remove: v.Actions.Remove,
			Follow: v.Actions.Follow,
		},
	}
}

func (m *WorkspaceMapper) ToProtoHistories(items []dto.ViewHistoryPreviewDTO) []*tqdpb.ViewHistoryPreview {
	result := make([]*tqdpb.ViewHistoryPreview, 0, len(items))
	for _, item := range items {
		result = append(result, m.ToProtoHistory(item))
	}
	return result
}

func (m *WorkspaceMapper) ToProtoReport(v dto.GeneratedReportPreviewDTO) *tqdpb.GeneratedReportPreview {
	return &tqdpb.GeneratedReportPreview{
		Id:         v.ID,
		Type:       v.Type,
		ReportId:   v.ReportID,
		UserId:     v.UserID,
		ReportType: v.ReportType,
		Status:     v.Status,
		Title:      v.Title,
		Subtitle:   v.Subtitle,
		Location:   m.ToProtoLocation(v.Location),
		Spatial: &tqdpb.ReportSpatialPreview{
			Bounds:   m.ToProtoBounds(v.Spatial.Bounds),
			Centroid: m.ToProtoPoint(v.Spatial.Centroid),
		},
		Assets: &tqdpb.ReportAssets{
			ThumbnailUrl: v.Assets.ThumbnailURL,
			ImageUrl:     v.Assets.ImageURL,
			PdfUrl:       v.Assets.PDFURL,
			ShareUrl:     v.Assets.ShareURL,
		},
		ReportMeta: &tqdpb.ReportMeta{
			CreatedAt: v.Meta.CreatedAt,
			UpdatedAt: v.Meta.UpdatedAt,
			ExpiresAt: v.Meta.ExpiresAt,
			FileSize:  v.Meta.FileSize,
			Format:    v.Meta.Format,
		},
		Comparison: &tqdpb.ReportComparison{
			FromPlanName: v.Comparison.FromPlanName,
			ToPlanName:   v.Comparison.ToPlanName,
			FromYear:     v.Comparison.FromYear,
			ToYear:       v.Comparison.ToYear,
		},
		Actions: &tqdpb.ReportActions{
			Focus:         v.Actions.Focus,
			DownloadImage: v.Actions.DownloadImage,
			DownloadPdf:   v.Actions.DownloadPDF,
			Share:         v.Actions.Share,
			Regenerate:    v.Actions.Regenerate,
			Remove:        v.Actions.Remove,
		},
	}
}

func (m *WorkspaceMapper) ToProtoReports(items []dto.GeneratedReportPreviewDTO) []*tqdpb.GeneratedReportPreview {
	result := make([]*tqdpb.GeneratedReportPreview, 0, len(items))
	for _, item := range items {
		result = append(result, m.ToProtoReport(item))
	}
	return result
}
