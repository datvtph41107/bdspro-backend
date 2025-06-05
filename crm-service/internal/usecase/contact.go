package usecase

import (
	_jwt "common/jwt"
	_routes "common/routes"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/enums"
	"crm/internal/repo"

	"github.com/gin-gonic/gin"
)

type ContactUsecase struct {
	repo repo.ContactRepo
}

func NewContactService(repo repo.ContactRepo) *ContactUsecase {
	return &ContactUsecase{
		repo: repo,
	}
}

func (s *ContactUsecase) CheckPermission(c *gin.Context, id uint64) error {
	entity, err := s.repo.GetByID(c, id)
	if err != nil {
		return err
	}
	if entity.OwnerType == enums.EOwnerTypeMember {
		profileId := _jwt.GetProfileId(c)
		if profileId != id {
			return &_routes.Except{
				Code:    403,
				Message: "Bạn không có quyền truy cập",
			}
		}
	}
	return nil
}

func (s *ContactUsecase) Search(c *gin.Context, dto dto.ContactSearchDTO) ([]domain.ContactEntity, int64, error) {
	contacts, total, err := s.repo.Search(c, dto)
	if err != nil {
		return nil, 0, err
	}
	return contacts, total, nil
}

func (s *ContactUsecase) Create(c *gin.Context, entity domain.ContactEntity) (*domain.ContactEntity, error) {
	if entity.OwnerType == enums.EOwnerTypeMember {
		profileId := _jwt.GetProfileId(c)
		entity.OwnerID = profileId
	}
	return s.repo.Create(c, &entity)
}

func (s *ContactUsecase) Update(c *gin.Context, id uint64, entity domain.ContactEntity) (*domain.ContactEntity, error) {
	err := s.CheckPermission(c, id)
	if err != nil {
		return nil, err
	}
	return s.repo.Update(c, &entity)
}

func (s *ContactUsecase) Delete(c *gin.Context, id uint64) error {
	err := s.CheckPermission(c, id)
	if err != nil {
		return err
	}
	return s.repo.Delete(c, id)
}
