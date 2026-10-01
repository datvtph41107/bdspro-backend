// Package repository provides GORM-backed persistence for modules/seo.
package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"userv4/modules/seo/dto"
	"userv4/modules/seo/entity"

	"gorm.io/gorm"
)

// ── NewsRepository ────────────────────────────────────────────────────────────

type NewsRepository interface {
	ListPlanningNews(ctx context.Context, filter dto.PlanningNewsListFilter) (*dto.PlanningNewsListResult, error)
	GetPlanningNews(ctx context.Context, id uint64) (*dto.PlanningNewsDetailResult, error)
	ListCategories(ctx context.Context) ([]dto.ContentCategory, error)
	GetCategory(ctx context.Context, id uint64) (*dto.ContentCategory, error)
	CreateCategory(ctx context.Context, req dto.CreateCategoryRequest) (*dto.ContentCategory, error)
	UpdateCategory(ctx context.Context, id uint64, req dto.UpdateCategoryRequest) (*dto.ContentCategory, error)
	DeleteCategory(ctx context.Context, id uint64) error
	ListSavedPlanningNews(ctx context.Context, userID uint64, filter dto.SavedNewsFilter) (*dto.PlanningNewsListResult, error)
	SavePlanningNews(ctx context.Context, userID, newsID uint64) (*dto.PlanningNewsItem, error)
	UnsavePlanningNews(ctx context.Context, userID, newsID uint64) error
	RecordClick(ctx context.Context, seoDomainID uint64, userID *uint64, ipAddr string) error
	SearchContent(ctx context.Context, query string, kinds []string, limit int) (*dto.ContentSearchResult, error)
}

type newsPostgres struct{ db *gorm.DB }

func NewNewsPostgres(db *gorm.DB) NewsRepository { return &newsPostgres{db: db} }

// ── Shared SQL fragments ──────────────────────────────────────────────────────

const planningNewsWhere = `
    sd.deleted_at IS NULL
    AND sd.published = TRUE
    AND sd.scope = 'public'
    AND (
        sd.canonical_url LIKE '/tin-tuc/%'
        OR sd.canonical_url LIKE '/quy-hoach/tin-tuc/%'
        OR sd.template_key IN ('planning-news', 'planning-news-detail')
        OR c.code = 'planning_news'
    )`

const planningNewsCategoryJoin = `
      LEFT JOIN seo_category c ON c.id = sd.category_id AND c.deleted_at IS NULL`

const planningNewsSelectBrief = `
    sd.id,
    sd.id AS seo_page_id,
    sd.slug AS slug,
    COALESCE(sd.title, '') AS title,
    COALESCE(sd.description, '') AS description,
    COALESCE(sd.summary, '') AS summary,
    '' AS "content",
    '' AS rendered_html,
    sd.canonical_url AS canonical_url,
    CASE
        WHEN sd.canonical_url LIKE 'http%' THEN sd.canonical_url
        ELSE 'https://qhpro.vn' || sd.canonical_url
    END AS public_url,
    COALESCE(sd.deep_link, '') AS deep_link,
    COALESCE(sd.image_url, '') AS image_url,
    COALESCE(sd.ref_source, '') AS ref_source,
    COALESCE(sd.ref_label, '') AS ref_label,
    COALESCE(sd.source_status, '') AS source_status,
    COALESCE(sd.visible_on, 30) AS visible_on,
    sd.published_at AS published_at,
    sd.updated_at AS updated_at,
    sd.category_id AS category_id,
    COALESCE(c.code, '') AS category_code,
    COALESCE(c.name, '') AS category_name,
    COALESCE(sd.metadata, '{}'::jsonb) AS metadata_json,
    COALESCE(sd.ref_snapshot_json, '{}'::jsonb) AS ref_snapshot_raw,
    COALESCE(sd.ref_type, 0) AS ref_type,
    sd.ref_id AS ref_id,
    sd.province_id AS province_id`

