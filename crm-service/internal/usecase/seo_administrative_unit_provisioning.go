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
	"crm/internal/enums"
	provider_interface "crm/internal/interface/provider"
	"crm/internal/repo"

	"gorm.io/datatypes"
)

const administrativeUnitProvisioningTimeout = 20 * time.Second

type SeoAdministrativeUnitProvisioningUsecase struct {
	seoDomainRepo repo.SeoDomainRepo
	tqdProvider   provider_interface.TqdProvider
	renderUsecase *SeoRenderUsecase
}

func NewSeoAdministrativeUnitProvisioningUsecase(
	seoDomainRepo repo.SeoDomainRepo,
	tqdProvider provider_interface.TqdProvider,
	renderUsecase *SeoRenderUsecase,
) *SeoAdministrativeUnitProvisioningUsecase {
	return &SeoAdministrativeUnitProvisioningUsecase{
		seoDomainRepo: seoDomainRepo,
		tqdProvider:   tqdProvider,
		renderUsecase: renderUsecase,
	}
}

func (u *SeoAdministrativeUnitProvisioningUsecase) Provision(
	ctx context.Context,
	req *dto.ProvisionAdministrativeUnitSeoRequest,
) (*dto.ProvisionAdministrativeUnitSeoResponse, error) {
	if u == nil || u.seoDomainRepo == nil || u.tqdProvider == nil {
		return nil, _errors.ReturnError(503, "Administrative Unit provisioning chưa sẵn sàng")
	}
	if req == nil {
		return nil, _errors.ReturnError(400, "request không hợp lệ")
	}

	identity, err := normalizeAdministrativeUnitProvisioningIdentity(req.Identity)
	if err != nil {
		return nil, err
	}

	projectionCtx, cancel := context.WithTimeout(ctx, administrativeUnitProvisioningTimeout)
	defer cancel()

	projection, err := u.tqdProvider.GetAdministrativeUnitProjection(projectionCtx, identity)
	if err != nil {
		return nil, _errors.ReturnError(503, fmt.Sprintf("Không thể tải Administrative Unit projection: %v", err))
	}
	if err := validateAdministrativeUnitProvisioningProjection(projection); err != nil {
		return nil, err
	}

	document := buildAdministrativeUnitPublicDocument(projection)
	warnings := make([]string, 0, 8)
	normalizeSeoPublicDocument(document, &warnings)

	// Evaluate source/content eligibility independently from the operator's
	// publication choice. Publication controls whether the page is public;
	// Answer Readiness and Index Decision remain backend-owned policy outputs.
	document.SEO.Indexable = true
	applySeoAnswerReadiness(document)
	readiness := cloneSeoAnswerReadiness(document.AnswerReadiness)
	indexDecision := cloneSeoIndexDecision(document.IndexDecision)
	if readiness == nil || indexDecision == nil {
		return nil, _errors.ReturnError(503, "Backend không tạo được Answer Readiness hoặc Index Decision")
	}

	projectionJSON, err := json.Marshal(projection)
	if err != nil {
		return nil, fmt.Errorf("encode Administrative Unit projection: %w", err)
	}

	// Persist source semantics and V2.U client/measurement contracts, but do not
	// persist a stale policy decision. composeSeoPublicDocument recalculates
	// readiness/indexability against the current domain lifecycle on every read.
	document.AnswerReadiness = nil
	document.IndexDecision = nil
	document.SEO = dto.SeoPublicSEO{}
	documentJSON, err := json.Marshal(document)
	if err != nil {
		return nil, fmt.Errorf("encode Administrative Unit PublicSeoDocument: %w", err)
	}

	refID := administrativeUnitReferenceID(projection)
	if refID == 0 {
		return nil, _errors.ReturnError(503, "Administrative Unit projection thiếu administrativeCode dạng số")
	}

	canonicalPath := strings.TrimSpace(projection.CanonicalPath)
	existing, err := u.findExistingAdministrativeUnitPage(ctx, projection.EntityID, canonicalPath, projection.Slug)
	if err != nil {
		return nil, err
	}
	if existingWarnings, ownershipErr := validateExistingAdministrativeUnitOwnership(existing, refID, projection.EntityID); ownershipErr != nil {
		return nil, ownershipErr
	} else {
		warnings = append(warnings, existingWarnings...)
	}
	if existing != nil && canonicalPathFromURL(existing.CanonicalURL, existing.Slug) != canonicalPath {
		return nil, _errors.ReturnError(409, "Canonical Administrative Unit đã thay đổi; cần tạo redirect mapping trước khi đồng bộ")
	}

	published, publishedAt := resolveAdministrativeUnitPublication(existing, req.Publish)
	indexable := published && indexDecision.Eligible
	if !published {
		warnings = append(warnings, "candidate_policy_evaluation_before_publication")
	}
	now := time.Now().UTC()
	if published && publishedAt == nil {
		publishedAt = &now
	}
	if !published {
		publishedAt = nil
	}

	snapshotHashBytes := sha256.Sum256(projectionJSON)
	snapshotHash := hex.EncodeToString(snapshotHashBytes[:])
	qualityScore := planningProjectQualityScore(projection.Source.DataQuality)
	isRobot := true
	mapPath := ""
	if document.MapTarget != nil {
		mapPath = strings.TrimSpace(document.MapTarget.Path)
	}
	pageStatus := seo_domain.SeoPageStatusDraft
	if published {
		pageStatus = seo_domain.SeoPageStatusPublished
	}

	created := false
	updated := false
	var page *seo_domain.SeoDomain
	if existing == nil {
		page = &seo_domain.SeoDomain{
			Slug:              projection.Slug,
			OriginURL:         canonicalPath,
			CanonicalURL:      canonicalPath,
			RefType:           enums.ESEORefTypeAdmUnit,
			RefID:             &refID,
			RefSource:         seo_domain.SeoRefSourceTqdAdministrativeUnit,
			RefLabel:          projection.Name,
			RefURL:            projection.EntityID,
			SourceStatus:      seo_domain.SeoSourceStatusLinked,
			RefMissing:        false,
			RefSnapshotJSON:   datatypes.JSON(projectionJSON),
			RefHash:           snapshotHash,
			RefLastSyncedAt:   &now,
			QualityScore:      qualityScore,
			Scope:             seo_domain.SeoScopePublic,
			Classify:          enums.ESeoClassifyA,
			PageStatus:        pageStatus,
			Title:             strings.TrimSpace(document.Heading.Title),
			Description:       strings.TrimSpace(document.Heading.Summary),
			Content:           strings.TrimSpace(document.Heading.Summary),
			Summary:           strings.TrimSpace(document.Heading.Summary),
			Published:         published,
			PublishedAt:       publishedAt,
			IsSiteMap:         indexable,
			IsIndex:           indexable,
			IsRobot:           &isRobot,
			SitemapPriority:   0.8,
			SitemapChangeFreq: seo_domain.SeoChangeFreqWeekly,
			NeedGenerate:      true,
			SourceUpdatedAt:   projection.Source.UpdatedAt,
			RenderStatus:      seo_domain.SeoRenderStatusPending,
			TemplateKey:       "administrative-unit",
			TemplateVersion:   dto.SeoPublicDocumentSchemaVersion,
			DeepLink:          mapPath,
			Metadata:          datatypes.JSON(documentJSON),
			Note:              "Provisioned by SEO Operating System V2.U typed Administrative Unit workflow.",
		}
		page, err = u.seoDomainRepo.Create(ctx, page)
		if err != nil {
			return nil, err
		}
		created = true
	} else {
		updates := map[string]interface{}{
			"slug":                projection.Slug,
			"origin_url":          canonicalPath,
			"canonical_url":       canonicalPath,
			"ref_type":            enums.ESEORefTypeAdmUnit,
			"ref_id":              refID,
			"ref_source":          seo_domain.SeoRefSourceTqdAdministrativeUnit,
			"ref_label":           projection.Name,
			"ref_url":             projection.EntityID,
			"source_status":       seo_domain.SeoSourceStatusLinked,
			"ref_missing":         false,
			"ref_snapshot_json":   datatypes.JSON(projectionJSON),
			"ref_hash":            snapshotHash,
			"ref_last_synced_at":  now,
			"quality_score":       qualityScore,
			"scope":               seo_domain.SeoScopePublic,
			"classify":            enums.ESeoClassifyA,
			"page_status":         pageStatus,
			"title":               strings.TrimSpace(document.Heading.Title),
			"description":         strings.TrimSpace(document.Heading.Summary),
			"content":             strings.TrimSpace(document.Heading.Summary),
			"summary":             strings.TrimSpace(document.Heading.Summary),
			"published":           published,
			"published_at":        publishedAt,
			"is_site_map":         indexable,
			"is_index":            indexable,
			"is_robot":            true,
			"sitemap_priority":    0.8,
			"sitemap_change_freq": seo_domain.SeoChangeFreqWeekly,
			"need_generate":       true,
			"render_status":       seo_domain.SeoRenderStatusPending,
			"last_render_error":   "",
			"source_updated_at":   projection.Source.UpdatedAt,
			"template_key":        "administrative-unit",
			"template_version":    dto.SeoPublicDocumentSchemaVersion,
			"deep_link":           mapPath,
			"metadata":            datatypes.JSON(documentJSON),
			"note":                "Synced by SEO Operating System V2.U typed Administrative Unit workflow.",
		}
		if err := u.seoDomainRepo.Update(ctx, existing.ID, updates); err != nil {
			return nil, err
		}
		page, err = u.seoDomainRepo.GetByID(ctx, existing.ID)
		if err != nil {
			return nil, err
		}
		updated = true
	}
	if page == nil {
		return nil, _errors.ReturnError(503, "Không thể đọc SEO domain sau provisioning")
	}

	generated := false
	if req.Generate {
		if !published {
			return nil, _errors.ReturnError(409, "Không thể generate public HTML khi page chưa published")
		}
		if u.renderUsecase == nil {
			return nil, _errors.ReturnError(503, "SEO render pipeline chưa sẵn sàng")
		}
		generation, generationErr := u.renderUsecase.GenerateSeoPage(ctx, &dto.GenerateSeoPageRequest{
			ID:          page.ID,
			TriggerType: "manual",
			Reason:      firstNonEmptyRender(strings.TrimSpace(req.Reason), "administrative_unit_v2u_provisioning"),
		})
		if generationErr != nil {
			return nil, generationErr
		}
		if generation == nil || generation.RenderStatus != seo_domain.SeoRenderStatusSuccess {
			message := "Administrative Unit render không thành công"
			if generation != nil && strings.TrimSpace(generation.ErrorMessage) != "" {
				message += ": " + generation.ErrorMessage
			}
			return nil, _errors.ReturnError(503, message)
		}
		generated = true
		page, err = u.seoDomainRepo.GetByID(ctx, page.ID)
		if err != nil {
			return nil, err
		}
		if page == nil {
			return nil, _errors.ReturnError(503, "Không thể đọc SEO domain sau generation")
		}
	}

	operation := "updated"
	if created {
		operation = "created"
	}
	if generated {
		operation += "_and_generated"
	}

	return &dto.ProvisionAdministrativeUnitSeoResponse{
		SchemaVersion:         dto.SeoAdministrativeUnitProvisioningSchemaVersion,
		Operation:             operation,
		SeoDomainID:           page.ID,
		EntityID:              projection.EntityID,
		CanonicalSlug:         projection.Slug,
		CanonicalURL:          canonicalPath,
		Created:               created,
		Updated:               updated,
		Generated:             generated,
		PageStatus:            page.PageStatus,
		Published:             page.Published,
		RenderStatus:          page.RenderStatus,
		NeedGenerate:          page.NeedGenerate,
		SourceStatus:          page.SourceStatus,
		DocumentVersion:       seoPublicDocumentVersion(page),
		ETag:                  page.StaticHtmlHash,
		AnswerReadiness:       readiness,
		IndexDecision:         indexDecision,
		PolicyEvaluationScope: "source-content-before-publication",
		FinalPublicIndexable:  page.IsPublicIndexable(),
		Warnings:              compactUniqueStrings(warnings),
	}, nil
}

