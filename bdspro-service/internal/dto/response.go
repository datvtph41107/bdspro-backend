package dto

type Response[T any] struct {
	Data  T     `json:"data"`
	Total int64 `json:"total"`
}
