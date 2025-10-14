package mapper

import (
	pb "pb/types/user"
	"user/internal/models"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type MembershipMapper struct{}

func NewMembershipMapper() *MembershipMapper {
	return &MembershipMapper{}
}

func (m *MembershipMapper) ToProto(entity *models.Membership) *pb.MembershipResponse {
	if entity == nil {
		return nil
	}

	resp := &pb.MembershipResponse{
		Id:        entity.ID,
		ProfileId: entity.ProfileID,
		PlanId:    entity.PlanID,
		OrderId:   entity.OrderID,
	}

	if entity.CreatedAt != nil {
		resp.CreatedAt = timestamppb.New(*entity.CreatedAt)
	}
	if entity.UpdatedAt != nil {
		resp.UpdatedAt = timestamppb.New(*entity.UpdatedAt)
	}

	if entity.FinishedDate != nil {
		resp.FinishedDate = timestamppb.New(*entity.FinishedDate)
	}

	if entity.PaidAt != nil {
		resp.PaidAt = timestamppb.New(*entity.PaidAt)
	}

	return resp
}

func (m *MembershipMapper) FromProto(proto *pb.MembershipResponse) *models.Membership {
	if proto == nil {
		return nil
	}

	entity := &models.Membership{
		ProfileID: proto.ProfileId,
		PlanID:    proto.PlanId,
		OrderID:   proto.OrderId,
	}

	entity.ID = proto.Id
	if proto.CreatedAt != nil {
		createdAt := proto.CreatedAt.AsTime()
		entity.CreatedAt = &createdAt
	}
	if proto.UpdatedAt != nil {
		updatedAt := proto.UpdatedAt.AsTime()
		entity.UpdatedAt = &updatedAt
	}

	if proto.FinishedDate != nil {
		finishedDate := proto.FinishedDate.AsTime()
		entity.FinishedDate = &finishedDate
	}

	if proto.PaidAt != nil {
		paidAt := proto.PaidAt.AsTime()
		entity.PaidAt = &paidAt
	}

	return entity
}

func (m *MembershipMapper) ToProtoList(entities []models.Membership) []*pb.MembershipResponse {
	if entities == nil {
		return []*pb.MembershipResponse{}
	}

	result := make([]*pb.MembershipResponse, len(entities))
	for i, entity := range entities {
		result[i] = m.ToProto(&entity)
	}
	return result
}
