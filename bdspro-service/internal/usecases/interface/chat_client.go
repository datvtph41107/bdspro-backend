package iusecase

import (
	"context"
)

type IChatClient interface {
	SendProductToReceiver(ctx context.Context, productID uint64, content string, receiverID *uint64) error
}
