package repository

import (
	"context"
	"strings"

	"userv4/modules/hub/dto"
	"userv4/modules/hub/entity"

	"gorm.io/gorm"
)

type FAQRepository interface {
	ListSimple(ctx context.Context, req dto.GetFAQSimpleRequest) ([]*entity.FAQ, int64, error)
	GetByID(ctx context.Context, id uint64) (*entity.FAQ, error)
	Create(ctx context.Context, faq *entity.FAQ) error
	Update(ctx context.Context, faq *entity.FAQ) error
	Delete(ctx context.Context, id uint64) error
}

type faqPostgres struct{ db *gorm.DB }

func NewFAQPostgres(db *gorm.DB) FAQRepository { return &faqPostgres{db: db} }

func (r *faqPostgres) ListSimple(ctx context.Context, req dto.GetFAQSimpleRequest) ([]*entity.FAQ, int64, error) {
	if r.db == nil {
		return nil, 0, nil
	}
	q := r.db.WithContext(ctx).Model(&entity.FAQ{}).Where("deleted_at IS NULL")
	if req.GroupKey != "" {
		q = q.Where("group_key = ?", req.GroupKey)
	}
	if text := strings.TrimSpace(req.Text); text != "" {
		like := "%" + text + "%"
		q = q.Where("question ILIKE ? OR answer ILIKE ?", like, like)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	size := int(req.Size)
	if size <= 0 {
		size = 20
	}
	page := int(req.Page)
	if page < 0 {
		page = 0
	}
	var list []*entity.FAQ
	err := q.Order("id ASC").Offset(page * size).Limit(size).Find(&list).Error
	return list, total, err
}

func (r *faqPostgres) GetByID(ctx context.Context, id uint64) (*entity.FAQ, error) {
	if r.db == nil {
		return nil, nil
	}
	var row entity.FAQ
	err := r.db.WithContext(ctx).Where("id = ? AND deleted_at IS NULL", id).First(&row).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *faqPostgres) Create(ctx context.Context, faq *entity.FAQ) error {
	return r.db.WithContext(ctx).Create(faq).Error
}

func (r *faqPostgres) Update(ctx context.Context, faq *entity.FAQ) error {
	return r.db.WithContext(ctx).Model(faq).Updates(map[string]interface{}{
		"question": faq.Question, "answer": faq.Answer, "group_key": faq.GroupKey,
	}).Error
}

func (r *faqPostgres) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&entity.FAQ{}, id).Error
}
