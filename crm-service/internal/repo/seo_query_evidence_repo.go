package repo

import (
	"context"
	"errors"
	"time"

	seo_strategy "crm/internal/domain/seo_strategy"
)

var (
	ErrSeoQueryEvidenceNotFound = errors.New(
		"seo query evidence không tồn tại",
	)

	ErrSeoQueryEvidenceConflict = errors.New(
		"seo query evidence xung đột với dữ liệu hiện tại",
	)

	ErrSeoQueryEvidenceInvalid = errors.New(
		"seo query evidence không hợp lệ",
	)

	ErrSeoQueryEvidenceUserNeedMismatch = errors.New(
		"query evidence và user need không cùng objective",
	)
)

type SeoQueryEvidenceFilter struct {
	ObjectiveID uint64
	UserNeedID  uint64

	Keyword      string
	SourceSystem string

	SourceKind      seo_strategy.QueryEvidenceSourceKind
	ObservationType seo_strategy.QueryObservationType
	DeviceCategory  seo_strategy.QueryDeviceCategory
	Confidence      seo_strategy.ObjectiveConfidence
	Status          seo_strategy.QueryEvidenceStatus

	// CurrentAt khác nil nghĩa là chỉ lấy evidence active và chưa hết hạn.
	CurrentAt *time.Time

	Page uint32
	Size uint32
}

type SeoQueryEvidenceRepo interface {
	// Upsert dùng identity:
	// user_need_id + source_system + source_record_key.
	Upsert(
		ctx context.Context,
		evidence *seo_strategy.SeoQueryEvidence,
	) error

	GetByID(
		ctx context.Context,
		id uint64,
	) (*seo_strategy.SeoQueryEvidence, error)

	GetBySourceRecordKey(
		ctx context.Context,
		userNeedID uint64,
		sourceSystem string,
		sourceRecordKey string,
	) (*seo_strategy.SeoQueryEvidence, error)

	List(
		ctx context.Context,
		filter SeoQueryEvidenceFilter,
	) ([]seo_strategy.SeoQueryEvidence, int64, error)

	// UpdateStatus dùng compare-and-set theo expectedStatus.
	UpdateStatus(
		ctx context.Context,
		id uint64,
		expectedStatus seo_strategy.QueryEvidenceStatus,
		targetStatus seo_strategy.QueryEvidenceStatus,
	) error
}
