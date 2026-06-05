package usecase

import (
	_dto "common/domain/dto"
	_utils "common/utils"
	"context"
	"encoding/json"
	"fmt"
	"time"

	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/dto"
	"tqd/internal/enums"
	"tqd/internal/interface/repo"

	"gorm.io/datatypes"
)

type MapWorkspaceUsecase interface {
	ListFollowedParcels(ctx context.Context, userID uint64, pagable _dto.IPagable) (*dto.WorkspaceListResponse[dto.FollowedParcelPreviewDTO], error)
	FollowParcel(ctx context.Context, req dto.FollowParcelRequestDTO) (uint64, error)
	RemoveFollowedParcel(ctx context.Context, req dto.RemoveFollowedParcelRequestDTO) error

	ListViewHistory(ctx context.Context, userID uint64, pagable _dto.IPagable) (*dto.WorkspaceListResponse[dto.ViewHistoryPreviewDTO], error)
	AddViewHistory(ctx context.Context, req dto.AddViewHistoryRequestDTO) (uint64, error)
	RemoveViewHistory(ctx context.Context, req dto.RemoveViewHistoryRequestDTO) error
	ClearViewHistory(ctx context.Context, userID uint64) error

	ListGeneratedReports(ctx context.Context, userID uint64, pagable _dto.IPagable) (*dto.WorkspaceListResponse[dto.GeneratedReportPreviewDTO], error)
	RemoveGeneratedReport(ctx context.Context, req dto.RemoveGeneratedReportRequestDTO) error
	RegenerateReport(ctx context.Context, userID, reportID uint64) error
	ShareReport(ctx context.Context, userID, reportID uint64) (*dto.ShareReportResultDTO, error)
}

type mapWorkspaceUsecase struct {
	repo repo.IMapWorkspaceRepo
}

func NewMapWorkspaceUsecase(repo repo.IMapWorkspaceRepo) MapWorkspaceUsecase {
	return &mapWorkspaceUsecase{
		repo: repo,
	}
}

func timeString(t time.Time) string {
	return _utils.FormatRFC3339(&t)
}

func timeStringPtr(t *time.Time) string {
	return _utils.FormatRFC3339(t)
}

func buildParcelGeometryDTO(row dto.ParcelWorkspacePreviewRow) dto.GeometryPreviewDTO {
	return dto.GeometryPreviewDTO{
		Type:    row.GeometryType,
		GeoJSON: row.GeometryGeoJSON,
		Centroid: dto.SpatialPointDTO{
			Lat: row.CentroidLat,
			Lon: row.CentroidLon,
		},
		Bounds: dto.SpatialBoundsDTO{
			MinLon: row.MinLon,
			MinLat: row.MinLat,
			MaxLon: row.MaxLon,
			MaxLat: row.MaxLat,
		},
		SRID: 4326,
	}
}

func buildPreviewHintDTO(g dto.GeometryPreviewDTO) dto.PreviewHintDTO {
	return dto.PreviewHintDTO{
		Mode:            enums.PreviewModeGeometry.Uint32(),
		Bounds:          g.Bounds,
		Centroid:        g.Centroid,
		FitPaddingRatio: 0.22,
		HighlightLayer:  "target-parcel",
	}
}

func buildGeometryFocusTargetDTO(g dto.GeometryPreviewDTO) dto.SpatialFocusTargetDTO {
	return dto.SpatialFocusTargetDTO{
		Type:     enums.FocusTypeGeometry.Uint32(),
		Bounds:   g.Bounds,
		Centroid: g.Centroid,
		GeoJSON:  g.GeoJSON,
		SRID:     g.SRID,
	}
}

func buildBoundsFocusTargetDTO(bounds dto.SpatialBoundsDTO, centroid dto.SpatialPointDTO) dto.SpatialFocusTargetDTO {
	return dto.SpatialFocusTargetDTO{
		Type:     enums.FocusTypeBounds.Uint32(),
		Bounds:   bounds,
		Centroid: centroid,
		SRID:     4326,
	}
}

