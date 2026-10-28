package services_test

import (
	"context"
	"testing"

	_dto "common/domain/dto"
	_utils "common/utils"
	"mdn4/modules/notification/dto"
	"mdn4/modules/notification/entity"
	"mdn4/modules/notification/services"
	shared_enum "pb/enums"
)

func TestLocalHistoryReaderContactHistory(t *testing.T) {
	repo := &stubRepo{histories: []*entity.ActivityHistory{{TargetID: 7, Title: "Chuyển giai đoạn"}}}
	items, total, err := services.NewLocalHistoryReader(repo).ContactHistory(context.Background(), 7, _dto.Pagable{})
	if err != nil {
		t.Fatalf("ContactHistory: %v", err)
	}
	if repo.gotTarget != 7 || len(repo.gotActions) != len(shared_enum.HistoryContactAction) {
		t.Fatalf("target=%d actions=%v", repo.gotTarget, repo.gotActions)
	}
	if total != 1 || len(items) != 1 || items[0].GetTitle() != "Chuyển giai đoạn" || items[0].GetActionUser() != "Hệ thống" {
		t.Fatalf("unexpected items=%v total=%d", items, total)
	}
}

type stubRepo struct {
	unread     int32
	histories  []*entity.ActivityHistory
	gotTarget  uint64
	gotActions []shared_enum.EHistory
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
func (s *stubRepo) ListActivityHistory(ctx context.Context, targetID uint64, actionTypes []shared_enum.EHistory, page _dto.Pagable) ([]*entity.ActivityHistory, int64, error) {
	s.gotTarget, s.gotActions = targetID, actionTypes
	return s.histories, int64(len(s.histories)), nil
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
