package client

import (
	_utils "common/utils"
	"context"
	"pb/clients"
	crmpb "pb/types/crm"
	sharepb "pb/types/shared"
)

type OrganizationClient struct {
	Client *clients.OrganizationClient
}

func NewOrganizationClient(rpcClient *clients.OrganizationClient) *OrganizationClient {
	return &OrganizationClient{Client: rpcClient}
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
				friend.DealMember = profile
			}
		} else {
			if profile, ok := dealMemberMap[friend.CreatedBy]; ok {
				friend.DealMember = profile
			}
		}
	}

	return nil
}

func (c *OrganizationClient) MapGroupMemberToContact(ctx context.Context, groupId uint64, contacts []*sharepb.ContactDTO) error {
	profileIDSet := make(map[uint64]struct{})
	for _, contact := range contacts {
		if contact.ProfileId != nil {
			profileIDSet[*contact.ProfileId] = struct{}{}
		}
	}

	userIds := make([]uint64, 0, len(profileIDSet))
	for userId := range profileIDSet {
		userIds = append(userIds, userId)
	}

	groupMemberMap, err := c.Client.GetGroupMemberByIds(ctx, groupId, userIds)
	if err != nil {
		return err
	}

	for _, contact := range contacts {
		if contact.ProfileId != nil {
			if profile, ok := groupMemberMap[*contact.ProfileId]; ok {
				contact.GroupMember = &sharepb.GroupMember{
					Id:      uint32(profile.Id),
					GroupId: uint32(profile.GroupId),
					UserId:  uint32(profile.UserId),
					Role:    profile.Role,
					RoleId:  profile.RoleId,
					Profile: profile.Profile,
				}
			}
		}
	}
	return nil
}
