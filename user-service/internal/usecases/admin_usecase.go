package usecases

import (
	_errors "common/errors"
	_utils "common/utils"
	"context"
	"fmt"
	"user/internal/dto"
	"user/internal/interface/providers"
	"user/internal/interface/repo"
	"user/internal/models"
)

type AdminUsecase struct {
	adminRepo    repo.IAdminRepo
	authProvider providers.AuthProvider
}

func NewAdminUsecase(
	adminRepo repo.IAdminRepo,
	authProvider providers.AuthProvider,
) *AdminUsecase {
	return &AdminUsecase{
		adminRepo:    adminRepo,
		authProvider: authProvider,
	}
}

// BlockUser khóa tài khoản user với lý do
func (uc *AdminUsecase) BlockUser(ctx context.Context, req *dto.BlockUserRequest) error {

	// Kiểm tra user có tồn tại không
	_, err := uc.adminRepo.GetByID(ctx, req.ProfileID)
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
	_, err := uc.adminRepo.GetByID(ctx, req.ProfileID)
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
	_, err := uc.adminRepo.GetByID(ctx, req.ProfileID)
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

// CreateAdmin tạo tài khoản admin mới bằng cách gọi auth service
func (u *AdminUsecase) CreateAdmin(ctx context.Context, req *models.AdminProfile) (*models.AdminProfile, error) {
	// Kiểm tra quyền admin
	profileID := _utils.GetProfileIdWithContext(ctx)
	if profileID == 0 {
		return nil, _errors.ReturnError(int32(401), "Không có quyền truy cập")
	}

	// Convert response từ auth service
	createAdminResp, err := u.adminRepo.CreateAdmin(ctx, req)
	if err != nil {
		return nil, _errors.ReturnError(int32(500), "Lỗi convert response từ auth service")
	}

	// Gọi auth service để tạo admin
	_, err = u.authProvider.Update(ctx, &providers.AuthMethodDomain{
		UserID:   createAdminResp.ID,
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

	return createAdminResp, nil
}

// UpdateAdmin cập nhật thông tin admin bằng cách gọi auth service
func (u *AdminUsecase) UpdateAdmin(ctx context.Context, req *models.AdminProfile) (*models.AdminProfile, error) {
	// Kiểm tra quyền admin
	profileID := _utils.GetProfileIdWithContext(ctx)
	if profileID == 0 {
		return nil, _errors.ReturnError(int32(401), "Không có quyền truy cập")
	}

	// Gọi auth service để cập nhật admin
	_, err := u.authProvider.Update(ctx, &providers.AuthMethodDomain{
		UserID:   req.ID,
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
		return nil, _errors.ReturnError(int32(500), "Lỗi cập nhật thông tin admin")
	}

	return result, nil
}

// DeleteAdmin xóa tài khoản admin bằng cách gọi auth service
func (u *AdminUsecase) DeleteAdmin(ctx context.Context, req uint64) error {
	// Kiểm tra quyền admin
	profileID := _utils.GetProfileIdWithContext(ctx)
	if profileID == 0 {
		return _errors.ReturnError(int32(401), "Không có quyền truy cập")
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

// ListAdmins lấy danh sách admin bằng cách gọi auth service
func (u *AdminUsecase) ListAdmins(ctx context.Context, req *dto.AdminListRequest) ([]models.AdminProfile, int64, error) {
	// Kiểm tra quyền admin
	profileID := _utils.GetProfileIdWithContext(ctx)
	if profileID == 0 {
		return nil, 0, _errors.ReturnError(int32(401), "Không có quyền truy cập")
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
				detail.Role = role
			}
		}
	}

	fmt.Printf("GetDetail admin ID %d with role: %+v\n", req, detail.Role)

	return detail, nil
}
