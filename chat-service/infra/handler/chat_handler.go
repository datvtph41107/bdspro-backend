package handler

import (
	_errors "common/errors"
	_utils "common/utils"
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	chatpb "pb/types/chat"
	sharepb "pb/types/shared"

	"chat/config"
	"chat/infra/client"
	"chat/infra/mapper"
	"chat/infra/redis"
	"chat/internal/constants"
	"chat/internal/domain"
	"chat/internal/usecases"

	"chat/utils"
)

type chatHandler struct {
	chatpb.UnimplementedChatServiceServer
	conversationUsecases    *usecases.ConversationUsecases
	messageUsecases         *usecases.MessageUsecases
	readMarkUsecases        *usecases.ReadMarkUsecases
	participantUsecases     *usecases.ParticipantUsecases
	messageReactionUsecases *usecases.MessageReactionUsecases
	eventChannel            chan any
	numberOfWorker          int
	maxParticipant          int
	redisClient             *redis.RedisClient
	userClient              *client.UserClient
	bdsproClient            *client.BdsproClient
	conversationMapper      *mapper.ConversationMapper
	backgroundMapper        *mapper.BackgroundMapper
	syncProvider            *_utils.SyncUtil
	workerCtx               context.Context
	cancelWorkers           context.CancelFunc
	workerWG                sync.WaitGroup
	startOnce               sync.Once
	closeOnce               sync.Once
}

func NewChatHandler(
	conversationUsecases *usecases.ConversationUsecases,
	messageUsecases *usecases.MessageUsecases,
	readMarkUsecases *usecases.ReadMarkUsecases,
	participantUsecases *usecases.ParticipantUsecases,
	messageReactionUsecases *usecases.MessageReactionUsecases,
	userClient *client.UserClient,
	bdsproClient *client.BdsproClient,
	conversationMapper *mapper.ConversationMapper,
	backgroundMapper *mapper.BackgroundMapper,
	syncProvider *_utils.SyncUtil,
) *chatHandler {
	workerCtx, cancelWorkers := context.WithCancel(context.Background())
	return &chatHandler{
		conversationUsecases:    conversationUsecases,
		messageUsecases:         messageUsecases,
		readMarkUsecases:        readMarkUsecases,
		participantUsecases:     participantUsecases,
		messageReactionUsecases: messageReactionUsecases,
		numberOfWorker:          config.AppProperties.Worker.Number,
		maxParticipant:          config.AppProperties.MaxParticipant,
		redisClient:             redis.NewRedisClient(),
		eventChannel:            make(chan any),
		userClient:              userClient,
		bdsproClient:            bdsproClient,
		conversationMapper:      conversationMapper,
		backgroundMapper:        backgroundMapper,
		syncProvider:            syncProvider,
		workerCtx:               workerCtx,
		cancelWorkers:           cancelWorkers,
	}
}

const TimeGap = 30 * time.Minute

// @Summary Tạo đoạn hội thoại
// @Description Tạo đoạn hội thoại với các thành viên được chỉ định
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param conversation body chatpb.CreateConversationRequest true "Đoạn hội thoại cần tạo"
// @Success 200 {object} chatpb.CreateConversationResponse
// @Router /conversations [post]
// func (s *chatHandler) CreateConversation(ctx context.Context, req *chatpb.CreateConversationRequest) (*chatpb.CreateConversationResponse, error) {
// 	if len(req.Members) > s.maxParticipant {
// 		return nil, errors.MaxParticipantsReached()
// 	}

// 	if len(req.Members) == 0 && req.ReceiverId == nil {
// 		return nil, errors.NotParticipant()
// 	}

// 	conversation := &domain.Conversation{
// 		Name:              &req.Name,
// 		Type:              domain.ConversationTypeEnum(req.Type),
// 		CreatedBy:         uint64(req.CreatedBy),
// 		IsBroadcast:       req.IsBroadcast,
// 		ForbidForward:     req.ForbidForward,
// 		Avatar:            req.Avatar,
// 		BackgroundImageID: req.BackgroundImageId,
// 		ReceiverId:        req.ReceiverId,
// 	}
// 	model, err := s.conversationUsecases.CreateConversation(ctx, conversation, utils.Int32SliceToUint64Slice(req.Members))
// 	if err != nil {
// 		return nil, err
// 	}
// 	event := &chatpb.CreateRoom{
// 		Type: constants.CREATE_ROOM_TYPE,
// 		Data: &chatpb.CreateRoomData{
// 			RoomId:  strconv.Itoa(int(model.ID)),
// 			Members: utils.Int32SliceToStringSlice(req.Members),
// 		},
// 	}
// 	s.eventChannel <- event
// 	return &chatpb.CreateConversationResponse{
// 		ConversationId: int32(model.ID),
// 	}, nil
// }

// @Summary Lấy danh sách đoạn hội thoại
// @Description Lấy danh sách đoạn hội thoại theo các tiêu chí được chỉ định
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param conversation query chatpb.GetConversationsRequest true "Tiêu chí lấy danh sách đoạn hội thoại"
// @Success 200 {object} chatpb.GetConversationsResponse
// @Router /conversations [get]
func (s *chatHandler) GetConversations(ctx context.Context, req *chatpb.GetConversationsRequest) (*chatpb.GetConversationsResponse, error) {
	var limit, offset int64
	if req.Limit != nil {
		limit = int64(*req.Limit)
	} else {
		limit = 10
	}

	if req.Offset != nil {
		offset = int64(*req.Offset)
	} else {
		offset = 0
	}

	conversations, err := s.conversationUsecases.GetConversations(ctx, limit, offset, req.Type, req.Keyword)
	if err != nil {
		return nil, err
	}

	pbConversations := s.conversationMapper.ConversationsToPb(conversations)
	s.userClient.MapUserProfile(ctx, pbConversations)

	return &chatpb.GetConversationsResponse{
		Conversations: pbConversations,
	}, nil
}

// @Summary Gửi tin nhắn
// @Description Gửi tin nhắn đến đoạn hội thoại được chỉ định
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param message body chatpb.SendMessageRequest true "Tin nhắn cần gửi"
// @Success 200 {object} chatpb.SendMessageResponse
// @Router /messages [post]
func (s *chatHandler) SendMessage(ctx context.Context, req *chatpb.SendMessageRequest) (*chatpb.SendMessageResponse, error) {
	message, err := s.messageUsecases.SendMessage(ctx, &domain.Message{
		ConversationID: uint64(req.ConversationId),
		SenderID:       uint64(req.SenderId),
		Content:        req.Content,
		ContentType:    domain.ContentTypeEnum(req.ContentType),
		FileURL:        req.FileUrl,
		ReplyId:        req.ReplyId,
		ExtraId:        req.ExtraId,
	})
	if err != nil {
		return nil, err
	}

	return s.ResponseMessage(ctx, message, req.CmId)
}

