package handler

import (
	"context"
	"crm/infra/mapper"
	"crm/internal/dto"
	"crm/internal/usecase"

	_dto "common/domain/dto"
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

// CreateSeoDomain tạo SEO domain mới.
// @Summary Tạo SEO domain
// @Description Tạo một SEO domain/page record.
// @Tags SEO Domain
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.SeoDomainCreateRequest true "Thông tin SEO domain"
// @Router /v2/crm/admin/seo-domains [post]
func (h *SeoDomainHandler) CreateSeoDomain(ctx context.Context, req *crmpb.CreateSeoDomainRequest) (*crmpb.SeoDomainResponse, error) {
	seoReq := &dto.SeoDomainCreateRequest{
		OriginURL:       req.OriginUrl,
		CanonicalURL:    req.CanonicalUrl,
		Classify:        req.Classify,
		SiteMapLastedAt: req.SiteMapLastedAt,
		IsSiteMap:       req.IsSiteMap,
		IsIndex:         req.IsIndex,
		IsRobot:         req.IsRobot,
		Metadata:        req.Metadata,
		Note:            req.Note,
		RefType:         req.RefType,
		RefID:           req.RefId,
	}

	created, err := h.seoDomainUsecase.CreateSeoDomain(ctx, seoReq)
	if err != nil {
		return nil, err
	}
	return h.seoDomainMapper.SeoDomainToPb(created), nil
}

// GetSeoDomain lấy chi tiết SEO domain.
// @Summary Lấy SEO domain theo ID
// @Tags SEO Domain
// @Accept json
// @Produce json
// @Param id path uint64 true "SEO domain ID"
// @Router /v2/crm/admin/seo-domains/{id} [get]
func (h *SeoDomainHandler) GetSeoDomain(ctx context.Context, req *sharepb.IdRequest) (*crmpb.SeoDomainResponse, error) {
	seoDomain, err := h.seoDomainUsecase.GetSeoDomain(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return h.seoDomainMapper.SeoDomainToPb(seoDomain), nil
}

// GetSeoDomainList lấy danh sách SEO domain.
// @Summary Lấy danh sách SEO domain
// @Tags SEO Domain
// @Accept json
// @Produce json
// @Param page query int false "Trang"
// @Param size query int false "Kích thước trang"
// @Router /v2/crm/admin/seo-domains [get]
func (h *SeoDomainHandler) GetSeoDomainList(ctx context.Context, req *crmpb.GetSeoDomainListRequest) (*crmpb.GetSeoDomainListResponse, error) {
	listReq := &dto.SeoDomainListRequest{
		Pagable: _dto.Pagable{
			Page: req.Page,
			Size: req.Size,
		},
		Keyword:   req.Keyword,
		Classify:  req.Classify,
		IsSiteMap: req.IsSiteMap,
		IsIndex:   req.IsIndex,
		IsRobot:   req.IsRobot,
		RefType:   req.RefType,
		RefID:     req.RefId,
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

// UpdateSeoDomain cập nhật SEO domain.
// @Summary Cập nhật SEO domain
// @Tags SEO Domain
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path uint64 true "SEO domain ID"
// @Param request body dto.SeoDomainUpdateRequest true "Thông tin cập nhật"
// @Router /v2/crm/admin/seo-domains/{id} [put]
func (h *SeoDomainHandler) UpdateSeoDomain(ctx context.Context, req *crmpb.UpdateSeoDomainRequest) (*crmpb.UpdateSeoDomainResponse, error) {
	updateReq := &dto.SeoDomainUpdateRequest{
		OriginURL:       req.OriginUrl,
		CanonicalURL:    req.CanonicalUrl,
		Classify:        req.Classify,
		SiteMapLastedAt: req.SiteMapLastedAt,
		IsSiteMap:       req.IsSiteMap,
		IsIndex:         req.IsIndex,
		IsRobot:         req.IsRobot,
		Metadata:        req.Metadata,
		Note:            req.Note,
		RefType:         req.RefType,
		RefID:           req.RefId,
	}

	updated, err := h.seoDomainUsecase.UpdateSeoDomain(ctx, req.Id, updateReq)
	if err != nil {
		return nil, err
	}
	return &crmpb.UpdateSeoDomainResponse{
		SeoDomain: h.seoDomainMapper.SeoDomainToPb(updated),
	}, nil
}

// DeleteSeoDomain xóa SEO domain.
// @Summary Xóa SEO domain
// @Tags SEO Domain
// @Security BearerAuth
// @Param id path uint64 true "SEO domain ID"
// @Router /v2/crm/admin/seo-domains/{id} [delete]
func (h *SeoDomainHandler) DeleteSeoDomain(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	if err := h.seoDomainUsecase.DeleteSeoDomain(ctx, req.Id); err != nil {
		return nil, err
	}
	return &sharepb.SubmitResponse{
		Id:      req.Id,
		Message: "Seo domain deleted successfully",
	}, nil
}

// CreateSeoInternalLink tạo internal link ParentSeo -> ChildSeo.
// @Summary Tạo SEO internal link
// @Tags SEO Domain
// @Router /v2/crm/admin/seo-domains/internal-links [post]
func (h *SeoDomainHandler) CreateSeoInternalLink(ctx context.Context, req *crmpb.CreateSeoInternalLinkRequest) (*crmpb.SeoInternalLinkResponse, error) {
	createReq := &dto.SeoInternalLinkCreateRequest{
		Link:        req.Link,
		Title:       req.Title,
		ParentSeoID: req.ParentSeoId,
		ChildSeoID:  req.ChildSeoId,
	}
	created, err := h.seoDomainUsecase.CreateSeoInternalLink(ctx, createReq)
	if err != nil {
		return nil, err
	}
	return h.seoDomainMapper.SeoInternalLinkToPb(created), nil
}

// GetSeoInternalLinkList lấy danh sách internal link.
// @Summary Lấy SEO internal links
// @Tags SEO Domain
// @Router /v2/crm/admin/seo-domains/internal-links [get]
func (h *SeoDomainHandler) GetSeoInternalLinkList(ctx context.Context, req *crmpb.GetSeoInternalLinkListRequest) (*crmpb.GetSeoInternalLinkListResponse, error) {
	listReq := &dto.SeoInternalLinkListRequest{
		Pagable: _dto.Pagable{
			Page: req.Page,
			Size: req.Size,
		},
		ParentSeoID: req.ParentSeoId,
		ChildSeoID:  req.ChildSeoId,
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

// UpdateSeoInternalLink cập nhật internal link.
// @Summary Cập nhật SEO internal link
// @Tags SEO Domain
// @Router /v2/crm/admin/seo-domains/internal-links/{id} [put]
func (h *SeoDomainHandler) UpdateSeoInternalLink(ctx context.Context, req *crmpb.UpdateSeoInternalLinkRequest) (*crmpb.UpdateSeoInternalLinkResponse, error) {
	updateReq := &dto.SeoInternalLinkUpdateRequest{
		Link:        req.Link,
		Title:       req.Title,
		ParentSeoID: req.ParentSeoId,
		ChildSeoID:  req.ChildSeoId,
	}
	updated, err := h.seoDomainUsecase.UpdateSeoInternalLink(ctx, req.Id, updateReq)
	if err != nil {
		return nil, err
	}
	return &crmpb.UpdateSeoInternalLinkResponse{
		SeoInternalLink: h.seoDomainMapper.SeoInternalLinkToPb(updated),
	}, nil
}

// DeleteSeoInternalLink xóa internal link.
// @Summary Xóa SEO internal link
// @Tags SEO Domain
// @Router /v2/crm/admin/seo-domains/internal-links/{id} [delete]
func (h *SeoDomainHandler) DeleteSeoInternalLink(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	if err := h.seoDomainUsecase.DeleteSeoInternalLink(ctx, req.Id); err != nil {
		return nil, err
	}
	return &sharepb.SubmitResponse{
		Id:      req.Id,
		Message: "Seo internal link deleted successfully",
	}, nil
}

// CreateSeoRelative tạo quan hệ SEO ParentSeo -> ChildSeo.
// @Summary Tạo SEO relative
// @Tags SEO Domain
// @Router /v2/crm/admin/seo-domains/relatives [post]
func (h *SeoDomainHandler) CreateSeoRelative(ctx context.Context, req *crmpb.CreateSeoRelativeRequest) (*crmpb.SeoRelativeResponse, error) {
	createReq := &dto.SeoRelativeCreateRequest{
		ParentSeoID: req.ParentSeoId,
		ChildSeoID:  req.ChildSeoId,
	}
	created, err := h.seoDomainUsecase.CreateSeoRelative(ctx, createReq)
	if err != nil {
		return nil, err
	}
	return h.seoDomainMapper.SeoRelativeToPb(created), nil
}

// GetSeoRelativeList lấy danh sách SEO relatives.
// @Summary Lấy SEO relatives
// @Tags SEO Domain
// @Router /v2/crm/admin/seo-domains/relatives [get]
func (h *SeoDomainHandler) GetSeoRelativeList(ctx context.Context, req *crmpb.GetSeoRelativeListRequest) (*crmpb.GetSeoRelativeListResponse, error) {
	listReq := &dto.SeoRelativeListRequest{
		Pagable: _dto.Pagable{
			Page: req.Page,
			Size: req.Size,
		},
		ParentSeoID: req.ParentSeoId,
		ChildSeoID:  req.ChildSeoId,
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

// DeleteSeoRelative xóa SEO relative.
// @Summary Xóa SEO relative
// @Tags SEO Domain
// @Router /v2/crm/admin/seo-domains/relatives/{id} [delete]
func (h *SeoDomainHandler) DeleteSeoRelative(ctx context.Context, req *sharepb.IdRequest) (*sharepb.SubmitResponse, error) {
	if err := h.seoDomainUsecase.DeleteSeoRelative(ctx, req.Id); err != nil {
		return nil, err
	}
	return &sharepb.SubmitResponse{
		Id:      req.Id,
		Message: "Seo relative deleted successfully",
	}, nil
}
