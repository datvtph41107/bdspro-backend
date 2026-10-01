package controller

import (
	"context"
	"net/http"
	"time"

	_dto "common/domain/dto"
	_errors "common/errors"
	_utils "common/utils"
	tqdpb "pb/types/tqd"
	"userv4/internal/httpapi"
	"userv4/modules/tqd/dto"
	qh "userv4/modules/tqd/entity/qh"
	"userv4/modules/tqd/enums"
	"userv4/modules/tqd/mapper"
	"userv4/modules/tqd/repository"
	"userv4/modules/tqd/services"

	"github.com/lib/pq"
)

// QHAdminLayerController owns /v2/tqd/qh/admin/layers*.
type QHAdminLayerController struct {
	layers services.LayerUsecase
}

func NewQHAdminLayerController(layers services.LayerUsecase) *QHAdminLayerController {
	return &QHAdminLayerController{layers: layers}
}

func (c *QHAdminLayerController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.POST(mux, "/v2/tqd/qh/admin/layers", c.create)
	httpapi.POST(mux, "/v2/tqd/qh/admin/layers/gen-seo", c.genSeo)
	httpapi.GET(mux, "/v2/tqd/qh/admin/layers", c.list)
	httpapi.GET(mux, "/v2/tqd/qh/admin/layers/{id}", c.get)
	httpapi.PUT(mux, "/v2/tqd/qh/admin/layers/{id}", c.update)
	httpapi.PATCH(mux, "/v2/tqd/qh/admin/layers/{id}/status", c.updateStatus)
	httpapi.POST(mux, "/v2/tqd/qh/admin/layers/{id}/toggle-basic-visibility", c.toggleVisibility)
	httpapi.DELETE(mux, "/v2/tqd/qh/admin/layers/{id}", c.delete)
	httpapi.DELETE(mux, "/v2/tqd/qh/admin/layers/{id}/delete", c.hardDelete)
}

type createLayerAdminReq struct {
	Name            string   `json:"name"`
	DisplayName     string   `json:"displayName"`
	Description     string   `json:"description"`
	Type            uint32   `json:"type"`
	Visible         int32    `json:"visible"`
	DisplayOrder    int32    `json:"displayOrder"`
	MinZoom         uint32   `json:"minZoom"`
	MaxZoom         uint32   `json:"maxZoom"`
	Avatar          string   `json:"avatar"`
	ImageURL        string   `json:"imageUrl"`
	ThumbnailURL    string   `json:"thumbnailUrl"`
	SourceType      string   `json:"sourceType"`
	EffectiveDate   *string  `json:"effectiveDate"`
	ExpiryDate      *string  `json:"expiryDate"`
	LegalStatus     *uint32  `json:"legalStatus"`
	TrustValue      *float64 `json:"trustValue"`
	FamilyID        *uint64  `json:"familyId"`
	LegalDoc        string   `json:"legalDoc"`
	PublishScopes   []int32  `json:"publishScopes"`
	ProvinceID      string   `json:"provinceId"`
	WardID          string   `json:"wardId"`
	Status          *uint32  `json:"status"`
	ReplaceLayerIDs []uint64 `json:"replaceLayerIds"`
}

type updateLayerAdminReq struct {
	ID              uint64   `path:"id" json:"-"`
	Name            *string  `json:"name"`
	DisplayName     *string  `json:"displayName"`
	Description     *string  `json:"description"`
	Type            *uint32  `json:"type"`
	Status          *uint32  `json:"status"`
	Visible         *int32   `json:"visible"`
	DisplayOrder    *int32   `json:"displayOrder"`
	SourceType      string   `json:"sourceType"`
	MinZoom         *uint32  `json:"minZoom"`
	MaxZoom         *uint32  `json:"maxZoom"`
	Avatar          *string  `json:"avatar"`
	ImageURL        *string  `json:"imageUrl"`
	ThumbnailURL    *string  `json:"thumbnailUrl"`
	EffectiveDate   *string  `json:"effectiveDate"`
	ExpiryDate      *string  `json:"expiryDate"`
	LegalStatus     *uint32  `json:"legalStatus"`
	TrustValue      *float64 `json:"trustValue"`
	FamilyID        *uint64  `json:"familyId"`
	LegalDoc        *string  `json:"legalDoc"`
	PublishScopes   []int32  `json:"publishScopes"`
	ProvinceID      *string  `json:"provinceId"`
	WardID          *string  `json:"wardId"`
	DefaultVisible  *bool    `json:"defaultVisible"`
}

