package usecase

import (
	"context"
	"organization/infrastructure/client"
	"organization/internal/domain/entity"
	"organization/internal/domain/repository"
	"organization/internal/dto"
	iusecase "organization/internal/interface"
	transactionpb "pb/types/transaction"
)

type DealActionUseCase struct {
	dealTransactionRepo repository.DealTransactionRepository
	transactionClient   *client.TransactionClient
	transaction         iusecase.ITransaction
}

// NewDealActionUseCase creates a new instance of DealActionUseCase
func NewDealActionUseCase(
	dealTransactionRepo repository.DealTransactionRepository,
	transactionClient *client.TransactionClient,
	transaction iusecase.ITransaction,
) *DealActionUseCase {
	return &DealActionUseCase{
		dealTransactionRepo: dealTransactionRepo,
		transactionClient:   transactionClient,
		transaction:         transaction,
	}
}

func (uc *DealActionUseCase) GetDealContractByDealID(ctx context.Context, dealID uint64, req dto.DealSearchRequest) ([]*entity.DealTransaction, int64, error) {
	dealTransaction, total, err := uc.dealTransactionRepo.GetListByDealID(ctx, dealID, req)
	if err != nil {
		return nil, 0, err
	}
	return dealTransaction, total, nil
}

// todo: map lại enum sau
func (uc *DealActionUseCase) GetActionWhenSave(ctx context.Context, req *transactionpb.DealContractPaymentRequest) (*transactionpb.CreateActionRequest, error) {
	action := &transactionpb.CreateActionRequest{
		FromId:   req.CustomerId,
		FromOf:   transactionpb.OwnerType_OWNER_TYPE_USER,
		Amount:   req.Amount,
		Currency: "VND",
		Note:     req.Note,
	}
	if req.TransactionStep == 10 {
		action.Action = 10
		action.ToId = req.DealId
		action.ToOf = transactionpb.OwnerType_OWNER_TYPE_DEAL
		return action, nil
	} else if req.TransactionStep == 20 {
		action.Action = 20
		action.ToId = req.DealId
		action.ToOf = transactionpb.OwnerType_OWNER_TYPE_DEAL
		return action, nil
	} else if req.TransactionStep == 30 {
		action.Action = 30
		action.ToId = req.DealId
		action.ToOf = transactionpb.OwnerType_OWNER_TYPE_DEAL
		return action, nil
	} else if req.TransactionStep == 40 {
		action.Action = 40
		action.ToId = req.DealId
		action.ToOf = transactionpb.OwnerType_OWNER_TYPE_DEAL
		return action, nil
	} else if req.TransactionStep == 50 {
		action.Action = 50
		action.ToId = req.DealId
		action.ToOf = transactionpb.OwnerType_OWNER_TYPE_DEAL
		return action, nil
	} else if req.TransactionStep == 60 {
		action.Action = 60
		action.ToId = req.ProductId
	}
	return nil, nil
}

// // TransferContract handles the transfer contract operation
// func (uc *DealActionUseCase) CreateContractTransfer(ctx context.Context, req *transactionpb.DealContractPaymentRequest) (*transactionpb.DealContractResp, error) {
// 	// Get user info from context
// 	organizationID := _utils.GetOrganizationIdFromContext(ctx)

// 	contract := &entity.DealTransaction{
// 		DealID:          req.DealId,
// 		Type:            enums.TransactionType(req.TransactionType),
// 		Amount:          req.Amount,
// 		Description:     req.Note,
// 		TransactionDate: _utils.ParseStringToTime(req.Timestamp),
// 		CustomerID:      req.CustomerId,
// 		Note:            req.Note,
// 		IsManual:        true,
// 		CreatedBy:       organizationID,
// 		Status:          enums.ApprovedStatusPending,
// 	}

// 	err := uc.transaction.WithTransaction(ctx, func(ctx context.Context) error {
// 		err := uc.dealTransactionRepo.Create(ctx, contract)
// 		if err != nil {
// 			return err
// 		}

// 		actionReq, err := uc.GetActionWhenSave(ctx, req)
// 		if err != nil {
// 			return err
// 		}
// 		// Create transaction in transaction service using Transfer method
// 		transactionReq := &transactionpb.CreateProductTransactionRequest{
// 			FromId:        req.CustomerId,
// 			FromOf:        transactionpb.OwnerType_OWNER_TYPE_USER,
// 			ToId:          contract.ID,
// 			ToOf:          transactionpb.OwnerType_OWNER_TYPE_DEAL,
// 			ActionRequest: actionReq,
// 		}

// 		transactionResp, err := uc.transactionClient.CreateProductTransaction(ctx, transactionReq)
// 		if err != nil {
// 			return err
// 		}

// 		contract.TransactionID = transactionResp.TransactionId
// 		err = uc.dealTransactionRepo.Update(ctx, contract)
// 		if err != nil {
// 			return err
// 		}

// 		return nil
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &transactionpb.DealContractResp{
// 		TransactionId: contract.TransactionID,
// 		ContractId:    contract.ID,
// 	}, nil
// }

// GetDealTransactionByID gets a deal transaction by ID
func (uc *DealActionUseCase) GetDealTransactionByID(ctx context.Context, id uint64) (*dto.DealTransactionResponse, error) {
	dealTransaction, err := uc.dealTransactionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dto.DealTransactionResponse{
		ID:              dealTransaction.ID,
		DealID:          dealTransaction.DealID,
		Type:            dealTransaction.Type,
		Amount:          dealTransaction.Amount,
		Description:     dealTransaction.Description,
		TransactionDate: dealTransaction.TransactionDate,
		Category:        dealTransaction.Category,
		Note:            dealTransaction.Note,
		IsManual:        dealTransaction.IsManual,
		CreatedBy:       dealTransaction.CreatedBy,
		Status:          dealTransaction.Status,
		ApprovedBy:      dealTransaction.ApprovedBy,
		ApprovedAt:      dealTransaction.ApprovedAt,
		RejectReason:    dealTransaction.RejectReason,
		// TransactionID:   dealTransaction.TransactionID,
	}, nil
}
