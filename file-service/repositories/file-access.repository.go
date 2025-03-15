package repositories

import (
	_db "common/db"
	cmodels "common/models"
	"file/models"

	"gorm.io/gorm"
)

// FileRepository thao tác với FileEntity trong DB
type AccessRepository struct {
	cmodels.BaseRepository
}

// NewUserRepository khởi tạo UserRepository
func NewFileAccessRepository(db *gorm.DB) *AccessRepository {
	return &AccessRepository{
		BaseRepository: cmodels.BaseRepository{DB: db},
	}
}

// Save lưu file vào DB
func (r *AccessRepository) Create(file *models.AccessEntity) error {
	return _db.DB.Create(file).Error
}

// existAccess kiểm tra sự tồn tại của quyền truy cập
func (r *AccessRepository) ExistsAccess(fileId int64, accessId int64) (bool, error) {
	var result bool
	// Sử dụng raw SQL query để kiểm tra sự tồn tại
	queryStr := "SELECT CASE WHEN COUNT(*) > 0 THEN TRUE ELSE FALSE END FROM file_access WHERE file_id = ? AND access_id = ?"
	err := r.DB.Raw(queryStr, fileId, accessId).Scan(&result).Error
	if err != nil {
		return false, err
	}
	return result, nil
}
