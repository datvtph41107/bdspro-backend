package routes

import (
	profession_service "user/internal/services/user"

	"github.com/gin-gonic/gin"
)

type ProfessionRouter struct {
	Service *profession_service.ProfessionService
}

func NewProfessionRouter(professionService *profession_service.ProfessionService) *ProfessionRouter {
	return &ProfessionRouter{
		Service: professionService,
	}
}

func (router *ProfessionRouter) CreateProfession(c *gin.Context) {
	router.Service.CreateProfession(c)
}
