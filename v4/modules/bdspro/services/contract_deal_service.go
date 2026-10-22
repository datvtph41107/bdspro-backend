package services

import (
	"context"
	"fmt"
	"time"

	_errors "common/errors"
	_utils "common/utils"
	bdspropb "pb/types/bdspro"

	"mdn4/modules/bdspro/dto"
	"mdn4/modules/bdspro/entity"
	"mdn4/modules/bdspro/enums"
	"mdn4/modules/bdspro/repository"
)

// ContractDealService owns contract-deal writes (reads delegate to DealService in the controller).
type ContractDealService interface {
	GetContract(ctx context.Context, id uint64) (*bdspropb.DealContractResponse, error)
	DeleteContract(ctx context.Context, id uint64) (uint64, error)
	CreateContractPayment(ctx context.Context, req *bdspropb.DealContractPaymentRequest) (*bdspropb.DealContractResponse, error)
	CreateActionPayment(ctx context.Context, req *bdspropb.PaymentActionRequest) (*bdspropb.ActionResponse, error)
	GetMyContracts(ctx context.Context, page, size uint32) (*bdspropb.DealContractListResponse, error)
}

// @bind: mdn4/modules/bdspro/services.ContractDealService
type LocalContractDealService struct {
	sync *_utils.SyncUtil
	repo repository.ContractDealRepository
}

func NewLocalContractDealService(sync *_utils.SyncUtil, repo repository.ContractDealRepository) *LocalContractDealService {
	return &LocalContractDealService{sync: sync, repo: repo}
}

func (s *LocalContractDealService) GetContract(ctx context.Context, id uint64) (*bdspropb.DealContractResponse, error) {
	if id == 0 {
		return nil, _errors.BadRequestException("contract id is required")
	}
	contract, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if contract == nil {
		return nil, _errors.NotFoundException("contract not found")
	}
	return contractToProto(ctx, s.repo, contract)
}

func (s *LocalContractDealService) DeleteContract(ctx context.Context, id uint64) (uint64, error) {
	if id == 0 {
		return 0, _errors.BadRequestException("contract id is required")
	}
	contract, _ := s.repo.GetByID(ctx, id)
	dealID := uint64(0)
	if contract != nil {
		dealID = contract.DealID
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return 0, err
	}
	s.bumpContractsSync(ctx, dealID)
	return id, nil
}

func (s *LocalContractDealService) CreateContractPayment(ctx context.Context, req *bdspropb.DealContractPaymentRequest) (*bdspropb.DealContractResponse, error) {
	if req == nil || req.GetDealId() == 0 {
		return nil, _errors.BadRequestException("deal id is required")
	}
	userID := _utils.GetProfileIdWithContext(ctx)
	if userID == 0 {
		return nil, _errors.UnauthorizedException("user not authenticated")
	}
	contract, err := s.repo.CreateContractPayment(ctx, repository.ContractPaymentInput{
		UserID:          userID,
		DealID:          req.GetDealId(),
		CustomerID:      req.GetCustomerId(),
		ProductID:       req.GetProductId(),
		Amount:          req.GetAmount(),
		Note:            req.GetNote(),
		TransactionType: req.GetTransactionType(),
		TransactionStep: req.GetTransactionStep(),
	})
	if err != nil {
		return nil, _errors.BadRequestException(err.Error())
	}
	s.bumpContractsSync(ctx, req.GetDealId())
	return contractToProto(ctx, s.repo, contract)
}

func (s *LocalContractDealService) CreateActionPayment(ctx context.Context, req *bdspropb.PaymentActionRequest) (*bdspropb.ActionResponse, error) {
	if req == nil || req.GetContractId() == 0 {
		return nil, _errors.BadRequestException("contract id is required")
	}
	var ts *time.Time
	if req.GetTimestamp() != "" {
		ts = _utils.ParseStringToTime(req.GetTimestamp())
	}
	action, err := s.repo.CreateActionPayment(ctx, repository.ActionPaymentInput{
		ContractID: req.GetContractId(),
		FromID:     req.GetFromId(),
		Action:     req.GetAction(),
		Amount:     req.GetAmount(),
		Note:       req.GetNote(),
		Timestamp:  ts,
	})
	if err != nil {
		return nil, _errors.BadRequestException(err.Error())
	}
	if contract, _ := s.repo.GetByID(ctx, req.GetContractId()); contract != nil {
		s.bumpContractsSync(ctx, contract.DealID)
	}
	return actionToProto(action), nil
}

