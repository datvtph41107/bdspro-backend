package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	_dto "common/domain/dto"
	_errors "common/errors"
	"common/httpapi"
	"mdn4/modules/tqd/dto"
	"mdn4/modules/tqd/mapper"
	repo "mdn4/modules/tqd/repository"
	"mdn4/modules/tqd/services"
	tqdpb "pb/types/tqd"
)

// ParcelUserController owns /v2/tqd/parcels/* + map-target + layer history/legal routes.
type ParcelUserController struct {
	parcels    services.ParcelUsecase
	layerLegal *services.LayerLegalUsecase
	mapper     *mapper.ParcelMapper
}

func NewParcelUserController(
	parcels services.ParcelUsecase,
	layerLegal *services.LayerLegalUsecase,
) *ParcelUserController {
	return &ParcelUserController{
		parcels:    parcels,
		layerLegal: layerLegal,
		mapper:     mapper.NewParcelMapper(),
	}
}

func (c *ParcelUserController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.POST(mux, "/v2/tqd/parcels/batch", c.createBatch)
	httpapi.POST(mux, "/v2/tqd/parcels/by-location", c.findByLocation)
	httpapi.POST(mux, "/v2/tqd/parcels/by-polygon", c.findByPolygon)
	httpapi.GET(mux, "/v2/tqd/parcels/{id}/preview", c.previewByID)
	httpapi.GET(mux, "/v2/tqd/parcels/{id}/public", c.publicByID)
	httpapi.GET(mux, "/v2/tqd/parcels/{parcelId}/around", c.around)
	httpapi.GET(mux, "/v2/tqd/parcels/{parcelId}/info", c.info)
	httpapi.GET(mux, "/v2/tqd/parcels/{parcelId}/detail", c.detail)
	httpapi.GET(mux, "/v2/tqd/parcels/{parcelId}/layers", c.layers)
	httpapi.GET(mux, "/v2/tqd/parcels/{parcelId}/layers/{layerId}/zones/{zoneId}/geometry", c.zoneGeometry)
	httpapi.GET(mux, "/v2/tqd/parcels/{parcelId}/legal-documents", c.legalDocuments)
	httpapi.GET(mux, "/v2/tqd/parcels/{parcelId}/planning", c.planning)
	httpapi.GET(mux, "/v2/tqd/parcels/{parcelId}/planning-eval", c.planningEval)
	httpapi.GET(mux, "/v2/tqd/parcels/{parcelId}/timeline", c.timeline)
	httpapi.GET(mux, "/v2/tqd/parcels/{parcelId}/history", c.history)

	httpapi.POST(mux, "/v2/tqd/map-target/by-location", c.mapTargetByLocation)
	httpapi.GET(mux, "/v2/tqd/map-target/by-poi/{poi_id}", c.mapTargetByPoi)

	httpapi.GET(mux, "/v2/tqd/layers/{layerId}/history", c.layerHistory)
	httpapi.POST(mux, "/v2/tqd/layers/{layerId}/legal-docs", c.addLegalDocs)
	httpapi.GET(mux, "/v2/tqd/layers/{layerId}/legal-docs", c.listLegalDocs)
	httpapi.DELETE(mux, "/v2/tqd/layers/legal-docs/{id}", c.deleteLegalDoc)
}

func (c *ParcelUserController) requireParcels() error {
	if c.parcels == nil {
		return _errors.ReturnError(501, "parcel service not configured")
	}
	return nil
}

func (c *ParcelUserController) createBatch(ctx context.Context, req dto.ParcelBatchHTTPRequest) (*httpapi.Result, error) {
	if err := c.requireParcels(); err != nil {
		return nil, err
	}
	if _, err := requireProfileID(ctx); err != nil {
		return nil, err
	}
	if len(req.Data) == 0 {
		return nil, _errors.ReturnError(400, "request data is required")
	}
	raw, err := json.Marshal(map[string]any{"data": req.Data})
	if err != nil {
		return nil, err
	}
	var listDTO tqdpb.ParcelListDTO
	if err := json.Unmarshal(raw, &listDTO); err != nil {
		return nil, _errors.ReturnError(400, "invalid parcel batch payload")
	}
	domains := c.mapper.ProtoToDomains(&listDTO)
	result, err := c.parcels.CreateBatch(ctx, &dto.ParcelList{Data: domains})
	if err != nil {
		return nil, err
	}
	return jsonOK(c.mapper.ToProtoFromDomainList(result))
}

