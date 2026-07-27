package usecase

import (
	"context"
	"testing"
	"time"

	seo_domain "crm/internal/domain/seo"
	"crm/internal/dto"
	"crm/internal/enums"
	provider_interface "crm/internal/interface/provider"
	"crm/internal/repo"
)

type administrativeUnitProvisioningRepo struct {
	repo.SeoDomainRepo
	page              *seo_domain.SeoDomain
	lookupRefIdentity string
}

func (r *administrativeUnitProvisioningRepo) Create(_ context.Context, page *seo_domain.SeoDomain) (*seo_domain.SeoDomain, error) {
	page.ID = 101
	page.NormalizeLifecycle()
	r.page = page
	return page, nil
}

func (r *administrativeUnitProvisioningRepo) GetByID(_ context.Context, id uint64) (*seo_domain.SeoDomain, error) {
	if r.page != nil && r.page.ID == id {
		return r.page, nil
	}
	return nil, nil
}

func (r *administrativeUnitProvisioningRepo) GetByRefSource(_ context.Context, _ uint32, _ string, _ uint64) (*seo_domain.SeoDomain, error) {
	return nil, nil
}

func (r *administrativeUnitProvisioningRepo) GetByRefIdentity(_ context.Context, _ uint32, _ string, refURL string) (*seo_domain.SeoDomain, error) {
	r.lookupRefIdentity = refURL
	if r.page != nil && r.page.RefURL == refURL {
		return r.page, nil
	}
	return nil, nil
}

func (r *administrativeUnitProvisioningRepo) GetByCanonicalURL(_ context.Context, _ string) (*seo_domain.SeoDomain, error) {
	return nil, nil
}

func (r *administrativeUnitProvisioningRepo) GetBySlug(_ context.Context, _ string) (*seo_domain.SeoDomain, error) {
	return nil, nil
}

type administrativeUnitProvisioningProvider struct {
	provider_interface.TqdProvider
	projection *dto.TqdAdministrativeUnitProjection
}

func (p *administrativeUnitProvisioningProvider) GetAdministrativeUnitProjection(_ context.Context, _ string) (*dto.TqdAdministrativeUnitProjection, error) {
	return p.projection, nil
}

func TestProvisionAdministrativeUnitCreatesTypedDraft(t *testing.T) {
	updatedAt := time.Date(2026, 7, 24, 8, 0, 0, 0, time.UTC)
	count := int64(1)
	repository := &administrativeUnitProvisioningRepo{}
	provider := &administrativeUnitProvisioningProvider{projection: &dto.TqdAdministrativeUnitProjection{
		SchemaVersion:      dto.TqdAdministrativeUnitProjectionSchemaVersion,
		EntityID:           "ward:26734",
		CanonicalKey:       "ward:26734",
		Slug:               "phuong-ben-nghe-26734",
		CanonicalPath:      "/dia-ban/phuong-ben-nghe-26734",
		UnitType:           "ward",
		UnitTypeName:       "Phường",
		Name:               "Phường Bến Nghé",
		DisplayName:        "Phường Bến Nghé",
		AdministrativeCode: "26734",
		Hierarchy: []dto.TqdAdministrativeUnitHierarchyItem{{
			EntityID:           "ward:26734",
			UnitType:           "ward",
			Name:               "Phường Bến Nghé",
			AdministrativeCode: "26734",
			CanonicalPath:      "/dia-ban/phuong-ben-nghe-26734",
		}},
		PlanningContext: dto.TqdAdministrativePlanningContext{
			Summary:              "Địa bàn có dữ liệu quy hoạch công khai để tra cứu.",
			PlanningProjectCount: &count,
			DataAvailability:     "available",
		},
		Source: dto.TqdPlanningProjectionSource{
			System:         "tqd",
			Dataset:        "ward_v2",
			Authority:      "official",
			AuthorityName:  "Cơ quan quản lý nhà nước",
			AuthorityClass: "government",
			DataQuality:    "verified",
			UpdatedAt:      &updatedAt,
		},
	}}

	uc := NewSeoAdministrativeUnitProvisioningUsecase(repository, provider, nil)
	response, err := uc.Provision(context.Background(), &dto.ProvisionAdministrativeUnitSeoRequest{
		Identity: "ward:26734",
	})
	if err != nil {
		t.Fatalf("Provision returned error: %v", err)
	}
	if response == nil || !response.Created || response.Generated {
		t.Fatalf("unexpected response: %#v", response)
	}
	if repository.lookupRefIdentity != "ward:26734" {
		t.Fatalf("expected typed identity lookup, got %q", repository.lookupRefIdentity)
	}
	if response.CanonicalSlug != "phuong-ben-nghe-26734" || response.CanonicalURL != "/dia-ban/phuong-ben-nghe-26734" {
		t.Fatalf("canonical mismatch: %#v", response)
	}
	if response.AnswerReadiness == nil || !response.AnswerReadiness.Ready {
		t.Fatalf("expected ready answer: %#v", response.AnswerReadiness)
	}
	if response.IndexDecision == nil || !response.IndexDecision.Eligible {
		t.Fatalf("expected source/content index eligibility: %#v", response.IndexDecision)
	}
	if response.PolicyEvaluationScope != "source-content-before-publication" || response.FinalPublicIndexable {
		t.Fatalf("draft provisioning must expose candidate policy separately from final public indexability: %#v", response)
	}
	if repository.page == nil {
		t.Fatal("expected persisted page")
	}
	if repository.page.RefSource != seo_domain.SeoRefSourceTqdAdministrativeUnit {
		t.Fatalf("unexpected ref source: %q", repository.page.RefSource)
	}
	if repository.page.RefURL != "ward:26734" {
		t.Fatalf("typed identity must be persisted in ref URL, got %q", repository.page.RefURL)
	}
	if repository.page.Published || repository.page.IsIndex || repository.page.IsSiteMap {
		t.Fatalf("new provisioning must remain draft until explicit publish: %#v", repository.page)
	}
	if repository.page.RenderStatus != seo_domain.SeoRenderStatusPending || !repository.page.NeedGenerate {
		t.Fatalf("new page must await generation: %#v", repository.page)
	}
}

