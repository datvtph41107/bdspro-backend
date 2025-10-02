package handler

import (
	_dto "common/domain/dto"
	_utils "common/utils"
	"context"
	sharepb "pb/types/shared"
	transactionpb "pb/types/transaction"

	"tx/infra/client"
	"tx/infra/mapper"
	tx_domain "tx/internal/domain"
	"tx/internal/dto"
	"tx/internal/enums"
	usecase "tx/internal/usecases"
)

type TransactionHandler struct {
	transactionpb.UnimplementedTransactionServiceServer
	TransactionUsecase usecase.TransactionUsecase
	UserClient         *client.UserClient
	TransactionMapper  *mapper.TransactionMapper
}

func NewTransactionHandler(
	transactionUsecase usecase.TransactionUsecase,
	userClient *client.UserClient,
	transactionMapper *mapper.TransactionMapper,
) *TransactionHandler {
	return &TransactionHandler{
		TransactionUsecase: transactionUsecase,
		UserClient:         userClient,
		TransactionMapper:  transactionMapper,
	}
}

// @Summary Tạo giao dịch sản phẩm
// @Description Tạo giao dịch sản phẩm
// @Tags TransactionService
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param req body transactionpb.CreateProductTransactionRequest true "CreateProductTransactionRequest"
// @Success 200 {object} transactionpb.TransactionResponse "Transaction created successfully"
// @Failure 400 {object} error "Bad Request"
// @Failure 500 {object} error "Internal Server Error"
// @Router /product [post]
func (h *TransactionHandler) CreateProductTransaction(ctx context.Context, req *transactionpb.CreateProductTransactionRequest) (*transactionpb.TransactionResponse, error) {
	// transaction, err := h.TransactionUsecase.CreateTransaction(ctx, req)
	// if err != nil {
	// 	return nil, err
	// }
	return &transactionpb.TransactionResponse{
		TransactionId: 0,
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
		// Note:            action.Note,
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
	transactions, total, err := h.TransactionUsecase.ListTransactions(ctx, &dto.TransactionFilters{
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
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
			// Note:            transaction.Note,
			// FromId:          transaction.FromID,
			// FromOf:          transactionpb.OwnerType(transaction.FromOf),
			// ToId:            transaction.ToID,
			// ToOf:            transactionpb.OwnerType(transaction.ToOf),
		}
	}

	h.UserClient.MapProfileToTransaction(ctx, transactionsResponse)

	return &transactionpb.TransactionListResponse{
		Data:  transactionsResponse,
		Total: total,
	}, nil
}

// @Summary Tạo giao dịch sản phẩm
// @Description Tạo giao dịch sản phẩm
// @Tags TransactionService
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param req body transactionpb.CreateProductActionRequest true "CreateProductActionRequest"
// @Success 200 {object} transactionpb.ActionResponse "Action created successfully"
// @Failure 400 {object} error "Bad Request"
// @Failure 500 {object} error "Internal Server Error"
// @Router /product/action [post]
func (h *TransactionHandler) CreateProductAction(ctx context.Context, req *transactionpb.CreateProductActionRequest) (*transactionpb.ActionResponse, error) {
	action, err := h.TransactionUsecase.CreateAction(ctx, &tx_domain.TxAction{
		TxId: req.TransactionId,
		// FromId:    req.FromId,
		FromOf: enums.EOwnerUser,
		// ToId:      req.ToId,
		ToOf:   enums.EOwnerDealContract,
		Action: enums.EActionHandOver,
		// Value:     req.ActionAmount,
		Note:      req.ActionNote,
		Timestamp: _utils.TimeNowPtr(),
	})
	if err != nil {
		return nil, err
	}
	return &transactionpb.ActionResponse{
		TransactionId: action.ID,
	}, nil
}

// @Summary Lấy danh sách giao dịch sản phẩm
// @Description Lấy danh sách giao dịch sản phẩm
// @Tags TransactionService
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Deal ID"
// @Success 200 {object} transactionpb.TransactionListResponse "Transaction list"
// @Failure 400 {object} error "Bad Request"
// @Failure 500 {object} error "Internal Server Error"
// @Router /deal/{id} [get]
func (h *TransactionHandler) ListProductTransactions(ctx context.Context, req *sharepb.IdRequest) (*transactionpb.TransactionListResponse, error) {
	filters := &dto.TransactionFilters{
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
		},
		ToID: &req.Id,
		ToOf: enums.EOwnerDealContract,
	}

	transactions, total, err := h.TransactionUsecase.ListTransactions(ctx, filters)
	if err != nil {
		return nil, err
	}
	transactionsResponse := make([]*transactionpb.TransactionResponse, len(transactions))
	for i, transaction := range transactions {
		transactionsResponse[i] = h.TransactionMapper.TxToPb(transaction)
	}
	return &transactionpb.TransactionListResponse{
		Data:  transactionsResponse,
		Total: total,
	}, nil
}
