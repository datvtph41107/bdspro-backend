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
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
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
}

func NewSeoDomainUsecase(
	seoDomainRepo repo.SeoDomainRepo,
	seoInternalRepo repo.SeoInternalLinkRepo,
	seoRelativeRepo repo.SeoRelativeRepo,
	tqdProvider provider.TqdProvider,
	transaction provider.ITransaction,
) *SeoDomainUsecase {
	return &SeoDomainUsecase{
		seoDomainRepo:   seoDomainRepo,
		seoInternalRepo: seoInternalRepo,
		seoRelativeRepo: seoRelativeRepo,
		tqdProvider:     tqdProvider,
		transaction:     transaction,
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

func (u *SeoDomainUsecase) CreateSeoFromParcel(ctx context.Context, req *dto.CreateSeoFromParcelRequest) (*seo_domain.SeoDomain, error) {
	if req == nil || req.ParcelID == 0 {
		return nil, _errors.ReturnError(400, "parcelId không được để trống")
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
	limit := uint32(0)
	if req != nil {
		limit = req.Limit
	}

	log.Printf("[SEO_GENERATE] start generate seo from parcels limit=%d", limit)
	sources, err := u.tqdProvider.GetParcelSeoSourcesForGenerate(ctx, limit)
	if err != nil {
		log.Printf("[SEO_GENERATE] get parcel seo sources failed: %v", err)
		return nil, nil, err
	}
	log.Printf("[SEO_GENERATE] received parcel seo sources total=%d", len(sources))

	result := &dto.GenerateSeoFromParcelsResult{
		Total: int64(len(sources)),
	}
	createdItems := make([]seo_domain.SeoDomain, 0, len(sources))

	for i := range sources {
		log.Printf("[SEO_GENERATE] processing parcelId=%d adrSearch=%q seoId=%v", sources[i].ParcelID, sources[i].AdrSearch, sources[i].SeoID)
		seoDomain, created, err := u.createSeoFromParcelSource(ctx, sources[i])
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("parcelId=%d: %v", sources[i].ParcelID, err))
			log.Printf("[SEO_GENERATE] failed parcelId=%d error=%v", sources[i].ParcelID, err)
			continue
		}
		if seoDomain == nil || !created {
			result.Skipped++
			if seoDomain == nil {
				log.Printf("[SEO_GENERATE] skipped parcelId=%d reason=seoDomain_nil", sources[i].ParcelID)
			} else {
				log.Printf("[SEO_GENERATE] skipped parcelId=%d reason=already_exists seoId=%d canonicalUrl=%s", sources[i].ParcelID, seoDomain.ID, seoDomain.CanonicalURL)
			}
			continue
		}
		result.Created++
		createdItems = append(createdItems, *seoDomain)
		log.Printf("[SEO_GENERATE] created parcelId=%d seoId=%d canonicalUrl=%s", sources[i].ParcelID, seoDomain.ID, seoDomain.CanonicalURL)
	}

	log.Printf("[SEO_GENERATE] done total=%d created=%d skipped=%d failed=%d", result.Total, result.Created, result.Skipped, result.Failed)
	return createdItems, result, nil
}

func (u *SeoDomainUsecase) createSeoFromParcelSource(ctx context.Context, source dto.ParcelSeoSource) (*seo_domain.SeoDomain, bool, error) {
	if source.SeoID != nil && *source.SeoID > 0 {
		log.Printf("[SEO_GENERATE] parcelId=%d has seoId=%d, checking existing seo_domain", source.ParcelID, *source.SeoID)
		existing, err := u.seoDomainRepo.GetByID(ctx, *source.SeoID)
		if err != nil {
			log.Printf("[SEO_GENERATE] get existing by seoId failed parcelId=%d seoId=%d error=%v", source.ParcelID, *source.SeoID, err)
			return nil, false, err
		}
		if existing != nil {
			log.Printf("[SEO_GENERATE] found existing by seoId parcelId=%d seoId=%d canonicalUrl=%s", source.ParcelID, existing.ID, existing.CanonicalURL)
			return existing, false, nil
		}
		log.Printf("[SEO_GENERATE] parcelId=%d seoId=%d not found in crm, will recreate", source.ParcelID, *source.SeoID)
	}

	parcelURL, err := buildParcelSeoURL(source.ParcelID, source.AdrSearch)
	if err != nil {
		log.Printf("[SEO_GENERATE] build parcel seo url failed parcelId=%d adrSearch=%q error=%v", source.ParcelID, source.AdrSearch, err)
		return nil, false, err
	}
	log.Printf("[SEO_GENERATE] built parcel url parcelId=%d url=%s", source.ParcelID, parcelURL)

	existingByURL, err := u.seoDomainRepo.GetByCanonicalURL(ctx, parcelURL)
	if err != nil {
		log.Printf("[SEO_GENERATE] get existing by canonical failed parcelId=%d canonicalUrl=%s error=%v", source.ParcelID, parcelURL, err)
		return nil, false, err
	}
	if existingByURL != nil {
		log.Printf("[SEO_GENERATE] existing canonical found parcelId=%d seoId=%d, updating tqd seo_id", source.ParcelID, existingByURL.ID)
		if err := u.tqdProvider.UpdateParcelSeoID(ctx, source.ParcelID, existingByURL.ID); err != nil {
			log.Printf("[SEO_GENERATE] update tqd seo_id failed parcelId=%d seoId=%d error=%v", source.ParcelID, existingByURL.ID, err)
			return nil, false, err
		}
		log.Printf("[SEO_GENERATE] updated tqd seo_id from existing canonical parcelId=%d seoId=%d", source.ParcelID, existingByURL.ID)
		return existingByURL, false, nil
	}

	metadataBytes, err := json.Marshal(map[string]interface{}{
		"parcelId":  source.ParcelID,
		"adrSearch": source.AdrSearch,
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
		Scope:             seo_domain.SeoScopePublic,
		Classify:          enums.ESeoClassifyA,
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
		log.Printf("[SEO_GENERATE] creating seo_domain parcelId=%d slug=%q canonicalUrl=%s", source.ParcelID, seoDomain.Slug, seoDomain.CanonicalURL)
		createdSeo, err := u.seoDomainRepo.Create(txCtx, seoDomain)
		if err != nil {
			log.Printf("[SEO_GENERATE] create seo_domain failed parcelId=%d canonicalUrl=%s error=%v", source.ParcelID, seoDomain.CanonicalURL, err)
			return err
		}
		if createdSeo == nil || createdSeo.ID == 0 {
			log.Printf("[SEO_GENERATE] create seo_domain returned empty parcelId=%d canonicalUrl=%s", source.ParcelID, seoDomain.CanonicalURL)
			return _errors.ReturnError(500, "tạo seo domain thất bại")
		}
		log.Printf("[SEO_GENERATE] create seo_domain success parcelId=%d seoId=%d", source.ParcelID, createdSeo.ID)
		if err := u.tqdProvider.UpdateParcelSeoID(txCtx, source.ParcelID, createdSeo.ID); err != nil {
			log.Printf("[SEO_GENERATE] update tqd seo_id after create failed parcelId=%d seoId=%d error=%v", source.ParcelID, createdSeo.ID, err)
			return err
		}
		log.Printf("[SEO_GENERATE] update tqd seo_id after create success parcelId=%d seoId=%d", source.ParcelID, createdSeo.ID)
		created = createdSeo
		return nil
	})
	if err != nil {
		log.Printf("[SEO_GENERATE] transaction failed parcelId=%d error=%v", source.ParcelID, err)
		return nil, false, err
	}

	return created, true, nil
}

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
	if !enums.IsValidSEORefType(req.RefType) || req.RefType == 0 {
		return _errors.ReturnError(400, "refType không hợp lệ")
	}
	if req.RefID == 0 {
		return _errors.ReturnError(400, "refId không hợp lệ")
	}
	return u.seoDomainRepo.MarkNeedGenerateByRef(ctx, req.RefType, req.RefID, time.Now())
}

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
// Internal links
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
		return nil, _errors.ReturnError(400, "internal link đã tồn tại")
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
		return nil, _errors.ReturnError(400, "internal link đã tồn tại")
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
// Relatives
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
		return nil, _errors.ReturnError(400, "seo relative đã tồn tại")
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
// Builders / validators
// -----------------------------------------------------------------------------

