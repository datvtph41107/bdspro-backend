package repositories

import (
	"context"

	"file/models"

	"gorm.io/gorm"
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

// FinalizeUpload atomically publishes ordinary File metadata and its optional
// private-access evidence. The physical artifact is already staged by the
// caller; a DB failure is compensated by removing that artifact and the
// allocated metadata row.
func (r *FileRepository) FinalizeUpload(
	ctx context.Context,
	file *models.FileEntity,
	access *models.AccessEntity,
) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(file).Error; err != nil {
			return err
		}
		if access != nil {
			if err := tx.Create(access).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
