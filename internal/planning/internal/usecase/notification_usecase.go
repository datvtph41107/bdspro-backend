package usecase

import (
	"context"
	"tqd/internal/domain"
	"tqd/internal/interface/repo"
)

type NotificationUsecase interface {
	List(ctx context.Context, userID uint64, notifType *uint32, read *bool, page, limit int) ([]domain.UserNotification, int64, int64, error)
	MarkRead(ctx context.Context, userID uint64, notificationID uint64) error
	MarkAllRead(ctx context.Context, userID uint64) error
}

type notificationUsecase struct {
	repo repo.NotificationRepository
}

func NewNotificationUsecase(repo repo.NotificationRepository) NotificationUsecase {
	return &notificationUsecase{repo: repo}
}

func (u *notificationUsecase) List(ctx context.Context, userID uint64, notifType *uint32, read *bool, page, limit int) ([]domain.UserNotification, int64, int64, error) {
	notifs, total, err := u.repo.ListByUser(ctx, userID, notifType, read, page, limit)
	if err != nil {
		return nil, 0, 0, err
	}
	unread, _ := u.repo.CountUnread(ctx, userID)
	return notifs, total, unread, nil
}

func (u *notificationUsecase) MarkRead(ctx context.Context, userID uint64, notificationID uint64) error {
	return u.repo.MarkAsRead(ctx, userID, notificationID)
}

func (u *notificationUsecase) MarkAllRead(ctx context.Context, userID uint64) error {
	return u.repo.MarkAllAsRead(ctx, userID)
}
