package iusecase

import "context"

type ChatClient interface {
	SendMessage(ctx context.Context, userId uint64, message string, newsFeedID uint64) error
}
