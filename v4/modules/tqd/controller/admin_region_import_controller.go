package controller

import (
	"context"
	"encoding/json"
	"net/http"

	_dto "common/domain/dto"
	_errors "common/errors"
	_utils "common/utils"
	"userv4/internal/httpapi"
	"userv4/modules/tqd/dto"
	"userv4/modules/tqd/enums"
	"userv4/modules/tqd/repository"
	"userv4/modules/tqd/services"
)

// AdminRegionImportController owns admin regions, region-extends, and import routes.
type AdminRegionImportController struct {
	regions services.RegionUsecase
	extends services.RegionExtendUsecase
	imports services.ImportUsecase
}

func NewAdminRegionImportController(
	regions services.RegionUsecase,
	extends services.RegionExtendUsecase,
	imports services.ImportUsecase,
) *AdminRegionImportController {
	return &AdminRegionImportController{regions: regions, extends: extends, imports: imports}
}

func (c *AdminRegionImportController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.POST(mux, "/v2/tqd/admin/regions", c.createRegion)
	httpapi.GET(mux, "/v2/tqd/admin/regions", c.listRegions)
	httpapi.GET(mux, "/v2/tqd/admin/regions/{id}", c.getRegion)
	httpapi.PUT(mux, "/v2/tqd/admin/regions/{id}", c.updateRegion)
	httpapi.DELETE(mux, "/v2/tqd/admin/regions/{id}", c.deleteRegion)
	httpapi.POST(mux, "/v2/tqd/admin/regions/sync-region", c.syncRegion)
	httpapi.PATCH(mux, "/v2/tqd/admin/regions/{regionId}/processing-status", c.processingStatus)
	httpapi.GET(mux, "/v2/tqd/admin/layers/{layer_id}/regions", c.listByLayer)

	httpapi.POST(mux, "/v2/tqd/admin/region-extends", c.createExtend)
	httpapi.GET(mux, "/v2/tqd/admin/region-extends/{id}", c.getExtend)
	httpapi.PUT(mux, "/v2/tqd/admin/region-extends/{id}", c.updateExtend)
	httpapi.DELETE(mux, "/v2/tqd/admin/region-extends/{id}", c.deleteExtend)
	httpapi.GET(mux, "/v2/tqd/admin/layers/{layer_id}/region-extends", c.listExtendsByLayer)

	httpapi.POST(mux, "/v2/tqd/admin/import/preview", c.importPreview)
	httpapi.GET(mux, "/v2/tqd/admin/import/batches/{batchId}", c.importBatch)
	httpapi.DELETE(mux, "/v2/tqd/admin/import/batches/{batchId}", c.importRollback)
	httpapi.POST(mux, "/v2/tqd/admin/import/errors/{errorId}/retry", c.importRetry)
	httpapi.GET(mux, "/v2/tqd/admin/layers/{layerId}/import/errors", c.importErrors)
}

type createRegionReq struct {
	LayerID  uint64          `json:"layerId"`
	Name     string          `json:"name"`
	LabelID  uint64          `json:"labelId"`
	Geometry json.RawMessage `json:"geometry"`
}

type updateRegionReq struct {
	ID          uint64  `path:"id" json:"-"`
	Name        *string `json:"name"`
	DisplayName *string `json:"displayName"`
	Description *string `json:"description"`
	Status      *uint32 `json:"status"`
	LabelID     *uint64 `json:"labelId"`
}

type listRegionReq struct {
	Page    uint32  `form:"page"`
	Size    uint32  `form:"size"`
	LayerID *uint64 `form:"layerId"`
	Search  *string `form:"search"`
}

func (c *AdminRegionImportController) createRegion(ctx context.Context, req createRegionReq) (*httpapi.Result, error) {
	out, err := c.regions.Create(ctx, req.LayerID, req.Name, []byte(req.Geometry), req.LabelID)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *AdminRegionImportController) getRegion(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	out, err := c.regions.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *AdminRegionImportController) updateRegion(ctx context.Context, req updateRegionReq) (*httpapi.Result, error) {
	var status *enums.RegionStatus
	if req.Status != nil {
		s := enums.RegionStatus(*req.Status)
		status = &s
	}
	out, err := c.regions.Update(ctx, services.RegionUpdatePatch{
		RegionID: req.ID, Name: req.Name, DisplayName: req.DisplayName, Description: req.Description,
		Status: status, LabelID: req.LabelID,
	})
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *AdminRegionImportController) deleteRegion(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	if err := c.regions.Delete(ctx, req.ID); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true})
}

func (c *AdminRegionImportController) listRegions(ctx context.Context, req listRegionReq) (*httpapi.Result, error) {
	pagable := _dto.Pagable{Page: req.Page, Size: req.Size}
	pagable.Normalize()
	rows, total, err := c.regions.List(ctx, &repository.RegionFilter{
		LayerID: req.LayerID, Search: req.Search, Pagable: &pagable,
	})
	if err != nil {
		return nil, err
	}
	return jsonListResult(rows, total, pagable.GetPage(), pagable.GetSize())
}

func (c *AdminRegionImportController) listByLayer(ctx context.Context, req struct {
	LayerID uint64 `path:"layer_id"`
	Page    uint32 `form:"page"`
	Size    uint32 `form:"size"`
}) (*httpapi.Result, error) {
	pagable := _dto.Pagable{Page: req.Page, Size: req.Size}
	pagable.Normalize()
	rows, total, err := c.regions.ListForAdminByLayer(ctx, req.LayerID, &repository.RegionFilter{Pagable: &pagable})
	if err != nil {
		return nil, err
	}
	return jsonListResult(rows, total, pagable.GetPage(), pagable.GetSize())
}

