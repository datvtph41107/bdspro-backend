package handlers

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hyperledger/fabric/common/flogging"
	"github.com/tuandq3/go-chat-service/configs"
	"github.com/tuandq3/go-chat-service/infrastructure/redis"
	"github.com/tuandq3/go-chat-service/internal/constants"
	"github.com/tuandq3/go-chat-service/internal/delivery/errors"
	"github.com/tuandq3/go-chat-service/models"
	pb "github.com/tuandq3/go-chat-service/proto/types/chat"
	messagepb "github.com/tuandq3/go-chat-service/proto/types/message"
	"github.com/tuandq3/go-chat-service/utils"
)

type chatHandler struct {
	pb.UnimplementedChatServiceServer
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
}

func NewChatHandler(
	conversationUsecases conversationUsecases,
	messageUsecases messageUsecases,
	readReceptUsecases readReceptUsecases,
	participantUsecases participantUsecases,
	messageReactionUsecases messageReactionUsecases,
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
	}
}

func (s *chatHandler) CreateConversation(ctx context.Context, req *pb.CreateConversationRequest) (*pb.CreateConversationResponse, error) {
	if len(req.Members) > s.maxParticipant {
		return nil, errors.MaxParticipantsReached()
	}

	if len(req.Members) == 0 {
		return nil, errors.NotParticipant()
	}

	model, err := s.conversationUsecases.CreateConversation(ctx, &models.ConversationModel{
		Name:            &req.Name,
		Type:            models.ConversationTypeEnum(req.Type),
		CreatedBy:       uint64(req.CreatedBy),
		IsBroadcast:     req.IsBroadcast,
		ForbidForward:   req.ForbidForward,
		Avatar:          req.Avatar,
		BackgroundImage: req.BackgroundImage,
	}, utils.Int32SliceToUint64Slice(req.Members))
	if err != nil {
		return nil, err
	}
	event := &messagepb.CreateRoom{
		Type: constants.CREATE_ROOM_TYPE,
		Data: &messagepb.CreateRoomData{
			RoomId:  strconv.Itoa(int(model.ID)),
			Members: utils.Int32SliceToStringSlice(req.Members),
		},
	}
	s.eventChannel <- event
	return &pb.CreateConversationResponse{
		ConversationId: int32(model.ID),
	}, nil
}

func (s *chatHandler) GetConversations(ctx context.Context, req *pb.GetConversationsRequest) (*pb.GetConversationsResponse, error) {
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

	conversations, err := s.conversationUsecases.GetConversations(ctx, limit, offset, models.ConversationTypeEnum(*req.Type), req.Keyword)
	if err != nil {
		return nil, err
	}

	var pbConversations []*pb.Conversation
	for _, conv := range conversations {
		if conv.LatestMessage.ID == 0 {
			pbConversations = append(pbConversations, &pb.Conversation{
				ConversationId:  int32(conv.ID),
				Name:            *conv.Name,
				Type:            pb.ConversationTypeEnum(conv.Type),
				UpdatedAt:       conv.UpdatedAt.Local().Format(time.RFC3339),
				Avatar:          conv.Avatar,
				UnreadCount:     int32(conv.UnreadCount),
				BackgroundImage: conv.BackgroundImage,
			})
		} else {
			pbConversations = append(pbConversations, &pb.Conversation{
				ConversationId: int32(conv.ID),
				Name:           *conv.Name,
				Type:           pb.ConversationTypeEnum(conv.Type),
				UpdatedAt:      conv.UpdatedAt.Format(time.DateTime),
				Avatar:         conv.Avatar,
				UnreadCount:    int32(conv.UnreadCount),
				LastMessage: &pb.Message{
					MessageId:   int32(conv.LatestMessage.ID),
					SenderId:    int32(conv.LatestMessage.SenderID),
					Content:     conv.LatestMessage.Content,
					ContentType: string(conv.LatestMessage.ContentType),
					FileUrl:     conv.LatestMessage.FileURL,
					CreatedAt:   conv.LatestMessage.CreatedAt.Local().Format(time.RFC3339),
				},
				BackgroundImage: conv.BackgroundImage,
			})
		}
	}

	return &pb.GetConversationsResponse{
		Conversations: pbConversations,
	}, nil
}

