package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/dto"
	"tqd/internal/interface/repo"

	"gorm.io/gorm"
)

type MapWorkspacePostgres struct {
	DB *gorm.DB
}

func NewMapWorkspacePostgres(db *gorm.DB) repo.IMapWorkspaceRepo {
	return &MapWorkspacePostgres{DB: db}
}

func (p *MapWorkspacePostgres) GetRegionWorkspacePreview(ctx context.Context, regionID uint64) (*dto.RegionWorkspacePreviewRow, error) {
	rows, err := p.GetRegionWorkspacePreviewsByIDs(ctx, []uint64{regionID})
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, nil
	}

	return &rows[0], nil
}

func (p *MapWorkspacePostgres) GetRegionWorkspacePreviewsByIDs(
	ctx context.Context,
	regionIDs []uint64,
) ([]dto.RegionWorkspacePreviewRow, error) {
	rows := make([]dto.RegionWorkspacePreviewRow, 0)

	if len(regionIDs) == 0 {
		return rows, nil
	}

	query := `
		WITH src AS (
			SELECT
				r.id AS region_id,
				r.layer_id,
				COALESCE(r.label_id, 0) AS label_id,
				COALESCE(r.land_use_id, 0) AS land_use_id,
				COALESCE(r.legend_id, 0) AS legend_id,

				COALESCE(r.name, '') AS name,
				COALESCE(r.display_name, '') AS display_name,
				COALESCE(r.description, '') AS description,

				COALESCE(l.name, '') AS layer_name,
				COALESCE(l.display_name, '') AS layer_display_name,

				COALESCE(lu.code, '') AS land_use_code,
				COALESCE(lu.name, '') AS land_use_name,
				COALESCE(lu.color, '') AS land_use_color,
				COALESCE(lu.can_build, false) AS can_build,

				COALESCE(lb.name, '') AS label_name,
				COALESCE(lb.color, '') AS label_color,

				COALESCE(leg.color, '') AS legend_color,
				COALESCE(leg.legend_type, '') AS legend_type,
				COALESCE(leg.geometry_type, '') AS geometry_type,

				COALESCE(r.legal_doc, '') AS legal_doc,
				COALESCE(r.planning_name, '') AS planning_name,

				COALESCE(r.area_sqm, 0) AS area_sqm,
				0 AS area_ha,
				COALESCE(r.perimeter_m, 0) AS perimeter_m,

				ST_Multi(
					ST_CollectionExtract(
						ST_MakeValid(r.geometry),
						3
					)
				)::geometry(MultiPolygon, 4326) AS geom,

				'' AS province,
				'' AS province_code,
				'' AS ward_code
			FROM qh_regions r
			JOIN qh_layers l
				ON l.id = r.layer_id
			   AND l.deleted_at IS NULL
			   AND l.status = 10
			LEFT JOIN qh_labels lb
				ON lb.id = r.label_id
			   AND lb.deleted_at IS NULL
			LEFT JOIN qh_land_use lu
				ON lu.id = r.land_use_id
			LEFT JOIN qh_legends leg
				ON leg.id = r.legend_id
			WHERE r.id IN ?
			  AND r.deleted_at IS NULL
			  AND r.status = 10
			  AND r.is_latest = true
			  AND r.geometry IS NOT NULL
			  AND NOT ST_IsEmpty(r.geometry)
		),
		metric AS (
			SELECT
				src.*,
				ST_Transform(src.geom, 3857) AS geom_m,
				Box2D(ST_Transform(src.geom, 3857)) AS box_m
			FROM src
		),
		simplified AS (
			SELECT
				metric.*,

				/*
				 Region có thể rất lớn.
				 List preview chỉ cần shape tổng thể.
				 - min 2m
				 - max 80m
				 - /220 giữ đủ hình dáng cho thumbnail.
				*/
				LEAST(
					80.0,
					GREATEST(
						2.0,
						GREATEST(
							ST_XMax(metric.box_m) - ST_XMin(metric.box_m),
							ST_YMax(metric.box_m) - ST_YMin(metric.box_m)
						) / 220.0
					)
				) AS tol_m
			FROM metric
		),
		preview AS (
			SELECT
				simplified.*,
				ST_Transform(
					ST_Multi(
						ST_CollectionExtract(
							ST_MakeValid(
								ST_SimplifyPreserveTopology(
									ST_RemoveRepeatedPoints(simplified.geom_m, 0.50),
									simplified.tol_m
								)
							),
							3
						)
					),
					4326
				)::geometry(MultiPolygon, 4326) AS preview_geom
			FROM simplified
		)
		SELECT
			region_id,
			layer_id,
			label_id,
			land_use_id,
			legend_id,

			name,
			display_name,
			description,

			layer_name,
			layer_display_name,

			land_use_code,
			land_use_name,
			land_use_color,
			can_build,

			label_name,
			label_color,

			legend_color,
			legend_type,
			geometry_type,

			legal_doc,
			planning_name,

			area_sqm,
			perimeter_m,

			COALESCE(ST_Y(ST_PointOnSurface(geom)), 0) AS center_lat,
			COALESCE(ST_X(ST_PointOnSurface(geom)), 0) AS center_lon,

			COALESCE(ST_XMin(Box2D(geom)), 0) AS min_lon,
			COALESCE(ST_YMin(Box2D(geom)), 0) AS min_lat,
			COALESCE(ST_XMax(Box2D(geom)), 0) AS max_lon,
			COALESCE(ST_YMax(Box2D(geom)), 0) AS max_lat,

			COALESCE(ST_AsGeoJSON(preview_geom, 6), '') AS geo_json,

			province,
			province_code,
			ward_code
		FROM preview
	`

	if err := p.DB.WithContext(ctx).Raw(query, regionIDs).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("get region workspace previews failed: %w", err)
	}

	return rows, nil
}

