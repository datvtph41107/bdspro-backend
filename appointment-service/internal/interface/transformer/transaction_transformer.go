package transformer

import (
	"time"

	pbTransaction "pb/types/transaction"

	"github.com/tuandq2112/appointment-service/internal/domain/entity/transaction"
)

type TransactionTransformer interface {
	CreateTransactionRequestToEntity(req *pbTransaction.CreateTransactionRequest) *transaction.Transaction
	UpdateTransactionRequestToEntity(req *pbTransaction.UpdateTransactionRequest) *transaction.Transaction
	EntityToTransactionResponse(tx *transaction.Transaction) *pbTransaction.Transaction
	EntityToListTransactionsResponse(txs []*transaction.Transaction, total int32) *pbTransaction.ListTransactionsResponse
	EntityToTransactionTypesResponse(types []string) *pbTransaction.GetTransactionTypesResponse
}

type transactionTransformer struct{}

func NewTransactionTransformer() TransactionTransformer {
	return &transactionTransformer{}
}

func (t *transactionTransformer) CreateTransactionRequestToEntity(req *pbTransaction.CreateTransactionRequest) *transaction.Transaction {
	transactionDate, _ := time.Parse(time.RFC3339, req.TransactionDate)
	return &transaction.Transaction{
		OwnerId:            req.OwnerId,
		OrganizationId:     req.OrganizationId,
		Amount:            req.Amount,
		Currency:          req.Currency,
		TransactionName:   req.TransactionName,
		Description:       req.Description,
		TransactionType:   transaction.TransactionType(req.TransactionType),
		CategoryId:        req.CategoryId,
		PaymentMethodId:   req.PaymentMethodId,
		RelatedDealId:     req.RelatedDealId,
		RelatedTransactionId: req.RelatedTransactionId,
		TransactionDate:   transactionDate,
		TransactionStatus: transaction.TransactionStatusDraft,
		ApprovalStatus:    transaction.ApprovalStatusPending,
	}
}

func (t *transactionTransformer) UpdateTransactionRequestToEntity(req *pbTransaction.UpdateTransactionRequest) *transaction.Transaction {
	transactionDate, _ := time.Parse(time.RFC3339, req.TransactionDate)
	return &transaction.Transaction{
		Id:                req.Id,
		Amount:            req.Amount,
		Currency:          req.Currency,
		TransactionName:   req.TransactionName,
		Description:       req.Description,
		CategoryId:        req.CategoryId,
		PaymentMethodId:   req.PaymentMethodId,
		RelatedDealId:     req.RelatedDealId,
		RelatedTransactionId: req.RelatedTransactionId,
		TransactionDate:   transactionDate,
	}
}

func (t *transactionTransformer) EntityToTransactionResponse(tx *transaction.Transaction) *pbTransaction.Transaction {
	if tx == nil {
		return nil
	}
	var approvalDate *string
	if tx.ApprovalDate != nil && !tx.ApprovalDate.IsZero() {
		date := tx.ApprovalDate.Format(time.RFC3339)
		approvalDate = &date
	}

	return &pbTransaction.Transaction{
		Id:                  tx.Id,
		OwnerId:            tx.OwnerId,
		OrganizationId:     tx.OrganizationId,
		Amount:            tx.Amount,
		Currency:          tx.Currency,
		TransactionName:   tx.TransactionName,
		Description:       tx.Description,
		TransactionType:   string(tx.TransactionType),
		CategoryId:        tx.CategoryId,
		PaymentMethodId:   tx.PaymentMethodId,
		RelatedDealId:     tx.RelatedDealId,
		RelatedTransactionId: tx.RelatedTransactionId,
		ApprovalStatus:    string(tx.ApprovalStatus),
		ApprovedBy:        tx.ApprovedBy,
		ApprovalDate:      approvalDate,
		TransactionStatus: string(tx.TransactionStatus),
		TransactionDate:   tx.TransactionDate.Format(time.RFC3339),
		IsDeleted:         tx.IsDeleted,
		CreatedAt:         tx.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         tx.UpdatedAt.Format(time.RFC3339),
	}
}

func (t *transactionTransformer) EntityToListTransactionsResponse(txs []*transaction.Transaction, total int32) *pbTransaction.ListTransactionsResponse {
	transactions := make([]*pbTransaction.Transaction, 0, len(txs))
	for _, tx := range txs {
		transactions = append(transactions, t.EntityToTransactionResponse(tx))
	}

	return &pbTransaction.ListTransactionsResponse{
		Transactions: transactions,
		Total:       total,
	}
}

func (t *transactionTransformer) EntityToTransactionTypesResponse(types []string) *pbTransaction.GetTransactionTypesResponse {
	return &pbTransaction.GetTransactionTypesResponse{
		Types: types,
	}
} 