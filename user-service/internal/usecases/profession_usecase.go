package usecases

import (
	"user/internal/interface/repo"

	"github.com/gin-gonic/gin"
)

type ProfessionUsecase struct {
	Repo repo.IProfessionRepo
}

func NewProfessionUsecase(repo repo.IProfessionRepo) *ProfessionUsecase {
	return &ProfessionUsecase{
		Repo: repo,
	}
}

func (s *ProfessionUsecase) CreateProfession(c *gin.Context) {

}
