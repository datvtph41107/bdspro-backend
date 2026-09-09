package handlers

import (
	_dto "common/domain/dto"
	_errors "common/errors"
	_utils "common/utils"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	chatpb "pb/types/chat"
	sharepb "pb/types/shared"

	"chat/config"
	"chat/infrastructure/client"
	"chat/infrastructure/delivery/errors"
	"chat/infrastructure/mapper"
	"chat/infrastructure/redis"
	"chat/internal/constants"
	"chat/internal/usecases"
	"chat/models"
	"chat/utils"

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
	crmClient               *client.CrmClient
	bdsproClient            *client.BdsproClient
	conversationMapper      *mapper.ConversationMapper
	messageMapper           *mapper.MessageMapper
	backgroundMapper        *mapper.BackgroundMapper
	syncProvider            *_utils.SyncUtil
	messageRepository       usecases.Repository
	systemMessageUsecase    *usecases.SystemMessageUsecase
	typingUsecases          *usecases.TypingUsecase
}

func (s *chatHandler) putTimestampFromTimePtr(ctx context.Context, key string, t *time.Time) {
	ts := int64(0)
	if t != nil {
		ts = t.UnixMilli()
	}
	if ts == 0 {
		ts = 1
	}
	_ = s.syncProvider.PutTimestamp(ctx, key, ts)
}

func (s *chatHandler) attachReadAtForCurrentUser(ctx context.Context, messages []*models.MessageModel) error {
	if len(messages) == 0 {
		return nil
	}
	messageIDs := make([]uint64, 0, len(messages))
	for _, msg := range messages {
		if msg == nil || msg.ID == 0 {
			continue
		}
		messageIDs = append(messageIDs, msg.ID)
	}
	if len(messageIDs) == 0 {
		return nil
	}
	receipts, err := s.readReceptUsecases.GetReadReceiptsByMessageIDsForCurrentUser(ctx, messageIDs)
	if err != nil {
		return err
	}
	readAtByMessageID := make(map[uint64]time.Time, len(receipts))
	for _, receipt := range receipts {
		if receipt == nil {
			continue
		}
		readAtByMessageID[receipt.MessageID] = receipt.ReadAt
	}
	for _, msg := range messages {
		if msg == nil {
			continue
		}
		if readAt, ok := readAtByMessageID[msg.ID]; ok {
			tmp := readAt
			msg.ReadAt = &tmp
		}
	}
	return nil
}

