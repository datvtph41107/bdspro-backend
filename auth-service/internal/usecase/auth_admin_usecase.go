package usecase

import (
	"auth/internal/domain"
	"auth/internal/dto"
	"auth/internal/enums"
	"auth/internal/interface/provider"
	"auth/internal/interface/repo"
	_errors "common/errors"
	_utils "common/utils"
	"context"
	"fmt"
	"strings"
	"time"
)

// AuthAdminUsecase - Quản lý admin trong auth-service
//
// TRÁCH NHIỆM:
// - Tạo, cập nhật, xóa tài khoản admin
// - Quản lý mật khẩu admin (change password, reset password)
// - Gán quyền cho admin
// - Lấy danh sách admin
// - Tất cả logic liên quan đến authentication và authorization của admin

type AuthAdminUsecase struct {
	AuthMethodRepo     repo.AuthMethodRepository
	ProfileProvider    provider.ProfileProvider
	OrganizationClient provider.OrganizationProvider
}

func NewAdminUsecase(
	authMethodRepo repo.AuthMethodRepository,
	profileProvider provider.ProfileProvider,
	organizationClient provider.OrganizationProvider,
) *AuthAdminUsecase {
	return &AuthAdminUsecase{
		AuthMethodRepo:     authMethodRepo,
		ProfileProvider:    profileProvider,
		OrganizationClient: organizationClient,
	}
}

// CreateAdmin tạo tài khoản admin mới
func (u *AuthAdminUsecase) CreateAdmin(ctx context.Context, req dto.CreateAdminRequest) (*dto.AdminResponse, error) {
	// Validate request
	if err := req.Validate(); err != nil {
		return nil, _errors.ReturnError(int32(400), err.Error())
	}

	// Kiểm tra quyền admin
	// profileID := _utils.GetProfileIdWithContext(ctx)
	// if profileID == 0 {
	// 	return nil, _errors.ReturnError(int32(401), "Không có quyền truy cập")
	// }

	// Kiểm tra username đã tồn tại chưa
	existingAuth, err := u.AuthMethodRepo.FindByUsernameAndProvider(ctx, req.Username, "ADMIN")
	if err == nil && existingAuth != nil {
		return nil, _errors.ReturnError(int32(400), "Username đã tồn tại")
	}

	// Kiểm tra email đã tồn tại chưa
	existingEmail, err := u.AuthMethodRepo.FindByEmail(ctx, req.Email)
	if err == nil && existingEmail != nil {
		return nil, _errors.ReturnError(int32(400), "Email đã tồn tại")
	}

	// Mã hóa password
	hashedPassword, err := _utils.HashPassword(req.Password)
	if err != nil {
		return nil, _errors.ReturnError(int32(500), "Lỗi mã hóa password")
	}

	// Convert roleKey string to uint32 (có thể cần mapping logic)
	var roleKey uint32 = 1 // Default role, cần implement logic mapping

	// Convert attachments array to JSON string
	attachmentsJSON := "[]"
	if len(req.Attachments) > 0 {
		attachmentsJSON = fmt.Sprintf(`["%s"]`, strings.Join(req.Attachments, `","`))
	}

	// Tạo auth method mới
	authMethod := &domain.AuthMethodDomain{
		Provider:         "ADMIN",
		Username:         req.Username,
		Password:         hashedPassword,
		Name:             req.FullName,
		Email:            req.Email,
		Phone:            req.Phone,
		Avatar:           req.Avatar,
		JobTitle:         req.JobTitle,
		RoleKey:          roleKey,
		Attachments:      attachmentsJSON,
		InternalNotes:    req.InternalNotes,
		SendNotification: req.SendNotification,
		CreatedAt:        time.Now(),
	}

	// Lưu auth method
	createdAuth, err := u.AuthMethodRepo.Create(ctx, authMethod)
	if err != nil {
		fmt.Printf("Error creating auth method: %v\n", err)
		return nil, _errors.ReturnError(int32(500), "Lỗi tạo tài khoản admin")
	}
	fmt.Printf("Created auth method with ID: %d\n", createdAuth.ID)

	// Tạo profile cho admin
	planID := uint64(1)
	profile := u.ProfileProvider.MakeUserProfileEntity(
		createdAuth.ID,
		&planID,
		req.FullName,
		req.Email,
		req.Phone,
		req.Avatar,
	)
	profile.RoleType = enums.ERoleAdmin

	// Lưu profile
	err = u.ProfileProvider.SaveProfile(ctx, profile)
	if err != nil {
		return nil, _errors.ReturnError(int32(500), "Lỗi tạo profile admin")
	}

	// Cập nhật UserID cho auth method
	createdAuth.UserID = profile.ProfileID
	_, err = u.AuthMethodRepo.Update(ctx, createdAuth)
	if err != nil {
		return nil, _errors.ReturnError(int32(500), "Lỗi cập nhật thông tin admin")
	}

	// Gửi thông báo tài khoản nếu được yêu cầu
	if req.SendNotification {
		// TODO: Implement gửi email/SMS thông báo tài khoản
		// Có thể gọi notification service để gửi email với thông tin tài khoản
		fmt.Printf("Sending notification to: %s, %s\n", req.Email, req.Phone)
		// u.NotificationProvider.SendAdminAccountNotification(ctx, req.Email, req.Phone, req.Username, req.Password)
	}

	return &dto.AdminResponse{
		AuthID:    createdAuth.ID,
		ProfileID: profile.ProfileID,
		Username:  req.Username,
		FullName:  req.FullName,
		Email:     req.Email,
		Phone:     req.Phone,
		Avatar:    req.Avatar,
		Role:      enums.RoleMap[enums.ERoleAdmin],
	}, nil
}

