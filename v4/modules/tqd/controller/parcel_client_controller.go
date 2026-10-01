package controller

import (
	"context"
	"net/http"
	"strings"

	_dto "common/domain/dto"
	_errors "common/errors"
	_utils "common/utils"
	tqdpb "pb/types/tqd"
	"userv4/internal/httpapi"
	"userv4/modules/tqd/dto"
	"userv4/modules/tqd/mapper"
	"userv4/modules/tqd/services"
)

// ParcelClientController owns parcel client routes under /v2/tqd/client.
type ParcelClientController struct {
	parcels services.ParcelUsecase
	mapper  *mapper.ParcelMapper
}

func NewParcelClientController(parcels services.ParcelUsecase) *ParcelClientController {
	return &ParcelClientController{
		parcels: parcels,
		mapper:  mapper.NewParcelMapper(),
	}
}

func (c *ParcelClientController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/tqd/client/parcels/search", c.searchParcels)
	httpapi.GET(mux, "/v2/tqd/client/administrative-units/{compositeId}/quick-layers", c.quickLayers)
}

func (c *ParcelClientController) searchParcels(ctx context.Context, req dto.ClientSearchParcelsRequest) (*httpapi.Result, error) {
	searchInfo := _utils.InferParcelSearch(req.Text)
	if searchInfo.CleanedText == "" && searchInfo.MapNumber == "" && searchInfo.LandNumber == "" {
		return nil, _errors.ReturnError(400, "search text is required")
	}
	pagable := _dto.Pagable{Page: req.Page, Size: req.Size}
	pagable.Normalize()
	rows, total, err := c.parcels.SearchParcelsByText(
		ctx,
		searchInfo.CleanedText,
		searchInfo.MapNumber,
		searchInfo.LandNumber,
		int(pagable.GetPage()),
		int(pagable.GetSize()),
	)
	if err != nil {
		return nil, err
	}
	items := make([]*tqdpb.ParcelInfo, 0, len(rows))
	for _, p := range rows {
		items = append(items, c.mapper.ToProtoParcelInfo(p))
	}
	return protoListResult(items, total, pagable.GetPage(), pagable.GetSize())
}

func (c *ParcelClientController) quickLayers(ctx context.Context, req dto.ClientCompositeIDRequest) (*httpapi.Result, error) {
	rows, err := c.parcels.GetAdministrativeUnitQuickInfo(ctx, req.CompositeID)
	if err != nil {
		if strings.Contains(err.Error(), "decode compositeId") || strings.Contains(err.Error(), "invalid compositeId kind") {
			return nil, _errors.ReturnError(400, err.Error())
		}
		return nil, err
	}
	return protoResult(c.mapper.BuildLayerQuickInfoFromRows(rows))
}
