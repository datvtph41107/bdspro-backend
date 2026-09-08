package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	_errors "common/errors"
	seo_domain "crm/internal/domain/seo"
	"crm/internal/dto"
	provider_interface "crm/internal/interface/provider"
	"crm/internal/repo"

	"gorm.io/datatypes"
)

type SeoRenderUsecase struct {
	seoDomainRepo        repo.SeoDomainRepo
	seoGenerationLogRepo repo.SeoGenerationLogRepo
	tqdProvider          provider_interface.TqdProvider
}

func NewSeoRenderUsecase(
	seoDomainRepo repo.SeoDomainRepo,
	seoGenerationLogRepo repo.SeoGenerationLogRepo,
	tqdProvider provider_interface.TqdProvider,
) *SeoRenderUsecase {
	return &SeoRenderUsecase{
		seoDomainRepo:        seoDomainRepo,
		seoGenerationLogRepo: seoGenerationLogRepo,
		tqdProvider:          tqdProvider,
	}
}

func (u *SeoRenderUsecase) PreviewSeoPage(ctx context.Context, req *dto.RenderSeoPreviewRequest) (*dto.RenderSeoPreviewResponse, error) {
	if req == nil {
		return nil, _errors.ReturnError(400, "request không hợp lệ")
	}

	warnings := validatePreviewWarnings(req.Title, req.Description, req.CanonicalURL, req.Content)

	htmlBody, documentWarnings, err := renderSeoHTMLWithWarnings(renderSeoInput{
		Slug:         strings.TrimSpace(req.Slug),
		Title:        strings.TrimSpace(req.Title),
		Description:  strings.TrimSpace(req.Description),
		CanonicalURL: strings.TrimSpace(req.CanonicalURL),
		Content:      req.Content,
		Summary:      req.Summary,
		Metadata:     req.Metadata,
		TemplateKey:  strings.TrimSpace(req.TemplateKey),
		ModuleID:     firstNonEmptyRender(strings.TrimSpace(req.TemplateKey), "default"),
		Environment:  "preview",
		Indexable:    false,
		GeneratedAt:  time.Now().UTC(),
	})
	if err != nil {
		return nil, err
	}
	warnings = append(warnings, documentWarnings...)

	hash := hashRenderedHTML(htmlBody)
	staticPath := buildStaticHTMLPath(0, req.Slug, hash)

	return &dto.RenderSeoPreviewResponse{
		HTML:       htmlBody,
		StaticPath: staticPath,
		Hash:       hash,
		Warnings:   warnings,
	}, nil
}

func (u *SeoRenderUsecase) PreviewExistingSeoPage(ctx context.Context, req *dto.RenderExistingSeoPreviewRequest) (*dto.RenderSeoPreviewResponse, error) {
	if req == nil || req.ID == 0 {
		return nil, _errors.ReturnError(400, "id không hợp lệ")
	}

	page, err := u.seoDomainRepo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if page == nil {
		return nil, _errors.ReturnError(404, "seo page không tồn tại")
	}

	content := page.Content
	if req.Content != nil {
		content = *req.Content
	}

	templateKey := page.TemplateKey
	if req.TemplateKey != nil {
		templateKey = strings.TrimSpace(*req.TemplateKey)
	}

	previewReq := &dto.RenderSeoPreviewRequest{
		Slug:         page.Slug,
		Title:        page.Title,
		Description:  page.Description,
		CanonicalURL: page.CanonicalURL,
		Content:      content,
		Summary:      page.Summary,
		Metadata:     string(page.Metadata),
		TemplateKey:  templateKey,
	}

	resp, err := u.PreviewSeoPage(ctx, previewReq)
	if err != nil {
		return nil, err
	}

	resp.StaticPath = buildStaticHTMLPath(page.ID, page.Slug, resp.Hash)
	return resp, nil
}

