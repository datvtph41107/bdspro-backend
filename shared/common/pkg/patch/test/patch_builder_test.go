package patch_test

import (
	"common/pkg/fieldmask"
	"common/pkg/patch"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ── helpers ───────────────────────────────────────────────────────────────────

func pStr(s string) *string   { return &s }
func pI32(i int32) *int32     { return &i }
func pU64(u uint64) *uint64   { return &u }
func pF64(f float64) *float64 { return &f }
func pBool(b bool) *bool      { return &b }
func pU32(u uint32) *uint32   { return &u }

func mustTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

var (
	allowAll = fieldmask.FieldMask{} // zero value = allowAll
	denyAll  = fieldmask.Only()      // no paths = ? — test below
)

// ─── New vs WithMask ──────────────────────────────────────────────────────────

func TestNew_AllowAll(t *testing.T) {
	// New() = không mask → allow all
	got := patch.New().
		String("any.path", "col", pStr("val"), "old").
		Build()
	assert.Equal(t, map[string]any{"col": "val"}, got)
}

func TestWithMask_OnlyAllowedPaths(t *testing.T) {
	mask := fieldmask.Only("info.title")
	got := patch.WithMask(mask).
		String("info.title", "title", pStr("new"), "old").
		String("info.note", "note", pStr("note"), ""). // not in mask
		Build()
	assert.Equal(t, map[string]any{"title": "new"}, got)
}

// ─── Has ─────────────────────────────────────────────────────────────────────

func TestHas_Empty(t *testing.T) {
	b := patch.New().String("x", "col", nil, "v")
	assert.False(t, b.Has())
}

func TestHas_NonEmpty(t *testing.T) {
	b := patch.New().String("x", "col", pStr("new"), "old")
	assert.True(t, b.Has())
}

// ─── String ───────────────────────────────────────────────────────────────────

func TestString(t *testing.T) {
	tests := []struct {
		name string
		in   *string
		ex   string
		want map[string]any
	}{
		{"nil=absent → skip", nil, "old", map[string]any{}},
		{"same value → skip (diff)", pStr("v"), "v", map[string]any{}},
		{"new value → update", pStr("new"), "old", map[string]any{"c": "new"}},
		{"\"\" là giá trị hợp lệ", pStr(""), "old", map[string]any{"c": ""}},
		{"\"\" same as existing → skip", pStr(""), "", map[string]any{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := patch.New().String("p", "c", tt.in, tt.ex).Build()
			assert.Equal(t, tt.want, got)
		})
	}
}

// ─── Int32 ────────────────────────────────────────────────────────────────────

func TestInt32(t *testing.T) {
	tests := []struct {
		name string
		in   *int32
		ex   int32
		want map[string]any
	}{
		{"nil → skip", nil, 1, map[string]any{}},
		{"same → skip", pI32(1), 1, map[string]any{}},
		{"0 là giá trị hợp lệ → update", pI32(0), 1, map[string]any{"c": int32(0)}},
		{"0 → 0 same → skip", pI32(0), 0, map[string]any{}},
		{"new value → update", pI32(2), 1, map[string]any{"c": int32(2)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := patch.New().Int32("p", "c", tt.in, tt.ex).Build()
			assert.Equal(t, tt.want, got)
		})
	}
}

// ─── Bool ─────────────────────────────────────────────────────────────────────

func TestBool(t *testing.T) {
	tests := []struct {
		name string
		in   *bool
		ex   bool
		want map[string]any
	}{
		{"nil → skip", nil, true, map[string]any{}},
		{"same true → skip", pBool(true), true, map[string]any{}},
		{"false là giá trị hợp lệ", pBool(false), true, map[string]any{"c": false}},
		{"false → false same → skip", pBool(false), false, map[string]any{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := patch.New().Bool("p", "c", tt.in, tt.ex).Build()
			assert.Equal(t, tt.want, got)
		})
	}
}

// ─── ForeignKey ───────────────────────────────────────────────────────────────

func TestForeignKey(t *testing.T) {
	tests := []struct {
		name string
		in   *uint64
		ex   *uint64
		want map[string]any
	}{
		{"nil=absent → skip", nil, pU64(5), map[string]any{}},
		{"0=clear, DB có giá trị → NULL", pU64(0), pU64(5), map[string]any{"c": nil}},
		{"0=clear, DB đã NULL → skip (no-op)", pU64(0), nil, map[string]any{}},
		{"set new value", pU64(9), pU64(5), map[string]any{"c": uint64(9)}},
		{"same value → skip (diff)", pU64(5), pU64(5), map[string]any{}},
		{"set value, DB đang NULL", pU64(9), nil, map[string]any{"c": uint64(9)}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := patch.New().ForeignKey("p", "c", tt.in, tt.ex).Build()
			assert.Equal(t, tt.want, got)
		})
	}
}

// ─── NullableFloat ────────────────────────────────────────────────────────────

