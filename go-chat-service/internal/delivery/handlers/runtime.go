package handlers

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	chatpb "github.com/tuandq3/go-chat-service/proto/types/chat"
	externalpb "github.com/tuandq3/go-chat-service/proto/types/external"
	"google.golang.org/grpc"
)

type Delivery struct {
	chatHandler     *chatHandler
	externalHandler *externalHandler
}

func NewDelivery(
	conversationUsecases conversationUsecases,
	messageUsecases messageUsecases,
	readReceptUsecases readReceptUsecases,
	participantUsecases participantUsecases,
	messageReactionUsecases messageReactionUsecases,

	backgroundImageUsecases backgroundImageUsecases,
) *Delivery {
	chatHandler := NewChatHandler(conversationUsecases, messageUsecases, readReceptUsecases, participantUsecases, messageReactionUsecases)
	go chatHandler.Start()

	externalHandler := NewExternalHandler(backgroundImageUsecases)
	return &Delivery{
		chatHandler:     chatHandler,
		externalHandler: externalHandler,
	}
}

func (d *Delivery) RegisterHTTPServer(ctx context.Context, mux *runtime.ServeMux) error {
	err := chatpb.RegisterChatServiceHandlerServer(ctx, mux, d.chatHandler)
	if err != nil {
		return err
	}
	err = externalpb.RegisterBackgroundImageServiceHandlerServer(ctx, mux, d.externalHandler)
	if err != nil {
		return err
	}
	return nil
}

func (d *Delivery) RegisterGRPCServer(server *grpc.Server) {
	chatpb.RegisterChatServiceServer(server, d.chatHandler)
	externalpb.RegisterBackgroundImageServiceServer(server, d.externalHandler)
}
