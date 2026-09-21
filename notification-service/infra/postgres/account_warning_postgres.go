package postgres

import (
	"context"
	"errors"
	"time"

	"notification/internal/domain"
	"notification/internal/enums"

	"gorm.io/gorm"
)

// @bind: notification/internal/usecase.AccountWarningStore
type AccountWarningPostgresRepo struct{ db *gorm.DB }

func NewAccountWarningPostgresRepo(db *gorm.DB) *AccountWarningPostgresRepo {
	return &AccountWarningPostgresRepo{db: db}
}

// Create tạo cảnh báo mới
func (r *AccountWarningPostgresRepo) Create(c context.Context, warning *domain.AccountWarningEntity) (*domain.AccountWarningEntity, error) {
	err := r.db.WithContext(c).Create(warning).Error
	if err != nil {
		return nil, err
	}
	return warning, nil
}

// Update cập nhật cảnh báo
func (r *AccountWarningPostgresRepo) Update(c context.Context, warning *domain.AccountWarningEntity) (*domain.AccountWarningEntity, error) {
	err := r.db.WithContext(c).Save(warning).Error
	if err != nil {
		return nil, err
	}
	return warning, nil
}

// Delete xóa cảnh báo
func (r *AccountWarningPostgresRepo) Delete(c context.Context, id uint64) error {
	return r.db.WithContext(c).Delete(&domain.AccountWarningEntity{}, id).Error
}

// SoftDelete soft delete cảnh báo
func (r *AccountWarningPostgresRepo) SoftDelete(c context.Context, id uint64) error {
	now := time.Now()
	return r.db.WithContext(c).Model(&domain.AccountWarningEntity{}).Where("id = ?", id).Update("deleted_at", now).Error
}

// FindByID tìm cảnh báo theo ID
func (r *AccountWarningPostgresRepo) FindByID(c context.Context, id uint64) (*domain.AccountWarningEntity, error) {
	var warning domain.AccountWarningEntity
	err := r.db.WithContext(c).Where("id = ? AND deleted_at IS NULL", id).First(&warning).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &warning, nil
}

// FindAll lấy danh sách cảnh báo với filter
func (r *AccountWarningPostgresRepo) FindAll(c context.Context, page, size int, filters map[string]interface{}) ([]*domain.AccountWarningEntity, int64, error) {
	var warnings []*domain.AccountWarningEntity
	var total int64

	query := r.db.WithContext(c).Model(&domain.AccountWarningEntity{}).Where("deleted_at IS NULL")

	// Apply filters
	if targetID, ok := filters["target_id"]; ok {
		query = query.Where("target_id = ?", targetID)
	}
	if targetType, ok := filters["target_type"]; ok {
		query = query.Where("target_type = ?", targetType)
	}
	if warningType, ok := filters["warning_type"]; ok {
		query = query.Where("warning_type = ?", warningType)
	}
	if severity, ok := filters["severity"]; ok {
		query = query.Where("severity = ?", severity)
	}
	if status, ok := filters["status"]; ok {
		query = query.Where("status = ?", status)
	}
	if sentBy, ok := filters["sent_by"]; ok {
		query = query.Where("sent_by = ?", sentBy)
	}
	if isRead, ok := filters["is_read"]; ok {
		query = query.Where("is_read = ?", isRead)
	}

	// Count total
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := page * size
	err = query.Debug().Order("sent_at DESC").Offset(offset).Limit(size).Find(&warnings).Error
	if err != nil {
		return nil, 0, err
	}

	return warnings, total, nil
}

// FindByTargetID tìm cảnh báo theo TargetID
func (r *AccountWarningPostgresRepo) FindByTargetID(c context.Context, targetID uint64, targetType string) ([]*domain.AccountWarningEntity, error) {
	var warnings []*domain.AccountWarningEntity
	err := r.db.WithContext(c).Where("target_id = ? AND target_type = ? AND deleted_at IS NULL", targetID, targetType).Order("sent_at DESC").Find(&warnings).Error
	if err != nil {
		return nil, err
	}
	return warnings, nil
}

// FindBySentBy tìm cảnh báo theo người gửi
func (r *AccountWarningPostgresRepo) FindBySentBy(c context.Context, sentBy uint64, page, size int) ([]*domain.AccountWarningEntity, int64, error) {
	var warnings []*domain.AccountWarningEntity
	var total int64

	err := r.db.WithContext(c).Model(&domain.AccountWarningEntity{}).
		Where("sent_by = ? AND deleted_at IS NULL", sentBy).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	err = r.db.WithContext(c).Where("sent_by = ? AND deleted_at IS NULL", sentBy).
		Order("sent_at DESC").
		Offset(offset).
		Limit(size).
		Find(&warnings).Error
	if err != nil {
		return nil, 0, err
	}

	return warnings, total, nil
}

// FindUnreadByTargetID tìm cảnh báo chưa đọc theo TargetID
func (r *AccountWarningPostgresRepo) FindUnreadByTargetID(c context.Context, targetID uint64) ([]*domain.AccountWarningEntity, error) {
	var warnings []*domain.AccountWarningEntity
	err := r.db.WithContext(c).Where("target_id = ? AND is_read = ? AND deleted_at IS NULL", targetID, false).Order("sent_at DESC").Find(&warnings).Error
	if err != nil {
		return nil, err
	}
	return warnings, nil
}

// MarkAsRead đánh dấu đã đọc
func (r *AccountWarningPostgresRepo) MarkAsRead(c context.Context, id uint64, targetID uint64) error {
	now := time.Now()
	return r.db.WithContext(c).Model(&domain.AccountWarningEntity{}).
		Where("id = ? AND target_id = ? AND deleted_at IS NULL", id, targetID).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": now,
			"status":  enums.WarningStatusRead,
		}).Error
}

