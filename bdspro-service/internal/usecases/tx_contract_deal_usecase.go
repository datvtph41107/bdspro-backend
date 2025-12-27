package usecases

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"bdspro/internal/provider"
	"bdspro/internal/repo"
	_utils "common/utils"
	"context"
	"errors"
	"time"
)

type DealContractUsecase struct {
	dealTransactionRepo repo.TxContractDealRepository
	actionRepo          repo.TxActionRepository
	transaction         provider.TransactionProvider
	transactionUsecase  TxTransactionUsecase
	// orgClient           provider.IOrganizationClient
	dealUsecase  DealUsecase
	dealCostRepo repo.TxDealCostRepository
}

// NewDealContractUsecase creates a new instance of DealContractUsecase
func NewDealContractUsecase(
	dealTransactionRepo repo.TxContractDealRepository,
	actionRepo repo.TxActionRepository,
	transaction provider.TransactionProvider,
	transactionUsecase TxTransactionUsecase,
	// orgClient provider.IOrganizationClient,
	dealUsecase DealUsecase,
	dealCostRepo repo.TxDealCostRepository,
) *DealContractUsecase {
	return &DealContractUsecase{
		dealTransactionRepo: dealTransactionRepo,
		actionRepo:          actionRepo,
		transaction:         transaction,
		transactionUsecase:  transactionUsecase,
		// orgClient:           orgClient,
		dealUsecase:  dealUsecase,
		dealCostRepo: dealCostRepo,
	}
}

func (uc *DealContractUsecase) GetDealContractProcess(ctx context.Context, dealID uint64) ([]dto.TxProcessDTO, error) {
	processes, err := uc.dealTransactionRepo.GetDealContractProcess(ctx, dealID)
	if err != nil {
		return nil, err
	}
	for i, process := range processes {
		action := enums.TxAction(process.Action)
		style := enums.GetTxActionStyle(action)
		processes[i].StatusName = enums.TxActionNames[action]
		processes[i].Color = style.Color
		processes[i].BorderColor = style.BorderColor
		processes[i].BgColor = style.BgColor
	}
	return processes, nil
}

func (uc *DealContractUsecase) GetDealContractByDealID(ctx context.Context, dealID uint64, req dto.TxDealSearch) ([]*dto.TxContractDealListResponse, int64, error) {
	dealTransaction, total, err := uc.dealTransactionRepo.GetListByDealID(ctx, dealID, req)
	if err != nil {
		return nil, 0, err
	}
	return dealTransaction, total, nil
}

func (uc *DealContractUsecase) GetDealContractByProfileID(ctx context.Context, profileID uint64, req dto.TxDealSearch) ([]*dto.TxContractDealListResponse, int64, error) {
	dealTransaction, total, err := uc.dealTransactionRepo.GetListByProfileID(ctx, profileID, req)
	if err != nil {
		return nil, 0, err
	}
	return dealTransaction, total, nil
}

func (uc *DealContractUsecase) GetDealContractByTransactionIds(ctx context.Context, transactionIds []uint64) ([]*dto.TxContractDealListResponse, error) {
	dealTransactions, err := uc.dealTransactionRepo.GetListByTransactionIds(ctx, transactionIds)
	if err != nil {
		return nil, err
	}
	return dealTransactions, nil
}

