package repository

import (
	"context"
	"strings"

	"userv4/modules/bdspro/dto"
	"userv4/modules/bdspro/entity"
	"userv4/modules/bdspro/enums"

	"gorm.io/gorm"
)

func (r *dealPostgres) ListDealProductIDs(ctx context.Context, dealID uint64, page, size uint32) ([]uint64, error) {
	if r.db == nil || dealID == 0 {
		return nil, nil
	}
	limit := int(size)
	if limit <= 0 {
		limit = 20
	}
	offset := int(page) * limit
	var ids []uint64
	err := r.db.WithContext(ctx).
		Model(&entity.DealProduct{}).
		Where("deal_id = ? AND deleted_at IS NULL", dealID).
		Order("id ASC").
		Offset(offset).Limit(limit).
		Pluck("product_id", &ids).Error
	return ids, err
}

func (r *dealPostgres) ListContractsByDealID(ctx context.Context, dealID uint64) ([]dto.TxContractDealListRow, int64, error) {
	if r.db == nil || dealID == 0 {
		return nil, 0, nil
	}
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
		Where(`tx_contract_deals.deal_id = ? AND tx_contract_deals.deleted_at IS NULL`, dealID).
		Group("tx_contract_deals.id, tx_contract_deals.deal_id, tx_contract_deals.created_at, tx_contract_deals.customer_id, txs.method, txs.transaction_name, txs.status, txs.timestamp").
		Order("tx_contract_deals.created_at DESC")

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []dto.TxContractDealListRow
	if err := q.Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	return rows, total, nil
}

