package handler

import (
	"bdspro/infra/mapper"
	"bdspro/infra/validator"
	"bdspro/internal/domain"
	"bdspro/internal/enums"
	"bdspro/internal/repo"
	usecase "bdspro/internal/usecases"
	_utils "common/utils"
	"context"
	bdspropb "pb/types/bdspro"
	sharepb "pb/types/shared"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type InternalTransactionHandler struct {
	bdspropb.UnimplementedInternalTransactionServiceServer
	TransactionUsecase  usecase.TxTransactionUsecase
	DealContractUsecase *usecase.DealContractUsecase
	Validator           validator.TransactionActionValidator
	Mapper              *mapper.TxMapper
	DealMapper          *mapper.ContractDealMapper
	DealCostRepo        repo.TxDealCostRepository
}

func NewInternalTransactionHandler(
	transactionUsecase usecase.TxTransactionUsecase,
	dealContractUsecase *usecase.DealContractUsecase,
	validator validator.TransactionActionValidator,
	mapper *mapper.TxMapper,
	dealMapper *mapper.ContractDealMapper,
	dealCostRepo repo.TxDealCostRepository,
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
func (h *InternalTransactionHandler) CreateAction(ctx context.Context, req *bdspropb.CreateActionRequest) (*bdspropb.ActionResponse, error) {
	// Validate request
	if err := h.Validator.ValidateCreateActionRequest(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	tx, err := h.TransactionUsecase.CreateAction(ctx, &domain.TxAction{
		FromId: req.FromId,
		TxId:   req.TransactionId,
		FromOf: enums.TxOwnerType(req.FromOf),
		ToId:   req.ToId,
		ToOf:   enums.TxOwnerType(req.ToOf),
		Action: enums.TxAction(req.Action),
		Value:  req.Amount,
		Note:   req.Note,
	})
	if err != nil {
		return nil, err
	}
	return &bdspropb.ActionResponse{
		TransactionId: tx.ID,
	}, nil
}

func (h *InternalTransactionHandler) CreateProductTransaction(ctx context.Context, req *bdspropb.CreateProductTransactionRequest) (*bdspropb.TransactionResponse, error) {
	tx, err := h.TransactionUsecase.CreateTransaction(ctx, &domain.Tx{
		FromID:     req.FromId,
		FromOf:     enums.TxOwnerUser,
		ToID:       req.ToId,
		ToOf:       enums.TxOwnerDealContract,
		Value:      req.Amount,
		Note:       req.Note,
		Method:     enums.TxMethodBuy,
		Timestamp:  time.Now(),
		LastAction: h.Mapper.CreateActionRequestToTxAction(req.ActionRequest),
	})
	if err != nil {
		return nil, err
	}
	return &bdspropb.TransactionResponse{
		TransactionId: tx.ID,
	}, nil
}

func (h *InternalTransactionHandler) CreateDepositAction(ctx context.Context, req *bdspropb.CreateActionRequest) (*bdspropb.ActionResponse, error) {
	// Validate request
	if err := h.Validator.ValidateCreateActionRequest(req); err != nil {
		return nil, err
	}

	// Create deposit action
	action := &domain.TxAction{
		TxId:      req.TransactionId,
		FromId:    req.FromId,
		FromOf:    enums.TxOwnerType(req.FromOf),
		ToId:      req.ToId,
		ToOf:      enums.TxOwnerType(req.ToOf),
		Action:    enums.TxActionDeposite,
		Value:     req.Amount,
		Note:      req.Note,
		Timestamp: _utils.TimeNowPtr(),
	}

	createdAction, err := h.TransactionUsecase.CreateAction(ctx, action)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create deposit action: %v", err)
	}

	return &bdspropb.ActionResponse{
		TransactionId: createdAction.TxId,
	}, nil
}

// @bind: internal/usecases.TransactionUsecase
func (h *InternalTransactionHandler) CreateSignAction(ctx context.Context, req *bdspropb.CreateActionRequest) (*bdspropb.ActionResponse, error) {
	// Validate request
	if err := h.Validator.ValidateCreateActionRequest(req); err != nil {
		return nil, err
	}

	// Create sign action
	action := &domain.TxAction{
		TxId:   req.TransactionId,
		FromId: req.FromId,
		FromOf: enums.TxOwnerType(req.FromOf),
		ToId:   req.ToId,
		ToOf:   enums.TxOwnerType(req.ToOf),
		Action: enums.TxActionSign,
		Value:  req.Amount,
		Note:   req.Note,
		// Timestamp: time.Now(),
	}

	createdAction, err := h.TransactionUsecase.CreateAction(ctx, action)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create sign action: %v", err)
	}

	return &bdspropb.ActionResponse{
		TransactionId: createdAction.TxId,
	}, nil
}

// @bind: internal/usecases.TransactionUsecase
func (h *InternalTransactionHandler) CreatePaymentAction(ctx context.Context, req *bdspropb.CreateActionRequest) (*bdspropb.ActionResponse, error) {
	// Validate request
	if err := h.Validator.ValidateCreateActionRequest(req); err != nil {
		return nil, err
	}

	// Create payment action
	action := &domain.TxAction{
		TxId:      req.TransactionId,
		FromId:    req.FromId,
		FromOf:    enums.TxOwnerType(req.FromOf),
		ToId:      req.ToId,
		ToOf:      enums.TxOwnerType(req.ToOf),
		Action:    enums.TxActionPay,
		Value:     req.Amount,
		Note:      req.Note,
		Timestamp: _utils.TimeNowPtr(),
	}

	createdAction, err := h.TransactionUsecase.CreateAction(ctx, action)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create payment action: %v", err)
	}

	return &bdspropb.ActionResponse{
		TransactionId: createdAction.TxId,
	}, nil
}

