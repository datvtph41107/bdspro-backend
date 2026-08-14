package usecase

import (
	_errors "common/errors"
	"context"
	seo_domain "crm/internal/domain/seo"
	"crm/internal/dto"
	"strings"
)

func (u *SeoDomainUsecase) GetInternalSeoRenderDataBySlug(ctx context.Context, slug string) (*dto.InternalSeoRenderDataResponse, error) {
	slug = strings.Trim(strings.TrimSpace(slug), "/")
	if slug == "" {
		return nil, _errors.ReturnError(400, "slug không được để trống")
	}

	seoDomain, err := u.seoDomainRepo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if seoDomain == nil {
		return nil, _errors.ReturnError(404, "seo domain không tồn tại")
	}
	if !seoDomain.Published {
		return nil, _errors.ReturnError(404, "seo domain chưa publish")
	}
	if !seoDomain.IsIndex {
		return nil, _errors.ReturnError(404, "seo domain không cho index")
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
