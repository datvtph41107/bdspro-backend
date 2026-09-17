// Package fabric chuyển log cũ một chiều vào slog; không cấu hình sink riêng.
package fabric

import (
	"context"
	"encoding/json"
	"github.com/hyperledger/fabric/common/flogging"
	"log/slog"
)

// Bridge phải gọi sau khi logging khởi tạo tại process root.
func Bridge() {
	flogging.Init(flogging.Config{Format: "json", LogSpec: "debug", Writer: fabricWriter{}})
}
func parseLevel(value string) (slog.Level, error) {
	var l slog.Level
	err := l.UnmarshalText([]byte(value))
	return l, err
}

type fabricWriter struct{}

func (fabricWriter) Write(data []byte) (int, error) {
	var event struct {
		Level   string `json:"level"`
		Message string `json:"msg"`
		Name    string `json:"name"`
		Caller  string `json:"caller"`
	}
	if err := json.Unmarshal(data, &event); err != nil {
		slog.Warn("Legacy logger emitted non-JSON record", "category", "application", "event_name", "logging.legacy.invalid")
		return len(data), nil
	}
	level, err := parseLevel(event.Level)
	if err != nil {
		level = slog.LevelError
	}
	slog.Log(context.Background(), level, event.Message, "category", "application", "logger", event.Name, "legacy_source", event.Caller)
	return len(data), nil
}
