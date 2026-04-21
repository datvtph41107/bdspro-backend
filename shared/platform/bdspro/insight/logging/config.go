package logging

import "errors"

// Config cấu hình Logger — tường minh hơn variadic options.
//
// Tại sao struct config thay vì functional options?
// - Functional options đẹp nhưng khó validate
// - Struct config: zero value = sensible defaults, validation rõ ràng
// - Wire/DI inject struct dễ hơn func options (F2)
type Config struct {
	// ServiceName xuất hiện trong mọi log entry: {"service": "gateway"}
	// Required — không có service name = không biết log từ đâu
	ServiceName string

	// Environment: "production" | "development" | "test"
	// - production: JSON output, no color, Info level
	// - development: text output, color, Debug level
	// - test: no output (discard) — test không bị polluted bởi log
	Environment string

	// Level override — default là Info (production) hoặc Debug (dev)
	// Để trống để dùng default theo environment
	Level string
}

// Validate kiểm tra config hợp lệ.
func (c Config) Validate() error {
	if c.ServiceName == "" {
		return errors.New("logging: ServiceName is required")
	}
	return nil
}
