package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"userv4/modules/bdspro/dto"
	"userv4/modules/bdspro/entity"
	"userv4/modules/bdspro/enums"

	"gorm.io/gorm"
)

// ContractDealRepository is tx_contract_deals + tx_transaction writes on the v4 DSN.
type ContractDealRepository interface {
	GetByID(ctx context.Context, id uint64) (*entity.TxContractDeal, error)
	Delete(ctx context.Context, id uint64) error
	ListByProfileID(ctx context.Context, profileID uint64, page, size uint32) ([]dto.TxContractDealListRow, int64, error)
	GetTransaction(ctx context.Context, id uint64) (*entity.TxTransaction, float64, error)
	CreateContractPayment(ctx context.Context, in ContractPaymentInput) (*entity.TxContractDeal, error)
	CreateActionPayment(ctx context.Context, in ActionPaymentInput) (*entity.TxTransactionAction, error)
	EnsureDealProduct(ctx context.Context, dealID, productID uint64) error
}

type ContractPaymentInput struct {
	UserID          uint64
	DealID          uint64
	CustomerID      uint64
	ProductID       uint64
	Amount          float64
	Note            string
	TransactionType uint32
	TransactionStep uint32
}

type ActionPaymentInput struct {
	ContractID uint64
	FromID     uint64
	Action     uint32
	Amount     float64
	Note       string
	Timestamp  *time.Time
}

type contractDealPostgres struct {
	db *gorm.DB
}

func NewContractDealPostgres(db *gorm.DB) ContractDealRepository {
	return &contractDealPostgres{db: db}
}

func (r *contractDealPostgres) GetByID(ctx context.Context, id uint64) (*entity.TxContractDeal, error) {
	if r.db == nil || id == 0 {
		return nil, nil
	}
	var row entity.TxContractDeal
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *contractDealPostgres) Delete(ctx context.Context, id uint64) error {
	if r.db == nil || id == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var contract entity.TxContractDeal
		if err := tx.Where("id = ? AND deleted_at IS NULL", id).First(&contract).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return err
		}
		now := time.Now()
		if err := tx.Model(&entity.TxContractDeal{}).Where("id = ?", id).Update("deleted_at", now).Error; err != nil {
			return err
		}
		if contract.TransactionID != nil {
			if err := tx.Model(&entity.TxTransaction{}).Where("id = ?", *contract.TransactionID).Update("deleted_at", now).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *contractDealPostgres) ListByProfileID(ctx context.Context, profileID uint64, page, size uint32) ([]dto.TxContractDealListRow, int64, error) {
	if r.db == nil || profileID == 0 {
		return nil, 0, nil
	}
	limit := int(size)
	if limit <= 0 {
		limit = 20
	}
	offset := int(page) * limit

	q := r.db.WithContext(ctx).
		Table("tx_contract_deals").
		Select(`
			tx_contract_deals.id as contract_id,
			tx_contract_deals.deal_id,
			txs.transaction_name AS transaction_name,
			txs.method as method,
			tx_contract_deals.customer_id AS customer_id,
			txs.status as status,
			COALESCE(SUM(tx_actions.value), 0) AS amount,
			txs.timestamp
		`).
		Joins(`LEFT JOIN tx_transaction txs ON txs.id = tx_contract_deals.transaction_id`).
		Joins(`LEFT JOIN tx_transaction_action tx_actions ON tx_actions.tx_id = txs.id AND tx_actions.deleted_at IS NULL`).
		Where(`txs.from_id = ? AND tx_contract_deals.deleted_at IS NULL AND txs.deleted_at IS NULL`, profileID).
		Group("tx_contract_deals.id, tx_contract_deals.deal_id, tx_contract_deals.created_at, tx_contract_deals.customer_id, txs.method, txs.transaction_name, txs.status, txs.timestamp").
		Order("tx_contract_deals.created_at DESC")

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []dto.TxContractDealListRow
	if err := q.Offset(offset).Limit(limit).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *contractDealPostgres) GetTransaction(ctx context.Context, id uint64) (*entity.TxTransaction, float64, error) {
	if r.db == nil || id == 0 {
		return nil, 0, nil
	}
	var txn entity.TxTransaction
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&txn).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, 0, nil
	}
	if err != nil {
		return nil, 0, err
	}
	var total float64
	err = r.db.WithContext(ctx).Table("tx_transaction_action").
		Select("COALESCE(SUM(value), 0)").
		Where("tx_id = ? AND deleted_at IS NULL", id).
		Scan(&total).Error
	return &txn, total, err
}

