package controller

import (
	"context"
	"net/http"
	"strings"
	"time"

	_errors "common/errors"
	"common/httpapi"
	"mdn4/modules/tqd/dto"
	domain "mdn4/modules/tqd/entity/discovery"
	"mdn4/modules/tqd/mapper"
	"mdn4/modules/tqd/provider"
	discoveryapp "mdn4/modules/tqd/services/discovery/application"
)

// DiscoveryController owns /v2/tqd/discovery/* routes.
type DiscoveryController struct {
	discovery *discoveryapp.Service
	mapper    *mapper.DiscoveryMapper
	analytics provider.SearchAnalytics
}

func NewDiscoveryController(discovery *discoveryapp.Service, analytics provider.SearchAnalytics) *DiscoveryController {
	if analytics == nil {
		analytics = provider.NoopSearchAnalytics{}
	}
	return &DiscoveryController{
		discovery: discovery,
		mapper:    mapper.NewDiscoveryMapper(),
		analytics: analytics,
	}
}

func (c *DiscoveryController) RegisterHTTP(mux *http.ServeMux) {
	httpapi.POST(mux, "/v2/tqd/discovery/identify", c.identify)
	httpapi.GET(mux, "/v2/tqd/discovery/search", c.search)
	httpapi.GET(mux, "/v2/tqd/discovery/entities/{key}", c.getEntity)
}

func (c *DiscoveryController) require() error {
	if c.discovery == nil {
		return _errors.ReturnError(501, "discovery service not configured")
	}
	return nil
}

func (c *DiscoveryController) identify(ctx context.Context, req dto.DiscoveryIdentifyHTTPRequest) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	if req.Point == nil {
		return nil, _errors.ReturnError(400, "point is required")
	}
	domainReq := domain.IdentifyRequest{
		Point: domain.Point{
			Latitude:  req.Point.Latitude,
			Longitude: req.Point.Longitude,
		},
	}
	if req.Context != nil {
		domainReq.Context = domain.MapContext{
			Zoom:           req.Context.Zoom,
			MapMode:        req.Context.MapMode,
			ActiveLayerIDs: req.Context.ActiveLayerIDs,
		}
		for _, k := range req.Context.PreferredKinds {
			domainReq.Context.PreferredKinds = append(domainReq.Context.PreferredKinds, domain.EntityKind(k))
		}
		if req.Context.Viewport != nil {
			domainReq.Context.Viewport = &domain.Bounds{
				MinLongitude: req.Context.Viewport.MinLon,
				MinLatitude:  req.Context.Viewport.MinLat,
				MaxLongitude: req.Context.Viewport.MaxLon,
				MaxLatitude:  req.Context.Viewport.MaxLat,
			}
		}
	}
	if req.Options != nil {
		domainReq.Options = domain.IdentifyOptions{
			ToleranceMeters:        req.Options.ToleranceMeters,
			Limit:                  int(req.Options.Limit),
			IncludeGeometryPreview: req.Options.IncludeGeometryPreview,
		}
	}
	resp, err := c.discovery.Identify(ctx, domainReq)
	if err != nil {
		return nil, err
	}
	proto, err := c.mapper.ToProtoIdentifyResponse(resp)
	if err != nil {
		return nil, err
	}
	return protoResult(proto)
}

func (c *DiscoveryController) search(ctx context.Context, req dto.DiscoverySearchHTTPRequest) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	q := strings.TrimSpace(req.Q)
	if q == "" {
		q = strings.TrimSpace(req.Text)
	}
	domainReq := domain.SearchRequest{
		Query:  q,
		Cursor: req.Cursor,
		Limit:  int(req.Limit),
	}
	if kinds := strings.TrimSpace(req.Kinds); kinds != "" {
		for _, part := range strings.Split(kinds, ",") {
			part = strings.TrimSpace(part)
			if part != "" {
				domainReq.EntityKinds = append(domainReq.EntityKinds, domain.EntityKind(part))
			}
		}
	}
	start := time.Now()
	resp, err := c.discovery.Search(ctx, domainReq)
	durationMs := time.Since(start).Milliseconds()
	if q != "" && c.analytics != nil {
		c.analytics.TrackKeyword(q, "tqd")
		c.analytics.TrackLatency("GET", "/v2/tqd/discovery/search", "tqd", durationMs)
	}
	if err != nil {
		return nil, err
	}
	proto, err := c.mapper.ToProtoSearchResponse(resp)
	if err != nil {
		return nil, err
	}
	return protoResult(proto)
}

func (c *DiscoveryController) getEntity(ctx context.Context, req dto.DiscoveryEntityKeyHTTPRequest) (*httpapi.Result, error) {
	if err := c.require(); err != nil {
		return nil, err
	}
	if strings.TrimSpace(req.Key) == "" {
		return nil, _errors.ReturnError(400, "key is required")
	}
	entity, err := c.discovery.GetEntity(ctx, req.Key)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, _errors.ReturnError(404, "entity not found")
	}
	proto, err := c.mapper.ToProtoEntityCandidate(entity)
	if err != nil {
		return nil, err
	}
	return protoResult(proto)
}
