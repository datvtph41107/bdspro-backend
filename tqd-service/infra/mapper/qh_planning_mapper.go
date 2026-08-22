package mapper

import (
	"encoding/json"
	"strings"
	"time"

	_utils "common/utils"
	tqdpb "pb/types/tqd"
	qh_domain "tqd/internal/domain/qh"
	qh_dto "tqd/internal/dto/qh"
	"tqd/internal/enums"

	"github.com/lib/pq"
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
		ResearchScope:  project.ResearchScope,
		Indicators:     project.Indicators,

		SourceFolderName:  project.SourceFolderName,
		SourceFolderPath:  project.SourceFolderPath,
		ProcessStatus:     uint32(project.ProcessStatus),
		ProcessStatusName: enums.GetPlanningProcessStatusLabel(uint32(project.ProcessStatus)),
		DocumentCount:     project.DocumentCount,
	}
	clientView := project.ClientView
	result.ClientView = &clientView
	if project.SeoID != nil && *project.SeoID > 0 {
		seoID := *project.SeoID
		result.SeoId = &seoID
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

func (m *QHPlanningMapper) ToLayerProtoList(layers []qh_domain.QHLayer) []*tqdpb.LayerResponse {
	result := make([]*tqdpb.LayerResponse, 0, len(layers))
	for i := range layers {
		proto, err := ToProtoLayer(&layers[i])
		if err != nil || proto == nil {
			continue
		}
		proto.Docs = m.toLayerAttachedDocs(layers[i].Docs)
		result = append(result, proto)
	}
	return result
}

func (m *QHPlanningMapper) toLayerAttachedDocs(docs []qh_domain.QHPlanningDocument) []*tqdpb.LayerAttachedDoc {
	result := make([]*tqdpb.LayerAttachedDoc, 0, len(docs))
	for _, doc := range docs {
		result = append(result, &tqdpb.LayerAttachedDoc{
			Id:           doc.ID,
			Title:        doc.Title,
			Filepath:     doc.Filepath,
			Thumbnail:    doc.Thumbnail,
			Ext:          doc.Ext,
			RelativePath: doc.RelativePath,
			Code:         doc.Code,
			DocumentType: enums.GetPlanningDocumentTypeLabel(uint32(doc.DocumentType)),
		})
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
		IssuedDate:        _utils.FormatTimeToString(document.IssueDate),
		EffectiveDate:     _utils.FormatTimeToString(document.EffectiveDate),
		Metadata:          document.Metadata,
		Indicators:        document.Indicators,
		Authority:         document.Authority,
		Sector:            document.Sector,
		Field:             document.Field,
		Signer:            document.Signer,
		CreatedAt:         _utils.FormatTimeToString(document.CreatedAt),
		UpdatedAt:         _utils.FormatTimeToString(document.UpdatedAt),
		RelativePath:      document.RelativePath,
		ProcessStatus:     uint32(document.ProcessStatus),
		ProcessStatusName: enums.GetPlanningProcessStatusLabel(uint32(document.ProcessStatus)),
		ClassifyError:     document.ClassifyError,
		Ext:               document.Ext,
		SeoId:             document.SeoID,
		Status:            document.Status,
		StatusName:        enums.GetPlanningDocumentStatusLabel(document.Status),
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
		EventType:         event.EventType,
		EventDate:         formatPlanningTime(event.EventDate),
		Description:       event.Description,
		VerNo:             event.VerNo,
		LegalStatus:       uint32(event.LegalStatus),
		LegalStatusName:   enums.LegalStatusMap[event.LegalStatus],
		CreatedAt:         formatPlanningTime(event.CreatedAt),
		UpdatedAt:         formatPlanningTime(event.UpdatedAt),
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

func (m *QHPlanningMapper) ToDocumentHistoryProto(event *qh_domain.QHPlanningDocumentEvent) *tqdpb.PlanningDocumentHistory {
	if event == nil {
		return nil
	}
	return &tqdpb.PlanningDocumentHistory{
		Id:                 event.ID,
		PlanningDocumentId: event.PlanningDocumentID,
		EventType:          event.EventType,
		EventName:          event.EventName,
		Description:        event.Description,
		EventDate:          formatPlanningTime(event.EventDate),
		CreatedAt:          formatPlanningTime(event.CreatedAt),
		UpdatedAt:          formatPlanningTime(event.UpdatedAt),
	}
}

func (m *QHPlanningMapper) ToDocumentHistoryProtoList(events []qh_domain.QHPlanningDocumentEvent) []*tqdpb.PlanningDocumentHistory {
	result := make([]*tqdpb.PlanningDocumentHistory, len(events))
	for i, event := range events {
		result[i] = m.ToDocumentHistoryProto(&event)
	}
	return result
}

func (m *QHPlanningMapper) ToDocumentHistoryDomain(proto *tqdpb.PlanningDocumentHistory) *qh_domain.QHPlanningDocumentEvent {
	if proto == nil {
		return nil
	}
	return &qh_domain.QHPlanningDocumentEvent{
		PlanningDocumentID: proto.PlanningDocumentId,
		EventType:          proto.EventType,
		EventName:          proto.EventName,
		Description:        proto.Description,
		EventDate:          _utils.ParseStringToTime(proto.EventDate),
		Metadata:           []byte("{}"),
	}
}

func (m *QHPlanningMapper) ToEventDomain(proto *tqdpb.PlanningEvent) *qh_domain.QHPlanningEvent {
	if proto == nil {
		return nil
	}

	legalStatus := enums.LegalStatus(proto.LegalStatus)
	if legalStatus == 0 {
		legalStatus = enums.LegalStatusUnknown
	}

	return &qh_domain.QHPlanningEvent{
		PlanningProjectID: proto.PlanningProjectId,
		EventType:         proto.EventType,
		EventName:         proto.EventName,
		Description:       proto.Description,
		EventDate:         _utils.ParseStringToTime(proto.EventDate),
		VerNo:             proto.VerNo,
		LegalStatus:       legalStatus,
		DocIDs:            toEventDocIDArray(proto.DocIds),
		Metadata:          []byte("{}"),
	}
}

func toEventDocIDArray(ids []uint64) pq.Int64Array {
	if len(ids) == 0 {
		return nil
	}
	result := make(pq.Int64Array, len(ids))
	for i, id := range ids {
		result[i] = int64(id)
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
	if item.Relation != nil {
		result.Id = item.Relation.ID
		result.FromId = item.Relation.FromID
		result.ToId = item.Relation.ToID
		result.IsPeer = item.Relation.IsPeer
		result.RelationType = item.Relation.RelationType
		result.CreatedAt = _utils.FormatTimeToString(item.Relation.CreatedAt)
		result.UpdatedAt = _utils.FormatTimeToString(item.Relation.UpdatedAt)
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

func (m *QHPlanningMapper) ToRelationAdminProto(item *qh_domain.QHPlanningRelationItem) *tqdpb.PlanningRelation {
	if item == nil {
		return nil
	}
	return m.ToRelationProto(item, 0)
}

func (m *QHPlanningMapper) ToRelationAdminProtoList(items []qh_domain.QHPlanningRelationItem) []*tqdpb.PlanningRelation {
	result := make([]*tqdpb.PlanningRelation, len(items))
	for i, item := range items {
		result[i] = m.ToRelationAdminProto(&item)
	}
	return result
}

func (m *QHPlanningMapper) ToRelationDomain(proto *tqdpb.PlanningRelation) *qh_domain.QHPlanningRelation {
	if proto == nil {
		return nil
	}
	relationType := proto.RelationType
	if relationType == "" {
		if proto.IsPeer {
			relationType = "peer"
		} else {
			relationType = "related"
		}
	}
	return &qh_domain.QHPlanningRelation{
		FromID:                   proto.FromId,
		ToID:                     proto.ToId,
		IsPeer:                   proto.IsPeer,
		RelationType:             relationType,
		PlanningProjectID:        proto.FromId,
		RelatedPlanningProjectID: proto.ToId,
		Metadata:                 []byte("{}"),
	}
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
	indicators := proto.Indicators
	if indicators == "" {
		indicators = "[]"
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
		ResearchScope:  proto.ResearchScope,
		Indicators:     indicators,
		ClientView:     planningProjectClientView(proto),
	}
}

func planningProjectClientView(proto *tqdpb.PlanningProject) bool {
	if proto == nil || proto.ClientView == nil {
		return true
	}
	return *proto.ClientView
}

func (m *QHPlanningMapper) ToDocumentDomain(proto *tqdpb.PlanningDocument) *qh_domain.QHPlanningDocument {
	if proto == nil {
		return nil
	}

	metadata := proto.Metadata
	if metadata == "" {
		metadata = "{}"
	}
	indicators := proto.Indicators
	if indicators == "" {
		indicators = "[]"
	}

	status := proto.Status
	if status == 0 {
		status = enums.PlanningDocumentStatusNormal
	}

	issueDate := firstNonEmpty(
		proto.IssueDate,
		proto.IssuedDate,
		jsonMetaString(metadata, "issuedDate", "issueDate"),
	)
	effectiveDate := firstNonEmpty(
		proto.EffectiveDate,
		jsonMetaString(metadata, "effectiveDate"),
	)
	validityStatus := firstNonEmpty(
		proto.ValidityStatus,
		jsonMetaString(metadata, "validityStatus"),
	)

	return &qh_domain.QHPlanningDocument{
		PlanningProjectID: proto.PlanningProjectId,
		DocumentType:      enums.PlanningDocumentType(enums.GetPlanningDocumentTypeValue(proto.DocumentType)),
		Status:            status,
		Code:              proto.Code,
		Title:             proto.Title,
		Description:       proto.Description,
		Filepath:          proto.Filepath,
		Ext:               proto.Ext,
		Thumbnail:         proto.Thumbnail,
		VersionNo:         proto.VersionNo,
		ValidityStatus:    validityStatus,
		IssueDate:         _utils.ParseStringToTime(issueDate),
		EffectiveDate:     _utils.ParseStringToTime(effectiveDate),
		Metadata:          metadata,
		Indicators:        indicators,
		Authority:         proto.Authority,
		Sector:            proto.Sector,
		Field:             proto.Field,
		Signer:            proto.Signer,
		RelativePath:      proto.RelativePath,
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
		ClientView:     true,
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

func (m *QHPlanningMapper) ToProAIJobProto(job *qh_domain.ProAIJob) *tqdpb.ProAIJob {
	if job == nil {
		return nil
	}
	var projectID uint64
	if job.PlanningProjectID != nil {
		projectID = *job.PlanningProjectID
	}
	return &tqdpb.ProAIJob{
		Id:                job.ID,
		JobType:           job.JobType,
		JobTypeName:       enums.GetProAIJobTypeLabel(job.JobType),
		Name:              job.Name,
		SourceFolderName:  job.SourceFolderName,
		ProcessStatus:     uint32(job.ProcessStatus),
		ProcessStatusName: enums.GetPlanningProcessStatusLabel(uint32(job.ProcessStatus)),
		PlanningProjectId: projectID,
		ClassifyError:     job.ClassifyError,
		Metadata:          job.Metadata,
		CreatedAt:         _utils.FormatTimeToString(job.CreatedAt),
		UpdatedAt:         _utils.FormatTimeToString(job.UpdatedAt),
		CompletedAt:       _utils.FormatTimeToString(job.CompletedAt),
		ApprovedAt:        _utils.FormatTimeToString(job.ApprovedAt),
	}
}

func (m *QHPlanningMapper) ToProAIJobProtoList(jobs []qh_domain.ProAIJob) []*tqdpb.ProAIJob {
	result := make([]*tqdpb.ProAIJob, len(jobs))
	for i := range jobs {
		result[i] = m.ToProAIJobProto(&jobs[i])
	}
	return result
}

func formatPlanningTime(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	if formatted := _utils.FormatTimeToString(t); formatted != "" {
		return formatted
	}
	return t.UTC().Format(time.RFC3339)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			return trimmed
		}
	}
	return ""
}

func jsonMetaString(raw string, keys ...string) string {
	var payload map[string]any
	if err := json.Unmarshal([]byte(raw), &payload); err != nil {
		return ""
	}
	for _, key := range keys {
		if value, ok := payload[key].(string); ok {
			if trimmed := strings.TrimSpace(value); trimmed != "" {
				return trimmed
			}
		}
	}
	return ""
}
