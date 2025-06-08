package repo

import (
	"context"
	"social/internal/domain"
)

type MediaRepo interface {
	Create(ctx context.Context, media *domain.NewsFeedMedia) (*domain.NewsFeedMedia, error)
}
