package handler

import (
	"chat/config"
	"chat/internal/constants"
	"chat/internal/dto"
	"chat/internal/enums"
	"chat/internal/usecases"
	"chat/utils"
	_errors "common/errors"
	_utils "common/utils"
	"context"
	chatpbv2 "pb/types/chat/v2"
	"strconv"

	"github.com/hyperledger/fabric/common/flogging"
)

type ConversationHandler struct {
	logger               *flogging.FabricLogger
	requiredParticipants int
	conversationUC       usecases.ConversationUsecases
}

func NewConversationHandler(
	conversationUsecase usecases.ConversationUsecases,
) *ConversationHandler {
	return &ConversationHandler{
		conversationUC:       conversationUsecase,
		logger:               flogging.MustGetLogger("chat_handler"),
		requiredParticipants: config.AppProperties.RequiredParticipants,
	}
}

func (s *ConversationHandler) CreateConversation(ctx context.Context, req *chatpbv2.CreateConversationRequest) (*chatpbv2.CreateConversationResponse, error) {
	userID := _utils.GetProfileIdWithContext(ctx)
	if userID == 0 {
		return nil, _errors.UnauthorizedException("unauthorized")
	}
	if len(req.Participants) < s.requiredParticipants {
		return nil, _errors.BadRequestException("not enough participants to create conversation")
	}

	reqDto := &dto.CreateConversationDTO{
		Name:           &req.Name,
		Type:           enums.ConversationType(req.Type),
		ParticipantIDs: req.Participants,
	}
	rsl, err := s.conversationUC.CreateConversation(ctx, reqDto)
	if err != nil {
		return nil, err
	}
	event := &chatpb.CreateRoom{
		Type: constants.CREATE_ROOM_TYPE,
		Data: &chatpb.CreateRoomData{
			RoomId:  strconv.Itoa(int(model.ID)),
			Members: utils.Int32SliceToStringSlice(req.Members),
		},
	}
	s.eventChannel <- event
	return &chatpb.CreateConversationResponse{
		ConversationId: int32(model.ID),
	}, nil
}
