package controller

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	_dto "common/domain/dto"
	_errors "common/errors"
	_utils "common/utils"
	"userv4/internal/httpapi"
	"userv4/modules/tqd/dto"
	qh_dto "userv4/modules/tqd/dto/qh"
	planningdomain "userv4/modules/tqd/entity/planningclient"
	"userv4/modules/tqd/mapper"
	planningapp "userv4/modules/tqd/services/planningclient/application"
	qh_usecase "userv4/modules/tqd/services/qh"
)

// PlanningClientController owns /v2/tqd/client planning-project and planning-document routes.
type PlanningClientController struct {
	planning *planningapp.Service
	projects qh_usecase.IQHPlanningProjectUsecase
	events   qh_usecase.IQHPlanningEventUsecase
	docs     qh_usecase.IQHPlanningDocumentUsecase
	qhMapper *mapper.QHPlanningMapper
}

func NewPlanningClientController(
	planning *planningapp.Service,
	projects qh_usecase.IQHPlanningProjectUsecase,
	events qh_usecase.IQHPlanningEventUsecase,
	docs qh_usecase.IQHPlanningDocumentUsecase,
) *PlanningClientController {
	return &PlanningClientController{
		planning: planning,
		projects: projects,
		events:   events,
		docs:     docs,
		qhMapper: mapper.NewQHPlanningMapper(),
	}
}

func (c *PlanningClientController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/tqd/client/planning-projects", c.listProjects)
	httpapi.GET(mux, "/v2/tqd/client/planning-projects/followed", c.listFollowed)
	httpapi.GET(mux, "/v2/tqd/client/planning-projects/relation-parcel", c.listByParcelRelation)
	httpapi.GET(mux, "/v2/tqd/client/planning-projects/{project_id}", c.getProject)
	httpapi.GET(mux, "/v2/tqd/client/planning-projects/{project_id}/documents", c.listDocuments)
	httpapi.GET(mux, "/v2/tqd/client/planning-projects/{project_id}/events", c.listEvents)
	httpapi.POST(mux, "/v2/tqd/client/planning-projects/{project_id}/follow", c.follow)
	httpapi.DELETE(mux, "/v2/tqd/client/planning-projects/{project_id}/follow", c.unfollow)
	httpapi.GET(mux, "/v2/tqd/client/planning-projects/{project_id}/follow-status", c.followStatus)
	httpapi.GET(mux, "/v2/tqd/client/planning-projects/{project_id}/layers", c.listLayers)
	httpapi.GET(mux, "/v2/tqd/client/planning-projects/{project_id}/timeline", c.timeline)
	httpapi.GET(mux, "/v2/tqd/client/planning-projects/{project_id}/relations", c.relations)
	httpapi.GET(mux, "/v2/tqd/client/planning-projects/{project_id}/count", c.count)
	httpapi.GET(mux, "/v2/tqd/client/planning-projects/{value}/projection", c.projection)

	// Static leaf first; two-segment paths share one pattern so Go ServeMux
	// does not reject relation-layer|view|{id}/histories as mutual conflicts.
	httpapi.GET(mux, "/v2/tqd/client/planning-documents/by-adm", c.documentsByAdm)
	httpapi.GET(mux, "/v2/tqd/client/planning-documents/{seg1}/{seg2}", c.dispatchPlanningDocumentsTwoSeg)

	httpapi.GET(mux, "/v2/tqd/client/nearby-areas", c.nearbyAreas)
}

func (c *PlanningClientController) requirePlanning() error {
	if c.planning == nil {
		return _errors.ReturnError(501, "planning client service not configured")
	}
	return nil
}

func requireProfileID(ctx context.Context) (uint64, error) {
	userID := _utils.GetProfileIdWithContext(ctx)
	if userID == 0 {
		return 0, _errors.ReturnError(401, "authentication is required")
	}
	return userID, nil
}

func pageSize(size int32, fallback int) int {
	if size > 0 {
		return int(size)
	}
	return fallback
}