type updateLayerStatusReq struct {
	ID     uint64 `path:"id" json:"-"`
	Status uint32 `json:"status"`
}

type listLayersAdminReq struct {
	Page              uint32  `form:"page"`
	Size              uint32  `form:"size"`
	Search            *string `form:"search"`
	Type              *uint32 `form:"type"`
	Status            *uint32 `form:"status"`
	FamilyID          *uint64 `form:"familyId"`
	PlanningProjectID *uint64 `form:"planningProjectId"`
}

func (c *QHAdminLayerController) create(ctx context.Context, req createLayerAdminReq) (*httpapi.Result, error) {
	userID := _utils.GetProfileIdWithContext(ctx)
	if userID == 0 {
		return nil, _errors.ReturnError(401, "unauthorized")
	}
	if req.DisplayName == "" {
		return nil, _errors.ReturnError(400, "display_name is required")
	}
	layer := &qh.QHLayer{
		Name: req.Name, DisplayName: req.DisplayName, Description: req.Description,
		Type: enums.LayerType(req.Type), Visible: req.Visible, DisplayOrder: int(req.DisplayOrder),
		MinZoom: req.MinZoom, MaxZoom: req.MaxZoom, Avatar: req.Avatar,
		ImageURL: req.ImageURL, ThumbnailURL: req.ThumbnailURL, SourceType: req.SourceType,
		LegalDoc: req.LegalDoc, LegalStatus: enums.LegalStatusUnknown,
	}
	if req.EffectiveDate != nil && *req.EffectiveDate != "" {
		if t, err := time.Parse("2006-01-02", *req.EffectiveDate); err == nil {
			layer.EffectiveDate = &t
		}
	}
	if req.ExpiryDate != nil && *req.ExpiryDate != "" {
		if t, err := time.Parse("2006-01-02", *req.ExpiryDate); err == nil {
			layer.ExpiryDate = &t
		}
	}
	if req.LegalStatus != nil {
		layer.LegalStatus = enums.LegalStatus(*req.LegalStatus)
	}
	if req.TrustValue != nil {
		layer.TrustValue = float32(*req.TrustValue)
	}
	if req.FamilyID != nil && *req.FamilyID != 0 {
		layer.FamilyID = req.FamilyID
	}
	if len(req.PublishScopes) > 0 {
		layer.PublishScopes = pq.Int32Array(req.PublishScopes)
	}
	layer.ProvinceID = mapper.NormalizeLocationID("province", req.ProvinceID)
	layer.WardID = mapper.NormalizeLocationID("ward", req.WardID)
	if req.Status != nil {
		layer.Status = enums.LayerStatus(*req.Status)
	} else {
		layer.Status = enums.LayerStatusDraft
	}
	created, err := c.layers.Create(ctx, layer, userID, req.ReplaceLayerIDs)
	if err != nil {
		return nil, err
	}
	protoLayer, err := mapper.ToProtoLayer(created)
	if err != nil {
		return nil, err
	}
	return protoResult(protoLayer)
}

