package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/dto"
	"tqd/internal/interface/repo"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MapWorkspacePostgres struct {
	DB *gorm.DB
}

func NewMapWorkspacePostgres(db *gorm.DB) repo.IMapWorkspaceRepo {
	return &MapWorkspacePostgres{DB: db}
}

func (p *MapWorkspacePostgres) GetParcelWorkspacePreview(ctx context.Context, parcelID uint64) (*dto.ParcelWorkspacePreviewRow, error) {
	var row dto.ParcelWorkspacePreviewRow

	query := `
		SELECT
			p.id AS parcel_id,
			COALESCE(p.lat, 0) AS lat,
			COALESCE(p.lon, 0) AS lon,

			COALESCE(pi.total_area_sqm, 0) AS area_sqm,
			COALESCE(pi.address_text, '') AS address,

			COALESCE(pi.map_number, '') AS map_number,
			COALESCE(pi.land_number, '') AS land_number,

			COALESCE(pi.property_code, '') AS property_code,
			COALESCE(pi.property_uuid::text, '') AS property_uuid,

			ST_AsGeoJSON(p.geometry) AS geometry_geojson,
			GeometryType(p.geometry) AS geometry_type,

			ST_XMin(Box2D(p.geometry)) AS min_lon,
			ST_YMin(Box2D(p.geometry)) AS min_lat,
			ST_XMax(Box2D(p.geometry)) AS max_lon,
			ST_YMax(Box2D(p.geometry)) AS max_lat,

			ST_Y(ST_PointOnSurface(p.geometry)) AS centroid_lat,
			ST_X(ST_PointOnSurface(p.geometry)) AS centroid_lon,

			COALESCE(prov.full_name, '') AS province,
			COALESCE(pi.province_code, '') AS province_code,
			COALESCE(pi.ward_code, '') AS ward_code
		FROM parcels p
		LEFT JOIN qh_parcel_info pi ON pi.parcel_id = p.id
		LEFT JOIN province_v2 prov ON pi.province_code = prov.code
		WHERE p.id = $1
		  AND p.deleted_at IS NULL
		LIMIT 1
	`

	err := p.DB.WithContext(ctx).Raw(query, parcelID).Scan(&row).Error
	if err != nil {
		return nil, fmt.Errorf("get parcel workspace preview failed: %w", err)
	}

	if row.ParcelID == 0 {
		return nil, nil
	}

	return &row, nil
}

func (p *MapWorkspacePostgres) GetParcelWorkspacePreviewsByIDs(ctx context.Context, parcelIDs []uint64) ([]dto.ParcelWorkspacePreviewRow, error) {
	rows := make([]dto.ParcelWorkspacePreviewRow, 0)

	if len(parcelIDs) == 0 {
		return rows, nil
	}

	query := `
		SELECT
			p.id AS parcel_id,
			COALESCE(p.lat, 0) AS lat,
			COALESCE(p.lon, 0) AS lon,

			COALESCE(pi.total_area_sqm, 0) AS area_sqm,
			COALESCE(pi.address_text, '') AS address,

			COALESCE(pi.map_number, '') AS map_number,
			COALESCE(pi.land_number, '') AS land_number,

			COALESCE(pi.property_code, '') AS property_code,
			COALESCE(pi.property_uuid::text, '') AS property_uuid,

			ST_AsGeoJSON(p.geometry) AS geometry_geojson,
			GeometryType(p.geometry) AS geometry_type,

			ST_XMin(Box2D(p.geometry)) AS min_lon,
			ST_YMin(Box2D(p.geometry)) AS min_lat,
			ST_XMax(Box2D(p.geometry)) AS max_lon,
			ST_YMax(Box2D(p.geometry)) AS max_lat,

			ST_Y(ST_PointOnSurface(p.geometry)) AS centroid_lat,
			ST_X(ST_PointOnSurface(p.geometry)) AS centroid_lon,

			COALESCE(prov.full_name, '') AS province,
			COALESCE(pi.province_code, '') AS province_code,
			COALESCE(pi.ward_code, '') AS ward_code
		FROM parcels p
		LEFT JOIN qh_parcel_info pi ON pi.parcel_id = p.id
		LEFT JOIN province_v2 prov ON pi.province_code = prov.code
		WHERE p.id IN ?
		  AND p.deleted_at IS NULL
	`

	err := p.DB.WithContext(ctx).Raw(query, parcelIDs).Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("get parcel workspace previews failed: %w", err)
	}

	return rows, nil
}

