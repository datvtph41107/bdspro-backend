package mapper

import (
	tx_domain "bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"bdspro/internal/utils"
	_utils "common/utils"
	bdspropb "pb/types/bdspro"
	"time"
)

type TxMapper struct{}

func NewTxMapper() *TxMapper {
	return &TxMapper{}
}

func (m *TxMapper) CreateActionRequestToTxAction(req *bdspropb.CreateActionRequest) *tx_domain.TxAction {
	if req == nil {
		return nil
	}
	return &tx_domain.TxAction{
		TxId:   req.TransactionId,
		FromId: req.FromId,
		FromOf: enums.TxOwnerType(req.FromOf),
		ToId:   req.ToId,
		ToOf:   enums.TxOwnerType(req.ToOf),
		Action: enums.TxAction(req.Action),
		Value:  req.Amount,
		Note:   req.Note,
		// Timestamp: time.Now(),
	}
}

func (m *TxMapper) CreateActionRequestToDepositAction(req *bdspropb.CreateActionRequest) *tx_domain.TxAction {
	return &tx_domain.TxAction{
		TxId:   req.TransactionId,
		FromId: req.FromId,
		FromOf: enums.TxOwnerType(req.FromOf),
		ToId:   req.ToId,
		ToOf:   enums.TxOwnerType(req.ToOf),
		Action: enums.TxActionDeposite,
		Value:  req.Amount,
		Note:   req.Note,
		// Timestamp: time.Now(),
	}
}

func (m *TxMapper) CreateActionRequestToSignAction(req *bdspropb.CreateActionRequest) *tx_domain.TxAction {
	return &tx_domain.TxAction{
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
}

func (m *TxMapper) CreateActionRequestToPaymentAction(req *bdspropb.CreateActionRequest) *tx_domain.TxAction {
	return &tx_domain.TxAction{
		TxId:   req.TransactionId,
		FromId: req.FromId,
		FromOf: enums.TxOwnerType(req.FromOf),
		ToId:   req.ToId,
		ToOf:   enums.TxOwnerType(req.ToOf),
		Action: enums.TxActionPay,
		Value:  req.Amount,
		Note:   req.Note,
		// Timestamp: req.Timestamp,
	}
}

func (m *TxMapper) CreateActionRequestToHandoverAction(req *bdspropb.CreateActionRequest) *tx_domain.TxAction {
	return &tx_domain.TxAction{
		TxId:   req.TransactionId,
		FromId: req.FromId,
		FromOf: enums.TxOwnerType(req.FromOf),
		ToId:   req.ToId,
		ToOf:   enums.TxOwnerType(req.ToOf),
		Action: enums.TxActionHandOver,
		Value:  req.Amount,
		Note:   req.Note,
		// Timestamp: req.Timestamp,
	}
}

func (m *TxMapper) CreateProductTransactionRequestToTx(req *bdspropb.CreateProductTransactionRequest) *tx_domain.Tx {
	return &tx_domain.Tx{
		FromID:          req.FromId,
		ToID:            req.ToId,
		Value:           req.Amount,
		Note:            req.Note,
		TransactionName: "Transfer Transaction",
		Method:          enums.TxMethodTransfer,
		Status:          enums.TxStatusPending,
		Timestamp:       time.Now(),
	}
}

func (m *TxMapper) TxActionToActionResponse(action *tx_domain.TxAction) *bdspropb.ActionResponse {
	return &bdspropb.ActionResponse{
		TransactionId: action.TxId,
		CreatedAt:     action.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (m *TxMapper) TxToActionResponse(tx *tx_domain.Tx) *bdspropb.ActionResponse {
	return &bdspropb.ActionResponse{
		TransactionId: tx.ID,
		CreatedAt:     tx.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (m *TxMapper) TxToPb(tx *tx_domain.Tx) *bdspropb.TransactionResponse {
	result := &bdspropb.TransactionResponse{
		TransactionId:   tx.ID,
		TransactionName: tx.TransactionName,
		Status:          bdspropb.TransactionStatus(tx.Status),
		Method:          bdspropb.TransactionMethod(tx.Method),
		Value:           tx.Value,
		FromId:          tx.FromID,
		FromOf:          bdspropb.OwnerType(tx.FromOf),
		ToId:            tx.ToID,
		ToOf:            bdspropb.OwnerType(tx.ToOf),
		CreatedAt:       _utils.FormatTimeToString(tx.CreatedAt),
		Actions:         m.TxActionsToPb(tx.Actions),
	}
	if tx.LastAction != nil {
		result.TransactionName = enums.TxActionNames[tx.LastAction.Action]
	}
	if result.TransactionName == "" {
		result.TransactionName = enums.TxMethodNames[tx.Method]
	}
	return result
}

func (m *TxMapper) TxActionToResponse(action *tx_domain.TxAction) *bdspropb.ActionResponse {
	return &bdspropb.ActionResponse{
		TransactionId: action.TxId,
		CreatedAt:     _utils.FormatTimeToString(action.Timestamp),
		FromId:        action.FromId,
		FromOf:        bdspropb.OwnerType(action.FromOf),
		ToId:          action.ToId,
		ToOf:          bdspropb.OwnerType(action.ToOf),
		Action:        uint32(action.Action),
		ActionName:    enums.TxActionNames[action.Action],
		Value:         action.Value,
		Note:          action.Note,
		Message:       action.Message,
	}
}

func (m *TxMapper) TxActionsToPb(actions []tx_domain.TxAction) []*bdspropb.ActionResponse {
	actionsPb := make([]*bdspropb.ActionResponse, len(actions))
	for i, action := range actions {
		actionsPb[i] = m.TxActionToResponse(&action)
	}
	return actionsPb
}

func (m *TxMapper) ProcessDTOsToPb(processes []dto.TxProcessDTO) []*bdspropb.ContractProcess {
	processesPb := make([]*bdspropb.ContractProcess, len(processes))
	for i, process := range processes {
		processesPb[i] = m.ProcessDTOToPb(&process)
	}
	return processesPb
}
func (m *TxMapper) ProcessDTOToPb(process *dto.TxProcessDTO) *bdspropb.ContractProcess {
	return &bdspropb.ContractProcess{
		StatusName:    process.StatusName,
		Timestamp:     _utils.FormatTimeToString(process.Timestamp),
		Code:          utils.GenerateTransactionCode(process.TransactionId),
		TransactionId: process.TransactionId,
		Color:         process.Color,
		BgColor:       process.BgColor,
		BorderColor:   process.BorderColor,
	}
}
