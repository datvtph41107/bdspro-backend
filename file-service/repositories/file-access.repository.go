package repositories

import (
	cmodels "common/models"
	"file/models"

	"gorm.io/gorm"
)

// AccessRepository lưu quyền truy cập file.
type AccessRepository struct {
	cmodels.BaseRepository
}

// NewFileAccessRepository tạo repository với database hiện tại.
func NewFileAccessRepository(db *gorm.DB) *AccessRepository {
	return &AccessRepository{
		BaseRepository: cmodels.BaseRepository{DB: db},
	}
}

// Create lưu một quyền truy cập file.
func (r *AccessRepository) Create(file *models.AccessEntity) error {
	return r.DB.Create(file).Error
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
