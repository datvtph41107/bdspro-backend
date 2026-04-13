package dto

import "tqd/internal/domain"

//	type ParcelDTO struct {
//		ID       uint64  `json:"id"`
//		Name     string  `json:"name"`
//		Address  string  `json:"address"`
//		Category uint64  `json:"categoryId"`
//		Latitude  float64 `json:"latitude"`
//		Longitude float64 `json:"longitude"`
//	}
type ParcelList struct {
	Data []domain.Parcel `json:"data"`
}

type ParcelResponse struct {
	ID         uint64  `json:"id"`
	Name       string  `json:"name"`
	Address    string  `json:"address"`
	CategoryID uint64  `json:"categoryId"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
}
