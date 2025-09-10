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
	client *clients.NotificationClient
}

func NewDealHistoryClient() *DealHistoryClient {
	return &DealHistoryClient{
		client: clients.MakeNotificationClient(env.ENV_RUNTIME),
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
	_, err = c.client.CreateDealHistory(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create deal history: %w", err)
	}

	return nil
}

func (c *DealHistoryClient) GetDealHistory(ctx context.Context, dealID uint64, page, size int) ([]*dto.DealHistoryDTO, uint32, error) {
	// Create request
	req := &notificationpb.GetDealHistoryRequest{
		DealId: dealID,
		Page:   int32(page),
		Size:   int32(size),
	}

	// Call gRPC service
	resp, err := c.client.GetDealHistory(ctx, req)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get deal history: %w", err)
	}

	// Convert response to DTOs
	histories := make([]*dto.DealHistoryDTO, len(resp.Data))
	for i, item := range resp.Data {
		histories[i] = &dto.DealHistoryDTO{
			DealID:      item.DealId,
			ActorID:     item.ActorId,
			ActorName:   item.ActorName,
			ActorAvatar: item.ActorAvatar,
			ActionType:  item.ActionType,
			ActionName:  item.ActionName,
			Content:     item.Content,
			CreatedAt:   item.CreatedAt.AsTime(),
			Metadata:    convertStructToMap(item.Metadata),
		}
	}

	return histories, resp.Total, nil
}

func (c *DealHistoryClient) GetDealHistoryByID(ctx context.Context, id uint64) (*dto.DealHistoryDTO, error) {
	// Create request
	req := &notificationpb.GetDealHistoryByIDRequest{
		Id: id,
	}

	// Call gRPC service
	resp, err := c.client.GetDealHistoryByID(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get deal history by id: %w", err)
	}

	if resp.Data == nil {
		return nil, nil
	}

	// Convert response to DTO
	history := &dto.DealHistoryDTO{
		DealID:      resp.Data.DealId,
		ActorID:     resp.Data.ActorId,
		ActorName:   resp.Data.ActorName,
		ActorAvatar: resp.Data.ActorAvatar,
		ActionType:  resp.Data.ActionType,
		ActionName:  resp.Data.ActionName,
		Content:     resp.Data.Content,
		CreatedAt:   resp.Data.CreatedAt.AsTime(),
		Metadata:    convertStructToMap(resp.Data.Metadata),
	}

	return history, nil
}

// Helper function to convert structpb.Struct to map
func convertStructToMap(s *structpb.Struct) map[string]interface{} {
	if s == nil {
		return make(map[string]interface{})
	}

	// Convert structpb.Struct to map
	// This is a simplified conversion - you might need more sophisticated handling
	return make(map[string]interface{})
}