func (p *MapWorkspacePostgres) GetParcelWorkspacePreview(ctx context.Context, parcelID uint64) (*dto.ParcelWorkspacePreviewRow, error) {
	rows, err := p.GetParcelWorkspacePreviewsByIDs(ctx, []uint64{parcelID})
	if err != nil {
		return nil, err
	}

	if len(rows) == 0 {
		return nil, nil
	}

	return &rows[0], nil
}

func (p *MapWorkspacePostgres) GetParcelWorkspacePreviewsByIDs(
	ctx context.Context,
	parcelIDs []uint64,
) ([]dto.ParcelWorkspacePreviewRow, error) {
	rows := make([]dto.ParcelWorkspacePreviewRow, 0)

	if len(parcelIDs) == 0 {
		return rows, nil
	}

	query := `
		WITH src AS (
			SELECT
				p.id AS parcel_id,
				COALESCE(info.lat, p.lat, 0) AS lat,
				COALESCE(info.lon, p.lon, 0) AS lon,

				COALESCE(info.total_area_sqm, 0) AS area_sqm,
				COALESCE(info.address_text, '') AS address,

				COALESCE(info.map_number, '') AS map_number,
				COALESCE(info.land_number, '') AS land_number,

				ST_Multi(
					ST_CollectionExtract(
						ST_MakeValid(p.geometry),
						3
					)
				)::geometry(MultiPolygon, 4326) AS geom,

				COALESCE(prov.full_name, '') AS province,
				COALESCE(info.province_code, '') AS province_code,
				COALESCE(info.ward_code, '') AS ward_code
			FROM parcels p
			LEFT JOIN LATERAL (
				SELECT parcel_info.*
				FROM qh_parcel_info parcel_info
				WHERE parcel_info.parcel_id = p.id
				ORDER BY parcel_info.updated_at DESC NULLS LAST, parcel_info.id DESC
				LIMIT 1
			) info ON TRUE
			LEFT JOIN province_v2 prov
				ON info.province_code = prov.code
			WHERE p.id IN ?
			  AND p.deleted_at IS NULL
			  AND p.geometry IS NOT NULL
			  AND NOT ST_IsEmpty(p.geometry)
		),
		metric AS (
			SELECT
				src.*,
				ST_Transform(src.geom, 3857) AS geom_m,
				Box2D(ST_Transform(src.geom, 3857)) AS box_m
			FROM src
		),
		simplified AS (
			SELECT
				metric.*,

				/*
				 Dynamic tolerance theo kích thước parcel.

				 Thumbnail trong SubPanel chỉ khoảng 64px,
				 nên không cần trả full geometry.

				 - min 0.03m: parcel nhỏ vẫn giữ shape.
				 - max 1.50m: parcel lớn không làm payload quá nặng.
				 - /128: tương ứng độ phân giải preview nhỏ.
				*/
				LEAST(
					1.50,
					GREATEST(
						0.03,
						GREATEST(
							ST_XMax(metric.box_m) - ST_XMin(metric.box_m),
							ST_YMax(metric.box_m) - ST_YMin(metric.box_m)
						) / 128.0
					)
				) AS tol_m
			FROM metric
		),
		preview AS (
			SELECT
				simplified.*,
				ST_Transform(
					ST_Multi(
						ST_CollectionExtract(
							ST_MakeValid(
								ST_SimplifyPreserveTopology(
									ST_RemoveRepeatedPoints(simplified.geom_m, 0.01),
									simplified.tol_m
								)
							),
							3
						)
					),
					4326
				)::geometry(MultiPolygon, 4326) AS preview_geom
			FROM simplified
		)
		SELECT
			parcel_id,
			lat,
			lon,

			area_sqm,
			address,

			map_number,
			land_number,

			/*
			 Đây là dòng sửa quan trọng nhất.

			 Cũ:
			   AS geometry_geojson

			 Mới:
			   AS geometry_geo_json

			 để scan đúng vào DTO.GeometryGeoJSON.
			*/
			COALESCE(ST_AsGeoJSON(preview_geom, 6), '') AS geometry_geo_json,

			GeometryType(geom) AS geometry_type,

			COALESCE(ST_XMin(Box2D(geom)), 0) AS min_lon,
			COALESCE(ST_YMin(Box2D(geom)), 0) AS min_lat,
			COALESCE(ST_XMax(Box2D(geom)), 0) AS max_lon,
			COALESCE(ST_YMax(Box2D(geom)), 0) AS max_lat,

			COALESCE(ST_Y(ST_PointOnSurface(geom)), 0) AS centroid_lat,
			COALESCE(ST_X(ST_PointOnSurface(geom)), 0) AS centroid_lon,

			province,
			province_code,
			ward_code
		FROM preview
	`

	if err := p.DB.WithContext(ctx).Raw(query, parcelIDs).Scan(&rows).Error; err != nil {
		slog.InfoContext(ctx, fmt.Sprintf("[MapWorkspaceRepo][GetParcelWorkspacePreviewsByIDs] inputIDsLen=%d err=%v",
			len(parcelIDs),
			err),
		)

		return nil, fmt.Errorf("get parcel workspace previews failed: %w", err)
	}

	return rows, nil
}

