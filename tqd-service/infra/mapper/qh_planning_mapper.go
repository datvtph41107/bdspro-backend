package mapper

import (
	_utils "common/utils"
	tqdpb "pb/types/tqd"
	qh_domain "tqd/internal/domain/qh"
	qh_dto "tqd/internal/dto/qh"
	"tqd/internal/enums"
)

type QHPlanningMapper struct{}

func NewQHPlanningMapper() *QHPlanningMapper {
	return &QHPlanningMapper{}
}

func (m *QHPlanningMapper) ToProjectProto(project *qh_domain.QHPlanningProject) *tqdpb.PlanningProject {
	if project == nil {
		return nil
	}

	jurisdictionId := uint64(0)
	if project.JurisdictionID != nil {
		jurisdictionId = *project.JurisdictionID
	}

	result := &tqdpb.PlanningProject{
		Id:           project.ID,
		Code:         project.Code,
		Name:         project.Name,
		PlanningType: project.PlanningType,

		PlanningTypeName: enums.GetPlanningTypeLabel(project.PlanningType),

		PlanningLevel:     project.PlanningLevel,
		PlanningLevelName: enums.GetPlanningLevelLabel(project.PlanningLevel),

		LegalStatus:     uint32(project.LegalStatus),
		LegalStatusName: enums.LegalStatusMap[project.LegalStatus],

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

		SourceFolderName:  project.SourceFolderName,
		SourceFolderPath:  project.SourceFolderPath,
		ProcessStatus:     uint32(project.ProcessStatus),
		ProcessStatusName: enums.GetPlanningProcessStatusLabel(uint32(project.ProcessStatus)),
	}

	if project.Jurisdiction != nil {
		result.Jurisdiction = &tqdpb.Jurisdiction{
			Id:   project.Jurisdiction.ID,
			Name: project.Jurisdiction.Name,
		}
	}

	if len(project.Layers) > 0 {
		result.Layers = make([]*tqdpb.LayerResponse, len(project.Layers))
		for i, layer := range project.Layers {
			result.Layers[i] = &tqdpb.LayerResponse{
				Id:   layer.ID,
				Name: layer.Name,
			}
		}
	}

	return result
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
		DocumentType:      enums.GetPlanningDocumentTypeLabel(uint32(document.DocumentType)),
		Code:              document.Code,
		Title:             document.Title,
		Description:       document.Description,
		Filepath:          document.Filepath,
		Thumbnail:         document.Thumbnail,
		VersionNo:         document.VersionNo,
		ValidityStatus:    document.ValidityStatus,
		IssueDate:         _utils.FormatTimeToString(document.IssueDate),
		EffectiveDate:     _utils.FormatTimeToString(document.EffectiveDate),
		Metadata:          document.Metadata,
		CreatedAt:         _utils.FormatTimeToString(document.CreatedAt),
		UpdatedAt:         _utils.FormatTimeToString(document.UpdatedAt),

		RelativePath:      document.RelativePath,
		ProcessStatus:     uint32(document.ProcessStatus),
		ProcessStatusName: enums.GetPlanningProcessStatusLabel(uint32(document.ProcessStatus)),
		ClassifyError:     document.ClassifyError,
	}
}

func (m *QHPlanningMapper) ToDocumentProtoList(documents []qh_domain.QHPlanningDocument) []*tqdpb.PlanningDocument {
	result := make([]*tqdpb.PlanningDocument, len(documents))
	for i, document := range documents {
		result[i] = m.ToDocumentProto(&document)
	}
	return result
}

func (m *QHPlanningMapper) ToEventProto(event *qh_domain.QHPlanningEvent) *tqdpb.PlanningEvent {
	if event == nil {
		return nil
	}

	return &tqdpb.PlanningEvent{
		Id:                event.ID,
		PlanningProjectId: event.PlanningProjectID,
		EventName:         event.EventName,
		EventDate:         _utils.FormatTimeToString(event.EventDate),
		CreatedAt:         _utils.FormatTimeToString(event.CreatedAt),
		UpdatedAt:         _utils.FormatTimeToString(event.UpdatedAt),
	}
}

func (m *QHPlanningMapper) ToEventProtoList(events []qh_domain.QHPlanningEvent) []*tqdpb.PlanningEvent {
	result := make([]*tqdpb.PlanningEvent, len(events))
	for i, event := range events {
		result[i] = m.ToEventProto(&event)
	}
	return result
}

func (m *QHPlanningMapper) ToProjectDomain(proto *tqdpb.PlanningProject) *qh_domain.QHPlanningProject {
	if proto == nil {
		return nil
	}

	var jurisdictionId *uint64
	if proto.JurisdictionId != 0 {
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
		DocumentType:      enums.PlanningDocumentType(enums.GetPlanningDocumentTypeValue(proto.DocumentType)),
		Code:              proto.Code,
		Title:             proto.Title,
		Description:       proto.Description,
		Filepath:          proto.Filepath,
		Thumbnail:         proto.Thumbnail,
		VersionNo:         proto.VersionNo,
		ValidityStatus:    proto.ValidityStatus,
		IssueDate:         _utils.ParseStringToTime(proto.IssueDate),
		EffectiveDate:     _utils.ParseStringToTime(proto.EffectiveDate),
		Metadata:          metadata,
	}
}

// ToProjectDomainFromFolderRequest chuyển metadata đồ án trong CreatePlanningProjectFromFolderRequest thành domain.
// Không set SourceFolderName/ProcessStatus ở đây — usecase CreateFromFolder sẽ set trước khi Create.
func (m *QHPlanningMapper) ToProjectDomainFromFolderRequest(req *tqdpb.CreatePlanningProjectFromFolderRequest) *qh_domain.QHPlanningProject {
	if req == nil {
		return nil
	}

	var jurisdictionId *uint64
	if req.JurisdictionId != 0 {
		jurisdictionId = &req.JurisdictionId
	}

	return &qh_domain.QHPlanningProject{
		Code:           req.Code,
		Name:           req.Name,
		PlanningType:   req.PlanningType,
		PlanningLevel:  req.PlanningLevel,
		Authority:      req.Authority,
		JurisdictionID: jurisdictionId,
		Summary:        req.Summary,
		ValidityStatus: req.ValidityStatus,
		Metadata:       "{}",
	}
}

// ToFolderFileInputs chuyển danh sách FolderFileEntry (proto) thành input cho usecase.
func (m *QHPlanningMapper) ToFolderFileInputs(entries []*tqdpb.FolderFileEntry) []qh_dto.FolderFileInput {
	inputs := make([]qh_dto.FolderFileInput, 0, len(entries))
	for _, e := range entries {
		if e == nil || e.RelativePath == "" {
			continue
		}
		inputs = append(inputs, qh_dto.FolderFileInput{
			RelativePath: e.RelativePath,
			Content:      e.Content,
			ContentType:  e.ContentType,
		})
	}
	return inputs
}