// @Summary Gửi tin nhắn đến người nhận
// @Description Gửi tin nhắn đến người nhận
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param message body chatpb.SendToReceiverRequest true "Tin nhắn cần gửi"
// @Success 200 {object} chatpb.SendToReceiverResponse
// @Router /messages/to-receiver [post]
func (s *chatHandler) SendToReceiver(ctx context.Context, req *chatpb.SendMessageRequest) (*chatpb.SendMessageResponse, error) {
	if req.ReceiverId == nil {
		return nil, _errors.ReturnError(_errors.RequestValidationFailed, _errors.WithPublicMessage("receiverId is required"))
	}

	message, err := s.messageUsecases.SendToReceiver(
		ctx,
		&domain.Message{
			ConversationID: uint64(req.ConversationId),
			SenderID:       uint64(req.SenderId),
			Content:        req.Content,
			ContentType:    domain.ContentTypeEnum(req.ContentType),
			FileURL:        req.FileUrl,
			ReplyId:        req.ReplyId,
			ExtraId:        req.ExtraId,
		},
		req.ReceiverId,
	)
	if err != nil {
		return nil, err
	}

	return s.ResponseMessage(ctx, message, req.CmId)
}

func (s *chatHandler) ResponseMessage(ctx context.Context, message *domain.Message, cmId int32) (*chatpb.SendMessageResponse, error) {
	createdAtStr := ""
	if message.CreatedAt != nil {
		createdAtStr = message.CreatedAt.Local().Format(time.RFC3339)
	}
	event := &chatpb.SendMessage{
		Type: constants.MESSAGE_TYPE,
		Data: &chatpb.Message{
			MessageId:   int32(message.ID),
			SenderId:    int32(message.SenderID),
			Content:     message.Content,
			ContentType: string(message.ContentType),
			FileUrl:     message.FileURL,
			RoomId:      strconv.FormatUint(message.ConversationID, 10),
			CreatedAt:   createdAtStr,
			IsRecall:    message.Recall,
			ExtraId:     message.ExtraId,
			CmId:        cmId,
			IndexKey:    message.IndexKey,
		},
	}
	if message.ReplyMessage != nil {
		replyCreatedAtStr := ""
		if message.ReplyMessage.CreatedAt != nil {
			replyCreatedAtStr = message.ReplyMessage.CreatedAt.Local().Format(time.RFC3339)
		}
		event.Data.ReplyMessage = &chatpb.Message{
			MessageId:   int32(message.ReplyMessage.ID),
			SenderId:    int32(message.ReplyMessage.SenderID),
			Content:     message.ReplyMessage.Content,
			ContentType: string(message.ReplyMessage.ContentType),
			FileUrl:     message.ReplyMessage.FileURL,
			RoomId:      strconv.FormatUint(message.ReplyMessage.ConversationID, 10),
			CreatedAt:   replyCreatedAtStr,
			IsRecall:    message.ReplyMessage.Recall,
		}
	}

	user, _ := s.userClient.GetUserById(ctx, message.SenderID)
	if user != nil {
		event.Data.FullName = &user.FullName
		event.Data.Avatar = &user.Avatar
	}

	s.eventChannel <- event

	updateRoomEvent := &chatpb.UpdateRoomMessage{
		Type: constants.UPDATE_ROOM_TYPE,
		Data: &chatpb.UpdateRoomMessageData{
			RoomId: strconv.FormatUint(message.ConversationID, 10),
		},
	}
	s.eventChannel <- updateRoomEvent
	if message.CreatedAt != nil {
		createdAtStr = message.CreatedAt.Local().Format(time.RFC3339)
	}
	return &chatpb.SendMessageResponse{
		MessageId: int32(message.ID),
		CreatedAt: createdAtStr,
		CmId:      cmId,
		IndexKey:  message.IndexKey,
	}, nil
}

func (s *chatHandler) GetMessagesByIndex(ctx context.Context, req *chatpb.GetMessagesByIndexRequest) (*chatpb.GetMessagesResponse, error) {
	messages, err := s.messageUsecases.GetMessagesByIndex(
		ctx,
		req.ConversationId,
		req.FromIdx,
		req.ToIdx,
	)
	if err != nil {
		return nil, err
	}

	return s.buildGetMessagesResponse(ctx, messages), nil
}

func (s *chatHandler) GetStateTyping(ctx context.Context, req *chatpb.StateTypingRequest) (*chatpb.StateTypingResponse, error) {
	userId := utils.GetCurrentUserID(ctx)
	if userId == 0 {
		return nil, _errors.ReturnError(_errors.AuthenticationRequired, _errors.WithPublicMessage("missing user_id"), _errors.WithLegacyCode(401))
	}

	err := s.conversationUsecases.ValidateConversationAndCurrentUser(ctx, req.ConversationId)
	if err != nil {
		return nil, err
	}

	count, err := s.redisClient.UpdateTyping(ctx, req.ConversationId, userId, req.IsTyping)
	if err != nil {
		return nil, fmt.Errorf("update typing state: %w", err)
	}

	event := &chatpb.TypingIndicator{
		Type: constants.TYPING_MESSAGE_TYPE,
		Data: &chatpb.TypingIndicatorData{
			RoomId:   strconv.FormatUint(req.ConversationId, 10),
			UserId:   strconv.FormatUint(userId, 10),
			IsTyping: req.IsTyping,
			Count:    count,
		},
	}
	s.eventChannel <- event

	return &chatpb.StateTypingResponse{Success: true, Count: count}, nil
}

// @Summary Lấy danh sách tin nhắn
// @Description Lấy danh sách tin nhắn theo các tiêu chí được chỉ định
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param message query chatpb.GetMessagesRequest true "Tiêu chí lấy danh sách tin nhắn"
// @Success 200 {object} chatpb.GetMessagesResponse
// @Router /messages [get]
func (s *chatHandler) GetMessages(ctx context.Context, req *chatpb.GetMessagesRequest) (*chatpb.GetMessagesResponse, error) {
	limit := int64(10)
	offset := int64(0)
	sort := "desc"

	if req.Size > 0 {
		limit = int64(req.Size)
	}
	if req.Page > 0 {
		offset = int64(req.Page) * limit
	}
	if strings.EqualFold(req.Sort, "asc") {
		sort = "asc"
	}

	var fromDate, toDate *time.Time
	if req.FromDate != nil {
		fromDate = _utils.ParseStringToTime(*req.FromDate)
	}
	if req.ToDate != nil {
		toDate = _utils.ParseStringToTime(*req.ToDate)
	}

	messages, err := s.messageUsecases.GetMessages(
		ctx,
		limit,
		offset,
		uint64(req.ConversationId),
		sort,
		fromDate,
		toDate,
	)
	if err != nil {
		return nil, err
	}

	return s.buildGetMessagesResponse(ctx, messages), nil
}

