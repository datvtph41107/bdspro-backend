package client

import (
	"chat/models"
	"context"
	"pb/clients"
	chatpb "pb/types/chat"
	sharepb "pb/types/shared"
)

type CrmClient struct {
	Client *clients.CrmGrpcClient
}

func NewCrmClient(rpcClient *clients.CrmGrpcClient) *CrmClient {
	return &CrmClient{Client: rpcClient}
}

func (c *CrmClient) MapContactToMessagePbList(ctx context.Context, messages []*chatpb.Message) {
	contactIdSet := make(map[uint64]struct{})
	for _, msg := range messages {
		if msg.ContentType == string(models.ContentTypeContact) && msg.ExtraId != nil {
			contactIdSet[*msg.ExtraId] = struct{}{}
		}
	}

	contactMap := make(map[uint64]*sharepb.ContactDTO)
	for contactId := range contactIdSet {
		contact, err := c.Client.GetContactById(ctx, contactId)
		if err != nil || contact == nil {
			continue
		}
		contactMap[contactId] = contact
	}
}
