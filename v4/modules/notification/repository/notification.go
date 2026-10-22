package repository

import (
	"context"
	"fmt"

	_dto "common/domain/dto"
	"mdn4/modules/notification/dto"
	"mdn4/modules/notification/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// NotificationRepository is the local inbox + related notification store.
type NotificationRepository interface {
	CountUnread(ctx context.Context, userID uint64) (int32, error)
	List(ctx context.Context, userID uint64, req dto.ListRequest) ([]*entity.Notification, int32, error)
	MarkRead(ctx context.Context, userID, id uint64) (int64, error)
	MarkReadMany(ctx context.Context, userID uint64, ids []uint64) (int64, error)
	MarkReadAll(ctx context.Context, userID uint64) (int64, error)
	Remove(ctx context.Context, userID, id uint64) error
	Create(ctx context.Context, n *entity.Notification) error

	CreateActivityHistory(ctx context.Context, h *entity.ActivityHistory) error
	CreateAdminHistory(ctx context.Context, h *entity.AdminHistory) error
	SearchAdminHistory(ctx context.Context, filter dto.AdminHistoryFilter) ([]*entity.AdminHistory, int64, error)
	CreateHistoryAuth(ctx context.Context, h *entity.HistoryAuth) error
	ListPersonConfigs(ctx context.Context, userID uint64) ([]*entity.PersonConfig, error)
	BatchUpsertPersonConfigs(ctx context.Context, configs []*entity.PersonConfig) error
}

type notificationPostgres struct{ db *gorm.DB }

func NewNotificationPostgres(db *gorm.DB) NotificationRepository {
	return &notificationPostgres{db: db}
}

func (r *notificationPostgres) visibleScope(q *gorm.DB) *gorm.DB {
	return q.Where("deleted_at IS NULL AND (visible_at IS NULL OR visible_at < CURRENT_TIMESTAMP)")
}

func (r *notificationPostgres) CountUnread(ctx context.Context, userID uint64) (int32, error) {
	if r.db == nil || userID == 0 {
		return 0, nil
	}
	var count int64
	err := r.visibleScope(r.db.WithContext(ctx).Model(&entity.Notification{})).
		Where("owner_id = ? AND (is_read IS NULL OR is_read = FALSE)", userID).
		Count(&count).Error
	return int32(count), err
}

func (r *notificationPostgres) List(ctx context.Context, userID uint64, req dto.ListRequest) ([]*entity.Notification, int32, error) {
	if r.db == nil || userID == 0 {
		return nil, 0, nil
	}
	page := _dto.Pagable{Page: uint32(req.Page), Size: uint32(req.Size)}

	base := func() *gorm.DB {
		q := r.visibleScope(r.db.WithContext(ctx).Model(&entity.Notification{})).
			Where("owner_id = ?", userID)
		if req.IsRead {
			q = q.Where("is_read = true")
		}
		return q
	}

	var total int64
	if err := base().Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var rows []*entity.Notification
	err := base().
		Order("created_at DESC").
		Offset(page.GetOffset()).
		Limit(page.GetLimit()).
		Find(&rows).Error
	return rows, int32(total), err
}

func (r *notificationPostgres) MarkRead(ctx context.Context, userID, id uint64) (int64, error) {
	if r.db == nil {
		return 0, nil
	}
	res := r.visibleScope(r.db.WithContext(ctx).Model(&entity.Notification{})).
		Where("id = ? AND owner_id = ?", id, userID).
		Update("is_read", true)
	return res.RowsAffected, res.Error
}

func (r *notificationPostgres) MarkReadMany(ctx context.Context, userID uint64, ids []uint64) (int64, error) {
	if r.db == nil || len(ids) == 0 {
		return 0, nil
	}
	res := r.visibleScope(r.db.WithContext(ctx).Model(&entity.Notification{})).
		Where("id IN ? AND owner_id = ?", ids, userID).
		Update("is_read", true)
	return res.RowsAffected, res.Error
}

func (r *notificationPostgres) MarkReadAll(ctx context.Context, userID uint64) (int64, error) {
	if r.db == nil || userID == 0 {
		return 0, nil
	}
	res := r.visibleScope(r.db.WithContext(ctx).Model(&entity.Notification{})).
		Where("owner_id = ? AND (is_read = false OR is_read IS NULL)", userID).
		Update("is_read", true)
	return res.RowsAffected, res.Error
}

func (r *notificationPostgres) Remove(ctx context.Context, userID, id uint64) error {
	if r.db == nil {
		return nil
	}
	return r.db.WithContext(ctx).Unscoped().
		Where("owner_id = ? AND id = ? AND (visible_at IS NULL OR visible_at < CURRENT_TIMESTAMP)", userID, id).
		Delete(&entity.Notification{}).Error
}

func (r *notificationPostgres) Create(ctx context.Context, n *entity.Notification) error {
	if r.db == nil || n == nil {
		return nil
	}
	return r.db.WithContext(ctx).Create(n).Error
}

func (r *notificationPostgres) CreateActivityHistory(ctx context.Context, h *entity.ActivityHistory) error {
	if r.db == nil || h == nil {
		return nil
	}
	return r.db.WithContext(ctx).Create(h).Error
}

func (r *notificationPostgres) CreateAdminHistory(ctx context.Context, h *entity.AdminHistory) error {
	if r.db == nil || h == nil {
		return nil
	}
	return r.db.WithContext(ctx).Create(h).Error
}

func (r *notificationPostgres) SearchAdminHistory(ctx context.Context, filter dto.AdminHistoryFilter) ([]*entity.AdminHistory, int64, error) {
	if r.db == nil {
		return []*entity.AdminHistory{}, 0, nil
	}
	q := r.db.WithContext(ctx).Model(&entity.AdminHistory{}).Where("deleted_at IS NULL")
	if filter.AdminID != nil {
		q = q.Where("admin_id = ?", *filter.AdminID)
	}
	if filter.OwnerID != nil {
		q = q.Where("owner_id = ?", *filter.OwnerID)
	}
	if filter.OwnerType != nil {
		q = q.Where("owner_of = ?", *filter.OwnerType)
	}
	if filter.TargetID != nil {
		q = q.Where("target_id = ?", *filter.TargetID)
	}
	if filter.TargetType != nil {
		q = q.Where("target_type = ?", *filter.TargetType)
	}
	if filter.ActionType != nil {
		q = q.Where("action_type = ?", *filter.ActionType)
	}
	if filter.AdminRole != nil && *filter.AdminRole != "" {
		q = q.Where("admin_role = ?", *filter.AdminRole)
	}
	if filter.FromDate != nil && *filter.FromDate != "" {
		q = q.Where("created_at >= ?", *filter.FromDate)
	}
	if filter.ToDate != nil && *filter.ToDate != "" {
		q = q.Where("created_at <= ?", *filter.ToDate)
	}
	if filter.IsInternal != nil {
		q = q.Where("is_internal = ?", *filter.IsInternal)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	size := filter.Size
	if size <= 0 {
		size = 20
	}
	page := filter.Page
	if page < 0 {
		page = 0
	}
	var rows []*entity.AdminHistory
	err := q.Order("created_at DESC").Offset(page * size).Limit(size).Find(&rows).Error
	return rows, total, err
}

func (r *notificationPostgres) CreateHistoryAuth(ctx context.Context, h *entity.HistoryAuth) error {
	if r.db == nil || h == nil {
		return nil
	}
	return r.db.WithContext(ctx).Create(h).Error
}

func (r *notificationPostgres) ListPersonConfigs(ctx context.Context, userID uint64) ([]*entity.PersonConfig, error) {
	if r.db == nil || userID == 0 {
		return []*entity.PersonConfig{}, nil
	}
	var rows []*entity.PersonConfig
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("id asc").Find(&rows).Error
	return rows, err
}

func (r *notificationPostgres) BatchUpsertPersonConfigs(ctx context.Context, configs []*entity.PersonConfig) error {
	if r.db == nil || len(configs) == 0 {
		return nil
	}
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	for _, cfg := range configs {
		if cfg == nil {
			continue
		}
		if cfg.ID > 0 {
			res := tx.Model(&entity.PersonConfig{}).
				Where("id = ? AND user_id = ?", cfg.ID, cfg.UserID).
				Updates(map[string]interface{}{
					"key":        cfg.Key,
					"checked":    cfg.Checked,
					"is_default": cfg.IsDefault,
					"channel":    cfg.Channel,
				})
			if res.Error != nil {
				tx.Rollback()
				return res.Error
			}
			if res.RowsAffected > 0 {
				continue
			}
			cfg.ID = 0
		}
		err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "channel"}, {Name: "user_id"}, {Name: "key"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"checked":    cfg.Checked,
				"is_default": cfg.IsDefault,
			}),
		}).Create(cfg).Error
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("upsert person_config: %w", err)
		}
	}
	return tx.Commit().Error
}
