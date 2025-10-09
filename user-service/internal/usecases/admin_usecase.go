package usecases

import (
	_errors "common/errors"
	_utils "common/utils"
	"context"
	"fmt"
	"strings"
	"time"
	"user/enums"
	"user/internal/dto"
	"user/internal/interface/providers"
	"user/internal/interface/repo"
	"user/internal/models"
)

type AdminUsecase struct {
	adminRepo    repo.IAdminRepo
	authProvider providers.AuthProvider
}

func NewAdminUsecase(adminRepo repo.IAdminRepo, authProvider providers.AuthProvider) *AdminUsecase {
	return &AdminUsecase{
		adminRepo:    adminRepo,
		authProvider: authProvider,
	}
}

// ListAllUsers lấy danh sách tất cả user với phân trang và tìm kiếm
// func (uc *AdminUsecase) ListAllUsers(ctx context.Context, req *dto.AdminListRequest) (*dto.AdminListResponse, error) {
// 	// Kiểm tra quyền admin (sẽ được implement sau)
// 	// profileID := _utils.GetProfileIdWithContext(ctx)
// 	// if !uc.isAdmin(profileID) {
// 	//     return nil, errors.New("unauthorized: admin access required")
// 	// }

// 	// Lấy dữ liệu từ repository
// 	users, total, err := uc.adminRepo.ListAllUsers(ctx, req)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Tính toán thông tin phân trang
// 	totalPages := uint32(math.Ceil(float64(total) / float64(req.Size)))

// 	response := &dto.AdminListResponse{
// 		Data:       users,
// 		Total:      int32(total),
// 		Page:       req.Page,
// 		Size:       req.Size,
// 		TotalPages: totalPages,
// 	}

// 	return response, nil
// }

// BlockUser khóa tài khoản user với lý do
func (uc *AdminUsecase) BlockUser(ctx context.Context, req *dto.BlockUserRequest) error {

	// Kiểm tra user có tồn tại không
	_, err := uc.adminRepo.GetUserByProfileID(req.ProfileID)
	if err != nil {
		return fmt.Errorf("user not found: %v", err)
	}

	// Cập nhật trạng thái user thành blocked
	err = uc.adminRepo.UpdateUserStatus(req.ProfileID, 0) // 0 = blocked
	if err != nil {
		fmt.Printf("Failed to block user: %v\n", err)
		return fmt.Errorf("failed to block user: %v", err)
	}
	fmt.Printf("User blocked successfully\n")

	// TODO: Log action và gửi notification
	// TODO: Lưu lý do khóa vào bảng audit log

	return nil
}

// UnblockUser mở khóa tài khoản user với lý do
func (uc *AdminUsecase) UnblockUser(ctx context.Context, req *dto.UnblockUserRequest) error {
	// Kiểm tra user có tồn tại không
	_, err := uc.adminRepo.GetUserByProfileID(req.ProfileID)
	if err != nil {
		return fmt.Errorf("user not found: %v", err)
	}

	// Cập nhật trạng thái user thành active
	err = uc.adminRepo.UpdateUserStatus(req.ProfileID, 1) // 1 = active
	if err != nil {
		return fmt.Errorf("failed to unblock user: %v", err)
	}

	// TODO: Log action và gửi notification
	// TODO: Lưu lý do mở khóa vào bảng audit log

	return nil
}

// ApproveUser duyệt tài khoản user
func (uc *AdminUsecase) ApproveUser(ctx context.Context, req *dto.ApproveUserRequest) error {
	// Kiểm tra user có tồn tại không
	_, err := uc.adminRepo.GetUserByProfileID(req.ProfileID)
	if err != nil {
		return fmt.Errorf("user not found: %v", err)
	}

	// Cập nhật trạng thái user thành approved
	err = uc.adminRepo.UpdateUserStatus(req.ProfileID, 2) // 2 = approved
	if err != nil {
		return fmt.Errorf("failed to approve user: %v", err)
	}

	// TODO: Log action và gửi notification
	// TODO: Lưu ghi chú duyệt vào bảng audit log

	return nil
}

