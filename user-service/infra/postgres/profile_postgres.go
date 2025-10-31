package postgres

import (
	_dto "common/domain/dto"
	_enum "common/domain/enum"
	_utils "common/utils"
	"context"
	"user/internal/dto"
	models "user/internal/models"

	"gorm.io/gorm"
)

// @bind: user/internal/interface/repo.IProfileRepo
type ProfilePostgres struct {
	DB *gorm.DB
}

func NewProfilePostgres(db *gorm.DB) *ProfilePostgres {
	return &ProfilePostgres{
		DB: db,
	}
}

func (r *ProfilePostgres) GlobalSearchProfile(c context.Context, text string, pageable _dto.Pagable) (*[]models.UserProfileSearch, int64, error) {
	var profiles []models.UserProfileSearch
	var total int64
	currentUserId := _utils.GetProfileIdWithContext(c)

	query := r.DB.WithContext(c).
		Model(&models.UserProfileEntity{}).
		Select("profile_id, full_name, avatar, email").
		Where("deleted_at is null AND profile_id != ?", currentUserId)

	if text != "" {
		query = query.Where("LOWER(full_name) LIKE ? OR email LIKE ?", "%"+text+"%", "%"+text+"%")
	}

	// Get total count first (use Session to avoid affecting the original query)
	err := query.Session(&gorm.Session{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Then fetch data with pagination
	err = query.Debug().
		Offset(pageable.GetOffset()).
		Limit(pageable.GetLimit()).
		Order("updated_at DESC").
		Find(&profiles).
		Error
	if err != nil {
		return nil, 0, err
	}

	return &profiles, total, nil
}

func (r *ProfilePostgres) FriendSearchProfile(c context.Context, profileId uint64, dto dto.ProfileSearch) (*[]models.UserProfileSearch, int64, error) {
	var friends []models.UserProfileSearch
	var total int64

	query := r.DB.
		Model(&models.FriendEntity{}).
		Select("user_profile.profile_id, user_profile.full_name, user_profile.avatar, user_profile.email").
		Preload("GroupUser").
		Preload("GroupReceiver").
		Joins("LEFT JOIN user_profile ON (friend.created_by = user_profile.profile_id OR friend.receiver_id = user_profile.profile_id) "+
			"						AND user_profile.profile_id <> ? ", profileId).
		Where("friend.status = 20 and friend.deleted_at is null")

	if dto.Text != "" {
		query = query.Where("LOWER(user_profile.full_name) LIKE ? OR user_profile.email LIKE ?", "%"+dto.Text+"%", "%"+dto.Text+"%")
	}

	// Get total count first (use Session to avoid affecting the original query)
	err := query.Session(&gorm.Session{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Then fetch data with pagination
	err = query.Debug().Offset(dto.GetOffset()).Limit(dto.GetLimit()).Scan(&friends).Error
	if err != nil {
		return nil, 0, err
	}

	return &friends, total, err
}

func (r *ProfilePostgres) UpdateProfile(c context.Context, profileId uint64, dto *dto.UserInfoRequest) error {
	// Chỉ cập nhật các trường cụ thể
	return r.DB.WithContext(c).
		Model(&models.UserProfileEntity{}).
		Where("profile_id = ?", profileId).
		Updates(map[string]interface{}{
			"full_name":          dto.FullName,
			"email":              dto.Email,
			"avatar":             dto.Avatar,
			"background_image":   dto.BackgroundImage,
			"phone":              dto.Phone,
			"birth":              dto.Birth,
			"start_date":         dto.StartDate,
			"address":            dto.Address,
			"phone2":             dto.Phone2,
			"tax_code":           dto.TaxCode,
			"position":           dto.Position,
			"department_id":      dto.DepartmentID,
			"slogan":             dto.Slogan,
			"website":            dto.Website,
			"facebook":           dto.Facebook,
			"instagram":          dto.Instagram,
			"twitter":            dto.Twitter,
			"linkedin":           dto.Linkedin,
			"youtube":            dto.Youtube,
			"introduction":       dto.Introduction,
			"visibility":         dto.Visibility,
			"profile_visibility": dto.ProfileVisibility,
			"status_online":      dto.StatusOnline,
			"front_identify":     dto.FrontIdentify,
			"back_identify":      dto.BackIdentify,
			"tab_default":        dto.TabDefault,
			"role_real_estate":   dto.RoleRealEstate,
		}).Error
}

// GetByProfileId truy vấn UserProfileEntity theo profileId
func (r *ProfilePostgres) GetByProfileID(profileId uint64) (*models.UserProfileEntity, error) {
	var userProfile models.UserProfileEntity
	if err := r.DB.Where("profile_id = ? and deleted_at is null", profileId).First(&userProfile).Error; err != nil {
		return nil, err
	}
	return &userProfile, nil
}

func (r *ProfilePostgres) GetProfileByIds(ids []uint64) ([]*models.UserProfileEntity, error) {
	var userProfiles []*models.UserProfileEntity
	if err := r.DB.Where("profile_id IN (?) and deleted_at is null", ids).Find(&userProfiles).Error; err != nil {
		return nil, err
	}
	return userProfiles, nil
}

func (r *ProfilePostgres) GetProfileByPhones(phones []string) ([]*models.UserProfileEntity, error) {
	var userProfiles []*models.UserProfileEntity
	if err := r.DB.Where("phone IN (?) and deleted_at is null", phones).Find(&userProfiles).Error; err != nil {
		return nil, err
	}
	return userProfiles, nil
}

func (r *ProfilePostgres) ListAllUsers(ctx context.Context, req *dto.ListUsersRequest) ([]models.UserProfileEntity, uint32, error) {
	var users []models.UserProfileEntity
	var total int64

	// Lấy adminID từ context nếu có
	adminID := _utils.GetProfileIdWithContext(ctx)

	// Luôn JOIN bookmark table để trả về thông tin bookmark
	needBookmarkInfo := true
	hasBookmarkFilter := req.Bookmark != nil

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
		COALESCE(up.warning, 0) as warning,
		up.role_id,
		up.role_key,
		up.position,
		up.department_id,
		up.tick_verified,
		NULL as verified_at,
		NULL as locked_at,
		NULL as last_login_at,
		NULL as total_login,
		NULL as job_title,
		NULL as attachments,
		NULL as internal_notes,
		NULL as send_notification,
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

	// Filter: deleted records
	query = query.Where("up.deleted_at IS NULL")

	// Filter: search by name, email, phone
	if req.Search != "" {
		searchTerm := "%" + req.Search + "%"
		query = query.Where(`
			up.full_name ILIKE ? OR
			up.email ILIKE ? OR
			up.phone LIKE ?
		`, searchTerm, searchTerm, searchTerm)
	}

	// Filter: status
	if req.Status != nil && *req.Status != 0 {
		query = query.Where("COALESCE(up.status, 10) = ?", *req.Status)
	}

	// Filter: role ID
	if req.RoleID != nil && *req.RoleID > 0 {
		query = query.Where("up.role_id = ?", *req.RoleID)
	}

	// Filter: role type (nếu cần - hiện tại không có trường role_type trong user_profile)
	// if req.RoleType != "" {
	// 	query = query.Where("up.role_type = ?", req.RoleType)
	// }

	// Filter: date range
	if req.StartDate != nil {
		query = query.Where("up.created_at >= ?", req.StartDate)
	}
	if req.EndDate != nil {
		query = query.Where("up.created_at <= ?", req.EndDate)
	}

	// Filter theo bookmark nếu có filter
	if hasBookmarkFilter {
		if *req.Bookmark {
			query = query.Where("bu.user_id IS NOT NULL")
		} else {
			query = query.Where("bu.user_id IS NULL")
		}
	}

	// Get total count (use Session to avoid affecting the original query)
	err := query.Session(&gorm.Session{}).Count(&total).Error
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
		Debug().
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Find(&users).Error

	if err != nil {
		return nil, 0, err
	}

	// Lấy thông tin từ auth-service cho các user
	// if len(users) > 0 && r.AuthClient != nil {
	// 	profileIDs := make([]uint64, len(users))
	// 	for i, user := range users {
	// 		profileIDs[i] = user.ProfileID
	// 	}

	// 	// Gọi auth client để lấy thông tin status
	// 	statusMap, err := r.AuthClient.GetUserStatusBatch(ctx, profileIDs)
	// 	if err == nil {
	// 		// Cập nhật thông tin verifiedAt và lockedAt từ auth-service
	// 		for i := range users {
	// 			if statusInfo, exists := statusMap[users[i].ProfileID]; exists {
	// 				if statusInfo.VerifiedAt != nil {
	// 					// Parse string time to time.Time nếu cần
	// 					// users[i].VerifiedAt = parseTime(*statusInfo.VerifiedAt)
	// 				}
	// 				if statusInfo.LockedUntil != nil {
	// 					// Parse string time to time.Time nếu cần
	// 					// users[i].LockedAt = parseTime(*statusInfo.LockedUntil)
	// 				}
	// 			}
	// 		}
	// 	}
	// }

	return users, uint32(total), nil
}

// GetUserDetailByID lấy chi tiết user theo profile ID
func (r *ProfilePostgres) GetUserDetailByID(ctx context.Context, profileID uint64) (*models.UserProfileEntity, error) {
	var user models.UserProfileEntity
	err := r.DB.WithContext(ctx).
		Where("profile_id = ? AND deleted_at IS NULL", profileID).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUserProfile tạo user profile mới
func (r *ProfilePostgres) CreateUserProfile(ctx context.Context, profile *models.UserProfileEntity) (*models.UserProfileEntity, error) {
	err := r.DB.WithContext(ctx).Create(profile).Error
	if err != nil {
		return nil, err
	}
	return profile, nil
}

// UpdateUserProfile cập nhật thông tin user profile
func (r *ProfilePostgres) UpdateUserProfile(ctx context.Context, profile *models.UserProfileEntity) (*models.UserProfileEntity, error) {
	err := r.DB.WithContext(ctx).
		Where("profile_id = ? AND deleted_at IS NULL", profile.ProfileID).
		Updates(profile).Error
	if err != nil {
		return nil, err
	}
	return profile, nil
}

// DeleteUserProfile xóa vĩnh viễn user profile (hard delete)
// Trước khi xóa, lưu thông tin vào bảng profile_deleted
func (r *ProfilePostgres) DeleteUserProfile(ctx context.Context, profileID uint64, deletedBy *uint64, reason string) error {
	// Bắt đầu transaction
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Lấy thông tin profile trước khi xóa
		var profile models.UserProfileEntity
		err := tx.Where("profile_id = ?", profileID).First(&profile).Error
		if err != nil {
			return err
		}

		// 2. Tạo bản sao vào profile_deleted
		profileDeleted := models.NewProfileDeleted(&profile, deletedBy, reason)
		err = tx.Create(profileDeleted).Error
		if err != nil {
			return err
		}

		// 3. Xóa profile (hard delete)
		err = tx.Unscoped().
			Where("profile_id = ?", profileID).
			Delete(&models.UserProfileEntity{}).Error
		if err != nil {
			return err
		}

		return nil
	})
}

// UpdateUserProfileStatus cập nhật trạng thái user profile
func (r *ProfilePostgres) UpdateUserProfileStatus(ctx context.Context, profileID uint64, status _enum.EUserStatus) error {
	err := r.DB.WithContext(ctx).
		Model(&models.UserProfileEntity{}).
		Where("profile_id = ? AND deleted_at IS NULL", profileID).
		Update("status", status).Error
	return err
}
