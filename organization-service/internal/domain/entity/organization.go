package entity

import (
	"time"
)

type Organization struct {
	ID                 uint32
	Name               string
	TaxCode            string
	BusinessLicenseUrl *string
	Address            *string
	Phone              *string
	Email              *string
	Website            *string
	LogoUrl            *string
	Status             string
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