func (s *chatHandler) buildGetMessagesResponse(
	ctx context.Context,
	messages []*domain.Message,
) *chatpb.GetMessagesResponse {

	var pbMessages []*chatpb.Message
	var lastTime *time.Time

	for _, msg := range messages {
		var msgCreatedAt time.Time
		if msg.CreatedAt != nil {
			msgCreatedAt = *msg.CreatedAt
		}

		if lastTime == nil || (msg.CreatedAt != nil && msgCreatedAt.Sub(*lastTime) >= TimeGap) {
			createdAtStr := ""
			if msg.CreatedAt != nil {
				createdAtStr = msgCreatedAt.Format(time.RFC3339)
			}
			pbMessages = append(pbMessages, &chatpb.Message{
				ContentType: "timeline",
				Content:     msgCreatedAt.Format("15:04, 02/01"),
				CreatedAt:   createdAtStr,
			})
		}

		if msg.Recall {
			createdAtStr := ""
			if msg.CreatedAt != nil {
				createdAtStr = msgCreatedAt.Format(time.RFC3339)
			}
			pbMessages = append(pbMessages, &chatpb.Message{
				MessageId: int32(msg.ID),
				SenderId:  int32(msg.SenderID),
				IsRecall:  true,
				CreatedAt: createdAtStr,
			})
			if msg.CreatedAt != nil {
				lastTime = msg.CreatedAt
			}
			continue
		}

		createdAtStr := ""
		if msg.CreatedAt != nil {
			createdAtStr = msgCreatedAt.Format(time.RFC3339)
		}

		pbMsg := &chatpb.Message{
			MessageId:   int32(msg.ID),
			SenderId:    int32(msg.SenderID),
			Content:     msg.Content,
			ContentType: string(msg.ContentType),
			FileUrl:     msg.FileURL,
			RoomId:      strconv.FormatUint(msg.ConversationID, 10),
			CreatedAt:   createdAtStr,
			IsRecall:    msg.Recall,
			ExtraId:     msg.ExtraId,
			IndexKey:    msg.IndexKey,
		}

		if msg.ReplyMessage != nil {
			replyCreatedAtStr := ""
			if msg.ReplyMessage.CreatedAt != nil {
				replyCreatedAtStr = msg.ReplyMessage.CreatedAt.Format(time.RFC3339)
			}
			pbMsg.ReplyMessage = &chatpb.Message{
				MessageId:   int32(msg.ReplyMessage.ID),
				SenderId:    int32(msg.ReplyMessage.SenderID),
				Content:     msg.ReplyMessage.Content,
				ContentType: string(msg.ReplyMessage.ContentType),
				FileUrl:     msg.ReplyMessage.FileURL,
				RoomId:      strconv.FormatUint(msg.ReplyMessage.ConversationID, 10),
				CreatedAt:   replyCreatedAtStr,
			}
		}

		if len(msg.Reactions) > 0 {
			for _, r := range msg.Reactions {
				reactionCreatedAtStr := ""
				if r.CreatedAt != nil {
					reactionCreatedAtStr = r.CreatedAt.Format(time.RFC3339)
				}
				pbMsg.Reactions = append(pbMsg.Reactions, &chatpb.MessageReaction{
					Id:        r.ID,
					UserId:    int32(r.UserID),
					MessageId: int32(r.MessageID),
					Reaction:  r.Reaction,
					CreatedAt: reactionCreatedAtStr,
				})
			}
		}

		pbMessages = append(pbMessages, pbMsg)
		if msg.CreatedAt != nil {
			lastTime = msg.CreatedAt
		}
	}

	// 6. Enrich user + product
	s.userClient.MapAvatarAndNameToMessagePbList(ctx, pbMessages)
	s.bdsproClient.MapProductPbByIds(ctx, pbMessages)

	return &chatpb.GetMessagesResponse{
		Messages: pbMessages,
	}
}

// @Summary Chuyển tiếp tin nhắn
// @Description Chuyển tiếp tin nhắn đến đoạn hội thoại được chỉ định
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param message body chatpb.ForwardMessageRequest true "Tin nhắn cần chuyển tiếp"
// @Success 200 {object} chatpb.ForwardMessageResponse
// @Router /messages/forward [post]
func (s *chatHandler) ForwardMessage(ctx context.Context, req *chatpb.ForwardMessageRequest) (*chatpb.ForwardMessageResponse, error) {
	messages, err := s.messageUsecases.ForwardMessage(ctx, uint64(req.MessageId), req.ConversationIds)
	if err != nil {
		return nil, err
	}

	pbMessages := make([]*chatpb.Message, 0)
	messageIds := make([]uint64, 0)
	for _, message := range messages {
		createdAtStr := ""
		if message.CreatedAt != nil {
			createdAtStr = message.CreatedAt.Local().Format(time.RFC3339)
		}
		messagePb := &chatpb.Message{
			MessageId:   int32(message.ID),
			SenderId:    int32(message.SenderID),
			Content:     message.Content,
			ContentType: string(message.ContentType),
			FileUrl:     message.FileURL,
			RoomId:      strconv.FormatUint(message.ConversationID, 10),
			CreatedAt:   createdAtStr,
		}
		pbMessages = append(pbMessages, messagePb)
		messageIds = append(messageIds, message.ID)

		// user, _ := s.userClient.GetUserById(ctx, message.SenderID)
		// if user != nil {
		// 	event.Data.FullName = &user.FullName
		// 	event.Data.Avatar = &user.Avatar
		// }
	}
	s.userClient.MapAvatarAndNameToMessagePbList(ctx, pbMessages)
	for _, message := range pbMessages {
		event := &chatpb.SendMessage{
			Type: constants.MESSAGE_TYPE,
			Data: message,
		}
		s.eventChannel <- event
	}

	// s.eventChannel <- event
	return &chatpb.ForwardMessageResponse{
		MessageIds: messageIds,
	}, nil
}

