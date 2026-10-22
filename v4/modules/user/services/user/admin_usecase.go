package services

import (
	_enum "common/domain/enum"
	_errors "common/errors"
	"context"
	"fmt"
	"mdn4/internal/interface/providers"
	"mdn4/internal/interface/repo"
	authdto "mdn4/modules/auth/dto"
	"mdn4/modules/user/domain/user"
	"mdn4/modules/user/dto"
	"mdn4/modules/user/entity"
	"mdn4/modules/user/services/useradmin"
	"time"
)

type AdminUsecase struct {
	adminRepo          repo.IAdminRepo
	profileRepo        repo.IProfileRepo
	profileAccountRepo repo.ProfileAccountRepository
	authProvider       providers.AuthProvider
}

func adminActorID(ctx context.Context) (uint64, error) {
	actorID, err := useradmin.ActorIDFromContext(ctx)
	if err != nil {
		return 0, _errors.ReturnError(int32(401), "Không có quyền truy cập")
	}
	return actorID, nil
}

func NewAdminUsecase(
	adminRepo repo.IAdminRepo,
	profileRepo repo.IProfileRepo,
	profileAccountRepo repo.ProfileAccountRepository,
	authProvider providers.AuthProvider,
) *AdminUsecase {
	return &AdminUsecase{
		adminRepo:          adminRepo,
		profileRepo:        profileRepo,
		profileAccountRepo: profileAccountRepo,
		authProvider:       authProvider,
	}
}

// BlockUser khóa tài khoản user với lý do
func (uc *AdminUsecase) LockUser(ctx context.Context, req *dto.LockUserRequest) error {
	if err := uc.rejectSystemRootMutation(ctx, req.ProfileID); err != nil {
		return err
	}

	_, err := uc.profileRepo.GetUserDetailByID(ctx, req.ProfileID)
	if err != nil {
		return fmt.Errorf("user not found: %v", err)
	}

	actorID, err := adminActorID(ctx)
	if err != nil {
		actorID = 0
	}

	// Temporary lock without explicit duration stays locked until unlock.
	durationHours := 0
	if req.LockType == _enum.EUserStatusTemporaryLocked {
		durationHours = 24
	}

	if err := uc.profileAccountRepo.LockUser(ctx, req.ProfileID, req.LockType, durationHours, req.Reason, actorID); err != nil {
		return fmt.Errorf("failed to update user profile status: %v", err)
	}
	return nil
}

// UnblockUser mở khóa tài khoản user với lý do
func (uc *AdminUsecase) UnLockUser(ctx context.Context, req *dto.UnLockUserRequest) error {
	if err := uc.rejectSystemRootMutation(ctx, req.ProfileID); err != nil {
		return err
	}
	_, err := uc.profileRepo.GetUserDetailByID(ctx, req.ProfileID)
	if err != nil {
		return fmt.Errorf("user not found: %v", err)
	}

	actorID, err := adminActorID(ctx)
	if err != nil {
		actorID = 0
	}

	if err := uc.profileAccountRepo.UnlockUser(ctx, req.ProfileID, req.Reason, actorID); err != nil {
		return fmt.Errorf("failed to update user profile status: %v", err)
	}
	return nil
}

// ApproveUser duyệt tài khoản user
func (uc *AdminUsecase) ApproveUser(ctx context.Context, req *dto.ApproveUserRequest) error {
	if err := uc.rejectSystemRootMutation(ctx, req.ProfileID); err != nil {
		return err
	}
	_, err := uc.profileRepo.GetUserDetailByID(ctx, req.ProfileID)
	if err != nil {
		return fmt.Errorf("user not found: %v", err)
	}

	if err := uc.profileAccountRepo.UnlockUser(ctx, req.ProfileID, req.Note, 0); err != nil {
		return fmt.Errorf("failed to update user profile status: %v", err)
	}
	return nil
}

