package _jwt

import "testing"

func TestMatchRouteUsesSegmentBoundary(t *testing.T) {
	t.Parallel()

	if !matchRoute("/v2/auth/login/key", "/v2/auth/login/key") {
		t.Fatal("exact route did not match")
	}
	if !matchRoute("/v2/bdspro/post/detail/", "/v2/bdspro/post/detail/123") {
		t.Fatal("nested route did not match")
	}
	if matchRoute("/v2/auth/login/key", "/v2/auth/login/key-malicious") {
		t.Fatal("partial path segment matched public route")
	}
	if matchRoute("/", "/private") {
		t.Fatal("root route matched every path")
	}
}
