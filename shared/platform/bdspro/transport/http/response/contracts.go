package http_resp

import "bdspro/core/pkg/query"

type Item[T any] struct {
	Data T        `json:"data"`
	Meta ItemMeta `json:"meta"`
}

type List[T any] struct {
	Data []T      `json:"data"`
	Meta ListMeta `json:"meta"`
}

func NewList[T any](items []T, p query.PageInfo, meta ItemMeta) List[T] {
	if items == nil {
		items = []T{}
	}
	return List[T]{
		Data: items,
		Meta: ListMeta{
			Page:       p.Page,
			PerPage:    p.PerPage,
			Total:      p.Total,
			TotalPages: p.TotalPages(),
			HasNext:    p.HasNext(),
			HasPrev:    p.HasPrev(),
		},
	}
}

type Async struct {
	Data AsyncTask `json:"data"`
	Meta ItemMeta  `json:"meta"`
}

type AsyncTask struct {
	TaskID       string `json:"task_id"`
	StatusURL    string `json:"status_url"`
	EstimatedSec int    `json:"estimated_sec,omitempty"`
}

// ─── Meta ─────────────────────────────────────────────────────────────────────

type ItemMeta struct {
	RequestID string `json:"request_id"`
	Timestamp string `json:"timestamp"`
}

type ListMeta struct {
	RequestID  string `json:"request_id"`
	Timestamp  string `json:"timestamp"`
	Page       int    `json:"page"`
	PerPage    int    `json:"per_page"`
	Total      int64  `json:"total"`
	TotalPages int    `json:"total_pages"`
	HasNext    bool   `json:"has_next"`
	HasPrev    bool   `json:"has_prev"`
}

// ─── Error ────────────────────────────────────────────────────────────────────

type ErrorEnvelope struct {
	Error ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code       string       `json:"code"`
	Message    string       `json:"message"`
	Fields     []FieldError `json:"fields,omitempty"`
	RetryAfter int          `json:"retry_after,omitempty"`
	RequestID  string       `json:"request_id"`
	DocURL     string       `json:"doc_url,omitempty"`
}

type FieldError struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}