func (u *SeoAdministrativeUnitProvisioningUsecase) findExistingAdministrativeUnitPage(
	ctx context.Context,
	typedIdentity string,
	canonicalPath string,
	slug string,
) (*seo_domain.SeoDomain, error) {
	if page, err := u.seoDomainRepo.GetByRefIdentity(ctx, uint32(enums.ESEORefTypeAdmUnit), seo_domain.SeoRefSourceTqdAdministrativeUnit, strings.TrimSpace(typedIdentity)); err != nil {
		return nil, err
	} else if page != nil {
		return page, nil
	}
	if page, err := u.seoDomainRepo.GetByCanonicalURL(ctx, canonicalPath); err != nil {
		return nil, err
	} else if page != nil {
		return page, nil
	}
	if page, err := u.seoDomainRepo.GetBySlug(ctx, slug); err != nil {
		return nil, err
	} else if page != nil {
		return page, nil
	}
	return nil, nil
}

func validateExistingAdministrativeUnitOwnership(
	page *seo_domain.SeoDomain,
	refID uint64,
	typedEntityID string,
) ([]string, error) {
	if page == nil {
		return nil, nil
	}

	warnings := make([]string, 0, 2)
	if page.RefType != 0 && page.RefType != enums.ESEORefTypeAdmUnit {
		return nil, _errors.ReturnError(409, "Slug hoặc canonical đã thuộc page family khác")
	}
	if page.RefID != nil && *page.RefID != refID {
		return nil, _errors.ReturnError(409, "Slug hoặc canonical đã thuộc Administrative Unit khác")
	}

	source := strings.ToLower(strings.TrimSpace(page.RefSource))
	switch source {
	case "", seo_domain.SeoRefSourceManual, "tqd", seo_domain.SeoRefSourceTqdAdministrativeUnit:
		// Empty/manual and the legacy short TQD key may be migrated safely once
		// identity and canonical ownership have been verified above.
	case seo_domain.SeoRefSourceBdsproRegion:
		warnings = append(warnings, "legacy_bdspro_region_source_migrated")
	default:
		return nil, _errors.ReturnError(409, "Slug hoặc canonical đã thuộc source owner khác")
	}

	existingTypedIdentity := strings.TrimSpace(page.RefURL)
	if strings.HasPrefix(existingTypedIdentity, "province:") || strings.HasPrefix(existingTypedIdentity, "ward:") {
		if existingTypedIdentity != strings.TrimSpace(typedEntityID) {
			return nil, _errors.ReturnError(409, "Typed Administrative Unit identity xung đột với page hiện tại")
		}
	}

	return warnings, nil
}

