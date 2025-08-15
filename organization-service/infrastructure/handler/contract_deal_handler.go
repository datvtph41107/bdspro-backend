package handler

import (
	"organization/infrastructure/client"
	"organization/internal/domain/repository"
	"organization/internal/usecase"
)

type ContractDealHandler struct {
	// contractpb.UnimplementedDealContractServiceServer
	contractDealUsecase usecase.DealTransactionUsecase
	dealActionUsecase   *usecase.DealActionUseCase
	// validator           validator.DealTransactionValidator
	// transformer         transformer.DealTransactionTransformer
	transactionClient   *client.TransactionClient
	bdsproClient        *client.BdsproClient
	userClient          *client.UserClient
	dealTransactionRepo repository.DealTransactionRepository
}

func NewContractDealHandler(
	contractDealUsecase usecase.DealTransactionUsecase,
	dealActionUsecase *usecase.DealActionUseCase,
	// validator validator.DealTransactionValidator,
	// transformer transformer.DealTransactionTransformer,
	transactionClient *client.TransactionClient,
	bdsproClient *client.BdsproClient,
	userClient *client.UserClient,
	dealTransactionRepo repository.DealTransactionRepository,
) *ContractDealHandler {
	return &ContractDealHandler{
		contractDealUsecase: contractDealUsecase,
		dealActionUsecase:   dealActionUsecase,
		// validator:           validator,
		// transformer:         transformer,
		transactionClient:   transactionClient,
		bdsproClient:        bdsproClient,
		userClient:          userClient,
		dealTransactionRepo: dealTransactionRepo,
	}
}

// // @Summary Tạo giao dịch thương vụ
// // @Description Tạo giao dịch thương vụ
// // @Tags Contract Deal
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param dealId path int true "ID của deal"
// // @Param contract body pb.CreateDealContractRequest true "Yêu cầu tạo giao dịch thương vụ"
// // @Success 200 {object} pb.DealTransactionResponse "Giao dịch thương vụ đã tạo"
// // @Failure 400 {object} status.Error "Lỗi yêu cầu không hợp lệ"
// // @Failure 500 {object} status.Error "Lỗi server"
// // @Router /contract/deal/{dealId} [post]
// func (h *ContractDealHandler) CreateTransaction(ctx context.Context, req *pb.CreateDealTransactionRequest) (*pb.DealContractResponse, error) {
// 	// Validate request
// 	if err := h.validator.ValidateCreateRequest(req); err != nil {
// 		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
// 	}

// 	// Transform request
// 	createReq, err := h.transformer.TransformCreateRequest(req)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "failed to transform request: %v", err)
// 	}
// 	// Get user ID from context
// 	userID, err := getUserIDFromContext(ctx)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated: %v", err)
// 	}

// 	// Create transaction
// 	_, err = h.contractDealUsecase.CreateTransaction(ctx, createReq, userID)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "failed to create transaction: %v", err)
// 	}
// 	// Transform response
// 	// pbResponse := h.transformer.TransformToProtoResponse(response)
// 	// return pbResponse, nil
// 	return nil, nil
// }

// // @Summary Cập nhật giao dịch thương vụ
// // @Description Cập nhật giao dịch thương vụ
// // @Tags Contract Deal
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param id path int true "ID của giao dịch"
// // @Param contract body pb.UpdateDealContractRequest true "Yêu cầu cập nhật giao dịch thương vụ"
// // @Success 200 {object} pb.DealTransactionResponse "Giao dịch thương vụ đã cập nhật"
// // @Failure 400 {object} status.Error "Lỗi yêu cầu không hợp lệ"
// // @Failure 500 {object} status.Error "Lỗi server"
// // @Router /contract/deal/{id} [put]
// func (h *ContractDealHandler) UpdateTransaction(ctx context.Context, req *pb.UpdateDealTransactionRequest) (*pb.DealContractResponse, error) {
// 	// Validate request
// 	if err := h.validator.ValidateUpdateRequest(req); err != nil {
// 		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
// 	}

// 	// Transform request
// 	updateReq, err := h.transformer.TransformUpdateRequest(req)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "failed to transform request: %v", err)
// 	}

// 	// Get user ID from context
// 	userID, err := getUserIDFromContext(ctx)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated: %v", err)
// 	}

// 	// Update transaction
// 	_, err = h.contractDealUsecase.UpdateTransaction(ctx, updateReq, userID)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "failed to update transaction: %v", err)
// 	}

