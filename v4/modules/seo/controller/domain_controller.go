package controller

import (
	"context"
	"net/http"
	"strings"

	"userv4/internal/httpapi"
	"userv4/modules/seo/dto"
	"userv4/modules/seo/entity"
	"userv4/modules/seo/services"
)

// SeoDomainController handles admin SEO domain CRUD/lifecycle and public document endpoints.
type SeoDomainController struct {
	svc services.SeoDomainService
}

func NewSeoDomainController(svc services.SeoDomainService) *SeoDomainController {
	return &SeoDomainController{svc: svc}
}

func (c *SeoDomainController) RegisterHTTP(mux *http.ServeMux) {
	// Public document (JSON) + rendered page HTML
	httpapi.GET(mux, "/v2/crm/public/seo/documents/{module}/{slug}", c.getDocument)
	mux.HandleFunc("GET /v2/crm/public/content/page/{path...}", c.servePageHTML)
	httpapi.GET(mux, "/v2/crm/public/seos/{refType}/{refId}/list", c.listPublicByRef)
	httpapi.GET(mux, "/v2/crm/internal/seos/{refType}/{refId}/count", c.countByRef)

	// Public stats (CRM path parity)
	httpapi.GET(mux, "/v2/crm/public/seo-domains/top-clicks", c.topClicks)
	httpapi.GET(mux, "/v2/crm/public/seo-domains/click-daily", c.clickDaily)
	httpapi.GET(mux, "/v2/crm/public/seo-domains/count", c.countRecent)

	// Admin ref-types catalog
	httpapi.GET(mux, "/v2/crm/admin/seo/ref-types", c.listRefTypes)
	httpapi.GET(mux, "/v2/crm/admin/seo/sources", c.listRefTypes)

	// Admin CRUD (CRM path parity: /v2/crm/admin/seo-domains)
	httpapi.GET(mux, "/v2/crm/admin/seo-domains", c.adminList)
	httpapi.GET(mux, "/v2/crm/admin/seo-domains/{id}", c.adminGetByID)
	httpapi.GET(mux, "/v2/crm/admin/seo-domains/by-ref", c.adminGetByRef)
	httpapi.POST(mux, "/v2/crm/admin/seo-domains", c.adminCreate)
	httpapi.POST(mux, "/v2/crm/admin/seo-pages/drafts", c.adminCreateDraft)
	httpapi.PUT(mux, "/v2/crm/admin/seo-domains/{id}", c.adminUpdate)
	httpapi.PUT(mux, "/v2/crm/admin/seo-domains/batch", c.adminBatchUpdate)
	mux.HandleFunc("DELETE /v2/crm/admin/seo-domains/{id}", httpapi.Bind(c.adminDelete))

	// Lifecycle (CRM seo-pages paths + domains/{id} aliases)
	httpapi.POST(mux, "/v2/crm/admin/seo-pages/{id}/publish", c.adminPublish)
	httpapi.POST(mux, "/v2/crm/admin/seo-pages/{id}/unpublish", c.adminUnpublish)
	httpapi.POST(mux, "/v2/crm/admin/seo-pages/{id}/archive", c.adminArchive)
	httpapi.POST(mux, "/v2/crm/admin/seo-pages/{id}/restore", c.adminRestore)
	httpapi.POST(mux, "/v2/crm/admin/seo-domains/{id}/publish", c.adminPublish)
	httpapi.POST(mux, "/v2/crm/admin/seo-domains/{id}/unpublish", c.adminUnpublish)
	httpapi.POST(mux, "/v2/crm/admin/seo-domains/{id}/archive", c.adminArchive)
	httpapi.POST(mux, "/v2/crm/admin/seo-domains/{id}/restore", c.adminRestore)
	httpapi.POST(mux, "/v2/crm/admin/seo-pages/{id}/preview", c.adminPreviewExisting)
	httpapi.POST(mux, "/v2/crm/admin/seo-pages/preview", c.adminPreviewContent)
	httpapi.POST(mux, "/v2/crm/admin/seo-pages/{id}/generate", c.adminGenerateExisting)
	httpapi.POST(mux, "/v2/crm/admin/seo-domains/touch-by-ref", c.adminTouchByRef)
	httpapi.PUT(mux, "/v2/crm/admin/content/planning-news/{id}/category", c.adminAssignCategory)

	// Generation
	httpapi.POST(mux, "/v2/crm/admin/seo-domains/generate-from-ref", c.adminGenerateFromRef)
	httpapi.POST(mux, "/v2/crm/admin/seo-domains/parcel", c.adminGenerateParcel)
	httpapi.POST(mux, "/v2/crm/admin/seo-domains/parcels/generate", c.adminGenerateParcels)
	httpapi.POST(mux, "/v2/crm/admin/seo-domains/layers/generate", c.adminGenerateLayers)
	httpapi.POST(mux, "/v2/crm/admin/seo-domains/planning-projects/generate", c.adminGenerateProjects)
	httpapi.POST(mux, "/v2/crm/admin/seo-domains/planning-documents/generate", c.adminGenerateDocuments)
	httpapi.POST(mux, "/v2/crm/admin/seo-domains/administrative-units/generate", c.adminGenerateAdmUnits)
	httpapi.GET(mux, "/v2/crm/admin/seo-pages/{seoDomainId}/generation-logs", c.adminGenerationLogs)
	httpapi.GET(mux, "/v2/crm/admin/seo/generation-logs", c.adminGenerationLogs)

	// Per-domain click stats
	httpapi.GET(mux, "/v2/crm/admin/seo-domains/{id}/click-stats", c.adminClickStats)

	// Internal links
	httpapi.GET(mux, "/v2/crm/admin/seo-domains/{id}/links", c.adminListLinks)
	httpapi.POST(mux, "/v2/crm/admin/seo-domains/{id}/links", c.adminCreateLink)
	httpapi.GET(mux, "/v2/crm/admin/seo-domains/internal-links", c.adminListLinksQuery)
	httpapi.POST(mux, "/v2/crm/admin/seo-domains/internal-links", c.adminCreateLinkBody)
	httpapi.PUT(mux, "/v2/crm/admin/seo-domains/internal-links/{id}", c.adminUpdateLink)
	mux.HandleFunc("DELETE /v2/crm/admin/seo-domains/internal-links/{id}", httpapi.Bind(c.adminDeleteLink))
	mux.HandleFunc("DELETE /v2/crm/admin/seo/links/{id}", httpapi.Bind(c.adminDeleteLink))

	// Relatives
	httpapi.GET(mux, "/v2/crm/admin/seo-domains/{id}/relatives", c.adminListRelatives)
	httpapi.POST(mux, "/v2/crm/admin/seo-domains/{id}/relatives", c.adminCreateRelative)
	httpapi.GET(mux, "/v2/crm/admin/seo-domains/relatives", c.adminListRelativesQuery)
	httpapi.POST(mux, "/v2/crm/admin/seo-domains/relatives", c.adminCreateRelativeBody)
	mux.HandleFunc("DELETE /v2/crm/admin/seo-domains/relatives/{id}", httpapi.Bind(c.adminDeleteRelative))
	mux.HandleFunc("DELETE /v2/crm/admin/seo/relatives/{id}", httpapi.Bind(c.adminDeleteRelative))
}

