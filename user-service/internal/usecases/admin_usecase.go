package usecases

import (
	_enum "common/domain/enum"
	_errors "common/errors"
	"context"
	"fmt"
	"user/internal"
	"user/internal/dto"
	"user/internal/interface/providers"
	"user/internal/interface/repo"
	"user/internal/models"
	"user/internal/usecase/useradmin"
)

type AdminUsecase struct {
	adminRepo    repo.IAdminRepo
	profileRepo  repo.IProfileRepo
	authProvider providers.AuthProvider
}

func adminActorID(ctx context.Context) (uint64, error) {
	actorID, err := useradmin.ActorIDFromContext(ctx)
	if err != nil {
		return 0, _errors.ReturnError(service.AccessDenied)
	}
	return actorID, nil
}

func NewAdminUsecase(
	adminRepo repo.IAdminRepo,
	profileRepo repo.IProfileRepo,
	authProvider providers.AuthProvider,
) *AdminUsecase {
	return &AdminUsecase{
		adminRepo:    adminRepo,
		profileRepo:  profileRepo,
		authProvider: authProvider,
	}
}

// BlockUser khóa tài khoản user với lý do
func (uc *AdminUsecase) LockUser(ctx context.Context, req *dto.LockUserRequest) error {
	if err := uc.rejectSystemRootMutation(ctx, req.ProfileID); err != nil {
		return err
	}

	// Kiểm tra user có tồn tại không
	_, err := uc.profileRepo.GetUserDetailByID(ctx, req.ProfileID)
	if err != nil {
		return fmt.Errorf("user not found: %v", err)
	}

	// Cập nhật trạng thái user_profile thành khóa tạm thời
	err = uc.profileRepo.UpdateUserProfileStatus(ctx, req.ProfileID, req.LockType)
	if err != nil {
		fmt.Printf("Failed to block user in user_profile: %v\n", err)
		return fmt.Errorf("failed to update user profile status: %v", err)
	}

	fmt.Printf("User blocked successfully (Temporary Locked)\n")

	// TODO: Log action và gửi notification
	// TODO: Lưu lý do khóa vào bảng audit log

	return nil
}

// UnblockUser mở khóa tài khoản user với lý do
func (uc *AdminUsecase) UnLockUser(ctx context.Context, req *dto.UnLockUserRequest) error {
	if err := uc.rejectSystemRootMutation(ctx, req.ProfileID); err != nil {
		return err
	}
	// Kiểm tra user có tồn tại không
	_, err := uc.profileRepo.GetUserDetailByID(ctx, req.ProfileID)
	if err != nil {
		return fmt.Errorf("user not found: %v", err)
	}

	// Cập nhật trạng thái user_profile thành active
	err = uc.profileRepo.UpdateUserProfileStatus(ctx, req.ProfileID, _enum.EUserStatusActive)
	if err != nil {
		return fmt.Errorf("failed to update user profile status: %v", err)
	}

	fmt.Printf("User unblocked successfully (Active)\n")

	// TODO: Log action và gửi notification
	// TODO: Lưu lý do mở khóa vào bảng audit log

	return nil
}

// ApproveUser duyệt tài khoản user
func (uc *AdminUsecase) ApproveUser(ctx context.Context, req *dto.ApproveUserRequest) error {
	if err := uc.rejectSystemRootMutation(ctx, req.ProfileID); err != nil {
		return err
	}
	// Kiểm tra user có tồn tại không
	_, err := uc.profileRepo.GetUserDetailByID(ctx, req.ProfileID)
	if err != nil {
		return fmt.Errorf("user not found: %v", err)
	}

	// Cập nhật trạng thái user_profile thành active
	err = uc.profileRepo.UpdateUserProfileStatus(ctx, req.ProfileID, _enum.EUserStatusActive)
	if err != nil {
		return fmt.Errorf("failed to update user profile status: %v", err)
	}

	fmt.Printf("User approved successfully (Active)\n")

	// TODO: Log action và gửi notification
	// TODO: Lưu ghi chú duyệt vào bảng audit log

	return nil
}

// CreateAdmin tạo tài khoản admin mới bằng cách gọi auth service
func (u *AdminUsecase) CreateAdmin(ctx context.Context, req *models.AdminProfile) (*models.AdminProfile, error) {
	// Kiểm tra quyền admin
	if _, err := adminActorID(ctx); err != nil {
		return nil, err
	}

	if req.Username == "" {
		return nil, _errors.ReturnError(service.UsernameRequired)
	}
	if req.Password == "" {
		return nil, _errors.ReturnError(service.PasswordRequired)
	}

	// Convert response từ auth service
	createAdminResp, err := u.adminRepo.CreateAdmin(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("create admin profile record: %w", err)
	}

	// Link auth.user_id ↔ admin_profiles.id (BaseEntity.ID). Fallback AuthID nếu ID chưa được GORM fill.
	linkUserID := createAdminResp.ID
	if linkUserID == 0 {
		linkUserID = createAdminResp.AuthID
	}
	if linkUserID == 0 {
		return nil, fmt.Errorf("admin profile id missing after create")
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
		return nil, fmt.Errorf("create admin login info: %w", err)
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
			return nil, fmt.Errorf("assign role after admin creation: %w", err)
		}
	}

	return createAdminResp, nil
}