func normalizeAdministrativeUnitProvisioningIdentity(value string) (string, error) {
	identity := strings.TrimSpace(value)
	lower := strings.ToLower(identity)
	if identity == "" ||
		strings.ContainsAny(identity, `/\\?#`) ||
		strings.Contains(identity, "..") ||
		lower == "undefined" ||
		lower == "null" ||
		strings.Contains(lower, "[object object]") {
		return "", _errors.ReturnError(400, "Administrative Unit identity không hợp lệ")
	}
	if prefix, rawID, ok := strings.Cut(identity, ":"); ok {
		prefix = strings.ToLower(strings.TrimSpace(prefix))
		rawID = strings.TrimSpace(rawID)
		if (prefix != "province" && prefix != "ward") || rawID == "" || rawID == "0" {
			return "", _errors.ReturnError(400, "Typed identity phải có dạng province:<source-id> hoặc ward:<source-id>")
		}
		return prefix + ":" + rawID, nil
	}
	if _, err := normalizePublicSeoSlug(identity); err != nil {
		return "", _errors.ReturnError(400, "Administrative Unit identity phải là typed identity hoặc canonical slug")
	}
	return identity, nil
}

func validateAdministrativeUnitProvisioningProjection(projection *dto.TqdAdministrativeUnitProjection) error {
	if projection == nil {
		return _errors.ReturnError(503, "TQD trả projection rỗng")
	}
	if projection.SchemaVersion != dto.TqdAdministrativeUnitProjectionSchemaVersion ||
		strings.TrimSpace(projection.EntityID) == "" ||
		strings.TrimSpace(projection.Slug) == "" ||
		strings.TrimSpace(projection.Name) == "" ||
		strings.TrimSpace(projection.AdministrativeCode) == "" {
		return _errors.ReturnError(503, "Administrative Unit projection thiếu identity bắt buộc")
	}
	expectedCanonical := "/tin-tuc/dia-ban/" + strings.TrimSpace(projection.Slug)
	if !strings.Contains(strings.TrimSpace(projection.Slug), "/") {
		return _errors.ReturnError(503, "Administrative Unit projection slug phải dạng {code}/{name}")
	}
	if strings.TrimSpace(projection.CanonicalPath) != expectedCanonical {
		return _errors.ReturnError(503, "Administrative Unit projection canonical không hợp lệ")
	}
	prefix, rawID, ok := strings.Cut(strings.TrimSpace(projection.EntityID), ":")
	prefix = strings.ToLower(strings.TrimSpace(prefix))
	rawID = strings.TrimSpace(rawID)
	if !ok || (prefix != "province" && prefix != "ward") || rawID == "" {
		return _errors.ReturnError(503, "Administrative Unit projection entityId không phải typed identity")
	}
	if strings.TrimSpace(projection.UnitType) != "" && prefix != strings.ToLower(strings.TrimSpace(projection.UnitType)) {
		return _errors.ReturnError(503, "Administrative Unit projection entityId không khớp unitType")
	}
	if administrativeUnitReferenceID(projection) == 0 {
		return _errors.ReturnError(503, "Administrative Unit projection administrativeCode không hợp lệ")
	}
	return nil
}

