package mapper

import (
	"bdspro/internal/domain"
	"bdspro/internal/enums"
	bdspropb "pb/types/bdspro"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type DealMilestoneTransformer interface {
	CreateRequestToEntity(request *bdspropb.CreateDealMilestoneRequest) *domain.DealMilestone
	UpdateRequestToEntity(request *bdspropb.UpdateDealMilestoneRequest) *domain.DealMilestone
	EntityToProto(milestone *domain.DealMilestone) *bdspropb.DealMilestone
}

type dealMilestoneTransformer struct{}

func NewDealMilestoneTransformer() DealMilestoneTransformer {
	return &dealMilestoneTransformer{}
}

func (t *dealMilestoneTransformer) CreateRequestToEntity(request *bdspropb.CreateDealMilestoneRequest) *domain.DealMilestone {
	expectedDate := request.ExpectedDate.AsTime()
	return &domain.DealMilestone{
		DealID:       request.DealId,
		Title:        request.Title,
		Description:  request.Description,
		ExpectedDate: &expectedDate,
		Status:       enums.MilestoneStatusPending,
	}
}

func (t *dealMilestoneTransformer) UpdateRequestToEntity(request *bdspropb.UpdateDealMilestoneRequest) *domain.DealMilestone {
	expectedDate := request.ExpectedDate.AsTime()
	milestone := &domain.DealMilestone{
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

func (t *dealMilestoneTransformer) EntityToProto(milestone *domain.DealMilestone) *bdspropb.DealMilestone {
	result := &bdspropb.DealMilestone{
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
