package transformer

import (
	"organization/internal/dto"
	"organization/internal/enums"
	pb "pb/types/organization"
	transactionpb "pb/types/transaction"
	"time"
)

type DealTransactionTransformer interface {
	TransformCreateRequest(req *pb.CreateDealTransactionRequest) (*dto.CreateDealTransactionRequest, error)
	TransformUpdateRequest(req *pb.UpdateDealTransactionRequest) (*dto.UpdateDealTransactionRequest, error)
	TransformToProtoResponse(response *dto.DealTransactionResponse) *pb.DealContractResponse
	TransactionResponseToPb(response *transactionpb.TransactionResponse) *pb.DealContractResponse
	TransformToProtoListResponse(response *dto.DealTransactionListResponse) *pb.DealContractListResponse
	TransformToProtoFinancialSummary(summary *dto.DealFinancialSummary) *pb.FinancialSummaryResponse
	TransformTransactionType(protoType pb.TransactionType) enums.TransactionType
}

type dealTransactionTransformer struct{}

func NewDealTransactionTransformer() DealTransactionTransformer {
	return &dealTransactionTransformer{}
}

func (t *dealTransactionTransformer) TransformCreateRequest(req *pb.CreateDealTransactionRequest) (*dto.CreateDealTransactionRequest, error) {
	// Parse transaction date
	var transactionDate *time.Time
	if req.TransactionDate != "" {
		parsedTime, err := time.Parse(time.RFC3339, req.TransactionDate)
		if err != nil {
			return nil, err
		}
		transactionDate = &parsedTime
	}

	return &dto.CreateDealTransactionRequest{
		DealID:          req.DealId,
		Type:            t.TransformTransactionType(req.Type),
		Amount:          req.Amount,
		Description:     req.Description,
		TransactionDate: transactionDate,
		Category:        req.Category,
		Note:            req.Note,
		DocumentIds:     req.DocumentIds,
	}, nil
}

func (t *dealTransactionTransformer) TransformUpdateRequest(req *pb.UpdateDealTransactionRequest) (*dto.UpdateDealTransactionRequest, error) {
	// Parse transaction date
	var transactionDate *time.Time
	if req.TransactionDate != "" {
		parsedTime, err := time.Parse(time.RFC3339, req.TransactionDate)
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

func (t *dealTransactionTransformer) TransformToProtoResponse(response *dto.DealTransactionResponse) *pb.DealContractResponse {
	// Format dates
	// transactionDate := ""
	// if response.TransactionDate != nil {
	// 	transactionDate = response.TransactionDate.Format(time.RFC3339)
	// }

	// approvedAt := ""
	// if response.ApprovedAt != nil {
	// 	approvedAt = response.ApprovedAt.Format(time.RFC3339)
	// }

	return &pb.DealContractResponse{
		TransactionId:   response.ID,
		DealId:          response.DealID,
		TransactionName: response.TypeName,
		Amount:          response.Amount,
		// Description:  response.Description,
		// ContractDate: transactionDate,
		// Category:     response.Category,
		// Note:         response.Note,
		// CreatedBy:    response.CreatedBy,
		// Status:       pb.ContractStatus(response.Status),
		// StatusName:   response.StatusName,
		// ApprovedBy:   response.ApprovedBy,
		// ApprovedAt:   &approvedAt,
		// RejectReason: response.RejectReason,
		// DocumentIds:  response.DocumentIds,
		// CreatedAt:    response.CreatedAt.Format(time.RFC3339),
		// UpdatedAt:    response.UpdatedAt.Format(time.RFC3339),
	}
}

func (t *dealTransactionTransformer) TransformToProtoListResponse(response *dto.DealTransactionListResponse) *pb.DealContractListResponse {
	transactions := make([]*pb.DealContractResponse, len(response.Transactions))
	for i, transaction := range response.Transactions {
		transactions[i] = t.TransformToProtoResponse(&transaction)
	}

	return &pb.DealContractListResponse{
		Data:  transactions,
		Total: int32(response.Total),
	}
}

func (t *dealTransactionTransformer) TransformToProtoFinancialSummary(summary *dto.DealFinancialSummary) *pb.FinancialSummaryResponse {
	return &pb.FinancialSummaryResponse{
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

func (t *dealTransactionTransformer) TransformTransactionType(protoType pb.TransactionType) enums.TransactionType {
	switch protoType {
	case pb.TransactionType_TRANSACTION_TYPE_EXPENSE:
		return enums.TransactionTypeExpense
	case pb.TransactionType_TRANSACTION_TYPE_REVENUE:
		return enums.TransactionTypeRevenue
	default:
		return enums.TransactionTypeExpense
	}
}

func (t *dealTransactionTransformer) TransactionResponseToPb(response *transactionpb.TransactionResponse) *pb.DealContractResponse {
	result := &pb.DealContractResponse{
		TransactionId:   response.TransactionId,
		DealId:          response.FromId,
		TransactionName: response.TransactionName,
		Amount:          response.Value,
		CustomerId:      response.ToId,
		ProductId:       response.ToId,
		StaffId:         response.ToId,
		// Status:          pb.ContractStatus(response.Status),
		StatusName: response.Status.String(),
		// Note:            response.Note,
	}

	if response.Profile != nil {
		result.CustomerName = response.Profile.FullName
	}

	// if response.ToOf == transactionpb.OwnerType_OWNER_TYPE_STAFF {
	// 	result.StaffName = response.Profile.FullName
	// }

	return result
}
