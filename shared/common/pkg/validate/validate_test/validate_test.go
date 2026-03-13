package validate_test

import (
	"common/pkg/validate"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func pStr(s string) *string   { return &s }
func pI32(i int32) *int32     { return &i }
func pU64(u uint64) *uint64   { return &u }
func pF64(f float64) *float64 { return &f }
func pU32(u uint32) *uint32   { return &u }

// ─── Error ────────────────────────────────────────────────────────────────────

func TestError_ZeroValue_NoErrors(t *testing.T) {
	e := validate.New()
	assert.False(t, e.HasErrors())
	assert.Nil(t, e.AsError())
	assert.Equal(t, "", e.Error())
}

func TestError_Add_HasErrors(t *testing.T) {
	e := validate.New()
	e.Add("title", "required")
	assert.True(t, e.HasErrors())
	assert.NotNil(t, e.AsError())
	assert.Contains(t, e.Error(), "title: required")
}

func TestError_Merge(t *testing.T) {
	e1 := validate.New()
	e1.Add("field1", "error1")

	e2 := validate.New()
	e2.Add("field2", "error2")

	e1.Merge(e2)

	msg := e1.Error()
	assert.Contains(t, msg, "field1: error1")
	assert.Contains(t, msg, "field2: error2")
}

func TestError_Merge_Nil_NoOp(t *testing.T) {
	e := validate.New()
	e.Add("f", "m")
	e.Merge(nil)
	assert.Len(t, e.Fields(), 1)
}

func TestError_Merge_EmptyOther_NoOp(t *testing.T) {
	e := validate.New()
	e.Add("f", "m")
	e.Merge(validate.New())
	assert.Len(t, e.Fields(), 1)
}

func TestError_Fields(t *testing.T) {
	e := validate.New()
	e.Add("title", "required")
	e.Add("scope", "invalid")

	fields := e.Fields()
	assert.Len(t, fields, 2)
	assert.Equal(t, "title", fields[0].Field)
	assert.Equal(t, "required", fields[0].Message)
}

// Error implement error interface
func TestError_ImplementsError(t *testing.T) {
	var err error = validate.New()
	assert.NotNil(t, err) // compiles = implements error
}

// ─── Required ─────────────────────────────────────────────────────────────────

func TestRequired(t *testing.T) {
	tests := []struct {
		name    string
		v       *string
		wantErr bool
	}{
		{"nil → error", nil, true},
		{"empty → error", pStr(""), true},
		{"spaces only → error", pStr("  "), true},
		{"value → ok", pStr("x"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := validate.New()
			validate.Required(e, "f", tt.v)
			assert.Equal(t, tt.wantErr, e.HasErrors(), e.Error())
		})
	}
}

func TestRequiredID(t *testing.T) {
	tests := []struct {
		name    string
		v       *uint64
		wantErr bool
	}{
		{"nil → error", nil, true},
		{"0 → error", pU64(0), true},
		{">0 → ok", pU64(1), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := validate.New()
			validate.RequiredID(e, "id", tt.v)
			assert.Equal(t, tt.wantErr, e.HasErrors())
		})
	}
}

// ─── StringMaxLen ─────────────────────────────────────────────────────────────

func TestStringMaxLen(t *testing.T) {
	tests := []struct {
		name    string
		v       *string
		max     int
		wantErr bool
	}{
		{"nil → skip", nil, 5, false},
		{"len < max → ok", pStr("abc"), 5, false},
		{"len == max → ok", pStr("abcde"), 5, false},
		{"len > max → error", pStr("abcdef"), 5, true},
		{"empty <= max → ok", pStr(""), 5, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := validate.New()
			validate.StringMaxLen(e, "f", tt.v, tt.max)
			assert.Equal(t, tt.wantErr, e.HasErrors())
		})
	}
}

// ─── Int32OneOf ───────────────────────────────────────────────────────────────

