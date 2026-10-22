package controller

import (
	"context"
	"net/http"

	_dto "common/domain/dto"
	_errors "common/errors"
	"common/httpapi"
	_utils "common/utils"
	"mdn4/modules/tqd/dto"
	"mdn4/modules/tqd/mapper"
	"mdn4/modules/tqd/services"
	tqdpb "pb/types/tqd"
)

// LocationController owns TQD client location + administrative-unit HTTP routes under /v2/tqd/client.
type LocationController struct {
	location   services.LocationService
	locationUC *services.LocationUsecase
	locMapper  *mapper.LocationMapper
}

func NewLocationController(location services.LocationService, locationUC *services.LocationUsecase) *LocationController {
	return &LocationController{
		location:   location,
		locationUC: locationUC,
		locMapper:  mapper.NewLocationMapper(),
	}
}

func (c *LocationController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/tqd/client/provinces", c.listProvinces)
	httpapi.GET(mux, "/v2/tqd/client/wards", c.listWards)

	httpapi.GET(mux, "/v2/tqd/client/administrative-units", c.listAdministrativeUnits)
	httpapi.GET(mux, "/v2/tqd/client/administrative-units/{compositeId}", c.getAdministrativeUnit)
	httpapi.GET(mux, "/v2/tqd/client/administrative-units/{compositeId}/neighbors", c.listAdministrativeUnitNeighbors)

	httpapi.GET(mux, "/v2/tqd/client/followed-administrative-units", c.listFollowedAdministrativeUnits)
	httpapi.POST(mux, "/v2/tqd/client/followed-administrative-units/{compositeId}", c.followAdministrativeUnit)
	httpapi.DELETE(mux, "/v2/tqd/client/followed-administrative-units/{compositeId}", c.unfollowAdministrativeUnit)
}

func (c *LocationController) listProvinces(ctx context.Context, req dto.ListProvincesRequest) (*httpapi.Result, error) {
	result, err := c.location.ListProvinces(ctx, req)
	if err != nil {
		return nil, err
	}
	return listProvincesResult(result)
}

func (c *LocationController) listWards(ctx context.Context, req dto.ListWardsRequest) (*httpapi.Result, error) {
	if req.Page != nil && *req.Page < 1 {
		return nil, _errors.ReturnError(400, "page must be >= 1")
	}
	result, err := c.location.ListWards(ctx, req)
	if err != nil {
		return nil, err
	}
	return listWardsResult(result)
}

func (c *LocationController) requireLocationUC() error {
	if c.locationUC == nil {
		return _errors.ReturnError(501, "location usecase not configured")
	}
	return nil
}

func (c *LocationController) listAdministrativeUnits(ctx context.Context, req dto.ClientListAdministrativeUnitsRequest) (*httpapi.Result, error) {
	if err := c.requireLocationUC(); err != nil {
		return nil, err
	}
	dtoReq := &dto.ListAdministrativeUnitsRequest{
		Pagable: _dto.Pagable{
			Page: req.Page,
			Size: req.Size,
			Sort: req.Sort,
		},
		Search:            req.Search,
		PlanningProjectID: req.PlanningProjectID,
		CompositeIDs:      _utils.ParseUint64Array(req.CompositeIDs),
	}
	resp, err := c.locationUC.ClientListAdministrativeUnits(ctx, dtoReq)
	if err != nil {
		return nil, err
	}
	items := make([]*tqdpb.AdministrativeUnitListItem, 0)
	total := int64(0)
	if resp != nil {
		total = resp.Total
		proto := c.locMapper.ToProtoListAdministrativeUnitsResponse(resp)
		if proto != nil {
			items = proto.GetData()
		}
	}
	return protoListResult(items, total, req.Page, req.Size)
}

func (c *LocationController) getAdministrativeUnit(ctx context.Context, req dto.ClientCompositeIDRequest) (*httpapi.Result, error) {
	if err := c.requireLocationUC(); err != nil {
		return nil, err
	}
	province, ward, err := c.locationUC.ClientGetAdministrativeUnitByCompositeId(ctx, req.CompositeID)
	if err != nil {
		return nil, err
	}
	if province == nil && ward == nil {
		return nil, _errors.ReturnError(404, "administrative unit not found")
	}
	return protoResult(&tqdpb.ClientGetAdministrativeUnitDetailByCompositeIdResponse{
		Province: c.locMapper.ToProtoProvince(province),
		Ward:     c.locMapper.ToProtoWard(ward),
	})
}