func (c *PlanningClientController) listProjects(ctx context.Context, req dto.ClientListPlanningProjectsRequest) (*httpapi.Result, error) {
	if err := c.requirePlanning(); err != nil {
		return nil, err
	}
	search := strings.TrimSpace(req.Keyword)
	if search == "" {
		search = strings.TrimSpace(req.Search)
	}
	layerIDs := _utils.ParseUint64Array(req.LayerIDs)
	page := int(req.Page)
	size := pageSize(req.Size, 20)

	if len(layerIDs) > 0 {
		result, err := c.planning.ListProjectsByLayerRelation(ctx, layerIDs, page, size)
		if err != nil {
			return nil, err
		}
		return jsonOK(result)
	}
	if req.ParcelID > 0 {
		result, err := c.planning.ListProjectsByParcelRelation(ctx, req.ParcelID, page, size)
		if err != nil {
			return nil, err
		}
		return jsonOK(result)
	}

	var hasLayer *bool
	if req.HasLayer != nil {
		switch *req.HasLayer {
		case 0:
			v := false
			hasLayer = &v
		case 1:
			v := true
			hasLayer = &v
		default:
			return nil, _errors.ReturnError(400, "hasLayer must be 0 or 1")
		}
	}
	updatedFrom, err := parseOptionalPlanningTime(req.UpdatedFrom, false)
	if err != nil {
		return nil, err
	}
	updatedTo, err := parseOptionalPlanningTime(req.UpdatedTo, true)
	if err != nil {
		return nil, err
	}
	filter := planningdomain.ProjectListFilter{
		Page:           page,
		Size:           size,
		Search:         strings.ToLower(search),
		PlanningType:   strings.TrimSpace(req.PlanningType),
		PlanningLevel:  strings.TrimSpace(req.PlanningLevel),
		ValidityStatus: strings.TrimSpace(req.ValidityStatus),
		ProcessStatus:  strings.TrimSpace(req.ProcessStatus),
		JurisdictionID: req.JurisdictionID,
		Area:           strings.TrimSpace(req.Area),
		LegalStatus:    strings.TrimSpace(req.LegalStatus),
		HasLayer:       hasLayer,
		SortBy:         strings.TrimSpace(req.SortBy),
		SortOrder:      strings.TrimSpace(req.SortOrder),
		Decision:       strings.TrimSpace(req.Decision),
		UpdatedFrom:    updatedFrom,
		UpdatedTo:      updatedTo,
		TagIDs:         strings.TrimSpace(req.TagIDs),
		ProvinceIDs:    strings.TrimSpace(req.ProvinceIDs),
		AdmID:          req.AdmID,
	}
	result, err := c.planning.ListProjects(ctx, filter)
	if err != nil {
		return nil, err
	}
	return jsonOK(result)
}

func parseOptionalPlanningTime(raw string, endOfDay bool) (*time.Time, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, nil
	}
	if t, err := time.Parse(time.RFC3339, raw); err == nil {
		return &t, nil
	}
	t, err := time.Parse("2006-01-02", raw)
	if err != nil {
		return nil, _errors.ReturnError(400, "updatedFrom/updatedTo must be RFC3339 or YYYY-MM-DD")
	}
	if endOfDay {
		t = t.Add(24*time.Hour - time.Nanosecond)
	}
	return &t, nil
}

func (c *PlanningClientController) getProject(ctx context.Context, req dto.ClientPlanningProjectIDRequest) (*httpapi.Result, error) {
	if err := c.requirePlanning(); err != nil {
		return nil, err
	}
	if req.ProjectID == 0 {
		return nil, _errors.ReturnError(400, "project_id is required")
	}
	result, err := c.planning.GetProject(ctx, req.ProjectID)
	if err != nil {
		return nil, err
	}
	return jsonOK(result)
}

func (c *PlanningClientController) listDocuments(ctx context.Context, req dto.ClientPlanningProjectIDRequest) (*httpapi.Result, error) {
	if err := c.requirePlanning(); err != nil {
		return nil, err
	}
	if req.ProjectID == 0 {
		return nil, _errors.ReturnError(400, "project_id is required")
	}
	result, err := c.planning.ListDocuments(ctx, req.ProjectID, int(req.Page), pageSize(req.Size, 20), req.Keyword, req.DocType, req.Status)
	if err != nil {
		return nil, err
	}
	return jsonOK(result)
}

func (c *PlanningClientController) listEvents(ctx context.Context, req dto.ClientPlanningProjectIDRequest) (*httpapi.Result, error) {
	if err := c.requirePlanning(); err != nil {
		return nil, err
	}
	if req.ProjectID == 0 {
		return nil, _errors.ReturnError(400, "project_id is required")
	}
	result, err := c.planning.ListEvents(ctx, req.ProjectID, int(req.Page), pageSize(req.Size, 20))
	if err != nil {
		return nil, err
	}
	return jsonOK(result)
}

func (c *PlanningClientController) follow(ctx context.Context, req dto.ClientFollowPlanningProjectRequest) (*httpapi.Result, error) {
	if err := c.requirePlanning(); err != nil {
		return nil, err
	}
	userID, err := requireProfileID(ctx)
	if err != nil {
		return nil, err
	}
	if req.ProjectID == 0 {
		return nil, _errors.ReturnError(400, "project_id is required")
	}
	result, err := c.planning.Follow(ctx, userID, req.ProjectID, strings.TrimSpace(req.Note))
	if err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"followed": true, "data": result})
}

