package repo

import (
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/enums"
)

type HistoryRepo interface {
	Search(c context.Context,
		ownerID uint64,
		ownerType enums.EOwnerType,
		dto dto.HistorySearchDTO,
	) ([]domain.HistoryEntity, int64, error)
	CreateHistory(c context.Context, entity *domain.HistoryEntity) error
}
