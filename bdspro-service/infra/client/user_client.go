package client

import (
	"bdspro/env"
	"context"
	clients "pb/clients"
	bdspropb "pb/types/bdspro"
	notificationpb "pb/types/notification"
	sharepb "pb/types/shared"
)

// @bind: bdspro/internal/provider.IUserProvider
type UserClient struct {
	*clients.UserGrpcClient
}

func NewUserClient() *UserClient {
	return &UserClient{
		UserGrpcClient: clients.MakeUserGrpcClient(env.ENV_RUNTIME),
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

func (c *UserClient) MapUserToSharingAccess(ctx context.Context, dtos []*bdspropb.SharingAccess) {
	ids := make(map[uint64]struct{})
	for _, dto := range dtos {
		ids[dto.ToId] = struct{}{}
	}
	profiles, err := c.UserGrpcClient.GetMapByIDs(ctx, ids)
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
	profile, err := c.UserGrpcClient.Client.GetProfileByIds(ctx, &sharepb.GetProfileByIdsRequest{
		Ids: []uint64{id},
	})
	if err != nil {
		return nil
	}
	if len(profile.Profiles) == 0 {
		return nil
	}
	return profile.Profiles[0]
}

func (c *UserClient) MapToInvestPb(ctx context.Context, invests []*bdspropb.Investment) {
	ids := make([]uint64, len(invests))
	for i, invest := range invests {
		if invest.MemberId != 0 {
			ids[i] = invest.MemberId
		}
	}

	profiles, err := c.UserGrpcClient.Client.GetProfileByIds(ctx, &sharepb.GetProfileByIdsRequest{Ids: ids})
	if err != nil {
		return
	}
	if len(profiles.Profiles) == 0 {
		return
	}

	profileMap := make(map[uint64]*sharepb.ProfileItem)
	for _, p := range profiles.Profiles {
		p.Phone = ""
		profileMap[p.Id] = p
	}
	for i, invest := range invests {
		if invest.MemberId != 0 {
			if p, ok := profileMap[invest.MemberId]; ok {
				invests[i].CreatedUser = p
			}
		}

	}
}

func (c *UserClient) MapToDealInvitationPb(ctx context.Context, invitations []*bdspropb.DealMember) {
	ids := make([]uint64, len(invitations))
	for i, invitation := range invitations {
		ids[i] = uint64(invitation.MemberId)
	}
	profiles, err := c.Client.GetProfileByIds(ctx, &sharepb.GetProfileByIdsRequest{Ids: ids})
	if err != nil {
		return
	}
	profileMap := make(map[uint64]*sharepb.ProfileItem)
	for _, p := range profiles.Profiles {
		p.Phone = ""
		profileMap[p.Id] = p
	}
	for _, invitation := range invitations {
		if p, ok := profileMap[uint64(invitation.MemberId)]; ok {
			invitation.Profile = p
		}
	}
}

func (c *UserClient) MapProfileToTransaction(ctx context.Context, transactions []*bdspropb.TransactionResponse) error {
	ids := make(map[uint64]struct{})
	for _, transaction := range transactions {
		if transaction.FromOf == bdspropb.OwnerType_OWNER_TYPE_USER {
			ids[transaction.FromId] = struct{}{}
		}
	}

	profiles, err := c.GetMapByIDs(ctx, ids)
	if err != nil {
		return err
	}

	for _, transaction := range transactions {
		profile, ok := profiles[transaction.FromId]
		if ok && transaction.FromOf == bdspropb.OwnerType_OWNER_TYPE_USER {
			transaction.Profile = profile
		}
	}

	return nil
}

func (c *UserClient) GetProfileByIds(ctx context.Context, in *sharepb.GetProfileByIdsRequest) (*sharepb.GetProfileByIdsResponse, error) {
	return c.UserGrpcClient.Client.GetProfileByIds(ctx, in)
}

func (c *UserClient) GetMapProfileByIds(ctx context.Context, in *sharepb.GetProfileByIdsRequest) (map[uint64]*sharepb.ProfileItem, error) {
	profiles, err := c.UserGrpcClient.Client.GetProfileByIds(ctx, in)
	if err != nil {
		return nil, err
	}
	profileMap := make(map[uint64]*sharepb.ProfileItem)
	for _, profile := range profiles.Profiles {
		profileMap[profile.Id] = profile
	}
	return profileMap, nil
}
