package usecase

import (
	_routes "common/routes"
	_utils "common/utils"
	"context"
	"crm/internal/domain"
	"crm/internal/dto"
	"crm/internal/enums"
	"crm/internal/repo"
)

type ContactUsecase struct {
	repo repo.ContactRepo
}

func NewContactUsecase(repo repo.ContactRepo) *ContactUsecase {
	return &ContactUsecase{
		repo: repo,
	}
}

func (s *ContactUsecase) CheckPermission(c context.Context, id uint64) error {
	entity, err := s.repo.GetByID(c, id)
	if err != nil {
		return err
	}
	if entity.OwnerType == enums.EOwnerTypeMember {
		profileId := _utils.GetProfileIdWithContext(c)
		if profileId != id {
			return &_routes.Except{
				Code:    403,
				Message: "Bạn không có quyền truy cập",
			}
		}
	}
	return nil
}

func (s *ContactUsecase) Search(c context.Context, dto dto.ContactSearchDTO) ([]domain.ContactEntity, int64, error) {
	profileId := _utils.GetProfileIdWithContext(c)
	contacts, total, err := s.repo.Search(c, profileId, dto)
	if err != nil {
		return nil, 0, err
	}
	return contacts, total, nil
}

func (s *ContactUsecase) Create(c context.Context, entity *domain.ContactEntity) (*domain.ContactEntity, error) {
	if entity.OwnerType == enums.EOwnerTypeMember {
		profileId := _utils.GetProfileIdWithContext(c)
		entity.OwnerID = profileId
	}
	return s.repo.Create(c, entity)
}

func (s *ContactUsecase) Update(c context.Context, id uint64, entity *domain.ContactEntity) (*domain.ContactEntity, error) {
	err := s.CheckPermission(c, id)
	if err != nil {
		return nil, err
	}
	return s.repo.Update(c, entity)
}

func (s *ContactUsecase) Delete(c context.Context, id uint64) error {
	err := s.CheckPermission(c, id)
	if err != nil {
		return err
	}
	return s.repo.Delete(c, id)
}
