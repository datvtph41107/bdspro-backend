package rule

import (
	"strings"

	"tqd/internal/usecase/resolver/layer_resolver/types"
)

// LandUseCodeResolver — Engine chuẩn hóa ngữ nghĩa land use code
type LandUseCodeResolver struct {
	canonicalMap map[string]string // raw code → canonical code
	labels       map[string]string // canonical code → display name
}

// MakeLandUseCodeResolver tạo LandUseCodeResolver
func MakeLandUseCodeResolver(canonicalMap map[string]string, labels map[string]string) *LandUseCodeResolver {
	if canonicalMap == nil {
		canonicalMap = make(map[string]string)
	}
	if labels == nil {
		labels = make(map[string]string)
	}
	return &LandUseCodeResolver{
		canonicalMap: canonicalMap,
		labels:       labels,
	}
}

// Name trả về tên engine
func (n *LandUseCodeResolver) Name() string { return "LandUseCodeResolver" }

// Normalize chuẩn hóa raw code → canonical code + display name
func (n *LandUseCodeResolver) Normalize(rawCode string) (canonical string, label string) {
	canonical, ok := n.canonicalMap[strings.TrimSpace(rawCode)]
	if !ok {
		canonical = rawCode
	}
	label = n.labels[canonical]
	return
}

// NormalizeCandidate chuẩn hóa LandUseCode trong LayerCandidate
func (n *LandUseCodeResolver) NormalizeCandidate(c *types.LayerCandidate) {
	code, name := n.Normalize(c.LandUseCode)
	c.LandUseCode = code
	if name != "" {
		c.LandUseName = name
	}
}

// Validate kiểm tra config
func (n *LandUseCodeResolver) Validate() error {
	return nil
}
