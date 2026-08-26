package contentcatalog

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"crm/internal/modules/contentcatalog/domain"

	"gorm.io/gorm"
)

type Repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) *Repository { return &Repository{db: db} }

const planningNewsWhere = `
    sd.deleted_at IS NULL
    AND sd.published = TRUE
    AND sd.scope = 'public'
    AND (
        sd.canonical_url LIKE '/quy-hoach/tin-tuc/%'
        OR sd.template_key IN ('planning-news', 'planning-news-detail')
    )`

type planningNewsRow struct {
	ID                                                                                                                             uint64
	SEOPageID                                                                                                                      uint64
	Slug, Title, Description, Summary, Content, RenderedHTML, CanonicalURL, PublicURL, DeepLink, RefSource, RefLabel, SourceStatus string
	PublishedAt, UpdatedAt                                                                                                         *time.Time
	SavedAt                                                                                                                        *time.Time
	MetadataJSON, RefSnapshotRaw                                                                                                   []byte
}

func decodeJSONObject(raw []byte) map[string]any {
	out := map[string]any{}
	if len(raw) > 0 {
		_ = json.Unmarshal(raw, &out)
	}
	return out
}

func newsFromRow(row planningNewsRow) domain.PlanningNewsItem {
	return domain.PlanningNewsItem{
		ID:              row.ID,
		SEOPageID:       row.ID,
		Slug:            row.Slug,
		Title:           row.Title,
		Description:     row.Description,
		Summary:         row.Summary,
		Content:         row.Content,
		RenderedHTML:    row.RenderedHTML,
		CanonicalURL:    row.CanonicalURL,
		PublicURL:       row.PublicURL,
		DeepLink:        row.DeepLink,
		RefSource:       row.RefSource,
		RefLabel:        row.RefLabel,
		SourceStatus:    row.SourceStatus,
		PublishedAt:     row.PublishedAt,
		UpdatedAt:       row.UpdatedAt,
		SavedAt:         row.SavedAt,
		Metadata:        decodeJSONObject(row.MetadataJSON),
		RefSnapshotJSON: decodeJSONObject(row.RefSnapshotRaw),
	}
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

func (r *Repository) ListCategories(ctx context.Context) ([]domain.ContentCategory, error) {
	items := []domain.ContentCategory{}
	err := r.db.WithContext(ctx).Raw(`
		SELECT c.id, c.code, c.slug, c.name,
		       COALESCE(c.description,'') AS description, c.status, c.sort_order,
		       COUNT(DISTINCT CASE WHEN sd.deleted_at IS NULL AND sd.published = TRUE AND sd.scope='public' THEN sd.id END) AS item_count
		FROM seo_category c
		LEFT JOIN seo_domain sd ON sd.category_id = c.id
		WHERE c.deleted_at IS NULL AND c.status = 'active'
		GROUP BY c.id
		ORDER BY c.sort_order, c.name, c.id`).Scan(&items).Error
	if err != nil {
		return nil, fmt.Errorf("list content categories: %w", err)
	}
	return items, nil
}

func (r *Repository) GetCategory(ctx context.Context, id uint64) (*domain.ContentCategory, error) {
	if id == 0 {
		return nil, fmt.Errorf("category id is required")
	}
	item := domain.ContentCategory{}
	err := r.db.WithContext(ctx).Raw(`
		SELECT c.id, c.code, c.slug, c.name,
		       COALESCE(c.description,'') AS description, c.status, c.sort_order,
		       COUNT(DISTINCT CASE WHEN sd.deleted_at IS NULL AND sd.published = TRUE AND sd.scope='public' THEN sd.id END) AS item_count
		FROM seo_category c
		LEFT JOIN seo_domain sd ON sd.category_id = c.id
		WHERE c.deleted_at IS NULL AND c.id = ?
		GROUP BY c.id
		LIMIT 1`, id).Scan(&item).Error
	if err != nil {
		return nil, fmt.Errorf("get seo category: %w", err)
	}
	if item.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &item, nil
}

func (r *Repository) CreateCategory(ctx context.Context, input domain.CreateSeoCategoryInput) (*domain.ContentCategory, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, fmt.Errorf("name is required")
	}
	code := normalizeCategoryCode(input.Code)
	if code == "" {
		code = normalizeCategoryCode(name)
	}
	if code == "" {
		return nil, fmt.Errorf("code is required")
	}
	slug := strings.TrimSpace(input.Slug)
	if slug == "" {
		slug = slugifyCategory(firstNonEmpty(name, code))
	} else {
		slug = slugifyCategory(slug)
	}
	if slug == "" {
		return nil, fmt.Errorf("slug is required")
	}
	status, err := normalizeCategoryStatus(input.Status, "active")
	if err != nil {
		return nil, err
	}
	sortOrder := input.SortOrder
	if sortOrder < 0 {
		sortOrder = 0
	}
	description := strings.TrimSpace(input.Description)

	var id uint64
	err = r.db.WithContext(ctx).Raw(`
		INSERT INTO seo_category (
			code, slug, name, description, status, sort_order, metadata, created_at, updated_at
		) VALUES (?, ?, ?, NULLIF(?, ''), ?, ?, '{}'::jsonb, NOW(), NOW())
		RETURNING id
	`, code, slug, name, description, status, sortOrder).Scan(&id).Error
	if err != nil {
		if isUniqueViolation(err) {
			return nil, fmt.Errorf("category code or slug already exists")
		}
		return nil, fmt.Errorf("create seo category: %w", err)
	}
	return r.GetCategory(ctx, id)
}

