package repo

import (
	"context"
	"errors"

	seo_strategy "crm/internal/domain/seo_strategy"
)

var (
	ErrSeoAudienceNotFound = errors.New(
		"seo audience không tồn tại",
	)

	ErrSeoAudienceConflict = errors.New(
		"seo audience đã được cập nhật",
	)
)

type SeoAudienceFilter struct {
	Keyword      string
	Status       string
	DecisionRole string
	Page         uint32
	Size         uint32
}

type SeoAudienceRepo interface {
	Create(
		ctx context.Context,
		audience *seo_strategy.SeoAudienceSegment,
	) error

	GetByID(
		ctx context.Context,
		id uint64,
	) (*seo_strategy.SeoAudienceSegment, error)

	GetByCode(
		ctx context.Context,
		code string,
	) (*seo_strategy.SeoAudienceSegment, error)

	List(
		ctx context.Context,
		filter SeoAudienceFilter,
	) ([]seo_strategy.SeoAudienceSegment, int64, error)

	UpdateDefinition(
		ctx context.Context,
		audience *seo_strategy.SeoAudienceSegment,
		expectedVersion uint32,
	) error

	UpdateStatus(
		ctx context.Context,
		id uint64,
		status seo_strategy.AudienceStatus,
		expectedVersion uint32,
	) error
}
