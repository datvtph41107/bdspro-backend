package admin_handler

import (
	"bdspro/infra/client"
	"bdspro/infra/mapper"
	admin_usecases "bdspro/internal/usecases/admin"
	_routes "common/routes"
	"context"
	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"
)

type AdminTransactionHandler struct {
	bdspropb.UnimplementedAdminTransactionServiceServer
	AdminTransactionUsecase *admin_usecases.AdminTransactionUsecase
	TransactionMapper       *mapper.AdminTransactionMapper
	AuthGrpcClient          *client.AuthClient
}

func NewAdminTransactionHandler(
	uc *admin_usecases.AdminTransactionUsecase,
	transactionMapper *mapper.AdminTransactionMapper,
	authGrpcClient *client.AuthClient,
) *AdminTransactionHandler {
	return &AdminTransactionHandler{
		AdminTransactionUsecase: uc,
		TransactionMapper:       transactionMapper,
		AuthGrpcClient:          authGrpcClient,
	}
}

// @Summary Lấy danh sách thương vụ
// @Description Lấy danh sách tất cả thương vụ trong hệ thống (dành cho admin)
// @Tags AdminTransaction
// @Accept json
// @Produce json
// @Param page query int true "Trang hiện tại"
// @Param size query int true "Số lượng items mỗi trang"
// @Param keyword query string false "Từ khóa tìm kiếm"
// @Param transactionType query int false "Loại giao dịch (10: Bán, 20: Thuê)"
// @Param transactionStatus query int false "Trạng thái giao dịch (10: Nháp, 20: Đặt cọc, 30: Đã ký, 40: Hoàn thành, 50: Hủy)"
// @Param approvalStatus query string false "Trạng thái phê duyệt (pending, approved, rejected)"
// @Param startDate query string false "Ngày bắt đầu (RFC3339)"
// @Param endDate query string false "Ngày kết thúc (RFC3339)"
// @Param productId query int false "ID sản phẩm"
// @Param contactId query int false "ID liên hệ"
// @Param ownerId query int false "ID chủ sở hữu"
// @Param ownerType query int false "Loại chủ sở hữu (10: Member, 20: Group, 30: Organization)"
// @Success 200 {object} bdspropb.AdminTransactionSearchResponse
// @Router /v2/bdspro/admin/transactions [get]
func (h *AdminTransactionHandler) GetTransactions(ctx context.Context, req *bdspropb.AdminTransactionSearchRequest) (*bdspropb.AdminTransactionSearchResponse, error) {
	// Check permission
	if err := h.AuthGrpcClient.HasPermissions(ctx, []string{"ADMIN_TRANSACTION_VIEW"}); err != nil {
		return nil, err
	}

	// Map request to filter
	filter := h.TransactionMapper.RequestToFilter(req)

	// Get list from usecase
	transactions, total, err := h.AdminTransactionUsecase.GetList(ctx, filter)
	if err != nil {
		return nil, &_routes.Except{
			Code:    500,
			Message: "Lỗi khi lấy danh sách thương vụ: " + err.Error(),
		}
	}

	// Map to response
	return h.TransactionMapper.EntitiesToListResponse(ctx, transactions, total), nil
}

// @Summary Lấy chi tiết thương vụ
// @Description Lấy thông tin chi tiết của một thương vụ (dành cho admin)
// @Tags AdminTransaction
// @Accept json
// @Produce json
// @Param id path int true "ID thương vụ"
// @Success 200 {object} bdspropb.AdminTransactionDetail
// @Router /v2/bdspro/admin/transactions/{id} [get]
func (h *AdminTransactionHandler) GetTransactionDetail(ctx context.Context, req *sharepb.IdRequest) (*bdspropb.AdminTransactionDetail, error) {
	// Check permission
	if err := h.AuthGrpcClient.HasPermissions(ctx, []string{"ADMIN_TRANSACTION_VIEW"}); err != nil {
		return nil, err
	}

	// Get detail from usecase
	transaction, err := h.AdminTransactionUsecase.GetDetail(ctx, req.Id)
	if err != nil {
		return nil, &_routes.Except{
			Code:    404,
			Message: "Không tìm thấy thương vụ",
		}
	}

	// Map to response
	return h.TransactionMapper.EntityToDetail(ctx, transaction), nil
}

