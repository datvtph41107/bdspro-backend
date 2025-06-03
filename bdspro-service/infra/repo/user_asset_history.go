package postgre_repo

import (
	"bdspro/internal/enums"
	"bdspro/v1/models"
	_db "common/db"
	"context"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo.UAssetHistoryRepo
type GormAssetHistoryRepo struct {
	DB *gorm.DB
}

func NewGormAssetSplitHistoryRepo(DB *gorm.DB) *GormAssetHistoryRepo {
	return &GormAssetHistoryRepo{
		DB: DB,
	}
}

func (*GormAssetHistoryRepo) CreateHistory(
	c context.Context,
	action enums.AssetHistory,
	assetID *uint64,
	// parentID uint64,
	BeforeData *string,
	AfterData *string,
) error {
	return _db.SaveWithContext(c, &models.AssetHistory{
		AssetID:    assetID,
		ActionType: action,
		BeforeData: BeforeData,
		AfterData:  AfterData,
		// ParentID: &parentID,
		// ChildAssetIDs: ConvertToPQArray(childIDs),
	})
}

// func (r *GormAssetHistoryRepo) Create(c context.Context, entity *domain.AssetHistory) error {
// 	return r.DB.WithContext(c).Create(entity).Error
// }

// func (r *GormAssetHistoryRepo) Update(c context.Context, id uint64, entity *domain.AssetHistory) error {
// 	return r.DB.WithContext(c).
// 		Model(entity).
// 		Where("id = ? and deleted_at is null", id).
// 		Updates(entity).Error
// }

// func (r *GormAssetHistoryRepo) Delete(c context.Context, id uint64) error {
// 	return r.DB.WithContext(c).
// 		Model(&domain.AssetHistory{}).
// 		Where("id = ?", id).
// 		Update("deleted_at", time.Now()).Error
// }

// func (r *GormAssetHistoryRepo) GetData(c context.Context) ([]domain.AssetHistory, error) {
// 	var entities []domain.AssetHistory
// 	err := r.DB.Where("deleted_at IS NULL").Find(&entities).Error
// 	return entities, err
// }

// func (r *GormAssetHistoryRepo) GetByID(c context.Context, id uint64) (*domain.AssetHistory, error) {
// 	var entity domain.AssetHistory
// 	err := r.DB.Where("id = ? AND deleted_at IS NULL", id).First(&entity).Error
// 	if err != nil {
// 		return nil, err
// 	}
// 	return &entity, nil
// }
