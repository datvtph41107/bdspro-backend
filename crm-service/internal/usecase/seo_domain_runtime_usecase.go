package usecase

import (
	_errors "common/errors"
	"context"
	"crm/internal"
	seo_domain "crm/internal/domain/seo"
	"crm/internal/dto"
	"strings"
)

func (u *SeoDomainUsecase) GetInternalSeoRenderDataBySlug(ctx context.Context, slug string) (*dto.InternalSeoRenderDataResponse, error) {
	slug = strings.Trim(strings.TrimSpace(slug), "/")
	if slug == "" {
		return nil, _errors.ReturnError(service.SEOSlugRequired)
	}

	seoDomain, err := u.seoDomainRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if seoDomain == nil {
		return nil, _errors.ReturnError(service.SEODomainNotFound)
	}
	if !seoDomain.Published {
		return nil, _errors.ReturnError(service.SEODomainUnpublished)
	}
	if !seoDomain.IsIndex {
		return nil, _errors.ReturnError(service.SEODomainIndexDenied)
	}

	return &dto.InternalSeoRenderDataResponse{
		SeoDomain:     seoDomain,
		Entity:        nil,
		InternalLinks: seoDomain.InternalLinks,
		Relatives:     seoDomain.Relatives,
	}, nil
}

func (u *SeoDomainUsecase) GetInternalSeoStaticPaths(ctx context.Context, page uint32, size uint32) ([]seo_domain.SeoDomain, int64, error) {
	return u.seoDomainRepo.GetPublishedStaticPaths(ctx, page, size)
}
