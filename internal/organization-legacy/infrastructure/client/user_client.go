package client

import (
	"context"
	"pb/clients"
	organizationpb "pb/types/organization"
	sharepb "pb/types/shared"
	transactionpb "pb/types/transaction"
	userpb "pb/types/user"
)

// @bind: organization/internal/interface.IUserClient
type UserClient struct {
	client userpb.ProfileServiceClient
}

func NewUserClient(rpcClient *clients.UserGrpcClient) *UserClient {
	return &UserClient{client: rpcClient.Client}
}

func (c *UserClient) GetProfileByIds(ctx context.Context, in *sharepb.GetProfileByIdsRequest) (*sharepb.GetProfileByIdsResponse, error) {
	return c.client.GetProfileByIds(ctx, in)
}

func (c *UserClient) GetMapProfileByIds(ctx context.Context, in *sharepb.GetProfileByIdsRequest) (map[uint64]*sharepb.ProfileItem, error) {
	profiles, err := c.client.GetProfileByIds(ctx, in)
	if err != nil {
		return nil, err
	}
	profileMap := make(map[uint64]*sharepb.ProfileItem)
	for _, profile := range profiles.Profiles {
		profileMap[profile.Id] = profile
	}
	return profileMap, nil
}

func (c *UserClient) MapMemberToDealPb(ctx context.Context, deals []*organizationpb.GroupDeal) {
	profileIds := make([]uint64, 0)
	for _, deal := range deals {
		if deal.ChargePersonId != 0 {
			profileIds = append(profileIds, uint64(deal.ChargePersonId))
		}
		if deal.Members != nil {
			for _, member := range deal.Members {
				profileIds = append(profileIds, uint64(member.Id))
			}
		}
		if deal.Customers != nil {
			for _, customer := range deal.Customers {
				profileIds = append(profileIds, uint64(customer.Id))
			}
		}
		if deal.Partners != nil {
			for _, partner := range deal.Partners {
				profileIds = append(profileIds, uint64(partner.Id))
			}
		}
	}

	profiles, err := c.client.GetProfileByIds(ctx, &sharepb.GetProfileByIdsRequest{Ids: profileIds})
	if err != nil {
		return
	}
	profileMap := make(map[uint64]*sharepb.ProfileItem)
	for _, profile := range profiles.Profiles {
		profileMap[profile.Id] = profile
	}

	for _, deal := range deals {
		if deal.Members != nil {
			members := make([]*sharepb.ProfileItem, len(deal.Members))
			for i, member := range deal.Members {
				members[i] = profileMap[uint64(member.Id)]
			}
			deal.Members = members
		}
		if deal.Customers != nil {
			customers := make([]*sharepb.ProfileItem, len(deal.Customers))
			for i, customer := range deal.Customers {
				customers[i] = profileMap[uint64(customer.Id)]
			}
			deal.Customers = customers
		}
		if deal.Partners != nil {
			partners := make([]*sharepb.ProfileItem, len(deal.Partners))
			for i, partner := range deal.Partners {
				partners[i] = profileMap[uint64(partner.Id)]
			}
			deal.Partners = partners
		}
		chargePerson := profileMap[deal.ChargePersonId]
		if chargePerson != nil {
			deal.ChargePerson = chargePerson
		}
	}

}

