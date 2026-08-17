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

// GetSeoSourceHealth returns lightweight source health based on registry readiness.
// Full production health can later call adapter-specific health checks.
func (u *SeoDomainUsecase) GetSeoSourceHealth(ctx context.Context) []seo_domain.SeoSourceHealth {
	configs := seo_domain.SeoRefTypeConfigs()
	items := make([]seo_domain.SeoSourceHealth, 0, len(configs))

	for _, cfg := range configs {
		cfg = seo_domain.NormalizeSeoRefTypeConfig(cfg)

		status := seo_domain.SeoSourceHealthHealthy
		message := ""

		switch {
		case cfg.Status == seo_domain.SeoSourceStatusComingSoon:
			status = seo_domain.SeoSourceHealthNotImplemented
			message = cfg.DisabledReason
		case cfg.Status == seo_domain.SeoSourceStatusDisabled:
			status = seo_domain.SeoSourceHealthUnavailable
			message = cfg.DisabledReason
		case !cfg.Capabilities.Resolvable && cfg.Key != "manual":
			status = seo_domain.SeoSourceHealthNotImplemented
			message = "source chưa hỗ trợ resolve"
		case cfg.Key != "manual" && !seo_domain.IsSeoResolverReady(cfg.ResolverKey):
			status = seo_domain.SeoSourceHealthNotImplemented
			if message == "" {
				message = "resolver chưa có adapter production"
			}
		}

		items = append(items, seo_domain.SeoSourceHealth{
			SourceKey:      cfg.Key,
			SourceService:  cfg.SourceService,
			ResolverKey:    cfg.ResolverKey,
			Status:         status,
			LatencyMs:      0,
			LastCheckedAt:  "",
			DisabledReason: cfg.DisabledReason,
			ErrorMessage:   message,
		})
	}

	return items
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

// GetSeoRefs trả danh sách SEO page public theo refType/refId:
// page gắn trực tiếp trên seo_domain, cộng page liên kết qua bảng seo_refs.
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

// CountSeoRefs đếm số bài SEO public theo refType/refId, cùng điều kiện với GetSeoRefs.
func (u *SeoDomainUsecase) CountSeoRefs(ctx context.Context, req *dto.GetSeoRefsRequest) (int64, error) {
	if req == nil {
		req = &dto.GetSeoRefsRequest{}
	}

	if !req.RefType.Validate() {
		return 0, _errors.ReturnError(400, "refType không hợp lệ")
	}
	if req.RefID == 0 {
		return 0, _errors.ReturnError(400, "refId không hợp lệ")
	}

	return u.seoDomainRepo.CountByLinkedRef(ctx, req)
}

func (u *SeoDomainUsecase) CreateSeoRef(ctx context.Context, req *dto.SeoRefCreateRequest) (*seo_domain.SeoRef, error) {
	if req == nil {
		return nil, _errors.ReturnError(400, "request không hợp lệ")
	}

	seoDomainID, refType, refID, err := u.prepareSeoRef(ctx, req.SeoDomainID, req.RefType, req.RefID)
	if err != nil {
		return nil, err
	}

	exist, err := u.seoRefRepo.Exist(ctx, seoDomainID, uint32(refType), refID, nil)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, _errors.ReturnError(409, "seo ref đã tồn tại")
	}

	return u.seoRefRepo.Create(ctx, &seo_domain.SeoRef{
		SeoID:   seoDomainID,
		RefType: refType,
		RefID:   refID,
	})
}

func (u *SeoDomainUsecase) GetSeoRef(ctx context.Context, id uint64) (*seo_domain.SeoRef, error) {
	if id == 0 {
		return nil, _errors.ReturnError(400, "id không hợp lệ")
	}
	item, err := u.seoRefRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, _errors.ReturnError(404, "seo ref không tồn tại")
	}
	return item, nil
}

func (u *SeoDomainUsecase) GetSeoRefList(ctx context.Context, req *dto.SeoRefListRequest) ([]seo_domain.SeoRef, int64, error) {
	if req == nil {
		req = &dto.SeoRefListRequest{}
	}
	if !req.RefType.Validate() {
		return nil, 0, _errors.ReturnError(400, "refType không hợp lệ")
	}
	if req.RefID == 0 {
		return nil, 0, _errors.ReturnError(400, "refId không hợp lệ")
	}
	req.Normalize()
	return u.seoRefRepo.GetList(ctx, req)
}

func (u *SeoDomainUsecase) UpdateSeoRef(ctx context.Context, id uint64, req *dto.SeoRefUpdateRequest) (*seo_domain.SeoRef, error) {
	if req == nil {
		return nil, _errors.ReturnError(400, "request không hợp lệ")
	}

	existing, err := u.GetSeoRef(ctx, id)
	if err != nil {
		return nil, err
	}

	seoDomainID := existing.SeoID
	refType := existing.RefType
	refID := existing.RefID
	if req.SeoDomainID != nil {
		seoDomainID = *req.SeoDomainID
	}
	if req.RefType != nil {
		refType = *req.RefType
	}
	if req.RefID != nil {
		refID = *req.RefID
	}

	seoDomainID, refType, refID, err = u.prepareSeoRef(ctx, seoDomainID, refType, refID)
	if err != nil {
		return nil, err
	}

	exist, err := u.seoRefRepo.Exist(ctx, seoDomainID, uint32(refType), refID, &id)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, _errors.ReturnError(409, "seo ref đã tồn tại")
	}

	if err := u.seoRefRepo.Update(ctx, id, map[string]interface{}{
		"seo_domain_id": seoDomainID,
		"ref_type":      refType,
		"ref_id":        refID,
	}); err != nil {
		return nil, err
	}
	return u.GetSeoRef(ctx, id)
}

func (u *SeoDomainUsecase) DeleteSeoRef(ctx context.Context, id uint64) error {
	if _, err := u.GetSeoRef(ctx, id); err != nil {
		return err
	}
	return u.seoRefRepo.Delete(ctx, id)
}

func (u *SeoDomainUsecase) prepareSeoRef(ctx context.Context, seoDomainID uint64, refType enums.ESEORefType, refID uint64) (uint64, enums.ESEORefType, uint64, error) {
	if seoDomainID == 0 {
		return 0, 0, 0, _errors.ReturnError(400, "seoDomainId không hợp lệ")
	}
	if !refType.Validate() {
		return 0, 0, 0, _errors.ReturnError(400, "refType không hợp lệ")
	}
	if refID == 0 {
		return 0, 0, 0, _errors.ReturnError(400, "refId không hợp lệ")
	}

	seo, err := u.GetSeoDomain(ctx, seoDomainID)
	if err != nil {
		return 0, 0, 0, err
	}
	if seo == nil {
		return 0, 0, 0, _errors.ReturnError(404, "seo domain không tồn tại")
	}

	return seo.ID, refType, refID, nil
}