// CreateAdmin tạo tài khoản admin mới bằng cách gọi auth service
func (u *AdminUsecase) CreateAdmin(ctx context.Context, req *entity.AdminProfile) (*entity.AdminProfile, error) {
	// Kiểm tra quyền admin
	if _, err := adminActorID(ctx); err != nil {
		return nil, err
	}

	if req.Username == "" {
		return nil, _errors.ReturnError(int32(400), "Username không được để trống")
	}
	if req.Password == "" {
		return nil, _errors.ReturnError(int32(400), "Password không được để trống")
	}

	// Convert response từ auth service
	createAdminResp, err := u.adminRepo.CreateAdmin(ctx, req)
	if err != nil {
		return nil, _errors.ReturnError(int32(500), "Lỗi tạo hồ sơ admin")
	}

	// Link auth.user_id ↔ admin_profiles.id (BaseEntity.ID). Fallback AuthID nếu ID chưa được GORM fill.
	linkUserID := createAdminResp.ID
	if linkUserID == 0 {
		linkUserID = createAdminResp.AuthID
	}
	if linkUserID == 0 {
		return nil, _errors.ReturnError(int32(500), "Không lấy được ID hồ sơ admin sau khi tạo")
	}

	// Gọi auth service để tạo/cập nhật credential đăng nhập (auth_name + password)
	createdAuth, err := u.authProvider.Update(ctx, &providers.AuthMethodDomain{
		UserID:   linkUserID,
		AuthName: req.Username,
		Password: req.Password,
		FullName: req.FullName,
		Email:    req.Email,
		Phone:    req.Phone,
		Avatar:   req.Avatar,
		Provider: "ADMIN",
	})
	if err != nil {
		fmt.Printf("Error updating admin via auth service: %v\n", err)
		return nil, _errors.ReturnError(int32(500), "Lỗi tạo thông tin đăng nhập admin")
	}

	createAdminResp.Username = req.Username
	if createdAuth != nil && createdAuth.ID != 0 {
		createAdminResp.AuthID = createdAuth.ID
		if err := u.adminRepo.UpdateAuthID(ctx, linkUserID, createdAuth.ID); err != nil {
			fmt.Printf("WARN: failed to persist auth_id on admin profile: %v\n", err)
		}
	}

	// Gán role vào role_profiles — bắt buộc để permission/me và HasPermissions hoạt động
	if req.RoleID != nil && *req.RoleID > 0 {
		if err := u.authProvider.AssignRoleToUser(ctx, linkUserID, *req.RoleID); err != nil {
			fmt.Printf("Error assigning role via auth service: %v\n", err)
			return nil, _errors.ReturnError(int32(500), "Tạo tài khoản thành công nhưng gán role thất bại: "+err.Error())
		}
	}

	return createAdminResp, nil
}

// UpdateAdmin cập nhật thông tin admin bằng cách gọi auth service
func (u *AdminUsecase) UpdateAdmin(ctx context.Context, req *entity.AdminProfile) (*entity.AdminProfile, error) {
	// Kiểm tra quyền admin
	if _, err := adminActorID(ctx); err != nil {
		return nil, err
	}
	if err := u.rejectSystemRootAdminMutation(ctx, req.ID); err != nil {
		return nil, err
	}

	linkUserID := req.ID
	if linkUserID == 0 {
		linkUserID = req.AuthID
	}

	// Gọi auth service để cập nhật admin
	_, err := u.authProvider.Update(ctx, &providers.AuthMethodDomain{
		UserID:   linkUserID,
		AuthName: req.Username,
		Password: req.Password,
		FullName: req.FullName,
		Email:    req.Email,
		Phone:    req.Phone,
		Avatar:   req.Avatar,
		Provider: "ADMIN",
	})
	if err != nil {
		fmt.Printf("Error updating admin via auth service: %v\n", err)
		return nil, _errors.ReturnError(int32(500), "Lỗi cập nhật thông tin admin")
	}

	result, err := u.adminRepo.UpdateAdmin(ctx, req)
	if err != nil {
		fmt.Printf("Error updating admin profile in DB: %v\n", err)
		return nil, _errors.ReturnError(int32(500), "Lỗi cập nhật hồ sơ admin: "+err.Error())
	}

	if req.RoleID != nil && *req.RoleID > 0 && linkUserID > 0 {
		if err := u.authProvider.AssignRoleToUser(ctx, linkUserID, *req.RoleID); err != nil {
			fmt.Printf("Error assigning role via auth service: %v\n", err)
			return nil, _errors.ReturnError(int32(500), "Cập nhật thành công nhưng gán role thất bại: "+err.Error())
		}
	}

	return result, nil
}

