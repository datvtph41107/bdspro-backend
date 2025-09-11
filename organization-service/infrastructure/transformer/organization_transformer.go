package transformer

import (
	_models "common/models"
	_utils "common/utils"
	organizationpb "pb/types/organization"

	"organization/internal/domain/entity"
	"organization/internal/dto"
)

type OrganizationTransformer interface {
	EntityToCreateOrganizationResponse(organization *entity.Organization) *organizationpb.CreateOrganizationResponse
	CreateOrganizationRequestToEntity(request *organizationpb.CreateOrganizationRequest) *entity.Organization
	UpdateOrganizationRequestToEntity(request *organizationpb.UpdateOrganizationRequest) *entity.Organization
	EntityToUpdateOrganizationResponse(organization *entity.Organization) *organizationpb.UpdateOrganizationResponse
	EntityToGetOrganizationsResponse(organizations []*entity.Organization, total uint32) *organizationpb.GetOrganizationsResponse
	EntitiesToGetOrganizationByUserIdResponse(organizations []*entity.Organization) *organizationpb.GetOrganizationByUserIdResponse
	EntityToGetOwnedOrganizationResponse(organization *entity.Organization) *organizationpb.GetOwnedOrganizationResponse
	EntitiesToGetAllOrganizationsWithMembersForAdminResponse(organizationsWithMembers []*dto.OrganizationWithMembers, total uint32) *organizationpb.GetAllOrganizationsWithMembersForAdminResponse
	EntityToPb(organization *entity.Organization) *organizationpb.Organization
}

type organizationTransformer struct{}

func NewOrganizationTransformer() OrganizationTransformer {
	return &organizationTransformer{}
}

func (o *organizationTransformer) EntityToCreateOrganizationResponse(organization *entity.Organization) *organizationpb.CreateOrganizationResponse {
	return &organizationpb.CreateOrganizationResponse{
		Id: uint32(organization.ID),
	}
}

func (o *organizationTransformer) EntityToPb(organization *entity.Organization) *organizationpb.Organization {
	return &organizationpb.Organization{
		Id:                 organization.ID,
		Name:               organization.Name,
		TaxCode:            organization.TaxCode,
		BusinessLicenseUrl: organization.BusinessLicenseUrl,
		Address:            organization.Address,
		Phone:              organization.Phone,
		Email:              organization.Email,
		Website:            organization.Website,
		LogoUrl:            organization.LogoUrl,
		BusinessDomains:    o.EntityToBusinessDomains(organization.BusinessDomains),
	}
}

func (o *organizationTransformer) CreateOrganizationRequestToEntity(request *organizationpb.CreateOrganizationRequest) *entity.Organization {
	return &entity.Organization{
		Name:               request.Name,
		TaxCode:            request.TaxCode,
		BusinessLicenseUrl: request.BusinessLicenseUrl,
		Address:            request.Address,
		DetailAddress:      request.DetailAddress,
		Phone:              request.Phone,
		Email:              request.Email,
		Website:            request.Website,
		LogoUrl:            request.LogoUrl,
		BusinessDomainIds:  request.BusinessDomainIds,
		Description:        request.Description,
		FoundedAt:          _utils.ParseStringToTime(request.FoundedAt),
		Status:             request.Status,
		ApproveInvestor:    *request.ApproveInvestor,
	}
}

func (o *organizationTransformer) UpdateOrganizationRequestToEntity(request *organizationpb.UpdateOrganizationRequest) *entity.Organization {
	result := &entity.Organization{
		BaseEntity: _models.BaseEntity{
			ID: uint64(request.Id),
		},
		Name:               request.Name,
		TaxCode:            request.TaxCode,
		BusinessLicenseUrl: request.BusinessLicenseUrl,
		Address:            request.Address,
		DetailAddress:      request.DetailAddress,
		Phone:              request.Phone,
		Email:              request.Email,
		Website:            request.Website,
		LogoUrl:            request.LogoUrl,
		BusinessDomainIds:  request.BusinessDomainIds,
		Description:        request.Description,
		FoundedAt:          _utils.ParseStringToTime(request.FoundedAt),
	}

	if request.ApproveInvestor != nil {
		result.ApproveInvestor = *request.ApproveInvestor
	}

	if request.Status != nil {
		result.Status = *request.Status
	}

	return result
}