func (s *chatHandler) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {
	message, err := s.messageUsecases.SendMessage(ctx, &models.MessageModel{
		ConversationID: uint64(req.ConversationId),
		SenderID:       uint64(req.SenderId),
		Content:        req.Content,
		ContentType:    models.ContentTypeEnum(req.ContentType),
		FileURL:        req.FileUrl,
		ReplyId:        req.ReplyId,
	})
	if err != nil {
		return nil, err
	}

	event := &messagepb.SendMessage{
		Type: constants.MESSAGE_TYPE,
		Data: &pb.Message{
			MessageId:   int32(message.ID),
			SenderId:    int32(message.SenderID),
			Content:     message.Content,
			ContentType: string(message.ContentType),
			FileUrl:     message.FileURL,
			RoomId:      strconv.FormatUint(message.ConversationID, 10),
			CreatedAt:   message.CreatedAt.Local().Format(time.RFC3339),
		},
	}
	s.eventChannel <- event

	updateRoomEvent := &messagepb.UpdateRoomMessage{
		Type: constants.UPDATE_ROOM_TYPE,
		Data: &messagepb.UpdateRoomMessageData{
			RoomId: strconv.FormatUint(message.ConversationID, 10),
		},
	}
	s.eventChannel <- updateRoomEvent
	return &pb.SendMessageResponse{
		MessageId: int32(message.ID),
		CreatedAt: message.CreatedAt.Local().Format(time.RFC3339),
	}, nil
}