// @Summary Đánh dấu tin nhắn đã đọc
// @Description Đánh dấu tin nhắn đã đọc trong đoạn hội thoại được chỉ định
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param message body chatpb.MarkAsReadRequest true "Đoạn hội thoại và tin nhắn cần đánh dấu"
// @Success 200 {object} chatpb.MarkAsReadResponse
// @Router /messages/read [post]
func (s *chatHandler) MarkAsRead(ctx context.Context, req *chatpb.MarkAsReadRequest) (*chatpb.MarkAsReadResponse, error) {
	pos, updated, err := s.readMarkUsecases.MarkAsRead(ctx,
		uint64(req.ConversationId),
		uint64(req.LastMessageId))
	if err != nil {
		return nil, err
	}

	if !updated {
		return &chatpb.MarkAsReadResponse{Success: true}, nil
	}

	userId := utils.GetCurrentUserID(ctx)
	userIdStr := strconv.FormatUint(userId, 10)

	user, _ := s.userClient.GetUserById(ctx, userId)

	event := &chatpb.ReadMessage{
		Type: constants.READ_TYPE,
		Data: &chatpb.ReadMessageData{
			RoomId:     strconv.Itoa(int(req.ConversationId)),
			UserId:     userIdStr,
			MessageIds: []string{strconv.FormatUint(pos.LastMessageID, 10)},
		},
	}

	if user != nil {
		event.Data.FullName = &user.FullName
		event.Data.Avatar = &user.Avatar
	}

	s.eventChannel <- event

	return &chatpb.MarkAsReadResponse{Success: true}, nil
}

// @Summary Đánh dấu tin nhắn đã đọc
// @Description Đánh dấu tin nhắn đã đọc trong đoạn hội thoại được chỉ định
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param message body chatpb.PinMessageRequest true "Đoạn hội thoại và tin nhắn cần đánh dấu"
// @Success 200 {object} chatpb.PinMessageResponse
// @Router /messages/pin [post]
func (s *chatHandler) PinMessage(ctx context.Context, req *chatpb.PinMessageRequest) (*chatpb.PinMessageResponse, error) {
	_, err := s.messageUsecases.PinMessage(ctx, uint64(req.MessageId), uint64(req.ConversationId), req.DoPin)
	if err != nil {
		return nil, err
	}

	event := &chatpb.PinMessage{
		Type: constants.PIN_MESSAGE_TYPE,
		Data: &chatpb.Message{
			RoomId:    strconv.Itoa(int(req.ConversationId)),
			MessageId: req.MessageId,
		},
	}

	user, _ := s.userClient.GetUserById(ctx, utils.GetCurrentUserID(ctx))
	if user != nil {
		event.Data.FullName = &user.FullName
		event.Data.Avatar = &user.Avatar
	}

	s.eventChannel <- event

	return &chatpb.PinMessageResponse{
		Success: true,
	}, nil
}

// @Summary Gỡ tin nhắn
// @Description Gỡ tin nhắn đã được gửi trong đoạn hội thoại được chỉ định
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param message body chatpb.RecallMessageRequest true "Đoạn hội thoại và tin nhắn cần gỡ"
// @Success 200 {object} chatpb.RecallMessageResponse
// @Router /messages/recall [post]
func (s *chatHandler) RecallMessage(ctx context.Context, req *chatpb.RecallMessageRequest) (*chatpb.RecallMessageResponse, error) {
	message, err := s.messageUsecases.RecallMessage(ctx, uint64(req.MessageId))
	if err != nil {
		return nil, err
	}
	event := &chatpb.RecallMessage{
		Type: constants.RECALL_MESSAGE_TYPE,
		Data: &chatpb.RecallMessageData{
			RoomId:    strconv.Itoa(int(message.ConversationID)),
			MessageId: strconv.Itoa(int(req.MessageId)),
		},
	}
	s.eventChannel <- event

	return &chatpb.RecallMessageResponse{
		Success: true,
	}, nil
}

// @Summary Xóa thành viên
// @Description Xóa thành viên khỏi đoạn hội thoại được chỉ định
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param message body chatpb.RemoveUserRequest true "Đoạn hội thoại và thành viên cần xóa"
// @Success 200 {object} chatpb.RemoveUserResponse
// @Router /conversation/remove-member [post]
func (s *chatHandler) RemoveUser(ctx context.Context, req *chatpb.RemoveUserRequest) (*chatpb.RemoveUserResponse, error) {
	err := s.participantUsecases.RemoveUser(ctx, uint64(req.ConversationId), uint64(req.TargetUserId))
	if err != nil {
		return nil, err
	}
	event := &chatpb.RemoveMember{
		Type: constants.REMOVE_MEMBER_TYPE,
		Data: &chatpb.RemoveMemberData{
			RoomId: strconv.Itoa(int(req.ConversationId)),
			UserId: strconv.Itoa(int(req.TargetUserId)),
		},
	}
	s.eventChannel <- event
	return &chatpb.RemoveUserResponse{
		Success: true,
	}, nil
}

// @Summary Chỉnh sửa đoạn hội thoại
// @Description Chỉnh sửa đoạn hội thoại được chỉ định
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param message body chatpb.EditChannelRequest true "Đoạn hội thoại cần chỉnh sửa"
// @Success 200 {object} chatpb.EditChannelResponse
// @Param conversationId path int true "ID đoạn hội thoại"
// @Router /conversation/{conversationId}/edit [put]
func (s *chatHandler) EditChannel(ctx context.Context, req *chatpb.EditChannelRequest) (*chatpb.EditChannelResponse, error) {
	err := s.conversationUsecases.EditChannel(ctx, uint64(req.ConversationId), req.IsBroadcast, req.ForbidForward, req.Name, req.BackgroundImageId)
	if err != nil {
		return nil, err
	}
	return &chatpb.EditChannelResponse{
		Success: true,
	}, nil
}

// @Summary Tạo đoạn hội thoại tổ chức
// @Description Tạo đoạn hội thoại tổ chức
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param message body chatpb.CreateOrgChannelRequest true "Đoạn hội thoại tổ chức cần tạo"
// @Success 200 {object} chatpb.CreateOrgChannelResponse
// @Router /conversation/create-org [post]
func (s *chatHandler) CreateOrgChannel(ctx context.Context, req *chatpb.CreateOrgChannelRequest) (*chatpb.CreateOrgChannelResponse, error) {
	channel, err := s.conversationUsecases.CreateOrgChannel(ctx, &domain.Conversation{})
	if err != nil {
		return nil, err
	}

	return &chatpb.CreateOrgChannelResponse{
		ConversationId: int32(channel.ID),
	}, nil
}