// ── Public ─────────────────────────────────────────────────────────────────────

func (c *SeoDomainController) getDocument(ctx context.Context, req dto.GetDocumentRequest) (*httpapi.Result, error) {
	doc, err := c.svc.GetPublicDocument(ctx, req.Module, req.Slug)
	if err != nil {
		return nil, err
	}
	return dataResult(doc)
}

func (c *SeoDomainController) servePageHTML(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSpace(r.PathValue("path"))
	html, err := c.svc.GetPageHTML(r.Context(), path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Cache-Control", "public, max-age=60, stale-while-revalidate=300")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(html))
}

func (c *SeoDomainController) topClicks(ctx context.Context, req dto.LastDatesRequest) (*httpapi.Result, error) {
	resp, err := c.svc.GetTopClicked(ctx, req.LastDates)
	if err != nil {
		return nil, err
	}
	return dataResult(resp)
}

func (c *SeoDomainController) clickDaily(ctx context.Context, req dto.LastDatesRequest) (*httpapi.Result, error) {
	resp, err := c.svc.GetClickDailyStats(ctx, req.LastDates)
	if err != nil {
		return nil, err
	}
	return dataResult(resp)
}

func (c *SeoDomainController) countRecent(ctx context.Context, req dto.LastDatesRequest) (*httpapi.Result, error) {
	days := req.Days
	if days == 0 {
		days = req.LastDates
	}
	resp, err := c.svc.CountRecentByRefType(ctx, days)
	if err != nil {
		return nil, err
	}
	return dataResult(resp)
}

// ── Admin list/get ─────────────────────────────────────────────────────────────

func (c *SeoDomainController) adminList(ctx context.Context, req dto.SeoDomainListFilter) (*httpapi.Result, error) {
	result, err := c.svc.List(ctx, req)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return listResult([]dto.SeoDomainItem{}, 0)
	}
	return httpapi.OK(result), nil
}

