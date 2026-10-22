package postgresrender

import (
	"context"
	"errors"
	"fmt"
	qh "mdn4/modules/tqd/entity/qh"
	"mdn4/modules/tqd/services/generatedreport/rendering"
	"strings"

	"gorm.io/gorm"
)

// Source loads only the accepted immutable Report snapshot needed by rendering.
type Source struct {
	db *gorm.DB
}

func NewSource(db *gorm.DB) *Source { return &Source{db: db} }

func (s *Source) LoadForRender(ctx context.Context, reportID, userID uint64, jobID string) (rendering.Document, error) {
	if s == nil || s.db == nil {
		return rendering.Document{}, errors.New("report render database is not configured")
	}
	jobID = strings.TrimSpace(jobID)
	if reportID == 0 || userID == 0 || jobID == "" {
		return rendering.Document{}, errors.New("report render identity is invalid")
	}

	var row qh.QHUserReported
	err := s.db.WithContext(ctx).
		Where("id = ? AND user_id = ? AND job_id = ? AND deleted_at IS NULL", reportID, userID, jobID).
		First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return rendering.Document{}, fmt.Errorf("accepted report snapshot not found")
		}
		return rendering.Document{}, fmt.Errorf("load accepted report snapshot: %w", err)
	}

	return rendering.Document{
		ReportID:       row.ID,
		UserID:         row.UserID,
		JobID:          row.JobID,
		ReportType:     uint32(row.ReportType),
		Title:          row.Title,
		Subtitle:       row.Subtitle,
		Address:        row.Address,
		Province:       row.Province,
		ProvinceCode:   row.ProvinceCode,
		WardCode:       row.WardCode,
		ParcelID:       row.ParcelID,
		RegionID:       row.RegionID,
		MinLon:         row.MinLon,
		MinLat:         row.MinLat,
		MaxLon:         row.MaxLon,
		MaxLat:         row.MaxLat,
		CenterLat:      row.CenterLat,
		CenterLon:      row.CenterLon,
		ComparisonJSON: append([]byte(nil), row.Comparison...),
		MetadataJSON:   append([]byte(nil), row.Metadata...),
		CreatedAt:      row.CreatedAt,
	}, nil
}
