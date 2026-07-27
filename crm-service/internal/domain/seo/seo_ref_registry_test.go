package seo_domain

import (
	"testing"

	"crm/internal/enums"
)

func TestAdministrativeUnitSourceRegistryUsesTQDAndFailsClosedUntilTypedProvisioning(t *testing.T) {
	cfg, ok := GetSeoRefTypeConfig(uint32(enums.ESEORefTypeAdmUnit), "", "")
	if !ok {
		t.Fatal("administrative unit source config is missing")
	}
	cfg = NormalizeSeoRefTypeConfig(cfg)

	if cfg.SourceService != "tqd" || cfg.ResolverKey != SeoRefSourceTqdAdministrativeUnit {
		t.Fatalf("administrative unit source owner must be TQD: %#v", cfg)
	}
	if cfg.PublicURLPattern != "/dia-ban/{slug}" {
		t.Fatalf("administrative unit canonical pattern mismatch: %q", cfg.PublicURLPattern)
	}
	if cfg.Status != SeoSourceStatusComingSoon || cfg.Capabilities.Resolvable || cfg.Capabilities.Searchable {
		t.Fatalf("generic numeric resolver must stay fail-closed until typed provisioning exists: %#v", cfg)
	}
	if IsSeoResolverReady(cfg.ResolverKey) {
		t.Fatal("typed Administrative Unit resolver must not be reported ready before implementation")
	}
}
