package dto

import qh_domain "tqd/internal/domain/qh"

//	type ParcelDTO struct {
//		ID       uint64  `json:"id"`
//		Name     string  `json:"name"`
//		Address  string  `json:"address"`
//		Category uint64  `json:"categoryId"`
//		Latitude  float64 `json:"latitude"`
//		Longitude float64 `json:"longitude"`
//	}
type ParcelList struct {
	Data []qh_domain.Parcel `json:"data"`
}

type ParcelResponse struct {
	ID         uint64  `json:"id"`
	Name       string  `json:"name"`
	Address    string  `json:"address"`
	CategoryID uint64  `json:"categoryId"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
}