// DeleteAdmin xóa tài khoản admin bằng cách gọi auth service
func (u *AdminUsecase) DeleteAdmin(ctx context.Context, req uint64) error {
	// Kiểm tra quyền admin
	if _, err := adminActorID(ctx); err != nil {
		return err
	}
	if err := u.rejectSystemRootAdminMutation(ctx, req); err != nil {
		return err
	}

	err := u.authProvider.SoftDelete(ctx, req)
	if err != nil {
		return _errors.ReturnError(int32(500), "Lỗi xóa tài khoản admin")
	}

	err = u.adminRepo.DeleteAdmin(ctx, req)
	if err != nil {
		return _errors.ReturnError(int32(500), "Lỗi xóa tài khoản admin")
	}

	return nil
}

func (u *AdminUsecase) rejectSystemRootAdminMutation(ctx context.Context, id uint64) error {
	if u == nil || u.adminRepo == nil {
		return _errors.ReturnError(int32(503), "Kho dữ liệu quản trị chưa sẵn sàng")
	}
	isRoot, err := u.adminRepo.IsSystemRootAdmin(ctx, id)
	if err != nil {
		return _errors.ReturnError(int32(500), "Không thể xác minh tài khoản quản trị gốc")
	}
	if isRoot {
		return _errors.ReturnError(int32(403), "Không thể thay đổi tài khoản quản trị gốc qua API")
	}
	return nil
}

// ListAdmins lấy danh sách admin bằng cách gọi auth service
func (u *AdminUsecase) ListAdmins(ctx context.Context, req *dto.AdminListRequest) ([]entity.AdminProfile, int64, error) {
	// Kiểm tra quyền admin
	if _, err := adminActorID(ctx); err != nil {
		return nil, 0, err
	}

	// Lấy danh sách admin từ repo
	authResp, total, err := u.adminRepo.ListAdmin(ctx, req)
	if err != nil {
		fmt.Printf("Error listing admins via auth service: %v\n", err)
		return nil, 0, _errors.ReturnError(int32(500), "Lỗi lấy danh sách admin")
	}

	// Lấy danh sách roleIds từ admins (loại bỏ duplicate và nil)
	roleIdsMap := make(map[uint64]bool)
	for _, admin := range authResp {
		if admin.RoleID != nil && *admin.RoleID > 0 {
			roleIdsMap[*admin.RoleID] = true
		}
	}

	// Convert map sang slice
	roleIds := make([]uint64, 0, len(roleIdsMap))
	for roleId := range roleIdsMap {
		roleIds = append(roleIds, roleId)
	}

	// Gọi auth service để lấy roles theo IDs
	var roleMap map[uint64]*authdto.ProviderRoleDTO
	if len(roleIds) > 0 {
		roleMap, err = u.authProvider.GetRolesByIds(ctx, roleIds)
		if err != nil {
			fmt.Printf("Error getting roles by IDs: %v\n", err)
			roleMap = make(map[uint64]*authdto.ProviderRoleDTO)
		}
	}

	for i := range authResp {
		if authResp[i].RoleID != nil && *authResp[i].RoleID > 0 {
			if role, ok := roleMap[*authResp[i].RoleID]; ok {
				authResp[i].Role = providerRoleToRoleDTO(role)
			}
		}
	}

	fmt.Printf("Listed %d admins with %d roles\n", len(authResp), len(roleMap))

	return authResp, total, nil
}

func (u *AdminUsecase) GetDetail(ctx context.Context, req uint64) (*entity.AdminProfile, error) {
	// Kiểm tra quyền admin
	authResp, _ := u.authProvider.FindByID(ctx, req)

	detail, err := u.adminRepo.GetDetail(ctx, req)
	if err != nil {
		return nil, _errors.ReturnError(int32(500), "Lỗi lấy thông tin admin")
	}
	if authResp != nil {
		detail.Username = authResp.AuthName
	}

	// Lấy thông tin role nếu có roleId
	if detail.RoleID != nil && *detail.RoleID > 0 {
		roleMap, err := u.authProvider.GetRolesByIds(ctx, []uint64{*detail.RoleID})
		if err != nil {
			fmt.Printf("Error getting role by ID: %v\n", err)
			// Không return error, chỉ log warning vì role là thông tin bổ sung
		} else {
			if role, ok := roleMap[*detail.RoleID]; ok {
				detail.Role = providerRoleToRoleDTO(role)
			}
		}
	}

	fmt.Printf("GetDetail admin ID %d with role: %+v\n", req, detail.Role)

	return detail, nil
}

