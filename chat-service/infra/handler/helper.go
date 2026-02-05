package handler

import (
	"context"

	"chat/infra/client"
	"chat/infra/mapper"
	chatpb "pb/types/chat"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

// Helper functions to export handlers and clients

func NewUserClient() *client.UserClient {
	return client.NewUserClient()
}

func NewBdsproClient() *client.BdsproClient {
	return client.NewBdsproClient()
}

func NewConversationMapper(backgroundMapper *mapper.BackgroundMapper) *mapper.ConversationMapper {
	return mapper.NewConversationMapper(backgroundMapper)
}

func NewBackgroundMapper() *mapper.BackgroundMapper {
	return mapper.NewBackgroundMapper()
}

func RegisterGRPCServer(server *grpc.Server, chatHandler *chatHandler, backgroundHandler *BackgroundHandler) {
	chatpb.RegisterChatServiceServer(server, chatHandler)
	chatpb.RegisterBackgroundImageServiceServer(server, backgroundHandler)
}

func RegisterHTTPServer(ctx context.Context, mux *runtime.ServeMux, chatHandler *chatHandler, backgroundHandler *BackgroundHandler) error {
	err := chatpb.RegisterChatServiceHandlerServer(ctx, mux, chatHandler)
	if err != nil {
		return err
	}
	err = chatpb.RegisterBackgroundImageServiceHandlerServer(ctx, mux, backgroundHandler)
	if err != nil {
		return err
	}
	return nil
}
