package providers

import "context"

// ISmsProvider định nghĩa các phương thức để gửi tin nhắn SMS tới nhà cung cấp bên ngoài
type ISmsProvider interface {
	// SendSMS gửi tin nhắn SMS với nội dung và brand cụ thể tới số điện thoại đích
	SendSMS(
		ctx context.Context,
		phone string,
		brand string,
		message string,
	) (map[string]interface{}, error)
}