// UpdateAdmin cập nhật thông tin admin
func (u *AuthAdminUsecase) UpdateAdmin(ctx context.Context, req dto.UpdateAdminRequest) (*dto.AdminResponse, error) {
	// Kiểm tra quyền admin
	profileID := _utils.GetProfileIdWithContext(ctx)
	if profileID == 0 {
		return nil, _errors.ReturnError(int32(401), "Không có quyền truy cập")
	}

	// Tìm admin cần cập nhật
	authMethod, err := u.AuthMethodRepo.FindByID(ctx, req.AuthID)
	if err != nil || authMethod == nil {
		return nil, _errors.ReturnError(int32(404), "Không tìm thấy tài khoản admin")
	}

	// Kiểm tra provider phải là ADMIN
	if authMethod.Provider != "ADMIN" {
		return nil, _errors.ReturnError(int32(400), "Tài khoản không phải admin")
	}

	// Kiểm tra email đã tồn tại chưa (nếu thay đổi)
	if req.Email != authMethod.Email {
		existingEmail, err := u.AuthMethodRepo.FindByEmail(ctx, req.Email)
		if err == nil && existingEmail != nil && existingEmail.ID != req.AuthID {
			return nil, _errors.ReturnError(int32(400), "Email đã tồn tại")
		}
	}

	// Cập nhật thông tin auth method
	authMethod.Name = req.FullName
	authMethod.Email = req.Email
	authMethod.Phone = req.Phone
	authMethod.Avatar = req.Avatar

	// Lưu cập nhật
	updatedAuth, err := u.AuthMethodRepo.Update(ctx, authMethod)
	if err != nil {
		return nil, _errors.ReturnError(int32(500), "Lỗi cập nhật thông tin admin")
	}

	// Cập nhật profile
	profile, err := u.ProfileProvider.GetByProfileID(ctx, authMethod.UserID)
	if err == nil && profile != nil {
		profile.FullName = req.FullName
		profile.Email = req.Email
		profile.Phone = req.Phone
		profile.Avatar = req.Avatar
		u.ProfileProvider.UpdateProfile(ctx, profile)
	}

	return &dto.AdminResponse{
		AuthID:    updatedAuth.ID,
		ProfileID: updatedAuth.UserID,
		Username:  updatedAuth.Username,
		FullName:  req.FullName,
		Email:     req.Email,
		Phone:     req.Phone,
		Avatar:    req.Avatar,
		Role:      enums.RoleMap[enums.ERoleAdmin],
	}, nil
}

// DeleteAdmin xóa tài khoản admin
func (u *AuthAdminUsecase) DeleteAdmin(ctx context.Context, req dto.DeleteAdminRequest) error {
	// // Kiểm tra quyền admin
	// profileID := _utils.GetProfileIdWithContext(ctx)
	// if profileID == 0 {
	// 	return _errors.ReturnError(int32(401), "Không có quyền truy cập")
	// }

	// Tìm admin cần xóa
	authMethod, err := u.AuthMethodRepo.FindByID(ctx, req.AuthID)
	if err != nil || authMethod == nil {
		return _errors.ReturnError(int32(404), "Không tìm thấy tài khoản admin")
	}

	// Kiểm tra provider phải là ADMIN
	if authMethod.Provider != "ADMIN" {
		return _errors.ReturnError(int32(400), "Tài khoản không phải admin")
	}

	// // Không cho phép xóa chính mình
	// if authMethod.UserID == profileID {
	// 	return _errors.ReturnError(int32(400), "Không thể xóa tài khoản của chính mình")
	// }

	// Soft delete auth method
	err = u.AuthMethodRepo.SoftDelete(ctx, req.AuthID)
	if err != nil {
		return _errors.ReturnError(int32(500), "Lỗi xóa tài khoản admin")
	}

	// Soft delete profile
	if authMethod.UserID > 0 {
		u.ProfileProvider.SoftDeleteProfile(ctx, authMethod.UserID)
	}

	return nil
}