const planningNewsSelectFull = `
    sd.id,
    sd.id AS seo_page_id,
    sd.slug AS slug,
    COALESCE(sd.title, '') AS title,
    COALESCE(sd.description, '') AS description,
    COALESCE(sd.summary, '') AS summary,
    COALESCE(sd.content, '') AS "content",
    COALESCE(sd.rendered_html, '') AS rendered_html,
    sd.canonical_url AS canonical_url,
    CASE
        WHEN sd.canonical_url LIKE 'http%' THEN sd.canonical_url
        ELSE 'https://qhpro.vn' || sd.canonical_url
    END AS public_url,
    COALESCE(sd.deep_link, '') AS deep_link,
    COALESCE(sd.image_url, '') AS image_url,
    COALESCE(sd.ref_source, '') AS ref_source,
    COALESCE(sd.ref_label, '') AS ref_label,
    COALESCE(sd.source_status, '') AS source_status,
    COALESCE(sd.visible_on, 30) AS visible_on,
    sd.published_at AS published_at,
    sd.updated_at AS updated_at,
    sd.category_id AS category_id,
    COALESCE(c.code, '') AS category_code,
    COALESCE(c.name, '') AS category_name,
    COALESCE(sd.metadata, '{}'::jsonb) AS metadata_json,
    COALESCE(sd.ref_snapshot_json, '{}'::jsonb) AS ref_snapshot_raw,
    COALESCE(sd.ref_type, 0) AS ref_type,
    sd.ref_id AS ref_id,
    sd.province_id AS province_id`

type newsRow struct {
	ID             uint64
	SeoPageID      uint64
	Slug           string
	Title          string
	Description    string
	Summary        string
	Content        string
	RenderedHTML   string
	CanonicalURL   string
	PublicURL      string
	DeepLink       string
	ImageURL       string
	RefSource      string
	RefLabel       string
	SourceStatus   string
	VisibleOn      uint32
	RefType        uint32
	RefID          *uint64
	PublishedAt    *time.Time
	UpdatedAt      *time.Time
	SavedAt        *time.Time
	CategoryID     *uint64
	CategoryCode   string
	CategoryName   string
	ProvinceID     *uint64
	MetadataJSON   []byte
	RefSnapshotRaw []byte
}

func decodeJSONMap(raw []byte) map[string]any {
	out := map[string]any{}
	if len(raw) > 0 {
		_ = json.Unmarshal(raw, &out)
	}
	return out
}

func newsItemFromRow(row newsRow) dto.PlanningNewsItem {
	item := dto.PlanningNewsItem{
		ID:           row.ID,
		SEOPageID:    row.ID,
		Slug:         row.Slug,
		Title:        row.Title,
		Description:  row.Description,
		Summary:      row.Summary,
		Content:      row.Content,
		RenderedHTML: row.RenderedHTML,
		CanonicalURL: row.CanonicalURL,
		PublicURL:    row.PublicURL,
		DeepLink:     row.DeepLink,
		ImageURL:     row.ImageURL,
		RefSource:    row.RefSource,
		RefLabel:     row.RefLabel,
		SourceStatus: row.SourceStatus,
		VisibleOn:    row.VisibleOn,
		RefType:      row.RefType,
		RefTypeName:  entity.SEORefType(row.RefType).Name(),
		RefID:        row.RefID,
		PublishedAt:  row.PublishedAt,
		UpdatedAt:    row.UpdatedAt,
		SavedAt:      row.SavedAt,
		CategoryID:   row.CategoryID,
		CategoryCode: row.CategoryCode,
		CategoryName: row.CategoryName,
		ProvinceID:   row.ProvinceID,
		Metadata:     decodeJSONMap(row.MetadataJSON),
	}
	if row.CategoryID != nil && *row.CategoryID > 0 {
		item.Category = &dto.PlanningNewsCategory{ID: *row.CategoryID, Name: row.CategoryName}
	}
	return item
}

func normalizePage(page, size int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 20
	}
	if size > 100 {
		size = 100
	}
	return page, size
}

func likeKeyword(keyword string) string {
	keyword = strings.ToLower(strings.TrimSpace(strings.ReplaceAll(keyword, "+", " ")))
	parts := strings.Fields(keyword)
	if len(parts) == 0 {
		return ""
	}
	return "%" + strings.Join(parts, "%") + "%"
}