func TestInt32OneOf(t *testing.T) {
	allowed := []int32{1, 2, 3}
	tests := []struct {
		name    string
		v       *int32
		wantErr bool
	}{
		{"nil → skip", nil, false},
		{"valid → ok", pI32(1), false},
		{"valid zero... wait", pI32(0), true}, // 0 không trong [1,2,3]
		{"invalid → error", pI32(5), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := validate.New()
			validate.Int32OneOf(e, "scope", tt.v, allowed)
			assert.Equal(t, tt.wantErr, e.HasErrors())
		})
	}
}

// ─── Float64NonNegative ───────────────────────────────────────────────────────

func TestFloat64NonNegative(t *testing.T) {
	tests := []struct {
		name    string
		v       *float64
		wantErr bool
	}{
		{"nil → skip", nil, false},
		{"0.0 → ok", pF64(0.0), false},
		{"positive → ok", pF64(85.5), false},
		{"negative → error", pF64(-0.1), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := validate.New()
			validate.Float64NonNegative(e, "area", tt.v)
			assert.Equal(t, tt.wantErr, e.HasErrors())
		})
	}
}

// ─── DateStringOpt ────────────────────────────────────────────────────────────

func TestDateStringOpt(t *testing.T) {
	tests := []struct {
		name    string
		v       *string
		wantErr bool
	}{
		{"nil → skip", nil, false},
		{"\"\" → skip (clear intent ok)", pStr(""), false},
		{"valid RFC3339 → ok", pStr("2024-01-15T00:00:00Z"), false},
		{"date only no time → error", pStr("2024-01-15"), true},
		{"random string → error", pStr("not-a-date"), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := validate.New()
			validate.DateStringOpt(e, "expired_land", tt.v)
			assert.Equal(t, tt.wantErr, e.HasErrors())
		})
	}
}

// ─── DateOrder ────────────────────────────────────────────────────────────────

func TestDateOrder(t *testing.T) {
	tests := []struct {
		name    string
		start   *string
		end     *string
		wantErr bool
	}{
		{"both nil → skip", nil, nil, false},
		{"start nil → skip", nil, pStr("2025-01-01T00:00:00Z"), false},
		{"end empty → skip", pStr("2024-01-01T00:00:00Z"), pStr(""), false},
		{"start < end → ok", pStr("2024-01-01T00:00:00Z"), pStr("2025-01-01T00:00:00Z"), false},
		{"start == end → error", pStr("2024-01-01T00:00:00Z"), pStr("2024-01-01T00:00:00Z"), true},
		{"start > end → error", pStr("2025-01-01T00:00:00Z"), pStr("2024-01-01T00:00:00Z"), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := validate.New()
			validate.DateOrder(e, "start", "end", tt.start, tt.end)
			assert.Equal(t, tt.wantErr, e.HasErrors())
		})
	}
}

// ─── Geo ──────────────────────────────────────────────────────────────────────

func TestLatitudeLongitude(t *testing.T) {
	tests := []struct {
		name    string
		lat     *float64
		lng     *float64
		wantErr bool
	}{
		{"nil → skip", nil, nil, false},
		{"valid → ok", pF64(10.5), pF64(106.7), false},
		{"lat out of range → error", pF64(91.0), pF64(0), true},
		{"lng out of range → error", pF64(0), pF64(181.0), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := validate.New()
			validate.Latitude(e, "lat", tt.lat)
			validate.Longitude(e, "lng", tt.lng)
			assert.Equal(t, tt.wantErr, e.HasErrors())
		})
	}
}

// ─── Multi-field accumulation ─────────────────────────────────────────────────

func TestMultipleErrors_AllAccumulated(t *testing.T) {
	e := validate.New()
	validate.Required(e, "title", nil)
	validate.Float64NonNegative(e, "area", pF64(-1))
	validate.DateStringOpt(e, "expired", pStr("bad-date"))

	assert.True(t, e.HasErrors())
	msg := e.Error()
	// Tất cả lỗi đều có mặt — không dừng ở lỗi đầu
	assert.True(t, strings.Contains(msg, "title"))
	assert.True(t, strings.Contains(msg, "area"))
	assert.True(t, strings.Contains(msg, "expired"))
}
