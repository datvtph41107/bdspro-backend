package controller

import (
	"context"
	"net/http"

	_dto "common/domain/dto"
	_errors "common/errors"
	"common/httpapi"
	"mdn4/modules/tqd/dto"
	"mdn4/modules/tqd/mapper"
	"mdn4/modules/tqd/services"
	tqdpb "pb/types/tqd"
)

// AdminLocationController owns /v2/tqd/admin provinces/wards/administrative-units.
type AdminLocationController struct {
	uc     *services.LocationUsecase
	mapper *mapper.LocationMapper
}

func NewAdminLocationController(uc *services.LocationUsecase) *AdminLocationController {
	return &AdminLocationController{uc: uc, mapper: mapper.NewLocationMapper()}
}

func (c *AdminLocationController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/tqd/admin/provinces", c.listProvinces)
	httpapi.GET(mux, "/v2/tqd/admin/provinces/{compositeId}", c.getProvince)
	httpapi.POST(mux, "/v2/tqd/admin/provinces", c.createProvince)
	httpapi.PUT(mux, "/v2/tqd/admin/provinces/{id}", c.updateProvince)
	httpapi.DELETE(mux, "/v2/tqd/admin/provinces/{id}", c.deleteProvince)

	httpapi.GET(mux, "/v2/tqd/admin/wards", c.listWards)
	httpapi.POST(mux, "/v2/tqd/admin/wards", c.createWard)
	httpapi.PUT(mux, "/v2/tqd/admin/wards/{id}", c.updateWard)
	httpapi.DELETE(mux, "/v2/tqd/admin/wards/{id}", c.deleteWard)
	httpapi.GET(mux, "/v2/tqd/admin/wards/{ward_id}/neighbors", c.listNeighbors)
	httpapi.POST(mux, "/v2/tqd/admin/wards/{ward_id}/neighbors", c.createNeighbor)
	httpapi.DELETE(mux, "/v2/tqd/admin/wards/{ward_id}/neighbors/{neighbor_ward_id}", c.deleteNeighbor)

	httpapi.GET(mux, "/v2/tqd/admin/administrative-units", c.listAdminUnits)
	httpapi.GET(mux, "/v2/tqd/admin/administrative-units/{compositeId}", c.getAdminUnit)
	httpapi.POST(mux, "/v2/tqd/admin/administrative-units/seo/generate", c.genSeo)
	httpapi.POST(mux, "/v2/tqd/admin/administrative-units/area-sync-geometry", c.syncArea)
}

type adminListPage struct {
	Page   uint32 `form:"page"`
	Size   uint32 `form:"size"`
	Search string `form:"search"`
}

func (c *AdminLocationController) listProvinces(ctx context.Context, req adminListPage) (*httpapi.Result, error) {
	pagable := _dto.Pagable{Page: req.Page, Size: req.Size}
	pagable.Normalize()
	out, err := c.uc.AdminListProvinces(ctx, &dto.AdminListProvincesRequest{
		Pagable: pagable, Search: req.Search,
	})
	if err != nil {
		return nil, err
	}
	items := make([]*tqdpb.Province, 0)
	total := int64(0)
	if out != nil {
		total = out.Total
		items = make([]*tqdpb.Province, 0, len(out.Data))
		for _, p := range out.Data {
			items = append(items, c.mapper.ToProtoProvince(p))
		}
	}
	return protoListResult(items, total, pagable.GetPage(), pagable.GetSize())
}

func (c *AdminLocationController) getProvince(ctx context.Context, req dto.CompositeIDPath) (*httpapi.Result, error) {
	out, err := c.uc.AdminGetProvinceByCompositeId(ctx, req.CompositeID)
	if err != nil {
		return nil, err
	}
	return protoResult(c.mapper.ToProtoProvince(out))
}

func (c *AdminLocationController) createProvince(ctx context.Context, req dto.ProvinceDTO) (*httpapi.Result, error) {
	out, err := c.uc.AdminCreateProvince(ctx, &req)
	if err != nil {
		return nil, err
	}
	return protoResult(c.mapper.ToProtoProvince(out))
}