func (r *Repository) UpdateCategory(ctx context.Context, id uint64, input domain.UpdateSeoCategoryInput) (*domain.ContentCategory, error) {
	if id == 0 {
		return nil, fmt.Errorf("category id is required")
	}
	current, err := r.GetCategory(ctx, id)
	if err != nil {
		return nil, err
	}

	updates := map[string]any{"updated_at": time.Now()}
	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		if name == "" {
			return nil, fmt.Errorf("name is required")
		}
		updates["name"] = name
	}
	if input.Code != nil {
		code := normalizeCategoryCode(*input.Code)
		if code == "" {
			return nil, fmt.Errorf("code is required")
		}
		updates["code"] = code
	}
	if input.Slug != nil {
		slug := slugifyCategory(*input.Slug)
		if slug == "" {
			return nil, fmt.Errorf("slug is required")
		}
		updates["slug"] = slug
	}
	if input.Description != nil {
		updates["description"] = strings.TrimSpace(*input.Description)
	}
	if input.Status != nil {
		status, statusErr := normalizeCategoryStatus(*input.Status, "")
		if statusErr != nil {
			return nil, statusErr
		}
		updates["status"] = status
	}
	if input.SortOrder != nil {
		sortOrder := *input.SortOrder
		if sortOrder < 0 {
			sortOrder = 0
		}
		updates["sort_order"] = sortOrder
	}
	if len(updates) == 1 {
		return current, nil
	}

	result := r.db.WithContext(ctx).Table("seo_category").
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(updates)
	if result.Error != nil {
		if isUniqueViolation(result.Error) {
			return nil, fmt.Errorf("category code or slug already exists")
		}
		return nil, fmt.Errorf("update seo category: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return r.GetCategory(ctx, id)
}

func (r *Repository) DeleteCategory(ctx context.Context, id uint64) error {
	if id == 0 {
		return fmt.Errorf("category id is required")
	}
	now := time.Now()
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		if recoverValue := recover(); recoverValue != nil {
			tx.Rollback()
			panic(recoverValue)
		}
	}()

	result := tx.Table("seo_category").
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{"deleted_at": now, "updated_at": now})
	if result.Error != nil {
		tx.Rollback()
		return fmt.Errorf("delete seo category: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		tx.Rollback()
		return gorm.ErrRecordNotFound
	}
	if err := tx.Table("seo_domain").
		Where("category_id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{"category_id": nil, "updated_at": now}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("clear seo_domain.category_id: %w", err)
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) AssignPlanningNewsCategory(ctx context.Context, id uint64, categoryID uint64) (*domain.PlanningNewsItem, error) {
	if id == 0 || categoryID == 0 {
		return nil, fmt.Errorf("planning news id and categoryId are required")
	}
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	defer func() {
		if recoverValue := recover(); recoverValue != nil {
			tx.Rollback()
			panic(recoverValue)
		}
	}()

	var pageCount int64
	if err := tx.Table("seo_domain").Where("id=? AND deleted_at IS NULL", id).Count(&pageCount).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if pageCount == 0 {
		tx.Rollback()
		return nil, gorm.ErrRecordNotFound
	}

	var foundCategoryID uint64
	if err := tx.Table("seo_category").
		Select("id").
		Where("id=? AND deleted_at IS NULL AND status='active'", categoryID).
		Limit(1).
		Scan(&foundCategoryID).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if foundCategoryID == 0 {
		tx.Rollback()
		return nil, gorm.ErrRecordNotFound
	}

	if err := tx.Table("seo_domain").
		Where("id=? AND deleted_at IS NULL", id).
		Updates(map[string]any{"category_id": foundCategoryID, "updated_at": time.Now()}).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("set seo_domain.category_id: %w", err)
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	detail, err := r.GetPlanningNews(ctx, strconv.FormatUint(id, 10))
	if err != nil {
		return nil, err
	}
	return &detail.SEODomain, nil
}

func normalizeCategoryStatus(value, fallback string) (string, error) {
	status := strings.ToLower(strings.TrimSpace(value))
	if status == "" {
		status = strings.ToLower(strings.TrimSpace(fallback))
	}
	switch status {
	case "active", "inactive", "archived":
		return status, nil
	default:
		return "", fmt.Errorf("status must be active, inactive or archived")
	}
}

func normalizeCategoryCode(value string) string {
	value = strings.TrimSpace(strings.ToLower(value))
	value = strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z':
			return r
		case r >= '0' && r <= '9':
			return r
		case r == '_' || r == '-':
			return '_'
		default:
			return '_'
		}
	}, value)
	for strings.Contains(value, "__") {
		value = strings.ReplaceAll(value, "__", "_")
	}
	return strings.Trim(value, "_")
}

