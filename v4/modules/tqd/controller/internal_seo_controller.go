package controller

import (
	"context"
	"net/http"

	_errors "common/errors"
	"userv4/internal/httpapi"
	"userv4/modules/tqd/dto"
	"userv4/modules/tqd/services"
	qh_usecase "userv4/modules/tqd/services/qh"
)

// InternalSEOController owns /v2/tqd/internal/*/seo-* projection routes.
type InternalSEOController struct {
	location  *services.LocationUsecase
	layers    services.LayerUsecase
	parcels   services.ParcelUsecase
	projects  qh_usecase.IQHPlanningProjectUsecase
	documents qh_usecase.IQHPlanningDocumentUsecase
}

func NewInternalSEOController(
	location *services.LocationUsecase,
	layers services.LayerUsecase,
	parcels services.ParcelUsecase,
	projects qh_usecase.IQHPlanningProjectUsecase,
	documents qh_usecase.IQHPlanningDocumentUsecase,
) *InternalSEOController {
	return &InternalSEOController{location: location, layers: layers, parcels: parcels, projects: projects, documents: documents}
}

func (c *InternalSEOController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/tqd/internal/administrative-units/seo-source", c.adminUnitSeoSource)
	httpapi.PUT(mux, "/v2/tqd/internal/administrative-units/seo-id", c.adminUnitSeoID)

	httpapi.GET(mux, "/v2/tqd/internal/layers/{id}/seo-source", c.layerSeoSource)
	httpapi.PUT(mux, "/v2/tqd/internal/layers/{layerId}/seo-id", c.layerSeoID)

	httpapi.GET(mux, "/v2/tqd/internal/parcels/{parcelId}/seo-source", c.parcelSeoSource)
	httpapi.PUT(mux, "/v2/tqd/internal/parcels/{parcelId}/seo-id", c.parcelSeoID)
	httpapi.POST(mux, "/v2/tqd/internal/parcels/seo-sources", c.parcelSeoSources)

	httpapi.GET(mux, "/v2/tqd/internal/planning-projects/{id}/seo-source", c.projectSeoSource)
	httpapi.PUT(mux, "/v2/tqd/internal/planning-projects/{planningProjectId}/seo-id", c.projectSeoID)

	httpapi.GET(mux, "/v2/tqd/internal/planning-documents/{id}/seo-source", c.documentSeoSource)
	httpapi.PUT(mux, "/v2/tqd/internal/planning-documents/{planningDocumentId}/seo-id", c.documentSeoID)
}

type seoSourceQuery struct {
	ID uint64 `form:"id"`
}

func (c *InternalSEOController) adminUnitSeoSource(ctx context.Context, req seoSourceQuery) (*httpapi.Result, error) {
	if req.ID == 0 {
		return nil, _errors.ReturnError(400, "id is required")
	}
	out, err := c.location.GetAdministrativeUnitSeoSource(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *InternalSEOController) adminUnitSeoID(ctx context.Context, req struct {
	ID    uint64 `json:"id"`
	SeoID uint64 `json:"seoId"`
}) (*httpapi.Result, error) {
	if err := c.location.UpdateAdministrativeUnitSeoID(ctx, req.ID, req.SeoID); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true})
}

func (c *InternalSEOController) layerSeoSource(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	out, err := c.layers.GetLayerSeoSource(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *InternalSEOController) layerSeoID(ctx context.Context, req dto.SeoIDBody) (*httpapi.Result, error) {
	if err := c.layers.UpdateLayerSeoID(ctx, req.ID, req.SeoID); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true})
}

func (c *InternalSEOController) parcelSeoSource(ctx context.Context, req dto.ParcelIDPath) (*httpapi.Result, error) {
	out, err := c.parcels.GetParcelSeoSource(ctx, req.ParcelID)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *InternalSEOController) parcelSeoID(ctx context.Context, req dto.SeoIDBody) (*httpapi.Result, error) {
	if err := c.parcels.UpdateParcelSeoID(ctx, req.ID, req.SeoID); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true})
}

func (c *InternalSEOController) parcelSeoSources(ctx context.Context, req dto.ListParcelSeoSourcesRequest) (*httpapi.Result, error) {
	out, err := c.parcels.ListParcelSeoSourcesForGenerate(ctx, req.Limit)
	if err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"data": out})
}

func (c *InternalSEOController) projectSeoSource(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	out, err := c.projects.GetPlanningProjectSeoSource(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *InternalSEOController) projectSeoID(ctx context.Context, req dto.SeoIDBody) (*httpapi.Result, error) {
	if err := c.projects.UpdatePlanningProjectSeoID(ctx, req.ID, req.SeoID); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true})
}

func (c *InternalSEOController) documentSeoSource(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	out, err := c.documents.GetPlanningDocumentSeoSource(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *InternalSEOController) documentSeoID(ctx context.Context, req dto.SeoIDBody) (*httpapi.Result, error) {
	if err := c.documents.UpdatePlanningDocumentSeoID(ctx, req.ID, req.SeoID); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true})
}
