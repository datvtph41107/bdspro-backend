package controller

import (
	"context"
	"net/http"

	_dto "common/domain/dto"
	_errors "common/errors"
	"userv4/internal/httpapi"
	"userv4/modules/tqd/dto"
	qh_dto "userv4/modules/tqd/dto/qh"
	qh "userv4/modules/tqd/entity/qh"
	qh_usecase "userv4/modules/tqd/services/qh"
)

// AdminPlanningController owns admin planning projects/documents/events/relations/layer-docs/ai-jobs.
type AdminPlanningController struct {
	projects  qh_usecase.IQHPlanningProjectUsecase
	documents qh_usecase.IQHPlanningDocumentUsecase
	events    qh_usecase.IQHPlanningEventUsecase
	aiJobs    qh_usecase.IProAIJobUsecase
}

func NewAdminPlanningController(
	projects qh_usecase.IQHPlanningProjectUsecase,
	documents qh_usecase.IQHPlanningDocumentUsecase,
	events qh_usecase.IQHPlanningEventUsecase,
	aiJobs qh_usecase.IProAIJobUsecase,
) *AdminPlanningController {
	return &AdminPlanningController{projects: projects, documents: documents, events: events, aiJobs: aiJobs}
}

func (c *AdminPlanningController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.POST(mux, "/v2/tqd/admin/planning-projects", c.createProject)
	httpapi.GET(mux, "/v2/tqd/admin/planning-projects", c.listProjects)
	httpapi.GET(mux, "/v2/tqd/admin/planning-projects/{id}", c.getProject)
	httpapi.PUT(mux, "/v2/tqd/admin/planning-projects/{id}", c.updateProject)
	httpapi.DELETE(mux, "/v2/tqd/admin/planning-projects/{id}", c.deleteProject)
	httpapi.POST(mux, "/v2/tqd/admin/planning-projects/from-folder", c.fromFolder)
	httpapi.POST(mux, "/v2/tqd/admin/planning-projects/gen-seo", c.genProjectSeo)
	httpapi.POST(mux, "/v2/tqd/admin/planning-projects/{id}/approve", c.approveProject)
	httpapi.POST(mux, "/v2/tqd/admin/planning-projects/{id}/retry-failed", c.retryFailed)
	httpapi.GET(mux, "/v2/tqd/admin/planning-projects/{id}/count", c.projectCount)
	httpapi.GET(mux, "/v2/tqd/admin/planning-projects/{planningId}/layers", c.projectLayers)
	httpapi.DELETE(mux, "/v2/tqd/admin/planning-projects/{planningId}/layers/{layerId}", c.unassignLayer)
	httpapi.POST(mux, "/v2/tqd/admin/planning-projects/{planningId}/provinces", c.assignProvinces)
	httpapi.POST(mux, "/v2/tqd/admin/planning-projects/{planningId}/provinces:remove", c.unassignProvinces)
	httpapi.POST(mux, "/v2/tqd/admin/planning-projects/{planningId}/wards", c.assignWards)
	httpapi.POST(mux, "/v2/tqd/admin/planning-projects/{planningId}/wards:remove", c.unassignWards)

	httpapi.POST(mux, "/v2/tqd/admin/planning-documents", c.createDocument)
	httpapi.GET(mux, "/v2/tqd/admin/planning-documents", c.listDocuments)
	httpapi.GET(mux, "/v2/tqd/admin/planning-documents/{id}", c.getDocument)
	httpapi.PUT(mux, "/v2/tqd/admin/planning-documents/{id}", c.updateDocument)
	httpapi.DELETE(mux, "/v2/tqd/admin/planning-documents/{id}", c.deleteDocument)
	httpapi.POST(mux, "/v2/tqd/admin/planning-documents/gen-seo", c.genDocumentSeo)
	httpapi.POST(mux, "/v2/tqd/admin/planning-documents/{id}/approve", c.approveDocument)
	httpapi.POST(mux, "/v2/tqd/admin/planning-documents/{id}/retry", c.retryDocument)
	httpapi.POST(mux, "/v2/tqd/admin/planning-documents/{planningDocumentId}/histories", c.createHistory)
	httpapi.DELETE(mux, "/v2/tqd/admin/planning-documents/{planningDocumentId}/histories/{id}", c.deleteHistory)

	httpapi.POST(mux, "/v2/tqd/admin/planning-events", c.createEvent)
	httpapi.GET(mux, "/v2/tqd/admin/planning-events", c.listEvents)
	httpapi.GET(mux, "/v2/tqd/admin/planning-events/{id}", c.getEvent)
	httpapi.PUT(mux, "/v2/tqd/admin/planning-events/{id}", c.updateEvent)
	httpapi.DELETE(mux, "/v2/tqd/admin/planning-events/{id}", c.deleteEvent)

	httpapi.POST(mux, "/v2/tqd/admin/planning-relations", c.createRelation)
	httpapi.GET(mux, "/v2/tqd/admin/planning-relations", c.listRelations)
	httpapi.GET(mux, "/v2/tqd/admin/planning-relations/{id}", c.getRelation)
	httpapi.PUT(mux, "/v2/tqd/admin/planning-relations/{id}", c.updateRelation)
	httpapi.DELETE(mux, "/v2/tqd/admin/planning-relations/{id}", c.deleteRelation)

	httpapi.POST(mux, "/v2/tqd/admin/layer-docs/assign", c.assignDocs)
	httpapi.POST(mux, "/v2/tqd/admin/layer-docs/unassign", c.unassignDocs)

	httpapi.GET(mux, "/v2/tqd/admin/ai-jobs", c.listAIJobs)
	httpapi.GET(mux, "/v2/tqd/admin/ai-jobs/{id}", c.getAIJob)
}

