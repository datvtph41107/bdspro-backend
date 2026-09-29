package client

import (
	_utils "common/utils"
	"context"
	"crm/internal/domain"
	"crm/internal/enums"
	"errors"
	crmpb "pb/types/crm"
	sharepb "pb/types/shared"
	userpb "pb/types/user"
	"strconv"

	"pb/clients"
)

// @bind: crm/internal/interface/provider.UserClient
type UserClient struct {
	*clients.UserGrpcClient
}

func NewUserClient(rpcClient *clients.UserGrpcClient) *UserClient {
	return &UserClient{
		UserGrpcClient: rpcClient,
	}
}

func (c *UserClient) CheckOrganizationMembers(ctx context.Context, organizationID uint64, profileIDs []uint64) ([]*userpb.OrganizationMemberCheck, error) {
	if c == nil || c.OrganizationMembershipClient == nil {
		return nil, errors.New("user organization membership client not available")
	}
	response, err := c.OrganizationMembershipClient.CheckMembers(ctx, &userpb.CheckOrganizationMembersRequest{
		OrganizationId: organizationID,
		ProfileIds:     profileIDs,
	})
	if err != nil {
		return nil, err
	}
	return response.Data, nil
}

// MapOrganizationMembersToContacts enriches CRM contacts from the
// User-owned organization membership projection. CRM no longer needs to ask
// the legacy Organization service for organization identity or membership.
func (c *UserClient) MapOrganizationMembersToContacts(ctx context.Context, organizationID uint64, contacts []*sharepb.ContactDTO) error {
	profileIDs := make([]uint64, 0, len(contacts))
	seen := make(map[uint64]struct{}, len(contacts))
	for _, contact := range contacts {
		if contact == nil || contact.ProfileId == nil || *contact.ProfileId == 0 {
			continue
		}
		if _, exists := seen[*contact.ProfileId]; exists {
			continue
		}
		seen[*contact.ProfileId] = struct{}{}
		profileIDs = append(profileIDs, *contact.ProfileId)
	}
	checks, err := c.CheckOrganizationMembers(ctx, organizationID, profileIDs)
	if err != nil {
		return err
	}
	byProfile := make(map[uint64]*userpb.OrganizationMemberCheck, len(checks))
	for _, check := range checks {
		if check != nil && check.IsMember {
			byProfile[check.ProfileId] = check
		}
	}
	for _, contact := range contacts {
		if contact == nil || contact.ProfileId == nil {
			continue
		}
		check := byProfile[*contact.ProfileId]
		if check == nil {
			continue
		}
		roleKey, roleName := legacyOrganizationRole(check.RoleKey)
		contact.OrganizationMember = &sharepb.OrganizationMember{
			Id: uint32(check.MemberId), UserId: check.ProfileId,
			Role: roleKey, RoleKey: roleKey, RoleName: roleName,
			Status: 10, Profile: contact.ProfileInfo,
		}
	}
	return nil
}

func legacyOrganizationRole(role string) (uint32, string) {
	switch role {
	case "owner":
		return 510, "Chủ sở hữu"
	case "admin":
		return 520, "Quản trị viên"
	case "billing":
		return 520, "Quản lý thanh toán"
	default:
		return 530, "Thành viên"
	}
}

func (c *UserClient) GetProfileByPhones(ctx context.Context, phones []string) (map[string]*domain.Profile, error) {
	request := &userpb.GetProfileByPhonesRequest{
		Phones: phones,
	}
	response, err := c.Client.GetProfileByPhones(ctx, request)
	if err != nil {
		return nil, err
	}

	profileMap := make(map[string]*domain.Profile)
	for _, profile := range response.Profiles {
		profileMap[profile.Phone] = &domain.Profile{
			ProfileId:    profile.Id,
			FullName:     profile.FullName,
			Avatar:       profile.Avatar,
			TickVerified: profile.TickVerified,
		}
	}
	return profileMap, nil
}

func (c *UserClient) GetProfileById(ctx context.Context, profileId uint64) (*domain.Profile, error) {
	result, err := c.GetProfileWithPhoneById(ctx, profileId)
	if err != nil {
		return nil, err
	}
	result.Phone = ""

	return result, nil
}

