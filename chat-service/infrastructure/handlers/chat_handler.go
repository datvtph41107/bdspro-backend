package handlers

import (
	_utils "common/utils"
	"context"
	"fmt"
	"strconv"
	"time"

	chatpb "pb/types/chat"
	sharepb "pb/types/shared"

	"chat/configs"
	"chat/infrastructure/client"
	"chat/infrastructure/delivery/errors"
	"chat/infrastructure/mapper"
	"chat/infrastructure/redis"
	"chat/internal/constants"
	"chat/models"
	"chat/utils"
	_errors "errors"

	"github.com/hyperledger/fabric/common/flogging"
)

type chatHandler struct {
	chatpb.UnimplementedChatServiceServer
	logger                  *flogging.FabricLogger
	conversationUsecases    conversationUsecases
	messageUsecases         messageUsecases
	readReceptUsecases      readReceptUsecases
	participantUsecases     participantUsecases
	messageReactionUsecases messageReactionUsecases
	eventChannel            chan any
	numberOfWorker          int
	maxParticipant          int
	redisClient             redisClient
	userClient              *client.UserClient
	bdsproClient            *client.BdsproClient
	conversationMapper      *mapper.ConversationMapper
	backgroundMapper        *mapper.BackgroundMapper
}

func NewChatHandler(
	conversationUsecases conversationUsecases,
	messageUsecases messageUsecases,
	readReceptUsecases readReceptUsecases,
	participantUsecases participantUsecases,
	messageReactionUsecases messageReactionUsecases,
	userClient *client.UserClient,
	bdsproClient *client.BdsproClient,
	conversationMapper *mapper.ConversationMapper,
	backgroundMapper *mapper.BackgroundMapper,
) *chatHandler {
	return &chatHandler{
		conversationUsecases:    conversationUsecases,
		messageUsecases:         messageUsecases,
		readReceptUsecases:      readReceptUsecases,
		participantUsecases:     participantUsecases,
		messageReactionUsecases: messageReactionUsecases,
		logger:                  flogging.MustGetLogger("chat_handler"),
		numberOfWorker:          configs.NUMBER_OF_WORKER,
		maxParticipant:          configs.MAX_PARTICIPANT,
		redisClient:             redis.NewRedisClient(),
		eventChannel:            make(chan any),
		userClient:              userClient,
		bdsproClient:            bdsproClient,
		conversationMapper:      conversationMapper,
		backgroundMapper:        backgroundMapper,
	}
}

// @Summary Tạo đoạn hội thoại
// @Description Tạo đoạn hội thoại với các thành viên được chỉ định
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param conversation body chatpb.CreateConversationRequest true "Đoạn hội thoại cần tạo"
// @Success 200 {object} chatpb.CreateConversationResponse
// @Router /conversations [post]
func (s *chatHandler) CreateConversation(ctx context.Context, req *chatpb.CreateConversationRequest) (*chatpb.CreateConversationResponse, error) {
	if len(req.Members) > s.maxParticipant {
		return nil, errors.MaxParticipantsReached()
	}

	if len(req.Members) == 0 && req.ReceiverId == nil {
		return nil, errors.NotParticipant()
	}

	conversation := &models.ConversationModel{
		Name:              &req.Name,
		Type:              models.ConversationTypeEnum(req.Type),
		CreatedBy:         uint64(req.CreatedBy),
		IsBroadcast:       req.IsBroadcast,
		ForbidForward:     req.ForbidForward,
		Avatar:            req.Avatar,
		BackgroundImageID: req.BackgroundImageId,
		ReceiverId:        req.ReceiverId,
	}
	model, err := s.conversationUsecases.CreateConversation(ctx, conversation, utils.Int32SliceToUint64Slice(req.Members))
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
	fmt.Println("req", req)
	message, err := s.messageUsecases.SendMessage(ctx, &models.MessageModel{
		ConversationID: uint64(req.ConversationId),
		SenderID:       uint64(req.SenderId),
		Content:        req.Content,
		ContentType:    models.ContentTypeEnum(req.ContentType),
		FileURL:        req.FileUrl,
		ReplyId:        req.ReplyId,
		ExtraId:        req.ExtraId,
	})
	if err != nil {
		return nil, err
	}

	return s.ResponseMessage(ctx, message)
}

