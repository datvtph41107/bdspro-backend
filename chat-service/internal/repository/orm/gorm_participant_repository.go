package orm

import (
	"context"

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

func (r *gormRepository) AddParticipant(ctx context.Context, participant *models.ParticipantModel) error {
	return r.db.WithContext(ctx).Create(participant).Error
}

func (r *gormRepository) HasExactConversationWithUserIDs(ctx context.Context, userIDs []uint64) (uint64, error) {
	var conversationID uint64
	err := r.db.WithContext(ctx).
		Model(&models.ParticipantModel{}).
		Select("conversation_id").
		Group("conversation_id").
		Having("COUNT(*) = ? AND COUNT(CASE WHEN user_id IN ? THEN 1 END) = ?", len(userIDs), userIDs, len(userIDs)).
		Limit(1).
		Scan(&conversationID).Error
	if err != nil {
		return 0, err
	}
	return conversationID, nil
}