// 	// Transform response
// 	// pbResponse := h.transformer.TransformToProtoResponse(response)
// 	// return pbResponse, nil
// 	return nil, nil
// }

// // @Summary Xóa giao dịch thương vụ
// // @Description Xóa giao dịch thương vụ
// // @Tags Contract Deal
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param id path int true "ID của giao dịch"
// // @Success 200 {object} emptypb.Empty "Giao dịch thương vụ đã xóa"
// // @Failure 400 {object} status.Error "Lỗi yêu cầu không hợp lệ"
// // @Failure 500 {object} status.Error "Lỗi server"
// // @Router /contract/deal/{id} [delete]
// func (h *ContractDealHandler) DeleteContract(ctx context.Context, req *pb.DeleteDealContractRequest) (*emptypb.Empty, error) {
// 	// Validate request
// 	if req.Id == 0 {
// 		return nil, status.Errorf(codes.InvalidArgument, "transaction ID is required")
// 	}

// 	// Get user ID from context
// 	userID, err := getUserIDFromContext(ctx)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated: %v", err)
// 	}

// 	// Delete transaction
// 	err = h.contractDealUsecase.DeleteTransaction(ctx, req.Id, userID)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "failed to delete transaction: %v", err)
// 	}

// 	return &emptypb.Empty{}, nil
// }

// // @Summary Lấy chi tiết giao dịch thương vụ
// // @Description Lấy chi tiết giao dịch thương vụ
// // @Tags Contract Deal
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param id path int true "ID của giao dịch"
// // @Success 200 {object} pb.DealContractResponse "Giao dịch thương vụ đã lấy"
// // @Failure 400 {object} status.Error "Lỗi yêu cầu không hợp lệ"
// // @Failure 500 {object} status.Error "Lỗi server"
// // @Router /contract/deal/detail/{id} [get]
// func (h *ContractDealHandler) GetContract(ctx context.Context, req *sharepb.IdRequest) (*pb.DealContractResponse, error) {
// 	// Validate request
// 	if req.Id == 0 {
// 		return nil, status.Errorf(codes.InvalidArgument, "transaction ID is required")
// 	}

// 	transactionPb, _ := h.transactionClient.GetTransactionDetail(ctx, &sharepb.IdRequest{
// 		Id: req.Id,
// 	})
// 	contract, _ := h.contractDealUsecase.GetContractByID(ctx, req.Id)
// 	dealContractPb := h.transformer.DealContractToPb(contract)
// 	dealContractPb.Transaction = transactionPb

// 	h.bdsproClient.MapDealContractToPb(ctx, dealContractPb)
// 	h.userClient.MapProfileToDealContractPb(ctx, dealContractPb)
// 	return dealContractPb, nil
// }

// // @Summary Lấy danh sách giao dịch thương vụ
// // @Description Lấy danh sách giao dịch thương vụ
// // @Tags Contract Deal
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param dealId path int true "ID của deal"
// // @Param type query int false "Loại giao dịch"
// // @Param page query int false "Trang"
// // @Param size query int false "Kích thước trang"
// // @Success 200 {object} pb.DealTransactionListResponse "Danh sách giao dịch thương vụ"
// // @Failure 400 {object} status.Error "Lỗi yêu cầu không hợp lệ"
// // @Failure 500 {object} status.Error "Lỗi server"
// // @Router /contract/deal/{dealId} [get]
// func (h *ContractDealHandler) GetTransactions(ctx context.Context, req *pb.GetDealTransactionsRequest) (*pb.DealContractListResponse, error) {
// 	// Validate request
// 	if req.DealId == 0 {
// 		return nil, status.Errorf(codes.InvalidArgument, "deal ID is required")
// 	}
// 	page := 0
// 	size := 10
// 	if req.Page > 0 {
// 		page = int(req.Page)
// 	}
// 	if req.Size > 0 {
// 		size = int(req.Size)
// 	}

// 	// Transform request
// 	getReq := &dto.GetDealTransactionsRequest{
// 		DealID: req.DealId,
// 		Page:   page,
// 		Size:   size,
// 	}

// 	// Set transaction type filter if provided
// 	if req.Type != nil {
// 		transactionType := h.transformer.TransformTransactionType(*req.Type)
// 		getReq.Type = &transactionType
// 	}

// 	// Get user ID from context
// 	userID, err := getUserIDFromContext(ctx)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated: %v", err)
// 	}

