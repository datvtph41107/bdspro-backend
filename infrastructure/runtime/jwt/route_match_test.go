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

func TestMatchRouteSupportsNamedParametersAndCatchAll(t *testing.T) {
	t.Parallel()

	if !matchRoute("/v2/items/:id", "/v2/items/42") {
		t.Fatal("named parameter route did not match")
	}
	if matchRoute("/v2/items/:id", "/v2/items/42/history") {
		t.Fatal("named parameter route matched extra segments")
	}
	if !matchRoute("/v2/files/*path", "/v2/files/a/b/c") {
		t.Fatal("catch-all route did not match remaining segments")
	}
	if matchRoute("/v2/files/*path/more", "/v2/files/a/more") {
		t.Fatal("non-terminal catch-all route was accepted")
	}
}
