package mapper

import (
	"bdspro/internal/domain"
	bdspropb "pb/types/bdspro"
)

type BuildMapper struct {
}

func NewBuildMapper() *BuildMapper {
	return &BuildMapper{}
}

func (m *BuildMapper) MapBuildPbItem(build *domain.ProjectBuild) *bdspropb.ProjectBuild {
	return &bdspropb.ProjectBuild{
		Id:           build.ID,
		Name:         build.Name,
		NumFloor:     uint32(build.NumFloor),
		NumApartment: uint32(build.NumApartment),
		Note:         build.Note,
		ProjectId:    build.ProjectID,
		// Apartments:   build.Apartments,
		// Attributes:   build.Attributes,
	}
}
