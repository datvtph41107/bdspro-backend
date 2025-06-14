package handler

import (
	"context"

	pbTransaction "pb/types/transaction"

	"github.com/tuandq2112/appointment-service/internal/interface/transformer"
	"github.com/tuandq2112/appointment-service/internal/interface/validator"
	"github.com/tuandq2112/appointment-service/internal/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TransactionHandler struct {
	pbTransaction.UnimplementedTransactionServiceServer
	TransactionUsecase usecase.TransactionUsecase
	Transformer        transformer.TransactionTransformer
	Validator          validator.TransactionValidator
}

func NewTransactionHandler(
	transactionUsecase usecase.TransactionUsecase,
	transformer transformer.TransactionTransformer,
	validator validator.TransactionValidator,
) *TransactionHandler {
	return &TransactionHandler{
		TransactionUsecase: transactionUsecase,
		Transformer:        transformer,
		Validator:          validator,
	}
}

// ListTransactions handles GET /transaction/list
func (h *TransactionHandler) ListTransactions(ctx context.Context, req *pbTransaction.ListTransactionsRequest) (*pbTransaction.ListTransactionsResponse, error) {
	if err := h.Validator.ValidateListTransactionsRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	transactions, total, err := h.TransactionUsecase.ListTransactions(ctx, nil)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return h.Transformer.EntityToListTransactionsResponse(transactions, int32(total)), nil
}

// GetTransactionDetail handles GET /transaction/detail/{id}
func (h *TransactionHandler) GetTransactionDetail(ctx context.Context, req *pbTransaction.GetTransactionDetailRequest) (*pbTransaction.Transaction, error) {
	if err := h.Validator.ValidateGetTransactionRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	transaction, err := h.TransactionUsecase.GetTransaction(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return h.Transformer.EntityToTransactionResponse(transaction), nil
}

// CreateTransaction handles POST /transaction/create
func (h *TransactionHandler) CreateTransaction(ctx context.Context, req *pbTransaction.CreateTransactionRequest) (*pbTransaction.Transaction, error) {
	if err := h.Validator.ValidateCreateTransactionRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	transaction := h.Transformer.CreateTransactionRequestToEntity(req)

	err := h.TransactionUsecase.CreateTransaction(ctx, transaction)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return h.Transformer.EntityToTransactionResponse(transaction), nil
}

// UpdateTransaction handles PUT /transaction/update/{id}
func (h *TransactionHandler) UpdateTransaction(ctx context.Context, req *pbTransaction.UpdateTransactionRequest) (*pbTransaction.Transaction, error) {
	if err := h.Validator.ValidateUpdateTransactionRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	transaction := h.Transformer.UpdateTransactionRequestToEntity(req)

	err := h.TransactionUsecase.UpdateTransaction(ctx, transaction)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return h.Transformer.EntityToTransactionResponse(transaction), nil
}

// DeleteTransaction handles DELETE /transaction/delete/{id}
func (h *TransactionHandler) DeleteTransaction(ctx context.Context, req *pbTransaction.DeleteTransactionRequest) (*pbTransaction.DeleteTransactionResponse, error) {
	if err := h.Validator.ValidateDeleteTransactionRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := h.TransactionUsecase.DeleteTransaction(ctx, req.Id); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pbTransaction.DeleteTransactionResponse{
		Id: req.Id,
	}, nil
}

// ApproveTransaction handles POST /transaction/approve/{id}
func (h *TransactionHandler) ApproveTransaction(ctx context.Context, req *pbTransaction.ApproveTransactionRequest) (*pbTransaction.Transaction, error) {
	if err := h.Validator.ValidateApproveTransactionRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := h.TransactionUsecase.ApproveTransaction(ctx, req.Id, 0); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	transaction, err := h.TransactionUsecase.GetTransaction(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return h.Transformer.EntityToTransactionResponse(transaction), nil
}

// RejectTransaction handles POST /transaction/reject/{id}
func (h *TransactionHandler) RejectTransaction(ctx context.Context, req *pbTransaction.RejectTransactionRequest) (*pbTransaction.Transaction, error) {
	if err := h.Validator.ValidateRejectTransactionRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := h.TransactionUsecase.RejectTransaction(ctx, req.Id); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	transaction, err := h.TransactionUsecase.GetTransaction(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return h.Transformer.EntityToTransactionResponse(transaction), nil
}

// GetTransactionTypes handles GET /transaction/types
func (h *TransactionHandler) GetTransactionTypes(ctx context.Context, req *pbTransaction.GetTransactionTypesRequest) (*pbTransaction.GetTransactionTypesResponse, error) {
	if err := h.Validator.ValidateGetTransactionTypesRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	types, err := h.TransactionUsecase.GetTransactionTypes(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return h.Transformer.EntityToTransactionTypesResponse(types), nil
} 