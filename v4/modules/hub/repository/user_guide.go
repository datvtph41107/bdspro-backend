package repository

import (
	"context"
	"strings"
	"time"

	_utils "common/utils"
	"mdn4/modules/hub/dto"
	"mdn4/modules/hub/entity"

	"gorm.io/gorm"
)

type UserGuideGroupCount struct {
	GroupKey string `gorm:"column:group_key"`
	Count    int64  `gorm:"column:count"`
}

type UserGuideRepository interface {
	ListSimple(ctx context.Context, req dto.GetUserGuideSimpleRequest) ([]*entity.UserGuide, int64, error)
	GetByID(ctx context.Context, id uint64) (*entity.UserGuide, error)
	GetByKey(ctx context.Context, key string) (*entity.UserGuide, error)
	ListAllByKey(ctx context.Context) ([]*entity.UserGuide, error)
	ListGroups(ctx context.Context) ([]UserGuideGroupCount, error)
	Create(ctx context.Context, ug *entity.UserGuide) error
	Update(ctx context.Context, ug *entity.UserGuide) error
	Delete(ctx context.Context, id uint64) error
}

type userGuidePostgres struct{ db *gorm.DB }

func NewUserGuidePostgres(db *gorm.DB) UserGuideRepository {
	return &userGuidePostgres{db: db}
}

func (r *userGuidePostgres) ListSimple(ctx context.Context, req dto.GetUserGuideSimpleRequest) ([]*entity.UserGuide, int64, error) {
	if r.db == nil {
		return nil, 0, nil
	}
	q := r.db.WithContext(ctx).Model(&entity.UserGuide{}).Where("deleted_at IS NULL")
	if req.GroupKey != "" {
		q = q.Where("group_key = ?", req.GroupKey)
	}
	if text := strings.TrimSpace(req.Text); text != "" {
		like := "%" + text + "%"
		q = q.Where("title ILIKE ? OR description ILIKE ?", like, like)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset, size := pageOffset(int(req.Page), int(req.Size))
	var list []*entity.UserGuide
	err := q.Order("id ASC").Offset(offset).Limit(size).Find(&list).Error
	return list, total, err
}

func (r *userGuidePostgres) GetByID(ctx context.Context, id uint64) (*entity.UserGuide, error) {
	if r.db == nil {
		return nil, nil
	}
	var row entity.UserGuide
	err := r.db.WithContext(ctx).
		Preload("Steps", func(db *gorm.DB) *gorm.DB {
			return db.Where("deleted_at IS NULL").Order("step_order ASC")
		}).
		Where("id = ? AND deleted_at IS NULL", id).First(&row).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *userGuidePostgres) GetByKey(ctx context.Context, key string) (*entity.UserGuide, error) {
	if r.db == nil || key == "" {
		return nil, nil
	}
	var row entity.UserGuide
	err := r.db.WithContext(ctx).
		Preload("Steps", func(db *gorm.DB) *gorm.DB {
			return db.Where("deleted_at IS NULL").Order("step_order ASC")
		}).
		Where("key = ? AND deleted_at IS NULL", key).First(&row).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *userGuidePostgres) ListAllByKey(ctx context.Context) ([]*entity.UserGuide, error) {
	if r.db == nil {
		return nil, nil
	}
	var list []*entity.UserGuide
	err := r.db.WithContext(ctx).
		Where("deleted_at IS NULL AND key IS NOT NULL AND key <> ''").
		Order("key ASC").Find(&list).Error
	return list, err
}

func (r *userGuidePostgres) ListGroups(ctx context.Context) ([]UserGuideGroupCount, error) {
	if r.db == nil {
		return nil, nil
	}
	var groups []UserGuideGroupCount
	err := r.db.WithContext(ctx).Model(&entity.UserGuide{}).
		Select("group_key, COUNT(*) as count").
		Where("deleted_at IS NULL").
		Group("group_key").
		Order("group_key ASC").
		Scan(&groups).Error
	return groups, err
}

func (r *userGuidePostgres) Create(ctx context.Context, ug *entity.UserGuide) error {
	return r.db.WithContext(ctx).Create(ug).Error
}

func (r *userGuidePostgres) Update(ctx context.Context, ug *entity.UserGuide) error {
	return r.db.WithContext(ctx).Model(ug).Updates(map[string]interface{}{
		"title": ug.Title, "description": ug.Description,
		"group_key": ug.GroupKey, "key": ug.Key, "mode": ug.Mode,
	}).Error
}

func (r *userGuidePostgres) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&entity.UserGuide{}, id).Error
}

// Touch sync helper for callers that bump Redis after guide reads.
func EnsureUserGuideSync(ctx context.Context, sync *_utils.SyncUtil) int64 {
	if sync == nil {
		return time.Now().UnixMilli()
	}
	return sync.EnsureTimestamp(ctx, _utils.SyncKeyUserGuideByKey)
}
