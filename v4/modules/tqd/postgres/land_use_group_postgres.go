package postgres

// import (
// 	"context"
// 	"errors"
// 	"fmt"

// 	qh "userv4/modules/tqd/entity/qh"
// 	repo "userv4/modules/tqd/repository"

// 	"gorm.io/gorm"
// )

// type landUseGroupRepo struct {
// 	db *gorm.DB
// }

// func NewLandUseGroupRepository(db *gorm.DB) repo.LandUseGroupRepository {
// 	return &landUseGroupRepo{db: db}
// }

// func (r *landUseGroupRepo) GetByID(ctx context.Context, id uint64) (*qh.LandUseGroup, error) {
// 	if id == 0 {
// 		return nil, nil
// 	}
// 	var row qh.LandUseGroup
// 	err := r.db.WithContext(ctx).Where("id = ?", id).First(&row).Error
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, nil
// 		}
// 		return nil, fmt.Errorf("get land_use_group by id: %w", err)
// 	}
// 	return &row, nil
// }