// IsSuperAdmin reports whether the admin profile for profileID has durable super_admin=1.
func (u *AdminUsecase) IsSuperAdmin(ctx context.Context, profileID uint64) (bool, error) {
	if u == nil || u.adminRepo == nil || profileID == 0 {
		return false, nil
	}
	return u.adminRepo.IsSuperAdminByProfileID(ctx, profileID)
}

// CreateUser tạo user mới giống với tài khoản admin
func (u *AdminUsecase) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.CreateUserResponse, error) {
	// Kiểm tra quyền admin
	if _, err := adminActorID(ctx); err != nil {
		return nil, err
	}

	// Kiểm tra email đã tồn tại chưa
	existingAuth, _ := u.authProvider.FindByEmail(ctx, req.Email)
	if existingAuth != nil {
		return nil, _errors.ReturnError(int32(400), "Email đã tồn tại trong hệ thống")
	}
	// Tự động generate username từ email nếu không có
	username := req.Username
	if username == "" {
		// Lấy phần trước @ của email làm username
		atIndex := 0
		for i, c := range req.Email {
			if c == '@' {
				atIndex = i
				break
			}
		}
		if atIndex > 0 {
			username = req.Email[:atIndex]
		} else {
			username = req.Email
		}
	}
	if existingAuth, _ := u.authProvider.FindByAuthNameAndProvider(ctx, username, "ADMIN"); existingAuth != nil {
		return nil, _errors.ReturnError(int32(400), "Tên đăng nhập đã tồn tại trong hệ thống")
	}

	password := req.Password
	if password == "" {
		return nil, _errors.ReturnError(int32(400), "Mật khẩu là bắt buộc")
	}

	// Tạo auth method trong auth service
	authMethod := &providers.AuthMethodDomain{
		AuthName: username,
		Password: password,
		FullName: req.FullName,
		Email:    req.Email,
		Phone:    req.Phone,
		Avatar:   req.Avatar,
		Provider: "ADMIN",
	}

	createdAuth, err := u.authProvider.Create(ctx, authMethod)
	if err != nil {
		fmt.Printf("Error creating auth method: %v\n", err)
		return nil, _errors.ReturnError(int32(500), "Lỗi tạo tài khoản xác thực")
	}

	// Tạo user profile trong user service (dùng UserProfileEntity)
	userProfile := &user.UserProfileEntity{
		// PostgreSQL owns profile identity generation. Reusing the operator's
		// identity here would overwrite or alias the admin account.
		ProfileID: 0,
		FullName:  req.FullName,
		Email:     req.Email,
		Phone:     req.Phone,
		Avatar:    req.Avatar,
		RoleID:    &req.RoleID,
		Gender:    0,
		Address:   req.Address,
		Birth:     req.Birth,
		Status:    _enum.EUserStatusActive,
	}

	// Set gender if provided
	if req.Gender != nil && *req.Gender > 0 {
		userProfile.Gender = *req.Gender
	}

	// Set tất cả visibility thành public mặc định
	user.SetAllVisibilityPublic(userProfile)

	createdProfile, err := u.profileRepo.CreateUserProfile(ctx, userProfile)
	if err != nil {
		fmt.Printf("Error creating user profile: %v\n", err)
		// Rollback: xóa auth method đã tạo
		_ = u.authProvider.SoftDelete(ctx, createdAuth.ID)
		return nil, _errors.ReturnError(int32(500), "Lỗi tạo hồ sơ người dùng")
	}

	// Link the credential to the canonical profile generated by PostgreSQL.
	createdAuth.UserID = createdProfile.ProfileID
	createdAuth, err = u.authProvider.Update(ctx, createdAuth)
	if err != nil {
		_ = u.authProvider.SoftDelete(ctx, createdAuth.ID)
		return nil, _errors.ReturnError(int32(500), "Lỗi liên kết tài khoản xác thực với hồ sơ người dùng")
	}
	if req.RoleID > 0 {
		if err := u.authProvider.AssignRoleToUser(ctx, createdProfile.ProfileID, req.RoleID); err != nil {
			return nil, _errors.ReturnError(int32(500), "Tạo người dùng thành công nhưng gán vai trò thất bại")
		}
	}

	return &dto.CreateUserResponse{
		ProfileID: createdProfile.ProfileID,
		AuthID:    createdAuth.ID,
		Username:  createdAuth.AuthName,
		FullName:  createdProfile.FullName,
		Email:     createdProfile.Email,
		Phone:     createdProfile.Phone,
		Avatar:    createdProfile.Avatar,
		RoleID:    *createdProfile.RoleID,
	}, nil
}

