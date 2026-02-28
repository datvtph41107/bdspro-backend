// file/internal/infrastructure/persistence/gorm/file_repository.go
package gorm

import (
	"context"
	"fmt"
	"time"

	"file/internal/core/domain"
	"file/internal/ports"

	"gorm.io/gorm"
)

type fileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) ports.FileRepository {
	return &fileRepository{db: db}
}

func (r *fileRepository) Create(ctx context.Context, file *domain.File) error {
	return r.db.WithContext(ctx).Create(file).Error
}

func (r *fileRepository) FindByID(ctx context.Context, id uint64) (*domain.File, error) {
	var file domain.File
	err := r.db.WithContext(ctx).First(&file, id).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *fileRepository) FindByUUID(ctx context.Context, uuid string) (*domain.File, error) {
	var file domain.File
	err := r.db.WithContext(ctx).Where("uuid = ?", uuid).First(&file).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *fileRepository) FindByHash(ctx context.Context, hash string) (*domain.File, error) {
	var file domain.File
	err := r.db.WithContext(ctx).Where("hash = ? AND status = ?", hash, domain.FileStatusActive).First(&file).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *fileRepository) FindByReference(ctx context.Context, serviceID, referenceType, referenceID string) ([]*domain.File, error) {
	var files []*domain.File
	query := r.db.WithContext(ctx).Where("status = ?", domain.FileStatusActive)

	if serviceID != "" {
		query = query.Where("service_id = ?", serviceID)
	}
	if referenceType != "" {
		query = query.Where("reference_type = ?", referenceType)
	}
	if referenceID != "" {
		query = query.Where("reference_id = ?", referenceID)
	}

	err := query.Find(&files).Error
	return files, err
}

func (r *fileRepository) Update(ctx context.Context, file *domain.File) error {
	return r.db.WithContext(ctx).Save(file).Error
}

func (r *fileRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Model(&domain.File{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":     domain.FileStatusDeleted,
			"deleted_at": time.Now(),
		}).Error
}

func (r *fileRepository) List(ctx context.Context, filter ports.FileFilter, pagination ports.Pagination) ([]*domain.File, int64, error) {
	var files []*domain.File
	var total int64

	query := r.db.WithContext(ctx).Model(&domain.File{}).Where("status = ?", domain.FileStatusActive)

	// Apply filters
	if filter.ServiceID != "" {
		query = query.Where("service_id = ?", filter.ServiceID)
	}
	if filter.ReferenceType != "" {
		query = query.Where("reference_type = ?", filter.ReferenceType)
	}
	if filter.ReferenceID != "" {
		query = query.Where("reference_id = ?", filter.ReferenceID)
	}
	if filter.MimeType != "" {
		query = query.Where("mime_type LIKE ?", filter.MimeType+"%")
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.UploadedBy > 0 {
		query = query.Where("uploaded_by = ?", filter.UploadedBy)
	}
	if !filter.StartDate.IsZero() {
		query = query.Where("created_at >= ?", filter.StartDate)
	}
	if !filter.EndDate.IsZero() {
		query = query.Where("created_at <= ?", filter.EndDate)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (pagination.Page - 1) * pagination.PageSize
	query = query.Offset(offset).Limit(pagination.PageSize)

	// Apply sorting
	if pagination.SortBy != "" {
		order := fmt.Sprintf("%s %s", pagination.SortBy, pagination.SortDir)
		query = query.Order(order)
	} else {
		query = query.Order("created_at DESC")
	}

	err := query.Find(&files).Error
	return files, total, err
}

// Additional methods for better querying
func (r *fileRepository) FindByServiceAndReference(ctx context.Context, serviceID, referenceType string) ([]*domain.File, error) {
	var files []*domain.File
	err := r.db.WithContext(ctx).
		Where("service_id = ? AND reference_type = ? AND status = ?",
			serviceID, referenceType, domain.FileStatusActive).
		Find(&files).Error
	return files, err
}

func (r *fileRepository) FindExpiredFiles(ctx context.Context, expiryTime time.Time) ([]*domain.File, error) {
	var files []*domain.File
	err := r.db.WithContext(ctx).
		Where("expires_at < ? AND status = ?", expiryTime, domain.FileStatusActive).
		Find(&files).Error
	return files, err
}

func (r *fileRepository) UpdateThumbnail(ctx context.Context, fileID uint64, thumbnailPath string) error {
	return r.db.WithContext(ctx).
		Model(&domain.File{}).
		Where("id = ?", fileID).
		Update("thumbnail_path", thumbnailPath).Error
}

func (r *fileRepository) UpdateStatus(ctx context.Context, fileID uint64, status string) error {
	return r.db.WithContext(ctx).
		Model(&domain.File{}).
		Where("id = ?", fileID).
		Update("status", status).Error
}
