package orm

import (
	_dto "common/domain/dto"
	"context"
	"time"

	"chat/models"

	"gorm.io/gorm"
)

func (r *gormRepository) CreateMessage(ctx context.Context, message *models.MessageModel) (*models.MessageModel, error) {
	return message, r.db.WithContext(ctx).Create(message).Error
}

func (r *gormRepository) CreateMessages(ctx context.Context, messages []*models.MessageModel) ([]*models.MessageModel, error) {
	return messages, r.db.WithContext(ctx).Create(messages).Error
}

func (r *gormRepository) GetMessages(ctx context.Context, pagable _dto.Pagable, conversationId uint64, sort string, fromDate, toDate *time.Time) ([]*models.MessageModel, error) {
	var messages []*models.MessageModel
	err := r.db.WithContext(ctx).
		Preload("ReplyMessage").
		Preload("Reactions").
		Where("conversation_id = ? and (? is null or updated_at >= ?) and (? is null or updated_at <= ?)", conversationId, fromDate, fromDate, toDate, toDate).
		Limit(pagable.GetLimit()).
		Offset(pagable.GetOffset()).
		Order(pagable.Sort).
		Find(&messages).Error
	return messages, err
}

func (r *gormRepository) GetLatestMessageCreatedAt(ctx context.Context, conversationId uint64) (*time.Time, error) {
	var t time.Time
	err := r.db.WithContext(ctx).
		Model(&models.MessageModel{}).
		Where("conversation_id = ?", conversationId).
		Select("updated_at").
		Order("updated_at DESC").
		Limit(1).
		Scan(&t).Error
	if err != nil {
		return nil, err
	}
	if t.IsZero() {
		return nil, nil
	}
	return &t, nil
}

func (r *gormRepository) GetLatestPinnedMessageUpdatedAt(ctx context.Context, conversationId uint64) (*time.Time, error) {
	var t time.Time
	err := r.db.WithContext(ctx).
		Model(&models.MessageModel{}).
		Where("conversation_id = ? AND pinned_at IS NOT NULL", conversationId).
		Select("pinned_at").
		Order("pinned_at DESC").
		Limit(1).
		Scan(&t).Error
	if err != nil {
		return nil, err
	}
	if t.IsZero() {
		return nil, nil
	}
	return &t, nil
}

func (r *gormRepository) GetLatestMessageUpdatedAtByTypes(ctx context.Context, conversationId uint64, types []string) (*time.Time, error) {
	if len(types) == 0 {
		return nil, nil
	}
	var t time.Time
	err := r.db.WithContext(ctx).
		Model(&models.MessageModel{}).
		Where("conversation_id = ? AND content_type IN (?)", conversationId, types).
		Select("updated_at").
		Order("updated_at DESC").
		Limit(1).
		Scan(&t).Error
	if err != nil {
		return nil, err
	}
	if t.IsZero() {
		return nil, nil
	}
	return &t, nil
}

func (r *gormRepository) FindMessageByID(ctx context.Context, messageId uint64) (*models.MessageModel, error) {
	var message models.MessageModel
	err := r.db.WithContext(ctx).Where("id = ?", messageId).First(&message).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &message, err
}

func (r *gormRepository) UpdatePinMessage(ctx context.Context, messageID uint64, isPin bool) (*models.MessageModel, error) {
	var message *models.MessageModel
	updates := map[string]interface{}{
		"pined": isPin,
	}
	if isPin {
		now := time.Now()
		updates["pinned_at"] = &now
	} else {
		updates["pinned_at"] = nil
	}
	err := r.db.WithContext(ctx).Model(&models.MessageModel{}).Where("id = ?", messageID).Updates(updates).Error
	return message, err
}

func (r *gormRepository) UpdateRecallMessage(ctx context.Context, messageID uint64, isRecall bool) (*models.MessageModel, error) {
	var message *models.MessageModel
	err := r.db.WithContext(ctx).Model(&models.MessageModel{}).Where("id = ?", messageID).Update("recall", isRecall).Error
	return message, err
}

func (r *gormRepository) GetPinnedMessages(ctx context.Context, conversationID uint64) ([]*models.MessageModel, error) {
	var messages []*models.MessageModel
	err := r.db.WithContext(ctx).
		Preload("ReplyMessage").
		Preload("Reactions").
		Where("conversation_id = ? AND pined = ?", conversationID, true).
		Order("pinned_at DESC").
		Find(&messages).
		Error
	return messages, err
}

func (r *gormRepository) UpdateMessage(ctx context.Context, messageID uint64, content string) (*models.MessageModel, error) {
	var message *models.MessageModel
	err := r.db.WithContext(ctx).Model(&models.MessageModel{}).Where("id = ?", messageID).Update("content", content).Error
	return message, err
}

func (r *gormRepository) DeleteViolation(ctx context.Context, messageID uint64) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&models.MessageModel{}).Where("id = ?", messageID).Update("deleted_at", now).Error
}

func (r *gormRepository) DeleteAllMessagesByConversationID(ctx context.Context, conversationID uint64) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&models.MessageModel{}).Where("conversation_id = ?", conversationID).Update("deleted_at", now).Error
}

func (r *gormRepository) GetMessageByTypeWithPagination(ctx context.Context, page, size int64, conversationId uint64, types []string) ([]*models.MessageModel, int64, error) {
	var messages []*models.MessageModel
	var count int64
	err := r.db.WithContext(ctx).
		Where("conversation_id = ? AND content_type IN (?)", conversationId, types).
		Order("created_at desc"). // Sort by created_at in ascending order
		Limit(int(size)).
		Offset(int(page * size)).
		Find(&messages).Error
	if err != nil {
		return nil, 0, err
	}
	err = r.db.WithContext(ctx).
		Model(&models.MessageModel{}).
		Where("conversation_id = ? AND content_type IN (?)", conversationId, types).
		Count(&count).Error
	return messages, count, err
}
