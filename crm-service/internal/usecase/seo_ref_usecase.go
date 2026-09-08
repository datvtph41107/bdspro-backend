package usecase

import (
	_dto "common/domain/dto"
	"context"
	"strings"

	_errors "common/errors"
	seo_domain "crm/internal/domain/seo"
	"crm/internal/dto"
	"crm/internal/enums"
)

// GetSeoRefTypes is the legacy registry accessor.
// It now returns normalized config so old response fields and V2 fields stay in sync.
func (u *SeoDomainUsecase) GetSeoRefTypes(ctx context.Context) []seo_domain.SeoRefTypeConfig {
	return seo_domain.SeoRefTypeConfigs()
}

// GetSeoSources is the V2 source registry accessor.
// New frontend should prefer this over GetSeoRefTypes.
func (u *SeoDomainUsecase) GetSeoSources(ctx context.Context) []seo_domain.SeoRefTypeConfig {
	return seo_domain.SeoRefTypeConfigs()
}

// GetSeoSource returns one source definition by sourceKey.
// sourceKey is the official V2 identifier used by frontend.
func (u *SeoDomainUsecase) GetSeoSource(ctx context.Context, sourceKey string) (seo_domain.SeoRefTypeConfig, error) {
	cfg, ok := seo_domain.GetSeoRefTypeConfigByKey(sourceKey)
	if !ok {
		return seo_domain.SeoRefTypeConfig{}, _errors.ReturnError(404, "seo source không tồn tại: "+strings.TrimSpace(sourceKey))
	}
	return cfg, nil
}

// SearchSeoRefValues keeps the legacy search behavior for /seo/ref-values,
// but now enforces registry readiness/capability before calling resolver.
func (u *SeoDomainUsecase) SearchSeoRefValues(ctx context.Context, req *dto.SearchSeoRefValuesRequest) ([]dto.SeoRefValueResponse, int64, error) {
	if req == nil {
		return nil, 0, _errors.ReturnError(400, "request không hợp lệ")
	}
	req.Normalize()
	if !enums.IsValidSEORefType(req.RefType) {
		return nil, 0, _errors.ReturnError(400, "refType không hợp lệ")
	}

	cfg, err := resolveSeoRefConfig(req.RefType, req.ResolverKey, req.SourceService)
	if err != nil {
		return nil, 0, err
	}
	cfg = seo_domain.NormalizeSeoRefTypeConfig(cfg)

	if !seo_domain.IsSeoRefTypeSearchable(cfg) {
		return []dto.SeoRefValueResponse{}, 0, nil
	}
	if !seo_domain.IsSeoResolverReady(cfg.ResolverKey) {
		return nil, 0, resolverNotImplemented(cfg.ResolverKey)
	}

	items, total, err := newSeoRefResolver(cfg, u.bdsproProvider, u.tqdProvider).Search(ctx, req)
	if err != nil {
		return nil, 0, err
	}

	if err := u.attachSeoRefExistingState(ctx, items); err != nil {
		return nil, 0, err
	}

	for i := range items {
		items[i].SourceKey = cfg.Key
	}

	return items, total, nil
}

// SearchSeoSourceValues is the V2 search behavior for /seo/sources/{sourceKey}/values.
func (u *SeoDomainUsecase) SearchSeoSourceValues(ctx context.Context, req *dto.SearchSeoSourceValuesRequest) ([]dto.SeoRefValueResponse, int64, error) {
	if req == nil {
		return nil, 0, _errors.ReturnError(400, "request không hợp lệ")
	}
	req.Normalize()

	cfg, err := u.GetSeoSource(ctx, req.SourceKey)
	if err != nil {
		return nil, 0, err
	}
	cfg = seo_domain.NormalizeSeoRefTypeConfig(cfg)

	if !seo_domain.IsSeoRefTypeSearchable(cfg) {
		return []dto.SeoRefValueResponse{}, 0, nil
	}
	if !seo_domain.IsSeoResolverReady(cfg.ResolverKey) {
		return nil, 0, resolverNotImplemented(cfg.ResolverKey)
	}

	legacyReq := &dto.SearchSeoRefValuesRequest{
		Pagable: _dto.Pagable{
			Text: req.Text,
			Page: req.Page,
			Size: req.Size,
		},
		RefType:       cfg.Value,
		ResolverKey:   cfg.ResolverKey,
		SourceService: cfg.SourceService,
	}

	items, total, err := newSeoRefResolver(cfg, u.bdsproProvider, u.tqdProvider).Search(ctx, legacyReq)
	if err != nil {
		return nil, 0, err
	}

	if err := u.attachSeoRefExistingState(ctx, items); err != nil {
		return nil, 0, err
	}

	for i := range items {
		items[i].SourceKey = cfg.Key
	}

	return items, total, nil
}

