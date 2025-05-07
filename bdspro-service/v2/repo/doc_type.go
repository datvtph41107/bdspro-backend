package repo

import (
	"bdspro/v2/internal/domain"
)

type DocTypeRepo interface {
	InterText(text string, limit int) ([]uint64, error)
	InterTextToItem(text string, limit int) ([]domain.DocTypeItem, error)
}
