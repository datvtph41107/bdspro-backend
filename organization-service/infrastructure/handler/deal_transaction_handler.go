package handler

import (
	"context"
	"organization/infrastructure/transformer"
	"organization/infrastructure/validator"
	"organization/internal/dto"
	"organization/internal/usecase"
	pb "pb/types/organization"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DealTransactionHandler struct {
	pb.UnimplementedDealTransactionServiceServer
	dealTransactionUsecase usecase.DealTransactionUsecase
	validator              validator.DealTransactionValidator
	transformer            transformer.DealTransactionTransformer
}

func NewDealTransactionHandler(
	dealTransactionUsecase usecase.DealTransactionUsecase,
	validator validator.DealTransactionValidator,
	transformer transformer.DealTransactionTransformer,
) *DealTransactionHandler {
	return &DealTransactionHandler{
		dealTransactionUsecase: dealTransactionUsecase,
		validator:              validator,
		transformer:            transformer,
	}
}

func (h *DealTransactionHandler) CreateTransaction(ctx context.Context, req *pb.CreateDealTransactionRequest) (*pb.DealTransactionResponse, error) {
	// Validate request
	if err := h.validator.ValidateCreateRequest(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	// Transform request
	createReq, err := h.transformer.TransformCreateRequest(req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to transform request: %v", err)
	}

	// Get user ID from context
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated: %v", err)
	}

	// Create transaction
	response, err := h.dealTransactionUsecase.CreateTransaction(ctx, createReq, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create transaction: %v", err)
	}

	// Transform response
	pbResponse := h.transformer.TransformToProtoResponse(response)
	return pbResponse, nil
}

func (h *DealTransactionHandler) UpdateTransaction(ctx context.Context, req *pb.UpdateDealTransactionRequest) (*pb.DealTransactionResponse, error) {
	// Validate request
	if err := h.validator.ValidateUpdateRequest(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid request: %v", err)
	}

	// Transform request
	updateReq, err := h.transformer.TransformUpdateRequest(req)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to transform request: %v", err)
	}

	// Get user ID from context
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated: %v", err)
	}

	// Update transaction
	response, err := h.dealTransactionUsecase.UpdateTransaction(ctx, updateReq, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update transaction: %v", err)
	}

	// Transform response
	pbResponse := h.transformer.TransformToProtoResponse(response)
	return pbResponse, nil
}

func (h *DealTransactionHandler) DeleteTransaction(ctx context.Context, req *pb.DeleteDealTransactionRequest) (*pb.DealTransactionResponse, error) {
	// Validate request
	if req.Id == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "transaction ID is required")
	}

	// Get user ID from context
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated: %v", err)
	}

	// Delete transaction
	err = h.dealTransactionUsecase.DeleteTransaction(ctx, req.Id, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete transaction: %v", err)
	}

	return &pb.DealTransactionResponse{
		Id: req.Id,
	}, nil
}

func (h *DealTransactionHandler) GetTransaction(ctx context.Context, req *pb.GetDealTransactionRequest) (*pb.DealTransactionResponse, error) {
	// Validate request
	if req.Id == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "transaction ID is required")
	}

	// Get user ID from context
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated: %v", err)
	}

	// Get transaction
	response, err := h.dealTransactionUsecase.GetTransaction(ctx, req.Id, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get transaction: %v", err)
	}

	// Transform response
	pbResponse := h.transformer.TransformToProtoResponse(response)
	return pbResponse, nil
}

func (h *DealTransactionHandler) GetTransactions(ctx context.Context, req *pb.GetDealTransactionsRequest) (*pb.DealTransactionListResponse, error) {
	// Validate request
	if req.DealId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "deal ID is required")
	}

	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 1 || req.Limit > 100 {
		req.Limit = 20
	}

	// Transform request
	getReq := &dto.GetDealTransactionsRequest{
		DealID: req.DealId,
		Page:   int(req.Page),
		Limit:  int(req.Limit),
	}

	// Set transaction type filter if provided
	if req.Type != nil {
		transactionType := h.transformer.TransformTransactionType(*req.Type)
		getReq.Type = &transactionType
	}

	// Get user ID from context
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated: %v", err)
	}

	// Get transactions
	response, err := h.dealTransactionUsecase.GetTransactions(ctx, getReq, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get transactions: %v", err)
	}

	// Transform response
	pbResponse := h.transformer.TransformToProtoListResponse(response)
	return pbResponse, nil
}

func (h *DealTransactionHandler) ApproveTransaction(ctx context.Context, req *pb.ApproveDealTransactionRequest) (*pb.DealTransactionResponse, error) {
	// Validate request
	if req.Id == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "transaction ID is required")
	}

	// Transform request
	approveReq := &dto.ApproveDealTransactionRequest{
		ID:           req.Id,
		RejectReason: req.RejectReason,
	}

	// Get user ID from context
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated: %v", err)
	}

	// Approve transaction
	err = h.dealTransactionUsecase.ApproveTransaction(ctx, approveReq, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to approve transaction: %v", err)
	}

		return &pb.DealTransactionResponse{
		Id: req.Id,
	}, nil
}

func (h *DealTransactionHandler) GetFinancialSummary(ctx context.Context, req *pb.GetFinancialSummaryRequest) (*pb.FinancialSummaryResponse, error) {
	// Validate request
	if req.DealId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "deal ID is required")
	}

	// Get user ID from context
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "user not authenticated: %v", err)
	}

	// Get financial summary
	summary, err := h.dealTransactionUsecase.GetFinancialSummary(ctx, req.DealId, userID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get financial summary: %v", err)
	}

	// Transform response
	pbResponse := h.transformer.TransformToProtoFinancialSummary(summary)
	return pbResponse, nil
}

// Helper function to get user ID from context
func getUserIDFromContext(ctx context.Context) (uint64, error) {
	// Implementation depends on your authentication middleware
	// This is a placeholder - you need to implement based on your auth system
	userID, ok := ctx.Value("user_id").(uint64)
	if !ok {
		return 0, status.Errorf(codes.Unauthenticated, "user ID not found in context")
	}
	return userID, nil
} 