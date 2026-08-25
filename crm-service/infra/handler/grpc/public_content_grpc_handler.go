package handler_grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"common/injection"
	_middleware "common/middleware"
	"crm/config"
	catalogapp "crm/internal/modules/contentcatalog/application"
	catalogdomain "crm/internal/modules/contentcatalog/domain"
	"crm/internal/usecase"
	"pb/clients"

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
	seoPublic             *usecase.SeoPublicUsecase
	sitemap               *usecase.SeoSitemapUsecase
	provisioning          *usecase.SeoAdministrativeUnitProvisioningUsecase
	catalog               *catalogapp.Service
	runtimeEvents         *usecase.SeoRuntimeEventUsecase
	seoCount              *usecase.SeoCountUsecase
	runtimeEventIngestKey string
	hubClient             *clients.HubGrpcClient
}

func NewPublicContentGrpcHandler(seoPublic *usecase.SeoPublicUsecase, sitemap *usecase.SeoSitemapUsecase, provisioning *usecase.SeoAdministrativeUnitProvisioningUsecase, catalog *catalogapp.Service, runtimeEvents *usecase.SeoRuntimeEventUsecase, seoCount *usecase.SeoCountUsecase, runtimeEventIngestKey string) *PublicContentGrpcHandler {
	return &PublicContentGrpcHandler{
		seoPublic:             seoPublic,
		sitemap:               sitemap,
		provisioning:          provisioning,
		catalog:               catalog,
		runtimeEvents:         runtimeEvents,
		seoCount:              seoCount,
		runtimeEventIngestKey: strings.TrimSpace(runtimeEventIngestKey),
		hubClient:             injection.NewHubClient(),
	}
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

// TASK9_VISIBLE_ON_VALIDATOR_BEGIN
func validatePlanningNewsVisibleOn(value uint32) error {
	visibleOn := catalogdomain.PlanningNewsVisibleOn(value)
	if visibleOn.IsValid() {
		return nil
	}
	return status.Error(
		codes.InvalidArgument,
		"visibleOn chỉ chấp nhận 10, 20, 30, 40 hoặc 50",
	)
}

func parsePlanningNewsVisibleOns(values ...string) ([]uint32, error) {
	seen := make(map[uint32]struct{})
	result := make([]uint32, 0)

	for _, raw := range values {
		raw = strings.TrimSpace(raw)
		if raw == "" {
			continue
		}
		for _, part := range strings.Split(raw, ",") {
			part = strings.TrimSpace(part)
			if part == "" {
				continue
			}
			parsed, err := strconv.ParseUint(part, 10, 32)
			if err != nil || parsed == 0 {
				return nil, status.Error(
					codes.InvalidArgument,
					"visibleOns chỉ chấp nhận danh sách giá trị 10, 20, 30, 40 hoặc 50",
				)
			}
			value := uint32(parsed)
			if err := validatePlanningNewsVisibleOn(value); err != nil {
				return nil, err
			}
			if _, ok := seen[value]; ok {
				continue
			}
			seen[value] = struct{}{}
			result = append(result, value)
		}
	}

	return result, nil
}

func appendUniqueVisibleOns(values []uint32, items ...uint32) []uint32 {
	if len(items) == 0 {
		return values
	}
	seen := make(map[uint32]struct{}, len(values)+len(items))
	for _, value := range values {
		seen[value] = struct{}{}
	}
	for _, item := range items {
		if item == 0 {
			continue
		}
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		values = append(values, item)
	}
	return values
}

// TASK9_VISIBLE_ON_VALIDATOR_END

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

	m := requestMap(req)
	visibleOns, err := parsePlanningNewsVisibleOns(
		req.GetVisibleOns(),
		stringValue(m, "visibleOns", "visible_ons"),
	)
	if err != nil {
		return nil, err
	}
	for _, part := range stringSlice(m, "visibleOns") {
		more, parseErr := parsePlanningNewsVisibleOns(part)
		if parseErr != nil {
			return nil, parseErr
		}
		visibleOns = appendUniqueVisibleOns(visibleOns, more...)
	}
	if req.VisibleOn != nil && req.GetVisibleOn() > 0 {
		visibleOns = appendUniqueVisibleOns(visibleOns, req.GetVisibleOn())
	}
	if single := uintValue(m, "visibleOn"); single > 0 {
		if validateErr := validatePlanningNewsVisibleOn(single); validateErr != nil {
			return nil, validateErr
		}
		visibleOns = appendUniqueVisibleOns(visibleOns, single)
	}

	filter := catalogdomain.PlanningNewsListFilter{
		Page:       intValue(m, "page", 0),
		Size:       intValue(m, "size", 20),
		Keyword:    stringValue(m, "keyword", "q", "search"),
		CategoryID: req.GetCategoryId(),
		ProvinceID: req.GetProvinceId(),
		Location:   stringValue(m, "location"),
		Range:      stringValue(m, "range"),
		VisibleOns: visibleOns,
		Sort:       stringValue(m, "sort"),
		SortBy:     stringValue(m, "sort_by", "sortBy"),
		SortOrder:  stringValue(m, "sort_order", "sortOrder"),
		RefType:    req.GetRefType(),
	}
	if keyword := strings.TrimSpace(filter.Keyword); keyword != "" && h.hubClient != nil {
		h.hubClient.TrackKeywordSearch(ctx, keyword, "crm")
	}
	if filter.Sort == "" {
		filter.Sort = req.GetSort()
	}
	if filter.SortBy == "" {
		filter.SortBy = req.GetSortBy()
	}
	if filter.SortOrder == "" {
		filter.SortOrder = req.GetSortOrder()
	}
	if filter.RefType == 0 {
		if v := uintValue(m, "refType"); v > 0 {
			filter.RefType = v
		} else if v := uintValue(m, "ref_type"); v > 0 {
			filter.RefType = v
		}
	}
	if filter.CategoryID == 0 {
		if id := uint64Value(m, "categoryId"); id > 0 {
			filter.CategoryID = id
		} else if id := uint64Value(m, "category_id"); id > 0 {
			filter.CategoryID = id
		}
	}
	if filter.ProvinceID == 0 {
		if id := uint64Value(m, "provinceId"); id > 0 {
			filter.ProvinceID = id
		} else if id := uint64Value(m, "province_id"); id > 0 {
			filter.ProvinceID = id
		}
	}
	if len(filter.VisibleOns) == 1 {
		filter.VisibleOn = filter.VisibleOns[0]
	}
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

// ListFollowedSeoCategories @Summary Danh sách chủ đề đang theo dõi
// @Description Lấy danh sách chủ đề (seo_category) user đã theo dõi
// @Tags CRM Client Content
// @Produce json
// @Param page query int false "page (0-based)"
// @Param size query int false "size"
// @Success 200 {object} map[string]any
// @Router /v2/crm/client/content/categories/followed [get]
func (h *PublicContentGrpcHandler) ListFollowedSeoCategories(ctx context.Context, req *crmpb.FollowedSeoCategoryQueryRequest) (*wrapperspb.BytesValue, error) {
	userID, err := authenticatedProfileID(ctx)
	if err != nil {
		return nil, err
	}
	if req == nil {
		req = &crmpb.FollowedSeoCategoryQueryRequest{}
	}
	result, err := h.catalog.ListFollowedSeoCategories(ctx, userID, catalogdomain.FollowedSeoCategoryListFilter{
		Page: int(req.GetPage()),
		Size: int(req.GetSize()),
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list followed seo categories: %v", err)
	}
	setHTTPMetadata(ctx, 200, "application/json; charset=utf-8", "private, no-store", "", "")
	return jsonBytes(result)
}

// FollowSeoCategory @Summary Theo dõi chủ đề
// @Description User theo dõi một chủ đề (seo_category)
// @Tags CRM Client Content
// @Produce json
// @Param id path uint64 true "category id"
// @Success 200 {object} map[string]any
// @Router /v2/crm/client/content/categories/{id}/followed [post]
func (h *PublicContentGrpcHandler) FollowSeoCategory(ctx context.Context, req *crmpb.SeoCategoryFollowRequest) (*wrapperspb.BytesValue, error) {
	userID, err := authenticatedProfileID(ctx)
	if err != nil {
		return nil, err
	}
	if req == nil || req.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "category id is required")
	}
	item, err := h.catalog.FollowSeoCategory(ctx, userID, req.GetId())
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "category was not found or is not active")
		}
		return nil, status.Errorf(codes.Internal, "follow seo category: %v", err)
	}
	setHTTPMetadata(ctx, 200, "application/json; charset=utf-8", "private, no-store", "", "")
	return jsonBytes(map[string]any{"followed": true, "data": item})
}