func (u *SeoRenderUsecase) GenerateSeoPage(ctx context.Context, req *dto.GenerateSeoPageRequest) (*dto.GenerateSeoPageResponse, error) {
	if req == nil || req.ID == 0 {
		return nil, _errors.ReturnError(400, "id không hợp lệ")
	}

	page, err := u.seoDomainRepo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	if page == nil {
		return nil, _errors.ReturnError(404, "seo page không tồn tại")
	}

	if err := validatePageCanGenerate(page); err != nil {
		return nil, err
	}

	triggerType := normalizeGenerationTrigger(req.TriggerType)
	metadata := buildGenerationMetadata(req.Reason)
	logEntry := seo_domain.NewSeoGenerationLog(page.ID, triggerType, metadata)

	logEntry, err = u.seoGenerationLogRepo.Create(ctx, logEntry)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	if err := u.seoGenerationLogRepo.MarkRunning(ctx, logEntry.ID, now); err != nil {
		return nil, err
	}
	if err := u.hydrateSeoSource(ctx, page, now); err != nil {
		failedAt := time.Now().UTC()
		_ = u.seoGenerationLogRepo.MarkFailed(ctx, logEntry.ID, err.Error(), failedAt)
		_ = u.seoDomainRepo.MarkRenderFailed(ctx, page.ID, err.Error())
		return &dto.GenerateSeoPageResponse{
			ID:           page.ID,
			RenderStatus: seo_domain.SeoRenderStatusFailed,
			ErrorMessage: err.Error(),
		}, nil
	}
	if err := u.seoDomainRepo.MarkRenderRunning(ctx, page.ID); err != nil {
		_ = u.seoGenerationLogRepo.MarkFailed(ctx, logEntry.ID, err.Error(), time.Now().UTC())
		return nil, err
	}

	htmlBody, err := renderSeoHTML(renderSeoInput{
		Page:         page,
		Slug:         page.Slug,
		Title:        page.Title,
		Description:  page.Description,
		CanonicalURL: page.CanonicalURL,
		Content:      page.Content,
		Summary:      page.Summary,
		Metadata:     string(page.Metadata),
		TemplateKey:  page.TemplateKey,
		ModuleID:     firstNonEmptyRender(page.TemplateKey, page.GetSchema(), "default"),
		Environment:  "crm-rendered",
		Indexable:    page.IsPublicIndexable(),
		GeneratedAt:  now,
	})
	if err != nil {
		failedAt := time.Now().UTC()
		_ = u.seoGenerationLogRepo.MarkFailed(ctx, logEntry.ID, err.Error(), failedAt)
		_ = u.seoDomainRepo.MarkRenderFailed(ctx, page.ID, err.Error())
		return &dto.GenerateSeoPageResponse{
			ID:           page.ID,
			RenderStatus: seo_domain.SeoRenderStatusFailed,
			ErrorMessage: err.Error(),
		}, nil
	}

	hash := hashRenderedHTML(htmlBody)
	staticPath := buildStaticHTMLPath(page.ID, page.Slug, hash)
	finishedAt := time.Now().UTC()

	if err := u.seoDomainRepo.MarkRenderSuccess(ctx, page.ID, htmlBody, staticPath, hash, finishedAt); err != nil {
		_ = u.seoGenerationLogRepo.MarkFailed(ctx, logEntry.ID, err.Error(), finishedAt)
		_ = u.seoDomainRepo.MarkRenderFailed(ctx, page.ID, err.Error())
		return nil, err
	}

	if err := u.seoGenerationLogRepo.MarkSuccess(ctx, logEntry.ID, staticPath, hash, finishedAt); err != nil {
		return nil, err
	}

	return &dto.GenerateSeoPageResponse{
		ID:             page.ID,
		RenderStatus:   seo_domain.SeoRenderStatusSuccess,
		StaticHTMLPath: staticPath,
		StaticHTMLHash: hash,
		GeneratedAt:    finishedAt.Format(time.RFC3339),
	}, nil
}

