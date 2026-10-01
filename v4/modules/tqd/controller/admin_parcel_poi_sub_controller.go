package controller

import (
	"context"
	"net/http"

	_dto "common/domain/dto"
	_errors "common/errors"
	"userv4/internal/httpapi"
	"userv4/modules/tqd/dto"
	"userv4/modules/tqd/services"
)

// AdminParcelPoiSubController owns admin parcels, pois, and subscriptions.
type AdminParcelPoiSubController struct {
	parcels services.ParcelUsecase
	pois    services.PoiUsecase
	subs    services.SubscriptionUsecase
}

func NewAdminParcelPoiSubController(
	parcels services.ParcelUsecase,
	pois services.PoiUsecase,
	subs services.SubscriptionUsecase,
) *AdminParcelPoiSubController {
	return &AdminParcelPoiSubController{parcels: parcels, pois: pois, subs: subs}
}

func (c *AdminParcelPoiSubController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/tqd/admin/parcels", c.listParcels)
	httpapi.GET(mux, "/v2/tqd/admin/parcels/{parcelId}", c.getParcel)
	httpapi.GET(mux, "/v2/tqd/admin/parcels/{parcelId}/neighbors", c.parcelNeighbors)
	httpapi.PUT(mux, "/v2/tqd/admin/parcels/{parcelId}", c.updateParcel)
	httpapi.POST(mux, "/v2/tqd/admin/parcels/gen-seo", c.genParcelSeo)

	httpapi.GET(mux, "/v2/tqd/admin/pois", c.listPois)
	httpapi.GET(mux, "/v2/tqd/admin/pois/{id}", c.getPoi)
	httpapi.PUT(mux, "/v2/tqd/admin/pois/{id}", c.updatePoi)
	httpapi.DELETE(mux, "/v2/tqd/admin/pois/{id}", c.deletePoi)

	httpapi.GET(mux, "/v2/tqd/admin/subscriptions", c.listSubs)
	httpapi.GET(mux, "/v2/tqd/admin/subscriptions/parcels", c.listParcelSubs)
	httpapi.DELETE(mux, "/v2/tqd/admin/subscriptions/{subscriptionId}", c.deleteSub)
}

type adminPageSearch struct {
	Page   uint32 `form:"page"`
	Size   uint32 `form:"size"`
	Search string `form:"search"`
}

func (c *AdminParcelPoiSubController) listParcels(ctx context.Context, req adminPageSearch) (*httpapi.Result, error) {
	pagable := _dto.Pagable{Page: req.Page, Size: req.Size}
	pagable.Normalize()
	rows, total, err := c.parcels.AdminListParcels(ctx, req.Search, int(pagable.GetPage()), int(pagable.GetSize()))
	if err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"data": rows, "total": total})
}

func (c *AdminParcelPoiSubController) getParcel(ctx context.Context, req dto.ParcelIDPath) (*httpapi.Result, error) {
	out, err := c.parcels.AdminGetParcel(ctx, req.ParcelID)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *AdminParcelPoiSubController) parcelNeighbors(ctx context.Context, req dto.ParcelIDPath) (*httpapi.Result, error) {
	out, err := c.parcels.AdminGetParcelNeighbors(ctx, req.ParcelID)
	if err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"data": out})
}

type updateParcelHTTP struct {
	ParcelID     uint64  `path:"parcelId" json:"-"`
	Lat          float64 `json:"lat"`
	Lon          float64 `json:"lon"`
	AreaSqm      float64 `json:"areaSqm"`
	Address      string  `json:"address"`
	MapNumber    string  `json:"mapNumber"`
	LandNumber   string  `json:"landNumber"`
	PropertyCode string  `json:"propertyCode"`
	ProvinceCode string  `json:"provinceCode"`
	WardCode     string  `json:"wardCode"`
	IsSeo        bool    `json:"isSeo"`
}

func (c *AdminParcelPoiSubController) updateParcel(ctx context.Context, req updateParcelHTTP) (*httpapi.Result, error) {
	out, err := c.parcels.AdminUpdateParcel(ctx, &dto.AdminUpdateParcelRequest{
		ParcelID: req.ParcelID, Lat: req.Lat, Lon: req.Lon, AreaSqm: req.AreaSqm, Address: req.Address,
		MapNumber: req.MapNumber, LandNumber: req.LandNumber, PropertyCode: req.PropertyCode,
		ProvinceCode: req.ProvinceCode, WardCode: req.WardCode, IsSeo: req.IsSeo,
	})
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *AdminParcelPoiSubController) genParcelSeo(ctx context.Context, req dto.AdminGenSeoIDsRequest) (*httpapi.Result, error) {
	out, err := c.parcels.AdminGenParcelSeo(ctx, req.IDs)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *AdminParcelPoiSubController) listPois(ctx context.Context, req adminPageSearch) (*httpapi.Result, error) {
	pagable := _dto.Pagable{Page: req.Page, Size: req.Size}
	pagable.Normalize()
	out, err := c.pois.List(ctx, &dto.PoiFilter{Pagable: pagable, Search: req.Search})
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *AdminParcelPoiSubController) getPoi(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	out, err := c.pois.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *AdminParcelPoiSubController) updatePoi(ctx context.Context, req struct {
	ID uint64 `path:"id" json:"-"`
	dto.UpdatePoiRequest
}) (*httpapi.Result, error) {
	out, err := c.pois.Update(ctx, req.ID, &req.UpdatePoiRequest)
	if err != nil {
		return nil, err
	}
	return jsonOK(out)
}

func (c *AdminParcelPoiSubController) deletePoi(ctx context.Context, req dto.IDPath) (*httpapi.Result, error) {
	if err := c.pois.Delete(ctx, req.ID); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true})
}

type adminSubListReq struct {
	Page       uint32  `form:"page"`
	Size       uint32  `form:"size"`
	UserID     *uint64 `form:"userId"`
	TargetType *string `form:"targetType"`
	Status     *uint32 `form:"status"`
}

func (c *AdminParcelPoiSubController) listSubs(ctx context.Context, req adminSubListReq) (*httpapi.Result, error) {
	pagable := _dto.Pagable{Page: req.Page, Size: req.Size}
	pagable.Normalize()
	rows, total, err := c.subs.AdminList(ctx, req.UserID, req.TargetType, req.Status, int(pagable.GetPage()), int(pagable.GetSize()))
	if err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"data": rows, "total": total})
}

func (c *AdminParcelPoiSubController) listParcelSubs(ctx context.Context, req adminSubListReq) (*httpapi.Result, error) {
	pagable := _dto.Pagable{Page: req.Page, Size: req.Size}
	pagable.Normalize()
	rows, total, err := c.subs.AdminListParcelSubscriptions(ctx, req.UserID, req.Status, int(pagable.GetPage()), int(pagable.GetSize()))
	if err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"data": rows, "total": total})
}

func (c *AdminParcelPoiSubController) deleteSub(ctx context.Context, req dto.SubscriptionIDPath) (*httpapi.Result, error) {
	if err := c.subs.AdminDelete(ctx, req.SubscriptionID); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true})
}

var _ = _errors.ReturnError
