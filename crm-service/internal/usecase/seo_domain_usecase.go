package usecase

import (
	_errors "common/errors"
	_utils "common/utils"
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
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/spf13/viper"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
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

	existing, err := u.findExistingSeoForCreate(ctx, seo)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		// Upsert trên bản ghi đã gắn refType/refId: cập nhật field khác, không đụng category_id (chủ đề).
		updates := seoDomainToCreateUpdates(seo)
		if err := u.updateSeoDomainPreserveCategory(ctx, existing.ID, updates); err != nil {
			return nil, err
		}
		updated, err := u.seoDomainRepo.GetByID(ctx, existing.ID)
		if err != nil {
			return nil, err
		}
		return u.finalizeManualTinTucPaths(ctx, updated)
	}

	created, err := u.seoDomainRepo.Create(ctx, seo)
	if err != nil {
		return nil, err
	}
	return u.finalizeManualTinTucPaths(ctx, created)
}

// findExistingSeoForCreate tìm SEO đã tồn tại theo refType+refId (ưu tiên), rồi refSource, rồi canonical URL.
func (u *SeoDomainUsecase) findExistingSeoForCreate(ctx context.Context, seo *seo_domain.SeoDomain) (*seo_domain.SeoDomain, error) {
	if seo == nil {
		return nil, nil
	}

	if seo.RefID != nil && *seo.RefID > 0 && uint32(seo.RefType) != 0 {
		existing, err := u.seoDomainRepo.GetByRef(ctx, uint32(seo.RefType), *seo.RefID)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return existing, nil
		}
	}

	if seo.RefSource != "" && seo.RefSource != seo_domain.SeoRefSourceManual && seo.RefID != nil && *seo.RefID > 0 {
		existing, err := u.seoDomainRepo.GetByRefSource(ctx, uint32(seo.RefType), seo.RefSource, *seo.RefID)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return existing, nil
		}
	}

	if canonical := strings.TrimSpace(seo.CanonicalURL); canonical != "" {
		// Manual create dùng slug segment tạm thời — không upsert theo canonical tạm.
		if !enums.IsManualSEORefType(uint32(seo.RefType)) {
			existing, err := u.seoDomainRepo.GetByCanonicalURL(ctx, canonical)
			if err != nil {
				return nil, err
			}
			if existing != nil {
				return existing, nil
			}
		}
	}

	return nil, nil
}

// updateSeoDomainPreserveCategory updates seo_domain on re-create/sync but never
// overwrites chủ đề (category_id) or publish lifecycle (published / published_at / page_status).
func (u *SeoDomainUsecase) updateSeoDomainPreserveCategory(ctx context.Context, id uint64, updates map[string]interface{}) error {
	if len(updates) == 0 {
		return nil
	}
	safe := make(map[string]interface{}, len(updates))
	for k, v := range updates {
		switch k {
		case "category_id", "CategoryID",
			"published", "Published",
			"published_at", "PublishedAt",
			"page_status", "PageStatus":
			continue
		}
		safe[k] = v
	}
	if len(safe) == 0 {
		return nil
	}
	return u.seoDomainRepo.Update(ctx, id, safe)
}

func seoDomainToCreateUpdates(seo *seo_domain.SeoDomain) map[string]interface{} {
	if seo == nil {
		return map[string]interface{}{}
	}
	return map[string]interface{}{
		"slug":                seo.Slug,
		"origin_url":          seo.OriginURL,
		"canonical_url":       seo.CanonicalURL,
		"ref_type":            seo.RefType,
		"ref_id":              seo.RefID,
		"ref_source":          seo.RefSource,
		"ref_label":           seo.RefLabel,
		"ref_url":             seo.RefURL,
		"source_status":       seo.SourceStatus,
		"ref_missing":         seo.RefMissing,
		"ref_snapshot_json":   seo.RefSnapshotJSON,
		"ref_hash":            seo.RefHash,
		"ref_last_synced_at":  seo.RefLastSyncedAt,
		"scope":               seo.Scope,
		"page_status":         seo.PageStatus,
		"title":               seo.Title,
		"description":         seo.Description,
		"content":             seo.Content,
		"summary":             seo.Summary,
		"published":           seo.Published,
		"published_at":        seo.PublishedAt,
		"is_site_map":         seo.IsSiteMap,
		"is_index":            seo.IsIndex,
		"is_robot":            seo.IsRobot,
		"site_map_lasted_at":  seo.SiteMapLastedAt,
		"sitemap_priority":    seo.SitemapPriority,
		"sitemap_change_freq": seo.SitemapChangeFreq,
		"need_generate":       seo.NeedGenerate,
		"generated_at":        seo.GeneratedAt,
		"source_updated_at":   seo.SourceUpdatedAt,
		"render_status":       seo.RenderStatus,
		"rendered_html":       seo.RenderedHTML,
		"last_render_error":   seo.LastRenderError,
		"static_html_path":    seo.StaticHtmlPath,
		"static_html_hash":    seo.StaticHtmlHash,
		"template_key":        seo.TemplateKey,
		"template_version":    seo.TemplateVersion,
		"deep_link":           seo.DeepLink,
		"image_url":           seo.ImageURL,
		"metadata":            seo.Metadata,
		"note":                seo.Note,
		"visible_on":          seo.VisibleOn,
	}
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

// CountRecentSeoByRefType đếm SEO theo từng refType (8 type) có updatedAt trong N ngày gần nhất (mặc định 5).
func (u *SeoDomainUsecase) CountRecentSeoByRefType(ctx context.Context, req *dto.CountRecentSeoByRefTypeRequest) ([]dto.SeoRefTypeCountItem, uint32, error) {
	if req == nil {
		req = &dto.CountRecentSeoByRefTypeRequest{}
	}
	req.Normalize()
	items, err := u.seoDomainRepo.CountRecentByRefType(ctx, req.Days)
	if err != nil {
		return nil, 0, err
	}
	return items, req.Days, nil
}

const seoTopClickedLimit = 15

func seoClickStatsFromTime(lastDates uint32) time.Time {
	now := time.Now()
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		loc = now.Location()
	}
	today := time.Date(now.In(loc).Year(), now.In(loc).Month(), now.In(loc).Day(), 0, 0, 0, 0, loc)
	// Inclusive window: lastDates days ending today.
	return today.AddDate(0, 0, -int(lastDates)+1)
}

func buildSeoClickDailySeries(fromTime time.Time, lastDates uint32, rows []dto.SeoClickDailyItem) []dto.SeoClickDailyItem {
	byDate := make(map[string]uint64, len(rows))
	for _, row := range rows {
		byDate[row.Date] = row.Count
	}

	data := make([]dto.SeoClickDailyItem, 0, lastDates)
	for i := 0; i < int(lastDates); i++ {
		day := fromTime.AddDate(0, 0, i)
		date := day.Format("2006-01-02")
		data = append(data, dto.SeoClickDailyItem{
			Date:  date,
			Count: byDate[date],
		})
	}
	return data
}

// GetSeoTopClicked thống kê tổng click và top trang click nhiều nhất từ seo_count_log theo lastDates ngày.
func (u *SeoDomainUsecase) GetSeoTopClicked(ctx context.Context, req *dto.GetSeoTopClickedRequest) (*dto.GetSeoTopClickedResponse, error) {
	if req == nil {
		req = &dto.GetSeoTopClickedRequest{}
	}
	req.Normalize()

	fromTime := seoClickStatsFromTime(req.LastDates)

	totalClicks, err := u.seoDomainRepo.CountClicksFromTime(ctx, fromTime)
	if err != nil {
		return nil, err
	}

	items, err := u.seoDomainRepo.GetTopClicked(ctx, fromTime, seoTopClickedLimit)
	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []dto.SeoTopClickedItem{}
	}

	return &dto.GetSeoTopClickedResponse{
		TotalClicks: totalClicks,
		Data:        items,
		LastDates:   req.LastDates,
	}, nil
}

