package usecase

import (
	_dto "common/domain/dto"
	_utils "common/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
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
	TrackViewHistory(ctx context.Context, req dto.TrackViewHistoryRequestDTO) (*dto.TrackViewHistoryResultDTO, error)
	RemoveViewHistory(ctx context.Context, req dto.RemoveViewHistoryRequestDTO) error
	ClearViewHistory(ctx context.Context, userID uint64) error

	ListGeneratedReports(ctx context.Context, userID uint64, pagable _dto.IPagable) (*dto.WorkspaceListResponse[dto.GeneratedReportPreviewDTO], error)
	GetGeneratedReport(ctx context.Context, userID, reportID uint64) (*dto.GeneratedReportPreviewDTO, error)
	RemoveGeneratedReport(ctx context.Context, req dto.RemoveGeneratedReportRequestDTO) error
	RegenerateReport(ctx context.Context, userID, reportID uint64) error
	ShareReport(ctx context.Context, userID, reportID uint64) (*dto.ShareReportResultDTO, error)
}

type mapWorkspaceUsecase struct {
	repo repo.IMapWorkspaceRepo
}

func NewMapWorkspaceUsecase(repo repo.IMapWorkspaceRepo) MapWorkspaceUsecase {
	return &mapWorkspaceUsecase{repo: repo}
}

func timeString(t time.Time) string {
	return _utils.FormatRFC3339(&t)
}

func timeStringPtr(t *time.Time) string {
	return _utils.FormatRFC3339(t)
}

func buildViewDedupeKey(userID uint64, entityType uint32, entityID uint64, source uint32) string {
	return fmt.Sprintf("%d:%d:%d:%d", userID, entityType, entityID, source)
}

func uniqueUint64(values []uint64) []uint64 {
	seen := map[uint64]bool{}
	out := make([]uint64, 0, len(values))

	for _, v := range values {
		if v == 0 || seen[v] {
			continue
		}

		seen[v] = true
		out = append(out, v)
	}

	return out
}

const maxListPreviewGeoJSONBytes = 16 * 1024

func normalizePreviewGeoJSON(s string) (string, string, uint32) {
	s = strings.TrimSpace(s)

	if s == "" || s == "null" {
		return "", "bounds", 0
	}

	size := uint32(len(s))

	if len(s) > maxListPreviewGeoJSONBytes {
		return "", "bounds", size
	}

	return s, "simplified", size
}

func buildParcelGeometryDTO(row dto.ParcelWorkspacePreviewRow) dto.GeometryPreviewDTO {
	geoJSON, quality, byteSize := normalizePreviewGeoJSON(row.GeometryGeoJSON)

	return dto.GeometryPreviewDTO{
		Type:    row.GeometryType,
		GeoJSON: geoJSON,
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
		SRID:     4326,
		Quality:  quality,
		ByteSize: byteSize,
	}
}

func buildRegionGeometryDTO(row dto.RegionWorkspacePreviewRow) dto.GeometryPreviewDTO {
	geoJSON, quality, byteSize := normalizePreviewGeoJSON(row.GeoJSON)

	return dto.GeometryPreviewDTO{
		Type:    row.GeometryType,
		GeoJSON: geoJSON,
		Centroid: dto.SpatialPointDTO{
			Lat: row.CenterLat,
			Lon: row.CenterLon,
		},
		Bounds: dto.SpatialBoundsDTO{
			MinLon: row.MinLon,
			MinLat: row.MinLat,
			MaxLon: row.MaxLon,
			MaxLat: row.MaxLat,
		},
		SRID:     4326,
		Quality:  quality,
		ByteSize: byteSize,
	}
}

