package usecase

import (
	_errors "common/errors"
	"context"
	"crm/config"
	seo_domain "crm/internal/domain/seo"
	"crm/internal/dto"
	"crm/internal/enums"
	"crm/internal/interface/provider"
	"crm/internal/repo"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/spf13/viper"
	"gorm.io/datatypes"
)

type SeoDomainUsecase struct {
	seoDomainRepo   repo.SeoDomainRepo
	seoInternalRepo repo.SeoInternalLinkRepo
	seoRelativeRepo repo.SeoRelativeRepo
	tqdProvider     provider.TqdProvider
	transaction     provider.ITransaction
	bdsproProvider  provider.BdsproProvider
}

func NewSeoDomainUsecase(
	seoDomainRepo repo.SeoDomainRepo,
	seoInternalRepo repo.SeoInternalLinkRepo,
	seoRelativeRepo repo.SeoRelativeRepo,
	tqdProvider provider.TqdProvider,
	transaction provider.ITransaction,
	bdsproProvider provider.BdsproProvider,
) *SeoDomainUsecase {
	return &SeoDomainUsecase{
		seoDomainRepo:   seoDomainRepo,
		seoInternalRepo: seoInternalRepo,
		seoRelativeRepo: seoRelativeRepo,
		tqdProvider:     tqdProvider,
		transaction:     transaction,
		bdsproProvider:  bdsproProvider,
	}
}

func (u *SeoDomainUsecase) CreateSeoDomain(ctx context.Context, req *dto.SeoDomainCreateRequest) (*seo_domain.SeoDomain, error) {
	seo, err := u.buildSeoDomainForCreate(ctx, req)
	if err != nil {
		return nil, err
	}
	return u.seoDomainRepo.Create(ctx, seo)
}

func (u *SeoDomainUsecase) GetSeoDomain(ctx context.Context, id uint64) (*seo_domain.SeoDomain, error) {
	seoDomain, err := u.seoDomainRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if seoDomain == nil {
		return nil, _errors.ReturnError(404, "seo domain không tồn tại")
	}
	return seoDomain, nil
}

func (u *SeoDomainUsecase) GetSeoDomainByRef(ctx context.Context, refType uint32, refID uint64) (*seo_domain.SeoDomain, error) {
	if !enums.IsValidSEORefType(refType) {
		return nil, _errors.ReturnError(400, "refType không hợp lệ")
	}
	if refID == 0 {
		return nil, _errors.ReturnError(400, "refId không hợp lệ")
	}

	seoDomain, err := u.seoDomainRepo.GetByRef(ctx, refType, refID)
	if err != nil {
		return nil, err
	}
	if seoDomain == nil {
		return nil, _errors.ReturnError(404, "seo domain không tồn tại")
	}
	return seoDomain, nil
}

func (u *SeoDomainUsecase) GetSeoDomainList(ctx context.Context, req *dto.SeoDomainListRequest) ([]seo_domain.SeoDomain, int64, error) {
	return u.seoDomainRepo.GetList(ctx, req)
}

func (u *SeoDomainUsecase) GetSitemapXML(ctx context.Context, req *dto.SitemapRequest) ([]byte, error) {
	baseURL := strings.TrimSpace(config.AppProperties.Seo.PublicBaseURL)
	if req != nil && req.BaseURL != nil && strings.TrimSpace(*req.BaseURL) != "" {
		baseURL = strings.TrimSpace(*req.BaseURL)
	}

	xmlText, err := u.BuildSitemapXML(ctx, &dto.SeoSitemapRequest{}, baseURL)
	if err != nil {
		return nil, err
	}
	return []byte(xmlText), nil
}

// -----------------------------------------------------------------------------
// Parcel generation legacy flow.
// Kept for backward compatibility. The generic source-aware create path now
// lives in buildSeoDomainForCreate + prepareSourceBinding.
// -----------------------------------------------------------------------------

func (u *SeoDomainUsecase) CreateSeoFromParcel(ctx context.Context, req *dto.CreateSeoFromParcelRequest) (*seo_domain.SeoDomain, error) {
	if req == nil || req.ParcelID == 0 {
		return nil, _errors.ReturnError(400, "parcelId không được để trống")
	}
	if u.tqdProvider == nil {
		return nil, _errors.ReturnError(503, "TQD provider chưa được cấu hình")
	}

	source, err := u.tqdProvider.GetParcelSeoSource(ctx, req.ParcelID)
	if err != nil {
		return nil, err
	}
	if source == nil {
		return nil, _errors.ReturnError(404, "parcel không tồn tại")
	}

	seoDomain, _, err := u.createSeoFromParcelSource(ctx, *source)
	return seoDomain, err
}

func (u *SeoDomainUsecase) GenerateSeoFromParcels(ctx context.Context, req *dto.GenerateSeoFromParcelsRequest) ([]seo_domain.SeoDomain, *dto.GenerateSeoFromParcelsResult, error) {
	if u.tqdProvider == nil {
		return nil, nil, _errors.ReturnError(503, "TQD provider chưa được cấu hình")
	}

	limit := uint32(0)
	if req != nil {
		limit = req.Limit
	}
	slog.InfoContext(ctx, fmt.Sprintf("[SEO_GENERATE] start generate seo from parcels limit=%d", limit))
	sources, err := u.tqdProvider.GetParcelSeoSourcesForGenerate(ctx, limit)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[SEO_GENERATE] get parcel seo sources failed: %v", err))
		return nil, nil, err
	}
	slog.InfoContext(ctx, fmt.Sprintf("[SEO_GENERATE] received parcel seo sources total=%d", len(sources)))

	result := &dto.GenerateSeoFromParcelsResult{
		Total: int64(len(sources)),
	}
	createdItems := make([]seo_domain.SeoDomain, 0, len(sources))

	for i := range sources {
		slog.InfoContext(ctx, fmt.Sprintf("[SEO_GENERATE] processing parcelId=%d adrSearch=%q seoId=%v", sources[i].ParcelID, sources[i].AdrSearch, sources[i].SeoID))

		seoDomain, created, err := u.createSeoFromParcelSource(ctx, sources[i])
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("parcelId=%d: %v", sources[i].ParcelID, err))
			slog.ErrorContext(ctx, fmt.Sprintf("[SEO_GENERATE] failed parcelId=%d error=%v", sources[i].ParcelID, err))
			continue
		}
		if seoDomain == nil || !created {
			result.Skipped++
			if seoDomain == nil {
				slog.WarnContext(ctx, fmt.Sprintf("[SEO_GENERATE] skipped parcelId=%d reason=seoDomain_nil", sources[i].ParcelID))
			} else {
				slog.WarnContext(ctx, fmt.Sprintf("[SEO_GENERATE] skipped parcelId=%d reason=already_exists seoId=%d canonicalUrl=%s", sources[i].ParcelID, seoDomain.ID, seoDomain.CanonicalURL))
			}
			continue
		}

		result.Created++
		createdItems = append(createdItems, *seoDomain)
		slog.InfoContext(ctx, fmt.Sprintf("[SEO_GENERATE] created parcelId=%d seoId=%d canonicalUrl=%s", sources[i].ParcelID, seoDomain.ID, seoDomain.CanonicalURL))
	}
	slog.ErrorContext(ctx, fmt.Sprintf("[SEO_GENERATE] done total=%d created=%d skipped=%d failed=%d", result.Total, result.Created, result.Skipped, result.Failed))
	return createdItems, result, nil
}