func (c *ParcelUserController) findByLocation(ctx context.Context, req dto.FindByLocationHTTPRequest) (*httpapi.Result, error) {
	if err := c.requireParcels(); err != nil {
		return nil, err
	}
	if req.Latitude == 0 && req.Longitude == 0 {
		return nil, _errors.ReturnError(400, "latitude and longitude are required")
	}
	result, err := c.parcels.FindByLocation(ctx, req.Latitude, req.Longitude)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, _errors.ReturnError(404, "parcel not found")
	}
	return protoResult(c.mapper.DomainToProto(result))
}

func (c *ParcelUserController) findByPolygon(ctx context.Context, req dto.FindByPolygonHTTPRequest) (*httpapi.Result, error) {
	if err := c.requireParcels(); err != nil {
		return nil, err
	}
	points, filter, err := polygonRequestParts(req)
	if err != nil {
		return nil, err
	}
	pagable := _dto.Pagable{Page: req.Page, Size: req.Size}
	pagable.Normalize()
	result, err := c.parcels.FindMapTargetsByPolygon(
		ctx, points, req.Intersect, req.MaxAreaHa,
		int(pagable.GetPage()), int(pagable.GetSize()),
		req.Z, filter, req.IncludeAnalysis,
	)
	if err != nil {
		return nil, err
	}
	return protoResult(c.mapper.PolygonMapTargetsToProto(result))
}

func polygonRequestParts(req dto.FindByPolygonHTTPRequest) ([]repo.Point, *dto.PolygonTargetFilter, error) {
	if len(req.Coordinates) == 0 {
		return nil, nil, _errors.ReturnError(400, "coordinates are required")
	}
	points := make([]repo.Point, 0, len(req.Coordinates))
	for _, pt := range req.Coordinates {
		if pt.Lat < -90 || pt.Lat > 90 || pt.Lng < -180 || pt.Lng > 180 {
			return nil, nil, _errors.ReturnError(400, "invalid coordinates")
		}
		points = append(points, repo.Point{Lat: pt.Lat, Lng: pt.Lng})
	}
	if len(points) < 3 {
		return nil, nil, _errors.ReturnError(400, "at least 3 coordinates are required")
	}
	filter := &dto.PolygonTargetFilter{
		LayerIDs:     req.LayerIDs,
		LabelIDs:     req.LabelIDs,
		LandUseIDs:   req.LandUseIDs,
		LandUseCodes: req.LandUseCodes,
		Keyword:      strings.TrimSpace(req.Keyword),
		CanBuild:     req.CanBuild,
	}
	if !filter.HasPlanningFilters() {
		filter = nil
	}
	return points, filter, nil
}

func (c *ParcelUserController) previewByID(ctx context.Context, req dto.ParcelPreviewIDHTTPRequest) (*httpapi.Result, error) {
	if err := c.requireParcels(); err != nil {
		return nil, err
	}
	info, err := c.parcels.GetParcelInfoByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, _errors.ReturnError(404, "parcel not found")
	}
	return protoResult(c.mapper.ToProtoParcelInfo(info))
}

func (c *ParcelUserController) publicByID(ctx context.Context, req dto.ParcelPreviewIDHTTPRequest) (*httpapi.Result, error) {
	return c.previewByID(ctx, req)
}

func (c *ParcelUserController) around(ctx context.Context, req dto.ParcelIDHTTPRequest) (*httpapi.Result, error) {
	if err := c.requireParcels(); err != nil {
		return nil, err
	}
	limit := int(req.Limit)
	if limit <= 0 {
		limit = 5
	}
	rows, err := c.parcels.GetParcelsAround(ctx, req.ParcelID, limit)
	if err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{
		"data":  c.mapper.ToProtoNearbyParcelItems(rows),
		"total": len(rows),
	})
}

func (c *ParcelUserController) info(ctx context.Context, req dto.ParcelIDHTTPRequest) (*httpapi.Result, error) {
	if err := c.requireParcels(); err != nil {
		return nil, err
	}
	data, err := c.parcels.GetParcelInfo(ctx, req.ParcelID)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, _errors.ReturnError(404, "parcel not found")
	}
	return protoResult(c.mapper.ToProtoParcelInfoDetail(data))
}