// GetSeoClickDailyStats trả về tổng click theo từng ngày để vẽ biểu đồ.
func (u *SeoDomainUsecase) GetSeoClickDailyStats(ctx context.Context, req *dto.GetSeoClickDailyStatsRequest) (*dto.GetSeoClickDailyStatsResponse, error) {
	if req == nil {
		req = &dto.GetSeoClickDailyStatsRequest{}
	}
	req.Normalize()

	fromTime := seoClickStatsFromTime(req.LastDates)
	rows, err := u.seoDomainRepo.GetClickDailyStats(ctx, fromTime)
	if err != nil {
		return nil, err
	}

	return &dto.GetSeoClickDailyStatsResponse{
		Data:      buildSeoClickDailySeries(fromTime, req.LastDates, rows),
		LastDates: req.LastDates,
	}, nil
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
	parcelIDs := []uint64(nil)
	if req != nil {
		limit = req.Limit
		parcelIDs = req.ParcelIDs
	}

	var sources []dto.ParcelSeoSource
	var err error
	if len(parcelIDs) > 0 {
		log.Printf("[SEO_GENERATE] start generate seo from explicit parcelIds total=%d", len(parcelIDs))
		sources = make([]dto.ParcelSeoSource, 0, len(parcelIDs))
		for _, parcelID := range parcelIDs {
			if parcelID == 0 {
				continue
			}
			source, sourceErr := u.tqdProvider.GetParcelSeoSource(ctx, parcelID)
			if sourceErr != nil {
				log.Printf("[SEO_GENERATE] get parcel seo source failed parcelId=%d error=%v", parcelID, sourceErr)
				return nil, nil, sourceErr
			}
			if source == nil {
				log.Printf("[SEO_GENERATE] parcel seo source not found parcelId=%d", parcelID)
				continue
			}
			sources = append(sources, *source)
		}
	} else {
		log.Printf("[SEO_GENERATE] start generate seo from parcels limit=%d", limit)
		sources, err = u.tqdProvider.GetParcelSeoSourcesForGenerate(ctx, limit)
		if err != nil {
			log.Printf("[SEO_GENERATE] get parcel seo sources failed: %v", err)
			return nil, nil, err
		}
	}
	log.Printf("[SEO_GENERATE] received parcel seo sources total=%d", len(sources))

	result := &dto.GenerateSeoFromParcelsResult{
		Total: int64(len(sources)),
	}
	processedItems := make([]seo_domain.SeoDomain, 0, len(sources))

	for i := range sources {
		log.Printf("[SEO_GENERATE] processing parcelId=%d adrSearch=%q seoId=%v", sources[i].ParcelID, sources[i].AdrSearch, sources[i].SeoID)

		seoDomain, created, err := u.createSeoFromParcelSource(ctx, sources[i])
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("parcelId=%d: %v", sources[i].ParcelID, err))
			log.Printf("[SEO_GENERATE] failed parcelId=%d error=%v", sources[i].ParcelID, err)
			continue
		}
		if seoDomain == nil {
			result.Skipped++
			log.Printf("[SEO_GENERATE] skipped parcelId=%d reason=seoDomain_nil", sources[i].ParcelID)
			continue
		}
		if created {
			result.Created++
			log.Printf("[SEO_GENERATE] created parcelId=%d seoId=%d canonicalUrl=%s", sources[i].ParcelID, seoDomain.ID, seoDomain.CanonicalURL)
		} else {
			log.Printf("[SEO_GENERATE] updated parcelId=%d seoId=%d canonicalUrl=%s", sources[i].ParcelID, seoDomain.ID, seoDomain.CanonicalURL)
		}
		processedItems = append(processedItems, *seoDomain)
	}

	log.Printf("[SEO_GENERATE] done total=%d created=%d skipped=%d failed=%d", result.Total, result.Created, result.Skipped, result.Failed)
	return processedItems, result, nil
}

func (u *SeoDomainUsecase) GenerateSeoFromLayers(ctx context.Context, req *dto.GenerateSeoFromLayersRequest) ([]seo_domain.SeoDomain, *dto.GenerateSeoFromLayersResult, error) {
	if u.tqdProvider == nil {
		return nil, nil, _errors.ReturnError(503, "TQD provider chưa được cấu hình")
	}

	layerIDs := []uint64(nil)
	if req != nil {
		layerIDs = req.LayerIDs
	}
	if len(layerIDs) == 0 {
		return nil, nil, _errors.ReturnError(400, "layerIds không được để trống")
	}

	log.Printf("[SEO_GENERATE_LAYER] start generate seo from explicit layerIds total=%d", len(layerIDs))
	sources := make([]dto.LayerSeoSource, 0, len(layerIDs))
	for _, layerID := range layerIDs {
		if layerID == 0 {
			continue
		}
		source, sourceErr := u.tqdProvider.GetLayerSeoSource(ctx, layerID)
		if sourceErr != nil {
			log.Printf("[SEO_GENERATE_LAYER] get layer seo source failed layerId=%d error=%v", layerID, sourceErr)
			return nil, nil, sourceErr
		}
		if source == nil {
			log.Printf("[SEO_GENERATE_LAYER] layer seo source not found layerId=%d", layerID)
			continue
		}
		sources = append(sources, *source)
	}
	log.Printf("[SEO_GENERATE_LAYER] received layer seo sources total=%d", len(sources))

	result := &dto.GenerateSeoFromLayersResult{
		Total: int64(len(sources)),
	}
	processedItems := make([]seo_domain.SeoDomain, 0, len(sources))

	for i := range sources {
		log.Printf("[SEO_GENERATE_LAYER] processing layerId=%d displayName=%q seoId=%v", sources[i].LayerID, sources[i].DisplayName, sources[i].SeoID)

		seoDomain, created, err := u.createSeoFromLayerSource(ctx, sources[i])
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("layerId=%d: %v", sources[i].LayerID, err))
			log.Printf("[SEO_GENERATE_LAYER] failed layerId=%d error=%v", sources[i].LayerID, err)
			continue
		}
		if seoDomain == nil {
			result.Skipped++
			continue
		}
		if created {
			result.Created++
		} else {
			result.Skipped++
		}
		processedItems = append(processedItems, *seoDomain)
	}

	log.Printf("[SEO_GENERATE_LAYER] done total=%d created=%d skipped=%d failed=%d", result.Total, result.Created, result.Skipped, result.Failed)
	return processedItems, result, nil
}

func (u *SeoDomainUsecase) GenerateSeoFromPlanningProjects(ctx context.Context, req *dto.GenerateSeoFromPlanningProjectsRequest) ([]seo_domain.SeoDomain, *dto.GenerateSeoFromPlanningProjectsResult, error) {
	if u.tqdProvider == nil {
		return nil, nil, _errors.ReturnError(503, "TQD provider chưa được cấu hình")
	}

	planningProjectIDs := []uint64(nil)
	if req != nil {
		planningProjectIDs = req.PlanningProjectIDs
	}
	if len(planningProjectIDs) == 0 {
		return nil, nil, _errors.ReturnError(400, "planningProjectIds không được để trống")
	}

	log.Printf("[SEO_GENERATE_PLANNING_PROJECT] start generate seo from explicit planningProjectIds total=%d", len(planningProjectIDs))
	sources := make([]dto.PlanningProjectSeoSource, 0, len(planningProjectIDs))
	for _, planningProjectID := range planningProjectIDs {
		if planningProjectID == 0 {
			continue
		}
		source, sourceErr := u.tqdProvider.GetPlanningProjectSeoSource(ctx, planningProjectID)
		if sourceErr != nil {
			log.Printf("[SEO_GENERATE_PLANNING_PROJECT] get planning project seo source failed planningProjectId=%d error=%v", planningProjectID, sourceErr)
			return nil, nil, sourceErr
		}
		if source == nil {
			log.Printf("[SEO_GENERATE_PLANNING_PROJECT] planning project seo source not found planningProjectId=%d", planningProjectID)
			continue
		}
		sources = append(sources, *source)
	}
	log.Printf("[SEO_GENERATE_PLANNING_PROJECT] received planning project seo sources total=%d", len(sources))

	result := &dto.GenerateSeoFromPlanningProjectsResult{
		Total: int64(len(sources)),
	}
	processedItems := make([]seo_domain.SeoDomain, 0, len(sources))

	for i := range sources {
		seoDomain, created, err := u.createSeoFromPlanningProjectSource(ctx, sources[i])
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("planningProjectId=%d: %v", sources[i].PlanningProjectID, err))
			continue
		}
		if seoDomain == nil {
			result.Skipped++
			continue
		}
		if created {
			result.Created++
		} else {
			result.Skipped++
		}
		processedItems = append(processedItems, *seoDomain)
	}

	log.Printf("[SEO_GENERATE_PLANNING_PROJECT] done total=%d created=%d skipped=%d failed=%d", result.Total, result.Created, result.Skipped, result.Failed)
	return processedItems, result, nil
}

func (u *SeoDomainUsecase) GenerateSeoFromPlanningDocuments(ctx context.Context, req *dto.GenerateSeoFromPlanningDocumentsRequest) ([]seo_domain.SeoDomain, *dto.GenerateSeoFromPlanningDocumentsResult, error) {
	if u.tqdProvider == nil {
		return nil, nil, _errors.ReturnError(503, "TQD provider chưa được cấu hình")
	}

	docIDs := []uint64(nil)
	if req != nil {
		docIDs = req.DocIDs
	}
	if len(docIDs) == 0 {
		return nil, nil, _errors.ReturnError(400, "docIds không được để trống")
	}

	log.Printf("[SEO_GENERATE_PLANNING_DOCUMENT] start generate seo from explicit docIds total=%d", len(docIDs))
	sources := make([]dto.PlanningDocumentSeoSource, 0, len(docIDs))
	for _, documentID := range docIDs {
		if documentID == 0 {
			continue
		}
		source, sourceErr := u.tqdProvider.GetPlanningDocumentSeoSource(ctx, documentID)
		if sourceErr != nil {
			log.Printf("[SEO_GENERATE_PLANNING_DOCUMENT] get planning document seo source failed docId=%d error=%v", documentID, sourceErr)
			return nil, nil, sourceErr
		}
		if source == nil {
			log.Printf("[SEO_GENERATE_PLANNING_DOCUMENT] planning document seo source not found docId=%d", documentID)
			continue
		}
		sources = append(sources, *source)
	}
	log.Printf("[SEO_GENERATE_PLANNING_DOCUMENT] received planning document seo sources total=%d", len(sources))

	result := &dto.GenerateSeoFromPlanningDocumentsResult{
		Total: int64(len(sources)),
	}
	processedItems := make([]seo_domain.SeoDomain, 0, len(sources))

	for i := range sources {
		seoDomain, created, err := u.createSeoFromPlanningDocumentSource(ctx, sources[i])
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("docId=%d: %v", sources[i].PlanningDocumentID, err))
			continue
		}
		if seoDomain == nil {
			result.Skipped++
			continue
		}
		if created {
			result.Created++
		} else {
			result.Skipped++
		}
		processedItems = append(processedItems, *seoDomain)
	}

	log.Printf("[SEO_GENERATE_PLANNING_DOCUMENT] done total=%d created=%d skipped=%d failed=%d", result.Total, result.Created, result.Skipped, result.Failed)
	return processedItems, result, nil
}

