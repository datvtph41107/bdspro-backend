package orm

import (
	"context"

	"chat/models"

	"gorm.io/gorm"
)

func (r *gormRepository) CreateBatchRead(ctx context.Context, models []*models.ReadReceptModel) error {
	return r.db.WithContext(ctx).Create(models).Error
}

func (r *gormRepository) GetReadReceipts(ctx context.Context, page, size int64, messageID uint64) ([]*models.ReadReceptModel, int64, error) {
	var readReceipts []*models.ReadReceptModel
	var count int64
	err := r.db.WithContext(ctx).Where("message_id = ?", messageID).Limit(int(size)).Offset(int(page * size)).Find(&readReceipts).Error
	if err != nil {
		return nil, 0, err
	}
	err = r.db.WithContext(ctx).Model(&models.ReadReceptModel{}).Where("message_id = ?", messageID).Count(&count).Error
	return readReceipts, count, err
}

func (r *gormRepository) GetReadReceiptsByUserId(ctx context.Context, userId uint64) ([]*models.ReadReceptModel, error) {
	var readReceipts []*models.ReadReceptModel
	err := r.db.WithContext(ctx).Where("user_id = ?", userId).Find(&readReceipts).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return readReceipts, nil
}

func (r *gormRepository) GetReadReceiptsByUserAndMessageIDs(ctx context.Context, userId uint64, messageIDs []uint64) ([]*models.ReadReceptModel, error) {
	if len(messageIDs) == 0 {
		return []*models.ReadReceptModel{}, nil
	}
	var readReceipts []*models.ReadReceptModel
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND message_id IN (?)", userId, messageIDs).
		Find(&readReceipts).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return []*models.ReadReceptModel{}, nil
		}
		return nil, err
	}
	return readReceipts, nil
}