// @Summary Phê duyệt thương vụ
// @Description Phê duyệt một thương vụ (dành cho admin)
// @Tags AdminTransaction
// @Accept json
// @Produce json
// @Param id path int true "ID thương vụ"
// @Success 200 {object} bdspropb.Response
// @Router /v2/bdspro/admin/transactions/{id}/approve [post]
func (h *AdminTransactionHandler) ApproveTransaction(ctx context.Context, req *sharepb.IdRequest) (*bdspropb.Response, error) {
	// Check permission
	if err := h.AuthGrpcClient.HasPermissions(ctx, []string{"ADMIN_TRANSACTION_APPROVE"}); err != nil {
		return nil, err
	}

	// Approve transaction
	err := h.AdminTransactionUsecase.ApproveTransaction(ctx, req.Id)
	if err != nil {
		return nil, &_routes.Except{
			Code:    500,
			Message: "Lỗi khi phê duyệt thương vụ: " + err.Error(),
		}
	}

	return &bdspropb.Response{
		Id:      req.Id,
		Message: "Phê duyệt thương vụ thành công",
	}, nil
}

// @Summary Từ chối thương vụ
// @Description Từ chối một thương vụ (dành cho admin)
// @Tags AdminTransaction
// @Accept json
// @Produce json
// @Param id path int true "ID thương vụ"
// @Param body body bdspropb.AdminRejectTransactionRequest true "Thông tin từ chối"
// @Success 200 {object} bdspropb.Response
// @Router /v2/bdspro/admin/transactions/{id}/reject [post]
func (h *AdminTransactionHandler) RejectTransaction(ctx context.Context, req *bdspropb.AdminRejectTransactionRequest) (*bdspropb.Response, error) {
	// Check permission
	if err := h.AuthGrpcClient.HasPermissions(ctx, []string{"ADMIN_TRANSACTION_APPROVE"}); err != nil {
		return nil, err
	}

	// Validate reason
	if req.Reason == "" {
		return nil, &_routes.Except{
			Code:    400,
			Message: "Vui lòng nhập lý do từ chối",
		}
	}

	// Reject transaction
	err := h.AdminTransactionUsecase.RejectTransaction(ctx, req.Id)
	if err != nil {
		return nil, &_routes.Except{
			Code:    500,
			Message: "Lỗi khi từ chối thương vụ: " + err.Error(),
		}
	}

	return &bdspropb.Response{
		Id:      req.Id,
		Message: "Từ chối thương vụ thành công",
	}, nil
}

// @Summary Xóa thương vụ
// @Description Xóa một thương vụ (soft delete) (dành cho admin)
// @Tags AdminTransaction
// @Accept json
// @Produce json
// @Param id path int true "ID thương vụ"
// @Success 200 {object} bdspropb.Response
// @Router /v2/bdspro/admin/transactions/{id} [delete]
func (h *AdminTransactionHandler) DeleteTransaction(ctx context.Context, req *sharepb.IdRequest) (*bdspropb.Response, error) {
	// Check permission
	if err := h.AuthGrpcClient.HasPermissions(ctx, []string{"ADMIN_TRANSACTION_DELETE"}); err != nil {
		return nil, err
	}

	// Delete transaction
	err := h.AdminTransactionUsecase.DeleteTransaction(ctx, req.Id)
	if err != nil {
		return nil, &_routes.Except{
			Code:    500,
			Message: "Lỗi khi xóa thương vụ: " + err.Error(),
		}
	}

	return &bdspropb.Response{
		Id:      req.Id,
		Message: "Xóa thương vụ thành công",
	}, nil
}
