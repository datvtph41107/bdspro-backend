package usecase

import (
	_enum "common/domain/enum"
	"context"
	"time"

	"notification/internal/domain"
	"notification/internal/dto"
	shared_enum "pb/enums"
	sharepb "pb/types/shared"
	userpb "pb/types/user"
)

// NotificationStore is the persistence capability consumed by NotificationUsecase.
// It belongs to the application owner; PostgreSQL only implements it.
type NotificationStore interface {
	Search(ctx context.Context, profileID uint64, dto *dto.SearchNotiRequest) ([]domain.NotificationEntity, int32, error)
	Create(ctx context.Context, notification *domain.NotificationEntity) error
	CreateBatch(ctx context.Context, notifications []domain.NotificationEntity) error
	MarkAsRead(ctx context.Context, profileID uint64, id uint64) (int64, error)
	MarkAsReadMany(ctx context.Context, profileID uint64, ids []uint64) (int64, error)
	ReadAll(ctx context.Context, profileID uint64) (int64, error)
	Count(ctx context.Context, profileID uint64) (int64, error)
	Delete(ctx context.Context, profileID uint64, id uint64) error
	ListNotiWithTargetIdAndOwnerType(ctx context.Context, targetID uint64, ownerType int32, hour int) ([]domain.NotificationEntity, error)
	DeleteIds(ctx context.Context, ids []uint64) error
	RemoveUnique(ctx context.Context, ownerID, targetID uint64, notificationType _enum.ENotificationType, ownerOf _enum.EOwnerOf) (*domain.NotificationEntity, error)
	RemoveUnique2(ctx context.Context, ownerID uint64, targetID *uint64, notificationType _enum.ENotificationType, ownerOf _enum.EOwnerOf, attachData []string) (*domain.NotificationEntity, error)
}

// PushTokenResolver is the narrow User capability Notification needs for push delivery.
type PushTokenResolver interface {
	GetPushTokensByProfileId(ctx context.Context, profileID uint64) ([]string, error)
}

// PushSender is the external push capability consumed by Notification application code.
type PushSender interface {
	SendPushNotification(ctx context.Context, tokens []string, title string, message string, data map[string]string) error
}

// UserProfileReader is the narrow profile projection capability consumed by audit/history usecases.
type UserProfileReader interface {
	GetProfileByIds(ctx context.Context, in *sharepb.GetProfileByIdsRequest) (*sharepb.GetProfileByIdsResponse, error)
	GetMapProfileByIds(ctx context.Context, in *sharepb.GetProfileByIdsRequest) (map[uint64]*sharepb.ProfileItem, error)
	GetProfileByID(ctx context.Context, in *userpb.GetProfileByIDRequest) (*userpb.ProfileResponse, error)
}

type AccountWarningStore interface {
	Create(ctx context.Context, warning *domain.AccountWarningEntity) (*domain.AccountWarningEntity, error)
	Update(ctx context.Context, warning *domain.AccountWarningEntity) (*domain.AccountWarningEntity, error)
	Delete(ctx context.Context, id uint64) error
	SoftDelete(ctx context.Context, id uint64) error
	FindByID(ctx context.Context, id uint64) (*domain.AccountWarningEntity, error)
	FindAll(ctx context.Context, page, size int, filters map[string]interface{}) ([]*domain.AccountWarningEntity, int64, error)
	FindByTargetID(ctx context.Context, targetID uint64, targetType string) ([]*domain.AccountWarningEntity, error)
	FindBySentBy(ctx context.Context, sentBy uint64, page, size int) ([]*domain.AccountWarningEntity, int64, error)
	FindUnreadByTargetID(ctx context.Context, targetID uint64) ([]*domain.AccountWarningEntity, error)
	MarkAsRead(ctx context.Context, id uint64, targetID uint64) error
	MarkAsAcknowledged(ctx context.Context, id uint64, targetID uint64) error
	CheckDuplicateWarning(ctx context.Context, targetID uint64, content string, hours int) (bool, error)
	CountByTargetID(ctx context.Context, targetID uint64) (int64, error)
	CreateTemplate(ctx context.Context, template *domain.AccountWarningTemplateEntity) (*domain.AccountWarningTemplateEntity, error)
	UpdateTemplate(ctx context.Context, template *domain.AccountWarningTemplateEntity) (*domain.AccountWarningTemplateEntity, error)
	DeleteTemplate(ctx context.Context, id uint64) error
	FindTemplateByID(ctx context.Context, id uint64) (*domain.AccountWarningTemplateEntity, error)
	FindTemplatesByType(ctx context.Context, warningType string) ([]*domain.AccountWarningTemplateEntity, error)
	FindAllTemplates(ctx context.Context, page, size int) ([]*domain.AccountWarningTemplateEntity, int64, error)
	CreateLog(ctx context.Context, log *domain.AccountWarningLogEntity) error
	FindLogsByWarningID(ctx context.Context, warningID uint64, page, size int) ([]*domain.AccountWarningLogEntity, int64, error)
	FindLogsByTargetID(ctx context.Context, targetID uint64, page, size int) ([]*domain.AccountWarningLogEntity, int64, error)
}

