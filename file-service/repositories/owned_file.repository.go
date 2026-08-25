package repositories

import (
	"context"

	"file/models"

	"gorm.io/gorm/clause"
)

// TryCreateOwned lets PostgreSQL arbitrate one logical File effect.
// Concurrent callers for the same owner converge through the unique index.
func (r *FileRepository) TryCreateOwned(
	ctx context.Context,
	file *models.FileEntity,
) (bool, error) {
	result := r.DB.
		WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(file)

	if result.Error != nil {
		return false, result.Error
	}

	return result.RowsAffected == 1, nil
}

func (r *FileRepository) FindByOwner(
	ctx context.Context,
	namespace string,
	key string,
) (*models.FileEntity, bool, error) {
	var row models.FileEntity

	result := r.DB.
		WithContext(ctx).
		Where(
			"owner_namespace = ? AND owner_key = ?",
			namespace,
			key,
		).
		Limit(1).
		Find(&row)

	if result.Error != nil {
		return nil, false, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, false, nil
	}

	return &row, true, nil
}

func (r *FileRepository) SaveWithContext(
	ctx context.Context,
	file *models.FileEntity,
) error {
	return r.DB.WithContext(ctx).Save(file).Error
}