func slugifyCategory(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	replacer := strings.NewReplacer(
		"à", "a", "á", "a", "ạ", "a", "ả", "a", "ã", "a",
		"â", "a", "ầ", "a", "ấ", "a", "ậ", "a", "ẩ", "a", "ẫ", "a",
		"ă", "a", "ằ", "a", "ắ", "a", "ặ", "a", "ẳ", "a", "ẵ", "a",
		"è", "e", "é", "e", "ẹ", "e", "ẻ", "e", "ẽ", "e",
		"ê", "e", "ề", "e", "ế", "e", "ệ", "e", "ể", "e", "ễ", "e",
		"ì", "i", "í", "i", "ị", "i", "ỉ", "i", "ĩ", "i",
		"ò", "o", "ó", "o", "ọ", "o", "ỏ", "o", "õ", "o",
		"ô", "o", "ồ", "o", "ố", "o", "ộ", "o", "ổ", "o", "ỗ", "o",
		"ơ", "o", "ờ", "o", "ớ", "o", "ợ", "o", "ở", "o", "ỡ", "o",
		"ù", "u", "ú", "u", "ụ", "u", "ủ", "u", "ũ", "u",
		"ư", "u", "ừ", "u", "ứ", "u", "ự", "u", "ử", "u", "ữ", "u",
		"ỳ", "y", "ý", "y", "ỵ", "y", "ỷ", "y", "ỹ", "y",
		"đ", "d",
	)
	value = replacer.Replace(value)
	value = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			return r
		}
		return '-'
	}, value)
	for strings.Contains(value, "--") {
		value = strings.ReplaceAll(value, "--", "-")
	}
	return strings.Trim(value, "-")
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "duplicate key") ||
		strings.Contains(message, "unique constraint") ||
		strings.Contains(message, "uq_seo_category")
}