func (u *SeoDomainUsecase) createSeoFromParcelSource(ctx context.Context, source dto.ParcelSeoSource) (*seo_domain.SeoDomain, bool, error) {
	if u.tqdProvider == nil {
		return nil, false, _errors.ReturnError(503, "TQD provider chưa được cấu hình")
	}

	if source.SeoID != nil && *source.SeoID > 0 {
		slog.InfoContext(ctx, fmt.Sprintf("[SEO_GENERATE] parcelId=%d has seoId=%d, checking existing seo_domain", source.ParcelID, *source.SeoID))

		existing, err := u.seoDomainRepo.GetByID(ctx, *source.SeoID)
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("[SEO_GENERATE] get existing by seoId failed parcelId=%d seoId=%d error=%v", source.ParcelID, *source.SeoID, err))
			return nil, false, err
		}
		if existing != nil {
			slog.InfoContext(ctx, fmt.Sprintf("[SEO_GENERATE] found existing by seoId parcelId=%d seoId=%d canonicalUrl=%s", source.ParcelID, existing.ID, existing.CanonicalURL))
			return existing, false, nil
		}
		slog.WarnContext(ctx, fmt.Sprintf("[SEO_GENERATE] parcelId=%d seoId=%d not found in crm, will recreate", source.ParcelID, *source.SeoID))
	}

	parcelURL, err := buildParcelSeoURL(source.ParcelID, source.AdrSearch)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[SEO_GENERATE] build parcel seo url failed parcelId=%d adrSearch=%q error=%v", source.ParcelID, source.AdrSearch, err))
		return nil, false, err
	}
	slog.InfoContext(ctx, fmt.Sprintf("[SEO_GENERATE] built parcel url parcelId=%d url=%s", source.ParcelID, parcelURL))

	existingByURL, err := u.seoDomainRepo.GetByCanonicalURL(ctx, parcelURL)
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[SEO_GENERATE] get existing by canonical failed parcelId=%d canonicalUrl=%s error=%v", source.ParcelID, parcelURL, err))
		return nil, false, err
	}
	if existingByURL != nil {
		slog.InfoContext(ctx, fmt.Sprintf("[SEO_GENERATE] existing canonical found parcelId=%d seoId=%d, updating tqd seo_id", source.ParcelID, existingByURL.ID))
		if err := u.tqdProvider.UpdateParcelSeoID(ctx, source.ParcelID, existingByURL.ID); err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("[SEO_GENERATE] update tqd seo_id failed parcelId=%d seoId=%d error=%v", source.ParcelID, existingByURL.ID, err))
			return nil, false, err
		}
		slog.InfoContext(ctx, fmt.Sprintf("[SEO_GENERATE] updated tqd seo_id from existing canonical parcelId=%d seoId=%d", source.ParcelID, existingByURL.ID))
		return existingByURL, false, nil
	}

	metadataBytes, err := json.Marshal(map[string]interface{}{
		"parcelId":  source.ParcelID,
		"adrSearch": source.AdrSearch,
		"source":    seo_domain.SeoRefSourceTqdParcel,
		"sourceKey": "parcel",
	})
	if err != nil {
		return nil, false, err
	}

	now := time.Now()
	refID := source.ParcelID
	seoDomain := &seo_domain.SeoDomain{
		Slug:              resolveParcelSeoSlug(source.AdrSearch, source.ParcelID),
		OriginURL:         parcelURL,
		CanonicalURL:      parcelURL,
		RefType:           enums.ESEORefTypeParcel,
		RefID:             &refID,
		RefSource:         seo_domain.SeoRefSourceTqdParcel,
		RefLabel:          firstNonEmpty(source.AdrSearch, fmt.Sprintf("Thửa đất #%d", source.ParcelID)),
		RefURL:            parcelURL,
		SourceStatus:      seo_domain.SeoSourceStatusLinked,
		RefMissing:        false,
		RefSnapshotJSON:   datatypes.JSON(metadataBytes),
		RefHash:           hashString(string(metadataBytes)),
		RefLastSyncedAt:   &now,
		Scope:             seo_domain.SeoScopePublic,
		Title:             buildParcelSeoTitle(source.AdrSearch, source.ParcelID),
		Description:       buildParcelSeoDescription(source.AdrSearch, source.ParcelID),
		Published:         true,
		PublishedAt:       &now,
		IsSiteMap:         true,
		IsIndex:           true,
		SiteMapLastedAt:   &now,
		SitemapPriority:   0.8,
		SitemapChangeFreq: seo_domain.SeoChangeFreqDaily,
		NeedGenerate:      true,
		SourceUpdatedAt:   &now,
		Metadata:          datatypes.JSON(metadataBytes),
	}

	var created *seo_domain.SeoDomain
	err = u.transaction.WithTransaction(ctx, func(txCtx context.Context) error {
		slog.InfoContext(ctx, fmt.Sprintf("[SEO_GENERATE] creating seo_domain parcelId=%d slug=%q canonicalUrl=%s", source.ParcelID, seoDomain.Slug, seoDomain.CanonicalURL))

		createdSeo, err := u.seoDomainRepo.Create(txCtx, seoDomain)
		if err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("[SEO_GENERATE] create seo_domain failed parcelId=%d canonicalUrl=%s error=%v", source.ParcelID, seoDomain.CanonicalURL, err))
			return err
		}
		if createdSeo == nil || createdSeo.ID == 0 {
			slog.InfoContext(ctx, fmt.Sprintf("[SEO_GENERATE] create seo_domain returned empty parcelId=%d canonicalUrl=%s", source.ParcelID, seoDomain.CanonicalURL))
			return _errors.ReturnError(500, "tạo seo domain thất bại")
		}
		slog.InfoContext(ctx, fmt.Sprintf("[SEO_GENERATE] create seo_domain success parcelId=%d seoId=%d", source.ParcelID, createdSeo.ID))
		if err := u.tqdProvider.UpdateParcelSeoID(txCtx, source.ParcelID, createdSeo.ID); err != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("[SEO_GENERATE] update tqd seo_id after create failed parcelId=%d seoId=%d error=%v", source.ParcelID, createdSeo.ID, err))
			return err
		}
		slog.InfoContext(ctx, fmt.Sprintf("[SEO_GENERATE] update tqd seo_id after create success parcelId=%d seoId=%d", source.ParcelID, createdSeo.ID))

		created = createdSeo
		return nil
	})
	if err != nil {
		slog.ErrorContext(ctx, fmt.Sprintf("[SEO_GENERATE] transaction failed parcelId=%d error=%v", source.ParcelID, err))
		return nil, false, err
	}

	return created, true, nil
}

// -----------------------------------------------------------------------------
// Update / delete.
// -----------------------------------------------------------------------------

func (u *SeoDomainUsecase) UpdateSeoDomain(ctx context.Context, id uint64, req *dto.SeoDomainUpdateRequest) (*seo_domain.SeoDomain, error) {
	existing, err := u.GetSeoDomain(ctx, id)
	if err != nil {
		return nil, err
	}

	updates, err := u.buildSeoDomainUpdates(ctx, existing, req)
	if err != nil {
		return nil, err
	}
	if len(updates) == 0 {
		return existing, nil
	}

	if err := u.seoDomainRepo.Update(ctx, id, updates); err != nil {
		return nil, err
	}
	return u.GetSeoDomain(ctx, id)
}

func (u *SeoDomainUsecase) PublishSeoPage(ctx context.Context, id uint64) (*seo_domain.SeoDomain, error) {
	existing, err := u.GetSeoDomain(ctx, id)
	if err != nil {
		return nil, err
	}

	state := seoState{
		Slug:       existing.Slug,
		Canonical:  existing.CanonicalURL,
		Title:      existing.Title,
		Published:  true,
		IsSiteMap:  existing.IsSiteMap,
		IsIndex:    existing.IsIndex,
		PageStatus: seo_domain.SeoPageStatusPublished,
	}
	if err := validateSeoState(state); err != nil {
		return nil, err
	}

	now := time.Now()
	updates := map[string]interface{}{
		"published":         true,
		"page_status":       seo_domain.SeoPageStatusPublished,
		"published_at":      &now,
		"need_generate":     true,
		"render_status":     seo_domain.SeoRenderStatusPending,
		"last_render_error": "",
	}
	if existing.PublishedAt != nil {
		updates["published_at"] = existing.PublishedAt
	}

	if err := u.seoDomainRepo.Update(ctx, id, updates); err != nil {
		return nil, err
	}
	return u.GetSeoDomain(ctx, id)
}

