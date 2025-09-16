package transformer

import (
	organizationpb "pb/types/organization"

	"organization/internal/domain/entity"
)

type InternalOrganizationMapper struct{}

func NewInternalOrganizationMapper() *InternalOrganizationMapper {
	return &InternalOrganizationMapper{}
}

func (m *InternalOrganizationMapper) OrganizationToPb(organization *entity.Organization) *organizationpb.OrganizationItem {
	result := &organizationpb.OrganizationItem{
		Id:   uint64(organization.ID),
		Name: organization.Name,
	}

	if organization.LogoUrl != nil {
		result.Avatar = *organization.LogoUrl
	}

	return result
}

func (m *InternalOrganizationMapper) OrganizationsToPb(organizations []*entity.Organization) []*organizationpb.OrganizationItem {
	pbOrganizations := make([]*organizationpb.OrganizationItem, len(organizations))
	for i, organization := range organizations {
		pbOrganizations[i] = m.OrganizationToPb(organization)
	}
	return pbOrganizations
}

func (m *InternalOrganizationMapper) GroupToPb(group *entity.Group) *organizationpb.GroupItem {
	result := &organizationpb.GroupItem{
		Id:     uint64(group.Id),
		Name:   group.Name,
		Avatar: group.AvatarUrl,
	}

	return result
}

func (m *InternalOrganizationMapper) GroupsToPb(groups []*entity.Group) []*organizationpb.GroupItem {
	pbGroups := make([]*organizationpb.GroupItem, len(groups))
	for i, group := range groups {
		pbGroups[i] = m.GroupToPb(group)
	}
	return pbGroups
}

func (m *InternalOrganizationMapper) DealInvitationsToPb(dealInvitations []*entity.DealMember) []*organizationpb.DealMember {
	pbDealInvitations := make([]*organizationpb.DealMember, len(dealInvitations))
	for i, dealInvitation := range dealInvitations {
		pbDealInvitations[i] = m.DealInvitationToPb(dealInvitation)
	}
	return pbDealInvitations
}

func (m *InternalOrganizationMapper) DealInvitationToPb(dealInvitation *entity.DealMember) *organizationpb.DealMember {
	return &organizationpb.DealMember{
		MemberId: uint64(dealInvitation.MemberID),
		Status:   uint32(dealInvitation.Status),
		RoleId:   &dealInvitation.RoleID,
		RoleKey:  uint32(dealInvitation.RoleKey),
		// JoinedAt: dealInvitation.RespondedAt,
	}
}

func (m *InternalOrganizationMapper) GroupMembersToPb(groupMembers []*entity.GroupMember) []*organizationpb.GroupMember {
	pbGroupMembers := make([]*organizationpb.GroupMember, len(groupMembers))
	for i, groupMember := range groupMembers {
		pbGroupMembers[i] = m.GroupMemberToPb(groupMember)
	}
	return pbGroupMembers
}

func (m *InternalOrganizationMapper) GroupMemberToPb(groupMember *entity.GroupMember) *organizationpb.GroupMember {
	return &organizationpb.GroupMember{
		Id:      uint32(groupMember.Id),
		GroupId: uint32(groupMember.GroupId),
		UserId:  uint32(groupMember.UserId),
		RoleId:  groupMember.RoleId,
	}
}

func (m *InternalOrganizationMapper) OrganizationMembersToPb(organizationMembers []*entity.OrganizationMember) []*organizationpb.OrganizationMember {
	pbOrganizationMembers := make([]*organizationpb.OrganizationMember, len(organizationMembers))
	for i, organizationMember := range organizationMembers {
		pbOrganizationMembers[i] = m.OrganizationMemberToPb(organizationMember)
	}
	return pbOrganizationMembers
}

func (m *InternalOrganizationMapper) OrganizationMemberToPb(organizationMember *entity.OrganizationMember) *organizationpb.OrganizationMember {
	return &organizationpb.OrganizationMember{
		Id:      organizationMember.ID,
		UserId:  organizationMember.UserID,
		Status:  uint32(organizationMember.Status),
		RoleId:  organizationMember.RoleId,
		RoleKey: organizationMember.RoleKey,
		// Role:     organizationMember.Role,
		// Profile:  m.ProfileToPb(organizationMember.Profile),
		// JoinedAt: organizationMember.JoinedAt,
	}
}

func (m *InternalOrganizationMapper) GroupDealToPb(deal *entity.Deal) *organizationpb.GroupDeal {
	return &organizationpb.GroupDeal{
		Id:        uint32(deal.ID),
		Name:      deal.Name,
		Status:    uint32(deal.Status),
		OwnerId:   uint32(deal.OwnerId),
		OwnerType: uint32(deal.OwnerType),
	}
}

// func (m *InternalOrganizationMapper) DealTransactionsToPb(dealTransactions []*entity.DealTransaction) []*organizationpb.DealTransactionResponse {
// 	pbDealTransactions := make([]*organizationpb.DealTransactionResponse, len(dealTransactions))
// 	for i, dealTransaction := range dealTransactions {
// 		pbDealTransactions[i] = m.DealTransactionToPb(dealTransaction)
// 	}
// 	return pbDealTransactions
// }

// func (m *InternalOrganizationMapper) DealTransactionToPb(dealTransaction *entity.DealTransaction) *organizationpb.DealTransactionResponse {
// 	return &organizationpb.DealTransactionResponse{
// 		Id:          dealTransaction.ID,
// 		DealId:      dealTransaction.DealID,
// 		Type:        organizationpb.TransactionType(dealTransaction.Type),
// 		Amount:      dealTransaction.Amount,
// 		Description: dealTransaction.Description,
// 		Category:    dealTransaction.Category,
// 		Note:        dealTransaction.Note,
// 		CustomerId:  dealTransaction.CustomerID,
// 		Status:      organizationpb.ApprovedStatus(dealTransaction.Status),
// 		CreatedBy:   dealTransaction.CreatedBy,
// 		CreatedAt:   dealTransaction.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
// 		UpdatedAt:   dealTransaction.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
// 	}
// }

func (m *InternalOrganizationMapper) DealMembersToPb(dealMembers []*entity.DealMember) []*organizationpb.DealMember {
	pbDealMembers := make([]*organizationpb.DealMember, len(dealMembers))
	for i, dealMember := range dealMembers {
		pbDealMembers[i] = m.DealMemberToPb(dealMember)
	}
	return pbDealMembers
}

func (m *InternalOrganizationMapper) DealMemberToPb(dealMember *entity.DealMember) *organizationpb.DealMember {
	return &organizationpb.DealMember{
		DealId:   uint64(dealMember.DealID),
		MemberId: uint64(dealMember.MemberID),
		// MemberType: uint32(dealMember.MemberType),
		RoleId:  &dealMember.RoleID,
		RoleKey: uint32(dealMember.RoleKey),
	}
}

func (m *InternalOrganizationMapper) BranchMembersToPb(dealMembers []*entity.OrganizationBranchMember) []*organizationpb.OrganizationBranchMember {
	pbBranchMembers := make([]*organizationpb.OrganizationBranchMember, len(dealMembers))
	for i, dealMember := range dealMembers {
		pbBranchMembers[i] = m.BranchMemberToPb(dealMember)
	}
	return pbBranchMembers
}

func (m *InternalOrganizationMapper) BranchMemberToPb(branchMember *entity.OrganizationBranchMember) *organizationpb.OrganizationBranchMember {
	return &organizationpb.OrganizationBranchMember{
		Id:     branchMember.Id,
		UserId: branchMember.UserId,
		RoleId: branchMember.RoleId,
	}
}
