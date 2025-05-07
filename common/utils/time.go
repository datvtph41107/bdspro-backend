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