type pageSearch struct {
	Page   uint32 `form:"page"`
	Size   uint32 `form:"size"`
	Search string `form:"search"`
}

func (c *AdminPlanningController) createProject(ctx context.Context, req qh.QHPlanningProject) (*httpapi.Result, error) {
	out, err := c.projects.Create(ctx, &req)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *AdminPlanningController) listProjects(ctx context.Context, req pageSearch) (*httpapi.Result, error) {
	pagable := _dto.Pagable{Page: req.Page, Size: req.Size}
	pagable.Normalize()
	rows, total, err := c.projects.GetList(ctx, &qh_dto.ListPlanningProjectsRequest{Pagable: pagable, Search: req.Search})
	if err != nil {
		return nil, err
	}
	return jsonListResult(rows, total, pagable.GetPage(), pagable.GetSize())
}

func (c *AdminPlanningController) getProject(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	out, err := c.projects.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *AdminPlanningController) updateProject(ctx context.Context, req qh.QHPlanningProject) (*httpapi.Result, error) {
	if req.ID == 0 {
		return nil, _errors.ReturnError(400, "id is required")
	}
	out, err := c.projects.Update(ctx, req.ID, &req)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *AdminPlanningController) deleteProject(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	if _, err := c.projects.Delete(ctx, req.ID); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true})
}

func (c *AdminPlanningController) fromFolder(ctx context.Context, _ dto.EmptyJSON) (*httpapi.Result, error) {
	return nil, _errors.ReturnError(501, "planning from-folder requires file provider wiring")
}