func buildParcelTitle(landNumber string) string {
	if landNumber == "" || landNumber == "0" {
		return "Thửa đất"
	}
	return "Thửa " + landNumber
}

func buildParcelSubtitle(mapNumber string, areaSqm float64) string {
	if mapNumber == "" || mapNumber == "0" {
		return fmt.Sprintf("Diện tích %.2f m²", areaSqm)
	}
	return fmt.Sprintf("Số tờ %s · Diện tích %.2f m²", mapNumber, areaSqm)
}

func buildFollowedParcelPreview(f qh_domain.QHUserFollowedParcel, row dto.ParcelWorkspacePreviewRow) dto.FollowedParcelPreviewDTO {
	geometry := buildParcelGeometryDTO(row)
	preview := buildPreviewHintDTO(geometry)

	return dto.FollowedParcelPreviewDTO{
		ID:       fmt.Sprintf("follow_%d", f.ID),
		Type:     enums.WorkspaceEntityTypeParcel.Uint32(),
		FollowID: f.ID,
		UserID:   f.UserID,
		ParcelID: row.ParcelID,

		Title:    buildParcelTitle(row.LandNumber),
		Subtitle: buildParcelSubtitle(row.MapNumber, row.AreaSqm),

		Parcel: dto.WorkspaceParcelMetaDTO{
			MapNumber:  row.MapNumber,
			LandNumber: row.LandNumber,
			AreaSqm:    row.AreaSqm,
		},
		Location: dto.WorkspaceLocationDTO{
			Address:      row.Address,
			Province:     row.Province,
			ProvinceCode: row.ProvinceCode,
			WardCode:     row.WardCode,
		},

		Lat: row.Lat,
		Lon: row.Lon,

		GeometryPreview: geometry,
		Preview:         preview,
		FocusTarget:     buildGeometryFocusTargetDTO(geometry),

		FollowedAt: timeString(f.CreatedAt),

		Actions: dto.FollowedParcelActionsDTO{
			Focus:        true,
			Remove:       true,
			Compare:      true,
			CreateReport: true,
		},
	}
}

func buildHistoryPreview(h qh_domain.QHUserViewHistory, row dto.ParcelWorkspacePreviewRow) dto.ViewHistoryPreviewDTO {
	geometry := buildParcelGeometryDTO(row)
	preview := buildPreviewHintDTO(geometry)

	return dto.ViewHistoryPreviewDTO{
		ID:        fmt.Sprintf("history_%d", h.ID),
		Type:      enums.WorkspaceEntityTypeParcel.Uint32(),
		HistoryID: h.ID,
		UserID:    h.UserID,

		EntityType: uint32(h.EntityType),
		EntityID:   h.EntityID,
		ParcelID:   row.ParcelID,

		Title:    buildParcelTitle(row.LandNumber),
		Subtitle: buildParcelSubtitle(row.MapNumber, row.AreaSqm),

		Parcel: dto.WorkspaceParcelMetaDTO{
			MapNumber:  row.MapNumber,
			LandNumber: row.LandNumber,
			AreaSqm:    row.AreaSqm,
		},
		Location: dto.WorkspaceLocationDTO{
			Address:      row.Address,
			Province:     row.Province,
			ProvinceCode: row.ProvinceCode,
			WardCode:     row.WardCode,
		},

		Lat: row.Lat,
		Lon: row.Lon,

		GeometryPreview: geometry,
		Preview:         preview,
		FocusTarget:     buildGeometryFocusTargetDTO(geometry),

		ViewedAt: timeString(h.ViewedAt),
		ViewContext: dto.ViewHistoryContextDTO{
			Source: uint32(h.Source),
		},

		Actions: dto.ViewHistoryActionsDTO{
			Focus:  true,
			Remove: true,
			Follow: true,
		},
	}
}

func decodeReportComparison(raw datatypes.JSON) dto.ReportComparisonDTO {
	var result dto.ReportComparisonDTO

	if len(raw) == 0 || string(raw) == "null" {
		return result
	}

	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		return result
	}

	return result
}

