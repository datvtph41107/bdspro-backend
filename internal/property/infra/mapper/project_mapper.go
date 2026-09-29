package mapper

import (
	"bdspro/internal/domain"
	_models "common/models"
	_utils "common/utils"
	bdspropb "pb/types/bdspro"
)

type ProjectMapper struct {
}

func NewProjectMapper() *ProjectMapper {
	return &ProjectMapper{}
}

func (m *ProjectMapper) MapProjectPbItem(project *domain.Project) *bdspropb.Project {
	return &bdspropb.Project{
		Id:          project.ID,
		Name:        project.Name,
		DeveloperId: project.DeveloperID,
		Description: project.Description,
		CreatedAt:   _utils.FormatTimeToString(project.CreatedAt),
		UpdatedAt:   _utils.FormatTimeToString(project.UpdatedAt),
		Active:      project.DeletedAt == nil,
	}
}

// MapProjectPbList map danh sách project domain sang pb
func (m *ProjectMapper) MapProjectPbList(projects []domain.Project) []*bdspropb.Project {
	var result []*bdspropb.Project
	for _, project := range projects {
		result = append(result, m.MapProjectPb(&project))
	}
	return result
}

// MapProjectPb map project domain sang pb
func (m *ProjectMapper) MapProjectPb(project *domain.Project) *bdspropb.Project {
	return &bdspropb.Project{
		Id:          project.ID,
		Name:        project.Name,
		DeveloperId: project.DeveloperID,
		Description: project.Description,
		CreatedAt:   project.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   project.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		Active:      project.DeletedAt == nil,
	}
}

// ProjectPbToDomain map project pb sang domain
func (m *ProjectMapper) ProjectPbToDomain(project *bdspropb.Project) *domain.Project {
	return &domain.Project{
		BaseEntity: _models.BaseEntity{
			ID: project.Id,
		},
		Name:        project.Name,
		DeveloperID: project.DeveloperId,
		Description: project.Description,
	}
}