func (u *SeoDomainUsecase) GenerateSeoFromAdministrativeUnits(ctx context.Context, req *dto.GenerateSeoFromAdministrativeUnitsRequest) ([]seo_domain.SeoDomain, *dto.GenerateSeoFromAdministrativeUnitsResult, error) {
	if u.tqdProvider == nil {
		return nil, nil, _errors.ReturnError(503, "TQD provider chưa được cấu hình")
	}

	compositeIDs := []uint64(nil)
	if req != nil {
		compositeIDs = req.CompositeIDs
	}
	if len(compositeIDs) == 0 {
		return nil, nil, _errors.ReturnError(400, "compositeIds không được để trống")
	}

	log.Printf("[SEO_GENERATE_ADM_UNIT] start generate seo from compositeIds total=%d", len(compositeIDs))
	sources := make([]dto.AdministrativeUnitSeoSource, 0, len(compositeIDs))
	for _, compositeID := range compositeIDs {
		if _, _, err := _utils.DecodeLocationID(compositeID); err != nil {
			log.Printf("[SEO_GENERATE_ADM_UNIT] invalid compositeId=%d error=%v", compositeID, err)
			return nil, nil, _errors.ReturnError(400, "compositeId không hợp lệ")
		}
		source, sourceErr := u.tqdProvider.GetAdministrativeUnitSeoSource(ctx, compositeID)
		if sourceErr != nil {
			log.Printf("[SEO_GENERATE_ADM_UNIT] get seo source failed compositeId=%d error=%v", compositeID, sourceErr)
			return nil, nil, sourceErr
		}
		if source == nil {
			log.Printf("[SEO_GENERATE_ADM_UNIT] seo source not found compositeId=%d", compositeID)
			continue
		}
		sources = append(sources, *source)
	}
	log.Printf("[SEO_GENERATE_ADM_UNIT] received seo sources total=%d", len(sources))

	result := &dto.GenerateSeoFromAdministrativeUnitsResult{
		Total: int64(len(sources)),
	}
	processedItems := make([]seo_domain.SeoDomain, 0, len(sources))

	for i := range sources {
		seoDomain, created, err := u.createSeoFromAdministrativeUnitSource(ctx, sources[i])
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("%d: %v", sources[i].CompositeID, err))
			continue
		}
		if seoDomain == nil {
			result.Skipped++
			continue
		}
		if created {
			result.Created++
		} else {
			result.Skipped++
		}
		processedItems = append(processedItems, *seoDomain)
	}

	log.Printf("[SEO_GENERATE_ADM_UNIT] done total=%d created=%d skipped=%d failed=%d", result.Total, result.Created, result.Skipped, result.Failed)
	return processedItems, result, nil
}

