package _provider

import (
	"context"
)

type ChatProvider interface {
	CreateConversation(ctx context.Context, name string, createdBy uint64, members []uint64) (uint64, error)
}
