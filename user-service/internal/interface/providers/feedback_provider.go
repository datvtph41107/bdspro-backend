package providers

import (
	"context"
	models "user/internal/models"
)

type FeedbackProvider interface {
	GetUserRateStats(ctx context.Context, userId uint64) (*models.RateStats, error)
}
