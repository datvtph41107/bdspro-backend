package postgre

import (
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/enums"

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
	return nil, 0, nil
}

func (r *HistoryRepo) CreateHistory(c context.Context, entity *domain.HistoryEntity) error {
	return r.DB.Create(entity).Error
}
