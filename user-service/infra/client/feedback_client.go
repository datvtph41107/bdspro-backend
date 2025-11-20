package client

import (
	clients "pb/clients"
	"user/env"
)

type FeedbackClient struct {
	*clients.FeedbackGrpcClient
}

func NewFeedbackClient() *FeedbackClient {
	return &FeedbackClient{
		FeedbackGrpcClient: clients.MakeFeedbackGrpcClient(env.ENV_RUNTIME),
	}
}
