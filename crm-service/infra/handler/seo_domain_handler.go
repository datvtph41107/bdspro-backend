package handler

import (
	_dto "common/domain/dto"
	_errors "common/errors"
	"context"
	"crm/config"
	"crm/infra/mapper"
	"crm/infra/utils"
	"crm/internal"
	"crm/internal/dto"
	"crm/internal/enums"
	"crm/internal/usecase"
	"strings"

	crmpb "pb/types/crm"
	sharepb "pb/types/shared"
)

type SeoDomainHandler struct {
	crmpb.UnimplementedSeoDomainServiceServer

	seoDomainUsecase *usecase.SeoDomainUsecase
	seoRenderUsecase *usecase.SeoRenderUsecase
	seoDomainMapper  *mapper.SeoDomainMapper
}

var (
	_ crmpb.SeoDomainServiceServer = (*SeoDomainHandler)(nil)
)

// @bind: crm/infra/handler.SeoDomainHandler
func NewSeoDomainHandler(
	seoDomainUsecase *usecase.SeoDomainUsecase,
	seoRenderUsecase *usecase.SeoRenderUsecase,
	seoDomainMapper *mapper.SeoDomainMapper,
) *SeoDomainHandler {
	return &SeoDomainHandler{
		seoDomainUsecase: seoDomainUsecase,
		seoRenderUsecase: seoRenderUsecase,
		seoDomainMapper:  seoDomainMapper,
	}
}

// -----------------------------------------------------------------------------
// Legacy source APIs.
// -----------------------------------------------------------------------------

func (h *SeoDomainHandler) GetSeoRefTypes(ctx context.Context, req *crmpb.GetSeoRefTypesRequest) (*crmpb.GetSeoRefTypesResponse, error) {
	return &crmpb.GetSeoRefTypesResponse{
		Data: h.seoDomainMapper.SeoRefTypeConfigsToPb(h.seoDomainUsecase.GetSeoRefTypes(ctx)),
	}, nil
}

func (h *SeoDomainHandler) SearchSeoRefValues(ctx context.Context, req *crmpb.SearchSeoRefValuesRequest) (*crmpb.SearchSeoRefValuesResponse, error) {
	items, total, err := h.seoDomainUsecase.SearchSeoRefValues(ctx, &dto.SearchSeoRefValuesRequest{
		Pagable: _dto.Pagable{
			Text: req.GetText(),
			Page: req.GetPage(),
			Size: req.GetSize(),
		},
		RefType:       req.GetRefType(),
		ResolverKey:   req.GetResolverKey(),
		SourceService: req.GetSourceService(),
	})
	if err != nil {
		return nil, err
	}

	responses := make([]*crmpb.SeoRefValueResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, seoRefValueToPb(item))
	}

	return &crmpb.SearchSeoRefValuesResponse{
		Data:  responses,
		Total: total,
	}, nil
}

func (h *SeoDomainHandler) ResolveSeoRefValue(ctx context.Context, req *crmpb.ResolveSeoRefValueRequest) (*crmpb.SeoRefSnapshotResponse, error) {
	item, err := h.seoDomainUsecase.ResolveSeoRefValue(ctx, &dto.ResolveSeoRefValueRequest{
		RefType:       req.RefType,
		RefID:         req.RefId,
		ResolverKey:   req.ResolverKey,
		SourceService: req.SourceService,
	})
	if err != nil {
		return nil, err
	}

	return seoRefSnapshotToPb(item), nil
}

// -----------------------------------------------------------------------------
// Source Registry V2 APIs.
// -----------------------------------------------------------------------------

func (h *SeoDomainHandler) GetSeoSources(ctx context.Context, req *crmpb.GetSeoSourcesRequest) (*crmpb.GetSeoSourcesResponse, error) {
	return &crmpb.GetSeoSourcesResponse{
		Data: h.seoDomainMapper.SeoSourceConfigsToPb(h.seoDomainUsecase.GetSeoSources(ctx)),
	}, nil
}

func (h *SeoDomainHandler) GetSeoSource(ctx context.Context, req *crmpb.GetSeoSourceRequest) (*crmpb.SeoSourceResponse, error) {
	cfg, err := h.seoDomainUsecase.GetSeoSource(ctx, req.SourceKey)
	if err != nil {
		return nil, err
	}
	return h.seoDomainMapper.SeoSourceConfigToPb(cfg), nil
}

func (h *SeoDomainHandler) SearchSeoSourceValues(ctx context.Context, req *crmpb.SearchSeoSourceValuesRequest) (*crmpb.SearchSeoSourceValuesResponse, error) {
	items, total, err := h.seoDomainUsecase.SearchSeoSourceValues(ctx, &dto.SearchSeoSourceValuesRequest{
		Pagable: _dto.Pagable{
			Text: req.GetText(),
			Page: req.GetPage(),
			Size: req.GetSize(),
		},
		SourceKey: req.GetSourceKey(),
	})
	if err != nil {
		return nil, err
	}

	responses := make([]*crmpb.SeoSourceValueResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, seoSourceValueToPb(item))
	}

	return &crmpb.SearchSeoSourceValuesResponse{
		Data:  responses,
		Total: total,
	}, nil
}

