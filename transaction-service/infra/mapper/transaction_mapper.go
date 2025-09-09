package mapper

import (
	_utils "common/utils"
	transactionpb "pb/types/transaction"
	"time"
	"tx/infra/utils"
	tx_domain "tx/internal/domain"
	"tx/internal/dto"
	"tx/internal/enums"
)

type TransactionMapper struct{}

func NewTransactionMapper() *TransactionMapper {
	return &TransactionMapper{}
}

func (m *TransactionMapper) CreateActionRequestToTxAction(req *transactionpb.CreateActionRequest) *tx_domain.TxAction {
	if req == nil {
		return nil
	}
	return &tx_domain.TxAction{
		TxId:   req.TransactionId,
		FromId: req.FromId,
		FromOf: enums.EOwnerType(req.FromOf),
		ToId:   req.ToId,
		ToOf:   enums.EOwnerType(req.ToOf),
		Action: enums.EAction(req.Action),
		Value:  req.Amount,
		Note:   req.Note,
		// Timestamp: time.Now(),
	}
}

func (m *TransactionMapper) CreateActionRequestToDepositAction(req *transactionpb.CreateActionRequest) *tx_domain.TxAction {
	return &tx_domain.TxAction{
		TxId:   req.TransactionId,
		FromId: req.FromId,
		FromOf: enums.EOwnerType(req.FromOf),
		ToId:   req.ToId,
		ToOf:   enums.EOwnerType(req.ToOf),
		Action: enums.EActionDeposite,
		Value:  req.Amount,
		Note:   req.Note,
		// Timestamp: time.Now(),
	}
}

func (m *TransactionMapper) CreateActionRequestToSignAction(req *transactionpb.CreateActionRequest) *tx_domain.TxAction {
	return &tx_domain.TxAction{
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
}

func (m *TransactionMapper) CreateActionRequestToPaymentAction(req *transactionpb.CreateActionRequest) *tx_domain.TxAction {
	return &tx_domain.TxAction{
		TxId:   req.TransactionId,
		FromId: req.FromId,
		FromOf: enums.EOwnerType(req.FromOf),
		ToId:   req.ToId,
		ToOf:   enums.EOwnerType(req.ToOf),
		Action: enums.EActionPay,
		Value:  req.Amount,
		Note:   req.Note,
		// Timestamp: req.Timestamp,
	}
}

func (m *TransactionMapper) CreateActionRequestToHandoverAction(req *transactionpb.CreateActionRequest) *tx_domain.TxAction {
	return &tx_domain.TxAction{
		TxId:   req.TransactionId,
		FromId: req.FromId,
		FromOf: enums.EOwnerType(req.FromOf),
		ToId:   req.ToId,
		ToOf:   enums.EOwnerType(req.ToOf),
		Action: enums.EActionHandOver,
		Value:  req.Amount,
		Note:   req.Note,
		// Timestamp: req.Timestamp,
	}
}

func (m *TransactionMapper) CreateProductTransactionRequestToTx(req *transactionpb.CreateProductTransactionRequest) *tx_domain.Tx {
	return &tx_domain.Tx{
		FromID:          req.FromId,
		ToID:            req.ToId,
		Value:           req.Amount,
		Note:            req.Note,
		TransactionName: "Transfer Transaction",
		Method:          enums.EMethodTransfer,
		Status:          enums.TxStatusPending,
		Timestamp:       time.Now(),
	}
}

func (m *TransactionMapper) TxActionToActionResponse(action *tx_domain.TxAction) *transactionpb.ActionResponse {
	return &transactionpb.ActionResponse{
		TransactionId: action.TxId,
		CreatedAt:     action.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (m *TransactionMapper) TxToActionResponse(tx *tx_domain.Tx) *transactionpb.ActionResponse {
	return &transactionpb.ActionResponse{
		TransactionId: tx.ID,
		CreatedAt:     tx.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (m *TransactionMapper) TxToPb(tx *tx_domain.Tx) *transactionpb.TransactionResponse {
	result := &transactionpb.TransactionResponse{
		TransactionId:   tx.ID,
		TransactionName: tx.TransactionName,
		Status:          transactionpb.TransactionStatus(tx.Status),
		Method:          transactionpb.TransactionMethod(tx.Method),
		Value:           tx.Value,
		FromId:          tx.FromID,
		FromOf:          transactionpb.OwnerType(tx.FromOf),
		ToId:            tx.ToID,
		ToOf:            transactionpb.OwnerType(tx.ToOf),
		CreatedAt:       _utils.FormatTimeToString(tx.CreatedAt),
		Actions:         m.TxActionsToPb(tx.Actions),
	}
	if tx.LastAction != nil {
		result.TransactionName = enums.ActionNames[tx.LastAction.Action]
	}
	if result.TransactionName == "" {
		result.TransactionName = enums.MethodNames[tx.Method]
	}
	return result
}

func (m *TransactionMapper) TxActionToResponse(action *tx_domain.TxAction) *transactionpb.ActionResponse {
	return &transactionpb.ActionResponse{
		TransactionId: action.TxId,
		CreatedAt:     _utils.FormatTimeToString(action.Timestamp),
		FromId:        action.FromId,
		FromOf:        transactionpb.OwnerType(action.FromOf),
		ToId:          action.ToId,
		ToOf:          transactionpb.OwnerType(action.ToOf),
		Action:        uint32(action.Action),
		ActionName:    enums.ActionNames[action.Action],
		Value:         action.Value,
		Note:          action.Note,
		Message:       action.Message,
	}
}

func (m *TransactionMapper) TxActionsToPb(actions []tx_domain.TxAction) []*transactionpb.ActionResponse {
	actionsPb := make([]*transactionpb.ActionResponse, len(actions))
	for i, action := range actions {
		actionsPb[i] = m.TxActionToResponse(&action)
	}
	return actionsPb
}

func (m *TransactionMapper) ProcessDTOsToPb(processes []dto.ProcessDTO) []*transactionpb.ContractProcess {
	processesPb := make([]*transactionpb.ContractProcess, len(processes))
	for i, process := range processes {
		processesPb[i] = m.ProcessDTOToPb(&process)
	}
	return processesPb
}
func (m *TransactionMapper) ProcessDTOToPb(process *dto.ProcessDTO) *transactionpb.ContractProcess {
	return &transactionpb.ContractProcess{
		StatusName:    process.StatusName,
		Timestamp:     _utils.FormatTimeToString(process.Timestamp),
		Code:          utils.GenerateTransactionCode(process.TransactionId),
		TransactionId: process.TransactionId,
		Color:         process.Color,
		BgColor:       process.BgColor,
		BorderColor:   process.BorderColor,
	}
}
