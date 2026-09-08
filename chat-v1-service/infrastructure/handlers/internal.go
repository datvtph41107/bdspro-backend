package handlers

import (
	_dto "common/domain/dto"
	"context"
	"time"

	"chat/models"
)

type redisClient interface {
	Publish(ctx context.Context, channel string, message string) error
}
type conversationUsecases interface {
	CreateConversation(ctx context.Context, conversation *models.ConversationModel, members []uint64) (*models.ConversationModel, error)
	GetConversations(ctx context.Context, limit, offset int64, conversationType *int32, keyword string) ([]*models.ConversationModel, error)
	GetUnreadCount(ctx context.Context) (int64, error)

	EditChannel(ctx context.Context, conversationId uint64, isBroadcast, forbidForward *bool, name *string, backgroundImageId *uint64) error
	CreateOrgChannel(ctx context.Context, conversation *models.ConversationModel) (*models.ConversationModel, error)
	Search(ctx context.Context, keyword string, searchType int, limit, offset int64) ([]*models.ConversationModel, []*models.MessageModel, error)

	GetConversationById(ctx context.Context, conversationId uint64, includesMember bool) (*models.ConversationModel, error)
	LeaveRoom(ctx context.Context, conversationId uint64) error
	DeleteRoom(ctx context.Context, conversationId uint64) error
	GetSetting(ctx context.Context, conversationId uint64) (*models.ConversationModel, *models.ParticipantModel, error)
	ValidateConversationAndCurrentUser(ctx context.Context, conversationId uint64) error
	GetMembersOfRoom(ctx context.Context, conversationId uint64) ([]uint64, error)
	SetBackgroundImage(ctx context.Context, conversationId uint64, backgroundImageId uint64) (*models.ConversationModel, error)
}

type messageUsecases interface {
	SendMessage(ctx context.Context, message *models.MessageModel) (*models.MessageModel, error)
	SendToReceiver(ctx context.Context, message *models.MessageModel, receiverID *uint64) (*models.MessageModel, error)
	GetMessages(ctx context.Context, pagable _dto.Pagable, conversationId uint64, sort string, fromDate, toDate *time.Time) ([]*models.MessageModel, error)
	ForwardMessage(ctx context.Context, messageId uint64, conversationIds []uint64) ([]*models.MessageModel, error)
	CreateSystemMessage(ctx context.Context, conversationId uint64, content string) (*models.MessageModel, error)
	PinMessage(ctx context.Context, messageId uint64, conversationId uint64, pin bool) (*models.MessageModel, error)
	RecallMessage(ctx context.Context, messageId uint64) (*models.MessageModel, error)
	GetPinnedMessages(ctx context.Context, conversationId uint64) ([]*models.MessageModel, error)
	EditMessage(ctx context.Context, messageId uint64, content string) (*models.MessageModel, error)
	DeleteViolation(ctx context.Context, messageId uint64) (*models.MessageModel, error)
	DeleteConversation(ctx context.Context, conversationId uint64) (uint64, error)

	GetMessageByTypeWithPagination(ctx context.Context, page, size int64, conversationId uint64, types []string) ([]*models.MessageModel, int64, error)
}

type participantUsecases interface {
	RemoveUser(ctx context.Context, conversationId uint64, userId uint64) error
	MuteConversation(ctx context.Context, conversationId uint64, mute bool) error
	GetListParticipant(ctx context.Context, conversationId uint64) ([]*models.ParticipantModel, int64, error)
	AddUsers(ctx context.Context, conversationId uint64, userIds []uint64) (int, time.Time, error)
	InternalGetListParticipant(ctx context.Context, conversationId string) ([]*models.ParticipantModel, error)
}

type readReceptUsecases interface {
	MarkAsRead(ctx context.Context, conversationId uint64, messageIds []uint64) error
	GetReadReceipt(ctx context.Context, limit, offset int64, messageId uint64) ([]*models.ReadReceptModel, int64, error)
	GetReadReceiptsByMessageIDsForCurrentUser(ctx context.Context, messageIDs []uint64) ([]*models.ReadReceptModel, error)
}

type messageReactionUsecases interface {
	CreateMessageReaction(ctx context.Context, messageId uint64, reaction string) (*models.MessageReactionModel, error)
}

type backgroundImageUsecases interface {
	CreateBackgroundImage(ctx context.Context, backgroundImage *models.BackgroundImageModel) (*models.BackgroundImageModel, error)
	GetBackgroundImages(ctx context.Context, page, size *int64) ([]*models.BackgroundImageModel, int64, error)
	UpdateBackgroundImage(ctx context.Context, backgroundImage *models.BackgroundImageModel) (*models.BackgroundImageModel, error)
	DeleteBackgroundImage(ctx context.Context, id uint64) error
}
