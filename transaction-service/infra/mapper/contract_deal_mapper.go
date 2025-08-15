package mapper

import (
	_utils "common/utils"
	transactionpb "pb/types/transaction"
	"time"
	"tx/internal/domain"
	"tx/internal/domain/entity"
	"tx/internal/dto"
	"tx/internal/enums"
)

type ContractDealMapper struct {
	transactionMapper *TransactionMapper
}

func NewContractDealMapper(transactionMapper *TransactionMapper) *ContractDealMapper {
	return &ContractDealMapper{
		transactionMapper: transactionMapper,
	}
}

// func (m *ContractDealMapper) CreateDealTransactionRequestToEntity(req *dto.DealContractPaymentRequest) *entity.ContractDeal {
// 	return &entity.ContractDeal{
// 		DealID: req.DealId,
// 	}
// }

// func (t *ContractDealMapper) DealContractToPb(protoType transactionpb.TransactionType) enums.TransactionType {
// 	switch protoType {
// 	case transactionpb.TransactionType_TRANSACTION_TYPE_EXPENSE:
// 		return enums.TransactionTypeExpense
// 	case transactionpb.TransactionType_TRANSACTION_TYPE_REVENUE:
// 		return enums.TransactionTypeRevenue
// 	}
// }

func (t *ContractDealMapper) DealContractPaymentReqToEntity(req *transactionpb.DealContractPaymentRequest) (*dto.DealContractPaymentRequest, error) {
	// Parse transaction date
	// var transactionDate *time.Time
	// if req.Timestamp != "" {
	// 	parsedTime, err := time.Parse(time.RFC3339, req.Timestamp)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	transactionDate = &parsedTime
	// }

	return &dto.DealContractPaymentRequest{
		DealId: req.DealId,
		// Type:            t.TransformTransactionType(transactionpb.ContractType(req.TransactionType)),
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

func (t *ContractDealMapper) TransformUpdateRequest(req *transactionpb.UpdateDealContractRequest) (*dto.UpdateDealTransactionRequest, error) {
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

func (t *ContractDealMapper) DealContractToPb(contract *entity.ContractDeal) *transactionpb.DealContractResponse {
	result := &transactionpb.DealContractResponse{
		ContractId:      contract.ID,
		DealId:          contract.DealID,
		CustomerId:      contract.CustomerID,
		CreatedAt:       _utils.FormatTimeToString(contract.CreatedAt),
		TransactionType: transactionpb.TransactionType(contract.Type),
		TypeName:        enums.TransactionTypeMap[contract.Type],
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

func (t *ContractDealMapper) DealContractToPbs(response []*entity.ContractDeal) []*transactionpb.DealContractResponse {
	transactions := make([]*transactionpb.DealContractResponse, len(response))
	for i, transaction := range response {
		transactions[i] = t.DealContractToPb(transaction)
	}

	return transactions
}

func (t *ContractDealMapper) TransformToProtoFinancialSummary(summary *dto.DealFinancialSummary) *transactionpb.FinancialSummaryResponse {
	return &transactionpb.FinancialSummaryResponse{
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

func (t *ContractDealMapper) TransformTransactionType(protoType transactionpb.ContractType) enums.TransactionType {
	switch protoType {
	case transactionpb.ContractType_CONTRACT_TYPE_EXPENSE:
		return enums.TransactionTypeExpense
	case transactionpb.ContractType_CONTRACT_TYPE_REVENUE:
		return enums.TransactionTypeRevenue
	default:
		return enums.TransactionTypeExpense
	}
}

func (t *ContractDealMapper) DealContractItemToPb(contract *dto.ContractDealListResponse) *transactionpb.DealContractItem {
	result := &transactionpb.DealContractItem{
		Id:         contract.ContractID,
		CustomerId: contract.CustomerID,
		CreatedAt:  _utils.FormatTimeToString(contract.Timestamp),
		ContractId: contract.ContractID,
		Amount:     contract.Amount,
		Status:     uint64(contract.Status),
	}

	if contract.TransactionName == "" {
		result.TransactionName = enums.MethodNames[enums.EMethod(contract.Method)]
		result.Status = uint64(contract.Status)
		result.StatusName = enums.ApprovedStatusMap[enums.ApprovedStatus(contract.Status)]
	}

	return result
}

func (t *ContractDealMapper) DealContractItemToPbs(contracts []*dto.ContractDealListResponse) []*transactionpb.DealContractItem {
	items := make([]*transactionpb.DealContractItem, len(contracts))
	for i, contract := range contracts {
		items[i] = t.DealContractItemToPb(contract)
	}
	return items
}

// DealContractItemToPbsFromEntity maps a slice of ContractDeal entities to protobuf items
func (t *ContractDealMapper) DealContractItemToPbsFromEntity(contracts []*entity.ContractDeal) []*transactionpb.DealContractItem {
	items := make([]*transactionpb.DealContractItem, len(contracts))
	for i, contract := range contracts {
		items[i] = t.DealContractItemToPbFromEntity(contract)
	}
	return items
}

// DealContractItemToPbFromEntity maps a single ContractDeal entity to protobuf item
func (t *ContractDealMapper) DealContractItemToPbFromEntity(contract *entity.ContractDeal) *transactionpb.DealContractItem {
	result := &transactionpb.DealContractItem{
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
		result.StatusName = enums.ApprovedStatusMap[enums.ApprovedStatus(contract.Transaction.Status)]
		result.ContractId = contract.ID
	}

	return result
}

func (t *ContractDealMapper) PaymentActionReqToDTO(req *transactionpb.PaymentActionRequest) *dto.ActionPaymentRequest {
	return &dto.ActionPaymentRequest{
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

func (t *ContractDealMapper) ActionToPb(action *domain.TxAction) *transactionpb.ActionResponse {
	return &transactionpb.ActionResponse{
		Action:        uint32(action.Action),
		Note:          action.Note,
		FromId:        action.FromId,
		FromOf:        transactionpb.OwnerType(action.FromOf),
		ToId:          action.ToId,
		ToOf:          transactionpb.OwnerType(action.ToOf),
		Value:         action.Value,
		TransactionId: action.TxId,
		CreatedAt:     _utils.FormatTimeToString(action.CreatedAt),
		Message:       action.Message,
	}
}

// StatisticsToPb maps DealComprehensiveStatisticsDTO to protobuf
func (t *ContractDealMapper) StatisticsToPb(statistics *dto.DealCostStatisticsDTO) *transactionpb.ContractDealStatisticsDTO {
	result := &transactionpb.ContractDealStatisticsDTO{
		DealId: statistics.DealID,
	}

	if statistics.Payment != nil {
		result.Payment = &transactionpb.StatsDTO{
			Approved: statistics.Payment.Approved,
			Reject:   statistics.Payment.Reject,
			Pending:  statistics.Payment.Pending,
			Total:    statistics.Payment.Total,
		}
	}

	if statistics.AmountCost != nil {
		result.AmountCost = &transactionpb.StatsDTO{
			Approved: statistics.AmountCost.Approved,
			Reject:   statistics.AmountCost.Reject,
			Pending:  statistics.AmountCost.Pending,
			Total:    statistics.AmountCost.Total,
		}
	}

	if statistics.AmountRevenue != nil {
		result.AmountRevenue = &transactionpb.StatsDTO{
			Approved: statistics.AmountRevenue.Approved,
			Reject:   statistics.AmountRevenue.Reject,
			Pending:  statistics.AmountRevenue.Pending,
			Total:    statistics.AmountRevenue.Total,
		}
	}

	if statistics.NumCost != nil {
		result.NumCost = &transactionpb.StatsDTO{
			Approved: statistics.NumCost.Approved,
			Reject:   statistics.NumCost.Reject,
			Pending:  statistics.NumCost.Pending,
			Total:    statistics.NumCost.Total,
		}
	}

	if statistics.NumRevenue != nil {
		result.NumRevenue = &transactionpb.StatsDTO{
			Approved: statistics.NumRevenue.Approved,
			Reject:   statistics.NumRevenue.Reject,
			Pending:  statistics.NumRevenue.Pending,
			Total:    statistics.NumRevenue.Total,
		}
	}

	if statistics.NumPayment != nil {
		result.NumPayment = &transactionpb.StatsDTO{
			Approved: statistics.NumPayment.Approved,
			Reject:   statistics.NumPayment.Reject,
			Pending:  statistics.NumPayment.Pending,
			Total:    statistics.NumPayment.Total,
		}
	}

	result.TotalAmount = statistics.TotalAmount
	result.TotalCost = statistics.TotalCost
	result.TotalProfit = statistics.TotalProfit

	return result
}
