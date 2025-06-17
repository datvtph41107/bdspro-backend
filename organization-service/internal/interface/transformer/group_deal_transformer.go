package transformer

import (
	organizationpb "pb/types/organization"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
)

type GroupDealTransformer interface {
	CreateGroupDealRequestToEntity(request *organizationpb.CreateGroupDealRequest) *entity.Deal
	UpdateGroupDealRequestToEntity(request *organizationpb.UpdateGroupDealRequest) *entity.Deal
	EntityToCreateGroupDealResponse(deal *entity.Deal) *organizationpb.CreateGroupDealResponse
	EntityToUpdateGroupDealResponse(deal *entity.Deal) *organizationpb.UpdateGroupDealResponse
	EntityToGetGroupDealResponse(deal *entity.Deal) *organizationpb.GetGroupDealResponse
	EntityToGetGroupDealsResponse(deals []*entity.Deal, total uint32) *organizationpb.GetGroupDealsResponse
}

type groupDealTransformer struct{}

func NewGroupDealTransformer() GroupDealTransformer {
	return &groupDealTransformer{}
}

func (t *groupDealTransformer) CreateGroupDealRequestToEntity(request *organizationpb.CreateGroupDealRequest) *entity.Deal {
	return &entity.Deal{
		Name:    request.Name,
		Amount:  request.Amount,
		Status:  entity.DealStatus(request.Status),
		GroupId: request.GroupId,
	}
}

func (t *groupDealTransformer) UpdateGroupDealRequestToEntity(request *organizationpb.UpdateGroupDealRequest) *entity.Deal {
	return &entity.Deal{
		Id:     uint32(request.Id),
		Name:   request.Name,
		Amount: request.Amount,
		Status: entity.DealStatus(request.Status),
	}
}

func (t *groupDealTransformer) EntityToCreateGroupDealResponse(deal *entity.Deal) *organizationpb.CreateGroupDealResponse {
	return &organizationpb.CreateGroupDealResponse{
		Id: deal.Id,
	}
}

func (t *groupDealTransformer) EntityToUpdateGroupDealResponse(deal *entity.Deal) *organizationpb.UpdateGroupDealResponse {
	return &organizationpb.UpdateGroupDealResponse{
		Id: deal.Id,
	}
}

func (t *groupDealTransformer) EntityToGetGroupDealResponse(deal *entity.Deal) *organizationpb.GetGroupDealResponse {
	return &organizationpb.GetGroupDealResponse{
		Deal: &organizationpb.GroupDeal{
			Id:     deal.Id,
			Name:   deal.Name,
			Amount: deal.Amount,
			Status: string(deal.Status),
		},
	}
}

func (t *groupDealTransformer) EntityToGetGroupDealsResponse(deals []*entity.Deal, total uint32) *organizationpb.GetGroupDealsResponse {
	dealsResponse := make([]*organizationpb.GroupDeal, len(deals))
	for i, deal := range deals {
		dealsResponse[i] = &organizationpb.GroupDeal{
			Id:     deal.Id,
			Name:   deal.Name,
			Amount: deal.Amount,
		}
	}
	return &organizationpb.GetGroupDealsResponse{
		Deals: dealsResponse,
		Total: total,
	}
}