func (c *LocationController) listAdministrativeUnitNeighbors(ctx context.Context, req dto.ClientCompositeIDRequest) (*httpapi.Result, error) {
	if err := c.requireLocationUC(); err != nil {
		return nil, err
	}
	resp, err := c.locationUC.ClientListAdministrativeUnitNeighbors(ctx, req.CompositeID)
	if err != nil {
		return nil, err
	}
	items := make([]*tqdpb.Ward, 0)
	total := int64(0)
	if resp != nil {
		total = resp.Total
		items = make([]*tqdpb.Ward, 0, len(resp.Data))
		for _, w := range resp.Data {
			items = append(items, c.locMapper.ToProtoWard(w))
		}
	}
	data, err := protoMessages(items)
	if err != nil {
		return nil, err
	}
	return httpapi.OK(data, httpapi.WithPagination(&httpapi.Pagination{Total: total})), nil
}

func (c *LocationController) listFollowedAdministrativeUnits(ctx context.Context, req dto.ClientFollowedAdministrativeUnitsRequest) (*httpapi.Result, error) {
	if err := c.requireLocationUC(); err != nil {
		return nil, err
	}
	if _utils.GetProfileIdWithContext(ctx) == 0 {
		return nil, _errors.ReturnError(401, "authentication is required")
	}
	resp, err := c.locationUC.ListFollowedAdministrativeUnits(ctx, req.Page, req.Size)
	if err != nil {
		return nil, err
	}
	items := make([]*tqdpb.AdministrativeUnitListItem, 0)
	total := int64(0)
	if resp != nil {
		total = resp.Total
		proto := c.locMapper.ToProtoListAdministrativeUnitsResponse(resp)
		if proto != nil {
			items = proto.GetData()
		}
	}
	return protoListResult(items, total, req.Page, req.Size)
}

func (c *LocationController) followAdministrativeUnit(ctx context.Context, req dto.ClientCompositeIDRequest) (*httpapi.Result, error) {
	if err := c.requireLocationUC(); err != nil {
		return nil, err
	}
	if _utils.GetProfileIdWithContext(ctx) == 0 {
		return nil, _errors.ReturnError(401, "authentication is required")
	}
	if err := c.locationUC.FollowAdministrativeUnit(ctx, req.CompositeID); err != nil {
		return nil, err
	}
	return protoResult(&tqdpb.FollowAdministrativeUnitResponse{
		Followed:    true,
		CompositeId: req.CompositeID,
	})
}

func (c *LocationController) unfollowAdministrativeUnit(ctx context.Context, req dto.ClientCompositeIDRequest) (*httpapi.Result, error) {
	if err := c.requireLocationUC(); err != nil {
		return nil, err
	}
	if _utils.GetProfileIdWithContext(ctx) == 0 {
		return nil, _errors.ReturnError(401, "authentication is required")
	}
	if err := c.locationUC.UnfollowAdministrativeUnit(ctx, req.CompositeID); err != nil {
		return nil, err
	}
	return protoResult(&tqdpb.FollowAdministrativeUnitResponse{
		Followed:    false,
		CompositeId: req.CompositeID,
	})
}

func listProvincesResult(resp *tqdpb.ListProvincesResponse) (*httpapi.Result, error) {
	if resp == nil {
		return httpapi.OK([]any{}), nil
	}
	data, err := protoMessages(resp.GetProvinces())
	if err != nil {
		return nil, err
	}
	data = normalizeClientLocationCatalogJSON(data)
	var opts []httpapi.OKOption
	if resp.GetTimestamp() > 0 {
		opts = append(opts, httpapi.WithTimestamp(resp.GetTimestamp()))
	}
	return httpapi.OK(data, opts...), nil
}

func listWardsResult(resp *tqdpb.ClientListWardsResponse) (*httpapi.Result, error) {
	if resp == nil {
		return httpapi.OK([]any{}), nil
	}
	data, err := protoMessages(resp.GetData())
	if err != nil {
		return nil, err
	}
	data = normalizeClientLocationCatalogJSON(data)
	var opts []httpapi.OKOption
	if p := resp.GetPagination(); p != nil {
		opts = append(opts, httpapi.WithPagination(&httpapi.Pagination{
			Page:  p.GetPage(),
			Size:  p.GetSize(),
			Total: p.GetTotal(),
		}))
	}
	if resp.GetTimestamp() > 0 {
		opts = append(opts, httpapi.WithTimestamp(resp.GetTimestamp()))
	}
	return httpapi.OK(data, opts...), nil
}