func (p *MapWorkspacePostgres) ListFollowedParcels(ctx context.Context, userID uint64, limit, offset int) ([]qh_domain.QHUserFollowedParcel, int64, error) {
	rows := make([]qh_domain.QHUserFollowedParcel, 0)
	var total int64
	slog.DebugContext(ctx, fmt.Sprintf("[DEBUG][Repo][ListFollowedParcels][START] userID=%d limit=%d offset=%d",
		userID,
		limit,
		offset),
	)

	q := p.DB.WithContext(ctx).
		Model(&qh_domain.QHUserFollowedParcel{}).
		Where("user_id = ? AND deleted_at IS NULL", userID)

	if err := q.Count(&total).Error; err != nil {
		slog.DebugContext(ctx, fmt.Sprintf("[DEBUG][Repo][ListFollowedParcels][COUNT_ERROR] userID=%d err=%v",
			userID,
			err),
		)
		return nil, 0, fmt.Errorf("count followed parcels failed: %w", err)
	}
	slog.DebugContext(ctx, fmt.Sprintf("[DEBUG][Repo][ListFollowedParcels][COUNT] userID=%d total=%d",
		userID,
		total),
	)

	if err := q.Order("created_at DESC").Limit(limit).Offset(offset).Find(&rows).Error; err != nil {
		slog.DebugContext(ctx, fmt.Sprintf("[DEBUG][Repo][ListFollowedParcels][FIND_ERROR] userID=%d err=%v",
			userID,
			err),
		)
		return nil, 0, fmt.Errorf("list followed parcels failed: %w", err)
	}

	parcelIDs := make([]uint64, 0, len(rows))
	for _, row := range rows {
		parcelIDs = append(parcelIDs, row.ParcelID)
	}
	slog.DebugContext(ctx, fmt.Sprintf("[DEBUG][Repo][ListFollowedParcels][RETURN] userID=%d rowsNil=%v rowsLen=%d total=%d parcelIDs=%v",
		userID,
		rows == nil,
		len(rows),
		total,
		parcelIDs),
	)

	return rows, total, nil
}

