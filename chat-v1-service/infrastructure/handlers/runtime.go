package handlers

import (
	"context"

	"chat/infrastructure/client"
	"chat/infrastructure/mapper"
	"chat/internal/usecases"
	_utils "common/utils"
	chatpb "pb/types/chat"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
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
	userClient *client.UserClient,
	crmClient *client.CrmClient,
	bdsproClient *client.BdsproClient,
	conversationMapper *mapper.ConversationMapper,
	messageMapper *mapper.MessageMapper,
	backgroundMapper *mapper.BackgroundMapper,
	syncProvider *_utils.SyncUtil,
	messageRepository usecases.Repository,
	systemMessageUsecase *usecases.SystemMessageUsecase,
) *Delivery {
	chatHandler := NewChatHandler(conversationUsecases,
		messageUsecases,
		readReceptUsecases,
		participantUsecases,
		messageReactionUsecases,
		userClient,
		crmClient,
		bdsproClient,
		conversationMapper,
		messageMapper,
		backgroundMapper,
		syncProvider,
		messageRepository,
		systemMessageUsecase,
	)
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
	err = chatpb.RegisterBackgroundImageServiceHandlerServer(ctx, mux, d.externalHandler)
	if err != nil {
		return err
	}
	return nil
}

func (d *Delivery) RegisterGRPCServer(server *grpc.Server) {
	chatpb.RegisterChatServiceServer(server, d.chatHandler)
	chatpb.RegisterBackgroundImageServiceServer(server, d.externalHandler)
}
