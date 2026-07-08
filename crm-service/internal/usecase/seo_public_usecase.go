package usecase

import (
	_errors "common/errors"
	"context"
	"strings"

	seo_domain "crm/internal/domain/seo"
	"crm/internal/dto"
	"crm/internal/repo"
)

type SeoPublicUsecase struct {
	seoDomainRepo repo.SeoDomainRepo
}

func NewSeoPublicUsecase(seoDomainRepo repo.SeoDomainRepo) *SeoPublicUsecase {
	return &SeoPublicUsecase{seoDomainRepo: seoDomainRepo}
}

func (u *SeoPublicUsecase) GetPublicSeoPageBySlug(ctx context.Context, slug string) (*dto.SeoPublicPageResponse, error) {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return nil, _errors.ReturnError(400, "slug không hợp lệ")
	}

	page, err := u.seoDomainRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	return buildPublicSeoPageResponse(page)
}

func (u *SeoPublicUsecase) GetPublicSeoPageByCanonicalURL(ctx context.Context, canonicalURL string) (*dto.SeoPublicPageResponse, error) {
	canonicalURL = strings.TrimSpace(canonicalURL)
	if canonicalURL == "" {
		return nil, _errors.ReturnError(400, "canonical URL không hợp lệ")
	}

	page, err := u.seoDomainRepo.GetByCanonicalURL(ctx, canonicalURL)
	if err != nil {
		return nil, err
	}

	return buildPublicSeoPageResponse(page)
}

func buildPublicSeoPageResponse(page *seo_domain.SeoDomain) (*dto.SeoPublicPageResponse, error) {
	decision := seo_domain.CanServePublic(page)
	if !decision.Allowed {
		switch decision.Code {
		case "page_archived":
			return nil, _errors.ReturnError(410, decision.Message)
		case "page_nil":
			return nil, _errors.ReturnError(404, decision.Message)
		case "page_not_published", "not_public_scope":
			return nil, _errors.ReturnError(404, decision.Message)
		default:
			return nil, _errors.ReturnError(503, decision.Message)
		}
	}

	noIndex := !page.IsIndex

	return &dto.SeoPublicPageResponse{
		ID:             page.ID,
		Slug:           page.Slug,
		CanonicalURL:   page.CanonicalURL,
		HTML:           page.RenderedHTML,
		ETag:           page.StaticHtmlHash,
		CacheControl:   buildSeoCacheControl(page),
		StatusCode:     200,
		ContentType:    "text/html; charset=utf-8",
		NoIndex:        noIndex,
		RenderStatus:   page.RenderStatus,
		PageStatus:     page.PageStatus,
		Warnings:       decision.Warnings,
		StaticHTMLPath: page.StaticHtmlPath,
	}, nil
}

func buildSeoCacheControl(page *seo_domain.SeoDomain) string {
	if page == nil {
		return "no-store"
	}
	if page.RenderStatus == seo_domain.SeoRenderStatusSuccess && page.StaticHtmlHash != "" {
		return "public, max-age=300, stale-while-revalidate=3600"
	}
	return "no-cache"
}
