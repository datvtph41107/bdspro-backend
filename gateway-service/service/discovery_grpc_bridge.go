package service

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	_middleware "common/middleware"
	tqdpb "pb/types/tqd"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const discoveryRequestBodyLimit = 2 << 20 // 2 MiB

func hasPathPrefix(path, prefix string) bool {
	return path == prefix || strings.HasPrefix(path, prefix+"/")
}

func discoveryOutgoingContext(req *http.Request) context.Context {
	md := _middleware.InjectGrpcMetadataContextMiddleware(req.Context(), req)
	if existing, ok := metadata.FromOutgoingContext(req.Context()); ok {
		md = metadata.Join(existing, md)
	}
	return metadata.NewOutgoingContext(req.Context(), md)
}

func discoveryHTTPStatus(err error) int {
	switch status.Code(err) {
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists, codes.Aborted:
		return http.StatusConflict
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.Unavailable, codes.DeadlineExceeded:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}

func writeDiscoveryError(c *gin.Context, err error) {
	statusCode := discoveryHTTPStatus(err)
	code := status.Code(err).String()
	message := err.Error()
	if grpcStatus, ok := status.FromError(err); ok {
		message = grpcStatus.Message()
	}
	c.AbortWithStatusJSON(statusCode, gin.H{
		"code":    code,
		"message": message,
	})
}

func apiEntityKind(kind tqdpb.EntityKind) string {
	switch kind {
	case tqdpb.EntityKind_ENTITY_KIND_PARCEL:
		return "parcel"
	case tqdpb.EntityKind_ENTITY_KIND_PLANNING_REGION:
		return "planning_region"
	case tqdpb.EntityKind_ENTITY_KIND_ADMINISTRATIVE_UNIT:
		return "administrative_unit"
	case tqdpb.EntityKind_ENTITY_KIND_POI:
		return "poi"
	case tqdpb.EntityKind_ENTITY_KIND_PLANNING_PROJECT:
		return "planning_project"
	case tqdpb.EntityKind_ENTITY_KIND_LEGAL_DOCUMENT:
		return "legal_document"
	case tqdpb.EntityKind_ENTITY_KIND_PLANNING_MAP:
		return "planning_map"
	case tqdpb.EntityKind_ENTITY_KIND_MAP_LAYER:
		return "map_layer"
	case tqdpb.EntityKind_ENTITY_KIND_ANALYSIS_REPORT:
		return "analysis_report"
	case tqdpb.EntityKind_ENTITY_KIND_NEWS_ARTICLE:
		return "news_article"
	default:
		return ""
	}
}

func protoEntityKind(value string) (tqdpb.EntityKind, bool) {
	normalized := strings.ToLower(strings.TrimSpace(value))
	normalized = strings.TrimPrefix(normalized, "entity_kind_")
	switch normalized {
	case "parcel":
		return tqdpb.EntityKind_ENTITY_KIND_PARCEL, true
	case "planning_region":
		return tqdpb.EntityKind_ENTITY_KIND_PLANNING_REGION, true
	case "administrative_unit":
		return tqdpb.EntityKind_ENTITY_KIND_ADMINISTRATIVE_UNIT, true
	case "poi":
		return tqdpb.EntityKind_ENTITY_KIND_POI, true
	case "planning_project":
		return tqdpb.EntityKind_ENTITY_KIND_PLANNING_PROJECT, true
	case "legal_document":
		return tqdpb.EntityKind_ENTITY_KIND_LEGAL_DOCUMENT, true
	case "planning_map":
		return tqdpb.EntityKind_ENTITY_KIND_PLANNING_MAP, true
	case "map_layer":
		return tqdpb.EntityKind_ENTITY_KIND_MAP_LAYER, true
	case "analysis_report":
		return tqdpb.EntityKind_ENTITY_KIND_ANALYSIS_REPORT, true
	case "news_article":
		return tqdpb.EntityKind_ENTITY_KIND_NEWS_ARTICLE, true
	default:
		return tqdpb.EntityKind_ENTITY_KIND_UNSPECIFIED, false
	}
}

func apiPoint(value *tqdpb.DiscoveryPoint) map[string]any {
	if value == nil {
		return nil
	}
	return map[string]any{
		"latitude":  value.GetLatitude(),
		"longitude": value.GetLongitude(),
	}
}

func apiBounds(value *tqdpb.DiscoveryBounds) map[string]any {
	if value == nil {
		return nil
	}
	return map[string]any{
		"minLongitude": value.GetMinLongitude(),
		"minLatitude":  value.GetMinLatitude(),
		"maxLongitude": value.GetMaxLongitude(),
		"maxLatitude":  value.GetMaxLatitude(),
	}
}

func apiCandidate(value *tqdpb.EntityCandidate) map[string]any {
	if value == nil {
		return nil
	}
	result := map[string]any{}
	if entity := value.GetEntity(); entity != nil {
		result["entity"] = map[string]any{
			"kind": apiEntityKind(entity.GetKind()),
			"id":   entity.GetId(),
			"key":  entity.GetKey(),
		}
	}
	if presentation := value.GetPresentation(); presentation != nil {
		result["presentation"] = map[string]any{
			"title":       presentation.GetTitle(),
			"subtitle":    presentation.GetSubtitle(),
			"description": presentation.GetDescription(),
			"status":      presentation.GetStatus(),
			"badges":      presentation.GetBadges(),
		}
	}
	if spatial := value.GetSpatial(); spatial != nil {
		spatialPayload := map[string]any{
			"geometryType": spatial.GetGeometryType(),
			"geometryRef":  spatial.GetGeometryRef(),
		}
		if centroid := apiPoint(spatial.GetCentroid()); centroid != nil {
			spatialPayload["centroid"] = centroid
		}
		if bounds := apiBounds(spatial.GetBounds()); bounds != nil {
			spatialPayload["bounds"] = bounds
		}
		if preview := spatial.GetGeometryPreview(); preview != nil {
			spatialPayload["geometryPreview"] = preview.AsMap()
		}
		result["spatial"] = spatialPayload
	}
	if location := value.GetLocation(); location != nil {
		result["location"] = map[string]any{
			"address":      location.GetAddress(),
			"province":     location.GetProvince(),
			"provinceCode": location.GetProvinceCode(),
			"ward":         location.GetWard(),
			"wardCode":     location.GetWardCode(),
		}
	}
	if source := value.GetSource(); source != nil {
		result["source"] = map[string]any{
			"system":      source.GetSystem(),
			"dataset":     source.GetDataset(),
			"authority":   source.GetAuthority(),
			"updatedAt":   source.GetUpdatedAt(),
			"dataQuality": source.GetDataQuality(),
		}
	}
	if capabilities := value.GetCapabilities(); capabilities != nil {
		result["capabilities"] = map[string]any{
			"canOpenDetail":   capabilities.GetCanOpenDetail(),
			"canFocusMap":     capabilities.GetCanFocusMap(),
			"canFollow":       capabilities.GetCanFollow(),
			"canCreateReport": capabilities.GetCanCreateReport(),
			"canCompare":      capabilities.GetCanCompare(),
			"canShare":        capabilities.GetCanShare(),
			"canDownload":     capabilities.GetCanDownload(),
		}
	}
	if links := value.GetLinks(); links != nil {
		result["links"] = map[string]any{
			"canonicalUrl": links.GetCanonicalUrl(),
			"detailUrl":    links.GetDetailUrl(),
			"geometryUrl":  links.GetGeometryUrl(),
		}
	}
	if match := value.GetMatch(); match != nil {
		result["match"] = map[string]any{
			"type":          match.GetType(),
			"matchedFields": match.GetMatchedFields(),
			"highlights":    match.GetHighlights(),
			"reasonCodes":   match.GetReasonCodes(),
		}
	}
	if rank := value.GetRank(); rank != nil {
		result["rank"] = map[string]any{
			"score":       rank.GetScore(),
			"reasonCodes": rank.GetReasonCodes(),
		}
	}
	if attributes := value.GetAttributes(); attributes != nil {
		result["attributes"] = attributes.AsMap()
	}
	return result
}

func apiCandidates(values []*tqdpb.EntityCandidate) []map[string]any {
	result := make([]map[string]any, 0, len(values))
	for _, value := range values {
		if candidate := apiCandidate(value); candidate != nil {
			result = append(result, candidate)
		}
	}
	return result
}

func writeIdentifyResponse(c *gin.Context, response *tqdpb.IdentifyAtPointResponse) {
	if response == nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"code": "EMPTY_DISCOVERY_RESPONSE", "message": "TQD Discovery returned an empty identify response"})
		return
	}
	payload := map[string]any{
		"requestId":  response.GetRequestId(),
		"candidates": apiCandidates(response.GetCandidates()),
		"partial":    response.GetPartial(),
		"warnings":   response.GetWarnings(),
	}
	if resolution := response.GetResolution(); resolution != nil {
		payload["resolution"] = map[string]any{
			"status":     resolution.GetStatus(),
			"primaryKey": resolution.GetPrimaryKey(),
			"confidence": resolution.GetConfidence(),
			"ambiguous":  resolution.GetAmbiguous(),
			"strategy":   resolution.GetStrategy(),
		}
	}
	if primary := apiCandidate(response.GetPrimary()); primary != nil {
		payload["primary"] = primary
	}
	c.JSON(http.StatusOK, payload)
	c.Abort()
}