func (r *Repository) ListPlanningNews(ctx context.Context, filter domain.PlanningNewsListFilter) (*domain.PlanningNewsListResult, error) {
	filter.Page, filter.Size = normalizePage(filter.Page, filter.Size)
	base := ` FROM seo_domain sd WHERE ` + planningNewsWhere
	args := []any{}
	var total int64
	if err := r.db.WithContext(ctx).Raw("SELECT COUNT(DISTINCT sd.id)"+base, args...).Scan(&total).Error; err != nil {
		return nil, fmt.Errorf("count planning news: %w", err)
	}
	rows := []planningNewsRow{}
	selectSQL := `SELECT sd.id,sd.id AS seo_page_id,sd.slug,COALESCE(sd.title,'') title,COALESCE(sd.description,'') description,COALESCE(sd.summary,'') summary,COALESCE(sd.content,'') content,COALESCE(sd.rendered_html,'') rendered_html,sd.canonical_url,CASE WHEN sd.canonical_url LIKE 'http%' THEN sd.canonical_url ELSE 'https://qhpro.vn'||sd.canonical_url END public_url,COALESCE(sd.deep_link,'') deep_link,COALESCE(sd.ref_source,'') ref_source,COALESCE(sd.ref_label,'') ref_label,COALESCE(sd.source_status,'') source_status,sd.published_at,sd.updated_at,COALESCE(sd.metadata,'{}'::jsonb) metadata_json,COALESCE(sd.ref_snapshot_json,'{}'::jsonb) ref_snapshot_raw` + base + ` ORDER BY sd.published_at DESC NULLS LAST,sd.updated_at DESC,sd.id DESC LIMIT ? OFFSET ?`
	listArgs := append(append([]any{}, args...), filter.Size, (filter.Page-1)*filter.Size)
	if err := r.db.WithContext(ctx).Raw(selectSQL, listArgs...).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("list planning news: %w", err)
	}
	data := make([]domain.PlanningNewsItem, 0, len(rows))
	for _, row := range rows {
		data = append(data, newsFromRow(row))
	}
	return &domain.PlanningNewsListResult{Data: data, Total: total, Page: filter.Page, Size: filter.Size}, nil
}

func (r *Repository) GetPlanningNews(ctx context.Context, slug string) (*domain.PlanningNewsDetailResult, error) {
	slug = strings.Trim(strings.TrimSpace(slug), "/")
	if strings.Contains(slug, "/") {
		parts := strings.Split(slug, "/")
		slug = parts[len(parts)-1]
	}
	row := planningNewsRow{}
	query := `SELECT sd.id,sd.id AS seo_page_id,sd.slug,COALESCE(sd.title,'') title,COALESCE(sd.description,'') description,COALESCE(sd.summary,'') summary,COALESCE(sd.content,'') content,COALESCE(sd.rendered_html,'') rendered_html,sd.canonical_url,CASE WHEN sd.canonical_url LIKE 'http%' THEN sd.canonical_url ELSE 'https://qhpro.vn'||sd.canonical_url END public_url,COALESCE(sd.deep_link,'') deep_link,COALESCE(sd.ref_source,'') ref_source,COALESCE(sd.ref_label,'') ref_label,COALESCE(sd.source_status,'') source_status,sd.published_at,sd.updated_at,COALESCE(sd.metadata,'{}'::jsonb) metadata_json,COALESCE(sd.ref_snapshot_json,'{}'::jsonb) ref_snapshot_raw
      FROM seo_domain sd WHERE ` + planningNewsWhere + ` AND (sd.slug=? OR sd.id::text=?) LIMIT 1`
	if err := r.db.WithContext(ctx).Raw(query, slug, slug).Scan(&row).Error; err != nil {
		return nil, fmt.Errorf("get planning news: %w", err)
	}
	if row.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	relatedRows := []planningNewsRow{}
	relQuery := `SELECT sd.id,sd.id AS seo_page_id,sd.slug,COALESCE(sd.title,'') title,COALESCE(sd.description,'') description,COALESCE(sd.summary,'') summary,'' content,'' rendered_html,sd.canonical_url,CASE WHEN sd.canonical_url LIKE 'http%' THEN sd.canonical_url ELSE 'https://qhpro.vn'||sd.canonical_url END public_url,COALESCE(sd.deep_link,'') deep_link,COALESCE(sd.ref_source,'') ref_source,COALESCE(sd.ref_label,'') ref_label,COALESCE(sd.source_status,'') source_status,sd.published_at,sd.updated_at,COALESCE(sd.metadata,'{}'::jsonb) metadata_json,COALESCE(sd.ref_snapshot_json,'{}'::jsonb) ref_snapshot_raw
      FROM seo_domain sd WHERE ` + planningNewsWhere + ` AND sd.id<>? ORDER BY sd.published_at DESC NULLS LAST LIMIT 8`
	if err := r.db.WithContext(ctx).Raw(relQuery, row.ID).Scan(&relatedRows).Error; err != nil {
		return nil, fmt.Errorf("list related planning news: %w", err)
	}
	related := make([]domain.PlanningNewsItem, 0, len(relatedRows))
	for _, rr := range relatedRows {
		related = append(related, newsFromRow(rr))
	}
	return &domain.PlanningNewsDetailResult{SEODomain: newsFromRow(row), RelatedNews: related}, nil
}

