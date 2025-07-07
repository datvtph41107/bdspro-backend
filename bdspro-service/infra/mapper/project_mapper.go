package mapper

import (
	"bdspro/internal/domain"
	pb_bdspro "pb/types/bdspro"
)

type ProjectMapper struct {
}

func NewProjectMapper() *ProjectMapper {
	return &ProjectMapper{}
}

func (m *ProjectMapper) MapProjectPbItem(project *domain.Project) *pb_bdspro.Project {
	// todo: map thêm developer, description, builds
	return &pb_bdspro.Project{
		Id:   project.ID,
		Name: project.Name,
		// Developer:   project.Developer,
		// Description: project.Description,
		// Builds:      project.Builds,

	}
}