func (u *SeoDomainUsecase) UnpublishSeoPage(ctx context.Context, id uint64) (*seo_domain.SeoDomain, error) {
	if _, err := u.GetSeoDomain(ctx, id); err != nil {
		return nil, err
	}

	updates := map[string]interface{}{
		"published":     false,
		"page_status":   seo_domain.SeoPageStatusDraft,
		"is_site_map":   false,
		"need_generate": true,
		"render_status": seo_domain.SeoRenderStatusPending,
	}
	if err := u.seoDomainRepo.Update(ctx, id, updates); err != nil {
		return nil, err
	}
	return u.GetSeoDomain(ctx, id)
}

func (u *SeoDomainUsecase) ArchiveSeoPage(ctx context.Context, id uint64) (*seo_domain.SeoDomain, error) {
	if _, err := u.GetSeoDomain(ctx, id); err != nil {
		return nil, err
	}

	updates := map[string]interface{}{
		"published":        false,
		"page_status":      seo_domain.SeoPageStatusArchived,
		"is_site_map":      false,
		"is_index":         false,
		"need_generate":    false,
		"render_status":    seo_domain.SeoRenderStatusNone,
		"static_html_hash": "",
	}
	if err := u.seoDomainRepo.Update(ctx, id, updates); err != nil {
		return nil, err
	}
	return u.GetSeoDomain(ctx, id)
}

func (u *SeoDomainUsecase) RestoreSeoPage(ctx context.Context, id uint64) (*seo_domain.SeoDomain, error) {
	if _, err := u.GetSeoDomain(ctx, id); err != nil {
		return nil, err
	}

	updates := map[string]interface{}{
		"published":     false,
		"page_status":   seo_domain.SeoPageStatusDraft,
		"is_site_map":   false,
		"need_generate": true,
		"render_status": seo_domain.SeoRenderStatusPending,
	}
	if err := u.seoDomainRepo.Update(ctx, id, updates); err != nil {
		return nil, err
	}
	return u.GetSeoDomain(ctx, id)
}

func (u *SeoDomainUsecase) DeleteSeoDomain(ctx context.Context, id uint64) error {
	if _, err := u.GetSeoDomain(ctx, id); err != nil {
		return err
	}

	// Xóa page SEO phải xóa mềm link/relative liên quan để không còn đường ray trỏ tới ga đã xóa.
	return u.transaction.WithTransaction(ctx, func(ctx context.Context) error {
		if err := u.seoInternalRepo.DeleteBySeoDomainID(ctx, id); err != nil {
			return err
		}
		if err := u.seoRelativeRepo.DeleteBySeoDomainID(ctx, id); err != nil {
			return err
		}
		return u.seoDomainRepo.Delete(ctx, id)
	})
}

func (u *SeoDomainUsecase) TouchByRef(ctx context.Context, req *dto.TouchSeoDomainByRefRequest) error {
	if req == nil {
		return _errors.ReturnError(400, "request không hợp lệ")
	}
	if !enums.IsValidSEORefType(req.RefType) || req.RefType == 0 {
		return _errors.ReturnError(400, "refType không hợp lệ")
	}
	if req.RefID == 0 {
		return _errors.ReturnError(400, "refId không hợp lệ")
	}

	cfg, err := resolveSeoRefConfig(req.RefType, req.RefSource, "")
	if err != nil {
		return err
	}
	cfg = seo_domain.NormalizeSeoRefTypeConfig(cfg)
	if !seo_domain.IsSeoRefTypeResolvable(cfg) || !seo_domain.IsSeoResolverReady(cfg.ResolverKey) {
		return resolverNotImplemented(cfg.ResolverKey)
	}

	existing, err := u.seoDomainRepo.GetByRefSource(ctx, req.RefType, cfg.ResolverKey, req.RefID)
	if err != nil {
		return err
	}

	now := time.Now()
	snapshot, err := newSeoRefResolver(cfg, u.bdsproProvider, u.tqdProvider).Resolve(ctx, &dto.ResolveSeoRefValueRequest{
		RefType:     req.RefType,
		RefID:       req.RefID,
		ResolverKey: cfg.ResolverKey,
	})
	if err != nil {
		updates := map[string]interface{}{
			"source_status":     seo_domain.SeoSourceStatusMissing,
			"ref_missing":       true,
			"source_updated_at": now,
			"need_generate":     true,
		}
		if existing != nil {
			return u.seoDomainRepo.UpdateSourceStateByRefSource(ctx, req.RefType, cfg.ResolverKey, req.RefID, updates)
		}
		return err
	}

	needGenerate := existing == nil || existing.RefHash != snapshot.Hash
	status := seo_domain.SeoSourceStatusLinked
	if needGenerate && existing != nil {
		status = seo_domain.SeoSourceStatusStale
	}

	updates := map[string]interface{}{
		"ref_label":          snapshot.Label,
		"ref_url":            firstNonEmpty(snapshot.PublicURL, snapshot.AdminURL),
		"source_status":      status,
		"ref_missing":        false,
		"ref_snapshot_json":  datatypes.JSON([]byte(validJSONOrEmpty(snapshot.DataJSON))),
		"ref_hash":           snapshot.Hash,
		"ref_last_synced_at": now,
		"source_updated_at":  now,
		"need_generate":      needGenerate,
	}

	if existing == nil {
		return u.seoDomainRepo.MarkNeedGenerateByRefSource(ctx, req.RefType, cfg.ResolverKey, req.RefID, now)
	}
	return u.seoDomainRepo.UpdateSourceStateByRefSource(ctx, req.RefType, cfg.ResolverKey, req.RefID, updates)
}

// -----------------------------------------------------------------------------
// Sitemap.
// -----------------------------------------------------------------------------

func (u *SeoDomainUsecase) BuildSitemapXML(ctx context.Context, req *dto.SeoSitemapRequest, baseURL string) (string, error) {
	items, err := u.seoDomainRepo.GetSitemapItems(ctx, req.Scope)
	if err != nil {
		return "", err
	}

	urlset := sitemapURLSet{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
		URLs:  make([]sitemapURL, 0, len(items)),
	}

	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	for _, item := range items {
		loc := buildAbsoluteSitemapLoc(item.CanonicalURL, baseURL)
		lastmod := item.SiteMapLastedAt
		if lastmod == "" {
			lastmod = item.UpdatedAt
		}

		urlset.URLs = append(urlset.URLs, sitemapURL{
			Loc:        loc,
			Lastmod:    lastmod,
			ChangeFreq: item.SitemapChangeFreq,
			Priority:   item.SitemapPriority,
		})
	}

	output, err := xml.MarshalIndent(urlset, "", "  ")
	if err != nil {
		return "", err
	}
	return xml.Header + string(output), nil
}

func buildAbsoluteSitemapLoc(canonicalURL string, baseURL string) string {
	loc := strings.TrimSpace(canonicalURL)
	if loc == "" || strings.HasPrefix(loc, "http://") || strings.HasPrefix(loc, "https://") {
		return loc
	}

	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	if baseURL == "" {
		baseURL = strings.TrimRight(strings.TrimSpace(config.AppProperties.Seo.PublicBaseURL), "/")
	}

	return baseURL + "/" + strings.TrimLeft(loc, "/")
}

func (u *SeoDomainUsecase) GetInternalRenderData(ctx context.Context, req *dto.InternalSeoRenderDataRequest) (*seo_domain.SeoDomain, error) {
	return u.GetSeoDomainByRef(ctx, req.RefType, req.RefID)
}