func (c *AdminPlanningController) genProjectSeo(ctx context.Context, req dto.AdminGenSeoIDsRequest) (*httpapi.Result, error) {
	out, err := c.projects.AdminGenPlanningProjectSeo(ctx, req.IDs)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *AdminPlanningController) approveProject(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	out, err := c.projects.Approve(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *AdminPlanningController) retryFailed(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	out, err := c.projects.RetryFailedDocuments(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *AdminPlanningController) projectCount(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	out, err := c.projects.Count(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *AdminPlanningController) projectLayers(ctx context.Context, req struct {
	PlanningID uint64 `path:"planningId"`
	Page       uint32 `form:"page"`
	Size       uint32 `form:"size"`
}) (*httpapi.Result, error) {
	pagable := _dto.Pagable{Page: req.Page, Size: req.Size}
	pagable.Normalize()
	rows, total, err := c.projects.ListLayers(ctx, &qh_dto.ListPlanningProjectLayersRequest{Pagable: pagable, PlanningID: req.PlanningID})
	if err != nil {
		return nil, err
	}
	return jsonListResult(rows, total, pagable.GetPage(), pagable.GetSize())
}

func (c *AdminPlanningController) unassignLayer(ctx context.Context, req dto.PlanningLayerPath) (*httpapi.Result, error) {
	if err := c.documents.UnassignLayer(ctx, &qh_dto.UnassignPlanningProjectLayerRequest{
		PlanningID: req.PlanningID, LayerID: req.LayerID,
	}); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true})
}

type idsBody struct {
	IDs []uint64 `json:"ids"`
	ProvinceIDs []uint64 `json:"provinceIds"`
	WardIDs []uint64 `json:"wardIds"`
}

func (c *AdminPlanningController) assignProvinces(ctx context.Context, req struct {
	PlanningID  uint64   `path:"planningId" json:"-"`
	ProvinceIDs []uint64 `json:"provinceIds"`
	IDs         []uint64 `json:"ids"`
}) (*httpapi.Result, error) {
	ids := req.ProvinceIDs
	if len(ids) == 0 {
		ids = req.IDs
	}
	if err := c.projects.AssignProvinces(ctx, req.PlanningID, ids); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true})
}

func (c *AdminPlanningController) unassignProvinces(ctx context.Context, req struct {
	PlanningID  uint64   `path:"planningId" json:"-"`
	ProvinceIDs []uint64 `json:"provinceIds"`
	IDs         []uint64 `json:"ids"`
}) (*httpapi.Result, error) {
	ids := req.ProvinceIDs
	if len(ids) == 0 {
		ids = req.IDs
	}
	if err := c.projects.UnassignProvinces(ctx, req.PlanningID, ids); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true})
}

func (c *AdminPlanningController) assignWards(ctx context.Context, req struct {
	PlanningID uint64   `path:"planningId" json:"-"`
	WardIDs    []uint64 `json:"wardIds"`
	IDs        []uint64 `json:"ids"`
}) (*httpapi.Result, error) {
	ids := req.WardIDs
	if len(ids) == 0 {
		ids = req.IDs
	}
	if err := c.projects.AssignWards(ctx, req.PlanningID, ids); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true})
}

func (c *AdminPlanningController) unassignWards(ctx context.Context, req struct {
	PlanningID uint64   `path:"planningId" json:"-"`
	WardIDs    []uint64 `json:"wardIds"`
	IDs        []uint64 `json:"ids"`
}) (*httpapi.Result, error) {
	ids := req.WardIDs
	if len(ids) == 0 {
		ids = req.IDs
	}
	if err := c.projects.UnassignWards(ctx, req.PlanningID, ids); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true})
}

func (c *AdminPlanningController) createDocument(ctx context.Context, req qh.QHPlanningDocument) (*httpapi.Result, error) {
	out, err := c.documents.Create(ctx, &req)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *AdminPlanningController) listDocuments(ctx context.Context, req pageSearch) (*httpapi.Result, error) {
	pagable := _dto.Pagable{Page: req.Page, Size: req.Size}
	pagable.Normalize()
	rows, total, err := c.documents.GetList(ctx, &qh_dto.ListPlanningDocumentsRequest{Pagable: pagable, Keyword: req.Search})
	if err != nil {
		return nil, err
	}
	return jsonListResult(rows, total, pagable.GetPage(), pagable.GetSize())
}

func (c *AdminPlanningController) getDocument(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	out, err := c.documents.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *AdminPlanningController) updateDocument(ctx context.Context, req qh.QHPlanningDocument) (*httpapi.Result, error) {
	out, err := c.documents.Update(ctx, req.ID, &req)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *AdminPlanningController) deleteDocument(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	if _, err := c.documents.Delete(ctx, req.ID); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true})
}

