package handler

import (
	"fmt"
	"strings"
	"time"
)

// parseOptionalStrategyTime không âm thầm biến ngày sai thành nil.
// Dữ liệu transport sai phải trả 400 thay vì trở thành "chưa xác định".
func parseOptionalStrategyTime(value *string) (*time.Time, error) {
	if value == nil || strings.TrimSpace(*value) == "" {
		return nil, nil
	}

	parsed, err := time.Parse(time.RFC3339, strings.TrimSpace(*value))
	if err != nil {
		return nil, fmt.Errorf("thời gian phải theo RFC3339: %w", err)
	}

	return &parsed, nil
}
