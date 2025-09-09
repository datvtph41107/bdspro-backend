package handler

import (
	_utils "common/utils"
	"context"
	sharepb "pb/types/shared"
	transactionpb "pb/types/transaction"
	"time"
	"tx/infra/mapper"
	"tx/infra/validator"
	tx_domain "tx/internal/domain"
	"tx/internal/enums"
	"tx/internal/repo"
	usecase "tx/internal/usecases"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InternalTransactionHandler struct {
	transactionpb.UnimplementedInternalTransactionServiceServer
	TransactionUsecase  usecase.TransactionUsecase
	DealContractUsecase *usecase.DealContractUsecase
	Validator           validator.TransactionActionValidator
	Mapper              *mapper.TransactionMapper
	DealMapper          *mapper.ContractDealMapper
	DealCostRepo        repo.DealCostRepository
}

func NewInternalTransactionHandler(
	transactionUsecase usecase.TransactionUsecase,
	dealContractUsecase *usecase.DealContractUsecase,
	validator validator.TransactionActionValidator,
	mapper *mapper.TransactionMapper,
	dealMapper *mapper.ContractDealMapper,
	dealCostRepo repo.DealCostRepository,
) *InternalTransactionHandler {
	return &InternalTransactionHandler{
		TransactionUsecase:  transactionUsecase,
		DealContractUsecase: dealContractUsecase,
		Validator:           validator,
		Mapper:              mapper,
		DealMapper:          dealMapper,
		DealCostRepo:        dealCostRepo,
	}
}

// @bind: internal/usecases.TransactionUsecase
func (h *InternalTransactionHandler) CreateAction(ctx context.Context, req *transactionpb.CreateActionRequest) (*transactionpb.ActionResponse, error) {
	// Validate request
	if err := h.Validator.ValidateCreateActionRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	tx, err := h.TransactionUsecase.CreateAction(ctx, &tx_domain.TxAction{
		FromId: req.FromId,
		TxId:   req.TransactionId,
		FromOf: enums.EOwnerType(req.FromOf),
		ToId:   req.ToId,
		ToOf:   enums.EOwnerType(req.ToOf),
		Action: enums.EAction(req.Action),
		Value:  req.Amount,
		Note:   req.Note,
	})
	if err != nil {
		return nil, err
	}
	return &transactionpb.ActionResponse{
		TransactionId: tx.ID,
	}, nil
}

func (h *InternalTransactionHandler) CreateProductTransaction(ctx context.Context, req *transactionpb.CreateProductTransactionRequest) (*transactionpb.TransactionResponse, error) {
	tx, err := h.TransactionUsecase.CreateTransaction(ctx, &tx_domain.Tx{
		FromID:     req.FromId,
		FromOf:     enums.EOwnerUser,
		ToID:       req.ToId,
		ToOf:       enums.EOwnerDealContract,
		Value:      req.Amount,
		Note:       req.Note,
		Method:     enums.EMethodBuy,
		Timestamp:  time.Now(),
		LastAction: h.Mapper.CreateActionRequestToTxAction(req.ActionRequest),
	})
	if err != nil {
		return nil, err
	}
	return &transactionpb.TransactionResponse{
		TransactionId: tx.ID,
	}, nil
}

func (h *InternalTransactionHandler) CreateDepositAction(ctx context.Context, req *transactionpb.CreateActionRequest) (*transactionpb.ActionResponse, error) {
	// Validate request
	if err := h.Validator.ValidateCreateActionRequest(req); err != nil {
		return nil, err
	}

	// Create deposit action
	action := &tx_domain.TxAction{
		TxId:      req.TransactionId,
		FromId:    req.FromId,
		FromOf:    enums.EOwnerType(req.FromOf),
		ToId:      req.ToId,
		ToOf:      enums.EOwnerType(req.ToOf),
		Action:    enums.EActionDeposite,
		Value:     req.Amount,
		Note:      req.Note,
		Timestamp: _utils.TimeNowPtr(),
	}

	createdAction, err := h.TransactionUsecase.CreateAction(ctx, action)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create deposit action: %v", err)
	}

	return &transactionpb.ActionResponse{
		TransactionId: createdAction.TxId,
	}, nil
}