func (o *organizationTransformer) EntityToUpdateOrganizationResponse(organization *entity.Organization) *organizationpb.UpdateOrganizationResponse {
	return &organizationpb.UpdateOrganizationResponse{
		Id: uint32(organization.ID),
	}
}

func (o *organizationTransformer) EntityToGetOrganizationsResponse(organizations []*entity.Organization, total uint32) *organizationpb.GetOrganizationsResponse {
	organizationsResponse := make([]*organizationpb.Organization, len(organizations))
	for i, org := range organizations {
		businessDomains := make([]*organizationpb.BusinessDomain, len(org.BusinessDomains))
		for j, domain := range org.BusinessDomains {
			businessDomains[j] = &organizationpb.BusinessDomain{
				Id:          domain.ID,
				Name:        domain.Name,
				Description: domain.Description,
				Code:        domain.Code,
				IsActive:    domain.IsActive,
			}
		}

		organizationsResponse[i] = &organizationpb.Organization{
			Id:                 org.ID,
			Name:               org.Name,
			TaxCode:            org.TaxCode,
			BusinessLicenseUrl: org.BusinessLicenseUrl,
			Address:            org.Address,
			Phone:              org.Phone,
			Email:              org.Email,
			Website:            org.Website,
			LogoUrl:            org.LogoUrl,
			BusinessDomains:    businessDomains,
		}
	}
	return &organizationpb.GetOrganizationsResponse{
		Data:  organizationsResponse,
		Total: total,
	}
}

func (o *organizationTransformer) EntitiesToGetOrganizationByUserIdResponse(organizations []*entity.Organization) *organizationpb.GetOrganizationByUserIdResponse {
	organizationsResponse := make([]*organizationpb.Organization, len(organizations))
	for i, org := range organizations {
		businessDomains := make([]*organizationpb.BusinessDomain, len(org.BusinessDomains))
		for j, domain := range org.BusinessDomains {
			businessDomains[j] = &organizationpb.BusinessDomain{
				Id:          domain.ID,
				Name:        domain.Name,
				Description: domain.Description,
				Code:        domain.Code,
				IsActive:    domain.IsActive,
			}
		}

		organizationsResponse[i] = &organizationpb.Organization{
			Id:                 org.ID,
			Name:               org.Name,
			TaxCode:            org.TaxCode,
			BusinessLicenseUrl: org.BusinessLicenseUrl,
			Address:            org.Address,
			Phone:              org.Phone,
			Email:              org.Email,
			LogoUrl:            org.LogoUrl,
			Website:            org.Website,
			BusinessDomains:    businessDomains,
		}
	}
	return &organizationpb.GetOrganizationByUserIdResponse{
		Data: organizationsResponse,
	}
}

func (o *organizationTransformer) EntityToGetOwnedOrganizationResponse(organization *entity.Organization) *organizationpb.GetOwnedOrganizationResponse {
	businessDomains := make([]*organizationpb.BusinessDomain, len(organization.BusinessDomains))
	for j, domain := range organization.BusinessDomains {
		businessDomains[j] = &organizationpb.BusinessDomain{
			Id:          domain.ID,
			Name:        domain.Name,
			Description: domain.Description,
			Code:        domain.Code,
			IsActive:    domain.IsActive,
		}
	}

	foundedAt := _utils.FormatTimeToString(organization.FoundedAt)
	return &organizationpb.GetOwnedOrganizationResponse{
		Organization: &organizationpb.Organization{
			Id:                 organization.ID,
			Name:               organization.Name,
			TaxCode:            organization.TaxCode,
			BusinessLicenseUrl: organization.BusinessLicenseUrl,
			Address:            organization.Address,
			Phone:              organization.Phone,
			Email:              organization.Email,
			Website:            organization.Website,
			LogoUrl:            organization.LogoUrl,
			BusinessDomains:    businessDomains,
			FoundedAt:          &foundedAt,
			Description:        organization.Description,
			ApproveInvestor:    &organization.ApproveInvestor,
			DetailAddress:      organization.DetailAddress,
		},
	}
}

