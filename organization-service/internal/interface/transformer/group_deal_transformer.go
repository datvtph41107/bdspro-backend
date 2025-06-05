package transformer

import (
	grouppb "pb/types/group"

	"github.com/tuandq2112/organization-service/internal/domain/entity"
)

type GroupDealTransformer interface {
	CreateGroupDealRequestToEntity(request *grouppb.CreateGroupDealRequest) *entity.Deal
	UpdateGroupDealRequestToEntity(request *grouppb.UpdateGroupDealRequest) *entity.Deal
	EntityToCreateGroupDealResponse(deal *entity.Deal) *grouppb.CreateGroupDealResponse
	EntityToUpdateGroupDealResponse(deal *entity.Deal) *grouppb.UpdateGroupDealResponse
	EntityToGetGroupDealResponse(deal *entity.Deal) *grouppb.GetGroupDealResponse
	EntityToGetGroupDealsResponse(deals []*entity.Deal) *grouppb.GetGroupDealsResponse
}

type groupDealTransformer struct{}

func NewGroupDealTransformer() GroupDealTransformer {
	return &groupDealTransformer{}
}

func (t *groupDealTransformer) CreateGroupDealRequestToEntity(request *grouppb.CreateGroupDealRequest) *entity.Deal {
	return &entity.Deal{
		Name:   request.Name,
		Amount: request.Amount,
		Status: entity.DealStatus(request.Status),
	}
}

func (t *groupDealTransformer) UpdateGroupDealRequestToEntity(request *grouppb.UpdateGroupDealRequest) *entity.Deal {
	return &entity.Deal{
		Id:     uint32(request.Id),
		Name:   request.Name,
		Amount: request.Amount,
		Status: entity.DealStatus(request.Status),
	}
}

func (t *groupDealTransformer) EntityToCreateGroupDealResponse(deal *entity.Deal) *grouppb.CreateGroupDealResponse {
	return &grouppb.CreateGroupDealResponse{
		Id: deal.Id,
	}
}

func (t *groupDealTransformer) EntityToUpdateGroupDealResponse(deal *entity.Deal) *grouppb.UpdateGroupDealResponse {
	return &grouppb.UpdateGroupDealResponse{
		Id: deal.Id,
	}
}

func (t *groupDealTransformer) EntityToGetGroupDealResponse(deal *entity.Deal) *grouppb.GetGroupDealResponse {
	return &grouppb.GetGroupDealResponse{
		Id: deal.Id,
	}
}

func (t *groupDealTransformer) EntityToGetGroupDealsResponse(deals []*entity.Deal) *grouppb.GetGroupDealsResponse {
	dealsResponse := make([]*grouppb.GroupDeal, len(deals))
	for i, deal := range deals {
		dealsResponse[i] = &grouppb.GroupDeal{
			Id:     deal.Id,
			Name:   deal.Name,
			Amount: deal.Amount,
		}
	}
	return &grouppb.GetGroupDealsResponse{
		Deals: dealsResponse,
	}
}
