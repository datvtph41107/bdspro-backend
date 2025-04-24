package models

import "bdspro/v1/enums"

type OwnershipDoc struct {
	ID        uint64 `gorm:"primaryKey;autoIncrement" json:"id"`
	ImageUrls string
	Type      enums.OwnershipDocType
}