// AssignRole gán quyền cho admin
func (u *AuthAdminUsecase) AssignRole(ctx context.Context, req dto.AssignRoleRequest) error {
	// Kiểm tra quyền admin
	profileID := _utils.GetProfileIdWithContext(ctx)
	if profileID == 0 {
		return _errors.ReturnError(int32(401), "Không có quyền truy cập")
	}

	// Tìm admin cần gán quyền
	authMethod, err := u.AuthMethodRepo.FindByID(ctx, req.AuthID)
	if err != nil || authMethod == nil {
		return _errors.ReturnError(int32(404), "Không tìm thấy tài khoản admin")
	}

	// Kiểm tra provider phải là ADMIN
	if authMethod.Provider != "ADMIN" {
		return _errors.ReturnError(int32(400), "Tài khoản không phải admin")
	}

	// Lấy danh sách role admin từ organization service
	adminRoles, err := u.OrganizationClient.GetAdminRoles(ctx)
	if err != nil {
		return _errors.ReturnError(int32(500), "Lỗi lấy danh sách role")
	}

	// Kiểm tra role key có hợp lệ không
	validRole := false
	for _, role := range adminRoles {
		if role.Key == req.RoleKey {
			validRole = true
			break
		}
	}

	if !validRole {
		return _errors.ReturnError(int32(400), "Role key không hợp lệ")
	}

	// Cập nhật role trong profile
	profile, err := u.ProfileProvider.GetByProfileID(ctx, authMethod.UserID)
	if err == nil && profile != nil {
		// Có thể cần thêm logic để map role key với RoleType
		// Hiện tại chỉ set mặc định là ERoleAdmin
		profile.RoleType = enums.ERoleAdmin
		u.ProfileProvider.UpdateProfile(ctx, profile)
	}

	return nil
}

// ListAdmins lấy danh sách admin
func (u *AuthAdminUsecase) ListAdmins(ctx context.Context, req dto.AdminListRequest) (*dto.AdminListResponse, error) {
	// Kiểm tra quyền admin
	// profileID := _utils.GetProfileIdWithContext(ctx)
	// if profileID == 0 {
	// 	return nil, _errors.ReturnError(int32(401), "Không có quyền truy cập")
	// }

	// Lấy danh sách admin từ database
	admins, total, err := u.AuthMethodRepo.FindAdminsByProvider(ctx, "ADMIN", int(req.Page), int(req.Size))
	if err != nil {
		return nil, _errors.ReturnError(int32(500), "Lỗi lấy danh sách admin")
	}

	// Lấy danh sách role admin từ organization service
	// adminRoles, err := u.OrganizationClient.GetAdminRoles(ctx)
	// if err != nil {
	// 	return nil, _errors.ReturnError(int32(500), "Lỗi lấy danh sách role")
	// }

	// Tạo map role để lookup nhanh
	// roleMap := make(map[string]string)
	// for _, role := range adminRoles {
	// 	roleMap[role.Key] = role.Name
	// }

	// Chuyển đổi sang response
	var adminItems []dto.AdminItem
	for _, admin := range admins {
		// Lấy thông tin profile
		var profile *dto.ProfileDTO
		if admin.UserID > 0 {
			profile, _ = u.ProfileProvider.GetByProfileID(ctx, admin.UserID)
		}

		adminItem := dto.AdminItem{
			AuthID:    admin.ID,
			ProfileID: admin.UserID,
			Username:  admin.Username,
			FullName:  admin.Name,
			Email:     admin.Email,
			Phone:     admin.Phone,
			Avatar:    admin.Avatar,
			Role:      enums.RoleMap[enums.ERoleAdmin],
			RoleType:  enums.ERoleAdmin,
			CreatedAt: admin.CreatedAt,
		}

		// Cập nhật thông tin từ profile nếu có
		if profile != nil {
			adminItem.FullName = profile.FullName
			adminItem.Email = profile.Email
			adminItem.Phone = profile.Phone
			adminItem.Avatar = profile.Avatar
			adminItem.RoleType = profile.RoleType
			adminItem.Role = enums.RoleMap[profile.RoleType]
		}

		adminItems = append(adminItems, adminItem)
	}

	return &dto.AdminListResponse{
		Data:  adminItems,
		Total: int32(total),
	}, nil
}

