package client

import (
	_utils "common/utils"
	"context"
	"crm/env"
	"crm/internal/domain"
	crmpb "pb/types/crm"
	sharepb "pb/types/shared"

	"pb/clients"
)

// @bind: crm/internal/interface.UserClient.go
type UserClient struct {
	*clients.UserGrpcClient
}

func NewUserClient() *UserClient {
	return &UserClient{UserGrpcClient: clients.MakeUserGrpcClient(env.ENV_RUNTIME)}
}

func (c *UserClient) GetProfileById(ctx context.Context, profileId uint64) (*domain.Profile, error) {
	request := &sharepb.GetProfileByIdsRequest{
		Ids: []uint64{profileId},
	}
	response, err := c.Client.GetProfileByIds(ctx, request)
	if err != nil {
		return nil, err
	}

	if len(response.Profiles) == 0 {
		return nil, nil
	}

	return &domain.Profile{
		ProfileId:    profileId,
		FullName:     response.Profiles[0].FullName,
		Avatar:       response.Profiles[0].Avatar,
		TickVerified: response.Profiles[0].TickVerified,
	}, nil
}

func (c *UserClient) PbUserToContact(ctx context.Context, contact *crmpb.ContactDTO) {
	profileIDSet := make(map[uint64]struct{})
	if contact.ProfileId != nil {
		profileIDSet[*contact.ProfileId] = struct{}{}
	}

	profileMap, err := c.GetMapByIDs(ctx, profileIDSet)
	if err != nil {
		return
	}

	// Gán Author và FriendTags cho từng newsFeed
	if contact.ProfileId != nil {
		if profile, ok := profileMap[*contact.ProfileId]; ok {
			contact.ProfileInfo = profile
		}
	}
}

func (c *UserClient) PbUserToContacts(ctx context.Context, contacts []*crmpb.ContactDTO) {
	profileIDSet := make(map[uint64]struct{})
	for _, contact := range contacts {
		if contact.ProfileId != nil {
			profileIDSet[*contact.ProfileId] = struct{}{}
		}
	}

	profileMap, err := c.GetMapByIDs(ctx, profileIDSet)
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

func (c *UserClient) PbUserToLead(ctx context.Context, lead *crmpb.LeadDTO) {
	profileIDSet := make(map[uint64]struct{})
	if lead.ChargeId != nil {
		profileIDSet[*lead.ChargeId] = struct{}{}
	}

	profileMap, err := c.GetMapByIDs(ctx, profileIDSet)
	if err != nil {
		return
	}

	if profile, ok := profileMap[*lead.ChargeId]; ok {
		lead.ChargeUser = profile
	}
}

func (c *UserClient) PbUserToFriends(ctx context.Context, friends []*crmpb.Friend, isCreatedBy bool) {
	profileIDSet := make(map[uint64]struct{})
	for _, friend := range friends {
		if isCreatedBy {
			profileIDSet[friend.CreatedBy] = struct{}{}
		} else {
			profileIDSet[friend.ReceiverId] = struct{}{}
		}
	}

	profileMap, err := c.GetMapByIDs(ctx, profileIDSet)
	if err != nil {
		return
	}

	for _, friend := range friends {
		if isCreatedBy {
			if profile, ok := profileMap[friend.CreatedBy]; ok {
				friend.CreatedUser = profile
			}
		} else {
			if profile, ok := profileMap[friend.ReceiverId]; ok {
				friend.ReceiverUser = profile
			}
		}
	}
}

func (c *UserClient) FriendToUserPbs(ctx context.Context, friends []*crmpb.Friend) ([]*sharepb.ProfileItem, error) {
	profileId := _utils.GetProfileIdWithContext(ctx)
	userIDs := make([]uint64, 0, len(friends))
	for _, friend := range friends {
		if friend.CreatedBy == profileId {
			userIDs = append(userIDs, friend.ReceiverId)
		} else {
			userIDs = append(userIDs, friend.CreatedBy)
		}
	}

	profileMap, err := c.Client.GetProfileByIds(ctx, &sharepb.GetProfileByIdsRequest{Ids: userIDs})
	if err != nil {
		return nil, err
	}

	return profileMap.Profiles, nil
}
