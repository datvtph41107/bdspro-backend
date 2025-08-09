package handler

import (
	_dto "common/dto"
	"context"
	sharepb "pb/types/shared"
	transactionpb "pb/types/transaction"
	"tx/infra/client"
	"tx/internal/repo"
	usecase "tx/internal/usecases"
)

type TransactionHandler struct {
	transactionpb.UnimplementedTransactionServiceServer
	TransactionUsecase usecase.TransactionUsecase
	UserClient         *client.UserClient
}

func NewTransactionHandler(transactionUsecase usecase.TransactionUsecase, userClient *client.UserClient) *TransactionHandler {
	return &TransactionHandler{
		TransactionUsecase: transactionUsecase,
		UserClient:         userClient,
	}
}

func (h *TransactionHandler) SaveTransaction(ctx context.Context, req *transactionpb.SaveTransactionRequest) (*transactionpb.SaveTransactionResponse, error) {
	// transaction, err := h.TransactionUsecase.CreateTransaction(ctx, req)
	// if err != nil {
	// 	return nil, err
	// }
	return &transactionpb.SaveTransactionResponse{
		TransactionId: "",
	}, nil
}

// @Summary Duyệt giao dịch
// @Description Duyệt giao dịch
// @Tags TransactionService
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Transaction ID"
// @Success 200 {object} transactionpb.TransactionResponse "Transaction approved successfully"
// @Failure 400 {object} error "Bad Request"
// @Failure 404 {object} error "Transaction not found"
// @Failure 500 {object} error "Internal Server Error"
// @Router /{id}/approve [post]
func (h *TransactionHandler) ApproveTransaction(ctx context.Context, req *sharepb.IdRequest) (*transactionpb.TransactionResponse, error) {
	action, err := h.TransactionUsecase.ApproveTransaction(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &transactionpb.TransactionResponse{
		TransactionId:   action.ID,
		TransactionName: action.TransactionName,
		Status:          transactionpb.TransactionStatus(action.Status),
		Method:          transactionpb.TransactionMethod(action.Method),
		Value:           action.Value,
		Note:            action.Note,
	}, nil
}

// @Summary Lấy danh sách giao dịch
// @Description Lấy danh sách giao dịch
// @Tags TransactionService
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query transactionpb.TransactionListRequest true "query"
// @Success 200 {object} transactionpb.TransactionResponse "Transaction list"
// @Failure 400 {object} error "Bad Request"
// @Failure 500 {object} error "Internal Server Error"
// @Router /list [get]
func (h *TransactionHandler) ListTransactions(ctx context.Context, req *transactionpb.TransactionListRequest) (*transactionpb.TransactionListResponse, error) {
	transactions, total, err := h.TransactionUsecase.ListTransactions(ctx, &repo.TransactionFilters{
		Pagable: _dto.Pagable{
			Page: int(req.Page),
			Size: int(req.Size),
		},
	})
	if err != nil {
		return nil, err
	}
	transactionsResponse := make([]*transactionpb.TransactionResponse, len(transactions))
	for i, transaction := range transactions {
		transactionsResponse[i] = &transactionpb.TransactionResponse{
			TransactionId:   transaction.ID,
			TransactionName: transaction.TransactionName,
			Status:          transactionpb.TransactionStatus(transaction.Status),
			Method:          transactionpb.TransactionMethod(transaction.Method),
			Value:           transaction.Value,
			Note:            transaction.Note,
			FromId:          transaction.FromID,
			FromOf:          transactionpb.OwnerType(transaction.FromOf),
			ToId:            transaction.ToID,
			ToOf:            transactionpb.OwnerType(transaction.ToOf),
		}
	}

	h.UserClient.MapProfileToTransaction(ctx, transactionsResponse)

	return &transactionpb.TransactionListResponse{
		Data:  transactionsResponse,
		Total: total,
	}, nil
}