// UnfollowSeoCategory @Summary Bỏ theo dõi chủ đề
// @Description User bỏ theo dõi một chủ đề (seo_category)
// @Tags CRM Client Content
// @Produce json
// @Param id path uint64 true "category id"
// @Success 200 {object} map[string]any
// @Router /v2/crm/client/content/categories/{id}/followed [delete]
func (h *PublicContentGrpcHandler) UnfollowSeoCategory(ctx context.Context, req *crmpb.SeoCategoryFollowRequest) (*wrapperspb.BytesValue, error) {
	userID, err := authenticatedProfileID(ctx)
	if err != nil {
		return nil, err
	}
	if req == nil || req.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "category id is required")
	}
	if err := h.catalog.UnfollowSeoCategory(ctx, userID, req.GetId()); err != nil {
		return nil, status.Errorf(codes.Internal, "unfollow seo category: %v", err)
	}
	setHTTPMetadata(ctx, 200, "application/json; charset=utf-8", "private, no-store", "", "")
	return jsonBytes(map[string]any{"followed": false, "categoryId": req.GetId()})
}

func (h *PublicContentGrpcHandler) UpdatePlanningNewsVisibility(ctx context.Context, req *crmpb.UpdatePlanningNewsVisibilityRequest) (*wrapperspb.BytesValue, error) {

	// TASK9_UPDATE_VISIBLE_ON_VALIDATION_BEGIN
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}
	if err := validatePlanningNewsVisibleOn(req.GetVisibleOn()); err != nil {
		return nil, err
	}
	// TASK9_UPDATE_VISIBLE_ON_VALIDATION_END

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	id := req.GetId()
	visibleOn := uint32(req.GetVisibleOn())
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

