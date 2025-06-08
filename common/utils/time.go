package _utils

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
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

	tm := t.AsTime() // convert to time.Time
	return &tm
}

func FormatTimeToString(t *time.Time) string {
	if t != nil {
		return t.Format(time.RFC3339)
	}
	return ""
}

// ParseISOTime parses a string in ISO 8601 format (e.g., "1990-05-20T00:00:00Z") into time.Time
func ParseStringToTime(value string) *time.Time {
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return nil
	}
	return &t
}

func NowPtr() *time.Time {
	t := time.Now()
	return &t
}
