package postgres

import (
	"bdspro/internal/domain"
	property_repo "bdspro/internal/repo/property"
	"context"
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
)

type WardV2Repo struct {
	DB *gorm.DB
}

func NewWardV2Repository(db *gorm.DB) property_repo.WardV2Repository {
	return &WardV2Repo{DB: db}
}
func (r *WardV2Repo) FirstOrCreate(ctx context.Context, ward *domain.WardV2) error {
	now := time.Now()
	ward.CreatedAt = &now
	ward.UpdatedAt = &now

	if originID, ok := ctx.Value("origin_id").(uint64); ok && originID > 0 {
		ward.CreatedBy = &originID
		ward.UpdatedBy = &originID
	}

	return r.DB.WithContext(ctx).
		Where(domain.WardV2{TQDID: ward.TQDID}).
		Assign(domain.WardV2{
			Name:          ward.Name,
			Code:          ward.Code,
			Codename:      ward.Codename,
			DivisionType:  ward.DivisionType,
			ShortCodename: ward.ShortCodename,
			Lat:           ward.Lat,
			Lng:           ward.Lng,
			ProvinceID:    ward.ProvinceID,
		}).
		FirstOrCreate(ward).Error
}
func (r *WardV2Repo) GetByID(ctx context.Context, id uint64) (*domain.WardV2, error) {
	log.Printf("Getting WardV2 by ID: %d", id)
	var entity domain.WardV2
	err := GetDB(ctx, r.DB).
		First(&entity, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *WardV2Repo) GetByCode(ctx context.Context, code string) (*domain.WardV2, error) {
	var entity domain.WardV2
	err := GetDB(ctx, r.DB).
		Where("code = ?", code).First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *WardV2Repo) GetByTQDID(ctx context.Context, tqdID string) (*domain.WardV2, error) {
	var entity domain.WardV2
	err := GetDB(ctx, r.DB).
		Where("tqd_id = ?", tqdID).First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *WardV2Repo) Create(ctx context.Context, entity *domain.WardV2) error {
	return GetDB(ctx, r.DB).Create(entity).Error
}

func (r *WardV2Repo) UpdateFields(ctx context.Context, id uint64, fields map[string]any) error {
	return GetDB(ctx, r.DB).
		Model(&domain.WardV2{}).
		Where("id = ?", id).
		Updates(fields).Error
}

// Thêm method để lấy danh sách wards theo province nếu cần
func (r *WardV2Repo) ListByProvinceID(ctx context.Context, provinceID uint64) ([]*domain.WardV2, error) {
	var entities []*domain.WardV2
	err := GetDB(ctx, r.DB).
		Where("province_id = ?", provinceID).
		Find(&entities).Error
	return entities, err
}
