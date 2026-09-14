package httpresponse

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	_request "common/request"
)

const ProblemMediaType = "application/problem+json"

type FieldProblem struct {
	Field  string `json:"field"`
	Detail string `json:"detail"`
}

// Problem is the single public HTTP error shape. It follows RFC 9457 and uses
// extensions for stable application identity and request correlation.
type Problem struct {
	Type        string            `json:"type"`
	Title       string            `json:"title"`
	Status      int               `json:"status"`
	Detail      string            `json:"detail,omitempty"`
	Instance    string            `json:"instance,omitempty"`
	Code        string            `json:"code,omitempty"`
	Reason      string            `json:"reason,omitempty"`
	Domain      string            `json:"domain,omitempty"`
	RequestID   string            `json:"request_id,omitempty"`
	OperationID string            `json:"operation_id,omitempty"`
	Errors      []FieldProblem    `json:"errors,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// NewProblem constructs the canonical public error identity. Callers supply
// semantics (status, stable code and public detail); this package owns the
// RFC 9457 shape and type URI convention.
func NewProblem(status int, code, detail string) Problem {
	return Problem{
		Type:   TypeForCode(code),
		Title:  http.StatusText(status),
		Status: status,
		Detail: strings.TrimSpace(detail),
		Code:   strings.TrimSpace(code),
	}
}

// TypeForCode maps a stable application error code to the repository-owned
// problem type namespace. Unknown/empty codes intentionally fall back to
// about:blank instead of inventing an unstable identifier.
func TypeForCode(code string) string {
	code = strings.TrimSpace(code)
	if code == "" {
		return "about:blank"
	}
	var result strings.Builder
	for _, value := range code {
		switch {
		case value >= 'a' && value <= 'z',
			value >= 'A' && value <= 'Z',
			value >= '0' && value <= '9',
			value == '.', value == '-', value == '_', value == ':':
			result.WriteRune(value)
		default:
			result.WriteByte('-')
		}
	}
	return "urn:qhpro:error:" + result.String()
}

// WriteProblem is the only canonical HTTP problem serializer. Callers decide
// semantics; this boundary owns correlation enrichment, media type and status.
func WriteProblem(ctx context.Context, w http.ResponseWriter, problem Problem) {
	if problem.Status < 400 || problem.Status > 599 {
		problem.Status = http.StatusInternalServerError
	}
	if strings.TrimSpace(problem.Type) == "" {
		problem.Type = TypeForCode(problem.Code)
	}
	if strings.TrimSpace(problem.Title) == "" {
		problem.Title = http.StatusText(problem.Status)
	}

	if requestID, ok := _request.RequestIDFromContext(ctx); ok {
		problem.RequestID = requestID
		if problem.Instance == "" {
			problem.Instance = "urn:qhpro:request:" + requestID
		}
	}
	if operationID, ok := _request.OperationIDFromContext(ctx); ok {
		problem.OperationID = operationID
	}

	w.Header().Set("Content-Type", ProblemMediaType)
	w.WriteHeader(problem.Status)
	_ = json.NewEncoder(w).Encode(problem)
}
