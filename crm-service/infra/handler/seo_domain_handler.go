package handler

import (
	_dto "common/domain/dto"
	_errors "common/errors"
	"context"
	"crm/config"
	"crm/infra/mapper"
	"crm/infra/utils"
	"crm/internal/dto"
	"crm/internal/usecase"

	crmpb "pb/types/crm"
	sharepb "pb/types/shared"
)

type SeoDomainHandler struct {
	crmpb.UnimplementedSeoDomainServiceServer
	seoDomainUsecase *usecase.SeoDomainUsecase
	seoDomainMapper  *mapper.SeoDomainMapper
}

// @bind: crm/infra/handler.SeoDomainHandler
func NewSeoDomainHandler(seoDomainUsecase *usecase.SeoDomainUsecase, seoDomainMapper *mapper.SeoDomainMapper) *SeoDomainHandler {
	return &SeoDomainHandler{
		seoDomainUsecase: seoDomainUsecase,
		seoDomainMapper:  seoDomainMapper,
	}
}

func (h *SeoDomainHandler) CreateSeoDomain(ctx context.Context, req *crmpb.CreateSeoDomainRequest) (*crmpb.SeoDomainResponse, error) {
	seoReq := &dto.SeoDomainCreateRequest{
		Slug:              req.Slug,
		OriginURL:         req.OriginUrl,
		CanonicalURL:      req.CanonicalUrl,
		RefType:           req.RefType,
		RefID:             req.RefId,
		Scope:             req.Scope,
		Classify:          req.Classify,
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

func (h *SeoDomainHandler) GetSeoDomainList(ctx context.Context, req *crmpb.GetSeoDomainListRequest) (*crmpb.GetSeoDomainListResponse, error) {
	listReq := &dto.SeoDomainListRequest{
		Pagable: _dto.Pagable{
			Page: req.Page,
			Size: req.Size,
		},
		Keyword:      req.Keyword,
		RefType:      req.RefType,
		RefID:        req.RefId,
		Scope:        req.Scope,
		Classify:     req.Classify,
		Published:    req.Published,
		IsSiteMap:    req.IsSiteMap,
		IsIndex:      req.IsIndex,
		IsRobot:      req.IsRobot,
		NeedGenerate: req.NeedGenerate,
	}

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
		Scope:             req.Scope,
		Classify:          req.Classify,
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

func (h *SeoDomainHandler) TouchSeoDomainByRef(ctx context.Context, req *crmpb.TouchSeoDomainByRefRequest) (*sharepb.SubmitResponse, error) {
	if err := h.seoDomainUsecase.TouchByRef(ctx, &dto.TouchSeoDomainByRefRequest{RefType: req.RefType, RefID: req.RefId}); err != nil {
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
		return nil, _errors.ReturnError(401, "secretKey không hợp lệ")
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