func (h *PublicContentGrpcHandler) SetPlanningNewsTop(ctx context.Context, req *crmpb.SetPlanningNewsTopRequest) (*wrapperspb.BytesValue, error) {
	if req == nil || req.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "planning news id is required")
	}
	item, err := h.catalog.SetPlanningNewsTop(ctx, req.GetId())
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "planning news was not found")
		}
		return nil, status.Errorf(codes.Internal, "set planning news top: %v", err)
	}
	setHTTPMetadata(ctx, 200, "application/json; charset=utf-8", "private, no-store", "", "")
	return jsonBytes(map[string]any{"data": item})
}

func (h *PublicContentGrpcHandler) AssignPlanningNewsCategory(ctx context.Context, req *crmpb.AssignPlanningNewsCategoryRequest) (*wrapperspb.BytesValue, error) {
	if req == nil || req.GetId() == 0 || req.GetCategoryId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "id and categoryId are required")
	}
	item, err := h.catalog.AssignPlanningNewsCategory(ctx, req.GetId(), req.GetCategoryId())
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "planning news or category was not found")
		}
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	setHTTPMetadata(ctx, 200, "application/json; charset=utf-8", "private, no-store", "", "")
	return jsonBytes(map[string]any{"data": item})
}

// CreateSeoCategory @Summary Tạo chủ đề SEO
// @Description Admin tạo chủ đề (seo_category)
// @Tags CRM Admin Content
// @Accept json
// @Produce json
// @Param body body crmpb.CreateSeoCategoryRequest true "payload"
// @Success 201 {object} map[string]any
// @Router /v2/crm/admin/content/categories [post]
func (h *PublicContentGrpcHandler) CreateSeoCategory(ctx context.Context, req *crmpb.CreateSeoCategoryRequest) (*wrapperspb.BytesValue, error) {
	m := requestMap(req)
	input := catalogdomain.CreateSeoCategoryInput{
		DomainID:    uint64Value(m, "domainId"),
		DomainCode:  stringValue(m, "domain", "domainCode"),
		Code:        stringValue(m, "code"),
		Slug:        stringValue(m, "slug"),
		Name:        stringValue(m, "name"),
		Description: stringValue(m, "description"),
		Status:      stringValue(m, "status"),
		SortOrder:   intValue(m, "sortOrder", 0),
	}
	item, err := h.catalog.CreateCategory(ctx, input)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "content domain was not found")
		}
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	setHTTPMetadata(ctx, 201, "application/json; charset=utf-8", "private, no-store", "", "")
	return jsonBytes(map[string]any{"data": item})
}

