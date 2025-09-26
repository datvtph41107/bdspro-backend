package providers

import (
	"context"
	user_model "user/internal/models/user"
)

type FeedbackProvider interface {
	GetUserRateStats(ctx context.Context, userId uint64) (*user_model.RateStats, error)
}
