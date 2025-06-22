package orm

import (
	"context"
	"time"

	"github.com/tuandq3/go-chat-service/models"
	"gorm.io/gorm"
)

var GET_CONVERSATIONS_QUERY = `
SELECT 
	c.*,
	m.id              AS latest_message_id,
	m.content         AS latest_message_content,
	m.content_type    AS latest_message_content_type,
	m.file_url        AS latest_message_file_url,
	m.forward_from    AS latest_message_forward_from,
	m.sender_id       AS latest_message_sender_id,
	m.created_at      AS latest_message_created_at,
	m.is_deleted      AS latest_message_is_deleted,
	m.pined           AS latest_message_pined,
	m.recall          AS latest_message_recall,
	COALESCE(unread.unread_count, 0) AS unread_count
FROM conversations c
LEFT JOIN LATERAL (
	SELECT *
	FROM messages
	WHERE conversation_id = c.id AND is_deleted = FALSE
	ORDER BY created_at DESC
	LIMIT 1
) m ON true
LEFT JOIN participants p ON p.conversation_id = c.id
LEFT JOIN (
	SELECT 
		m.conversation_id,
		COUNT(m.id) AS unread_count
	FROM messages m
	LEFT JOIN read_recepts rr ON rr.message_id = m.id AND rr.user_id = $1
	WHERE m.is_deleted = FALSE AND rr.message_id IS NULL
	GROUP BY m.conversation_id
) unread ON unread.conversation_id = c.id
WHERE c.type = $2 AND p.user_id = $3 AND LOWER((TRIM(c.name))) LIKE '%' || LOWER((TRIM($4))) || '%' AND c.is_deleted = FALSE
ORDER BY c.created_at DESC
LIMIT $5 OFFSET $6
`

func (r *gormRepository) CreateConversation(ctx context.Context, conversation *models.ConversationModel) (*models.ConversationModel, error) {
	conversation.IsDeleted = false
	return conversation, r.db.WithContext(ctx).Create(conversation).Error
}

func (r *gormRepository) GetConversations(ctx context.Context, limit, offset int64, conversationType models.ConversationTypeEnum, userId uint64, keyword string) ([]*models.ConversationModel, error) {
	var results []models.ConversationWithMessage
	err := r.db.WithContext(ctx).Raw(GET_CONVERSATIONS_QUERY, userId, conversationType, userId, keyword, limit, offset).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	var conversations []*models.ConversationModel
	for _, res := range results {
		conv := res.ConversationModel
		conv.UnreadCount = res.UnreadCount
		if res.LatestMessageID != 0 {
			conv.LatestMessage = models.MessageModel{
				ID:          res.LatestMessageID,
				Content:     res.LatestMessageContent,
				ContentType: models.ContentTypeEnum((res.LatestMessageContentType)),
				FileURL:     res.LatestMessageFileURL,
				ForwardFrom: res.LatestMessageForwardFrom,
				SenderID:    res.LatestMessageSenderID,
				CreatedAt:   res.LatestMessageCreatedAt,
				IsDeleted:   res.LatestMessageIsDeleted,
				Pined:       res.LatestMessagePined,
				Recall:      res.LatestMessageRecall,
			}
		}
		conversations = append(conversations, &conv)
	}

	return conversations, nil
}

func (r *gormRepository) GetConversationByID(ctx context.Context, conversationID uint64) (*models.ConversationModel, error) {
	var conversation models.ConversationModel
	err := r.db.WithContext(ctx).Where("id = ? AND is_deleted = ?", conversationID, false).First(&conversation).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &conversation, err
}

func (r *gormRepository) UpdateConversation(ctx context.Context, id uint64, isBoardCast, forbidForward bool, name string, backgroundImage string) error {
	return r.db.WithContext(ctx).Model(&models.ConversationModel{}).Where("id = ? AND is_deleted = ?", id, false).Updates(map[string]any{
		"is_broadcast":     isBoardCast,
		"forbid_forward":   forbidForward,
		"name":             name,
		"background_image": backgroundImage,
	}).Error
}

func (r *gormRepository) GetMessageOrConversation(ctx context.Context, keyword string, searchType int, limit int64, offset int64) ([]*models.ConversationModel, []*models.MessageModel, error) {
	var conversations []*models.ConversationModel
	var messages []*models.MessageModel
	err := r.db.WithContext(ctx).Where("name LIKE ? AND is_deleted = ?", "%"+keyword+"%", false).Limit(int(limit)).Offset(int(offset)).Find(&conversations).Error
	if err != nil {
		return nil, nil, err
	}
	err = r.db.WithContext(ctx).Where("content LIKE ? AND is_deleted = ?", "%"+keyword+"%", false).Limit(int(limit)).Offset(int(offset)).Find(&messages).Error
	if err != nil {
		return nil, nil, err
	}
	if searchType == 0 {
		return conversations, nil, nil
	}
	return conversations, messages, nil
}

func (r *gormRepository) GetConversationByIDWithMember(ctx context.Context, conversationID uint64) (*models.ConversationModel, error) {
	var conversation models.ConversationModel
	err := r.db.WithContext(ctx).Where("id = ? AND is_deleted = ?", conversationID, false).Preload("Participants").First(&conversation).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &conversation, err
}

func (r *gormRepository) DeleteConversation(ctx context.Context, conversationID uint64) error {
	// Soft delete by setting is_deleted to true and updating deleted_at
	return r.db.WithContext(ctx).Model(&models.ConversationModel{}).Where("id = ?", conversationID).Updates(map[string]interface{}{
		"is_deleted": true,
		"deleted_at": time.Now(),
	}).Error
}