// @Summary Ngưng thông báo
// @Description Ngưng thông báo trong đoạn hội thoại được chỉ định
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param message body chatpb.MuteConversationRequest true "Đoạn hội thoại và thông báo cần ngưng"
// @Success 200 {object} chatpb.MuteConversationResponse
// @Router /conversation/mute [post]
func (s *chatHandler) MuteConversation(ctx context.Context, req *chatpb.MuteConversationRequest) (*chatpb.MuteConversationResponse, error) {
	err := s.participantUsecases.MuteConversation(ctx, uint64(req.ConversationId), req.Mute)
	if err != nil {
		return nil, err
	}

	return &chatpb.MuteConversationResponse{
		Success: true,
	}, nil
}

// @Summary Tìm kiếm
// @Description Tìm kiếm tin nhắn và đoạn hội thoại theo các tiêu chí được chỉ định
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param message query chatpb.SearchRequest true "Tiêu chí tìm kiếm"
// @Success 200 {object} chatpb.SearchResponse
// @Router /search [get]
func (s *chatHandler) Search(ctx context.Context, req *chatpb.SearchRequest) (*chatpb.SearchResponse, error) {
	var searchType, limit, offset int
	if req.Type != nil {
		searchType = int(*req.Type)
	}
	if req.Limit != nil {
		limit = int(*req.Limit)
	}
	if req.Offset != nil {
		offset = int(*req.Offset)
	}

	conversations, messages, err := s.conversationUsecases.Search(ctx, req.Keyword, searchType, int64(limit), int64(offset))
	if err != nil {
		return nil, err
	}

	var pbMessages []*chatpb.Message
	for _, msg := range messages {
		pbMessages = append(pbMessages, &chatpb.Message{
			MessageId:   int32(msg.ID),
			SenderId:    int32(msg.SenderID),
			Content:     msg.Content,
			ContentType: string(msg.ContentType),
			FileUrl:     msg.FileURL,
			RoomId:      strconv.FormatUint(msg.ConversationID, 10),
			CreatedAt:   msg.CreatedAt.Local().Format(time.RFC3339),
		})
	}

	var pbConversations []*chatpb.Conversation
	for _, conv := range conversations {
		pbConversations = append(pbConversations, &chatpb.Conversation{
			ConversationId: int32(conv.ID),
			Name:           *conv.Name,
			Type:           chatpb.ConversationTypeEnum(conv.Type),
		})
	}

	return &chatpb.SearchResponse{
		Messages:      pbMessages,
		Conversations: pbConversations,
	}, nil
}

// @Summary Lấy đoạn hội thoại theo ID
// @Description Lấy đoạn hội thoại theo ID
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param message query chatpb.GetConversationByIdRequest true "ID đoạn hội thoại cần lấy"
// @Success 200 {object} chatpb.GetConversationByIdResponse
// @Param id path int true "ID đoạn hội thoại"
// @Router /conversation/{id} [get]
func (s *chatHandler) GetConversationById(ctx context.Context, req *chatpb.GetConversationByIdRequest) (*chatpb.GetConversationByIdResponse, error) {
	req.IncludeMembers = true
	conversation, err := s.conversationUsecases.GetConversationById(ctx, uint64(req.ConversationId), req.IncludeMembers)
	if err != nil {
		return nil, err
	}

	result := s.conversationMapper.ConversationToPbGetConversationById(conversation, req.IncludeMembers)

	s.userClient.MapAvatarAndNameToPb(ctx, result)

	return result, nil
}

// @Summary Lấy danh sách thành viên
// @Description Lấy danh sách thành viên trong đoạn hội thoại được chỉ định
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param message query chatpb.GetListParticipantRequest true "Đoạn hội thoại cần lấy danh sách thành viên"
// @Success 200 {object} chatpb.GetListParticipantResponse
// @Router /conversation/members [get]
func (s *chatHandler) GetListParticipant(ctx context.Context, req *chatpb.GetListParticipantRequest) (*chatpb.GetListParticipantResponse, error) {
	participants, total, err := s.participantUsecases.GetListParticipant(ctx, req.ConversationId)
	if err != nil {
		return nil, err
	}

	// Lấy danh sách userId
	userIds := make([]uint64, 0, len(participants))
	for _, participant := range participants {
		userIds = append(userIds, participant.UserID)
	}

	// Gọi user-service để lấy thông tin user
	userMap := make(map[uint64]*sharepb.ProfileItem)
	if len(userIds) > 0 {
		users, err := s.userClient.Client.GetProfileByIds(ctx, userIds)
		if err == nil && users != nil {
			for _, user := range users.Profiles {
				userMap[user.Id] = user
			}
		}
	}

	// Map data
	var pbParticipants []*chatpb.Participant
	for _, participant := range participants {
		pbParticipant := &chatpb.Participant{
			UserId:    participant.UserID,
			Role:      chatpb.ParticipantRoleEnum(participant.Role),
			JoinedAt:  participant.JoinedAt.Format(time.RFC3339),
			MuteNotif: participant.MuteNotif,
		}

		// Thêm thông tin user nếu có
		if user, ok := userMap[participant.UserID]; ok {
			pbParticipant.Profile = user
		}

		pbParticipants = append(pbParticipants, pbParticipant)
	}

	return &chatpb.GetListParticipantResponse{
		Data:  pbParticipants,
		Total: int32(total),
	}, nil
}

// @Summary Lấy danh sách đọc tin nhắn
// @Description Lấy danh sách đọc tin nhắn trong đoạn hội thoại được chỉ định
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param message query chatpb.GetReadReceiptRequest true "Đoạn hội thoại và tin nhắn cần lấy danh sách đọc"
// @Success 200 {object} chatpb.GetReadReceiptResponse
// @Router /messages/read-receipt [get]
func (s *chatHandler) GetReadReceipt(ctx context.Context, req *chatpb.GetReadReceiptRequest) (*chatpb.GetReadReceiptResponse, error) {
	messages, total, err := s.readMarkUsecases.GetReadReceipt(ctx, req.Page, req.Size, req.MessageId)
	if err != nil {
		return nil, err
	}
	var pbReadReceipts []*chatpb.ReadReceipt
	for _, readReceipt := range messages {
		pbReadReceipts = append(pbReadReceipts, &chatpb.ReadReceipt{
			UserId: readReceipt.UserID,
			ReadAt: readReceipt.ReadAt.Format(time.RFC3339),
		})
	}
	return &chatpb.GetReadReceiptResponse{
		ReadReceipts: pbReadReceipts,
		Total:        total,
	}, nil
}

