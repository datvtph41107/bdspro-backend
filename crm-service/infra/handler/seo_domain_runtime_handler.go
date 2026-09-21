package handler

import (
	_errors "common/errors"
	"context"
	"crm/infra/utils"
	"crm/internal"

	crmpb "pb/types/crm"
)

func (h *SeoDomainHandler) GetInternalSeoRenderDataBySlug(ctx context.Context, req *crmpb.GetInternalSeoRenderDataBySlugRequest) (*crmpb.InternalSeoRenderDataResponse, error) {
	if !utils.ValidateSeoInternalSecret(ctx) {
		return nil, _errors.ReturnError(service.SEOSecretKeyInvalid)
	}

	data, err := h.seoDomainUsecase.GetInternalSeoRenderDataBySlug(ctx, req.Slug)
	if err != nil {
		return nil, err
	}

	return &crmpb.InternalSeoRenderDataResponse{
		SeoDomain:     h.seoDomainMapper.SeoDomainToPb(data.SeoDomain),
		InternalLinks: h.seoDomainMapper.SeoInternalLinkListToPb(data.InternalLinks),
		Relatives:     h.seoDomainMapper.SeoRelativeListToPb(data.Relatives),
	}, nil
}

func (h *SeoDomainHandler) GetInternalSeoStaticPaths(ctx context.Context, req *crmpb.GetInternalSeoStaticPathsRequest) (*crmpb.GetInternalSeoStaticPathsResponse, error) {
	if !utils.ValidateSeoInternalSecret(ctx) {
		return nil, _errors.ReturnError(service.SEOSecretKeyInvalid)
	}

	items, total, err := h.seoDomainUsecase.GetInternalSeoStaticPaths(ctx, req.Page, req.Size)
	if err != nil {
		return nil, err
	}

	return &crmpb.GetInternalSeoStaticPathsResponse{
		Data:  h.seoDomainMapper.SeoDomainStaticPathListToPb(items),
		Total: total,
	}, nil
}
