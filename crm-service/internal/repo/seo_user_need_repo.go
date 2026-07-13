package repo

import (
	"context"
	"errors"

	seo_strategy "crm/internal/domain/seo_strategy"
)

var (
	ErrSeoUserNeedNotFound = errors.New(
		"seo user need không tồn tại",
	)

	ErrSeoUserNeedConflict = errors.New(
		"seo user need đã được cập nhật",
	)

	ErrSeoUserNeedEvidenceAlreadyLinked = errors.New(
		"evidence đã được liên kết",
	)

	ErrSeoUserNeedEvidenceNotLinked = errors.New(
		"evidence chưa được liên kết với user need",
	)
)

type SeoUserNeedFilter struct {
	ObjectiveID uint64
	AudienceID  uint64

	Keyword    string
	Status     string
	Confidence string
	Priority   string

	Page uint32
	Size uint32
}

type SeoUserNeedRepo interface {
	Create(
		ctx context.Context,
		need *seo_strategy.SeoUserNeed,
	) error

	GetByID(
		ctx context.Context,
		id uint64,
	) (*seo_strategy.SeoUserNeed, error)

	GetByCode(
		ctx context.Context,
		objectiveID uint64,
		code string,
	) (*seo_strategy.SeoUserNeed, error)

	List(
		ctx context.Context,
		filter SeoUserNeedFilter,
	) ([]seo_strategy.SeoUserNeed, int64, error)

	UpdateDefinition(
		ctx context.Context,
		need *seo_strategy.SeoUserNeed,
		expectedVersion uint32,
	) error

	UpdateStatus(
		ctx context.Context,
		id uint64,
		status seo_strategy.UserNeedStatus,
		expectedVersion uint32,
	) error

	LinkEvidence(
		ctx context.Context,
		link *seo_strategy.SeoUserNeedEvidenceLink,
	) error

	UnlinkEvidence(
		ctx context.Context,
		userNeedID uint64,
		evidenceID uint64,
	) error

	ListLinkedEvidence(
		ctx context.Context,
		userNeedID uint64,
	) ([]seo_strategy.LinkedUserNeedEvidence, error)
}
