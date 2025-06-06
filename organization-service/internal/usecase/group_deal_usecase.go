package usecase

import (
	"context"

	"github.com/tuandq2112/organization-service/internal/custom_error"
	"github.com/tuandq2112/organization-service/internal/domain/entity"
	"github.com/tuandq2112/organization-service/internal/domain/repository"
)

type GroupDealUsecase interface {
	CreateGroupDeal(ctx context.Context, deal *entity.Deal) (*entity.Deal, error)
	UpdateGroupDeal(ctx context.Context, deal *entity.Deal) (*entity.Deal, error)
	DeleteGroupDeal(ctx context.Context, id uint32) error
	GetGroupDealByID(ctx context.Context, id uint32) (*entity.Deal, error)
	GetGroupDealsByGroupID(ctx context.Context, groupID uint32) ([]*entity.Deal, error)
	UpdateGroupDealStatus(ctx context.Context, id uint32, status entity.DealStatus) error
}

type groupDealUsecase struct {
	groupDealRepository repository.GroupDealRepository
	logActivityUsecase GroupLogActivityUsecase
}

func NewGroupDealUsecase(groupDealRepository repository.GroupDealRepository, logActivityUsecase GroupLogActivityUsecase) GroupDealUsecase {
	return &groupDealUsecase{
		groupDealRepository: groupDealRepository,
		logActivityUsecase: logActivityUsecase,
	}
}

func (u *groupDealUsecase) CreateGroupDeal(ctx context.Context, deal *entity.Deal) (*entity.Deal, error) {
	createdDeal, err := u.groupDealRepository.Create(ctx, deal)
	if err != nil {
		return nil, err
	}

	go func() {
		log := &entity.GroupLogActivity{
			GroupId: createdDeal.GroupId,
			ActorId: createdDeal.CreatedBy,
			LogType: "CREATE_DEAL",
			LogData: "Created new deal",
		}
		u.logActivityUsecase.CreateGroupLogActivity(context.Background(), log)
	}()

	return createdDeal, nil
}

func (u *groupDealUsecase) UpdateGroupDeal(ctx context.Context, deal *entity.Deal) (*entity.Deal, error) {
	existedDeal, err := u.groupDealRepository.GetByID(ctx, deal.Id)
	if err != nil {
		return nil, err
	}
	
	if existedDeal == nil {
		return nil, custom_error.RecordNotFound("deal not found")
	}

	if deal.CreatedBy != existedDeal.CreatedBy {
		return nil, custom_error.Forbidden("you are not allowed to update this deal")
	}
	
	updatedDeal, err := u.groupDealRepository.Update(ctx, deal)
	if err != nil {
		return nil, err
	}

	go func() {
		log := &entity.GroupLogActivity{
			GroupId: updatedDeal.GroupId,
			ActorId: updatedDeal.CreatedBy,
			LogType: "UPDATE_DEAL",
			LogData: "Updated deal information",
		}
		u.logActivityUsecase.CreateGroupLogActivity(context.Background(), log)
	}()

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

	err = u.groupDealRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	go func() {
		log := &entity.GroupLogActivity{
			GroupId: deal.GroupId,
			ActorId: deal.CreatedBy,
			LogType: "DELETE_DEAL",
			LogData: "Deleted deal",
		}
		u.logActivityUsecase.CreateGroupLogActivity(context.Background(), log)
	}()

	return nil
}

func (u *groupDealUsecase) GetGroupDealByID(ctx context.Context, id uint32) (*entity.Deal, error) {
	return u.groupDealRepository.GetByID(ctx, id)
}

func (u *groupDealUsecase) GetGroupDealsByGroupID(ctx context.Context, groupID uint32) ([]*entity.Deal, error) {
	return u.groupDealRepository.GetByGroupID(ctx, groupID)
}

func (u *groupDealUsecase) UpdateGroupDealStatus(ctx context.Context, id uint32, status entity.DealStatus) error {
	deal, err := u.groupDealRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if deal == nil {
		return custom_error.RecordNotFound("deal not found")
	}

	err = u.groupDealRepository.UpdateStatus(ctx, id, status)
	if err != nil {
		return err
	}

	go func() {
		log := &entity.GroupLogActivity{
			GroupId: deal.GroupId,
			ActorId: deal.CreatedBy,
			LogType: "UPDATE_DEAL_STATUS",
			LogData: "Updated deal status to " + string(status),
		}
		u.logActivityUsecase.CreateGroupLogActivity(context.Background(), log)
	}()

	return nil
}