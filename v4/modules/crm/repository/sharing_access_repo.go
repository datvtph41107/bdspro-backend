package repository

import (
	"context"
	domain "userv4/modules/crm/entity"
	"userv4/modules/crm/dto"
)

type SharingAccessRepo interface {
	Bulk(c context.Context, deleteIds []uint64, sharing []domain.SharingEntity) ([]domain.SharingEntity, error)
	List(c context.Context, contactId uint64, dto *dto.SharingAccessSearchDTO) ([]domain.SharingEntity, int64, error)
}