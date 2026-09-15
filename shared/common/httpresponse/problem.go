package httpresponse

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"common/fault"
	_request "common/request"
)

const (
	ProblemMediaType    = "application/problem+json"
	LegacyJSONMediaType = "application/json"
)

type FieldProblem struct {
	Field  string `json:"field"`
	Detail string `json:"detail"`
}

// Problem is the single public HTTP error representation.
// It follows RFC 9457 and uses extensions for stable application identity
// and request correlation.
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

type writeConfig struct {
	legacyJSONEnvelope bool
}

// WriteOption configures representation at the HTTP boundary without
// changing the canonical Problem semantics.
type WriteOption func(*writeConfig)

// WithLegacyJSONEnvelope preserves the historical {code,message} contract
// while adding stable application identity for gradual client migration.
func WithLegacyJSONEnvelope() WriteOption {
	return func(cfg *writeConfig) {
		cfg.legacyJSONEnvelope = true
	}
}

type legacyJSONEnvelope struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	ErrorCode   string `json:"error_code,omitempty"`
	Status      int    `json:"status"`
	RequestID   string `json:"request_id,omitempty"`
	OperationID string `json:"operation_id,omitempty"`
}

// NewProblem constructs one canonical public HTTP problem.
func NewProblem(
	status int,
	code string,
	detail string,
) Problem {
	return Problem{
		Type:   TypeForCode(code),
		Title:  http.StatusText(status),
		Status: status,
		Detail: strings.TrimSpace(detail),
		Code:   strings.TrimSpace(code),
	}
}

// TypeForCode maps stable application identity into the repository-owned
// RFC 9457 problem type namespace.
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
			value == '.',
			value == '-',
			value == '_',
			value == ':':
			result.WriteRune(value)

		default:
			result.WriteByte('-')
		}
	}

	return "urn:qhpro:error:" + result.String()
}

// StatusForKind owns the transport-neutral fault -> HTTP mapping.
func StatusForKind(kind fault.Kind) int {
	switch kind {
	case fault.KindValidation:
		return http.StatusBadRequest

	case fault.KindUnauthenticated:
		return http.StatusUnauthorized

	case fault.KindPermissionDenied:
		return http.StatusForbidden

	case fault.KindNotFound:
		return http.StatusNotFound

	case fault.KindConflict:
		return http.StatusConflict

	case fault.KindPrecondition:
		return http.StatusPreconditionFailed

	case fault.KindAborted:
		return http.StatusConflict

	case fault.KindResourceExhausted:
		return http.StatusTooManyRequests

	case fault.KindUnavailable:
		return http.StatusServiceUnavailable

	case fault.KindInternal:
		return http.StatusInternalServerError

	default:
		return http.StatusInternalServerError
	}
}

// ProblemFromError projects canonical error identity into public HTTP
// representation. Unknown errors fail closed and never expose their cause.
func ProblemFromError(err error) Problem {
	failure, ok := fault.As(err)

	if !ok {
		failure = fault.Wrap(
			err,
			fault.KindInternal,
			"common.internal",
			"internal server error",
		)
	}

	problem := NewProblem(
		StatusForKind(failure.Kind()),
		failure.Code(),
		failure.PublicMessage(),
	)

	for _, violation := range failure.Violations() {
		problem.Errors = append(
			problem.Errors,
			FieldProblem{
				Field:  violation.Field,
				Detail: violation.Description,
			},
		)
	}

	problem.Metadata = failure.Metadata()

	return problem
}

// WriteProblem is the ONE canonical HTTP problem serializer.
// Semantic decisions happen before this boundary. Legacy representation is
// an explicit migration mode owned here rather than by middleware callers.
func WriteProblem(
	ctx context.Context,
	w http.ResponseWriter,
	problem Problem,
	options ...WriteOption,
) {
	cfg := writeConfig{}

	for _, option := range options {
		if option != nil {
			option(&cfg)
		}
	}

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
			problem.Instance =
				"urn:qhpro:request:" + requestID
		}
	}

	if operationID, ok :=
		_request.OperationIDFromContext(ctx); ok {
		problem.OperationID = operationID
	}

	if cfg.legacyJSONEnvelope {
		w.Header().Set(
			"Content-Type",
			LegacyJSONMediaType,
		)

		w.WriteHeader(problem.Status)

		_ = json.NewEncoder(w).Encode(
			legacyJSONEnvelope{
				Code:        problem.Status,
				Message:     problem.Detail,
				ErrorCode:   problem.Code,
				Status:      problem.Status,
				RequestID:   problem.RequestID,
				OperationID: problem.OperationID,
			},
		)

		return
	}

	w.Header().Set(
		"Content-Type",
		ProblemMediaType,
	)

	w.WriteHeader(problem.Status)

	_ = json.NewEncoder(w).Encode(problem)
}
