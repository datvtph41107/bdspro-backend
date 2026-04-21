package validate

import (
	"bdspro/core/pkg/validate"
	"testing"
)

func TestRequired(t *testing.T) {
	tests := []struct {
		name string
		val  *string
		want int
	}{
		{name: "nil", val: nil, want: 1},
		{name: "empty", val: ptr(""), want: 1},
		{name: "space", val: ptr("   "), want: 1},
		{name: "ok", val: ptr("abc"), want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var v validate.Violations
			validate.Required(&v, "title", tt.val)
			if got := v.Len(); got != tt.want {
				t.Fatalf("got %d violations, want %d", got, tt.want)
			}
		})
	}
}

func TestStringMaxLen(t *testing.T) {
	tests := []struct {
		name string
		val  *string
		max  int
		want int
	}{
		{name: "nil", val: nil, max: 3, want: 0},
		{name: "ok", val: ptr("abcd"), max: 4, want: 0},
		{name: "too long", val: ptr("abcde"), max: 4, want: 1},
		{name: "unicode count", val: ptr("你好"), max: 1, want: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var v validate.Violations
			validate.StringMaxLen(&v, "title", tt.val, tt.max)
			if got := v.Len(); got != tt.want {
				t.Fatalf("got %d violations, want %d", got, tt.want)
			}
		})
	}
}

func TestInt32Between(t *testing.T) {
	tests := []struct {
		name string
		val  *int32
		min  int32
		max  int32
		want int
	}{
		{name: "nil", val: nil, min: 0, max: 10, want: 0},
		{name: "ok", val: ptr(int32(5)), min: 0, max: 10, want: 0},
		{name: "low", val: ptr(int32(-1)), min: 0, max: 10, want: 1},
		{name: "high", val: ptr(int32(11)), min: 0, max: 10, want: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var v validate.Violations
			validate.Int32Between(&v, "page", tt.val, tt.min, tt.max)
			if got := v.Len(); got != tt.want {
				t.Fatalf("got %d violations, want %d", got, tt.want)
			}
		})
	}
}

func TestDateStringRequired(t *testing.T) {
	tests := []struct {
		name string
		val  *string
		want int
	}{
		{name: "nil", val: nil, want: 1},
		{name: "empty", val: ptr(""), want: 1},
		{name: "invalid", val: ptr("abc"), want: 1},
		{name: "ok", val: ptr("2024-01-15T00:00:00Z"), want: 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var v validate.Violations
			validate.DateStringRequired(&v, "published_at", tt.val)
			if got := v.Len(); got != tt.want {
				t.Fatalf("got %d violations, want %d", got, tt.want)
			}
		})
	}
}

func TestDateOrder(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		var v validate.Violations
		start := ptr("2024-01-01T00:00:00Z")
		end := ptr("2024-01-02T00:00:00Z")
		validate.DateOrder(&v, "start", "end", start, end)
		if v.HasErrors() {
			t.Fatalf("expected no violations, got %v", v.List())
		}
	})

	t.Run("invalid order", func(t *testing.T) {
		var v validate.Violations
		start := ptr("2024-01-02T00:00:00Z")
		end := ptr("2024-01-01T00:00:00Z")
		validate.DateOrder(&v, "start", "end", start, end)
		if v.Len() != 1 {
			t.Fatalf("got %d violations, want 1", v.Len())
		}
	})
}

func TestURLOpt(t *testing.T) {
	tests := []struct {
		name string
		val  *string
		want int
	}{
		{name: "nil", val: nil, want: 0},
		{name: "empty", val: ptr(""), want: 0},
		{name: "ok http", val: ptr("http://example.com"), want: 0},
		{name: "ok https", val: ptr("https://example.com"), want: 0},
		{name: "bad", val: ptr("ftp://example.com"), want: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var v validate.Violations
			validate.URLOpt(&v, "website", tt.val)
			if got := v.Len(); got != tt.want {
				t.Fatalf("got %d violations, want %d", got, tt.want)
			}
		})
	}
}

func TestMediaType(t *testing.T) {
	tests := []struct {
		name string
		val  *string
		want int
	}{
		{name: "nil", val: nil, want: 0},
		{name: "ok image", val: ptr("image"), want: 0},
		{name: "ok video", val: ptr("video"), want: 0},
		{name: "bad", val: ptr("audio"), want: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var v validate.Violations
			validate.MediaType(&v, "media_type", tt.val)
			if got := v.Len(); got != tt.want {
				t.Fatalf("got %d violations, want %d", got, tt.want)
			}
		})
	}
}