func TestNullableFloat(t *testing.T) {
	tests := []struct {
		name string
		in   *float64
		ex   *float64
		want map[string]any
	}{
		{"nil → skip", nil, pF64(85.5), map[string]any{}},
		{"same value → skip", pF64(85.5), pF64(85.5), map[string]any{}},
		{"0.0 là diện tích hợp lệ", pF64(0.0), pF64(85.5), map[string]any{"c": 0.0}},
		{"new value → update", pF64(90.0), pF64(85.5), map[string]any{"c": 90.0}},
		{"set value, DB đang NULL", pF64(85.5), nil, map[string]any{"c": 85.5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := patch.New().NullableFloat("p", "c", tt.in, tt.ex).Build()
			assert.Equal(t, tt.want, got)
		})
	}
}

// ─── DateString ───────────────────────────────────────────────────────────────

func TestDateString(t *testing.T) {
	t1 := mustTime("2025-01-01T00:00:00Z")
	t2 := mustTime("2026-01-01T00:00:00Z")

	tests := []struct {
		name    string
		in      *string
		ex      *time.Time
		want    map[string]any
		wantErr bool
	}{
		{
			name: "nil → skip",
			in:   nil, ex: &t1,
			want: map[string]any{},
		},
		{
			name: "\"\"=clear, DB có date → NULL",
			in:   pStr(""), ex: &t1,
			want: map[string]any{"c": nil},
		},
		{
			name: "\"\"=clear, DB đã NULL → skip",
			in:   pStr(""), ex: nil,
			want: map[string]any{},
		},
		{
			name: "same date → skip",
			in:   pStr("2025-01-01T00:00:00Z"), ex: &t1,
			want: map[string]any{},
		},
		{
			name: "new date → update",
			in:   pStr("2026-01-01T00:00:00Z"), ex: &t1,
			want: map[string]any{"c": t2},
		},
		{
			name: "set date, DB đang NULL",
			in:   pStr("2025-01-01T00:00:00Z"), ex: nil,
			want: map[string]any{"c": t1},
		},
		{
			name: "invalid date → error propagated, không panic",
			in:   pStr("bad-date"), ex: &t1,
			want:    map[string]any{}, // không emit field khi error
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := patch.New().DateString("p", "c", tt.in, tt.ex).Result()
			assert.Equal(t, tt.want, m)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// ─── FieldMask integration ────────────────────────────────────────────────────

func TestFieldMask_Granular_Update(t *testing.T) {
	// Scenario: client chỉ muốn update title và area_total
	mask := fieldmask.Only("info.title", "landInfo.area.area_total")

	got := patch.WithMask(mask).
		String("info.title", "title", pStr("new"), "old").                               // allowed
		String("info.note", "note", pStr("note"), "").                                   // not in mask
		Int32("info.scope", "scope", pI32(2), 1).                                        // not in mask
		NullableFloat("landInfo.area.area_total", "area_total", pF64(90.0), pF64(85.0)). // allowed
		NullableFloat("landInfo.area.area_land", "area_land", pF64(50.0), pF64(40.0)).   // not in mask
		Build()

	assert.Equal(t, map[string]any{
		"title":      "new",
		"area_total": 90.0,
	}, got)
}

func TestFieldMask_PrefixMatch(t *testing.T) {
	// "landInfo.area" cover tất cả children
	mask := fieldmask.Only("landInfo.area")

	got := patch.WithMask(mask).
		NullableFloat("landInfo.area.area_total", "area_total", pF64(90.0), nil).
		NullableFloat("landInfo.area.area_land", "area_land", pF64(50.0), nil).
		NullableFloat("landInfo.expiry.expired", "expired", pF64(1.0), nil). // not under area
		Build()

	assert.Equal(t, map[string]any{
		"area_total": 90.0,
		"area_land":  50.0,
		// expired không có
	}, got)
}

// ─── Chain all types ──────────────────────────────────────────────────────────

func TestChain_AllTypes_ZeroValues(t *testing.T) {
	// Verify: zero value của mỗi type đều được update khi có trong mask
	// và thực sự khác existing
	got := patch.New().
		String("a", "title", pStr(""), "old").            // "" hợp lệ
		Int32("b", "scope", pI32(0), 1).                  // 0 hợp lệ
		Bool("c", "active", pBool(false), true).          // false hợp lệ
		ForeignKey("d", "proj_id", pU64(0), pU64(5)).     // 0 = clear
		NullableFloat("e", "area", pF64(0.0), pF64(1.0)). // 0.0 hợp lệ
		Build()

	assert.Equal(t, map[string]any{
		"title":   "",
		"scope":   int32(0),
		"active":  false,
		"proj_id": nil, // clear FK
		"area":    0.0,
	}, got)
}

// ─── Result vs Build ─────────────────────────────────────────────────────────

func TestResult_NoError(t *testing.T) {
	m, err := patch.New().
		String("p", "c", pStr("new"), "old").
		Result()
	require.NoError(t, err)
	assert.Equal(t, map[string]any{"c": "new"}, m)
}

func TestResult_WithDateError_PropagatesError(t *testing.T) {
	m, err := patch.New().
		String("p1", "title", pStr("ok"), "old").
		DateString("p2", "date", pStr("BAD"), nil).
		Result()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "p2") // path trong error message
	// map vẫn có field hợp lệ
	assert.Equal(t, map[string]any{"title": "ok"}, m)
}
