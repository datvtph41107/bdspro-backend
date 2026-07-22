package contentcatalog

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"crm/internal/modules/contentcatalog/domain"

	"gorm.io/gorm"
)

type Repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) *Repository { return &Repository{db: db} }

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
		       COALESCE(c.description,'') AS description, c.sort_order
		FROM seo_content_category c
		JOIN seo_content_domain d ON d.id = c.domain_id AND d.deleted_at IS NULL
		WHERE c.deleted_at IS NULL AND c.status = 'active' AND d.status = 'active'`
	args := []any{}
	if strings.TrimSpace(domainCode) != "" {
		query += " AND d.code = ?"
		args = append(args, strings.TrimSpace(domainCode))
	}
	query += " ORDER BY d.sort_order, c.sort_order, c.name, c.id"
	if err := r.db.WithContext(ctx).Raw(query, args...).Scan(&items).Error; err != nil {
		return nil, fmt.Errorf("list content categories: %w", err)
	}
	return items, nil
}

func (r *Repository) ListFeed(ctx context.Context, feed, domainCode, categoryCode string, limit int) (*domain.FeedResult, error) {
	feed = strings.ToLower(strings.TrimSpace(feed))
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	switch feed {
	case "featured", "attention", "breaking", "pinned":
		return r.listEditorialFeed(ctx, feed, domainCode, categoryCode, limit)
	case "latest":
		return r.listLatestFeed(ctx, domainCode, categoryCode, limit)
	case "most_read", "trending":
		return &domain.FeedResult{
			Feed:     feed,
			Computed: true,
			Items:    []domain.ContentFeedItem{},
			Warnings: []string{"Feed cần dữ liệu analytics; không được giả lập bằng visible_on hoặc placement cố định."},
		}, nil
	default:
		return nil, fmt.Errorf("unsupported content feed %q", feed)
	}
}

func (r *Repository) listEditorialFeed(ctx context.Context, feed, domainCode, categoryCode string, limit int) (*domain.FeedResult, error) {
	items := []domain.ContentFeedItem{}
	query := `
		SELECT sd.id AS seo_page_id, sd.slug, COALESCE(sd.title,'') AS title,
		       COALESCE(sd.description,'') AS description, COALESCE(sd.summary,'') AS summary,
		       sd.canonical_url, sd.published_at, sd.updated_at,
		       COALESCE(c.code,'') AS category_code, COALESCE(c.name,'') AS category_name,
		       p.placement_code, p.priority
		FROM seo_content_placement p
		JOIN seo_domain sd ON sd.id = p.seo_domain_id
		LEFT JOIN seo_content_category_assignment a
		  ON a.seo_domain_id = sd.id AND a.deleted_at IS NULL AND a.is_primary = TRUE
		LEFT JOIN seo_content_category c ON c.id = a.category_id AND c.deleted_at IS NULL
		LEFT JOIN seo_content_domain d ON d.id = c.domain_id AND d.deleted_at IS NULL
		WHERE p.deleted_at IS NULL
		  AND p.status = 'active'
		  AND p.placement_code = ?
		  AND (p.starts_at IS NULL OR p.starts_at <= NOW())
		  AND (p.ends_at IS NULL OR p.ends_at > NOW())
		  AND sd.deleted_at IS NULL AND sd.published = TRUE AND sd.scope = 'public'`
	args := []any{feed}
	query, args = addTaxonomyFilter(query, args, domainCode, categoryCode)
	query += " ORDER BY p.priority DESC, sd.published_at DESC NULLS LAST, sd.id DESC LIMIT ?"
	args = append(args, limit)
	if err := r.db.WithContext(ctx).Raw(query, args...).Scan(&items).Error; err != nil {
		return nil, fmt.Errorf("list editorial feed: %w", err)
	}
	return &domain.FeedResult{Feed: feed, Computed: false, Items: items}, nil
}

func (r *Repository) listLatestFeed(ctx context.Context, domainCode, categoryCode string, limit int) (*domain.FeedResult, error) {
	items := []domain.ContentFeedItem{}
	query := `
		SELECT sd.id AS seo_page_id, sd.slug, COALESCE(sd.title,'') AS title,
		       COALESCE(sd.description,'') AS description, COALESCE(sd.summary,'') AS summary,
		       sd.canonical_url, sd.published_at, sd.updated_at,
		       COALESCE(c.code,'') AS category_code, COALESCE(c.name,'') AS category_name,
		       '' AS placement_code, 0 AS priority
		FROM seo_domain sd
		LEFT JOIN seo_content_category_assignment a
		  ON a.seo_domain_id = sd.id AND a.deleted_at IS NULL AND a.is_primary = TRUE
		LEFT JOIN seo_content_category c ON c.id = a.category_id AND c.deleted_at IS NULL
		LEFT JOIN seo_content_domain d ON d.id = c.domain_id AND d.deleted_at IS NULL
		WHERE sd.deleted_at IS NULL AND sd.published = TRUE AND sd.scope = 'public'`
	args := []any{}
	query, args = addTaxonomyFilter(query, args, domainCode, categoryCode)
	query += " ORDER BY sd.published_at DESC NULLS LAST, sd.updated_at DESC, sd.id DESC LIMIT ?"
	args = append(args, limit)
	if err := r.db.WithContext(ctx).Raw(query, args...).Scan(&items).Error; err != nil {
		return nil, fmt.Errorf("list latest feed: %w", err)
	}
	return &domain.FeedResult{Feed: "latest", Computed: true, Items: items}, nil
}

func addTaxonomyFilter(query string, args []any, domainCode, categoryCode string) (string, []any) {
	if value := strings.TrimSpace(domainCode); value != "" {
		query += " AND d.code = ?"
		args = append(args, value)
	}
	if value := strings.TrimSpace(categoryCode); value != "" {
		query += " AND c.code = ?"
		args = append(args, value)
	}
	return query, args
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
	CategoryCode string
	CategoryName string
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
		       COALESCE(c.code,'') AS category_code, COALESCE(c.name,'') AS category_name,
		       CASE WHEN LOWER(COALESCE(sd.title,'')) = ? THEN 100
		            WHEN LOWER(COALESCE(sd.title,'')) LIKE ? THEN 70 ELSE 35 END::float8 AS score
		FROM seo_domain sd
		LEFT JOIN seo_content_category_assignment a ON a.seo_domain_id=sd.id AND a.deleted_at IS NULL AND a.is_primary=TRUE
		LEFT JOIN seo_content_category c ON c.id=a.category_id AND c.deleted_at IS NULL
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
		item.Attributes = map[string]any{"categoryCode": row.CategoryCode, "categoryName": row.CategoryName, "deepLink": row.DeepLink, "publishedAt": row.PublishedAt}
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