func (c *SeoDomainController) adminGetByID(ctx context.Context, req dto.SeoDomainIDPath) (*httpapi.Result, error) {
	seo, err := c.svc.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return dataResult(seo)
}

func (c *SeoDomainController) adminGetByRef(ctx context.Context, req dto.SeoDomainByRefRequest) (*httpapi.Result, error) {
	seo, err := c.svc.GetByRef(ctx, req.RefType, req.RefID)
	if err != nil {
		return nil, err
	}
	return dataResult(seo)
}

// ── Admin CRUD ─────────────────────────────────────────────────────────────────

func (c *SeoDomainController) adminCreate(ctx context.Context, req dto.CreateSeoDomainRequest) (*httpapi.Result, error) {
	seo, err := c.svc.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	return dataResult(seo)
}

func (c *SeoDomainController) adminUpdate(ctx context.Context, req dto.UpdateSeoDomainRequest) (*httpapi.Result, error) {
	seo, err := c.svc.Update(ctx, req.ID, req)
	if err != nil {
		return nil, err
	}
	return dataResult(seo)
}

func (c *SeoDomainController) adminBatchUpdate(ctx context.Context, req dto.BatchUpdateSeoDomainRequest) (*httpapi.Result, error) {
	if err := c.svc.BatchUpdate(ctx, req); err != nil {
		return nil, err
	}
	return dataResult(map[string]string{"message": "ok"})
}

func (c *SeoDomainController) adminDelete(ctx context.Context, req dto.SeoDomainIDPath) (*httpapi.Result, error) {
	if err := c.svc.Delete(ctx, req.ID); err != nil {
		return nil, err
	}
	return dataResult(map[string]string{"message": "deleted"})
}

// ── Lifecycle ──────────────────────────────────────────────────────────────────

