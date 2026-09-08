package transformer

// type DealTransactionTransformer interface {
// 	TransformCreateRequest(req *pb.CreateDealTransactionRequest) (*dto.CreateDealTransactionRequest, error)
// 	TransformUpdateRequest(req *pb.UpdateDealTransactionRequest) (*dto.UpdateDealTransactionRequest, error)
// 	// TransformToProtoResponse(response *entity.DealTransaction) *pb.DealContractResponse
// 	DealContractToPb(contract *entity.DealTransaction) *transactionpb.DealContractResponse
// 	DealContractToPbs(response []*entity.DealTransaction) []*transactionpb.DealContractResponse
// 	TransformToProtoFinancialSummary(summary *dto.DealFinancialSummary) *pb.FinancialSummaryResponse
// 	TransformTransactionType(protoType pb.TransactionType) enums.TransactionType
// }

// type dealTransactionTransformer struct{}

// func NewDealTransactionTransformer() DealTransactionTransformer {
// 	return &dealTransactionTransformer{}
// }

// func (t *dealTransactionTransformer) TransformCreateRequest(req *pb.CreateDealTransactionRequest) (*dto.CreateDealTransactionRequest, error) {
// 	// Parse transaction date
// 	var transactionDate *time.Time
// 	if req.TransactionDate != "" {
// 		parsedTime, err := time.Parse(time.RFC3339, req.TransactionDate)
// 		if err != nil {
// 			return nil, err
// 		}
// 		transactionDate = &parsedTime
// 	}

// 	return &dto.CreateDealTransactionRequest{
// 		DealID:          req.DealId,
// 		Type:            t.TransformTransactionType(req.Type),
// 		Amount:          req.Amount,
// 		Description:     req.Description,
// 		TransactionDate: transactionDate,
// 		Category:        req.Category,
// 		Note:            req.Note,
// 		DocumentIds:     req.DocumentIds,
// 		TransactionType: t.TransformTransactionType(req.Type),
// 	}, nil
// }

// func (t *dealTransactionTransformer) TransformUpdateRequest(req *pb.UpdateDealTransactionRequest) (*dto.UpdateDealTransactionRequest, error) {
// 	// Parse transaction date
// 	var transactionDate *time.Time
// 	if req.TransactionDate != "" {
// 		parsedTime, err := time.Parse(time.RFC3339, req.TransactionDate)
// 		if err != nil {
// 			return nil, err
// 		}
// 		transactionDate = &parsedTime
// 	}

// 	return &dto.UpdateDealTransactionRequest{
// 		ID:              req.Id,
// 		Amount:          req.Amount,
// 		Description:     req.Description,
// 		TransactionDate: transactionDate,
// 		Category:        req.Category,
// 		Note:            req.Note,
// 		DocumentIds:     req.DocumentIds,
// 	}, nil
// }

// func (t *dealTransactionTransformer) DealContractToPb(contract *entity.DealTransaction) *transactionpb.DealContractResponse {
// 	return &transactionpb.DealContractResponse{
// 		TransactionId:   contract.ID,
// 		DealId:          contract.DealID,
// 		Amount:          contract.Amount,
// 		CustomerId:      contract.CustomerID,
// 		CreatedAt:       _utils.FormatTimeToString(contract.CreatedAt),
// 		TransactionType: transactionpb.TransactionType(contract.Type),
// 		TypeName:        enums.TransactionTypeMap[contract.Type],
// 	}
// }

// func (t *dealTransactionTransformer) DealContractToPbs(response []*entity.DealTransaction) []*transactionpb.DealContractResponse {
// 	transactions := make([]*transactionpb.DealContractResponse, len(response))
// 	for i, transaction := range response {
// 		transactions[i] = t.DealContractToPb(transaction)
// 	}

// 	return transactions
// }

// func (t *dealTransactionTransformer) TransformToProtoFinancialSummary(summary *dto.DealFinancialSummary) *pb.FinancialSummaryResponse {
// 	return &pb.FinancialSummaryResponse{
// 		DealId:          summary.DealID,
// 		TotalRevenue:    summary.TotalRevenue,
// 		TotalExpense:    summary.TotalExpense,
// 		NetProfit:       summary.NetProfit,
// 		ProfitMargin:    summary.ProfitMargin,
// 		TargetProfit:    summary.TargetProfit,
// 		ProfitStatus:    summary.ProfitStatus,
// 		WarningMessages: summary.WarningMessages,
// 	}
// }

// func (t *dealTransactionTransformer) TransformTransactionType(protoType pb.TransactionType) enums.TransactionType {
// 	switch protoType {
// 	case pb.TransactionType_TRANSACTION_TYPE_EXPENSE:
// 		return enums.TransactionTypeExpense
// 	case pb.TransactionType_TRANSACTION_TYPE_REVENUE:
// 		return enums.TransactionTypeRevenue
// 	default:
// 		return enums.TransactionTypeExpense
// 	}
// }

// // func (t *dealTransactionTransformer) DealContractToPb(contract *entity.DealTransaction) *pb.DealContractResponse {
// // 	// result := t.TransformToProtoResponse(contract.Transaction)
// // 	result := &pb.DealContractResponse{
// // 		TransactionId: contract.ID,
// // 		DealId:        contract.DealID,
// // 		// TransactionName: contract.TypeName,
// // 		Amount:     contract.Amount,
// // 		CustomerId: contract.CustomerID,
// // 		// ProductId:       contract.ProductID,
// // 		// StaffId:         contract.StaffID,
// // 		// Status:          pb.ContractStatus(contract.Status),
// // 		StatusName: enums.ApprovedStatusMap[enums.ApprovedStatus(contract.Status)],
// // 		// Transaction: contract,
// // 		// Note:            response.Note,
// // 	}

// // 	// if contract != nil {
// // 	// 	contract.StaffID = response.ToId
// // 	// }

// // 	// result := t.TransformToProtoResponse(contract)

// // 	// if response.Actions != nil {
// // 	// 	result.Actions = make([]*pb.ActionResponse, len(response.Actions))
// // 	// 	for i, action := range response.Actions {
// // 	// 		result.Actions[i] = &pb.ActionResponse{
// // 	// 			ActionId: action.ActionId,
// // 	// 		}
// // 	// 	}
// // 	// }

// // 	// if response.ToOf == transactionpb.OwnerType_OWNER_TYPE_STAFF {
// // 	// 	result.StaffName = response.Profile.FullName
// // 	// }

// // 	return result
// // }