func (c *UserClient) MapToInvestPb(ctx context.Context, invests []*organizationpb.Investment) {
	ids := make([]uint64, len(invests))
	for i, invest := range invests {
		if invest.MemberId != 0 {
			ids[i] = invest.MemberId
		}
	}

	profiles, err := c.client.GetProfileByIds(ctx, &sharepb.GetProfileByIdsRequest{Ids: ids})
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

func (c *UserClient) MapToOrganizationBranchDetailPb(ctx context.Context, branch *organizationpb.OrganizationBranchDetail) {
	if branch.ManagerId < 1 {
		return
	}
	ids := make([]uint64, 1)
	ids[0] = uint64(branch.ManagerId)
	profiles, err := c.client.GetProfileByIds(ctx, &sharepb.GetProfileByIdsRequest{Ids: ids})
	if err != nil {
		return
	}
	if len(profiles.Profiles) > 0 {
		branch.ManagerUser = profiles.Profiles[0]
	}
}

func (c *UserClient) MapToOrganizationBranchMemberPb(ctx context.Context, members []*organizationpb.OrganizationBranchMember) {
	if len(members) < 1 {
		return
	}
	ids := make([]uint64, len(members))
	for i, member := range members {
		ids[i] = uint64(member.UserId)
	}
	profiles, err := c.client.GetProfileByIds(ctx, &sharepb.GetProfileByIdsRequest{Ids: ids})
	if err != nil {
		return
	}
	profileMap := make(map[uint64]*sharepb.ProfileItem)
	for _, p := range profiles.Profiles {
		p.Phone = ""
		profileMap[p.Id] = p
	}
	if len(profiles.Profiles) > 0 {
		for i, member := range members {
			if p, ok := profileMap[uint64(member.UserId)]; ok {
				members[i].User = p
			}
		}
	}
}

func (c *UserClient) MapToDealInvitationPb(ctx context.Context, invitations []*organizationpb.DealMember) {
	ids := make([]uint64, len(invitations))
	for i, invitation := range invitations {
		ids[i] = uint64(invitation.MemberId)
	}
	profiles, err := c.client.GetProfileByIds(ctx, &sharepb.GetProfileByIdsRequest{Ids: ids})
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

func (c *UserClient) MapProfileToDealMember(ctx context.Context, members []*organizationpb.DealMember) {
	ids := make([]uint64, len(members))
	for i, member := range members {
		ids[i] = uint64(member.MemberId)
	}
	profiles, err := c.client.GetProfileByIds(ctx, &sharepb.GetProfileByIdsRequest{Ids: ids})
	if err != nil {
		return
	}
	profileMap := make(map[uint64]*sharepb.ProfileItem)
	for _, p := range profiles.Profiles {
		p.Phone = ""
		profileMap[p.Id] = p
	}

	for _, member := range members {
		if p, ok := profileMap[uint64(member.MemberId)]; ok {
			// member.MemberName = p.FullName
			// member.MemberAvatar = p.Avatar
			// member.MemberPhone = p.Phone
			member.Profile = p
		}
	}
}

// func (c *UserClient) MapProfileToDealContractPb(ctx context.Context, contract *organizationpb.DealContractResponse) {
// 	ids := make([]uint64, 1)
// 	ids[0] = uint64(contract.CustomerId)
// 	profiles, err := c.client.GetProfileByIds(ctx, &sharepb.GetProfileByIdsRequest{Ids: ids})
// 	if err != nil {
// 		return
// 	}
// 	if len(profiles.Profiles) > 0 {
// 		profile := profiles.Profiles[0]
// 		profile.Phone = ""
// 		contract.Customer = profile
// 	}
// }

func (c *UserClient) MapProfileToDealContractPbList(ctx context.Context, contracts []*transactionpb.DealContractResponse) {
	ids := make([]uint64, len(contracts))
	for i, contract := range contracts {
		ids[i] = uint64(contract.CustomerId)
	}
	profiles, err := c.client.GetProfileByIds(ctx, &sharepb.GetProfileByIdsRequest{Ids: ids})
	if err != nil {
		return
	}
	profileMap := make(map[uint64]*sharepb.ProfileItem)
	for _, p := range profiles.Profiles {
		p.Phone = ""
		profileMap[p.Id] = p
	}
	for _, contract := range contracts {
		if p, ok := profileMap[uint64(contract.CustomerId)]; ok {
			contract.Customer = p
		}
	}
}

func (c *UserClient) MapProfileToOrganizationPb(ctx context.Context, organizations []*organizationpb.Organization) {
	ids := make([]uint64, len(organizations))
	for i, organization := range organizations {
		ids[i] = uint64(organization.OwnerId)
	}
	profiles, err := c.client.GetProfileByIds(ctx, &sharepb.GetProfileByIdsRequest{Ids: ids})
	if err != nil {
		return
	}
	profileMap := make(map[uint64]*sharepb.ProfileItem)
	for _, p := range profiles.Profiles {
		p.Phone = ""
		profileMap[p.Id] = p
	}
	for _, organization := range organizations {
		if p, ok := profileMap[uint64(organization.OwnerId)]; ok {
			organization.Owner = p
		}
	}
}

func (c *UserClient) MapProfileToOrganizationPbWithMembers(ctx context.Context, organizations []*organizationpb.OrganizationWithMembers) {
	ids := make([]uint64, len(organizations))
	for i, organization := range organizations {
		if organization.Organization != nil {
			ids[i] = uint64(organization.Organization.OwnerId)
		}
	}
	profiles, err := c.client.GetProfileByIds(ctx, &sharepb.GetProfileByIdsRequest{Ids: ids})
	if err != nil {
		return
	}
	profileMap := make(map[uint64]*sharepb.ProfileItem)
	for _, p := range profiles.Profiles {
		// p.Phone = ""
		profileMap[p.Id] = p
	}
	for _, organization := range organizations {
		if organization.Organization != nil {
			if p, ok := profileMap[uint64(organization.Organization.OwnerId)]; ok {
				organization.Organization.Owner = p
			}
		}
	}
}