// UpdateAdmin cập nhật thông tin admin bằng cách gọi auth service
func (u *AdminUsecase) UpdateAdmin(ctx context.Context, req *models.AdminProfile) (*models.AdminProfile, error) {
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
		return nil, fmt.Errorf("update admin auth service: %w", err)
	}

	result, err := u.adminRepo.UpdateAdmin(ctx, req)
	if err != nil {
		fmt.Printf("Error updating admin profile in DB: %v\n", err)
		return nil, fmt.Errorf("update admin profile: %w", err)
	}

	if req.RoleID != nil && *req.RoleID > 0 && linkUserID > 0 {
		if err := u.authProvider.AssignRoleToUser(ctx, linkUserID, *req.RoleID); err != nil {
			fmt.Printf("Error assigning role via auth service: %v\n", err)
			return nil, fmt.Errorf("assign role after admin update: %w", err)
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
		return fmt.Errorf("delete admin account: %w", err)
	}

	err = u.adminRepo.DeleteAdmin(ctx, req)
	if err != nil {
		return fmt.Errorf("delete admin account: %w", err)
	}

	return nil
}

func (u *AdminUsecase) rejectSystemRootAdminMutation(ctx context.Context, id uint64) error {
	if u == nil || u.adminRepo == nil {
		return _errors.ReturnError(service.AdminDataStoreUnavailable)
	}
	isRoot, err := u.adminRepo.IsSystemRootAdmin(ctx, id)
	if err != nil {
		return fmt.Errorf("verify system root admin: %w", err)
	}
	if isRoot {
		return _errors.ReturnError(service.RootAdminMutationDenied)
	}
	return nil
}

// ListAdmins lấy danh sách admin bằng cách gọi auth service
func (u *AdminUsecase) ListAdmins(ctx context.Context, req *dto.AdminListRequest) ([]models.AdminProfile, int64, error) {
	// Kiểm tra quyền admin
	if _, err := adminActorID(ctx); err != nil {
		return nil, 0, err
	}

	// Lấy danh sách admin từ repo
	authResp, total, err := u.adminRepo.ListAdmin(ctx, req)
	if err != nil {
		fmt.Printf("Error listing admins via auth service: %v\n", err)
		return nil, 0, fmt.Errorf("list admins: %w", err)
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
	var roleMap map[uint64]*models.RoleDTO
	if len(roleIds) > 0 {
		roleMap, err = u.authProvider.GetRolesByIds(ctx, roleIds)
		if err != nil {
			fmt.Printf("Error getting roles by IDs: %v\n", err)
			// Không return error, chỉ log warning vì role là thông tin bổ sung
			roleMap = make(map[uint64]*models.RoleDTO)
		}
	}

	// Map thông tin role vào từng admin
	for i := range authResp {
		if authResp[i].RoleID != nil && *authResp[i].RoleID > 0 {
			if role, ok := roleMap[*authResp[i].RoleID]; ok {
				authResp[i].Role = role
			}
		}
	}

	fmt.Printf("Listed %d admins with %d roles\n", len(authResp), len(roleMap))

	return authResp, total, nil
}

func (u *AdminUsecase) GetDetail(ctx context.Context, req uint64) (*models.AdminProfile, error) {
	// Kiểm tra quyền admin
	authResp, _ := u.authProvider.FindByID(ctx, req)

	detail, err := u.adminRepo.GetDetail(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("load admin information: %w", err)
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
				detail.Role = role
			}
		}
	}

	fmt.Printf("GetDetail admin ID %d with role: %+v\n", req, detail.Role)

	return detail, nil
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
		return nil, _errors.ReturnError(service.EmailAlreadyExistsInSystem)
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
		return nil, _errors.ReturnError(service.UsernameAlreadyExistsInSystem)
	}

	password := req.Password
	if password == "" {
		return nil, _errors.ReturnError(service.PasswordRequiredVI)
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
		return nil, fmt.Errorf("create auth account: %w", err)
	}

	// Tạo user profile trong user service (dùng UserProfileEntity)
	userProfile := &models.UserProfileEntity{
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
	models.SetAllVisibilityPublic(userProfile)

	createdProfile, err := u.profileRepo.CreateUserProfile(ctx, userProfile)
	if err != nil {
		fmt.Printf("Error creating user profile: %v\n", err)
		// Rollback: xóa auth method đã tạo
		_ = u.authProvider.SoftDelete(ctx, createdAuth.ID)
		return nil, fmt.Errorf("create user profile: %w", err)
	}

	// Link the credential to the canonical profile generated by PostgreSQL.
	createdAuth.UserID = createdProfile.ProfileID
	createdAuth, err = u.authProvider.Update(ctx, createdAuth)
	if err != nil {
		_ = u.authProvider.SoftDelete(ctx, createdAuth.ID)
		return nil, fmt.Errorf("link auth account to profile: %w", err)
	}
	if req.RoleID > 0 {
		if err := u.authProvider.AssignRoleToUser(ctx, createdProfile.ProfileID, req.RoleID); err != nil {
			return nil, fmt.Errorf("assign role to created user: %w", err)
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
		return nil, _errors.ReturnError(service.UserNotFound)
	}

	// Cập nhật thông tin trong user profile
	updatedProfile := &models.UserProfileEntity{
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
		return nil, fmt.Errorf("update user profile: %w", err)
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
		return _errors.ReturnError(service.UserNotFound)
	}

	// Xóa user profile (soft delete)
	// Lấy thông tin admin đang thực hiện xóa
	reason := "Xóa bởi admin"
	if req.Reason != "" {
		reason = req.Reason
	}
	err = u.profileRepo.DeleteUserProfile(ctx, req.ProfileID, &adminProfileID, reason)
	if err != nil {
		return fmt.Errorf("delete user profile: %w", err)
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
		return fmt.Errorf("user data store is unavailable")
	}
	isRoot, err := u.profileRepo.IsSystemRootProfile(ctx, profileID)
	if err != nil {
		return fmt.Errorf("verify system root admin: %w", err)
	}
	if isRoot {
		return _errors.ReturnError(service.RootAdminMutationDenied)
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
		return nil, _errors.ReturnError(service.UserNotFound)
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
