package contentcatalog

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"crm/internal/enums"
	"crm/internal/modules/contentcatalog/domain"

	"gorm.io/gorm"
)

type Repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) *Repository { return &Repository{db: db} }

// Client/public seo/news: chỉ trả trang có published = true.
const planningNewsWhere = `
    sd.deleted_at IS NULL
    AND sd.published = TRUE
    AND sd.scope = 'public'
    AND (
        sd.canonical_url LIKE '/tin-tuc/%'
        OR sd.template_key IN ('planning-news', 'planning-news-detail')
        OR c.code = 'planning_news'
    )`

const planningNewsCategoryJoin = `
      LEFT JOIN seo_content_category_assignment a ON a.seo_domain_id=sd.id AND a.deleted_at IS NULL AND a.is_primary=TRUE
      LEFT JOIN seo_category c ON c.id=COALESCE(sd.category_id, a.category_id) AND c.deleted_at IS NULL`

const planningNewsSelectCols = `sd.id,sd.id AS seo_page_id,sd.slug,COALESCE(sd.title,'') title,COALESCE(sd.description,'') description,COALESCE(sd.summary,'') summary,COALESCE(sd.content,'') AS "content",COALESCE(sd.rendered_html,'') rendered_html,sd.canonical_url,sd.canonical_url AS public_url,COALESCE(sd.deep_link,'') deep_link,COALESCE(sd.image_url,'') image_url,COALESCE(sd.ref_source,'') ref_source,COALESCE(sd.ref_label,'') ref_label,COALESCE(sd.source_status,'') source_status,COALESCE(sd.quality_score,0) quality_score,sd.visible_on,sd.published_at,sd.updated_at,COALESCE(sd.category_id, c.id) AS category_id,COALESCE(c.code,'') category_code,COALESCE(c.name,'') category_name,COALESCE(sd.metadata,'{}'::jsonb) metadata_json,COALESCE(sd.ref_snapshot_json,'{}'::jsonb) ref_snapshot_raw,COALESCE(sd.ref_type,0) ref_type,sd.ref_id`

const planningNewsSelectColsLite = `sd.id,sd.id AS seo_page_id,sd.slug,COALESCE(sd.title,'') title,COALESCE(sd.description,'') description,COALESCE(sd.summary,'') summary,'' AS "content",'' rendered_html,sd.canonical_url,sd.canonical_url AS public_url,COALESCE(sd.deep_link,'') deep_link,COALESCE(sd.image_url,'') image_url,COALESCE(sd.ref_source,'') ref_source,COALESCE(sd.ref_label,'') ref_label,COALESCE(sd.source_status,'') source_status,COALESCE(sd.quality_score,0) quality_score,sd.visible_on,sd.published_at,sd.updated_at,COALESCE(sd.category_id, c.id) AS category_id,COALESCE(c.code,'') category_code,COALESCE(c.name,'') category_name,COALESCE(sd.metadata,'{}'::jsonb) metadata_json,COALESCE(sd.ref_snapshot_json,'{}'::jsonb) ref_snapshot_raw,COALESCE(sd.ref_type,0) ref_type,sd.ref_id`

func (r *Repository) ListDomains(ctx context.Context) ([]domain.ContentDomain, error) {
	items := []domain.ContentDomain{}
	err := r.db.WithContext(ctx).Raw(`
        SELECT id, code, name, COALESCE(description,'') AS description, sort_order
        FROM seo_content_domain
        WHERE deleted_at IS NULL AND status = 'active'
        ORDER BY sort_order, name, id`).Scan(&items).Error
	if err != nil {
		return nil, fmt.Errorf("list content domains: %w", err)
	}
	return items, nil
}

