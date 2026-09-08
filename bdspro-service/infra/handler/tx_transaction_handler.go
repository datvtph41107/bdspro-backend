package handler

import (
	_dto "common/domain/dto"
	_utils "common/utils"
	"context"
	sharepb "pb/types/shared"

	"bdspro/infra/client"
	"bdspro/infra/mapper"
	tx_domain "bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	usecase "bdspro/internal/usecases"
	bdspropb "pb/types/bdspro"
)

type TxHandler struct {
	bdspropb.UnimplementedTxServiceServer
	TransactionUsecase usecase.TxTransactionUsecase
	UserClient         *client.UserClient
	TransactionMapper  *mapper.TxMapper
}

func NewTxHandler(
	transactionUsecase usecase.TxTransactionUsecase,
	userClient *client.UserClient,
	transactionMapper *mapper.TxMapper,
) *TxHandler {
	return &TxHandler{
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
// @Param req body bdspropb.CreateProductTransactionRequest true "CreateProductTransactionRequest"
// @Success 200 {object} bdspropb.TransactionResponse "Transaction created successfully"
// @Failure 400 {object} error "Bad Request"
// @Failure 500 {object} error "Internal Server Error"
// @Router /product [post]
func (h *TxHandler) CreateProductTransaction(ctx context.Context, req *bdspropb.CreateProductTransactionRequest) (*bdspropb.TransactionResponse, error) {
	// transaction, err := h.TransactionUsecase.CreateTransaction(ctx, req)
	// if err != nil {
	// 	return nil, err
	// }
	return &bdspropb.TransactionResponse{
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
// @Success 200 {object} bdspropb.TransactionResponse "Transaction approved successfully"
// @Failure 400 {object} error "Bad Request"
// @Failure 404 {object} error "Transaction not found"
// @Failure 500 {object} error "Internal Server Error"
// @Router /{id}/approve [post]
func (h *TxHandler) ApproveTransaction(ctx context.Context, req *sharepb.IdRequest) (*bdspropb.TransactionResponse, error) {
	action, err := h.TransactionUsecase.ApproveTransaction(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &bdspropb.TransactionResponse{
		TransactionId:   action.ID,
		TransactionName: action.TransactionName,
		Status:          bdspropb.TransactionStatus(action.Status),
		Method:          bdspropb.TransactionMethod(action.Method),
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
// @Param page query bdspropb.TransactionListRequest true "query"
// @Success 200 {object} bdspropb.TransactionResponse "Transaction list"
// @Failure 400 {object} error "Bad Request"
// @Failure 500 {object} error "Internal Server Error"
// @Router /list [get]
func (h *TxHandler) ListTransactions(ctx context.Context, req *bdspropb.TransactionListRequest) (*bdspropb.TransactionListResponse, error) {
	transactions, total, err := h.TransactionUsecase.ListTransactions(ctx, &dto.TxTransactionFilters{
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
		},
	})
	if err != nil {
		return nil, err
	}
	transactionsResponse := make([]*bdspropb.TransactionResponse, len(transactions))
	for i, transaction := range transactions {
		transactionsResponse[i] = &bdspropb.TransactionResponse{
			TransactionId:   transaction.ID,
			TransactionName: transaction.TransactionName,
			Status:          bdspropb.TransactionStatus(transaction.Status),
			Method:          bdspropb.TransactionMethod(transaction.Method),
			Value:           transaction.Value,
			// Note:            transaction.Note,
			// FromId:          transaction.FromID,
			// FromOf:          bdspropb.OwnerType(transaction.FromOf),
			// ToId:            transaction.ToID,
			// ToOf:            bdspropb.OwnerType(transaction.ToOf),
		}
	}

	h.UserClient.MapProfileToTransaction(ctx, transactionsResponse)

	return &bdspropb.TransactionListResponse{
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
// @Param req body bdspropb.CreateProductActionRequest true "CreateProductActionRequest"
// @Success 200 {object} bdspropb.ActionResponse "Action created successfully"
// @Failure 400 {object} error "Bad Request"
// @Failure 500 {object} error "Internal Server Error"
// @Router /product/action [post]
func (h *TxHandler) CreateProductAction(ctx context.Context, req *bdspropb.CreateProductActionRequest) (*bdspropb.ActionResponse, error) {
	action, err := h.TransactionUsecase.CreateAction(ctx, &tx_domain.TxAction{
		TxId: req.TransactionId,
		// FromId:    req.FromId,
		FromOf: enums.TxOwnerUser,
		// ToId:      req.ToId,
		ToOf:   enums.TxOwnerDealContract,
		Action: enums.TxActionHandOver,
		// Value:     req.ActionAmount,
		Note:      req.ActionNote,
		Timestamp: _utils.TimeNowPtr(),
	})
	if err != nil {
		return nil, err
	}
	return &bdspropb.ActionResponse{
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
// @Success 200 {object} bdspropb.TransactionListResponse "Transaction list"
// @Failure 400 {object} error "Bad Request"
// @Failure 500 {object} error "Internal Server Error"
// @Router /deal/{id} [get]
func (h *TxHandler) ListProductTransactions(ctx context.Context, req *sharepb.IdRequest) (*bdspropb.TransactionListResponse, error) {
	filters := &dto.TxTransactionFilters{
		Pagable: _dto.Pagable{
			Page: uint32(req.Page),
			Size: uint32(req.Size),
		},
		ToID: &req.Id,
		ToOf: enums.TxOwnerDealContract,
	}

	transactions, total, err := h.TransactionUsecase.ListTransactions(ctx, filters)
	if err != nil {
		return nil, err
	}
	transactionsResponse := make([]*bdspropb.TransactionResponse, len(transactions))
	for i, transaction := range transactions {
		transactionsResponse[i] = h.TransactionMapper.TxToPb(transaction)
	}
	return &bdspropb.TransactionListResponse{
		Data:  transactionsResponse,
		Total: total,
	}, nil
}
