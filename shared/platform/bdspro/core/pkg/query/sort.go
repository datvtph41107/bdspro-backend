package query

import (
	"net/http"
	"strings"
)

type SortDir string

const (
	SortAsc  SortDir = "asc"
	SortDesc SortDir = "desc"
)

type Sort struct {
	Field string
	Dir   SortDir
}

// ParseSort parses ?sort=name:asc hoặc ?sort=-name (prefix minus = desc).
func ParseSort(r *http.Request, allowed []string) (Sort, bool) {
	raw := strings.TrimSpace(r.URL.Query().Get("sort"))
	if raw == "" {
		return Sort{}, false
	}

	var field string
	dir := SortAsc

	if strings.HasPrefix(raw, "-") {
		field = raw[1:]
		dir = SortDesc
	} else if idx := strings.Index(raw, ":"); idx > 0 {
		field = raw[:idx]
		d := raw[idx+1:]
		if d == "desc" {
			dir = SortDesc
		}
	} else {
		field = raw
	}

	for _, a := range allowed {
		if a == field {
			return Sort{Field: field, Dir: dir}, true
		}
	}
	return Sort{}, false // unknown field — ignore silently OR return error
}

func (s Sort) IsZero() bool { return s.Field == "" }
func (s Sort) IsDesc() bool { return s.Dir == SortDesc }
func (s Sort) SQL() string {
	if s.IsZero() {
		return ""
	}
	// NOTE: field phải từ whitelist, không dùng trực tiếp từ user
	return s.Field + " " + string(s.Dir)
}
