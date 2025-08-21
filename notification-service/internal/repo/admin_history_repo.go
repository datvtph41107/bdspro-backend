package repo

import (
	"context"
	shared_enum "pb/enums"
	"time"

	"noti/internal/domain"
	"noti/internal/dto"
)

type AdminHistoryRepo interface {
	Search(c context.Context,
		adminID uint64,
		dto dto.AdminHistorySearchDTO,
	) ([]domain.AdminHistoryEntity, int64, error)

	SearchByOwner(c context.Context,
		ownerID uint64,
		ownerType shared_enum.EOwnerType,
		dto dto.AdminHistorySearchDTO,
	) ([]domain.AdminHistoryEntity, int64, error)

	CreateAdminHistory(c context.Context, entity *domain.AdminHistoryEntity) error
	CreateAdminHistoryBatch(c context.Context, entities []*domain.AdminHistoryEntity) error

	GetAdminActions(c context.Context, adminID uint64, fromDate time.Time, toDate time.Time) ([]domain.AdminHistoryEntity, error)

	GetRecentAdminActions(c context.Context, limit int) ([]domain.AdminHistoryEntity, error)

	GetAdminActionStats(c context.Context, adminID uint64, fromDate time.Time, toDate time.Time) (map[string]int64, error)

	SearchInternal(c context.Context, dto dto.AdminHistorySearchDTO) ([]domain.AdminHistoryEntity, int64, error)
}