func (c *AdminPlanningController) genDocumentSeo(ctx context.Context, req dto.AdminGenSeoIDsRequest) (*httpapi.Result, error) {
	out, err := c.documents.AdminGenPlanningDocumentSeo(ctx, req.IDs)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *AdminPlanningController) approveDocument(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	out, err := c.documents.Approve(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *AdminPlanningController) retryDocument(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	out, err := c.documents.Retry(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *AdminPlanningController) createHistory(ctx context.Context, req qh.QHPlanningDocumentEvent) (*httpapi.Result, error) {
	out, err := c.documents.CreateHistory(ctx, &req)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *AdminPlanningController) deleteHistory(ctx context.Context, req dto.PlanningDocHistoryPath) (*httpapi.Result, error) {
	if _, err := c.documents.DeleteHistory(ctx, req.PlanningDocumentID, req.ID); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true})
}

func (c *AdminPlanningController) createEvent(ctx context.Context, req qh.QHPlanningEvent) (*httpapi.Result, error) {
	out, err := c.events.Create(ctx, &req)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *AdminPlanningController) listEvents(ctx context.Context, req pageSearch) (*httpapi.Result, error) {
	pagable := _dto.Pagable{Page: req.Page, Size: req.Size}
	pagable.Normalize()
	rows, total, err := c.events.GetList(ctx, &qh_dto.ListPlanningEventsRequest{Pagable: pagable})
	if err != nil {
		return nil, err
	}
	return jsonListResult(rows, total, pagable.GetPage(), pagable.GetSize())
}

func (c *AdminPlanningController) getEvent(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	out, err := c.events.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *AdminPlanningController) updateEvent(ctx context.Context, req qh.QHPlanningEvent) (*httpapi.Result, error) {
	out, err := c.events.Update(ctx, req.ID, &req)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *AdminPlanningController) deleteEvent(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	if _, err := c.events.Delete(ctx, req.ID); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true})
}

func (c *AdminPlanningController) createRelation(ctx context.Context, req qh.QHPlanningRelation) (*httpapi.Result, error) {
	out, err := c.events.CreateRelation(ctx, &req)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *AdminPlanningController) listRelations(ctx context.Context, req pageSearch) (*httpapi.Result, error) {
	pagable := _dto.Pagable{Page: req.Page, Size: req.Size}
	pagable.Normalize()
	rows, total, err := c.events.GetRelations(ctx, &qh_dto.ListPlanningEventsRequest{Pagable: pagable})
	if err != nil {
		return nil, err
	}
	return jsonListResult(rows, total, pagable.GetPage(), pagable.GetSize())
}

func (c *AdminPlanningController) getRelation(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	out, err := c.events.GetRelationByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *AdminPlanningController) updateRelation(ctx context.Context, req qh.QHPlanningRelation) (*httpapi.Result, error) {
	out, err := c.events.UpdateRelation(ctx, req.ID, &req)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *AdminPlanningController) deleteRelation(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	if _, err := c.events.DeleteRelation(ctx, req.ID); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true})
}

func (c *AdminPlanningController) assignDocs(ctx context.Context, req qh_dto.AssignPlanningDocumentsToLayerRequest) (*httpapi.Result, error) {
	if err := c.documents.AssignToLayer(ctx, &req); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true})
}

func (c *AdminPlanningController) unassignDocs(ctx context.Context, req qh_dto.AssignPlanningDocumentsToLayerRequest) (*httpapi.Result, error) {
	if err := c.documents.UnassignFromLayer(ctx, &req); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true})
}

func (c *AdminPlanningController) listAIJobs(ctx context.Context, req pageSearch) (*httpapi.Result, error) {
	pagable := _dto.Pagable{Page: req.Page, Size: req.Size}
	pagable.Normalize()
	rows, total, err := c.aiJobs.GetList(ctx, &qh_dto.ListProAIJobsRequest{Pagable: pagable, Search: req.Search})
	if err != nil {
		return nil, err
	}
	return jsonListResult(rows, total, pagable.GetPage(), pagable.GetSize())
}

func (c *AdminPlanningController) getAIJob(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	out, err := c.aiJobs.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}
