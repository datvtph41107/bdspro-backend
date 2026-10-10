package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"userv4/modules/bdspro/dto"
	"userv4/modules/bdspro/entity"
	"userv4/modules/bdspro/enums"

	"gorm.io/gorm"
)

// InvestmentRepository is deal_investments CRUD on the v4 DSN.
type InvestmentRepository interface {
	Create(ctx context.Context, row *entity.DealInvestment) error
	Update(ctx context.Context, row *entity.DealInvestment) error
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*entity.DealInvestment, error)
	List(ctx context.Context, q dto.InvestmentListQuery) ([]*entity.DealInvestment, int64, error)
	ChangeStatus(ctx context.Context, id uint64, status enums.TxApprovedStatus) error
	ChangeConfirmation(ctx context.Context, id uint64, confirmed bool) error
	ListConfirmedByDealMember(ctx context.Context, dealID, memberID uint64) ([]*entity.DealInvestment, error)
}

type investmentPostgres struct {
	db *gorm.DB
}

func NewInvestmentPostgres(db *gorm.DB) InvestmentRepository {
	return &investmentPostgres{db: db}
}

func (r *investmentPostgres) Create(ctx context.Context, row *entity.DealInvestment) error {
	if r.db == nil || row == nil {
		return nil
	}
	if row.Status == 0 ||
		(row.Status != uint32(enums.TxApprovedStatusPending) &&
			row.Status != uint32(enums.TxApprovedStatusRejected) &&
			row.Status != uint32(enums.TxApprovedStatusApproved)) {
		row.Status = uint32(enums.TxApprovedStatusPending)
	}
	row.ID = 0
	err := r.db.WithContext(ctx).Omit("ID").Create(row).Error
	if err == nil {
		return nil
	}
	if strings.Contains(err.Error(), "deal_investments_pkey") {
		_ = r.db.WithContext(ctx).Exec(`
			SELECT setval(
				pg_get_serial_sequence('deal_investments', 'id'),
				COALESCE((SELECT MAX(id) FROM deal_investments), 1)
			)`).Error
		row.ID = 0
		return r.db.WithContext(ctx).Omit("ID").Create(row).Error
	}
	return err
}

func (r *investmentPostgres) Update(ctx context.Context, row *entity.DealInvestment) error {
	if r.db == nil || row == nil || row.ID == 0 {
		return nil
	}
	current, err := r.GetByID(ctx, row.ID)
	if err != nil {
		return err
	}
	if current == nil {
		return gorm.ErrRecordNotFound
	}
	if current.Status == uint32(enums.TxApprovedStatusApproved) {
		return errors.New("cannot update approved investment")
	}
	status := row.Status
	if status == 0 || !enums.TxApprovedStatus(status).IsValid() {
		status = current.Status
	}
	if current.Status == uint32(enums.TxApprovedStatusRejected) {
		status = uint32(enums.TxApprovedStatusPending)
	}
	now := time.Now()
	return r.db.WithContext(ctx).Model(&entity.DealInvestment{}).
		Where("id = ? AND deleted_at IS NULL", row.ID).
		Updates(map[string]interface{}{
			"member_id":            row.MemberID,
			"invest_type":          row.InvestType,
			"amount":               row.Amount,
			"transfer_time":        row.TransferTime,
			"note":                 row.Note,
			"transfer_proof_image": row.TransferProofImage,
			"status":               status,
			"changed_at":           now,
			"proxy_user_id":        row.ProxyUserID,
		}).Error
}

func (r *investmentPostgres) Delete(ctx context.Context, id uint64) error {
	if r.db == nil || id == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity.DealInvestment{}).Error
}

func (r *investmentPostgres) GetByID(ctx context.Context, id uint64) (*entity.DealInvestment, error) {
	if r.db == nil || id == 0 {
		return nil, nil
	}
	var row entity.DealInvestment
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *investmentPostgres) List(ctx context.Context, q dto.InvestmentListQuery) ([]*entity.DealInvestment, int64, error) {
	if r.db == nil || q.DealID == 0 {
		return nil, 0, nil
	}
	query := r.db.WithContext(ctx).Model(&entity.DealInvestment{}).
		Where("deal_id = ? AND deleted_at IS NULL", q.DealID)
	if q.MemberID != nil && *q.MemberID > 0 {
		query = query.Where("member_id = ?", *q.MemberID)
	}
	if q.InvestType > 0 {
		query = query.Where("invest_type = ?", q.InvestType)
	}
	if q.Status > 0 {
		query = query.Where("status = ?", q.Status)
	}
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []*entity.DealInvestment
	err := query.Order("updated_at DESC").Find(&rows).Error
	return rows, total, err
}

func (r *investmentPostgres) ChangeStatus(ctx context.Context, id uint64, status enums.TxApprovedStatus) error {
	if r.db == nil || id == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Model(&entity.DealInvestment{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]interface{}{
			"status":     status,
			"changed_at": time.Now(),
		}).Error
}

func (r *investmentPostgres) ChangeConfirmation(ctx context.Context, id uint64, confirmed bool) error {
	if r.db == nil || id == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Model(&entity.DealInvestment{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]interface{}{
			"confirmed":  confirmed,
			"changed_at": time.Now(),
		}).Error
}

func (r *investmentPostgres) ListConfirmedByDealMember(ctx context.Context, dealID, memberID uint64) ([]*entity.DealInvestment, error) {
	if r.db == nil || dealID == 0 {
		return nil, nil
	}
	var rows []*entity.DealInvestment
	err := r.db.WithContext(ctx).
		Where("deal_id = ? AND member_id = ? AND status = ? AND deleted_at IS NULL",
			dealID, memberID, enums.TxApprovedStatusApproved).
		Order("transfer_time ASC").
		Find(&rows).Error
	return rows, err
}