// -----------------------------------------------------------------------------
// Internal links.
// -----------------------------------------------------------------------------

func (u *SeoDomainUsecase) CreateSeoInternalLink(ctx context.Context, req *dto.SeoInternalLinkCreateRequest) (*seo_domain.SeoInternalLink, error) {
	parent, child, linkType, priority, err := u.prepareInternalLink(ctx, req.ParentSeoID, req.ChildSeoID, req.LinkType, req.Priority)
	if err != nil {
		return nil, err
	}

	linkURL := strings.TrimSpace(req.Link)
	if linkURL == "" {
		linkURL = child.CanonicalURL
	}

	title := strings.TrimSpace(req.Title)
	if title == "" {
		title = child.Title
	}
	if title == "" {
		return nil, _errors.ReturnError(400, "title là bắt buộc")
	}

	exist, err := u.seoInternalRepo.Exist(ctx, parent.ID, child.ID, linkType, nil)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, _errors.ReturnError(409, "internal link đã tồn tại")
	}

	link := &seo_domain.SeoInternalLink{
		Link:        linkURL,
		Title:       title,
		LinkType:    linkType,
		Priority:    priority,
		ParentSeoID: parent.ID,
		ChildSeoID:  child.ID,
	}
	return u.seoInternalRepo.Create(ctx, link)
}

func (u *SeoDomainUsecase) GetSeoInternalLinkList(ctx context.Context, req *dto.SeoInternalLinkListRequest) ([]seo_domain.SeoInternalLink, int64, error) {
	return u.seoInternalRepo.GetList(ctx, req)
}

func (u *SeoDomainUsecase) UpdateSeoInternalLink(ctx context.Context, id uint64, req *dto.SeoInternalLinkUpdateRequest) (*seo_domain.SeoInternalLink, error) {
	existing, err := u.seoInternalRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, _errors.ReturnError(404, "seo internal link không tồn tại")
	}

	parentID := existing.ParentSeoID
	childID := existing.ChildSeoID
	linkType := existing.LinkType
	priority := existing.Priority

	if req.ParentSeoID != nil {
		parentID = *req.ParentSeoID
	}
	if req.ChildSeoID != nil {
		childID = *req.ChildSeoID
	}
	if req.LinkType != nil {
		linkType = strings.TrimSpace(*req.LinkType)
	}
	if req.Priority != nil {
		priority = *req.Priority
	}

	parent, child, linkType, priority, err := u.prepareInternalLink(ctx, parentID, childID, linkType, priority)
	if err != nil {
		return nil, err
	}

	exist, err := u.seoInternalRepo.Exist(ctx, parent.ID, child.ID, linkType, &id)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, _errors.ReturnError(409, "internal link đã tồn tại")
	}

	updates := map[string]interface{}{
		"parent_seo_id": parent.ID,
		"child_seo_id":  child.ID,
		"link_type":     linkType,
		"priority":      priority,
	}

	if req.Link != nil {
		linkURL := strings.TrimSpace(*req.Link)
		if linkURL == "" {
			linkURL = child.CanonicalURL
		}
		updates["link"] = linkURL
	}
	if req.Title != nil {
		title := strings.TrimSpace(*req.Title)
		if title == "" {
			return nil, _errors.ReturnError(400, "title không được để trống")
		}
		updates["title"] = title
	}

	if err := u.seoInternalRepo.Update(ctx, id, updates); err != nil {
		return nil, err
	}
	return u.seoInternalRepo.GetByID(ctx, id)
}

func (u *SeoDomainUsecase) DeleteSeoInternalLink(ctx context.Context, id uint64) error {
	existing, err := u.seoInternalRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return _errors.ReturnError(404, "seo internal link không tồn tại")
	}
	return u.seoInternalRepo.Delete(ctx, id)
}

// -----------------------------------------------------------------------------
// Relatives.
// -----------------------------------------------------------------------------

func (u *SeoDomainUsecase) CreateSeoRelative(ctx context.Context, req *dto.SeoRelativeCreateRequest) (*seo_domain.SeoRelative, error) {
	parent, child, relationType, priority, err := u.prepareRelative(ctx, req.ParentSeoID, req.ChildSeoID, req.RelationType, req.Priority)
	if err != nil {
		return nil, err
	}

	exist, err := u.seoRelativeRepo.Exist(ctx, parent.ID, child.ID, relationType)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, _errors.ReturnError(409, "seo relative đã tồn tại")
	}

	relative := &seo_domain.SeoRelative{
		ParentSeoID:  parent.ID,
		ChildSeoID:   child.ID,
		RelationType: relationType,
		Priority:     priority,
	}
	return u.seoRelativeRepo.Create(ctx, relative)
}

func (u *SeoDomainUsecase) GetSeoRelativeList(ctx context.Context, req *dto.SeoRelativeListRequest) ([]seo_domain.SeoRelative, int64, error) {
	return u.seoRelativeRepo.GetList(ctx, req)
}

func (u *SeoDomainUsecase) DeleteSeoRelative(ctx context.Context, id uint64) error {
	existing, err := u.seoRelativeRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if existing == nil {
		return _errors.ReturnError(404, "seo relative không tồn tại")
	}
	return u.seoRelativeRepo.Delete(ctx, id)
}

// -----------------------------------------------------------------------------
// Builders / validators.
// -----------------------------------------------------------------------------

