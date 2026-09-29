package handler

import (
	"bdspro/infra/mapper"
	"bdspro/infra/validator"
	"bdspro/internal/usecases"
	"context"
	bdspropb "pb/types/bdspro"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TransactionHandler struct {
	bdspropb.UnimplementedTransactionServiceServer
	// TransactionUsecase transaction_usecase.TransactionUsecase
	// Transformer mapper.TransactionTransformer
	Validator             validator.TransactionValidator
	TransferTaxFeeUsecase *usecases.TransferTaxFeeUsecase
	LoanPaymentUsecase    *usecases.LoanPaymentUsecase
}

func NewTransactionHandler(
	// transactionUsecase transaction_usecase.TransactionUsecase,
	// transformer mapper.TransactionTransformer,
	validator validator.TransactionValidator,
	transferTaxFeeUsecase *usecases.TransferTaxFeeUsecase,
	loanPaymentUsecase *usecases.LoanPaymentUsecase,
) *TransactionHandler {
	return &TransactionHandler{
		// TransactionUsecase: transactionUsecase,
		// Transformer: transformer,
		Validator:             validator,
		TransferTaxFeeUsecase: transferTaxFeeUsecase,
		LoanPaymentUsecase:    loanPaymentUsecase,
	}
}

// // @Summary Lấy danh sách giao dịch
// // @Tags Giao dịch
// // @Produce json
// // @Param page query int false "Trang"
// // @Param size query int false "Kích thước trang"
// // @Param transactionType query string false "Loại giao dịch"
// // @Param approvalStatus query string false "Trạng thái phê duyệt"
// // @Param transactionStatus query string false "Trạng thái giao dịch"
// // @Param startDate query string false "Ngày bắt đầu"
// // @Param endDate query string false "Ngày kết thúc"
// // @Security BearerAuth
// // @Router /v2/bdspro/v2/transaction/list [get]
// func (h *TransactionHandler) ListTransactions(ctx context.Context, req *bdspropb.ListTransactionsRequest) (*bdspropb.ListTransactionsResponse, error) {
// 	if err := h.Validator.ValidateListTransactionsRequest(req); err != nil {
// 		return nil, status.Error(codes.InvalidArgument, err.Error())
// 	}
// 	filters := &repo.TransactionFilters{
// 		Page: 0,
// 		Size: 10,
// 	}

// 	if req.TransactionType != nil {
// 		filters.TransactionType = req.TransactionType
// 	}
// 	if req.ApprovalStatus != nil {
// 		filters.ApprovalStatus = req.ApprovalStatus
// 	}
// 	if req.TransactionStatus != nil {
// 		filters.TransactionStatus = req.TransactionStatus
// 	}
// 	if req.StartDate != nil {
// 		date, err := time.Parse(time.RFC3339, *req.StartDate)
// 		if err != nil {
// 			return nil, status.Error(codes.InvalidArgument, err.Error())
// 		}
// 		filters.StartDate = &date
// 	}
// 	if req.EndDate != nil {
// 		date, err := time.Parse(time.RFC3339, *req.EndDate)
// 		if err != nil {
// 			return nil, status.Error(codes.InvalidArgument, err.Error())
// 		}
// 		filters.EndDate = &date
// 	}
// 	if req.Page != nil {
// 		filters.Page = int(*req.Page)
// 	}
// 	if req.Size != nil {
// 		filters.Size = int(*req.Size)
// 	}
// 	transactions, total, err := h.TransactionUsecase.ListTransactions(ctx, filters)
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return h.Transformer.EntityToListTransactionsResponse(transactions, int32(total)), nil
// }

// // @Summary Lấy chi tiết giao dịch
// // @Tags Giao dịch
// // @Produce json
// // @Param id path uint32 true "ID"
// // @Security BearerAuth
// // @Router /v2/bdspro/v2/transaction/{id} [get]
// func (h *TransactionHandler) GetTransactionDetail(ctx context.Context, req *bdspropb.GetTransactionDetailRequest) (*bdspropb.Transaction, error) {
// 	if err := h.Validator.ValidateGetTransactionRequest(req); err != nil {
// 		return nil, status.Error(codes.InvalidArgument, err.Error())
// 	}

// 	transaction, err := h.TransactionUsecase.GetTransaction(ctx, uint64(req.Id))
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return h.Transformer.EntityToTransactionResponse(transaction), nil
// }

// // @Summary Tạo giao dịch
// // @Tags Giao dịch
// // @Produce json
// // @Param transaction body bdspropb.CreateTransactionRequest true "Thông tin giao dịch"
// // @Security BearerAuth
// // @Router /v2/bdspro/v2/transaction [post]
// func (h *TransactionHandler) CreateTransaction(ctx context.Context, req *bdspropb.CreateTransactionRequest) (*bdspropb.Transaction, error) {
// 	if err := h.Validator.ValidateCreateTransactionRequest(req); err != nil {
// 		return nil, status.Error(codes.InvalidArgument, err.Error())
// 	}

// 	transaction := h.Transformer.CreateTransactionRequestToEntity(req)

// 	entity, err := h.TransactionUsecase.CreateTransaction(ctx, transaction)
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return h.Transformer.EntityToTransactionResponse(entity), nil
// }

// // @Summary Cập nhật giao dịch
// // @Tags Giao dịch
// // @Produce json
// // @Param id path uint32 true "ID"
// // @Param transaction body bdspropb.UpdateTransactionRequest true "Thông tin giao dịch"
// // @Security BearerAuth
// // @Router /v2/bdspro/v2/transaction/{id} [put]
// func (h *TransactionHandler) UpdateTransaction(ctx context.Context, req *bdspropb.UpdateTransactionRequest) (*bdspropb.Transaction, error) {
// 	if err := h.Validator.ValidateUpdateTransactionRequest(req); err != nil {
// 		return nil, status.Error(codes.InvalidArgument, err.Error())
// 	}

// 	transaction := h.Transformer.UpdateTransactionRequestToEntity(req)

// 	entity, err := h.TransactionUsecase.UpdateTransaction(ctx, transaction)
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return h.Transformer.EntityToTransactionResponse(entity), nil
// }

// // @Summary Xóa giao dịch
// // @Tags Giao dịch
// // @Produce json
// // @Param id path uint32 true "ID"
// // @Security BearerAuth
// // @Router /v2/bdspro/v2/transaction/{id} [delete]
// func (h *TransactionHandler) DeleteTransaction(ctx context.Context, req *bdspropb.DeleteTransactionRequest) (*bdspropb.DeleteTransactionResponse, error) {
// 	if err := h.Validator.ValidateDeleteTransactionRequest(req); err != nil {
// 		return nil, status.Error(codes.InvalidArgument, err.Error())
// 	}

// 	if err := h.TransactionUsecase.DeleteTransaction(ctx, uint64(req.Id)); err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return &bdspropb.DeleteTransactionResponse{
// 		Id: req.Id,
// 	}, nil
// }

// // @Summary Phê duyệt giao dịch
// // @Tags Giao dịch
// // @Produce json
// // @Param id path uint32 true "ID"
// // @Security BearerAuth
// // @Router /v2/bdspro/v2/transaction/{id}/approve [put]
// func (h *TransactionHandler) ApproveTransaction(ctx context.Context, req *bdspropb.ApproveTransactionRequest) (*bdspropb.Transaction, error) {
// 	if err := h.Validator.ValidateApproveTransactionRequest(req); err != nil {
// 		return nil, status.Error(codes.InvalidArgument, err.Error())
// 	}

// 	if err := h.TransactionUsecase.ApproveTransaction(ctx, uint64(req.Id), 0); err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	transaction, err := h.TransactionUsecase.GetTransaction(ctx, uint64(req.Id))
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return h.Transformer.EntityToTransactionResponse(transaction), nil
// }

// // @Summary Từ chối giao dịch
// // @Tags Giao dịch
// // @Produce json
// // @Param id path uint32 true "ID"
// // @Security BearerAuth
// // @Router /v2/bdspro/v2/transaction/{id}/reject [put]
// func (h *TransactionHandler) RejectTransaction(ctx context.Context, req *bdspropb.RejectTransactionRequest) (*bdspropb.Transaction, error) {
// 	if err := h.Validator.ValidateRejectTransactionRequest(req); err != nil {
// 		return nil, status.Error(codes.InvalidArgument, err.Error())
// 	}

// 	if err := h.TransactionUsecase.RejectTransaction(ctx, uint64(req.Id)); err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	transaction, err := h.TransactionUsecase.GetTransaction(ctx, uint64(req.Id))
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return h.Transformer.EntityToTransactionResponse(transaction), nil
// }

// // @Summary Lấy danh sách loại giao dịch
// // @Tags Giao dịch
// // @Produce json
// // @Security BearerAuth
// // @Router /v2/bdspro/v2/transaction/types [get]
// func (h *TransactionHandler) GetTransactionTypes(ctx context.Context, req *bdspropb.GetTransactionTypesRequest) (*bdspropb.GetTransactionTypesResponse, error) {
// 	if err := h.Validator.ValidateGetTransactionTypesRequest(req); err != nil {
// 		return nil, status.Error(codes.InvalidArgument, err.Error())
// 	}

// 	types, err := h.TransactionUsecase.GetTransactionTypes(ctx)
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, err.Error())
// 	}