func (h *SeoDomainHandler) ResolveSeoSourceValue(ctx context.Context, req *crmpb.ResolveSeoSourceValueRequest) (*crmpb.SeoSourceSnapshotResponse, error) {
	item, err := h.seoDomainUsecase.ResolveSeoSourceValue(ctx, &dto.ResolveSeoSourceValueRequest{
		SourceKey: req.SourceKey,
		RefID:     req.RefId,
	})
	if err != nil {
		return nil, err
	}
	return seoSourceSnapshotToPb(item), nil
}

// -----------------------------------------------------------------------------
// SEO domain CRUD.
// -----------------------------------------------------------------------------

func (h *SeoDomainHandler) CreateSeoDomain(ctx context.Context, req *crmpb.CreateSeoDomainRequest) (*crmpb.SeoDomainResponse, error) {
	seoReq := &dto.SeoDomainCreateRequest{
		Slug:              req.Slug,
		OriginURL:         req.OriginUrl,
		CanonicalURL:      req.CanonicalUrl,
		RefType:           req.RefType,
		RefID:             req.RefId,
		RefSource:         req.RefSource,
		RefLabel:          req.RefLabel,
		RefURL:            req.RefUrl,
		SourceStatus:      req.SourceStatus,
		RefMissing:        req.RefMissing,
		RefSnapshotJSON:   req.RefSnapshotJson,
		RefHash:           req.RefHash,
		RefLastSyncedAt:   req.RefLastSyncedAt,
		Scope:             req.Scope,
		Title:             req.Title,
		Description:       req.Description,
		Content:           req.Content,
		Summary:           req.Summary,
		Published:         req.Published,
		PublishedAt:       req.PublishedAt,
		IsSiteMap:         req.IsSiteMap,
		IsIndex:           req.IsIndex,
		IsRobot:           req.IsRobot,
		SiteMapLastedAt:   req.SiteMapLastedAt,
		SitemapPriority:   req.SitemapPriority,
		SitemapChangeFreq: req.SitemapChangeFreq,
		NeedGenerate:      req.NeedGenerate,
		GeneratedAt:       req.GeneratedAt,
		SourceUpdatedAt:   req.SourceUpdatedAt,
		StaticHtmlPath:    req.StaticHtmlPath,
		StaticHtmlHash:    req.StaticHtmlHash,
		DeepLink:          req.DeepLink,
		Metadata:          req.Metadata,
		Note:              req.Note,
		PageStatus:        req.PageStatus,
		RenderStatus:      req.RenderStatus,
		RenderedHTML:      req.RenderedHtml,
		LastRenderError:   req.LastRenderError,
		TemplateKey:       req.TemplateKey,
		TemplateVersion:   req.TemplateVersion,
	}

	created, err := h.seoDomainUsecase.CreateSeoDomain(ctx, seoReq)
	if err != nil {
		return nil, err
	}
	return h.seoDomainMapper.SeoDomainToPb(created), nil
}

