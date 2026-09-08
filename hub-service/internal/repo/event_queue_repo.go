package repo

import (
	_crud "common/domain/crud"
	"context"
	"hub/internal/domain"
	"hub/internal/dto"
)

type IEventQueueRepo interface {
	_crud.ICrudRepo[domain.EventQueueEntity]
	Search(c context.Context, searchDTO dto.EventQueueSearchDTO) ([]domain.EventQueueEntity, int64, error)
	UpdateStatus(c context.Context, id uint64, status uint32, errorMessage *string) (*domain.EventQueueEntity, error)
	GetPendingEvents(c context.Context, limit int) ([]domain.EventQueueEntity, error)
}