func (c *AdminRegionImportController) syncRegion(ctx context.Context, req dto.AdminSyncRegionRequest) (*httpapi.Result, error) {
	out, err := c.regions.SyncRegion(ctx, req.LayerID, req.LabelID)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *AdminRegionImportController) processingStatus(ctx context.Context, req dto.AdminRegionProcessingStatusRequest) (*httpapi.Result, error) {
	userID := _utils.GetProfileIdWithContext(ctx)
	if err := c.regions.UpdateProcessingStatus(ctx, req.RegionID, enums.ProcessingStatus(req.Status), userID, req.Notes); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true})
}

type createExtendReq struct {
	LayerID            uint64          `json:"layerId"`
	Name               string          `json:"name"`
	DisplayName        string          `json:"displayName"`
	Description        string          `json:"description"`
	LabelID            *uint64         `json:"labelId"`
	Geometry           json.RawMessage `json:"geometry"`
	LegalDoc           string          `json:"legalDoc"`
	PlanningName       string          `json:"planningName"`
	OriginalProperties json.RawMessage `json:"originalProperties"`
	SourceFile         string          `json:"sourceFile"`
	ImportBatchID      string          `json:"importBatchId"`
	Status             *uint32         `json:"status"`
	MergeSourceIDs     []uint64        `json:"mergeSourceIds"`
}

func (c *AdminRegionImportController) createExtend(ctx context.Context, req createExtendReq) (*httpapi.Result, error) {
	out, err := c.extends.Create(ctx, services.RegionExtendCreateInput{
		LayerID: req.LayerID, Name: req.Name, DisplayName: req.DisplayName, Description: req.Description,
		LabelID: req.LabelID, Geometry: []byte(req.Geometry), LegalDoc: req.LegalDoc, PlanningName: req.PlanningName,
		OriginalProperties: []byte(req.OriginalProperties), SourceFile: req.SourceFile, ImportBatchID: req.ImportBatchID,
		Status: req.Status, MergeSourceIDs: req.MergeSourceIDs,
	})
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *AdminRegionImportController) getExtend(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	out, err := c.extends.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *AdminRegionImportController) updateExtend(ctx context.Context, req struct {
	ID                 uint64          `path:"id" json:"-"`
	LayerID            *uint64         `json:"layerId"`
	Name               *string         `json:"name"`
	DisplayName        *string         `json:"displayName"`
	Description        *string         `json:"description"`
	LabelID            *uint64         `json:"labelId"`
	Geometry           json.RawMessage `json:"geometry"`
	LegalDoc           *string         `json:"legalDoc"`
	PlanningName       *string         `json:"planningName"`
	OriginalProperties json.RawMessage `json:"originalProperties"`
	SourceFile         *string         `json:"sourceFile"`
	ImportBatchID      *string         `json:"importBatchId"`
	Status             *uint32         `json:"status"`
}) (*httpapi.Result, error) {
	out, err := c.extends.Update(ctx, services.RegionExtendUpdatePatch{
		ID: req.ID, LayerID: req.LayerID, Name: req.Name, DisplayName: req.DisplayName, Description: req.Description,
		LabelID: req.LabelID, Geometry: []byte(req.Geometry), LegalDoc: req.LegalDoc, PlanningName: req.PlanningName,
		OriginalProperties: []byte(req.OriginalProperties), SourceFile: req.SourceFile, ImportBatchID: req.ImportBatchID,
		Status: req.Status,
	})
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *AdminRegionImportController) deleteExtend(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	if err := c.extends.Delete(ctx, req.ID); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true})
}

func (c *AdminRegionImportController) listExtendsByLayer(ctx context.Context, req struct {
	LayerID uint64 `path:"layer_id"`
	Page    uint32 `form:"page"`
	Size    uint32 `form:"size"`
}) (*httpapi.Result, error) {
	pagable := _dto.Pagable{Page: req.Page, Size: req.Size}
	pagable.Normalize()
	rows, total, err := c.extends.List(ctx, req.LayerID, nil, &pagable)
	if err != nil {
		return nil, err
	}
	return jsonListResult(rows, total, pagable.GetPage(), pagable.GetSize())
}

func (c *AdminRegionImportController) importPreview(ctx context.Context, req dto.AdminImportPreviewRequest) (*httpapi.Result, error) {
	out, err := c.imports.Preview(ctx, req.GeoJSON, req.LayerID)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *AdminRegionImportController) importBatch(ctx context.Context, req dto.BatchIDPath) (*httpapi.Result, error) {
	out, err := c.imports.GetBatchStatus(ctx, req.BatchID)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *AdminRegionImportController) importRollback(ctx context.Context, req dto.BatchIDPath) (*httpapi.Result, error) {
	n, err := c.imports.Rollback(ctx, req.BatchID)
	if err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"deleted": n})
}

func (c *AdminRegionImportController) importRetry(ctx context.Context, req dto.ErrorIDPath) (*httpapi.Result, error) {
	out, err := c.imports.RetryImportError(ctx, req.ErrorID)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *AdminRegionImportController) importErrors(ctx context.Context, req dto.AdminListImportErrorsRequest) (*httpapi.Result, error) {
	pagable := _dto.Pagable{Page: req.Page, Size: req.Size}
	pagable.Normalize()
	rows, total, err := c.imports.ListImportRegionErrors(ctx, req.LayerID, &pagable)
	if err != nil {
		return nil, err
	}
	return jsonListResult(rows, total, pagable.GetPage(), pagable.GetSize())
}

var _ = _errors.ReturnError
