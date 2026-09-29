package usecase

import (
	_errors "common/errors"
	"context"
	"crm/internal"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
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
		return nil, _errors.ReturnError(service.SEOSlugInvalid)
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
		return nil, _errors.ReturnError(service.SEOCanonicalURLInvalid)
	}

	page, err := u.seoDomainRepo.GetByCanonicalURL(ctx, canonicalURL)
	if err != nil {
		return nil, err
	}

	return buildPublicSeoPageResponse(page)
}

func buildPublicSeoPageResponse(page *seo_domain.SeoDomain) (*dto.SeoPublicPageResponse, error) {
	if err := validatePublicSeoPage(page); err != nil {
		return nil, err
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
		Warnings:       nil,
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

func (u *SeoPublicUsecase) GetPublicSeoDocument(
	ctx context.Context,
	moduleID string,
	slug string,
) (*dto.SeoPublicDocumentResponse, error) {
	moduleID = normalizeSeoModuleID(moduleID)
	if moduleID == "" || moduleID == "default" {
		return nil, _errors.ReturnError(service.SEOModuleUnsupported)
	}
	slug, err := normalizePublicSeoSlug(slug)
	if err != nil {
		return nil, err
	}

	canonicalPath, err := publicCanonicalPath(moduleID, slug)
	if err != nil {
		return nil, err
	}
	page, err := u.resolvePublicDocumentPage(ctx, canonicalPath, slug)
	if err != nil {
		return nil, err
	}
	if err := validatePublicSeoPage(page); err != nil {
		return nil, err
	}

	generatedAt := page.UpdatedAt.UTC()
	if page.GeneratedAt != nil {
		generatedAt = page.GeneratedAt.UTC()
	}
	document, warnings, err := composeSeoPublicDocument(page, generatedAt)
	if err != nil {
		return nil, fmt.Errorf("compose SEO public document: %w", err)
	}
	if document.ModuleID != moduleID {
		return nil, _errors.ReturnError(service.SEOModuleContractMismatch, _errors.WithPublicMessage(fmt.Sprintf("SEO module contract mismatch: expected %s, got %s", moduleID, document.ModuleID)))
	}
	if document.Identity.CanonicalPath != canonicalPath || document.SEO.CanonicalPath != canonicalPath {
		return nil, _errors.ReturnError(service.SEOCanonicalContractMismatch)
	}
	payload, err := json.Marshal(document)
	if err != nil {
		return nil, fmt.Errorf("encode public SEO document: %w", err)
	}
	hash := sha256.Sum256(payload)
	return &dto.SeoPublicDocumentResponse{
		Document:     document,
		ETag:         hex.EncodeToString(hash[:]),
		CacheControl: buildSeoCacheControl(page),
		StatusCode:   200,
		Warnings:     compactUniqueStrings(warnings),
	}, nil
}

func validatePublicSeoPage(page *seo_domain.SeoDomain) error {
	if page == nil {
		return _errors.ReturnError(service.SEOPageNotFound)
	}
	page.NormalizeLifecycle()
	if page.PageStatus == seo_domain.SeoPageStatusArchived {
		return _errors.ReturnError(service.SEOPageArchived)
	}
	if page.PageStatus != seo_domain.SeoPageStatusPublished || !page.Published || page.Scope != seo_domain.SeoScopePublic {
		return _errors.ReturnError(service.SEOPageNotPublic)
	}
	if page.RenderStatus != seo_domain.SeoRenderStatusSuccess || strings.TrimSpace(page.RenderedHTML) == "" {
		return _errors.ReturnError(service.SEOPageRenderIncomplete)
	}
	return nil
}

func (u *SeoPublicUsecase) resolvePublicDocumentPage(
	ctx context.Context,
	canonicalPath string,
	slug string,
) (*seo_domain.SeoDomain, error) {
	for _, candidate := range compactUniqueStrings([]string{
		canonicalPath,
		strings.TrimPrefix(canonicalPath, "/"),
		slug,
	}) {
		page, err := u.seoDomainRepo.GetBySlug(ctx, candidate)
		if err != nil {
			return nil, err
		}
		if page != nil {
			return page, nil
		}
	}
	page, err := u.seoDomainRepo.GetByCanonicalURL(ctx, canonicalPath)
	if err != nil {
		return nil, err
	}
	if page == nil {
		return nil, _errors.ReturnError(service.SEOPublicPageNotFound)
	}
	return page, nil
}

func normalizePublicSeoSlug(value string) (string, error) {
	slug := strings.TrimSpace(value)
	lower := strings.ToLower(slug)
	if slug == "" ||
		strings.ContainsAny(slug, `/\`) ||
		strings.Contains(slug, "..") ||
		strings.ContainsAny(slug, "?#") ||
		lower == "undefined" ||
		lower == "null" ||
		strings.Contains(lower, "[object object]") {
		return "", _errors.ReturnError(service.SEOSlugInvalid)
	}
	return slug, nil
}

func publicCanonicalPath(moduleID string, slug string) (string, error) {
	prefixes := map[string]string{
		"administrative-unit": "/dia-ban",
		"parcel":              "/thua-dat",
		"planning-region":     "/vung-quy-hoach",
		"planning-project":    "/do-an-quy-hoach",
		"planning-news":       "/quy-hoach/tin-tuc",
		"planning-report":     "/quy-hoach/bao-cao",
	}
	prefix := prefixes[moduleID]
	if prefix == "" {
		return "", _errors.ReturnError(service.SEOPublicRouteContractMissing)
	}
	return prefix + "/" + strings.Trim(slug, "/"), nil
}