// @Summary Lấy danh sách tin nhắn được đánh dấu
// @Description Lấy danh sách tin nhắn được đánh dấu trong đoạn hội thoại được chỉ định
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param message query chatpb.GetMessagesRequest true "Đoạn hội thoại cần lấy danh sách tin nhắn được đánh dấu"
// @Success 200 {object} chatpb.GetMessagesResponse
// @Router /messages/pinned [get]
func (s *chatHandler) GetPinnedMessages(ctx context.Context, req *chatpb.GetMessagesRequest) (*chatpb.GetMessagesResponse, error) {
	messages, err := s.messageUsecases.GetPinnedMessages(ctx, uint64(req.ConversationId))
	if err != nil {
		return nil, err
	}

	var pbMessages []*chatpb.Message
	for _, msg := range messages {
		pbMessages = append(pbMessages, &chatpb.Message{
			MessageId:   int32(msg.ID),
			SenderId:    int32(msg.SenderID),
			Content:     msg.Content,
			ContentType: string(msg.ContentType),
			FileUrl:     msg.FileURL,
			RoomId:      strconv.FormatUint(msg.ConversationID, 10),
			CreatedAt:   msg.CreatedAt.Local().Format(time.RFC3339),
		})
	}

	s.userClient.MapAvatarAndNameToMessagePbList(ctx, pbMessages)

	return &chatpb.GetMessagesResponse{
		Messages: pbMessages,
		Total:    int64(len(messages)),
	}, nil
}

// @Summary Thêm thành viên
// @Description Thêm một hoặc nhiều thành viên vào đoạn hội thoại được chỉ định
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param message body chatpb.AddUserRequest true "Đoạn hội thoại và danh sách thành viên cần thêm"
// @Success 200 {object} chatpb.AddUserResponse
// @Router /conversation/add-user [post]
func (s *chatHandler) AddUsers(ctx context.Context, req *chatpb.AddUserRequest) (*chatpb.AddUserResponse, error) {
	// Kiểm tra danh sách userIds
	if len(req.UserIds) == 0 {
		return &chatpb.AddUserResponse{
			Success:    false,
			TotalAdded: 0,
		}, nil
	}

	// Thêm users vào conversation
	totalAdded, err := s.participantUsecases.AddUsers(ctx, req.ConversationId, req.UserIds)
	if err != nil {
		return nil, err
	}

	// Lấy thông tin users để gửi event join room
	users, _ := s.userClient.Client.GetProfileByIds(ctx, req.UserIds)
	userMap := make(map[uint64]*sharepb.ProfileItem)
	if users != nil {
		for _, user := range users.Profiles {
			userMap[user.Id] = user
		}
	}

	// Gửi event join room cho từng user
	roomId := strconv.FormatUint(req.ConversationId, 10)
	for _, userId := range req.UserIds {
		event := &chatpb.JoinRoom{
			Type: constants.JOIN_ROOM_TYPE,
			Data: &chatpb.JoinRoomData{
				RoomId: roomId,
				UserId: strconv.FormatUint(userId, 10),
			},
		}
		if user, ok := userMap[userId]; ok {
			event.Data.FullName = &user.FullName
			event.Data.Avatar = &user.Avatar
		}
		s.eventChannel <- event
	}

	return &chatpb.AddUserResponse{
		Success:    totalAdded > 0,
		TotalAdded: int32(totalAdded),
	}, nil
}

// @Summary Rời đoạn hội thoại
// @Description Rời đoạn hội thoại được chỉ định
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param message body chatpb.LeaveRoomRequest true "Đoạn hội thoại cần rời"
func (s *chatHandler) LeaveRoom(ctx context.Context, req *chatpb.LeaveRoomRequest) (*chatpb.LeaveRoomResponse, error) {
	err := s.conversationUsecases.LeaveRoom(ctx, uint64(req.ConversationId))
	if err != nil {
		return nil, err
	}
	event := &chatpb.LeaveRoom{
		Type: constants.LEAVE_ROOM_TYPE,
		Data: &chatpb.LeaveRoomData{
			RoomId: strconv.Itoa(int(req.ConversationId)),
		},
	}
	s.eventChannel <- event
	return &chatpb.LeaveRoomResponse{
		Success: true,
	}, nil
}

// @Summary Chỉnh sửa tin nhắn
// @Description Chỉnh sửa tin nhắn đã được gửi trong đoạn hội thoại được chỉ định
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param message body chatpb.EditMessageRequest true "Đoạn hội thoại và tin nhắn cần chỉnh sửa"
// @Success 200 {object} chatpb.EditMessageResponse
// @Router /messages/edit [post]
func (s *chatHandler) EditMessage(ctx context.Context, req *chatpb.EditMessageRequest) (*chatpb.EditMessageResponse, error) {
	message, err := s.messageUsecases.EditMessage(ctx, uint64(req.MessageId), req.Content)
	if err != nil {
		return nil, err
	}
	event := &chatpb.EditMessage{
		Type: constants.EDIT_MESSAGE_TYPE,
		Data: &chatpb.EditMessageData{
			RoomId:     strconv.Itoa(int(message.ConversationID)),
			MessageId:  strconv.Itoa(int(req.MessageId)),
			NewMessage: message.Content,
		},
	}
	s.eventChannel <- event
	return &chatpb.EditMessageResponse{
		Success: true,
	}, nil
}

// @Summary Xóa vi phạm
// @Description Xóa vi phạm tin nhắn đã được gửi trong đoạn hội thoại được chỉ định
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param message body chatpb.DeleteViolationRequest true "Đoạn hội thoại và tin nhắn cần xóa vi phạm"
// @Success 200 {object} chatpb.DeleteViolationResponse
// @Router /messages/delete-violation [post]
func (s *chatHandler) DeleteViolation(ctx context.Context, req *chatpb.DeleteViolationRequest) (*chatpb.DeleteViolationResponse, error) {
	err := s.messageUsecases.DeleteViolation(ctx, uint64(req.MessageId))
	if err != nil {
		return nil, err
	}
	event := &chatpb.DeleteViolationMessage{
		Type: constants.DELETE_VIOLATION_MESSAGE_TYPE,
		Data: &chatpb.DeleteViolationMessageData{
			MessageId: strconv.Itoa(int(req.MessageId)),
		},
	}
	s.eventChannel <- event
	return &chatpb.DeleteViolationResponse{
		Success: true,
	}, nil
}

