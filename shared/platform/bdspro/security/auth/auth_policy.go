package auth

import (
	"strings"

	"gateway/internal/auth"
	"gateway/internal/principal"
)

// Requirement mô tả yêu cầu auth cho một route
type Requirement uint8

const (
	// RequireAuth route cần token hợp lệ
	RequireAuth Requirement = iota + 1
	// Public route không cần token
	Public
	// RequireTempOnly chỉ chấp nhận TempToken (e.g. onboarding flow)
	RequireTempOnly
)

// RoutePolicy quyết định requirement cho một path
// Tách khỏi middleware — testable độc lập, thay đổi không ảnh hưởng middleware chain
type RoutePolicy struct {
	// publicPrefixes là prefix của các route public
	publicPrefixes []string
	// tempOnlyPrefixes là prefix của các route chỉ dùng temp token
	tempOnlyPrefixes []string
}

// NewRoutePolicy tạo policy từ config — không hardcode trong main
func NewRoutePolicy(publicPrefixes []string, tempOnlyPrefixes []string) *RoutePolicy {
	return &RoutePolicy{
		publicPrefixes:   publicPrefixes,
		tempOnlyPrefixes: tempOnlyPrefixes,
	}
}

// Decide trả về requirement cho path — O(n) nhưng n nhỏ và được cache ở layer trên
func (p *RoutePolicy) Decide(path string) Requirement {
	for _, prefix := range p.tempOnlyPrefixes {
		if matchPrefix(prefix, path) {
			return RequireTempOnly
		}
	}
	for _, prefix := range p.publicPrefixes {
		if matchPrefix(prefix, path) {
			return Public
		}
	}
	return RequireAuth
}

// Authorize kiểm tra principal có đáp ứng requirement không
// Trả về reason để log — không phải bool mà là typed decision
type Decision struct {
	Allowed bool
	Reason  string
}

func (p *RoutePolicy) Authorize(path string, pr *principal.Principal) Decision {
	req := p.Decide(path)

	switch req {
	case Public:
		return Decision{Allowed: true, Reason: "public_route"}

	case RequireTempOnly:
		if pr == nil {
			return Decision{Allowed: false, Reason: "missing_token_for_temp_route"}
		}
		if !auth.IsTempToken(pr) {
			return Decision{Allowed: false, Reason: "non_temp_token_on_temp_route"}
		}
		return Decision{Allowed: true, Reason: "temp_token_allowed"}

	case RequireAuth:
		if pr == nil || !pr.IsAuthenticated() {
			return Decision{Allowed: false, Reason: "missing_or_invalid_token"}
		}
		if auth.IsTempToken(pr) {
			return Decision{Allowed: false, Reason: "temp_token_on_protected_route"}
		}
		return Decision{Allowed: true, Reason: "authenticated"}
	}

	return Decision{Allowed: false, Reason: "unknown_requirement"}
}

// matchPrefix hỗ trợ exact match và prefix match với "/"
// "/v2/auth" match "/v2/auth/otp" và "/v2/auth"
// Không dùng regex để tránh overhead — prefix matching đủ cho gateway
func matchPrefix(prefix, path string) bool {
	if !strings.HasPrefix(path, prefix) {
		return false
	}
	// Tránh false positive: "/v1/user" không match "/v1/users"
	if len(path) > len(prefix) {
		next := path[len(prefix)]
		return next == '/' || next == '?'
	}
	return true
}
