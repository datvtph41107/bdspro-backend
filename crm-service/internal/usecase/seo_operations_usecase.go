package usecase

import (
	_commonDto "common/domain/dto"
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	seo_domain "crm/internal/domain/seo"
	"crm/internal/dto"
	"crm/internal/enums"
)

func (u *SeoDomainUsecase) GetSeoOperationsDashboard(ctx context.Context, req *dto.SeoOperationsDashboardRequest) (*dto.SeoOperationsDashboardResponse, error) {
	limit := uint32(500)
	if req != nil && req.Limit > 0 {
		limit = req.Limit
	}
	if limit > 1000 {
		limit = 1000
	}

	items, total, err := u.seoDomainRepo.GetList(ctx, &dto.SeoDomainListRequest{
		Pagable: _commonDto.Pagable{Page: 0, Size: limit},
	})
	if err != nil {
		return nil, err
	}

	needGenerateItems, err := u.seoDomainRepo.GetNeedGenerateItems(ctx, 100)
	if err != nil {
		return nil, err
	}

	from, to := seoOperationsDateRange(req)
	trend := []dto.SeoOperationTrendPoint{}
	topQueries := []dto.SeoOperationTopQuery{}
	if u.seoTrendRepo != nil {
		trend, err = u.seoTrendRepo.GetDailySummary(ctx, from, to)
		if err != nil {
			return nil, err
		}
		topQueries, err = u.seoTrendRepo.GetTopPages(ctx, from, to, 10)
		if err != nil {
			return nil, err
		}
	}

	return buildSeoOperationsDashboard(items, total, needGenerateItems, trend, topQueries), nil
}

func seoOperationsDateRange(req *dto.SeoOperationsDashboardRequest) (time.Time, time.Time) {
	now := time.Now()
	days := 28
	if req != nil {
		switch strings.TrimSpace(req.Range) {
		case "7d":
			days = 7
		case "14d":
			days = 14
		case "90d":
			days = 90
		}
	}
	to := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	from := to.AddDate(0, 0, -(days - 1))
	return from, to
}

func buildSeoOperationsDashboard(items []seo_domain.SeoDomain, total int64, needGenerateItems []seo_domain.SeoDomain, trend []dto.SeoOperationTrendPoint, topQueries []dto.SeoOperationTopQuery) *dto.SeoOperationsDashboardResponse {
	published := 0
	draft := 0
	needGenerate := 0
	notIndexed := 0
	sourceIssues := 0
	qualitySum := 0
	sourceGroups := map[string]*dto.SeoSourceHealthItem{}

	for i := range items {
		item := &items[i]
		if item.Published {
			published++
		} else {
			draft++
		}
		if item.NeedGenerate {
			needGenerate++
		}
		if !item.IsIndex {
			notIndexed++
		}
		if item.SourceStatus == seo_domain.SeoSourceStatusMissing || item.SourceStatus == seo_domain.SeoSourceStatusStale || item.RefMissing {
			sourceIssues++
		}
		qualitySum += int(item.QualityScore)

		sourceKey := strings.TrimSpace(item.RefSource)
		if sourceKey == "" {
			sourceKey = seo_domain.SeoRefSourceManual
		}
		group := sourceGroups[sourceKey]
		if group == nil {
			group = &dto.SeoSourceHealthItem{
				Source:  sourceKey,
				Label:   sourceLabel(sourceKey),
				Service: sourceServiceFromResolver(sourceKey),
				Status:  seo_domain.SeoSourceStatusLinked,
			}
			sourceGroups[sourceKey] = group
		}
		group.TotalPages++
		if item.SourceStatus == seo_domain.SeoSourceStatusMissing || item.SourceStatus == seo_domain.SeoSourceStatusStale || item.RefMissing {
			group.NeedAttention++
			group.Status = item.SourceStatus
		}
	}

	avgQuality := 0
	if len(items) > 0 {
		avgQuality = qualitySum / len(items)
	}

	return &dto.SeoOperationsDashboardResponse{
		KPIs: []dto.SeoOperationKPI{
			{Label: "Trang SEO", Value: strconv.FormatInt(total, 10), Change: "đang quản lý", Tone: "good"},
			{Label: "Đã publish", Value: strconv.Itoa(published), Change: strconv.Itoa(draft) + " bản nháp", Tone: "good"},
			{Label: "Cần Render", Value: strconv.Itoa(needGenerate), Change: "dữ liệu đổi", Tone: toneByCount(needGenerate)},
			{Label: "Chặn Index", Value: strconv.Itoa(notIndexed), Change: "cần kiểm tra", Tone: toneByCount(notIndexed)},
			{Label: "Điểm chất lượng", Value: strconv.Itoa(avgQuality), Change: "trung bình", Tone: qualityTone(avgQuality)},
			{Label: "Nguồn có vấn đề", Value: strconv.Itoa(sourceIssues), Change: "thiếu/cũ", Tone: toneByCount(sourceIssues)},
		},
		Trend:         trend,
		Opportunities: buildSeoOpportunities(items),
		TopQueries:    topQueries,
		SystemEvents:  buildSeoEvents(items, needGenerate, sourceIssues, notIndexed),
		Requirements:  buildSeoRequirementGroups(items),
		RenderQueue:   buildSeoRenderQueue(needGenerateItems),
		IndexingRows:  buildSeoIndexingRows(items),
		SourceHealth:  flattenSourceHealth(sourceGroups),
	}
}

