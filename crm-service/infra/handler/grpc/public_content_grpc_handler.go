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
	seoPublic *usecase.SeoPublicUsecase
	sitemap   *usecase.SeoSitemapUsecase
	catalog   *catalogapp.Service
	seoCount  *usecase.SeoCountUsecase
	hubClient *clients.HubGrpcClient
}

func NewPublicContentGrpcHandler(
	seoPublic *usecase.SeoPublicUsecase,
	sitemap *usecase.SeoSitemapUsecase,
	catalog *catalogapp.Service,
	seoCount *usecase.SeoCountUsecase,
	hubClient *clients.HubGrpcClient,
) *PublicContentGrpcHandler {
	return &PublicContentGrpcHandler{
		seoPublic: seoPublic,
		sitemap:   sitemap,
		catalog:   catalog,
		seoCount:  seoCount,
		hubClient: hubClient,
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

func parseVisibleOnsCSV(value string) []uint32 {
	ids := parseUint64CSV(value)
	if len(ids) == 0 {
		return nil
	}
	out := make([]uint32, 0, len(ids))
	for _, id := range ids {
		out = append(out, uint32(id))
	}
	return out
}

func parseUint64CSV(value string) []uint64 {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	out := make([]uint64, 0, len(parts))
	seen := make(map[uint64]struct{}, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		n, err := strconv.ParseUint(part, 10, 64)
		if err != nil || n == 0 {
			continue
		}
		if _, ok := seen[n]; ok {
			continue
		}
		seen[n] = struct{}{}
		out = append(out, n)
	}
	return out
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
	keyword := stringValue(m, "q", "query", "keyword")
	if keyword != "" && h.hubClient != nil {
		h.hubClient.TrackKeywordSearch(ctx, keyword, "crm")
	}
	result, err := h.catalog.SearchPublishedContent(ctx, keyword, stringSlice(m, "kinds"), intValue(m, "limit", 20))
	if err != nil {
		return nil, status.Errorf(codes.Internal, "search public content: %v", err)
	}
	setHTTPMetadata(ctx, 200, "application/json; charset=utf-8", "no-store", "", "")
	return jsonBytes(result)
}

// ListContentCategories @Summary Danh sách chủ đề nội dung
// @Description Lấy danh sách chủ đề từ seo_category (không phụ thuộc domain)
// @Tags CRM Public Content
// @Produce json
// @Success 200 {object} map[string]any
// @Router /v2/crm/public/content/categories [get]
func (h *PublicContentGrpcHandler) ListContentCategories(ctx context.Context, _ *crmpb.ContentQueryRequest) (*wrapperspb.BytesValue, error) {
	items, err := h.catalog.ListCategories(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list content categories: %v", err)
	}
	setHTTPMetadata(ctx, 200, "application/json; charset=utf-8", "public, max-age=300", "", "")
	return jsonBytes(map[string]any{"data": items})
}

// ListPlanningNews @Summary Danh sách tin SEO công khai
// @Tags CRM Public Content
// @Produce json
// @Param page query int false "Trang (bắt đầu từ 1)"
// @Param size query int false "Số bản ghi / trang"
// @Param visibleOns query string false "CSV vị trí hiển thị, ví dụ 40 hoặc 20,40"
// @Param lastDates query int false "Số ngày gần nhất theo updatedAt (vd 10, 20); 0 = không lọc"
// @Param refType query int false "Lọc theo refType (vd 800)"
// @Param sort query string false "Sắp xếp: updated_desc, published_desc, ..."
// @Param provinceIds query string false "CSV tỉnh/thành, ví dụ 1 hoặc 1,2,3"
// @Param categoryId query int false "Lọc theo chủ đề seo_category.id (vd 4)"
// @Success 200 {object} map[string]any
// @Router /v2/crm/public/seo/news [get]
func (h *PublicContentGrpcHandler) ListPlanningNews(ctx context.Context, req *crmpb.PlanningNewsQueryRequest) (*wrapperspb.BytesValue, error) {
	if req == nil {
		req = &crmpb.PlanningNewsQueryRequest{}
	}

	filter := catalogdomain.PlanningNewsListFilter{
		Page:        int(req.GetPage()),
		Size:        int(req.GetSize()),
		VisibleOns:  parseVisibleOnsCSV(req.GetVisibleOns()),
		LastDates:   req.GetLastDates(),
		RefType:     req.GetRefType(),
		Sort:        strings.TrimSpace(req.GetSort()),
		ProvinceIDs: parseUint64CSV(req.GetProvinceIds()),
		CategoryID:  req.GetCategoryId(),
	}
	result, err := h.catalog.ListPlanningNews(ctx, filter)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "list planning news: %v", err)
	}
	setHTTPMetadata(ctx, 200, "application/json; charset=utf-8", "public, max-age=60, stale-while-revalidate=300", "", "")
	return jsonBytes(result)
}

