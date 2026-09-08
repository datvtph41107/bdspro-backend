package repository

// type repository interface {
// 	CreateConversation(ctx context.Context, conversation *models.ConversationModel) (*models.ConversationModel, error)
// 	GetConversations(ctx context.Context, limit, offset int64, conversationType *int32, userId uint64, keyword string) ([]*models.ConversationModel, error)
// 	GetConversationByID(ctx context.Context, conversationID uint64) (*models.ConversationModel, error)
// 	UpdateConversation(ctx context.Context, id uint64, isBoardCast, forbidForward bool, name string, backgroundImageId *uint64) error
// 	GetMessageOrConversation(ctx context.Context, keyword string, searchType int, limit int64, offset int64) ([]*models.ConversationModel, []*models.MessageModel, error)
// 	GetConversationByIDWithMember(ctx context.Context, conversationID uint64) (*models.ConversationModel, error)
// 	DeleteConversation(ctx context.Context, conversationID uint64) error
// 	GetMessageByTypeWithPagination(ctx context.Context, page, size int64, conversationId uint64, types []string) ([]*models.MessageModel, int64, error)

// 	CreateMessage(ctx context.Context, message *models.MessageModel) (*models.MessageModel, error)
// 	GetMessages(ctx context.Context, limit, offset int64, conversationId uint64, sort string, fromDate, toDate *time.Time) ([]*models.MessageModel, error)
// 	FindMessageByID(ctx context.Context, messageId uint64) (*models.MessageModel, error)
// 	UpdatePinMessage(ctx context.Context, messageID uint64, isPin bool) (*models.MessageModel, error)
// 	UpdateRecallMessage(ctx context.Context, messageID uint64, isRecall bool) (*models.MessageModel, error)
// 	GetPinnedMessages(ctx context.Context, conversationID uint64) ([]*models.MessageModel, error)
// 	UpdateMessage(ctx context.Context, messageID uint64, content string) (*models.MessageModel, error)
// 	DeleteViolation(ctx context.Context, messageID uint64) error

// 	ExistParticipant(ctx context.Context, conversationID uint64, userID uint64) (bool, error)
// 	FindParticipant(ctx context.Context, conversationID uint64, userID uint64) (*models.ParticipantModel, error)
// 	RemoveParticipant(ctx context.Context, conversationID uint64, userID uint64) error
// 	UpdateMuteNotification(ctx context.Context, conversationID uint64, userID uint64, isMute bool) error
// 	GetListParticipant(ctx context.Context, conversationId uint64) ([]*models.ParticipantModel, int64, error)
// 	AddParticipant(ctx context.Context, participant *models.ParticipantModel) error
// 	HasExactConversationWithUserIDs(ctx context.Context, userIDs []uint64) (*models.ConversationModel, error)
// 	GetConversationWithReceiverID(ctx context.Context, currentUserId uint64, receiverID uint64) (*models.ConversationModel, error)

// 	CreateBatchRead(ctx context.Context, models []*models.ReadReceptModel) error
// 	GetReadReceipts(ctx context.Context, limit, offset int64, messageID uint64) ([]*models.ReadReceptModel, int64, error)
// 	GetReadReceiptsByUserId(ctx context.Context, userId uint64) ([]*models.ReadReceptModel, error)

// 	CreateMessageReaction(ctx context.Context, messageReaction *models.MessageReactionModel) (*models.MessageReactionModel, error)

// 	GetBackgroundImages(ctx context.Context, page, size *int64) ([]*models.BackgroundImageModel, int64, error)
// 	CreateBackgroundImage(ctx context.Context, model *models.BackgroundImageModel) (*models.BackgroundImageModel, error)
// 	GetBackgroundImageByID(ctx context.Context, id string) (*models.BackgroundImageModel, error)
// 	UpdateBackgroundImage(ctx context.Context, model *models.BackgroundImageModel) (*models.BackgroundImageModel, error)
// 	DeleteBackgroundImage(ctx context.Context, id uint64) error
// 	SetBackgroundImage(ctx context.Context, conversationId uint64, backgroundImageId uint64) error

// 	GetRoomMemberIDs(ctx context.Context, conversationId uint64) ([]uint64, error)
// }
