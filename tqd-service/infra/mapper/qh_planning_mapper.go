package mapper

import (
	_utils "common/utils"
	tqdpb "pb/types/tqd"
	qh_domain "tqd/internal/domain/qh"
)

type QHPlanningMapper struct{}

func NewQHPlanningMapper() *QHPlanningMapper {
	return &QHPlanningMapper{}
}

func (m *QHPlanningMapper) ToProjectProto(project *qh_domain.QHPlanningProject) *tqdpb.PlanningProject {
	if project == nil {
		return nil
	}

	jurisdictionId := ""
	if project.JurisdictionID != nil {
		jurisdictionId = *project.JurisdictionID
	}

	return &tqdpb.PlanningProject{
		Id:             project.ID,
		Code:           project.Code,
		Name:           project.Name,
		PlanningType:   project.PlanningType,
		PlanningLevel:  project.PlanningLevel,
		Authority:      project.Authority,
		JurisdictionId: jurisdictionId,
		Summary:        project.Summary,
		ApprovalDate:   _utils.FormatTimeToString(project.ApprovalDate),
		EffectiveDate:  _utils.FormatTimeToString(project.EffectiveDate),
		ExpiryDate:     _utils.FormatTimeToString(project.ExpiryDate),
		ValidityStatus: project.ValidityStatus,
		CurrentVersion: project.CurrentVersion,
		Metadata:       project.Metadata,
		CreatedAt:      _utils.FormatTimeToString(project.CreatedAt),
		UpdatedAt:      _utils.FormatTimeToString(project.UpdatedAt),
	}
}

func (m *QHPlanningMapper) ToProjectProtoList(projects []qh_domain.QHPlanningProject) []*tqdpb.PlanningProject {
	result := make([]*tqdpb.PlanningProject, len(projects))
	for i, project := range projects {
		result[i] = m.ToProjectProto(&project)
	}
	return result
}

func (m *QHPlanningMapper) ToDocumentProto(document *qh_domain.QHPlanningDocument) *tqdpb.PlanningDocument {
	if document == nil {
		return nil
	}

	return &tqdpb.PlanningDocument{
		Id:                document.ID,
		PlanningProjectId: document.PlanningProjectID,
		DocumentType:      document.DocumentType,
		Code:              document.Code,
		Title:             document.Title,
		Description:       document.Description,
		Filepath:          document.Filepath,
		VersionNo:         document.VersionNo,
		ValidityStatus:    document.ValidityStatus,
		IssueDate:         _utils.FormatTimeToString(document.IssueDate),
		EffectiveDate:     _utils.FormatTimeToString(document.EffectiveDate),
		Metadata:          document.Metadata,
		CreatedAt:         _utils.FormatTimeToString(document.CreatedAt),
		UpdatedAt:         _utils.FormatTimeToString(document.UpdatedAt),
	}
}

func (m *QHPlanningMapper) ToDocumentProtoList(documents []qh_domain.QHPlanningDocument) []*tqdpb.PlanningDocument {
	result := make([]*tqdpb.PlanningDocument, len(documents))
	for i, document := range documents {
		result[i] = m.ToDocumentProto(&document)
	}
	return result
}

func (m *QHPlanningMapper) ToProjectDomain(proto *tqdpb.PlanningProject) *qh_domain.QHPlanningProject {
	if proto == nil {
		return nil
	}

	var jurisdictionId *string
	if proto.JurisdictionId != "" {
		jurisdictionId = &proto.JurisdictionId
	}

	metadata := proto.Metadata
	if metadata == "" {
		metadata = "{}"
	}

	return &qh_domain.QHPlanningProject{
		Code:           proto.Code,
		Name:           proto.Name,
		PlanningType:   proto.PlanningType,
		PlanningLevel:  proto.PlanningLevel,
		Authority:      proto.Authority,
		JurisdictionID: jurisdictionId,
		Summary:        proto.Summary,
		ApprovalDate:   _utils.ParseStringToTime(proto.ApprovalDate),
		EffectiveDate:  _utils.ParseStringToTime(proto.EffectiveDate),
		ExpiryDate:     _utils.ParseStringToTime(proto.ExpiryDate),
		ValidityStatus: proto.ValidityStatus,
		CurrentVersion: proto.CurrentVersion,
		Metadata:       metadata,
	}
}

func (m *QHPlanningMapper) ToDocumentDomain(proto *tqdpb.PlanningDocument) *qh_domain.QHPlanningDocument {
	if proto == nil {
		return nil
	}

	metadata := proto.Metadata
	if metadata == "" {
		metadata = "{}"
	}

	return &qh_domain.QHPlanningDocument{
		PlanningProjectID: proto.PlanningProjectId,
		DocumentType:      proto.DocumentType,
		Code:              proto.Code,
		Title:             proto.Title,
		Description:       proto.Description,
		Filepath:          proto.Filepath,
		VersionNo:         proto.VersionNo,
		ValidityStatus:    proto.ValidityStatus,
		IssueDate:         _utils.ParseStringToTime(proto.IssueDate),
		EffectiveDate:     _utils.ParseStringToTime(proto.EffectiveDate),
		Metadata:          metadata,
	}
}
