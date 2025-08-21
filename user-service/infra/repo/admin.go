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

	query := r.DB.WithContext(ctx).
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
			CASE WHEN deleted_at IS NULL THEN 1 ELSE 0 END as status,
			role_type,
			role_real_estate,
			position,
			department_id,
			tick_verified
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
		if *req.Status == 1 {
			query = query.Where("deleted_at IS NULL")
		} else {
			query = query.Where("deleted_at IS NOT NULL")
		}
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
