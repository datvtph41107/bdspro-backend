package postgres

import (
	"context"
	"time"
	accessdomain "user/internal/domain/access"

	"gorm.io/gorm"
)

// @bind: user/internal/interface/repo.AdminAccessRepository
type AdminAccessPostgresRepo struct {
	db *gorm.DB
}

func NewAdminAccessPostgresRepo(db *gorm.DB) *AdminAccessPostgresRepo {
	return &AdminAccessPostgresRepo{db: db}
}

// Create tạo quyền truy cập mới
func (r *AdminAccessPostgresRepo) Create(c context.Context, access *accessdomain.AdminAccessDomain) (*accessdomain.AdminAccessDomain, error) {
	err := r.db.WithContext(c).Create(access).Error
	if err != nil {
		return nil, err
	}
	return access, nil
}

// Update cập nhật quyền truy cập
func (r *AdminAccessPostgresRepo) Update(c context.Context, access *accessdomain.AdminAccessDomain) (*accessdomain.AdminAccessDomain, error) {
	err := r.db.WithContext(c).Save(access).Error
	if err != nil {
		return nil, err
	}
	return access, nil
}

// Delete xóa quyền truy cập
func (r *AdminAccessPostgresRepo) Delete(c context.Context, id uint64) error {
	return r.db.WithContext(c).Delete(&accessdomain.AdminAccessDomain{}, id).Error
}

// SoftDelete soft delete quyền truy cập
func (r *AdminAccessPostgresRepo) SoftDelete(c context.Context, id uint64) error {
	now := time.Now()
	return r.db.WithContext(c).Model(&accessdomain.AdminAccessDomain{}).Where("id = ?", id).Update("deleted_at", now).Error
}

// FindByID tìm quyền truy cập theo ID
func (r *AdminAccessPostgresRepo) FindByID(c context.Context, id uint64) (*accessdomain.AdminAccessDomain, error) {
	var access accessdomain.AdminAccessDomain
	err := r.db.WithContext(c).Where("id = ? AND deleted_at IS NULL", id).First(&access).Error
	if err != nil {
		return nil, err
	}
	return &access, nil
}

// FindAll lấy danh sách quyền truy cập với filter
func (r *AdminAccessPostgresRepo) FindAll(c context.Context, page, size int, filters map[string]interface{}) ([]*accessdomain.AdminAccessDomain, int64, error) {
	var accesses []*accessdomain.AdminAccessDomain
	var total int64

	query := r.db.WithContext(c).Model(&accessdomain.AdminAccessDomain{}).Where("deleted_at IS NULL")

	// Apply filters
	if userID, ok := filters["user_id"]; ok {
		query = query.Where("user_id = ?", userID)
	}
	if status, ok := filters["status"]; ok {
		query = query.Where("status = ?", status)
	}
	if authType, ok := filters["auth_type"]; ok {
		query = query.Where("auth_type = ?", authType)
	}
	if ipAddress, ok := filters["ip_address"]; ok {
		query = query.Where("ip_address LIKE ?", "%"+ipAddress.(string)+"%")
	}

	// Count total
	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * size
	err = query.Order("created_at DESC").Offset(offset).Limit(size).Find(&accesses).Error
	if err != nil {
		return nil, 0, err
	}

	return accesses, total, nil
}

// FindByUserID tìm quyền truy cập theo UserID
func (r *AdminAccessPostgresRepo) FindByUserID(c context.Context, userID uint64) ([]*accessdomain.AdminAccessDomain, error) {
	var accesses []*accessdomain.AdminAccessDomain
	err := r.db.WithContext(c).Where("user_id = ? AND deleted_at IS NULL", userID).Find(&accesses).Error
	if err != nil {
		return nil, err
	}
	return accesses, nil
}

// FindByIPAddress tìm quyền truy cập theo IP
func (r *AdminAccessPostgresRepo) FindByIPAddress(c context.Context, ipAddress string) ([]*accessdomain.AdminAccessDomain, error) {
	var accesses []*accessdomain.AdminAccessDomain
	err := r.db.WithContext(c).Where("(ip_address = ? OR ip_range LIKE ?) AND deleted_at IS NULL", ipAddress, "%"+ipAddress+"%").Find(&accesses).Error
	if err != nil {
		return nil, err
	}
	return accesses, nil
}

