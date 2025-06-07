package postgre

import (
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/enums"
	"time"

	"gorm.io/gorm"
)

// @bind: crm/internal/repo.HistoryRepo
type HistoryRepo struct {
	DB *gorm.DB
}

func NewHistoryRepo(db *gorm.DB) *HistoryRepo {
	return &HistoryRepo{DB: db}
}

func (r *HistoryRepo) Search(c context.Context,
	ownerID uint64,
	ownerType enums.EOwnerType,
	dto dto.HistorySearchDTO,
) ([]domain.HistoryEntity, int64, error) {
	query := r.DB.WithContext(c).Model(&domain.HistoryEntity{}).
		Where("owner_id = ?", ownerID).
		Where("owner_type = ?", ownerType)

	if dto.TargetType != nil {
		query = query.Where("target_id = ? and target_type = ?", dto.TargetId, dto.TargetType)
	}

	if dto.ActionType != nil {
		query = query.Where("action_type = ?", dto.ActionType)
	}

	if dto.FromDate != nil {
		query = query.Where("created_at >= ?", dto.FromDate)
	}

	if dto.ToDate != nil {
		query = query.Where("created_at <= ?", dto.ToDate)
	}

	query = query.Order("created_at DESC")

	var histories []domain.HistoryEntity
	var total int64

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Find(&histories).Error
	if err != nil {
		return nil, 0, err
	}

	return histories, total, nil
}

func (r *HistoryRepo) CreateHistory(c context.Context, entity *domain.HistoryEntity) error {
	return r.DB.WithContext(c).Create(entity).Error
}

func (r *HistoryRepo) ExistedFromDate(c context.Context, customerID uint64, fromDate time.Time) (bool, error) {
	var count int64
	err := r.DB.WithContext(c).Model(&domain.HistoryEntity{}).
		Where(`target_id = ? 
		and target_type = ? 
		and created_at >= ?
		and deleted_at is null
		`, customerID, enums.TargetHistoryCustomer, fromDate).
		Count(&count).Error
	return count > 0, err
}