func (u *SeoRenderUsecase) GetSeoGenerationLogs(ctx context.Context, req *dto.SeoGenerationLogListRequest) ([]seo_domain.SeoGenerationLog, int64, error) {
	if req == nil {
		req = &dto.SeoGenerationLogListRequest{}
	}
	req.Normalize()
	return u.seoGenerationLogRepo.GetList(ctx, req)
}

type renderSeoInput struct {
	Page         *seo_domain.SeoDomain
	Slug         string
	Title        string
	Description  string
	CanonicalURL string
	Content      string
	Summary      string
	Metadata     string
	TemplateKey  string
	ModuleID     string
	Environment  string
	Indexable    bool
	GeneratedAt  time.Time
}

func renderSeoHTML(input renderSeoInput) (string, error) {
	htmlBody, _, err := renderSeoHTMLWithWarnings(input)
	return htmlBody, err
}

func renderSeoHTMLWithWarnings(input renderSeoInput) (string, []string, error) {
	generatedAt := input.GeneratedAt
	if generatedAt.IsZero() {
		generatedAt = time.Now().UTC()
	}

	page := input.Page
	if page == nil {
		metadata := datatypes.JSON([]byte(strings.TrimSpace(input.Metadata)))
		if len(metadata) == 0 || !json.Valid(metadata) {
			metadata = datatypes.JSON([]byte("{}"))
		}
		page = &seo_domain.SeoDomain{
			Slug:         strings.TrimSpace(input.Slug),
			CanonicalURL: strings.TrimSpace(input.CanonicalURL),
			Title:        strings.TrimSpace(input.Title),
			Description:  strings.TrimSpace(input.Description),
			Content:      input.Content,
			Summary:      input.Summary,
			TemplateKey:  firstNonEmptyRender(input.TemplateKey, input.ModuleID, "default"),
			Metadata:     metadata,
			Scope:        seo_domain.SeoScopePublic,
			PageStatus:   seo_domain.SeoPageStatusDraft,
			RenderStatus: seo_domain.SeoRenderStatusNone,
			Published:    false,
			IsIndex:      input.Indexable,
		}
	}

	document, warnings, err := composeSeoPublicDocument(page, generatedAt)
	if err != nil {
		return "", warnings, err
	}

	// Explicit render input controls preview/runtime delivery details. Domain
	// fields remain authoritative for persisted generation.
	if strings.TrimSpace(input.Title) != "" {
		document.Heading.Title = strings.TrimSpace(input.Title)
		document.SEO.Title = strings.TrimSpace(input.Title)
	}
	if strings.TrimSpace(input.Description) != "" {
		document.SEO.Description = strings.TrimSpace(input.Description)
	}
	if strings.TrimSpace(input.Summary) != "" {
		document.Heading.Summary = strings.TrimSpace(input.Summary)
	}
	if strings.TrimSpace(input.CanonicalURL) != "" {
		document.Identity.CanonicalPath = canonicalPathFromURL(input.CanonicalURL, input.Slug)
		document.SEO.CanonicalPath = document.Identity.CanonicalPath
	}
	if input.Page == nil {
		document.SEO.Indexable = false
	}

	htmlBody, err := renderSeoPublicDocumentHTML(document, seoDocumentRenderOptions{
		CanonicalURL: firstNonEmptyRender(input.CanonicalURL, page.CanonicalURL, document.SEO.CanonicalPath),
		TemplateKey:  firstNonEmptyRender(input.TemplateKey, page.TemplateKey, document.ModuleID),
		ModuleID:     firstNonEmptyRender(input.ModuleID, document.ModuleID),
		Environment:  firstNonEmptyRender(input.Environment, "crm-rendered"),
		GeneratedAt:  generatedAt,
	})
	return htmlBody, warnings, err
}

