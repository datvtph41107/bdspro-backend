package usecase

import (
	"auth/internal/domain"
	"auth/internal/dto"
	"auth/internal/enums"
	"auth/internal/interface/provider"
	"auth/internal/interface/repo"
	"auth/internal/utils"
	_errors "common/errors"
	_utils "common/utils"
	"context"
	"fmt"
	"time"
)

type AdminUsecase struct {
	AuthMethodRepo     repo.AuthMethodRepository
	ProfileProvider    provider.ProfileProvider
	OrganizationClient provider.OrganizationProvider
}

func NewAdminUsecase(
	authMethodRepo repo.AuthMethodRepository,
	profileProvider provider.ProfileProvider,
	organizationClient provider.OrganizationProvider,
) *AdminUsecase {
	return &AdminUsecase{
		AuthMethodRepo:     authMethodRepo,
		ProfileProvider:    profileProvider,
		OrganizationClient: organizationClient,
	}
}

// CreateAdmin tạo tài khoản admin mới
func (u *AdminUsecase) CreateAdmin(ctx context.Context, req dto.CreateAdminRequest) (*dto.AdminResponse, error) {
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
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, _errors.ReturnError(int32(500), "Lỗi mã hóa password")
	}

	// Tạo auth method mới
	authMethod := &domain.AuthMethodDomain{
		Provider:  "ADMIN",
		Username:  req.Username,
		Password:  hashedPassword,
		Name:      req.FullName,
		Email:     req.Email,
		Phone:     req.Phone,
		Avatar:    req.Avatar,
		CreatedAt: time.Now(),
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
func (u *AdminUsecase) UpdateAdmin(ctx context.Context, req dto.UpdateAdminRequest) (*dto.AdminResponse, error) {
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
func (u *AdminUsecase) DeleteAdmin(ctx context.Context, req dto.DeleteAdminRequest) error {
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
func (u *AdminUsecase) AssignRole(ctx context.Context, req dto.AssignRoleRequest) error {
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
func (u *AdminUsecase) ListAdmins(ctx context.Context, req dto.AdminListRequest) (*dto.AdminListResponse, error) {
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