func buildSeoOpportunities(items []seo_domain.SeoDomain) []dto.SeoOperationOpportunity {
	opportunities := []dto.SeoOperationOpportunity{}
	for i := range items {
		item := items[i]
		issue, description, tone := seoPrimaryIssue(item)
		if issue == "" {
			continue
		}
		opportunities = append(opportunities, dto.SeoOperationOpportunity{
			ID:          fmt.Sprintf("%d", item.ID),
			Title:       firstNonEmpty(item.Title, item.Slug, item.CanonicalURL),
			Description: description,
			Issue:       issue,
			Impact:      impactByClassify(item.Classify),
			PageCount:   1,
			Owner:       sourceLabel(item.RefSource),
			Tone:        tone,
		})
	}
	sort.SliceStable(opportunities, func(i, j int) bool {
		return impactRank(opportunities[i].Impact) > impactRank(opportunities[j].Impact)
	})
	if len(opportunities) > 8 {
		return opportunities[:8]
	}
	return opportunities
}

func buildSeoEvents(items []seo_domain.SeoDomain, needGenerate int, sourceIssues int, notIndexed int) []dto.SeoOperationEvent {
	now := time.Now().Format(time.RFC3339)
	events := []dto.SeoOperationEvent{
		{ID: "seo-event-source", Title: "Nguồn dữ liệu cần kiểm tra", Description: fmt.Sprintf("%d trang có nguồn thiếu, cũ hoặc cần đồng bộ.", sourceIssues), CreatedAt: now, Tone: eventTone(sourceIssues)},
		{ID: "seo-event-render", Title: "Hàng đợi Render", Description: fmt.Sprintf("%d trang cần Render lại trước khi publish hoặc gửi sitemap.", needGenerate), CreatedAt: now, Tone: eventTone(needGenerate)},
		{ID: "seo-event-index", Title: "Indexing", Description: fmt.Sprintf("%d trang đang tắt Index hoặc bị loại khỏi sitemap.", notIndexed), CreatedAt: now, Tone: eventTone(notIndexed)},
	}
	if len(items) == 0 {
		events = append(events, dto.SeoOperationEvent{ID: "seo-event-empty", Title: "Chưa có trang SEO", Description: "Hãy tạo trang SEO đầu tiên từ nguồn dữ liệu hệ thống.", CreatedAt: now, Tone: "info"})
	}
	return events
}

