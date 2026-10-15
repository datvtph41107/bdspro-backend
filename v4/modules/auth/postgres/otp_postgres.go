package postgres

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"userv4/modules/auth/domain/auth"
)

// OTPPostgres implementation của OTP repository
// @bind: user/internal/interface/repo.OTPRepository
type OTPPostgres struct {
	db *gorm.DB
}

// NewOTPRepository tạo mới OTPRepository
func NewOTPRepository(db *gorm.DB) *OTPPostgres {
	return &OTPPostgres{db: db}
}

// GetByID lấy OTP theo authID. Nếu bản ghi bị miss thì tạo skeleton mới.
func (r *OTPPostgres) GetByID(ctx context.Context, authID uint64) (*auth.UserOTPEntity, error) {
	var otp auth.UserOTPEntity
	err := r.db.WithContext(ctx).Where("auth_id = ?", authID).First(&otp).Error
	if err == nil {
		return &otp, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	otp = auth.UserOTPEntity{
		IdDomain: auth.IdDomain{AuthID: authID},
	}
	if err := r.db.WithContext(ctx).Create(&otp).Error; err != nil {
		// Concurrent create: đọc lại bản ghi vừa được tạo.
		if readErr := r.db.WithContext(ctx).Where("auth_id = ?", authID).First(&otp).Error; readErr != nil {
			return nil, err
		}
		return &otp, nil
	}
	return &otp, nil
}

// CreateOTP tạo mới OTP
func (r *OTPPostgres) CreateOTP(ctx context.Context, otp *auth.UserOTPEntity) error {
	return r.db.WithContext(ctx).Create(otp).Error
}

// UpdateOTP cập nhật OTP, nếu chưa có thì tạo mới
func (r *OTPPostgres) UpdateOTP(ctx context.Context, otp *auth.UserOTPEntity) error {
	var existingOTP auth.UserOTPEntity
	err := r.db.WithContext(ctx).Where("auth_id = ?", otp.AuthID).First(&existingOTP).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return r.db.WithContext(ctx).Create(otp).Error
		}
		return err
	}

	return r.db.WithContext(ctx).Model(&auth.UserOTPEntity{}).
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
