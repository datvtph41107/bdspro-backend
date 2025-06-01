package user_service

import (
	"user/repo"

	"github.com/gin-gonic/gin"
)

type ProfessionService struct {
	Repo *repo.ProfessionRepo
}

func NewProfessionService(repo *repo.ProfessionRepo) *ProfessionService {
	return &ProfessionService{
		Repo: repo,
	}
}

func (s *ProfessionService) CreateProfession(c *gin.Context) {

}