func (r *Repository) ListSavedPlanningNews(ctx context.Context, userID uint64, filter domain.SavedPlanningNewsListFilter) (*domain.PlanningNewsListResult, error) {
	if userID == 0 {
		return nil, fmt.Errorf("user id is required")
	}
	filter.Page, filter.Size = normalizePage(filter.Page, filter.Size)

	base := ` FROM user_saved_planning_news saved
	      JOIN seo_domain sd ON sd.id=saved.seo_domain_id
	      WHERE saved.user_id=? AND saved.deleted_at IS NULL AND ` + planningNewsWhere

	var total int64
	if err := r.db.WithContext(ctx).Raw("SELECT COUNT(DISTINCT sd.id)"+base, userID).Scan(&total).Error; err != nil {
		return nil, fmt.Errorf("count saved planning news: %w", err)
	}

	rows := []planningNewsRow{}
	selectSQL := `SELECT sd.id,sd.id AS seo_page_id,sd.slug,COALESCE(sd.title,'') title,COALESCE(sd.description,'') description,COALESCE(sd.summary,'') summary,'' content,'' rendered_html,sd.canonical_url,CASE WHEN sd.canonical_url LIKE 'http%' THEN sd.canonical_url ELSE 'https://qhpro.vn'||sd.canonical_url END public_url,COALESCE(sd.deep_link,'') deep_link,COALESCE(sd.ref_source,'') ref_source,COALESCE(sd.ref_label,'') ref_label,COALESCE(sd.source_status,'') source_status,sd.published_at,sd.updated_at,saved.updated_at AS saved_at,COALESCE(sd.metadata,'{}'::jsonb) metadata_json,COALESCE(sd.ref_snapshot_json,'{}'::jsonb) ref_snapshot_raw` + base + ` ORDER BY saved.updated_at DESC,saved.id DESC LIMIT ? OFFSET ?`
	if err := r.db.WithContext(ctx).Raw(selectSQL, userID, filter.Size, (filter.Page-1)*filter.Size).Scan(&rows).Error; err != nil {
		return nil, fmt.Errorf("list saved planning news: %w", err)
	}

	data := make([]domain.PlanningNewsItem, 0, len(rows))
	for _, row := range rows {
		data = append(data, newsFromRow(row))
	}
	return &domain.PlanningNewsListResult{Data: data, Total: total, Page: filter.Page, Size: filter.Size}, nil
}