func newsOrderBy(filter dto.PlanningNewsListFilter) string {
	token := strings.ToLower(strings.TrimSpace(filter.Sort))
	if token == "" {
		by := strings.ToLower(strings.TrimSpace(filter.SortBy))
		order := strings.ToLower(strings.TrimSpace(filter.SortOrder))
		switch {
		case strings.HasSuffix(by, "_desc") || strings.HasSuffix(by, "_asc"):
			token = by
		case by != "" && order != "":
			token = by + "_" + order
		case by != "":
			token = by + "_desc"
		}
	}
	switch strings.ReplaceAll(token, "-", "_") {
	case "updated_desc", "updated_at_desc":
		return "sd.updated_at DESC NULLS LAST, sd.id DESC"
	case "updated_asc", "updated_at_asc":
		return "sd.updated_at ASC NULLS LAST, sd.id ASC"
	case "published_asc", "published_at_asc":
		return "sd.published_at ASC NULLS LAST, sd.updated_at ASC, sd.id ASC"
	default:
		return "sd.published_at DESC NULLS LAST, sd.updated_at DESC, sd.id DESC"
	}
}

func (r *newsPostgres) ListPlanningNews(ctx context.Context, filter dto.PlanningNewsListFilter) (*dto.PlanningNewsListResult, error) {
	filter.Page, filter.Size = normalizePage(filter.Page, filter.Size)

	base := ` FROM seo_domain sd` + planningNewsCategoryJoin + ` WHERE ` + planningNewsWhere
	args := []any{}

	if kw := likeKeyword(filter.Keyword); kw != "" {
		base += ` AND (LOWER(COALESCE(sd.title,'')) LIKE ? OR LOWER(COALESCE(sd.description,'')) LIKE ? OR LOWER(COALESCE(sd.slug,'')) LIKE ?)`
		args = append(args, kw, kw, kw)
	}
	if filter.CategoryID > 0 {
		base += ` AND sd.category_id = ?`
		args = append(args, filter.CategoryID)
	}
	if filter.ProvinceID > 0 {
		base += ` AND sd.province_id = ?`
		args = append(args, filter.ProvinceID)
	}
	if filter.RefType > 0 {
		base += ` AND sd.ref_type = ?`
		args = append(args, filter.RefType)
	}
	if filter.VisibleOn > 0 {
		base += ` AND sd.visible_on = ?`
		args = append(args, filter.VisibleOn)
	}
	if filter.LastDates > 0 {
		days := filter.LastDates
		if days > 365 {
			days = 365
		}
		since := time.Now().AddDate(0, 0, -int(days))
		base += ` AND sd.updated_at >= ?`
		args = append(args, since)
	}

	var total int64
	if err := r.db.WithContext(ctx).Raw("SELECT COUNT(DISTINCT sd.id)"+base, args...).Scan(&total).Error; err != nil {
		return nil, fmt.Errorf("count planning news: %w", err)
	}
	rows := []newsRow{}
	selectSQL := `SELECT ` + planningNewsSelectBrief + base + ` ORDER BY ` + newsOrderBy(filter) + ` LIMIT ? OFFSET ?`
	listArgs := append(append([]any{}, args...), filter.Size, (filter.Page-1)*filter.Size)
	if err := r.db.WithContext(ctx).Raw(selectSQL, listArgs...).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("list planning news: %w", err)
	}
	data := make([]dto.PlanningNewsItem, 0, len(rows))
	for _, row := range rows {
		data = append(data, newsItemFromRow(row))
	}
	return &dto.PlanningNewsListResult{Data: data, Total: total, Page: filter.Page, Size: filter.Size}, nil
}

func (r *newsPostgres) GetPlanningNews(ctx context.Context, id uint64) (*dto.PlanningNewsDetailResult, error) {
	if id == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	row := newsRow{}
	query := `SELECT ` + planningNewsSelectFull + `
      FROM seo_domain sd` + planningNewsCategoryJoin + `
      WHERE ` + planningNewsWhere + ` AND sd.id = ? LIMIT 1`
	if err := r.db.WithContext(ctx).Raw(query, id).Scan(&row).Error; err != nil {
		return nil, fmt.Errorf("get planning news: %w", err)
	}
	if row.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	relatedRows := []newsRow{}
	catID := uint64(0)
	if row.CategoryID != nil {
		catID = *row.CategoryID
	}
	relQuery := `SELECT ` + planningNewsSelectBrief + `
      FROM seo_domain sd` + planningNewsCategoryJoin + `
      WHERE ` + planningNewsWhere + ` AND sd.id <> ? AND (? = 0 OR sd.category_id = ?)
      ORDER BY sd.published_at DESC NULLS LAST LIMIT 8`
	if err := r.db.WithContext(ctx).Raw(relQuery, row.ID, catID, catID).Scan(&relatedRows).Error; err != nil {
		return nil, fmt.Errorf("list related planning news: %w", err)
	}
	related := make([]dto.PlanningNewsItem, 0, len(relatedRows))
	for _, rr := range relatedRows {
		related = append(related, newsItemFromRow(rr))
	}
	return &dto.PlanningNewsDetailResult{SEODomain: newsItemFromRow(row), RelatedNews: related}, nil
}

