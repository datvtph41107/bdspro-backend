package admin_handler

import (
	_dto "common/domain/dto"
	"context"
	bdspropb "pb/types/bdspro"

	"bdspro/infra/mapper"
	"bdspro/internal/dto"
	"bdspro/internal/usecases"
)

type AdminTxHandler struct {
	bdspropb.UnimplementedAdminTxServiceServer
	TransactionUsecase usecases.TxTransactionUsecase
	TransactionMapper  *mapper.TxMapper
}

func NewAdminTxHandler(
	transactionUsecase usecases.TxTransactionUsecase,
	transactionMapper *mapper.TxMapper,
) *AdminTxHandler {
	return &AdminTxHandler{
		TransactionUsecase: transactionUsecase,
		TransactionMapper:  transactionMapper,
	}
}

// @Summary Lấy danh sách giao dịch (Admin)
// @Description Admin lấy danh sách tất cả giao dịch trong hệ thống
// @Tags AdminTransactionService
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query uint64 false "Số trang" default(1)
// @Param size query uint64 false "Số lượng bản ghi" default(20)
// @Param organizationId query uint64 false "ID tổ chức"
// @Param status query int32 false "Trạng thái giao dịch"
// @Param method query int32 false "Phương thức giao dịch"
// @Param search query string false "Tìm kiếm theo tên giao dịch"
// @Success 200 {object} bdspropb.TransactionListResponse "Danh sách giao dịch"
// @Failure 400 {object} error "Bad Request"
// @Failure 500 {object} error "Internal Server Error"
// @Router /admin/transaction/list [get]
func (h *AdminTxHandler) AdminListTransactions(ctx context.Context, req *bdspropb.AdminTransactionListRequest) (*bdspropb.TransactionListResponse, error) {
	filters := &dto.TxTransactionFilters{
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
		},
	}

	// Áp dụng các filter nếu có
	if req.Status != bdspropb.TransactionStatus_TRANSACTION_STATUS_UNKNOWN {
		status := uint32(req.Status)
		filters.TransactionStatus = &status
	}

	transactions, total, err := h.TransactionUsecase.ListTransactions(ctx, filters)
	if err != nil {
		return nil, err
	}

	transactionsResponse := make([]*bdspropb.TransactionResponse, len(transactions))
	for i, transaction := range transactions {
		transactionsResponse[i] = h.TransactionMapper.TxToPb(transaction)
	}

	return &bdspropb.TransactionListResponse{
		Data:  transactionsResponse,
		Total: total,
	}, nil
}