func writeSearchResponse(c *gin.Context, response *tqdpb.SearchDiscoveryResponse) {
	if response == nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"code": "EMPTY_DISCOVERY_RESPONSE", "message": "TQD Discovery returned an empty search response"})
		return
	}
	groups := make([]map[string]any, 0, len(response.GetGroups()))
	for _, group := range response.GetGroups() {
		if group == nil {
			continue
		}
		groups = append(groups, map[string]any{
			"kind":  apiEntityKind(group.GetKind()),
			"label": group.GetLabel(),
			"total": group.GetTotal(),
			"items": apiCandidates(group.GetItems()),
		})
	}
	payload := map[string]any{
		"requestId":  response.GetRequestId(),
		"groups":     groups,
		"nextCursor": response.GetNextCursor(),
		"partial":    response.GetPartial(),
		"warnings":   response.GetWarnings(),
	}
	if interpretation := response.GetInterpretation(); interpretation != nil {
		payload["interpretation"] = map[string]any{
			"originalQuery":   interpretation.GetOriginalQuery(),
			"normalizedQuery": interpretation.GetNormalizedQuery(),
			"detectedIntents": interpretation.GetDetectedIntents(),
			"structuredTerms": interpretation.GetStructuredTerms(),
		}
	}
	c.JSON(http.StatusOK, payload)
	c.Abort()
}

