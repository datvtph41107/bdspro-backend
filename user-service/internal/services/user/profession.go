package user_service

import (
	user_repo "user/infra/repo"

	"github.com/gin-gonic/gin"
)

type ProfessionService struct {
	Repo *user_repo.ProfessionRepo
}

func NewProfessionService(repo *user_repo.ProfessionRepo) *ProfessionService {
	return &ProfessionService{
		Repo: repo,
	}
}

func (s *ProfessionService) CreateProfession(c *gin.Context) {

}