func (p *MapWorkspacePostgres) ListFollowedParcels(ctx context.Context, userID uint64, limit, offset int) ([]qh_domain.QHUserFollowedParcel, int64, error) {
	var rows []qh_domain.QHUserFollowedParcel
	var total int64

	q := p.DB.WithContext(ctx).
		Model(&qh_domain.QHUserFollowedParcel{}).
		Where("user_id = ? AND deleted_at IS NULL", userID)

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("count followed parcels failed: %w", err)
	}

	err := q.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&rows).Error

	if err != nil {
		return nil, 0, fmt.Errorf("list followed parcels failed: %w", err)
	}

	return rows, total, nil
}

func (p *MapWorkspacePostgres) UpsertFollowedParcel(ctx context.Context, userID, parcelID uint64, note string) (*qh_domain.QHUserFollowedParcel, error) {
	row := qh_domain.QHUserFollowedParcel{
		UserID:   userID,
		ParcelID: parcelID,
		Note:     note,
	}

	err := p.DB.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "user_id"},
				{Name: "parcel_id"},
			},
			DoUpdates: clause.Assignments(map[string]any{
				"note":       note,
				"deleted_at": nil,
				"updated_at": time.Now(),
			}),
		}).
		Create(&row).Error

	if err != nil {
		return nil, fmt.Errorf("upsert followed parcel failed: %w", err)
	}

	if row.ID == 0 {
		if err := p.DB.WithContext(ctx).
			Where("user_id = ? AND parcel_id = ? AND deleted_at IS NULL", userID, parcelID).
			First(&row).Error; err != nil {
			return nil, fmt.Errorf("reload followed parcel failed: %w", err)
		}
	}

	return &row, nil
}

func (p *MapWorkspacePostgres) RemoveFollowedParcel(ctx context.Context, userID, followID, parcelID uint64) error {
	q := p.DB.WithContext(ctx).
		Model(&qh_domain.QHUserFollowedParcel{}).
		Where("user_id = ?", userID)

	if followID > 0 {
		q = q.Where("id = ?", followID)
	} else if parcelID > 0 {
		q = q.Where("parcel_id = ?", parcelID)
	} else {
		return fmt.Errorf("follow_id or parcel_id is required")
	}

	if err := q.Update("deleted_at", time.Now()).Error; err != nil {
		return fmt.Errorf("remove followed parcel failed: %w", err)
	}

	return nil
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

	err := q.
		Order("viewed_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&rows).Error

	if err != nil {
		return nil, 0, fmt.Errorf("list view history failed: %w", err)
	}

	return rows, total, nil
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

func (p *MapWorkspacePostgres) RemoveViewHistory(ctx context.Context, userID, historyID uint64) error {
	if historyID == 0 {
		return fmt.Errorf("history_id is required")
	}

	err := p.DB.WithContext(ctx).
		Model(&qh_domain.QHUserViewHistory{}).
		Where("user_id = ? AND id = ?", userID, historyID).
		Update("deleted_at", time.Now()).Error

	if err != nil {
		return fmt.Errorf("remove view history failed: %w", err)
	}

	return nil
}

func (p *MapWorkspacePostgres) ClearViewHistory(ctx context.Context, userID uint64) error {
	err := p.DB.WithContext(ctx).
		Model(&qh_domain.QHUserViewHistory{}).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Update("deleted_at", time.Now()).Error

	if err != nil {
		return fmt.Errorf("clear view history failed: %w", err)
	}

	return nil
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

	err := q.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&rows).Error

	if err != nil {
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

	err := p.DB.WithContext(ctx).
		Model(&qh_domain.QHUserReported{}).
		Where("user_id = ? AND id = ?", userID, reportID).
		Update("deleted_at", time.Now()).Error

	if err != nil {
		return fmt.Errorf("remove generated report failed: %w", err)
	}

	return nil
}

func (p *MapWorkspacePostgres) UpdateGeneratedReportStatus(ctx context.Context, userID, reportID uint64, status uint32) error {
	err := p.DB.WithContext(ctx).
		Model(&qh_domain.QHUserReported{}).
		Where("user_id = ? AND id = ? AND deleted_at IS NULL", userID, reportID).
		Updates(map[string]any{
			"status":     status,
			"updated_at": time.Now(),
		}).Error

	if err != nil {
		return fmt.Errorf("update generated report status failed: %w", err)
	}

	return nil
}

func (p *MapWorkspacePostgres) UpdateGeneratedReportShareURL(ctx context.Context, userID, reportID uint64, shareURL string) error {
	err := p.DB.WithContext(ctx).
		Model(&qh_domain.QHUserReported{}).
		Where("user_id = ? AND id = ? AND deleted_at IS NULL", userID, reportID).
		Updates(map[string]any{
			"share_url":  shareURL,
			"updated_at": time.Now(),
		}).Error

	if err != nil {
		return fmt.Errorf("update generated report share url failed: %w", err)
	}

	return nil
}
