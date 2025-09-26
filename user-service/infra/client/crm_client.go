package client

import (
	"context"
	"pb/clients"
	userpb "pb/types/user"
	"user/env"
	user_model "user/internal/models/user"
)

// @bind: user/internal/interface/providers.FriendProvider
type CrmClient struct {
	*clients.CrmGrpcClient
}

func NewCrmClient() *CrmClient {
	return &CrmClient{
		clients.MakeCrmGrpcClient(env.ENV_RUNTIME),
	}
}

func (c *CrmClient) MapFriendStatusToSearchItemPb(ctx context.Context, searchItems []*userpb.UserSearchItem) {
	profileIds := make(map[uint64]struct{})
	for _, searchItem := range searchItems {
		profileIds[searchItem.ProfileId] = struct{}{}
	}

	friendStatusMap, err := c.CrmGrpcClient.GetMapFriendStatusByProfileIds(ctx, profileIds)
	if err != nil {
		return
	}

	for i, searchItem := range searchItems {
		searchItems[i].FriendStatus = friendStatusMap[searchItem.ProfileId]
	}
}

func (c *CrmClient) GetRelationShip(ctx context.Context, currentId uint64, targetId uint64) (*user_model.FriendItem, error) {
	return nil, nil
}
