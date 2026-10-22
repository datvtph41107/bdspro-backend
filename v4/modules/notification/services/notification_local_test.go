package services_test

import (
	"context"
	"testing"

	_utils "common/utils"
	"mdn4/modules/notification/dto"
	"mdn4/modules/notification/entity"
	"mdn4/modules/notification/services"
)

type stubRepo struct {
	unread int32
}

func (s *stubRepo) CountUnread(ctx context.Context, userID uint64) (int32, error) {
	return s.unread, nil
}
func (s *stubRepo) List(ctx context.Context, userID uint64, req dto.ListRequest) ([]*entity.Notification, int32, error) {
	return nil, 0, nil
}
func (s *stubRepo) MarkRead(ctx context.Context, userID, id uint64) (int64, error) { return 1, nil }
func (s *stubRepo) MarkReadMany(ctx context.Context, userID uint64, ids []uint64) (int64, error) {
	return int64(len(ids)), nil
}
func (s *stubRepo) MarkReadAll(ctx context.Context, userID uint64) (int64, error) { return 3, nil }
func (s *stubRepo) Remove(ctx context.Context, userID, id uint64) error           { return nil }
func (s *stubRepo) Create(ctx context.Context, n *entity.Notification) error      { return nil }
func (s *stubRepo) CreateActivityHistory(ctx context.Context, h *entity.ActivityHistory) error {
	return nil
}
func (s *stubRepo) CreateAdminHistory(ctx context.Context, h *entity.AdminHistory) error { return nil }
func (s *stubRepo) SearchAdminHistory(ctx context.Context, filter dto.AdminHistoryFilter) ([]*entity.AdminHistory, int64, error) {
	return nil, 0, nil
}
func (s *stubRepo) CreateHistoryAuth(ctx context.Context, h *entity.HistoryAuth) error { return nil }
func (s *stubRepo) ListPersonConfigs(ctx context.Context, userID uint64) ([]*entity.PersonConfig, error) {
	return nil, nil
}
func (s *stubRepo) BatchUpsertPersonConfigs(ctx context.Context, configs []*entity.PersonConfig) error {
	return nil
}

func TestLocalNotificationServiceCount(t *testing.T) {
	svc := services.NewLocalNotificationService(&stubRepo{unread: 7}, (*_utils.SyncUtil)(nil), nil)
	resp, err := svc.Count(context.Background(), dto.CountRequest{})
	if err != nil {
		t.Fatalf("Count: %v", err)
	}
	if resp.GetCount() != 7 || resp.GetUnread() != 7 {
		t.Fatalf("unexpected count=%d unread=%d", resp.GetCount(), resp.GetUnread())
	}
}
