package usecase

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	_fault "common/fault"

	"user/internal/domain/auth"
	"user/internal/dto"
	"user/internal/enums"
)

func TestOtpCooldownFaultsCarryCanonicalIdentityAndLegacyCompatibilityMetadata(
	t *testing.T,
) {
	tests := []struct {
		name       string
		code       enums.AuthCodeEnum
		wantCode   string
		legacyCode string
		message    func(string) string
	}{
		{
			name:       "next send limited",
			code:       enums.LIMIT_NEXT_TIME,
			wantCode:   "user.otp.next_send_limited",
			legacyCode: "1006",
			message: func(second string) string {
				return fmt.Sprintf("Hãy thử lại sau %ss", second)
			},
		},
		{
			name:       "request limited",
			code:       enums.LIMIT_REQUEST_TIME,
			wantCode:   "user.otp.request_limited",
			legacyCode: "1017",
			message: func(second string) string {
				return fmt.Sprintf(
					"Bạn đã yêu cầu OTP quá nhiều, hãy thử lại sau %ss",
					second,
				)
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

			err := service.Validate(
				context.Background(),
				tt.code,
				param,
			)

			failure, ok := _fault.As(err)
			if !ok {
				t.Fatalf(
					"error = %T %v, want typed fault",
					err,
					err,
				)
			}

			if failure.Kind() != _fault.KindResourceExhausted {
				t.Fatalf(
					"kind = %q, want %q",
					failure.Kind(),
					_fault.KindResourceExhausted,
				)
			}

			if failure.Code() != tt.wantCode {
				t.Fatalf(
					"code = %q, want %q",
					failure.Code(),
					tt.wantCode,
				)
			}

			metadata := failure.Metadata()

			if metadata["legacy_code"] != tt.legacyCode {
				t.Fatalf(
					"legacy_code = %q, want %q",
					metadata["legacy_code"],
					tt.legacyCode,
				)
			}

			secondRaw := metadata["second"]

			second, err := strconv.Atoi(secondRaw)
			if err != nil {
				t.Fatalf(
					"second = %q: %v",
					secondRaw,
					err,
				)
			}

			if second < 0 || second > 60 {
				t.Fatalf(
					"second = %d, want 0..60",
					second,
				)
			}

			wantMessage := tt.message(secondRaw)

			if failure.PublicMessage() != wantMessage {
				t.Fatalf(
					"message = %q, want %q",
					failure.PublicMessage(),
					wantMessage,
				)
			}
		})
	}
}