func (u *SeoDomainUsecase) buildSeoDomainForCreate(ctx context.Context, req *dto.SeoDomainCreateRequest) (*seo_domain.SeoDomain, error) {
	if req == nil {
		return nil, _errors.ReturnError(400, "request không hợp lệ")
	}
	if err := validateRef(req.RefType, req.RefID); err != nil {
		return nil, err
	}

	sourceBinding, err := u.prepareSourceBinding(ctx, req.RefType, req.RefID, strings.TrimSpace(req.RefSource))
	if err != nil {
		return nil, err
	}

	originURL := strings.TrimSpace(req.OriginURL)
	if originURL == "" && sourceBinding.RefURL != "" {
		originURL = sourceBinding.RefURL
	}
	if originURL == "" {
		originURL = "/seo/" + firstNonEmpty(strings.TrimSpace(req.Slug), sourceBinding.SuggestedSlug, slugify(req.Title), "seo-page")
	}

	canonicalURL := originURL
	if req.CanonicalURL != nil && strings.TrimSpace(*req.CanonicalURL) != "" {
		canonicalURL = strings.TrimSpace(*req.CanonicalURL)
	}

	slug := strings.TrimSpace(req.Slug)
	if slug == "" {
		slug = firstNonEmpty(sourceBinding.SuggestedSlug, slugify(canonicalURL))
	}

	scope := normalizeScope(req.Scope)
	publishedAt, err := parseSeoTime(req.PublishedAt)
	if err != nil {
		return nil, _errors.ReturnError(400, "publishedAt phải theo RFC3339 hoặc YYYY-MM-DD")
	}
	sitemapLastedAt, err := parseSeoTime(req.SiteMapLastedAt)
	if err != nil {
		return nil, _errors.ReturnError(400, "siteMapLastedAt phải theo RFC3339 hoặc YYYY-MM-DD")
	}
	generatedAt, err := parseSeoTime(req.GeneratedAt)
	if err != nil {
		return nil, _errors.ReturnError(400, "generatedAt phải theo RFC3339 hoặc YYYY-MM-DD")
	}
	sourceUpdatedAt, err := parseSeoTime(req.SourceUpdatedAt)
	if err != nil {
		return nil, _errors.ReturnError(400, "sourceUpdatedAt phải theo RFC3339 hoặc YYYY-MM-DD")
	}
	refLastSyncedAt, err := parseSeoTime(req.RefLastSyncedAt)
	if err != nil {
		return nil, _errors.ReturnError(400, "refLastSyncedAt phải theo RFC3339 hoặc YYYY-MM-DD")
	}

	metadata, err := parseMetadata(req.Metadata)
	if err != nil {
		return nil, err
	}

	refSnapshotJSON := sourceBinding.RefSnapshotJSON
	if req.RefSnapshotJSON != nil {
		refSnapshotJSON, err = parseJSONPayload(req.RefSnapshotJSON, "refSnapshotJson")
		if err != nil {
			return nil, err
		}
	}
	if refLastSyncedAt == nil {
		refLastSyncedAt = sourceBinding.RefLastSyncedAt
	}

	needGenerate := true
	if req.NeedGenerate != nil {
		needGenerate = *req.NeedGenerate
	}

	published := req.Published
	pageStatus, err := normalizePageStatusForCreate(req.PageStatus, published)
	if err != nil {
		return nil, err
	}
	if pageStatus == seo_domain.SeoPageStatusPublished {
		published = true
	}
	if pageStatus == seo_domain.SeoPageStatusArchived {
		published = false
	}
	if published && publishedAt == nil {
		now := time.Now()
		publishedAt = &now
	}

	staticHash := strings.TrimSpace(req.StaticHtmlHash)
	renderStatus, err := normalizeRenderStatusForCreate(req.RenderStatus, needGenerate, generatedAt != nil || staticHash != "")
	if err != nil {
		return nil, err
	}

	sitemapPriority := req.SitemapPriority
	if sitemapPriority == 0 {
		sitemapPriority = 0.5
	}
	if sitemapPriority < 0 || sitemapPriority > 1 {
		return nil, _errors.ReturnError(400, "sitemapPriority phải nằm trong khoảng 0..1")
	}

	changeFreq := normalizeChangeFreq(req.SitemapChangeFreq)

	title := firstNonEmpty(strings.TrimSpace(req.Title), sourceBinding.SuggestedTitle)
	description := firstNonEmpty(strings.TrimSpace(req.Description), sourceBinding.SuggestedDescription)

	state := seoState{
		Slug:       slug,
		Canonical:  canonicalURL,
		Title:      title,
		Published:  published,
		IsSiteMap:  req.IsSiteMap,
		IsIndex:    req.IsIndex,
		PageStatus: pageStatus,
	}
	if err := validateSeoState(state); err != nil {
		return nil, err
	}

	exist, err := u.seoDomainRepo.ExistByCanonicalURL(ctx, canonicalURL, nil)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, _errors.ReturnError(409, "canonicalUrl đã tồn tại")
	}

	if sourceBinding.RefSource != seo_domain.SeoRefSourceManual && req.RefID != nil {
		existingByRef, err := u.seoDomainRepo.GetByRefSource(ctx, req.RefType, sourceBinding.RefSource, *req.RefID)
		if err != nil {
			return nil, err
		}
		if existingByRef != nil {
			return nil, _errors.ReturnError(409, "SEO page cho source này đã tồn tại")
		}
	}

	refMissing := sourceBinding.RefMissing
	if req.RefMissing != nil {
		refMissing = *req.RefMissing
	}

	return &seo_domain.SeoDomain{
		Slug:              slug,
		OriginURL:         originURL,
		CanonicalURL:      canonicalURL,
		RefType:           enums.ESEORefType(req.RefType),
		RefID:             req.RefID,
		RefSource:         sourceBinding.RefSource,
		RefLabel:          firstNonEmpty(strings.TrimSpace(req.RefLabel), sourceBinding.RefLabel),
		RefURL:            firstNonEmpty(strings.TrimSpace(req.RefURL), sourceBinding.RefURL),
		SourceStatus:      normalizeSourceStatus(firstNonEmpty(strings.TrimSpace(req.SourceStatus), sourceBinding.SourceStatus), req.RefType == 0),
		RefMissing:        refMissing,
		RefSnapshotJSON:   refSnapshotJSON,
		RefHash:           firstNonEmpty(strings.TrimSpace(req.RefHash), sourceBinding.RefHash),
		RefLastSyncedAt:   refLastSyncedAt,
		Scope:             scope,
		PageStatus:        pageStatus,
		Title:             title,
		Description:       description,
		Content:           req.Content,
		Summary:           strings.TrimSpace(req.Summary),
		Published:         published,
		PublishedAt:       publishedAt,
		IsSiteMap:         req.IsSiteMap,
		IsIndex:           req.IsIndex,
		IsRobot:           req.IsRobot,
		SiteMapLastedAt:   sitemapLastedAt,
		SitemapPriority:   sitemapPriority,
		SitemapChangeFreq: changeFreq,
		NeedGenerate:      needGenerate,
		GeneratedAt:       generatedAt,
		SourceUpdatedAt:   sourceUpdatedAt,
		RenderStatus:      renderStatus,
		RenderedHTML:      req.RenderedHTML,
		LastRenderError:   strings.TrimSpace(req.LastRenderError),
		StaticHtmlPath:    strings.TrimSpace(req.StaticHtmlPath),
		StaticHtmlHash:    staticHash,
		TemplateKey:       strings.TrimSpace(req.TemplateKey),
		TemplateVersion:   strings.TrimSpace(req.TemplateVersion),
		DeepLink:          strings.TrimSpace(req.DeepLink),
		Metadata:          metadata,
		Note:              req.Note,
	}, nil
}

