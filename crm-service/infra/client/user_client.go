package client

import (
	"context"
	"crm/env"
	crmpb "pb/types/crm"

	"pb/clients"
)

type UserClient struct {
	client *clients.UserGrpcClient
}

func NewUserClient() *UserClient {
	return &UserClient{client: clients.MakeUserGrpcClient(env.ENV_RUNTIME)}
}

func (c *UserClient) PbUserToContact(ctx context.Context, contacts []*crmpb.ContactDTO) {
	profileIDSet := make(map[uint64]struct{})
	for _, contact := range contacts {
		if contact.ProfileId != nil {
			profileIDSet[*contact.ProfileId] = struct{}{}
		}
	}

	profileMap, err := c.client.GetMapByIDs(ctx, profileIDSet)
	if err != nil {
		return
	}

	// Gán Author và FriendTags cho từng newsFeed
	for _, contact := range contacts {
		// Gán Author
		if contact.ProfileId != nil {
			if profile, ok := profileMap[*contact.ProfileId]; ok {
				contact.ProfileInfo = profile
			}
		}
	}

}
