package provider

import (
	"context"
)

type ChatProvider interface {
	SendProductToReceiver(ctx context.Context, productID uint64, content string, receiverID *uint64) error
}
