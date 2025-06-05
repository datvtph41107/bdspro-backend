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
}

func NewGroupDealUsecase(groupDealRepository repository.GroupDealRepository) GroupDealUsecase {
	return &groupDealUsecase{groupDealRepository: groupDealRepository}
}

func (u *groupDealUsecase) CreateGroupDeal(ctx context.Context, deal *entity.Deal) (*entity.Deal, error) {
	return u.groupDealRepository.Create(ctx, deal)
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
	
	return u.groupDealRepository.Update(ctx, deal)
}

func (u *groupDealUsecase) DeleteGroupDeal(ctx context.Context, id uint32) error {
	return u.groupDealRepository.Delete(ctx, id)
}

func (u *groupDealUsecase) GetGroupDealByID(ctx context.Context, id uint32) (*entity.Deal, error) {
	return u.groupDealRepository.GetByID(ctx, id)
}

func (u *groupDealUsecase) GetGroupDealsByGroupID(ctx context.Context, groupID uint32) ([]*entity.Deal, error) {
	return u.groupDealRepository.GetByGroupID(ctx, groupID)
}

func (u *groupDealUsecase) UpdateGroupDealStatus(ctx context.Context, id uint32, status entity.DealStatus) error {
	return u.groupDealRepository.UpdateStatus(ctx, id, status)
}