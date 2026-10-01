package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	_dto "common/domain/dto"
	_errors "common/errors"
	_utils "common/utils"
	tqdpb "pb/types/tqd"
	"userv4/internal/httpapi"
	"userv4/modules/tqd/dto"
	relateddomain "userv4/modules/tqd/entity/relatedentity"
	"userv4/modules/tqd/mapper"
	publiccontentpg "userv4/modules/tqd/postgres/publiccontent"
	"userv4/modules/tqd/provider"
	"userv4/modules/tqd/services"
	planningapp "userv4/modules/tqd/services/planningclient/application"
	relatedapp "userv4/modules/tqd/services/relatedentity/application"
)

var qhConfigNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// PublicExtendedController owns remaining /v2/tqd/public/* and /v2/tqd/qh/public/* routes.
type PublicExtendedController struct {
	parcels       services.ParcelUsecase
	locationUC    *services.LocationUsecase
	legend        services.QHLayerLegendUsecase
	poi           services.PoiUsecase
	planning      *planningapp.Service
	related       *relatedapp.Service
	publicContent *publiccontentpg.Repository
	crmSeo        provider.CrmSeoProvider
	analytics     provider.SearchAnalytics
	parcelMapper  *mapper.ParcelMapper

	processedConfigs map[string]bool
	mu               sync.Mutex
	domainAPI        string
	domainFile       string
	domainPMTile     string
	domainStyle      string
}

func NewPublicExtendedController(
	parcels services.ParcelUsecase,
	locationUC *services.LocationUsecase,
	legend services.QHLayerLegendUsecase,
	poi services.PoiUsecase,
	planning *planningapp.Service,
	related *relatedapp.Service,
	publicContent *publiccontentpg.Repository,
	crmSeo provider.CrmSeoProvider,
	analytics provider.SearchAnalytics,
) *PublicExtendedController {
	if analytics == nil {
		analytics = provider.NoopSearchAnalytics{}
	}
	return &PublicExtendedController{
		parcels:          parcels,
		locationUC:       locationUC,
		legend:           legend,
		poi:              poi,
		planning:         planning,
		related:          related,
		publicContent:    publicContent,
		crmSeo:           crmSeo,
		analytics:        analytics,
		parcelMapper:     mapper.NewParcelMapper(),
		processedConfigs: map[string]bool{},
		domainAPI:        strings.TrimRight(os.Getenv("TQD_DOMAIN_API"), "/"),
		domainFile:       os.Getenv("TQD_DOMAIN_FILE"),
		domainPMTile:     os.Getenv("TQD_DOMAIN_PMTILE"),
		domainStyle:      os.Getenv("TQD_DOMAIN_STYLE"),
	}
}

func (c *PublicExtendedController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.GET(mux, "/v2/tqd/public/maps", c.getMaps)
	httpapi.POST(mux, "/v2/tqd/public/seos/default", c.createDefaultSeo)
	httpapi.GET(mux, "/v2/tqd/public/administrative-units/{value}/projection", c.adminUnitProjection)
	httpapi.GET(mux, "/v2/tqd/public/parcels/{parcel_id}/quick-view", c.parcelQuickView)
	httpapi.GET(mux, "/v2/tqd/public/planning-projects/{identity}/projection", c.planningProjection)
	httpapi.GET(mux, "/v2/tqd/public/planning-projects/{identity}", c.planningProjection)
	httpapi.GET(mux, "/v2/tqd/public/entities/{key}/related", c.relatedEntities)
	httpapi.GET(mux, "/v2/tqd/public/pois", c.listPois)
	httpapi.POST(mux, "/v2/tqd/public/qh/checking-polygon", c.checkingPolygon)
	httpapi.GET(mux, "/v2/tqd/public/parcels/search", c.searchPublicParcels)
	httpapi.GET(mux, "/v2/tqd/public/client/search", c.clientSearch)

	httpapi.GET(mux, "/v2/tqd/qh/public/json/{config}", c.qhConfig)
	httpapi.GET(mux, "/v2/tqd/qh/public/font/{fontstack}/{range}", c.qhFont)
	httpapi.GET(mux, "/v2/tqd/qh/public/layer-legends/all", c.legendsAll)
	httpapi.GET(mux, "/v2/tqd/qh/public/parcels/{parcelId}/quick-layers", c.parcelQuickLayers)
}

func (c *PublicExtendedController) getMaps(ctx context.Context, _ dto.PublicMapsHTTPRequest) (*httpapi.Result, error) {
	return jsonOK([]*tqdpb.MapSource{
		{Name: "Giao thông", Source: c.domainAPI + "/v2/tqd/qh/public/json/style", Type: "style"},
		{Name: "Vệ tinh", Source: "https://mt1.google.com/vt/lyrs=s&x={x}&y={y}&z={z}", Type: "raster"},
	})
}

