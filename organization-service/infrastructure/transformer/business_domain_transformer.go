package transformer

import (
	organizationpb "pb/types/organization"

	"organization/internal/domain/entity"
)

type BusinessDomainTransformer interface {
	CreateBusinessDomainRequestToEntity(req *organizationpb.CreateBusinessDomainRequest) *entity.BusinessDomain
	UpdateBusinessDomainRequestToEntity(req *organizationpb.UpdateBusinessDomainRequest) *entity.BusinessDomain
	EntityToCreateBusinessDomainResponse(entity *entity.BusinessDomain) *organizationpb.CreateBusinessDomainResponse
	EntityToUpdateBusinessDomainResponse(entity *entity.BusinessDomain) *organizationpb.UpdateBusinessDomainResponse
	EntityToDeleteBusinessDomainResponse(entity *entity.BusinessDomain) *organizationpb.DeleteBusinessDomainResponse
	EntityToGetBusinessDomainResponse(entity *entity.BusinessDomain) *organizationpb.GetBusinessDomainResponse
	EntitiesToGetAllBusinessDomainsResponse(entities []*entity.BusinessDomain) *organizationpb.GetAllBusinessDomainsResponse
	EntitiesToGetActiveBusinessDomainsResponse(entities []*entity.BusinessDomain) *organizationpb.GetActiveBusinessDomainsResponse
	EntityToBusinessDomain(entity *entity.BusinessDomain) *organizationpb.BusinessDomain
}

type businessDomainTransformer struct{}

func NewBusinessDomainTransformer() BusinessDomainTransformer {
	return &businessDomainTransformer{}
}

func (t *businessDomainTransformer) CreateBusinessDomainRequestToEntity(req *organizationpb.CreateBusinessDomainRequest) *entity.BusinessDomain {
	return &entity.BusinessDomain{
		Name:        req.Name,
		Description: req.Description,
		Code:        req.Code,
		IsActive:    req.IsActive,
	}
}

func (t *businessDomainTransformer) UpdateBusinessDomainRequestToEntity(req *organizationpb.UpdateBusinessDomainRequest) *entity.BusinessDomain {
	return &entity.BusinessDomain{
		ID:          req.Id,
		Name:        req.Name,
		Description: req.Description,
		Code:        req.Code,
		IsActive:    req.IsActive,
	}
}

func (t *businessDomainTransformer) EntityToCreateBusinessDomainResponse(entity *entity.BusinessDomain) *organizationpb.CreateBusinessDomainResponse {
	return &organizationpb.CreateBusinessDomainResponse{
		Id: entity.ID,
	}
}

func (t *businessDomainTransformer) EntityToUpdateBusinessDomainResponse(entity *entity.BusinessDomain) *organizationpb.UpdateBusinessDomainResponse {
	return &organizationpb.UpdateBusinessDomainResponse{
		Id: entity.ID,
	}
}

func (t *businessDomainTransformer) EntityToDeleteBusinessDomainResponse(entity *entity.BusinessDomain) *organizationpb.DeleteBusinessDomainResponse {
	return &organizationpb.DeleteBusinessDomainResponse{
		Id: entity.ID,
	}
}

func (t *businessDomainTransformer) EntityToGetBusinessDomainResponse(entity *entity.BusinessDomain) *organizationpb.GetBusinessDomainResponse {
	return &organizationpb.GetBusinessDomainResponse{
		BusinessDomain: t.EntityToBusinessDomain(entity),
	}
}

func (t *businessDomainTransformer) EntitiesToGetAllBusinessDomainsResponse(entities []*entity.BusinessDomain) *organizationpb.GetAllBusinessDomainsResponse {
	businessDomains := make([]*organizationpb.BusinessDomain, len(entities))
	for i, entity := range entities {
		businessDomains[i] = t.EntityToBusinessDomain(entity)
	}
	return &organizationpb.GetAllBusinessDomainsResponse{
		Data:  businessDomains,
		Total: uint32(len(entities)),
	}
}

func (t *businessDomainTransformer) EntitiesToGetActiveBusinessDomainsResponse(entities []*entity.BusinessDomain) *organizationpb.GetActiveBusinessDomainsResponse {
	businessDomains := make([]*organizationpb.BusinessDomain, len(entities))
	for i, entity := range entities {
		businessDomains[i] = t.EntityToBusinessDomain(entity)
	}
	return &organizationpb.GetActiveBusinessDomainsResponse{
		Data: businessDomains,
	}
}

func (t *businessDomainTransformer) EntityToBusinessDomain(entity *entity.BusinessDomain) *organizationpb.BusinessDomain {
	return &organizationpb.BusinessDomain{
		Id:          entity.ID,
		Name:        entity.Name,
		Description: entity.Description,
		Code:        entity.Code,
		IsActive:    entity.IsActive,
		CreatedAt:   entity.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   entity.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		CreatedBy:   entity.CreatedBy,
		UpdatedBy:   entity.UpdatedBy,
	}
}
