package repo

import (
	"bdspro/internal/domain"
	"context"
)

type PropertyStatisticRepo interface {
	Save(ctx context.Context, statistic *domain.PropertyStatistic) error
}
