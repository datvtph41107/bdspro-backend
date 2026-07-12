package repo

import (
	"context"
	"errors"

	seo_strategy "crm/internal/domain/seo_strategy"
)

var (
	ErrSeoObjectiveNotFound = errors.New(
		"seo objective không tồn tại",
	)

	ErrSeoObjectiveConflict = errors.New(
		"seo objective đã được cập nhật bởi thao tác khác",
	)
)

type SeoObjectiveFilter struct {
	Keyword    string
	Status     string
	Confidence string
	Page       uint32
	Size       uint32
}

type SeoObjectiveRepo interface {
	CreateObjective(ctx context.Context, objective *seo_strategy.SeoObjective) error
	GetObjectiveByID(ctx context.Context, id uint64) (*seo_strategy.SeoObjective, error)
	GetObjectiveByCode(ctx context.Context, code string) (*seo_strategy.SeoObjective, error)
	ListObjectives(ctx context.Context, filter SeoObjectiveFilter) ([]seo_strategy.SeoObjective, int64, error)

	UpdateObjectiveDefinition(ctx context.Context, objective *seo_strategy.SeoObjective, expectedVersion uint32) error
	UpdateObjectiveStatus(ctx context.Context, id uint64, status seo_strategy.ObjectiveStatus, expectedVersion uint32) error

	CreateEvidence(ctx context.Context, evidence *seo_strategy.SeoObjectiveEvidence) error
	ListEvidence(ctx context.Context, objectiveID uint64) ([]seo_strategy.SeoObjectiveEvidence, error)

	UpsertMetric(ctx context.Context, metric *seo_strategy.SeoObjectiveMetric) error
	ListMetrics(ctx context.Context, objectiveID uint64) ([]seo_strategy.SeoObjectiveMetric, error)
}
