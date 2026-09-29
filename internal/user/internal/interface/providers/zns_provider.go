package providers

import "context"

// IZnsProvider interface định nghĩa các phương thức để gửi tin nhắn qua Zalo Notification Service (ZNS)
type IZnsProvider interface {
	// SendZNSMessage gửi tin nhắn ZNS với template tùy chỉnh
	// Trả về response data và error nếu có
	SendZNSMessage(
		ctx context.Context,
		phone string,
		templateData map[string]interface{},
		templateID *int,
		trackingID *string,
	) (map[string]interface{}, error)

	// SendOTPZNS gửi OTP qua ZNS
	// Tự động cập nhật thông tin OTP vào auth method sau khi gửi thành công
	SendOTPZNS(
		ctx context.Context,
		phone string,
		otp string,
	) (map[string]interface{}, error)

	// RefreshAccessToken làm mới access token cho ZNS sử dụng refresh token
	// Trả về access token và refresh token mới
	RefreshAccessToken(ctx context.Context) (accessToken string, refreshToken string, err error)
}