func NewChatHandler(
	conversationUsecases conversationUsecases,
	messageUsecases messageUsecases,
	readReceptUsecases readReceptUsecases,
	participantUsecases participantUsecases,
	messageReactionUsecases messageReactionUsecases,
	userClient *client.UserClient,
	crmClient *client.CrmClient,
	bdsproClient *client.BdsproClient,
	conversationMapper *mapper.ConversationMapper,
	messageMapper *mapper.MessageMapper,
	backgroundMapper *mapper.BackgroundMapper,
	syncProvider *_utils.SyncUtil,
	messageRepository usecases.Repository,
	systemMessageUsecase *usecases.SystemMessageUsecase,
) *chatHandler {
	redisCli := redis.NewRedisClient()
	return &chatHandler{
		conversationUsecases:    conversationUsecases,
		messageUsecases:         messageUsecases,
		readReceptUsecases:      readReceptUsecases,
		participantUsecases:     participantUsecases,
		messageReactionUsecases: messageReactionUsecases,
		logger:                  flogging.MustGetLogger("chat_handler"),
		numberOfWorker:          config.AppProperties.Worker.Number,
		maxParticipant:          config.AppProperties.MaxParticipant,
		redisClient:             redisCli,
		typingUsecases:          usecases.NewTypingUsecase(conversationUsecases, redisCli),
		eventChannel:            make(chan any),
		userClient:              userClient,
		crmClient:               crmClient,
		bdsproClient:            bdsproClient,
		conversationMapper:      conversationMapper,
		messageMapper:           messageMapper,
		backgroundMapper:        backgroundMapper,
		syncProvider:            syncProvider,
		messageRepository:       messageRepository,
		systemMessageUsecase:    systemMessageUsecase,
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
	if s.syncProvider != nil {
		ts := time.Now().UnixMilli()
		_ = s.syncProvider.PutTimestamp(ctx, s.syncProvider.GetKey(ctx, _utils.SyncKeyConversationDetail, model.ID), ts)
		_ = s.syncProvider.PutTimestamp(ctx, s.syncProvider.GetKey(ctx, _utils.SyncKeyConversationCommon, 0), ts)
		memberIds := append(utils.Int32SliceToUint64Slice(req.Members), uint64(req.CreatedBy))
		s.invalidateConversationsListSyncKeyForUsers(ctx, memberIds)
	}
	return &chatpb.CreateConversationResponse{
		ConversationId: int32(model.ID),
	}, nil
}

// @Summary Lấy danh sách đoạn hội thoại
// @Description Lấy danh sách đoạn hội thoại theo các tiêu chí được chỉ định (sync: gửi timestamp lần trước, unchanged=true nếu dùng cache)
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param conversation query chatpb.GetConversationsRequest true "Tiêu chí lấy danh sách đoạn hội thoại"
// @Success 200 {object} chatpb.GetConversationsResponse
// @Router /conversations [get]
func (s *chatHandler) GetConversations(ctx context.Context, req *chatpb.GetConversationsRequest) (*chatpb.GetConversationsResponse, error) {
	if s.syncProvider != nil && req != nil && req.Timestamp > 0 {
		key := s.syncProvider.GetKey(ctx, _utils.SyncKeyConversationsMe, 0)
		updated := s.syncProvider.HasUpdated(ctx, key, req.Timestamp)
		s.syncProvider.PutTimeRequest(ctx, key, req.Timestamp)
		if !updated {
			return &chatpb.GetConversationsResponse{Conversations: nil, Unchanged: true}, nil
		}
	}

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
	latestMessages := make([]*models.MessageModel, 0, len(conversations))
	for _, conversation := range conversations {
		if conversation == nil || conversation.LatestMessage.ID == 0 {
			continue
		}
		latestMessages = append(latestMessages, &conversation.LatestMessage)
	}
	if err = s.attachReadAtForCurrentUser(ctx, latestMessages); err != nil {
		return nil, err
	}

	pbConversations := s.conversationMapper.ConversationsToPb(conversations)
	s.userClient.MapUserProfile(ctx, pbConversations)

	ts := time.Now().UnixMilli()
	return &chatpb.GetConversationsResponse{
		Conversations: pbConversations,
		Timestamp:     ts,
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

	go func() {
		ctxClone := _utils.CloneContext(ctx)
		key := s.syncProvider.GetKeyWithoutMe(ctxClone, _utils.SyncKeyConversationCommon, message.ConversationID)
		s.syncProvider.PutTimestamp(ctxClone, key, message.UpdatedAt.UnixMilli())
	}()
	return s.ResponseMessage(ctx, message)
}

func (s *chatHandler) ResponseMessage(ctx context.Context, message *models.MessageModel) (*chatpb.SendMessageResponse, error) {
	pbMsg := s.messageMapper.MessageToPb(message)
	event := &chatpb.SendMessage{
		Type: constants.MESSAGE_TYPE,
		Data: pbMsg,
	}

	// Giống GetMessages: avatar, fullName, contactUser (tin contact) qua batch profile
	toEnrich := []*chatpb.Message{pbMsg}
	if pbMsg.ReplyMessage != nil {
		toEnrich = append(toEnrich, pbMsg.ReplyMessage)
	}
	s.userClient.MapAvatarAndNameToMessagePbList(ctx, toEnrich)
	s.bdsproClient.MapProductPbByIds(ctx, toEnrich)

	s.eventChannel <- event

	updateRoomEvent := &chatpb.UpdateRoomMessage{
		Type: constants.UPDATE_ROOM_TYPE,
		Data: &chatpb.UpdateRoomMessageData{
			RoomId: strconv.FormatUint(message.ConversationID, 10),
		},
	}
	s.eventChannel <- updateRoomEvent
	s.invalidateMessageTypeSyncKeys(ctx, message.ConversationID, string(message.ContentType))
	s.invalidateUnreadSyncKeysForReceivers(ctx, message.ConversationID, message.SenderID)
	return &chatpb.SendMessageResponse{
		MessageId: int32(message.ID),
		CreatedAt: message.CreatedAt.Local().Format(time.RFC3339),
	}, nil
}

// invalidateMessageTypeSyncKeys cập nhật timestamp khi có message file/image/link mới
func (s *chatHandler) invalidateMessageTypeSyncKeys(ctx context.Context, conversationID uint64, contentType string) {
	if s.syncProvider == nil {
		return
	}
	ct := strings.ToLower(contentType)
	ts := time.Now().UnixMilli()
	switch ct {
	case "image":
		_ = s.syncProvider.PutTimestamp(ctx, s.syncProvider.GetKey(ctx, _utils.SyncKeyConversationImages, conversationID), ts)
	case "link":
		_ = s.syncProvider.PutTimestamp(ctx, s.syncProvider.GetKey(ctx, _utils.SyncKeyConversationLinks, conversationID), ts)
	case "file", "video", "audio", "pdf":
		_ = s.syncProvider.PutTimestamp(ctx, s.syncProvider.GetKey(ctx, _utils.SyncKeyConversationFiles, conversationID), ts)
	}
	_ = s.syncProvider.PutTimestamp(ctx, s.syncProvider.GetKey(ctx, _utils.SyncKeyConversationCommon, 0), ts)
}

// invalidateUnreadSyncKeysForReceivers cập nhật timestamp unread của các thành viên nhận tin nhắn (trừ sender)
func (s *chatHandler) invalidateUnreadSyncKeysForReceivers(ctx context.Context, conversationID uint64, senderID uint64) {
	if s.syncProvider == nil {
		return
	}
	participants, _, err := s.participantUsecases.GetListParticipant(ctx, conversationID)
	if err != nil {
		return
	}
	ts := time.Now().UnixMilli()
	for _, p := range participants {
		if p.UserID != senderID {
			key := s.syncProvider.GetKeyForProfile(p.UserID, 0, _utils.SyncKeyChatUnreadMe)
			_ = s.syncProvider.PutTimestamp(ctx, key, ts)
			// Danh sách conversations cũng thay đổi (lastMessage)
			keyConv := s.syncProvider.GetKeyForProfile(p.UserID, 0, _utils.SyncKeyConversationsMe)
			_ = s.syncProvider.PutTimestamp(ctx, keyConv, ts)
		}
	}
}

// invalidateConversationsListSyncKeyForUsers cập nhật timestamp danh sách conversations và conversation common của các user
func (s *chatHandler) invalidateConversationsListSyncKeyForUsers(ctx context.Context, profileIds []uint64) {
	if s.syncProvider == nil {
		return
	}
	ts := time.Now().UnixMilli()
	for _, id := range profileIds {
		key := s.syncProvider.GetKeyForProfile(id, 0, _utils.SyncKeyConversationsMe)
		_ = s.syncProvider.PutTimestamp(ctx, key, ts)
		keyCommon := s.syncProvider.GetKeyForProfile(id, 0, _utils.SyncKeyConversationCommon)
		_ = s.syncProvider.PutTimestamp(ctx, keyCommon, ts)
	}
}

// invalidateConversationsListSyncKeyForParticipants cập nhật timestamp danh sách conversations của tất cả thành viên
func (s *chatHandler) invalidateConversationsListSyncKeyForParticipants(ctx context.Context, conversationID uint64) {
	if s.syncProvider == nil {
		return
	}
	participants, _, err := s.participantUsecases.GetListParticipant(ctx, conversationID)
	if err != nil {
		return
	}
	profileIds := make([]uint64, 0, len(participants))
	for _, p := range participants {
		profileIds = append(profileIds, p.UserID)
	}
	s.invalidateConversationsListSyncKeyForUsers(ctx, profileIds)
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
		return nil, fmt.Errorf("receiverId is required")
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
	s.invalidateMessageTypeSyncKeys(ctx, message.ConversationID, string(message.ContentType))
	s.invalidateUnreadSyncKeysForReceivers(ctx, message.ConversationID, message.SenderID)
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
	key := s.syncProvider.GetKeyWithoutMe(ctx, _utils.SyncKeyConversationCommon, uint64(req.ConversationId))
	updated, cacheExists := s.syncProvider.IsCacheChanged(ctx, key, req.Timestamp)
	fmt.Println("cacheExists", cacheExists)

	if !cacheExists {
		go func() {
			ctxClone := _utils.CloneContext(ctx)
			keyCommon := s.syncProvider.GetKeyWithoutMe(ctxClone, _utils.SyncKeyConversationCommon, uint64(req.ConversationId))
			latestAt, errLatest := s.messageRepository.GetLatestMessageCreatedAt(ctxClone, uint64(req.ConversationId))
			if errLatest == nil && latestAt != nil && !latestAt.IsZero() {
				_ = s.syncProvider.PutTimestamp(ctxClone, keyCommon, latestAt.UnixMilli())
			}
		}()
	}

	// timestamp == 0: lấy mới; != 0: check key cnv:c (last), nếu chưa có thay đổi thì trả ErrCacheUnchanged
	if !updated {
		return &chatpb.GetMessagesResponse{Messages: nil, Total: 0, Unchanged: true}, nil
	}

	// var limit, offset int64
	var sort string
	var fromDate, toDate *time.Time
	// if req.Size != nil {
	// 	limit = int64(*req.Size)
	// } else {
	// 	limit = 10
	// }

	// if req.Offset != nil {
	// 	offset = int64(*req.Offset)
	// } else {
	// 	offset = 0
	// }

	pagable := _dto.Pagable{
		Page: req.Page,
		Size: req.Size,
		Sort: "created_at desc",
	}

	// if req.Sort != nil {
	// 	sort = *req.Sort
	// } else {
	// 	sort = "desc"
	// }

	if req.Timestamp != 0 {
		tempFromDate := time.UnixMilli(req.Timestamp) // _utils.ParseStringToTime(strconv.FormatInt(req.Timestamp, 10))
		fromDate = &tempFromDate
	}

	if req.ToDate != nil {
		tempToDate := _utils.ParseStringToTime(*req.ToDate)
		toDate = tempToDate
	}

	messages, err := s.messageUsecases.GetMessages(ctx, pagable, uint64(req.ConversationId), sort, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	if err = s.attachReadAtForCurrentUser(ctx, messages); err != nil {
		return nil, err
	}

	var pbMessages []*chatpb.Message
	for _, msg := range messages {
		pbMessages = append(pbMessages, s.messageMapper.MessageToPb(msg))
	}

	s.userClient.MapAvatarAndNameToMessagePbList(ctx, pbMessages)
	// s.crmClient.MapContactToMessagePbList(ctx, pbMessages)
	s.bdsproClient.MapProductPbByIds(ctx, pbMessages)

	return &chatpb.GetMessagesResponse{
		Messages:  pbMessages,
		Total:     int64(len(pbMessages)),
		Unchanged: false,
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
	s.crmClient.MapContactToMessagePbList(ctx, pbMessages)
	for _, message := range messages {
		s.invalidateMessageTypeSyncKeys(ctx, message.ConversationID, string(message.ContentType))
		s.invalidateUnreadSyncKeysForReceivers(ctx, message.ConversationID, message.SenderID)
	}
	for _, message := range pbMessages {
		event := &chatpb.SendMessage{
			Type: constants.MESSAGE_TYPE,
			Data: message,
		}
		s.eventChannel <- event
	}

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
	if s.syncProvider != nil {
		ts := time.Now().UnixMilli()
		_ = s.syncProvider.PutTimestamp(ctx, s.syncProvider.GetKey(ctx, _utils.SyncKeyChatUnreadMe, 0), ts)
		s.invalidateConversationsListSyncKeyForUsers(ctx, []uint64{utils.GetCurrentUserID(ctx)})
	}
	return &chatpb.MarkAsReadResponse{
		Success: true,
	}, nil
}

// @Summary Ghim tin nhắn
// @Description Ghim hoặc bỏ ghim tin nhắn trong đoạn hội thoại
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param message body chatpb.PinMessageRequest true "Thông tin ghim tin nhắn"
// @Success 200 {object} chatpb.PinMessageResponse
// @Router /messages/pin [post]
func (s *chatHandler) PinMessage(ctx context.Context, req *chatpb.PinMessageRequest) (*chatpb.PinMessageResponse, error) {
	// Existing logic
	_, err := s.messageUsecases.PinMessage(ctx, uint64(req.MessageId), uint64(req.ConversationId), req.DoPin)
	if err != nil {
		return nil, err
	}

	// content := "Đã bỏ ghim tin nhắn"
	// if req.DoPin {
	// 	content = "Đã ghim tin nhắn"
	// }
	message, err := s.systemMessageUsecase.PinSystemMessage(ctx, uint64(req.ConversationId), uint64(req.MessageId), req.DoPin)
	if err != nil {
		return nil, err
	}

	pbMsg := s.messageMapper.MessageToPb(message)
	event := &chatpb.PinMessage{
		Type: constants.PIN_MESSAGE_TYPE,
		Data: pbMsg,
	}

	user, _ := s.userClient.GetUserById(ctx, utils.GetCurrentUserID(ctx))
	if user != nil {
		event.Data.FullName = &user.FullName
		event.Data.Avatar = &user.Avatar
	}

	s.eventChannel <- event

	if s.syncProvider != nil {
		ts := time.Now().UnixMilli()
		_ = s.syncProvider.PutTimestamp(ctx, s.syncProvider.GetKeyWithoutMe(ctx, _utils.SyncKeyConversationPinned, uint64(req.ConversationId)), ts)
		_ = s.syncProvider.PutTimestamp(ctx, s.syncProvider.GetKeyWithoutMe(ctx, _utils.SyncKeyConversationCommon, 0), ts)
	}

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
	s.invalidateMessageTypeSyncKeys(ctx, message.ConversationID, string(message.ContentType))

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
	if s.syncProvider != nil {
		ts := time.Now().UnixMilli()
		_ = s.syncProvider.PutTimestamp(ctx, s.syncProvider.GetKey(ctx, _utils.SyncKeyConversationDetail, uint64(req.ConversationId)), ts)
		_ = s.syncProvider.PutTimestamp(ctx, s.syncProvider.GetKey(ctx, _utils.SyncKeyConversationCommon, 0), ts)
		s.invalidateConversationsListSyncKeyForUsers(ctx, []uint64{uint64(req.TargetUserId)})
	}
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
	if s.syncProvider != nil {
		key := s.syncProvider.GetKey(ctx, _utils.SyncKeyConversationDetail, uint64(req.ConversationId))
		ts := time.Now().UnixMilli()
		_ = s.syncProvider.PutTimestamp(ctx, key, ts)
		_ = s.syncProvider.PutTimestamp(ctx, s.syncProvider.GetKey(ctx, _utils.SyncKeyConversationSetting, uint64(req.ConversationId)), ts)
		_ = s.syncProvider.PutTimestamp(ctx, s.syncProvider.GetKeyWithoutMe(ctx, _utils.SyncKeyConversationCommon, 0), ts)
		s.invalidateConversationsListSyncKeyForParticipants(ctx, uint64(req.ConversationId))
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
	if s.syncProvider != nil {
		ts := time.Now().UnixMilli()
		_ = s.syncProvider.PutTimestamp(ctx, s.syncProvider.GetKey(ctx, _utils.SyncKeyConversationDetail, channel.ID), ts)
		_ = s.syncProvider.PutTimestamp(ctx, s.syncProvider.GetKeyWithoutMe(ctx, _utils.SyncKeyConversationCommon, 0), ts)
		s.invalidateConversationsListSyncKeyForParticipants(ctx, channel.ID)
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
	if s.syncProvider != nil {
		ts := time.Now().UnixMilli()
		_ = s.syncProvider.PutTimestamp(ctx, s.syncProvider.GetKey(ctx, _utils.SyncKeyConversationSetting, uint64(req.ConversationId)), ts)
		_ = s.syncProvider.PutTimestamp(ctx, s.syncProvider.GetKeyWithoutMe(ctx, _utils.SyncKeyConversationCommon, 0), ts)
		s.invalidateConversationsListSyncKeyForUsers(ctx, []uint64{utils.GetCurrentUserID(ctx)})
	}

	return &chatpb.MuteConversationResponse{
		Success: true,
	}, nil
}

// @Summary Lấy timestamps các loại data của conversation
// @Description API check sync - lấy timestamp từ Redis cho từng key, FE so sánh với cache và chỉ gọi API tương ứng khi có thay đổi. Keys: detail, files, links, images, members, setting, pinned, messages, last (cnv:c - cập nhật khi có event conversation)
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Conversation ID"
// @Success 200 {object} chatpb.ConversationTimestampsResponse "Map key -> timestamp (ms)"
// @Router /conversation/{id}/timestamp [get]
func (s *chatHandler) GetConversationTimestamps(ctx context.Context, req *sharepb.IdRequest) (*chatpb.ConversationTimestampsResponse, error) {
	if req == nil || req.Id == 0 {
		return nil, _errors.BadRequestException("conversation id is required")
	}

	if s.syncProvider == nil {
		return &chatpb.ConversationTimestampsResponse{Timestamps: map[string]int64{}}, nil
	}

	convID := req.Id
	redisKeys := []string{
		s.syncProvider.GetKeyWithoutMe(ctx, _utils.SyncKeyConversationDetail, convID),
		s.syncProvider.GetKeyWithoutMe(ctx, _utils.SyncKeyConversationFiles, convID),
		s.syncProvider.GetKeyWithoutMe(ctx, _utils.SyncKeyConversationLinks, convID),
		s.syncProvider.GetKeyWithoutMe(ctx, _utils.SyncKeyConversationImages, convID),
		s.syncProvider.GetKeyWithoutMe(ctx, _utils.SyncKeyConversationMembers, convID),
		s.syncProvider.GetKey(ctx, _utils.SyncKeyConversationSetting, convID),
		s.syncProvider.GetKeyWithoutMe(ctx, _utils.SyncKeyConversationPinned, convID),
		s.syncProvider.GetKey(ctx, _utils.SyncKeyConversationMessages, convID),
		s.syncProvider.GetKeyWithoutMe(ctx, _utils.SyncKeyConversationCommon, convID),
	}

	vals, err := s.syncProvider.MGet(ctx, redisKeys)
	if err != nil {
		return nil, fmt.Errorf("%s", err.Error())
	}

	names := []string{
		"detail", "files", "links", "images", "members", "setting", "pinned", "messages", "last",
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
	key := s.syncProvider.GetKeyWithoutMe(ctx, _utils.SyncKeyConversationDetail, uint64(req.ConversationId))
	updated, cacheExists := s.syncProvider.IsCacheChanged(ctx, key, req.Timestamp)
	if !updated {
		return &chatpb.GetConversationByIdResponse{}, nil
	}

	req.IncludeMembers = true
	conversation, err := s.conversationUsecases.GetConversationById(ctx, uint64(req.ConversationId), req.IncludeMembers)
	if err != nil {
		return nil, err
	}

	result := s.conversationMapper.ConversationToPbGetConversationById(conversation, req.IncludeMembers)

	s.userClient.MapAvatarAndNameToPb(ctx, result)

	if !cacheExists {
		go func() {
			ctxClone := _utils.CloneContext(ctx)
			s.syncProvider.PutTimestamp(ctxClone, key, conversation.UpdatedAt.UnixMilli())
		}()
	}

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
	key := s.syncProvider.GetKeyWithoutMe(ctx, _utils.SyncKeyConversationMembers, uint64(req.ConversationId))
	updated, cacheExists := s.syncProvider.IsCacheChanged(ctx, key, req.Timestamp)

	if !cacheExists {
		go func() {
			ctxClone := _utils.CloneContext(ctx)
			latestJoinedAt, errLatest := s.messageRepository.GetLatestParticipantJoinedAt(ctxClone, uint64(req.ConversationId))
			if errLatest == nil && latestJoinedAt != nil && !latestJoinedAt.IsZero() {
				_ = s.syncProvider.PutTimestamp(ctxClone, key, latestJoinedAt.UnixMilli())
			}
		}()
	}

	if !updated {
		return &chatpb.GetListParticipantResponse{Data: []*chatpb.Participant{}, Total: 0}, nil
	}
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
			JoinedAt:  _utils.FormatTimeToString(&participant.JoinedAt),
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
	key := s.syncProvider.GetKeyWithoutMe(ctx, _utils.SyncKeyConversationPinned, uint64(req.ConversationId))
	updated, cacheExists := s.syncProvider.IsCacheChanged(ctx, key, req.Timestamp)

	if !cacheExists {
		go func() {
			ctxClone := _utils.CloneContext(ctx)
			latestAt, errLatest := s.messageRepository.GetLatestPinnedMessageUpdatedAt(ctxClone, uint64(req.ConversationId))
			if errLatest == nil {
				s.putTimestampFromTimePtr(ctxClone, key, latestAt)
			}
		}()
	}

	if !updated {
		return &chatpb.GetMessagesResponse{Messages: []*chatpb.Message{}, Total: 0}, nil
	}
	messages, err := s.messageUsecases.GetPinnedMessages(ctx, uint64(req.ConversationId))
	if err != nil {
		return nil, err
	}

	var pbMessages []*chatpb.Message
	for _, msg := range messages {
		pbMessages = append(pbMessages, s.messageMapper.MessageToPb(msg))
	}

	s.userClient.MapAvatarAndNameToMessagePbList(ctx, pbMessages)
	s.crmClient.MapContactToMessagePbList(ctx, pbMessages)
	s.bdsproClient.MapProductPbByIds(ctx, pbMessages)

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
	if len(req.UserIds) == 0 {
		return &chatpb.AddUserResponse{
			Success:    false,
			TotalAdded: 0,
		}, nil
	}

	totalAdded, joinedAt, err := s.participantUsecases.AddUsers(ctx, req.ConversationId, req.UserIds)
	if err != nil {
		return nil, err
	}

	if totalAdded > 0 {
		if err := s.systemMessageUsecase.AddMemberSystemMessage(ctx, req.ConversationId, req.UserIds); err != nil {
			s.logger.Errorf("Failed to create system message for add users: %v", err)
		}
	}

	users, _ := s.userClient.Client.GetProfileByIds(ctx, req.UserIds)
	userMap := make(map[uint64]*sharepb.ProfileItem)
	if users != nil {
		for _, user := range users.Profiles {
			userMap[user.Id] = user
		}
	}

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
	if s.syncProvider != nil {
		ts := joinedAt.UnixMilli()
		_ = s.syncProvider.PutTimestamp(ctx, s.syncProvider.GetKeyWithoutMe(ctx, _utils.SyncKeyConversationCommon, req.ConversationId), ts)
		_ = s.syncProvider.PutTimestamp(ctx, s.syncProvider.GetKeyWithoutMe(ctx, _utils.SyncKeyConversationMembers, req.ConversationId), ts)
		s.invalidateConversationsListSyncKeyForUsers(ctx, req.UserIds)
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
			UserId: strconv.FormatUint(utils.GetCurrentUserID(ctx), 10),
			RoomId: strconv.FormatUint(uint64(req.ConversationId), 10),
		},
	}
	s.eventChannel <- event
	if s.syncProvider != nil {
		ts := time.Now().UnixMilli()
		_ = s.syncProvider.PutTimestamp(ctx, s.syncProvider.GetKey(ctx, _utils.SyncKeyConversationDetail, uint64(req.ConversationId)), ts)
		s.invalidateConversationsListSyncKeyForUsers(ctx, []uint64{utils.GetCurrentUserID(ctx)})
	}
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
	message, err := s.messageUsecases.DeleteViolation(ctx, uint64(req.MessageId))
	if err != nil {
		return nil, err
	}
	if message != nil {
		s.invalidateMessageTypeSyncKeys(ctx, message.ConversationID, string(message.ContentType))
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
	if s.syncProvider != nil {
		ts := time.Now().UnixMilli()
		_ = s.syncProvider.PutTimestamp(ctx, s.syncProvider.GetKey(ctx, _utils.SyncKeyConversationDetail, uint64(req.ConversationId)), ts)
		_ = s.syncProvider.PutTimestamp(ctx, s.syncProvider.GetKey(ctx, _utils.SyncKeyConversationSetting, uint64(req.ConversationId)), ts)
		s.invalidateConversationsListSyncKeyForParticipants(ctx, uint64(req.ConversationId))
	}
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
	key := s.syncProvider.GetKey(ctx, _utils.SyncKeyConversationSetting, req.Id)
	updated := s.syncProvider.HasUpdated(ctx, key, req.Timestamp)
	s.syncProvider.PutTimeRequest(ctx, key, req.Timestamp)
	if !updated {
		return &chatpb.GetSettingResponse{}, nil
	}

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
	key := s.syncProvider.GetKeyWithoutMe(ctx, _utils.SyncKeyConversationFiles, uint64(req.ConversationId))
	updated, cacheExists := s.syncProvider.IsCacheChanged(ctx, key, req.Timestamp)

	if !cacheExists {
		go func() {
			ctxClone := _utils.CloneContext(ctx)
			latestAt, errLatest := s.messageRepository.GetLatestMessageUpdatedAtByTypes(ctxClone, uint64(req.ConversationId), []string{"FILE", "video", "audio", "pdf"})
			if errLatest == nil {
				s.putTimestampFromTimePtr(ctxClone, key, latestAt)
			}
		}()
	}

	if !updated {
		return &chatpb.GetFilesResponse{Messages: []*chatpb.Message{}, Total: 0}, nil
	}
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
			CreatedAt:   _utils.FormatTimeToString(&file.CreatedAt),
			UpdatedAt:   _utils.FormatTimeToString(&file.UpdatedAt),
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
	key := s.syncProvider.GetKeyWithoutMe(ctx, _utils.SyncKeyConversationImages, uint64(req.ConversationId))
	updated, cacheExists := s.syncProvider.IsCacheChanged(ctx, key, req.Timestamp)

	if !cacheExists {
		go func() {
			ctxClone := _utils.CloneContext(ctx)
			latestAt, errLatest := s.messageRepository.GetLatestMessageUpdatedAtByTypes(ctxClone, uint64(req.ConversationId), []string{"IMAGE"})
			if errLatest == nil {
				s.putTimestampFromTimePtr(ctxClone, key, latestAt)
			}
		}()
	}

	if !updated {
		return &chatpb.GetImagesResponse{Data: []*chatpb.Message{}, Total: 0}, nil
	}
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
			CreatedAt:   _utils.FormatTimeToString(&file.CreatedAt),
			UpdatedAt:   _utils.FormatTimeToString(&file.UpdatedAt),
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
	key := s.syncProvider.GetKeyWithoutMe(ctx, _utils.SyncKeyConversationLinks, uint64(req.ConversationId))
	updated, cacheExists := s.syncProvider.IsCacheChanged(ctx, key, req.Timestamp)

	if !cacheExists {
		go func() {
			ctxClone := _utils.CloneContext(ctx)
			latestAt, errLatest := s.messageRepository.GetLatestMessageUpdatedAtByTypes(ctxClone, uint64(req.ConversationId), []string{"LINK"})
			if errLatest == nil {
				s.putTimestampFromTimePtr(ctxClone, key, latestAt)
			}
		}()
	}

	if !updated {
		return &chatpb.GetLinksResponse{Messages: []*chatpb.Message{}, Total: 0}, nil
	}
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
			CreatedAt:   _utils.FormatTimeToString(&file.CreatedAt),
			UpdatedAt:   _utils.FormatTimeToString(&file.UpdatedAt),
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
			CreatedAt: _utils.FormatTimeToString(&model.CreatedAt),
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

// @Summary Cập nhật trạng thái đang gõ
// @Description isTyping=true INCR, false DECR (Redis key theo phòng + user); fanout WebSocket type typing — data.count là giá trị Redis; không gửi notification khi nhận offline.
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param message body chatpb.StateTypingRequest true "conversationId và isTyping"
// @Success 200 {object} chatpb.StateTypingResponse
// @Router /typing [post]
func (s *chatHandler) GetStateTyping(ctx context.Context, req *chatpb.StateTypingRequest) (*chatpb.StateTypingResponse, error) {
	if req == nil {
		return nil, _errors.BadRequestException("request is required")
	}
	counter, err := s.typingUsecases.ApplyTypingState(ctx, req.ConversationId, req.IsTyping)
	if err != nil {
		return nil, err
	}
	userID := utils.GetCurrentUserID(ctx)
	event := &chatpb.TypingIndicator{
		Type: constants.TYPING_INDICATOR_TYPE,
		Data: &chatpb.TypingIndicatorData{
			RoomId:   strconv.FormatUint(req.ConversationId, 10),
			UserId:   strconv.FormatUint(userID, 10),
			IsTyping: req.IsTyping,
			Count:    counter,
		},
	}
	user, _ := s.userClient.GetUserById(ctx, userID)
	if user != nil {
		event.Data.FullName = &user.FullName
		event.Data.Avatar = &user.Avatar
	}
	s.eventChannel <- event
	return &chatpb.StateTypingResponse{Success: true, Count: counter}, nil
}

// @Summary Lấy số lượng tin nhắn chưa đọc
// @Description Lấy số lượng tin nhắn chưa đọc của người dùng hiện tại (sync: gửi timestamp lần trước, unchanged=true nếu dùng cache)
// @Tags Chat
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param timestamp query int false "Client timestamp lần trước; 0 = luôn lấy từ DB"
// @Success 200 {object} chatpb.UnreadCountResponse
// @Router /unread-count [get]
func (s *chatHandler) UnreadCount(ctx context.Context, req *chatpb.UnreadCountRequest) (*chatpb.UnreadCountResponse, error) {
	if s.syncProvider != nil && req != nil && req.Timestamp > 0 {
		key := s.syncProvider.GetKey(ctx, _utils.SyncKeyChatUnreadMe, 0)
		updated := s.syncProvider.HasUpdated(ctx, key, req.Timestamp)
		s.syncProvider.PutTimeRequest(ctx, key, req.Timestamp)
		if !updated {
			return &chatpb.UnreadCountResponse{UnreadCount: 0, Unchanged: true}, nil
		}
	}
	unreadCount, err := s.conversationUsecases.GetUnreadCount(ctx)
	if err != nil {
		return nil, err
	}
	ts := time.Now().UnixMilli()
	return &chatpb.UnreadCountResponse{UnreadCount: unreadCount, Timestamp: ts}, nil
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
	if s.syncProvider != nil {
		ts := conversation.UpdatedAt.UnixMilli()
		_ = s.syncProvider.PutTimestamp(ctx, s.syncProvider.GetKeyWithoutMe(ctx, _utils.SyncKeyConversationDetail, uint64(req.ConversationId)), ts)
		_ = s.syncProvider.PutTimestamp(ctx, s.syncProvider.GetKey(ctx, _utils.SyncKeyConversationsMe, uint64(req.ConversationId)), ts)
		s.invalidateConversationsListSyncKeyForParticipants(ctx, uint64(req.ConversationId))
	}
	return &chatpb.SetBackgroundImageResponse{
		Success: true,
	}, nil
}