// UpdateSeoCategory @Summary Cập nhật chủ đề SEO
// @Description Admin sửa chủ đề (seo_category)
// @Tags CRM Admin Content
// @Accept json
// @Produce json
// @Param id path uint64 true "category id"
// @Param body body crmpb.UpdateSeoCategoryRequest true "payload"
// @Success 200 {object} map[string]any
// @Router /v2/crm/admin/content/categories/{id} [put]
func (h *PublicContentGrpcHandler) UpdateSeoCategory(ctx context.Context, req *crmpb.UpdateSeoCategoryRequest) (*wrapperspb.BytesValue, error) {
	if req == nil || req.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "category id is required")
	}
	input := catalogdomain.UpdateSeoCategoryInput{}
	if req.DomainId != nil {
		value := req.GetDomainId()
		input.DomainID = &value
	}
	if req.Domain != nil {
		value := strings.TrimSpace(req.GetDomain())
		input.DomainCode = &value
	}
	if req.Code != nil {
		value := strings.TrimSpace(req.GetCode())
		input.Code = &value
	}
	if req.Slug != nil {
		value := strings.TrimSpace(req.GetSlug())
		input.Slug = &value
	}
	if req.Name != nil {
		value := strings.TrimSpace(req.GetName())
		input.Name = &value
	}
	if req.Description != nil {
		value := strings.TrimSpace(req.GetDescription())
		input.Description = &value
	}
	if req.Status != nil {
		value := strings.TrimSpace(req.GetStatus())
		input.Status = &value
	}
	if req.SortOrder != nil {
		value := int(req.GetSortOrder())
		input.SortOrder = &value
	}
	item, err := h.catalog.UpdateCategory(ctx, req.GetId(), input)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "category or domain was not found")
		}
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	setHTTPMetadata(ctx, 200, "application/json; charset=utf-8", "private, no-store", "", "")
	return jsonBytes(map[string]any{"data": item})
}

// DeleteSeoCategory @Summary Xóa chủ đề SEO
// @Description Admin soft-delete chủ đề (seo_category)
// @Tags CRM Admin Content
// @Produce json
// @Param id path uint64 true "category id"
// @Success 200 {object} map[string]any
// @Router /v2/crm/admin/content/categories/{id} [delete]
func (h *PublicContentGrpcHandler) DeleteSeoCategory(ctx context.Context, req *crmpb.DeleteSeoCategoryRequest) (*wrapperspb.BytesValue, error) {
	if req == nil || req.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "category id is required")
	}
	if err := h.catalog.DeleteCategory(ctx, req.GetId()); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, status.Error(codes.NotFound, "category was not found")
		}
		return nil, status.Errorf(codes.Internal, "delete seo category: %v", err)
	}
	setHTTPMetadata(ctx, 200, "application/json; charset=utf-8", "private, no-store", "", "")
	return jsonBytes(map[string]any{"data": map[string]any{"id": req.GetId(), "deleted": true}})
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

// TrackSeoClick records a public SEO page click into the in-memory counter.
// Counts are flushed to seo_count every 10s by a single pending job.
func (h *PublicContentGrpcHandler) TrackSeoClick(ctx context.Context, req *crmpb.TrackSeoClickRequest) (*wrapperspb.BytesValue, error) {
	if h.seoCount == nil {
		return nil, status.Error(codes.FailedPrecondition, "seo click counter is not configured")
	}
	refID := req.GetRefId()
	refType := req.GetRefType()
	if refID == 0 {
		return nil, status.Error(codes.InvalidArgument, "refId is required")
	}
	if refType == 0 {
		return nil, status.Error(codes.InvalidArgument, "refType is required")
	}
	if err := h.seoCount.IncrementClick(ctx, refType, refID); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	setHTTPMetadata(ctx, 202, "application/json; charset=utf-8", "no-store", "", "")
	return jsonBytes(map[string]any{
		"accepted": true,
		"refId":    refID,
		"refType":  refType,
	})
}