func (u *SeoDomainUsecase) createSeoFromParcelSource(ctx context.Context, source dto.ParcelSeoSource) (*seo_domain.SeoDomain, bool, error) {
	if u.tqdProvider == nil {
		return nil, false, _errors.ReturnError(503, "TQD provider chưa được cấu hình")
	}

	parcelURL, err := buildParcelSeoURL(source.ParcelID, source.AdrSearch)
	if err != nil {
		log.Printf("[SEO_GENERATE] build parcel seo url failed parcelId=%d adrSearch=%q error=%v", source.ParcelID, source.AdrSearch, err)
		return nil, false, err
	}
	log.Printf("[SEO_GENERATE] built parcel url parcelId=%d url=%s", source.ParcelID, parcelURL)

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
	parcelUpdates := map[string]interface{}{
		"slug":                resolveParcelSeoSlug(source.AdrSearch, source.ParcelID),
		"origin_url":          parcelURL,
		"canonical_url":       parcelURL,
		"ref_type":            uint32(enums.ESEORefTypeParcel),
		"ref_id":              refID,
		"ref_source":          seo_domain.SeoRefSourceTqdParcel,
		"ref_label":           firstNonEmpty(source.AdrSearch, fmt.Sprintf("Thửa đất #%d", source.ParcelID)),
		"ref_url":             parcelURL,
		"source_status":       seo_domain.SeoSourceStatusLinked,
		"ref_missing":         false,
		"ref_snapshot_json":   datatypes.JSON(metadataBytes),
		"ref_hash":            hashString(string(metadataBytes)),
		"ref_last_synced_at":  now,
		"scope":               seo_domain.SeoScopePublic,
		"title":               buildParcelSeoTitle(source.AdrSearch, source.ParcelID),
		"description":         buildParcelSeoDescription(source.AdrSearch, source.ParcelID),
		"published":           true,
		"published_at":        now,
		"is_site_map":         true,
		"is_index":            true,
		"site_map_lasted_at":  now,
		"sitemap_priority":    0.8,
		"sitemap_change_freq": seo_domain.SeoChangeFreqDaily,
		"need_generate":       true,
		"source_updated_at":   now,
		"metadata":            datatypes.JSON(metadataBytes),
	}

	// Nếu parcel đã có seoId, gen lại và update bản ghi đó.
	if source.SeoID != nil && *source.SeoID > 0 {
		log.Printf("[SEO_GENERATE] parcelId=%d has seoId=%d, checking existing seo_domain", source.ParcelID, *source.SeoID)

		existing, err := u.seoDomainRepo.GetByID(ctx, *source.SeoID)
		if err != nil {
			log.Printf("[SEO_GENERATE] get existing by seoId failed parcelId=%d seoId=%d error=%v", source.ParcelID, *source.SeoID, err)
			return nil, false, err
		}
		if existing != nil {
			if err := u.updateSeoDomainPreserveCategory(ctx, existing.ID, parcelUpdates); err != nil {
				log.Printf("[SEO_GENERATE] update existing by seoId failed parcelId=%d seoId=%d error=%v", source.ParcelID, existing.ID, err)
				return nil, false, err
			}
			if err := u.tqdProvider.UpdateParcelSeoID(ctx, source.ParcelID, existing.ID); err != nil {
				log.Printf("[SEO_GENERATE] update tqd seo_id after regenerate failed parcelId=%d seoId=%d error=%v", source.ParcelID, existing.ID, err)
				return nil, false, err
			}
			updated, err := u.seoDomainRepo.GetByID(ctx, existing.ID)
			if err != nil {
				return nil, false, err
			}
			log.Printf("[SEO_GENERATE] updated existing by seoId parcelId=%d seoId=%d canonicalUrl=%s", source.ParcelID, existing.ID, updated.CanonicalURL)
			return updated, false, nil
		}
		log.Printf("[SEO_GENERATE] parcelId=%d seoId=%d not found in crm, will recreate", source.ParcelID, *source.SeoID)
	}

	existingByURL, err := u.seoDomainRepo.GetByCanonicalURL(ctx, parcelURL)
	if err != nil {
		log.Printf("[SEO_GENERATE] get existing by canonical failed parcelId=%d canonicalUrl=%s error=%v", source.ParcelID, parcelURL, err)
		return nil, false, err
	}
	if existingByURL != nil {
		log.Printf("[SEO_GENERATE] existing canonical found parcelId=%d seoId=%d, regenerating seo_domain", source.ParcelID, existingByURL.ID)
		if err := u.updateSeoDomainPreserveCategory(ctx, existingByURL.ID, parcelUpdates); err != nil {
			log.Printf("[SEO_GENERATE] update existing by canonical failed parcelId=%d seoId=%d error=%v", source.ParcelID, existingByURL.ID, err)
			return nil, false, err
		}
		if err := u.tqdProvider.UpdateParcelSeoID(ctx, source.ParcelID, existingByURL.ID); err != nil {
			log.Printf("[SEO_GENERATE] update tqd seo_id failed parcelId=%d seoId=%d error=%v", source.ParcelID, existingByURL.ID, err)
			return nil, false, err
		}
		updated, err := u.seoDomainRepo.GetByID(ctx, existingByURL.ID)
		if err != nil {
			return nil, false, err
		}
		log.Printf("[SEO_GENERATE] updated existing by canonical parcelId=%d seoId=%d", source.ParcelID, existingByURL.ID)
		return updated, false, nil
	}

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

func (u *SeoDomainUsecase) createSeoFromLayerSource(ctx context.Context, source dto.LayerSeoSource) (*seo_domain.SeoDomain, bool, error) {
	if u.tqdProvider == nil {
		return nil, false, _errors.ReturnError(503, "TQD provider chưa được cấu hình")
	}

	layerURL, err := buildLayerSeoURL(source.LayerID, source.DisplayName)
	if err != nil {
		return nil, false, err
	}

	metadataBytes, err := json.Marshal(map[string]interface{}{
		"layerId":     source.LayerID,
		"displayName": source.DisplayName,
		"source":      seo_domain.SeoRefSourceTqdLayer,
		"sourceKey":   "layer",
	})
	if err != nil {
		return nil, false, err
	}

	now := time.Now()
	refID := source.LayerID
	imageURL := resolveLayerSeoImageURL(source)

	layerUpdates := map[string]interface{}{
		"slug":                resolveLayerSeoSlug(source.DisplayName, source.LayerID),
		"origin_url":          layerURL,
		"canonical_url":       layerURL,
		"ref_type":            uint32(enums.ESEORefTypeMap),
		"ref_id":              refID,
		"ref_source":          seo_domain.SeoRefSourceTqdLayer,
		"ref_label":           firstNonEmpty(source.DisplayName, fmt.Sprintf("Lớp quy hoạch #%d", source.LayerID)),
		"ref_url":             layerURL,
		"source_status":       seo_domain.SeoSourceStatusLinked,
		"ref_missing":         false,
		"ref_snapshot_json":   datatypes.JSON(metadataBytes),
		"ref_hash":            hashString(string(metadataBytes)),
		"ref_last_synced_at":  now,
		"scope":               seo_domain.SeoScopePublic,
		"title":               buildLayerSeoTitle(source.DisplayName, source.LayerID),
		"description":         buildLayerSeoDescription(source.DisplayName, source.LayerID),
		"published":           true,
		"published_at":        now,
		"is_site_map":         true,
		"is_index":            true,
		"site_map_lasted_at":  now,
		"sitemap_priority":    0.7,
		"sitemap_change_freq": seo_domain.SeoChangeFreqWeekly,
		"need_generate":       true,
		"source_updated_at":   now,
		"image_url":           imageURL,
		"metadata":            datatypes.JSON(metadataBytes),
	}
	if source.ProvinceID != nil && *source.ProvinceID > 0 {
		layerUpdates["province_id"] = *source.ProvinceID
	}

	// Nếu layer đã có seoId, update lại bản ghi đó.
	if source.SeoID != nil && *source.SeoID > 0 {
		existing, err := u.seoDomainRepo.GetByID(ctx, *source.SeoID)
		if err != nil {
			return nil, false, err
		}
		if existing != nil {
			if err := u.updateSeoDomainPreserveCategory(ctx, existing.ID, layerUpdates); err != nil {
				return nil, false, err
			}
			if err := u.tqdProvider.UpdateLayerSeoID(ctx, source.LayerID, existing.ID); err != nil {
				return nil, false, err
			}
			updated, err := u.seoDomainRepo.GetByID(ctx, existing.ID)
			if err != nil {
				return nil, false, err
			}
			return updated, false, nil
		}
	}

	// Kiểm tra theo canonical URL.
	existingByURL, err := u.seoDomainRepo.GetByCanonicalURL(ctx, layerURL)
	if err != nil {
		return nil, false, err
	}
	if existingByURL != nil {
		if err := u.updateSeoDomainPreserveCategory(ctx, existingByURL.ID, layerUpdates); err != nil {
			return nil, false, err
		}
		if err := u.tqdProvider.UpdateLayerSeoID(ctx, source.LayerID, existingByURL.ID); err != nil {
			return nil, false, err
		}
		updated, err := u.seoDomainRepo.GetByID(ctx, existingByURL.ID)
		if err != nil {
			return nil, false, err
		}
		return updated, false, nil
	}

	seoDomain := &seo_domain.SeoDomain{
		Slug:              resolveLayerSeoSlug(source.DisplayName, source.LayerID),
		OriginURL:         layerURL,
		CanonicalURL:      layerURL,
		RefType:           enums.ESEORefTypeMap,
		RefID:             &refID,
		RefSource:         seo_domain.SeoRefSourceTqdLayer,
		RefLabel:          firstNonEmpty(source.DisplayName, fmt.Sprintf("Lớp quy hoạch #%d", source.LayerID)),
		RefURL:            layerURL,
		SourceStatus:      seo_domain.SeoSourceStatusLinked,
		RefMissing:        false,
		RefSnapshotJSON:   datatypes.JSON(metadataBytes),
		RefHash:           hashString(string(metadataBytes)),
		RefLastSyncedAt:   &now,
		Scope:             seo_domain.SeoScopePublic,
		Title:             buildLayerSeoTitle(source.DisplayName, source.LayerID),
		Description:       buildLayerSeoDescription(source.DisplayName, source.LayerID),
		Published:         true,
		PublishedAt:       &now,
		IsSiteMap:         true,
		IsIndex:           true,
		SiteMapLastedAt:   &now,
		SitemapPriority:   0.7,
		SitemapChangeFreq: seo_domain.SeoChangeFreqWeekly,
		NeedGenerate:      true,
		SourceUpdatedAt:   &now,
		ImageURL:          imageURL,
		Metadata:          datatypes.JSON(metadataBytes),
	}
	if source.ProvinceID != nil && *source.ProvinceID > 0 {
		seoDomain.ProvinceID = source.ProvinceID
	}

	var created *seo_domain.SeoDomain
	err = u.transaction.WithTransaction(ctx, func(txCtx context.Context) error {
		createdSeo, err := u.seoDomainRepo.Create(txCtx, seoDomain)
		if err != nil {
			return err
		}
		if createdSeo == nil || createdSeo.ID == 0 {
			return _errors.ReturnError(500, "tạo seo domain thất bại")
		}
		if err := u.tqdProvider.UpdateLayerSeoID(txCtx, source.LayerID, createdSeo.ID); err != nil {
			return err
		}
		created = createdSeo
		return nil
	})
	if err != nil {
		return nil, false, err
	}

	return created, true, nil
}

func (u *SeoDomainUsecase) createSeoFromPlanningProjectSource(ctx context.Context, source dto.PlanningProjectSeoSource) (*seo_domain.SeoDomain, bool, error) {
	if u.tqdProvider == nil {
		return nil, false, _errors.ReturnError(503, "TQD provider chưa được cấu hình")
	}

	projectURL, err := buildPlanningProjectSeoURL(source.PlanningProjectID, source.Slug, source.Name)
	if err != nil {
		return nil, false, err
	}

	metadataBytes, err := json.Marshal(map[string]interface{}{
		"planningProjectId": source.PlanningProjectID,
		"name":              source.Name,
		"slug":              source.Slug,
		"source":            seo_domain.SeoRefSourceTqdPlanningProject,
		"sourceKey":         "planning_project",
	})
	if err != nil {
		return nil, false, err
	}

	now := time.Now()
	refID := source.PlanningProjectID
	projectUpdates := map[string]interface{}{
		"slug":                resolvePlanningProjectSeoSlug(source.Slug, source.Name, source.PlanningProjectID),
		"origin_url":          projectURL,
		"canonical_url":       projectURL,
		"ref_type":            uint32(enums.ESEORefTypeProject),
		"ref_id":              refID,
		"ref_source":          seo_domain.SeoRefSourceTqdPlanningProject,
		"ref_label":           firstNonEmpty(source.Name, fmt.Sprintf("Đồ án quy hoạch #%d", source.PlanningProjectID)),
		"ref_url":             projectURL,
		"source_status":       seo_domain.SeoSourceStatusLinked,
		"ref_missing":         false,
		"ref_snapshot_json":   datatypes.JSON(metadataBytes),
		"ref_hash":            hashString(string(metadataBytes)),
		"ref_last_synced_at":  now,
		"scope":               seo_domain.SeoScopePublic,
		"title":               buildPlanningProjectSeoTitle(source.Name, source.PlanningProjectID),
		"description":         buildPlanningProjectSeoDescription(source.Name, source.PlanningProjectID),
		"published":           true,
		"published_at":        now,
		"is_site_map":         true,
		"is_index":            true,
		"site_map_lasted_at":  now,
		"sitemap_priority":    0.75,
		"sitemap_change_freq": seo_domain.SeoChangeFreqWeekly,
		"need_generate":       true,
		"source_updated_at":   now,
		"metadata":            datatypes.JSON(metadataBytes),
	}
	if source.ProvinceID != nil && *source.ProvinceID > 0 {
		projectUpdates["province_id"] = *source.ProvinceID
	}

	if source.SeoID != nil && *source.SeoID > 0 {
		existing, err := u.seoDomainRepo.GetByID(ctx, *source.SeoID)
		if err != nil {
			return nil, false, err
		}
		if existing != nil {
			if err := u.updateSeoDomainPreserveCategory(ctx, existing.ID, projectUpdates); err != nil {
				return nil, false, err
			}
			if err := u.tqdProvider.UpdatePlanningProjectSeoID(ctx, source.PlanningProjectID, existing.ID); err != nil {
				return nil, false, err
			}
			updated, err := u.seoDomainRepo.GetByID(ctx, existing.ID)
			if err != nil {
				return nil, false, err
			}
			log.Printf("[SEO_GENERATE_PLANNING_PROJECT] updated existing by seoId planningProjectId=%d seoId=%d canonicalUrl=%s", source.PlanningProjectID, existing.ID, updated.CanonicalURL)
			return updated, false, nil
		}
	}

	existingByURL, err := u.seoDomainRepo.GetByCanonicalURL(ctx, projectURL)
	if err != nil {
		return nil, false, err
	}
	if existingByURL != nil {
		if err := u.updateSeoDomainPreserveCategory(ctx, existingByURL.ID, projectUpdates); err != nil {
			return nil, false, err
		}
		if err := u.tqdProvider.UpdatePlanningProjectSeoID(ctx, source.PlanningProjectID, existingByURL.ID); err != nil {
			return nil, false, err
		}
		updated, err := u.seoDomainRepo.GetByID(ctx, existingByURL.ID)
		if err != nil {
			return nil, false, err
		}
		log.Printf("[SEO_GENERATE_PLANNING_PROJECT] updated existing by canonical planningProjectId=%d seoId=%d canonicalUrl=%s", source.PlanningProjectID, existingByURL.ID, updated.CanonicalURL)
		return updated, false, nil
	}

	seoDomain := &seo_domain.SeoDomain{
		Slug:              resolvePlanningProjectSeoSlug(source.Slug, source.Name, source.PlanningProjectID),
		OriginURL:         projectURL,
		CanonicalURL:      projectURL,
		RefType:           enums.ESEORefTypeProject,
		RefID:             &refID,
		RefSource:         seo_domain.SeoRefSourceTqdPlanningProject,
		RefLabel:          firstNonEmpty(source.Name, fmt.Sprintf("Đồ án quy hoạch #%d", source.PlanningProjectID)),
		RefURL:            projectURL,
		SourceStatus:      seo_domain.SeoSourceStatusLinked,
		RefMissing:        false,
		RefSnapshotJSON:   datatypes.JSON(metadataBytes),
		RefHash:           hashString(string(metadataBytes)),
		RefLastSyncedAt:   &now,
		Scope:             seo_domain.SeoScopePublic,
		Title:             buildPlanningProjectSeoTitle(source.Name, source.PlanningProjectID),
		Description:       buildPlanningProjectSeoDescription(source.Name, source.PlanningProjectID),
		Published:         true,
		PublishedAt:       &now,
		IsSiteMap:         true,
		IsIndex:           true,
		SiteMapLastedAt:   &now,
		SitemapPriority:   0.75,
		SitemapChangeFreq: seo_domain.SeoChangeFreqWeekly,
		NeedGenerate:      true,
		SourceUpdatedAt:   &now,
		Metadata:          datatypes.JSON(metadataBytes),
	}
	if source.ProvinceID != nil && *source.ProvinceID > 0 {
		seoDomain.ProvinceID = source.ProvinceID
	}

	var created *seo_domain.SeoDomain
	err = u.transaction.WithTransaction(ctx, func(txCtx context.Context) error {
		createdSeo, err := u.seoDomainRepo.Create(txCtx, seoDomain)
		if err != nil {
			return err
		}
		if createdSeo == nil || createdSeo.ID == 0 {
			return _errors.ReturnError(500, "tạo seo domain thất bại")
		}
		if err := u.tqdProvider.UpdatePlanningProjectSeoID(txCtx, source.PlanningProjectID, createdSeo.ID); err != nil {
			return err
		}
		created = createdSeo
		return nil
	})
	if err != nil {
		return nil, false, err
	}

	return created, true, nil
}

func (u *SeoDomainUsecase) createSeoFromPlanningDocumentSource(ctx context.Context, source dto.PlanningDocumentSeoSource) (*seo_domain.SeoDomain, bool, error) {
	if u.tqdProvider == nil {
		return nil, false, _errors.ReturnError(503, "TQD provider chưa được cấu hình")
	}

	documentURL, err := buildPlanningDocumentSeoURL(source.PlanningDocumentID, source.Slug, source.Title)
	if err != nil {
		return nil, false, err
	}

	metadataBytes, err := json.Marshal(map[string]interface{}{
		"planningDocumentId": source.PlanningDocumentID,
		"title":              source.Title,
		"slug":               source.Slug,
		"source":             seo_domain.SeoRefSourceTqdDocument,
		"sourceKey":          "document",
	})
	if err != nil {
		return nil, false, err
	}

	now := time.Now()
	refID := source.PlanningDocumentID
	imageURL := resolvePlanningDocumentSeoImageURL(source)
	documentUpdates := map[string]interface{}{
		"slug":                resolvePlanningDocumentSeoSlug(source.Slug, source.Title, source.PlanningDocumentID),
		"origin_url":          documentURL,
		"canonical_url":       documentURL,
		"ref_type":            uint32(enums.ESEORefTypeDocument),
		"ref_id":              refID,
		"ref_source":          seo_domain.SeoRefSourceTqdDocument,
		"ref_label":           firstNonEmpty(source.Title, fmt.Sprintf("Tài liệu quy hoạch #%d", source.PlanningDocumentID)),
		"ref_url":             documentURL,
		"source_status":       seo_domain.SeoSourceStatusLinked,
		"ref_missing":         false,
		"ref_snapshot_json":   datatypes.JSON(metadataBytes),
		"ref_hash":            hashString(string(metadataBytes)),
		"ref_last_synced_at":  now,
		"scope":               seo_domain.SeoScopePublic,
		"title":               buildPlanningDocumentSeoTitle(source.Title, source.PlanningDocumentID),
		"description":         buildPlanningDocumentSeoDescription(source.Title, source.PlanningDocumentID),
		"published":           true,
		"published_at":        now,
		"is_site_map":         true,
		"is_index":            true,
		"site_map_lasted_at":  now,
		"sitemap_priority":    0.7,
		"sitemap_change_freq": seo_domain.SeoChangeFreqWeekly,
		"need_generate":       true,
		"source_updated_at":   now,
		"image_url":           imageURL,
		"metadata":            datatypes.JSON(metadataBytes),
	}
	if source.ProvinceID != nil && *source.ProvinceID > 0 {
		documentUpdates["province_id"] = *source.ProvinceID
	}

	if source.SeoID != nil && *source.SeoID > 0 {
		existing, err := u.seoDomainRepo.GetByID(ctx, *source.SeoID)
		if err != nil {
			return nil, false, err
		}
		if existing != nil {
			if err := u.updateSeoDomainPreserveCategory(ctx, existing.ID, documentUpdates); err != nil {
				return nil, false, err
			}
			if err := u.tqdProvider.UpdatePlanningDocumentSeoID(ctx, source.PlanningDocumentID, existing.ID); err != nil {
				return nil, false, err
			}
			updated, err := u.seoDomainRepo.GetByID(ctx, existing.ID)
			if err != nil {
				return nil, false, err
			}
			return updated, false, nil
		}
	}

	existingByURL, err := u.seoDomainRepo.GetByCanonicalURL(ctx, documentURL)
	if err != nil {
		return nil, false, err
	}
	if existingByURL != nil {
		if err := u.updateSeoDomainPreserveCategory(ctx, existingByURL.ID, documentUpdates); err != nil {
			return nil, false, err
		}
		if err := u.tqdProvider.UpdatePlanningDocumentSeoID(ctx, source.PlanningDocumentID, existingByURL.ID); err != nil {
			return nil, false, err
		}
		updated, err := u.seoDomainRepo.GetByID(ctx, existingByURL.ID)
		if err != nil {
			return nil, false, err
		}
		return updated, false, nil
	}

	seoDomain := &seo_domain.SeoDomain{
		Slug:              resolvePlanningDocumentSeoSlug(source.Slug, source.Title, source.PlanningDocumentID),
		OriginURL:         documentURL,
		CanonicalURL:      documentURL,
		RefType:           enums.ESEORefTypeDocument,
		RefID:             &refID,
		RefSource:         seo_domain.SeoRefSourceTqdDocument,
		RefLabel:          firstNonEmpty(source.Title, fmt.Sprintf("Tài liệu quy hoạch #%d", source.PlanningDocumentID)),
		RefURL:            documentURL,
		SourceStatus:      seo_domain.SeoSourceStatusLinked,
		RefMissing:        false,
		RefSnapshotJSON:   datatypes.JSON(metadataBytes),
		RefHash:           hashString(string(metadataBytes)),
		RefLastSyncedAt:   &now,
		Scope:             seo_domain.SeoScopePublic,
		Title:             buildPlanningDocumentSeoTitle(source.Title, source.PlanningDocumentID),
		Description:       buildPlanningDocumentSeoDescription(source.Title, source.PlanningDocumentID),
		Published:         true,
		PublishedAt:       &now,
		IsSiteMap:         true,
		IsIndex:           true,
		SiteMapLastedAt:   &now,
		SitemapPriority:   0.7,
		SitemapChangeFreq: seo_domain.SeoChangeFreqWeekly,
		NeedGenerate:      true,
		SourceUpdatedAt:   &now,
		ImageURL:          imageURL,
		Metadata:          datatypes.JSON(metadataBytes),
	}
	if source.ProvinceID != nil && *source.ProvinceID > 0 {
		seoDomain.ProvinceID = source.ProvinceID
	}

	var created *seo_domain.SeoDomain
	err = u.transaction.WithTransaction(ctx, func(txCtx context.Context) error {
		createdSeo, err := u.seoDomainRepo.Create(txCtx, seoDomain)
		if err != nil {
			return err
		}
		if createdSeo == nil || createdSeo.ID == 0 {
			return _errors.ReturnError(500, "tạo seo domain thất bại")
		}
		if err := u.tqdProvider.UpdatePlanningDocumentSeoID(txCtx, source.PlanningDocumentID, createdSeo.ID); err != nil {
			return err
		}
		created = createdSeo
		return nil
	})
	if err != nil {
		return nil, false, err
	}

	return created, true, nil
}

func administrativeUnitCompositeID(source dto.AdministrativeUnitSeoSource) (uint64, uint32, error) {
	compositeID := source.CompositeID
	if compositeID == 0 {
		compositeID = source.Identity
	}
	_, kind, err := _utils.DecodeLocationID(compositeID)
	if err != nil {
		return 0, 0, _errors.ReturnError(400, "compositeId không hợp lệ")
	}
	return compositeID, kind, nil
}

func administrativeUnitSeoRefType(kind uint32) enums.ESEORefType {
	return enums.ESEORefTypeAdmUnit
}

func (u *SeoDomainUsecase) createSeoFromAdministrativeUnitSource(ctx context.Context, source dto.AdministrativeUnitSeoSource) (*seo_domain.SeoDomain, bool, error) {
	if u.tqdProvider == nil {
		return nil, false, _errors.ReturnError(503, "TQD provider chưa được cấu hình")
	}

	compositeID, kind, err := administrativeUnitCompositeID(source)
	if err != nil {
		return nil, false, err
	}
	refID := compositeID
	refType := administrativeUnitSeoRefType(kind)
	identityLabel := strconv.FormatUint(compositeID, 10)

	slug, err := resolveAdministrativeUnitSeoSlug(source.Slug, source.Name, source.AdministrativeCode, kind)
	if err != nil {
		return nil, false, err
	}
	canonicalPath, err := buildAdministrativeUnitSeoURL(source.Slug, source.Name, source.AdministrativeCode, kind)
	if err != nil {
		return nil, false, err
	}

	metadataBytes, err := json.Marshal(map[string]interface{}{
		"identity":           compositeID,
		"compositeId":        compositeID,
		"unitType":           kind,
		"sourceId":           source.SourceID,
		"name":               source.Name,
		"slug":               slug,
		"administrativeCode": source.AdministrativeCode,
		"source":             seo_domain.SeoRefSourceTqdAdministrativeUnit,
		"sourceKey":          "administrative_unit",
	})
	if err != nil {
		return nil, false, err
	}

	now := time.Now()
	updates := map[string]interface{}{
		"slug":                slug,
		"origin_url":          canonicalPath,
		"canonical_url":       canonicalPath,
		"ref_type":            uint32(refType),
		"ref_id":              refID,
		"ref_source":          seo_domain.SeoRefSourceTqdAdministrativeUnit,
		"ref_label":           firstNonEmpty(source.Name, identityLabel),
		"ref_url":             identityLabel,
		"source_status":       seo_domain.SeoSourceStatusLinked,
		"ref_missing":         false,
		"ref_snapshot_json":   datatypes.JSON(metadataBytes),
		"ref_hash":            hashString(string(metadataBytes)),
		"ref_last_synced_at":  now,
		"scope":               seo_domain.SeoScopePublic,
		"title":               buildAdministrativeUnitSeoTitle(source.Name, identityLabel),
		"description":         buildAdministrativeUnitSeoDescription(source.Name, identityLabel),
		"published":           true,
		"published_at":        now,
		"is_site_map":         true,
		"is_index":            true,
		"site_map_lasted_at":  now,
		"sitemap_priority":    0.8,
		"sitemap_change_freq": seo_domain.SeoChangeFreqWeekly,
		"need_generate":       true,
		"source_updated_at":   now,
		"template_key":        "administrative-unit",
		"metadata":            datatypes.JSON(metadataBytes),
	}
	if source.ProvinceID != nil && *source.ProvinceID > 0 {
		updates["province_id"] = *source.ProvinceID
	}

	persistSeoID := func(ctx context.Context, seoID uint64) error {
		return u.tqdProvider.UpdateAdministrativeUnitSeoID(ctx, compositeID, seoID)
	}

	if source.SeoID != nil && *source.SeoID > 0 {
		existing, err := u.seoDomainRepo.GetByID(ctx, *source.SeoID)
		if err != nil {
			return nil, false, err
		}
		if existing != nil {
			if err := u.updateSeoDomainPreserveCategory(ctx, existing.ID, updates); err != nil {
				return nil, false, err
			}
			if err := persistSeoID(ctx, existing.ID); err != nil {
				return nil, false, err
			}
			updated, err := u.seoDomainRepo.GetByID(ctx, existing.ID)
			if err != nil {
				return nil, false, err
			}
			return updated, false, nil
		}
	}

	existingByRef, err := u.seoDomainRepo.GetByRefSource(ctx, uint32(refType), seo_domain.SeoRefSourceTqdAdministrativeUnit, refID)
	if err != nil {
		return nil, false, err
	}
	if existingByRef != nil {
		if err := u.updateSeoDomainPreserveCategory(ctx, existingByRef.ID, updates); err != nil {
			return nil, false, err
		}
		if err := persistSeoID(ctx, existingByRef.ID); err != nil {
			return nil, false, err
		}
		updated, err := u.seoDomainRepo.GetByID(ctx, existingByRef.ID)
		if err != nil {
			return nil, false, err
		}
		return updated, false, nil
	}

	existingByURL, err := u.seoDomainRepo.GetByCanonicalURL(ctx, canonicalPath)
	if err != nil {
		return nil, false, err
	}
	if existingByURL != nil {
		if err := u.updateSeoDomainPreserveCategory(ctx, existingByURL.ID, updates); err != nil {
			return nil, false, err
		}
		if err := persistSeoID(ctx, existingByURL.ID); err != nil {
			return nil, false, err
		}
		updated, err := u.seoDomainRepo.GetByID(ctx, existingByURL.ID)
		if err != nil {
			return nil, false, err
		}
		return updated, false, nil
	}

	seoDomain := &seo_domain.SeoDomain{
		Slug:              slug,
		OriginURL:         canonicalPath,
		CanonicalURL:      canonicalPath,
		RefType:           refType,
		RefID:             &refID,
		RefSource:         seo_domain.SeoRefSourceTqdAdministrativeUnit,
		RefLabel:          firstNonEmpty(source.Name, identityLabel),
		RefURL:            identityLabel,
		SourceStatus:      seo_domain.SeoSourceStatusLinked,
		RefMissing:        false,
		RefSnapshotJSON:   datatypes.JSON(metadataBytes),
		RefHash:           hashString(string(metadataBytes)),
		RefLastSyncedAt:   &now,
		Scope:             seo_domain.SeoScopePublic,
		Title:             buildAdministrativeUnitSeoTitle(source.Name, identityLabel),
		Description:       buildAdministrativeUnitSeoDescription(source.Name, identityLabel),
		Published:         true,
		PublishedAt:       &now,
		IsSiteMap:         true,
		IsIndex:           true,
		SiteMapLastedAt:   &now,
		SitemapPriority:   0.8,
		SitemapChangeFreq: seo_domain.SeoChangeFreqWeekly,
		NeedGenerate:      true,
		SourceUpdatedAt:   &now,
		TemplateKey:       "administrative-unit",
		Metadata:          datatypes.JSON(metadataBytes),
	}
	if source.ProvinceID != nil && *source.ProvinceID > 0 {
		seoDomain.ProvinceID = source.ProvinceID
	}

	var created *seo_domain.SeoDomain
	err = u.transaction.WithTransaction(ctx, func(txCtx context.Context) error {
		createdSeo, err := u.seoDomainRepo.Create(txCtx, seoDomain)
		if err != nil {
			return err
		}
		if createdSeo == nil || createdSeo.ID == 0 {
			return _errors.ReturnError(500, "tạo seo domain thất bại")
		}
		if err := persistSeoID(txCtx, createdSeo.ID); err != nil {
			return err
		}
		created = createdSeo
		return nil
	})
	if err != nil {
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

func (u *SeoDomainUsecase) BatchUpdateSeoDomain(ctx context.Context, req *dto.BatchUpdateSeoDomainRequest) error {
	if req == nil || len(req.IDs) == 0 {
		return _errors.ReturnError(400, "ids là bắt buộc")
	}
	if len(req.Fields) == 0 {
		return _errors.ReturnError(400, "fields là bắt buộc")
	}

	ids := make([]uint64, 0, len(req.IDs))
	seen := make(map[uint64]struct{}, len(req.IDs))
	for _, id := range req.IDs {
		if id == 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return _errors.ReturnError(400, "ids không hợp lệ")
	}

	updates := map[string]interface{}{}
	fieldSeen := make(map[string]struct{}, len(req.Fields))
	for _, raw := range req.Fields {
		field := strings.TrimSpace(raw)
		if field == "" {
			continue
		}
		key := strings.ToLower(field)
		if _, ok := fieldSeen[key]; ok {
			continue
		}
		fieldSeen[key] = struct{}{}

		switch key {
		case "issitemap", "is_site_map":
			if req.IsSiteMap == nil {
				return _errors.ReturnError(400, "isSiteMap là bắt buộc khi fields chứa isSiteMap")
			}
			updates["is_site_map"] = *req.IsSiteMap
		default:
			return _errors.ReturnError(400, fmt.Sprintf("field không được hỗ trợ: %s", field))
		}
	}
	if len(updates) == 0 {
		return _errors.ReturnError(400, "fields không hợp lệ")
	}

	return u.seoDomainRepo.UpdateByIDs(ctx, ids, updates)
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
	if !enums.IsValidSEORefType(req.RefType) || enums.IsManualSEORefType(req.RefType) {
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
	req.RefType = enums.CanonicalSEORefType(req.RefType)

	sourceBinding, err := u.prepareSourceBinding(ctx, req.RefType, req.RefID, strings.TrimSpace(req.RefSource))
	if err != nil {
		return nil, err
	}

	slug := strings.TrimSpace(req.Slug)
	if slug == "" {
		slug = firstNonEmpty(sourceBinding.SuggestedSlug, slugify(req.Title), "seo-page")
	}
	isManual := enums.IsManualSEORefType(req.RefType)
	if isManual {
		slug = extractManualTinTucSlugSegment(slug)
	}

	originURL := strings.TrimSpace(req.OriginURL)
	if originURL == "" && sourceBinding.RefURL != "" {
		originURL = sourceBinding.RefURL
	}
	if originURL == "" {
		if isManual {
			originURL = slug
		} else {
			originURL = "/seo/" + firstNonEmpty(strings.TrimSpace(req.Slug), sourceBinding.SuggestedSlug, slugify(req.Title), "seo-page")
		}
	} else if isManual {
		originURL = extractManualTinTucSlugSegment(originURL)
	}

	canonicalURL := originURL
	if req.CanonicalURL != nil && strings.TrimSpace(*req.CanonicalURL) != "" {
		canonicalURL = strings.TrimSpace(*req.CanonicalURL)
		if isManual {
			canonicalURL = extractManualTinTucSlugSegment(canonicalURL)
		}
	} else if isManual {
		canonicalURL = slug
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

	visibleOn, err := normalizeVisibleOn(req.VisibleOn)
	if err != nil {
		return nil, err
	}

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
	// Manual: canonical tạm là slug segment; path đầy đủ /tin-tuc/tin/{id}/{slug} gắn sau khi có ID.
	if exist && !isManual {
		return nil, _errors.ReturnError(409, "canonicalUrl đã tồn tại")
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
		SourceStatus:      normalizeSourceStatus(firstNonEmpty(strings.TrimSpace(req.SourceStatus), sourceBinding.SourceStatus), enums.IsManualSEORefType(req.RefType)),
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
		ImageURL:          strings.TrimSpace(req.ImageURL),
		Metadata:          metadata,
		Note:              req.Note,
		CategoryID:        req.CategoryID,
		VisibleOn:         visibleOn,
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

	manualRefType := refType
	if req.RefType != nil {
		manualRefType = enums.CanonicalSEORefType(*req.RefType)
	}
	isManual := enums.IsManualSEORefType(manualRefType)

	if req.Slug != nil {
		slug := strings.TrimSpace(*req.Slug)
		if slug == "" {
			return nil, _errors.ReturnError(400, "slug không được để trống")
		}
		if isManual {
			segment := extractManualTinTucSlugSegment(slug)
			path := buildManualTinTucCanonical(existing.ID, segment)
			slugField := buildManualTinTucSlugField(existing.ID, segment)
			state.Slug = slugField
			state.Canonical = path
			updates["slug"] = slugField
			updates["origin_url"] = path
			updates["canonical_url"] = path
		} else {
			state.Slug = slug
			updates["slug"] = slug
		}
	}

	if req.OriginURL != nil {
		originURL := strings.TrimSpace(*req.OriginURL)
		if originURL == "" {
			return nil, _errors.ReturnError(400, "originUrl không được để trống")
		}
		if isManual {
			// Nếu đã rebuild từ slug ở trên thì giữ path đó; không thì chuẩn hóa theo id.
			if _, ok := updates["origin_url"]; !ok {
				segment := extractManualTinTucSlugSegment(originURL)
				path := buildManualTinTucCanonical(existing.ID, segment)
				updates["origin_url"] = path
				if _, hasCanonical := updates["canonical_url"]; !hasCanonical {
					updates["canonical_url"] = path
					state.Canonical = path
				}
				if _, hasSlug := updates["slug"]; !hasSlug {
					updates["slug"] = buildManualTinTucSlugField(existing.ID, segment)
					state.Slug = updates["slug"].(string)
				}
			}
		} else {
			updates["origin_url"] = originURL
		}
	}

	if req.CanonicalURL != nil {
		canonicalURL := strings.TrimSpace(*req.CanonicalURL)
		if canonicalURL == "" {
			canonicalURL = existing.OriginURL
		}
		if isManual {
			if _, ok := updates["canonical_url"]; !ok {
				segment := extractManualTinTucSlugSegment(canonicalURL)
				path := buildManualTinTucCanonical(existing.ID, segment)
				canonicalURL = path
				updates["canonical_url"] = path
				if _, hasOrigin := updates["origin_url"]; !hasOrigin {
					updates["origin_url"] = path
				}
				if _, hasSlug := updates["slug"]; !hasSlug {
					updates["slug"] = buildManualTinTucSlugField(existing.ID, segment)
					state.Slug = updates["slug"].(string)
				}
			} else {
				canonicalURL = updates["canonical_url"].(string)
			}
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
		refType = enums.CanonicalSEORefType(*req.RefType)
		updates["ref_type"] = enums.ESEORefType(refType)
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
		status := normalizeSourceStatus(*req.SourceStatus, enums.IsManualSEORefType(refType))
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

	if req.CategoryID != nil {
		if *req.CategoryID == 0 {
			updates["category_id"] = nil
		} else {
			updates["category_id"] = *req.CategoryID
		}
	}
	if req.VisibleOn != nil {
		visibleOn, err := normalizeVisibleOn(req.VisibleOn)
		if err != nil {
			return nil, err
		}
		updates["visible_on"] = visibleOn
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
	if req.ImageURL != nil {
		updates["image_url"] = strings.TrimSpace(*req.ImageURL)
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

	if enums.IsManualSEORefType(refType) {
		return seoSourceBinding{
			RefSource:       seo_domain.SeoRefSourceManual,
			SourceStatus:    seo_domain.SeoSourceStatusManual,
			RefMissing:      false,
			RefSnapshotJSON: emptyJSON,
		}, nil
	}

	if refID == nil || *refID == 0 {
		return seoSourceBinding{}, _errors.ReturnError(400, "refId là bắt buộc khi refType khác 800")
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
	if enums.IsManualSEORefType(refType) {
		return nil
	}
	if refID == nil || *refID == 0 {
		return _errors.ReturnError(400, "refId là bắt buộc khi refType khác 800")
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
	case seo_domain.SeoSourceStatusManual,
		"Thủ công": // legacy label — map về giá trị DB
		return seo_domain.SeoSourceStatusManual
	case seo_domain.SeoSourceStatusLinked,
		seo_domain.SeoSourceStatusStale,
		seo_domain.SeoSourceStatusMissing,
		seo_domain.SeoSourceStatusPermissionDenied:
		return strings.TrimSpace(status)
	default:
		return seo_domain.SeoSourceStatusLinked
	}
}

func normalizeVisibleOn(value *uint32) (uint32, error) {
	if value == nil || *value == 0 {
		return seo_domain.SeoVisibleOnLatest, nil
	}
	switch *value {
	case seo_domain.SeoVisibleOnFeatured,
		seo_domain.SeoVisibleOnAttention,
		seo_domain.SeoVisibleOnLatest,
		seo_domain.SeoVisibleOnMostRead,
		seo_domain.SeoVisibleOnBreaking:
		return *value, nil
	default:
		return 0, _errors.ReturnError(400, "visibleOn không hợp lệ (cho phép 10,20,30,40,50)")
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

// removeVietnameseDiacritics chuyển chuỗi Unicode sang ASCII bằng NFD decomposition rồi loại bỏ combining marks.
func removeVietnameseDiacritics(s string) string {
	s = strings.NewReplacer("đ", "d", "Đ", "D").Replace(s)
	t := transform.Chain(norm.NFD, transform.RemoveFunc(func(r rune) bool {
		return unicode.Is(unicode.Mn, r)
	}), norm.NFC)
	result, _, err := transform.String(t, s)
	if err != nil {
		return s
	}
	return result
}

func slugify(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	value = removeVietnameseDiacritics(value)
	value = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '/' || r == '_' || r == '-' {
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

// extractManualTinTucSlugSegment lấy segment slug cuối từ input (slug thuần hoặc path đầy đủ).
// Ví dụ: "bai-viet-a", "/tin-tuc/tin/8/bai-viet-a", "/tin-tuc/bai-viet/bai-viet-a" → "bai-viet-a".
func extractManualTinTucSlugSegment(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	value = strings.Trim(value, "/")
	if value == "" {
		return "seo-page"
	}
	parts := strings.Split(value, "/")
	segment := parts[len(parts)-1]
	segment = slugify(segment)
	segment = strings.ReplaceAll(segment, "/", "-")
	segment = strings.Trim(segment, "-")
	if segment == "" {
		return "seo-page"
	}
	return segment
}

func buildManualTinTucCanonical(id uint64, slugSegment string) string {
	return buildSeoURLFromTemplate(manualTinTucURLTemplate(), id, slugSegment)
}

func buildManualTinTucSlugField(id uint64, slugSegment string) string {
	return fmt.Sprintf("%d/%s", id, slugSegment)
}

func manualTinTucURLTemplate() string {
	if template := strings.TrimSpace(viper.GetString("seo.manualTinTucUrlTemplate")); template != "" {
		return template
	}
	return "/tin-tuc/tin/{id}/{slug}"
}

// finalizeManualTinTucPaths gắn canonical/origin/slug theo /tin-tuc/tin/{id}/{slug} sau khi có ID.
func (u *SeoDomainUsecase) finalizeManualTinTucPaths(ctx context.Context, seo *seo_domain.SeoDomain) (*seo_domain.SeoDomain, error) {
	if seo == nil || seo.ID == 0 || !enums.IsManualSEORefType(uint32(seo.RefType)) {
		return seo, nil
	}
	segment := extractManualTinTucSlugSegment(seo.Slug)
	if segment == "" || segment == "seo-page" {
		if alt := extractManualTinTucSlugSegment(seo.CanonicalURL); alt != "" {
			segment = alt
		}
	}
	path := buildManualTinTucCanonical(seo.ID, segment)
	slugField := buildManualTinTucSlugField(seo.ID, segment)
	if seo.Slug == slugField && seo.CanonicalURL == path && seo.OriginURL == path {
		return seo, nil
	}
	if err := u.seoDomainRepo.Update(ctx, seo.ID, map[string]interface{}{
		"slug":          slugField,
		"origin_url":    path,
		"canonical_url": path,
	}); err != nil {
		return nil, err
	}
	return u.seoDomainRepo.GetByID(ctx, seo.ID)
}

func seoSlugSegment(label string, fallbackID uint64) string {
	slug := slugify(label)
	slug = strings.ReplaceAll(slug, "/", "-")
	slug = strings.Trim(slug, "-")
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}
	if slug == "" || slug == "seo-page" {
		return fmt.Sprintf("%d", fallbackID)
	}
	return slug
}

func buildSeoURLFromTemplate(template string, id uint64, slug string) string {
	return strings.NewReplacer(
		"{id}", fmt.Sprintf("%d", id),
		"{slug}", slug,
	).Replace(template)
}

func parcelSeoURLTemplate() string {
	if template := strings.TrimSpace(viper.GetString("seo.parcelUrlTemplate")); template != "" {
		return template
	}
	return "/tin-tuc/thua-dat/{id}/{slug}"
}

func layerSeoURLTemplate() string {
	if template := strings.TrimSpace(viper.GetString("seo.layerUrlTemplate")); template != "" {
		return template
	}
	return "/tin-tuc/lop/{id}/{slug}"
}

func buildParcelSeoURL(parcelID uint64, adrSearch string) (string, error) {
	return buildSeoURLFromTemplate(
		parcelSeoURLTemplate(),
		parcelID,
		seoSlugSegment(adrSearch, parcelID),
	), nil
}

func resolveParcelSeoSlug(adrSearch string, parcelID uint64) string {
	return fmt.Sprintf("%d/%s", parcelID, seoSlugSegment(adrSearch, parcelID))
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

func buildLayerSeoURL(layerID uint64, displayName string) (string, error) {
	return buildSeoURLFromTemplate(
		layerSeoURLTemplate(),
		layerID,
		seoSlugSegment(displayName, layerID),
	), nil
}

func resolveLayerSeoSlug(displayName string, layerID uint64) string {
	return fmt.Sprintf("%d/%s", layerID, seoSlugSegment(displayName, layerID))
}

func buildLayerSeoTitle(displayName string, layerID uint64) string {
	title := strings.TrimSpace(displayName)
	if title == "" {
		return fmt.Sprintf("Lớp quy hoạch %d", layerID)
	}
	return title
}

func buildLayerSeoDescription(displayName string, layerID uint64) string {
	title := buildLayerSeoTitle(displayName, layerID)
	return fmt.Sprintf("Thông tin lớp quy hoạch %s", title)
}

func resolveLayerSeoImageURL(source dto.LayerSeoSource) string {
	return firstNonEmpty(source.ThumbnailURL, source.ImageURL, source.Avatar)
}

func isPlanningDocumentImageExt(ext string) bool {
	switch strings.ToLower(strings.TrimPrefix(strings.TrimSpace(ext), ".")) {
	case "jpg", "jpeg", "png", "gif", "webp", "bmp", "svg", "heic", "tif", "tiff":
		return true
	default:
		return false
	}
}

func resolvePlanningDocumentSeoImageURL(source dto.PlanningDocumentSeoSource) string {
	if !isPlanningDocumentImageExt(source.Ext) {
		return ""
	}
	return firstNonEmpty(source.Filepath, source.Thumbnail)
}

func planningProjectSeoURLTemplate() string {
	if template := strings.TrimSpace(viper.GetString("seo.planningProjectUrlTemplate")); template != "" {
		return template
	}
	return "/tin-tuc/do-an-quy-hoach/{id}/{slug}"
}

func buildPlanningProjectSeoURL(planningProjectID uint64, slug, name string) (string, error) {
	template := planningProjectSeoURLTemplate()
	resolvedSlug := planningProjectNameSlug(slug, name, planningProjectID)

	return strings.NewReplacer(
		"{id}", fmt.Sprintf("%d", planningProjectID),
		"{slug}", resolvedSlug,
	).Replace(template), nil
}

func planningProjectNameSlug(slug, name string, planningProjectID uint64) string {
	label := strings.Trim(strings.TrimSpace(slug), "/")
	if label == "" {
		label = strings.TrimSpace(name)
	}
	// Avoid embedding id/ again if caller already passed compound slug.
	if parts := strings.Split(label, "/"); len(parts) == 2 && parts[0] == fmt.Sprintf("%d", planningProjectID) {
		label = parts[1]
	}
	return seoSlugSegment(label, planningProjectID)
}

func resolvePlanningProjectSeoSlug(slug, name string, planningProjectID uint64) string {
	return fmt.Sprintf("%d/%s", planningProjectID, planningProjectNameSlug(slug, name, planningProjectID))
}

func buildPlanningProjectSeoTitle(name string, planningProjectID uint64) string {
	title := strings.TrimSpace(name)
	if title == "" {
		return fmt.Sprintf("Đồ án quy hoạch %d", planningProjectID)
	}
	return title
}

func buildPlanningProjectSeoDescription(name string, planningProjectID uint64) string {
	title := buildPlanningProjectSeoTitle(name, planningProjectID)
	return fmt.Sprintf("Thông tin đồ án quy hoạch %s", title)
}

func planningDocumentSeoURLTemplate() string {
	if template := strings.TrimSpace(viper.GetString("seo.planningDocumentUrlTemplate")); template != "" {
		return template
	}
	return "/tin-tuc/tai-lieu/{id}/{slug}"
}

func buildPlanningDocumentSeoURL(documentID uint64, slug, title string) (string, error) {
	template := planningDocumentSeoURLTemplate()
	resolvedSlug := planningDocumentNameSlug(slug, title, documentID)
	return strings.NewReplacer(
		"{id}", fmt.Sprintf("%d", documentID),
		"{slug}", resolvedSlug,
	).Replace(template), nil
}

func planningDocumentNameSlug(slug, title string, documentID uint64) string {
	label := strings.Trim(strings.TrimSpace(slug), "/")
	if label == "" {
		label = strings.TrimSpace(title)
	}
	if parts := strings.Split(label, "/"); len(parts) == 2 && parts[0] == fmt.Sprintf("%d", documentID) {
		label = parts[1]
	}
	return seoSlugSegment(label, documentID)
}

func resolvePlanningDocumentSeoSlug(slug, title string, documentID uint64) string {
	return fmt.Sprintf("%d/%s", documentID, planningDocumentNameSlug(slug, title, documentID))
}

func buildPlanningDocumentSeoTitle(title string, documentID uint64) string {
	name := strings.TrimSpace(title)
	if name == "" {
		return fmt.Sprintf("Tài liệu quy hoạch %d", documentID)
	}
	return name
}

func buildPlanningDocumentSeoDescription(title string, documentID uint64) string {
	name := buildPlanningDocumentSeoTitle(title, documentID)
	return fmt.Sprintf("Thông tin tài liệu quy hoạch %s", name)
}

func administrativeUnitSeoURLTemplate() string {
	if template := strings.TrimSpace(viper.GetString("seo.administrativeUnitUrlTemplate")); template != "" {
		return template
	}
	return "/tin-tuc/dia-ban/{id}/{slug}"
}

func administrativeUnitNameSlug(slug, name string) string {
	resolved := strings.Trim(strings.TrimSpace(slug), "/")
	if resolved != "" {
		if parts := strings.Split(resolved, "/"); len(parts) >= 2 {
			resolved = parts[len(parts)-1]
		}
	}
	if resolved == "" {
		resolved = seoSlugSegment(strings.TrimSpace(name), 0)
	} else {
		resolved = seoSlugSegment(resolved, 0)
	}
	return resolved
}

func resolveAdministrativeUnitSeoSlug(slug, name, administrativeCode string, kind uint32) (string, error) {
	if kind != _utils.LocationIDKindProvince && kind != _utils.LocationIDKindWard {
		return "", _errors.ReturnError(400, "unitType địa bàn không hợp lệ")
	}
	code := strings.Trim(slugify(strings.TrimSpace(administrativeCode)), "-")
	nameSlug := administrativeUnitNameSlug(slug, name)
	if code == "" || nameSlug == "" || nameSlug == "seo-page" || nameSlug == "0" {
		return "", _errors.ReturnError(400, "không tạo được slug cho đơn vị hành chính")
	}
	// id = {administrativeCode}{1|2} theo quy tắc compositeId
	return code + strconv.FormatUint(uint64(kind), 10) + "/" + nameSlug, nil
}

func buildAdministrativeUnitSeoURL(slug, name, administrativeCode string, kind uint32) (string, error) {
	compound, err := resolveAdministrativeUnitSeoSlug(slug, name, administrativeCode, kind)
	if err != nil {
		return "", err
	}
	parts := strings.SplitN(compound, "/", 2)
	return strings.NewReplacer(
		"{id}", parts[0],
		"{slug}", parts[1],
	).Replace(administrativeUnitSeoURLTemplate()), nil
}

func buildAdministrativeUnitSeoTitle(name, identity string) string {
	title := strings.TrimSpace(name)
	if title == "" {
		return firstNonEmpty(identity, "Đơn vị hành chính")
	}
	return title
}

func buildAdministrativeUnitSeoDescription(name, identity string) string {
	title := buildAdministrativeUnitSeoTitle(name, identity)
	return fmt.Sprintf("Thông tin địa bàn %s", title)
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