func TestProvisionAdministrativeUnitRejectsInvalidTypedIdentity(t *testing.T) {
	uc := NewSeoAdministrativeUnitProvisioningUsecase(
		&administrativeUnitProvisioningRepo{},
		&administrativeUnitProvisioningProvider{},
		nil,
	)
	for _, identity := range []string{"", "district:1", "ward:0", "undefined", "../ward:1"} {
		if _, err := uc.Provision(context.Background(), &dto.ProvisionAdministrativeUnitSeoRequest{Identity: identity}); err == nil {
			t.Fatalf("expected invalid identity %q to fail", identity)
		}
	}
}

func TestValidateExistingAdministrativeUnitOwnershipRejectsConflicts(t *testing.T) {
	refID := uint64(26734)
	otherRefID := uint64(26735)

	cases := []struct {
		name string
		page *seo_domain.SeoDomain
	}{
		{
			name: "other page family",
			page: &seo_domain.SeoDomain{RefType: enums.ESEORefTypeProject, RefID: &refID, RefSource: seo_domain.SeoRefSourceTqdAdministrativeUnit},
		},
		{
			name: "other administrative unit",
			page: &seo_domain.SeoDomain{RefType: enums.ESEORefTypeAdmUnit, RefID: &otherRefID, RefSource: seo_domain.SeoRefSourceTqdAdministrativeUnit},
		},
		{
			name: "other source owner",
			page: &seo_domain.SeoDomain{RefType: enums.ESEORefTypeAdmUnit, RefID: &refID, RefSource: seo_domain.SeoRefSourceCrmReport},
		},
		{
			name: "other typed identity",
			page: &seo_domain.SeoDomain{RefType: enums.ESEORefTypeAdmUnit, RefID: &refID, RefSource: seo_domain.SeoRefSourceTqdAdministrativeUnit, RefURL: "ward:26735"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := validateExistingAdministrativeUnitOwnership(tc.page, refID, "ward:26734"); err == nil {
				t.Fatalf("expected ownership conflict for %#v", tc.page)
			}
		})
	}
}

func TestValidateExistingAdministrativeUnitOwnershipAllowsLegacyMigration(t *testing.T) {
	refID := uint64(26734)
	warnings, err := validateExistingAdministrativeUnitOwnership(&seo_domain.SeoDomain{
		RefType:   enums.ESEORefTypeAdmUnit,
		RefID:     &refID,
		RefSource: seo_domain.SeoRefSourceBdsproRegion,
	}, refID, "ward:26734")
	if err != nil {
		t.Fatalf("legacy Administrative Unit source should migrate safely: %v", err)
	}
	if len(warnings) != 1 || warnings[0] != "legacy_bdspro_region_source_migrated" {
		t.Fatalf("unexpected migration warnings: %#v", warnings)
	}
}
