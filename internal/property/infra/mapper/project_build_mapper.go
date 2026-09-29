package mapper

import (
	"bdspro/internal/domain"
	_models "common/models"
	bdspropb "pb/types/bdspro"
)

type ProjectBuildMapper struct {
	ApartmentMapper     *ApartmentMapper
	ApartmentAttrMapper *ApartmentAttributeMapper
}

func NewProjectBuildMapper(apartmentMapper *ApartmentMapper, apartmentAttrMapper *ApartmentAttributeMapper) *ProjectBuildMapper {
	return &ProjectBuildMapper{
		ApartmentMapper:     apartmentMapper,
		ApartmentAttrMapper: apartmentAttrMapper,
	}
}

func (m *ProjectBuildMapper) ProjectBuildItemToPbList(projectBuilds []domain.ProjectBuild) []*bdspropb.ItemResponse {
	pbs := make([]*bdspropb.ItemResponse, len(projectBuilds))
	for i, projectBuild := range projectBuilds {
		pbs[i] = m.ProjectBuildItemToPb(&projectBuild)
	}
	return pbs
}

func (m *ProjectBuildMapper) ProjectBuildItemToPb(projectBuild *domain.ProjectBuild) *bdspropb.ItemResponse {
	return &bdspropb.ItemResponse{
		Id:   projectBuild.ID,
		Name: projectBuild.Name,
	}
}

// func (m *ProjectBuildMapper) ApartmentToPb(apartments []domain.ApartmentItem) []*bdspropb.Apartment {
// 	pbs := make([]*bdspropb.Apartment, len(apartments))
// 	for i, apartment := range apartments {
// 		pbs[i] = m.ApartmentMapper.MapApartmentPbItem(&apartment)
// 	}
// 	return pbs
// }

func (m *ProjectBuildMapper) ApartmentAttributeToPb(attributes []domain.ApartmentAttrItem) []*bdspropb.ApartmentAttribute {
	return nil
}

func (m *ProjectBuildMapper) ProjectBuildToPb(projectBuild *domain.ProjectBuild) *bdspropb.ProjectBuild {
	result := &bdspropb.ProjectBuild{
		Id:           projectBuild.ID,
		Name:         projectBuild.Name,
		NumFloor:     uint32(projectBuild.NumFloor),
		NumApartment: uint32(projectBuild.NumApartment),
		Note:         projectBuild.Note,
		ProjectId:    projectBuild.ProjectID,
		Apartments:   m.ApartmentMapper.MapApartmentPbItems(projectBuild.Apartments),
		Attributes:   m.ApartmentAttrMapper.MapApartmentAttributePbItems(projectBuild.Attributes),
	}

	return result
}

func (m *ProjectBuildMapper) ProjectBuildToDomain(projectBuild *bdspropb.ProjectBuild) *domain.ProjectBuild {
	return &domain.ProjectBuild{
		BaseEntity: _models.BaseEntity{
			ID: projectBuild.Id,
		},
		Name:         projectBuild.Name,
		NumFloor:     uint(projectBuild.NumFloor),
		NumApartment: uint(projectBuild.NumApartment),
		Note:         projectBuild.Note,
		ProjectID:    projectBuild.ProjectId,
	}
}

func (m *ProjectBuildMapper) ProjectBuildDetailToPb(projectBuild *domain.ProjectBuild) *bdspropb.ProjectBuild {
	result := m.ProjectBuildToPb(projectBuild)
	result.Apartments = m.ApartmentMapper.MapApartmentPbItems(projectBuild.Apartments)
	result.Attributes = m.ApartmentAttributeToPb(projectBuild.Attributes)
	return result
}