// UpdateUser cập nhật thông tin user
func (u *AdminUsecase) UpdateUser(ctx context.Context, req *dto.UpdateUserRequest) (*dto.UpdateUserResponse, error) {
	// Kiểm tra quyền admin
	if _, err := adminActorID(ctx); err != nil {
		return nil, err
	}
	if err := u.rejectSystemRootMutation(ctx, req.ProfileID); err != nil {
		return nil, err
	}

	// Kiểm tra user có tồn tại không
	existingUser, err := u.profileRepo.GetUserDetailByID(ctx, req.ProfileID)
	if err != nil {
		return nil, _errors.ReturnError(int32(404), "Không tìm thấy người dùng")
	}

	// Cập nhật thông tin trong user profile
	updatedProfile := &user.UserProfileEntity{
		ProfileID: req.ProfileID,
		RoleID:    req.RoleID,
		Address:   req.Address,
		Birth:     req.Birth,
	}

	// Set gender if provided
	if req.Gender != nil && *req.Gender > 0 {
		updatedProfile.Gender = *req.Gender
	} else {
		updatedProfile.Gender = existingUser.Gender
	}

	// Chỉ cập nhật các trường không rỗng
	if req.FullName != "" {
		updatedProfile.FullName = req.FullName
	} else {
		updatedProfile.FullName = existingUser.FullName
	}

	if req.Email != "" {
		updatedProfile.Email = req.Email
	} else {
		updatedProfile.Email = existingUser.Email
	}

	if req.Phone != "" {
		updatedProfile.Phone = req.Phone
	} else {
		updatedProfile.Phone = existingUser.Phone
	}

	if req.Avatar != "" {
		updatedProfile.Avatar = req.Avatar
	} else {
		updatedProfile.Avatar = existingUser.Avatar
	}

	// Set status if provided
	if req.Status != 0 {
		updatedProfile.Status = req.Status
	} else {
		updatedProfile.Status = existingUser.Status
	}

	result, err := u.profileRepo.UpdateUserProfile(ctx, updatedProfile)
	if err != nil {
		return nil, _errors.ReturnError(int32(500), "Lỗi cập nhật thông tin người dùng")
	}

	// Cập nhật thông tin trong auth service
	_, err = u.authProvider.Update(ctx, &providers.AuthMethodDomain{
		UserID:   req.ProfileID,
		FullName: result.FullName,
		Email:    result.Email,
		Phone:    result.Phone,
		Avatar:   result.Avatar,
	})
	if err != nil {
		fmt.Printf("Error updating auth method: %v\n", err)
		// Không return error vì profile đã được cập nhật
	}

	roleID := uint64(0)
	if result.RoleID != nil {
		roleID = *result.RoleID
	}

	return &dto.UpdateUserResponse{
		ProfileID: result.ProfileID,
		FullName:  result.FullName,
		Email:     result.Email,
		Phone:     result.Phone,
		Avatar:    result.Avatar,
		RoleID:    roleID,
		Status:    result.Status,
		UpdatedAt: result.UpdatedAt,
	}, nil
}

// DeleteUser xóa user
func (u *AdminUsecase) DeleteUser(ctx context.Context, req *dto.DeleteUserRequest) error {
	// Kiểm tra quyền admin
	adminProfileID, err := adminActorID(ctx)
	if err != nil {
		return err
	}
	if err := u.rejectSystemRootMutation(ctx, req.ProfileID); err != nil {
		return err
	}

	// Kiểm tra user có tồn tại không
	_, err = u.profileRepo.GetUserDetailByID(ctx, req.ProfileID)
	if err != nil {
		return _errors.ReturnError(int32(404), "Không tìm thấy người dùng")
	}

	// Xóa user profile (soft delete)
	// Lấy thông tin admin đang thực hiện xóa
	reason := "Xóa bởi admin"
	if req.Reason != "" {
		reason = req.Reason
	}
	err = u.profileRepo.DeleteUserProfile(ctx, req.ProfileID, &adminProfileID, reason)
	if err != nil {
		return _errors.ReturnError(int32(500), "Lỗi xóa hồ sơ người dùng")
	}

	// Xóa TẤT CẢ auth method của user trong auth service
	err = u.authProvider.DeleteAllAuthMethodsByUserId(ctx, req.ProfileID)
	if err != nil {
		fmt.Printf("Error deleting all auth methods for user %d: %v\n", req.ProfileID, err)
		// Không return error vì profile đã được xóa
	}

	return nil
}