func (s *LocalContractDealService) GetMyContracts(ctx context.Context, page, size uint32) (*bdspropb.DealContractListResponse, error) {
	userID := _utils.GetProfileIdWithContext(ctx)
	if userID == 0 {
		return nil, _errors.UnauthorizedException("user not authenticated")
	}
	rows, total, err := s.repo.ListByProfileID(ctx, userID, page, size)
	if err != nil {
		return nil, err
	}
	return &bdspropb.DealContractListResponse{
		Data:  contractRowsToProto(rows),
		Total: int32(total),
	}, nil
}

func (s *LocalContractDealService) bumpContractsSync(ctx context.Context, dealID uint64) {
	if s.sync == nil || dealID == 0 {
		return
	}
	now := time.Now().UnixMilli()
	_ = s.sync.PutTimestamp(ctx, s.sync.GetKey(ctx, _utils.SyncKeyDealContracts, dealID), now)
}

func contractRowsToProto(rows []dto.TxContractDealListRow) []*bdspropb.DealContractItem {
	out := make([]*bdspropb.DealContractItem, 0, len(rows))
	for _, row := range rows {
		item := &bdspropb.DealContractItem{
			Id:         row.ContractID,
			CustomerId: row.CustomerID,
			CreatedAt:  _utils.FormatTimeToString(row.Timestamp),
			ContractId: row.ContractID,
			Amount:     row.Amount,
			Status:     uint64(row.Status),
			Code:       fmt.Sprintf("TXT%06d", row.ContractID),
		}
		if row.TransactionName == "" {
			item.TransactionName = enums.TxMethodNames[enums.TxMethod(row.Method)]
		} else {
			item.TransactionName = row.TransactionName
		}
		item.StatusName = enums.TxApprovedStatusMap[enums.TxApprovedStatus(row.Status)]
		out = append(out, item)
	}
	return out
}

func contractToProto(ctx context.Context, repo repository.ContractDealRepository, contract *entity.TxContractDeal) (*bdspropb.DealContractResponse, error) {
	if contract == nil {
		return nil, nil
	}
	out := &bdspropb.DealContractResponse{
		ContractId:      contract.ID,
		DealId:          contract.DealID,
		CustomerId:      contract.CustomerID,
		CreatedAt:       _utils.FormatTimeToString(contract.CreatedAt),
		TransactionType: bdspropb.TransactionType(contract.Type),
		TypeName:        enums.TxTransactionTypeMap[enums.TxTransactionType(contract.Type)],
		Note:            contract.Note,
		Code:            fmt.Sprintf("TXT%06d", contract.ID),
	}
	if contract.ProductID != nil {
		out.ProductId = *contract.ProductID
	}
	if contract.TransactionID == nil {
		return out, nil
	}
	txn, sum, err := repo.GetTransaction(ctx, *contract.TransactionID)
	if err != nil {
		return nil, err
	}
	if txn == nil {
		return out, nil
	}
	out.TransactionId = txn.ID
	out.TransactionName = txn.TransactionName
	out.Amount = sum
	out.Status = uint64(txn.Status)
	out.StatusName = enums.TxApprovedStatusMap[enums.TxApprovedStatus(txn.Status)]
	out.Transaction = &bdspropb.TransactionResponse{
		TransactionId:   txn.ID,
		TransactionName: txn.TransactionName,
		Value:           sum,
		Status:          bdspropb.TransactionStatus(txn.Status),
		Method:          bdspropb.TransactionMethod(txn.Method),
	}
	return out, nil
}

func actionToProto(action *entity.TxTransactionAction) *bdspropb.ActionResponse {
	if action == nil {
		return &bdspropb.ActionResponse{}
	}
	return &bdspropb.ActionResponse{
		TransactionId: action.TxID,
		CreatedAt:     _utils.FormatTimeToString(action.Timestamp),
		FromId:        action.FromID,
		FromOf:        bdspropb.OwnerType(action.FromOf),
		ToId:          action.ToID,
		ToOf:          bdspropb.OwnerType(action.ToOf),
		Action:        action.Action,
		Value:         action.Value,
		Note:          action.Note,
		ActionName:    enums.TxActionNames[enums.TxAction(action.Action)],
		Message:       action.Message,
	}
}
