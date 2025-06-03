package transformer

import (
	organizationpb "pb/types/organization"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
)

type OrganizationTransformer interface {
	EntityToCreateOrganizationResponse(organization *entity.Organization) *organizationpb.CreateOrganizationResponse
	CreateOrganizationRequestToEntity(request *organizationpb.CreateOrganizationRequest) *entity.Organization
	UpdateOrganizationRequestToEntity(request *organizationpb.UpdateOrganizationRequest) *entity.Organization
	EntityToUpdateOrganizationResponse(organization *entity.Organization) *organizationpb.UpdateOrganizationResponse
	EntityToGetOrganizationsResponse(organizations []*entity.Organization) *organizationpb.GetOrganizationsResponse
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
	}
}

func (o *organizationTransformer) UpdateOrganizationRequestToEntity(request *organizationpb.UpdateOrganizationRequest) *entity.Organization {
	return &entity.Organization{
		ID: request.Id,
		Name: request.Name,
		TaxCode: request.TaxCode,
		BusinessLicenseUrl: request.BusinessLicenseUrl,
		Address: request.Address,
		Phone: request.Phone,	
		Email: request.Email,
		Website: request.Website,
		LogoUrl: request.LogoUrl,
	}
}

func (o *organizationTransformer) EntityToUpdateOrganizationResponse(organization *entity.Organization) *organizationpb.UpdateOrganizationResponse {
	return &organizationpb.UpdateOrganizationResponse{
		Id: organization.ID,
	}
}

func (o *organizationTransformer) EntityToGetOrganizationsResponse(organizations []*entity.Organization) *organizationpb.GetOrganizationsResponse {
	organizationsResponse := make([]*organizationpb.Organization, len(organizations))
	for i, organization := range organizations {
		organizationsResponse[i] = &organizationpb.Organization{
			Id: organization.ID,
			Name: organization.Name,
			TaxCode: organization.TaxCode,
			BusinessLicenseUrl: organization.BusinessLicenseUrl,
			Address: organization.Address,
			Phone: organization.Phone,
			Email: organization.Email,
			Website: organization.Website,
			LogoUrl: organization.LogoUrl,			
		}
	}
	return &organizationpb.GetOrganizationsResponse{
		Organizations: organizationsResponse,
	}
}