package repo

import (
	"context"
	"user/internal/domain/auth"
)

// OTPRepository contains the OTP state operations consumed by serving usecases.
type OTPRepository interface {
	GetByID(ctx context.Context, authID uint64) (*auth.UserOTPEntity, error)
	CreateOTP(ctx context.Context, otp *auth.UserOTPEntity) error
	UpdateOTP(ctx context.Context, otp *auth.UserOTPEntity) error
}
