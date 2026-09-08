package postgres

import (
	"context"
	"time"

	"chat/internal/domain"
	_repo "chat/internal/repo"
	_db "common/db"

	"gorm.io/gorm"
)

type ReaMarkRepo struct {
	*_db.TransactionRepo
}

func NewReadPositionRepo(db *_db.TransactionRepo) _repo.ReadMarkRepo {
	return &ReaMarkRepo{db}

}

func (r *ReaMarkRepo) GetReadPosition(ctx context.Context, convId, userId uint64) (*domain.ReadMark, error) {
	var pos domain.ReadMark
	err := r.GetDB(ctx).
		Where("conversation_id = ? AND user_id = ?", convId, userId).
		First(&pos).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &pos, err
}

func (r *ReaMarkRepo) UpsertReadMark(
	ctx context.Context,
	conversationID uint64,
	userID uint64,
	lastMessageID uint64,
) (uint64, bool, time.Time, error) {

	type Result struct {
		LastMessageID uint64
		ReadAt        time.Time
		Updated       bool
	}

	var res Result

	sql := `
    INSERT INTO read_positions (conversation_id, user_id, last_message_id, read_at)
    VALUES (?, ?, ?, NOW())
    ON CONFLICT (conversation_id, user_id)
    DO UPDATE SET
        last_message_id = GREATEST(read_positions.last_message_id, EXCLUDED.last_message_id),
        read_at = CASE
            WHEN EXCLUDED.last_message_id > read_positions.last_message_id
            THEN NOW()
            ELSE read_positions.read_at
        END
    RETURNING last_message_id,
              read_at,
              (EXCLUDED.last_message_id > read_positions.last_message_id) AS updated
    `

	err := r.GetDB(ctx).Raw(
		sql,
		conversationID,
		userID,
		lastMessageID,
	).Scan(&res).Error

	if err != nil {
		return 0, false, time.Time{}, err
	}

	return res.LastMessageID, res.Updated, res.ReadAt, nil
}

func (r *ReaMarkRepo) GetReadPositions(
	ctx context.Context,
	conversationID uint64,
) ([]*domain.ReadMark, error) {

	var rows []domain.ReadMark

	err := r.GetDB(ctx).
		Table("read_positions").
		Where("conversation_id = ?", conversationID).
		Find(&rows).Error

	if err != nil {
		return nil, err
	}

	res := make([]*domain.ReadMark, 0, len(rows))
	for i := range rows {
		res = append(res, &rows[i])
	}

	return res, nil
}

func (r *ReaMarkRepo) GetReadPositionsByConversation(ctx context.Context, convId uint64) ([]*domain.ReadMark, error) {
	var res []*domain.ReadMark
	err := r.GetDB(ctx).
		Where("conversation_id = ?", convId).
		Find(&res).Error
	return res, err
}

func (r *ReaMarkRepo) CreateBatchRead(ctx context.Context, models []*domain.ReadMark) error {
	return r.GetDB(ctx).Create(models).Error
}

func (r *ReaMarkRepo) GetReadReceipts(ctx context.Context, limit, offset int64, messageID uint64) ([]*domain.ReadMark, int64, error) {
	var receipts []*domain.ReadMark
	var count int64

	err := r.GetDB(ctx).
		Model(&domain.ReadMark{}).
		Where("last_message_id = ?", messageID).
		Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.GetDB(ctx).
		Where("last_message_id = ?", messageID).
		Limit(int(limit)).
		Offset(int(offset)).
		Find(&receipts).Error
	return receipts, count, err
}

func (r *ReaMarkRepo) GetReadReceiptsByUserId(ctx context.Context, userId uint64) ([]*domain.ReadMark, error) {
	var receipts []*domain.ReadMark
	err := r.GetDB(ctx).
		Where("user_id = ?", userId).
		Find(&receipts).Error
	return receipts, err
}