func (c *PlanningClientController) unfollow(ctx context.Context, req dto.ClientPlanningProjectIDRequest) (*httpapi.Result, error) {
	if err := c.requirePlanning(); err != nil {
		return nil, err
	}
	userID, err := requireProfileID(ctx)
	if err != nil {
		return nil, err
	}
	if req.ProjectID == 0 {
		return nil, _errors.ReturnError(400, "project_id is required")
	}
	if err := c.planning.Unfollow(ctx, userID, req.ProjectID); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"followed": false, "planningProjectId": req.ProjectID})
}

func (c *PlanningClientController) followStatus(ctx context.Context, req dto.ClientPlanningProjectIDRequest) (*httpapi.Result, error) {
	if err := c.requirePlanning(); err != nil {
		return nil, err
	}
	userID, err := requireProfileID(ctx)
	if err != nil {
		return nil, err
	}
	if req.ProjectID == 0 {
		return nil, _errors.ReturnError(400, "project_id is required")
	}
	result, err := c.planning.GetFollow(ctx, userID, req.ProjectID)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return jsonOK(map[string]any{"followed": false, "planningProjectId": req.ProjectID})
		}
		return nil, err
	}
	return jsonOK(map[string]any{"followed": result != nil, "data": result})
}

func (c *PlanningClientController) listFollowed(ctx context.Context, req dto.ClientListFollowedPlanningProjectsRequest) (*httpapi.Result, error) {
	if err := c.requirePlanning(); err != nil {
		return nil, err
	}
	userID, err := requireProfileID(ctx)
	if err != nil {
		return nil, err
	}
	result, err := c.planning.ListFollowed(ctx, userID, int(req.Page), pageSize(req.Size, 20), req.Search)
	if err != nil {
		return nil, err
	}
	return jsonOK(result)
}

func (c *PlanningClientController) listByParcelRelation(ctx context.Context, req dto.ClientListParcelPlanningRelationRequest) (*httpapi.Result, error) {
	if err := c.requirePlanning(); err != nil {
		return nil, err
	}
	if req.ParcelID == 0 {
		return nil, _errors.ReturnError(400, "parcelId is required")
	}
	result, err := c.planning.ListProjectsByParcelRelation(ctx, req.ParcelID, int(req.Page), pageSize(req.Size, 20))
	if err != nil {
		return nil, err
	}
	return jsonOK(result)
}

func (c *PlanningClientController) listLayers(ctx context.Context, req dto.ClientPlanningProjectIDRequest) (*httpapi.Result, error) {
	if err := c.requirePlanning(); err != nil {
		return nil, err
	}
	if req.ProjectID == 0 {
		return nil, _errors.ReturnError(400, "project_id is required")
	}
	result, err := c.planning.ListProjectLayers(ctx, req.ProjectID)
	if err != nil {
		return nil, err
	}
	return jsonOK(result)
}

func (c *PlanningClientController) projection(ctx context.Context, req dto.ClientPlanningProjectProjectionRequest) (*httpapi.Result, error) {
	if err := c.requirePlanning(); err != nil {
		return nil, err
	}
	if strings.TrimSpace(req.Value) == "" {
		return nil, _errors.ReturnError(400, "value is required")
	}
	result, err := c.planning.GetPlanningProjectProjection(ctx, req.Value)
	if err != nil {
		return nil, err
	}
	return jsonOK(result)
}

func (c *PlanningClientController) dispatchPlanningDocumentsTwoSeg(ctx context.Context, req dto.ClientPlanningDocumentsTwoSegRequest) (*httpapi.Result, error) {
	switch {
	case req.Seg1 == "relation-layer":
		parcelID, err := parsePathUint64(req.Seg2, "parcel_id")
		if err != nil {
			return nil, err
		}
		return c.listDocsByParcelRelation(ctx, dto.ClientLayerPlanDocRelationRequest{
			ParcelID: parcelID,
			Page:     req.Page,
			Size:     req.Size,
			DocType:  req.DocType,
		})
	case req.Seg1 == "view":
		documentID, err := parsePathUint64(req.Seg2, "document_id")
		if err != nil {
			return nil, err
		}
		return c.getDocument(ctx, dto.ClientPlanningDocumentIDRequest{DocumentID: documentID})
	case req.Seg2 == "histories":
		docID, err := parsePathUint64(req.Seg1, "planningDocumentId")
		if err != nil {
			return nil, err
		}
		return c.listDocumentHistories(ctx, dto.ClientPlanningDocumentHistoryRequest{
			PlanningDocumentID: docID,
			Page:               req.Page,
			Size:               uint32(pageSize(req.Size, 20)),
		})
	default:
		return nil, _errors.ReturnError(404, "planning document route not found")
	}
}

func parsePathUint64(raw, name string) (uint64, error) {
	raw = strings.TrimSpace(raw)
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || id == 0 {
		return 0, _errors.ReturnError(400, name+" is required")
	}
	return id, nil
}

