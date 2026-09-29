package seo_domain

import (
	"testing"

	"crm/internal/enums"
)

func TestPlanningProjectRegistryUsesTqdResolver(t *testing.T) {
	cfg, ok := GetSeoRefTypeConfig(
		uint32(enums.ESEORefTypeProject),
		SeoRefSourceTqdPlanningProject,
		"tqd",
	)
	if !ok {
		t.Fatal("planning_project registry config is missing")
	}

	if cfg.Key != "planning_project" {
		t.Fatalf("unexpected source key: %q", cfg.Key)
	}
	if cfg.ResolverKey != SeoRefSourceTqdPlanningProject {
		t.Fatalf("unexpected resolver: %q", cfg.ResolverKey)
	}
	if cfg.SourceService != "tqd" {
		t.Fatalf("unexpected source service: %q", cfg.SourceService)
	}
	if cfg.Capabilities.Searchable {
		t.Fatal("planning_project must not claim search support")
	}
	if !cfg.Capabilities.Resolvable {
		t.Fatal("planning_project must be resolvable")
	}
	if !IsSeoResolverReady(cfg.ResolverKey) {
		t.Fatalf("resolver is not ready: %q", cfg.ResolverKey)
	}
}

func TestProjectRefTypeSupportsBdsproAndTqd(t *testing.T) {
	bdspro, bdsproOK := GetSeoRefTypeConfig(
		uint32(enums.ESEORefTypeProject),
		SeoRefSourceBdsproProject,
		"bdspro",
	)

	tqd, tqdOK := GetSeoRefTypeConfig(
		uint32(enums.ESEORefTypeProject),
		SeoRefSourceTqdPlanningProject,
		"tqd",
	)

	if !bdsproOK || !tqdOK {
		t.Fatalf(
			"expected both configs: bdspro=%v tqd=%v",
			bdsproOK,
			tqdOK,
		)
	}

	if bdspro.ResolverKey == tqd.ResolverKey {
		t.Fatal("the two sources must use different resolver keys")
	}
}
