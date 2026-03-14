package fieldmask_test

import (
	"common/pkg/fieldmask"
	"testing"

	"github.com/stretchr/testify/assert"
)

// ─── Zero value ───────────────────────────────────────────────────────────────

func TestZeroValue_AllowAll(t *testing.T) {
	var mask fieldmask.FieldMask // zero value

	assert.True(t, mask.IsAllowAll())
	assert.True(t, mask.Allows("info.title"))
	assert.True(t, mask.Allows("anything.at.all"))
}

// ─── FromPaths ────────────────────────────────────────────────────────────────

func TestFromPaths_Nil_AllowAll(t *testing.T) {
	mask := fieldmask.FromPaths(nil)
	assert.True(t, mask.IsAllowAll())
}

func TestFromPaths_Empty_AllowAll(t *testing.T) {
	mask := fieldmask.FromPaths([]string{})
	assert.True(t, mask.IsAllowAll())
}

func TestFromPaths_BlankStrings_AllowAll(t *testing.T) {
	mask := fieldmask.FromPaths([]string{"", ""})
	assert.True(t, mask.IsAllowAll())
}

func TestFromPaths_WithPaths_StrictWhitelist(t *testing.T) {
	mask := fieldmask.FromPaths([]string{"info.title", "info.note"})

	assert.False(t, mask.IsAllowAll())
	assert.True(t, mask.Allows("info.title"))
	assert.True(t, mask.Allows("info.note"))
	assert.False(t, mask.Allows("info.scope")) // không có trong mask
	assert.False(t, mask.Allows("landInfo.area_total"))
}

// ─── Only ─────────────────────────────────────────────────────────────────────

func TestOnly(t *testing.T) {
	mask := fieldmask.Only("info.title", "landInfo.area")

	assert.True(t, mask.Allows("info.title"))
	assert.False(t, mask.Allows("info.note"))
}

// ─── Allows — exact match ─────────────────────────────────────────────────────

func TestAllows_ExactMatch(t *testing.T) {
	mask := fieldmask.Only("info.title")

	assert.True(t, mask.Allows("info.title"))
	assert.False(t, mask.Allows("info.titl")) // prefix của path, không phải ngược lại
	assert.False(t, mask.Allows("info.title2"))
}

// ─── Allows — prefix match ────────────────────────────────────────────────────

func TestAllows_PrefixMatch(t *testing.T) {
	// mask có "landInfo.area" → cover tất cả children
	mask := fieldmask.Only("landInfo.area")

	assert.True(t, mask.Allows("landInfo.area"))            // exact
	assert.True(t, mask.Allows("landInfo.area.area_total")) // child
	assert.True(t, mask.Allows("landInfo.area.area_land"))  // child
	assert.False(t, mask.Allows("landInfo.areaXXX"))        // không phải child
	assert.False(t, mask.Allows("landInfo.expiry"))         // sibling
	assert.False(t, mask.Allows("info.title"))              // khác domain
}

// ─── AllowsAny ────────────────────────────────────────────────────────────────

func TestAllowsAny(t *testing.T) {
	mask := fieldmask.Only("info.title")

	assert.True(t, mask.AllowsAny("info.note", "info.title")) // 1 trong 2 ok
	assert.False(t, mask.AllowsAny("info.note", "info.scope"))
	assert.False(t, mask.AllowsAny()) // empty → false
}

func TestAllowsAny_AllowAll(t *testing.T) {
	var mask fieldmask.FieldMask
	assert.True(t, mask.AllowsAny("anything"))
}

// ─── Intersect ────────────────────────────────────────────────────────────────

func TestIntersect_StripUnauthorized(t *testing.T) {
	// user gửi mask muốn update lineage_meta, nhưng không có quyền
	userMask := fieldmask.Only("info.title", "lineage_meta.national_id")
	adminOnly := []string{"info.title", "info.note"}

	result := userMask.Intersect(adminOnly)

	assert.True(t, result.Allows("info.title"))
	assert.False(t, result.Allows("info.note"))                // user không gửi
	assert.False(t, result.Allows("lineage_meta.national_id")) // stripped
}

func TestIntersect_AllowAll_ClipsToAllowed(t *testing.T) {
	var mask fieldmask.FieldMask // allowAll
	allowed := []string{"info.title", "info.note"}

	result := mask.Intersect(allowed)

	assert.True(t, result.Allows("info.title"))
	assert.True(t, result.Allows("info.note"))
	assert.False(t, result.Allows("info.scope")) // không trong allowed
}

func TestIntersect_NoOverlap_Empty(t *testing.T) {
	mask := fieldmask.Only("info.title")
	result := mask.Intersect([]string{"landInfo.area_total"})

	assert.False(t, result.Allows("info.title"))
	assert.False(t, result.Allows("landInfo.area_total"))
}

// ─── Paths — deterministic order ─────────────────────────────────────────────

func TestPaths_Sorted(t *testing.T) {
	mask := fieldmask.Only("info.title", "info.note", "landInfo.area_total")
	paths := mask.Paths()

	// Sorted → deterministic, test không bao giờ flaky
	assert.Equal(t, []string{
		"info.note",
		"info.title",
		"landInfo.area_total",
	}, paths)
}

func TestPaths_AllowAll(t *testing.T) {
	var mask fieldmask.FieldMask
	assert.Equal(t, []string{"*"}, mask.Paths())
}
