package repo

import (
	"context"
	"errors"

	seo_strategy "crm/internal/domain/seo_strategy"
)

var (
	ErrSeoIntentHypothesisNotFound    = errors.New("seo intent hypothesis không tồn tại")
	ErrSeoIntentHypothesisConflict    = errors.New("seo intent hypothesis đã được cập nhật")
	ErrSeoIntentHypothesisInvalid     = errors.New("seo intent hypothesis không hợp lệ")
	ErrSeoIntentEvidenceAlreadyLinked = errors.New("query evidence đã được liên kết với intent hypothesis")
	ErrSeoIntentEvidenceNotLinked     = errors.New("query evidence chưa được liên kết với intent hypothesis")
	ErrSeoIntentEvidenceMismatch      = errors.New("intent hypothesis và query evidence không cùng user need/objective")
)

type SeoIntentHypothesisFilter struct {
	ObjectiveID uint64
	UserNeedID  uint64

	Keyword      string
	SearchAction seo_strategy.IntentSearchAction
	SubjectType  seo_strategy.IntentSubjectType
	ScopeType    seo_strategy.IntentScopeType
	Status       seo_strategy.IntentHypothesisStatus
	Confidence   seo_strategy.ObjectiveConfidence

	Page uint32
	Size uint32
}

type SeoIntentHypothesisRepo interface {
	Create(ctx context.Context, intent *seo_strategy.SeoIntentHypothesis) error

	GetByID(ctx context.Context, id uint64) (*seo_strategy.SeoIntentHypothesis, error)

	GetByCode(
		ctx context.Context,
		userNeedID uint64,
		code string,
	) (*seo_strategy.SeoIntentHypothesis, error)

	List(
		ctx context.Context,
		filter SeoIntentHypothesisFilter,
	) ([]seo_strategy.SeoIntentHypothesis, int64, error)

	UpdateDefinition(
		ctx context.Context,
		intent *seo_strategy.SeoIntentHypothesis,
		expectedVersion uint32,
	) error

	UpdateStatus(
		ctx context.Context,
		id uint64,
		status seo_strategy.IntentHypothesisStatus,
		expectedVersion uint32,
	) error

	LinkEvidence(ctx context.Context, link *seo_strategy.SeoIntentEvidenceLink) error

	UnlinkEvidence(ctx context.Context, intentID uint64, queryEvidenceID uint64) error

	ListLinkedEvidence(
		ctx context.Context,
		intentID uint64,
	) ([]seo_strategy.LinkedIntentQueryEvidence, error)
}
