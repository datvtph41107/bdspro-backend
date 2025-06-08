package repo

import (
	"context"
	"social/internal/domain"
)

type ReportRepo interface {
	Create(ctx context.Context, report *domain.Report) (*domain.Report, error)
	ExistByNewsFeedIdAndUserId(ctx context.Context, newsFeedId uint64, userId uint64) (bool, error)
}