func (u *AdminUsecase) rejectSystemRootMutation(ctx context.Context, profileID uint64) error {
	if u == nil || u.profileRepo == nil {
		return _errors.ReturnError(int32(503), "Kho dữ liệu người dùng chưa sẵn sàng")
	}
	isRoot, err := u.profileRepo.IsSystemRootProfile(ctx, profileID)
	if err != nil {
		return _errors.ReturnError(int32(500), "Không thể xác minh tài khoản quản trị gốc")
	}
	if isRoot {
		return _errors.ReturnError(int32(403), "Không thể thay đổi tài khoản quản trị gốc qua API")
	}
	return nil
}

// GetUserDetail lấy chi tiết thông tin user
func (u *AdminUsecase) GetUserDetail(ctx context.Context, profileID uint64) (*dto.GetUserDetailResponse, error) {
	// Kiểm tra quyền admin
	if _, err := adminActorID(ctx); err != nil {
		return nil, err
	}

	// Lấy thông tin user từ profile repo
	user, err := u.profileRepo.GetUserDetailByID(ctx, profileID)
	if err != nil {
		return nil, _errors.ReturnError(int32(404), "Không tìm thấy người dùng")
	}

	// Lấy thông tin role nếu có
	var roleName, roleColor, roleBgColor string
	var roleKey uint32
	if user.RoleID != nil && *user.RoleID > 0 {
		roleMap, err := u.authProvider.GetRolesByIds(ctx, []uint64{*user.RoleID})
		if err == nil {
			if role, ok := roleMap[*user.RoleID]; ok {
				roleName = role.RoleName
				roleKey = role.RoleKey
				// Assuming role has color fields
			}
		}
	}

	gender := uint32(0)
	if user.Gender > 0 {
		gender = user.Gender
	}

	return &dto.GetUserDetailResponse{
		ProfileID:       user.ProfileID,
		FullName:        user.FullName,
		Email:           user.Email,
		Phone:           user.Phone,
		Phone2:          user.Phone2,
		Avatar:          user.Avatar,
		Address:         user.Address,
		Gender:          gender,
		Birth:           user.Birth,
		Status:          user.Status,
		Position:        user.Position,
		Workplace:       user.Workplace,
		DepartmentID:    user.DepartmentID,
		TickVerified:    user.TickVerified,
		BackgroundImage: user.BackgroundImage,
		TaxCode:         user.TaxCode,
		FrontIdentify:   user.FrontIdentify,
		BackIdentify:    user.BackIdentify,
		Slogan:          user.Slogan,
		Website:         user.WebsiteURL,
		Facebook:        user.FacebookURL,
		Instagram:       user.InstagramURL,
		Twitter:         user.TwitterURL,
		Linkedin:        user.LinkedinURL,
		Youtube:         user.YoutubeURL,
		Introduction:    user.Introduction,
		RoleID:          user.RoleID,
		RoleName:        roleName,
		RoleKey:         roleKey,
		RoleColor:       roleColor,
		RoleBgColor:     roleBgColor,
	}, nil
}

