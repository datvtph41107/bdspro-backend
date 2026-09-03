package postgres

import (
	"context"

	"gorm.io/gorm"
	"user/internal/domain/auth"
)

// OTPPostgres implementation của OTP repository
// @bind: user/internal/interface/repo.OTPRepository
type OTPPostgres struct {
	// TODO: Add database connection
	db *gorm.DB
}

// NewOTPRepository tạo mới OTPRepository
func NewOTPRepository(db *gorm.DB) *OTPPostgres {
	return &OTPPostgres{db: db}
}

// GetByID lấy OTP theo authID
func (r *OTPPostgres) GetByID(ctx context.Context, authID uint64) (*auth.UserOTPEntity, error) {
	var otp auth.UserOTPEntity
	err := r.db.Where("auth_id = ?", authID).First(&otp).Error
	if err != nil {
		return nil, err
	}
	return &otp, nil
}

// CreateOTP tạo mới OTP
func (r *OTPPostgres) CreateOTP(ctx context.Context, otp *auth.UserOTPEntity) error {
	return r.db.Create(otp).Error
}

// UpdateOTP cập nhật OTP, nếu chưa có thì tạo mới
func (r *OTPPostgres) UpdateOTP(ctx context.Context, otp *auth.UserOTPEntity) error {
	// Kiểm tra xem OTP đã tồn tại chưa
	var existingOTP auth.UserOTPEntity
	err := r.db.Where("auth_id = ?", otp.AuthID).First(&existingOTP).Error

	if err != nil {
		// Nếu không tìm thấy, tạo mới
		if err == gorm.ErrRecordNotFound {
			return r.db.Create(otp).Error
		}
		return err
	}

	// Nếu đã tồn tại, cập nhật
	return r.db.Model(&auth.UserOTPEntity{}).
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