// MarkAsAcknowledged đánh dấu đã xác nhận
func (r *AccountWarningPostgresRepo) MarkAsAcknowledged(c context.Context, id uint64, targetID uint64) error {
	return r.db.WithContext(c).Model(&domain.AccountWarningEntity{}).
		Where("id = ? AND target_id = ? AND deleted_at IS NULL", id, targetID).
		Update("status", enums.WarningStatusAcknowledged).Error
}

// CheckDuplicateWarning kiểm tra cảnh báo trùng lặp
func (r *AccountWarningPostgresRepo) CheckDuplicateWarning(c context.Context, targetID uint64, content string, hours int) (bool, error) {
	var count int64
	timeLimit := time.Now().Add(-time.Duration(hours) * time.Hour)

	err := r.db.WithContext(c).Model(&domain.AccountWarningEntity{}).
		Where("target_id = ? AND content = ? AND sent_at > ? AND deleted_at IS NULL", targetID, content, timeLimit).
		Count(&count).Error

	return count > 0, err
}

// CountByTargetID đếm số cảnh báo theo TargetID
func (r *AccountWarningPostgresRepo) CountByTargetID(c context.Context, targetID uint64) (int64, error) {
	var count int64
	err := r.db.WithContext(c).Model(&domain.AccountWarningEntity{}).
		Where("target_id = ? AND deleted_at IS NULL", targetID).
		Count(&count).Error
	return count, err
}

// CreateTemplate tạo mẫu cảnh báo
func (r *AccountWarningPostgresRepo) CreateTemplate(c context.Context, template *domain.AccountWarningTemplateEntity) (*domain.AccountWarningTemplateEntity, error) {
	err := r.db.WithContext(c).Create(template).Error
	if err != nil {
		return nil, err
	}
	return template, nil
}

// UpdateTemplate cập nhật mẫu cảnh báo
func (r *AccountWarningPostgresRepo) UpdateTemplate(c context.Context, template *domain.AccountWarningTemplateEntity) (*domain.AccountWarningTemplateEntity, error) {
	err := r.db.WithContext(c).Save(template).Error
	if err != nil {
		return nil, err
	}
	return template, nil
}

// DeleteTemplate xóa mẫu cảnh báo
func (r *AccountWarningPostgresRepo) DeleteTemplate(c context.Context, id uint64) error {
	now := time.Now()
	return r.db.WithContext(c).Model(&domain.AccountWarningTemplateEntity{}).Where("id = ?", id).Update("deleted_at", now).Error
}

// FindTemplateByID tìm mẫu cảnh báo theo ID
func (r *AccountWarningPostgresRepo) FindTemplateByID(c context.Context, id uint64) (*domain.AccountWarningTemplateEntity, error) {
	var template domain.AccountWarningTemplateEntity
	err := r.db.WithContext(c).Where("id = ? AND deleted_at IS NULL", id).First(&template).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &template, nil
}

// FindTemplatesByType tìm mẫu cảnh báo theo loại
func (r *AccountWarningPostgresRepo) FindTemplatesByType(c context.Context, warningType string) ([]*domain.AccountWarningTemplateEntity, error) {
	var templates []*domain.AccountWarningTemplateEntity
	err := r.db.WithContext(c).Where("warning_type = ? AND is_active = ? AND deleted_at IS NULL", warningType, true).Find(&templates).Error
	if err != nil {
		return nil, err
	}
	return templates, nil
}

// FindAllTemplates lấy tất cả mẫu cảnh báo
func (r *AccountWarningPostgresRepo) FindAllTemplates(c context.Context, page, size int) ([]*domain.AccountWarningTemplateEntity, int64, error) {
	var templates []*domain.AccountWarningTemplateEntity
	var total int64

	err := r.db.WithContext(c).Model(&domain.AccountWarningTemplateEntity{}).
		Where("deleted_at IS NULL").
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	err = r.db.WithContext(c).Where("deleted_at IS NULL").
		Order("created_at DESC").
		Offset(offset).
		Limit(size).
		Find(&templates).Error
	if err != nil {
		return nil, 0, err
	}

	return templates, total, nil
}

// CreateLog tạo log cảnh báo
func (r *AccountWarningPostgresRepo) CreateLog(c context.Context, log *domain.AccountWarningLogEntity) error {
	return r.db.WithContext(c).Create(log).Error
}

// FindLogsByWarningID tìm log theo WarningID
func (r *AccountWarningPostgresRepo) FindLogsByWarningID(c context.Context, warningID uint64, page, size int) ([]*domain.AccountWarningLogEntity, int64, error) {
	var logs []*domain.AccountWarningLogEntity
	var total int64

	err := r.db.WithContext(c).Model(&domain.AccountWarningLogEntity{}).
		Where("warning_id = ?", warningID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	err = r.db.WithContext(c).Where("warning_id = ?", warningID).
		Order("performed_at DESC").
		Offset(offset).
		Limit(size).
		Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// FindLogsByTargetID tìm log theo TargetID
func (r *AccountWarningPostgresRepo) FindLogsByTargetID(c context.Context, targetID uint64, page, size int) ([]*domain.AccountWarningLogEntity, int64, error) {
	var logs []*domain.AccountWarningLogEntity
	var total int64

	err := r.db.WithContext(c).Model(&domain.AccountWarningLogEntity{}).
		Where("target_id = ?", targetID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	err = r.db.WithContext(c).Where("target_id = ?", targetID).
		Order("performed_at DESC").
		Offset(offset).
		Limit(size).
		Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