func (s *chatHandler) ResponseMessage(ctx context.Context, message *models.MessageModel) (*chatpb.SendMessageResponse, error) {
	event := &chatpb.SendMessage{
		Type: constants.MESSAGE_TYPE,
		Data: &chatpb.Message{
			MessageId:   int32(message.ID),
			SenderId:    int32(message.SenderID),
			Content:     message.Content,
			ContentType: string(message.ContentType),
			FileUrl:     message.FileURL,
			RoomId:      strconv.FormatUint(message.ConversationID, 10),
			CreatedAt:   message.CreatedAt.Local().Format(time.RFC3339),
			IsRecall:    message.Recall,
			ExtraId:     message.ExtraId,
		},
	}
	if message.ReplyMessage != nil {
		event.Data.ReplyMessage = &chatpb.Message{
			MessageId:   int32(message.ReplyMessage.ID),
			SenderId:    int32(message.ReplyMessage.SenderID),
			Content:     message.ReplyMessage.Content,
			ContentType: string(message.ReplyMessage.ContentType),
			FileUrl:     message.ReplyMessage.FileURL,
			RoomId:      strconv.FormatUint(message.ReplyMessage.ConversationID, 10),
			CreatedAt:   message.ReplyMessage.CreatedAt.Local().Format(time.RFC3339),
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
	return &chatpb.SendMessageResponse{
		MessageId: int32(message.ID),
		CreatedAt: message.CreatedAt.Local().Format(time.RFC3339),
	}, nil
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
	receiverID := req.ReceiverId
	if receiverID == nil {
		return nil, _errors.New("receiverId is required")
	}
	message, err := s.messageUsecases.SendToReceiver(ctx, &models.MessageModel{
		ConversationID: uint64(req.ConversationId),
		SenderID:       uint64(req.SenderId),
		Content:        req.Content,
		ContentType:    models.ContentTypeEnum(req.ContentType),
		FileURL:        req.FileUrl,
		ReplyId:        req.ReplyId,
		ExtraId:        req.ExtraId,
	}, req.ReceiverId)
	if err != nil {
		return nil, err
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
			CreatedAt:   message.CreatedAt.Local().Format(time.RFC3339),
			IsRecall:    message.Recall,
		},
	}
	if message.ReplyMessage != nil {
		event.Data.ReplyMessage = &chatpb.Message{
			MessageId:   int32(message.ReplyMessage.ID),
			SenderId:    int32(message.ReplyMessage.SenderID),
			Content:     message.ReplyMessage.Content,
			ContentType: string(message.ReplyMessage.ContentType),
			FileUrl:     message.ReplyMessage.FileURL,
			RoomId:      strconv.FormatUint(message.ReplyMessage.ConversationID, 10),
			CreatedAt:   message.ReplyMessage.CreatedAt.Local().Format(time.RFC3339),
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
	return &chatpb.SendMessageResponse{
		MessageId: int32(message.ID),
		CreatedAt: message.CreatedAt.Local().Format(time.RFC3339),
	}, nil
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
	var limit, offset int64
	var sort string
	var fromDate, toDate *time.Time
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

	if req.Sort != nil {
		sort = *req.Sort
	} else {
		sort = "desc"
	}

	if req.FromDate != nil {
		tempFromDate := _utils.ParseStringToTime(*req.FromDate)
		fromDate = tempFromDate
	}

	if req.ToDate != nil {
		tempToDate := _utils.ParseStringToTime(*req.ToDate)
		toDate = tempToDate
	}

	messages, err := s.messageUsecases.GetMessages(ctx, limit, offset, uint64(req.ConversationId), sort, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	var pbMessages []*chatpb.Message
	for _, msg := range messages {
		if msg.Recall {
			pbMessages = append(pbMessages, &chatpb.Message{
				MessageId: int32(msg.ID),
				SenderId:  int32(msg.SenderID),
				IsRecall:  true,
				CreatedAt: _utils.FormatTimeToString(&msg.CreatedAt),
			})
			continue
		}
		newMessage := &chatpb.Message{
			MessageId:   int32(msg.ID),
			SenderId:    int32(msg.SenderID),
			Content:     msg.Content,
			ContentType: string(msg.ContentType),
			FileUrl:     msg.FileURL,
			RoomId:      strconv.FormatUint(msg.ConversationID, 10),
			CreatedAt:   msg.CreatedAt.Local().Format(time.RFC3339),
			IsRecall:    msg.Recall,
			ExtraId:     msg.ExtraId,
		}
		if msg.ReplyMessage != nil {
			newMessage.ReplyMessage = &chatpb.Message{
				MessageId:   int32(msg.ReplyMessage.ID),
				SenderId:    int32(msg.ReplyMessage.SenderID),
				Content:     msg.ReplyMessage.Content,
				ContentType: string(msg.ReplyMessage.ContentType),
				FileUrl:     msg.ReplyMessage.FileURL,
				RoomId:      strconv.FormatUint(msg.ReplyMessage.ConversationID, 10),
				CreatedAt:   msg.ReplyMessage.CreatedAt.Local().Format(time.RFC3339),
			}
		}
		if msg.Reactions != nil {
			newMessage.Reactions = make([]*chatpb.MessageReaction, len(msg.Reactions))
			for i, reaction := range msg.Reactions {
				newMessage.Reactions[i] = &chatpb.MessageReaction{
					UserId:    int32(reaction.UserID),
					MessageId: int32(reaction.MessageID),
					CreatedAt: reaction.CreatedAt.Local().Format(time.RFC3339),
					Reaction:  reaction.Reaction,
					Id:        uint64(reaction.ID),
				}
			}
		}

		pbMessages = append(pbMessages, newMessage)
	}

	s.userClient.MapAvatarAndNameToMessagePbList(ctx, pbMessages)
	s.bdsproClient.MapProductPbByIds(ctx, pbMessages)

	return &chatpb.GetMessagesResponse{
		Messages: pbMessages,
	}, nil
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
		messagePb := &chatpb.Message{
			MessageId:   int32(message.ID),
			SenderId:    int32(message.SenderID),
			Content:     message.Content,
			ContentType: string(message.ContentType),
			FileUrl:     message.FileURL,
			RoomId:      strconv.FormatUint(message.ConversationID, 10),
			CreatedAt:   message.CreatedAt.Local().Format(time.RFC3339),
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
	err := s.readReceptUsecases.MarkAsRead(ctx, uint64(req.ConversationId), utils.Int32SliceToUint64Slice(req.MessageIds))
	if err != nil {
		return nil, err
	}
	userIdStr := strconv.FormatUint(utils.GetCurrentUserID(ctx), 10)
	event := &chatpb.ReadMessage{
		Type: constants.READ_TYPE,
		Data: &chatpb.ReadMessageData{
			RoomId:     strconv.Itoa(int(req.ConversationId)),
			UserId:     userIdStr,
			MessageIds: utils.Int32SliceToStringSlice(req.MessageIds),
		},
	}

	user, _ := s.userClient.GetUserById(ctx, utils.GetCurrentUserID(ctx))
	if user != nil {
		event.Data.FullName = &user.FullName
		event.Data.Avatar = &user.Avatar
	}

	s.eventChannel <- event
	return &chatpb.MarkAsReadResponse{
		Success: true,
	}, nil
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
		Data: &chatpb.PinMessageData{
			RoomId:    strconv.Itoa(int(req.ConversationId)),
			MessageId: strconv.Itoa(int(req.MessageId)),
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
	channel, err := s.conversationUsecases.CreateOrgChannel(ctx, &models.ConversationModel{})
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
	var pbParticipants []*chatpb.Participant
	for _, participant := range participants {
		pbParticipants = append(pbParticipants, &chatpb.Participant{
			UserId:    (participant.UserID),
			Role:      chatpb.ParticipantRoleEnum(participant.Role),
			JoinedAt:  participant.JoinedAt.Local().Format(time.RFC3339),
			MuteNotif: participant.MuteNotif,
		})
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
	messages, total, err := s.readReceptUsecases.GetReadReceipt(ctx, req.Page, req.Size, req.MessageId)
	if err != nil {
		return nil, err
	}
	var pbReadReceipts []*chatpb.ReadReceipt
	for _, readReceipt := range messages {
		pbReadReceipts = append(pbReadReceipts, &chatpb.ReadReceipt{
			UserId: (readReceipt.UserID),
			ReadAt: readReceipt.ReadAt.Local().Format(time.RFC3339),
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
// @Description Thêm thành viên vào đoạn hội thoại được chỉ định
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param message body chatpb.AddUserRequest true "Đoạn hội thoại và thành viên cần thêm"
// @Success 200 {object} chatpb.AddUserResponse
// @Router /conversation/add-user [post]
func (s *chatHandler) AddUser(ctx context.Context, req *chatpb.AddUserRequest) (*chatpb.AddUserResponse, error) {
	err := s.participantUsecases.AddUser(ctx, uint64(req.ConversationId), uint64(req.UserId))
	if err != nil {
		return nil, err
	}
	event := &chatpb.JoinRoom{
		Type: constants.JOIN_ROOM_TYPE,
		Data: &chatpb.JoinRoomData{
			RoomId: strconv.Itoa(int(req.ConversationId)),
			UserId: strconv.Itoa(int(req.UserId)),
		},
	}
	user, _ := s.userClient.GetUserById(ctx, uint64(req.UserId))
	if user != nil {
		event.Data.FullName = &user.FullName
		event.Data.Avatar = &user.Avatar
	}
	s.eventChannel <- event
	return &chatpb.AddUserResponse{
		Success: true,
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
		IsDeleted:       conversation.IsDeleted,
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
		pbMessage = append(pbMessage, &chatpb.Message{
			MessageId:   int32(file.ID),
			FileUrl:     file.FileURL,
			Content:     file.Content,
			CreatedAt:   file.CreatedAt.Local().Format(time.RFC3339),
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
		pbMessage = append(pbMessage, &chatpb.Message{
			MessageId:   int32(file.ID),
			FileUrl:     file.FileURL,
			Content:     file.Content,
			CreatedAt:   file.CreatedAt.Local().Format(time.RFC3339),
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
		pbMessage = append(pbMessage, &chatpb.Message{
			MessageId:   int32(file.ID),
			FileUrl:     file.FileURL,
			Content:     file.Content,
			CreatedAt:   file.CreatedAt.Local().Format(time.RFC3339),
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
	event := &chatpb.ReactionMessage{
		Type: constants.REACT_MESSAGE_TYPE,
		Data: &chatpb.MessageReaction{
			MessageId: int32(model.MessageID),
			UserId:    int32(model.UserID),
			Reaction:  model.Reaction,
			CreatedAt: model.CreatedAt.Local().Format(time.RFC3339),
			Id:        uint64(model.ID),
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
func (s *chatHandler) UnreadCount(ctx context.Context, req *sharepb.Empty) (*chatpb.UnreadCountResponse, error) {
	unreadCount, err := s.conversationUsecases.GetUnreadCount(ctx)
	if err != nil {
		return nil, err
	}
	return &chatpb.UnreadCountResponse{UnreadCount: unreadCount}, nil
}

func (s *chatHandler) Start() {
	for i := 0; i < s.numberOfWorker; i++ {
		redisClientWorker := NewRedisClientWorker(s, s.redisClient, s.logger)
		go redisClientWorker.serve()
	}
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
