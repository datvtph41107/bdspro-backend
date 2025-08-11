package mapper

import (
	_utils "common/utils"
	transactionpb "pb/types/transaction"
	"time"
	tx_domain "tx/internal/domain"
	"tx/internal/enums"
)

type TransactionMapper struct{}

func NewTransactionMapper() *TransactionMapper {
	return &TransactionMapper{}
}

func (m *TransactionMapper) CreateActionRequestToTxAction(req *transactionpb.CreateActionRequest) *tx_domain.TxAction {
	return &tx_domain.TxAction{
		TxId:      req.TransactionId,
		FromId:    req.FromId,
		FromOf:    enums.EOwnerType(req.FromOf),
		ToId:      req.ToId,
		ToOf:      enums.EOwnerType(req.ToOf),
		Action:    enums.EAction(req.Action),
		Value:     req.Amount,
		Note:      req.Note,
		Timestamp: time.Now(),
	}
}

func (m *TransactionMapper) CreateActionRequestToDepositAction(req *transactionpb.CreateActionRequest) *tx_domain.TxAction {
	return &tx_domain.TxAction{
		TxId:      req.TransactionId,
		FromId:    req.FromId,
		FromOf:    enums.EOwnerType(req.FromOf),
		ToId:      req.ToId,
		ToOf:      enums.EOwnerType(req.ToOf),
		Action:    enums.EActionDeposite,
		Value:     req.Amount,
		Note:      req.Note,
		Timestamp: time.Now(),
	}
}

func (m *TransactionMapper) CreateActionRequestToSignAction(req *transactionpb.CreateActionRequest) *tx_domain.TxAction {
	return &tx_domain.TxAction{
		TxId:      req.TransactionId,
		FromId:    req.FromId,
		FromOf:    enums.EOwnerType(req.FromOf),
		ToId:      req.ToId,
		ToOf:      enums.EOwnerType(req.ToOf),
		Action:    enums.EActionSign,
		Value:     req.Amount,
		Note:      req.Note,
		Timestamp: time.Now(),
	}
}

func (m *TransactionMapper) CreateActionRequestToPaymentAction(req *transactionpb.CreateActionRequest) *tx_domain.TxAction {
	return &tx_domain.TxAction{
		TxId:      req.TransactionId,
		FromId:    req.FromId,
		FromOf:    enums.EOwnerType(req.FromOf),
		ToId:      req.ToId,
		ToOf:      enums.EOwnerType(req.ToOf),
		Action:    enums.EActionPay,
		Value:     req.Amount,
		Note:      req.Note,
		Timestamp: time.Now(),
	}
}

func (m *TransactionMapper) CreateActionRequestToHandoverAction(req *transactionpb.CreateActionRequest) *tx_domain.TxAction {
	return &tx_domain.TxAction{
		TxId:      req.TransactionId,
		FromId:    req.FromId,
		FromOf:    enums.EOwnerType(req.FromOf),
		ToId:      req.ToId,
		ToOf:      enums.EOwnerType(req.ToOf),
		Action:    enums.EActionHandOver,
		Value:     req.Amount,
		Note:      req.Note,
		Timestamp: time.Now(),
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
		Timestamp:     action.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (m *TransactionMapper) TxToActionResponse(tx *tx_domain.Tx) *transactionpb.ActionResponse {
	return &transactionpb.ActionResponse{
		TransactionId: tx.ID,
		Timestamp:     tx.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (m *TransactionMapper) TxToTransactionResponse(tx *tx_domain.Tx) *transactionpb.TransactionResponse {
	result := &transactionpb.TransactionResponse{
		TransactionId:   tx.ID,
		TransactionName: tx.TransactionName,
		Status:          transactionpb.TransactionStatus(tx.Status),
		Method:          transactionpb.TransactionMethod(tx.Method),
		Value:           tx.Value,
		CreatedAt:       _utils.FormatTimeToString(tx.CreatedAt),
	}
	if tx.LastAction != nil {
		result.TransactionName = enums.ActionNames[tx.LastAction.Action]
	}
	return result
}
