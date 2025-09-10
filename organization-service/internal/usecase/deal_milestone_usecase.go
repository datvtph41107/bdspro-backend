package usecase

import (
	_utils "common/utils"
	"context"
	"organization/internal/custom_error"
	"organization/internal/domain/entity"
	"organization/internal/domain/repository"
	"organization/internal/enums"
)

type DealMilestoneUsecase interface {
	CreateMilestone(ctx context.Context, milestone *entity.DealMilestone) (*entity.DealMilestone, error)
	UpdateMilestone(ctx context.Context, id uint64, milestone *entity.DealMilestone) (*entity.DealMilestone, error)
	DeleteMilestone(ctx context.Context, id uint64) error
	GetMilestonesByDealID(ctx context.Context, dealID uint64) ([]*entity.DealMilestone, error)
	UpdateMilestoneOrder(ctx context.Context, dealID uint64, milestoneIDs []uint64) error
}

type dealMilestoneUsecase struct {
	milestoneRepository repository.DealMilestoneRepository
	dealRepository      repository.DealRepository
	groupAuthUsecase    GroupAuthUsecase
	logWorker           *LogWorker
}

func NewDealMilestoneUsecase(
	milestoneRepository repository.DealMilestoneRepository,
	dealRepository repository.DealRepository,
	groupAuthUsecase GroupAuthUsecase,
	logWorker *LogWorker,
) DealMilestoneUsecase {
	return &dealMilestoneUsecase{
		milestoneRepository: milestoneRepository,
		dealRepository:      dealRepository,
		groupAuthUsecase:    groupAuthUsecase,
		logWorker:           logWorker,
	}
}

func (u *dealMilestoneUsecase) CreateMilestone(ctx context.Context, milestone *entity.DealMilestone) (*entity.DealMilestone, error) {
	// Check if deal exists and user has permission
	deal, err := u.dealRepository.GetByID(ctx, milestone.DealID)
	if err != nil {
		return nil, err
	}

	currentUserID := _utils.GetProfileIdWithContext(ctx)
	err = u.groupAuthUsecase.IsLeaderGroup(ctx, uint32(deal.OwnerId), uint32(currentUserID))
	if err != nil {
		return nil, custom_error.Forbidden("you are not allowed to create milestone")
	}

	// Create milestone
	createdMilestone, err := u.milestoneRepository.Create(ctx, milestone)
	if err != nil {
		return nil, err
	}

	// Log the action
	u.logWorker.Push(LogEvent{
		GroupId: uint32(deal.OwnerId),
		ActorId: uint32(currentUserID),
		LogType: "CREATE_MILESTONE",
		LogData: "Created new milestone",
	})

	return createdMilestone, nil
}

func (u *dealMilestoneUsecase) UpdateMilestone(ctx context.Context, id uint64, milestone *entity.DealMilestone) (*entity.DealMilestone, error) {
	// Get existing milestone
	existingMilestone, err := u.milestoneRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check if deal exists and user has permission
	deal, err := u.dealRepository.GetByID(ctx, existingMilestone.DealID)
	if err != nil {
		return nil, err
	}

	currentUserID := _utils.GetProfileIdWithContext(ctx)
	err = u.groupAuthUsecase.IsLeaderGroup(ctx, uint32(deal.OwnerId), uint32(currentUserID))
	if err != nil {
		return nil, custom_error.Forbidden("you are not allowed to update milestone")
	}

	// Update milestone
	milestone.ID = id
	milestone.DealID = existingMilestone.DealID
	updatedMilestone, err := u.milestoneRepository.Update(ctx, milestone)
	if err != nil {
		return nil, err
	}

	// Log the action
	logData := "Updated milestone"
	if milestone.Status == enums.MilestoneStatusCompleted {
		logData = "Completed milestone"
	}
	u.logWorker.Push(LogEvent{
		GroupId: uint32(deal.OwnerId),
		ActorId: uint32(currentUserID),
		LogType: "UPDATE_MILESTONE",
		LogData: logData,
	})

	return updatedMilestone, nil
}

func (u *dealMilestoneUsecase) DeleteMilestone(ctx context.Context, id uint64) error {
	// Get existing milestone
	milestone, err := u.milestoneRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Check if deal exists and user has permission
	deal, err := u.dealRepository.GetByID(ctx, milestone.DealID)
	if err != nil {
		return err
	}

	currentUserID := _utils.GetProfileIdWithContext(ctx)
	err = u.groupAuthUsecase.IsLeaderGroup(ctx, uint32(deal.OwnerId), uint32(currentUserID))
	if err != nil {
		return custom_error.Forbidden("you are not allowed to delete milestone")
	}

	// Delete milestone
	err = u.milestoneRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Log the action
	u.logWorker.Push(LogEvent{
		GroupId: uint32(deal.OwnerId),
		ActorId: uint32(currentUserID),
		LogType: "DELETE_MILESTONE",
		LogData: "Deleted milestone",
	})

	return nil
}

func (u *dealMilestoneUsecase) GetMilestonesByDealID(ctx context.Context, dealID uint64) ([]*entity.DealMilestone, error) {
	// Check if deal exists and user has permission
	deal, err := u.dealRepository.GetByID(ctx, dealID)
	if err != nil {
		return nil, err
	}

	currentUserID := _utils.GetProfileIdWithContext(ctx)
	err = u.groupAuthUsecase.IsMemberGroup(ctx, uint32(deal.OwnerId), uint32(currentUserID))
	if err != nil {
		return nil, custom_error.Forbidden("you are not allowed to view milestones")
	}

	// Get milestones
	return u.milestoneRepository.GetByDealID(ctx, dealID)
}

func (u *dealMilestoneUsecase) UpdateMilestoneOrder(ctx context.Context, dealID uint64, milestoneIDs []uint64) error {
	// Check if deal exists and user has permission
	deal, err := u.dealRepository.GetByID(ctx, dealID)
	if err != nil {
		return err
	}

	currentUserID := _utils.GetProfileIdWithContext(ctx)
	err = u.groupAuthUsecase.IsLeaderGroup(ctx, uint32(deal.OwnerId), uint32(currentUserID))
	if err != nil {
		return custom_error.Forbidden("you are not allowed to reorder milestones")
	}

	// Update order
	return u.milestoneRepository.UpdateOrder(ctx, milestoneIDs)
}
