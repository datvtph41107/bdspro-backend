package client_dto

type GroupDTO struct {
	Name string `json:"name" validate:"required"`
}
