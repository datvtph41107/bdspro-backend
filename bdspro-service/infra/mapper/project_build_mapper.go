package mapper

import (
	"bdspro/internal/domain"
	_models "common/models"
	pb_bdspro "pb/types/bdspro"
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

func (m *ProjectBuildMapper) ProjectBuildItemToPbList(projectBuilds []domain.ProjectBuild) []*pb_bdspro.ItemResponse {
	pbs := make([]*pb_bdspro.ItemResponse, len(projectBuilds))
	for i, projectBuild := range projectBuilds {
		pbs[i] = m.ProjectBuildItemToPb(&projectBuild)
	}
	return pbs
}

func (m *ProjectBuildMapper) ProjectBuildItemToPb(projectBuild *domain.ProjectBuild) *pb_bdspro.ItemResponse {
	return &pb_bdspro.ItemResponse{
		Id:   projectBuild.ID,
		Name: projectBuild.Name,
	}
}

// func (m *ProjectBuildMapper) ApartmentToPb(apartments []domain.ApartmentItem) []*pb_bdspro.Apartment {
// 	pbs := make([]*pb_bdspro.Apartment, len(apartments))
// 	for i, apartment := range apartments {
// 		pbs[i] = m.ApartmentMapper.MapApartmentPbItem(&apartment)
// 	}
// 	return pbs
// }

func (m *ProjectBuildMapper) ApartmentAttributeToPb(attributes []domain.ApartmentAttrItem) []*pb_bdspro.ApartmentAttribute {
	return nil
}

func (m *ProjectBuildMapper) ProjectBuildToPb(projectBuild *domain.ProjectBuild) *pb_bdspro.ProjectBuild {
	result := &pb_bdspro.ProjectBuild{
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

func (m *ProjectBuildMapper) ProjectBuildToDomain(projectBuild *pb_bdspro.ProjectBuild) *domain.ProjectBuild {
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

func (m *ProjectBuildMapper) ProjectBuildDetailToPb(projectBuild *domain.ProjectBuild) *pb_bdspro.ProjectBuild {
	result := m.ProjectBuildToPb(projectBuild)
	result.Apartments = m.ApartmentMapper.MapApartmentPbItems(projectBuild.Apartments)
	result.Attributes = m.ApartmentAttributeToPb(projectBuild.Attributes)
	return result
}
