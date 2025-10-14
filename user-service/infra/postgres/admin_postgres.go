package postgres

import (
	_enum "common/domain/enum"
	_utils "common/utils"
	"context"
	"fmt"
	"user/infra/client"
	"user/internal/dto"
	"user/internal/interface/repo"
	models "user/internal/models"

	"gorm.io/gorm"
)

type AdminPostgres struct {
	DB         *gorm.DB
	AuthClient *client.AuthClient
}

func NewAdminPostgres(db *gorm.DB, authClient *client.AuthClient) repo.IAdminRepo {
	return &AdminPostgres{
		DB:         db,
		AuthClient: authClient,
	}
}

func (r *AdminPostgres) ListAdmin(ctx context.Context, req *dto.AdminListRequest) ([]models.AdminProfile, int64, error) {
	var admins []models.AdminProfile
	var total int64

	// Lấy adminID từ context nếu có
	adminID := _utils.GetProfileIdWithContext(ctx)

	// Kiểm tra xem có cần JOIN bookmark table không
	needBookmarkInfo := false

	// Base select fields
	selectFields := `
		up.id,
		up.full_name,
		up.email,
		up.phone,
		up.auth_id,
		up.avatar,
		up.address,
		up.gender,
		up.birth,
		up.created_at,
		up.updated_at,
		COALESCE(up.status, 10) as status,
		up.verified_at,
		up.locked_at,
		up.last_login_at,
		up.total_login,
		up.job_title,
		up.role_id,
		up.role_key,
		up.role_type,
		up.attachments,
		up.internal_notes,
		up.send_notification,
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
	query := r.DB.Debug().WithContext(ctx).Table("admin_profiles up").Select(selectFields)

	// JOIN bookmark table nếu cần
	if needBookmarkInfo {
		query = query.Joins("LEFT JOIN bookmark_user bu ON up.profile_id = bu.user_id AND bu.admin_id = ?", adminID)
	}

	// Apply name filter
	if req.Name != "" {
		like := "%" + req.Name + "%"
		query = query.Where("up.full_name ILIKE ?", like)
	}

	// // Apply status filter
	// if req.Status != nil {
	// 	query = query.Where("COALESCE(up.status, 1) = ?", *req.Status)
	// }

	// Apply roleId filter
	if req.RoleID > 0 {
		query = query.Where("up.role_id = ?", req.RoleID)
	}

	// // Apply date range filter
	// if req.StartDate != nil {
	// 	query = query.Where("up.created_at >= ?", req.StartDate)
	// }
	// if req.EndDate != nil {
	// 	query = query.Where("up.created_at <= ?", req.EndDate)
	// }

	// Get total count
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Apply sorting
	if false {
		// sortOrder := "ASC"
		// if strings.ToUpper(req.SortOrder) == "DESC" {
		// 	sortOrder = "DESC"
		// }

		// // Map sortBy to actual column names
		// sortColumn := req.SortBy
		// switch req.SortBy {
		// case "fullName":
		// 	sortColumn = "up.full_name"
		// case "createdAt":
		// 	sortColumn = "up.created_at"
		// case "updatedAt":
		// 	sortColumn = "up.updated_at"
		// default:
		// 	sortColumn = "up." + sortColumn
		// }

		// query = query.Order(fmt.Sprintf("%s %s", sortColumn, sortOrder))
	} else {
		// Default sorting by created_at DESC
		query = query.Order("up.created_at DESC")
	}

	// Apply pagination
	err = query.
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Find(&admins).Error

	if err != nil {
		return nil, 0, err
	}

	// Lấy thông tin username từ auth-service cho các admin
	if len(admins) > 0 && r.AuthClient != nil {
		// Lấy danh sách auth_id để gọi auth service
		authIDs := make([]uint64, 0, len(admins))
		for _, admin := range admins {
			if admin.AuthID > 0 {
				authIDs = append(authIDs, admin.AuthID)
			}
		}

		// Gọi auth client để lấy thông tin batch
		if len(authIDs) > 0 {
			authAdminsResp, err := r.AuthClient.GetAuthAdminByIds(ctx, authIDs)
			if err == nil && authAdminsResp != nil {
				// Tạo map để lookup nhanh
				authAdminMap := make(map[uint64]string)
				for _, authAdmin := range authAdminsResp {
					authAdminMap[authAdmin.ProfileID] = authAdmin.Username
				}

				// Cập nhật username cho từng admin
				for i := range admins {
					if username, ok := authAdminMap[admins[i].ID]; ok {
						admins[i].Username = username
					}
				}
			}
		}
	}

	return admins, total, nil
}

// GetUserByProfileID lấy thông tin user theo profile ID
func (r *AdminPostgres) GetByID(ctx context.Context, profileID uint64) (*models.AdminProfile, error) {

	// Kiểm tra xem có dữ liệu trong bảng không
	var count int64
	r.DB.Model(&models.AdminProfile{}).Count(&count)

	// Kiểm tra xem có user nào với profile_id này không (kể cả đã bị soft delete)
	var user models.AdminProfile
	err := r.DB.Unscoped().Where("profile_id = ?", profileID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUserStatus cập nhật trạng thái user
func (r *AdminPostgres) UpdateUserStatus(profileID uint64, status _enum.EUserStatus) error {
	err := r.DB.Model(&models.AdminProfile{}).
		Where("profile_id = ?", profileID).
		Update("status", status).Error

	if err != nil {
		return err
	}

	return nil
}

// CreateAdmin tạo admin profile mới
func (r *AdminPostgres) CreateAdmin(ctx context.Context, admin *models.AdminProfile) (*models.AdminProfile, error) {
	err := r.DB.WithContext(ctx).Create(admin).Error
	if err != nil {
		return nil, fmt.Errorf("failed to create admin profile: %v", err)
	}

	return admin, nil
}

// UpdateAdmin cập nhật thông tin admin
func (r *AdminPostgres) UpdateAdmin(ctx context.Context, req *models.AdminProfile) (*models.AdminProfile, error) {
	updates := map[string]interface{}{
		"full_name":         req.FullName,
		"email":             req.Email,
		"phone":             req.Phone,
		"avatar":            req.Avatar,
		"address":           req.Address,
		"gender":            req.Gender,
		"birth":             req.Birth,
		"job_title":         req.JobTitle,
		"work_at":           req.WorkAt,
		"role_id":           req.RoleID,
		"role_key":          req.RoleKey,
		"attachments":       req.Attachments,
		"internal_notes":    req.InternalNotes,
		"send_notification": req.SendNotification,
	}

	err := r.DB.WithContext(ctx).
		Model(&models.AdminProfile{}).
		Where("profile_id = ?", req.ProfileID).
		Updates(updates).Error

	if err != nil {
		return nil, fmt.Errorf("failed to update admin: %v", err)
	}

	return req, nil
}

// DeleteAdmin xóa admin profile (soft delete)
func (r *AdminPostgres) DeleteAdmin(ctx context.Context, authID uint64) error {
	err := r.DB.WithContext(ctx).
		Where("auth_id = ?", authID).
		Delete(&models.AdminProfile{}).Error

	if err != nil {
		return fmt.Errorf("failed to delete admin: %v", err)
	}

	return nil
}

// AssignRole gán quyền cho admin
func (r *AdminPostgres) AssignRole(ctx context.Context, authID uint64, roleKey string) error {
	err := r.DB.WithContext(ctx).
		Model(&models.AdminProfile{}).
		Where("auth_id = ?", authID).
		Update("role_key", roleKey).Error

	if err != nil {
		return fmt.Errorf("failed to assign role: %v", err)
	}

	return nil
}

func (r *AdminPostgres) GetDetail(ctx context.Context, req uint64) (*models.AdminProfile, error) {
	var admin models.AdminProfile
	err := r.DB.WithContext(ctx).Where("id = ?", req).First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}