func (c *QHAdminLayerController) update(ctx context.Context, req updateLayerAdminReq) (*httpapi.Result, error) {
	userID := _utils.GetOriginIdFromContext(ctx)
	existing, err := c.layers.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, _errors.ReturnError(404, "layer not found")
	}
	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.DisplayName != nil {
		existing.DisplayName = *req.DisplayName
	}
	if req.Description != nil {
		existing.Description = *req.Description
	}
	if req.Type != nil {
		existing.Type = enums.LayerType(*req.Type)
	}
	if req.Status != nil {
		existing.Status = enums.LayerStatus(*req.Status)
	}
	if req.Visible != nil {
		existing.Visible = *req.Visible
	}
	if req.DisplayOrder != nil {
		existing.DisplayOrder = int(*req.DisplayOrder)
	}
	if req.SourceType != "" {
		existing.SourceType = req.SourceType
	}
	if req.MinZoom != nil {
		existing.MinZoom = *req.MinZoom
	}
	if req.MaxZoom != nil {
		existing.MaxZoom = *req.MaxZoom
	}
	if req.Avatar != nil {
		existing.Avatar = *req.Avatar
	}
	if req.ImageURL != nil {
		existing.ImageURL = *req.ImageURL
	}
	if req.ThumbnailURL != nil {
		existing.ThumbnailURL = *req.ThumbnailURL
	}
	if req.EffectiveDate != nil && *req.EffectiveDate != "" {
		if t, err := time.Parse("2006-01-02", *req.EffectiveDate); err == nil {
			existing.EffectiveDate = &t
		}
	}
	if req.ExpiryDate != nil && *req.ExpiryDate != "" {
		if t, err := time.Parse("2006-01-02", *req.ExpiryDate); err == nil {
			existing.ExpiryDate = &t
		}
	}
	if req.LegalStatus != nil {
		existing.LegalStatus = enums.LegalStatus(*req.LegalStatus)
	}
	if req.TrustValue != nil {
		existing.TrustValue = float32(*req.TrustValue)
	}
	if req.FamilyID != nil {
		existing.FamilyID = req.FamilyID
	}
	if req.LegalDoc != nil {
		existing.LegalDoc = *req.LegalDoc
	}
	if req.PublishScopes != nil {
		existing.PublishScopes = pq.Int32Array(req.PublishScopes)
	}
	if req.ProvinceID != nil {
		existing.ProvinceID = mapper.NormalizeLocationID("province", *req.ProvinceID)
	}
	if req.WardID != nil {
		existing.WardID = mapper.NormalizeLocationID("ward", *req.WardID)
	}
	if req.DefaultVisible != nil {
		existing.DefaultVisible = *req.DefaultVisible
	}
	updated, err := c.layers.Update(ctx, req.ID, existing, userID)
	if err != nil {
		return nil, err
	}
	protoLayer, err := mapper.ToProtoLayer(updated)
	if err != nil {
		return nil, err
	}
	return protoResult(protoLayer)
}

func (c *QHAdminLayerController) get(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	if req.ID == 0 {
		return nil, _errors.ReturnError(400, "id is required")
	}
	layer, err := c.layers.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	protoLayer, err := mapper.ToProtoLayer(layer)
	if err != nil {
		return nil, err
	}
	return protoResult(protoLayer)
}

func (c *QHAdminLayerController) list(ctx context.Context, req listLayersAdminReq) (*httpapi.Result, error) {
	filter := &repository.LayerFilter{Pagable: _dto.Pagable{Page: req.Page, Size: req.Size}}
	filter.Search = req.Search
	if req.Type != nil {
		t := enums.LayerType(*req.Type)
		filter.Type = &t
	}
	if req.Status != nil {
		s := enums.LayerStatus(*req.Status)
		filter.Status = &s
	}
	filter.PlanningProjectID = req.PlanningProjectID
	rows, total, err := c.layers.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	items := make([]*tqdpb.LayerResponse, 0, len(rows))
	for i := range rows {
		p, err := mapper.ToProtoLayer(&rows[i])
		if err != nil {
			continue
		}
		items = append(items, p)
	}
	return protoListResult(items, total, req.Page, req.Size)
}

func (c *QHAdminLayerController) updateStatus(ctx context.Context, req updateLayerStatusReq) (*httpapi.Result, error) {
	userID := _utils.GetProfileIdWithContext(ctx)
	if err := c.layers.UpdateStatus(ctx, req.ID, enums.LayerStatus(req.Status), userID); err != nil {
		return nil, err
	}
	return jsonOK(map[string]string{"message": "updated"})
}

func (c *QHAdminLayerController) toggleVisibility(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	userID := _utils.GetProfileIdWithContext(ctx)
	out, err := c.layers.ToggleLayerBasicVisibility(ctx, req.ID, userID)
	if err != nil {
		return nil, err
	}
	protoLayer, err := mapper.ToProtoLayer(out)
	if err != nil {
		return nil, err
	}
	return protoResult(protoLayer)
}

func (c *QHAdminLayerController) delete(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	userID := _utils.GetProfileIdWithContext(ctx)
	if err := c.layers.Delete(ctx, req.ID, userID); err != nil {
		return nil, err
	}
	return jsonOK(map[string]string{"message": "deleted"})
}

func (c *QHAdminLayerController) hardDelete(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	userID := _utils.GetProfileIdWithContext(ctx)
	if err := c.layers.HardDelete(ctx, req.ID, userID); err != nil {
		return nil, err
	}
	return jsonOK(map[string]string{"message": "deleted"})
}

func (c *QHAdminLayerController) genSeo(ctx context.Context, req dto.AdminGenSeoIDsRequest) (*httpapi.Result, error) {
	out, err := c.layers.AdminGenLayerSeo(ctx, req.IDs)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}
