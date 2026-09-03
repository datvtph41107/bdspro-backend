package postgresreport

import (
	"context"
	"errors"
	"time"

	reportadmin "tqd/internal/usecase/generatedreport/admin"
)

type adminReportRow struct {
	ReportID   uint64  `gorm:"column:report_id"`
	ProfileID  uint64  `gorm:"column:profile_id"`
	ReportType uint32  `gorm:"column:report_type"`
	Status     uint32  `gorm:"column:status"`
	Title      string  `gorm:"column:title"`
	Subtitle   string  `gorm:"column:subtitle"`
	ParcelID   *uint64 `gorm:"column:parcel_id"`
	RegionID   *uint64 `gorm:"column:region_id"`
	Format     string  `gorm:"column:format"`

	ThumbnailURL string `gorm:"column:thumbnail_url"`
	ImageURL     string `gorm:"column:image_url"`
	PDFURL       string `gorm:"column:pdf_url"`
	ShareURL     string `gorm:"column:share_url"`
	FileSize     uint64 `gorm:"column:file_size"`
	ErrorMessage string `gorm:"column:error_message"`

	JobID           string     `gorm:"column:job_id"`
	JobStatus       string     `gorm:"column:job_status"`
	JobAttempts     int32      `gorm:"column:job_attempts"`
	JobAvailableAt  *time.Time `gorm:"column:job_available_at"`
	JobLockedBy     string     `gorm:"column:job_locked_by"`
	JobClaimVersion int64      `gorm:"column:job_claim_version"`
	JobLastError    string     `gorm:"column:job_last_error"`

	CreatedAt time.Time  `gorm:"column:created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at"`
	ExpiresAt *time.Time `gorm:"column:expires_at"`
}

func (s *ReportStore) ListAdminGeneratedReports(
	ctx context.Context,
	query reportadmin.Query,
) (reportadmin.Page, error) {
	if s == nil || s.db == nil {
		return reportadmin.Page{}, errors.New("report database is not configured")
	}

	base := s.db.WithContext(ctx).
		Table("user_reported AS r").
		Joins("LEFT JOIN report_jobs AS j ON j.report_id = r.id").
		Where("r.deleted_at IS NULL")
	if query.ProfileID != 0 {
		base = base.Where("r.user_id = ?", query.ProfileID)
	}
	if query.Status != 0 {
		base = base.Where("r.status = ?", query.Status)
	}
	if query.JobStatus != "" {
		base = base.Where("j.status = ?", query.JobStatus)
	}

	var total int64
	if err := base.Count(&total).Error; err != nil {
		return reportadmin.Page{}, err
	}

	rows := make([]adminReportRow, 0, query.PageSize)
	err := base.
		Select(`
            r.id AS report_id,
            r.user_id AS profile_id,
            r.report_type,
            r.status,
            r.title,
            r.subtitle,
            r.parcel_id,
            r.region_id,
            r.format,
            r.thumbnail_url,
            r.image_url,
            r.pdf_url,
            r.share_url,
            r.file_size,
            r.error_message,
            COALESCE(j.id, '') AS job_id,
            COALESCE(j.status, '') AS job_status,
            COALESCE(j.attempts, 0) AS job_attempts,
            j.available_at AS job_available_at,
            COALESCE(j.locked_by, '') AS job_locked_by,
            COALESCE(j.claim_version, 0) AS job_claim_version,
            COALESCE(j.last_error, '') AS job_last_error,
            r.created_at,
            r.updated_at,
            r.expires_at
        `).
		Order("r.created_at DESC, r.id DESC").
		Offset(int((query.Page - 1) * query.PageSize)).
		Limit(int(query.PageSize)).
		Scan(&rows).Error
	if err != nil {
		return reportadmin.Page{}, err
	}

	result := reportadmin.Page{
		Reports:  make([]reportadmin.Report, 0, len(rows)),
		Total:    uint64(total),
		Page:     query.Page,
		PageSize: query.PageSize,
	}
	for _, row := range rows {
		item := reportadmin.Report{
			ReportID: row.ReportID, ProfileID: row.ProfileID,
			ReportType: row.ReportType, Status: row.Status,
			Title: row.Title, Subtitle: row.Subtitle, Format: row.Format,
			ThumbnailURL: row.ThumbnailURL, ImageURL: row.ImageURL,
			PDFURL: row.PDFURL, ShareURL: row.ShareURL,
			FileSize: row.FileSize, ErrorMessage: row.ErrorMessage,
			JobID: row.JobID, JobStatus: row.JobStatus,
			JobAttempts: row.JobAttempts, JobAvailableAt: row.JobAvailableAt,
			JobLockedBy: row.JobLockedBy, JobClaimVersion: row.JobClaimVersion,
			JobLastError: row.JobLastError,
			CreatedAt:    row.CreatedAt, UpdatedAt: row.UpdatedAt, ExpiresAt: row.ExpiresAt,
		}
		if row.ParcelID != nil {
			item.ParcelID = *row.ParcelID
		}
		if row.RegionID != nil {
			item.RegionID = *row.RegionID
		}
		result.Reports = append(result.Reports, item)
	}
	return result, nil
}

var _ reportadmin.Repository = (*ReportStore)(nil)
