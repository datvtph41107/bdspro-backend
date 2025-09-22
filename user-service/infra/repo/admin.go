package repo

import (
	_utils "common/utils"
	"context"
	"fmt"
	"strings"
	"user/infra/client"
	"user/internal/dto"
	user_model "user/internal/models/user"

	"gorm.io/gorm"
)

// @bind: user/internal/repo.AdminRepository
type AdminRepo struct {
	DB         *gorm.DB
	AuthClient *client.AuthClient
}

func NewAdminRepo(db *gorm.DB, authClient *client.AuthClient) *AdminRepo {
	return &AdminRepo{
		DB:         db,
		AuthClient: authClient,
	}
}

func (r *AdminRepo) ListAllUsers(ctx context.Context, req *dto.ListAllUsersRequest) ([]dto.AdminUserItem, int64, error) {
	var users []dto.AdminUserItem
	var total int64

	// Lấy adminID từ context nếu có
	adminID := _utils.GetProfileIdWithContext(ctx)

	// Kiểm tra xem có cần JOIN bookmark table không
	needBookmarkInfo := req.Bookmark != nil && *req.Bookmark

	// Base select fields
	selectFields := `
		up.profile_id,
		up.full_name,
		up.email,
		up.phone,
		up.avatar,
		up.address,
		up.gender,
		up.birth,
		up.created_at,
		up.updated_at,
		COALESCE(up.status, 10) as status,
		up.role_type,
		up.role_real_estate,
		up.position,
		up.department_id,
		up.tick_verified,
		NULL as last_login_at,
		0 as total_login`

	// Thêm bookmark field tùy theo trường hợp
	if needBookmarkInfo {
		selectFields += `,
		CASE WHEN bu.user_id IS NOT NULL THEN true ELSE false END as bookmark`
	} else {
		selectFields += `,
		false as bookmark`
	}

	// Tạo query
	query := r.DB.Debug().WithContext(ctx).Table("user_profile up").Select(selectFields)

	// JOIN bookmark table nếu cần
	if needBookmarkInfo {
		query = query.Joins("LEFT JOIN bookmark_user bu ON up.profile_id = bu.user_id AND bu.admin_id = ?", adminID)
	}

	// Apply search filter
	if req.Search != "" {
		searchTerm := "%" + strings.ToLower(req.Search) + "%"
		query = query.Where(`
			LOWER(up.full_name) LIKE ? OR 
			LOWER(up.email) LIKE ? OR 
			up.phone LIKE ?
		`, searchTerm, searchTerm, searchTerm)
	}

	// Apply status filter
	if req.Status != nil {
		query = query.Where("COALESCE(up.status, 1) = ?", *req.Status)
	}

	// Apply role type filter
	if req.RoleType != "" {
		query = query.Where("up.role_type = ?", req.RoleType)
	}

	// Apply date range filter
	if req.StartDate != nil {
		query = query.Where("up.created_at >= ?", req.StartDate)
	}
	if req.EndDate != nil {
		query = query.Where("up.created_at <= ?", req.EndDate)
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
			sortColumn = "up.full_name"
		case "createdAt":
			sortColumn = "up.created_at"
		case "updatedAt":
			sortColumn = "up.updated_at"
		default:
			sortColumn = "up." + sortColumn
		}

		query = query.Order(fmt.Sprintf("%s %s", sortColumn, sortOrder))
	} else {
		// Default sorting by created_at DESC
		query = query.Order("up.created_at DESC")
	}

	// Apply pagination
	err = query.
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Find(&users).Error

	if err != nil {
		return nil, 0, err
	}

	// Lấy thông tin từ auth-service cho các user
	if len(users) > 0 && r.AuthClient != nil {
		profileIDs := make([]uint64, len(users))
		for i, user := range users {
			profileIDs[i] = user.ProfileID
		}

		// Gọi auth client để lấy thông tin status
		statusMap, err := r.AuthClient.GetUserStatusBatch(ctx, profileIDs)
		if err == nil {
			// Cập nhật thông tin verifiedAt và lockedAt từ auth-service
			for i := range users {
				if statusInfo, exists := statusMap[users[i].ProfileID]; exists {
					if statusInfo.VerifiedAt != nil {
						// Parse string time to time.Time nếu cần
						// users[i].VerifiedAt = parseTime(*statusInfo.VerifiedAt)
					}
					if statusInfo.LockedUntil != nil {
						// Parse string time to time.Time nếu cần
						// users[i].LockedAt = parseTime(*statusInfo.LockedUntil)
					}
				}
			}
		}
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