func buildSeoRequirementGroups(items []seo_domain.SeoDomain) []dto.SeoRequirementGroup {
	total := len(items)
	return []dto.SeoRequirementGroup{
		{
			ID:          "content",
			Title:       "Nội dung trước khi publish",
			Description: "Tiêu đề, description, slug và chất lượng nội dung.",
			Items: []dto.SeoRequirementItem{
				requirementItem("seo-title", "SEO title và description đã đủ", "Trang cần có title và description để kết quả tìm kiếm rõ ràng.", total, countWhere(items, func(s seo_domain.SeoDomain) bool {
					return strings.TrimSpace(s.Title) != "" && strings.TrimSpace(s.Description) != ""
				}), "Đội nội dung", "Thiếu title/description làm CTR thấp và khó kiểm soát snippet."),
				requirementItem("seo-canonical", "Canonical URL hợp lệ", "Canonical giúp tránh trùng lặp nội dung giữa nhiều URL.", total, countWhere(items, func(s seo_domain.SeoDomain) bool { return strings.TrimSpace(s.CanonicalURL) != "" }), "Platform", "Canonical sai khiến tín hiệu SEO bị chia nhỏ."),
			},
		},
		{
			ID:          "source",
			Title:       "Nguồn dữ liệu hệ thống",
			Description: "Đảm bảo nội dung SEO dựa trên dữ liệu gốc và biết khi nào dữ liệu thay đổi.",
			Items: []dto.SeoRequirementItem{
				requirementItem("source-linked", "Nguồn chính đã liên kết", "Trang có nguồn dữ liệu rõ ràng hoặc được đánh dấu thủ công.", total, countWhere(items, func(s seo_domain.SeoDomain) bool {
					return s.RefSource != "" && s.SourceStatus != seo_domain.SeoSourceStatusMissing
				}), "Data Ops", "Admin không cần thấy refType/refId, nhưng backend phải giữ binding này."),
				requirementItem("source-fresh", "Snapshot dữ liệu còn mới", "Trang không bị đánh dấu dữ liệu cũ hoặc cần Render lại.", total, countWhere(items, func(s seo_domain.SeoDomain) bool {
					return !s.NeedGenerate && s.SourceStatus != seo_domain.SeoSourceStatusStale && !s.RefMissing
				}), "Data Ops", "Dữ liệu cũ phải chặn publish hoặc đưa vào hàng đợi Render."),
			},
		},
		{
			ID:          "technical",
			Title:       "Indexing, sitemap và Render",
			Description: "Các công tắc quyết định trang có thể được Google đọc và đưa vào Index hay không.",
			Items: []dto.SeoRequirementItem{
				requirementItem("render-clean", "Render đã hoàn tất", "Trang không còn nằm trong hàng đợi Render.", total, countWhere(items, func(s seo_domain.SeoDomain) bool { return !s.NeedGenerate }), "Publisher", "Render chưa xong thì chưa nên publish nội dung public."),
				requirementItem("index-enabled", "Indexing đang bật", "Trang public cần bật Index và nằm trong sitemap.", total, countWhere(items, func(s seo_domain.SeoDomain) bool { return s.IsIndex && s.IsSiteMap }), "Platform", "Không bật Index/sitemap thì Google khó phát hiện URL."),
			},
		},
	}
}

func buildSeoRenderQueue(items []seo_domain.SeoDomain) []dto.SeoRenderQueueItem {
	rows := make([]dto.SeoRenderQueueItem, 0, len(items))
	for i := range items {
		item := items[i]
		rows = append(rows, dto.SeoRenderQueueItem{
			ID:          fmt.Sprintf("render-%d", item.ID),
			SeoDomainID: item.ID,
			Title:       firstNonEmpty(item.Title, item.Slug, item.CanonicalURL),
			Type:        "Render",
			Status:      "Queued",
			Worker:      "seo-render",
			Duration:    "-",
			Priority:    priorityByClassify(item.Classify),
		})
	}
	return rows
}

func buildSeoIndexingRows(items []seo_domain.SeoDomain) []dto.SeoIndexingRow {
	rows := []dto.SeoIndexingRow{}
	for i := range items {
		item := items[i]
		url := firstNonEmpty(item.CanonicalURL, item.OriginURL, item.Slug)
		if url == "" {
			continue
		}
		state := "Indexed"
		if !item.IsIndex {
			state = "Blocked"
		} else if !item.Published {
			state = "Waiting"
		}
		sitemap := "Included"
		if !item.IsSiteMap {
			sitemap = "Excluded"
		}
		rows = append(rows, dto.SeoIndexingRow{
			URL:       url,
			Canonical: "Self",
			Robots:    robotsState(item),
			Sitemap:   sitemap,
			State:     state,
			LastPing:  formatSeoOperationTime(item.SiteMapLastedAt, item.UpdatedAt),
		})
	}
	if len(rows) > 100 {
		return rows[:100]
	}
	return rows
}

func flattenSourceHealth(groups map[string]*dto.SeoSourceHealthItem) []dto.SeoSourceHealthItem {
	items := make([]dto.SeoSourceHealthItem, 0, len(groups))
	for _, item := range groups {
		if item.NeedAttention == 0 {
			item.Status = seo_domain.SeoSourceStatusLinked
		}
		items = append(items, *item)
	}
	sort.SliceStable(items, func(i, j int) bool {
		if items[i].NeedAttention == items[j].NeedAttention {
			return items[i].Source < items[j].Source
		}
		return items[i].NeedAttention > items[j].NeedAttention
	})
	return items
}

