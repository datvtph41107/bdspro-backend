package mapper

import (
	pb "pb/types/user"
	"user/internal/models"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type PlanMapper struct{}

func NewPlanMapper() *PlanMapper {
	return &PlanMapper{}
}

func (m *PlanMapper) ToProto(entity *models.Plan) *pb.PlanResponse {
	if entity == nil {
		return nil
	}

	resp := &pb.PlanResponse{
		Id:                 entity.ID,
		PackageName:        entity.PackageName,
		Duration:           entity.Duration,
		Price:              entity.Price,
		CreateProduct:      int32(entity.CreateProduct),
		LimitPost:          entity.LimitPost,
		PublicProductLimit: int32(entity.PublicProductLimit),
		CreateProductAI:    int32(entity.CreateProductAI),
		FriendLimit:        int32(entity.FriendLimit),
		CustomerLimit:      int32(entity.CustomerLimit),
		JoinOrg:            int32(entity.JoinOrg),
		CreateOrganization: int32(entity.CreateOrganization),
		OrgMemberLimit:     int32(entity.OrgMemberLimit),
		OrgProductLimit:    int32(entity.OrgProductLimit),
		Color:              entity.Color,
		Features:           entity.Features,
		Visibility:         entity.Visibility,
	}

	if entity.CreatedAt != nil {
		resp.CreatedAt = timestamppb.New(*entity.CreatedAt)
	}
	if entity.UpdatedAt != nil {
		resp.UpdatedAt = timestamppb.New(*entity.UpdatedAt)
	}

	return resp
}

func (m *PlanMapper) FromCreateRequest(req *pb.CreatePlanRequest) *models.Plan {
	if req == nil {
		return nil
	}

	return &models.Plan{
		PackageName:        req.PackageName,
		Duration:           req.Duration,
		Price:              req.Price,
		CreateProduct:      int(req.CreateProduct),
		LimitPost:          req.LimitPost,
		PublicProductLimit: int(req.PublicProductLimit),
		CreateProductAI:    int(req.CreateProductAI),
		FriendLimit:        int(req.FriendLimit),
		CustomerLimit:      int(req.CustomerLimit),
		JoinOrg:            int(req.JoinOrg),
		CreateOrganization: int(req.CreateOrganization),
		OrgMemberLimit:     int(req.OrgMemberLimit),
		OrgProductLimit:    int(req.OrgProductLimit),
		Color:              req.Color,
		Features:           req.Features,
		Visibility:         req.Visibility,
	}
}

func (m *PlanMapper) FromUpdateRequest(req *pb.UpdatePlanRequest) *models.Plan {
	if req == nil {
		return nil
	}

	entity := &models.Plan{
		PackageName:        req.PackageName,
		Duration:           req.Duration,
		Price:              req.Price,
		CreateProduct:      int(req.CreateProduct),
		LimitPost:          req.LimitPost,
		PublicProductLimit: int(req.PublicProductLimit),
		CreateProductAI:    int(req.CreateProductAI),
		FriendLimit:        int(req.FriendLimit),
		CustomerLimit:      int(req.CustomerLimit),
		JoinOrg:            int(req.JoinOrg),
		CreateOrganization: int(req.CreateOrganization),
		OrgMemberLimit:     int(req.OrgMemberLimit),
		OrgProductLimit:    int(req.OrgProductLimit),
		Color:              req.Color,
		Features:           req.Features,
		Visibility:         req.Visibility,
	}
	entity.ID = req.Id
	return entity
}

func (m *PlanMapper) ToProtoList(entities []models.Plan) []*pb.PlanResponse {
	if entities == nil {
		return []*pb.PlanResponse{}
	}

	result := make([]*pb.PlanResponse, len(entities))
	for i, entity := range entities {
		result[i] = m.ToProto(&entity)
	}
	return result
}
