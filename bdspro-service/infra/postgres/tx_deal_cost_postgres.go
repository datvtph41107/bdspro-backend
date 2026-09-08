package postgres

import (
	"context"

	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	"bdspro/internal/repo"

	"gorm.io/gorm"
)

type DealCostPostgreRepo struct {
	db *gorm.DB
}

func NewDealCostRepository(db *gorm.DB) repo.TxDealCostRepository {
	return &DealCostPostgreRepo{db: db}
}

func (r *DealCostPostgreRepo) Create(ctx context.Context, dealCost *domain.TxDealCost) error {
	return r.db.WithContext(ctx).Create(dealCost).Error
}

func (r *DealCostPostgreRepo) Update(ctx context.Context, dealCost *domain.TxDealCost) error {
	return r.db.WithContext(ctx).Save(dealCost).Error
}

func (r *DealCostPostgreRepo) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).
		Model(&domain.TxDealCost{}).
		Where("id = ?", id).
		Delete(&domain.TxDealCost{}).Error
}

func (r *DealCostPostgreRepo) GetByID(ctx context.Context, id uint64) (*domain.TxDealCost, error) {
	var cost domain.TxDealCost
	err := r.db.WithContext(ctx).
		Preload("CostType").
		Where("id = ?", id).
		First(&cost).Error
	if err != nil {
		return nil, err
	}
	return &cost, nil
}

func (r *DealCostPostgreRepo) ListByDealID(ctx context.Context, search *dto.TxDealCostSearchDTO) ([]*domain.TxDealCost, error) {
	var costs []*domain.TxDealCost
	query := r.db.WithContext(ctx).
		Preload("CostType").
		Where("deal_id = ?", search.DealId)

	query = query.Offset(search.GetOffset()).Limit(search.GetLimit())

	err := query.Find(&costs).Error
	if err != nil {
		return nil, err
	}
	return costs, nil
}

func (r *DealCostPostgreRepo) GetStatisticsByDealID(ctx context.Context, dealID uint64) (*dto.TxDealCostStatisticsDTO, error) {
	// Query for cost transactions (EMethodCost)
	var costStats struct {
		// Status enums.ApprovedStatus
		Amount float64
		Count  int64
	}

	err := r.db.WithContext(ctx).
		Table("tx_deal_costs").
		Select("SUM(amount) as amount, COUNT(*) as count").
		Joins("JOIN tx_cost_types ON tx_deal_costs.cost_type_id = tx_cost_types.id").
		Where("deal_id = ? AND tx_deal_costs.deleted_at IS NULL AND tx_cost_types.cost_type = ?", dealID, enums.TxCostTypeExpense).
		// Group("status").
		Scan(&costStats).Error
	if err != nil {
		return nil, err
	}

	// Query for revenue transactions (EMethodRevenue)
	var revenueStats struct {
		// Status enums.ApprovedStatus
		Amount float64
		Count  int64
	}

	err = r.db.WithContext(ctx).
		Table("tx_deal_costs").
		Select("SUM(amount) as amount, COUNT(*) as count").
		Joins("JOIN tx_cost_types ON tx_deal_costs.cost_type_id = tx_cost_types.id").
		Where("deal_id = ? AND tx_deal_costs.deleted_at IS NULL AND tx_cost_types.cost_type = ?", dealID, enums.TxCostTypeRevenue).
		// Group("status").
		Scan(&revenueStats).Error
	if err != nil {
		return nil, err
	}

	// Initialize result
	result := &dto.TxDealCostStatisticsDTO{
		DealID:        dealID,
		AmountCost:    &dto.TxStats{},
		AmountRevenue: &dto.TxStats{},
		NumCost:       &dto.TxStats{},
		NumRevenue:    &dto.TxStats{},
	}

	// // Process cost statistics
	// for _, stat := range costStats {
	// 	switch stat.Status {
	// 	case enums.ApprovedStatusApproved:
	// 		result.AmountCost.Approved = stat.Amount
	// 		result.NumCost.Approved = float64(stat.Count)
	// 	case enums.ApprovedStatusRejected:
	// 		result.AmountCost.Reject = stat.Amount
	// 		result.NumCost.Reject = float64(stat.Count)
	// 	case enums.ApprovedStatusPending:
	// 		result.AmountCost.Pending = stat.Amount
	// 		result.NumCost.Pending = float64(stat.Count)
	// 	}
	// }

	result.AmountCost.Approved = costStats.Amount
	result.NumCost.Approved = float64(costStats.Count)
	result.AmountRevenue.Approved = revenueStats.Amount
	result.NumRevenue.Approved = float64(revenueStats.Count)
	// Process revenue statistics
	// for _, stat := range revenueStats {
	// 	switch stat.Status {
	// 	case enums.ApprovedStatusApproved:
	// 		result.AmountRevenue.Approved = stat.Amount
	// 		result.NumRevenue.Approved = float64(stat.Count)
	// 	case enums.ApprovedStatusRejected:
	// 		result.AmountRevenue.Reject = stat.Amount
	// 		result.NumRevenue.Reject = float64(stat.Count)
	// 	case enums.ApprovedStatusPending:
	// 		result.AmountRevenue.Pending = stat.Amount
	// 		result.NumRevenue.Pending = float64(stat.Count)
	// 	}
	// }

	return result, nil
}

func (r *DealCostPostgreRepo) GetStatisticsByDealIDs(ctx context.Context, dealIDs []uint64) (*dto.TxStatisticDealsDTO, error) {
	var totalTx, totalRevenue, totalExpense float64

	// 1. Transactions (cộng theo tx_transaction_action.value)
	err := r.db.WithContext(ctx).
		Debug().
		Table("tx_contract_deals cd").
		Select("COALESCE(SUM(ta.value),0)").
		Joins("JOIN tx_transaction t ON t.id = cd.transaction_id").
		Joins("JOIN tx_transaction_action ta ON ta.tx_id = t.id").
		Where("cd.deal_id IN ? AND cd.deleted_at IS NULL AND ta.deleted_at IS NULL", dealIDs).
		Scan(&totalTx).Error
	if err != nil {
		return nil, err
	}

	// 2. Revenue (cost_type = 20)
	err = r.db.WithContext(ctx).
		Debug().
		Table("tx_deal_costs dc").
		Select("COALESCE(SUM(dc.amount),0)").
		Joins("JOIN tx_cost_types ct ON ct.id = dc.cost_type_id").
		Where("dc.deal_id IN ? AND dc.deleted_at IS NULL AND ct.cost_type = ?", dealIDs, enums.TxCostTypeRevenue).
		Scan(&totalRevenue).Error
	if err != nil {
		return nil, err
	}

	// 3. Expense (cost_type = 10)
	err = r.db.WithContext(ctx).
		Debug().
		Table("tx_deal_costs dc").
		Select("COALESCE(SUM(dc.amount),0)").
		Joins("JOIN tx_cost_types ct ON ct.id = dc.cost_type_id").
		Where("dc.deal_id IN ? AND dc.deleted_at IS NULL AND ct.cost_type = ?", dealIDs, enums.TxCostTypeExpense).
		Scan(&totalExpense).Error
	if err != nil {
		return nil, err
	}

	// 4. Tổng hợp kết quả
	result := &dto.TxStatisticDealsDTO{
		DealIDs:       dealIDs,
		AmountTxt:     totalTx,
		AmountRevenue: totalRevenue,
		AmountCost:    totalExpense,
		TotalProfit:   totalTx + totalRevenue - totalExpense,
	}

	return result, nil
}
