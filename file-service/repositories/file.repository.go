package repositories

import (
	cmodels "common/models"
	"file/models"

	"gorm.io/gorm"
)

// FileRepository thao tác với FileEntity trong DB
type FileRepository struct {
	cmodels.BaseRepository
}

// NewUserRepository khởi tạo UserRepository
func NewFileRepository(db *gorm.DB) *FileRepository {
	return &FileRepository{
		BaseRepository: cmodels.BaseRepository{DB: db},
	}
}

// Save lưu file vào DB
func (r *FileRepository) Save(file *models.FileEntity) error {
	return r.DB.Save(file).Error
}

// DeleteByID is used to compensate a failed filesystem/access write.
func (r *FileRepository) DeleteByID(id uint64) error {
	if id == 0 {
		return nil
	}
	return r.DB.Unscoped().Delete(&models.FileEntity{}, id).Error
}