// todo: map lại enum sau
func (uc *DealContractUsecase) GetActionWhenSave(ctx context.Context, req *dto.TxDealContractPaymentRequest) (*domain.TxAction, error) {
	action := &domain.TxAction{
		TxId:      req.Id,
		FromId:    req.CustomerId,
		FromOf:    enums.TxOwnerType(enums.TxOwnerUser),
		ToId:      req.DealId,
		ToOf:      enums.TxOwnerDealContract,
		Action:    enums.TxAction(req.TransactionStep),
		Value:     req.Amount,
		Note:      req.Note,
		Timestamp: req.Timestamp,
	}
	if req.TransactionStep == 10 {
		action.Action = 10
		action.ToId = req.DealId
		action.ToOf = enums.TxOwnerDealContract
		return action, nil
	} else if req.TransactionStep == 20 {
		action.Action = 20
		action.ToId = req.DealId
		action.ToOf = enums.TxOwnerDealContract
		return action, nil
	} else if req.TransactionStep == 30 {
		action.Action = 30
		action.ToId = req.DealId
		action.ToOf = enums.TxOwnerDealContract
		return action, nil
	} else if req.TransactionStep == 40 {
		action.Action = 40
		action.ToId = req.DealId
		action.ToOf = enums.TxOwnerDealContract
		return action, nil
	} else if req.TransactionStep == 50 {
		action.Action = 50
		action.ToId = req.DealId
		action.ToOf = enums.TxOwnerDealContract
		return action, nil
	} else if req.TransactionStep == 60 {
		action.Action = 60
		action.ToId = req.ProductId
	}
	return action, nil
}

func (uc *DealContractUsecase) CreateTransactionBuy(ctx context.Context, req *dto.TxDealContractPaymentRequest) error {
	// currentUserId := _utils.GetProfileIdWithContext(ctx)
	_, err := uc.dealUsecase.SaveDealProduct(ctx, req.DealId, []uint64{req.ProductId})
	if err != nil {
		return err
	}
	return nil
}

