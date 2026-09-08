package repo

import (
	"context"
	"tqd/internal/domain"
)

type ReportRepository interface {
	Create(ctx context.Context, report *domain.Report) error
	Update(ctx context.Context, report *domain.Report) error
	UpdateAssignee(ctx context.Context, id, assigneeID uint64, qaStatus uint32) error
	UpdateStatus(ctx context.Context, id uint64, status uint32, errorMsg *string) error
	UpdateFileInfo(ctx context.Context, id uint64, fileURL *string, fileSize *int64, fileHash *string, completedAt *string) error
	GetByID(ctx context.Context, id uint64) (*domain.Report, error)
	ListByUser(ctx context.Context, userID uint64, reportType *uint32, status *uint32, page, limit int) ([]domain.Report, int64, error)
	AdminList(ctx context.Context, filter AdminReportFilter) ([]domain.Report, int64, error)
	AdminSummary(ctx context.Context) (map[uint32]int64, error)
	AdminQueueSummary(ctx context.Context, assigneeID uint64) (QueueSummary, error)
	Delete(ctx context.Context, id uint64) error
	DeleteExpired(ctx context.Context) error
	// LookupParcelLabels: id → "địa chỉ (Tờ x · Thửa y)" từ bảng parcels.
	LookupParcelLabels(ctx context.Context, ids []uint64) (map[uint64]string, error)
	// LookupLayerLabels: id → display_name/name từ qh_layers.
	LookupLayerLabels(ctx context.Context, ids []uint64) (map[uint64]string, error)

	AddEvent(ctx context.Context, event *domain.ReportEvent) error
	ListEventsByReport(ctx context.Context, reportID uint64, limit int) ([]domain.ReportEvent, error)
	ListEvents(ctx context.Context, filter AdminEventFilter) ([]domain.ReportEvent, int64, error)
}

type AdminReportFilter struct {
	UserID        *uint64
	ReportType    *uint32
	Status        *uint32
	ProblemReport *uint32
	QaStatus      *uint32
	Severity      *uint32
	AssigneeID    *uint64
	// UnassignedOnly: assignee_id IS NULL
	UnassignedOnly bool
	// QaStatuses: lọc theo nhiều QA status (ưu tiên hơn QaStatus nếu set)
	QaStatuses []uint32
	// ExcludeQaStatuses: loại trừ (vd Closed/Rejected)
	ExcludeQaStatuses []uint32
	Q                 string
	Page              int
	Limit             int
}

// QueueSummary đếm 3 bucket hàng đợi (10.8.3).
type QueueSummary struct {
	Mine       int64 `json:"mine"`
	Unassigned int64 `json:"unassigned"`
	DataFix    int64 `json:"dataFix"`
}

type AdminEventFilter struct {
	ReportID *uint64
	ActorID  *uint64
	Action   string
	From     *string // RFC3339
	To       *string
	Page     int
	Limit    int
}
