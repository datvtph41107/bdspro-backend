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
