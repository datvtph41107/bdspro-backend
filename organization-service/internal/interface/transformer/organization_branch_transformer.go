package transformer

import (
	organizationpb "pb/types/organization"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
)

type OrganizationBranchTransformer interface {
	CreateOrganizationBranchRequestToEntity(request *organizationpb.CreateOrganizationBranchRequest) *entity.OrganizationBranch
	UpdateOrganizationBranchRequestToEntity(request *organizationpb.UpdateOrganizationBranchRequest) *entity.OrganizationBranch
}

type organizationBranchTransformer struct{}

func NewOrganizationBranchTransformer() OrganizationBranchTransformer {
	return &organizationBranchTransformer{}
}

func (t *organizationBranchTransformer) CreateOrganizationBranchRequestToEntity(request *organizationpb.CreateOrganizationBranchRequest) *entity.OrganizationBranch {
	return &entity.OrganizationBranch{
		Name:           request.Name,
		Address:        request.Address,
		Phone:          request.Phone,
		Email:          request.Email,
		ManagerId:      request.ManagerId,
		OrganizationId: request.OrganizationId,
	}
}

func (t *organizationBranchTransformer) UpdateOrganizationBranchRequestToEntity(request *organizationpb.UpdateOrganizationBranchRequest) *entity.OrganizationBranch {
	return &entity.OrganizationBranch{
		Id:      request.Id,
		Name:    request.Name,
		Address: request.Address,
		Phone:   request.Phone,
		Email:   request.Email,
	}
}