// GetUserGrowthDailyStats thống kê user mới (created_at) và user hoạt động (login/refresh) theo gap dates|hours.
func (u *AdminUsecase) GetUserGrowthDailyStats(ctx context.Context, req *dto.GetUserGrowthDailyStatsRequest) (*dto.GetUserGrowthDailyStatsResponse, error) {
	if req == nil {
		req = &dto.GetUserGrowthDailyStatsRequest{}
	}
	req.Normalize()

	fromTime := userGrowthFromTime(req.LastDates)

	var (
		newRows    []dto.UserGrowthDailyItem
		activeRows []dto.UserGrowthDailyItem
		err        error
		data       []dto.UserGrowthDailyItem
	)

	if req.Gap == dto.UserGrowthGapHours {
		newRows, err = u.profileRepo.GetNewUsersHourlyStats(ctx, fromTime)
		if err != nil {
			return nil, _errors.ReturnError(int32(500), "Lỗi khi thống kê user mới theo giờ: "+err.Error())
		}
		activeRows, err = u.profileRepo.GetActiveUsersHourlyStats(ctx, fromTime)
		if err != nil {
			return nil, _errors.ReturnError(int32(500), "Lỗi khi thống kê user hoạt động theo giờ: "+err.Error())
		}
		data = buildUserGrowthHourlySeries(fromTime, newRows, activeRows)
	} else {
		newRows, err = u.profileRepo.GetNewUsersDailyStats(ctx, fromTime)
		if err != nil {
			return nil, _errors.ReturnError(int32(500), "Lỗi khi thống kê user mới: "+err.Error())
		}
		activeRows, err = u.profileRepo.GetActiveUsersDailyStats(ctx, fromTime)
		if err != nil {
			return nil, _errors.ReturnError(int32(500), "Lỗi khi thống kê user hoạt động: "+err.Error())
		}
		data = buildUserGrowthDailySeries(fromTime, req.LastDates, newRows, activeRows)
	}

	return &dto.GetUserGrowthDailyStatsResponse{
		Data:      data,
		LastDates: req.LastDates,
		Gap:       req.Gap,
	}, nil
}

func userGrowthFromTime(lastDates uint32) time.Time {
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		loc = time.FixedZone("Asia/Ho_Chi_Minh", 7*60*60)
	}
	now := time.Now().In(loc)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	return today.AddDate(0, 0, -int(lastDates)+1)
}

func buildUserGrowthDailySeries(
	fromTime time.Time,
	lastDates uint32,
	newRows []dto.UserGrowthDailyItem,
	activeRows []dto.UserGrowthDailyItem,
) []dto.UserGrowthDailyItem {
	newByDate := make(map[string]uint64, len(newRows))
	for _, row := range newRows {
		newByDate[row.Date] = row.NewUsers
	}

	activeByDate := make(map[string]uint64, len(activeRows))
	for _, row := range activeRows {
		activeByDate[row.Date] = row.ActiveUsers
	}

	data := make([]dto.UserGrowthDailyItem, 0, lastDates)
	for i := 0; i < int(lastDates); i++ {
		day := fromTime.AddDate(0, 0, i)
		date := day.Format("2006-01-02")
		data = append(data, dto.UserGrowthDailyItem{
			Date:        date,
			NewUsers:    newByDate[date],
			ActiveUsers: activeByDate[date],
		})
	}
	return data
}

func buildUserGrowthHourlySeries(
	fromTime time.Time,
	newRows []dto.UserGrowthDailyItem,
	activeRows []dto.UserGrowthDailyItem,
) []dto.UserGrowthDailyItem {
	newByHour := make(map[string]uint64, len(newRows))
	for _, row := range newRows {
		newByHour[row.Date] = row.NewUsers
	}

	activeByHour := make(map[string]uint64, len(activeRows))
	for _, row := range activeRows {
		activeByHour[row.Date] = row.ActiveUsers
	}

	loc := fromTime.Location()
	now := time.Now().In(loc)
	endHour := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, loc)
	startHour := time.Date(fromTime.Year(), fromTime.Month(), fromTime.Day(), 0, 0, 0, 0, loc)

	hours := int(endHour.Sub(startHour).Hours()) + 1
	if hours < 1 {
		hours = 1
	}

	data := make([]dto.UserGrowthDailyItem, 0, hours)
	for cursor := startHour; !cursor.After(endHour); cursor = cursor.Add(time.Hour) {
		key := cursor.Format("2006-01-02 15:00")
		data = append(data, dto.UserGrowthDailyItem{
			Date:        key,
			NewUsers:    newByHour[key],
			ActiveUsers: activeByHour[key],
		})
	}
	return data
}

func providerRoleToRoleDTO(role *authdto.ProviderRoleDTO) *authdto.RoleDTO {
	if role == nil {
		return nil
	}
	return &authdto.RoleDTO{
		ID:              role.ID,
		RoleName:        role.RoleName,
		RoleKey:         role.RoleKey,
		RoleDescription: role.RoleDescription,
	}
}