// 	// Get transactions
// 	_, err = h.contractDealUsecase.GetTransactions(ctx, getReq, userID)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "failed to get transactions: %v", err)
// 	}

// 	// Transform response
// 	// pbResponse := h.transformer.TransformToProtoListResponse(response)
// 	return nil, nil
// }

// // @Summary Phê duyệt giao dịch thương vụ
// // @Description Phê duyệt giao dịch thương vụ
// // @Tags Contract Deal
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param id path int true "ID của giao dịch"
// // @Param rejectReason query string false "Lý do từ chối"
// // @Success 200 {object} emptypb.Empty "Giao dịch thương vụ đã phê duyệt"
// // @Failure 400 {object} status.Error "Lỗi yêu cầu không hợp lệ"
// // @Failure 500 {object} status.Error "Lỗi server"
// // @Router /contract/deal/{id}/approve [put]
// func (h *ContractDealHandler) ApproveTransaction(ctx context.Context, req *pb.ApproveDealTransactionRequest) (*emptypb.Empty, error) {
// 	// Validate request
// 	if req.Id == 0 {
// 		return nil, status.Errorf(codes.InvalidArgument, "transaction ID is required")
// 	}

// 	// Transform request
// 	approveReq := &dto.ApproveDealTransactionRequest{
// 		ID:           req.Id,
// 		RejectReason: req.RejectReason,
// 	}

// 	// Get user ID from context
// 	userID, err := getUserIDFromContext(ctx)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated: %v", err)
// 	}

// 	// Approve transaction
// 	err = h.contractDealUsecase.ApproveTransaction(ctx, approveReq, userID)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "failed to approve transaction: %v", err)
// 	}

// 	return &emptypb.Empty{}, nil
// }

// // @Summary Lấy thống kê tài chính
// // @Description Lấy thống kê tài chính
// // @Tags Contract Deal
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param dealId path int true "ID của deal"
// // @Success 200 {object} pb.FinancialSummaryResponse "Thống kê tài chính"
// // @Failure 400 {object} status.Error "Lỗi yêu cầu không hợp lệ"
// // @Failure 500 {object} status.Error "Lỗi server"
// // @Router /contract/deal/{dealId}/financial-summary [get]
// func (h *ContractDealHandler) GetFinancialSummary(ctx context.Context, req *pb.GetFinancialSummaryRequest) (*pb.FinancialSummaryResponse, error) {
// 	// Validate request
// 	if req.DealId == 0 {
// 		return nil, status.Errorf(codes.InvalidArgument, "deal ID is required")
// 	}

// 	// Get user ID from context
// 	userID, err := getUserIDFromContext(ctx)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated: %v", err)
// 	}

// 	// Get financial summary
// 	summary, err := h.contractDealUsecase.GetFinancialSummary(ctx, req.DealId, userID)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "failed to get financial summary: %v", err)
// 	}

// 	// Transform response
// 	pbResponse := h.transformer.TransformToProtoFinancialSummary(summary)
// 	return pbResponse, nil
// }

// // Helper function to get user ID from context
// func getUserIDFromContext(ctx context.Context) (uint64, error) {
// 	userID := _utils.GetProfileIdWithContext(ctx)
// 	return userID, nil
// }

// // @Summary Ký hợp đồng
// // @Description Ký hợp đồng
// // @Tags Contract Deal
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param signContract body pb.SignContractRequest true "Yêu cầu ký hợp đồng"
// // @Success 200 {object} pb.DealContractResponse "Giao dịch thương vụ đã ký"
// // @Failure 400 {object} status.Error "Lỗi yêu cầu không hợp lệ"
// // @Failure 500 {object} status.Error "Lỗi server"
// // @Router /contract/deal/sign [post]
// func (h *ContractDealHandler) SignContract(ctx context.Context, req *pb.SignContractRequest) (*pb.DealContractResp, error) {
// 	// // Validate request
// 	// if req.DealId == 0 || req.SignedBy == 0 {
// 	// 	return nil, status.Errorf(codes.InvalidArgument, "deal ID and signed by are required")
// 	// }

// 	// // Get user ID from context
// 	// // userID, err := getUserIDFromContext(ctx)
// 	// // if err != nil {
// 	// // 	return nil, status.Errorf(codes.Unauthenticated, "user not authenticated: %v", err)
// 	// // }