func (r *Repository) SavePlanningNews(ctx context.Context, userID, newsID uint64) (*domain.PlanningNewsItem, error) {
	if userID == 0 || newsID == 0 {
		return nil, fmt.Errorf("user id and planning news id are required")
	}
	result := r.db.WithContext(ctx).Exec(`
		INSERT INTO user_saved_planning_news (user_id,seo_domain_id,created_at,updated_at,deleted_at)
		SELECT ?,sd.id,NOW(),NOW(),NULL
		FROM seo_domain sd
		WHERE `+planningNewsWhere+` AND sd.id=?
		ON CONFLICT (user_id,seo_domain_id)
		DO UPDATE SET updated_at=NOW(),deleted_at=NULL`, userID, newsID)
	if result.Error != nil {
		return nil, fmt.Errorf("save planning news: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	detail, err := r.GetPlanningNews(ctx, strconv.FormatUint(newsID, 10))
	if err != nil {
		return nil, err
	}
	now := time.Now()
	item := detail.SEODomain
	item.SavedAt = &now
	return &item, nil
}

func (r *Repository) UnsavePlanningNews(ctx context.Context, userID, newsID uint64) error {
	if userID == 0 || newsID == 0 {
		return fmt.Errorf("user id and planning news id are required")
	}
	if err := r.db.WithContext(ctx).Table("user_saved_planning_news").
		Where("user_id=? AND seo_domain_id=? AND deleted_at IS NULL", userID, newsID).
		Updates(map[string]any{"deleted_at": time.Now(), "updated_at": time.Now()}).Error; err != nil {
		return fmt.Errorf("unsave planning news: %w", err)
	}
	return nil
}

type contentSearchRow struct {
	ID           uint64
	Slug         string
	Title        string
	Description  string
	Summary      string
	CanonicalURL string
	DeepLink     string
	PublishedAt  *time.Time
	UpdatedAt    *time.Time
	Score        float64
}

func contentKind(canonicalURL string) (string, string) {
	path := strings.ToLower(strings.TrimSpace(canonicalURL))
	switch {
	case strings.HasPrefix(path, "/quy-hoach/tin-tuc/"):
		return "news_article", "Tin quy hoạch"
	case strings.HasPrefix(path, "/quy-hoach/bao-cao/"):
		return "analysis_report", "Báo cáo và phân tích"
	case strings.HasPrefix(path, "/do-an-quy-hoach/"):
		return "planning_project", "Đồ án quy hoạch"
	case strings.HasPrefix(path, "/van-ban/"):
		return "legal_document", "Văn bản pháp lý"
	case strings.HasPrefix(path, "/vung-quy-hoach/"):
		return "planning_region", "Vùng quy hoạch"
	default:
		return "news_article", "Nội dung quy hoạch"
	}
}

func (r *Repository) SearchPublishedContent(ctx context.Context, query string, kinds []string, limit int) (*domain.ContentSearchResult, error) {
	query = strings.Join(strings.Fields(strings.ToLower(strings.TrimSpace(query))), " ")
	if query == "" {
		return &domain.ContentSearchResult{Groups: []domain.SearchGroup{}}, nil
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	pattern := "%" + query + "%"
	rows := []contentSearchRow{}
	err := r.db.WithContext(ctx).Raw(`
		SELECT sd.id, sd.slug, COALESCE(sd.title,'') AS title,
		       COALESCE(sd.description,'') AS description, COALESCE(sd.summary,'') AS summary,
		       sd.canonical_url, COALESCE(sd.deep_link,'') AS deep_link,
		       sd.published_at, sd.updated_at,
		       CASE WHEN LOWER(COALESCE(sd.title,'')) = ? THEN 100
		            WHEN LOWER(COALESCE(sd.title,'')) LIKE ? THEN 70 ELSE 35 END::float8 AS score
		FROM seo_domain sd
		WHERE sd.deleted_at IS NULL AND sd.published=TRUE AND sd.scope='public'
		  AND (LOWER(COALESCE(sd.title,'')) LIKE ? OR LOWER(COALESCE(sd.description,'')) LIKE ?
		       OR LOWER(COALESCE(sd.summary,'')) LIKE ? OR LOWER(COALESCE(sd.slug,'')) LIKE ?)
		ORDER BY score DESC, sd.published_at DESC NULLS LAST, sd.id DESC LIMIT ?`,
		query, query+"%", pattern, pattern, pattern, pattern, limit*2).Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("search published content: %w", err)
	}
	allowed := map[string]bool{}
	for _, k := range kinds {
		allowed[strings.TrimSpace(k)] = true
	}
	byKind := map[string][]domain.SearchCandidate{}
	labels := map[string]string{}
	for _, row := range rows {
		kind, label := contentKind(row.CanonicalURL)
		if len(allowed) > 0 && !allowed[kind] {
			continue
		}
		item := domain.SearchCandidate{}
		id := strconv.FormatUint(row.ID, 10)
		item.Entity = domain.SearchEntityRef{Kind: kind, ID: id, Key: kind + ":" + id}
		item.Presentation.Title = row.Title
		item.Presentation.Description = row.Summary
		item.Source.System = "crm-service"
		item.Source.Dataset = "seo_domain"
		item.Capabilities.CanOpenDetail = true
		item.Capabilities.CanShare = true
		item.Capabilities.CanFocusMap = row.DeepLink != ""
		item.Links.CanonicalURL = row.CanonicalURL
		item.Links.DetailURL = row.CanonicalURL
		item.Match.Type = "full_text"
		item.Match.MatchedFields = []string{"title", "description", "summary", "slug"}
		item.Rank.Score = row.Score
		item.Attributes = map[string]any{"deepLink": row.DeepLink, "publishedAt": row.PublishedAt}
		byKind[kind] = append(byKind[kind], item)
		labels[kind] = label
	}
	order := []string{"news_article", "analysis_report", "legal_document", "planning_project", "planning_region"}
	groups := []domain.SearchGroup{}
	for _, kind := range order {
		if items := byKind[kind]; len(items) > 0 {
			if len(items) > limit {
				items = items[:limit]
			}
			groups = append(groups, domain.SearchGroup{Kind: kind, Label: labels[kind], Total: len(items), Items: items})
		}
	}
	result := &domain.ContentSearchResult{RequestID: fmt.Sprintf("crm-%d", time.Now().UnixNano()), Groups: groups}
	result.Interpretation.OriginalQuery = query
	result.Interpretation.NormalizedQuery = query
	return result, nil
}