// @Summary Xóa đoạn hội thoại
// @Description Xóa đoạn hội thoại được chỉ định
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param message body chatpb.DeleteRoomRequest true "Đoạn hội thoại cần xóa"
// @Success 200 {object} chatpb.DeleteRoomResponse
// @Router /conversation/delete [post]
func (s *chatHandler) DeleteRoom(ctx context.Context, req *chatpb.DeleteRoomRequest) (*chatpb.DeleteRoomResponse, error) {
	err := s.conversationUsecases.DeleteRoom(ctx, uint64(req.ConversationId))
	if err != nil {
		return nil, err
	}
	event := &chatpb.DeleteRoom{
		Type: constants.DELETE_ROOM_TYPE,
		Data: &chatpb.DeleteRoomData{
			RoomId: strconv.Itoa(int(req.ConversationId)),
		},
	}
	s.eventChannel <- event
	return &chatpb.DeleteRoomResponse{
		Success: true,
	}, nil
}

// @Summary Lấy cài đặt đoạn hội thoại
// @Description Lấy cài đặt đoạn hội thoại được chỉ định
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID đoạn hội thoại"
// @Success 200 {object} chatpb.GetSettingResponse
// @Router /conversation/{id}/setting [get]
func (s *chatHandler) GetSetting(ctx context.Context, req *chatpb.GetSettingRequest) (*chatpb.GetSettingResponse, error) {
	conversation, participant, err := s.conversationUsecases.GetSetting(ctx, uint64(req.Id))
	if err != nil {
		return nil, err
	}
	return &chatpb.GetSettingResponse{
		Id:              conversation.ID,
		Type:            chatpb.ConversationTypeEnum(conversation.Type),
		Name:            conversation.Name,
		IsBroadcast:     conversation.IsBroadcast,
		ForbidForward:   conversation.ForbidForward,
		IsDeleted:       conversation.DeletedAt != nil,
		MuteNotif:       participant.MuteNotif,
		BackgroundImage: s.backgroundMapper.BackgroundImageToPb(conversation.BackgroundImage),
	}, nil
}

// @Summary Lấy danh sách tệp tin
// @Description Lấy danh sách tệp tin trong đoạn hội thoại được chỉ định
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param message query chatpb.GetFilesRequest true "Đoạn hội thoại cần lấy danh sách tệp tin"
// @Success 200 {object} chatpb.GetFilesResponse
// @Router /messages/files [get]
func (s *chatHandler) GetFiles(ctx context.Context, req *chatpb.GetFilesRequest) (*chatpb.GetFilesResponse, error) {
	files, total, err := s.messageUsecases.GetMessageByTypeWithPagination(ctx, int64(req.Page), int64(req.Size), uint64(req.ConversationId), []string{"FILE", "video", "audio", "pdf"})
	if err != nil {
		return nil, err
	}

	var pbMessage []*chatpb.Message
	for _, file := range files {
		createdAtStr := ""
		if file.CreatedAt != nil {
			createdAtStr = file.CreatedAt.Local().Format(time.RFC3339)
		}
		pbMessage = append(pbMessage, &chatpb.Message{
			MessageId:   int32(file.ID),
			FileUrl:     file.FileURL,
			Content:     file.Content,
			CreatedAt:   createdAtStr,
			SenderId:    int32(file.SenderID),
			ContentType: string(file.ContentType),
			RoomId:      strconv.FormatUint(file.ConversationID, 10),
			IsRecall:    file.Recall,
		})
	}

	return &chatpb.GetFilesResponse{
		Messages: pbMessage,
		Total:    total,
	}, nil
}

// @Summary Lấy danh sách hình ảnh
// @Description Lấy danh sách hình ảnh trong đoạn hội thoại được chỉ định
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param message query chatpb.GetImagesRequest true "Đoạn hội thoại cần lấy danh sách hình ảnh"
// @Success 200 {object} chatpb.GetImagesResponse
// @Router /messages/images [get]
func (s *chatHandler) GetImages(ctx context.Context, req *chatpb.GetImagesRequest) (*chatpb.GetImagesResponse, error) {
	files, total, err := s.messageUsecases.GetMessageByTypeWithPagination(ctx, int64(req.Page), int64(req.Size), uint64(req.ConversationId), []string{"IMAGE"})
	if err != nil {
		return nil, err
	}

	var pbMessage []*chatpb.Message
	for _, file := range files {
		createdAtStr := ""
		if file.CreatedAt != nil {
			createdAtStr = file.CreatedAt.Local().Format(time.RFC3339)
		}
		pbMessage = append(pbMessage, &chatpb.Message{
			MessageId:   int32(file.ID),
			FileUrl:     file.FileURL,
			Content:     file.Content,
			CreatedAt:   createdAtStr,
			SenderId:    int32(file.SenderID),
			ContentType: string(file.ContentType),
			RoomId:      strconv.FormatUint(file.ConversationID, 10),
			IsRecall:    file.Recall,
		})
	}

	return &chatpb.GetImagesResponse{
		Data:  pbMessage,
		Total: total,
	}, nil
}

// @Summary Lấy danh sách liên kết
// @Description Lấy danh sách liên kết trong đoạn hội thoại được chỉ định
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param message query chatpb.GetLinksRequest true "Đoạn hội thoại cần lấy danh sách liên kết"
// @Success 200 {object} chatpb.GetLinksResponse
// @Router /messages/links [get]
func (s *chatHandler) GetLinks(ctx context.Context, req *chatpb.GetLinksRequest) (*chatpb.GetLinksResponse, error) {
	files, total, err := s.messageUsecases.GetMessageByTypeWithPagination(ctx, int64(req.Page), int64(req.Size), uint64(req.ConversationId), []string{"LINK"})
	if err != nil {
		return nil, err
	}

	var pbMessage []*chatpb.Message
	for _, file := range files {
		createdAtStr := ""
		if file.CreatedAt != nil {
			createdAtStr = file.CreatedAt.Local().Format(time.RFC3339)
		}
		pbMessage = append(pbMessage, &chatpb.Message{
			MessageId:   int32(file.ID),
			FileUrl:     file.FileURL,
			Content:     file.Content,
			CreatedAt:   createdAtStr,
			SenderId:    int32(file.SenderID),
			ContentType: string(file.ContentType),
			RoomId:      strconv.FormatUint(file.ConversationID, 10),
			IsRecall:    file.Recall,
		})
	}

	return &chatpb.GetLinksResponse{
		Messages: pbMessage,
		Total:    total,
	}, nil
}

