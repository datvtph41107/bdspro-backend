package client

import (
	_dto "common/domain/dto"
	"context"
	clients "pb/clients"
	feedbackpb "pb/types/feedback"
	sharepb "pb/types/shared"
	"user/env"
	models "user/internal/models"
)

// @bind: user/internal/interface/providers.FeedbackProvider
type FeedbackClient struct {
	*clients.FeedbackGrpcClient
}

func NewFeedbackClient() *FeedbackClient {
	return &FeedbackClient{
		FeedbackGrpcClient: clients.MakeFeedbackGrpcClient(env.ENV_RUNTIME),
	}
}

func (c *FeedbackClient) RegisterService(ctx context.Context) error {
	_, err := c.FeedbackGrpcClient.RegisterEntity(ctx, &feedbackpb.FeedbackRegistry{
		RegistryName: "user",
		Port:         8201,
		Service:      "user",
	})

	return err
}

func (c *FeedbackClient) GetUserRateStats(ctx context.Context, userId uint64) (*models.RateStats, error) {
	// rateStats, err := c.FeedbackGrpcClient.GetRateStats(ctx, userId, 10)
	// if err != nil {
	// 	return nil, err
	// }
	return &models.RateStats{
		TotalReviews: 0,
		AverageScore: 0,
		Star1Count:   0,
		Star2Count:   0,
		Star3Count:   0,
		Star4Count:   0,
		Star5Count:   0,
	}, nil
}

func (c *FeedbackClient) GetUserRateStatsV3(ctx context.Context, userId uint64) (*_dto.RateStateV3DTO, error) {

	resp, err := c.FeedbackGrpcClient.GetRateStats(ctx, &sharepb.IdRequest{Id: userId})
	if err != nil {
		return &_dto.RateStateV3DTO{
			TotalReviews: 0,
			AverageScore: 0,
			Star1Count:   0,
			Star2Count:   0,
			Star3Count:   0,
			Star4Count:   0,
			Star5Count:   0,
		}, nil
	}

	return &_dto.RateStateV3DTO{
		TotalReviews: uint64(resp.TotalReviews),
		AverageScore: float32(resp.AverageScore),
		Star1Count:   uint64(resp.Star1Count),
		Star2Count:   uint64(resp.Star2Count),
		Star3Count:   uint64(resp.Star3Count),
		Star4Count:   uint64(resp.Star4Count),
		Star5Count:   uint64(resp.Star5Count),
	}, nil
}