func buildGeneratedReportPreview(r qh_domain.QHUserReported) dto.GeneratedReportPreviewDTO {
	status := enums.GeneratedReportStatus(r.Status)

	bounds := dto.SpatialBoundsDTO{
		MinLon: r.MinLon,
		MinLat: r.MinLat,
		MaxLon: r.MaxLon,
		MaxLat: r.MaxLat,
	}

	centroid := dto.SpatialPointDTO{
		Lat: r.CenterLat,
		Lon: r.CenterLon,
	}

	canDownload := status.CanDownload()

	return dto.GeneratedReportPreviewDTO{
		ID:       fmt.Sprintf("report_%d", r.ID),
		Type:     enums.WorkspaceEntityTypeReport.Uint32(),
		ReportID: r.ID,
		UserID:   r.UserID,

		ReportType: uint32(r.ReportType),
		Status:     status.Uint32(),

		Title:    r.Title,
		Subtitle: r.Subtitle,

		Location: dto.WorkspaceLocationDTO{
			Address:      r.Address,
			Province:     r.Province,
			ProvinceCode: r.ProvinceCode,
			WardCode:     r.WardCode,
		},
		Spatial: dto.ReportSpatialPreviewDTO{
			Bounds:   bounds,
			Centroid: centroid,
		},
		Assets: dto.ReportAssetsDTO{
			ThumbnailURL: r.ThumbnailURL,
			ImageURL:     r.ImageURL,
			PDFURL:       r.PDFURL,
			ShareURL:     r.ShareURL,
		},
		Meta: dto.ReportMetaDTO{
			CreatedAt: timeString(r.CreatedAt),
			UpdatedAt: timeString(r.UpdatedAt),
			ExpiresAt: timeStringPtr(r.ExpiresAt),
			FileSize:  r.FileSize,
			Format:    r.Format,
		},
		Comparison: decodeReportComparison(r.Comparison),
		Actions: dto.ReportActionsDTO{
			Focus:         r.MinLon != 0 && r.MinLat != 0 && r.MaxLon != 0 && r.MaxLat != 0,
			DownloadImage: canDownload && r.ImageURL != "",
			DownloadPDF:   canDownload && r.PDFURL != "",
			Share:         status.CanShare(),
			Regenerate:    status.CanRegenerate(),
			Remove:        true,
		},
	}
}

func buildParcelPreviewMap(rows []dto.ParcelWorkspacePreviewRow) map[uint64]dto.ParcelWorkspacePreviewRow {
	result := make(map[uint64]dto.ParcelWorkspacePreviewRow, len(rows))

	for _, row := range rows {
		result[row.ParcelID] = row
	}

	return result
}

func (u *mapWorkspaceUsecase) ListFollowedParcels(
	ctx context.Context,
	userID uint64,
	pagable _dto.IPagable,
) (*dto.WorkspaceListResponse[dto.FollowedParcelPreviewDTO], error) {
	if userID == 0 {
		return nil, fmt.Errorf("user_id is required")
	}

	follows, total, err := u.repo.ListFollowedParcels(
		ctx,
		userID,
		pagable.GetLimit(),
		pagable.GetOffset(),
	)
	if err != nil {
		return nil, err
	}

	parcelIDs := make([]uint64, 0, len(follows))
	for _, item := range follows {
		parcelIDs = append(parcelIDs, item.ParcelID)
	}

	parcelRows, err := u.repo.GetParcelWorkspacePreviewsByIDs(ctx, parcelIDs)
	if err != nil {
		return nil, err
	}

	parcelMap := buildParcelPreviewMap(parcelRows)

	items := make([]dto.FollowedParcelPreviewDTO, 0, len(follows))
	for _, follow := range follows {
		row, ok := parcelMap[follow.ParcelID]
		if !ok {
			continue
		}

		items = append(items, buildFollowedParcelPreview(follow, row))
	}

	return &dto.WorkspaceListResponse[dto.FollowedParcelPreviewDTO]{
		Items: items,
		Total: total,
	}, nil
}

