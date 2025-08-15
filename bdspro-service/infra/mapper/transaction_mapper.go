package mapper

// type TransactionTransformer interface {
// 	CreateTransactionRequestToEntity(req *bdspropb.CreateTransactionRequest) *tx_domain.Tx
// 	UpdateTransactionRequestToEntity(req *bdspropb.UpdateTransactionRequest) *tx_domain.Tx
// 	EntityToTransactionResponse(tx *tx_domain.Tx) *bdspropb.Transaction
// 	EntityToListTransactionsResponse(txs []*tx_domain.Tx, total int32) *bdspropb.ListTransactionsResponse
// 	EntityToTransactionTypesResponse(types []string) *bdspropb.GetTransactionTypesResponse
// }

// type transactionTransformer struct {
// 	DepositeMapper *DepositeMapper
// }

// func NewTransactionTransformer(
// 	depositeMapper *DepositeMapper,
// ) TransactionTransformer {
// 	return &transactionTransformer{
// 		DepositeMapper: depositeMapper,
// 	}
// }

// // func (t *transactionTransformer) GetStatusTransaction(e *domain.Transaction) uint32 {
// // 	if e.DepositeID != nil {
// // 		return 20
// // 	}
// // 	return 10
// // }

// func (t *transactionTransformer) CreateTransactionRequestToEntity(req *bdspropb.CreateTransactionRequest) *tx_domain.Tx {
// 	transactionDate, _ := time.Parse(time.RFC3339, req.TransactionDate)
// 	return &tx_domain.Tx{
// 		// OrganizationId:       req.OrganizationId,
// 		OwnerId:              uint64(req.OwnerId),
// 		Amount:               req.Amount,
// 		Currency:             req.Currency,
// 		TransactionName:      req.TransactionName,
// 		Description:          req.Description,
// 		TransactionType:      enums.TransactionType(req.TransactionType),
// 		CategoryId:           req.CategoryId,
// 		PaymentMethodId:      req.PaymentMethodId,
// 		RelatedDealId:        req.RelatedDealId,
// 		RelatedTransactionId: req.RelatedTransactionId,
// 		TransactionDate:      transactionDate,
// 		TransactionStatus:    enums.TransactionStatusDraft,
// 		ApprovalStatus:       enums.ApprovalStatusPending,
// 	}
// }

// func (t *transactionTransformer) UpdateTransactionRequestToEntity(req *bdspropb.UpdateTransactionRequest) *tx_domain.Tx {
// 	transactionDate, _ := time.Parse(time.RFC3339, req.TransactionDate)
// 	return &tx_domain.Tx{
// 		BaseEntity: _models.BaseEntity{
// 			ID: uint64(req.Id),
// 		},
// 		Amount:               req.Amount,
// 		Currency:             req.Currency,
// 		TransactionName:      req.TransactionName,
// 		Description:          req.Description,
// 		CategoryId:           req.CategoryId,
// 		PaymentMethodId:      req.PaymentMethodId,
// 		RelatedDealId:        req.RelatedDealId,
// 		RelatedTransactionId: req.RelatedTransactionId,
// 		TransactionDate:      transactionDate,
// 	}
// }

// func (t *transactionTransformer) EntityToTransactionResponse(tx *tx_domain.Tx) *bdspropb.Transaction {
// 	if tx == nil {
// 		return nil
// 	}
// 	var approvalDate *string
// 	if tx.ApprovalDate != nil && !tx.ApprovalDate.IsZero() {
// 		date := tx.ApprovalDate.Format(time.RFC3339)
// 		approvalDate = &date
// 	}

// 	result := &bdspropb.Transaction{
// 		Id:      tx.ID,
// 		OwnerId: uint32(tx.OwnerId),
// 		// OrganizationId:       tx.OrganizationId,
// 		Amount:               tx.Amount,
// 		Currency:             tx.Currency,
// 		TransactionName:      tx.TransactionName,
// 		Description:          tx.Description,
// 		TransactionType:      string(tx.TransactionType),
// 		CategoryId:           tx.CategoryId,
// 		PaymentMethodId:      tx.PaymentMethodId,
// 		RelatedDealId:        tx.RelatedDealId,
// 		RelatedTransactionId: tx.RelatedTransactionId,
// 		ApprovalStatus:       string(tx.ApprovalStatus),
// 		ApprovedBy:           tx.ApprovedBy,
// 		ApprovalDate:         approvalDate,
// 		TransactionStatus:    uint32(tx.TransactionStatus),
// 		TransactionDate:      tx.TransactionDate.Format(time.RFC3339),
// 		DepositeAmount:       tx.DepositeAmount,
// 		DepositeNote:         tx.DepositeNote,
// 		ContactId:            tx.ContactID,
// 	}

// 	return result
// }

// func (t *transactionTransformer) EntityToListTransactionsResponse(txs []*tx_domain.Tx, total int32) *bdspropb.ListTransactionsResponse {
// 	transactions := make([]*bdspropb.Transaction, 0, len(txs))
// 	for _, tx := range txs {
// 		transactions = append(transactions, t.EntityToTransactionResponse(tx))
// 	}

// 	return &bdspropb.ListTransactionsResponse{
// 		Data:  transactions,
// 		Total: total,
// 	}
// }

// func (t *transactionTransformer) EntityToTransactionTypesResponse(types []string) *bdspropb.GetTransactionTypesResponse {
// 	return &bdspropb.GetTransactionTypesResponse{
// 		Types: types,
// 	}
// }
