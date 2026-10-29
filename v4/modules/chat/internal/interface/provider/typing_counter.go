package provider

import "context"

// ITypingCounter điều chỉnh bộ đếm Redis (INCR/DECR) theo phòng và user đang gõ.
type ITypingCounter interface {
	AdjustTypingRef(ctx context.Context, conversationID, userID uint64, increment bool) (int64, error)
}