// 	// // Sign contract
// 	// signReq := &transactionpb.CreateActionRequest{
// 	// 	TransactionId: req.DealId,
// 	// 	FromId:        req.SignedBy,
// 	// 	FromOf:        transactionpb.OwnerType_OWNER_TYPE_USER,
// 	// 	ToId:          req.DealId,
// 	// 	ToOf:          transactionpb.OwnerType_OWNER_TYPE_DEAL,
// 	// 	Action:        transactionpb.Action_ACTION_SIGN,
// 	// 	Value:         req.Amount,
// 	// 	Note:          req.Note,
// 	// }
// 	// if req.ContractName != "" {
// 	// 	signReq.ContractName = req.ContractName
// 	// }
// 	// response, err := h.transactionClient.CreateAction(ctx, signReq)
// 	// if err != nil {
// 	// 	return nil, status.Errorf(codes.Internal, "failed to sign contract: %v", err)
// 	// }

// 	// // Transform response
// 	// return &pb.DealContractResp{
// 	// 	TransactionId: response.TransactionId,
// 	// }, nil
// 	return nil, nil
// }

// // // @Summary Mua/bán/thuê
// // // @Description Mua/bán/thuê
// // // @Tags Contract Deal
// // // @Accept json
// // // @Produce json
// // // @Security BearerAuth
// // // @Param request body pb.TransferContractRequest true "Yêu cầu mua/bán/thuê"
// // // @Success 200 {object} pb.DealContractResp "Giao dịch thương vụ đã mua/bán/thuê"
// // // @Failure 400 {object} status.Error "Lỗi yêu cầu không hợp lệ"
// // // @Failure 500 {object} status.Error "Lỗi server"
// // // @Router /contract/deal/transfer [post]
// // func (h *ContractDealHandler) TransferContract(ctx context.Context, req *pb.TransferContractRequest) (*pb.DealContractResp, error) {
// // 	// Validate request
// // 	if req.DealId == 0 || req.CustomerId == 0 || req.ProductId == 0 {
// // 		return nil, status.Errorf(codes.InvalidArgument, "deal ID, customer ID, and product IDs are required")
// // 	}

// // 	if req.Amount <= 0 {
// // 		return nil, status.Errorf(codes.InvalidArgument, "amount must be greater than 0")
// // 	}

// // 	// Call usecase to handle transfer contract logic
// // 	response, err := h.dealActionUsecase.CreateContractTransfer(ctx, req)
// // 	if err != nil {
// // 		return nil, status.Errorf(codes.Internal, "failed to transfer contract: %v", err)
// // 	}

// // 	return response, nil
// // }

// // @Summary Lấy danh sách giao dịch thương vụ
// // @Description Lấy danh sách giao dịch thương vụ
// // @Tags Contract Deal
// // @Accept json
// // @Produce json
// // @Security BearerAuth
// // @Param id path int true "ID của deal"
// // @Router /contract/deal/{id}/transactions [get]
// func (h *ContractDealHandler) GetDealContract(ctx context.Context, req *sharepb.IdRequest) (*pb.DealContractListResponse, error) {
// 	dealID := req.Id
// 	dealSearch := dto.DealSearch{
// 		Pagable: _dto.Pagable{
// 			Page: 0,
// 			Size: 10,
// 		},
// 	}
// 	dealTransactions, total, err := h.dealActionUsecase.GetDealContractByDealID(ctx, dealID, dealSearch)
// 	if err != nil {
// 		return nil, status.Errorf(codes.Internal, "failed to get deal contract: %v", err)
// 	}
// 	pbResponse := h.transformer.DealContractToPbs(dealTransactions)
// 	h.transactionClient.MapTransactionToDealContract(ctx, pbResponse)
// 	h.userClient.MapProfileToDealContractPbList(ctx, pbResponse)

// 	return &pb.DealContractListResponse{
// 		Data:  pbResponse,
// 		Total: int32(total),
// 	}, nil
// }

// func (h *ContractDealHandler) GetDealTransactionsByIds(ctx context.Context, req *sharepb.IdRequest) (*transactionpb.GetDealTransactionsByIdsResponse, error) {
// 	dealTransactions, err := h.dealTransactionRepo.GetByIds(ctx, req.Ids)
// 	if err != nil {
// 		return nil, err
// 	}
// 	pbDealTransactions := h.transformer.DealContractToPbs(dealTransactions)
// 	return &transactionpb.GetDealTransactionsByIdsResponse{
// 		DealTransactions: pbDealTransactions,
// 	}, nil
// }