func validatePreviewWarnings(title string, description string, canonicalURL string, content string) []string {
	warnings := make([]string, 0)

	if strings.TrimSpace(title) == "" {
		warnings = append(warnings, "Thiếu title")
	}
	if len([]rune(strings.TrimSpace(title))) > 70 {
		warnings = append(warnings, "Title dài hơn 70 ký tự")
	}
	if strings.TrimSpace(description) == "" {
		warnings = append(warnings, "Thiếu meta description")
	}
	if len([]rune(strings.TrimSpace(description))) > 170 {
		warnings = append(warnings, "Meta description dài hơn 170 ký tự")
	}
	if strings.TrimSpace(canonicalURL) == "" {
		warnings = append(warnings, "Thiếu canonicalUrl, hệ thống sẽ fallback theo slug")
	}
	if strings.TrimSpace(content) == "" {
		warnings = append(warnings, "Thiếu content")
	}

	return warnings
}

func validatePageCanGenerate(page *seo_domain.SeoDomain) error {
	if page == nil {
		return _errors.ReturnError(404, "seo page không tồn tại")
	}

	page.NormalizeLifecycle()

	if page.PageStatus != seo_domain.SeoPageStatusPublished || !page.Published {
		return _errors.ReturnError(409, "chỉ generate SEO page đã published")
	}
	if strings.TrimSpace(page.Title) == "" {
		return _errors.ReturnError(422, "title là bắt buộc để generate")
	}
	if strings.TrimSpace(page.CanonicalURL) == "" {
		return _errors.ReturnError(422, "canonicalUrl là bắt buộc để generate")
	}
	return nil
}

func normalizeGenerationTrigger(triggerType string) string {
	switch strings.TrimSpace(triggerType) {
	case seo_domain.SeoGenerationTriggerManual,
		seo_domain.SeoGenerationTriggerScheduler,
		seo_domain.SeoGenerationTriggerSourceSync,
		seo_domain.SeoGenerationTriggerPublish:
		return strings.TrimSpace(triggerType)
	default:
		return seo_domain.SeoGenerationTriggerManual
	}
}

func buildGenerationMetadata(reason string) datatypes.JSON {
	payload := map[string]string{}
	if strings.TrimSpace(reason) != "" {
		payload["reason"] = strings.TrimSpace(reason)
	}
	if len(payload) == 0 {
		return datatypes.JSON([]byte("{}"))
	}

	raw, err := json.Marshal(payload)
	if err != nil {
		return datatypes.JSON([]byte("{}"))
	}
	return datatypes.JSON(raw)
}

func buildMetadataScript(metadata string) string {
	metadata = strings.TrimSpace(metadata)
	if metadata == "" || metadata == "null" {
		return ""
	}

	var payload any
	if err := json.Unmarshal([]byte(metadata), &payload); err != nil {
		return ""
	}

	normalized, err := json.Marshal(payload)
	if err != nil {
		return ""
	}

	// encoding/json escapes <, > and &, which prevents a JSON value from
	// prematurely closing the script element while keeping valid JSON-LD.
	return `<script type="application/ld+json">` + string(normalized) + `</script>`
}

func sanitizeTrustedContent(content string) string {
	content = strings.TrimSpace(content)
	if content == "" {
		return "<p></p>"
	}
	return content
}

func hashRenderedHTML(htmlBody string) string {
	sum := sha256.Sum256([]byte(htmlBody))
	return hex.EncodeToString(sum[:])
}

func buildStaticHTMLPath(id uint64, slug string, hash string) string {
	slug = strings.Trim(strings.TrimSpace(slug), "/")
	if slug == "" {
		if id > 0 {
			slug = fmt.Sprintf("seo-page-%d", id)
		} else {
			slug = "preview"
		}
	}

	shortHash := hash
	if len(shortHash) > 12 {
		shortHash = shortHash[:12]
	}

	if id > 0 {
		return fmt.Sprintf("/seo-static/%d/%s-%s.html", id, slug, shortHash)
	}
	return fmt.Sprintf("/seo-static/preview/%s-%s.html", slug, shortHash)
}

func firstNonEmptyRender(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
