package repo

import (
	"context"
	"social/internal/domain"
)

type NewsFeedMediaRepo interface {
	UpdateMedia(ctx context.Context, newsFeedID uint64, medias []domain.NewsFeedMedia) error
}