// CreateAdmin tạo tài khoản admin mới
func (u *AdminUsecase) CreateAdmin(ctx context.Context, req *models.AdminProfile) (*models.AdminProfile, error) {
	// Validate request
	// if err := req.Validate(); err != nil {
	// 	return nil, _errors.ReturnError(int32(400), err.Error())
	// }

	// Kiểm tra quyền admin
	// profileID := _utils.GetProfileIdWithContext(ctx)
	// if profileID == 0 {
	// 	return nil, _errors.ReturnError(int32(401), "Không có quyền truy cập")
	// }

	// Kiểm tra username đã tồn tại chưa
	existingAuth, err := u.authProvider.FindByUsernameAndProvider(ctx, req.Username, "ADMIN")
	if err == nil && existingAuth != nil {
		return nil, _errors.ReturnError(int32(400), "Username đã tồn tại")
	}

	// Kiểm tra email đã tồn tại chưa
	existingEmail, err := u.authProvider.FindByEmail(ctx, req.Email)
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
	authMethod := &providers.AuthMethodDomain{
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
	createdAuth, err := u.authProvider.Create(ctx, authMethod)
	if err != nil {
		fmt.Printf("Error creating auth method: %v\n", err)
		return nil, _errors.ReturnError(int32(500), "Lỗi tạo tài khoản admin")
	}
	fmt.Printf("Created auth method with ID: %d\n", createdAuth.ID)

	// TODO: Tạo profile cho admin
	// planID := uint64(1)
	// profile := u.profileProvider.MakeUserProfileEntity(
	// 	createdAuth.ID,
	// 	&planID,
	// 	req.FullName,
	// 	req.Email,
	// 	req.Phone,
	// 	req.Avatar,
	// )
	// profile.RoleType = enums.ERoleAdmin

	// // Lưu profile
	// err = u.profileProvider.SaveProfile(ctx, profile)
	// if err != nil {
	// 	return nil, _errors.ReturnError(int32(500), "Lỗi tạo profile admin")
	// }

	// TODO: Cập nhật UserID cho auth method
	// createdAuth.UserID = profile.ProfileID
	// _, err = u.authProvider.Update(ctx, createdAuth)
	// if err != nil {
	// 	return nil, _errors.ReturnError(int32(500), "Lỗi cập nhật thông tin admin")
	// }

	// Gửi thông báo tài khoản nếu được yêu cầu
	if req.SendNotification {
		// TODO: Implement gửi email/SMS thông báo tài khoản
		// Có thể gọi notification service để gửi email với thông tin tài khoản
		fmt.Printf("Sending notification to: %s, %s\n", req.Email, req.Phone)
		// u.NotificationProvider.SendAdminAccountNotification(ctx, req.Email, req.Phone, req.Username, req.Password)
	}

	return &models.AdminProfile{
		AuthID:    createdAuth.ID,
		ProfileID: 0, // TODO: Set from profile when available
		Username:  req.Username,
		FullName:  req.FullName,
		Email:     req.Email,
		Phone:     req.Phone,
		Avatar:    req.Avatar,
		RoleKey:   enums.RoleMap[enums.ERoleAdmin],
	}, nil
}

// UpdateAdmin cập nhật thông tin admin
func (u *AdminUsecase) UpdateAdmin(ctx context.Context, req dto.UpdateAdminRequest) (*models.AdminProfile, error) {
	// Kiểm tra quyền admin
	profileID := _utils.GetProfileIdWithContext(ctx)
	if profileID == 0 {
		return nil, _errors.ReturnError(int32(401), "Không có quyền truy cập")
	}

	// Tìm admin cần cập nhật
	authMethod, err := u.authProvider.FindByID(ctx, req.AuthID)
	if err != nil || authMethod == nil {
		return nil, _errors.ReturnError(int32(404), "Không tìm thấy tài khoản admin")
	}

	// Kiểm tra provider phải là ADMIN
	if authMethod.Provider != "ADMIN" {
		return nil, _errors.ReturnError(int32(400), "Tài khoản không phải admin")
	}

	// Kiểm tra email đã tồn tại chưa (nếu thay đổi)
	if req.Email != authMethod.Email {
		existingEmail, err := u.authProvider.FindByEmail(ctx, req.Email)
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
	updatedAuth, err := u.authProvider.Update(ctx, authMethod)
	if err != nil {
		return nil, _errors.ReturnError(int32(500), "Lỗi cập nhật thông tin admin")
	}

	// Cập nhật profile
	profile, err := u.adminRepo.GetUserByProfileID(authMethod.UserID)
	if err == nil && profile != nil {
		profile.FullName = req.FullName
		profile.Email = req.Email
		profile.Phone = req.Phone
		profile.Avatar = req.Avatar
		u.adminRepo.UpdateAdmin(ctx, profile)
	}

	return &models.AdminProfile{
		AuthID:    updatedAuth.ID,
		ProfileID: updatedAuth.UserID,
		Username:  updatedAuth.Username,
		FullName:  req.FullName,
		Email:     req.Email,
		Phone:     req.Phone,
		Avatar:    req.Avatar,
		RoleKey:   enums.RoleMap[enums.ERoleAdmin],
	}, nil
}

// DeleteAdmin xóa tài khoản admin
func (u *AdminUsecase) DeleteAdmin(ctx context.Context, req dto.DeleteAdminRequest) error {
	// // Kiểm tra quyền admin
	// profileID := _utils.GetProfileIdWithContext(ctx)
	// if profileID == 0 {
	// 	return _errors.ReturnError(int32(401), "Không có quyền truy cập")
	// }

	// Tìm admin cần xóa
	authMethod, err := u.authProvider.FindByID(ctx, req.AuthID)
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
	err = u.authProvider.SoftDelete(ctx, req.AuthID)
	if err != nil {
		return _errors.ReturnError(int32(500), "Lỗi xóa tài khoản admin")
	}

	// Soft delete profile
	if authMethod.UserID > 0 {
		u.adminRepo.DeleteAdmin(ctx, authMethod.UserID)
	}

	return nil
}

// // AssignRole gán quyền cho admin
// func (u *AdminUsecase) AssignRole(ctx context.Context, req dto.AssignRoleRequest) error {
// 	// Kiểm tra quyền admin
// 	profileID := _utils.GetProfileIdWithContext(ctx)
// 	if profileID == 0 {
// 		return _errors.ReturnError(int32(401), "Không có quyền truy cập")
// 	}

// 	// Tìm admin cần gán quyền
// 	authMethod, err := u.authProvider.FindByID(ctx, req.AuthID)
// 	if err != nil || authMethod == nil {
// 		return _errors.ReturnError(int32(404), "Không tìm thấy tài khoản admin")
// 	}

// 	// Kiểm tra provider phải là ADMIN
// 	if authMethod.Provider != "ADMIN" {
// 		return _errors.ReturnError(int32(400), "Tài khoản không phải admin")
// 	}

// 	// Lấy danh sách role admin từ organization service
// 	adminRoles, err := u.organizationProvider.GetAdminRoles(ctx)
// 	if err != nil {
// 		return _errors.ReturnError(int32(500), "Lỗi lấy danh sách role")
// 	}

// 	// Kiểm tra role key có hợp lệ không
// 	validRole := false
// 	for _, role := range adminRoles {
// 		if role.Key == req.RoleKey {
// 			validRole = true
// 			break
// 		}
// 	}

// 	if !validRole {
// 		return _errors.ReturnError(int32(400), "Role key không hợp lệ")
// 	}

// 	// Cập nhật role trong profile
// 	profile, err := u.adminRepo.GetUserByProfileID(ctx, authMethod.UserID)
// 	if err == nil && profile != nil {
// 		// Có thể cần thêm logic để map role key với RoleType
// 		// Hiện tại chỉ set mặc định là ERoleAdmin
// 		profile.RoleType = enums.ERoleAdmin
// 		u.adminRepo.UpdateAdmin(ctx, profile)
// 	}

// 	return nil
// }

// ListAdmins lấy danh sách admin
func (u *AdminUsecase) ListAdmins(ctx context.Context, req *dto.AdminListRequest) ([]models.AdminProfile, int64, error) {
	// Kiểm tra quyền admin
	// profileID := _utils.GetProfileIdWithContext(ctx)
	// if profileID == 0 {
	// 	return nil, _errors.ReturnError(int32(401), "Không có quyền truy cập")
	// }

	// Lấy danh sách admin từ database
	admins, total, err := u.adminRepo.ListAllUsers(ctx, req)
	if err != nil {
		return nil, 0, _errors.ReturnError(500, "Lỗi lấy danh sách admin")
	}

	// Lấy danh sách role admin từ organization service
	// adminRoles, err := u.OrganizationClient.GetAdminRoles(ctx)
	// if err != nil {
	// 	return nil, _errors.ReturnError(int32(500), "Lỗi lấy danh sách role")
	// }

	return admins, total, nil
}
