package postgres

import (
	_dto "common/domain/dto"
	_enum "common/domain/enum"
	"common/identity"
	_utils "common/utils"
	"context"
	"fmt"
	hubpb "pb/types/hub"
	"strings"
	"time"
	"user/enums"
	"user/infra/client"
	"user/internal/dto"
	models "user/internal/models"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

// @bind: user/internal/interface/repo.IProfileRepo
type ProfilePostgres struct {
	DB        *gorm.DB
	HubClient *client.HubClient
}

func NewProfilePostgres(
	db *gorm.DB,
	hubClient *client.HubClient,
) *ProfilePostgres {
	return &ProfilePostgres{
		DB:        db,
		HubClient: hubClient,
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
		Order("updated_at DESC").
		Offset(pageable.GetOffset()).
		Limit(pageable.GetLimit()).
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
			"full_name":        dto.FullName,
			"email":            dto.Email,
			"avatar":           dto.Avatar,
			"background_image": dto.BackgroundImage,
			"phone":            dto.Phone,
			"birth":            dto.Birth,
			"start_date":       dto.StartDate,
			"address":          dto.Address,
			"phone2":           dto.Phone2,
			"tax_code":         dto.TaxCode,
			"position":         dto.Position,
			"workplace":        dto.Workplace,
			"role_title":       dto.RoleTitle,
			"department_id":    dto.DepartmentID,
			"slogan":           dto.Slogan,
			"website":          dto.Website,
			"facebook":         dto.Facebook,
			"instagram":        dto.Instagram,
			"twitter":          dto.Twitter,
			"linkedin":         dto.Linkedin,
			"youtube":          dto.Youtube,
			"introduction":     dto.Introduction,
			"front_identify":   dto.FrontIdentify,
			"back_identify":    dto.BackIdentify,
			"tab_default":      dto.TabDefault,
			"role_real_estate": dto.RoleRealEstate,

			// "visibility":         dto.Visibility,
			// "profile_visibility": dto.ProfileVisibility,
			// "status_online":      dto.StatusOnline,
		}).Error
}

