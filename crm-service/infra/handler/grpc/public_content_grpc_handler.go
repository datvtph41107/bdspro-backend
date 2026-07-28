package handler_grpc

import (
	"context"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"crm/config"
	"crm/internal/dto"
	catalogapp "crm/internal/modules/contentcatalog/application"
	catalogdomain "crm/internal/modules/contentcatalog/domain"
	"crm/internal/usecase"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	protojson "google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	structpb "google.golang.org/protobuf/types/known/structpb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
	"gorm.io/gorm"
	crmpb "pb/types/crm"
)

// PublicContentGrpcHandler is CRM's transport-neutral SEO/content boundary.
// HTTP status, caching and ETag values are returned as gRPC metadata so only
// gateway-service has to understand HTTP.
type PublicContentGrpcHandler struct {
	crmpb.UnimplementedCRMPublicContentServiceServer
	seoPublic             *usecase.SeoPublicUsecase
	sitemap               *usecase.SeoSitemapUsecase
	provisioning          *usecase.SeoAdministrativeUnitProvisioningUsecase
	catalog               *catalogapp.Service
	runtimeEvents         *usecase.SeoRuntimeEventUsecase
	runtimeEventIngestKey string
}

func NewPublicContentGrpcHandler(seoPublic *usecase.SeoPublicUsecase, sitemap *usecase.SeoSitemapUsecase, provisioning *usecase.SeoAdministrativeUnitProvisioningUsecase, catalog *catalogapp.Service, runtimeEvents *usecase.SeoRuntimeEventUsecase, runtimeEventIngestKey string) *PublicContentGrpcHandler {
	return &PublicContentGrpcHandler{seoPublic: seoPublic, sitemap: sitemap, provisioning: provisioning, catalog: catalog, runtimeEvents: runtimeEvents, runtimeEventIngestKey: strings.TrimSpace(runtimeEventIngestKey)}
}

func jsonBytes(value any) (*wrapperspb.BytesValue, error) {
	payload, err := json.Marshal(value)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "encode response: %v", err)
	}
	return wrapperspb.Bytes(payload), nil
}
func requestMap(req proto.Message) map[string]any {
	if req == nil {
		return map[string]any{}
	}
	if value, ok := req.(*structpb.Struct); ok {
		return value.AsMap()
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
func uintValue(m map[string]any, key string) uint32 {
	raw, ok := m[key]
	if !ok {
		return 0
	}
	switch value := raw.(type) {
	case float64:
		if value > 0 {
			return uint32(value)
		}
	case string:
		if parsed, err := strconv.ParseUint(strings.TrimSpace(value), 10, 32); err == nil {
			return uint32(parsed)
		}
	}
	return 0
}
func uint64Value(m map[string]any, key string) uint64 {
	raw, ok := m[key]
	if !ok {
		return 0
	}
	switch value := raw.(type) {
	case float64:
		if value > 0 {
			return uint64(value)
		}
	case string:
		if parsed, err := strconv.ParseUint(strings.TrimSpace(value), 10, 64); err == nil {
			return parsed
		}
	}
	return 0
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
func (h *PublicContentGrpcHandler) GetGovernanceSeoDocument(ctx context.Context, req *crmpb.SeoDocumentRequest) (*wrapperspb.BytesValue, error) {
	m := requestMap(req)
	response, err := h.seoPublic.GetSeoGovernanceDocument(ctx, stringValue(m, "module"), stringValue(m, "slug"))
	if err != nil {
		return nil, err
	}
	setHTTPMetadata(ctx, response.StatusCode, "application/json; charset=utf-8", "private, no-store", response.ETag, "")
	return jsonBytes(map[string]any{"code": response.StatusCode, "data": response.Document})
}
func (h *PublicContentGrpcHandler) GetOperatingSystemManifest(ctx context.Context, _ *emptypb.Empty) (*wrapperspb.BytesValue, error) {
	setHTTPMetadata(ctx, 200, "application/json; charset=utf-8", "private, no-store", "", "")
	return jsonBytes(map[string]any{"code": 200, "data": h.seoPublic.GetSeoOperatingSystemManifest()})
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
func (h *PublicContentGrpcHandler) ListContentDomains(ctx context.Context, _ *crmpb.ContentQueryRequest) (*wrapperspb.BytesValue, error) {
	items, err := h.catalog.ListDomains(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list content domains: %v", err)
	}
	setHTTPMetadata(ctx, 200, "application/json; charset=utf-8", "public, max-age=300", "", "")
	return jsonBytes(map[string]any{"data": items})
}
func (h *PublicContentGrpcHandler) ListContentCategories(ctx context.Context, req *crmpb.ContentQueryRequest) (*wrapperspb.BytesValue, error) {
	items, err := h.catalog.ListCategories(ctx, stringValue(requestMap(req), "domain", "domainCode"))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list content categories: %v", err)
	}
	setHTTPMetadata(ctx, 200, "application/json; charset=utf-8", "public, max-age=300", "", "")
	return jsonBytes(map[string]any{"data": items})
}
func (h *PublicContentGrpcHandler) ListContentFeed(ctx context.Context, req *crmpb.ContentFeedRequest) (*wrapperspb.BytesValue, error) {
	m := requestMap(req)
	result, err := h.catalog.ListFeed(ctx, stringValue(m, "feed"), stringValue(m, "domain", "domainCode"), stringValue(m, "category", "categoryCode"), intValue(m, "limit", 20))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	setHTTPMetadata(ctx, 200, "application/json; charset=utf-8", "public, max-age=60, stale-while-revalidate=300", "", "")
	return jsonBytes(result)
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
	m := requestMap(req)
	filter := catalogdomain.PlanningNewsListFilter{Page: intValue(m, "page", 1), Size: intValue(m, "size", 20), Keyword: stringValue(m, "keyword", "q", "search"), Topic: stringValue(m, "topic", "category"), Location: stringValue(m, "location"), Range: stringValue(m, "range"), VisibleOn: uintValue(m, "visibleOn")}
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

func (h *PublicContentGrpcHandler) UpdatePlanningNewsVisibility(ctx context.Context, req *crmpb.UpdatePlanningNewsVisibilityRequest) (*wrapperspb.BytesValue, error) {
	m := requestMap(req)
	id := uint64Value(m, "id")
	visibleOn := uintValue(m, "visibleOn")
	if id == 0 || visibleOn == 0 {
		return nil, status.Error(codes.InvalidArgument, "id and visibleOn are required")
	}
	item, err := h.catalog.UpdatePlanningNewsVisibility(ctx, id, visibleOn)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "planning news was not found")
		}
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	setHTTPMetadata(ctx, 200, "application/json; charset=utf-8", "private, no-store", "", "")
	return jsonBytes(map[string]any{"data": item})
}

func (h *PublicContentGrpcHandler) AssignPlanningNewsCategory(ctx context.Context, req *crmpb.AssignPlanningNewsCategoryRequest) (*wrapperspb.BytesValue, error) {
	m := requestMap(req)
	id := uint64Value(m, "id")
	category := stringValue(m, "category", "categoryCode", "categoryId")
	if id == 0 || category == "" {
		return nil, status.Error(codes.InvalidArgument, "id and category are required")
	}
	item, err := h.catalog.AssignPlanningNewsCategory(ctx, id, category)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "planning news or category was not found")
		}
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	setHTTPMetadata(ctx, 200, "application/json; charset=utf-8", "private, no-store", "", "")
	return jsonBytes(map[string]any{"data": item})
}

