package postgres

import (
	"bdspro/internal/domain"
	property_repo "bdspro/internal/repo/property"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProvinceV2Repo struct {
	DB *gorm.DB
}

func NewProvinceV2Repository(db *gorm.DB) property_repo.ProvinceV2Repository {
	return &ProvinceV2Repo{DB: db}
}

func (r *ProvinceV2Repo) GetByID(ctx context.Context, id uint64) (*domain.ProvinceV2, error) {
	var entity domain.ProvinceV2
	err := GetDB(ctx, r.DB).
		First(&entity, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *ProvinceV2Repo) GetByCode(ctx context.Context, code string) (*domain.ProvinceV2, error) {
	var entity domain.ProvinceV2
	err := GetDB(ctx, r.DB).
		Where("code = ?", code).First(&entity).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *ProvinceV2Repo) GetByFullName(ctx context.Context, name string) (*domain.ProvinceV2, error) {
	var entity domain.ProvinceV2
	err := GetDB(ctx, r.DB).
		Where("name LIKE ?", "%"+name+"%").
		First(&entity).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &entity, err
}

func (r *ProvinceV2Repo) GetByTQDID(ctx context.Context, tqdID string) (*domain.ProvinceV2, error) {
	var province domain.ProvinceV2
	err := r.DB.WithContext(ctx).
		Where("tqd_id = ?", tqdID).
		Where("deleted_at IS NULL").
		First(&province).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &province, err
}

// func (r *ProvinceV2Repo) Create(ctx context.Context, entity *domain.ProvinceV2) error {
// 	return GetDB(ctx, r.DB).Create(entity).Error
// }

func (r *ProvinceV2Repo) UpdateFields(ctx context.Context, id uint64, fields map[string]any) error {
	return GetDB(ctx, r.DB).
		Model(&domain.ProvinceV2{}).
		Where("id = ?", id).
		Updates(fields).Error
}

func (r *ProvinceV2Repo) Create(ctx context.Context, province *domain.ProvinceV2) error {
	// Set timestamps
	now := time.Now()
	province.CreatedAt = &now
	province.UpdatedAt = &now

	// Set created_by từ context
	if originID, ok := ctx.Value("origin_id").(uint64); ok && originID > 0 {
		province.CreatedBy = &originID
		province.UpdatedBy = &originID
	}

	// KHÔNG DÙNG OnConflict, dùng Create đơn giản
	return r.DB.WithContext(ctx).Create(province).Error
}

// Hoặc dùng FirstOrCreate nếu muốn
func (r *ProvinceV2Repo) FirstOrCreate(ctx context.Context, province *domain.ProvinceV2) error {
	// Set timestamps
	now := time.Now()
	province.CreatedAt = &now
	province.UpdatedAt = &now

	// FirstOrCreate sẽ tìm theo điều kiện, nếu không có thì tạo mới
	return r.DB.WithContext(ctx).
		Where(domain.ProvinceV2{TQDID: province.TQDID}).
		Assign(domain.ProvinceV2{
			Name:         province.Name,
			Code:         province.Code,
			Codename:     province.Codename,
			DivisionType: province.DivisionType,
			PhoneCode:    province.PhoneCode,
			Lat:          province.Lat,
			Lng:          province.Lng,
		}).
		FirstOrCreate(province).Error
}

// Thêm method GetByTQDIDWithLock nếu cần
func (r *ProvinceV2Repo) GetByTQDIDWithLock(ctx context.Context, tqdID string) (*domain.ProvinceV2, error) {
	var province domain.ProvinceV2
	err := r.DB.WithContext(ctx).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("tqd_id = ?", tqdID).
		First(&province).Error

	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &province, err
}