func (r *newsPostgres) ListCategories(ctx context.Context) ([]dto.ContentCategory, error) {
	items := []dto.ContentCategory{}
	err := r.db.WithContext(ctx).Raw(`
		SELECT c.id, c.code, c.slug, c.name,
		       COALESCE(c.description,'') AS description, c.status, c.sort_order,
		       COUNT(DISTINCT CASE WHEN sd.deleted_at IS NULL AND sd.published = TRUE AND sd.scope='public' THEN sd.id END) AS item_count
		FROM seo_category c
		LEFT JOIN seo_domain sd ON sd.category_id = c.id
		WHERE c.deleted_at IS NULL AND c.status = 'active'
		GROUP BY c.id, c.code, c.slug, c.name, c.description, c.status, c.sort_order
		ORDER BY c.sort_order, c.name, c.id`).Scan(&items).Error
	return items, err
}

func (r *newsPostgres) GetCategory(ctx context.Context, id uint64) (*dto.ContentCategory, error) {
	var item dto.ContentCategory
	err := r.db.WithContext(ctx).Raw(`
		SELECT c.id, c.code, c.slug, c.name,
		       COALESCE(c.description,'') AS description, c.status, c.sort_order,
		       COUNT(DISTINCT CASE WHEN sd.deleted_at IS NULL AND sd.published = TRUE AND sd.scope='public' THEN sd.id END) AS item_count
		FROM seo_category c
		LEFT JOIN seo_domain sd ON sd.category_id = c.id
		WHERE c.deleted_at IS NULL AND c.id = ?
		GROUP BY c.id, c.code, c.slug, c.name, c.description, c.status, c.sort_order
		LIMIT 1`, id).Scan(&item).Error
	if err != nil {
		return nil, err
	}
	if item.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &item, nil
}

func (r *newsPostgres) CreateCategory(ctx context.Context, req dto.CreateCategoryRequest) (*dto.ContentCategory, error) {
	cat := &entity.SeoCategory{
		Code:        req.Code,
		Slug:        req.Slug,
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
		SortOrder:   req.SortOrder,
	}
	if cat.Status == "" {
		cat.Status = "active"
	}
	if err := r.db.WithContext(ctx).Create(cat).Error; err != nil {
		return nil, fmt.Errorf("create seo category: %w", err)
	}
	return r.GetCategory(ctx, cat.ID)
}

func (r *newsPostgres) UpdateCategory(ctx context.Context, id uint64, req dto.UpdateCategoryRequest) (*dto.ContentCategory, error) {
	updates := map[string]any{"updated_at": time.Now()}
	if req.Code != nil {
		updates["code"] = *req.Code
	}
	if req.Slug != nil {
		updates["slug"] = *req.Slug
	}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}
	if err := r.db.WithContext(ctx).Model(&entity.SeoCategory{}).Where("id = ? AND deleted_at IS NULL", id).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("update seo category: %w", err)
	}
	return r.GetCategory(ctx, id)
}

func (r *newsPostgres) DeleteCategory(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&entity.SeoCategory{}, id).Error
}

