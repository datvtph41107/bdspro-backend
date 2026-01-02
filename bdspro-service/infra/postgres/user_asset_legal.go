package postgres

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"context"
	"time"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo.AssetLegalRepo
type GormAssetLegalRepo struct {
	DB *gorm.DB
}

func NewGormAssetLegalRepo(DB *gorm.DB) *GormAssetLegalRepo {
	return &GormAssetLegalRepo{
		DB: DB,
	}
}

func (r *GormAssetLegalRepo) Create(c context.Context, entity *domain.AssetLegal) error {
	return r.DB.WithContext(c).Create(entity).Error
}

func (r *GormAssetLegalRepo) Update(c context.Context, id uint64, entity *domain.AssetLegal) error {
	return r.DB.WithContext(c).
		Model(entity).
		Where("id = ? and deleted_at is null", id).
		Updates(entity).Error
}

func (r *GormAssetLegalRepo) Delete(c context.Context, id uint64) error {
	return r.DB.WithContext(c).
		Model(&domain.AssetLegal{}).
		Where("id = ?", id).
		Update("deleted_at", time.Now()).Error
}

func (r *GormAssetLegalRepo) GetDataQuery(c context.Context, dto dto.AssetLegalGetDTO) *gorm.DB {
	query := r.DB.Model(domain.AssetLegal{}).Where("deleted_at IS NULL")
	if dto.AssetID != 0 {
		query = query.Where("asset_id = ?", dto.AssetID)
	}
	return query
}
func (r *GormAssetLegalRepo) GetData(c context.Context, dto dto.AssetLegalGetDTO) ([]domain.AssetLegal, int64, error) {
	var entities []domain.AssetLegal
	query := r.GetDataQuery(c, dto)
	total := int64(0)

	err := query.Limit(dto.GetLimit()).
		Offset(dto.GetOffset()).
		Find(&entities).Error
	if err != nil {
		return nil, 0, err
	}

	query = r.GetDataQuery(c, dto)

	err = query.Count(&total).Error
	return entities, total, err
}

func (r *GormAssetLegalRepo) GetByID(c context.Context, id uint64) (*domain.AssetLegal, error) {
	var entity domain.AssetLegal
	err := r.DB.Where("id = ? AND deleted_at IS NULL", id).First(&entity).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *GormAssetLegalRepo) UpdateLegalItems(c context.Context, assetID uint64, items []domain.AssetLegal) error {
	// Bắt đầu transaction để đảm bảo toàn vẹn dữ liệu
	tx := r.DB.WithContext(c).Begin()

	// Xóa các bản ghi cũ trong bảng
	if err := tx.Where("asset_id = ?", assetID).Delete(&domain.AssetLegal{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Nếu danh sách amenityIDs rỗng -> không cần thêm mới, chỉ cần xóa là đủ
	if len(items) == 0 {
		tx.Commit()
		return nil
	}

	// Tạo danh sách mới để chèn vào
	var legalItems []domain.AssetLegal
	for _, item := range items {
		legalItems = append(legalItems, domain.AssetLegal{
			AssetID:      assetID,
			DocumentName: item.DocumentName,
			DocumentURL:  item.DocumentURL,
			DocumentType: item.DocumentType,
			IssuedDate:   item.IssuedDate,
			ExpiryDate:   item.ExpiryDate,
			Description:  item.Description,
		})
	}

	// Chèn danh sách mới vào bảng nối
	if err := tx.Create(&legalItems).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction
	tx.Commit()
	return nil
}
