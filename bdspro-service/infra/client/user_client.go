package client

import (
	"bdspro/env"
	"context"
	clients "pb/clients"
	notificationpb "pb/types/notification"
)

type UserClient struct {
	Client *clients.UserGrpcClient
}

func NewUserClient() *UserClient {
	return &UserClient{
		Client: clients.MakeUserGrpcClient(env.ENV_RELEASE),
	}
}

// todo: thêm thông tin user vào lịch sử
func (c *UserClient) MapProfileToHistory(ctx context.Context, history []*notificationpb.HistoryDTO) {
	ids := make(map[uint64]struct{})
	for _, h := range history {
		if h.OwnerId != nil {
			ids[*h.OwnerId] = struct{}{}
		}
	}
	// profiles, err := c.Client.GetMapByIDs(ctx, ids)
	// if err != nil {
	// 	return
	// }

	for _, h := range history {
		if h.OwnerId != nil {
			// if profile, ok := profiles[*h.OwnerId]; ok {
			// 	h.OwnerId = profile.Id
			// }
		}
	}
}
