package rest

type SearchDTO[DTO any] struct {
	Page int `form:"page"`
	Size int `form:"size"`
	DTO  DTO
}
