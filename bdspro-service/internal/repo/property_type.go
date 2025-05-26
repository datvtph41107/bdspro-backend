package repo

import (
	"bdspro/internal/domain"
)

type PropertyTypeRepo interface {
	InterText(text string, limit int) ([]uint64, error)
	InterTextToItem(text string, limit int) ([]domain.PropertyTypeItem, error)
}
