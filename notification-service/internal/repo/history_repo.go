package repo

import (
	"context"
	"noti/internal/domain"
	"noti/internal/dto"
	shared_enum "pb/enums"
	"time"
)

type HistoryRepo interface {
	Search(c context.Context,
		ownerID uint64,
		ownerType shared_enum.EOwnerType,
		dto dto.HistorySearchDTO,
	) ([]domain.HistoryEntity, int64, error)
	ContactHistory(c context.Context, contactId uint64, dto dto.HistorySearchDTO) ([]domain.HistoryEntity, int64, error)
	AssetHistory(c context.Context, assetId uint64, dto dto.HistorySearchDTO) ([]domain.HistoryEntity, int64, error)
	ProductHistory(c context.Context, productId uint64, dto dto.HistorySearchDTO) ([]domain.HistoryEntity, int64, error)
	ProductChildHistory(c context.Context, productId uint64, dto dto.HistorySearchDTO) ([]domain.HistoryEntity, int64, error)
	RuleEventHistory(c context.Context, ruleEventId uint64, dto dto.HistorySearchDTO) ([]domain.HistoryEntity, int64, error)
	CreateHistory(c context.Context, entity *domain.HistoryEntity) error
	ExistedFromDate(c context.Context, customerID uint64, fromDate time.Time) (bool, error)
	// ExistedFromDate(c context.Context, customerID uint64, fromDate time.Time) (bool, error)
}
