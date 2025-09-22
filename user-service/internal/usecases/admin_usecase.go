package usecases

import (
	"context"
	"fmt"
	"math"
	user_repo "user/infra/repo"
	"user/internal/dto"
	provider "user/internal/interface/provider"
)

type AdminUsecase struct {
	AdminRepo  *user_repo.AdminRepo
	AuthClient provider.AuthProvider
}

func NewAdminUsecase(adminRepo *user_repo.AdminRepo, authClient provider.AuthProvider) *AdminUsecase {
	return &AdminUsecase{
		AdminRepo:  adminRepo,
		AuthClient: authClient,
	}
}

// ListAllUsers lấy danh sách tất cả user với phân trang và tìm kiếm
func (uc *AdminUsecase) ListAllUsers(ctx context.Context, req *dto.ListAllUsersRequest) (*dto.ListAllUsersResponse, error) {
	// Kiểm tra quyền admin (sẽ được implement sau)
	// profileID := _utils.GetProfileIdWithContext(ctx)
	// if !uc.isAdmin(profileID) {
	//     return nil, errors.New("unauthorized: admin access required")
	// }

	// Lấy dữ liệu từ repository
	users, total, err := uc.AdminRepo.ListAllUsers(ctx, req)
	if err != nil {
		return nil, err
	}

	// Tính toán thông tin phân trang
	totalPages := uint32(math.Ceil(float64(total) / float64(req.Size)))

	response := &dto.ListAllUsersResponse{
		Data:       users,
		Total:      int32(total),
		Page:       req.Page,
		Size:       req.Size,
		TotalPages: totalPages,
	}

	return response, nil
}

// BlockUser khóa tài khoản user với lý do
func (uc *AdminUsecase) BlockUser(ctx context.Context, req *dto.BlockUserRequest) error {

	// Kiểm tra user có tồn tại không
	_, err := uc.AdminRepo.GetUserByProfileID(req.ProfileID)
	if err != nil {
		return fmt.Errorf("user not found: %v", err)
	}

	// Cập nhật trạng thái user thành blocked
	err = uc.AdminRepo.UpdateUserStatus(req.ProfileID, 0) // 0 = blocked
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
	_, err := uc.AdminRepo.GetUserByProfileID(req.ProfileID)
	if err != nil {
		return fmt.Errorf("user not found: %v", err)
	}

	// Cập nhật trạng thái user thành active
	err = uc.AdminRepo.UpdateUserStatus(req.ProfileID, 1) // 1 = active
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
	_, err := uc.AdminRepo.GetUserByProfileID(req.ProfileID)
	if err != nil {
		return fmt.Errorf("user not found: %v", err)
	}

	// Cập nhật trạng thái user thành approved
	err = uc.AdminRepo.UpdateUserStatus(req.ProfileID, 2) // 2 = approved
	if err != nil {
		return fmt.Errorf("failed to approve user: %v", err)
	}

	// TODO: Log action và gửi notification
	// TODO: Lưu ghi chú duyệt vào bảng audit log

	return nil
}