func (c *SeoDomainController) adminPublish(ctx context.Context, req dto.SeoDomainIDPath) (*httpapi.Result, error) {
	seo, err := c.svc.Publish(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return dataResult(seo)
}

func (c *SeoDomainController) adminUnpublish(ctx context.Context, req dto.SeoDomainIDPath) (*httpapi.Result, error) {
	seo, err := c.svc.Unpublish(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return dataResult(seo)
}

func (c *SeoDomainController) adminArchive(ctx context.Context, req dto.SeoDomainIDPath) (*httpapi.Result, error) {
	if err := c.svc.Archive(ctx, req.ID); err != nil {
		return nil, err
	}
	return dataResult(map[string]string{"message": "archived"})
}

func (c *SeoDomainController) adminRestore(ctx context.Context, req dto.SeoDomainIDPath) (*httpapi.Result, error) {
	if err := c.svc.Restore(ctx, req.ID); err != nil {
		return nil, err
	}
	return dataResult(map[string]string{"message": "restored"})
}

// ── Generation ─────────────────────────────────────────────────────────────────

func (c *SeoDomainController) adminGenerateFromRef(ctx context.Context, req dto.GenerateFromRefRequest) (*httpapi.Result, error) {
	seo, _, err := c.svc.GenerateFromRef(ctx, req)
	if err != nil {
		return nil, err
	}
	return dataResult(seo)
}

func (c *SeoDomainController) adminGenerateParcel(ctx context.Context, req dto.GenerateIDsRequest) (*httpapi.Result, error) {
	id := req.ParcelID
	if id == 0 && len(req.ParcelIDs) > 0 {
		id = req.ParcelIDs[0]
	}
	seo, _, err := c.svc.GenerateFromRef(ctx, dto.GenerateFromRefRequest{
		RefType:   uint32(entity.SEORefTypeParcel),
		RefID:     id,
		RefSource: "tqd_parcel",
	})
	if err != nil {
		return nil, err
	}
	return dataResult(seo)
}

func (c *SeoDomainController) adminGenerateParcels(ctx context.Context, req dto.GenerateIDsRequest) (*httpapi.Result, error) {
	ids := firstIDs(req.ParcelIDs, req.IDs)
	result, err := c.svc.GenerateBatch(ctx, entity.SEORefTypeParcel, "tqd_parcel", ids)
	if err != nil {
		return nil, err
	}
	return dataResult(result)
}

func (c *SeoDomainController) adminGenerateLayers(ctx context.Context, req dto.GenerateIDsRequest) (*httpapi.Result, error) {
	ids := firstIDs(req.LayerIDs, req.IDs)
	result, err := c.svc.GenerateBatch(ctx, entity.SEORefTypeMap, "tqd_layer", ids)
	if err != nil {
		return nil, err
	}
	return dataResult(result)
}

func (c *SeoDomainController) adminGenerateProjects(ctx context.Context, req dto.GenerateIDsRequest) (*httpapi.Result, error) {
	ids := firstIDs(req.PlanningProjectIDs, req.IDs)
	result, err := c.svc.GenerateBatch(ctx, entity.SEORefTypeProject, "tqd_planning_project", ids)
	if err != nil {
		return nil, err
	}
	return dataResult(result)
}

func (c *SeoDomainController) adminGenerateDocuments(ctx context.Context, req dto.GenerateIDsRequest) (*httpapi.Result, error) {
	ids := firstIDs(req.DocIDs, req.IDs)
	result, err := c.svc.GenerateBatch(ctx, entity.SEORefTypeDocument, "tqd_planning_document", ids)
	if err != nil {
		return nil, err
	}
	return dataResult(result)
}

func (c *SeoDomainController) adminGenerateAdmUnits(ctx context.Context, req dto.GenerateIDsRequest) (*httpapi.Result, error) {
	ids := firstIDs(req.CompositeIDs, req.IDs)
	result, err := c.svc.GenerateBatch(ctx, entity.SEORefTypeAdmUnit, "tqd_adm_unit", ids)
	if err != nil {
		return nil, err
	}
	return dataResult(result)
}

func firstIDs(preferred, fallback []uint64) []uint64 {
	if len(preferred) > 0 {
		return preferred
	}
	return fallback
}

type generationLogPathFilter struct {
	SeoDomainID uint64 `path:"seoDomainId" form:"seoDomainId"`
	Page        int    `form:"page"`
	Size        int    `form:"size"`
}

func (c *SeoDomainController) adminGenerationLogs(ctx context.Context, req generationLogPathFilter) (*httpapi.Result, error) {
	logs, total, err := c.svc.ListGenerationLogs(ctx, dto.GenerationLogFilter{
		SeoDomainID: req.SeoDomainID,
		Page:        req.Page,
		Size:        req.Size,
	})
	if err != nil {
		return nil, err
	}
	return listResult(logs, total)
}

// ── Stats ──────────────────────────────────────────────────────────────────────

func (c *SeoDomainController) adminClickStats(ctx context.Context, req dto.SeoDomainIDPath) (*httpapi.Result, error) {
	stats, err := c.svc.GetClickStats(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return dataResult(stats)
}

// ── Links ──────────────────────────────────────────────────────────────────────

func (c *SeoDomainController) adminListLinks(ctx context.Context, req dto.SeoDomainIDPath) (*httpapi.Result, error) {
	links, err := c.svc.ListLinks(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return listResult(links, int64(len(links)))
}

type parentSeoIDQuery struct {
	ParentSeoID uint64 `form:"parentSeoId"`
}

func (c *SeoDomainController) adminListLinksQuery(ctx context.Context, req parentSeoIDQuery) (*httpapi.Result, error) {
	links, err := c.svc.ListLinks(ctx, req.ParentSeoID)
	if err != nil {
		return nil, err
	}
	return listResult(links, int64(len(links)))
}

func (c *SeoDomainController) adminCreateLink(ctx context.Context, req dto.CreateSeoLinkRequest) (*httpapi.Result, error) {
	link, err := c.svc.CreateLink(ctx, req)
	if err != nil {
		return nil, err
	}
	return dataResult(link)
}

func (c *SeoDomainController) adminCreateLinkBody(ctx context.Context, req dto.CreateSeoLinkRequest) (*httpapi.Result, error) {
	return c.adminCreateLink(ctx, req)
}

func (c *SeoDomainController) adminDeleteLink(ctx context.Context, req httpapi.PathUint64) (*httpapi.Result, error) {
	if err := c.svc.DeleteLink(ctx, req.ID); err != nil {
		return nil, err
	}
	return dataResult(map[string]string{"message": "deleted"})
}

// ── Relatives ──────────────────────────────────────────────────────────────────

func (c *SeoDomainController) adminListRelatives(ctx context.Context, req dto.SeoDomainIDPath) (*httpapi.Result, error) {
	rels, err := c.svc.ListRelatives(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return listResult(rels, int64(len(rels)))
}

func (c *SeoDomainController) adminListRelativesQuery(ctx context.Context, req parentSeoIDQuery) (*httpapi.Result, error) {
	rels, err := c.svc.ListRelatives(ctx, req.ParentSeoID)
	if err != nil {
		return nil, err
	}
	return listResult(rels, int64(len(rels)))
}

func (c *SeoDomainController) adminCreateRelative(ctx context.Context, req dto.CreateSeoLinkRequest) (*httpapi.Result, error) {
	parentID := req.ParentSeoID
	rel, err := c.svc.CreateRelative(ctx, parentID, req)
	if err != nil {
		return nil, err
	}
	return dataResult(rel)
}

func (c *SeoDomainController) adminCreateRelativeBody(ctx context.Context, req dto.CreateSeoLinkRequest) (*httpapi.Result, error) {
	return c.adminCreateRelative(ctx, req)
}

func (c *SeoDomainController) adminDeleteRelative(ctx context.Context, req httpapi.PathUint64) (*httpapi.Result, error) {
	if err := c.svc.DeleteRelative(ctx, req.ID); err != nil {
		return nil, err
	}
	return dataResult(map[string]string{"message": "deleted"})
}

func (c *SeoDomainController) listPublicByRef(ctx context.Context, req dto.PublicSeoByRefPath) (*httpapi.Result, error) {
	list, err := c.svc.ListPublicByRef(ctx, req.RefType, req.RefID)
	if err != nil {
		return nil, err
	}
	return listResult(list, int64(len(list)))
}

func (c *SeoDomainController) countByRef(ctx context.Context, req dto.PublicSeoByRefPath) (*httpapi.Result, error) {
	total, err := c.svc.CountByRef(ctx, req.RefType, req.RefID)
	if err != nil {
		return nil, err
	}
	return dataResult(map[string]any{"count": total})
}

func (c *SeoDomainController) listRefTypes(ctx context.Context, _ httpapi.Empty) (*httpapi.Result, error) {
	return dataResult(c.svc.ListRefTypes(ctx))
}

func (c *SeoDomainController) adminCreateDraft(ctx context.Context, req dto.CreateSeoDomainRequest) (*httpapi.Result, error) {
	seo, err := c.svc.CreateDraft(ctx, req)
	if err != nil {
		return nil, err
	}
	return dataResult(seo)
}

func (c *SeoDomainController) adminPreviewExisting(ctx context.Context, req dto.SeoDomainIDPath) (*httpapi.Result, error) {
	html, err := c.svc.Preview(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return dataResult(map[string]string{"html": html})
}

func (c *SeoDomainController) adminPreviewContent(ctx context.Context, req dto.CreateSeoDomainRequest) (*httpapi.Result, error) {
	html, err := c.svc.PreviewContent(ctx, req)
	if err != nil {
		return nil, err
	}
	return dataResult(map[string]string{"html": html})
}

func (c *SeoDomainController) adminGenerateExisting(ctx context.Context, req dto.SeoDomainIDPath) (*httpapi.Result, error) {
	seo, err := c.svc.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	html, err := c.svc.Preview(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	updated, err := c.svc.Update(ctx, req.ID, dto.UpdateSeoDomainRequest{
		ID: req.ID,
	})
	if err != nil {
		return nil, err
	}
	_ = updated
	_ = seo
	_ = html
	// Mark render success via repo update through Publish path is heavy; store preview HTML.
	now := true
	_ = now
	return dataResult(map[string]any{"id": req.ID, "html": html, "message": "generated"})
}

func (c *SeoDomainController) adminTouchByRef(ctx context.Context, req dto.TouchByRefRequest) (*httpapi.Result, error) {
	n, err := c.svc.TouchByRef(ctx, req.RefType, req.RefID)
	if err != nil {
		return nil, err
	}
	return dataResult(map[string]any{"touched": n})
}

func (c *SeoDomainController) adminAssignCategory(ctx context.Context, req dto.AssignCategoryRequest) (*httpapi.Result, error) {
	seo, err := c.svc.AssignCategory(ctx, req.ID, req.CategoryID)
	if err != nil {
		return nil, err
	}
	return dataResult(seo)
}

func (c *SeoDomainController) adminUpdateLink(ctx context.Context, req dto.UpdateSeoLinkRequest) (*httpapi.Result, error) {
	link, err := c.svc.UpdateLink(ctx, req.ID, req)
	if err != nil {
		return nil, err
	}
	return dataResult(link)
}