func (c *AdminLocationController) updateProvince(ctx context.Context, req dto.ProvinceDTO) (*httpapi.Result, error) {
	if req.ID == 0 {
		return nil, _errors.ReturnError(400, "id is required")
	}
	out, err := c.uc.AdminUpdateProvince(ctx, &req)
	if err != nil {
		return nil, err
	}
	return protoResult(c.mapper.ToProtoProvince(out))
}

func (c *AdminLocationController) deleteProvince(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	if err := c.uc.AdminDeleteProvince(ctx, req.ID); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true})
}

func (c *AdminLocationController) listWards(ctx context.Context, req adminListPage) (*httpapi.Result, error) {
	pagable := _dto.Pagable{Page: req.Page, Size: req.Size}
	pagable.Normalize()
	out, err := c.uc.AdminListWards(ctx, &dto.AdminListWardsRequest{
		Pagable: pagable, Search: req.Search,
	})
	if err != nil {
		return nil, err
	}
	if out == nil {
		return jsonListResult([]*dto.WardDTO{}, 0, pagable.GetPage(), pagable.GetSize())
	}
	return jsonListResult(out.Data, out.Total, pagable.GetPage(), pagable.GetSize())
}

func (c *AdminLocationController) createWard(ctx context.Context, req dto.WardDTO) (*httpapi.Result, error) {
	out, err := c.uc.AdminCreateWard(ctx, &req)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *AdminLocationController) updateWard(ctx context.Context, req dto.WardDTO) (*httpapi.Result, error) {
	out, err := c.uc.AdminUpdateWard(ctx, &req)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *AdminLocationController) deleteWard(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	if err := c.uc.AdminDeleteWard(ctx, req.ID); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true})
}

type createNeighborReq struct {
	WardID         uint64 `path:"ward_id" json:"-"`
	NeighborWardID uint64 `json:"neighborWardId"`
}

func (c *AdminLocationController) createNeighbor(ctx context.Context, req createNeighborReq) (*httpapi.Result, error) {
	if err := c.uc.AdminCreateWardNeighbor(ctx, req.WardID, req.NeighborWardID); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true})
}

func (c *AdminLocationController) deleteNeighbor(ctx context.Context, req dto.WardNeighborPath) (*httpapi.Result, error) {
	if err := c.uc.AdminDeleteWardNeighbor(ctx, req.WardID, req.NeighborWardID); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true})
}

func (c *AdminLocationController) listNeighbors(ctx context.Context, req dto.WardIDPath) (*httpapi.Result, error) {
	out, err := c.uc.AdminListWardNeighbors(ctx, req.WardID)
	if err != nil {
		return nil, err
	}
	if out == nil {
		return jsonListResult([]*dto.WardDTO{}, 0, 1, 0)
	}
	return jsonListResult(out.Data, out.Total, 1, uint32(len(out.Data)))
}

func (c *AdminLocationController) listAdminUnits(ctx context.Context, req adminListPage) (*httpapi.Result, error) {
	pagable := _dto.Pagable{Page: req.Page, Size: req.Size}
	pagable.Normalize()
	out, err := c.uc.AdminListAdministrativeUnits(ctx, &dto.ListAdministrativeUnitsRequest{
		Pagable: pagable, Search: req.Search,
	})
	if err != nil {
		return nil, err
	}
	if out == nil {
		return jsonListResult([]*dto.AdministrativeUnitListItem{}, 0, pagable.GetPage(), pagable.GetSize())
	}
	return jsonListResult(out.Data, out.Total, pagable.GetPage(), pagable.GetSize())
}

func (c *AdminLocationController) getAdminUnit(ctx context.Context, req dto.CompositeIDPath) (*httpapi.Result, error) {
	province, ward, err := c.uc.AdminGetAdministrativeUnitByCompositeId(ctx, req.CompositeID)
	if err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"province": province, "ward": ward})
}

func (c *AdminLocationController) genSeo(ctx context.Context, req dto.AdminGenAdminUnitSeoRequest) (*httpapi.Result, error) {
	out, err := c.uc.AdminGenAdministrativeUnitSeo(ctx, req.CompositeIDs)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *AdminLocationController) syncArea(ctx context.Context, req dto.AdminSyncAreaRequest) (*httpapi.Result, error) {
	out, err := c.uc.AdminSyncAreaFromGeometry(ctx, req.CompositeIDs)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}