func (u *SeoDomainUsecase) buildSeoDomainForCreate(ctx context.Context, req *dto.SeoDomainCreateRequest) (*seo_domain.SeoDomain, error) {
	originURL := strings.TrimSpace(req.OriginURL)
	if originURL == "" {
		return nil, _errors.ReturnError(400, "originUrl không được để trống")
	}

	canonicalURL := originURL
	if req.CanonicalURL != nil && strings.TrimSpace(*req.CanonicalURL) != "" {
		canonicalURL = strings.TrimSpace(*req.CanonicalURL)
	}

	slug := strings.TrimSpace(req.Slug)
	if slug == "" {
		slug = slugify(canonicalURL)
	}

	scope := normalizeScope(req.Scope)
	classify := enums.ESeoClassifyN
	if req.Classify != 0 {
		if !enums.IsValidSeoClassify(req.Classify) {
			return nil, _errors.ReturnError(400, "classify không hợp lệ")
		}
		classify = enums.SeoClassifyFromUint32(req.Classify)
	}

	if err := validateRef(req.RefType, req.RefID); err != nil {
		return nil, err
	}

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
	metadata, err := parseMetadata(req.Metadata)
	if err != nil {
		return nil, err
	}

	needGenerate := true
	if req.NeedGenerate != nil {
		needGenerate = *req.NeedGenerate
	}

	sitemapPriority := req.SitemapPriority
	if sitemapPriority == 0 {
		sitemapPriority = 0.5
	}
	changeFreq := normalizeChangeFreq(req.SitemapChangeFreq)

	state := seoState{
		Slug:      slug,
		Canonical: canonicalURL,
		Title:     strings.TrimSpace(req.Title),
		Published: req.Published,
		IsSiteMap: req.IsSiteMap,
		IsIndex:   req.IsIndex,
		Classify:  classify,
	}
	if err := validateSeoState(state); err != nil {
		return nil, err
	}

	exist, err := u.seoDomainRepo.ExistByCanonicalURL(ctx, canonicalURL, nil)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, _errors.ReturnError(400, "canonicalUrl đã tồn tại")
	}

	return &seo_domain.SeoDomain{
		Slug:              slug,
		OriginURL:         originURL,
		CanonicalURL:      canonicalURL,
		RefType:           enums.ESEORefType(req.RefType),
		RefID:             req.RefID,
		Scope:             scope,
		Classify:          classify,
		Title:             strings.TrimSpace(req.Title),
		Description:       strings.TrimSpace(req.Description),
		Content:           req.Content,
		Summary:           strings.TrimSpace(req.Summary),
		Published:         req.Published,
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
		StaticHtmlPath:    strings.TrimSpace(req.StaticHtmlPath),
		StaticHtmlHash:    strings.TrimSpace(req.StaticHtmlHash),
		DeepLink:          strings.TrimSpace(req.DeepLink),
		Metadata:          metadata,
		Note:              req.Note,
	}, nil
}

