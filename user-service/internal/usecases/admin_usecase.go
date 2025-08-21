package usecases

import (
	"context"
	"math"
	user_repo "user/infra/repo"
	"user/internal/dto"
)

type AdminUsecase struct {
	AdminRepo *user_repo.AdminRepo
}

func NewAdminUsecase(adminRepo *user_repo.AdminRepo) *AdminUsecase {
	return &AdminUsecase{
		AdminRepo: adminRepo,
	}
}

// ListAllUsers lấy danh sách tất cả user với phân trang và tìm kiếm
func (uc *AdminUsecase) ListAllUsers(ctx context.Context, req *dto.ListAllUsersRequest) (*dto.ListAllUsersResponse, error) {
	// Kiểm tra quyền admin (sẽ được implement sau)
	// profileID := _utils.GetProfileIdWithContext(ctx)
	// if !uc.isAdmin(profileID) {
	//     return nil, errors.New("unauthorized: admin access required")
	// }

	// Validate request
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Size == 0 {
		req.Size = 20
	}
	if req.Size > 100 {
		req.Size = 100
	}

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