// @Summary Lấy timestamps các loại data của conversation
// @Description API check sync - lấy timestamp từ Redis cho từng key, FE so sánh với cache và chỉ gọi API tương ứng khi có thay đổi. Keys: detail, files, links, images, members, setting, pinned, messages
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Conversation ID"
// @Success 200 {object} chatpb.ConversationTimestampsResponse "Map key -> timestamp (ms)"
// @Router /conversation/{id}/timestamp [get]
func (s *chatHandler) GetConversationTimestamps(ctx context.Context, req *sharepb.IdRequest) (*chatpb.ConversationTimestampsResponse, error) {
	if req == nil || req.Id == 0 {
		return nil, _errors.ReturnError(_errors.RequestValidationFailed, _errors.WithPublicMessage("conversation id is required"))
	}

	if s.syncProvider == nil {
		return &chatpb.ConversationTimestampsResponse{Timestamps: map[string]int64{}}, nil
	}

	convID := req.Id
	redisKeys := []string{
		s.syncProvider.GetKey(ctx, _utils.SyncKeyConversationDetail, convID),
		s.syncProvider.GetKey(ctx, _utils.SyncKeyConversationFiles, convID),
		s.syncProvider.GetKey(ctx, _utils.SyncKeyConversationLinks, convID),
		s.syncProvider.GetKey(ctx, _utils.SyncKeyConversationImages, convID),
		s.syncProvider.GetKey(ctx, _utils.SyncKeyConversationMembers, convID),
		s.syncProvider.GetKey(ctx, _utils.SyncKeyConversationSetting, convID),
		s.syncProvider.GetKey(ctx, _utils.SyncKeyConversationPinned, convID),
		s.syncProvider.GetKey(ctx, _utils.SyncKeyConversationMessages, convID),
	}

	vals, err := s.syncProvider.MGet(ctx, redisKeys)
	if err != nil {
		return nil, fmt.Errorf("get conversation timestamps: %w", err)
	}

	names := []string{
		"detail", "files", "links", "images", "members", "setting", "pinned", "messages",
	}

	timestamps := make(map[string]int64, len(names))
	for i, v := range vals {
		if v == nil {
			timestamps[names[i]] = 0
			continue
		}
		ts, _ := strconv.ParseInt(fmt.Sprintf("%v", v), 10, 64)
		timestamps[names[i]] = ts
	}

	return &chatpb.ConversationTimestampsResponse{
		Timestamps: timestamps,
	}, nil
}

// @Summary Thêm phản hồi
// @Description Thêm phản hồi vào tin nhắn đã được gửi trong đoạn hội thoại được chỉ định
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param message body chatpb.ReactionMessageRequest true "Đoạn hội thoại và tin nhắn cần thêm phản hồi"
// @Success 200 {object} chatpb.ReactionMessageResponse
// @Router /messages/reaction [post]
func (s *chatHandler) ReactionMessage(ctx context.Context, req *chatpb.ReactionMessageRequest) (*chatpb.ReactionMessageResponse, error) {
	model, err := s.messageReactionUsecases.CreateMessageReaction(ctx, uint64(req.MessageId), req.Reaction)
	if err != nil {
		return nil, err
	}
	createdAtStr := ""
	if model.CreatedAt != nil {
		createdAtStr = model.CreatedAt.Local().Format(time.RFC3339)
	}
	event := &chatpb.ReactionMessage{
		Type: constants.REACT_MESSAGE_TYPE,
		Data: &chatpb.MessageReaction{
			MessageId: int32(model.MessageID),
			UserId:    int32(model.UserID),
			Reaction:  model.Reaction,
			CreatedAt: createdAtStr,
			Id:        model.ID,
			RoomId:    strconv.FormatUint(model.ConversationId, 10),
		},
	}
	user, _ := s.userClient.GetUserById(ctx, model.UserID)
	if user != nil {
		event.Data.FullName = &user.FullName
		event.Data.Avatar = &user.Avatar
	}
	s.eventChannel <- event
	return &chatpb.ReactionMessageResponse{
		Success: true,
	}, nil
}

// @Summary Lấy số lượng tin nhắn chưa đọc
// @Description Lấy số lượng tin nhắn chưa đọc của người dùng hiện tại
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} chatpb.UnreadCountResponse
// @Router /unread-count [get]
func (s *chatHandler) UnreadCount(ctx context.Context, _ *chatpb.UnreadCountRequest) (*chatpb.UnreadCountResponse, error) {
	unreadCount, err := s.conversationUsecases.GetUnreadCount(ctx)
	if err != nil {
		return nil, err
	}
	return &chatpb.UnreadCountResponse{
		UnreadCount: unreadCount,
		Timestamp:   time.Now().UTC().UnixMilli(),
	}, nil
}

func (s *chatHandler) Start() {
	s.startOnce.Do(func() {
		for i := 0; i < s.numberOfWorker; i++ {
			redisClientWorker := NewRedisClientWorker(s.workerCtx, s, s.redisClient)
			s.workerWG.Add(1)
			go func() {
				defer s.workerWG.Done()
				redisClientWorker.serve()
			}()
		}
	})
}

func (s *chatHandler) Close() error {
	var closeErr error
	s.closeOnce.Do(func() {
		s.cancelWorkers()
		s.workerWG.Wait()
		closeErr = s.redisClient.Close()
	})
	return closeErr
}

func (s *chatHandler) ValidConversationAndCurrentUser(ctx context.Context, req *sharepb.IdRequest) (*sharepb.Empty, error) {
	err := s.conversationUsecases.ValidateConversationAndCurrentUser(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &sharepb.Empty{}, nil
}

func (s *chatHandler) GetMembersOfRoom(ctx context.Context, req *sharepb.IdRequest) (*chatpb.GetRoomMemberIdResponse, error) {
	fmt.Println("GetMembersOfRoom", req.Id)
	members, err := s.conversationUsecases.GetMembersOfRoom(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &chatpb.GetRoomMemberIdResponse{
		ProfileIds: members,
	}, nil
}

func (s *chatHandler) SetBackgroundImage(ctx context.Context, req *chatpb.SetBackgroundImageRequest) (*chatpb.SetBackgroundImageResponse, error) {
	conversation, err := s.conversationUsecases.SetBackgroundImage(ctx, uint64(req.ConversationId), uint64(req.BackgroundImageId))
	if err != nil {
		return nil, err
	}

	event := &chatpb.UpdateSettings{
		Type: constants.UPDATE_SETTINGS_TYPE,
		Data: s.conversationMapper.ConversationToPb(conversation),
	}

	s.eventChannel <- event

	return &chatpb.SetBackgroundImageResponse{
		Success: true,
	}, nil
}