func (r *Repository) ListCategories(ctx context.Context, domainCode string) ([]domain.ContentCategory, error) {
	items := []domain.ContentCategory{}
	query := `
        SELECT c.id, c.domain_id, d.code AS domain_code, c.code, c.slug, c.name,
               COALESCE(c.description,'') AS description, c.status, c.sort_order,
               COUNT(DISTINCT CASE WHEN sd.deleted_at IS NULL AND sd.published = TRUE AND sd.scope='public' THEN sd.id END) AS item_count
        FROM seo_category c
        JOIN seo_content_domain d ON d.id = c.domain_id AND d.deleted_at IS NULL
        LEFT JOIN seo_content_category_assignment a ON a.category_id=c.id AND a.deleted_at IS NULL
        LEFT JOIN seo_domain sd ON sd.id=a.seo_domain_id
        WHERE c.deleted_at IS NULL AND c.status = 'active' AND d.status = 'active'`
	args := []any{}
	if value := strings.TrimSpace(domainCode); value != "" {
		query += " AND d.code = ?"
		args = append(args, value)
	}
	query += " GROUP BY c.id,d.code,d.sort_order ORDER BY d.sort_order,c.sort_order,c.name,c.id"
	if err := r.db.WithContext(ctx).Raw(query, args...).Scan(&items).Error; err != nil {
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
		SELECT c.id, c.domain_id, d.code AS domain_code, c.code, c.slug, c.name,
		       COALESCE(c.description,'') AS description, c.status, c.sort_order,
		       COUNT(DISTINCT CASE WHEN sd.deleted_at IS NULL AND sd.published = TRUE AND sd.scope='public' THEN sd.id END) AS item_count
		FROM seo_category c
		JOIN seo_content_domain d ON d.id = c.domain_id AND d.deleted_at IS NULL
		LEFT JOIN seo_content_category_assignment a ON a.category_id=c.id AND a.deleted_at IS NULL
		LEFT JOIN seo_domain sd ON sd.id=a.seo_domain_id
		WHERE c.deleted_at IS NULL AND c.id = ?
		GROUP BY c.id, d.code
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
	domainID, err := r.resolveDomainID(ctx, input.DomainID, input.DomainCode)
	if err != nil {
		return nil, err
	}
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
			domain_id, code, slug, name, description, status, sort_order, metadata, created_at, updated_at
		) VALUES (?, ?, ?, ?, NULLIF(?, ''), ?, ?, '{}'::jsonb, NOW(), NOW())
		RETURNING id
	`, domainID, code, slug, name, description, status, sortOrder).Scan(&id).Error
	if err != nil {
		if isUniqueViolation(err) {
			return nil, fmt.Errorf("category code or slug already exists in this domain")
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
	domainID := current.DomainID
	if input.DomainID != nil || (input.DomainCode != nil && strings.TrimSpace(*input.DomainCode) != "") {
		domainCode := ""
		if input.DomainCode != nil {
			domainCode = *input.DomainCode
		}
		var requestedDomainID uint64
		if input.DomainID != nil {
			requestedDomainID = *input.DomainID
		}
		resolved, resolveErr := r.resolveDomainID(ctx, requestedDomainID, domainCode)
		if resolveErr != nil {
			return nil, resolveErr
		}
		domainID = resolved
		updates["domain_id"] = domainID
	}
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
			return nil, fmt.Errorf("category code or slug already exists in this domain")
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
	if err := tx.Table("seo_content_category_assignment").
		Where("category_id = ? AND deleted_at IS NULL", id).
		Updates(map[string]any{"deleted_at": now}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("clear category assignments: %w", err)
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}

func (r *Repository) resolveDomainID(ctx context.Context, domainID uint64, domainCode string) (uint64, error) {
	domainCode = strings.TrimSpace(domainCode)
	if domainID == 0 && domainCode == "" {
		return 0, fmt.Errorf("domainId or domain is required")
	}
	query := r.db.WithContext(ctx).Table("seo_content_domain").
		Select("id").
		Where("deleted_at IS NULL AND status = 'active'")
	if domainID > 0 {
		query = query.Where("id = ?", domainID)
	} else {
		query = query.Where("code = ?", domainCode)
	}
	var resolved uint64
	if err := query.Limit(1).Scan(&resolved).Error; err != nil {
		return 0, fmt.Errorf("resolve content domain: %w", err)
	}
	if resolved == 0 {
		return 0, gorm.ErrRecordNotFound
	}
	return resolved, nil
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
		strings.Contains(message, "uq_seo_content_category")
}

// visibleOnForFeed là nguồn mapping của public feed hiện tại.
//
// Các alias chỉ phục vụ tương thích với frontend đang tồn tại;
// database vẫn chỉ sử dụng các enum 10, 20, 30, 40 và 50.
func visibleOnForFeed(feed string) (uint32, bool) {
	switch strings.ToLower(strings.TrimSpace(feed)) {
	case "top", "pinned", "featured":
		return 10, true
	case "notable", "attention":
		return 20, true
	case "latest":
		return 30, true
	case "most_read", "trending":
		return 40, true
	case "quick_news", "breaking":
		return 50, true
	default:
		return 0, false
	}
}

func normalizeFeedLimit(limit int) int {
	if limit <= 0 {
		return 20
	}
	if limit > 100 {
		return 100
	}
	return limit
}

func normalizePlanningFeedDomain(domainCode string) (string, error) {
	value := strings.ToLower(strings.TrimSpace(domainCode))
	if value == "" || value == "planning" {
		return "planning", nil
	}

	return "", fmt.Errorf(
		"unsupported content domain %q for planning news feed",
		domainCode,
	)
}

func (r *Repository) ListFeed(
	ctx context.Context,
	feed string,
	domainCode string,
	categoryCode string,
	limit int,
) (*domain.FeedResult, error) {
	feed = strings.ToLower(strings.TrimSpace(feed))
	if feed == "" {
		return nil, fmt.Errorf("content feed is required")
	}

	visibleOn, ok := visibleOnForFeed(feed)
	if !ok {
		return nil, fmt.Errorf("unsupported content feed %q", feed)
	}

	if _, err := normalizePlanningFeedDomain(domainCode); err != nil {
		return nil, err
	}

	return r.listVisibilityFeed(
		ctx,
		feed,
		visibleOn,
		strings.TrimSpace(categoryCode),
		normalizeFeedLimit(limit),
	)
}

func (r *Repository) listVisibilityFeed(
	ctx context.Context,
	feed string,
	visibleOn uint32,
	categoryCode string,
	limit int,
) (*domain.FeedResult, error) {
	items := []domain.ContentFeedItem{}

	/*
		visible_on là điều kiện lựa chọn bắt buộc của Planning News feed.

		Category assignment là metadata tùy chọn. Vì vậy:
		- không dùng d.code='planning' để lọc;
		- không biến LEFT JOIN category thành INNER JOIN ngầm;
		- không dùng placement để kéo tin có visible_on khác vào feed;
		- placement chỉ được dùng để ưu tiên thứ tự trong cùng visible_on.
	*/
	query := `
        SELECT DISTINCT ON (sd.id)
          sd.id AS seo_page_id,
          sd.slug,
          COALESCE(sd.title,'') AS title,
          COALESCE(sd.description,'') AS description,
          COALESCE(sd.summary,'') AS summary,
          sd.canonical_url,
          sd.published_at,
          sd.updated_at,
          COALESCE(sd.image_url,'') AS image_url,
          COALESCE(c.code,'') AS category_code,
          COALESCE(c.name,'') AS category_name,
          COALESCE(p.placement_code,'') AS placement_code,
          COALESCE(p.priority,0) AS priority,
          sd.visible_on
        FROM seo_domain sd
        LEFT JOIN seo_content_category_assignment a
          ON a.seo_domain_id = sd.id
         AND a.deleted_at IS NULL
         AND a.is_primary = TRUE
        LEFT JOIN seo_category c
          ON c.id = a.category_id
         AND c.deleted_at IS NULL
        LEFT JOIN seo_content_placement p
          ON p.seo_domain_id = sd.id
         AND p.deleted_at IS NULL
         AND p.status = 'active'
         AND p.placement_code = ?
         AND (p.starts_at IS NULL OR p.starts_at <= NOW())
         AND (p.ends_at IS NULL OR p.ends_at > NOW())
        WHERE ` + planningNewsWhere + `
          AND sd.visible_on = ?`

	args := []any{feed, visibleOn}

	if categoryCode != "" {
		query += `
          AND (
            c.code = ?
            OR c.slug = ?
            OR LOWER(c.name) = LOWER(?)
          )`
		args = append(args, categoryCode, categoryCode, categoryCode)
	}

	query += `
        ORDER BY
          sd.id,
          CASE WHEN p.id IS NULL THEN 1 ELSE 0 END,
          p.priority DESC,
          sd.published_at DESC NULLS LAST,
          sd.updated_at DESC,
          sd.id DESC`

	wrapped := `
        SELECT *
        FROM (` + query + `) AS feed_rows
        ORDER BY
          priority DESC,
          published_at DESC NULLS LAST,
          updated_at DESC,
          seo_page_id DESC
        LIMIT ?`

	args = append(args, limit)

	if err := r.db.WithContext(ctx).
		Raw(wrapped, args...).
		Scan(&items).Error; err != nil {
		return nil, fmt.Errorf(
			"list planning news feed %q visible_on=%d: %w",
			feed,
			visibleOn,
			err,
		)
	}

	warnings := []string{}
	if feed == "most_read" || feed == "trending" {
		warnings = append(
			warnings,
			"Feed hiện được phân loại biên tập bằng visible_on=40; viewCount chưa tham gia sắp xếp.",
		)
	}

	return &domain.FeedResult{
		Feed:     feed,
		Computed: false,
		Items:    items,
		Warnings: warnings,
	}, nil
}

type planningNewsRow struct {
	ID                                                                                                                             uint64
	SEOPageID                                                                                                                      uint64
	Slug, Title, Description, Summary, Content, RenderedHTML, CanonicalURL, PublicURL, DeepLink, ImageURL, RefSource, RefLabel, SourceStatus string
	QualityScore, VisibleOn                                                                                                        uint32
	RefType                                                                                                                        uint32
	RefID                                                                                                                          *uint64
	PublishedAt, UpdatedAt                                                                                                         *time.Time
	SavedAt                                                                                                                        *time.Time
	CategoryID                                                                                                                     *uint64
	CategoryCode, CategoryName                                                                                                     string
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
	item := domain.PlanningNewsItem{
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
		ImageURL:        row.ImageURL,
		RefSource:       row.RefSource,
		RefLabel:        row.RefLabel,
		SourceStatus:    row.SourceStatus,
		QualityScore:    row.QualityScore,
		VisibleOn:       row.VisibleOn,
		RefType:         row.RefType,
		RefTypeName:     enums.ESEORefType(row.RefType).Name(),
		RefID:           row.RefID,
		PublishedAt:     row.PublishedAt,
		UpdatedAt:       row.UpdatedAt,
		SavedAt:         row.SavedAt,
		CategoryID:      row.CategoryID,
		CategoryCode:    row.CategoryCode,
		CategoryName:    row.CategoryName,
		Metadata:        decodeJSONObject(row.MetadataJSON),
		RefSnapshotJSON: decodeJSONObject(row.RefSnapshotRaw),
	}
	if row.CategoryID != nil && *row.CategoryID > 0 {
		item.Category = &domain.PlanningNewsCategory{
			ID:   *row.CategoryID,
			Name: row.CategoryName,
		}
	}
	return item
}

func normalizePage(page, size int) (int, int) {
	if page <= 0 {
		page = 0
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
	keyword = strings.ToLower(strings.TrimSpace(keyword))
	keyword = strings.ReplaceAll(keyword, "+", " ")
	parts := strings.Fields(keyword)
	if len(parts) == 0 {
		return ""
	}
	return "%" + strings.Join(parts, "%") + "%"
}

func planningNewsOrderBy(filter domain.PlanningNewsListFilter) string {
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
	case "updated_desc", "updatedat_desc", "updated_at_desc":
		return "sd.updated_at DESC NULLS LAST,sd.id DESC"
	case "updated_asc", "updatedat_asc", "updated_at_asc":
		return "sd.updated_at ASC NULLS LAST,sd.id ASC"
	case "published_asc", "publishedat_asc", "published_at_asc":
		return "sd.published_at ASC NULLS LAST,sd.updated_at ASC,sd.id ASC"
	case "published_desc", "publishedat_desc", "published_at_desc", "":
		return "sd.published_at DESC NULLS LAST,sd.updated_at DESC,sd.id DESC"
	default:
		return "sd.published_at DESC NULLS LAST,sd.updated_at DESC,sd.id DESC"
	}
}

func (r *Repository) ListPlanningNews(ctx context.Context, filter domain.PlanningNewsListFilter) (*domain.PlanningNewsListResult, error) {
	filter.Page, filter.Size = normalizePage(filter.Page, filter.Size)

	base := ` FROM seo_domain sd` + planningNewsCategoryJoin + `
      WHERE ` + planningNewsWhere
	args := []any{}
	if p := likeKeyword(filter.Keyword); p != "" {
		base += " AND (LOWER(COALESCE(sd.title,'')) LIKE ? OR LOWER(COALESCE(sd.description,'')) LIKE ? OR LOWER(COALESCE(sd.summary,'')) LIKE ? OR LOWER(COALESCE(sd.content,'')) LIKE ? OR LOWER(COALESCE(sd.slug,'')) LIKE ?)"
		args = append(args, p, p, p, p, p)
	}
	if filter.CategoryID > 0 {
		base += " AND COALESCE(sd.category_id, a.category_id)=?"
		args = append(args, filter.CategoryID)
	}
	if filter.ProvinceID > 0 {
		base += " AND sd.province_id=?"
		args = append(args, filter.ProvinceID)
	}
	if filter.RefType > 0 {
		base += " AND sd.ref_type=?"
		args = append(args, filter.RefType)
	}
	if p := likeKeyword(filter.Location); p != "" {
		base += " AND (LOWER(COALESCE(sd.ref_label,'')) LIKE ? OR LOWER(COALESCE(sd.metadata->>'location','')) LIKE ? OR LOWER(COALESCE(sd.ref_snapshot_json->>'location','')) LIKE ?)"
		args = append(args, p, p, p)
	}
	if len(filter.VisibleOns) > 0 {
		base += " AND sd.visible_on IN ?"
		args = append(args, filter.VisibleOns)
	} else if filter.VisibleOn > 0 {
		base += " AND sd.visible_on=?"
		args = append(args, filter.VisibleOn)
	}
	switch strings.ToLower(strings.TrimSpace(filter.Range)) {
	case "today":
		base += " AND sd.published_at>=CURRENT_DATE"
	case "week", "7d":
		base += " AND sd.published_at>=NOW()-INTERVAL '7 days'"
	case "month", "30d":
		base += " AND sd.published_at>=NOW()-INTERVAL '30 days'"
	case "year":
		base += " AND sd.published_at>=NOW()-INTERVAL '1 year'"
	}
	var total int64
	if err := r.db.WithContext(ctx).Raw("SELECT COUNT(DISTINCT sd.id)"+base, args...).Scan(&total).Error; err != nil {
		return nil, fmt.Errorf("count planning news: %w", err)
	}
	rows := []planningNewsRow{}
	selectSQL := `SELECT ` + planningNewsSelectCols + base + ` ORDER BY ` + planningNewsOrderBy(filter) + ` LIMIT ? OFFSET ?`
	listArgs := append(append([]any{}, args...), filter.Size, filter.Page*filter.Size)
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
	query := `SELECT ` + planningNewsSelectCols + `
      FROM seo_domain sd` + planningNewsCategoryJoin + `
      WHERE ` + planningNewsWhere + ` AND (sd.slug=? OR sd.id::text=?) LIMIT 1`
	if err := r.db.WithContext(ctx).Raw(query, slug, slug).Scan(&row).Error; err != nil {
		return nil, fmt.Errorf("get planning news: %w", err)
	}
	if row.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	relatedRows := []planningNewsRow{}
	relQuery := `SELECT ` + planningNewsSelectColsLite + `
      FROM seo_domain sd` + planningNewsCategoryJoin + `
      WHERE ` + planningNewsWhere + ` AND sd.id<>? AND (?=0 OR COALESCE(sd.category_id, a.category_id)=?) ORDER BY sd.published_at DESC NULLS LAST LIMIT 8`
	categoryID := uint64(0)
	if row.CategoryID != nil {
		categoryID = *row.CategoryID
	}
	if err := r.db.WithContext(ctx).Raw(relQuery, row.ID, categoryID, categoryID).Scan(&relatedRows).Error; err != nil {
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
	      JOIN seo_domain sd ON sd.id=saved.seo_domain_id` + planningNewsCategoryJoin + `
	      WHERE saved.user_id=? AND saved.deleted_at IS NULL AND ` + planningNewsWhere

	var total int64
	if err := r.db.WithContext(ctx).Raw("SELECT COUNT(DISTINCT sd.id)"+base, userID).Scan(&total).Error; err != nil {
		return nil, fmt.Errorf("count saved planning news: %w", err)
	}

	rows := []planningNewsRow{}
	selectSQL := `SELECT ` + planningNewsSelectColsLite + `,saved.updated_at AS saved_at` + base + ` ORDER BY saved.updated_at DESC,saved.id DESC LIMIT ? OFFSET ?`
	if err := r.db.WithContext(ctx).Raw(selectSQL, userID, filter.Size, filter.Page*filter.Size).Scan(&rows).Error; err != nil {
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
		FROM seo_domain sd`+planningNewsCategoryJoin+`
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

func (r *Repository) ListFollowedSeoCategories(ctx context.Context, userID uint64, filter domain.FollowedSeoCategoryListFilter) (*domain.FollowedSeoCategoryListResult, error) {
	if userID == 0 {
		return nil, fmt.Errorf("user id is required")
	}
	filter.Page, filter.Size = normalizePage(filter.Page, filter.Size)

	base := ` FROM user_followed_seo_category f
	      JOIN seo_category c ON c.id = f.seo_category_id AND c.deleted_at IS NULL AND c.status = 'active'
	      JOIN seo_content_domain d ON d.id = c.domain_id AND d.deleted_at IS NULL AND d.status = 'active'
	      LEFT JOIN seo_content_category_assignment a ON a.category_id = c.id AND a.deleted_at IS NULL
	      LEFT JOIN seo_domain sd ON sd.id = a.seo_domain_id
	      WHERE f.user_id = ? AND f.deleted_at IS NULL`

	var total int64
	if err := r.db.WithContext(ctx).Raw("SELECT COUNT(DISTINCT c.id)"+base, userID).Scan(&total).Error; err != nil {
		return nil, fmt.Errorf("count followed seo categories: %w", err)
	}

	items := []domain.ContentCategory{}
	selectSQL := `SELECT c.id, c.domain_id, d.code AS domain_code, c.code, c.slug, c.name,
	       COALESCE(c.description,'') AS description, c.status, c.sort_order,
	       COUNT(DISTINCT CASE WHEN sd.deleted_at IS NULL AND sd.published = TRUE AND sd.scope='public' THEN sd.id END) AS item_count,
	       f.updated_at AS followed_at` + base + `
	      GROUP BY c.id, d.code, f.updated_at, f.id
	      ORDER BY f.updated_at DESC, f.id DESC
	      LIMIT ? OFFSET ?`
	if err := r.db.WithContext(ctx).Raw(selectSQL, userID, filter.Size, filter.Page*filter.Size).Scan(&items).Error; err != nil {
		return nil, fmt.Errorf("list followed seo categories: %w", err)
	}
	return &domain.FollowedSeoCategoryListResult{Data: items, Total: total, Page: filter.Page, Size: filter.Size}, nil
}

func (r *Repository) FollowSeoCategory(ctx context.Context, userID, categoryID uint64) (*domain.ContentCategory, error) {
	if userID == 0 || categoryID == 0 {
		return nil, fmt.Errorf("user id and category id are required")
	}
	result := r.db.WithContext(ctx).Exec(`
		INSERT INTO user_followed_seo_category (user_id, seo_category_id, created_at, updated_at, deleted_at)
		SELECT ?, c.id, NOW(), NOW(), NULL
		FROM seo_category c
		JOIN seo_content_domain d ON d.id = c.domain_id AND d.deleted_at IS NULL AND d.status = 'active'
		WHERE c.id = ? AND c.deleted_at IS NULL AND c.status = 'active'
		ON CONFLICT (user_id, seo_category_id)
		DO UPDATE SET updated_at = NOW(), deleted_at = NULL`, userID, categoryID)
	if result.Error != nil {
		return nil, fmt.Errorf("follow seo category: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	item, err := r.GetCategory(ctx, categoryID)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	item.FollowedAt = &now
	return item, nil
}

func (r *Repository) UnfollowSeoCategory(ctx context.Context, userID, categoryID uint64) error {
	if userID == 0 || categoryID == 0 {
		return fmt.Errorf("user id and category id are required")
	}
	if err := r.db.WithContext(ctx).Table("user_followed_seo_category").
		Where("user_id=? AND seo_category_id=? AND deleted_at IS NULL", userID, categoryID).
		Updates(map[string]any{"deleted_at": time.Now(), "updated_at": time.Now()}).Error; err != nil {
		return fmt.Errorf("unfollow seo category: %w", err)
	}
	return nil
}

func validPlanningNewsVisibility(value uint32) bool {
	return value == 10 || value == 20 || value == 30 || value == 40 || value == 50
}

func (r *Repository) UpdatePlanningNewsVisibility(ctx context.Context, id uint64, visibleOn uint32) (*domain.PlanningNewsItem, error) {
	if id == 0 {
		return nil, fmt.Errorf("planning news id is required")
	}
	if !validPlanningNewsVisibility(visibleOn) {
		return nil, fmt.Errorf("visibleOn must be one of 10,20,30,40,50")
	}
	result := r.db.WithContext(ctx).Table("seo_domain").
		Where("id=? AND deleted_at IS NULL", id).
		Updates(map[string]any{"visible_on": visibleOn, "updated_at": time.Now()})
	if result.Error != nil {
		return nil, fmt.Errorf("update planning news visibility: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	detail, err := r.GetPlanningNews(ctx, strconv.FormatUint(id, 10))
	if err != nil {
		return nil, err
	}
	return &detail.SEODomain, nil
}

func (r *Repository) SetPlanningNewsTop(ctx context.Context, id uint64) (*domain.PlanningNewsItem, error) {
	if id == 0 {
		return nil, fmt.Errorf("planning news id is required")
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

	var targetCount int64
	if err := tx.Table("seo_domain").Where("id=? AND deleted_at IS NULL", id).Count(&targetCount).Error; err != nil {
		tx.Rollback()
		return nil, err
	}
	if targetCount == 0 {
		tx.Rollback()
		return nil, gorm.ErrRecordNotFound
	}

	// Hạ bản ghi đang là top (visible_on=10) xuống notable (visible_on=20), trừ bản ghi đang được set
	if err := tx.Exec(
		`UPDATE seo_domain SET visible_on=20, updated_at=? WHERE deleted_at IS NULL AND visible_on=10 AND id<>?`,
		time.Now(), id,
	).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("demote current top planning news: %w", err)
	}

	// Đặt bản ghi mới thành top (visible_on=10)
	if err := tx.Exec(
		`UPDATE seo_domain SET visible_on=10, updated_at=? WHERE id=? AND deleted_at IS NULL`,
		time.Now(), id,
	).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("set planning news top: %w", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("commit set planning news top: %w", err)
	}

	detail, err := r.GetPlanningNews(ctx, strconv.FormatUint(id, 10))
	if err != nil {
		return nil, err
	}
	return &detail.SEODomain, nil
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

	if err := tx.Table("seo_content_category_assignment").
		Where("seo_domain_id=? AND deleted_at IS NULL AND is_primary=TRUE", id).
		Updates(map[string]any{"deleted_at": time.Now()}).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("clear primary planning news category: %w", err)
	}
	if err := tx.Exec(`
		INSERT INTO seo_content_category_assignment (
			seo_domain_id, category_id, is_primary,
			source_system, source_type, source_id,
			created_at, deleted_at
		) VALUES (?, ?, TRUE, 'admin', 'manual', ?, NOW(), NULL)
		ON CONFLICT (category_id, source_system, source_type, source_id)
		DO UPDATE SET
			seo_domain_id = EXCLUDED.seo_domain_id,
			is_primary = TRUE,
			created_at = NOW(),
			deleted_at = NULL
	`, id, foundCategoryID, strconv.FormatUint(id, 10)).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("assign planning news category: %w", err)
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
	case strings.HasPrefix(path, "/tin-tuc/bai-viet/"):
		return "news_article", "Tin quy hoạch"
	case strings.HasPrefix(path, "/tin-tuc/bao-cao/"):
		return "analysis_report", "Báo cáo và phân tích"
	case strings.HasPrefix(path, "/tin-tuc/do-an-quy-hoach/"):
		return "planning_project", "Đồ án quy hoạch"
	case strings.HasPrefix(path, "/tin-tuc/van-ban/"), strings.HasPrefix(path, "/tin-tuc/tai-lieu/"):
		return "legal_document", "Văn bản pháp lý"
	case strings.HasPrefix(path, "/tin-tuc/vung-quy-hoach/"):
		return "planning_region", "Vùng quy hoạch"
	case strings.HasPrefix(path, "/tin-tuc/lop/"):
		return "planning_layer", "Lớp quy hoạch"
	default:
		return "news_article", "Nội dung quy hoạch"
	}
}

func (r *Repository) SearchPublishedContent(ctx context.Context, query string, kinds []string, limit int) (*domain.ContentSearchResult, error) {
	pattern := likeKeyword(query)
	if pattern == "" {
		return &domain.ContentSearchResult{Groups: []domain.SearchGroup{}}, nil
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	query = strings.ToLower(strings.TrimSpace(strings.ReplaceAll(query, "+", " ")))
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
		LEFT JOIN seo_content_category_assignment a ON a.seo_domain_id=sd.id AND a.deleted_at IS NULL AND a.is_primary=TRUE
		LEFT JOIN seo_category c ON c.id=a.category_id AND c.deleted_at IS NULL
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
		item.Presentation.Subtitle = row.CategoryName
		item.Presentation.Description = row.Summary
		item.Source.System = "crm-service"
		item.Source.Dataset = "seo_domain"
		item.Source.Authority = "editorial"
		item.Source.DataQuality = "published"
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
	result.Interpretation.DetectedIntents = []string{"content_lookup", "free_text"}
	return result, nil
}
