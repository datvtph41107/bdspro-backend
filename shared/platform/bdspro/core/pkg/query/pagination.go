package query

import (
	"net/http"
	"strconv"
)

const (
	DefaultPage    = 1
	DefaultPerPage = 20
	MaxPerPage     = 100
)

// Pagination là value object — immutable sau khi Normalize.
type Pagination struct {
	Page    int `json:"page"    form:"page"    validate:"min=1"`
	PerPage int `json:"per_page" form:"per_page" validate:"min=1,max=100"`
}

// Normalize áp dụng default và clamp values.
// Luôn gọi trước khi dùng.
func (p *Pagination) Normalize() {
	if p.Page < 1 {
		p.Page = DefaultPage
	}
	if p.PerPage < 1 {
		p.PerPage = DefaultPerPage
	}
	if p.PerPage > MaxPerPage {
		p.PerPage = MaxPerPage
	}
}

func (p Pagination) Offset() int {
	return int((p.Page - 1) * p.PerPage)
}

func (p Pagination) Limit() int {
	return p.PerPage
}

// ParseFromRequest parses ?page=&per_page= với defaults.
func ParseFromRequest(r *http.Request) Pagination {
	p := Pagination{
		Page:    parseIntOr(r.URL.Query().Get("page"), DefaultPage),
		PerPage: parseIntOr(r.URL.Query().Get("per_page"), DefaultPerPage),
	}
	p.Normalize()
	return p
}

func parseIntOr(s string, def int) int {
	if v, err := strconv.Atoi(s); err == nil && v > 0 {
		return v
	}
	return def
}