func (u *SeoDomainUsecase) buildSeoDomainUpdates(ctx context.Context, existing *seo_domain.SeoDomain, req *dto.SeoDomainUpdateRequest) (map[string]interface{}, error) {
	if req == nil {
		return map[string]interface{}{}, nil
	}

	updates := map[string]interface{}{}

	pageStatus := existing.PageStatus
	if pageStatus == "" {
		pageStatus = derivePageStatusFromLegacy(existing.Published, existing.DeletedAt.Valid)
	}

	state := seoState{
		Slug:       existing.Slug,
		Canonical:  existing.CanonicalURL,
		Title:      existing.Title,
		Published:  existing.Published,
		IsSiteMap:  existing.IsSiteMap,
		IsIndex:    existing.IsIndex,
		PageStatus: pageStatus,
	}

	refType := uint32(existing.RefType)
	refID := existing.RefID
	refSource := existing.RefSource

	if req.Slug != nil {
		slug := strings.TrimSpace(*req.Slug)
		if slug == "" {
			return nil, _errors.ReturnError(400, "slug không được để trống")
		}
		state.Slug = slug
		updates["slug"] = slug
	}

	if req.OriginURL != nil {
		originURL := strings.TrimSpace(*req.OriginURL)
		if originURL == "" {
			return nil, _errors.ReturnError(400, "originUrl không được để trống")
		}
		updates["origin_url"] = originURL
	}

	if req.CanonicalURL != nil {
		canonicalURL := strings.TrimSpace(*req.CanonicalURL)
		if canonicalURL == "" {
			canonicalURL = existing.OriginURL
		}

		exist, err := u.seoDomainRepo.ExistByCanonicalURL(ctx, canonicalURL, &existing.ID)
		if err != nil {
			return nil, err
		}
		if exist {
			return nil, _errors.ReturnError(409, "canonicalUrl đã tồn tại")
		}

		state.Canonical = canonicalURL
		updates["canonical_url"] = canonicalURL
	}

	if req.RefType != nil {
		refType = *req.RefType
		updates["ref_type"] = enums.ESEORefType(*req.RefType)
	}
	if req.RefID != nil {
		refID = req.RefID
		updates["ref_id"] = req.RefID
	}
	if req.RefSource != nil {
		refSource = strings.TrimSpace(*req.RefSource)
	}

	if err := validateRef(refType, refID); err != nil {
		return nil, err
	}

	if req.RefType != nil || req.RefID != nil || req.RefSource != nil {
		sourceBinding, err := u.prepareSourceBinding(ctx, refType, refID, refSource)
		if err != nil {
			return nil, err
		}

		if sourceBinding.RefSource != seo_domain.SeoRefSourceManual && refID != nil {
			existingByRef, err := u.seoDomainRepo.GetByRefSource(ctx, refType, sourceBinding.RefSource, *refID)
			if err != nil {
				return nil, err
			}
			if existingByRef != nil && existingByRef.ID != existing.ID {
				return nil, _errors.ReturnError(409, "SEO page cho source này đã tồn tại")
			}
		}

		updates["ref_source"] = sourceBinding.RefSource
		updates["ref_label"] = sourceBinding.RefLabel
		updates["ref_url"] = sourceBinding.RefURL
		updates["source_status"] = sourceBinding.SourceStatus
		updates["ref_missing"] = sourceBinding.RefMissing
		updates["ref_snapshot_json"] = sourceBinding.RefSnapshotJSON
		updates["ref_hash"] = sourceBinding.RefHash
		updates["ref_last_synced_at"] = sourceBinding.RefLastSyncedAt
	}

	if req.Scope != nil {
		updates["scope"] = normalizeScope(*req.Scope)
	}

	if req.Title != nil {
		state.Title = strings.TrimSpace(*req.Title)
		updates["title"] = state.Title
	}
	if req.Description != nil {
		updates["description"] = strings.TrimSpace(*req.Description)
	}
	if req.Content != nil {
		updates["content"] = *req.Content
	}
	if req.Summary != nil {
		updates["summary"] = strings.TrimSpace(*req.Summary)
	}

	if req.Published != nil {
		state.Published = *req.Published
		updates["published"] = *req.Published
		if *req.Published {
			state.PageStatus = seo_domain.SeoPageStatusPublished
			updates["page_status"] = seo_domain.SeoPageStatusPublished
		} else {
			state.PageStatus = seo_domain.SeoPageStatusDraft
			updates["page_status"] = seo_domain.SeoPageStatusDraft
		}

		if *req.Published && existing.PublishedAt == nil && req.PublishedAt == nil {
			now := time.Now()
			updates["published_at"] = &now
		}
	}

	if req.PublishedAt != nil {
		parsed, err := parseSeoTime(req.PublishedAt)
		if err != nil {
			return nil, _errors.ReturnError(400, "publishedAt phải theo RFC3339 hoặc YYYY-MM-DD")
		}
		updates["published_at"] = parsed
	}

	if req.IsSiteMap != nil {
		state.IsSiteMap = *req.IsSiteMap
		updates["is_site_map"] = *req.IsSiteMap
	}
	if req.IsIndex != nil {
		state.IsIndex = *req.IsIndex
		updates["is_index"] = *req.IsIndex
	}
	if req.IsRobot != nil {
		updates["is_robot"] = *req.IsRobot
	}

	if req.SiteMapLastedAt != nil {
		parsed, err := parseSeoTime(req.SiteMapLastedAt)
		if err != nil {
			return nil, _errors.ReturnError(400, "siteMapLastedAt phải theo RFC3339 hoặc YYYY-MM-DD")
		}
		updates["site_map_lasted_at"] = parsed
	}

	if req.SitemapPriority != nil {
		if *req.SitemapPriority < 0 || *req.SitemapPriority > 1 {
			return nil, _errors.ReturnError(400, "sitemapPriority phải nằm trong khoảng 0..1")
		}
		updates["sitemap_priority"] = *req.SitemapPriority
	}

	if req.SitemapChangeFreq != nil {
		updates["sitemap_change_freq"] = normalizeChangeFreq(*req.SitemapChangeFreq)
	}
	if req.NeedGenerate != nil {
		updates["need_generate"] = *req.NeedGenerate
		if *req.NeedGenerate {
			updates["render_status"] = seo_domain.SeoRenderStatusPending
		}
	}
	if req.GeneratedAt != nil {
		parsed, err := parseSeoTime(req.GeneratedAt)
		if err != nil {
			return nil, _errors.ReturnError(400, "generatedAt phải theo RFC3339 hoặc YYYY-MM-DD")
		}
		updates["generated_at"] = parsed
	}
	if req.SourceUpdatedAt != nil {
		parsed, err := parseSeoTime(req.SourceUpdatedAt)
		if err != nil {
			return nil, _errors.ReturnError(400, "sourceUpdatedAt phải theo RFC3339 hoặc YYYY-MM-DD")
		}
		updates["source_updated_at"] = parsed
	}

	if req.RefLabel != nil {
		updates["ref_label"] = strings.TrimSpace(*req.RefLabel)
	}
	if req.RefURL != nil {
		updates["ref_url"] = strings.TrimSpace(*req.RefURL)
	}
	if req.SourceStatus != nil {
		status := normalizeSourceStatus(*req.SourceStatus, refType == 0)
		updates["source_status"] = status
	}
	if req.RefMissing != nil {
		updates["ref_missing"] = *req.RefMissing
	}
	if req.RefSnapshotJSON != nil {
		refSnapshotJSON, err := parseJSONPayload(req.RefSnapshotJSON, "refSnapshotJson")
		if err != nil {
			return nil, err
		}
		updates["ref_snapshot_json"] = refSnapshotJSON
	}
	if req.RefHash != nil {
		updates["ref_hash"] = strings.TrimSpace(*req.RefHash)
	}
	if req.RefLastSyncedAt != nil {
		parsed, err := parseSeoTime(req.RefLastSyncedAt)
		if err != nil {
			return nil, _errors.ReturnError(400, "refLastSyncedAt phải theo RFC3339 hoặc YYYY-MM-DD")
		}
		updates["ref_last_synced_at"] = parsed
	}
	if req.PageStatus != nil {
		status, err := normalizePageStatusStrict(*req.PageStatus)
		if err != nil {
			return nil, err
		}
		state.PageStatus = status
		updates["page_status"] = status

		switch status {
		case seo_domain.SeoPageStatusPublished:
			state.Published = true
			updates["published"] = true
			if existing.PublishedAt == nil && req.PublishedAt == nil {
				now := time.Now()
				updates["published_at"] = &now
			}
		case seo_domain.SeoPageStatusDraft:
			state.Published = false
			updates["published"] = false
		case seo_domain.SeoPageStatusArchived:
			state.Published = false
			state.IsSiteMap = false
			state.IsIndex = false
			updates["published"] = false
			updates["is_site_map"] = false
			updates["is_index"] = false
			updates["need_generate"] = false
			updates["render_status"] = seo_domain.SeoRenderStatusNone
		}
	}

	if req.RenderStatus != nil {
		status, err := normalizeRenderStatusStrict(*req.RenderStatus)
		if err != nil {
			return nil, err
		}
		updates["render_status"] = status
	}
	if req.RenderedHTML != nil {
		updates["rendered_html"] = *req.RenderedHTML
	}
	if req.LastRenderError != nil {
		updates["last_render_error"] = strings.TrimSpace(*req.LastRenderError)
	}
	if req.TemplateKey != nil {
		updates["template_key"] = strings.TrimSpace(*req.TemplateKey)
	}
	if req.TemplateVersion != nil {
		updates["template_version"] = strings.TrimSpace(*req.TemplateVersion)
	}

	if req.StaticHtmlPath != nil {
		updates["static_html_path"] = strings.TrimSpace(*req.StaticHtmlPath)
	}
	if req.StaticHtmlHash != nil {
		updates["static_html_hash"] = strings.TrimSpace(*req.StaticHtmlHash)
	}
	if req.DeepLink != nil {
		updates["deep_link"] = strings.TrimSpace(*req.DeepLink)
	}
	if req.Metadata != nil {
		metadata, err := parseMetadata(req.Metadata)
		if err != nil {
			return nil, err
		}
		updates["metadata"] = metadata
	}
	if req.Note != nil {
		updates["note"] = *req.Note
	}

	if err := validateSeoState(state); err != nil {
		return nil, err
	}

	return updates, nil
}

