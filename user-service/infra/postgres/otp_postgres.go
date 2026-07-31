package postgres

import (
	"context"
	"time"

	_db "common/db"
	"user/internal/domain"

	"gorm.io/gorm"
)

// OTPPostgres implementation của OTP repository
// @bind: auth/internal/interface/repo.OTPRepository
type OTPPostgres struct {
	// TODO: Add database connection
	db *gorm.DB
}

// NewOTPRepository tạo mới OTPRepository
func NewOTPRepository(db *gorm.DB) *OTPPostgres {
	return &OTPPostgres{db: db}
}

// GetByID lấy OTP theo authID
func (r *OTPPostgres) GetByID(ctx context.Context, authID uint64) (*domain.UserOTPEntity, error) {
	var otp domain.UserOTPEntity
	err := _db.DB.Where("auth_id = ?", authID).First(&otp).Error
	if err != nil {
		return nil, err
	}
	return &otp, nil
}

// CreateOTP tạo mới OTP
func (r *OTPPostgres) CreateOTP(ctx context.Context, otp *domain.UserOTPEntity) error {
	return _db.DB.Create(otp).Error
}

// UpdateOTP cập nhật OTP, nếu chưa có thì tạo mới
func (r *OTPPostgres) UpdateOTP(ctx context.Context, otp *domain.UserOTPEntity) error {
	// Kiểm tra xem OTP đã tồn tại chưa
	var existingOTP domain.UserOTPEntity
	err := _db.DB.Where("auth_id = ?", otp.AuthID).First(&existingOTP).Error

	if err != nil {
		// Nếu không tìm thấy, tạo mới
		if err == gorm.ErrRecordNotFound {
			return _db.DB.Create(otp).Error
		}
		return err
	}

	// Nếu đã tồn tại, cập nhật
	return r.db.Model(&domain.UserOTPEntity{}).
		Where("auth_id = ?", otp.AuthID).
		Updates(map[string]interface{}{
			"otp":            otp.OTP,
			"otp_send_time":  otp.OTPSendTime,
			"otp_check_time": otp.OTPCheckTime,
			"otp_date":       otp.OTPDate,
			"expired_time":   otp.ExpiredTime,
			"activate":       otp.Activate,
		}).
		Error
}

// DeleteOTP xóa OTP
func (r *OTPPostgres) DeleteOTP(ctx context.Context, authID uint64) error {
	return _db.DB.Where("auth_id = ?", authID).Delete(&domain.UserOTPEntity{}).Error
}

// GetOTPByPhone lấy OTP theo số điện thoại
func (r *OTPPostgres) GetOTPByPhone(ctx context.Context, phone string) (*domain.UserOTPEntity, error) {
	var otp domain.UserOTPEntity
	err := _db.DB.Joins("JOIN auth_method ON user_otp_entities.auth_id = auth_method.id").
		Where("auth_method.phone = ? AND auth_method.provider = ?", phone, "PHONE").
		First(&otp).Error
	if err != nil {
		return nil, err
	}
	return &otp, nil
}

// GetExpiredOTPs lấy danh sách OTP đã hết hạn
func (r *OTPPostgres) GetExpiredOTPs(ctx context.Context, before time.Time) ([]*domain.UserOTPEntity, error) {
	var otps []*domain.UserOTPEntity
	err := _db.DB.Where("expired_time < ?", before).Find(&otps).Error
	if err != nil {
		return nil, err
	}
	return otps, nil
}

// CleanupExpiredOTPs xóa các OTP đã hết hạn
func (r *OTPPostgres) CleanupExpiredOTPs(ctx context.Context, before time.Time) error {
	return _db.DB.Where("expired_time < ?", before).Delete(&domain.UserOTPEntity{}).Error
}

// MarkOTPAsUsed đánh dấu OTP đã được sử dụng
func (r *OTPPostgres) MarkOTPAsUsed(ctx context.Context, authID uint64) error {
	return _db.DB.Model(&domain.UserOTPEntity{}).
		Where("auth_id = ?", authID).
		Update("activate", true).Error
}

// ResetOTPCounters reset các counter của OTP
func (r *OTPPostgres) ResetOTPCounters(ctx context.Context, authID uint64) error {
	return _db.DB.Model(&domain.UserOTPEntity{}).
		Where("auth_id = ?", authID).
		Updates(map[string]interface{}{
			"otp_send_time":  0,
			"otp_check_time": 0,
		}).Error
}