func buildPreviewHintDTO(g dto.GeometryPreviewDTO, layer string) dto.PreviewHintDTO {
	return dto.PreviewHintDTO{
		Mode:            enums.PreviewModeGeometry.Uint32(),
		Bounds:          g.Bounds,
		Centroid:        g.Centroid,
		FitPaddingRatio: 0.22,
		HighlightLayer:  layer,
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

func buildListFocusTargetDTO(g dto.GeometryPreviewDTO) dto.SpatialFocusTargetDTO {
	return dto.SpatialFocusTargetDTO{
		Type:     enums.FocusTypeGeometry.Uint32(),
		Bounds:   g.Bounds,
		Centroid: g.Centroid,
		GeoJSON:  "",
		SRID:     g.SRID,
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
		if areaSqm <= 0 {
			return ""
		}
		return fmt.Sprintf("Diện tích %.2f m²", areaSqm)
	}

	if areaSqm <= 0 {
		return "Số tờ " + mapNumber
	}

	return fmt.Sprintf("Số tờ %s · Diện tích %.2f m²", mapNumber, areaSqm)
}

func buildRegionTitle(row dto.RegionWorkspacePreviewRow) string {
	if row.DisplayName != "" {
		return row.DisplayName
	}
	if row.LandUseName != "" {
		return row.LandUseName
	}
	if row.Name != "" {
		return row.Name
	}
	return "Vùng quy hoạch"
}

func buildRegionSubtitle(row dto.RegionWorkspacePreviewRow) string {
	parts := make([]string, 0, 3)

	if row.LayerDisplayName != "" {
		parts = append(parts, row.LayerDisplayName)
	} else if row.LayerName != "" {
		parts = append(parts, row.LayerName)
	}

	if row.LandUseName != "" {
		parts = append(parts, row.LandUseName)
	}

	if row.AreaSqm > 0 {
		parts = append(parts, fmt.Sprintf("Diện tích %.2f m²", row.AreaSqm))
	}

	out := ""
	for i, part := range parts {
		if i > 0 {
			out += " · "
		}
		out += part
	}

	return out
}

func buildParcelPreviewMap(rows []dto.ParcelWorkspacePreviewRow) map[uint64]dto.ParcelWorkspacePreviewRow {
	result := make(map[uint64]dto.ParcelWorkspacePreviewRow, len(rows))
	for _, row := range rows {
		result[row.ParcelID] = row
	}
	return result
}

func buildRegionPreviewMap(rows []dto.RegionWorkspacePreviewRow) map[uint64]dto.RegionWorkspacePreviewRow {
	result := make(map[uint64]dto.RegionWorkspacePreviewRow, len(rows))
	for _, row := range rows {
		result[row.RegionID] = row
	}
	return result
}

func buildFollowedParcelPreview(
	f qh_domain.QHUserFollowedParcel,
	row dto.ParcelWorkspacePreviewRow,
) dto.FollowedParcelPreviewDTO {
	geometry := buildParcelGeometryDTO(row)

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
		Preview:         buildPreviewHintDTO(geometry, "target-parcel"),
		FocusTarget:     buildListFocusTargetDTO(geometry),

		FollowedAt: timeString(f.CreatedAt),

		Actions: dto.FollowedParcelActionsDTO{
			Focus:        true,
			Remove:       true,
			Compare:      true,
			CreateReport: true,
		},
	}
}

func buildParcelHistoryPreview(
	h qh_domain.QHUserViewHistory,
	row dto.ParcelWorkspacePreviewRow,
) dto.ViewHistoryPreviewDTO {
	geometry := buildParcelGeometryDTO(row)

	return dto.ViewHistoryPreviewDTO{
		ID:        fmt.Sprintf("history_%d", h.ID),
		Type:      enums.WorkspaceEntityTypeParcel.Uint32(),
		HistoryID: h.ID,
		UserID:    h.UserID,

		EntityType: enums.WorkspaceEntityTypeParcel.Uint32(),
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
		Preview:         buildPreviewHintDTO(geometry, "target-parcel"),
		FocusTarget:     buildListFocusTargetDTO(geometry),

		ViewedAt: timeString(h.ViewedAt),
		ViewContext: dto.ViewHistoryContextDTO{
			Source: uint32(h.Source),
			Zoom:   h.LastZoom,
		},

		ViewCount:        h.ViewCount,
		CountedViewCount: h.CountedViewCount,

		Actions: dto.ViewHistoryActionsDTO{
			Focus:  true,
			Remove: true,
			Follow: true,
		},
	}
}

func buildRegionHistoryPreview(
	h qh_domain.QHUserViewHistory,
	row dto.RegionWorkspacePreviewRow,
) dto.ViewHistoryPreviewDTO {
	geometry := buildRegionGeometryDTO(row)

	locationLabel := row.PlanningName
	if locationLabel == "" {
		locationLabel = row.LayerDisplayName
	}
	if locationLabel == "" {
		locationLabel = row.Description
	}

	return dto.ViewHistoryPreviewDTO{
		ID:        fmt.Sprintf("history_%d", h.ID),
		Type:      enums.WorkspaceEntityTypeRegion.Uint32(),
		HistoryID: h.ID,
		UserID:    h.UserID,

		EntityType: enums.WorkspaceEntityTypeRegion.Uint32(),
		EntityID:   h.EntityID,
		RegionID:   row.RegionID,

		Title:    buildRegionTitle(row),
		Subtitle: buildRegionSubtitle(row),

		Location: dto.WorkspaceLocationDTO{
			Address:      locationLabel,
			Province:     row.Province,
			ProvinceCode: row.ProvinceCode,
			WardCode:     row.WardCode,
		},

		Lat: row.CenterLat,
		Lon: row.CenterLon,

		GeometryPreview: geometry,
		Preview:         buildPreviewHintDTO(geometry, "target-region"),
		FocusTarget:     buildListFocusTargetDTO(geometry),

		ViewedAt: timeString(h.ViewedAt),
		ViewContext: dto.ViewHistoryContextDTO{
			Source: uint32(h.Source),
			Zoom:   h.LastZoom,
		},

		ViewCount:        h.ViewCount,
		CountedViewCount: h.CountedViewCount,

		Actions: dto.ViewHistoryActionsDTO{
			Focus:  true,
			Remove: true,
			Follow: false,
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
	canDownload := status.CanDownload()

	entityType := enums.WorkspaceEntityTypeReport.Uint32()
	entityID := r.ID

	var parcelID uint64
	if r.ParcelID != nil {
		parcelID = *r.ParcelID
		entityType = enums.WorkspaceEntityTypeParcel.Uint32()
		entityID = parcelID
	}

	var regionID uint64
	if r.RegionID != nil {
		regionID = *r.RegionID
		entityType = enums.WorkspaceEntityTypeRegion.Uint32()
		entityID = regionID
	}

	return dto.GeneratedReportPreviewDTO{
		ID:       fmt.Sprintf("report_%d", r.ID),
		Type:     enums.WorkspaceEntityTypeReport.Uint32(),
		ReportID: r.ID,
		UserID:   r.UserID,

		ReportType: uint32(r.ReportType),
		Status:     status.Uint32(),

		Title:    r.Title,
		Subtitle: r.Subtitle,

		EntityType: entityType,
		EntityID:   entityID,
		ParcelID:   parcelID,
		RegionID:   regionID,

		Location: dto.WorkspaceLocationDTO{
			Address:      r.Address,
			Province:     r.Province,
			ProvinceCode: r.ProvinceCode,
			WardCode:     r.WardCode,
		},
		Spatial: dto.ReportSpatialPreviewDTO{
			Bounds: dto.SpatialBoundsDTO{
				MinLon: r.MinLon,
				MinLat: r.MinLat,
				MaxLon: r.MaxLon,
				MaxLat: r.MaxLat,
			},
			Centroid: dto.SpatialPointDTO{
				Lat: r.CenterLat,
				Lon: r.CenterLon,
			},
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
			Remove:        status.CanRemove(),
		},
	}
}

func (u *mapWorkspaceUsecase) ListFollowedParcels(
	ctx context.Context,
	userID uint64,
	pagable _dto.IPagable,
) (*dto.WorkspaceListResponse[dto.FollowedParcelPreviewDTO], error) {
	if userID == 0 {
		return nil, fmt.Errorf("user_id is required")
	}

	log.Printf(
		"[DEBUG][Usecase][ListFollowedParcels][START] userID=%d limit=%d offset=%d page=%d size=%d",
		userID,
		pagable.GetLimit(),
		pagable.GetOffset(),
		pagable.GetPage(),
		pagable.GetSize(),
	)

	follows, total, err := u.repo.ListFollowedParcels(ctx, userID, pagable.GetLimit(), pagable.GetOffset())
	if err != nil {
		log.Printf(
			"[DEBUG][Usecase][ListFollowedParcels][REPO_FOLLOWS_ERROR] userID=%d err=%v",
			userID,
			err,
		)
		return nil, err
	}

	log.Printf(
		"[DEBUG][Usecase][ListFollowedParcels][FOLLOWS] userID=%d followsNil=%v followsLen=%d total=%d",
		userID,
		follows == nil,
		len(follows),
		total,
	)

	parcelIDs := make([]uint64, 0, len(follows))
	for _, item := range follows {
		parcelIDs = append(parcelIDs, item.ParcelID)
	}

	uniqueIDs := uniqueUint64(parcelIDs)

	log.Printf(
		"[DEBUG][Usecase][ListFollowedParcels][PARCEL_IDS] userID=%d parcelIDsLen=%d uniqueIDsLen=%d uniqueIDs=%v",
		userID,
		len(parcelIDs),
		len(uniqueIDs),
		uniqueIDs,
	)

	parcelRows, err := u.repo.GetParcelWorkspacePreviewsByIDs(ctx, uniqueIDs)
	if err != nil {
		log.Printf(
			"[DEBUG][Usecase][ListFollowedParcels][REPO_PREVIEWS_ERROR] userID=%d err=%v",
			userID,
			err,
		)
		return nil, err
	}

	log.Printf(
		"[DEBUG][Usecase][ListFollowedParcels][PREVIEWS] userID=%d parcelRowsNil=%v parcelRowsLen=%d",
		userID,
		parcelRows == nil,
		len(parcelRows),
	)

	parcelMap := buildParcelPreviewMap(parcelRows)
	items := make([]dto.FollowedParcelPreviewDTO, 0, len(follows))

	missingParcelIDs := make([]uint64, 0)

	for _, follow := range follows {
		row, ok := parcelMap[follow.ParcelID]
		if !ok {
			missingParcelIDs = append(missingParcelIDs, follow.ParcelID)
			continue
		}

		items = append(items, buildFollowedParcelPreview(follow, row))
	}

	log.Printf(
		"[DEBUG][Usecase][ListFollowedParcels][RETURN] userID=%d followsLen=%d parcelRowsLen=%d itemsNil=%v itemsLen=%d total=%d missingPreviewLen=%d missingPreviewIDs=%v",
		userID,
		len(follows),
		len(parcelRows),
		items == nil,
		len(items),
		total,
		len(missingParcelIDs),
		missingParcelIDs,
	)

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

	histories, total, err := u.repo.ListViewHistory(ctx, userID, pagable.GetLimit(), pagable.GetOffset())
	if err != nil {
		return nil, err
	}

	parcelIDs := make([]uint64, 0, len(histories))
	regionIDs := make([]uint64, 0, len(histories))

	for _, item := range histories {
		if item.EntityType == enums.WorkspaceEntityTypeParcel {
			if item.ParcelID != nil && *item.ParcelID > 0 {
				parcelIDs = append(parcelIDs, *item.ParcelID)
			} else {
				parcelIDs = append(parcelIDs, item.EntityID)
			}
		}

		if item.EntityType == enums.WorkspaceEntityTypeRegion {
			if item.RegionID != nil && *item.RegionID > 0 {
				regionIDs = append(regionIDs, *item.RegionID)
			} else {
				regionIDs = append(regionIDs, item.EntityID)
			}
		}
	}

	parcelRows, err := u.repo.GetParcelWorkspacePreviewsByIDs(ctx, uniqueUint64(parcelIDs))
	if err != nil {
		return nil, err
	}

	regionRows, err := u.repo.GetRegionWorkspacePreviewsByIDs(ctx, uniqueUint64(regionIDs))
	if err != nil {
		return nil, err
	}

	parcelMap := buildParcelPreviewMap(parcelRows)
	regionMap := buildRegionPreviewMap(regionRows)
	items := make([]dto.ViewHistoryPreviewDTO, 0, len(histories))

	for _, history := range histories {
		switch history.EntityType {
		case enums.WorkspaceEntityTypeParcel:
			parcelID := history.EntityID
			if history.ParcelID != nil && *history.ParcelID > 0 {
				parcelID = *history.ParcelID
			}

			row, ok := parcelMap[parcelID]
			if !ok {
				continue
			}

			items = append(items, buildParcelHistoryPreview(history, row))

		case enums.WorkspaceEntityTypeRegion:
			regionID := history.EntityID
			if history.RegionID != nil && *history.RegionID > 0 {
				regionID = *history.RegionID
			}

			row, ok := regionMap[regionID]
			if !ok {
				continue
			}

			items = append(items, buildRegionHistoryPreview(history, row))
		}
	}

	return &dto.WorkspaceListResponse[dto.ViewHistoryPreviewDTO]{
		Items: items,
		Total: total,
	}, nil
}

func (u *mapWorkspaceUsecase) AddViewHistory(ctx context.Context, req dto.AddViewHistoryRequestDTO) (uint64, error) {
	result, err := u.TrackViewHistory(ctx, dto.TrackViewHistoryRequestDTO{
		UserID:        req.UserID,
		EntityType:    req.EntityType,
		EntityID:      req.EntityID,
		ParcelID:      req.ParcelID,
		RegionID:      req.RegionID,
		Source:        req.Source,
		Zoom:          req.Zoom,
		VisibleMs:     0,
		CountIntent:   false,
		ClientEventID: fmt.Sprintf("legacy_%d_%d_%d", req.UserID, req.EntityID, time.Now().UnixNano()),
	})
	if err != nil {
		return 0, err
	}

	return result.HistoryID, nil
}

func (u *mapWorkspaceUsecase) TrackViewHistory(ctx context.Context, req dto.TrackViewHistoryRequestDTO) (*dto.TrackViewHistoryResultDTO, error) {
	if req.UserID == 0 {
		return nil, fmt.Errorf("user_id is required")
	}
	if req.EntityID == 0 {
		return nil, fmt.Errorf("entity_id is required")
	}

	entityType := enums.WorkspaceEntityTypeFromUint32(req.EntityType)
	if !entityType.IsValid() {
		entityType = enums.WorkspaceEntityTypeParcel
		req.EntityType = entityType.Uint32()
	}

	source := enums.ViewHistorySourceFromUint32(req.Source)
	if !source.IsValid() {
		source = enums.ViewHistorySourceUnknown
		req.Source = source.Uint32()
	}

	if req.ClientEventID == "" {
		req.ClientEventID = fmt.Sprintf("%d-%d-%d-%d", req.UserID, req.EntityType, req.EntityID, time.Now().UnixNano())
	}

	var parcelID *uint64
	var regionID *uint64

	if entityType == enums.WorkspaceEntityTypeParcel {
		id := req.ParcelID
		if id == 0 {
			id = req.EntityID
		}

		parcel, err := u.repo.GetParcelWorkspacePreview(ctx, id)
		if err != nil {
			return nil, err
		}
		if parcel == nil || parcel.ParcelID == 0 {
			return nil, fmt.Errorf("parcel not found")
		}

		parcelID = &id
	}

	if entityType == enums.WorkspaceEntityTypeRegion {
		id := req.RegionID
		if id == 0 {
			id = req.EntityID
		}

		region, err := u.repo.GetRegionWorkspacePreview(ctx, id)
		if err != nil {
			return nil, err
		}
		if region == nil || region.RegionID == 0 {
			return nil, fmt.Errorf("region not found")
		}

		regionID = &id
	}

	metadata := datatypes.JSON([]byte("{}"))
	if req.MetadataJSON != "" {
		metadata = datatypes.JSON([]byte(req.MetadataJSON))
	}

	event := qh_domain.QHUserViewEvent{
		UserID:        req.UserID,
		ClientEventID: req.ClientEventID,
		SessionID:     req.SessionID,
		DeviceID:      req.DeviceID,

		EntityType: entityType,
		EntityID:   req.EntityID,
		ParcelID:   parcelID,
		RegionID:   regionID,

		Source:    source,
		SourceRef: req.SourceRef,
		RouteName: req.RouteName,

		Zoom:      req.Zoom,
		CenterLat: req.Center.Lat,
		CenterLon: req.Center.Lon,

		ViewportMinLon: req.ViewportBounds.MinLon,
		ViewportMinLat: req.ViewportBounds.MinLat,
		ViewportMaxLon: req.ViewportBounds.MaxLon,
		ViewportMaxLat: req.ViewportBounds.MaxLat,

		VisibleMs:   req.VisibleMs,
		CountIntent: req.CountIntent,
		DedupeKey:   buildViewDedupeKey(req.UserID, req.EntityType, req.EntityID, req.Source),

		Metadata: metadata,
		ViewedAt: _utils.TimeNowUTC(),
	}

	return u.repo.TrackViewHistory(ctx, event, 60*time.Second)
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

	reports, total, err := u.repo.ListGeneratedReports(ctx, userID, pagable.GetLimit(), pagable.GetOffset())
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

func (u *mapWorkspaceUsecase) GetGeneratedReport(ctx context.Context, userID, reportID uint64) (*dto.GeneratedReportPreviewDTO, error) {
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

	item := buildGeneratedReportPreview(*report)
	return &item, nil
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

	report, err := u.repo.GetGeneratedReport(ctx, userID, reportID)
	if err != nil {
		return err
	}
	if report == nil || report.ID == 0 {
		return fmt.Errorf("report not found")
	}

	status := enums.GeneratedReportStatus(report.Status)
	if !status.CanRegenerate() {
		return fmt.Errorf("report cannot regenerate in current status")
	}

	return u.repo.UpdateGeneratedReportStatus(ctx, userID, reportID, enums.GeneratedReportStatusProcessing.Uint32())
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

	status := enums.GeneratedReportStatus(report.Status)
	if !status.CanShare() {
		return nil, fmt.Errorf("report is not ready to share")
	}

	shareURL := report.ShareURL
	if shareURL == "" {
		shareURL = fmt.Sprintf("/share/reports/%d", report.ID)
		if err := u.repo.UpdateGeneratedReportShareURL(ctx, userID, reportID, shareURL); err != nil {
			return nil, err
		}
	}

	return &dto.ShareReportResultDTO{
		Success:  true,
		ShareURL: shareURL,
		Message:  "share url generated",
	}, nil
}
