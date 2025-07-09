package client

import (
	"bdspro/env"
	"context"
	clients "pb/clients"
	pb_bdspro "pb/types/bdspro"
	notificationpb "pb/types/notification"
	sharepb "pb/types/shared"
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

func (c *UserClient) MapUserToSharingAccess(ctx context.Context, dtos []*pb_bdspro.SharingAccess) {
	ids := make(map[uint64]struct{})
	for _, dto := range dtos {
		ids[dto.ToId] = struct{}{}
	}
	profiles, err := c.Client.GetMapByIDs(ctx, ids)
	if err != nil {
		return
	}

	for i, dto := range dtos {
		if profile, ok := profiles[dto.ToId]; ok {
			dtos[i].ToName = profile.FullName
			dtos[i].Avatar = profile.Avatar
		}
	}
}

func (c *UserClient) GetProfileById(ctx context.Context, id uint64) *sharepb.ProfileItem {
	if c.Client == nil {
		return nil
	}
	profile, err := c.Client.Client.GetProfileByIds(ctx, &sharepb.GetProfileByIdsRequest{
		Ids: []uint64{id},
	})
	if err != nil {
		return nil
	}
	return profile.Profiles[0]
}
