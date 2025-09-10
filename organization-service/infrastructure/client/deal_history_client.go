package client

import (
	"context"
	"fmt"
	"organization/env"
	"organization/internal/dto"
	clients "pb/clients"
	notificationpb "pb/types/notification"

	"google.golang.org/protobuf/types/known/structpb"
)

// @bind: organization/internal/interface.IDealHistoryClient
type DealHistoryClient struct {
	*clients.NotificationClient
}

func NewDealHistoryClient() *DealHistoryClient {
	return &DealHistoryClient{
		NotificationClient: clients.MakeNotificationClient(env.ENV_RUNTIME),
	}
}

func (c *DealHistoryClient) CreateDealHistory(ctx context.Context, history *dto.DealHistoryDTO) error {
	// Convert metadata to structpb.Struct
	metadata, err := structpb.NewStruct(history.Metadata)
	if err != nil {
		return fmt.Errorf("failed to convert metadata: %w", err)
	}

	// Create request
	req := &notificationpb.CreateDealHistoryRequest{
		DealId:      history.DealID,
		ActorId:     history.ActorID,
		ActorName:   history.ActorName,
		ActorAvatar: history.ActorAvatar,
		ActionType:  string(history.ActionType),
		ActionName:  history.ActionName,
		Content:     history.Content,
		Metadata:    metadata,
	}

	// Call gRPC service
	_, err = c.Client.CreateDealHistory(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create deal history: %w", err)
	}

	return nil
}
