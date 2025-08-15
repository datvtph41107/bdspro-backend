package postgre

import (
	"context"
	"tx/internal/domain/entity"
	"tx/internal/dto"
	"tx/internal/enums"
	"tx/internal/repo"

	"gorm.io/gorm"
)

type DealCostPostgreRepo struct {
	db *gorm.DB
}

func NewDealCostRepository(db *gorm.DB) repo.DealCostRepository {
	return &DealCostPostgreRepo{db: db}
}

func (r *DealCostPostgreRepo) Create(ctx context.Context, dealCost *entity.DealCost) error {
	return r.db.WithContext(ctx).Create(dealCost).Error
}

func (r *DealCostPostgreRepo) Update(ctx context.Context, dealCost *entity.DealCost) error {
	return r.db.WithContext(ctx).Save(dealCost).Error
}

func (r *DealCostPostgreRepo) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).
		Model(&entity.DealCost{}).
		Where("id = ?", id).
		Delete(&entity.DealCost{}).Error
}

func (r *DealCostPostgreRepo) GetByID(ctx context.Context, id uint64) (*entity.DealCost, error) {
	var cost entity.DealCost
	err := r.db.WithContext(ctx).
		Preload("CostType").
		Where("id = ?", id).
		First(&cost).Error
	if err != nil {
		return nil, err
	}
	return &cost, nil
}

func (r *DealCostPostgreRepo) ListByDealID(ctx context.Context, search *dto.DealCostSearchDTO) ([]*entity.DealCost, error) {
	var costs []*entity.DealCost
	query := r.db.WithContext(ctx).
		Preload("CostType").
		Where("deal_id = ?", search.DealId)

	if search.Page > 0 && search.Size > 0 {
		offset := (search.Page - 1) * search.Size
		query = query.Offset(offset).Limit(search.Size)
	}

	err := query.Find(&costs).Error
	if err != nil {
		return nil, err
	}
	return costs, nil
}

func (r *DealCostPostgreRepo) GetStatisticsByDealID(ctx context.Context, dealID uint64) (*dto.DealCostStatisticsDTO, error) {
	// Query for cost transactions (EMethodCost)
	var costStats struct {
		// Status enums.ApprovedStatus
		Amount float64
		Count  int64
	}

	err := r.db.WithContext(ctx).
		Table("deal_costs").
		Select("SUM(amount) as amount, COUNT(*) as count").
		Joins("JOIN cost_types ON deal_costs.cost_type_id = cost_types.id").
		Where("deal_id = ? AND deal_costs.deleted_at IS NULL AND cost_types.cost_type = ?", dealID, enums.ECostTypeExpense).
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
		Table("deal_costs").
		Select("SUM(amount) as amount, COUNT(*) as count").
		Joins("JOIN cost_types ON deal_costs.cost_type_id = cost_types.id").
		Where("deal_id = ? AND deal_costs.deleted_at IS NULL AND cost_types.cost_type = ?", dealID, enums.ECostTypeRevenue).
		// Group("status").
		Scan(&revenueStats).Error

	if err != nil {
		return nil, err
	}

	// Initialize result
	result := &dto.DealCostStatisticsDTO{
		DealID:        dealID,
		AmountCost:    &dto.Stats{},
		AmountRevenue: &dto.Stats{},
		NumCost:       &dto.Stats{},
		NumRevenue:    &dto.Stats{},
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
