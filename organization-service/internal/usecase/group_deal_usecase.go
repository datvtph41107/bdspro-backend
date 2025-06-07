package usecase

import (
	"context"

	"github.com/tuandq2112/organization-service/appconfig"
	"github.com/tuandq2112/organization-service/internal/custom_error"
	"github.com/tuandq2112/organization-service/internal/domain/entity"
	"github.com/tuandq2112/organization-service/internal/domain/repository"
	"github.com/tuandq2112/organization-service/pkg/utils"
)

type GroupDealUsecase interface {
	CreateGroupDeal(ctx context.Context, deal *entity.Deal) (*entity.Deal, error)
	UpdateGroupDeal(ctx context.Context, deal *entity.Deal) (*entity.Deal, error)
	DeleteGroupDeal(ctx context.Context, id uint32) error
	GetGroupDealByID(ctx context.Context, id uint32) (*entity.Deal, error)
	GetGroupDealsByGroupID(ctx context.Context, groupID uint32, page, size int) ([]*entity.Deal, uint32, error)
	UpdateGroupDealStatus(ctx context.Context, id uint32, status entity.DealStatus) error
}

type groupDealUsecase struct {
	groupDealRepository repository.GroupDealRepository
	groupAuthUsecase    GroupAuthUsecase
	logWorker           *LogWorker
}

func NewGroupDealUsecase(
	groupDealRepository repository.GroupDealRepository,
	groupAuthUsecase GroupAuthUsecase,
	logWorker *LogWorker,
) GroupDealUsecase {
	return &groupDealUsecase{
		groupDealRepository: groupDealRepository,
		groupAuthUsecase:    groupAuthUsecase,
		logWorker:           logWorker,
	}
}

func (u *groupDealUsecase) CreateGroupDeal(ctx context.Context, deal *entity.Deal) (*entity.Deal, error) {
	currentUserId := utils.GetUserID(ctx, appconfig.USER_CONTEXT)
	err := u.groupAuthUsecase.IsLeaderGroup(ctx, deal.GroupId, currentUserId)
	if err != nil {
		return nil, custom_error.Forbidden("you are not allowed to create deal")
	}
	deal.CreatedBy = currentUserId
	deal.UpdatedBy = currentUserId
	createdDeal, err := u.groupDealRepository.Create(ctx, deal)
	if err != nil {
		return nil, err
	}

	u.logWorker.Push(LogEvent{
		GroupId: createdDeal.GroupId,
		ActorId: createdDeal.CreatedBy,
		LogType: "CREATE_DEAL",
		LogData: "Created new deal",
	})

	return createdDeal, nil
}

func (u *groupDealUsecase) UpdateGroupDeal(ctx context.Context, deal *entity.Deal) (*entity.Deal, error) {
	existedDeal, err := u.groupDealRepository.GetByID(ctx, deal.Id)
	if err != nil {
		return nil, err
	}

	currentUserId := utils.GetUserID(ctx, appconfig.USER_CONTEXT)
	err = u.groupAuthUsecase.IsLeaderGroup(ctx, existedDeal.GroupId, currentUserId)
	if err != nil {
		return nil, custom_error.Forbidden("you are not allowed to update deal")
	}
	deal.UpdatedBy = currentUserId

	updatedDeal, err := u.groupDealRepository.Update(ctx, deal)
	if err != nil {
		return nil, err
	}

	u.logWorker.Push(LogEvent{
		GroupId: updatedDeal.GroupId,
		ActorId: updatedDeal.CreatedBy,
		LogType: "UPDATE_DEAL",
		LogData: "Updated deal information",
	})

	return updatedDeal, nil
}

func (u *groupDealUsecase) DeleteGroupDeal(ctx context.Context, id uint32) error {
	deal, err := u.groupDealRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if deal == nil {
		return custom_error.RecordNotFound("deal not found")
	}

	currentUserId := utils.GetUserID(ctx, appconfig.USER_CONTEXT)
	err = u.groupAuthUsecase.IsLeaderGroup(ctx, deal.GroupId, currentUserId)
	if err != nil {
		return custom_error.Forbidden("you are not allowed to delete deal")
	}

	err = u.groupDealRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	u.logWorker.Push(LogEvent{
		GroupId: deal.GroupId,
		ActorId: deal.CreatedBy,
		LogType: "DELETE_DEAL",
		LogData: "Deleted deal",
	})

	return nil
}

func (u *groupDealUsecase) GetGroupDealByID(ctx context.Context, id uint32) (*entity.Deal, error) {
	deal, err := u.groupDealRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	currentUserId := utils.GetUserID(ctx, appconfig.USER_CONTEXT)
	err = u.groupAuthUsecase.IsMemberGroup(ctx, deal.GroupId, currentUserId)
	if err != nil {
		return nil, custom_error.Forbidden("you are not allowed to get deal")
	}

	return deal, nil
}

func (u *groupDealUsecase) GetGroupDealsByGroupID(ctx context.Context, groupID uint32, page, size int) ([]*entity.Deal, uint32, error) {
	currentUserId := utils.GetUserID(ctx, appconfig.USER_CONTEXT)
	err := u.groupAuthUsecase.IsMemberGroup(ctx, groupID, currentUserId)
	if err != nil {
		return nil, 0, custom_error.Forbidden("you are not allowed to get deals")
	}
	return u.groupDealRepository.GetByGroupID(ctx, groupID, page, size)
}

func (u *groupDealUsecase) UpdateGroupDealStatus(ctx context.Context, id uint32, status entity.DealStatus) error {
	deal, err := u.groupDealRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if deal == nil {
		return custom_error.RecordNotFound("deal not found")
	}

	currentUserId := utils.GetUserID(ctx, appconfig.USER_CONTEXT)
	err = u.groupAuthUsecase.IsLeaderGroup(ctx, deal.GroupId, currentUserId)
	if err != nil {
		return custom_error.Forbidden("you are not allowed to update deal status")
	}

	err = u.groupDealRepository.UpdateStatus(ctx, id, status)
	if err != nil {
		return err
	}

	u.logWorker.Push(LogEvent{
		GroupId: deal.GroupId,
		ActorId: deal.CreatedBy,
		LogType: "UPDATE_DEAL_STATUS",
		LogData: "Updated deal status to " + string(status),
	})

	return nil
}