func (p *MapWorkspacePostgres) UpsertFollowedParcel(
	ctx context.Context,
	userID uint64,
	parcelID uint64,
	note string,
) (*qh_domain.QHUserFollowedParcel, error) {
	var row qh_domain.QHUserFollowedParcel
	err := p.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now()

		err := tx.
			Where("user_id = ? AND parcel_id = ? AND deleted_at IS NULL", userID, parcelID).
			First(&row).Error

		if err == nil && row.ID > 0 {
			if err := tx.
				Model(&qh_domain.QHUserFollowedParcel{}).
				Where("id = ?", row.ID).
				Updates(map[string]any{
					"note":       note,
					"updated_at": now,
				}).Error; err != nil {
				return err
			}

			return tx.
				Where("id = ? AND deleted_at IS NULL", row.ID).
				First(&row).Error
		}

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		err = tx.Unscoped().
			Where("user_id = ? AND parcel_id = ? AND deleted_at IS NOT NULL", userID, parcelID).
			Order("id DESC").
			First(&row).Error

		if err == nil && row.ID > 0 {
			if err := tx.Unscoped().
				Model(&qh_domain.QHUserFollowedParcel{}).
				Where("id = ?", row.ID).
				Updates(map[string]any{
					"note":       note,
					"deleted_at": nil,
					"updated_at": now,
				}).Error; err != nil {
				return err
			}

			return tx.
				Where("id = ? AND deleted_at IS NULL", row.ID).
				First(&row).Error
		}

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		row = qh_domain.QHUserFollowedParcel{
			UserID:   userID,
			ParcelID: parcelID,
			Note:     note,
		}

		return tx.Create(&row).Error
	})

	if err != nil {
		return nil, fmt.Errorf("upsert followed parcel failed: %w", err)
	}

	return &row, nil
}