func (r *newsPostgres) ListSavedPlanningNews(ctx context.Context, userID uint64, filter dto.SavedNewsFilter) (*dto.PlanningNewsListResult, error) {
	filter.Page, filter.Size = normalizePage(filter.Page, filter.Size)
	base := ` FROM user_saved_planning_news saved
	      JOIN seo_domain sd ON sd.id = saved.seo_domain_id` + planningNewsCategoryJoin + `
	      WHERE saved.user_id = ? AND saved.deleted_at IS NULL AND ` + planningNewsWhere
	var total int64
	if err := r.db.WithContext(ctx).Raw("SELECT COUNT(DISTINCT sd.id)"+base, userID).Scan(&total).Error; err != nil {
		return nil, fmt.Errorf("count saved planning news: %w", err)
	}
	rows := []newsRow{}
	selectSQL := `SELECT ` + planningNewsSelectBrief + `, saved.updated_at AS saved_at` + base + ` ORDER BY saved.updated_at DESC, saved.id DESC LIMIT ? OFFSET ?`
	if err := r.db.WithContext(ctx).Raw(selectSQL, userID, filter.Size, (filter.Page-1)*filter.Size).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("list saved planning news: %w", err)
	}
	data := make([]dto.PlanningNewsItem, 0, len(rows))
	for _, row := range rows {
		data = append(data, newsItemFromRow(row))
	}
	return &dto.PlanningNewsListResult{Data: data, Total: total, Page: filter.Page, Size: filter.Size}, nil
}

