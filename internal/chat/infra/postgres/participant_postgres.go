package postgres

import (
	"context"

	"chat/internal/domain"
	_repo "chat/internal/repo"
	_db "common/db"

	"gorm.io/gorm"
)

type ParticipantRepo struct {
	*_db.TransactionRepo
}

func NewParticipantRepo(db *_db.TransactionRepo) _repo.ParticipantRepo {
	return &ParticipantRepo{db}
}

func (r *ParticipantRepo) ExistParticipant(ctx context.Context, conversationID uint64, userID uint64) (bool, error) {
	var count int64
	err := r.GetDB(ctx).Model(&domain.Participant{}).Where("conversation_id = ? AND user_id = ?", conversationID, userID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *ParticipantRepo) ExistParticipants(ctx context.Context, conversationIDs []uint64, userID uint64) (bool, error) {
	var count int64
	err := r.GetDB(ctx).Model(&domain.Participant{}).Where("conversation_id IN (?) AND user_id = ?", conversationIDs, userID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count == int64(len(conversationIDs)), nil
}

func (r *ParticipantRepo) FindParticipant(ctx context.Context, conversationID uint64, userID uint64) (*domain.Participant, error) {
	var participant domain.Participant
	err := r.GetDB(ctx).Where("conversation_id = ? AND user_id = ?", conversationID, userID).First(&participant).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &participant, nil
}

func (r *ParticipantRepo) RemoveParticipant(ctx context.Context, conversationID uint64, userID uint64) error {
	return r.GetDB(ctx).Where("conversation_id = ? AND user_id = ?", conversationID, userID).Delete(&domain.Participant{}).Error
}

func (r *ParticipantRepo) UpdateMuteNotification(ctx context.Context, conversationID uint64, userID uint64, isMute bool) error {
	return r.GetDB(ctx).Model(&domain.Participant{}).Where("conversation_id = ? AND user_id = ?", conversationID, userID).Update("mute_notif", isMute).Error
}

func (r *ParticipantRepo) GetListParticipant(ctx context.Context, conversationID uint64) ([]*domain.Participant, int64, error) {
	var participants []*domain.Participant
	var count int64
	err := r.GetDB(ctx).Model(&domain.Participant{}).Where("conversation_id = ?", conversationID).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	err = r.GetDB(ctx).Where("conversation_id = ?", conversationID).Find(&participants).Error
	if err != nil {
		return nil, 0, err
	}
	return participants, count, nil
}

func (r *ParticipantRepo) AddParticipant(ctx context.Context, participant *domain.Participant) error {
	return r.GetDB(ctx).Create(participant).Error
}

func (r *ParticipantRepo) HasExactConversationWithUserIDs(ctx context.Context, userIDs []uint64) (*domain.Conversation, error) {
	var conversationID uint64
	err := r.GetDB(ctx).
		Model(&domain.Participant{}).
		Select("conversation_id").
		Where("user_id IN ?", userIDs).
		Group("conversation_id").
		Having("COUNT(DISTINCT user_id) = ? AND COUNT(CASE WHEN user_id IN ? THEN 1 END) = ?", len(userIDs), userIDs, len(userIDs)).
		Limit(1).
		Scan(&conversationID).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if conversationID == 0 {
		return nil, nil
	}

	// Get full conversation
	var conversation domain.Conversation
	err = r.GetDB(ctx).Where("id = ?", conversationID).First(&conversation).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &conversation, err
}