type AdminHistoryStore interface {
	Search(ctx context.Context, adminID uint64, dto dto.AdminHistorySearchDTO) ([]domain.AdminHistoryEntity, int64, error)
	SearchByOwner(ctx context.Context, ownerID uint64, ownerType shared_enum.EOwnerType, dto dto.AdminHistorySearchDTO) ([]domain.AdminHistoryEntity, int64, error)
	CreateAdminHistory(ctx context.Context, entity *domain.AdminHistoryEntity) error
	CreateAdminHistoryBatch(ctx context.Context, entities []*domain.AdminHistoryEntity) error
	GetAdminActions(ctx context.Context, adminID uint64, fromDate time.Time, toDate time.Time) ([]domain.AdminHistoryEntity, error)
	GetRecentAdminActions(ctx context.Context, limit int) ([]domain.AdminHistoryEntity, error)
	GetAdminActionStats(ctx context.Context, adminID uint64, fromDate time.Time, toDate time.Time) (map[string]int64, error)
	SearchInternal(ctx context.Context, dto dto.AdminHistorySearchDTO) ([]domain.AdminHistoryEntity, int64, error)
}

type DealHistoryStore interface {
	Create(ctx context.Context, history *domain.DealHistoryEntity) (*domain.DealHistoryEntity, error)
	GetByDealID(ctx context.Context, dealID uint64, page, size int) ([]*domain.DealHistoryEntity, uint32, error)
	GetByID(ctx context.Context, id uint64) (*domain.DealHistoryEntity, error)
	Update(ctx context.Context, history *domain.DealHistoryEntity) (*domain.DealHistoryEntity, error)
	Delete(ctx context.Context, id uint64) error
}

type HistoryAuthSearch struct {
	UserID         uint64
	OrganizationID *uint64
	ActionType     string
	Channel        string
	DeviceID       string
	Success        *bool
	FromDate       *time.Time
	ToDate         *time.Time
	Offset         int
	Limit          int
}

type HistoryAuthStore interface {
	CreateHistoryAuth(ctx context.Context, entity *domain.HistoryAuthEntity) error
	SearchHistoryAuth(ctx context.Context, params HistoryAuthSearch) ([]domain.HistoryAuthEntity, int64, error)
}

type HistoryStore interface {
	Search(ctx context.Context, ownerID uint64, ownerType shared_enum.EOwnerType, dto dto.HistorySearchDTO) ([]domain.HistoryEntity, int64, error)
	ContactHistory(ctx context.Context, contactID uint64, dto dto.HistorySearchDTO) ([]domain.HistoryEntity, int64, error)
	AssetHistory(ctx context.Context, assetID uint64, dto dto.HistorySearchDTO) ([]domain.HistoryEntity, int64, error)
	ProductHistory(ctx context.Context, productID uint64, dto dto.HistorySearchDTO) ([]domain.HistoryEntity, int64, error)
	ProductChildHistory(ctx context.Context, productID uint64, dto dto.HistorySearchDTO) ([]domain.HistoryEntity, int64, error)
	RuleEventHistory(ctx context.Context, ruleEventID uint64, dto dto.HistorySearchDTO) ([]domain.HistoryEntity, int64, error)
	CreateHistory(ctx context.Context, entity *domain.HistoryEntity) error
	CrmHistory(ctx context.Context, crmID uint64, dto dto.HistorySearchDTO) ([]domain.HistoryEntity, int64, error)
	ExistedFromDate(ctx context.Context, customerID uint64, fromDate time.Time) (bool, error)
	SearchInternal(ctx context.Context, dto dto.HistorySearchDTO) ([]domain.HistoryEntity, int64, error)
	CampaignHistory(ctx context.Context, campaignID uint64, dto dto.HistorySearchDTO) ([]domain.HistoryEntity, int64, error)
	PackageHistory(ctx context.Context, packageID uint64, dto dto.HistorySearchDTO) ([]domain.HistoryEntity, int64, error)
}

type PersonConfigStore interface {
	BatchUpsert(ctx context.Context, configs []*domain.PersonConfigEntity) ([]*domain.PersonConfigEntity, error)
	ListByUserID(ctx context.Context, userID uint64) ([]*domain.PersonConfigEntity, error)
}

type PropertyHistoryStore interface {
	Create(ctx context.Context, entity *domain.PropertyHistory) (*domain.PropertyHistory, error)
	GetByID(ctx context.Context, id uint64) (*domain.PropertyHistory, error)
	Delete(ctx context.Context, id uint64) error
	Search(ctx context.Context, query *dto.PropertyHistoryQueryDTO) ([]*domain.PropertyHistory, int64, error)
}
