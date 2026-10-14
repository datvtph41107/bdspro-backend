package mapper

import (
	_models "common/models"
	_utils "common/utils"
	bdspropb "pb/types/bdspro"
	"userv4/modules/bdspro/entity"
)

// ProjectBuildToPb maps a project_build row to the gateway contract.
func ProjectBuildToPb(build *entity.ProjectBuild) *bdspropb.ProjectBuild {
	if build == nil {
		return nil
	}
	out := &bdspropb.ProjectBuild{
		Id:           build.ID,
		Name:         build.Name,
		NumFloor:     uint32(build.NumFloor),
		NumApartment: uint32(build.NumApartment),
		Note:         build.Note,
		ProjectId:    build.ProjectID,
		Apartments:   apartmentsToPb(build.Apartments),
		Attributes:   apartmentAttrsToPb(build.Attributes),
	}
	return out
}

func ProjectBuildFromPb(pb *bdspropb.ProjectBuild) *entity.ProjectBuild {
	if pb == nil {
		return nil
	}
	return &entity.ProjectBuild{
		BaseEntity: _models.BaseEntity{ID: pb.GetId()},
		Name:         pb.GetName(),
		NumFloor:     uint(pb.GetNumFloor()),
		NumApartment: uint(pb.GetNumApartment()),
		Note:         pb.GetNote(),
		ProjectID:    pb.GetProjectId(),
	}
}

func ProjectBuildItemsToPbList(builds []entity.ProjectBuild) []*bdspropb.ItemResponse {
	out := make([]*bdspropb.ItemResponse, len(builds))
	for i := range builds {
		out[i] = &bdspropb.ItemResponse{Id: builds[i].ID, Name: builds[i].Name}
	}
	return out
}

func ProjectItemsToPbList(items []entity.ProjectItem) []*bdspropb.ItemResponse {
	out := make([]*bdspropb.ItemResponse, len(items))
	for i := range items {
		out[i] = &bdspropb.ItemResponse{Id: items[i].ID, Name: items[i].Name}
	}
	return out
}

func ProjectToPb(project *entity.Project) *bdspropb.Project {
	if project == nil {
		return nil
	}
	active := project.DeletedAt == nil
	return &bdspropb.Project{
		Id:          project.ID,
		Name:        project.Name,
		DeveloperId: project.DeveloperID,
		Description: project.Description,
		CreatedAt:   _utils.FormatTimeToString(project.CreatedAt),
		UpdatedAt:   _utils.FormatTimeToString(project.UpdatedAt),
		Active:      active,
	}
}

func ProjectsToPbList(projects []entity.Project) []*bdspropb.Project {
	out := make([]*bdspropb.Project, len(projects))
	for i := range projects {
		out[i] = ProjectToPb(&projects[i])
	}
	return out
}

func ProjectFromPb(pb *bdspropb.Project) *entity.Project {
	if pb == nil {
		return nil
	}
	return &entity.Project{
		BaseEntity:  _models.BaseEntity{ID: pb.GetId()},
		Name:        pb.GetName(),
		DeveloperID: pb.GetDeveloperId(),
		Description: pb.GetDescription(),
	}
}

func apartmentsToPb(items []entity.ApartmentItem) []*bdspropb.Apartment {
	if len(items) == 0 {
		return nil
	}
	out := make([]*bdspropb.Apartment, len(items))
	for i := range items {
		a := &items[i]
		pb := &bdspropb.Apartment{
			Id:              a.ID,
			Name:            a.Name,
			Ordinal:         int32(a.Ordinal),
			ApartmentAttrId: a.ApartmentAttrID,
			Status:          int32(a.Status),
			Archived:        int32(a.Archived),
			BuildId:         a.BuildID,
		}
		if a.Floor != nil {
			pb.Floor = int32(*a.Floor)
		}
		out[i] = pb
	}
	return out
}

func apartmentAttrsToPb(items []entity.ApartmentAttrItem) []*bdspropb.ApartmentAttribute {
	if len(items) == 0 {
		return nil
	}
	out := make([]*bdspropb.ApartmentAttribute, len(items))
	for i := range items {
		a := &items[i]
		pb := &bdspropb.ApartmentAttribute{
			Id:           a.ID,
			Furniture:    a.Furniture,
			BlueprintUrl: a.BlueprintUrl,
			Price:        a.Price,
		}
		if a.Area != nil {
			area := *a.Area
			pb.Area = &area
		}
		if a.NumBedroom != nil {
			n := uint32(*a.NumBedroom)
			pb.NumBedroom = &n
		}
		if a.NumBathroom != nil {
			n := uint32(*a.NumBathroom)
			pb.NumBathroom = &n
		}
		out[i] = pb
	}
	return out
}