func (c *PlanningClientController) listDocsByParcelRelation(ctx context.Context, req dto.ClientLayerPlanDocRelationRequest) (*httpapi.Result, error) {
	if err := c.requirePlanning(); err != nil {
		return nil, err
	}
	if req.ParcelID == 0 {
		return nil, _errors.ReturnError(400, "parcel_id is required")
	}
	result, err := c.planning.ListDocumentsByParcelRelation(ctx, req.ParcelID, int(req.Page), pageSize(req.Size, 20), req.DocType)
	if err != nil {
		return nil, err
	}
	return jsonOK(result)
}

func (c *PlanningClientController) getDocument(ctx context.Context, req dto.ClientPlanningDocumentIDRequest) (*httpapi.Result, error) {
	if err := c.requirePlanning(); err != nil {
		return nil, err
	}
	if req.DocumentID == 0 {
		return nil, _errors.ReturnError(400, "document_id is required")
	}
	result, err := c.planning.GetDocument(ctx, req.DocumentID)
	if err != nil {
		return nil, err
	}
	return jsonOK(result)
}

func (c *PlanningClientController) timeline(ctx context.Context, req dto.ClientPlanningProjectIDRequest) (*httpapi.Result, error) {
	if c.events == nil {
		return nil, _errors.ReturnError(501, "planning event usecase not configured")
	}
	if req.ProjectID == 0 {
		return nil, _errors.ReturnError(400, "project_id is required")
	}
	listReq := &qh_dto.ListPlanningEventsRequest{
		PlanningProjectID: req.ProjectID,
		Pagable:           _dto.Pagable{Page: req.Page, Size: uint32(pageSize(req.Size, 20))},
	}
	rows, total, err := c.events.GetList(ctx, listReq)
	if err != nil {
		return nil, err
	}
	return protoListResult(c.qhMapper.ToEventProtoList(rows), total, listReq.GetPage(), listReq.GetSize())
}

func (c *PlanningClientController) relations(ctx context.Context, req dto.ClientPlanningProjectIDRequest) (*httpapi.Result, error) {
	if c.events == nil {
		return nil, _errors.ReturnError(501, "planning event usecase not configured")
	}
	if req.ProjectID == 0 {
		return nil, _errors.ReturnError(400, "project_id is required")
	}
	listReq := &qh_dto.ListPlanningEventsRequest{
		PlanningProjectID: req.ProjectID,
		Pagable:           _dto.Pagable{Page: req.Page, Size: uint32(pageSize(req.Size, 20))},
	}
	rows, _, err := c.events.GetRelations(ctx, listReq)
	if err != nil {
		return nil, err
	}
	items := c.qhMapper.ToRelationProtoList(rows, req.ProjectID)
	return protoListResult(items, int64(len(items)), listReq.GetPage(), listReq.GetSize())
}

func (c *PlanningClientController) count(ctx context.Context, req dto.ClientPlanningProjectIDRequest) (*httpapi.Result, error) {
	if c.projects == nil {
		return nil, _errors.ReturnError(501, "planning project usecase not configured")
	}
	if req.ProjectID == 0 {
		return nil, _errors.ReturnError(400, "project_id is required")
	}
	result, err := c.projects.Count(ctx, req.ProjectID)
	if err != nil {
		return nil, err
	}
	return jsonOK(result)
}

func (c *PlanningClientController) listDocumentHistories(ctx context.Context, req dto.ClientPlanningDocumentHistoryRequest) (*httpapi.Result, error) {
	if c.docs == nil {
		return nil, _errors.ReturnError(501, "planning document usecase not configured")
	}
	if req.PlanningDocumentID == 0 {
		return nil, _errors.ReturnError(400, "planningDocumentId is required")
	}
	listReq := &qh_dto.ListPlanningDocumentHistoriesRequest{
		PlanningDocumentID: req.PlanningDocumentID,
		Pagable:            _dto.Pagable{Page: req.Page, Size: req.Size},
	}
	rows, total, err := c.docs.ListHistories(ctx, listReq)
	if err != nil {
		return nil, err
	}
	return protoListResult(c.qhMapper.ToDocumentHistoryProtoList(rows), total, listReq.GetPage(), listReq.GetSize())
}

func (c *PlanningClientController) nearbyAreas(ctx context.Context, req dto.ClientNearbyAreasRequest) (*httpapi.Result, error) {
	return nil, _errors.ReturnError(501, "not implemented")
}

func (c *PlanningClientController) documentsByAdm(ctx context.Context, req dto.ClientPlanningDocumentsByAdmRequest) (*httpapi.Result, error) {
	return nil, _errors.ReturnError(501, "not implemented")
}