// ResolveSeoRefValue keeps the legacy resolve behavior for /seo/ref-values/{refType}/{refId},
// but now enforces registry readiness/capability before calling resolver.
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
	cfg = seo_domain.NormalizeSeoRefTypeConfig(cfg)

	if !seo_domain.IsSeoRefTypeResolvable(cfg) {
		return nil, resolverNotImplemented(cfg.ResolverKey)
	}
	if !seo_domain.IsSeoResolverReady(cfg.ResolverKey) {
		return nil, resolverNotImplemented(cfg.ResolverKey)
	}

	item, err := newSeoRefResolver(cfg, u.bdsproProvider, u.tqdProvider).Resolve(ctx, req)
	if err != nil {
		return nil, err
	}

	item.SourceKey = cfg.Key
	item.Visibility = inferSeoSnapshotVisibility(item)
	item.Indexable = inferSeoSnapshotIndexable(item)
	item.DisabledReason = cfg.DisabledReason

	return item, nil
}

// ResolveSeoSourceValue is the V2 resolve behavior for /seo/sources/{sourceKey}/values/{refId}.
func (u *SeoDomainUsecase) ResolveSeoSourceValue(ctx context.Context, req *dto.ResolveSeoSourceValueRequest) (*dto.SeoRefSnapshotResponse, error) {
	if req == nil {
		return nil, _errors.ReturnError(400, "request không hợp lệ")
	}
	if req.RefID == 0 {
		return nil, _errors.ReturnError(400, "refId không hợp lệ")
	}

	cfg, err := u.GetSeoSource(ctx, req.SourceKey)
	if err != nil {
		return nil, err
	}
	cfg = seo_domain.NormalizeSeoRefTypeConfig(cfg)

	if !seo_domain.IsSeoRefTypeResolvable(cfg) {
		return nil, resolverNotImplemented(cfg.ResolverKey)
	}
	if !seo_domain.IsSeoResolverReady(cfg.ResolverKey) {
		return nil, resolverNotImplemented(cfg.ResolverKey)
	}

	legacyReq := &dto.ResolveSeoRefValueRequest{
		RefType:       cfg.Value,
		RefID:         req.RefID,
		ResolverKey:   cfg.ResolverKey,
		SourceService: cfg.SourceService,
	}

	item, err := newSeoRefResolver(cfg, u.bdsproProvider, u.tqdProvider).Resolve(ctx, legacyReq)
	if err != nil {
		return nil, err
	}

	item.SourceKey = cfg.Key
	item.Visibility = inferSeoSnapshotVisibility(item)
	item.Indexable = inferSeoSnapshotIndexable(item)
	item.DisabledReason = cfg.DisabledReason

	return item, nil
}

func (u *SeoDomainUsecase) attachSeoRefExistingState(ctx context.Context, items []dto.SeoRefValueResponse) error {
	for i := range items {
		existing, err := u.seoDomainRepo.GetByRefSource(ctx, items[i].RefType, items[i].ResolverKey, items[i].RefID)
		if err != nil {
			return err
		}
		if existing != nil {
			items[i].AlreadyHasSeoPage = true
			items[i].SeoPageID = existing.ID
		}
	}
	return nil
}

func inferSeoSnapshotVisibility(item *dto.SeoRefSnapshotResponse) string {
	if item == nil {
		return ""
	}
	// Current resolvers do not yet expose restricted/private source visibility.
	// Default to public; future source adapters should override this using source data.
	return "public"
}

func inferSeoSnapshotIndexable(item *dto.SeoRefSnapshotResponse) bool {
	if item == nil {
		return false
	}
	return inferSeoSnapshotVisibility(item) == "public"
}

// GetSeoRefs trả danh sách SEO page public được liên kết qua bảng seo_refs.
func (u *SeoDomainUsecase) GetSeoRefs(ctx context.Context, req *dto.GetSeoRefsRequest) ([]seo_domain.SeoDomain, int64, error) {
	if req == nil {
		req = &dto.GetSeoRefsRequest{}
	}

	refType := req.RefType.Validate()
	if !refType {
		return nil, 0, _errors.ReturnError(400, "refType không hợp lệ")
	}
	if req.RefID == 0 {
		return nil, 0, _errors.ReturnError(400, "refId không hợp lệ")
	}

	req.Normalize()

	return u.seoDomainRepo.GetListByLinkedRef(ctx, req)
}
