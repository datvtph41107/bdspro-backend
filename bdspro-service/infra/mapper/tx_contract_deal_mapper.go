package mapper

import (
	"context"
	"time"

	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	_usecase "common/domain/usecase"
	_utils "common/utils"
	bdspropb "pb/types/bdspro"
)

type ContractDealMapper struct {
	transactionMapper *TxMapper
	codeDataUsecase   *_usecase.CodeDataUsecase
}

func NewContractDealMapper(transactionMapper *TxMapper, codeDataUsecase *_usecase.CodeDataUsecase) *ContractDealMapper {
	return &ContractDealMapper{
		transactionMapper: transactionMapper,
		codeDataUsecase:   codeDataUsecase,
	}
}

// func (m *ContractDealMapper) CreateDealTransactionRequestToEntity(req *dto.DealContractPaymentRequest) *domain.ContractDeal {
// 	return &domain.ContractDeal{
// 		DealID: req.DealId,
// 	}
// }

// func (t *ContractDealMapper) DealContractToPb(protoType bdspropb.TransactionType) enums.TransactionType {
// 	switch protoType {
// 	case bdspropb.TransactionType_TRANSACTION_TYPE_EXPENSE:
// 		return enums.TransactionTypeExpense
// 	case bdspropb.TransactionType_TRANSACTION_TYPE_REVENUE:
// 		return enums.TransactionTypeRevenue
// 	}
// }

func (t *ContractDealMapper) DealContractPaymentReqToEntity(req *bdspropb.DealContractPaymentRequest) (*dto.TxDealContractPaymentRequest, error) {
	// Parse transaction date
	// var transactionDate *time.Time
	// if req.Timestamp != "" {
	// 	parsedTime, err := time.Parse(time.RFC3339, req.Timestamp)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	transactionDate = &parsedTime
	// }

	return &dto.TxDealContractPaymentRequest{
		DealId: req.DealId,
		// Type:            t.TransformTransactionType(bdspropb.ContractType(req.TransactionType)),
		Amount: req.Amount,
		// Description:     req.Description,
		// TransactionDate: transactionDate,
		// Category:        req.Category,
		Note: req.Note,
		// DocumentIds:     req.DocumentIds,
		TransactionType: req.TransactionType,
		CustomerId:      req.CustomerId,
		ProductId:       req.ProductId,
		TransactionStep: req.TransactionStep,
	}, nil
}

func (t *ContractDealMapper) TransformUpdateRequest(req *bdspropb.UpdateDealContractRequest) (*dto.UpdateDealTransactionRequest, error) {
	// Parse transaction date
	var transactionDate *time.Time
	if req.ContractDate != "" {
		parsedTime, err := time.Parse(time.RFC3339, req.ContractDate)
		if err != nil {
			return nil, err
		}
		transactionDate = &parsedTime
	}

	return &dto.UpdateDealTransactionRequest{
		ID:              req.Id,
		Amount:          req.Amount,
		Description:     req.Description,
		TransactionDate: transactionDate,
		Category:        req.Category,
		Note:            req.Note,
		DocumentIds:     req.DocumentIds,
	}, nil
}

func (t *ContractDealMapper) DealContractToPb(contract *domain.TxContractDeal) *bdspropb.DealContractResponse {
	result := &bdspropb.DealContractResponse{
		ContractId:      contract.ID,
		DealId:          contract.DealID,
		CustomerId:      contract.CustomerID,
		CreatedAt:       _utils.FormatTimeToString(contract.CreatedAt),
		TransactionType: bdspropb.TransactionType(contract.Type),
		TypeName:        enums.TxTransactionTypeMap[contract.Type],
	}

	if contract.ProductID != nil {
		result.ProductId = *contract.ProductID
	}

	if contract.Transaction != nil {
		txt := t.transactionMapper.TxToPb(contract.Transaction)
		result.Transaction = txt
		result.TransactionName = txt.TransactionName
		result.TransactionId = txt.TransactionId
		result.Amount = txt.Value
	}

	return result
}