func (u *mapWorkspaceUsecase) FollowParcel(ctx context.Context, req dto.FollowParcelRequestDTO) (uint64, error) {
	if req.UserID == 0 {
		return 0, fmt.Errorf("user_id is required")
	}
	if req.ParcelID == 0 {
		return 0, fmt.Errorf("parcel_id is required")
	}

	parcelPreview, err := u.repo.GetParcelWorkspacePreview(ctx, req.ParcelID)
	if err != nil {
		return 0, err
	}
	if parcelPreview == nil || parcelPreview.ParcelID == 0 {
		return 0, fmt.Errorf("parcel not found")
	}

	row, err := u.repo.UpsertFollowedParcel(ctx, req.UserID, req.ParcelID, req.Note)
	if err != nil {
		return 0, err
	}

	return row.ID, nil
}

func (u *mapWorkspaceUsecase) RemoveFollowedParcel(ctx context.Context, req dto.RemoveFollowedParcelRequestDTO) error {
	if req.UserID == 0 {
		return fmt.Errorf("user_id is required")
	}

	if req.FollowID == 0 && req.ParcelID == 0 {
		return fmt.Errorf("follow_id or parcel_id is required")
	}

	return u.repo.RemoveFollowedParcel(ctx, req.UserID, req.FollowID, req.ParcelID)
}

func (u *mapWorkspaceUsecase) ListViewHistory(
	ctx context.Context,
	userID uint64,
	pagable _dto.IPagable,
) (*dto.WorkspaceListResponse[dto.ViewHistoryPreviewDTO], error) {
	if userID == 0 {
		return nil, fmt.Errorf("user_id is required")
	}

	histories, total, err := u.repo.ListViewHistory(
		ctx,
		userID,
		pagable.GetLimit(),
		pagable.GetOffset(),
	)
	if err != nil {
		return nil, err
	}

	parcelIDs := make([]uint64, 0, len(histories))
	for _, item := range histories {
		if item.ParcelID != nil && *item.ParcelID > 0 {
			parcelIDs = append(parcelIDs, *item.ParcelID)
			continue
		}

		if item.EntityType == enums.WorkspaceEntityType(enums.WorkspaceEntityTypeParcel) {
			parcelIDs = append(parcelIDs, item.EntityID)
		}
	}

	parcelRows, err := u.repo.GetParcelWorkspacePreviewsByIDs(ctx, parcelIDs)
	if err != nil {
		return nil, err
	}

	parcelMap := buildParcelPreviewMap(parcelRows)

	items := make([]dto.ViewHistoryPreviewDTO, 0, len(histories))
	for _, history := range histories {
		var parcelID uint64

		if history.ParcelID != nil && *history.ParcelID > 0 {
			parcelID = *history.ParcelID
		} else if history.EntityType == enums.WorkspaceEntityType(enums.WorkspaceEntityTypeParcel) {
			parcelID = history.EntityID
		}

		row, ok := parcelMap[parcelID]
		if !ok {
			continue
		}

		items = append(items, buildHistoryPreview(history, row))
	}

	return &dto.WorkspaceListResponse[dto.ViewHistoryPreviewDTO]{
		Items: items,
		Total: total,
	}, nil
}

func (u *mapWorkspaceUsecase) AddViewHistory(ctx context.Context, req dto.AddViewHistoryRequestDTO) (uint64, error) {
	if req.UserID == 0 {
		return 0, fmt.Errorf("user_id is required")
	}
	if req.EntityID == 0 {
		return 0, fmt.Errorf("entity_id is required")
	}
	if req.EntityType == 0 {
		req.EntityType = enums.WorkspaceEntityTypeParcel.Uint32()
	}
	if req.Source == 0 {
		req.Source = enums.ViewHistorySourceUnknown.Uint32()
	}

	var parcelID *uint64
	if req.ParcelID > 0 {
		parcelID = &req.ParcelID
	} else if req.EntityType == enums.WorkspaceEntityTypeParcel.Uint32() {
		parcelID = &req.EntityID
	}

	row, err := u.repo.AddViewHistory(ctx, qh_domain.QHUserViewHistory{
		UserID:     req.UserID,
		EntityType: enums.WorkspaceEntityType(req.EntityType),
		EntityID:   req.EntityID,
		ParcelID:   parcelID,
		Source:     enums.ViewHistorySource(req.Source),
		ViewedAt:   _utils.TimeNowUTC(),
	})
	if err != nil {
		return 0, err
	}

	return row.ID, nil
}

