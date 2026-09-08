package repositories

import (
	"context"

	cmodels "common/models"
	"file/models"

	"gorm.io/gorm"
)

// FileRepository is the PostgreSQL adapter for durable File metadata.
type FileRepository struct {
	cmodels.BaseRepository
}

func NewFileRepository(db *gorm.DB) *FileRepository {
	return &FileRepository{
		BaseRepository: cmodels.BaseRepository{DB: db},
	}
}

// DeleteByIDWithContext hard-deletes a not-yet-published File row while
// compensating a failed storage/access effect. Published/owned records are not
// removed through this method.
func (r *FileRepository) DeleteByIDWithContext(ctx context.Context, id uint64) error {
	if id == 0 {
		return nil
	}
	return r.DB.WithContext(ctx).Unscoped().Delete(&models.FileEntity{}, id).Error
}