// GetAdminByID lấy thông tin admin theo ID
func (u *AuthAdminUsecase) GetAdminByID(ctx context.Context, authID uint64) (*dto.AdminResponse, error) {
	// Tìm admin theo ID
	authMethod, err := u.AuthMethodRepo.FindByID(ctx, authID)
	if err != nil || authMethod == nil {
		return nil, _errors.ReturnError(int32(404), "Không tìm thấy tài khoản admin")
	}

	// Kiểm tra provider phải là ADMIN
	if authMethod.Provider != "ADMIN" {
		return nil, _errors.ReturnError(int32(400), "Tài khoản không phải admin")
	}

	// Lấy thông tin profile
	var profile *dto.ProfileDTO
	if authMethod.UserID > 0 {
		profile, _ = u.ProfileProvider.GetByProfileID(ctx, authMethod.UserID)
	}

	response := &dto.AdminResponse{
		AuthID:    authMethod.ID,
		ProfileID: authMethod.UserID,
		Username:  authMethod.Username,
		FullName:  authMethod.Name,
		Email:     authMethod.Email,
		Phone:     authMethod.Phone,
		Avatar:    authMethod.Avatar,
		Role:      enums.RoleMap[enums.ERoleAdmin],
	}

	// Cập nhật thông tin từ profile nếu có
	if profile != nil {
		response.FullName = profile.FullName
		response.Email = profile.Email
		response.Phone = profile.Phone
		response.Avatar = profile.Avatar
		response.Role = enums.RoleMap[profile.RoleType]
	}

	return response, nil
}

// ChangePassword thay đổi mật khẩu admin
func (u *AuthAdminUsecase) ChangePassword(ctx context.Context, req dto.ChangePasswordRequest) error {
	// Kiểm tra quyền admin
	profileID := _utils.GetProfileIdWithContext(ctx)
	if profileID == 0 {
		return _errors.ReturnError(int32(401), "Không có quyền truy cập")
	}

	// Tìm admin theo ID
	authMethod, err := u.AuthMethodRepo.FindByID(ctx, req.AuthID)
	if err != nil || authMethod == nil {
		return _errors.ReturnError(int32(404), "Không tìm thấy tài khoản admin")
	}

	// Kiểm tra provider phải là ADMIN
	if authMethod.Provider != "ADMIN" {
		return _errors.ReturnError(int32(400), "Tài khoản không phải admin")
	}

	// Kiểm tra mật khẩu cũ
	if !_utils.CheckPasswordHash(req.OldPassword, authMethod.Password) {
		return _errors.ReturnError(int32(400), "Mật khẩu cũ không đúng")
	}

	// Mã hóa mật khẩu mới
	hashedPassword, err := _utils.HashPassword(req.NewPassword)
	if err != nil {
		return _errors.ReturnError(int32(500), "Lỗi mã hóa mật khẩu")
	}

	// Cập nhật mật khẩu
	authMethod.Password = hashedPassword
	_, err = u.AuthMethodRepo.Update(ctx, authMethod)
	if err != nil {
		return _errors.ReturnError(int32(500), "Lỗi cập nhật mật khẩu")
	}

	return nil
}

// ResetPassword reset mật khẩu admin
func (u *AuthAdminUsecase) ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error {
	// Kiểm tra quyền admin
	profileID := _utils.GetProfileIdWithContext(ctx)
	if profileID == 0 {
		return _errors.ReturnError(int32(401), "Không có quyền truy cập")
	}

	// Tìm admin theo ID
	authMethod, err := u.AuthMethodRepo.FindByID(ctx, req.AuthID)
	if err != nil || authMethod == nil {
		return _errors.ReturnError(int32(404), "Không tìm thấy tài khoản admin")
	}

	// Kiểm tra provider phải là ADMIN
	if authMethod.Provider != "ADMIN" {
		return _errors.ReturnError(int32(400), "Tài khoản không phải admin")
	}

	// Mã hóa mật khẩu mới
	hashedPassword, err := _utils.HashPassword(req.NewPassword)
	if err != nil {
		return _errors.ReturnError(int32(500), "Lỗi mã hóa mật khẩu")
	}

	// Cập nhật mật khẩu
	authMethod.Password = hashedPassword
	_, err = u.AuthMethodRepo.Update(ctx, authMethod)
	if err != nil {
		return _errors.ReturnError(int32(500), "Lỗi reset mật khẩu")
	}

	return nil
}