func (u *mapWorkspaceUsecase) RemoveViewHistory(ctx context.Context, req dto.RemoveViewHistoryRequestDTO) error {
	if req.UserID == 0 {
		return fmt.Errorf("user_id is required")
	}
	if req.HistoryID == 0 {
		return fmt.Errorf("history_id is required")
	}

	return u.repo.RemoveViewHistory(ctx, req.UserID, req.HistoryID)
}

func (u *mapWorkspaceUsecase) ClearViewHistory(ctx context.Context, userID uint64) error {
	if userID == 0 {
		return fmt.Errorf("user_id is required")
	}

	return u.repo.ClearViewHistory(ctx, userID)
}

func (u *mapWorkspaceUsecase) ListGeneratedReports(
	ctx context.Context,
	userID uint64,
	pagable _dto.IPagable,
) (*dto.WorkspaceListResponse[dto.GeneratedReportPreviewDTO], error) {
	if userID == 0 {
		return nil, fmt.Errorf("user_id is required")
	}

	reports, total, err := u.repo.ListGeneratedReports(
		ctx,
		userID,
		pagable.GetLimit(),
		pagable.GetOffset(),
	)
	if err != nil {
		return nil, err
	}

	items := make([]dto.GeneratedReportPreviewDTO, 0, len(reports))
	for _, report := range reports {
		items = append(items, buildGeneratedReportPreview(report))
	}

	return &dto.WorkspaceListResponse[dto.GeneratedReportPreviewDTO]{
		Items: items,
		Total: total,
	}, nil
}

func (u *mapWorkspaceUsecase) RemoveGeneratedReport(ctx context.Context, req dto.RemoveGeneratedReportRequestDTO) error {
	if req.UserID == 0 {
		return fmt.Errorf("user_id is required")
	}
	if req.ReportID == 0 {
		return fmt.Errorf("report_id is required")
	}

	return u.repo.RemoveGeneratedReport(ctx, req.UserID, req.ReportID)
}

func (u *mapWorkspaceUsecase) RegenerateReport(ctx context.Context, userID, reportID uint64) error {
	if userID == 0 {
		return fmt.Errorf("user_id is required")
	}
	if reportID == 0 {
		return fmt.Errorf("report_id is required")
	}

	return u.repo.UpdateGeneratedReportStatus(
		ctx,
		userID,
		reportID,
		enums.GeneratedReportStatusProcessing.Uint32(),
	)
}

func (u *mapWorkspaceUsecase) ShareReport(ctx context.Context, userID, reportID uint64) (*dto.ShareReportResultDTO, error) {
	if userID == 0 {
		return nil, fmt.Errorf("user_id is required")
	}
	if reportID == 0 {
		return nil, fmt.Errorf("report_id is required")
	}

	report, err := u.repo.GetGeneratedReport(ctx, userID, reportID)
	if err != nil {
		return nil, err
	}
	if report == nil || report.ID == 0 {
		return nil, fmt.Errorf("report not found")
	}

	shareURL := report.ShareURL
	// if shareURL == "" {
	// 	shareURL = fmt.Sprintf("/share/reports/%d", report.ID)

	// 	if err := u.repo.UpdateGeneratedReportShareURL(ctx, userID, reportID, shareURL); err != nil {
	// 		return nil, err
	// 	}
	// }

	return &dto.ShareReportResultDTO{
		Success:  true,
		ShareURL: shareURL,
		Message:  "share url generated",
	}, nil
}