func (c *UserClient) GetProfileWithPhoneById(ctx context.Context, profileId uint64) (*domain.Profile, error) {
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
		Phone:        response.Profiles[0].Phone,
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

func (c *UserClient) PbUserToContacts(ctx context.Context, contacts []*sharepb.ContactDTO) {
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

func (c *UserClient) PbChargePersonAndProfileToLeads(ctx context.Context, lead []*crmpb.LeadDTO) {
	profileIDSet := make(map[uint64]struct{})
	for _, lead := range lead {
		if lead.ChargeId != nil {
			profileIDSet[*lead.ChargeId] = struct{}{}
		}
		if lead.Contact != nil && lead.Contact.ProfileId != nil {
			profileIDSet[*lead.Contact.ProfileId] = struct{}{}
		}
	}

	profileMap, err := c.GetMapByIDs(ctx, profileIDSet)
	if err != nil {
		return
	}

	for _, lead := range lead {
		if lead.ChargeId != nil {
			if profile, ok := profileMap[*lead.ChargeId]; ok {
				profile.Phone = ""
				lead.ChargePerson = profile
			}
		}
		if lead.Contact != nil && lead.Contact.ProfileId != nil {
			if profile, ok := profileMap[*lead.Contact.ProfileId]; ok {
				lead.ProfileInfo = profile
			}
		}
	}
}

func (c *UserClient) PbLeadToManager(ctx context.Context, lead []*crmpb.LeadManagerItem) {
	profileIDSet := make(map[uint64]struct{})
	for _, lead := range lead {
		if lead.ChargePersonId != nil {
			profileIDSet[*lead.ChargePersonId] = struct{}{}
		}
		if lead.ProfileId != nil {
			profileIDSet[*lead.ProfileId] = struct{}{}
		}
	}

	profileMap, err := c.GetMapByIDs(ctx, profileIDSet)
	if err != nil {
		return
	}

	for _, lead := range lead {
		if lead.ChargePersonId != nil {
			if profile, ok := profileMap[*lead.ChargePersonId]; ok {
				profile.Phone = ""
				lead.ChargePersonName = profile.FullName
				lead.ChargePersonPhone = profile.Phone
				lead.ChargePersonAvatar = profile.Avatar
			}
		}
		if lead.ProfileId != nil {
			if profile, ok := profileMap[*lead.ProfileId]; ok {
				lead.ProfileName = profile.FullName
				lead.ProfilePhone = profile.Phone
				lead.ProfileAvatar = profile.Avatar
			}
		}
	}
}

func (c *UserClient) PbChargePersonToLeadDetail(ctx context.Context, lead *crmpb.LeadDetailDTO) {
	profileIDSet := make(map[uint64]struct{})
	if lead.ChargePersonId != nil {
		profileIDSet[*lead.ChargePersonId] = struct{}{}
	} else {
		return
	}

	profileMap, err := c.GetMapByIDs(ctx, profileIDSet)
	if err != nil {
		return
	}

	if profile, ok := profileMap[*lead.ChargePersonId]; ok {
		profile.Phone = ""
		lead.ChargePerson = profile
	}
}

func (c *UserClient) PbUserToReceiverFriends(ctx context.Context, friends []*crmpb.Friend) {
	profileId := _utils.GetProfileIdWithContext(ctx)
	profileIDSet := make(map[uint64]struct{})
	for _, friend := range friends {
		profileIDSet[friend.ReceiverId] = struct{}{}
		profileIDSet[friend.CreatedBy] = struct{}{}
	}

	profileMap, err := c.GetMapByIDs(ctx, profileIDSet)
	if err != nil {
		return
	}

	for _, friend := range friends {
		if friend.CreatedBy == profileId {
			if profile, ok := profileMap[friend.ReceiverId]; ok {
				friend.ReceiverUser = profile
			}
		} else {
			if profile, ok := profileMap[friend.CreatedBy]; ok {
				friend.ReceiverUser = profile
			}
		}
	}
}

func (c *UserClient) PbUserToFriends(ctx context.Context, friends []*crmpb.Friend, isCreatedBy bool) {
	profileIDSet := make(map[uint64]struct{})
	profileId := _utils.GetProfileIdWithContext(ctx)
	for _, friend := range friends {
		if friend.CreatedBy == profileId {
			profileIDSet[friend.ReceiverId] = struct{}{}
		} else {
			profileIDSet[friend.CreatedBy] = struct{}{}
		}
	}

	profileMap, err := c.GetMapByIDs(ctx, profileIDSet)
	if err != nil {
		return
	}

	for _, friend := range friends {
		if friend.CreatedBy == profileId {
			if profile, ok := profileMap[friend.ReceiverId]; ok {
				friend.ReceiverUser = profile
			}
		} else {
			if profile, ok := profileMap[friend.CreatedBy]; ok {
				friend.CreatedUser = profile
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

func (c *UserClient) ValidateProfileIds(ctx context.Context, profileIds []uint64) (bool, error) {
	request := &sharepb.GetProfileByIdsRequest{
		Ids: profileIds,
	}
	response, err := c.Client.GetProfileByIds(ctx, request)
	if err != nil {
		return false, err
	}

	return len(response.Profiles) == len(profileIds), nil
}

func (c *UserClient) GetMapByIDs(ctx context.Context, profileIDSet map[uint64]struct{}) (map[uint64]*sharepb.ProfileItem, error) {
	return c.UserGrpcClient.GetMapByIDs(ctx, profileIDSet)
}

func (c *UserClient) MapToRulePb(ctx context.Context, rules []*crmpb.RuleDTO) {
	profileIDSet := make(map[uint64]struct{})
	for _, rule := range rules {
		if rule.Trigger == int32(enums.RuleThenAssignLead) {
			key, err := strconv.ParseUint(rule.TriggerValue, 10, 64)
			if err != nil {
				continue
			}
			profileIDSet[key] = struct{}{}
		}
	}

	profileMap, err := c.GetMapByIDs(ctx, profileIDSet)
	if err != nil {
		return
	}

	for _, rule := range rules {
		if rule.Trigger == int32(enums.RuleThenAssignLead) {
			key, err := strconv.ParseUint(rule.TriggerValue, 10, 64)
			if err != nil {
				continue
			}
			if profile, ok := profileMap[key]; ok {
				rule.TriggerProfile = profile
			}
		}
	}
}

func (c *UserClient) GetProfileByIds(ctx context.Context, ids []uint64) ([]*sharepb.ProfileItem, error) {
	request := &sharepb.GetProfileByIdsRequest{
		Ids: ids,
	}
	response, err := c.Client.GetProfileByIds(ctx, request)
	if err != nil {
		return nil, err
	}

	return response.Profiles, nil
}