func (h *PublicContentGrpcHandler) ProvisionAdministrativeUnit(ctx context.Context, req *structpb.Struct) (*wrapperspb.BytesValue, error) {
	payload, err := json.Marshal(requestMap(req))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid provisioning request")
	}
	var input dto.ProvisionAdministrativeUnitSeoRequest
	if err := json.Unmarshal(payload, &input); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid provisioning request")
	}
	response, err := h.provisioning.Provision(ctx, &input)
	if err != nil {
		return nil, err
	}
	code := 200
	if response != nil && response.Created {
		code = 201
	}
	setHTTPMetadata(ctx, code, "application/json; charset=utf-8", "private, no-store", "", "")
	return jsonBytes(map[string]any{"code": code, "data": response})
}

func (h *PublicContentGrpcHandler) IngestRuntimeEvent(ctx context.Context, req *structpb.Struct) (*wrapperspb.BytesValue, error) {
	if h.runtimeEvents == nil || h.runtimeEventIngestKey == "" {
		return nil, status.Error(codes.FailedPrecondition, "runtime telemetry is not configured")
	}
	incoming, _ := metadata.FromIncomingContext(ctx)
	provided := ""
	if values := incoming.Get("x-qhpro-telemetry-key"); len(values) > 0 {
		provided = strings.TrimSpace(values[0])
	}
	if len(provided) != len(h.runtimeEventIngestKey) || subtle.ConstantTimeCompare([]byte(provided), []byte(h.runtimeEventIngestKey)) != 1 {
		return nil, status.Error(codes.Unauthenticated, "invalid telemetry key")
	}
	payload, err := json.Marshal(requestMap(req))
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid telemetry event")
	}
	var input dto.SeoRuntimeEventIngestRequest
	if err := json.Unmarshal(payload, &input); err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid telemetry event")
	}
	response, err := h.runtimeEvents.Ingest(ctx, &input)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	setHTTPMetadata(ctx, 202, "application/json; charset=utf-8", "no-store", "", "")
	return jsonBytes(response)
}
