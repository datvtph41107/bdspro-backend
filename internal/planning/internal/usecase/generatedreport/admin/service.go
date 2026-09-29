package admin

import (
	"context"
	"errors"
	"strings"
	"time"
)

const PermissionReportView = "COMMERCIAL_USAGE_VIEW"

type Query struct {
	Page      uint32
	PageSize  uint32
	ProfileID uint64
	Status    uint32
	JobStatus string
}

type Report struct {
	ReportID   uint64
	ProfileID  uint64
	ReportType uint32
	Status     uint32
	Title      string
	Subtitle   string
	ParcelID   uint64
	RegionID   uint64
	Format     string

	ThumbnailURL string
	ImageURL     string
	PDFURL       string
	ShareURL     string
	FileSize     uint64
	ErrorMessage string

	JobID           string
	JobStatus       string
	JobAttempts     int32
	JobAvailableAt  *time.Time
	JobLockedBy     string
	JobClaimVersion int64
	JobLastError    string

	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt *time.Time
}

type Page struct {
	Reports  []Report
	Total    uint64
	Page     uint32
	PageSize uint32
}

type Repository interface {
	ListAdminGeneratedReports(context.Context, Query) (Page, error)
}

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) List(ctx context.Context, query Query) (Page, error) {
	if s == nil || s.repository == nil {
		return Page{}, errors.New("generated report projection is unavailable")
	}
	if query.Page == 0 {
		query.Page = 1
	}
	if query.PageSize == 0 {
		query.PageSize = 20
	}
	if query.PageSize > 100 {
		return Page{}, errors.New("page size must be between 1 and 100")
	}
	query.JobStatus = strings.TrimSpace(strings.ToLower(query.JobStatus))
	switch query.JobStatus {
	case "", "pending", "running", "completed", "failed":
	default:
		return Page{}, errors.New("invalid report job status")
	}
	return s.repository.ListAdminGeneratedReports(ctx, query)
}
