package handler_grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	_middleware "common/middleware"
	"crm/config"
	catalogapp "crm/internal/modules/contentcatalog/application"
	catalogdomain "crm/internal/modules/contentcatalog/domain"
	"crm/internal/usecase"

	crmpb "pb/types/crm"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	protojson "google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
	"gorm.io/gorm"
)

// PublicContentGrpcHandler is CRM's transport-neutral SEO/content boundary.
// HTTP status, caching and ETag values are returned as gRPC metadata so only
// gateway-service has to understand HTTP.
type PublicContentGrpcHandler struct {
	crmpb.UnimplementedCRMPublicContentServiceServer
	seoPublic *usecase.SeoPublicUsecase
	sitemap   *usecase.SeoSitemapUsecase
	catalog   *catalogapp.Service
}

func NewPublicContentGrpcHandler(seoPublic *usecase.SeoPublicUsecase, sitemap *usecase.SeoSitemapUsecase, catalog *catalogapp.Service) *PublicContentGrpcHandler {
	return &PublicContentGrpcHandler{seoPublic: seoPublic, sitemap: sitemap, catalog: catalog}
}

func jsonBytes(value any) (*wrapperspb.BytesValue, error) {
	payload, err := json.Marshal(value)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "encode response: %v", err)
	}
	return wrapperspb.Bytes(payload), nil
}

func authenticatedProfileID(ctx context.Context) (uint64, error) {
	principal, err := _middleware.PrincipalFromContext(ctx)
	if err != nil || principal == nil || principal.ProfileId == 0 {
		return 0, status.Error(codes.Unauthenticated, "authentication is required")
	}
	return principal.ProfileId, nil
}

func requestMap(req proto.Message) map[string]any {
	if req == nil {
		return map[string]any{}
	}
	payload, err := protojson.MarshalOptions{UseProtoNames: false}.Marshal(req)
	if err != nil {
		return map[string]any{}
	}
	result := map[string]any{}
	if err := json.Unmarshal(payload, &result); err != nil {
		return map[string]any{}
	}
	return result
}
func stringValue(m map[string]any, keys ...string) string {
	for _, key := range keys {
		if raw, ok := m[key]; ok {
			if value := strings.TrimSpace(fmt.Sprint(raw)); value != "" && value != "<nil>" {
				return value
			}
		}
	}
	return ""
}
func intValue(m map[string]any, key string, fallback int) int {
	raw, ok := m[key]
	if !ok {
		return fallback
	}
	switch value := raw.(type) {
	case float64:
		return int(value)
	case string:
		if parsed, err := strconv.Atoi(strings.TrimSpace(value)); err == nil {
			return parsed
		}
	}
	return fallback
}
func stringSlice(m map[string]any, key string) []string {
	raw, ok := m[key]
	if !ok {
		return nil
	}
	result := []string{}
	switch value := raw.(type) {
	case []any:
		for _, item := range value {
			if s := strings.TrimSpace(fmt.Sprint(item)); s != "" {
				result = append(result, s)
			}
		}
	case string:
		for _, item := range strings.Split(value, ",") {
			if s := strings.TrimSpace(item); s != "" {
				result = append(result, s)
			}
		}
	}
	return result
}
func incomingMetadataValue(ctx context.Context, key string) string {
	incoming, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	values := incoming.Get(strings.ToLower(strings.TrimSpace(key)))
	if len(values) == 0 {
		return ""
	}
	return strings.TrimSpace(values[0])
}

func requestETagMatches(ctx context.Context, etag string) bool {
	expected := normalizeETag(etag)
	provided := incomingMetadataValue(ctx, "if-none-match")
	if expected == "" || provided == "" {
		return false
	}
	for _, candidate := range strings.Split(provided, ",") {
		candidate = strings.TrimSpace(candidate)
		if candidate == "*" || candidate == expected {
			return true
		}
	}
	return false
}

func setHTTPMetadata(ctx context.Context, statusCode int, contentType, cacheControl, etag, robots string) {
	if statusCode >= 200 && statusCode < 300 && requestETagMatches(ctx, etag) {
		statusCode = 304
	}
	pairs := []string{"x-http-status", strconv.Itoa(statusCode)}
	if contentType != "" {
		pairs = append(pairs, "content-type", contentType)
	}
	if cacheControl != "" {
		pairs = append(pairs, "cache-control", cacheControl)
	}
	if etag != "" {
		pairs = append(pairs, "etag", normalizeETag(etag))
	}
	if robots != "" {
		pairs = append(pairs, "x-robots-tag", robots)
	}
	pairs = append(pairs, "x-content-type-options", "nosniff")
	_ = grpc.SetHeader(ctx, metadata.Pairs(pairs...))
}
func normalizeETag(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if strings.HasPrefix(value, `W/"`) || strings.HasPrefix(value, `"`) {
		return value
	}
	return `"` + strings.Trim(value, `"`) + `"`
}