func (c *ParcelUserController) detail(ctx context.Context, req dto.ParcelIDHTTPRequest) (*httpapi.Result, error) {
	if err := c.requireParcels(); err != nil {
		return nil, err
	}
	data, err := c.parcels.GetParcelDetail(ctx, req.ParcelID)
	if err != nil {
		return nil, err
	}
	if data == nil {
		return nil, _errors.ReturnError(404, "parcel not found")
	}
	return protoResult(c.mapper.ToProtoParcelDetail(data))
}

func (c *ParcelUserController) layers(ctx context.Context, req dto.ParcelIDHTTPRequest) (*httpapi.Result, error) {
	if err := c.requireParcels(); err != nil {
		return nil, err
	}
	pagable := _dto.Pagable{Page: req.Page, Size: req.Size}
	pagable.Normalize()
	rows, total, err := c.parcels.GetParcelLayers(ctx, req.ParcelID, pagable)
	if err != nil {
		return nil, err
	}
	return protoListResult(c.mapper.ToProtoLayerDetailItems(rows, req.ParcelID), total, pagable.GetPage(), pagable.GetSize())
}

func (c *ParcelUserController) zoneGeometry(ctx context.Context, req dto.ZoneGeometryHTTPRequest) (*httpapi.Result, error) {
	if err := c.requireParcels(); err != nil {
		return nil, err
	}
	geoJSON, props, err := c.parcels.GetZoneGeometry(ctx, req.ParcelID, req.LayerID, req.ZoneID)
	if err != nil {
		return nil, err
	}
	out := map[string]any{"geometry": geoJSON}
	if props != nil {
		out["zoneName"] = props.ZoneName
		out["landUseCode"] = props.LandUseCode
		out["landUseName"] = props.LandUseName
		out["landUseColor"] = props.LandUseColor
		out["buildStatus"] = props.BuildStatus
		out["alertLevel"] = props.AlertLevel
		out["legalDocumentIds"] = props.LegalDocumentIds
	}
	return jsonOK(out)
}

func (c *ParcelUserController) legalDocuments(ctx context.Context, req dto.ParcelIDHTTPRequest) (*httpapi.Result, error) {
	if err := c.requireParcels(); err != nil {
		return nil, err
	}
	pagable := _dto.Pagable{Page: req.Page, Size: req.Size}
	pagable.Normalize()
	docs, total, err := c.parcels.GetParcelLegalDocuments(ctx, req.ParcelID, req.DocType, req.Status, pagable)
	if err != nil {
		return nil, err
	}
	return protoListResult(c.mapper.ToProtoLegalDocumentItems(docs), total, pagable.GetPage(), pagable.GetSize())
}

func (c *ParcelUserController) planning(ctx context.Context, req dto.ParcelIDHTTPRequest) (*httpapi.Result, error) {
	if err := c.requireParcels(); err != nil {
		return nil, err
	}
	result, err := c.parcels.GetParcelPlanning(ctx, &dto.GetParcelPlanningRequest{
		ParcelID: req.ParcelID,
		Mode:     req.Mode,
		View:     req.View,
	})
	if err != nil {
		return nil, err
	}
	return jsonOK(result)
}

func (c *ParcelUserController) planningEval(ctx context.Context, req dto.ParcelIDHTTPRequest) (*httpapi.Result, error) {
	if req.ParcelID == 0 {
		return nil, _errors.ReturnError(400, "parcel_id is required")
	}
	return jsonOK(map[string]any{
		"parcelId":   req.ParcelID,
		"checkLevel": "HIGH",
		"mainAssessment": []string{
			"Thửa đất nằm trong vùng quy hoạch sử dụng đất giao thông (đường vành đai).",
			"Có lớp QH phân khu H2-4 chồng lấn, cần đối chiếu chi tiết.",
		},
	})
}

func (c *ParcelUserController) timeline(ctx context.Context, req dto.ParcelIDHTTPRequest) (*httpapi.Result, error) {
	if err := c.requireParcels(); err != nil {
		return nil, err
	}
	timeline, err := c.parcels.GetParcelTimeline(ctx, req.ParcelID)
	if err != nil {
		return nil, err
	}
	return jsonOK(timeline)
}

