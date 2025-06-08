package client

import (
	"context"
	"social/internal/domain"
)

type ServiceClient struct {
}

func (s *ServiceClient) GetPost(ctx context.Context, id uint64) (*domain.NewsFeed, error) {
	return nil, nil
}