func (o *organizationTransformer) EntitiesToGetAllOrganizationsWithMembersForAdminResponse(organizationsWithMembers []*dto.OrganizationWithMembers, total uint32) *organizationpb.GetAllOrganizationsWithMembersForAdminResponse {
	organizationsResponse := make([]*organizationpb.OrganizationWithMembers, len(organizationsWithMembers))
	for i, org := range organizationsWithMembers {
		businessDomains := make([]*organizationpb.BusinessDomain, len(org.Organization.BusinessDomains))
		for j, domain := range org.Organization.BusinessDomains {
			businessDomains[j] = &organizationpb.BusinessDomain{
				Id:          domain.ID,
				Name:        domain.Name,
				Description: domain.Description,
				Code:        domain.Code,
				IsActive:    domain.IsActive,
			}
		}

		// Convert members
		members := make([]*organizationpb.OrganizationMemberWithProfile, len(org.Members))
		for j, member := range org.Members {
			members[j] = &organizationpb.OrganizationMemberWithProfile{
				Id:             member.OrganizationMember.ID,
				OrganizationId: member.OrganizationMember.OrganizationID,
				UserId:         uint32(member.OrganizationMember.UserID),
				RoleId:         uint32(member.OrganizationMember.RoleId),
				Status:         uint32(member.OrganizationMember.Status),
				JoinedAt:       member.OrganizationMember.JoinedAt.Format("2006-01-02T15:04:05Z07:00"),
				RoleName:       member.OrganizationMember.RoleName,
				RoleKey:        member.OrganizationMember.GetRoleKey(),
				Permissions:    member.OrganizationMember.GetPermissions(),
				Profile:        member.Profile,
			}
			if member.OrganizationMember.RemovedAt != nil {
				removedAtStr := member.OrganizationMember.RemovedAt.Format("2006-01-02T15:04:05Z07:00")
				members[j].RemovedAt = &removedAtStr
			}
		}

		organizationsResponse[i] = &organizationpb.OrganizationWithMembers{
			Organization: &organizationpb.Organization{
				Id:                 org.Organization.ID,
				Name:               org.Organization.Name,
				TaxCode:            org.Organization.TaxCode,
				BusinessLicenseUrl: org.Organization.BusinessLicenseUrl,
				Address:            org.Organization.Address,
				Phone:              org.Organization.Phone,
				Email:              org.Organization.Email,
				Website:            org.Organization.Website,
				LogoUrl:            org.Organization.LogoUrl,
				BusinessDomains:    businessDomains,
			},
			Members:      members,
			TotalMembers: org.TotalMembers,
		}
	}
	return &organizationpb.GetAllOrganizationsWithMembersForAdminResponse{
		Data:  organizationsResponse,
		Total: total,
	}
}

func (o *organizationTransformer) EntityToBusinessDomains(businessDomains []*entity.BusinessDomain) []*organizationpb.BusinessDomain {
	businessDomainsResponse := make([]*organizationpb.BusinessDomain, len(businessDomains))
	for i, businessDomain := range businessDomains {
		businessDomainsResponse[i] = &organizationpb.BusinessDomain{
			Id:          businessDomain.ID,
			Name:        businessDomain.Name,
			Description: businessDomain.Description,
			Code:        businessDomain.Code,
			IsActive:    businessDomain.IsActive,
		}
	}
	return businessDomainsResponse
}