func (c *ParcelUserController) history(ctx context.Context, req dto.ParcelIDHTTPRequest) (*httpapi.Result, error) {
	if err := c.requireParcels(); err != nil {
		return nil, err
	}
	pagable := _dto.Pagable{Page: req.Page, Size: req.Size}
	pagable.Normalize()
	rows, total, err := c.parcels.GetParcelHistory(ctx, req.ParcelID, pagable)
	if err != nil {
		return nil, err
	}
	return protoListResult(c.mapper.ToProtoParcelHistoryItems(rows), total, pagable.GetPage(), pagable.GetSize())
}

func (c *ParcelUserController) mapTargetByLocation(ctx context.Context, req dto.ResolveMapTargetByLocationHTTPRequest) (*httpapi.Result, error) {
	if err := c.requireParcels(); err != nil {
		return nil, err
	}
	if req.Latitude == 0 && req.Longitude == 0 {
		return nil, _errors.ReturnError(400, "latitude and longitude are required")
	}
	result, err := c.parcels.ResolveMapTargetByLocation(ctx, dto.ResolveMapTargetRequestDTO{
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Radius:    req.Radius,
	})
	if err != nil {
		if err.Error() == "map target not found" {
			return nil, _errors.ReturnError(404, err.Error())
		}
		return nil, err
	}
	resp := &tqdpb.ResolveMapTargetByLocationResponse{
		TargetType: result.TargetType,
		TargetId:   result.TargetID,
	}
	if result.Parcel != nil {
		resp.Parcel = c.mapper.ToProtoParcelInfoDetail(result.Parcel)
	}
	if result.Region != nil {
		resp.Region = c.mapper.RegionInfoToProto(result.Region)
	}
	return protoResult(resp)
}

func (c *ParcelUserController) mapTargetByPoi(ctx context.Context, req dto.ResolveMapTargetByPoiHTTPRequest) (*httpapi.Result, error) {
	if err := c.requireParcels(); err != nil {
		return nil, err
	}
	if req.PoiID == 0 {
		return nil, _errors.ReturnError(400, "poi_id is required")
	}
	result, err := c.parcels.ResolveMapTargetByPoi(ctx, req.PoiID)
	if err != nil {
		msg := err.Error()
		if msg == "map target not found" || msg == "poi not found or missing coordinates" || msg == "poi id is required" {
			return nil, _errors.ReturnError(404, msg)
		}
		return nil, err
	}
	resp := &tqdpb.ResolveMapTargetByLocationResponse{
		TargetType: result.TargetType,
		TargetId:   result.TargetID,
	}
	if result.Parcel != nil {
		resp.Parcel = c.mapper.ToProtoParcelInfoDetail(result.Parcel)
	}
	if result.Region != nil {
		resp.Region = c.mapper.RegionInfoToProto(result.Region)
	}
	return protoResult(resp)
}

func (c *ParcelUserController) layerHistory(ctx context.Context, req dto.LayerHistoryHTTPRequest) (*httpapi.Result, error) {
	if err := c.requireParcels(); err != nil {
		return nil, err
	}
	lineage, err := c.parcels.GetLayerHistory(ctx, req.LayerID)
	if err != nil {
		return nil, err
	}
	return jsonOK(lineage)
}

func (c *ParcelUserController) addLegalDocs(ctx context.Context, req dto.LayerLegalDocsHTTPRequest) (*httpapi.Result, error) {
	if c.layerLegal == nil {
		return nil, _errors.ReturnError(501, "layer legal service not configured")
	}
	resp, err := c.layerLegal.AddLegalDocs(ctx, req.LayerID, req.LegalDocs)
	if err != nil {
		return nil, err
	}
	return jsonOK(resp)
}

func (c *ParcelUserController) listLegalDocs(ctx context.Context, req dto.LayerLegalDocsListHTTPRequest) (*httpapi.Result, error) {
	if c.layerLegal == nil {
		return nil, _errors.ReturnError(501, "layer legal service not configured")
	}
	resp, err := c.layerLegal.GetLegalDocsByLayer(ctx, req.LayerID)
	if err != nil {
		return nil, err
	}
	return jsonOK(resp)
}

func (c *ParcelUserController) deleteLegalDoc(ctx context.Context, req dto.LayerLegalDocIDHTTPRequest) (*httpapi.Result, error) {
	if c.layerLegal == nil {
		return nil, _errors.ReturnError(501, "layer legal service not configured")
	}
	if err := c.layerLegal.DeleteLegalDoc(ctx, req.ID); err != nil {
		return nil, err
	}
	return jsonOK(map[string]any{"success": true})
}
