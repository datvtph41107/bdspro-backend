package _utils

import (
	"strings"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	TimeFormat   = "2006-01-02 15:04:05"
	RFC3339Milli = "2006-01-02T15:04:05.000Z07:00"
)

func TimeToTimestamp(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}

func TimestampToTime(t *timestamppb.Timestamp) *time.Time {
	if t == nil {
		return nil
	}

	tm := t.AsTime()
	return &tm
}

func FormatTimeToString(t *time.Time) string {
	if t != nil {
		// In(time.Local): TIMESTAMPTZ từ DB thường là UTC; format theo Local (+07)
		// để offset thật (Z07:00), không hardcode literal "+00:00".
		return t.In(time.Local).Format(RFC3339Milli)
	}
	return ""
}

// Parse RFC3339 / RFC3339Nano / YYYY-MM-DD string -> time.Time.
// Empty or unrecognized value returns nil (GORM Updates will skip the field).
func ParseStringToTime(value string) *time.Time {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}

	// Offset có trong chuỗi → parse tuyệt đối.
	withZone := []string{
		time.RFC3339Nano,
		time.RFC3339,
		RFC3339Milli,
	}
	for _, layout := range withZone {
		t, err := time.Parse(layout, value)
		if err == nil {
			return &t
		}
	}

	// Không có offset → hiểu theo time.Local (không mặc định UTC).
	localLayouts := []string{
		"2006-01-02T15:04:05.000",
		"2006-01-02T15:04:05",
		TimeFormat,
		"2006-01-02",
	}
	for _, layout := range localLayouts {
		t, err := time.ParseInLocation(layout, value, time.Local)
		if err == nil {
			return &t
		}
	}

	return nil
}

// --- ADDITIONAL FUNCTIONS FOR INVESTMENT MODULE ---

func ParseStringToTimeCustom(timeStr string) *time.Time {
	if timeStr == "" {
		return nil
	}

	t, err := time.Parse(TimeFormat, timeStr)
	if err != nil {
		return nil
	}

	return &t
}

func FormatTimeToStringCustom(t *time.Time) string {
	if t == nil {
		return ""
	}

	return t.Format(TimeFormat)
}

func TimeNowPtr() *time.Time {
	now := time.Now()
	return &now
}

func TimeNowRFC3339() string {
	return time.Now().UTC().Format(time.RFC3339)
}

// Current time in UTC
func TimeNowUTC() time.Time {
	return time.Now().UTC()
}

// Current time pointer in UTC
func TimeNowUTCPtr() *time.Time {
	now := time.Now().UTC()
	return &now
}

// Format time -> RFC3339 (UTC)
func FormatRFC3339(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.UTC().Format(time.RFC3339)
}
