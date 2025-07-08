package transformer

import (
	"organization/internal/domain/entity"
	"organization/internal/enums"
	organizationpb "pb/types/organization"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type DealMilestoneTransformer interface {
	CreateRequestToEntity(request *organizationpb.CreateDealMilestoneRequest) *entity.DealMilestone
	UpdateRequestToEntity(request *organizationpb.UpdateDealMilestoneRequest) *entity.DealMilestone
	EntityToProto(milestone *entity.DealMilestone) *organizationpb.DealMilestone
}

type dealMilestoneTransformer struct{}

func NewDealMilestoneTransformer() DealMilestoneTransformer {
	return &dealMilestoneTransformer{}
}

func (t *dealMilestoneTransformer) CreateRequestToEntity(request *organizationpb.CreateDealMilestoneRequest) *entity.DealMilestone {
	expectedDate := request.ExpectedDate.AsTime()
	return &entity.DealMilestone{
		DealID:       request.DealId,
		Title:        request.Title,
		Description:  request.Description,
		ExpectedDate: &expectedDate,
		Status:       enums.MilestoneStatusPending,
	}
}

func (t *dealMilestoneTransformer) UpdateRequestToEntity(request *organizationpb.UpdateDealMilestoneRequest) *entity.DealMilestone {
	expectedDate := request.ExpectedDate.AsTime()
	milestone := &entity.DealMilestone{
		Title:        request.Title,
		Description:  request.Description,
		ExpectedDate: &expectedDate,
		Status:       enums.MilestoneStatus(request.Status),
	}

	if request.CompletedDate != nil {
		completedDate := request.CompletedDate.AsTime()
		milestone.CompletedDate = &completedDate
	}

	return milestone
}

func (t *dealMilestoneTransformer) EntityToProto(milestone *entity.DealMilestone) *organizationpb.DealMilestone {
	result := &organizationpb.DealMilestone{
		Id:          milestone.ID,
		DealId:      milestone.DealID,
		Title:       milestone.Title,
		Description: milestone.Description,
		Status:      uint32(milestone.Status),
		OrderIndex:  milestone.OrderIndex,
	}

	if milestone.ExpectedDate != nil {
		result.ExpectedDate = timestamppb.New(*milestone.ExpectedDate)
	}

	if milestone.CompletedDate != nil {
		result.CompletedDate = timestamppb.New(*milestone.CompletedDate)
	}

	if milestone.CreatedAt != nil {
		result.CreatedAt = milestone.CreatedAt.Format(time.RFC3339)
	}

	if milestone.UpdatedAt != nil {
		result.UpdatedAt = milestone.UpdatedAt.Format(time.RFC3339)
	}

	return result
} 