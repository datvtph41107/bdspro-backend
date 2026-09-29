package postgres_v2

import (
	"chat/internal/domain"
	"context"

	"gorm.io/gorm"
)

type MessagePosgres struct {
	db *gorm.DB
}

func NewMessagePostgres(db *gorm.DB) *MessagePosgres {
	return &MessagePosgres{db}
}

func (r *MessagePosgres) UpdateRead(ctx context.Context, userID, convID, seq uint64, ts int64) error {
	return r.db.WithContext(ctx).
		Model(&domain.Membership{}).
		Where("user_id=? AND conversation_id=?", userID, convID).
		Updates(map[string]interface{}{
			"last_read_seq_id": seq,
			"last_read_ts":     ts,
		}).Error
}

func (r *MessagePosgres) GlobalReadSeq(ctx context.Context, convID uint64) (uint64, error) {
	var seq uint64
	err := r.db.WithContext(ctx).
		Model(&domain.Membership{}).
		Where("conversation_id=?", convID).
		Select("MIN(last_read_seq_id)").
		Scan(&seq).Error
	return seq, err
}

func (r *MessagePosgres) Revoke(ctx context.Context, convID, msgID uint64) error {
	return r.db.WithContext(ctx).
		Model(&domain.ChatTimeline{}).
		Where("conversation_id=? AND id=?", convID, msgID).
		Update("is_deleted", true).Error
}
