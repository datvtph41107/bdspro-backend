package usecases

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"bdspro/internal/provider"
	"bdspro/internal/repo"
	_utils "common/utils"
	"context"
	"errors"
	"time"
)

var (
	ErrNotFound     = errors.New("transaction not found")
	ErrInvalidState = errors.New("giao dịch đã kết thúc")
)

type TxTransactionUsecase interface {
	GetTransaction(ctx context.Context, id uint64) (*domain.Tx, error)
	UpdateTransaction(ctx context.Context, tx *domain.Tx) (*domain.Tx, error)
	DeleteTransaction(ctx context.Context, id uint64) error
	ListTransactions(ctx context.Context, filters *dto.TxTransactionFilters) ([]*domain.Tx, int64, error)
	ApproveTransaction(ctx context.Context, id uint64) (*domain.Tx, error)
	RejectTransaction(ctx context.Context, id uint64) error
	GetTransactionTypes(ctx context.Context) ([]string, error)
	GetTransactionByIDs(ctx context.Context, ids []uint64) ([]*domain.Tx, error)
	GetMessageByAction(action *domain.TxAction) string

	// method
	CreateTransaction(ctx context.Context, tx *domain.Tx) (*domain.Tx, error)
	CreateAction(ctx context.Context, action *domain.TxAction) (*domain.TxAction, error)
}

type transactionUsecase struct {
	transactionRepo repo.TxTransactionRepository
	actionRepo      repo.TxActionRepository
	Transaction     provider.TransactionProvider
}

func NewTransactionUsecase(
	transactionRepo repo.TxTransactionRepository,
	actionRepo repo.TxActionRepository,
	transaction provider.TransactionProvider,
) TxTransactionUsecase {
	return &transactionUsecase{
		transactionRepo: transactionRepo,
		actionRepo:      actionRepo,
		Transaction:     transaction,
	}
}

func (uc *transactionUsecase) GetMessageByAction(action *domain.TxAction) string {
	switch action.Action {
	case enums.TxActionDeposite:
		return "Đặt cọc " + _utils.FormatVND(action.Value)
	case enums.TxActionSign:
		return "Ký hợp đồng " + _utils.FormatVND(action.Value)
	case enums.TxActionPay:
		return "Thanh toán " + _utils.FormatVND(action.Value)
	case enums.TxActionHandOver:
		return "Bàn giao " + _utils.FormatVND(action.Value)
	case enums.TxActionCancel:
		return "Hủy"
	}
	return ""
}