func (r *newsPostgres) SavePlanningNews(ctx context.Context, userID, newsID uint64) (*dto.PlanningNewsItem, error) {
	result := r.db.WithContext(ctx).Exec(`
		INSERT INTO user_saved_planning_news (user_id, seo_domain_id, created_at, updated_at, deleted_at)
		SELECT ?, sd.id, NOW(), NOW(), NULL
		FROM seo_domain sd`+planningNewsCategoryJoin+`
		WHERE `+planningNewsWhere+` AND sd.id = ?
		ON CONFLICT (user_id, seo_domain_id)
		DO UPDATE SET updated_at = NOW(), deleted_at = NULL`, userID, newsID)
	if result.Error != nil {
		return nil, fmt.Errorf("save planning news: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	detail, err := r.GetPlanningNews(ctx, newsID)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	item := detail.SEODomain
	item.SavedAt = &now
	return &item, nil
}

func (r *newsPostgres) UnsavePlanningNews(ctx context.Context, userID, newsID uint64) error {
	return r.db.WithContext(ctx).Table("user_saved_planning_news").
		Where("user_id = ? AND seo_domain_id = ? AND deleted_at IS NULL", userID, newsID).
		Updates(map[string]any{"deleted_at": time.Now(), "updated_at": time.Now()}).Error
}

func (r *newsPostgres) RecordClick(ctx context.Context, seoDomainID uint64, userID *uint64, ipAddr string) error {
	log := &entity.SeoCountLog{
		SeoDomainID: seoDomainID,
		EventType:   "click",
		IPAddr:      ipAddr,
		UserID:      userID,
	}
	if err := r.db.WithContext(ctx).Create(log).Error; err != nil {
		return err
	}
	return r.db.WithContext(ctx).Exec(`
		INSERT INTO seo_count (seo_domain_id, click_count, view_count, updated_at)
		VALUES (?, 1, 0, NOW())
		ON CONFLICT (seo_domain_id) DO UPDATE SET click_count = seo_count.click_count + 1, updated_at = NOW()
	`, seoDomainID).Error
}

type contentSearchRow struct {
	ID           uint64
	Slug         string
	Title        string
	Description  string
	Summary      string
	CanonicalURL string
	DeepLink     string
	ImageURL     string
	PublishedAt  *time.Time
	UpdatedAt    *time.Time
	CategoryCode string
	CategoryName string
	Score        float64
}

func contentKind(canonicalURL string) (string, string) {
	path := strings.ToLower(strings.TrimSpace(canonicalURL))
	switch {
	case strings.HasPrefix(path, "/tin-tuc/tin/"), strings.HasPrefix(path, "/tin-tuc/bai-viet/"):
		return "news_article", "Tin quy hoạch"
	case strings.HasPrefix(path, "/tin-tuc/bao-cao/"):
		return "analysis_report", "Báo cáo và phân tích"
	case strings.HasPrefix(path, "/tin-tuc/do-an-quy-hoach/"):
		return "planning_project", "Đồ án quy hoạch"
	case strings.HasPrefix(path, "/tin-tuc/van-ban/"), strings.HasPrefix(path, "/tin-tuc/tai-lieu/"):
		return "legal_document", "Văn bản pháp lý"
	case strings.HasPrefix(path, "/tin-tuc/vung-quy-hoach/"):
		return "planning_region", "Vùng quy hoạch"
	default:
		return "news_article", "Nội dung quy hoạch"
	}
}

func (r *newsPostgres) SearchContent(ctx context.Context, query string, kinds []string, limit int) (*dto.ContentSearchResult, error) {
	pattern := likeKeyword(query)
	if pattern == "" {
		return &dto.ContentSearchResult{Groups: []dto.SearchGroup{}}, nil
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	q := strings.ToLower(strings.TrimSpace(strings.ReplaceAll(query, "+", " ")))
	rows := []contentSearchRow{}
	err := r.db.WithContext(ctx).Raw(`
		SELECT sd.id, sd.slug, COALESCE(sd.title,'') AS title,
		       COALESCE(sd.description,'') AS description, COALESCE(sd.summary,'') AS summary,
		       sd.canonical_url, COALESCE(sd.deep_link,'') AS deep_link,
		       COALESCE(sd.image_url,'') AS image_url,
		       sd.published_at, sd.updated_at,
		       COALESCE(c.code,'') AS category_code, COALESCE(c.name,'') AS category_name,
		       CASE WHEN LOWER(COALESCE(sd.title,'')) = ? THEN 100
		            WHEN LOWER(COALESCE(sd.title,'')) LIKE ? THEN 70 ELSE 35 END::float8 AS score
		FROM seo_domain sd
		LEFT JOIN seo_category c ON c.id = sd.category_id AND c.deleted_at IS NULL
		WHERE sd.deleted_at IS NULL AND sd.published = TRUE AND sd.scope = 'public'
		  AND (LOWER(COALESCE(sd.title,'')) LIKE ? OR LOWER(COALESCE(sd.description,'')) LIKE ?
		       OR LOWER(COALESCE(sd.summary,'')) LIKE ? OR LOWER(COALESCE(sd.slug,'')) LIKE ?)
		ORDER BY score DESC, sd.published_at DESC NULLS LAST, sd.id DESC LIMIT ?`,
		q, q+"%", pattern, pattern, pattern, pattern, limit*2).Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("search content: %w", err)
	}
	allowed := map[string]bool{}
	for _, k := range kinds {
		allowed[strings.TrimSpace(k)] = true
	}
	byKind := map[string][]dto.SearchCandidate{}
	labels := map[string]string{}
	for _, row := range rows {
		kind, label := contentKind(row.CanonicalURL)
		if len(allowed) > 0 && !allowed[kind] {
			continue
		}
		item := dto.SearchCandidate{}
		item.Entity = dto.SearchEntityRef{Kind: kind, ID: fmt.Sprintf("%d", row.ID), Key: kind + ":" + fmt.Sprintf("%d", row.ID)}
		item.Presentation.Title = row.Title
		item.Presentation.Subtitle = row.CategoryName
		item.Presentation.Description = row.Summary
		item.Source.System = "seo"
		item.Source.Dataset = "seo_domain"
		item.Capabilities.CanOpenDetail = true
		item.Capabilities.CanShare = true
		item.Capabilities.CanFocusMap = row.DeepLink != ""
		item.Links.CanonicalURL = row.CanonicalURL
		item.Links.DetailURL = row.CanonicalURL
		item.Match.Type = "full_text"
		item.Match.MatchedFields = []string{"title", "description", "summary", "slug"}
		item.Rank.Score = row.Score
		item.Rank.ReasonCodes = []string{"published_content"}
		item.Attributes = map[string]any{"categoryCode": row.CategoryCode, "categoryName": row.CategoryName, "deepLink": row.DeepLink, "imageUrl": row.ImageURL, "publishedAt": row.PublishedAt}
		byKind[kind] = append(byKind[kind], item)
		labels[kind] = label
	}
	order := []string{"news_article", "analysis_report", "legal_document", "planning_project", "planning_region"}
	groups := []dto.SearchGroup{}
	for _, kind := range order {
		if items := byKind[kind]; len(items) > 0 {
			if len(items) > limit {
				items = items[:limit]
			}
			groups = append(groups, dto.SearchGroup{Kind: kind, Label: labels[kind], Total: len(items), Items: items})
		}
	}
	result := &dto.ContentSearchResult{RequestID: fmt.Sprintf("seo-%d", time.Now().UnixNano()), Groups: groups}
	result.Interpretation.OriginalQuery = query
	result.Interpretation.NormalizedQuery = q
	return result, nil
}
