package usecase

import (
	_errors "common/errors"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	seo_domain "crm/internal/domain/seo"
	"crm/internal/dto"
)

func (u *SeoPublicUsecase) GetSeoGovernanceDocument(
	ctx context.Context,
	moduleID string,
	slug string,
) (*dto.SeoGovernanceDocumentResponse, error) {
	moduleID = normalizeSeoModuleID(moduleID)
	if moduleID == "" || moduleID == "default" {
		return nil, _errors.ReturnError(400, "SEO module không được hỗ trợ")
	}

	normalizedSlug, err := normalizePublicSeoSlug(slug)
	if err != nil {
		return nil, err
	}

	canonicalPath, err := publicCanonicalPath(moduleID, normalizedSlug)
	if err != nil {
		return nil, err
	}

	page, err := u.resolvePublicDocumentPage(ctx, canonicalPath, normalizedSlug)
	if err != nil {
		return nil, err
	}
	if page == nil {
		return nil, _errors.ReturnError(404, "SEO page không tồn tại")
	}

	generatedAt := page.UpdatedAt.UTC()
	basis := "seo_domain.updated_at"
	if page.GeneratedAt != nil {
		generatedAt = page.GeneratedAt.UTC()
		basis = "seo_domain.generated_at"
	}

	publicDocument, warnings, err := composeSeoPublicDocument(page, generatedAt)
	if err != nil {
		return nil, _errors.ReturnError(503, err.Error())
	}
	if publicDocument.ModuleID != moduleID {
		return nil, _errors.ReturnError(503, fmt.Sprintf(
			"SEO module contract mismatch: expected %s, got %s",
			moduleID,
			publicDocument.ModuleID,
		))
	}
	if publicDocument.Identity.CanonicalPath != canonicalPath ||
		publicDocument.SEO.CanonicalPath != canonicalPath {
		return nil, _errors.ReturnError(503, "SEO canonical contract mismatch")
	}
	if publicDocument.AnswerReadiness == nil || publicDocument.IndexDecision == nil {
		return nil, _errors.ReturnError(503, "SEO governance decision is missing")
	}

	governance := buildSeoGovernanceDocument(
		page,
		publicDocument,
		basis,
		generatedAt,
		warnings,
	)
	payload, err := json.Marshal(governance)
	if err != nil {
		return nil, fmt.Errorf("encode SEO governance document: %w", err)
	}
	hash := sha256.Sum256(payload)
	etag := hex.EncodeToString(hash[:])
	governance.Evaluation.ETag = etag

	return &dto.SeoGovernanceDocumentResponse{
		Document:   governance,
		ETag:       etag,
		StatusCode: 200,
	}, nil
}

func buildSeoGovernanceDocument(
	page *seo_domain.SeoDomain,
	publicDocument *dto.SeoPublicDocument,
	evaluationBasis string,
	lastEvaluatedAt time.Time,
	warnings []string,
) *dto.SeoGovernanceDocument {
	factValues := make(map[string]string, len(publicDocument.Facts))
	for _, fact := range publicDocument.Facts {
		factValues[fact.ID] = strings.TrimSpace(fact.Value)
	}
	sectionValues := make(map[string]dto.SeoSection, len(publicDocument.Sections))
	for _, section := range publicDocument.Sections {
		sectionValues[section.ID] = section
	}

	requiredFacts := make([]dto.SeoGovernanceRequirementStatus, 0, len(publicDocument.AnswerReadiness.RequiredFacts))
	for _, id := range publicDocument.AnswerReadiness.RequiredFacts {
		value := factValues[id]
		requiredFacts = append(requiredFacts, dto.SeoGovernanceRequirementStatus{
			ID:      id,
			Present: isMeaningfulSeoValue(value),
			Value:   value,
		})
	}

	requiredSections := make([]dto.SeoGovernanceRequirementStatus, 0, len(publicDocument.AnswerReadiness.RequiredSections))
	for _, id := range publicDocument.AnswerReadiness.RequiredSections {
		requiredSections = append(requiredSections, dto.SeoGovernanceRequirementStatus{
			ID:      id,
			Present: hasMeaningfulSeoSection(sectionValues[id]),
		})
	}

	publishedAt := ""
	if page.PublishedAt != nil {
		publishedAt = page.PublishedAt.UTC().Format(time.RFC3339)
	}
	generatedAt := ""
	if page.GeneratedAt != nil {
		generatedAt = page.GeneratedAt.UTC().Format(time.RFC3339)
	}

	reviewState := "not_recorded"
	if strings.TrimSpace(page.Note) != "" {
		reviewState = "operator_note_present"
	}

	return &dto.SeoGovernanceDocument{
		SchemaVersion:              dto.SeoGovernanceDocumentSchemaVersion,
		ModuleID:                   publicDocument.ModuleID,
		PageFamily:                 publicDocument.ModuleID,
		SeoDomainID:                page.ID,
		EntityID:                   publicDocument.Identity.ID,
		CanonicalSlug:              publicDocument.Identity.Slug,
		CanonicalURL:               publicDocument.Identity.CanonicalPath,
		AnswerSpecificationVersion: publicDocument.AnswerReadiness.SpecificationVersion,
		AnswerReadiness:            publicDocument.AnswerReadiness,
		IndexDecision:              publicDocument.IndexDecision,
		RequiredFacts:              requiredFacts,
		RequiredSections:           requiredSections,
		Source:                     publicDocument.Source,
		Publication: dto.SeoGovernancePublication{
			PageStatus:      page.PageStatus,
			Published:       page.Published,
			PublishedAt:     publishedAt,
			Scope:           page.Scope,
			IndexConfigured: page.IsIndex,
			SitemapEnabled:  page.IsSiteMap,
			RenderStatus:    page.RenderStatus,
			NeedGenerate:    page.NeedGenerate,
			LastRenderError: strings.TrimSpace(page.LastRenderError),
		},
		Evaluation: dto.SeoGovernanceEvaluation{
			LastEvaluatedAt:   lastEvaluatedAt.UTC().Format(time.RFC3339),
			EvaluationBasis:   evaluationBasis,
			PublicDocumentVer: publicDocument.DocumentVersion,
			TemplateVersion:   page.TemplateVersion,
			RenderedHTMLHash:  page.StaticHtmlHash,
			GeneratedAt:       generatedAt,
			UpdatedAt:         page.UpdatedAt.UTC().Format(time.RFC3339),
		},
		Review: dto.SeoGovernanceReview{
			ReviewState:    reviewState,
			AuditState:     "not_available",
			OperatorNote:   strings.TrimSpace(page.Note),
			AuditAvailable: false,
		},
		OperatingSystem: publicDocument.OperatingSystem,
		ClientHandoff:   publicDocument.ClientHandoff,
		Measurement:     publicDocument.Measurement,
		GenerativeAI:    publicDocument.GenerativeAI,
		Warnings:        compactUniqueStrings(warnings),
	}
}

func (u *SeoPublicUsecase) GetSeoOperatingSystemManifest() *dto.SeoOperatingSystemManifest {
	return buildSeoOperatingSystemManifest()
}
