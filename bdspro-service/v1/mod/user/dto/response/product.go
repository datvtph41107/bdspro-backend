package user_dto

import "bdspro/v1/models"

type CreateResponse struct {
	Product *models.Product `json:"product"`
	Post    *models.Post    `json:"post"`
}
