package iusecase

import (
	"context"
	"organization/internal/dto"
)

type IChatClient interface {
	CreateConversation(ctx context.Context, req *dto.CreateConversationRequest) (*uint64, error)
	RemoveUserFromConversation(ctx context.Context, req *dto.RemoveUserRequest) error
	AddUserToConversation(ctx context.Context, req *dto.AddUserRequest) error
}
