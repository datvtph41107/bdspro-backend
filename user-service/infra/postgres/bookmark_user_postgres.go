package postgres

import (
	"context"
	"gorm.io/gorm"
	"user/internal/dto"
	"user/internal/interface/repo"
	models "user/internal/models"
)

// BookmarkUserPostgres implement BookmarkUserRepository
type BookmarkUserPostgres struct {
	db *gorm.DB
}

// NewBookmarkUserPostgres creates the bookmark persistence adapter.
func NewBookmarkUserPostgres(db *gorm.DB) repo.IBookmarkUserRepo {
	return &BookmarkUserPostgres{db: db}
}

// GetBookmarkUsers lấy danh sách bookmark user với phân trang
func (r *BookmarkUserPostgres) GetBookmarkUsers(ctx context.Context, adminID uint64, page, size int) ([]models.BookmarkUserEntity, int64, error) {
	var bookmarks []models.BookmarkUserEntity
	var total int64

	offset := (page - 1) * size

	// Đếm tổng số record
	if err := r.db.Model(&models.BookmarkUserEntity{}).
		Where("admin_id = ?", adminID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Lấy danh sách với phân trang
	if err := r.db.Preload("User").
		Where("admin_id = ?", adminID).
		Order("created_at DESC").
		Offset(offset).
		Limit(size).
		Find(&bookmarks).Error; err != nil {
		return nil, 0, err
	}

	return bookmarks, total, nil
}

// CreateBookmarkUser tạo bookmark user mới
func (r *BookmarkUserPostgres) CreateBookmarkUser(ctx context.Context, bookmark *models.BookmarkUserEntity) error {
	return r.db.Create(bookmark).Error
}

// DeleteBookmarkUser xóa bookmark user
func (r *BookmarkUserPostgres) DeleteBookmarkUser(ctx context.Context, adminID, userID uint64) error {
	return r.db.WithContext(ctx).
		Model(&models.BookmarkUserEntity{}).
		Where("admin_id = ? AND user_id = ?", adminID, userID).
		Delete(&models.BookmarkUserEntity{}).
		Error
}

// DeleteAllBookmarkUsers xóa tất cả bookmark user của admin
func (r *BookmarkUserPostgres) DeleteAllBookmarkUsers(ctx context.Context, adminID uint64) error {
	return r.db.WithContext(ctx).Where("admin_id = ?", adminID).
		Delete(&models.BookmarkUserEntity{}).Error
}

// CheckBookmarkUserExists kiểm tra bookmark user đã tồn tại chưa
func (r *BookmarkUserPostgres) CheckBookmarkUserExists(ctx context.Context, adminID, userID uint64) (bool, error) {
	var count int64
	err := r.db.Model(&models.BookmarkUserEntity{}).
		Where("admin_id = ? AND user_id = ?", adminID, userID).
		Count(&count).Error
	return count > 0, err
}

// GetBookmarkedUsers lấy danh sách user có bookmark với phân trang và tìm kiếm
func (r *BookmarkUserPostgres) GetBookmarkedUsers(ctx context.Context, adminID uint64, req *dto.BookmarkUsersRequest) ([]dto.AdminUserItem, int64, error) {
	var users []dto.AdminUserItem
	var total int64

	// Query bắt đầu từ bookmark_user rồi join sang user_profile để tối ưu
	query := r.db.Table("bookmark_user bu").
		Select(`u.profile_id, u.full_name, u.email, u.phone, u.avatar, u.address, u.gender, 
				u.birth, u.created_at, u.updated_at, u.status, u.role_type, 
				u.role_real_estate, u.position, u.department_id, u.tick_verified`).
		Joins("INNER JOIN user_profile u ON bu.user_id = u.profile_id").
		Where("bu.admin_id = ?", adminID)

	// Thêm điều kiện tìm kiếm
	if req.Search != "" {
		query = query.Where("(u.full_name LIKE ? OR u.email LIKE ? OR u.phone LIKE ?)",
			"%"+req.Search+"%", "%"+req.Search+"%", "%"+req.Search+"%")
	}

	// Thêm điều kiện lọc theo role type
	if req.RoleType != "" {
		query = query.Where("u.role_type = ?", req.RoleType)
	}

	// Thêm điều kiện lọc theo ngày tạo
	if req.StartDate != "" {
		query = query.Where("u.created_at >= ?", req.StartDate)
	}
	if req.EndDate != "" {
		query = query.Where("u.created_at <= ?", req.EndDate)
	}

	// Đếm tổng số record
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Lấy danh sách với phân trang
	if err := query.Order("u.created_at DESC").
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Scan(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}