func (u *SeoDomainUsecase) buildSeoDomainUpdates(ctx context.Context, existing *seo_domain.SeoDomain, req *dto.SeoDomainUpdateRequest) (map[string]interface{}, error) {
	updates := map[string]interface{}{}

	state := seoState{
		Slug:      existing.Slug,
		Canonical: existing.CanonicalURL,
		Title:     existing.Title,
		Published: existing.Published,
		IsSiteMap: existing.IsSiteMap,
		IsIndex:   existing.IsIndex,
		Classify:  existing.Classify,
	}
	refType := uint32(existing.RefType)
	refID := existing.RefID

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
			return nil, _errors.ReturnError(400, "canonicalUrl đã tồn tại")
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
	if err := validateRef(refType, refID); err != nil {
		return nil, err
	}
	if req.Scope != nil {
		updates["scope"] = normalizeScope(*req.Scope)
	}
	if req.Classify != nil {
		if !enums.IsValidSeoClassify(*req.Classify) {
			return nil, _errors.ReturnError(400, "classify không hợp lệ")
		}
		state.Classify = enums.SeoClassifyFromUint32(*req.Classify)
		updates["classify"] = state.Classify
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
	if parent.Classify == enums.ESeoClassifyN || child.Classify == enums.ESeoClassifyN {
		return nil, nil, "", 0, _errors.ReturnError(400, "không tạo internal link cho page SEO-N")
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
	Slug      string
	Canonical string
	Title     string
	Published bool
	IsSiteMap bool
	IsIndex   bool
	Classify  enums.ESeoClassify
}

func validateSeoState(state seoState) error {
	if state.Classify == enums.ESeoClassifyN {
		if state.Published || state.IsIndex || state.IsSiteMap {
			return _errors.ReturnError(400, "classify=N không được publish/index/sitemap")
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

func normalizeScope(scope string) string {
	switch strings.TrimSpace(scope) {
	case seo_domain.SeoScopePublic, seo_domain.SeoScopeInternal, seo_domain.SeoScopeApp, seo_domain.SeoScopeAdmin:
		return strings.TrimSpace(scope)
	default:
		return seo_domain.SeoScopePublic
	}
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
