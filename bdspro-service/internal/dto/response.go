package dto

type Response[T any] struct {
	Data         T     `json:"data"`
	TotalElement int64 `json:"totalElement"`
}
