package dto

import "bdspro/internal/domain"

type ApartmentRequest struct {
	BuildID   uint64                    `json:"buildId"`
	Attribute domain.ApartmentAttribute `json:"attribute"`
	Elements  []domain.Apartment        `json:"elements"`
}

type ApartmentStatusRequest struct {
	BuildID uint64 `json:"buildId"`
	Data    struct {
		Status   int `json:"status"`
		Archived int `json:"archived"`
	} `json:"data"`
	Apartments []domain.Apartment `json:"apartments"`
}