func (c *PublicExtendedController) createDefaultSeo(ctx context.Context, req dto.CreateDefaultSeoHTTPRequest) (*httpapi.Result, error) {
	if req.RefType == 0 || req.RefID == 0 {
		return nil, _errors.ReturnError(400, "refType và refId là bắt buộc")
	}
	if c.crmSeo == nil {
		return nil, _errors.ReturnError(501, "crm seo provider chưa sẵn sàng")
	}
	result, err := c.crmSeo.CreateOrGetDefaultSeoByRef(ctx, req.RefType, req.RefID)
	if err != nil {
		return nil, err
	}
	return jsonOK(result)
}

func (c *PublicExtendedController) adminUnitProjection(ctx context.Context, req dto.PublicProjectionValueHTTPRequest) (*httpapi.Result, error) {
	if c.publicContent == nil {
		return nil, _errors.ReturnError(501, "projection repository unavailable")
	}
	identity := strings.TrimSpace(req.Value)
	if identity == "" {
		return nil, _errors.ReturnError(400, "administrative unit identity is required")
	}
	projection, err := c.publicContent.GetAdministrativeUnitProjection(ctx, identity)
	if err != nil {
		if publiccontentpg.IsNotFound(err) {
			return nil, _errors.ReturnError(404, "administrative unit was not found")
		}
		return nil, err
	}
	return jsonOK(projection)
}

func (c *PublicExtendedController) parcelQuickView(ctx context.Context, req dto.PublicParcelQuickViewHTTPRequest) (*httpapi.Result, error) {
	if c.planning == nil {
		return nil, _errors.ReturnError(501, "projection service unavailable")
	}
	view, err := c.planning.GetParcelQuickView(ctx, req.ParcelID)
	if err != nil {
		return nil, err
	}
	return jsonOK(view)
}

func (c *PublicExtendedController) planningProjection(ctx context.Context, req dto.PublicPlanningIdentityHTTPRequest) (*httpapi.Result, error) {
	if c.planning == nil {
		return nil, _errors.ReturnError(501, "projection service unavailable")
	}
	identity := strings.TrimSpace(req.Identity)
	if identity == "" {
		return nil, _errors.ReturnError(400, "identity is required")
	}
	projection, err := c.planning.GetPlanningProjectProjection(ctx, identity)
	if err != nil {
		return nil, err
	}
	return jsonOK(projection)
}

func (c *PublicExtendedController) relatedEntities(ctx context.Context, req dto.PublicRelatedEntityHTTPRequest) (*httpapi.Result, error) {
	if c.related == nil {
		return nil, _errors.ReturnError(501, "related entity service unavailable")
	}
	resp, err := c.related.GetRelatedEntities(ctx, relateddomain.Request{Key: req.Key})
	if err != nil {
		return nil, err
	}
	return jsonOK(resp)
}

func (c *PublicExtendedController) listPois(ctx context.Context, req dto.PublicPoisHTTPRequest) (*httpapi.Result, error) {
	if c.poi == nil {
		return nil, _errors.ReturnError(501, "poi service not configured")
	}
	search := strings.TrimSpace(req.Search)
	listReq := &dto.PublicPoiListRequest{
		Pagable: _dto.Pagable{Page: req.Page, Size: req.Size},
		Search:  search,
		Mode:    req.Mode,
	}
	listReq.Normalize()
	start := time.Now()
	rows, err := c.poi.ListPublic(ctx, listReq)
	durationMs := time.Since(start).Milliseconds()
	if c.analytics != nil {
		c.analytics.TrackLatency("GET", "/v2/tqd/public/pois", "tqd", durationMs)
		if search != "" {
			c.analytics.TrackKeyword(search, "tqd")
		}
	}
	if err != nil {
		return nil, err
	}
	return jsonOK(rows)
}

func (c *PublicExtendedController) checkingPolygon(ctx context.Context, req dto.FindByPolygonHTTPRequest) (*httpapi.Result, error) {
	if c.parcels == nil {
		return nil, _errors.ReturnError(501, "parcel service not configured")
	}
	req.MaxAreaHa = 100
	points, filter, err := polygonRequestParts(req)
	if err != nil {
		return nil, err
	}
	pagable := _dto.Pagable{Page: req.Page, Size: req.Size}
	pagable.Normalize()
	result, err := c.parcels.CheckingPolygon(
		ctx, points, req.Intersect, req.MaxAreaHa,
		int(pagable.GetPage()), int(pagable.GetSize()),
		req.Z, filter, req.IncludeAnalysis,
	)
	if err != nil {
		return nil, err
	}
	return protoResult(c.parcelMapper.PolygonMapTargetsToProto(result))
}

func (c *PublicExtendedController) searchPublicParcels(ctx context.Context, req dto.PublicParcelsSearchHTTPRequest) (*httpapi.Result, error) {
	if c.parcels == nil {
		return nil, _errors.ReturnError(501, "parcel service not configured")
	}
	pagable := _dto.Pagable{Page: req.Page, Size: req.Size}
	pagable.Normalize()
	rows, total, err := c.parcels.SearchPublicParcels(ctx, int(pagable.GetPage()), int(pagable.GetSize()))
	if err != nil {
		return nil, err
	}
	items := make([]*tqdpb.ParcelInfo, 0, len(rows))
	for _, p := range rows {
		items = append(items, c.parcelMapper.ToProtoParcelInfo(p))
	}
	return protoListResult(items, total, pagable.GetPage(), pagable.GetSize())
}

