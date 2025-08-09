package handler

import (
	_models "common/models"
	"context"
	transactionpb "pb/types/transaction"
	tx_domain "tx/internal/domain"
	"tx/internal/enums"
	usecase "tx/internal/usecases"
)

type InternalTransactionHandler struct {
	transactionpb.UnimplementedInternalTransactionServiceServer
	TransactionUsecase usecase.TransactionUsecase
}

func NewInternalTransactionHandler(transactionUsecase usecase.TransactionUsecase) *InternalTransactionHandler {
	return &InternalTransactionHandler{
		TransactionUsecase: transactionUsecase,
	}
}

func (h *InternalTransactionHandler) Transfer(ctx context.Context, req *transactionpb.TransferRequest) (*transactionpb.ActionResponse, error) {
	return nil, nil
	// action, err := h.TransactionUsecase.CreateAction(ctx, &domain.TxAction{
	// 	TxId:      req.Id,
	// 	FromId:    req.FromId,
	// 	ToId:      req.ToId,
	// 	Method:    enums.EMethodTransfer,
	// 	Note:      req.Note,
	// 	Timestamp: time.Now(),
	// })
}

func (h *InternalTransactionHandler) Approve(ctx context.Context, req *transactionpb.ApproveRequest) (*transactionpb.ActionResponse, error) {
	return nil, nil
	// action, err := h.TransactionUsecase.CreateAction(ctx, &domain.TxAction{
	// 	TxId:      req.Id,
	// 	Method:    enums.EMethodApprove,
	// 	Note:      req.Note,
	// 	Timestamp: time.Now(),
	// })
}

func (h *InternalTransactionHandler) Reject(ctx context.Context, req *transactionpb.RejectRequest) (*transactionpb.ActionResponse, error) {
	return nil, nil
	// action, err := h.TransactionUsecase.CreateAction(ctx, &domain.TxAction{
	// 	TxId:      req.Id,
	// 	Method:    enums.EMethodReject,
	// 	Note:      req.Note,
	// 	Timestamp: time.Now(),
	// })
}

func (h *InternalTransactionHandler) Sign(ctx context.Context, req *transactionpb.SignRequest) (*transactionpb.ActionResponse, error) {
	return nil, nil
	// action, err := h.TransactionUsecase.CreateAction(ctx, &domain.TxAction{
	// 	TxId:      req.Id,
	// 	Method:    enums.EMethodSign,
	// 	Note:      req.Note,
	// 	Timestamp: time.Now(),
	// })
}

func (h *InternalTransactionHandler) CreateSignTransaction(ctx context.Context, req *transactionpb.SignRequest) (*transactionpb.ActionResponse, error) {
	tx, err := h.TransactionUsecase.CreateSignTransaction(ctx, &tx_domain.Tx{
		BaseEntity:      _models.BaseEntity{ID: req.Id},
		TransactionName: req.ContractName,
		Note:            req.Note,
		FromID:          req.FromId,
		FromOf:          enums.EOwnerType(req.FromOf),
		ToID:            req.ToId,
		ToOf:            enums.EOwnerType(req.ToOf),
		Value:           req.Amount,
	})
	if err != nil {
		return nil, err
	}
	return &transactionpb.ActionResponse{
		TransactionId: tx.ID,
	}, nil
}
