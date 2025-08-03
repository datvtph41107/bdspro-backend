package transformer

import (
	_utils "common/utils"
	organizationpb "pb/types/organization"

	"organization/internal/domain/entity"
	"organization/internal/enums"
)

type OrganizationTransformer interface {
	EntityToCreateOrganizationResponse(organization *entity.Organization) *organizationpb.CreateOrganizationResponse
	CreateOrganizationRequestToEntity(request *organizationpb.CreateOrganizationRequest) *entity.Organization
	UpdateOrganizationRequestToEntity(request *organizationpb.UpdateOrganizationRequest) *entity.Organization
	EntityToUpdateOrganizationResponse(organization *entity.Organization) *organizationpb.UpdateOrganizationResponse
	EntityToGetOrganizationsResponse(organizations []*entity.Organization, total uint32) *organizationpb.GetOrganizationsResponse
	EntitiesToGetOrganizationByUserIdResponse(organizations []*entity.Organization) *organizationpb.GetOrganizationByUserIdResponse
}

type organizationTransformer struct{}

func NewOrganizationTransformer() OrganizationTransformer {
	return &organizationTransformer{}
}

func (o *organizationTransformer) EntityToCreateOrganizationResponse(organization *entity.Organization) *organizationpb.CreateOrganizationResponse {
	return &organizationpb.CreateOrganizationResponse{
		Id: organization.ID,
	}
}

func (o *organizationTransformer) CreateOrganizationRequestToEntity(request *organizationpb.CreateOrganizationRequest) *entity.Organization {
	return &entity.Organization{
		Name:               request.Name,
		TaxCode:            request.TaxCode,
		BusinessLicenseUrl: request.BusinessLicenseUrl,
		Address:            request.Address,
		Phone:              request.Phone,
		Email:              request.Email,
		Website:            request.Website,
		LogoUrl:            request.LogoUrl,
		BusinessDomainIds:  request.BusinessDomainIds,
		Description:        request.Description,
		FoundedAt:          _utils.ParseStringToTime(request.FoundedAt),
		Status:             enums.EMemberStatus(request.Status),
		ApproveInvestor:    *request.ApproveInvestor,
	}
}

func (o *organizationTransformer) UpdateOrganizationRequestToEntity(request *organizationpb.UpdateOrganizationRequest) *entity.Organization {
	result := &entity.Organization{
		ID:                 request.Id,
		Name:               request.Name,
		TaxCode:            request.TaxCode,
		BusinessLicenseUrl: request.BusinessLicenseUrl,
		Address:            request.Address,
		Phone:              request.Phone,
		Email:              request.Email,
		Website:            request.Website,
		LogoUrl:            request.LogoUrl,
		BusinessDomainIds:  request.BusinessDomainIds,
		Description:        request.Description,
		FoundedAt:          _utils.ParseStringToTime(request.FoundedAt),
		ApproveInvestor:    *request.ApproveInvestor,
	}

	if request.Status != nil {
		result.Status = enums.EMemberStatus(*request.Status)
	}

	return result
}

func (o *organizationTransformer) EntityToUpdateOrganizationResponse(organization *entity.Organization) *organizationpb.UpdateOrganizationResponse {
	return &organizationpb.UpdateOrganizationResponse{
		Id: organization.ID,
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
