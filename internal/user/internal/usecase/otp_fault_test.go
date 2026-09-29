package usecase

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	_errors "common/errors"
	"google.golang.org/grpc/codes"

	"user/internal/domain/auth"
	"user/internal/dto"
	"user/internal/enums"
)

func TestOtpCooldownErrorsCarryCanonicalIdentityAndLegacyCompatibilityMetadata(t *testing.T) {
	tests := []struct {
		name       string
		code       enums.AuthCodeEnum
		wantKey    _errors.Key
		legacyCode int32
		message    func(string) string
	}{
		{
			name:       "next send limited",
			code:       enums.LIMIT_NEXT_TIME,
			wantKey:    "USER_OTP_NEXT_SEND_LIMITED",
			legacyCode: 1006,
			message: func(second string) string {
				return fmt.Sprintf("Hãy thử lại sau %ss", second)
			},
		},
		{
			name:       "request limited",
			code:       enums.LIMIT_REQUEST_TIME,
			wantKey:    "USER_OTP_REQUEST_LIMITED",
			legacyCode: 1017,
			message: func(second string) string {
				return fmt.Sprintf("Bạn đã yêu cầu OTP quá nhiều, hãy thử lại sau %ss", second)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := time.Now()
			service := &OtpUsecase{
				properties: dto.PropertiesDTO{
					LockSendAfter:    60,
					LockRequestAfter: 60,
				},
			}
			param := &dto.AuthParam{
				OTP: &auth.UserOTPEntity{
					OTPDate: &now,
				},
			}

			err := service.Validate(context.Background(), tt.code, param)
			application, ok := _errors.As(err)
			if !ok {
				t.Fatalf("error = %T %v, want canonical application error", err, err)
			}
			if application.Key() != tt.wantKey {
				t.Fatalf("key = %q, want %q", application.Key(), tt.wantKey)
			}
			if application.RPCCode() != codes.ResourceExhausted {
				t.Fatalf("rpc = %q, want %q", application.RPCCode(), codes.ResourceExhausted)
			}
			legacyCode, ok := application.LegacyCode()
			if !ok || legacyCode != tt.legacyCode {
				t.Fatalf("legacy code = %d, %v; want %d, true", legacyCode, ok, tt.legacyCode)
			}

			metadata := application.Metadata()
			secondRaw := metadata["second"]
			second, parseErr := strconv.Atoi(secondRaw)
			if parseErr != nil {
				t.Fatalf("second = %q: %v", secondRaw, parseErr)
			}
			if second < 0 || second > 60 {
				t.Fatalf("second = %d, want 0..60", second)
			}

			wantMessage := tt.message(secondRaw)
			if application.PublicMessage() != wantMessage {
				t.Fatalf("message = %q, want %q", application.PublicMessage(), wantMessage)
			}
		})
	}
}
