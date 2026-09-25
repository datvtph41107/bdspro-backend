package handler

import (
	_dto "common/domain/dto"
	_utils "common/utils"
	"context"
	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"
	"time"

	"bdspro/infra/client"
	"bdspro/infra/mapper"
	"bdspro/infra/validator"
	"bdspro/internal/dto"
	transaction_usecase "bdspro/internal/usecases"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TxContractDealHandler struct {
	bdspropb.UnimplementedDealContractServiceServer
	contractDealUsecase *transaction_usecase.DealContractUsecase
	transactionUsecase  transaction_usecase.TxTransactionUsecase
	dealCostUsecase     transaction_usecase.DealCostUsecase
	validator           *validator.ContractDealValidator
	mapper              *mapper.ContractDealMapper
	txtMapper           *mapper.TxMapper
	userClient          *client.UserClient
	SyncProvider        *_utils.SyncUtil
}

func NewTxContractDealHandler(
	contractDealUsecase *transaction_usecase.DealContractUsecase,
	transactionUsecase transaction_usecase.TxTransactionUsecase,
	dealCostUsecase transaction_usecase.DealCostUsecase,
	validator *validator.ContractDealValidator,
	mapper *mapper.ContractDealMapper,
	txtMapper *mapper.TxMapper,
	userClient *client.UserClient,
	syncProvider *_utils.SyncUtil,
) *TxContractDealHandler {
	return &TxContractDealHandler{
		contractDealUsecase: contractDealUsecase,
		transactionUsecase:  transactionUsecase,
		dealCostUsecase:     dealCostUsecase,
		validator:           validator,
		mapper:              mapper,
		txtMapper:           txtMapper,
		userClient:          userClient,
		SyncProvider:        syncProvider,
	}
}

func (h *TxContractDealHandler) bumpContractsSync(ctx context.Context, dealID uint64) {
	if h.SyncProvider == nil || dealID == 0 {
		return
	}
	now := time.Now().UnixMilli()
	_ = h.SyncProvider.PutTimestamp(ctx, h.SyncProvider.GetKey(ctx, _utils.SyncKeyDealContracts, dealID), now)
}

// @Summary Lấy quá trình giao dịch thương vụ
// @Description Lấy quá trình giao dịch thương vụ
// @Tags Contract Deal
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của giao dịch"
// @Success 200 {object} bdspropb.ContractProcessResponse "Quá trình giao dịch thương vụ"
// @Failure 400 {object} status.Error "Lỗi yêu cầu không hợp lệ"
// @Failure 500 {object} status.Error "Lỗi server"
// @Router /contract/deal/process/{id} [get]
func (h *TxContractDealHandler) DealContractProcess(ctx context.Context, req *sharepb.IdRequest) (*bdspropb.ContractProcessResponse, error) {
	processes, err := h.contractDealUsecase.GetDealContractProcess(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get deal contract process: %v", err)
	}

	data := h.txtMapper.ProcessDTOsToPb(processes)
	now := time.Now().UnixMilli()
	if h.SyncProvider != nil && req.Id > 0 {
		_ = h.SyncProvider.PutTimestamp(ctx, h.SyncProvider.GetKey(ctx, _utils.SyncKeyDealProcess, req.Id), now)
	}
	return &bdspropb.ContractProcessResponse{Data: data}, nil
}

// @Summary Tạo giao dịch thương vụ
// @Description Tạo giao dịch thương vụ
// @Tags Contract Deal
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body bdspropb.DealCostDTO true "Yêu cầu tạo giao dịch thương vụ"
// @Success 200 {object} bdspropb.DealContractResponse "Giao dịch thương vụ đã tạo"
// @Failure 400 {object} status.Error "Lỗi yêu cầu không hợp lệ"
// @Failure 500 {object} status.Error "Lỗi server"
// @Router /contract/deal/payment [post]
func (h *TxContractDealHandler) CreateContractPayment(ctx context.Context, req *bdspropb.DealContractPaymentRequest) (*bdspropb.DealContractResponse, error) {
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

	h.bumpContractsSync(ctx, req.DealId)
	result := h.mapper.DealContractToPb(response)

	return result, nil
}

// // UpdateContract updates an existing contract
// func (h *ContractDealHandler) UpdateContract(ctx context.Context, req *bdspropb.UpdateDealContractRequest) (*bdspropb.DealContractResponse, error) {
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
func (h *TxContractDealHandler) DeleteContract(ctx context.Context, req *sharepb.IdRequest) (*bdspropb.IdResponse, error) {
	var dealID uint64
	if contract, _ := h.contractDealUsecase.GetContractDealByTransactionID(ctx, req.Id); contract != nil {
		dealID = contract.DealID
	}

	if err := h.contractDealUsecase.DeleteContract(ctx, req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete contract: %v", err)
	}

	h.bumpContractsSync(ctx, dealID)
	return &bdspropb.IdResponse{Id: req.Id}, nil
}

// @Summary Lấy chi tiết giao dịch thương vụ
// @Description Lấy chi tiết giao dịch thương vụ
// @Tags Contract Deal
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của giao dịch"
// @Success 200 {object} bdspropb.DealContractResponse "Giao dịch thương vụ đã lấy"
// @Failure 400 {object} status.Error "Lỗi yêu cầu không hợp lệ"
// @Failure 500 {object} status.Error "Lỗi server"
// @Router /contract/deal/detail/{id} [get]
func (h *TxContractDealHandler) GetContract(ctx context.Context, req *sharepb.IdRequest) (*bdspropb.DealContractResponse, error) {
	// Validate request
	if req.Id == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "transaction ID is required")
	}

	contract, _ := h.contractDealUsecase.GetContractDealByTransactionID(ctx, req.Id)
	dealContractPb := h.mapper.DealContractToPb(contract)
	// todo:
	// h.bdsproClient.MapDealContractToPb(ctx, dealContractPb)
	// todo:
	// h.userClient.MapProfileToDealContractPb(ctx, dealContractPb)
	return dealContractPb, nil
}