func (h *SeoDomainHandler) GetSeoDomain(ctx context.Context, req *sharepb.IdRequest) (*crmpb.SeoDomainResponse, error) {
	seoDomain, err := h.seoDomainUsecase.GetSeoDomain(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return h.seoDomainMapper.SeoDomainToPb(seoDomain), nil
}

// GetSeoRefs @Summary Danh sách SEO liên kết theo ref
// @Description Lấy danh sách bài SEO public liên kết tới domain (refType/refId) qua bảng seo_refs
func (h *SeoDomainHandler) GetSeoRefs(ctx context.Context, req *crmpb.GetSeoRefsRequest) (*crmpb.GetSeoDomainListResponse, error) {
	if req == nil {
		req = &crmpb.GetSeoRefsRequest{}
	}

	items, total, err := h.seoDomainUsecase.GetSeoRefs(ctx, &dto.GetSeoRefsRequest{
		Pagable: _dto.Pagable{
			Page: req.GetPage(),
			Size: req.GetSize(),
		},
		RefType: enums.ESEORefType(req.GetRefType()),
		RefID:   req.GetRefId(),
	})
	if err != nil {
		return nil, err
	}

	return &crmpb.GetSeoDomainListResponse{
		Data:  h.seoDomainMapper.SeoDomainListToPb(items),
		Total: total,
	}, nil
}

func (h *SeoDomainHandler) GetSeoDomainList(ctx context.Context, req *crmpb.GetSeoDomainListRequest) (*crmpb.GetSeoDomainListResponse, error) {
	if req == nil {
		req = &crmpb.GetSeoDomainListRequest{}
	}

	listReq := &dto.SeoDomainListRequest{
		Pagable: _dto.Pagable{
			Text: req.GetKeyword(),
			Page: req.GetPage(),
			Size: req.GetSize(),
		},
		Keyword:      req.Keyword,
		RefType:      req.RefType,
		RefID:        req.RefId,
		RefSource:    req.RefSource,
		SourceStatus: req.SourceStatus,
		RefMissing:   req.RefMissing,
		Scope:        req.Scope,
		Published:    req.Published,
		IsSiteMap:    req.IsSiteMap,
		IsIndex:      req.IsIndex,
		IsRobot:      req.IsRobot,
		NeedGenerate: req.NeedGenerate,
		PageStatus:   req.PageStatus,
		RenderStatus: req.RenderStatus,
		TemplateKey:  req.TemplateKey,
	}
	listReq.Normalize()

	items, total, err := h.seoDomainUsecase.GetSeoDomainList(ctx, listReq)
	if err != nil {
		return nil, err
	}
	return &crmpb.GetSeoDomainListResponse{
		Data:  h.seoDomainMapper.SeoDomainListToPb(items),
		Total: total,
	}, nil
}

func (h *SeoDomainHandler) UpdateSeoDomain(ctx context.Context, req *crmpb.UpdateSeoDomainRequest) (*crmpb.UpdateSeoDomainResponse, error) {
	updateReq := &dto.SeoDomainUpdateRequest{
		Slug:              req.Slug,
		OriginURL:         req.OriginUrl,
		CanonicalURL:      req.CanonicalUrl,
		RefType:           req.RefType,
		RefID:             req.RefId,
		RefSource:         req.RefSource,
		RefLabel:          req.RefLabel,
		RefURL:            req.RefUrl,
		SourceStatus:      req.SourceStatus,
		RefMissing:        req.RefMissing,
		RefSnapshotJSON:   req.RefSnapshotJson,
		RefHash:           req.RefHash,
		RefLastSyncedAt:   req.RefLastSyncedAt,
		Scope:             req.Scope,
		Title:             req.Title,
		Description:       req.Description,
		Content:           req.Content,
		Summary:           req.Summary,
		Published:         req.Published,
		PublishedAt:       req.PublishedAt,
		IsSiteMap:         req.IsSiteMap,
		IsIndex:           req.IsIndex,
		IsRobot:           req.IsRobot,
		SiteMapLastedAt:   req.SiteMapLastedAt,
		SitemapPriority:   req.SitemapPriority,
		SitemapChangeFreq: req.SitemapChangeFreq,
		NeedGenerate:      req.NeedGenerate,
		GeneratedAt:       req.GeneratedAt,
		SourceUpdatedAt:   req.SourceUpdatedAt,
		StaticHtmlPath:    req.StaticHtmlPath,
		StaticHtmlHash:    req.StaticHtmlHash,
		DeepLink:          req.DeepLink,
		Metadata:          req.Metadata,
		Note:              req.Note,
		PageStatus:        req.PageStatus,
		RenderStatus:      req.RenderStatus,
		RenderedHTML:      req.RenderedHtml,
		LastRenderError:   req.LastRenderError,
		TemplateKey:       req.TemplateKey,
		TemplateVersion:   req.TemplateVersion,
	}

	updated, err := h.seoDomainUsecase.UpdateSeoDomain(ctx, req.Id, updateReq)
	if err != nil {
		return nil, err
	}
	return &crmpb.UpdateSeoDomainResponse{
		SeoDomain: h.seoDomainMapper.SeoDomainToPb(updated),
	}, nil
}

func (h *SeoDomainHandler) DeleteSeoDomain(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	if err := h.seoDomainUsecase.DeleteSeoDomain(ctx, req.Id); err != nil {
		return nil, err
	}
	return &sharepb.SubmitResponse{Id: req.Id, Message: "Seo domain deleted successfully"}, nil
}

// -----------------------------------------------------------------------------
// SEO page command APIs.
// These provide business-action endpoints while reusing existing usecase methods.
// A dedicated command usecase can replace these implementations later.
// -----------------------------------------------------------------------------

func (h *SeoDomainHandler) CreateSeoDraft(ctx context.Context, req *crmpb.CreateSeoDraftRequest) (*crmpb.SeoDomainResponse, error) {
	if req == nil || req.Source == nil || req.Seo == nil {
		return nil, _errors.ReturnError(service.RequestInvalid)
	}

	sourceKey := strings.TrimSpace(req.Source.SourceKey)
	if sourceKey == "" {
		return nil, _errors.ReturnError(service.SEOSourceKeyRequired)
	}

	cfg, err := h.seoDomainUsecase.GetSeoSource(ctx, sourceKey)
	if err != nil {
		return nil, err
	}

	seoReq := buildCreateSeoDraftRequest(req, cfg.Value, cfg.ResolverKey)
	seoReq.RefSource = cfg.ResolverKey
	seoReq.SourceStatus = "manual"

	if sourceKey != "manual" {
		if req.Source.RefId == 0 {
			return nil, _errors.ReturnError(service.SEOSourceRefIDRequired)
		}

		snapshot, err := h.seoDomainUsecase.ResolveSeoSourceValue(ctx, &dto.ResolveSeoSourceValueRequest{
			SourceKey: sourceKey,
			RefID:     req.Source.RefId,
		})
		if err != nil {
			return nil, err
		}

		refID := req.Source.RefId
		seoReq.RefID = &refID
		seoReq.RefSource = snapshot.ResolverKey
		seoReq.RefLabel = snapshot.Label
		seoReq.RefURL = firstNonEmpty(snapshot.PublicURL, snapshot.AdminURL)
		seoReq.SourceStatus = firstNonEmpty(snapshot.SourceStatus, "linked")
		seoReq.RefSnapshotJSON = &snapshot.DataJSON
		seoReq.RefHash = snapshot.Hash

		if strings.TrimSpace(seoReq.Slug) == "" {
			seoReq.Slug = snapshot.SuggestedSlug
		}
		if strings.TrimSpace(seoReq.Title) == "" {
			seoReq.Title = snapshot.SuggestedTitle
		}
		if strings.TrimSpace(seoReq.Description) == "" {
			seoReq.Description = snapshot.SuggestedMetaDescription
		}
		if strings.TrimSpace(seoReq.OriginURL) == "" {
			seoReq.OriginURL = firstNonEmpty(snapshot.PublicURL, snapshot.AdminURL)
		}
		if seoReq.CanonicalURL == nil || strings.TrimSpace(*seoReq.CanonicalURL) == "" {
			canonical := firstNonEmpty(snapshot.PublicURL, snapshot.AdminURL, seoReq.OriginURL)
			seoReq.CanonicalURL = &canonical
		}
	} else {
		refSnapshot := "{}"
		seoReq.RefSnapshotJSON = &refSnapshot
	}

	if strings.TrimSpace(seoReq.OriginURL) == "" {
		seoReq.OriginURL = "/seo/" + strings.Trim(strings.TrimSpace(seoReq.Slug), "/")
	}
	if seoReq.CanonicalURL == nil || strings.TrimSpace(*seoReq.CanonicalURL) == "" {
		canonical := seoReq.OriginURL
		seoReq.CanonicalURL = &canonical
	}

	created, err := h.seoDomainUsecase.CreateSeoDomain(ctx, seoReq)
	if err != nil {
		return nil, err
	}
	return h.seoDomainMapper.SeoDomainToPb(created), nil
}

func (h *SeoDomainHandler) PublishSeoPage(ctx context.Context, req *crmpb.SeoPageCommandRequest) (*crmpb.SeoDomainResponse, error) {
	if req == nil || req.Id == 0 {
		return nil, _errors.ReturnError(service.IDInvalid)
	}

	updated, err := h.seoDomainUsecase.PublishSeoPage(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return h.seoDomainMapper.SeoDomainToPb(updated), nil
}

func (h *SeoDomainHandler) UnpublishSeoPage(ctx context.Context, req *crmpb.SeoPageCommandRequest) (*crmpb.SeoDomainResponse, error) {
	if req == nil || req.Id == 0 {
		return nil, _errors.ReturnError(service.IDInvalid)
	}

	updated, err := h.seoDomainUsecase.UnpublishSeoPage(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return h.seoDomainMapper.SeoDomainToPb(updated), nil
}

func (h *SeoDomainHandler) ArchiveSeoPage(ctx context.Context, req *crmpb.SeoPageCommandRequest) (*crmpb.SeoDomainResponse, error) {
	if req == nil || req.Id == 0 {
		return nil, _errors.ReturnError(service.IDInvalid)
	}

	updated, err := h.seoDomainUsecase.ArchiveSeoPage(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return h.seoDomainMapper.SeoDomainToPb(updated), nil
}

func (h *SeoDomainHandler) RestoreSeoPage(ctx context.Context, req *crmpb.SeoPageCommandRequest) (*crmpb.SeoDomainResponse, error) {
	if req == nil || req.Id == 0 {
		return nil, _errors.ReturnError(service.IDInvalid)
	}

	updated, err := h.seoDomainUsecase.RestoreSeoPage(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return h.seoDomainMapper.SeoDomainToPb(updated), nil
}

// -----------------------------------------------------------------------------
// Render/generate command placeholders.
// Real implementation should move into SeoRenderUsecase after render subsystem.
// -----------------------------------------------------------------------------

func (h *SeoDomainHandler) PreviewSeoPage(ctx context.Context, req *crmpb.RenderSeoPreviewRequest) (*crmpb.RenderSeoPreviewResponse, error) {
	if h.seoRenderUsecase == nil {
		return nil, _errors.ReturnError(service.SEORenderUsecaseNotConfigured)
	}

	var refID *uint64
	if req.RefId != nil {
		v := req.GetRefId()
		refID = &v
	}

	resp, err := h.seoRenderUsecase.PreviewSeoPage(ctx, &dto.RenderSeoPreviewRequest{
		SourceKey:    req.GetSourceKey(),
		RefID:        refID,
		Slug:         req.GetSlug(),
		Title:        req.GetTitle(),
		Description:  req.GetDescription(),
		CanonicalURL: req.GetCanonicalUrl(),
		Content:      req.GetContent(),
		Summary:      req.GetSummary(),
		Metadata:     req.GetMetadata(),
		TemplateKey:  req.GetTemplateKey(),
	})
	if err != nil {
		return nil, err
	}

	return &crmpb.RenderSeoPreviewResponse{
		Html:       resp.HTML,
		StaticPath: resp.StaticPath,
		Hash:       resp.Hash,
		Warnings:   resp.Warnings,
	}, nil
}

func (h *SeoDomainHandler) PreviewExistingSeoPage(ctx context.Context, req *crmpb.RenderExistingSeoPreviewRequest) (*crmpb.RenderSeoPreviewResponse, error) {
	if h.seoRenderUsecase == nil {
		return nil, _errors.ReturnError(service.SEORenderUsecaseNotConfigured)
	}

	resp, err := h.seoRenderUsecase.PreviewExistingSeoPage(ctx, &dto.RenderExistingSeoPreviewRequest{
		ID:          req.GetId(),
		Content:     req.Content,
		TemplateKey: req.TemplateKey,
	})
	if err != nil {
		return nil, err
	}

	return &crmpb.RenderSeoPreviewResponse{
		Html:       resp.HTML,
		StaticPath: resp.StaticPath,
		Hash:       resp.Hash,
		Warnings:   resp.Warnings,
	}, nil
}

func (h *SeoDomainHandler) GenerateSeoPage(ctx context.Context, req *crmpb.GenerateSeoPageRequest) (*crmpb.GenerateSeoPageResponse, error) {
	if h.seoRenderUsecase == nil {
		return nil, _errors.ReturnError(service.SEORenderUsecaseNotConfigured)
	}

	resp, err := h.seoRenderUsecase.GenerateSeoPage(ctx, &dto.GenerateSeoPageRequest{
		ID:          req.GetId(),
		TriggerType: req.GetTriggerType(),
		Reason:      req.GetReason(),
	})
	if err != nil {
		return nil, err
	}

	return &crmpb.GenerateSeoPageResponse{
		Id:             resp.ID,
		RenderStatus:   resp.RenderStatus,
		StaticHtmlPath: resp.StaticHTMLPath,
		StaticHtmlHash: resp.StaticHTMLHash,
		GeneratedAt:    resp.GeneratedAt,
		ErrorMessage:   resp.ErrorMessage,
	}, nil
}

func (h *SeoDomainHandler) GetSeoGenerationLogs(ctx context.Context, req *crmpb.GetSeoGenerationLogsRequest) (*crmpb.GetSeoGenerationLogsResponse, error) {
	if h.seoRenderUsecase == nil {
		return nil, _errors.ReturnError(service.SEORenderUsecaseNotConfigured)
	}
	if req == nil {
		req = &crmpb.GetSeoGenerationLogsRequest{}
	}

	seoDomainID := req.GetSeoDomainId()
	listReq := &dto.SeoGenerationLogListRequest{
		Pagable: _dto.Pagable{
			Page: req.GetPage(),
			Size: req.GetSize(),
		},
		SeoDomainID: &seoDomainID,
	}
	listReq.Normalize()

	items, total, err := h.seoRenderUsecase.GetSeoGenerationLogs(ctx, listReq)
	if err != nil {
		return nil, err
	}

	return &crmpb.GetSeoGenerationLogsResponse{
		Data:  h.seoDomainMapper.SeoGenerationLogListToPb(items),
		Total: total,
	}, nil
}

// -----------------------------------------------------------------------------
// Existing sitemap / parcel generation / internal render APIs.
// -----------------------------------------------------------------------------

// GetSitemap trả XML sitemap public cho Google Search Console.
// @Summary Lấy sitemap XML
// @Description Trả sitemap.xml từ SEO domain public.
// @Tags SEO Domain
// @Accept json
// @Produce xml
// @Param baseUrl query string false "Base URL để ghép khi canonicalUrl/originUrl là path tương đối"
// @Router /v2/crm/sitemap.xml [get]
func (h *SeoDomainHandler) GetSitemap(ctx context.Context, req *crmpb.GetSitemapRequest) (*sharepb.CommonResponse, error) {
	xmlData, err := h.seoDomainUsecase.GetSitemapXML(ctx, &dto.SitemapRequest{
		BaseURL: req.BaseUrl,
	})
	if err != nil {
		return nil, err
	}

	return &sharepb.CommonResponse{
		Code:    8501,
		Message: "application/xml; charset=utf-8",
		Data:    xmlData,
	}, nil
}

// CreateSeoFromParcel tạo SEO domain từ parcel.
// @Summary Tạo SEO từ parcel
// @Description Sinh URL từ cấu hình seo.parcelUrlTemplate, replace {id} theo parcelId và {slug} theo adr_search.
// @Tags SEO Domain
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateSeoFromParcelRequest true "Parcel ID"
// @Router /v2/crm/admin/seo-domains/parcel [post]
func (h *SeoDomainHandler) CreateSeoFromParcel(ctx context.Context, req *crmpb.CreateSeoFromParcelRequest) (*crmpb.SeoDomainResponse, error) {
	seoDomain, err := h.seoDomainUsecase.CreateSeoFromParcel(ctx, &dto.CreateSeoFromParcelRequest{
		ParcelID: req.ParcelId,
	})
	if err != nil {
		return nil, err
	}
	return h.seoDomainMapper.SeoDomainToPb(seoDomain), nil
}

// GenerateSeoFromParcels tự động tạo SEO domain cho các parcel có bật isSeo và chưa có seoId.
// @Summary Tự động tạo SEO từ parcel
// @Description Quét parcels.is_seo = true, seo_id trống, sinh URL theo seo.parcelUrlTemplate rồi tạo seo_domain.
// @Tags SEO Domain
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.GenerateSeoFromParcelsRequest true "Giới hạn số parcel, 0 là toàn bộ"
// @Router /v2/crm/admin/seo-domains/parcels/generate [post]
func (h *SeoDomainHandler) GenerateSeoFromParcels(ctx context.Context, req *crmpb.GenerateSeoFromParcelsRequest) (*crmpb.GenerateSeoFromParcelsResponse, error) {
	items, result, err := h.seoDomainUsecase.GenerateSeoFromParcels(ctx, &dto.GenerateSeoFromParcelsRequest{
		Limit: req.Limit,
	})
	if err != nil {
		return nil, err
	}

	return &crmpb.GenerateSeoFromParcelsResponse{
		Data:    h.seoDomainMapper.SeoDomainListToPb(items),
		Total:   result.Total,
		Created: result.Created,
		Skipped: result.Skipped,
		Failed:  result.Failed,
		Errors:  result.Errors,
	}, nil
}

func (h *SeoDomainHandler) TouchSeoDomainByRef(ctx context.Context, req *crmpb.TouchSeoDomainByRefRequest) (*sharepb.SubmitResponse, error) {
	if err := h.seoDomainUsecase.TouchByRef(ctx, &dto.TouchSeoDomainByRefRequest{RefType: req.RefType, RefID: req.RefId, RefSource: req.RefSource}); err != nil {
		return nil, err
	}
	return &sharepb.SubmitResponse{Id: req.RefId, Message: "Seo domain marked for regeneration"}, nil
}

func (h *SeoDomainHandler) GetSeoSitemap(ctx context.Context, req *crmpb.GetSeoSitemapRequest) (*crmpb.SeoSitemapXmlResponse, error) {
	xmlText, err := h.seoDomainUsecase.BuildSitemapXML(ctx, &dto.SeoSitemapRequest{Scope: req.Scope}, config.AppProperties.Seo.PublicBaseURL)
	if err != nil {
		return nil, err
	}
	return &crmpb.SeoSitemapXmlResponse{
		Xml:         xmlText,
		ContentType: "application/xml; charset=utf-8",
	}, nil
}

func (h *SeoDomainHandler) GetInternalSeoRenderData(ctx context.Context, req *crmpb.GetInternalSeoRenderDataRequest) (*crmpb.InternalSeoRenderDataResponse, error) {
	if !utils.ValidateSeoInternalSecret(ctx) {
		return nil, _errors.ReturnError(service.SEOSecretKeyInvalid)
	}

	seoDomain, err := h.seoDomainUsecase.GetInternalRenderData(ctx, &dto.InternalSeoRenderDataRequest{RefType: req.RefType, RefID: req.RefId})
	if err != nil {
		return nil, err
	}
	pb := h.seoDomainMapper.SeoDomainToPb(seoDomain)
	return &crmpb.InternalSeoRenderDataResponse{
		SeoDomain:     pb,
		InternalLinks: pb.InternalLinks,
		Relatives:     pb.Relatives,
	}, nil
}

// -----------------------------------------------------------------------------
// Internal links.
// -----------------------------------------------------------------------------

func (h *SeoDomainHandler) CreateSeoInternalLink(ctx context.Context, req *crmpb.CreateSeoInternalLinkRequest) (*crmpb.SeoInternalLinkResponse, error) {
	createReq := &dto.SeoInternalLinkCreateRequest{
		Link:        req.Link,
		Title:       req.Title,
		LinkType:    req.LinkType,
		Priority:    req.Priority,
		ParentSeoID: req.ParentSeoId,
		ChildSeoID:  req.ChildSeoId,
	}
	created, err := h.seoDomainUsecase.CreateSeoInternalLink(ctx, createReq)
	if err != nil {
		return nil, err
	}
	return h.seoDomainMapper.SeoInternalLinkToPb(created), nil
}

func (h *SeoDomainHandler) GetSeoInternalLinkList(ctx context.Context, req *crmpb.GetSeoInternalLinkListRequest) (*crmpb.GetSeoInternalLinkListResponse, error) {
	listReq := &dto.SeoInternalLinkListRequest{
		Pagable:     _dto.Pagable{Page: req.Page, Size: req.Size},
		ParentSeoID: req.ParentSeoId,
		ChildSeoID:  req.ChildSeoId,
		LinkType:    req.LinkType,
		Keyword:     req.Keyword,
	}
	items, total, err := h.seoDomainUsecase.GetSeoInternalLinkList(ctx, listReq)
	if err != nil {
		return nil, err
	}
	return &crmpb.GetSeoInternalLinkListResponse{
		Data:  h.seoDomainMapper.SeoInternalLinkListToPb(items),
		Total: total,
	}, nil
}

func (h *SeoDomainHandler) UpdateSeoInternalLink(ctx context.Context, req *crmpb.UpdateSeoInternalLinkRequest) (*crmpb.UpdateSeoInternalLinkResponse, error) {
	updateReq := &dto.SeoInternalLinkUpdateRequest{
		Link:        req.Link,
		Title:       req.Title,
		LinkType:    req.LinkType,
		Priority:    req.Priority,
		ParentSeoID: req.ParentSeoId,
		ChildSeoID:  req.ChildSeoId,
	}
	updated, err := h.seoDomainUsecase.UpdateSeoInternalLink(ctx, req.Id, updateReq)
	if err != nil {
		return nil, err
	}
	return &crmpb.UpdateSeoInternalLinkResponse{SeoInternalLink: h.seoDomainMapper.SeoInternalLinkToPb(updated)}, nil
}

func (h *SeoDomainHandler) DeleteSeoInternalLink(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	if err := h.seoDomainUsecase.DeleteSeoInternalLink(ctx, req.Id); err != nil {
		return nil, err
	}
	return &sharepb.SubmitResponse{Id: req.Id, Message: "Seo internal link deleted successfully"}, nil
}

// -----------------------------------------------------------------------------
// SEO relatives.
// -----------------------------------------------------------------------------

func (h *SeoDomainHandler) CreateSeoRelative(ctx context.Context, req *crmpb.CreateSeoRelativeRequest) (*crmpb.SeoRelativeResponse, error) {
	createReq := &dto.SeoRelativeCreateRequest{
		ParentSeoID:  req.ParentSeoId,
		ChildSeoID:   req.ChildSeoId,
		RelationType: req.RelationType,
		Priority:     req.Priority,
	}
	created, err := h.seoDomainUsecase.CreateSeoRelative(ctx, createReq)
	if err != nil {
		return nil, err
	}
	return h.seoDomainMapper.SeoRelativeToPb(created), nil
}

func (h *SeoDomainHandler) GetSeoRelativeList(ctx context.Context, req *crmpb.GetSeoRelativeListRequest) (*crmpb.GetSeoRelativeListResponse, error) {
	listReq := &dto.SeoRelativeListRequest{
		Pagable:      _dto.Pagable{Page: req.Page, Size: req.Size},
		ParentSeoID:  req.ParentSeoId,
		ChildSeoID:   req.ChildSeoId,
		RelationType: req.RelationType,
	}
	items, total, err := h.seoDomainUsecase.GetSeoRelativeList(ctx, listReq)
	if err != nil {
		return nil, err
	}
	return &crmpb.GetSeoRelativeListResponse{
		Data:  h.seoDomainMapper.SeoRelativeListToPb(items),
		Total: total,
	}, nil
}

func (h *SeoDomainHandler) DeleteSeoRelative(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	if err := h.seoDomainUsecase.DeleteSeoRelative(ctx, req.Id); err != nil {
		return nil, err
	}
	return &sharepb.SubmitResponse{Id: req.Id, Message: "Seo relative deleted successfully"}, nil
}

// -----------------------------------------------------------------------------
// Local mapping helpers.
// -----------------------------------------------------------------------------

func seoRefValueToPb(item dto.SeoRefValueResponse) *crmpb.SeoRefValueResponse {
	return &crmpb.SeoRefValueResponse{
		RefType:           item.RefType,
		RefTypeKey:        item.RefTypeKey,
		RefId:             item.RefID,
		Label:             item.Label,
		Description:       item.Description,
		ImageUrl:          item.ImageURL,
		AdminUrl:          item.AdminURL,
		PublicUrl:         item.PublicURL,
		AlreadyHasSeoPage: item.AlreadyHasSeoPage,
		SeoPageId:         item.SeoPageID,
		SourceService:     item.SourceService,
		ResolverKey:       item.ResolverKey,
		SourceStatus:      item.SourceStatus,
		SourceKey:         item.SourceKey,
	}
}

func seoSourceValueToPb(item dto.SeoRefValueResponse) *crmpb.SeoSourceValueResponse {
	return &crmpb.SeoSourceValueResponse{
		SourceKey:         item.SourceKey,
		RefType:           item.RefType,
		RefTypeKey:        item.RefTypeKey,
		RefId:             item.RefID,
		Label:             item.Label,
		Description:       item.Description,
		ImageUrl:          item.ImageURL,
		AdminUrl:          item.AdminURL,
		PublicUrl:         item.PublicURL,
		SourceService:     item.SourceService,
		ResolverKey:       item.ResolverKey,
		SourceStatus:      item.SourceStatus,
		AlreadyHasSeoPage: item.AlreadyHasSeoPage,
		SeoPageId:         item.SeoPageID,
	}
}

func seoRefSnapshotToPb(item *dto.SeoRefSnapshotResponse) *crmpb.SeoRefSnapshotResponse {
	if item == nil {
		return nil
	}
	return &crmpb.SeoRefSnapshotResponse{
		RefType:                  item.RefType,
		RefTypeKey:               item.RefTypeKey,
		RefId:                    item.RefID,
		Label:                    item.Label,
		Description:              item.Description,
		ImageUrl:                 item.ImageURL,
		AdminUrl:                 item.AdminURL,
		PublicUrl:                item.PublicURL,
		SuggestedSlug:            item.SuggestedSlug,
		SuggestedTitle:           item.SuggestedTitle,
		SuggestedMetaTitle:       item.SuggestedMetaTitle,
		SuggestedMetaDescription: item.SuggestedMetaDescription,
		SourceService:            item.SourceService,
		ResolverKey:              item.ResolverKey,
		Hash:                     item.Hash,
		DataJson:                 item.DataJSON,
		SourceStatus:             item.SourceStatus,
		SourceKey:                item.SourceKey,
		Visibility:               item.Visibility,
		Indexable:                item.Indexable,
		DisabledReason:           item.DisabledReason,
	}
}

func seoSourceSnapshotToPb(item *dto.SeoRefSnapshotResponse) *crmpb.SeoSourceSnapshotResponse {
	if item == nil {
		return nil
	}
	return &crmpb.SeoSourceSnapshotResponse{
		SourceKey:                item.SourceKey,
		RefType:                  item.RefType,
		RefTypeKey:               item.RefTypeKey,
		RefId:                    item.RefID,
		Label:                    item.Label,
		Description:              item.Description,
		ImageUrl:                 item.ImageURL,
		AdminUrl:                 item.AdminURL,
		PublicUrl:                item.PublicURL,
		SuggestedSlug:            item.SuggestedSlug,
		SuggestedTitle:           item.SuggestedTitle,
		SuggestedMetaTitle:       item.SuggestedMetaTitle,
		SuggestedMetaDescription: item.SuggestedMetaDescription,
		SourceService:            item.SourceService,
		ResolverKey:              item.ResolverKey,
		Hash:                     item.Hash,
		DataJson:                 item.DataJSON,
		SourceStatus:             item.SourceStatus,
		Visibility:               item.Visibility,
		Indexable:                item.Indexable,
		DisabledReason:           item.DisabledReason,
	}
}

func buildCreateSeoDraftRequest(req *crmpb.CreateSeoDraftRequest, refType uint32, refSource string) *dto.SeoDomainCreateRequest {
	seo := req.Seo

	canonical := strings.TrimSpace(seo.CanonicalUrl)
	var canonicalPtr *string
	if canonical != "" {
		canonicalPtr = &canonical
	}

	scope := strings.TrimSpace(seo.Scope)
	if scope == "" {
		scope = "public"
	}

	needGenerate := true
	refMissing := false
	sourceStatus := "manual"

	return &dto.SeoDomainCreateRequest{
		Slug:            strings.TrimSpace(seo.Slug),
		OriginURL:       strings.TrimSpace(seo.OriginUrl),
		CanonicalURL:    canonicalPtr,
		RefType:         refType,
		RefSource:       refSource,
		SourceStatus:    sourceStatus,
		RefMissing:      &refMissing,
		Scope:           scope,
		Title:           strings.TrimSpace(seo.Title),
		Description:     strings.TrimSpace(seo.Description),
		Content:         seo.Content,
		Summary:         seo.Summary,
		Published:       false,
		IsSiteMap:       seo.IsSiteMap,
		IsIndex:         seo.IsIndex,
		IsRobot:         seo.IsRobot,
		SitemapPriority: 0.5,
		NeedGenerate:    &needGenerate,
		Metadata:        optionalStringFromTrimmed(seo.Metadata),
		TemplateKey:     strings.TrimSpace(seo.TemplateKey),
	}
}

func optionalStringFromTrimmed(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value != "" {
			return value
		}
	}
	return ""
}
