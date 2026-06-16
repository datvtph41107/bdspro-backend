package usecase

import (
	_errors "common/errors"
	"context"
	seo_domain "crm/internal/domain/seo"
	"crm/internal/dto"
	"crm/internal/enums"
)

func (u *SeoDomainUsecase) GetSeoRefTypes(ctx context.Context) []seo_domain.SeoRefTypeConfig {
	return seo_domain.SeoRefTypeConfigs()
}

func (u *SeoDomainUsecase) SearchSeoRefValues(ctx context.Context, req *dto.SearchSeoRefValuesRequest) ([]dto.SeoRefValueResponse, int64, error) {
	if req == nil {
		return nil, 0, _errors.ReturnError(400, "request không hợp lệ")
	}
	if !enums.IsValidSEORefType(req.RefType) {
		return nil, 0, _errors.ReturnError(400, "refType không hợp lệ")
	}
	cfg, err := resolveSeoRefConfig(req.RefType, req.ResolverKey, req.SourceService)
	if err != nil {
		return nil, 0, err
	}
	if !cfg.Searchable {
		return []dto.SeoRefValueResponse{}, 0, nil
	}

	items, total, err := newSeoRefResolver(cfg, u.bdsproProvider, u.tqdProvider).Search(ctx, req)
	if err != nil {
		return nil, 0, err
	}
	for i := range items {
		existing, err := u.seoDomainRepo.GetByRefSource(ctx, items[i].RefType, items[i].ResolverKey, items[i].RefID)
		if err != nil {
			return nil, 0, err
		}
		if existing != nil {
			items[i].AlreadyHasSeoPage = true
			items[i].SeoPageID = existing.ID
		}
	}
	return items, total, nil
}

func (u *SeoDomainUsecase) ResolveSeoRefValue(ctx context.Context, req *dto.ResolveSeoRefValueRequest) (*dto.SeoRefSnapshotResponse, error) {
	if req == nil {
		return nil, _errors.ReturnError(400, "request không hợp lệ")
	}
	if !enums.IsValidSEORefType(req.RefType) {
		return nil, _errors.ReturnError(400, "refType không hợp lệ")
	}
	if req.RefID == 0 {
		return nil, _errors.ReturnError(400, "refId không hợp lệ")
	}
	cfg, err := resolveSeoRefConfig(req.RefType, req.ResolverKey, req.SourceService)
	if err != nil {
		return nil, err
	}
	return newSeoRefResolver(cfg, u.bdsproProvider, u.tqdProvider).Resolve(ctx, req)
}