func (c *PublicExtendedController) clientSearch(ctx context.Context, req dto.PublicClientSearchHTTPRequest) (*httpapi.Result, error) {
	if c.locationUC == nil {
		return nil, _errors.ReturnError(501, "location service not configured")
	}
	q := strings.TrimSpace(req.Text)
	if q == "" {
		q = strings.TrimSpace(req.Q)
	}
	cleaned := _utils.CleanSearchText(q)
	start := time.Now()
	rows, err := c.locationUC.SearchTxtClient(ctx, cleaned, 20)
	durationMs := time.Since(start).Milliseconds()
	if c.analytics != nil {
		c.analytics.TrackLatency("GET", "/v2/tqd/public/client/search", "tqd", durationMs)
		if strings.TrimSpace(cleaned) != "" {
			c.analytics.TrackKeyword(cleaned, "tqd")
		}
	}
	if err != nil {
		return nil, err
	}
	items := make([]map[string]any, 0, len(rows))
	for _, item := range rows {
		row := map[string]any{
			"address":  item.Address,
			"title":    item.Title,
			"poiId":    item.PoiId,
			"parcelId": item.ParcelId,
		}
		if lat, lng, err := item.Geom.GetLatLng(); err == nil {
			row["latitude"] = lat
			row["longitude"] = lng
		}
		items = append(items, row)
	}
	return jsonOK(items)
}

func (c *PublicExtendedController) qhConfig(ctx context.Context, req dto.QhPublicConfigHTTPRequest) (*httpapi.Result, error) {
	if !qhConfigNamePattern.MatchString(req.Config) {
		return nil, _errors.ReturnError(400, "invalid config name")
	}
	jsonPath, err := c.ensureQhConfigProcessed(req.Config)
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, _errors.ReturnError(404, fmt.Sprintf("config file not found: %v", err))
	}
	var parsed any
	if err := json.Unmarshal(data, &parsed); err != nil {
		return nil, err
	}
	return jsonOK(parsed)
}

func (c *PublicExtendedController) ensureQhConfigProcessed(configName string) (string, error) {
	jsonPath := fmt.Sprintf("files/json/%s.json", configName)
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.processedConfigs[configName] {
		if _, err := os.Stat(jsonPath); err == nil {
			return jsonPath, nil
		}
		delete(c.processedConfigs, configName)
	}
	templatePath := fmt.Sprintf("template/%s.json", configName)
	templateData, err := os.ReadFile(templatePath)
	if err != nil {
		if _, statErr := os.Stat(jsonPath); statErr == nil {
			return jsonPath, nil
		}
		return "", _errors.ReturnError(404, fmt.Sprintf("template file not found: %v", err))
	}
	content := string(templateData)
	content = strings.ReplaceAll(content, "{{domain_api}}", c.domainAPI)
	content = strings.ReplaceAll(content, "{{domain_file}}", c.domainFile)
	content = strings.ReplaceAll(content, "{{domain_pmtile}}", c.domainPMTile)
	content = strings.ReplaceAll(content, "{{domain_style}}", c.domainStyle)
	if err := os.MkdirAll("files/json", 0755); err != nil {
		return "", err
	}
	if err := os.WriteFile(jsonPath, []byte(content), 0644); err != nil {
		return "", err
	}
	c.processedConfigs[configName] = true
	return jsonPath, nil
}

func (c *PublicExtendedController) qhFont(ctx context.Context, req dto.QhPublicFontHTTPRequest) (*httpapi.Result, error) {
	path := fmt.Sprintf("files/fonts/%s/%s", req.Fontstack, req.Range)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, _errors.ReturnError(404, fmt.Sprintf("font file not found: %v", err))
	}
	return jsonOK(map[string]any{"code": 8501, "data": data})
}

func (c *PublicExtendedController) legendsAll(ctx context.Context, _ dto.QhPublicLegendsAllHTTPRequest) (*httpapi.Result, error) {
	if c.legend == nil {
		return nil, _errors.ReturnError(501, "legend service not configured")
	}
	rows, total, err := c.legend.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return jsonListResult(rows, total, 1, uint32(total))
}

func (c *PublicExtendedController) parcelQuickLayers(ctx context.Context, req dto.QhPublicParcelQuickLayersHTTPRequest) (*httpapi.Result, error) {
	if c.parcels == nil {
		return nil, _errors.ReturnError(501, "parcel service not configured")
	}
	rows, err := c.parcels.GetParcelQuickInfo(ctx, req.ParcelID)
	if err != nil {
		return nil, err
	}
	return protoResult(c.parcelMapper.BuildLayerQuickInfoFromRows(rows))
}
