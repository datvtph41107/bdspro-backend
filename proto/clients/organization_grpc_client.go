package clients

import (
	"context"
	"errors"

	organizationpb "pb/types/organization"

	"google.golang.org/grpc"
)

type OrganizationClient struct {
	Client     organizationpb.InternalOrganizationServiceClient
	RoleClient organizationpb.OrganizationRoleServiceClient
}

func BindOrganizationClient(conn grpc.ClientConnInterface) *OrganizationClient {
	return &OrganizationClient{
		Client:     organizationpb.NewInternalOrganizationServiceClient(conn),
		RoleClient: organizationpb.NewOrganizationRoleServiceClient(conn),
	}
}

func (c *OrganizationClient) GetOrganizationMapByIDs(ctx context.Context, organizationIDSet map[uint64]struct{}) (map[uint64]*organizationpb.OrganizationItem, error) {
	if c.Client == nil {
		return nil, errors.New("organization client not avaiable")
	}

	organizationIDs := make([]uint64, 0, len(organizationIDSet))
	for id := range organizationIDSet {
		organizationIDs = append(organizationIDs, id)
	}

	organizations, err := c.Client.GetOrganizationByIds(ctx, &organizationpb.GetOrganizationByIdsRequest{Ids: organizationIDs})
	if err != nil {
		return nil, err
	}

	organizationMap := make(map[uint64]*organizationpb.OrganizationItem, len(organizations.Organizations))
	for _, organization := range organizations.Organizations {
		organizationMap[organization.Id] = organization
	}

	return organizationMap, nil
}

func (c *OrganizationClient) GetGroupMapByIDs(ctx context.Context, groupIDSet map[uint64]struct{}) (map[uint64]*organizationpb.GroupItem, error) {
	if c.Client == nil {
		return nil, errors.New("organization client not avaiable")
	}

	groupIDs := make([]uint64, 0, len(groupIDSet))
	for id := range groupIDSet {
		groupIDs = append(groupIDs, id)
	}

	groups, err := c.Client.GetGroupByIds(ctx, &organizationpb.GetGroupByIdsRequest{Ids: groupIDs})
	if err != nil {
		return nil, err
	}

	groupMap := make(map[uint64]*organizationpb.GroupItem, len(groups.Groups))
	for _, group := range groups.Groups {
		groupMap[group.Id] = group
	}

	return groupMap, nil
}

func (c *OrganizationClient) GetDealMemberByIds(ctx context.Context, dealId uint64, userIds []uint64) (map[uint64]*organizationpb.DealMember, error) {
	if c.Client == nil {
		return nil, errors.New("organization client not avaiable")
	}

	dealMembers, err := c.Client.GetDealMemberByIds(ctx, &organizationpb.GetMemberByIdsRequest{
		GroupId: dealId,
		Ids:     userIds,
	})
	if err != nil {
		return nil, err
	}

	dealMemberMap := make(map[uint64]*organizationpb.DealMember)
	for _, dealMember := range dealMembers.Data {
		dealMemberMap[dealMember.MemberId] = dealMember
	}

	return dealMemberMap, nil
}

func (c *OrganizationClient) GetGroupMemberByIds(ctx context.Context, groupId uint64, userIds []uint64) (map[uint64]*organizationpb.GroupMember, error) {
	if c.Client == nil {
		return nil, errors.New("organization client not avaiable")
	}

	groupMembers, err := c.Client.GetGroupMemberByIds(ctx, &organizationpb.GetMemberByIdsRequest{
		GroupId: groupId,
		Ids:     userIds,
	})
	if err != nil {
		return nil, err
	}

	groupMemberMap := make(map[uint64]*organizationpb.GroupMember)
	for _, groupMember := range groupMembers.Data {
		groupMemberMap[uint64(groupMember.UserId)] = groupMember
	}

	return groupMemberMap, nil
}

func (c *OrganizationClient) GetOrganizationMemberByIds(ctx context.Context, organizationId uint64, userIds []uint64) (map[uint64]*organizationpb.OrganizationMember, error) {
	if c.Client == nil {
		return nil, errors.New("organization client not avaiable")
	}

	organizationMembers, err := c.Client.GetOrganizationMemberByIds(ctx, &organizationpb.GetMemberByIdsRequest{
		GroupId: organizationId,
		Ids:     userIds,
	})
	if err != nil {
		return nil, err
	}

	organizationMemberMap := make(map[uint64]*organizationpb.OrganizationMember)
	for _, organizationMember := range organizationMembers.Data {
		organizationMemberMap[uint64(organizationMember.UserId)] = organizationMember
	}

	return organizationMemberMap, nil
}

