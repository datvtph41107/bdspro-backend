// file/internal/infrastructure/persistence/gorm/access_repository.go
package gorm

import (
	"context"
	"file/internal/core/domain"
	"file/internal/ports"
	"time"

	"gorm.io/gorm"
)

type accessRepository struct {
	db *gorm.DB
}

func NewAccessRepository(db *gorm.DB) ports.FileAccessRepository {
	return &accessRepository{db: db}
}

func (r *accessRepository) GrantAccess(ctx context.Context, access *domain.FileAccess) error {
	// Check if access already exists
	var existing domain.FileAccess
	err := r.db.WithContext(ctx).
		Where("file_id = ? AND user_id = ?", access.FileID, access.UserID).
		First(&existing).Error

	if err == nil {
		// Update existing access
		existing.Permission = access.Permission
		existing.ExpiresAt = access.ExpiresAt
		return r.db.WithContext(ctx).Save(&existing).Error
	}

	// Create new access
	return r.db.WithContext(ctx).Create(access).Error
}

func (r *accessRepository) RevokeAccess(ctx context.Context, fileID, userID uint64) error {
	return r.db.WithContext(ctx).
		Where("file_id = ? AND user_id = ?", fileID, userID).
		Delete(&domain.FileAccess{}).Error
}

func (r *accessRepository) HasAccess(ctx context.Context, fileID, userID uint64, permission string) (bool, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(&domain.FileAccess{}).
		Where("file_id = ? AND user_id = ?", fileID, userID)

	if permission != "" {
		query = query.Where("permission = ?", permission)
	}

	// Check expiry
	query = query.Where("expires_at > ? OR expires_at IS NULL", time.Now())

	err := query.Count(&count).Error
	return count > 0, err
}

func (r *accessRepository) GetAccessibleFiles(ctx context.Context, userID uint64, permission string) ([]uint64, error) {
	var fileIDs []uint64

	query := r.db.WithContext(ctx).Model(&domain.FileAccess{}).
		Select("file_id").
		Where("user_id = ?", userID).
		Where("expires_at > ? OR expires_at IS NULL", time.Now())

	if permission != "" {
		query = query.Where("permission = ?", permission)
	}

	err := query.Pluck("file_id", &fileIDs).Error
	return fileIDs, err
}

func (r *accessRepository) GetFileAccesses(ctx context.Context, fileID uint64) ([]*domain.FileAccess, error) {
	var accesses []*domain.FileAccess
	err := r.db.WithContext(ctx).
		Where("file_id = ?", fileID).
		Find(&accesses).Error
	return accesses, err
}

func (r *accessRepository) CleanupExpiredAccesses(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&domain.FileAccess{}).Error
}