// @Summary Lấy danh sách giao dịch thương vụ
// @Description Lấy danh sách giao dịch thương vụ
// @Tags Contract Deal
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID của giao dịch"
// @Success 200 {object} bdspropb.DealContractListResponse "Danh sách giao dịch thương vụ"
// @Failure 400 {object} status.Error "Lỗi yêu cầu không hợp lệ"
// @Failure 500 {object} status.Error "Lỗi server"
// @Router /contract/deal/list/{id} [get]
func (h *TxContractDealHandler) GetContracts(ctx context.Context, req *bdspropb.GetDealContractRequest) (*bdspropb.DealContractListResponse, error) {
	userID := _utils.GetProfileIdWithContext(ctx)
	if userID == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}

	if h.SyncProvider != nil && req.DealId > 0 {
		key := h.SyncProvider.GetKey(ctx, _utils.SyncKeyDealContracts, req.DealId)
		if !h.SyncProvider.HasUpdated(ctx, key, req.GetTimestamp()) && req.GetPage() == 0 {
			return &bdspropb.DealContractListResponse{
				Timestamp: h.SyncProvider.GetTimestamp(ctx, key),
			}, nil
		}
	}

	query := dto.TxDealSearch{
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
		},
	}
	response, total, err := h.contractDealUsecase.GetDealContractByDealID(ctx, req.DealId, query)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get contracts: %v", err)
	}

	data := h.mapper.DealContractItemToPbs(ctx, response)

	var ts int64
	if h.SyncProvider != nil && req.DealId > 0 {
		ts = h.SyncProvider.EnsureTimestamp(ctx, h.SyncProvider.GetKey(ctx, _utils.SyncKeyDealContracts, req.DealId))
	}

	return &bdspropb.DealContractListResponse{
		Data:      data,
		Total:     int32(total),
		Timestamp: ts,
	}, nil
}

// @Summary Tạo hành động thanh toán
// @Description Tạo hành động thanh toán
// @Tags Contract Deal
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body bdspropb.PaymentActionRequest true "Yêu cầu tạo hành động thanh toán"
// @Success 200 {object} bdspropb.ActionResponse "Hành động thanh toán đã tạo"
// @Failure 400 {object} status.Error "Lỗi yêu cầu không hợp lệ"
// @Failure 500 {object} status.Error "Lỗi server"
// @Router /contract/deal/payment/action [post]
func (h *TxContractDealHandler) CreateActionPayment(ctx context.Context, req *bdspropb.PaymentActionRequest) (*bdspropb.ActionResponse, error) {
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
// @Success 200 {object} bdspropb.ContractDealStatisticsDTO
// @Router /contract/deal/statistics/{dealId} [get]
func (h *TxContractDealHandler) GetStatistics(ctx context.Context, req *bdspropb.DealCostStatisticsRequestDTO) (*bdspropb.ContractDealStatisticsDTO, error) {
	if h.SyncProvider != nil && req.DealId > 0 {
		key := h.SyncProvider.GetKey(ctx, _utils.SyncKeyDealStatistics, req.DealId)
		if !h.SyncProvider.HasUpdated(ctx, key, req.GetTimestamp()) {
			return &bdspropb.ContractDealStatisticsDTO{}, nil
		}
	}

	statistics, err := h.contractDealUsecase.GetStatisticsByDealID(ctx, req.DealId, req.Payment, req.Cost)
	if err != nil {
		return nil, err
	}
	if h.SyncProvider != nil && req.DealId > 0 {
		_ = h.SyncProvider.EnsureTimestamp(ctx, h.SyncProvider.GetKey(ctx, _utils.SyncKeyDealStatistics, req.DealId))
	}
	return h.mapper.StatisticsToPb(statistics), nil
}

// @Summary Lấy danh sách contract của user hiện tại
// @Description Lấy danh sách contract của user hiện tại
// @Tags Contract Deal
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int32 false "Page number" default(1)
// @Param size query int32 false "Page size" default(20)
// @Success 200 {object} bdspropb.DealContractListResponse "Danh sách contract của user"
// @Failure 400 {object} status.Error "Lỗi yêu cầu không hợp lệ"
// @Failure 500 {object} status.Error "Lỗi server"
// @Router /contract/deal/me [get]
func (h *TxContractDealHandler) GetMyContracts(ctx context.Context, req *sharepb.IdRequest) (*bdspropb.DealContractListResponse, error) {
	// Get user ID from context
	userID := _utils.GetProfileIdWithContext(ctx)
	if userID == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated")
	}

	query := dto.TxDealSearch{
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
		},
	}

	// Get contracts
	response, total, err := h.contractDealUsecase.GetDealContractByProfileID(ctx, userID, query)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get contracts: %v", err)
	}

	data := h.mapper.DealContractItemToPbs(ctx, response)
	h.userClient.MapProfileToDealContractPbs(ctx, data)

	return &bdspropb.DealContractListResponse{
		Data:  data,
		Total: int32(total),
	}, nil
}
