package mapper

import (
	_utils "common/utils"
	tqdpb "pb/types/tqd"
	qh_domain "tqd/internal/domain/qh"
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

	result := &tqdpb.PlanningEvent{
		Id:                event.ID,
		PlanningProjectId: event.PlanningProjectID,
		DocIds:            toEventDocIDs(event.DocIDs),
		EventName:         event.EventName,
		EventDate:         _utils.FormatTimeToString(event.EventDate),
		Description:       event.Description,
		VerNo:             event.VerNo,
		LegalStatus:       uint32(event.LegalStatus),
		LegalStatusName:   enums.LegalStatusMap[event.LegalStatus],
		CreatedAt:         _utils.FormatTimeToString(event.CreatedAt),
		UpdatedAt:         _utils.FormatTimeToString(event.UpdatedAt),
	}

	if len(event.Docs) > 0 {
		result.Docs = m.ToDocumentProtoList(event.Docs)
	}

	return result
}

func toEventDocIDs(docIDs []int64) []uint64 {
	if len(docIDs) == 0 {
		return nil
	}
	result := make([]uint64, len(docIDs))
	for i, id := range docIDs {
		result[i] = uint64(id)
	}
	return result
}

func (m *QHPlanningMapper) ToEventProtoList(events []qh_domain.QHPlanningEvent) []*tqdpb.PlanningEvent {
	result := make([]*tqdpb.PlanningEvent, len(events))
	for i, event := range events {
		result[i] = m.ToEventProto(&event)
	}
	return result
}

func (m *QHPlanningMapper) ToRelationProto(item *qh_domain.QHPlanningRelationItem, planningProjectId uint64) *tqdpb.PlanningRelation {
	if item == nil {
		return nil
	}

	result := &tqdpb.PlanningRelation{
		Parent: m.ToProjectProto(item.Parent),
		Child:  m.ToProjectProto(item.Child),
		Peer:   m.ToProjectProto(item.Peer),
	}

	if item.Parent != nil && item.Parent.ID == planningProjectId {
		result.Parent = nil
	}

	if item.Child != nil && item.Child.ID == planningProjectId {
		result.Child = nil
	}

	if item.Peer != nil && item.Peer.ID == planningProjectId {
		result.Peer = nil
	}

	return result
}

func (m *QHPlanningMapper) ToRelationProtoList(items []qh_domain.QHPlanningRelationItem, planningProjectId uint64) []*tqdpb.PlanningRelation {
	result := make([]*tqdpb.PlanningRelation, len(items))
	for i, item := range items {
		result[i] = m.ToRelationProto(&item, planningProjectId)
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
