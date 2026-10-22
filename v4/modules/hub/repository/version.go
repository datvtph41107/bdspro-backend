package repository

import (
	"context"
	"errors"
	"time"

	"mdn4/modules/hub/entity"

	"gorm.io/gorm"
)

type VersionRepository interface {
	GetLatest(ctx context.Context, appName, platform string, activeOnly bool) (*entity.Version, error)
	Create(ctx context.Context, v *entity.Version) error
	ValidateAPIKey(ctx context.Context, apiKey string) error
	GetByID(ctx context.Context, id uint64) (*entity.Version, error)
	List(ctx context.Context, appName, platform string, page, size int) ([]*entity.Version, int64, error)
	Update(ctx context.Context, v *entity.Version) error
	Delete(ctx context.Context, id uint64) error
}

type versionPostgres struct{ db *gorm.DB }

func NewVersionPostgres(db *gorm.DB) VersionRepository {
	return &versionPostgres{db: db}
}

func (r *versionPostgres) GetLatest(ctx context.Context, appName, platform string, activeOnly bool) (*entity.Version, error) {
	if r.db == nil {
		return nil, nil
	}
	q := r.db.WithContext(ctx).
		Where("app_name = ? AND platform = ? AND deleted_at IS NULL", appName, platform)
	if activeOnly {
		q = q.Where("active = true")
	}
	var row entity.Version
	err := q.Order("build_number DESC, id DESC").First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *versionPostgres) Create(ctx context.Context, v *entity.Version) error {
	if r.db == nil || v == nil {
		return nil
	}
	var count int64
	if err := r.db.WithContext(ctx).Model(&entity.Version{}).
		Where("app_name = ? AND platform = ? AND version_name = ? AND build_number = ? AND deleted_at IS NULL",
			v.AppName, v.Platform, v.VersionName, v.BuildNumber).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("phiên bản đã tồn tại")
	}
	return r.db.WithContext(ctx).Create(v).Error
}

func (r *versionPostgres) GetByID(ctx context.Context, id uint64) (*entity.Version, error) {
	if r.db == nil {
		return nil, nil
	}
	var row entity.Version
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &row, err
}

func (r *versionPostgres) List(ctx context.Context, appName, platform string, page, size int) ([]*entity.Version, int64, error) {
	if r.db == nil {
		return nil, 0, nil
	}
	q := r.db.WithContext(ctx).Model(&entity.Version{}).Where("deleted_at IS NULL")
	if appName != "" {
		q = q.Where("app_name = ?", appName)
	}
	if platform != "" {
		q = q.Where("platform = ?", platform)
	}
	var total int64
	q.Count(&total)
	offset, size := pageOffset(page, size)
	var rows []*entity.Version
	err := q.Order("build_number DESC, id DESC").Offset(offset).Limit(size).Find(&rows).Error
	return rows, total, err
}

func (r *versionPostgres) Update(ctx context.Context, v *entity.Version) error {
	if r.db == nil {
		return nil
	}
	return r.db.WithContext(ctx).Save(v).Error
}

func (r *versionPostgres) Delete(ctx context.Context, id uint64) error {
	if r.db == nil {
		return nil
	}
	return r.db.WithContext(ctx).Delete(&entity.Version{}, id).Error
}

func (r *versionPostgres) ValidateAPIKey(ctx context.Context, apiKey string) error {
	if r.db == nil {
		return errors.New("api key store unavailable")
	}
	if apiKey == "" {
		return errors.New("missing api key")
	}
	var row entity.APIKey
	err := r.db.WithContext(ctx).
		Where("api_key = ? AND deleted_at IS NULL", apiKey).
		First(&row).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("invalid api key")
	}
	if err != nil {
		return err
	}
	if row.ExpiredAt != nil && row.ExpiredAt.Before(time.Now()) {
		return errors.New("api key expired")
	}
	return nil
}
