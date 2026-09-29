package orm

import (
	"context"
	"fmt"
	"strings"
	"time"

	"chat/models"

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
	m.deleted_at      AS latest_message_deleted_at,
	m.pined           AS latest_message_pined,
	m.recall          AS latest_message_recall,
	COALESCE(unread.unread_count, 0) AS unread_count,
	COALESCE(participant_ids.ids, ARRAY[]::BIGINT[]) AS participant_ids
FROM conversations c
LEFT JOIN LATERAL (
	SELECT *
	FROM messages
	WHERE conversation_id = c.id AND deleted_at IS NULL
	ORDER BY created_at DESC
	LIMIT 1
) m ON true
LEFT JOIN participants p ON p.conversation_id = c.id
LEFT JOIN LATERAL (
	SELECT 
		conversation_id, 
		array_agg(user_id) AS ids
	FROM participants
	WHERE conversation_id = c.id
	GROUP BY conversation_id
) participant_ids ON true
LEFT JOIN (
	SELECT 
		m.conversation_id,
		COUNT(m.id) AS unread_count
	FROM messages m
	LEFT JOIN read_recepts rr ON rr.message_id = m.id AND rr.user_id = $1
	WHERE m.deleted_at IS NULL AND rr.message_id IS NULL
	GROUP BY m.conversation_id
) unread ON unread.conversation_id = c.id
WHERE ?? p.user_id = $2 AND LOWER((TRIM(c.name))) LIKE '%' || LOWER((TRIM($3))) || '%' AND c.deleted_at IS NULL
ORDER BY COALESCE(m.created_at,c.created_at) DESC
LIMIT $4 OFFSET $5
`

var QUERY_COUNT_UNREAD_USER = `
	select count(p.user_id) from participants p
left join conversations c on p.conversation_id = c.id
left join lateral (
	select * from messages m 
	where m.conversation_id = c.id
	order by created_at desc nulls last
	limit 1
) m on true
left join read_recepts r on r.message_id = m.id and r.user_id = $1
where p.user_id = $2 and r.read_at is null and m.created_at is not null
`

func (r *gormRepository) CreateConversation(ctx context.Context, conversation *models.ConversationModel) (*models.ConversationModel, error) {
	err := r.db.WithContext(ctx).Create(conversation).Error
	if err != nil {
		return nil, err
	}
	return conversation, nil
}

func (r *gormRepository) GetConversations(ctx context.Context, limit, offset int64, conversationType *int32, userId uint64, keyword string) ([]*models.ConversationModel, error) {
	var results []models.ConversationWithMessage

	typeQuery := ""
	query := GET_CONVERSATIONS_QUERY
	if conversationType != nil {
		typeQuery = fmt.Sprintf("c.type = %d AND", *conversationType)
	}
	query = strings.ReplaceAll(query, "??", typeQuery)

	err := r.db.WithContext(ctx).Raw(query, userId, userId, keyword, limit, offset).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	var conversations []*models.ConversationModel
	var latestMessageIDs []uint64
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
				DeletedAt:   res.LatestMessageDeletedAt,
				Pined:       res.LatestMessagePined,
				Recall:      res.LatestMessageRecall,
			}
			latestMessageIDs = append(latestMessageIDs, res.LatestMessageID)
		}

		var participantIDs []uint64
		for _, id := range res.ParticipantIDs {
			participantIDs = append(participantIDs, uint64(id))
		}
		conv.ParticipantIDs = participantIDs
		conversations = append(conversations, &conv)
	}

	// Load full LatestMessage (ReplyMessage, Reactions) giống GetMessages
	if len(latestMessageIDs) > 0 {
		var fullMessages []*models.MessageModel
		err = r.db.WithContext(ctx).
			Preload("ReplyMessage").
			Preload("Reactions").
			Where("id IN ?", latestMessageIDs).
			Find(&fullMessages).Error
		if err != nil {
			return nil, err
		}
		messageByID := make(map[uint64]*models.MessageModel)
		for _, m := range fullMessages {
			messageByID[m.ID] = m
		}
		for _, conv := range conversations {
			if conv.LatestMessage.ID != 0 {
				if full, ok := messageByID[conv.LatestMessage.ID]; ok {
					conv.LatestMessage = *full
				}
			}
		}
	}

	return conversations, nil
}

func (r *gormRepository) GetConversationByID(ctx context.Context, conversationID uint64) (*models.ConversationModel, error) {
	var conversation models.ConversationModel
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", conversationID).
		Preload("BackgroundImage").
		First(&conversation).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &conversation, err
}

func (r *gormRepository) GetConversationsByIDs(ctx context.Context, conversationIDs []uint64) ([]*models.ConversationModel, error) {
	var conversations []*models.ConversationModel
	err := r.db.WithContext(ctx).
		Where("id IN (?) AND deleted_at IS NULL", conversationIDs).
		Preload("BackgroundImage").
		Find(&conversations).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return conversations, err
}

func (r *gormRepository) UpdateConversation(ctx context.Context, id uint64, isBoardCast, forbidForward bool, name string, backgroundImageId *uint64) error {
	return r.db.WithContext(ctx).Model(&models.ConversationModel{}).Where("id = ? AND deleted_at IS NULL", id).Updates(map[string]any{
		"is_broadcast":        isBoardCast,
		"forbid_forward":      forbidForward,
		"name":                name,
		"background_image_id": backgroundImageId,
	}).Error
}

func (r *gormRepository) GetMessageOrConversation(ctx context.Context, keyword string, searchType int, limit int64, offset int64) ([]*models.ConversationModel, []*models.MessageModel, error) {
	var conversations []*models.ConversationModel
	var messages []*models.MessageModel
	err := r.db.WithContext(ctx).Where("name LIKE ? AND deleted_at IS NULL", "%"+keyword+"%").Limit(int(limit)).Offset(int(offset)).Find(&conversations).Error
	if err != nil {
		return nil, nil, err
	}
	err = r.db.WithContext(ctx).Where("content LIKE ? AND deleted_at IS NULL", "%"+keyword+"%").Limit(int(limit)).Offset(int(offset)).Find(&messages).Error
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
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", conversationID).
		Preload("Participants").
		Preload("BackgroundImage").
		First(&conversation).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	conversation.ParticipantIDs = make([]uint64, len(conversation.Participants))
	for i, participant := range conversation.Participants {
		conversation.ParticipantIDs[i] = participant.UserID
	}
	// Load lastMessage (tin nhắn mới nhất) đủ như GetMessages: ReplyMessage, Reactions
	var lastMessage models.MessageModel
	err = r.db.WithContext(ctx).
		Preload("ReplyMessage").
		Preload("Reactions").
		Where("conversation_id = ? AND deleted_at IS NULL", conversationID).
		Order("created_at DESC").
		Limit(1).
		First(&lastMessage).Error
	if err == nil && lastMessage.ID != 0 {
		conversation.LatestMessage = lastMessage
	}
	return &conversation, nil
}

func (r *gormRepository) DeleteConversation(ctx context.Context, conversationID uint64) error {
	// Soft delete by setting deleted_at
	now := time.Now()
	return r.db.WithContext(ctx).Model(&models.ConversationModel{}).Where("id = ?", conversationID).Update("deleted_at", now).Error
}

func (r *gormRepository) GetConversationWithReceiverID(ctx context.Context, currentUserId uint64, receiverID uint64) (*models.ConversationModel, error) {
	var conversation models.ConversationModel
	err := r.db.WithContext(ctx).
		Model(&models.ConversationModel{}).
		Select("id").
		Where("(receiver_id = ? AND created_by = ?) OR (receiver_id = ? AND created_by = ?)",
			receiverID, currentUserId, currentUserId, receiverID).
		First(&conversation).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &conversation, err
}

func (r *gormRepository) SetBackgroundImage(ctx context.Context, conversationId uint64, backgroundImageId uint64) error {
	return r.db.WithContext(ctx).Model(&models.ConversationModel{}).Where("id = ? AND deleted_at IS NULL", conversationId).Updates(map[string]any{
		"background_image_id": backgroundImageId,
	}).Error
}

func (r *gormRepository) GetUnreadCount(ctx context.Context, userId uint64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Raw(QUERY_COUNT_UNREAD_USER, userId, userId).
		Scan(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
