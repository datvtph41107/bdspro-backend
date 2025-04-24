package u_list

import (
	"bdspro/v1/data"
	a_doc_type "bdspro/v1/mod/admin/doc_type"
	a_project "bdspro/v1/mod/admin/project"
	a_property_type "bdspro/v1/mod/admin/property_type"
	a_amenity "bdspro/v1/mod/global/amenity"
	g_region "bdspro/v1/mod/global/region"
	"bdspro/v1/models"
	_routes "common/routes"
	_utils "common/utils"

	"github.com/gin-gonic/gin"
)

type ListService struct {
	propertyTypeRepo *a_property_type.PropertyTypeRepo
	regionRepo       *g_region.RegionRepo
	docTypeRepo      *a_doc_type.DocTypeRepo
	projectRepo      *a_project.ProjectRepo
	amenityRepo      *a_amenity.AmenityRepo
}

func NewUserListService(propertyTypeRepo *a_property_type.PropertyTypeRepo,
	regionRepo *g_region.RegionRepo,
	docTypeRepo *a_doc_type.DocTypeRepo,
	projectRepo *a_project.ProjectRepo,
	amenityRepo *a_amenity.AmenityRepo,
) *ListService {
	return &ListService{
		propertyTypeRepo: propertyTypeRepo,
		regionRepo:       regionRepo,
		docTypeRepo:      docTypeRepo,
		projectRepo:      projectRepo,
		amenityRepo:      amenityRepo,
	}
}

func (s *ListService) ListPropertyType(c *gin.Context) ([]models.PropertyType, error) {
	return s.propertyTypeRepo.GetAll()
}

func (s *ListService) ListRegion(c *gin.Context) (_routes.ResponseDTO, error) {
	dto, _ := _utils.ParseQuery[data.RegionRequest](c)
	data, totalElement, err := s.regionRepo.Search(dto)
	return _routes.ResponseDTO{
		Code:          0,
		Data:          data,
		TotalElements: &totalElement,
	}, err
}

func (s *ListService) ListDocType(c *gin.Context) (_routes.ResponseDTO, error) {
	// dto, _ := _utils.ParseQuery[data.RegionRequest](c)
	data, err := s.docTypeRepo.GetAll()
	return _routes.ResponseDTO{
		Code: 0,
		Data: data,
	}, err
}

func (s *ListService) ListProject(c *gin.Context) (_routes.ResponseDTO, error) {
	dto, _ := _utils.ParseQuery[data.SearchName](c)
	data, err := s.projectRepo.Search(dto)
	return _routes.ResponseDTO{
		Code: 0,
		Data: data,
	}, err
}

func (s *ListService) ListAmenity(c *gin.Context) (_routes.ResponseDTO, error) {
	// dto, _ := _utils.ParseQuery[data.SearchName](c)
	data, err := s.amenityRepo.GetAll()
	return _routes.ResponseDTO{
		Code: 0,
		Data: data,
	}, err
}
