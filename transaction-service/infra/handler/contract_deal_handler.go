package handler

import (
	_dto "common/dto"
	_utils "common/utils"
	"context"
	sharepb "pb/types/shared"
	transactionpb "pb/types/transaction"

	"tx/infra/client"
	"tx/infra/mapper"
	"tx/infra/validator"
	"tx/internal/dto"
	transaction_usecase "tx/internal/usecases"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ContractDealHandler struct {
	transactionpb.UnimplementedDealContractServiceServer
	contractDealUsecase *transaction_usecase.DealContractUsecase
	transactionUsecase  transaction_usecase.TransactionUsecase
	dealCostUsecase     transaction_usecase.DealCostUsecase
	validator           *validator.ContractDealValidator
	mapper              *mapper.ContractDealMapper
	txtMapper           *mapper.TransactionMapper
	bdsproClient        *client.BdsproClient
	userClient          *client.UserClient
}

func NewContractDealHandler(
	contractDealUsecase *transaction_usecase.DealContractUsecase,
	transactionUsecase transaction_usecase.TransactionUsecase,
	dealCostUsecase transaction_usecase.DealCostUsecase,
	validator *validator.ContractDealValidator,
	mapper *mapper.ContractDealMapper,
	txtMapper *mapper.TransactionMapper,
	bdsproClient *client.BdsproClient,
	userClient *client.UserClient,
) *ContractDealHandler {
	return &ContractDealHandler{
		contractDealUsecase: contractDealUsecase,
		transactionUsecase:  transactionUsecase,
		dealCostUsecase:     dealCostUsecase,
		validator:           validator,
		mapper:              mapper,
		txtMapper:           txtMapper,
		bdsproClient:        bdsproClient,
		userClient:          userClient,
	}
}

// @Summary Lấy quá trình giao dịch thương vụ
// @Description Lấy quá trình giao dịch thương vụ
// @Tags Contract Deal
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của giao dịch"
// @Success 200 {object} transactionpb.ContractProcessResponse "Quá trình giao dịch thương vụ"
// @Failure 400 {object} status.Error "Lỗi yêu cầu không hợp lệ"
// @Failure 500 {object} status.Error "Lỗi server"
// @Router /contract/deal/process/{id} [get]
func (h *ContractDealHandler) DealContractProcess(ctx context.Context, req *sharepb.IdRequest) (*transactionpb.ContractProcessResponse, error) {
	processes, err := h.contractDealUsecase.GetDealContractProcess(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get deal contract process: %v", err)
	}

	data := h.txtMapper.ProcessDTOsToPb(processes)
	return &transactionpb.ContractProcessResponse{Data: data}, nil
}

// @Summary Tạo giao dịch thương vụ
// @Description Tạo giao dịch thương vụ
// @Tags Contract Deal
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body transactionpb.DealCostDTO true "Yêu cầu tạo giao dịch thương vụ"
// @Success 200 {object} transactionpb.DealContractResponse "Giao dịch thương vụ đã tạo"
// @Failure 400 {object} status.Error "Lỗi yêu cầu không hợp lệ"
// @Failure 500 {object} status.Error "Lỗi server"
// @Router /contract/deal/payment [post]
func (h *ContractDealHandler) CreateContractPayment(ctx context.Context, req *transactionpb.DealContractPaymentRequest) (*transactionpb.DealContractResponse, error) {
	// Get user ID from context
	userID := _utils.GetProfileIdWithContext(ctx)
	if userID == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}

	// Validate request
	if err := h.validator.ValidateDealContractPaymentRequest(ctx, req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	body, err := h.mapper.DealContractPaymentReqToEntity(req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	// Create contract transfer
	response, err := h.contractDealUsecase.CreateContractPayment(ctx, body)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create contract transfer: %v", err)
	}

	result := h.mapper.DealContractToPb(response)

	return result, nil
}

// // UpdateContract updates an existing contract
// func (h *ContractDealHandler) UpdateContract(ctx context.Context, req *transactionpb.UpdateDealContractRequest) (*transactionpb.DealContractResponse, error) {
// 	// Get user ID from context
// 	userID := _utils.GetProfileIdWithContext(ctx)
// 	if userID == 0 {
// 		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
// 	}

// 	// Validate request
// 	if err := h.validator.ValidateUpdateRequest(req); err != nil {
// 		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
// 	}

// 	// Update contract
// 	response, err := h.contractDealUsecase.UpdateContract(ctx, req, userID)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "failed to update contract: %v", err)
// 	}

// 	return response, nil
// }

// @Summary Xóa giao dịch thương vụ
// @Description Xóa giao dịch thương vụ
// @Tags Contract Deal
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của giao dịch"
// @Router /contract/deal/delete/{id} [delete]
func (h *ContractDealHandler) DeleteContract(ctx context.Context, req *sharepb.IdRequest) (*transactionpb.IdResponse, error) {
	// Delete contract
	if err := h.contractDealUsecase.DeleteContract(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete contract: %v", err)
	}

	return &transactionpb.IdResponse{Id: req.Id}, nil
}

// @Summary Lấy chi tiết giao dịch thương vụ
// @Description Lấy chi tiết giao dịch thương vụ
// @Tags Contract Deal
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của giao dịch"
// @Success 200 {object} transactionpb.DealContractResponse "Giao dịch thương vụ đã lấy"
// @Failure 400 {object} status.Error "Lỗi yêu cầu không hợp lệ"
// @Failure 500 {object} status.Error "Lỗi server"
// @Router /contract/deal/detail/{id} [get]
func (h *ContractDealHandler) GetContract(ctx context.Context, req *sharepb.IdRequest) (*transactionpb.DealContractResponse, error) {
	// Validate request
	if req.Id == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "transaction ID is required")
	}

	contract, _ := h.contractDealUsecase.GetContractDealByTransactionID(ctx, req.Id)
	dealContractPb := h.mapper.DealContractToPb(contract)

	h.bdsproClient.MapDealContractToPb(ctx, dealContractPb)
	h.userClient.MapProfileToDealContractPb(ctx, dealContractPb)
	return dealContractPb, nil
}

