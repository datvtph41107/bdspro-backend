package _utils

import (
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
		return t.Format("2006-01-02T15:04:05.000+00:00")
	}
	return ""
}

// Parse RFC3339 string -> time.Time
func ParseStringToTime(value string) *time.Time {
	if value == "" {
		return nil
	}

	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return nil
	}

	return &t
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