func (h *InternalTransactionHandler) CreateHandoverAction(ctx context.Context, req *bdspropb.CreateActionRequest) (*bdspropb.ActionResponse, error) {
	// Validate request
	if err := h.Validator.ValidateCreateActionRequest(req); err != nil {
		return nil, err
	}

	// Create handover action
	action := &domain.TxAction{
		TxId:      req.TransactionId,
		FromId:    req.FromId,
		FromOf:    enums.TxOwnerType(req.FromOf),
		ToId:      req.ToId,
		ToOf:      enums.TxOwnerType(req.ToOf),
		Action:    enums.TxActionHandOver,
		Value:     req.Amount,
		Note:      req.Note,
		Timestamp: _utils.TimeNowPtr(),
	}

	createdAction, err := h.TransactionUsecase.CreateAction(ctx, action)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create handover action: %v", err)
	}

	return &bdspropb.ActionResponse{
		TransactionId: createdAction.TxId,
	}, nil
}

func (h *InternalTransactionHandler) GetTransactionDetail(ctx context.Context, req *sharepb.IdRequest) (*bdspropb.TransactionResponse, error) {
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

func (h *InternalTransactionHandler) GetTransactionByIDs(ctx context.Context, req *sharepb.IdRequest) (*bdspropb.TransactionListResponse, error) {
	txs, err := h.TransactionUsecase.GetTransactionByIDs(ctx, req.Ids)
	if err != nil {
		return nil, err
	}
	transactions := make([]*bdspropb.TransactionResponse, len(txs))
	for i, tx := range txs {
		transactions[i] = h.Mapper.TxToPb(tx)
	}
	return &bdspropb.TransactionListResponse{
		Data:  transactions,
		Total: int64(len(txs)),
	}, nil
}

func (h *InternalTransactionHandler) GetDealStatisticByDealIds(ctx context.Context, req *sharepb.IdRequest) (*bdspropb.StatisticDealsResponse, error) {
	statistics, err := h.DealCostRepo.GetStatisticsByDealIDs(ctx, req.Ids)
	if err != nil {
		return nil, err
	}
	return &bdspropb.StatisticDealsResponse{
		DealIds:       req.Ids,
		AmountTxt:     statistics.AmountTxt,
		AmountRevenue: statistics.AmountRevenue,
		AmountCost:    statistics.AmountCost,
		TotalProfit:   statistics.TotalProfit,
	}, nil
}