func (p *MapWorkspacePostgres) RemoveFollowedParcel(ctx context.Context, userID, followID, parcelID uint64) error {
	q := p.DB.WithContext(ctx).
		Model(&qh_domain.QHUserFollowedParcel{}).
		Where("user_id = ? AND deleted_at IS NULL", userID)

	if followID > 0 {
		q = q.Where("id = ?", followID)
	} else if parcelID > 0 {
		q = q.Where("parcel_id = ?", parcelID)
	} else {
		return fmt.Errorf("follow_id or parcel_id is required")
	}

	tx := q.Update("deleted_at", time.Now())
	if tx.Error != nil {
		return fmt.Errorf("remove followed parcel failed: %w", tx.Error)
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func boolToUint64(v bool) uint64 {
	if v {
		return 1
	}
	return 0
}

func (p *MapWorkspacePostgres) ListViewHistory(ctx context.Context, userID uint64, limit, offset int) ([]qh_domain.QHUserViewHistory, int64, error) {
	var rows []qh_domain.QHUserViewHistory
	var total int64

	q := p.DB.WithContext(ctx).
		Model(&qh_domain.QHUserViewHistory{}).
		Where("user_id = ? AND deleted_at IS NULL", userID)

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count view history failed: %w", err)
	}

	if err := q.Order("viewed_at DESC").Limit(limit).Offset(offset).Find(&rows).Error; err != nil {
		return nil, 0, fmt.Errorf("list view history failed: %w", err)
	}

	return rows, total, nil
}

func (p *MapWorkspacePostgres) TrackViewHistory(
	ctx context.Context,
	event qh_domain.QHUserViewEvent,
	dedupeWindow time.Duration,
) (*dto.TrackViewHistoryResultDTO, error) {
	result := &dto.TrackViewHistoryResultDTO{}

	err := p.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if event.ClientEventID != "" {
			var existing qh_domain.QHUserViewEvent

			err := tx.
				Where("user_id = ? AND client_event_id = ? AND deleted_at IS NULL", event.UserID, event.ClientEventID).
				First(&existing).Error

			if err == nil && existing.ID > 0 {
				var history qh_domain.QHUserViewHistory
				_ = tx.
					Where("user_id = ? AND entity_type = ? AND entity_id = ? AND deleted_at IS NULL",
						event.UserID, event.EntityType, event.EntityID).
					First(&history).Error

				result.EventID = existing.ID
				result.HistoryID = history.ID
				result.Counted = existing.Counted
				result.ViewCount = history.CountedViewCount
				return nil
			}

			if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}
		}

		counted := false
		if event.CountIntent && event.VisibleMs >= 2500 {
			var recentCount int64
			since := event.ViewedAt.Add(-dedupeWindow)

			if err := tx.Model(&qh_domain.QHUserViewEvent{}).
				Where(`
					user_id = ?
					AND dedupe_key = ?
					AND counted = true
					AND viewed_at >= ?
					AND deleted_at IS NULL
				`, event.UserID, event.DedupeKey, since).
				Count(&recentCount).Error; err != nil {
				return err
			}

			counted = recentCount == 0
		}

		event.Counted = counted

		if err := tx.Create(&event).Error; err != nil {
			return err
		}

		var history qh_domain.QHUserViewHistory
		err := tx.
			Where("user_id = ? AND entity_type = ? AND entity_id = ? AND deleted_at IS NULL",
				event.UserID, event.EntityType, event.EntityID).
			First(&history).Error

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		if history.ID == 0 {
			history = qh_domain.QHUserViewHistory{
				UserID:            event.UserID,
				EntityType:        event.EntityType,
				EntityID:          event.EntityID,
				ParcelID:          event.ParcelID,
				RegionID:          event.RegionID,
				Source:            event.Source,
				ViewedAt:          event.ViewedAt,
				FirstViewedAt:     &event.ViewedAt,
				ViewCount:         1,
				CountedViewCount:  boolToUint64(counted),
				LastZoom:          event.Zoom,
				LastCenterLat:     event.CenterLat,
				LastCenterLon:     event.CenterLon,
				ViewportMinLon:    event.ViewportMinLon,
				ViewportMinLat:    event.ViewportMinLat,
				ViewportMaxLon:    event.ViewportMaxLon,
				ViewportMaxLat:    event.ViewportMaxLat,
				LastClientEventID: event.ClientEventID,
				LastVisibleMs:     event.VisibleMs,
				Metadata:          event.Metadata,
			}

			if err := tx.Create(&history).Error; err != nil {
				return err
			}
		} else {
			updates := map[string]any{
				"parcel_id":            event.ParcelID,
				"region_id":            event.RegionID,
				"source":               event.Source,
				"viewed_at":            event.ViewedAt,
				"view_count":           history.ViewCount + 1,
				"last_zoom":            event.Zoom,
				"last_center_lat":      event.CenterLat,
				"last_center_lon":      event.CenterLon,
				"viewport_min_lon":     event.ViewportMinLon,
				"viewport_min_lat":     event.ViewportMinLat,
				"viewport_max_lon":     event.ViewportMaxLon,
				"viewport_max_lat":     event.ViewportMaxLat,
				"last_client_event_id": event.ClientEventID,
				"last_visible_ms":      event.VisibleMs,
				"metadata":             event.Metadata,
				"updated_at":           time.Now(),
			}

			if counted {
				updates["counted_view_count"] = history.CountedViewCount + 1
			}

			if err := tx.Model(&qh_domain.QHUserViewHistory{}).
				Where("id = ?", history.ID).
				Updates(updates).Error; err != nil {
				return err
			}

			if err := tx.Where("id = ?", history.ID).First(&history).Error; err != nil {
				return err
			}
		}

		result.EventID = event.ID
		result.HistoryID = history.ID
		result.Counted = counted
		result.ViewCount = history.CountedViewCount
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("track view history failed: %w", err)
	}

	return result, nil
}