func (u *SeoDomainUsecase) prepareInternalLink(ctx context.Context, parentID, childID uint64, linkType string, priority uint32) (*seo_domain.SeoDomain, *seo_domain.SeoDomain, string, uint32, error) {
	if parentID == 0 || childID == 0 {
		return nil, nil, "", 0, _errors.ReturnError(400, "parentSeoId và childSeoId là bắt buộc")
	}
	if parentID == childID {
		return nil, nil, "", 0, _errors.ReturnError(400, "parentSeoId không được trùng childSeoId")
	}

	parent, err := u.GetSeoDomain(ctx, parentID)
	if err != nil {
		return nil, nil, "", 0, err
	}
	child, err := u.GetSeoDomain(ctx, childID)
	if err != nil {
		return nil, nil, "", 0, err
	}

	linkType = strings.TrimSpace(linkType)
	if linkType == "" {
		linkType = seo_domain.SeoLinkTypeNavigation
	}
	if priority == 0 {
		priority = 5
	}

	return parent, child, linkType, priority, nil
}

func (u *SeoDomainUsecase) prepareRelative(ctx context.Context, parentID, childID uint64, relationType string, priority uint32) (*seo_domain.SeoDomain, *seo_domain.SeoDomain, string, uint32, error) {
	if parentID == 0 || childID == 0 {
		return nil, nil, "", 0, _errors.ReturnError(400, "parentSeoId và childSeoId là bắt buộc")
	}
	if parentID == childID {
		return nil, nil, "", 0, _errors.ReturnError(400, "parentSeoId không được trùng childSeoId")
	}

	parent, err := u.GetSeoDomain(ctx, parentID)
	if err != nil {
		return nil, nil, "", 0, err
	}
	child, err := u.GetSeoDomain(ctx, childID)
	if err != nil {
		return nil, nil, "", 0, err
	}

	relationType = strings.TrimSpace(relationType)
	if relationType == "" {
		relationType = seo_domain.SeoRelationTypeParentChild
	}
	if priority == 0 {
		priority = 5
	}

	return parent, child, relationType, priority, nil
}

type seoSourceBinding struct {
	RefSource            string
	RefLabel             string
	RefURL               string
	SourceStatus         string
	RefMissing           bool
	RefSnapshotJSON      datatypes.JSON
	RefHash              string
	RefLastSyncedAt      *time.Time
	SuggestedSlug        string
	SuggestedTitle       string
	SuggestedDescription string
}

func (u *SeoDomainUsecase) prepareSourceBinding(ctx context.Context, refType uint32, refID *uint64, refSource string) (seoSourceBinding, error) {
	emptyJSON := datatypes.JSON([]byte("{}"))

	if refType == 0 {
		return seoSourceBinding{
			RefSource:       seo_domain.SeoRefSourceManual,
			SourceStatus:    seo_domain.SeoSourceStatusManual,
			RefMissing:      false,
			RefSnapshotJSON: emptyJSON,
		}, nil
	}

	if refID == nil || *refID == 0 {
		return seoSourceBinding{}, _errors.ReturnError(400, "refId là bắt buộc khi refType khác 0")
	}

	cfg, err := resolveSeoRefConfig(refType, refSource, "")
	if err != nil {
		return seoSourceBinding{}, err
	}
	cfg = seo_domain.NormalizeSeoRefTypeConfig(cfg)

	if !seo_domain.IsSeoRefTypeResolvable(cfg) {
		return seoSourceBinding{}, _errors.ReturnError(501, "seo source chưa hỗ trợ resolve: "+cfg.Key)
	}
	if !seo_domain.IsSeoResolverReady(cfg.ResolverKey) {
		return seoSourceBinding{}, resolverNotImplemented(cfg.ResolverKey)
	}

	snapshot, err := newSeoRefResolver(cfg, u.bdsproProvider, u.tqdProvider).Resolve(ctx, &dto.ResolveSeoRefValueRequest{
		RefType:       refType,
		RefID:         *refID,
		ResolverKey:   cfg.ResolverKey,
		SourceService: cfg.SourceService,
	})
	if err != nil {
		return seoSourceBinding{}, err
	}

	now := time.Now()
	dataJSON := validJSONOrEmpty(snapshot.DataJSON)

	return seoSourceBinding{
		RefSource:            cfg.ResolverKey,
		RefLabel:             snapshot.Label,
		RefURL:               firstNonEmpty(snapshot.PublicURL, snapshot.AdminURL),
		SourceStatus:         firstNonEmpty(snapshot.SourceStatus, seo_domain.SeoSourceStatusLinked),
		RefMissing:           false,
		RefSnapshotJSON:      datatypes.JSON([]byte(dataJSON)),
		RefHash:              snapshot.Hash,
		RefLastSyncedAt:      &now,
		SuggestedSlug:        snapshot.SuggestedSlug,
		SuggestedTitle:       snapshot.SuggestedTitle,
		SuggestedDescription: snapshot.SuggestedMetaDescription,
	}, nil
}

func validateRef(refType uint32, refID *uint64) error {
	if !enums.IsValidSEORefType(refType) {
		return _errors.ReturnError(400, "refType không hợp lệ")
	}
	if refType == 0 {
		return nil
	}
	if refID == nil || *refID == 0 {
		return _errors.ReturnError(400, "refId là bắt buộc khi refType khác 0")
	}
	return nil
}

type seoState struct {
	Slug       string
	Canonical  string
	Title      string
	Published  bool
	IsSiteMap  bool
	IsIndex    bool
	PageStatus string
}

func validateSeoState(state seoState) error {
	if state.PageStatus != "" {
		if _, err := normalizePageStatusStrict(state.PageStatus); err != nil {
			return err
		}
		if state.PageStatus == seo_domain.SeoPageStatusPublished && !state.Published {
			return _errors.ReturnError(400, "page_status=published yêu cầu published=true")
		}
		if state.PageStatus == seo_domain.SeoPageStatusArchived && (state.Published || state.IsSiteMap || state.IsIndex) {
			return _errors.ReturnError(400, "page_status=archived không được published/index/sitemap")
		}
	}

	if state.Published {
		if strings.TrimSpace(state.Slug) == "" || strings.TrimSpace(state.Canonical) == "" || strings.TrimSpace(state.Title) == "" {
			return _errors.ReturnError(400, "published page cần slug, canonicalUrl và title")
		}
	}

	if state.IsSiteMap && (!state.Published || !state.IsIndex) {
		return _errors.ReturnError(400, "page trong sitemap phải published=true và isIndex=true")
	}

	return nil
}

func parseMetadata(value *string) (datatypes.JSON, error) {
	if value == nil || strings.TrimSpace(*value) == "" {
		return datatypes.JSON([]byte("{}")), nil
	}

	payload := strings.TrimSpace(*value)
	if !json.Valid([]byte(payload)) {
		return nil, _errors.ReturnError(400, "metadata phải là JSON hợp lệ")
	}

	var obj map[string]interface{}
	if err := json.Unmarshal([]byte(payload), &obj); err != nil {
		return nil, _errors.ReturnError(400, "metadata phải là JSON object")
	}

	return datatypes.JSON([]byte(payload)), nil
}

func parseJSONPayload(value *string, fieldName string) (datatypes.JSON, error) {
	if value == nil || strings.TrimSpace(*value) == "" {
		return datatypes.JSON([]byte("{}")), nil
	}

	payload := strings.TrimSpace(*value)
	if !json.Valid([]byte(payload)) {
		return nil, _errors.ReturnError(400, fieldName+" phải là JSON hợp lệ")
	}

	return datatypes.JSON([]byte(payload)), nil
}

