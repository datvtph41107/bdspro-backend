package bank

import (
	_models "common/models"
)

type Bank struct {
	_models.BaseEntity
	Name        string
	Logo        string
	Code        string
	Active      bool
	Description string
}
