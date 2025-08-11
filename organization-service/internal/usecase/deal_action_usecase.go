package usecase

import (
	_utils "common/utils"
	"context"
	"organization/infrastructure/client"
	"organization/infrastructure/transformer"
	"organization/internal/domain/entity"
	"organization/internal/dto"
	"organization/internal/repo"
	pb "pb/types/organization"
	transactionpb "pb/types/transaction"
)

type DealActionUseCase struct {
	dealActionRepo    repo.DealActionRepository
	transactionClient *client.TransactionClient
}

// NewDealActionUseCase creates a new instance of DealActionUseCase
func NewDealActionUseCase(
	dealActionRepo repo.DealActionRepository,
	transactionClient *client.TransactionClient,
) *DealActionUseCase {
	return &DealActionUseCase{
		dealActionRepo:    dealActionRepo,
		transactionClient: transactionClient,
	}
}

// TransferContract handles the transfer contract operation
func (uc *DealActionUseCase) TransferContract(ctx context.Context, req *pb.TransferContractRequest) (*pb.DealContractResp, error) {
	// Get user info from context
	organizationID := _utils.GetOrganizationIdFromContext(ctx)

	// Create transaction in transaction service using Transfer method
	transactionReq := &transactionpb.CreateProductTransactionRequest{
		FromId: req.CustomerId,
		ToId:   req.DealId,
		// ToOf:     transactionpb.OwnerType_OWNER_TYPE_DEAL,
		Amount:   req.Amount,
		Currency: "VND", // Default currency
		Note:     req.Note,
	}

	transactionResp, err := uc.transactionClient.CreateProductTransaction(ctx, transactionReq)
	if err != nil {
		return nil, err
	}

	// Create deal action record
	dealAction := &entity.DealAction{
		DealID:          req.DealId,
		TransactionID:   transactionResp.TransactionId,
		CustomerID:      req.CustomerId,
		ProductID:       req.ProductId,
		Amount:          req.Amount,
		TransactionStep: req.TransactionStep,
		TransactionType: req.TransactionType,
		Note:            &req.Note,
		Timestamp:       _utils.ParseStringToTime(req.Timestamp),
		OrganizationID:  organizationID,
	}

	if err := uc.dealActionRepo.Create(ctx, dealAction); err != nil {
		return nil, err
	}

	return &pb.DealContractResp{
		TransactionId: transactionResp.TransactionId,
	}, nil
}

// GetDealActions gets deal actions with pagination
func (uc *DealActionUseCase) GetDealActions(ctx context.Context, req *dto.GetDealActionsRequest) (*dto.GetDealActionsResponse, error) {
	dealActions, total, err := uc.dealActionRepo.GetByDealID(ctx, *req.DealID, int32(req.Page), int32(req.Size))
	if err != nil {
		return nil, err
	}

	responses := make([]dto.DealActionResponse, len(dealActions))
	for i, dealAction := range dealActions {
		responses[i] = *transformer.ToDealActionResponse(dealAction)
	}

	return &dto.GetDealActionsResponse{
		Data:  responses,
		Total: total,
	}, nil
}

// GetDealActionByID gets a deal action by ID
func (uc *DealActionUseCase) GetDealActionByID(ctx context.Context, id uint64) (*dto.DealActionResponse, error) {
	dealAction, err := uc.dealActionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return transformer.ToDealActionResponse(dealAction), nil
}