func requirementItem(id, title, description string, total int, passed int, owner string, hint string) dto.SeoRequirementItem {
	coverage := uint32(100)
	if total > 0 {
		coverage = uint32((passed * 100) / total)
	}
	tone := "good"
	status := "Đạt"
	if coverage < 70 {
		tone = "bad"
		status = "Cần xử lý"
	} else if coverage < 90 {
		tone = "warn"
		status = "Cần duyệt"
	}
	return dto.SeoRequirementItem{ID: id, Title: title, Description: description, Status: status, Tone: tone, Owner: owner, Coverage: coverage, Hint: hint}
}

func countWhere(items []seo_domain.SeoDomain, fn func(seo_domain.SeoDomain) bool) int {
	count := 0
	for _, item := range items {
		if fn(item) {
			count++
		}
	}
	return count
}

func seoPrimaryIssue(item seo_domain.SeoDomain) (string, string, string) {
	if item.SourceStatus == seo_domain.SeoSourceStatusMissing || item.RefMissing {
		return "stale_source", "Nguồn dữ liệu thiếu hoặc không còn truy cập được.", "bad"
	}
	if item.SourceStatus == seo_domain.SeoSourceStatusStale {
		return "stale_source", "Nguồn dữ liệu đã thay đổi, cần kiểm tra snapshot và Render lại.", "warn"
	}
	if item.NeedGenerate {
		return "need_generate", "Nội dung cần Render lại trước khi publish hoặc gửi sitemap.", "warn"
	}
	if !item.IsIndex {
		return "not_indexed", "Trang đang tắt Index hoặc bị loại khỏi sitemap.", "bad"
	}
	if strings.TrimSpace(item.Title) == "" {
		return "missing_title", "Trang thiếu title nên khó kiểm soát kết quả tìm kiếm.", "bad"
	}
	if strings.TrimSpace(item.Description) == "" {
		return "missing_description", "Trang thiếu description nên snippet có thể kém hấp dẫn.", "warn"
	}
	return "", "", ""
}

func sourceLabel(source string) string {
	source = strings.TrimSpace(source)
	if source == "" || source == seo_domain.SeoRefSourceManual {
		return "Thủ công"
	}
	for _, cfg := range seo_domain.SeoRefTypeConfigs() {
		if cfg.ResolverKey == source {
			return cfg.Label
		}
	}
	return source
}

func toneByCount(count int) string {
	if count > 0 {
		return "warning"
	}
	return "good"
}

func qualityTone(score int) string {
	if score >= 85 {
		return "good"
	}
	if score >= 65 {
		return "warning"
	}
	return "danger"
}

func eventTone(count int) string {
	if count > 0 {
		return "warning"
	}
	return "success"
}

func impactByClassify(classify enums.ESeoClassify) string {
	switch classify {
	case enums.ESeoClassifyA:
		return "high"
	case enums.ESeoClassifyB, enums.ESeoClassifyC:
		return "medium"
	default:
		return "low"
	}
}

func impactRank(impact string) int {
	switch impact {
	case "high":
		return 3
	case "medium":
		return 2
	default:
		return 1
	}
}

func priorityByClassify(classify enums.ESeoClassify) string {
	switch classify {
	case enums.ESeoClassifyA:
		return "P0"
	case enums.ESeoClassifyB:
		return "P1"
	case enums.ESeoClassifyC:
		return "P2"
	default:
		return "P3"
	}
}

func percent(a uint64, b uint64) float64 {
	if b == 0 {
		return 0
	}
	return float64(a) * 100 / float64(b)
}

func avgPositionFromQuality(score uint32) float64 {
	if score >= 90 {
		return 4.2
	}
	if score >= 75 {
		return 8.8
	}
	if score >= 50 {
		return 14.6
	}
	return 22.4
}

func robotsState(item seo_domain.SeoDomain) string {
	if !item.IsIndex {
		return "noindex,follow"
	}
	if item.IsRobot != nil && !*item.IsRobot {
		return "noindex,nofollow"
	}
	return "index,follow"
}

func formatSeoOperationTime(primary *time.Time, fallback time.Time) string {
	if primary != nil && !primary.IsZero() {
		return primary.Format(time.RFC3339)
	}
	if !fallback.IsZero() {
		return fallback.Format(time.RFC3339)
	}
	return ""
}