// FindByDeviceID tìm quyền truy cập theo DeviceID
func (r *AdminAccessPostgresRepo) FindByDeviceID(c context.Context, deviceID string) ([]*accessdomain.AdminAccessDomain, error) {
	var accesses []*accessdomain.AdminAccessDomain
	err := r.db.WithContext(c).Where("device_id = ? AND deleted_at IS NULL", deviceID).Find(&accesses).Error
	if err != nil {
		return nil, err
	}
	return accesses, nil
}

// ValidateAccess kiểm tra quyền truy cập
func (r *AdminAccessPostgresRepo) ValidateAccess(c context.Context, userID uint64, ipAddress, deviceID string) (*accessdomain.AdminAccessValidationResult, error) {
	var access accessdomain.AdminAccessDomain

	// Tìm quyền truy cập phù hợp
	query := r.db.WithContext(c).Where("user_id = ? AND status = 'active' AND deleted_at IS NULL", userID)

	// Kiểm tra IP hoặc device
	if ipAddress != "" {
		query = query.Where("(ip_address = ? OR ip_range LIKE ?)", ipAddress, "%"+ipAddress+"%")
	}
	if deviceID != "" {
		query = query.Where("device_id = ?", deviceID)
	}

	err := query.First(&access).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &accessdomain.AdminAccessValidationResult{
				IsAllowed:  false,
				Reason:     "IP hoặc thiết bị không được phép truy cập",
				Require2FA: true,
			}, nil
		}
		return nil, err
	}

	// Kiểm tra thời gian hiệu lực
	now := time.Now()
	if access.EffectiveFrom != nil && now.Before(*access.EffectiveFrom) {
		return &accessdomain.AdminAccessValidationResult{
			IsAllowed:  false,
			Reason:     "Quyền truy cập chưa có hiệu lực",
			Require2FA: true,
		}, nil
	}

	if access.EffectiveTo != nil && now.After(*access.EffectiveTo) {
		return &accessdomain.AdminAccessValidationResult{
			IsAllowed:  false,
			Reason:     "Quyền truy cập đã hết hạn",
			Require2FA: true,
		}, nil
	}

	return &accessdomain.AdminAccessValidationResult{
		IsAllowed:  true,
		Reason:     "Truy cập được phép",
		Require2FA: access.Require2FA,
		AccessID:   access.ID,
		DeviceID:   access.DeviceID,
		IPAddress:  access.IPAddress,
	}, nil
}

// CountByUserID đếm số quyền truy cập theo UserID
func (r *AdminAccessPostgresRepo) CountByUserID(c context.Context, userID uint64) (int64, error) {
	var count int64
	err := r.db.WithContext(c).Model(&accessdomain.AdminAccessDomain{}).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Count(&count).Error
	return count, err
}

// CountByIPAddress đếm số quyền truy cập theo IP
func (r *AdminAccessPostgresRepo) CountByIPAddress(c context.Context, ipAddress string) (int64, error) {
	var count int64
	err := r.db.WithContext(c).Model(&accessdomain.AdminAccessDomain{}).
		Where("(ip_address = ? OR ip_range LIKE ?) AND deleted_at IS NULL", ipAddress, "%"+ipAddress+"%").
		Count(&count).Error
	return count, err
}

// CreateLog tạo log truy cập
func (r *AdminAccessPostgresRepo) CreateLog(c context.Context, log *accessdomain.AdminAccessLogDomain) error {
	return r.db.WithContext(c).Create(log).Error
}

// FindLogsByUserID tìm log theo UserID
func (r *AdminAccessPostgresRepo) FindLogsByUserID(c context.Context, userID uint64, page, size int) ([]*accessdomain.AdminAccessLogDomain, int64, error) {
	var logs []*accessdomain.AdminAccessLogDomain
	var total int64

	err := r.db.WithContext(c).Model(&accessdomain.AdminAccessLogDomain{}).
		Where("user_id = ?", userID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	err = r.db.WithContext(c).Where("user_id = ?", userID).
		Order("created_at DESC").
		Offset(offset).
		Limit(size).
		Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// FindLogsByIPAddress tìm log theo IP
func (r *AdminAccessPostgresRepo) FindLogsByIPAddress(c context.Context, ipAddress string, page, size int) ([]*accessdomain.AdminAccessLogDomain, int64, error) {
	var logs []*accessdomain.AdminAccessLogDomain
	var total int64

	err := r.db.WithContext(c).Model(&accessdomain.AdminAccessLogDomain{}).
		Where("ip_address = ?", ipAddress).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	err = r.db.WithContext(c).Where("ip_address = ?", ipAddress).
		Order("created_at DESC").
		Offset(offset).
		Limit(size).
		Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
