package repositories

import (
	cmodels "common/models"

	"gorm.io/gorm"
)

// AccessRepository is the PostgreSQL adapter for durable private-file access
// evidence.
type AccessRepository struct {
	cmodels.BaseRepository
}

func NewFileAccessRepository(db *gorm.DB) *AccessRepository {
	return &AccessRepository{
		BaseRepository: cmodels.BaseRepository{DB: db},
	}
}

func (r *AccessRepository) ExistsAccess(fileID int64, accessID int64) (bool, error) {
	var exists bool
	query := `SELECT EXISTS (
		SELECT 1 FROM file_access WHERE file_id = ? AND access_id = ?
	)`
	if err := r.DB.Raw(query, fileID, accessID).Scan(&exists).Error; err != nil {
		return false, err
	}
	return exists, nil
}
