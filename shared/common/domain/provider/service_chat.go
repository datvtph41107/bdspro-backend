package _provider

import (
	_dto "common/domain/dto"
	"context"
)

type ChatProvider interface {
	CreateConversation(ctx context.Context, name string, createdBy uint64, members []uint64) (uint64, error)
	MapParticipantToContact(ctx context.Context, conversationId uint64, contacts []_dto.ContactInfoV3DTO) error
	SendProductToReceiver(ctx context.Context, productID uint64, content string, receiverID *uint64) error
}