// @Summary Lấy danh sách giao dịch thương vụ
// @Description Lấy danh sách giao dịch thương vụ
// @Tags Contract Deal
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của giao dịch"
// @Success 200 {object} transactionpb.DealContractListResponse "Danh sách giao dịch thương vụ"
// @Failure 400 {object} status.Error "Lỗi yêu cầu không hợp lệ"
// @Failure 500 {object} status.Error "Lỗi server"
// @Router /contract/deal/list/{id} [get]
func (h *ContractDealHandler) GetContracts(ctx context.Context, req *transactionpb.GetDealContractRequest) (*transactionpb.DealContractListResponse, error) {
	// Get user ID from context
	userID := _utils.GetProfileIdWithContext(ctx)
	if userID == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}

	query := dto.DealSearch{
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
		},
	}
	// Get contracts
	response, total, err := h.contractDealUsecase.GetDealContractByDealID(ctx, req.DealId, query)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get contracts: %v", err)
	}

	data := h.mapper.DealContractItemToPbs(response)
	h.userClient.MapProfileToDealContractPbs(ctx, data)

	return &transactionpb.DealContractListResponse{
		Data:  data,
		Total: int32(total),
	}, nil
}

// @Summary Tạo hành động thanh toán
// @Description Tạo hành động thanh toán
// @Tags Contract Deal
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body transactionpb.PaymentActionRequest true "Yêu cầu tạo hành động thanh toán"
// @Success 200 {object} transactionpb.ActionResponse "Hành động thanh toán đã tạo"
// @Failure 400 {object} status.Error "Lỗi yêu cầu không hợp lệ"
// @Failure 500 {object} status.Error "Lỗi server"
// @Router /contract/deal/payment/action [post]
func (h *ContractDealHandler) CreateActionPayment(ctx context.Context, req *transactionpb.PaymentActionRequest) (*transactionpb.ActionResponse, error) {
	body := h.mapper.PaymentActionReqToDTO(req)
	action, err := h.contractDealUsecase.CreateActionPayment(ctx, body)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create action payment: %v", err)
	}

	result := h.mapper.ActionToPb(action)
	return result, nil
}

// @Summary Lấy thống kê tổng hợp Chi phí/Doanh thu và Payment
// @Description Lấy thống kê tổng hợp Chi phí/Doanh thu và Payment theo Deal ID với các tham số tùy chọn
// @Tags DealCost
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param dealId query uint64 true "Deal ID"
// @Param payment query bool false "Include payment statistics"
// @Param cost query bool false "Include cost statistics"
// @Success 200 {object} transactionpb.ContractDealStatisticsDTO
// @Router /contract/deal/statistics/{dealId} [get]
func (h *ContractDealHandler) GetStatistics(ctx context.Context, req *transactionpb.DealCostStatisticsRequestDTO) (*transactionpb.ContractDealStatisticsDTO, error) {
	statistics, err := h.contractDealUsecase.GetStatisticsByDealID(ctx, req.DealId, req.Payment, req.Cost)
	if err != nil {
		return nil, err
	}
	return h.mapper.StatisticsToPb(statistics), nil
}