func writeEntityResponse(c *gin.Context, response *tqdpb.EntityCandidate) {
	if response == nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"code": "NOT_FOUND", "message": "Discovery entity not found"})
		return
	}
	c.JSON(http.StatusOK, apiCandidate(response))
	c.Abort()
}

type identifyHTTPBody struct {
	Point struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"point"`
	Context *struct {
		Zoom           float64  `json:"zoom"`
		MapMode        string   `json:"mapMode"`
		ActiveLayerIDs []string `json:"activeLayerIds"`
		PreferredKinds []string `json:"preferredKinds"`
		Viewport       *struct {
			MinLongitude float64 `json:"minLongitude"`
			MinLatitude  float64 `json:"minLatitude"`
			MaxLongitude float64 `json:"maxLongitude"`
			MaxLatitude  float64 `json:"maxLatitude"`
		} `json:"viewport"`
	} `json:"context"`
	Options *struct {
		ToleranceMeters        float64 `json:"toleranceMeters"`
		Limit                  int32   `json:"limit"`
		IncludeGeometryPreview bool    `json:"includeGeometryPreview"`
	} `json:"options"`
}

func bodyAsIdentifyRequest(c *gin.Context) (*tqdpb.IdentifyAtPointRequest, error) {
	reader := http.MaxBytesReader(c.Writer, c.Request.Body, discoveryRequestBodyLimit)
	defer reader.Close()
	var body identifyHTTPBody
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&body); err != nil {
		if errors.Is(err, io.EOF) {
			return nil, errors.New("request body is required")
		}
		return nil, err
	}
	request := &tqdpb.IdentifyAtPointRequest{
		Point: &tqdpb.DiscoveryPoint{
			Latitude:  body.Point.Latitude,
			Longitude: body.Point.Longitude,
		},
	}
	if body.Context != nil {
		contextValue := &tqdpb.IdentifyContext{
			Zoom:           body.Context.Zoom,
			MapMode:        body.Context.MapMode,
			ActiveLayerIds: append([]string(nil), body.Context.ActiveLayerIDs...),
		}
		for _, value := range body.Context.PreferredKinds {
			if kind, ok := protoEntityKind(value); ok {
				contextValue.PreferredKinds = append(contextValue.PreferredKinds, kind)
			}
		}
		if body.Context.Viewport != nil {
			contextValue.Viewport = &tqdpb.DiscoveryBounds{
				MinLongitude: body.Context.Viewport.MinLongitude,
				MinLatitude:  body.Context.Viewport.MinLatitude,
				MaxLongitude: body.Context.Viewport.MaxLongitude,
				MaxLatitude:  body.Context.Viewport.MaxLatitude,
			}
		}
		request.Context = contextValue
	}
	if body.Options != nil {
		request.Options = &tqdpb.IdentifyOptions{
			ToleranceMeters:        body.Options.ToleranceMeters,
			Limit:                  body.Options.Limit,
			IncludeGeometryPreview: body.Options.IncludeGeometryPreview,
		}
	}
	return request, nil
}

func splitQueryList(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func searchRequestFromQuery(c *gin.Context) (*tqdpb.SearchDiscoveryRequest, error) {
	query := strings.TrimSpace(c.Query("q"))
	if query == "" {
		query = strings.TrimSpace(c.Query("query"))
	}
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil {
		return nil, errors.New("limit must be an integer")
	}
	request := &tqdpb.SearchDiscoveryRequest{
		Q:              query,
		Mode:           c.DefaultQuery("mode", "suggest"),
		MapMode:        strings.TrimSpace(c.Query("mapMode")),
		ActiveLayerIds: splitQueryList(c.Query("activeLayerIds")),
		Cursor:         strings.TrimSpace(c.Query("cursor")),
		Limit:          int32(limit),
	}
	for _, value := range splitQueryList(c.Query("kinds")) {
		kind, ok := protoEntityKind(value)
		if !ok {
			return nil, errors.New("unsupported entity kind: " + value)
		}
		request.Kinds = append(request.Kinds, kind)
	}
	return request, nil
}

func entityRequestFromPath(path string) (*tqdpb.GetDiscoveryEntityRequest, error) {
	const prefix = "/v2/tqd/discovery/entities/"
	if !strings.HasPrefix(path, prefix) {
		return nil, errors.New("invalid discovery entity path")
	}
	key, err := url.PathUnescape(strings.TrimPrefix(path, prefix))
	if err != nil {
		return nil, errors.New("invalid encoded canonical entity key")
	}
	key = strings.TrimSpace(key)
	if key == "" {
		return nil, errors.New("canonical entity key is required")
	}
	return &tqdpb.GetDiscoveryEntityRequest{Key: key}, nil
}

// DiscoveryOrGRPCHandler keeps Gin's single /v2/*any catch-all while routing
// the Discovery REST contract to the typed TQD Discovery gRPC client. Every
// other /v2 route continues to the normal grpc-gateway mux.
func DiscoveryOrGRPCHandler(grpcGateway http.Handler, discoveryClient tqdpb.DiscoveryServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		if !hasPathPrefix(path, "/v2/tqd/discovery") {
			// Gin NoRoute preset status=404. grpc-gateway thường Write() body
			// mà không WriteHeader(200) → client nhận HTTP 404 dù body code:0.
			// Reset trước khi forward; handler/error sẽ ghi đè khi cần.
			c.Status(http.StatusOK)
			grpcGateway.ServeHTTP(c.Writer, c.Request)
			return
		}
		if discoveryClient == nil {
			c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
				"code":    "TQD_DISCOVERY_CLIENT_UNAVAILABLE",
				"message": "TQD Discovery gRPC client is unavailable",
			})
			return
		}

		switch {
		case c.Request.Method == http.MethodPost && path == "/v2/tqd/discovery/identify":
			request, err := bodyAsIdentifyRequest(c)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": "INVALID_REQUEST", "message": err.Error()})
				return
			}
			ctx, cancel := context.WithTimeout(discoveryOutgoingContext(c.Request), 12*time.Second)
			defer cancel()
			response, err := discoveryClient.IdentifyAtPoint(ctx, request)
			if err != nil {
				writeDiscoveryError(c, err)
				return
			}
			writeIdentifyResponse(c, response)

		case c.Request.Method == http.MethodGet && path == "/v2/tqd/discovery/search":
			request, err := searchRequestFromQuery(c)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": "INVALID_REQUEST", "message": err.Error()})
				return
			}
			ctx, cancel := context.WithTimeout(discoveryOutgoingContext(c.Request), 12*time.Second)
			defer cancel()
			response, err := discoveryClient.SearchDiscovery(ctx, request)
			if err != nil {
				writeDiscoveryError(c, err)
				return
			}
			writeSearchResponse(c, response)

		case c.Request.Method == http.MethodGet && strings.HasPrefix(path, "/v2/tqd/discovery/entities/"):
			request, err := entityRequestFromPath(path)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"code": "INVALID_REQUEST", "message": err.Error()})
				return
			}
			ctx, cancel := context.WithTimeout(discoveryOutgoingContext(c.Request), 8*time.Second)
			defer cancel()
			response, err := discoveryClient.GetDiscoveryEntity(ctx, request)
			if err != nil {
				writeDiscoveryError(c, err)
				return
			}
			writeEntityResponse(c, response)

		default:
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"code":    "DISCOVERY_ROUTE_NOT_FOUND",
				"message": "Discovery route not found",
			})
		}
	}
}