func (h *PublicContentGrpcHandler) GetPublicSeoDocument(ctx context.Context, req *crmpb.SeoDocumentRequest) (*wrapperspb.BytesValue, error) {
	m := requestMap(req)
	response, err := h.seoPublic.GetPublicSeoDocument(ctx, stringValue(m, "module"), stringValue(m, "slug"))
	if err != nil {
		return nil, err
	}
	robots := ""
	if response.Document != nil && !response.Document.SEO.Indexable {
		robots = "noindex, follow"
	}
	setHTTPMetadata(ctx, response.StatusCode, "application/json; charset=utf-8", response.CacheControl, response.ETag, robots)
	return jsonBytes(response.Document)
}
func (h *PublicContentGrpcHandler) GetRenderedSeoPage(ctx context.Context, req *crmpb.RenderedSeoPageRequest) (*wrapperspb.BytesValue, error) {
	response, err := h.seoPublic.GetPublicSeoPageBySlug(ctx, strings.TrimSpace(req.GetPath()))
	if err != nil {
		return nil, err
	}
	robots := ""
	if response.NoIndex {
		robots = "noindex, follow"
	}
	setHTTPMetadata(ctx, response.StatusCode, response.ContentType, response.CacheControl, response.ETag, robots)
	return wrapperspb.Bytes([]byte(response.HTML)), nil
}
func (h *PublicContentGrpcHandler) GetSitemap(ctx context.Context, _ *emptypb.Empty) (*wrapperspb.BytesValue, error) {
	response, err := h.sitemap.BuildSitemapXML(ctx, config.AppProperties.Seo.PublicBaseURL)
	if err != nil {
		setHTTPMetadata(ctx, 503, "application/xml; charset=utf-8", "no-store, max-age=0", "", "")
		_ = grpc.SetHeader(ctx, metadata.Pairs("retry-after", "60"))
		return wrapperspb.Bytes([]byte(`<?xml version="1.0" encoding="UTF-8"?><urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"></urlset>`)), nil
	}
	setHTTPMetadata(ctx, 200, response.ContentType, "public, max-age=300, stale-while-revalidate=3600", "", "")
	return wrapperspb.Bytes([]byte(response.XML)), nil
}
func (h *PublicContentGrpcHandler) SearchContent(ctx context.Context, req *crmpb.ContentQueryRequest) (*wrapperspb.BytesValue, error) {
	m := requestMap(req)
	result, err := h.catalog.SearchPublishedContent(ctx, stringValue(m, "q", "query", "keyword"), stringSlice(m, "kinds"), intValue(m, "limit", 20))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "search public content: %v", err)
	}
	setHTTPMetadata(ctx, 200, "application/json; charset=utf-8", "no-store", "", "")
	return jsonBytes(result)
}
func (h *PublicContentGrpcHandler) ListPlanningNews(ctx context.Context, req *crmpb.PlanningNewsQueryRequest) (*wrapperspb.BytesValue, error) {
	if req == nil {
		req = &crmpb.PlanningNewsQueryRequest{}
	}

	filter := catalogdomain.PlanningNewsListFilter{Page: int(req.GetPage()), Size: int(req.GetSize())}
	result, err := h.catalog.ListPlanningNews(ctx, filter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list planning news: %v", err)
	}
	setHTTPMetadata(ctx, 200, "application/json; charset=utf-8", "public, max-age=60, stale-while-revalidate=300", "", "")
	return jsonBytes(result)
}
func (h *PublicContentGrpcHandler) GetPlanningNews(ctx context.Context, req *crmpb.PlanningNewsSlugRequest) (*wrapperspb.BytesValue, error) {
	result, err := h.catalog.GetPlanningNews(ctx, req.GetSlug())
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "planning news was not found")
		}
		return nil, status.Errorf(codes.Internal, "get planning news: %v", err)
	}
	setHTTPMetadata(ctx, 200, "application/json; charset=utf-8", "public, max-age=60, stale-while-revalidate=300", "", "")
	return jsonBytes(result)
}

func (h *PublicContentGrpcHandler) ListSavedPlanningNews(ctx context.Context, req *crmpb.SavedPlanningNewsQueryRequest) (*wrapperspb.BytesValue, error) {
	userID, err := authenticatedProfileID(ctx)
	if err != nil {
		return nil, err
	}
	if req == nil {
		req = &crmpb.SavedPlanningNewsQueryRequest{}
	}
	result, err := h.catalog.ListSavedPlanningNews(ctx, userID, catalogdomain.SavedPlanningNewsListFilter{
		Page: int(req.GetPage()),
		Size: int(req.GetSize()),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list saved planning news: %v", err)
	}
	setHTTPMetadata(ctx, 200, "application/json; charset=utf-8", "private, no-store", "", "")
	return jsonBytes(result)
}

func (h *PublicContentGrpcHandler) SavePlanningNews(ctx context.Context, req *crmpb.PlanningNewsSaveRequest) (*wrapperspb.BytesValue, error) {
	userID, err := authenticatedProfileID(ctx)
	if err != nil {
		return nil, err
	}
	if req == nil || req.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "planning news id is required")
	}
	item, err := h.catalog.SavePlanningNews(ctx, userID, req.GetId())
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "planning news was not found or is not public")
		}
		return nil, status.Errorf(codes.Internal, "save planning news: %v", err)
	}
	setHTTPMetadata(ctx, 200, "application/json; charset=utf-8", "private, no-store", "", "")
	return jsonBytes(map[string]any{"saved": true, "data": item})
}

func (h *PublicContentGrpcHandler) UnsavePlanningNews(ctx context.Context, req *crmpb.PlanningNewsSaveRequest) (*wrapperspb.BytesValue, error) {
	userID, err := authenticatedProfileID(ctx)
	if err != nil {
		return nil, err
	}
	if req == nil || req.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "planning news id is required")
	}
	if err := h.catalog.UnsavePlanningNews(ctx, userID, req.GetId()); err != nil {
		return nil, status.Errorf(codes.Internal, "unsave planning news: %v", err)
	}
	setHTTPMetadata(ctx, 200, "application/json; charset=utf-8", "private, no-store", "", "")
	return jsonBytes(map[string]any{"saved": false, "planningNewsId": req.GetId()})
}
