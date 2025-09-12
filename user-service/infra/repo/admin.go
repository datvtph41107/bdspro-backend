package repo

import (
	"context"
	"fmt"
	"strings"
	"user/internal/dto"
	user_model "user/internal/models/user"

	"gorm.io/gorm"
)

// @bind: user/internal/repo.AdminRepository
type AdminRepo struct {
	DB *gorm.DB
}

func NewAdminRepo(db *gorm.DB) *AdminRepo {
	return &AdminRepo{
		DB: db,
	}
}

func (r *AdminRepo) ListAllUsers(ctx context.Context, req *dto.ListAllUsersRequest) ([]dto.AdminUserItem, int64, error) {
	var users []dto.AdminUserItem
	var total int64

	query := r.DB.Debug().WithContext(ctx).
		Model(&user_model.UserProfileEntity{}).
		Select(`
			profile_id,
			full_name,
			email,
			phone,
			avatar,
			address,
			gender,
			birth,
			created_at,
			updated_at,
			COALESCE(status, 1) as status,
			role_type,
			role_real_estate,
			position,
			department_id,
			tick_verified,
			NULL as last_login_at,
			0 as total_login
		`)

	// Apply search filter
	if req.Search != "" {
		searchTerm := "%" + strings.ToLower(req.Search) + "%"
		query = query.Where(`
			LOWER(full_name) LIKE ? OR 
			LOWER(email) LIKE ? OR 
			phone LIKE ?
		`, searchTerm, searchTerm, searchTerm)
	}

	// Apply status filter
	if req.Status != nil {
		query = query.Where("COALESCE(status, 1) = ?", *req.Status)
	}

	// Apply role type filter
	if req.RoleType != "" {
		query = query.Where("role_type = ?", req.RoleType)
	}

	// Apply date range filter
	if req.StartDate != nil {
		query = query.Where("created_at >= ?", req.StartDate)
	}
	if req.EndDate != nil {
		query = query.Where("created_at <= ?", req.EndDate)
	}

	// Get total count
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Apply sorting
	if req.SortBy != "" {
		sortOrder := "ASC"
		if strings.ToUpper(req.SortOrder) == "DESC" {
			sortOrder = "DESC"
		}

		// Map sortBy to actual column names
		sortColumn := req.SortBy
		switch req.SortBy {
		case "fullName":
			sortColumn = "full_name"
		case "createdAt":
			sortColumn = "created_at"
		case "updatedAt":
			sortColumn = "updated_at"
		}

		query = query.Order(fmt.Sprintf("%s %s", sortColumn, sortOrder))
	} else {
		// Default sorting by created_at DESC
		query = query.Order("created_at DESC")
	}

	// Apply pagination
	err = query.
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Find(&users).Error

	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetUserByProfileID lấy thông tin user theo profile ID
func (r *AdminRepo) GetUserByProfileID(profileID uint64) (*user_model.UserProfileEntity, error) {

	// Kiểm tra xem có dữ liệu trong bảng không
	var count int64
	r.DB.Model(&user_model.UserProfileEntity{}).Count(&count)

	// Kiểm tra xem có user nào với profile_id này không (kể cả đã bị soft delete)
	var user user_model.UserProfileEntity
	err := r.DB.Unscoped().Where("profile_id = ?", profileID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUserStatus cập nhật trạng thái user
func (r *AdminRepo) UpdateUserStatus(profileID uint64, status uint8) error {
	// Status: 0 = blocked, 1 = active, 2 = approved

	err := r.DB.Model(&user_model.UserProfileEntity{}).
		Where("profile_id = ?", profileID).
		Update("status", status).Error

	if err != nil {
		return err
	}

	return nil
}
