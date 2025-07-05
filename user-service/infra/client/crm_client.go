package client

import (
	"context"
	"pb/clients"
	userpb "pb/types/user"
	"user/env"
)

type CrmClient struct {
	Client *clients.CrmGrpcClient
}

func NewCrmClient() *CrmClient {
	return &CrmClient{
		Client: clients.MakeCrmGrpcClient(env.ENV_RUNTIME),
	}
}

func (c *CrmClient) MapFriendStatusToSearchItemPb(ctx context.Context, searchItems []*userpb.UserSearchItem) {
	profileIds := make(map[uint64]struct{})
	for _, searchItem := range searchItems {
		profileIds[searchItem.ProfileId] = struct{}{}
	}

	friendStatusMap, err := c.Client.GetMapFriendStatusByProfileIds(ctx, profileIds)
	if err != nil {
		return
	}

	for i, searchItem := range searchItems {
		searchItems[i].FriendStatus = friendStatusMap[searchItem.ProfileId]
	}
}
