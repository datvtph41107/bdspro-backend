package connect

import "time"

// ServiceStatus mô tả một upstream gRPC mà gateway biết.
type ServiceStatus struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Up      bool   `json:"up"`
	Error   string `json:"error,omitempty"`
}

// DiscoverResponse là payload từ GET /internal/services của gateway.
type DiscoverResponse struct {
	Services []ServiceStatus `json:"services"`
}

// DialConfig cấu hình dial / reconnect.
type DialConfig struct {
	BaseDelay         time.Duration
	MaxDelay          time.Duration
	MinConnectTimeout time.Duration
	KeepaliveTime     time.Duration
	KeepaliveTimeout  time.Duration
}

// DefaultDialConfig backoff reconnect mặc định (giống gateway ConnectParams).
func DefaultDialConfig() DialConfig {
	return DialConfig{
		BaseDelay:         1 * time.Second,
		MaxDelay:          5 * time.Second,
		MinConnectTimeout: 5 * time.Second,
		KeepaliveTime:     30 * time.Second,
		KeepaliveTimeout:  10 * time.Second,
	}
}

// BootstrapConfig dùng khi service start: hỏi gateway rồi nối tới các upstream đang up.
type BootstrapConfig struct {
	GatewayURL string
	EnvRuntime string
	// Wanted: tên service cần kết nối (hub, notification, ...). Rỗng = tất cả đang up.
	Wanted []string
	// WaitReadyTimeout: chờ từng conn READY sau khi dial. 0 = không chờ.
	WaitReadyTimeout time.Duration
	// DiscoverTimeout: timeout gọi gateway. 0 = 3s.
	DiscoverTimeout time.Duration
}
