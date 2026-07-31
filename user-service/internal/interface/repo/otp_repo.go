package repo

import (
	"context"
	"time"
	"user/internal/domain"
)

// OTPRepository interface cho việc tương tác với dữ liệu OTP
type OTPRepository interface {
	// OTP methods
	GetByID(ctx context.Context, authID uint64) (*domain.UserOTPEntity, error)
	CreateOTP(ctx context.Context, otp *domain.UserOTPEntity) error
	UpdateOTP(ctx context.Context, otp *domain.UserOTPEntity) error
	DeleteOTP(ctx context.Context, authID uint64) error

	// OTP utilities
	GetOTPByPhone(ctx context.Context, phone string) (*domain.UserOTPEntity, error)
	GetExpiredOTPs(ctx context.Context, before time.Time) ([]*domain.UserOTPEntity, error)
	CleanupExpiredOTPs(ctx context.Context, before time.Time) error
	MarkOTPAsUsed(ctx context.Context, authID uint64) error
	ResetOTPCounters(ctx context.Context, authID uint64) error
}
