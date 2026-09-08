package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"chat/internal/domain"
	_repo "chat/internal/repo"
	_db "common/db"
	_models "common/domain/entity"

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

type ConversationRepo struct {
	*_db.TransactionRepo
}

func NewConversationRepo(db *_db.TransactionRepo) _repo.ConversationRepo {
	return &ConversationRepo{db}
}

func (r *ConversationRepo) CreateConversation(ctx context.Context, conversation *domain.Conversation) (*domain.Conversation, error) {
	err := r.GetDB(ctx).Create(conversation).Error
	if err != nil {
		return nil, err
	}
	return conversation, nil
}

func (r *ConversationRepo) GetConversations(ctx context.Context, limit, offset int64, conversationType *int32, userId uint64, keyword string) ([]*domain.Conversation, error) {
	var results []domain.ConversationWithMessage

	typeQuery := ""
	query := GET_CONVERSATIONS_QUERY
	if conversationType != nil {
		typeQuery = fmt.Sprintf("c.type = %d AND", *conversationType)
	}
	query = strings.ReplaceAll(query, "??", typeQuery)

	err := r.GetDB(ctx).Raw(query, userId, userId, keyword, limit, offset).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	var conversations []*domain.Conversation
	for _, res := range results {
		conv := res.Conversation
		conv.UnreadCount = res.UnreadCount
		if res.LatestMessageID != 0 {
			msg := domain.Message{
				BaseEntity:  _models.BaseEntity{ID: res.LatestMessageID},
				Content:     res.LatestMessageContent,
				ContentType: domain.ContentTypeEnum(res.LatestMessageContentType),
				FileURL:     res.LatestMessageFileURL,
				ForwardFrom: res.LatestMessageForwardFrom,
				SenderID:    res.LatestMessageSenderID,
				Pined:       res.LatestMessagePined,
				Recall:      res.LatestMessageRecall,
			}
			if !res.LatestMessageCreatedAt.IsZero() {
				createdAt := res.LatestMessageCreatedAt
				msg.CreatedAt = &createdAt
			}
			if res.LatestMessageDeletedAt != nil {
				msg.DeletedAt = &gorm.DeletedAt{
					Time:  *res.LatestMessageDeletedAt,
					Valid: true,
				}
			}
			conv.LatestMessage = msg
		}

		var result []uint64
		for _, id := range res.ParticipantIDs {
			result = append(result, uint64(id))
		}
		conv.ParticipantIDs = result
		conversations = append(conversations, &conv)
	}

	return conversations, nil
}

func (r *ConversationRepo) GetConversationByID(ctx context.Context, conversationID uint64) (*domain.Conversation, error) {
	var conversation domain.Conversation
	err := r.GetDB(ctx).
		Where("id = ? AND deleted_at IS NULL", conversationID).
		Preload("BackgroundImage").
		First(&conversation).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &conversation, err
}

func (r *ConversationRepo) GetConversationsByIDs(ctx context.Context, conversationIDs []uint64) ([]*domain.Conversation, error) {
	var conversations []*domain.Conversation
	err := r.GetDB(ctx).
		Where("id IN (?) AND deleted_at IS NULL", conversationIDs).
		Preload("BackgroundImage").
		Find(&conversations).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return conversations, err
}

func (r *ConversationRepo) UpdateConversation(ctx context.Context, id uint64, isBoardCast, forbidForward bool, name string, backgroundImageId *uint64) error {
	return r.GetDB(ctx).Model(&domain.Conversation{}).Where("id = ? AND deleted_at IS NULL", id).Updates(map[string]any{
		"is_broadcast":        isBoardCast,
		"forbid_forward":      forbidForward,
		"name":                name,
		"background_image_id": backgroundImageId,
	}).Error
}

func (r *ConversationRepo) GetMessageOrConversation(ctx context.Context, keyword string, searchType int, limit int64, offset int64) ([]*domain.Conversation, []*domain.Message, error) {
	var conversations []*domain.Conversation
	var messages []*domain.Message
	err := r.GetDB(ctx).Where("name LIKE ? AND deleted_at IS NULL", "%"+keyword+"%").Limit(int(limit)).Offset(int(offset)).Find(&conversations).Error
	if err != nil {
		return nil, nil, err
	}
	err = r.GetDB(ctx).Where("content LIKE ? AND deleted_at IS NULL", "%"+keyword+"%").Limit(int(limit)).Offset(int(offset)).Find(&messages).Error
	if err != nil {
		return nil, nil, err
	}
	if searchType == 0 {
		return conversations, nil, nil
	}
	return conversations, messages, nil
}

func (r *ConversationRepo) GetConversationByIDWithMember(ctx context.Context, conversationID uint64) (*domain.Conversation, error) {
	var conversation domain.Conversation
	err := r.GetDB(ctx).
		Where("id = ? AND deleted_at IS NULL", conversationID).
		Preload("Participants").
		Preload("BackgroundImage").
		First(&conversation).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	conversation.ParticipantIDs = make([]uint64, len(conversation.Participants))
	for i, participant := range conversation.Participants {
		conversation.ParticipantIDs[i] = participant.UserID
	}
	return &conversation, err
}

func (r *ConversationRepo) DeleteConversation(ctx context.Context, conversationID uint64) error {
	now := time.Now()
	return r.GetDB(ctx).Model(&domain.Conversation{}).Where("id = ?", conversationID).Update("deleted_at", now).Error
}

func (r *ConversationRepo) GetConversationWithReceiverID(ctx context.Context, currentUserId uint64, receiverID uint64) (*domain.Conversation, error) {
	var conversation domain.Conversation
	err := r.GetDB(ctx).
		Model(&domain.Conversation{}).
		Select("id").
		Where("(receiver_id = ? AND created_by = ?) OR (receiver_id = ? AND created_by = ?)",
			receiverID, currentUserId, currentUserId, receiverID).
		First(&conversation).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &conversation, err
}

func (r *ConversationRepo) SetBackgroundImage(ctx context.Context, conversationId uint64, backgroundImageId uint64) error {
	return r.GetDB(ctx).Model(&domain.Conversation{}).Where("id = ? AND deleted_at IS NULL", conversationId).Updates(map[string]any{
		"background_image_id": backgroundImageId,
	}).Error
}

func (r *ConversationRepo) GetUnreadCount(ctx context.Context, userId uint64) (int64, error) {
	var count int64
	err := r.GetDB(ctx).
		Raw(QUERY_COUNT_UNREAD_USER, userId, userId).
		Scan(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *ConversationRepo) GetRoomMemberIDs(ctx context.Context, conversationId uint64) ([]uint64, error) {
	var ids []uint64
	err := r.GetDB(ctx).
		Model(&domain.Participant{}).
		Where("conversation_id = ?", conversationId).
		Select("user_id").
		Pluck("user_id", &ids).Error
	if err != nil {
		return nil, err
	}
	return ids, nil
}
