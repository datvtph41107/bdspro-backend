package orm

import (
	"context"
	"time"

	"chat/models"
)

func (r *gormRepository) ExistParticipant(ctx context.Context, conversationID uint64, userID uint64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.ParticipantModel{}).Where("conversation_id = ? AND user_id = ?", conversationID, userID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *gormRepository) ExistParticipants(ctx context.Context, conversationIDs []uint64, userID uint64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.ParticipantModel{}).Where("conversation_id IN (?) AND user_id = ?", conversationIDs, userID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count == int64(len(conversationIDs)), nil
}

func (r *gormRepository) FindParticipant(ctx context.Context, conversationID uint64, userID uint64) (*models.ParticipantModel, error) {
	var participant models.ParticipantModel
	err := r.db.WithContext(ctx).Where("conversation_id = ? AND user_id = ?", conversationID, userID).First(&participant).Error
	if err != nil {
		return nil, err
	}
	return &participant, nil
}

func (r *gormRepository) RemoveParticipant(ctx context.Context, conversationID uint64, userID uint64) error {
	return r.db.WithContext(ctx).Where("conversation_id = ? AND user_id = ?", conversationID, userID).Delete(&models.ParticipantModel{}).Error
}

func (r *gormRepository) UpdateMuteNotification(ctx context.Context, conversationID uint64, userID uint64, isMute bool) error {
	return r.db.WithContext(ctx).Model(&models.ParticipantModel{}).Where("conversation_id = ? AND user_id = ?", conversationID, userID).Update("mute_notif", isMute).Error
}

func (r *gormRepository) GetListParticipant(ctx context.Context, conversationID uint64) ([]*models.ParticipantModel, int64, error) {
	var participants []*models.ParticipantModel
	var count int64
	err := r.db.WithContext(ctx).Model(&models.ParticipantModel{}).Where("conversation_id = ?", conversationID).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	err = r.db.WithContext(ctx).Where("conversation_id = ?", conversationID).Find(&participants).Error
	if err != nil {
		return nil, 0, err
	}
	return participants, count, nil
}

func (r *gormRepository) GetLatestParticipantJoinedAt(ctx context.Context, conversationID uint64) (*time.Time, error) {
	var t time.Time
	err := r.db.WithContext(ctx).
		Model(&models.ParticipantModel{}).
		Where("conversation_id = ?", conversationID).
		Select("joined_at").
		Order("joined_at DESC").
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

func (r *gormRepository) AddParticipant(ctx context.Context, participant *models.ParticipantModel) error {
	return r.db.WithContext(ctx).Create(participant).Error
}

func (r *gormRepository) HasExactConversationWithUserIDs(ctx context.Context, userIDs []uint64) (*models.ConversationModel, error) {
	var conversation *models.ConversationModel
	err := r.db.WithContext(ctx).
		Model(&models.ParticipantModel{}).
		// Select("conversation_id").
		Group("conversation_id, user_id").
		Having("COUNT(*) = ? AND COUNT(CASE WHEN user_id IN ? THEN 1 END) = ?", len(userIDs), userIDs, len(userIDs)).
		Limit(1).
		Scan(conversation).Error
	if err != nil {
		return nil, err
	}
	return conversation, nil
}