func (r *ProfilePostgres) UpdatePrivacySetting(ctx context.Context, profileID uint64, visibility uint8, profileVisibility uint8, statusOnline enums.EOnlineStatus) error {
	return r.DB.WithContext(ctx).
		Model(&models.UserProfileEntity{}).
		Where("profile_id = ? AND deleted_at IS NULL", profileID).
		Updates(map[string]interface{}{
			"visibility":         visibility,
			"profile_visibility": profileVisibility,
			"status_online":      statusOnline,
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

func (r *ProfilePostgres) SearchByPhoneOrName(ctx context.Context, text string) ([]models.UserProfileEntity, error) {
	var profiles []models.UserProfileEntity
	query := "%" + text + "%"
	if err := r.DB.WithContext(ctx).
		Where("(phone ILIKE ? OR full_name ILIKE ?) AND deleted_at IS NULL", query, query).
		Find(&profiles).Error; err != nil {
		return nil, err
	}
	return profiles, nil
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

	// Admin routes use only the canonical actor promoted by strict ingress.
	// A missing actor cannot match any bookmark owner.
	adminID := uint64(0)
	if actor, ok := identity.ActorFromContext(ctx); ok {
		adminID = actor.ProfileID
	}

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
		session_agg.last_login_at,
		COALESCE(session_agg.total_login, 0) as total_login,
		NULL as job_title,
		NULL as attachments,
		NULL as internal_notes,
		NULL as send_notification`

	// Thêm bookmark field tùy theo trường hợp
	if needBookmarkInfo {
		selectFields += `,
		CASE WHEN bu.user_id IS NOT NULL THEN true ELSE false END as bookmark`
	} else {
		selectFields += `,
		false as bookmark`
	}

	// Tạo query
	query := r.DB.Debug().WithContext(ctx).Table("user_profile up").Select(selectFields).
		Joins(`LEFT JOIN (
			SELECT am.user_id,
				MAX(us.last_login) AS last_login_at,
				COUNT(us.session_id) AS total_login
			FROM user_session us
			INNER JOIN auth_method am ON us.auth_id = am.id
			GROUP BY am.user_id
		) session_agg ON session_agg.user_id = up.profile_id`)

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
	orderClause := r.parseListUsersSort(req)
	if orderClause == "" {
		orderClause = "up.created_at DESC, session_agg.last_login_at DESC NULLS LAST"
	}
	query = query.Order(orderClause)

	// Apply pagination
	page := req.Page
	if page < 1 {
		page = 1
	}
	limit := req.GetLimit()
	offset := int(page-1) * limit
	err = query.
		Debug().
		Offset(offset).
		Limit(limit).
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
	// Set tất cả visibility thành public mặc định nếu chưa được set
	if profile.Visibility == 0 {
		profile.Visibility = _enum.EVisibilityDTOPublic
	}
	if profile.ProfileVisibility == 0 {
		profile.ProfileVisibility = _enum.EVisibilityDTOPublic
	}
	if profile.VisibilityIntroduce == 0 {
		profile.VisibilityIntroduce = _enum.EVisibilityDTOPublic
	}
	if profile.VisibilityProfession == 0 {
		profile.VisibilityProfession = _enum.EVisibilityDTOPublic
	}
	if profile.VisibilityMainArea == 0 {
		profile.VisibilityMainArea = _enum.EVisibilityDTOPublic
	}
	if profile.VisibilityFriends == 0 {
		profile.VisibilityFriends = _enum.EVisibilityDTOPublic
	}
	if profile.VisibilitySignature == 0 {
		profile.VisibilitySignature = _enum.EVisibilityDTOPublic
	}

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

// IsSystemRootProfile identifies the one bootstrap-owned IAM principal. This
// is queried by administrative mutations so ordinary USER_ADMIN_MANAGE access
// cannot disable or rewrite the root operator.
func (r *ProfilePostgres) IsSystemRootProfile(ctx context.Context, profileID uint64) (bool, error) {
	var exists bool
	err := r.DB.WithContext(ctx).Raw(`
SELECT EXISTS (
    SELECT 1
      FROM role_profiles assignment
      JOIN roles role ON role.id = assignment.role_id
     WHERE assignment.profile_id = ?
       AND assignment.deleted_at IS NULL
       AND role.organization_id = 0
       AND role.key = 'QHPRO_SYSTEM_ROOT'
       AND role.deleted_at IS NULL
)`, profileID).Scan(&exists).Error
	return exists, err
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

// GetByPhone lấy profile theo số điện thoại
func (r *ProfilePostgres) GetByPhone(phone string) (*models.UserProfileEntity, error) {
	var profile models.UserProfileEntity
	err := r.DB.Model(&models.UserProfileEntity{}).
		Where("phone = ? AND deleted_at IS NULL", phone).
		First(&profile).Error
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

// UpdateLastSeen cập nhật thời gian hoạt động gần nhất của user
func (r *ProfilePostgres) UpdateLastSeen(ctx context.Context, profileID uint64) error {
	now := time.Now()
	err := r.DB.WithContext(ctx).
		Model(&models.UserProfileEntity{}).
		Where("profile_id = ? AND deleted_at IS NULL", profileID).
		Update("last_seen_at", now).Error
	return err
}

// UpdatePerson cập nhật thông tin cá nhân cơ bản
func (r *ProfilePostgres) UpdatePerson(c context.Context, profileId uint64, req *_dto.UserV3DTO) error {
	updates := make(map[string]interface{})

	// Tạo set từ FieldSets để check nhanh hơn
	fieldSetMap := make(map[string]bool)
	for _, field := range req.FieldSets {
		fieldSetMap[field] = true
	}

	// Helper function để check xem field có trong FieldSets không
	// Nếu FieldSets rỗng, cho phép update tất cả (backward compatible)
	shouldUpdate := func(fieldName string) bool {
		// Nếu không có FieldSets thì cho phép update tất cả các trường (backward compatible)
		if len(fieldSetMap) == 0 {
			return true
		}
		// Nếu có FieldSets thì chỉ update nếu fieldName có trong đó
		_, ok := fieldSetMap[fieldName]
		return ok
	}

	// Chỉ update các field được gửi lên (không nil) và có trong FieldSets (nếu có)
	if req.FullName != "" && shouldUpdate("fullName") {
		updates["full_name"] = req.FullName
	}
	if req.Email != "" || shouldUpdate("email") {
		updates["email"] = req.Email
	}
	if req.ProvinceID != nil || shouldUpdate("provinceId") {
		updates["province_id"] = req.ProvinceID
	}
	if req.WardID != nil || shouldUpdate("wardId") {
		updates["ward_id"] = req.WardID
	}
	if req.Gender != 0 || shouldUpdate("gender") {
		updates["gender"] = req.Gender
	}
	if req.Birth != nil || shouldUpdate("birth") {
		updates["birth"] = req.Birth
	}
	if req.Phone2 != "" || shouldUpdate("phone2") {
		updates["phone2"] = req.Phone2
	}
	if req.RoleRealEstate != 0 || shouldUpdate("roleRealEstate") {
		updates["role_real_estate"] = req.RoleRealEstate
	}
	if req.SignatureVisible || shouldUpdate("signatureVisible") {
		updates["signature_visible"] = true
	}
	if req.ZaloURL != "" || shouldUpdate("zaloUrl") {
		updates["zalo_url"] = req.ZaloURL
	}
	if req.FacebookURL != "" || shouldUpdate("facebookUrl") {
		updates["facebook_url"] = req.FacebookURL
	}
	if req.InstagramURL != "" || shouldUpdate("instagramUrl") {
		updates["instagram_url"] = req.InstagramURL
	}
	if req.TwitterURL != "" || shouldUpdate("twitterUrl") {
		updates["twitter_url"] = req.TwitterURL
	}
	if req.LinkedinURL != "" || shouldUpdate("linkedinUrl") {
		updates["linkedin_url"] = req.LinkedinURL
	}
	if req.YoutubeURL != "" || shouldUpdate("youtubeUrl") {
		updates["youtube_url"] = req.YoutubeURL
	}
	if req.WebsiteURL != "" || shouldUpdate("websiteUrl") {
		updates["website_url"] = req.WebsiteURL
	}
	if req.Introduction != "" || shouldUpdate("introduction") {
		updates["introduction"] = req.Introduction
	}
	if req.Position != "" || shouldUpdate("position") {
		updates["position"] = req.Position
	}
	if req.Workplace != "" || shouldUpdate("workplace") {
		updates["workplace"] = req.Workplace
	}
	if req.RoleTitle != "" || shouldUpdate("roleTitle") {
		updates["role_title"] = req.RoleTitle
	}
	if req.DepartmentID != 0 || shouldUpdate("departmentId") {
		updates["department_id"] = req.DepartmentID
	}
	if req.Avatar != "" && shouldUpdate("avatar") {
		updates["avatar"] = req.Avatar
	}
	if req.BackgroundImage != "" && shouldUpdate("backgroundImage") {
		updates["background_image"] = req.BackgroundImage
	}
	if req.Address != nil {
		if req.Address.Detail != "" && shouldUpdate("address") {
			updates["address"] = req.Address.Detail
		}
		if req.Address.ProvinceID != nil && shouldUpdate("provinceId") {
			updates["province_id"] = *req.Address.ProvinceID
		}
		if req.Address.WardID != nil && shouldUpdate("wardId") {
			updates["ward_id"] = *req.Address.WardID
		}
	}

	// Cập nhật các visibility fields
	if req.VisibilityIntroduce != 0 && shouldUpdate("visibilityIntroduce") {
		updates["visibility_introduce"] = req.VisibilityIntroduce
	}
	if req.VisibilityProfession != 0 && shouldUpdate("visibilityProfession") {
		updates["visibility_profession"] = req.VisibilityProfession
	}
	if req.VisibilityMainArea != 0 && shouldUpdate("visibilityMainArea") {
		updates["visibility_main_area"] = req.VisibilityMainArea
	}
	if req.VisibilityFriends != 0 && shouldUpdate("visibilityFriends") {
		updates["visibility_friends"] = req.VisibilityFriends
	}
	if req.VisibilitySignature != 0 && shouldUpdate("visibilitySignature") {
		updates["visibility_signature"] = req.VisibilitySignature
	}
	// Cập nhật ViewRoles (integer array) - Convert []uint32 sang pq.Int32Array
	if len(req.ViewRoles) > 0 && shouldUpdate("viewRoles") {
		viewRoles := make(pq.Int32Array, len(req.ViewRoles))
		for i, role := range req.ViewRoles {
			viewRoles[i] = int32(role)
		}
		updates["view_roles"] = viewRoles
	}

	// Nếu không có gì để update
	if len(updates) == 0 {
		return nil
	}

	err := r.DB.WithContext(c).
		Model(&models.UserProfileEntity{}).
		Where("profile_id = ? AND deleted_at IS NULL", profileId).
		Updates(updates).Error

	return err
}

func (r *ProfilePostgres) PatchProfileV3(c context.Context, profileId uint64, req *_dto.UserV3DTO) error {
	updates := make(map[string]interface{})

	if req.RoleRealEstate != 0 {
		updates["role_real_estate"] = req.RoleRealEstate
	}

	return r.DB.WithContext(c).
		Model(&models.UserProfileEntity{}).
		Where("profile_id = ? AND deleted_at IS NULL", profileId).
		Updates(updates).Error
}

// GetUserGuideByKey lấy user guide theo key từ hub service
func (r *ProfilePostgres) GetUserGuideByKey(ctx context.Context, key string) (*hubpb.UserGuideDetail, error) {
	if r.HubClient == nil {
		return nil, fmt.Errorf("hub client not available")
	}
	return r.HubClient.GetUserGuideByKey(ctx, key)
}

// parseListUsersSort parses sort string (format: "field,order") to SQL ORDER BY clause.
// Ví dụ: "createdAt,desc" -> "up.created_at DESC, session_agg.last_login_at DESC NULLS LAST"
func (r *ProfilePostgres) parseListUsersSort(req *dto.ListUsersRequest) string {
	sort := strings.TrimSpace(req.GetSort())
	if sort == "" {
		sort = strings.TrimSpace(req.Sort)
	}
	if sort == "" && req.SortBy != "" {
		order := req.SortOrder
		if order == "" {
			order = "desc"
		}
		sort = req.SortBy + "," + order
	}
	if sort == "" {
		return ""
	}

	parts := strings.Split(sort, ",")
	if len(parts) != 2 {
		return ""
	}

	field := strings.TrimSpace(parts[0])
	order := strings.TrimSpace(strings.ToUpper(parts[1]))
	if order != "ASC" && order != "DESC" {
		return ""
	}

	fieldMap := map[string]string{
		"createdAt":   "up.created_at",
		"updatedAt":   "up.updated_at",
		"fullName":    "up.full_name",
		"email":       "up.email",
		"lastLogin":   "session_agg.last_login_at",
		"lastLoginAt": "session_agg.last_login_at",
	}

	dbField, exists := fieldMap[field]
	if !exists {
		return ""
	}

	nullsLast := ""
	if field == "lastLogin" || field == "lastLoginAt" {
		nullsLast = " NULLS LAST"
	}

	orderClause := fmt.Sprintf("%s %s%s", dbField, order, nullsLast)
	if field == "createdAt" {
		orderClause += ", session_agg.last_login_at DESC NULLS LAST"
	}

	return orderClause
}