func (p *MapWorkspacePostgres) RemoveViewHistory(ctx context.Context, userID, historyID uint64) error {
	if historyID == 0 {
		return fmt.Errorf("history_id is required")
	}

	tx := p.DB.WithContext(ctx).
		Model(&qh_domain.QHUserViewHistory{}).
		Where("user_id = ? AND id = ? AND deleted_at IS NULL", userID, historyID).
		Update("deleted_at", time.Now())

	if tx.Error != nil {
		return fmt.Errorf("remove view history failed: %w", tx.Error)
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (p *MapWorkspacePostgres) ClearViewHistory(ctx context.Context, userID uint64) error {
	tx := p.DB.WithContext(ctx).
		Model(&qh_domain.QHUserViewHistory{}).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Update("deleted_at", time.Now())

	if tx.Error != nil {
		return fmt.Errorf("clear view history failed: %w", tx.Error)
	}

	return nil
}

func (p *MapWorkspacePostgres) AddViewHistory(ctx context.Context, item qh_domain.QHUserViewHistory) (*qh_domain.QHUserViewHistory, error) {
	if item.ViewedAt.IsZero() {
		item.ViewedAt = time.Now()
	}

	err := p.DB.WithContext(ctx).Create(&item).Error
	if err != nil {
		return nil, fmt.Errorf("add view history failed: %w", err)
	}

	return &item, nil
}

func (p *MapWorkspacePostgres) ListGeneratedReports(ctx context.Context, userID uint64, limit, offset int) ([]qh_domain.QHUserReported, int64, error) {
	var rows []qh_domain.QHUserReported
	var total int64

	q := p.DB.WithContext(ctx).
		Model(&qh_domain.QHUserReported{}).
		Where("user_id = ? AND deleted_at IS NULL", userID)

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count generated reports failed: %w", err)
	}

	if err := q.Order("created_at DESC").Limit(limit).Offset(offset).Find(&rows).Error; err != nil {
		return nil, 0, fmt.Errorf("list generated reports failed: %w", err)
	}

	return rows, total, nil
}

func (p *MapWorkspacePostgres) GetGeneratedReport(ctx context.Context, userID, reportID uint64) (*qh_domain.QHUserReported, error) {
	var row qh_domain.QHUserReported

	err := p.DB.WithContext(ctx).
		Where("user_id = ? AND id = ? AND deleted_at IS NULL", userID, reportID).
		First(&row).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, fmt.Errorf("get generated report failed: %w", err)
	}

	return &row, nil
}

func (p *MapWorkspacePostgres) RemoveGeneratedReport(ctx context.Context, userID, reportID uint64) error {
	if reportID == 0 {
		return fmt.Errorf("report_id is required")
	}

	tx := p.DB.WithContext(ctx).
		Model(&qh_domain.QHUserReported{}).
		Where("user_id = ? AND id = ? AND deleted_at IS NULL", userID, reportID).
		Update("deleted_at", time.Now())

	if tx.Error != nil {
		return fmt.Errorf("remove generated report failed: %w", tx.Error)
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (p *MapWorkspacePostgres) UpdateGeneratedReportStatus(ctx context.Context, userID, reportID uint64, status uint32) error {
	tx := p.DB.WithContext(ctx).
		Model(&qh_domain.QHUserReported{}).
		Where("user_id = ? AND id = ? AND deleted_at IS NULL", userID, reportID).
		Updates(map[string]any{
			"status":     status,
			"updated_at": time.Now(),
		})

	if tx.Error != nil {
		return fmt.Errorf("update generated report status failed: %w", tx.Error)
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (p *MapWorkspacePostgres) UpdateGeneratedReportShareURL(ctx context.Context, userID, reportID uint64, shareURL string) error {
	tx := p.DB.WithContext(ctx).
		Model(&qh_domain.QHUserReported{}).
		Where("user_id = ? AND id = ? AND deleted_at IS NULL", userID, reportID).
		Updates(map[string]any{
			"share_url":  shareURL,
			"updated_at": time.Now(),
		})

	if tx.Error != nil {
		return fmt.Errorf("update generated report share url failed: %w", tx.Error)
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