func (t *ContractDealMapper) DealContractToPbs(response []*domain.TxContractDeal) []*bdspropb.DealContractResponse {
	transactions := make([]*bdspropb.DealContractResponse, len(response))
	for i, transaction := range response {
		transactions[i] = t.DealContractToPb(transaction)
	}

	return transactions
}

func (t *ContractDealMapper) TransformToProtoFinancialSummary(summary *dto.DealFinancialSummary) *bdspropb.FinancialSummaryResponse {
	return &bdspropb.FinancialSummaryResponse{
		DealId:          summary.DealID,
		TotalRevenue:    summary.TotalRevenue,
		TotalExpense:    summary.TotalExpense,
		NetProfit:       summary.NetProfit,
		ProfitMargin:    summary.ProfitMargin,
		TargetProfit:    summary.TargetProfit,
		ProfitStatus:    summary.ProfitStatus,
		WarningMessages: summary.WarningMessages,
	}
}

func (t *ContractDealMapper) TransformTransactionType(protoType bdspropb.ContractType) enums.TxTransactionType {
	switch protoType {
	case bdspropb.ContractType_CONTRACT_TYPE_EXPENSE:
		return enums.TxTransactionTypeExpense
	case bdspropb.ContractType_CONTRACT_TYPE_REVENUE:
		return enums.TxTransactionTypeRevenue
	default:
		return enums.TxTransactionTypeExpense
	}
}

func (t *ContractDealMapper) DealContractItemToPb(ctx context.Context, contract *dto.TxContractDealListResponse) *bdspropb.DealContractItem {
	result := &bdspropb.DealContractItem{
		Id:         contract.ContractID,
		CustomerId: contract.CustomerID,
		CreatedAt:  _utils.FormatTimeToString(contract.Timestamp),
		ContractId: contract.ContractID,
		Amount:     contract.Amount,
		Status:     uint64(contract.Status),
	}

	if t.codeDataUsecase != nil {
		result.Code = t.codeDataUsecase.GetContractCode(ctx, contract.ContractID)
	}

	if contract.TransactionName == "" {
		result.TransactionName = enums.TxMethodNames[enums.TxMethod(contract.Method)]
		result.Status = uint64(contract.Status)
		result.StatusName = enums.TxApprovedStatusMap[enums.TxApprovedStatus(contract.Status)]
	}

	return result
}

func (t *ContractDealMapper) DealContractItemToPbs(ctx context.Context, contracts []*dto.TxContractDealListResponse) []*bdspropb.DealContractItem {
	items := make([]*bdspropb.DealContractItem, len(contracts))
	for i, contract := range contracts {
		items[i] = t.DealContractItemToPb(ctx, contract)
	}
	return items
}

// DealContractItemToPbsFromEntity maps a slice of ContractDeal entities to protobuf items
func (t *ContractDealMapper) DealContractItemToPbsFromEntity(contracts []*domain.TxContractDeal) []*bdspropb.DealContractItem {
	items := make([]*bdspropb.DealContractItem, len(contracts))
	for i, contract := range contracts {
		items[i] = t.DealContractItemToPbFromEntity(contract)
	}
	return items
}

// DealContractItemToPbFromEntity maps a single ContractDeal entity to protobuf item
func (t *ContractDealMapper) DealContractItemToPbFromEntity(contract *domain.TxContractDeal) *bdspropb.DealContractItem {
	result := &bdspropb.DealContractItem{
		Id:         contract.ID,
		CustomerId: contract.CustomerID,
		CreatedAt:  _utils.FormatTimeToString(contract.CreatedAt),
		ContractId: contract.ID,
	}

	if contract.Transaction != nil && t.transactionMapper != nil {
		txt := t.transactionMapper.TxToPb(contract.Transaction)
		result.TransactionName = txt.TransactionName
		result.Amount = txt.Value
		result.Status = uint64(contract.Transaction.Status)
		result.StatusName = enums.TxApprovedStatusMap[enums.TxApprovedStatus(contract.Transaction.Status)]
		result.ContractId = contract.ID
	}

	return result
}

