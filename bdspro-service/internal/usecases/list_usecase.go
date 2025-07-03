package usecases

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/repo"
	"context"
)

type ListUsecase struct {
	// Repo repo.ListRepo
	propertyTypeRepo repo.PropertyTypeRepo
	regionRepo       repo.RegionRepo
	docTypeRepo      repo.DocTypeRepo
	projectRepo      repo.ProjectRepo
	amenityRepo      repo.AmenityRepo
}

func NewListUsecase(propertyTypeRepo repo.PropertyTypeRepo,
	regionRepo repo.RegionRepo,
	docTypeRepo repo.DocTypeRepo,
	projectRepo repo.ProjectRepo,
	amenityRepo repo.AmenityRepo,
) *ListUsecase {
	return &ListUsecase{propertyTypeRepo: propertyTypeRepo,
		regionRepo:  regionRepo,
		docTypeRepo: docTypeRepo,
		projectRepo: projectRepo,
		amenityRepo: amenityRepo,
	}
}

func (s *ListUsecase) ListPropertyType(c context.Context, dto *dto.PropertyTypeSearchDTO) ([]domain.PropertyTypeItem, error) {
	return s.propertyTypeRepo.GetAllItem(c, dto)
}

func (s *ListUsecase) ListRegion(c context.Context, dto *dto.RegionRequest) ([]domain.Region, int64, error) {
	return s.regionRepo.Search(c, dto)
}

func (s *ListUsecase) ListDocType(c context.Context, dto *dto.DocTypeSearchDTO) ([]domain.DocTypeItem, error) {
	data, err := s.docTypeRepo.GetAllItem(c)
	return data, err
}

func (s *ListUsecase) ListProject(c context.Context, dto *dto.ProjectSearchDTO) ([]domain.ProjectItem, int64, error) {
	return s.projectRepo.SearchItem(c, dto)
}

func (s *ListUsecase) ListAmenity(c context.Context, dto *dto.AmenitySearchDTO) ([]domain.AmenityItem, error) {
	return s.amenityRepo.GetAll(c, dto)
}
