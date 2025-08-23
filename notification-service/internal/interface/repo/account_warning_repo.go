package repo

import (
	"context"
	"noti/internal/domain"
)

type AccountWarningRepository interface {
	// CRUD operations
	Create(c context.Context, warning *domain.AccountWarningEntity) (*domain.AccountWarningEntity, error)
	Update(c context.Context, warning *domain.AccountWarningEntity) (*domain.AccountWarningEntity, error)
	Delete(c context.Context, id uint64) error
	SoftDelete(c context.Context, id uint64) error
	FindByID(c context.Context, id uint64) (*domain.AccountWarningEntity, error)

	// List operations
	FindAll(c context.Context, page, size int, filters map[string]interface{}) ([]*domain.AccountWarningEntity, int64, error)
	FindByTargetID(c context.Context, targetID uint64, targetType string) ([]*domain.AccountWarningEntity, error)
	FindBySentBy(c context.Context, sentBy uint64, page, size int) ([]*domain.AccountWarningEntity, int64, error)
	FindUnreadByTargetID(c context.Context, targetID uint64) ([]*domain.AccountWarningEntity, error)

	// Business operations
	MarkAsRead(c context.Context, id uint64, targetID uint64) error
	MarkAsAcknowledged(c context.Context, id uint64, targetID uint64) error
	CheckDuplicateWarning(c context.Context, targetID uint64, content string, hours int) (bool, error)
	CountByTargetID(c context.Context, targetID uint64) (int64, error)

	// Template operations
	CreateTemplate(c context.Context, template *domain.AccountWarningTemplateEntity) (*domain.AccountWarningTemplateEntity, error)
	UpdateTemplate(c context.Context, template *domain.AccountWarningTemplateEntity) (*domain.AccountWarningTemplateEntity, error)
	DeleteTemplate(c context.Context, id uint64) error
	FindTemplateByID(c context.Context, id uint64) (*domain.AccountWarningTemplateEntity, error)
	FindTemplatesByType(c context.Context, warningType string) ([]*domain.AccountWarningTemplateEntity, error)
	FindAllTemplates(c context.Context, page, size int) ([]*domain.AccountWarningTemplateEntity, int64, error)

	// Log operations
	CreateLog(c context.Context, log *domain.AccountWarningLogEntity) error
	FindLogsByWarningID(c context.Context, warningID uint64, page, size int) ([]*domain.AccountWarningLogEntity, int64, error)
	FindLogsByTargetID(c context.Context, targetID uint64, page, size int) ([]*domain.AccountWarningLogEntity, int64, error)
}