func (uc *transactionUsecase) CreateTransaction(ctx context.Context, tx *domain.Tx) (*domain.Tx, error) {
	if !enums.IsValidTxMethod(tx.Method) {
		return nil, errors.New("invalid method")
	}

	// Approve với ký và thu nhập
	tx.Status = enums.TxStatusPending
	if tx.Method == enums.TxMethodSign || tx.Method == enums.TxMethodRevenue {
		tx.Status = enums.TxStatusDone
	}

	err := uc.Transaction.WithTransaction(ctx, func(ctx context.Context) error {
		lastAction := tx.LastAction
		tx.LastAction = nil
		_, err := uc.transactionRepo.Create(ctx, tx)
		if err != nil {
			return err
		}
		action := lastAction
		if lastAction == nil {
			action = &domain.TxAction{
				FromId: tx.FromID,
				FromOf: tx.FromOf,
				ToId:   tx.ToID,
				ToOf:   tx.ToOf,
				Action: enums.TxActionDeposite,
				Note:   tx.Note,
				// Timestamp: time.Now(),
			}
		}
		action.TxId = tx.ID
		action.Message = uc.GetMessageByAction(action)
		action, err = uc.transactionRepo.CreateAction(ctx, action)
		if err != nil {
			return err
		}
		tx.LastActionID = &action.ID
		_, err = uc.transactionRepo.Update(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})

	return tx, err
}

func (uc *transactionUsecase) GetTransaction(ctx context.Context, id uint64) (*domain.Tx, error) {
	tx, err := uc.transactionRepo.GetByID(ctx, &id)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		return nil, ErrNotFound
	}
	tx.Value, err = uc.transactionRepo.SumByTransactionID(ctx, id)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (uc *transactionUsecase) UpdateTransaction(ctx context.Context, tx *domain.Tx) (*domain.Tx, error) {
	// Get existing transaction
	existingTx, err := uc.transactionRepo.GetByID(ctx, &tx.ID)
	if err != nil {
		return nil, err
	}
	if existingTx == nil {
		return nil, ErrNotFound
	}

	// Update only allowed fields
	existingTx.Value = tx.Value
	existingTx.TransactionName = tx.TransactionName
	existingTx.Method = tx.Method
	existingTx.Status = tx.Status
	existingTx.Timestamp = tx.Timestamp

	entity, err := uc.transactionRepo.Update(ctx, existingTx)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (uc *transactionUsecase) DeleteTransaction(ctx context.Context, id uint64) error {
	// Check if transaction exists
	existingTx, err := uc.transactionRepo.GetByID(ctx, &id)
	if err != nil {
		return err
	}
	if existingTx == nil {
		return ErrNotFound
	}

	return uc.transactionRepo.Delete(ctx, id)
}

func (uc *transactionUsecase) ListTransactions(ctx context.Context, filters *dto.TxTransactionFilters) ([]*domain.Tx, int64, error) {
	return uc.transactionRepo.List(ctx, filters)
}

func (uc *transactionUsecase) ApproveTransaction(ctx context.Context, id uint64) (*domain.Tx, error) {
	// Get existing transaction
	tx, err := uc.transactionRepo.GetByID(ctx, &id)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		return nil, ErrNotFound
	}

	// Check if transaction can be approved
	if tx.Status != enums.TxStatusPending {
		return nil, ErrInvalidState
	}

	err = uc.Transaction.WithTransaction(ctx, func(ctx context.Context) error {
		currentUserId := _utils.GetProfileIdWithContext(ctx)
		// action := &domain.TxAction{
		// 	TxId:   id,
		// 	FromId: currentUserId,
		// 	FromOf: enums.EOwnerUser,
		// 	ToId:   tx.ToID,
		// 	ToOf:   tx.ToOf,
		// 	Action: enums.EActionHandOver,
		// 	Note:   tx.Note,
		// 	// Timestamp: time.Now(),
		// }
		// _, err = uc.actionRepo.Approve(ctx, action)
		// if err != nil {
		// 	return err
		// }
		err = uc.transactionRepo.Approve(ctx, id, currentUserId)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (uc *transactionUsecase) RejectTransaction(ctx context.Context, id uint64) error {
	// Get existing transaction
	tx, err := uc.transactionRepo.GetByID(ctx, &id)
	if err != nil {
		return err
	}
	if tx == nil {
		return ErrNotFound
	}

	// Check if transaction can be rejected
	if tx.Status != enums.TxStatusPending {
		return ErrInvalidState
	}

	err = uc.transactionRepo.Reject(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (uc *transactionUsecase) GetTransactionTypes(ctx context.Context) ([]string, error) {
	return uc.transactionRepo.GetTransactionTypes(ctx)
}

func (uc *transactionUsecase) CreateAction(ctx context.Context, action *domain.TxAction) (*domain.TxAction, error) {
	tx, err := uc.transactionRepo.GetByID(ctx, &action.TxId)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		return nil, ErrNotFound
	}
	if tx.Status == enums.TxStatusCancel || tx.Status == enums.TxStatusDone {
		return nil, ErrInvalidState
	}
	// action.Timestamp = time.Now()
	err = uc.Transaction.WithTransaction(ctx, func(ctx context.Context) error {
		_, err = uc.transactionRepo.CreateAction(ctx, action)
		if err != nil {
			return err
		}
		tx.LastActionID = &action.ID
		_, err = uc.transactionRepo.Update(ctx, tx)
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return action, nil
}

func (uc *transactionUsecase) Approve(ctx context.Context, action *domain.TxAction) (*domain.TxAction, error) {
	tx, err := uc.transactionRepo.GetByID(ctx, &action.TxId)
	if err != nil {
		return nil, err
	}
	if tx == nil {
		return nil, ErrNotFound
	}
	// action.Timestamp = time.Now()
	action.Action = enums.TxAction(enums.TxMethodApprove)
	action.Note = tx.Note

	uc.Transaction.WithTransaction(ctx, func(ctx context.Context) error {
		_, err := uc.actionRepo.Approve(ctx, action)
		if err != nil {
			return err
		}

		if tx.Method == enums.TxMethodSign {
			tx.Status = enums.TxStatusDone
		}

		_, err = uc.transactionRepo.Update(ctx, tx)
		if err != nil {
			return err
		}

		return nil
	})
	return action, nil
}

func (uc *transactionUsecase) Reject(ctx context.Context, action *domain.TxAction) (*domain.TxAction, error) {
	return nil, nil
}

func (uc *transactionUsecase) Sign(ctx context.Context, action *domain.TxAction) (*domain.TxAction, error) {
	return nil, nil
}

func (uc *transactionUsecase) CreateSignTransaction(ctx context.Context, tx *domain.Tx) (*domain.Tx, error) {
	uc.Transaction.WithTransaction(ctx, func(ctx context.Context) error {
		tx.Method = enums.TxMethodSign
		tx.Status = enums.TxStatusPending
		_, err := uc.transactionRepo.Create(ctx, tx)
		if err != nil {
			return err
		}
		_, err = uc.actionRepo.Sign(ctx, &domain.TxAction{
			TxId:   tx.ID,
			FromId: tx.FromID,
			FromOf: enums.TxOwnerType(tx.FromOf),
			Action: enums.TxActionSign,
			Value:  tx.Value,
			Note:   tx.Note,
			// Timestamp: time.Now(),
		})
		if err != nil {
			return err
		}
		return nil
	})
	return tx, nil
}

func (uc *transactionUsecase) CreateTransferTransaction(ctx context.Context, req *domain.TransferTransactionRequest) (*domain.Tx, error) {
	// Create transaction with transfer method
	tx := &domain.Tx{
		TransactionName: "Chuyển tiền",
		FromID:          req.FromID,
		FromOf:          enums.TxOwnerUser, // User type
		ToID:            req.ToID,
		ToOf:            enums.TxOwnerDealContract, // Deal type
		Method:          enums.TxMethodTransfer,
		Value:           req.Amount,
		Status:          enums.TxStatusPending,
		Timestamp:       time.Now(),
		Note:            req.Note,
	}

	uc.Transaction.WithTransaction(ctx, func(ctx context.Context) error {
		_, err := uc.transactionRepo.Create(ctx, tx)
		if err != nil {
			return err
		}

		action := &domain.TxAction{
			TxId:   tx.ID,
			FromId: req.FromID,
			FromOf: enums.TxOwnerUser,
			ToId:   req.ToID,
			ToOf:   enums.TxOwnerDealContract,
			Action: enums.TxActionDeposite,
			Value:  req.Amount,
			Note:   req.Note,
			// Timestamp: time.Now(),
		}
		// Create transfer action
		_, err = uc.actionRepo.Transfer(ctx, action)
		if err != nil {
			return err
		}

		tx.LastActionID = &action.ID
		_, err = uc.transactionRepo.Update(ctx, tx)
		if err != nil {
			return err
		}
		return nil
	})
	return tx, nil
}

func (uc *transactionUsecase) GetTransactionByIDs(ctx context.Context, ids []uint64) ([]*domain.Tx, error) {
	return uc.transactionRepo.GetByIDs(ctx, ids)
}