func (c *OrganizationClient) CheckUsersInOrganization(ctx context.Context, organizationID uint32, userIds []uint64) ([]*organizationpb.CheckUserInOrganization, error) {
	if c.Client == nil {
		return nil, errors.New("organization client not avaiable")
	}

	users, err := c.Client.CheckUsersInOrganization(ctx, &organizationpb.CheckUsersInOrganizationRequest{
		OrganizationId: organizationID,
		UserIds:        userIds,
	})
	if err != nil {
		return nil, err
	}

	return users.Data, nil
}

func (c *OrganizationClient) GetDealById(ctx context.Context, dealId uint64) (*organizationpb.GroupDeal, error) {
	if c.Client == nil {
		return nil, errors.New("organization client not avaiable")
	}

	deal, err := c.Client.GetDealById(ctx, &organizationpb.GetDealByIdRequest{Id: dealId})
	if err != nil {
		return nil, err
	}

	return deal.Deal, nil
}

func (c *OrganizationClient) GetDealMember(ctx context.Context, dealId uint64, userId uint64) (*organizationpb.DealMember, error) {
	if c.Client == nil {
		return nil, errors.New("organization client not avaiable")
	}

	resp, err := c.Client.GetDealMembers(ctx, &organizationpb.DealMemberRequest{
		DealId:  dealId,
		UserIds: []uint64{userId},
	})
	if err != nil {
		return nil, err
	}
	if len(resp.Data) == 0 {
		return nil, nil
	}

	return resp.Data[0], nil
}

func (c *OrganizationClient) GetBranchMember(ctx context.Context, branchId uint64, userId uint64) (*organizationpb.OrganizationBranchMember, error) {
	if c.Client == nil {
		return nil, errors.New("organization client not avaiable")
	}

	resp, err := c.Client.GetBranchMembers(ctx, &organizationpb.GetBranchMembersRequest{
		BranchId: branchId,
		UserIds:  []uint64{userId},
	})
	if err != nil {
		return nil, err
	}

	if len(resp.Data) == 0 {
		return nil, nil
	}

	return resp.Data[0], nil
}

func (c *OrganizationClient) GetOrganizationMember(ctx context.Context, organizationId uint64, userId uint64) (*organizationpb.OrganizationMember, error) {
	if c.Client == nil {
		return nil, errors.New("organization client not avaiable")
	}

	resp, err := c.Client.GetOrganizationMemberByIds(ctx,
		&organizationpb.GetMemberByIdsRequest{GroupId: organizationId, Ids: []uint64{userId}},
	)
	if err != nil {
		return nil, err
	}

	if len(resp.Data) == 0 {
		return nil, nil
	}

	return resp.Data[0], nil
}

// AuthorizeOrganizationAction hỏi đúng policy owner thay vì để service gọi
// tự diễn giải role/permission của organization.
func (c *OrganizationClient) AuthorizeOrganizationAction(ctx context.Context, organizationID, actorProfileID uint64, action organizationpb.OrganizationAction) (bool, error) {
	if c == nil || c.Client == nil {
		return false, errors.New("organization client not available")
	}
	resp, err := c.Client.AuthorizeOrganizationAction(ctx, &organizationpb.AuthorizeOrganizationActionRequest{
		OrganizationId: organizationID,
		ActorProfileId: actorProfileID,
		Action:         action,
	})
	if err != nil {
		return false, err
	}
	return resp.GetAllowed(), nil
}

func (c *OrganizationClient) GetOrganizationsByUserId(ctx context.Context, userId uint64) ([]*organizationpb.Organization, error) {
	if c.Client == nil {
		return nil, errors.New("organization client not avaiable")
	}

	// resp, err := c.Client.GetOrganizationByUserId(ctx, &organizationpb.GetOrganizationByUserIdRequest{UserId: userId})
	// if err != nil {
	// 	return nil, err
	// }

	// if len(resp.Data) == 0 {
	// 	return nil, nil
	// }

	return nil, nil
}

// func (c *OrganizationClient) GetDealTransactionsByIds(ctx context.Context, dealIds []uint64) ([]*organizationpb.DealTransactionResponse, error) {
// 	if c.Client == nil {
// 		return nil, errors.New("organization client not avaiable")
// 	}

// 	dealTransactions, err := c.Client.GetDealTransactionsByIds(ctx, &organizationpb.GetDealTransactionsByIdsRequest{Ids: dealIds})
// 	if err != nil {
// 		return nil, err
// 	}

// 	return dealTransactions.DealTransactions, nil
// }