// GetPlanningNews @Summary Chi tiết tin SEO công khai theo id
// @Tags CRM Public Content
// @Produce json
// @Param id path uint64 true "seo_domain.id"
// @Success 200 {object} map[string]any
// @Router /v2/crm/public/seo/news/{id} [get]
func (h *PublicContentGrpcHandler) GetPlanningNews(ctx context.Context, req *crmpb.PlanningNewsByIDRequest) (*wrapperspb.BytesValue, error) {
	if req == nil || req.GetId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	result, err := h.catalog.GetPlanningNews(ctx, req.GetId())
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

// TrackSeoClick @Summary Ghi nhận click trang SEO công khai
// @Description Hỗ trợ 2 kiểu body: {"id":8} hoặc {"refType":100,"refId":123}. Ưu tiên id nếu có.
// @Tags CRM Public Content
// @Accept json
// @Produce json
// @Param body body crmpb.TrackSeoClickRequest true "id hoặc refType+refId"
// @Success 202 {object} map[string]any
// @Router /v2/crm/public/seo/click [post]
// TrackSeoClick records a public SEO page click into the in-memory counter.
// Counts are flushed to seo_count every ~10s by a single pending job.
func (h *PublicContentGrpcHandler) TrackSeoClick(ctx context.Context, req *crmpb.TrackSeoClickRequest) (*wrapperspb.BytesValue, error) {
	if h.seoCount == nil {
		return nil, status.Error(codes.FailedPrecondition, "seo click counter is not configured")
	}
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "id hoặc refType+refId là bắt buộc")
	}

	seoID := req.GetId()
	refID := req.GetRefId()
	refType := req.GetRefType()

	// Kiểu 1: chỉ id (seo_domain.id)
	if seoID > 0 {
		if err := h.seoCount.IncrementClickBySeoID(ctx, seoID); err != nil {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		setHTTPMetadata(ctx, 202, "application/json; charset=utf-8", "no-store", "", "")
		return jsonBytes(map[string]any{
			"accepted": true,
			"id":       seoID,
		})
	}

	// Kiểu 2: refType + refId
	if refID == 0 || refType == 0 {
		return nil, status.Error(codes.InvalidArgument, "cần truyền id hoặc cặp refType+refId")
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

// AssignPlanningNewsCategory @Summary Gán chủ đề cho tin quy hoạch
// @Tags CRM Admin Content
// @Accept json
// @Produce json
// @Param id path uint64 true "seo domain id"
// @Param body body crmpb.AssignPlanningNewsCategoryRequest true "payload"
// @Success 200 {object} map[string]any
// @Router /v2/crm/admin/content/planning-news/{id}/category [put]
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
// @Tags CRM Admin Content
// @Accept json
// @Produce json
// @Param body body crmpb.CreateSeoCategoryRequest true "payload"
// @Success 201 {object} map[string]any
// @Router /v2/crm/admin/content/categories [post]
func (h *PublicContentGrpcHandler) CreateSeoCategory(ctx context.Context, req *crmpb.CreateSeoCategoryRequest) (*wrapperspb.BytesValue, error) {
	m := requestMap(req)
	input := catalogdomain.CreateSeoCategoryInput{
		Code:        stringValue(m, "code"),
		Slug:        stringValue(m, "slug"),
		Name:        stringValue(m, "name"),
		Description: stringValue(m, "description"),
		Status:      stringValue(m, "status"),
		SortOrder:   intValue(m, "sortOrder", 0),
	}
	item, err := h.catalog.CreateCategory(ctx, input)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	setHTTPMetadata(ctx, 201, "application/json; charset=utf-8", "private, no-store", "", "")
	return jsonBytes(map[string]any{"data": item})
}

// UpdateSeoCategory @Summary Cập nhật chủ đề SEO
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
			return nil, status.Error(codes.NotFound, "category was not found")
		}
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	setHTTPMetadata(ctx, 200, "application/json; charset=utf-8", "private, no-store", "", "")
	return jsonBytes(map[string]any{"data": item})
}

// DeleteSeoCategory @Summary Xóa chủ đề SEO
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