func (s *chatHandler) GetMessages(ctx context.Context, req *pb.GetMessagesRequest) (*pb.GetMessagesResponse, error) {
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
		tempFromDate := req.FromDate.AsTime()
		fromDate = &tempFromDate
	}

	if req.ToDate != nil {
		tempToDate := req.ToDate.AsTime()
		toDate = &tempToDate
	}

	messages, err := s.messageUsecases.GetMessages(ctx, limit, offset, uint64(req.ConversationId), sort, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	var pbMessages []*pb.Message
	for _, msg := range messages {
		newMessage := &pb.Message{
			MessageId:   int32(msg.ID),
			SenderId:    int32(msg.SenderID),
			Content:     msg.Content,
			ContentType: string(msg.ContentType),
			FileUrl:     msg.FileURL,
			RoomId:      strconv.FormatUint(msg.ConversationID, 10),
			CreatedAt:   msg.CreatedAt.Local().Format(time.RFC3339),
			IsRecall:    msg.Recall,
		}
		if msg.ReplyMessage != nil {
			newMessage.ReplyMessage = &pb.Message{
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
			newMessage.Reactions = make([]*pb.MessageReaction, len(msg.Reactions))
			for i, reaction := range msg.Reactions {
				newMessage.Reactions[i] = &pb.MessageReaction{
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

	return &pb.GetMessagesResponse{
		Messages: pbMessages,
	}, nil
}

func (s *chatHandler) ForwardMessage(ctx context.Context, req *pb.ForwardMessageRequest) (*pb.ForwardMessageResponse, error) {
	message, err := s.messageUsecases.ForwardMessage(ctx, uint64(req.MessageId), uint64(req.TargetConversationId))
	if err != nil {
		return nil, err
	}
	event := &messagepb.SendMessage{
		Type: constants.MESSAGE_TYPE,
		Data: &pb.Message{
			MessageId:   int32(message.ID),
			SenderId:    int32(message.SenderID),
			Content:     message.Content,
			ContentType: string(message.ContentType),
			FileUrl:     message.FileURL,
			RoomId:      strconv.FormatUint(message.ConversationID, 10),
			CreatedAt:   message.CreatedAt.Local().Format(time.RFC3339),
		},
	}
	s.eventChannel <- event
	return &pb.ForwardMessageResponse{
		MessageId: int32(message.ID),
	}, nil
}

func (s *chatHandler) MarkAsRead(ctx context.Context, req *pb.MarkAsReadRequest) (*pb.MarkAsReadResponse, error) {
	err := s.readReceptUsecases.MarkAsRead(ctx, uint64(req.ConversationId), utils.Int32SliceToUint64Slice(req.MessageIds))
	if err != nil {
		return nil, err
	}
	userIdStr := strconv.FormatUint(utils.GetCurrentUserID(ctx), 10)
	event := &messagepb.ReadMessage{
		Type: constants.READ_TYPE,
		Data: &messagepb.ReadMessageData{
			RoomId:     strconv.Itoa(int(req.ConversationId)),
			UserId:     userIdStr,
			MessageIds: utils.Int32SliceToStringSlice(req.MessageIds),
		},
	}

	s.eventChannel <- event
	return &pb.MarkAsReadResponse{
		Success: true,
	}, nil
}

func (s *chatHandler) PinMessage(ctx context.Context, req *pb.PinMessageRequest) (*pb.PinMessageResponse, error) {
	_, err := s.messageUsecases.PinMessage(ctx, uint64(req.MessageId), uint64(req.ConversationId), req.DoPin)
	if err != nil {
		return nil, err
	}

	event := &messagepb.PinMessage{
		Type: constants.PIN_MESSAGE_TYPE,
		Data: &messagepb.PinMessageData{
			RoomId:    strconv.Itoa(int(req.ConversationId)),
			MessageId: strconv.Itoa(int(req.MessageId)),
		},
	}
	s.eventChannel <- event

	return &pb.PinMessageResponse{
		Success: true,
	}, nil
}

func (s *chatHandler) RecallMessage(ctx context.Context, req *pb.RecallMessageRequest) (*pb.RecallMessageResponse, error) {
	message, err := s.messageUsecases.RecallMessage(ctx, uint64(req.MessageId))
	if err != nil {
		return nil, err
	}
	event := &messagepb.RecallMessage{
		Type: constants.RECALL_MESSAGE_TYPE,
		Data: &messagepb.RecallMessageData{
			RoomId:    strconv.Itoa(int(message.ConversationID)),
			MessageId: strconv.Itoa(int(req.MessageId)),
		},
	}
	s.eventChannel <- event

	return &pb.RecallMessageResponse{
		Success: true,
	}, nil
}

func (s *chatHandler) RemoveUser(ctx context.Context, req *pb.RemoveUserRequest) (*pb.RemoveUserResponse, error) {
	err := s.participantUsecases.RemoveUser(ctx, uint64(req.ConversationId), uint64(req.TargetUserId))
	if err != nil {
		return nil, err
	}
	event := &messagepb.RemoveMember{
		Type: constants.REMOVE_MEMBER_TYPE,
		Data: &messagepb.RemoveMemberData{
			RoomId: strconv.Itoa(int(req.ConversationId)),
			UserId: strconv.Itoa(int(req.TargetUserId)),
		},
	}
	s.eventChannel <- event
	return &pb.RemoveUserResponse{
		Success: true,
	}, nil
}

func (s *chatHandler) EditChannel(ctx context.Context, req *pb.EditChannelRequest) (*pb.EditChannelResponse, error) {
	err := s.conversationUsecases.EditChannel(ctx, uint64(req.ConversationId), req.IsBroadcast, req.ForbidForward, req.Name)
	if err != nil {
		return nil, err
	}
	return &pb.EditChannelResponse{
		Success: true,
	}, nil
}

func (s *chatHandler) CreateOrgChannel(ctx context.Context, req *pb.CreateOrgChannelRequest) (*pb.CreateOrgChannelResponse, error) {
	channel, err := s.conversationUsecases.CreateOrgChannel(ctx, &models.ConversationModel{})
	if err != nil {
		return nil, err
	}

	return &pb.CreateOrgChannelResponse{
		ConversationId: int32(channel.ID),
	}, nil
}

func (s *chatHandler) MuteConversation(ctx context.Context, req *pb.MuteConversationRequest) (*pb.MuteConversationResponse, error) {
	err := s.participantUsecases.MuteConversation(ctx, uint64(req.ConversationId), req.Mute)
	if err != nil {
		return nil, err
	}

	return &pb.MuteConversationResponse{
		Success: true,
	}, nil
}

func (s *chatHandler) Search(ctx context.Context, req *pb.SearchRequest) (*pb.SearchResponse, error) {
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

	var pbMessages []*pb.Message
	for _, msg := range messages {
		pbMessages = append(pbMessages, &pb.Message{
			MessageId:   int32(msg.ID),
			SenderId:    int32(msg.SenderID),
			Content:     msg.Content,
			ContentType: string(msg.ContentType),
			FileUrl:     msg.FileURL,
			RoomId:      strconv.FormatUint(msg.ConversationID, 10),
			CreatedAt:   msg.CreatedAt.Local().Format(time.RFC3339),
		})
	}

	var pbConversations []*pb.Conversation
	for _, conv := range conversations {
		pbConversations = append(pbConversations, &pb.Conversation{
			ConversationId: int32(conv.ID),
			Name:           *conv.Name,
			Type:           pb.ConversationTypeEnum(conv.Type),
		})
	}

	return &pb.SearchResponse{
		Messages:      pbMessages,
		Conversations: pbConversations,
	}, nil
}

func (s *chatHandler) GetConversationById(ctx context.Context, req *pb.GetConversationByIdRequest) (*pb.GetConversationByIdResponse, error) {
	conversation, err := s.conversationUsecases.GetConversationById(ctx, uint64(req.ConversationId), req.IncludeMembers)
	if err != nil {
		return nil, err
	}
	var participants []*pb.Participant
	if req.IncludeMembers {
		for _, participant := range conversation.Participants {
			participants = append(participants, &pb.Participant{
				UserId:    (participant.UserID),
				Role:      pb.ParticipantRoleEnum(participant.Role),
				JoinedAt:  participant.JoinedAt.Local().Format(time.RFC3339),
				MuteNotif: participant.MuteNotif,
			})
		}
	}

	return &pb.GetConversationByIdResponse{
		ConversationId:  int32(conversation.ID),
		Name:            *conversation.Name,
		Type:            pb.ConversationTypeEnum(conversation.Type),
		CreatedBy:       int32(conversation.CreatedBy),
		IsBroadcast:     conversation.IsBroadcast,
		ForbidForward:   conversation.ForbidForward,
		CreatedAt:       conversation.CreatedAt.Local().Format(time.RFC3339),
		UpdatedAt:       conversation.UpdatedAt.Local().Format(time.RFC3339),
		Participants:    participants,
		BackgroundImage: conversation.BackgroundImage,
		Avatar:          conversation.Avatar,
	}, nil
}

func (s *chatHandler) GetListParticipant(ctx context.Context, req *pb.GetListParticipantRequest) (*pb.GetListParticipantResponse, error) {
	participants, total, err := s.participantUsecases.GetListParticipant(ctx, req.ConversationId)
	if err != nil {
		return nil, err
	}
	var pbParticipants []*pb.Participant
	for _, participant := range participants {
		pbParticipants = append(pbParticipants, &pb.Participant{
			UserId:    (participant.UserID),
			Role:      pb.ParticipantRoleEnum(participant.Role),
			JoinedAt:  participant.JoinedAt.Local().Format(time.RFC3339),
			MuteNotif: participant.MuteNotif,
		})
	}
	return &pb.GetListParticipantResponse{
		Participants: pbParticipants,
		Total:        total,
	}, nil
}

func (s *chatHandler) GetReadReceipt(ctx context.Context, req *pb.GetReadReceiptRequest) (*pb.GetReadReceiptResponse, error) {
	messages, total, err := s.readReceptUsecases.GetReadReceipt(ctx, req.Page, req.Size, req.MessageId)
	if err != nil {
		return nil, err
	}
	var pbReadReceipts []*pb.ReadReceipt
	for _, readReceipt := range messages {
		pbReadReceipts = append(pbReadReceipts, &pb.ReadReceipt{
			UserId: (readReceipt.UserID),
			ReadAt: readReceipt.ReadAt.Local().Format(time.RFC3339),
		})
	}
	return &pb.GetReadReceiptResponse{
		ReadReceipts: pbReadReceipts,
		Total:        total,
	}, nil
}

func (s *chatHandler) GetPinnedMessages(ctx context.Context, req *pb.GetMessagesRequest) (*pb.GetMessagesResponse, error) {
	messages, err := s.messageUsecases.GetPinnedMessages(ctx, uint64(req.ConversationId))
	if err != nil {
		return nil, err
	}

	var pbMessages []*pb.Message
	for _, msg := range messages {
		pbMessages = append(pbMessages, &pb.Message{
			MessageId:   int32(msg.ID),
			SenderId:    int32(msg.SenderID),
			Content:     msg.Content,
			ContentType: string(msg.ContentType),
			FileUrl:     msg.FileURL,
			RoomId:      strconv.FormatUint(msg.ConversationID, 10),
			CreatedAt:   msg.CreatedAt.Local().Format(time.RFC3339),
		})
	}

	return &pb.GetMessagesResponse{
		Messages: pbMessages,
		Total:    int64(len(messages)),
	}, nil
}

func (s *chatHandler) AddUser(ctx context.Context, req *pb.AddUserRequest) (*pb.AddUserResponse, error) {
	err := s.participantUsecases.AddUser(ctx, uint64(req.ConversationId), uint64(req.UserId))
	if err != nil {
		return nil, err
	}
	event := &messagepb.JoinRoom{
		Type: constants.JOIN_ROOM_TYPE,
		Data: &messagepb.JoinRoomData{
			RoomId: strconv.Itoa(int(req.ConversationId)),
			UserId: strconv.Itoa(int(req.UserId)),
		},
	}
	s.eventChannel <- event
	return &pb.AddUserResponse{
		Success: true,
	}, nil
}

func (s *chatHandler) LeaveRoom(ctx context.Context, req *pb.LeaveRoomRequest) (*pb.LeaveRoomResponse, error) {
	err := s.conversationUsecases.LeaveRoom(ctx, uint64(req.ConversationId))
	if err != nil {
		return nil, err
	}
	event := &messagepb.LeaveRoom{
		Type: constants.LEAVE_ROOM_TYPE,
		Data: &messagepb.LeaveRoomData{
			RoomId: strconv.Itoa(int(req.ConversationId)),
		},
	}
	s.eventChannel <- event
	return &pb.LeaveRoomResponse{
		Success: true,
	}, nil
}

func (s *chatHandler) EditMessage(ctx context.Context, req *pb.EditMessageRequest) (*pb.EditMessageResponse, error) {
	message, err := s.messageUsecases.EditMessage(ctx, uint64(req.MessageId), req.Content)
	if err != nil {
		return nil, err
	}
	event := &messagepb.EditMessage{
		Type: constants.EDIT_MESSAGE_TYPE,
		Data: &messagepb.EditMessageData{
			RoomId:     strconv.Itoa(int(message.ConversationID)),
			MessageId:  strconv.Itoa(int(req.MessageId)),
			NewMessage: message.Content,
		},
	}
	s.eventChannel <- event
	return &pb.EditMessageResponse{
		Success: true,
	}, nil
}

func (s *chatHandler) DeleteViolation(ctx context.Context, req *pb.DeleteViolationRequest) (*pb.DeleteViolationResponse, error) {
	err := s.messageUsecases.DeleteViolation(ctx, uint64(req.MessageId))
	if err != nil {
		return nil, err
	}
	event := &messagepb.DeleteViolationMessage{
		Type: constants.DELETE_VIOLATION_MESSAGE_TYPE,
		Data: &messagepb.DeleteViolationMessageData{
			MessageId: strconv.Itoa(int(req.MessageId)),
		},
	}
	s.eventChannel <- event
	return &pb.DeleteViolationResponse{
		Success: true,
	}, nil
}

func (s *chatHandler) DeleteRoom(ctx context.Context, req *pb.DeleteRoomRequest) (*pb.DeleteRoomResponse, error) {
	err := s.conversationUsecases.DeleteRoom(ctx, uint64(req.ConversationId))
	if err != nil {
		return nil, err
	}
	event := &messagepb.DeleteRoom{
		Type: constants.DELETE_ROOM_TYPE,
		Data: &messagepb.DeleteRoomData{
			RoomId: strconv.Itoa(int(req.ConversationId)),
		},
	}
	s.eventChannel <- event
	return &pb.DeleteRoomResponse{
		Success: true,
	}, nil
}

func (s *chatHandler) GetSetting(ctx context.Context, req *pb.GetSettingRequest) (*pb.GetSettingResponse, error) {
	conversation, participant, err := s.conversationUsecases.GetSetting(ctx, uint64(req.Id))
	if err != nil {
		return nil, err
	}
	return &pb.GetSettingResponse{
		Id:            conversation.ID,
		Type:          pb.ConversationTypeEnum(conversation.Type),
		Name:          conversation.Name,
		IsBroadcast:   conversation.IsBroadcast,
		ForbidForward: conversation.ForbidForward,
		IsDeleted:     conversation.IsDeleted,
		MuteNotif:     participant.MuteNotif,
	}, nil
}

func (s *chatHandler) GetFiles(ctx context.Context, req *pb.GetFilesRequest) (*pb.GetFilesResponse, error) {
	files, total, err := s.messageUsecases.GetMessageByTypeWithPagination(ctx, int64(req.Page), int64(req.Size), uint64(req.ConversationId), []string{"video", "audio", "pdf"})
	if err != nil {
		return nil, err
	}

	var pbMessage []*pb.Message
	for _, file := range files {
		pbMessage = append(pbMessage, &pb.Message{
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

	return &pb.GetFilesResponse{
		Messages: pbMessage,
		Total:    total,
	}, nil
}

func (s *chatHandler) GetImages(ctx context.Context, req *pb.GetImagesRequest) (*pb.GetImagesResponse, error) {
	files, total, err := s.messageUsecases.GetMessageByTypeWithPagination(ctx, int64(req.Page), int64(req.Size), uint64(req.ConversationId), []string{"image"})
	if err != nil {
		return nil, err
	}

	var pbMessage []*pb.Message
	for _, file := range files {
		pbMessage = append(pbMessage, &pb.Message{
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

	return &pb.GetImagesResponse{
		Messages: pbMessage,
		Total:    total,
	}, nil
}

func (s *chatHandler) GetLinks(ctx context.Context, req *pb.GetLinksRequest) (*pb.GetLinksResponse, error) {
	files, total, err := s.messageUsecases.GetMessageByTypeWithPagination(ctx, int64(req.Page), int64(req.Size), uint64(req.ConversationId), []string{"link"})
	if err != nil {
		return nil, err
	}

	var pbMessage []*pb.Message
	for _, file := range files {
		pbMessage = append(pbMessage, &pb.Message{
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

	return &pb.GetLinksResponse{
		Messages: pbMessage,
		Total:    total,
	}, nil
}

func (s *chatHandler) ReactionMessage(ctx context.Context, req *pb.ReactionMessageRequest) (*pb.ReactionMessageResponse, error) {
	fmt.Println("ReactionMessage", req.MessageId, req.Reaction)
	fmt.Println("gothere")
	_, err := s.messageReactionUsecases.CreateMessageReaction(ctx, uint64(req.MessageId), req.Reaction)
	fmt.Println("gothere2")

	if err != nil {
		return nil, err
	}

	return &pb.ReactionMessageResponse{
		Success: true,
	}, nil
}

func (s *chatHandler) Start() {
	for i := 0; i < s.numberOfWorker; i++ {
		redisClientWorker := NewRedisClientWorker(s, s.redisClient, s.logger)
		go redisClientWorker.serve()
	}
}
