package mapper

import (
	"mdn4/modules/user/entity"
	pb "pb/types/user"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type PriceTableMapper struct{}

func NewPriceTableMapper() *PriceTableMapper {
	return &PriceTableMapper{}
}

func (m *PriceTableMapper) ToProto(entity *entity.PriceTableDomain) *pb.PriceTableResponse {
	if entity == nil {
		return nil
	}

	resp := &pb.PriceTableResponse{
		Id:          entity.ID,
		Name:        entity.Name,
		Price:       entity.Price,
		Currency:    entity.Currency,
		Description: entity.Description,
		IsActive:    entity.IsActive,
	}

	if entity.CreatedAt != nil {
		resp.CreatedAt = timestamppb.New(*entity.CreatedAt)
	}
	if entity.UpdatedAt != nil {
		resp.UpdatedAt = timestamppb.New(*entity.UpdatedAt)
	}

	return resp
}

func (m *PriceTableMapper) FromCreateRequest(req *pb.CreatePriceTableRequest) *entity.PriceTableDomain {
	if req == nil {
		return nil
	}

	return &entity.PriceTableDomain{
		Name:        req.Name,
		Price:       req.Price,
		Currency:    req.Currency,
		Description: req.Description,
		IsActive:    req.IsActive,
	}
}

func (m *PriceTableMapper) FromUpdateRequest(req *pb.UpdatePriceTableRequest) *entity.PriceTableDomain {
	if req == nil {
		return nil
	}

	entity := &entity.PriceTableDomain{
		Name:        req.Name,
		Price:       req.Price,
		Currency:    req.Currency,
		Description: req.Description,
		IsActive:    req.IsActive,
	}
	entity.ID = req.Id
	return entity
}

func (m *PriceTableMapper) ToProtoList(entities []entity.PriceTableDomain) []*pb.PriceTableResponse {
	if entities == nil {
		return []*pb.PriceTableResponse{}
	}

	result := make([]*pb.PriceTableResponse, len(entities))
	for i, entity := range entities {
		result[i] = m.ToProto(&entity)
	}
	return result
}
