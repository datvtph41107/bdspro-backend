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
	RequestID   string            `json:"request_id,omitempty"`
	OperationID string            `json:"operation_id,omitempty"`
	Errors      []FieldProblem    `json:"errors,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// WriteProblem is the only canonical HTTP problem serializer. Callers decide
// semantics; this boundary owns correlation enrichment, media type and status.
func WriteProblem(ctx context.Context, w http.ResponseWriter, problem Problem) {
	if problem.Status < 400 || problem.Status > 599 {
		problem.Status = http.StatusInternalServerError
	}
	if strings.TrimSpace(problem.Type) == "" {
		problem.Type = "about:blank"
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