// 	return h.Transformer.EntityToTransactionTypesResponse(types), nil
// }

// @Summary Tính thuế và phí chuyển nhượng
// @Description API tính toán thuế và phí chuyển nhượng bất động sản bao gồm: thuế TNCN (nếu cá nhân), thuế TN doanh nghiệp (nếu doanh nghiệp), lệ phí trước bạ, phí công chứng
// @Tags Giao dịch
// @Accept json
// @Produce json
// @Param body body bdspropb.CalculateTransferTaxFeeRequest true "Thông tin tính thuế phí"
// @Security BearerAuth
// @Router /v2/bdspro/v2/transaction/calculate-tax-fee [post]
func (h *TransactionHandler) CalculateTransferTaxFee(ctx context.Context, req *bdspropb.CalculateTransferTaxFeeRequest) (*bdspropb.CalculateTransferTaxFeeResponse, error) {
	if req.TransferValue <= 0 {
		return nil, status.Error(codes.InvalidArgument, "Giá trị chuyển nhượng phải lớn hơn 0")
	}

	// Convert proto request to DTO
	dtoReq := mapper.MapCalculateTransferTaxFeeRequestToDTO(req)

	// Calculate tax and fee
	result, err := h.TransferTaxFeeUsecase.CalculateTransferTaxFee(ctx, dtoReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Convert DTO response to proto
	return mapper.MapCalculateTransferTaxFeeResponseToProto(result), nil
}

// @Summary Tính lịch trả nợ vay ngân hàng
// @Description API tính toán lịch trả nợ vay ngân hàng theo 2 phương thức: dư nợ giảm dần (10) và trả lãi định kỳ/gốc cuối kỳ (20)
// @Tags Giao dịch
// @Accept json
// @Produce json
// @Param body body bdspropb.CalculateLoanPaymentRequest true "Thông tin tính vay"
// @Security BearerAuth
// @Router /v2/bdspro/v2/transaction/calculate-loan-payment [post]
func (h *TransactionHandler) CalculateLoanPayment(ctx context.Context, req *bdspropb.CalculateLoanPaymentRequest) (*bdspropb.CalculateLoanPaymentResponse, error) {
	if req.LoanAmount <= 0 {
		return nil, status.Error(codes.InvalidArgument, "Số tiền vay phải lớn hơn 0")
	}
	if req.InterestRate <= 0 {
		return nil, status.Error(codes.InvalidArgument, "Lãi suất phải lớn hơn 0")
	}
	if req.Term <= 0 {
		return nil, status.Error(codes.InvalidArgument, "Kỳ hạn phải lớn hơn 0")
	}
	if req.PaymentMethod != 10 && req.PaymentMethod != 20 {
		return nil, status.Error(codes.InvalidArgument, "Phương thức trả nợ phải là 10 (dư nợ giảm dần) hoặc 20 (trả lãi định kỳ/gốc cuối kỳ)")
	}

	// Convert proto request to DTO
	dtoReq := mapper.MapCalculateLoanPaymentRequestToDTO(req)

	// Calculate loan payment schedule
	result, err := h.LoanPaymentUsecase.CalculateLoanPayment(ctx, dtoReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// Convert DTO response to proto
	return mapper.MapCalculateLoanPaymentResponseToProto(result), nil
}
