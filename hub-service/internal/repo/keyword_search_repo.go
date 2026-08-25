package repo

import (
	"context"
	"time"

	"hub/internal/domain"
	"hub/internal/dto"
)

type IKeywordSearchRepo interface {
	FlushCounts(ctx context.Context, logs []*domain.KeywordSearchLog) error
	CountFromTime(ctx context.Context, fromTime time.Time, source string) (uint64, error)
	GetTopKeywords(ctx context.Context, fromTime time.Time, source string, limit int) ([]dto.KeywordSearchItem, error)
}