func (r *dealPostgres) ListMembersByDealID(ctx context.Context, dealID uint64, page, size int) ([]*entity.DealMember, uint32, error) {
	if r.db == nil || dealID == 0 {
		return nil, 0, nil
	}
	if size <= 0 {
		size = 10
	}
	var total int64
	if err := r.db.WithContext(ctx).Model(&entity.DealMember{}).
		Where("deal_id = ? AND deleted_at IS NULL", dealID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []*entity.DealMember
	err := r.db.WithContext(ctx).
		Where("deal_id = ? AND deleted_at IS NULL", dealID).
		Order("CASE WHEN role_key = 410 THEN 0 ELSE 1 END ASC").
		Order("created_at DESC").
		Offset(page * size).Limit(size).
		Find(&list).Error
	return list, uint32(total), err
}

// ListAcceptedMembers returns accepted deal members (status=20), matching legacy GET .../members.
func (r *dealPostgres) ListAcceptedMembers(ctx context.Context, filter dto.DealMembersRequest) ([]*entity.DealMember, error) {
	if r.db == nil || filter.DealID == 0 {
		return nil, nil
	}
	q := r.db.WithContext(ctx).
		Where("deal_id = ? AND status = ? AND deleted_at IS NULL", filter.DealID, enums.DealMemberStatusAccepted)
	if filter.DoneInvestment != nil {
		q = q.Where("done_investment = ?", *filter.DoneInvestment)
	}
	var list []*entity.DealMember
	err := q.Order("CASE WHEN role_key = 410 THEN 0 ELSE 1 END ASC").
		Order("created_at DESC").
		Find(&list).Error
	return list, err
}

func (r *dealPostgres) ListInvestmentsByDealID(ctx context.Context, dealID uint64, memberID uint64) ([]*entity.DealInvestment, int64, error) {
	if r.db == nil || dealID == 0 {
		return nil, 0, nil
	}
	q := r.db.WithContext(ctx).Model(&entity.DealInvestment{}).
		Where("deal_id = ? AND deleted_at IS NULL", dealID)
	if memberID > 0 {
		q = q.Where("member_id = ?", memberID)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []*entity.DealInvestment
	err := q.Order("updated_at DESC").Find(&list).Error
	return list, total, err
}

func (r *dealPostgres) ListCostsByDealID(ctx context.Context, dealID uint64, page, size uint32) ([]*entity.TxDealCost, error) {
	if r.db == nil || dealID == 0 {
		return nil, nil
	}
	limit := int(size)
	if limit <= 0 {
		limit = 20
	}
	offset := int(page) * limit
	var list []*entity.TxDealCost
	err := r.db.WithContext(ctx).
		Preload("CostType").
		Where("deal_id = ? AND deleted_at IS NULL", dealID).
		Offset(offset).Limit(limit).
		Find(&list).Error
	return list, err
}

func (r *dealPostgres) CountDashboard(ctx context.Context, summary dto.DealSummaryRequest) (*dto.DealSummaryRow, error) {
	result := &dto.DealSummaryRow{}
	if r.db == nil || summary.DealID == 0 {
		return result, nil
	}

	var selectFields []string
	var args []interface{}
	all := summary.All

	if summary.AmountPending || all {
		selectFields = append(selectFields, `COALESCE(SUM(CASE WHEN status = ? THEN amount ELSE 0 END), 0) as amount_pending`)
		args = append(args, enums.TxApprovedStatusPending)
	}
	if summary.AmountApproved || all {
		selectFields = append(selectFields, `COALESCE(SUM(CASE WHEN status = ? THEN amount ELSE 0 END), 0) as amount_approved`)
		args = append(args, enums.TxApprovedStatusApproved)
	}
	if summary.AmountRejected || all {
		selectFields = append(selectFields, `COALESCE(SUM(CASE WHEN status = ? THEN amount ELSE 0 END), 0) as amount_rejected`)
		args = append(args, enums.TxApprovedStatusRejected)
	}
	if summary.NumPending || all {
		selectFields = append(selectFields, `COUNT(CASE WHEN status = ? THEN 1 END) as num_pending`)
		args = append(args, enums.TxApprovedStatusPending)
	}
	if summary.NumApproved || all {
		selectFields = append(selectFields, `COUNT(CASE WHEN status = ? THEN 1 END) as num_approved`)
		args = append(args, enums.TxApprovedStatusApproved)
	}
	if summary.NumRejected || all {
		selectFields = append(selectFields, `COUNT(CASE WHEN status = ? THEN 1 END) as num_rejected`)
		args = append(args, enums.TxApprovedStatusRejected)
	}
	if summary.NumInvestment || all {
		selectFields = append(selectFields, `COUNT(*) as num_investment`)
	}
	if summary.NumMemberSubmit || all {
		selectFields = append(selectFields, `COUNT(DISTINCT member_id) as num_member_submit`)
	}
	if summary.AmountCost || all {
		selectFields = append(selectFields, `
			(SELECT COALESCE(SUM(dc.amount), 0) FROM tx_deal_costs dc
			LEFT JOIN tx_cost_types ct ON ct.id = dc.cost_type_id
			WHERE dc.deal_id = ? AND ct.cost_type = 10 AND dc.deleted_at IS NULL
		) as amount_cost`)
		args = append(args, summary.DealID)
	}
	if len(selectFields) == 0 {
		return result, nil
	}

	err := r.db.WithContext(ctx).
		Model(&entity.DealInvestment{}).
		Where("deal_id = ? AND deleted_at IS NULL", summary.DealID).
		Select(strings.Join(selectFields, ", "), args...).
		Scan(result).Error
	return result, err
}

func (r *dealPostgres) GetBankAccountByID(ctx context.Context, id uint64) (*entity.BankAccount, error) {
	if r.db == nil || id == 0 {
		return nil, nil
	}
	var row entity.BankAccount
	err := r.db.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&row).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &row, nil
}

func (r *dealPostgres) GetCostStatisticsByDealID(ctx context.Context, dealID uint64) (*dto.DealCostStatsRow, error) {
	result := &dto.DealCostStatsRow{}
	if r.db == nil || dealID == 0 {
		return result, nil
	}

	var costStats struct {
		Amount float64
		Count  int64
	}
	if err := r.db.WithContext(ctx).
		Table("tx_deal_costs").
		Select("COALESCE(SUM(amount), 0) as amount, COUNT(*) as count").
		Joins("JOIN tx_cost_types ON tx_deal_costs.cost_type_id = tx_cost_types.id").
		Where("deal_id = ? AND tx_deal_costs.deleted_at IS NULL AND tx_cost_types.cost_type = ?", dealID, enums.TxCostTypeExpense).
		Scan(&costStats).Error; err != nil {
		return nil, err
	}

	var revenueStats struct {
		Amount float64
		Count  int64
	}
	if err := r.db.WithContext(ctx).
		Table("tx_deal_costs").
		Select("COALESCE(SUM(amount), 0) as amount, COUNT(*) as count").
		Joins("JOIN tx_cost_types ON tx_deal_costs.cost_type_id = tx_cost_types.id").
		Where("deal_id = ? AND tx_deal_costs.deleted_at IS NULL AND tx_cost_types.cost_type = ?", dealID, enums.TxCostTypeRevenue).
		Scan(&revenueStats).Error; err != nil {
		return nil, err
	}

	result.AmountCost.Approved = costStats.Amount
	result.NumCost.Approved = float64(costStats.Count)
	result.AmountRevenue.Approved = revenueStats.Amount
	result.NumRevenue.Approved = float64(revenueStats.Count)
	return result, nil
}

func (r *dealPostgres) GetPaymentStatisticsByDealID(ctx context.Context, dealID uint64) (*dto.DealPaymentStatsRow, error) {
	result := &dto.DealPaymentStatsRow{}
	if r.db == nil || dealID == 0 {
		return result, nil
	}

	var agg struct {
		TotalValue float64
		TotalCount float64
	}
	err := r.db.WithContext(ctx).Raw(`
		SELECT
			COUNT(sub.id) AS total_count,
			COALESCE(SUM(sub.v), 0) AS total_value
		FROM (
			SELECT
				t.id,
				SUM(ta.value) AS v
			FROM tx_contract_deals cd
			JOIN tx_transaction t ON cd.transaction_id = t.id
			JOIN tx_transaction_action ta ON ta.tx_id = t.id
			WHERE cd.deal_id = ?
			  AND t.method IN (?, ?)
			  AND t.status = ?
			  AND ta.deleted_at IS NULL
			  AND t.deleted_at IS NULL
			  AND cd.deleted_at IS NULL
			GROUP BY t.id
		) sub
	`, dealID, enums.TxMethodSale, enums.TxMethodRent, enums.TxApprovedStatusApproved).
		Scan(&agg).Error
	if err != nil {
		return nil, err
	}

	result.Payment = dto.TxStats{
		Approved: agg.TotalValue,
		Total:    agg.TotalValue,
	}
	result.NumPayment = dto.TxStats{
		Approved: agg.TotalCount,
		Total:    agg.TotalCount,
	}
	return result, nil
}