func (r *contractDealPostgres) CreateContractPayment(ctx context.Context, in ContractPaymentInput) (*entity.TxContractDeal, error) {
	if r.db == nil || in.DealID == 0 {
		return nil, fmt.Errorf("deal id is required")
	}
	method := enums.TxMethod(in.TransactionType)
	if !enums.IsValidTxMethod(method) {
		return nil, fmt.Errorf("invalid transaction type")
	}

	var out entity.TxContractDeal
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		now := time.Now()
		var productID *uint64
		if in.ProductID > 0 {
			productID = &in.ProductID
		}
		contract := entity.TxContractDeal{
			DealID:     in.DealID,
			Type:       in.TransactionType,
			Note:       in.Note,
			ProductID:  productID,
			CustomerID: in.CustomerID,
		}
		if in.UserID > 0 {
			contract.CreatedBy = &in.UserID
		}
		if err := tx.Create(&contract).Error; err != nil {
			return err
		}

		status := uint32(enums.TxStatusPending)
		if method == enums.TxMethodSign || method == enums.TxMethodRevenue {
			status = uint32(enums.TxStatusDone)
		}
		txn := entity.TxTransaction{
			TransactionName: enums.TxMethodNames[method],
			FromID:          in.UserID,
			FromOf:          uint32(enums.TxOwnerUser),
			ToID:            in.DealID,
			ToOf:            uint32(enums.TxOwnerDealContract),
			Method:          uint32(method),
			Status:          status,
			Timestamp:       &now,
		}
		if in.UserID > 0 {
			txn.CreatedBy = &in.UserID
		}
		if err := tx.Create(&txn).Error; err != nil {
			return err
		}

		action := uint32(enums.TxActionDeposite)
		if in.TransactionStep > 0 {
			action = in.TransactionStep
		}
		txAction := entity.TxTransactionAction{
			TxID:      txn.ID,
			FromID:    in.UserID,
			FromOf:    uint32(enums.TxOwnerUser),
			ToID:      in.DealID,
			ToOf:      uint32(enums.TxOwnerDealContract),
			Action:    action,
			Value:     in.Amount,
			Note:      in.Note,
			Message:   enums.TxActionNames[enums.TxAction(action)],
			Timestamp: &now,
		}
		if in.UserID > 0 {
			txAction.CreatedBy = &in.UserID
		}
		if err := tx.Create(&txAction).Error; err != nil {
			return err
		}
		if err := tx.Model(&txn).Update("last_action_id", txAction.ID).Error; err != nil {
			return err
		}

		txnID := txn.ID
		contract.TransactionID = &txnID
		if err := tx.Save(&contract).Error; err != nil {
			return err
		}
		out = contract
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *contractDealPostgres) CreateActionPayment(ctx context.Context, in ActionPaymentInput) (*entity.TxTransactionAction, error) {
	if r.db == nil || in.ContractID == 0 {
		return nil, fmt.Errorf("contract id is required")
	}
	var out entity.TxTransactionAction
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var contract entity.TxContractDeal
		if err := tx.Where("id = ? AND deleted_at IS NULL", in.ContractID).First(&contract).Error; err != nil {
			return err
		}
		if contract.TransactionID == nil {
			return fmt.Errorf("transaction not found")
		}
		var txn entity.TxTransaction
		if err := tx.Where("id = ? AND deleted_at IS NULL", *contract.TransactionID).First(&txn).Error; err != nil {
			return err
		}
		if txn.Status == uint32(enums.TxStatusCancel) || txn.Status == uint32(enums.TxStatusDone) {
			return fmt.Errorf("invalid transaction state")
		}

		ts := in.Timestamp
		if ts == nil {
			now := time.Now()
			ts = &now
		}
		action := entity.TxTransactionAction{
			TxID:      txn.ID,
			FromID:    in.FromID,
			FromOf:    uint32(enums.TxOwnerUser),
			ToID:      txn.ToID,
			ToOf:      txn.ToOf,
			Action:    in.Action,
			Value:     in.Amount,
			Note:      in.Note,
			Message:   enums.TxActionNames[enums.TxAction(in.Action)],
			Timestamp: ts,
		}
		if err := tx.Create(&action).Error; err != nil {
			return err
		}
		if err := tx.Model(&txn).Update("last_action_id", action.ID).Error; err != nil {
			return err
		}

		switch enums.TxAction(in.Action) {
		case enums.TxActionHandOver:
			if err := tx.Model(&txn).Update("status", uint32(enums.TxStatusDone)).Error; err != nil {
				return err
			}
			if contract.ProductID != nil && txn.Method == uint32(enums.TxMethodBuy) {
				if err := ensureDealProductTx(tx, contract.DealID, *contract.ProductID); err != nil {
					return err
				}
			}
		case enums.TxActionCancel:
			if err := tx.Model(&txn).Update("status", uint32(enums.TxStatusCancel)).Error; err != nil {
				return err
			}
		}
		out = action
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *contractDealPostgres) EnsureDealProduct(ctx context.Context, dealID, productID uint64) error {
	if r.db == nil || dealID == 0 || productID == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return ensureDealProductTx(tx, dealID, productID)
	})
}

func ensureDealProductTx(tx *gorm.DB, dealID, productID uint64) error {
	var count int64
	if err := tx.Model(&entity.DealProduct{}).
		Where("deal_id = ? AND product_id = ? AND deleted_at IS NULL", dealID, productID).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	return tx.Create(&entity.DealProduct{DealID: dealID, ProductID: productID}).Error
}
