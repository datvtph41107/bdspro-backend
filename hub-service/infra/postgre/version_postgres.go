package postgres

import (
	"context"
	"errors"

	_db "common/db"
	_dto "common/domain/dto"
	_provider "common/provider"
	"hub/internal/domain"
	"hub/internal/repo"

	"gorm.io/gorm"
)

type VersionPostgres struct {
	_provider.CrudRepo[domain.VersionEntity]
}

func NewVersionRepo(db *_db.TransactionRepo) repo.IVersionRepo {
	repo := &VersionPostgres{}
	repo.Init(repo, db)
	return repo
}

func (r *VersionPostgres) GetListWithFilter(ctx context.Context, appName string, platform string, pagable _dto.IPagable) ([]*domain.VersionEntity, int64, error) {
	var (
		versions []*domain.VersionEntity
		total    int64
	)

	db := r.GetDB(ctx).Model(&domain.VersionEntity{}).Where("deleted_at IS NULL")

	if appName != "" {
		db = db.Where("app_name = ?", appName)
	}

	if platform != "" {
		db = db.Where("platform = ?", platform)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Order("created_at DESC").Offset(pagable.GetOffset()).Limit(pagable.GetLimit()).Find(&versions).Error; err != nil {
		return nil, 0, err
	}

	return versions, total, nil
}

func (r *VersionPostgres) GetLatestByAppAndPlatform(ctx context.Context, appName string, platform string, versionName string, active bool) (*domain.VersionEntity, error) {
	var version domain.VersionEntity

	query := r.GetDB(ctx).
		Where("app_name = ? AND platform = ? AND version_name = ? AND deleted_at IS NULL", appName, platform, versionName)

	if active {
		query = query.Where("active = true")
	}
	query = query.Order("build_number DESC, id DESC")

	if err := query.First(&version).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &version, nil
}

func (r *VersionPostgres) BeforeSave(ctx context.Context, id *uint64, entity *domain.VersionEntity) error {
	if entity.AppName == "" {
		return errors.New("appName không được để trống")
	}

	if entity.Platform == "" {
		return errors.New("platform không được để trống")
	}

	if entity.VersionName == "" {
		return errors.New("versionName không được để trống")
	}

	if entity.BuildNumber == 0 {
		return errors.New("buildNumber không được để trống")
	}

	query := r.GetDB(ctx).
		Model(&domain.VersionEntity{}).
		Where("app_name = ? AND platform = ? AND version_name = ? AND build_number = ? AND deleted_at IS NULL", entity.AppName, entity.Platform, entity.VersionName, entity.BuildNumber)

	if id != nil {
		query = query.Where("id <> ?", *id)
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return errors.New("phiên bản đã tồn tại")
	}

	return nil
}

func (r *VersionPostgres) Update(ctx context.Context, id uint64, entity *domain.VersionEntity) error {
	err := r.GetDB(ctx).
		Model(&domain.VersionEntity{}).
		Where("id = ? and deleted_at is null", id).
		Updates(map[string]interface{}{
			"force_update":  entity.ForceUpdate,
			"active":        entity.Active,
			"release_notes": entity.ReleaseNotes,
		}).Error
	if err != nil {
		return err
	}
	return nil
}