func resolveAdministrativeUnitPublication(existing *seo_domain.SeoDomain, requested *bool) (bool, *time.Time) {
	if existing == nil {
		if requested == nil || !*requested {
			return false, nil
		}
		return true, nil
	}
	published := existing.Published && existing.PageStatus == seo_domain.SeoPageStatusPublished
	publishedAt := existing.PublishedAt
	if requested != nil {
		published = *requested
		if !published {
			publishedAt = nil
		}
	}
	return published, publishedAt
}

func cloneSeoAnswerReadiness(value *dto.SeoAnswerReadiness) *dto.SeoAnswerReadiness {
	if value == nil {
		return nil
	}
	copyValue := *value
	copyValue.RequiredFacts = append([]string(nil), value.RequiredFacts...)
	copyValue.MissingRequiredFacts = append([]string(nil), value.MissingRequiredFacts...)
	copyValue.RequiredSections = append([]string(nil), value.RequiredSections...)
	copyValue.RenderedSections = append([]string(nil), value.RenderedSections...)
	return &copyValue
}

func cloneSeoIndexDecision(value *dto.SeoIndexDecision) *dto.SeoIndexDecision {
	if value == nil {
		return nil
	}
	copyValue := *value
	copyValue.ReasonCodes = append([]string(nil), value.ReasonCodes...)
	return &copyValue
}
