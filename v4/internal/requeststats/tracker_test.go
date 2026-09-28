package requeststats

import "testing"

func TestNormalizePath(t *testing.T) {
	tests := map[string]string{
		"/v2/user/profile/info/me":                               "/v2/user/profile/info/me",
		"/v2/bdspro/v2/product/570/detail":                       "/v2/bdspro/v2/product/:id/detail",
		"/v2/bdspro/v2/property/undefined/timestamp":             "/v2/bdspro/v2/property/:id/timestamp",
		"/v2/hub/user-guides/42":                                 "/v2/hub/user-guides/:id",
		"/v2/user/profiles/a1b2c3d4-e5f6-7890-abcd-ef1234567890": "/v2/user/profiles/:id",
		"/v2/hub/request-stats/dashboard":                        "/v2/hub/request-stats/dashboard",
		"": "/",
	}
	for in, want := range tests {
		if got := normalizePath(in); got != want {
			t.Fatalf("normalizePath(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestShouldSkipDashboardAPIs(t *testing.T) {
	tracker := &Tracker{ignorePrefixes: mergeIgnorePrefixes(nil)}
	skip := []string{
		"/health",
		"/v2/hub/request-stats/dashboard",
		"/v2/hub/request-stats/users/dashboard",
		"/v2/hub/search-latency/dashboard",
		"/v2/hub/keyword-search-stats",
		"/v2/bdspro/dashboard/post-stats",
		"/v2/user/admin/users/dashboard/daily",
	}
	for _, path := range skip {
		if !tracker.shouldSkip(path) {
			t.Fatalf("expected %q to be skipped", path)
		}
	}
	if tracker.shouldSkip("/v2/hub/provinces") {
		t.Fatal("did not expect hub provinces to be skipped")
	}
}

func TestMergeIgnorePrefixesAddsConfigExtras(t *testing.T) {
	merged := mergeIgnorePrefixes([]string{"/v2/custom/dashboard", "/v2/hub/request-stats"})
	tracker := &Tracker{ignorePrefixes: merged}
	if !tracker.shouldSkip("/v2/custom/dashboard/x") {
		t.Fatal("expected config ignore prefix to apply")
	}
}

func TestIsSpecURL(t *testing.T) {
	tracker := &Tracker{specURLs: map[string]struct{}{"/v2/crm/sitemap.xml": {}}}
	if !tracker.isSpecURL("/v2/crm/sitemap.xml") {
		t.Fatal("expected configured spec URL to match")
	}
	if tracker.isSpecURL("/v2/other") {
		t.Fatal("did not expect unrelated path to match")
	}
}

func TestExtractService(t *testing.T) {
	if got := extractService("/v2/hub/provinces"); got != "hub" {
		t.Fatalf("extractService = %q, want hub", got)
	}
	if got := extractService("/v2.1/auth/otp/verify"); got != "auth" {
		t.Fatalf("extractService = %q, want auth", got)
	}
}