// @bind: internal/usecases.TransactionUsecase
func (h *InternalTransactionHandler) CreateSignAction(ctx context.Context, req *transactionpb.CreateActionRequest) (*transactionpb.ActionResponse, error) {
	// Validate request
	if err := h.Validator.ValidateCreateActionRequest(req); err != nil {
		return nil, err
	}

	// Create sign action
	action := &tx_domain.TxAction{
		TxId:   req.TransactionId,
		FromId: req.FromId,
		FromOf: enums.EOwnerType(req.FromOf),
		ToId:   req.ToId,
		ToOf:   enums.EOwnerType(req.ToOf),
		Action: enums.EActionSign,
		Value:  req.Amount,
		Note:   req.Note,
		// Timestamp: time.Now(),
	}

	createdAction, err := h.TransactionUsecase.CreateAction(ctx, action)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create sign action: %v", err)
	}

	return &transactionpb.ActionResponse{
		TransactionId: createdAction.TxId,
	}, nil
}

// @bind: internal/usecases.TransactionUsecase
func (h *InternalTransactionHandler) CreatePaymentAction(ctx context.Context, req *transactionpb.CreateActionRequest) (*transactionpb.ActionResponse, error) {
	// Validate request
	if err := h.Validator.ValidateCreateActionRequest(req); err != nil {
		return nil, err
	}

	// Create payment action
	action := &tx_domain.TxAction{
		TxId:      req.TransactionId,
		FromId:    req.FromId,
		FromOf:    enums.EOwnerType(req.FromOf),
		ToId:      req.ToId,
		ToOf:      enums.EOwnerType(req.ToOf),
		Action:    enums.EActionPay,
		Value:     req.Amount,
		Note:      req.Note,
		Timestamp: _utils.TimeNowPtr(),
	}

	createdAction, err := h.TransactionUsecase.CreateAction(ctx, action)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create payment action: %v", err)
	}

	return &transactionpb.ActionResponse{
		TransactionId: createdAction.TxId,
	}, nil
}

func (h *InternalTransactionHandler) CreateHandoverAction(ctx context.Context, req *transactionpb.CreateActionRequest) (*transactionpb.ActionResponse, error) {
	// Validate request
	if err := h.Validator.ValidateCreateActionRequest(req); err != nil {
		return nil, err
	}

	// Create handover action
	action := &tx_domain.TxAction{
		TxId:      req.TransactionId,
		FromId:    req.FromId,
		FromOf:    enums.EOwnerType(req.FromOf),
		ToId:      req.ToId,
		ToOf:      enums.EOwnerType(req.ToOf),
		Action:    enums.EActionHandOver,
		Value:     req.Amount,
		Note:      req.Note,
		Timestamp: _utils.TimeNowPtr(),
	}

	createdAction, err := h.TransactionUsecase.CreateAction(ctx, action)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create handover action: %v", err)
	}

	return &transactionpb.ActionResponse{
		TransactionId: createdAction.TxId,
	}, nil
}

func (h *InternalTransactionHandler) GetTransactionDetail(ctx context.Context, req *sharepb.IdRequest) (*transactionpb.TransactionResponse, error) {
	tx, err := h.TransactionUsecase.GetTransaction(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	result := h.Mapper.TxToPb(tx)

	for _, action := range tx.Actions {
		at := h.Mapper.TxActionToResponse(&action)
		result.Actions = append(result.Actions, at)
	}

	return result, nil
}

func (h *InternalTransactionHandler) GetTransactionByIDs(ctx context.Context, req *sharepb.IdRequest) (*transactionpb.TransactionListResponse, error) {
	txs, err := h.TransactionUsecase.GetTransactionByIDs(ctx, req.Ids)
	if err != nil {
		return nil, err
	}
	transactions := make([]*transactionpb.TransactionResponse, len(txs))
	for i, tx := range txs {
		transactions[i] = h.Mapper.TxToPb(tx)
	}
	return &transactionpb.TransactionListResponse{
		Data:  transactions,
		Total: int64(len(txs)),
	}, nil
}

func (h *InternalTransactionHandler) GetDealStatisticByDealIds(ctx context.Context, req *sharepb.IdRequest) (*transactionpb.StatisticDealsResponse, error) {
	statistics, err := h.DealCostRepo.GetStatisticsByDealIDs(ctx, req.Ids)
	if err != nil {
		return nil, err
	}
	return &transactionpb.StatisticDealsResponse{
		DealIds:       req.Ids,
		AmountTxt:     statistics.AmountTxt,
		AmountRevenue: statistics.AmountRevenue,
		AmountCost:    statistics.AmountCost,
		TotalProfit:   statistics.TotalProfit,
	}, nil
}