func parseSeoTime(value *string) (*time.Time, error) {
	if value == nil {
		return nil, nil
	}

	v := strings.TrimSpace(*value)
	if v == "" {
		return nil, nil
	}

	layouts := []string{time.RFC3339, "2006-01-02", "2006-01-02 15:04:05"}
	var lastErr error

	for _, layout := range layouts {
		parsed, err := time.Parse(layout, v)
		if err == nil {
			return &parsed, nil
		}
		lastErr = err
	}

	return nil, lastErr
}

func normalizePageStatusForCreate(status string, published bool) (string, error) {
	status = strings.TrimSpace(status)
	if status == "" {
		if published {
			return seo_domain.SeoPageStatusPublished, nil
		}
		return seo_domain.SeoPageStatusDraft, nil
	}
	return normalizePageStatusStrict(status)
}

func normalizePageStatusStrict(status string) (string, error) {
	switch strings.TrimSpace(status) {
	case seo_domain.SeoPageStatusDraft:
		return seo_domain.SeoPageStatusDraft, nil
	case seo_domain.SeoPageStatusPublished:
		return seo_domain.SeoPageStatusPublished, nil
	case seo_domain.SeoPageStatusArchived:
		return seo_domain.SeoPageStatusArchived, nil
	default:
		return "", _errors.ReturnError(400, "pageStatus không hợp lệ")
	}
}

func derivePageStatusFromLegacy(published bool, archived bool) string {
	if archived {
		return seo_domain.SeoPageStatusArchived
	}
	if published {
		return seo_domain.SeoPageStatusPublished
	}
	return seo_domain.SeoPageStatusDraft
}

func normalizeRenderStatusForCreate(status string, needGenerate bool, hasRenderedOutput bool) (string, error) {
	status = strings.TrimSpace(status)
	if status != "" {
		return normalizeRenderStatusStrict(status)
	}
	if needGenerate {
		return seo_domain.SeoRenderStatusPending, nil
	}
	if hasRenderedOutput {
		return seo_domain.SeoRenderStatusSuccess, nil
	}
	return seo_domain.SeoRenderStatusNone, nil
}

func normalizeRenderStatusStrict(status string) (string, error) {
	switch strings.TrimSpace(status) {
	case seo_domain.SeoRenderStatusNone:
		return seo_domain.SeoRenderStatusNone, nil
	case seo_domain.SeoRenderStatusPending:
		return seo_domain.SeoRenderStatusPending, nil
	case seo_domain.SeoRenderStatusRendering:
		return seo_domain.SeoRenderStatusRendering, nil
	case seo_domain.SeoRenderStatusSuccess:
		return seo_domain.SeoRenderStatusSuccess, nil
	case seo_domain.SeoRenderStatusFailed:
		return seo_domain.SeoRenderStatusFailed, nil
	default:
		return "", _errors.ReturnError(400, "renderStatus không hợp lệ")
	}
}

func normalizeSourceStatus(status string, isManual bool) string {
	if isManual {
		return seo_domain.SeoSourceStatusManual
	}

	switch strings.TrimSpace(status) {
	case seo_domain.SeoSourceStatusLinked,
		seo_domain.SeoSourceStatusStale,
		seo_domain.SeoSourceStatusMissing,
		seo_domain.SeoSourceStatusPermissionDenied:
		return strings.TrimSpace(status)
	default:
		return seo_domain.SeoSourceStatusLinked
	}
}

func normalizeScope(scope string) string {
	switch strings.TrimSpace(scope) {
	case seo_domain.SeoScopePublic, seo_domain.SeoScopeInternal, seo_domain.SeoScopeApp, seo_domain.SeoScopeAdmin:
		return strings.TrimSpace(scope)
	default:
		return seo_domain.SeoScopePublic
	}
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func normalizeChangeFreq(freq string) string {
	switch strings.TrimSpace(freq) {
	case seo_domain.SeoChangeFreqAlways,
		seo_domain.SeoChangeFreqHourly,
		seo_domain.SeoChangeFreqDaily,
		seo_domain.SeoChangeFreqWeekly,
		seo_domain.SeoChangeFreqMonthly,
		seo_domain.SeoChangeFreqYearly,
		seo_domain.SeoChangeFreqNever:
		return strings.TrimSpace(freq)
	default:
		return seo_domain.SeoChangeFreqDaily
	}
}

func resolveSeoRefConfig(refType uint32, resolverKey string, sourceService string) (seo_domain.SeoRefTypeConfig, error) {
	if !enums.IsValidSEORefType(refType) {
		return seo_domain.SeoRefTypeConfig{}, _errors.ReturnError(400, "refType không hợp lệ")
	}

	resolverKey = strings.TrimSpace(resolverKey)
	sourceService = strings.TrimSpace(sourceService)

	cfg, ok := seo_domain.GetSeoRefTypeConfig(refType, resolverKey, sourceService)
	if !ok {
		return seo_domain.SeoRefTypeConfig{}, _errors.ReturnError(400, "refType hoặc resolver không được hỗ trợ")
	}

	return seo_domain.NormalizeSeoRefTypeConfig(cfg), nil
}

func validJSONOrEmpty(value string) string {
	value = strings.TrimSpace(value)
	if value == "" || !json.Valid([]byte(value)) {
		return "{}"
	}
	return value
}

func hashString(value string) string {
	sum := sha256Sum(value)
	return sum
}

func sha256Sum(value string) string {
	h := sha256.New()
	_, _ = h.Write([]byte(value))
	return hex.EncodeToString(h.Sum(nil))
}

var nonSlugChars = regexp.MustCompile(`[^a-z0-9/_-]+`)

func slugify(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	value = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '/' || r == '_' || r == '-' {
			return r
		}
		return '-'
	}, value)

	value = nonSlugChars.ReplaceAllString(value, "-")
	value = strings.Trim(value, "-/ ")
	if value == "" {
		return "seo-page"
	}
	return value
}

func buildParcelSeoURL(parcelID uint64, adrSearch string) (string, error) {
	template := strings.TrimSpace(viper.GetString("seo.parcelUrlTemplate"))
	if template == "" {
		return "", _errors.ReturnError(500, "seo.parcelUrlTemplate chưa được cấu hình")
	}

	slug := slugify(adrSearch)
	if slug == "" || slug == "seo-page" {
		slug = fmt.Sprintf("%d", parcelID)
	}

	return strings.NewReplacer(
		"{id}", fmt.Sprintf("%d", parcelID),
		"{slug}", slug,
	).Replace(template), nil
}

func resolveParcelSeoSlug(adrSearch string, parcelID uint64) string {
	slug := strings.TrimSpace(adrSearch)
	if slug == "" {
		return fmt.Sprintf("%d", parcelID)
	}
	return slug
}

func buildParcelSeoTitle(adrSearch string, parcelID uint64) string {
	title := strings.TrimSpace(adrSearch)
	if title == "" {
		return fmt.Sprintf("Thửa đất %d", parcelID)
	}
	return title
}

func buildParcelSeoDescription(adrSearch string, parcelID uint64) string {
	title := buildParcelSeoTitle(adrSearch, parcelID)
	return fmt.Sprintf("Thông tin quy hoạch thửa đất %s", title)
}

type sitemapURLSet struct {
	XMLName xml.Name     `xml:"urlset"`
	Xmlns   string       `xml:"xmlns,attr"`
	URLs    []sitemapURL `xml:"url"`
}

type sitemapURL struct {
	Loc        string  `xml:"loc"`
	Lastmod    string  `xml:"lastmod,omitempty"`
	ChangeFreq string  `xml:"changefreq,omitempty"`
	Priority   float32 `xml:"priority,omitempty"`
}
