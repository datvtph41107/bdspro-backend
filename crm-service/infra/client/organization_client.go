package client

import (
	_utils "common/utils"
	"context"
	"pb/clients"
	crmpb "pb/types/crm"
	organizationpb "pb/types/organization"

	"crm/env"
)

type OrganizationClient struct {
	Client *clients.OrganizationClient
}

func NewOrganizationClient() *OrganizationClient {
	return &OrganizationClient{
		Client: clients.MakeOrganizationClient(env.ENV_RUNTIME),
	}
}

func (c *OrganizationClient) CheckUsersInOrganization(ctx context.Context, organizationID uint32, userIds []uint64) ([]*organizationpb.CheckUserInOrganization, error) {
	return c.Client.CheckUsersInOrganization(ctx, organizationID, userIds)
}

func (c *OrganizationClient) MapDealMemberToFriend(ctx context.Context, dealId uint64, userPbs []*crmpb.Friend) error {
	profileId := _utils.GetProfileIdWithContext(ctx)
	profileIDSet := make(map[uint64]struct{})
	for _, userPb := range userPbs {
		profileIDSet[userPb.ReceiverId] = struct{}{}
		profileIDSet[userPb.CreatedBy] = struct{}{}
	}

	userIds := make([]uint64, 0, len(profileIDSet))
	for userId := range profileIDSet {
		userIds = append(userIds, userId)
	}

	dealMemberMap, err := c.Client.GetDealMemberByIds(ctx, dealId, userIds)
	if err != nil {
		return err
	}

	for _, friend := range userPbs {
		if friend.CreatedBy == profileId {
			if profile, ok := dealMemberMap[friend.ReceiverId]; ok {
				friend.DealInvitation = profile
			}
		} else {
			if profile, ok := dealMemberMap[friend.CreatedBy]; ok {
				friend.DealInvitation = profile
			}
		}
	}

	return nil
}

func (c *OrganizationClient) MapGroupMemberToFriend(ctx context.Context, groupId uint64, userPbs []*crmpb.Friend) error {
	profileId := _utils.GetProfileIdWithContext(ctx)
	profileIDSet := make(map[uint64]struct{})
	for _, userPb := range userPbs {
		profileIDSet[userPb.ReceiverId] = struct{}{}
		profileIDSet[userPb.CreatedBy] = struct{}{}
	}

	userIds := make([]uint64, 0, len(profileIDSet))
	for userId := range profileIDSet {
		userIds = append(userIds, userId)
	}

	groupMemberMap, err := c.Client.GetGroupMemberByIds(ctx, groupId, userIds)
	if err != nil {
		return err
	}

	for _, friend := range userPbs {
		if friend.CreatedBy == profileId {
			if profile, ok := groupMemberMap[friend.ReceiverId]; ok {
				friend.GroupMember = profile
			}
		} else {
			if profile, ok := groupMemberMap[friend.CreatedBy]; ok {
				friend.GroupMember = profile
			}
		}
	}
	return nil
}

func (c *OrganizationClient) GetOrganizationMemberByIds(ctx context.Context, organizationId uint64, userPbs []*crmpb.Friend) error {
	profileId := _utils.GetProfileIdWithContext(ctx)
	profileIDSet := make(map[uint64]struct{})
	for _, userPb := range userPbs {
		profileIDSet[userPb.ReceiverId] = struct{}{}
		profileIDSet[userPb.CreatedBy] = struct{}{}
	}

	userIds := make([]uint64, 0, len(profileIDSet))
	for userId := range profileIDSet {
		userIds = append(userIds, userId)
	}

	organizationMemberMap, err := c.Client.GetOrganizationMemberByIds(ctx, organizationId, userIds)
	if err != nil {
		return err
	}

	for _, friend := range userPbs {
		if friend.CreatedBy == profileId {
			if profile, ok := organizationMemberMap[friend.ReceiverId]; ok {
				friend.OrganizationMember = profile
			}
		} else {
			if profile, ok := organizationMemberMap[friend.CreatedBy]; ok {
				friend.OrganizationMember = profile
			}
		}
	}
	return nil
}
