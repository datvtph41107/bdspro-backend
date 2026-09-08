package usecases

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/repo"
	_dto "common/domain/dto"
	"context"
)

type ListUsecase struct {
	// Repo repo.ListRepo
	propertyTypeRepo repo.PropertyTypeRepo
	regionRepo       repo.RegionRepo
	areaRegionRepo   repo.AreaRegionRepo
	docTypeRepo      repo.DocTypeRepo
	projectRepo      repo.ProjectRepo
	amenityRepo      repo.PropertyAmenityRepo
	projectBuildRepo repo.ProjectBuildRepo
}

func NewListUsecase(propertyTypeRepo repo.PropertyTypeRepo,
	regionRepo repo.RegionRepo,
	areaRegionRepo repo.AreaRegionRepo,
	docTypeRepo repo.DocTypeRepo,
	projectRepo repo.ProjectRepo,
	amenityRepo repo.PropertyAmenityRepo,
	projectBuildRepo repo.ProjectBuildRepo,
) *ListUsecase {
	return &ListUsecase{propertyTypeRepo: propertyTypeRepo,
		regionRepo:       regionRepo,
		areaRegionRepo:   areaRegionRepo,
		docTypeRepo:      docTypeRepo,
		projectRepo:      projectRepo,
		amenityRepo:      amenityRepo,
		projectBuildRepo: projectBuildRepo,
	}
}

func (s *ListUsecase) ListPropertyType(c context.Context, dto *dto.PropertyTypeSearchDTO) ([]domain.PropertyType, error) {
	return s.propertyTypeRepo.GetAllItem(c, dto)
}

func (s *ListUsecase) ListRegion(c context.Context, dto *dto.RegionRequest) ([]domain.Region, int64, error) {
	return s.regionRepo.Search(c, dto)
}

func (s *ListUsecase) ListAreaRegion(c context.Context, dto *_dto.Pagable) ([]domain.AreaRegion, int64, error) {
	return s.areaRegionRepo.Search(c, dto)
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

func (s *ListUsecase) ListProjectBuild(c context.Context, dto *dto.TextSearchRequest) ([]domain.ProjectBuild, int64, error) {
	return s.projectBuildRepo.Search(c, dto)
}
