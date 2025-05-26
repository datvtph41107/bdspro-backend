package repo

import (
	"bdspro/internal/domain"
)

type DocTypeRepo interface {
	InterText(text string, limit int) ([]uint64, error)
	InterTextToItem(text string, limit int) ([]domain.DocTypeItem, error)
}