func (t *ContractDealMapper) PaymentActionReqToDTO(req *bdspropb.PaymentActionRequest) *dto.TxActionPaymentRequest {
	return &dto.TxActionPaymentRequest{
		ID:            req.Id,
		DealId:        req.ToId,
		FromId:        req.FromId,
		Timestamp:     _utils.ParseStringToTime(req.Timestamp),
		Amount:        req.Amount,
		Note:          req.Note,
		TransactionId: req.TransactionId,
		ContractId:    req.ContractId,
		Action:        req.Action,
		// TransactionStep: req.TransactionStep,
		// ProductId:     req.TransactionId,
	}
}

func (t *ContractDealMapper) ActionToPb(action *domain.TxAction) *bdspropb.ActionResponse {
	return &bdspropb.ActionResponse{
		Action:        uint32(action.Action),
		Note:          action.Note,
		FromId:        action.FromId,
		FromOf:        bdspropb.OwnerType(action.FromOf),
		ToId:          action.ToId,
		ToOf:          bdspropb.OwnerType(action.ToOf),
		Value:         action.Value,
		TransactionId: action.TxId,
		CreatedAt:     _utils.FormatTimeToString(action.CreatedAt),
		Message:       action.Message,
	}
}

// StatisticsToPb maps DealComprehensiveStatisticsDTO to protobuf
func (t *ContractDealMapper) StatisticsToPb(statistics *dto.TxDealCostStatisticsDTO) *bdspropb.ContractDealStatisticsDTO {
	result := &bdspropb.ContractDealStatisticsDTO{
		DealId: statistics.DealID,
	}

	if statistics.Payment != nil {
		result.Payment = &bdspropb.StatsDTO{
			Approved: statistics.Payment.Approved,
			Reject:   statistics.Payment.Rejected,
			Pending:  statistics.Payment.Pending,
			Total:    statistics.Payment.Total,
		}
	}

	if statistics.AmountCost != nil {
		result.AmountCost = &bdspropb.StatsDTO{
			Approved: statistics.AmountCost.Approved,
			Reject:   statistics.AmountCost.Rejected,
			Pending:  statistics.AmountCost.Pending,
			Total:    statistics.AmountCost.Total,
		}
	}

	if statistics.AmountRevenue != nil {
		result.AmountRevenue = &bdspropb.StatsDTO{
			Approved: statistics.AmountRevenue.Approved,
			Reject:   statistics.AmountRevenue.Rejected,
			Pending:  statistics.AmountRevenue.Pending,
			Total:    statistics.AmountRevenue.Total,
		}
	}

	if statistics.NumCost != nil {
		result.NumCost = &bdspropb.StatsDTO{
			Approved: statistics.NumCost.Approved,
			Reject:   statistics.NumCost.Rejected,
			Pending:  statistics.NumCost.Pending,
			Total:    statistics.NumCost.Total,
		}
	}

	if statistics.NumRevenue != nil {
		result.NumRevenue = &bdspropb.StatsDTO{
			Approved: statistics.NumRevenue.Approved,
			Reject:   statistics.NumRevenue.Rejected,
			Pending:  statistics.NumRevenue.Pending,
			Total:    statistics.NumRevenue.Total,
		}
	}

	if statistics.NumPayment != nil {
		result.NumPayment = &bdspropb.StatsDTO{
			Approved: statistics.NumPayment.Approved,
			Reject:   statistics.NumPayment.Rejected,
			Pending:  statistics.NumPayment.Pending,
			Total:    statistics.NumPayment.Total,
		}
	}

	result.TotalAmount = statistics.TotalAmount
	result.TotalCost = statistics.TotalCost
	result.TotalProfit = statistics.TotalProfit

	return result
}