// TransferContract handles the transfer contract operation
func (uc *DealContractUsecase) CreateContractPayment(ctx context.Context, req *dto.TxDealContractPaymentRequest) (*domain.TxContractDeal, error) {
	// Get user info from context
	currentUserId := _utils.GetProfileIdWithContext(ctx)

	contract := &domain.TxContractDeal{
		DealID:     req.DealId,
		Type:       enums.TxTransactionType(req.TransactionType),
		CustomerID: req.CustomerId,
		Note:       req.Note,
		ProductID:  &req.ProductId,
	}

	err := uc.transaction.WithTransaction(ctx, func(ctx context.Context) error {
		err := uc.dealTransactionRepo.Create(ctx, contract)
		if err != nil {
			return err
		}

		// uc.transactionUsecase.CreateTransaction({})
		actionReq, err := uc.GetActionWhenSave(ctx, req)
		if err != nil {
			return err
		}
		// Create transaction in transaction service using Transfer method
		transactionReq := &domain.Tx{
			FromID:     currentUserId,
			FromOf:     enums.TxOwnerUser,
			ToID:       req.DealId,
			ToOf:       enums.TxOwnerDealContract,
			Value:      req.Amount,
			Note:       req.Note,
			Method:     enums.TxMethod(req.TransactionType),
			Timestamp:  time.Now(),
			LastAction: actionReq,
		}

		transactionResp, err := uc.transactionUsecase.CreateTransaction(ctx, transactionReq)
		if err != nil {
			return err
		}

		contract.TransactionID = &transactionResp.ID
		err = uc.dealTransactionRepo.Update(ctx, contract)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return contract, nil
}

func (uc *DealContractUsecase) GetContractDealByTransactionID(ctx context.Context, id uint64) (*domain.TxContractDeal, error) {
	contract, err := uc.dealTransactionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	txt, err := uc.transactionUsecase.GetTransaction(ctx, *contract.TransactionID)
	if err != nil {
		return nil, err
	}
	contract.Transaction = txt

	return contract, nil
}

func (uc *DealContractUsecase) GetMessageByAction(action enums.TxAction) string {
	switch action {
	case enums.TxActionDeposite:
		return "Đặt cọc"
	case enums.TxActionSign:
		return "Ký hợp đồng"
	case enums.TxActionPay:
		return "Thanh toán"
	case enums.TxActionHandOver:
		return "Bàn giao"
	case enums.TxActionCancel:
		return "Hủy"
	}
	return ""
}

func (uc *DealContractUsecase) CreateActionPayment(ctx context.Context, req *dto.TxActionPaymentRequest) (*domain.TxAction, error) {
	contract, err := uc.dealTransactionRepo.GetByID(ctx, req.ContractId)
	if err != nil {
		return nil, err
	}

	if contract.TransactionID == nil {
		return nil, errors.New("transaction not found")
	}

	transaction, err := uc.transactionUsecase.GetTransaction(ctx, *contract.TransactionID)
	if err != nil {
		return nil, err
	}

	action := &domain.TxAction{
		TxId:      *contract.TransactionID,
		FromId:    req.FromId,
		FromOf:    enums.TxOwnerUser,
		ToId:      transaction.ToID,
		ToOf:      transaction.ToOf,
		Action:    enums.TxAction(req.Action),
		Value:     req.Amount,
		Note:      req.Note,
		Timestamp: req.Timestamp,
	}

	action.Message = uc.transactionUsecase.GetMessageByAction(action)

	err = uc.transaction.WithTransaction(ctx, func(ctx context.Context) error {
		action, err = uc.transactionUsecase.CreateAction(ctx, action)
		if err != nil {
			return err
		}

		if action.Action == enums.TxActionHandOver {
			_, err = uc.transactionUsecase.ApproveTransaction(ctx, *contract.TransactionID)
			if err != nil {
				return err
			}
			if action.Action == enums.TxActionHandOver && transaction.Method == enums.TxMethodBuy {
				uc.CreateTransactionBuy(ctx, &dto.TxDealContractPaymentRequest{
					DealId:          contract.DealID,
					ProductId:       *contract.ProductID,
					TransactionType: uint32(enums.TxMethodBuy),
				})
			}
		} else if action.Action == enums.TxActionCancel {
			err = uc.transactionUsecase.RejectTransaction(ctx, *contract.TransactionID)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return action, nil
}

func (uc *DealContractUsecase) DeleteContract(ctx context.Context, contractID uint64) error {
	dealTransaction, err := uc.dealTransactionRepo.GetByID(ctx, contractID)
	if err != nil {
		return err
	}

	if dealTransaction.TransactionID == nil {
		return errors.New("không tìm thấy giao dịch")
	}
	// transaction, err := uc.transactionUsecase.GetTransaction(ctx, *dealTransaction.TransactionID)
	// if err != nil {
	// 	return err
	// }

	// if transaction.Status == enums.TxStatusDone {
	// 	return errors.New("giao dịch đã hoàn thành không thể xóa")
	// }

	uc.transaction.WithTransaction(ctx, func(ctx context.Context) error {
		err = uc.dealTransactionRepo.Delete(ctx, contractID)
		if err != nil {
			return err
		}

		err = uc.transactionUsecase.DeleteTransaction(ctx, *dealTransaction.TransactionID)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// GetStatisticsByDealID returns comprehensive statistics for a specific deal
func (uc *DealContractUsecase) GetStatisticsByDealID(ctx context.Context, dealID uint64, includePayment, includeCost bool) (*dto.TxDealCostStatisticsDTO, error) {
	result := &dto.TxDealCostStatisticsDTO{
		DealID: dealID,
	}

	// Get cost statistics if requested
	if includeCost {
		costStats, err := uc.dealCostRepo.GetStatisticsByDealID(ctx, dealID)
		if err != nil {
			return nil, err
		}
		result.AmountCost = costStats.AmountCost
		result.AmountRevenue = costStats.AmountRevenue
		result.NumCost = costStats.NumCost
		result.NumRevenue = costStats.NumRevenue
	}

	// Get payment statistics if requested
	if includePayment {
		paymentStats, err := uc.dealTransactionRepo.GetTotalAmountByDealID(ctx, dealID)
		if err != nil {
			return nil, err
		}
		result.Payment = paymentStats.Payment
		result.NumPayment = paymentStats.NumPayment
	}

	if includePayment && includeCost {
		result.TotalAmount = result.Payment.Approved + result.AmountRevenue.Approved
		result.TotalCost = result.AmountCost.Approved
		result.TotalProfit = result.TotalAmount - result.TotalCost
	}

	return result, nil
}
