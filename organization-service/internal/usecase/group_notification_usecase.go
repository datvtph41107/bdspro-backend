package usecase

import (
	"context"

	"organization/internal/domain/entity"
	"organization/internal/domain/repository"
)

type GroupNotificationUsecase interface {
	CreateGroupNotification(ctx context.Context, notification *entity.GroupNotification) (*entity.GroupNotification, error)
	GetGroupNotificationByID(ctx context.Context, id uint32) (*entity.GroupNotification, error)
	GetGroupNotificationsByGroupID(ctx context.Context, groupID uint32) ([]*entity.GroupNotification, error)
}

type groupNotificationUsecase struct {
	groupNotificationRepository repository.GroupNotificationRepository
}

func NewGroupNotificationUsecase(groupNotificationRepository repository.GroupNotificationRepository) GroupNotificationUsecase {
	return &groupNotificationUsecase{groupNotificationRepository: groupNotificationRepository}
}

func (u *groupNotificationUsecase) CreateGroupNotification(ctx context.Context, notification *entity.GroupNotification) (*entity.GroupNotification, error) {
	return u.groupNotificationRepository.Create(ctx, notification)
}

func (u *groupNotificationUsecase) GetGroupNotificationByID(ctx context.Context, id uint32) (*entity.GroupNotification, error) {
	return u.groupNotificationRepository.GetByID(ctx, id)
}

func (u *groupNotificationUsecase) GetGroupNotificationsByGroupID(ctx context.Context, groupID uint32) ([]*entity.GroupNotification, error) {
	return u.groupNotificationRepository.GetByGroupID(ctx, groupID)
}
