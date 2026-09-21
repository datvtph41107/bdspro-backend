package httpresponse

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	_errors "common/errors"
	_request "common/request"

	"google.golang.org/grpc/codes"
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
	legacyJSONEnvelope     bool
	legacyDataJSONEnvelope bool
}

type WriteOption func(*writeConfig)

func WithLegacyJSONEnvelope() WriteOption {
	return func(cfg *writeConfig) {
		cfg.legacyJSONEnvelope = true
	}
}

func WithLegacyDataJSONEnvelope() WriteOption {
	return func(cfg *writeConfig) {
		cfg.legacyDataJSONEnvelope = true
	}
}

type legacyDataJSONEnvelope struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Data    struct{} `json:"data"`
}

type legacyJSONEnvelope struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	ErrorCode   string `json:"error_code,omitempty"`
	Status      int    `json:"status"`
	RequestID   string `json:"request_id,omitempty"`
	OperationID string `json:"operation_id,omitempty"`
}

func NewProblem(status int, code string, detail string) Problem {
	return Problem{
		Type:   TypeForCode(code),
		Title:  http.StatusText(status),
		Status: status,
		Detail: strings.TrimSpace(detail),
		Code:   strings.TrimSpace(code),
	}
}

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

// StatusForRPC is the single canonical semantic RPC -> HTTP projection for
// direct HTTP boundaries. Gateway owns the same transport decision for gRPC.
func StatusForRPC(code codes.Code) int {
	switch code {
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.NotFound:
		return http.StatusNotFound
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.FailedPrecondition:
		return http.StatusPreconditionFailed
	case codes.Aborted:
		return http.StatusConflict
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}

// ProblemFromError projects canonical application identity into the public
// HTTP representation. Unknown technical failures deliberately have no
// application identity and fail closed.
func ProblemFromError(err error) Problem {
	application, ok := _errors.As(err)
	if !ok {
		return NewProblem(
			http.StatusInternalServerError,
			"",
			"internal server error",
		)
	}

	spec := application.Spec()
	publicCode := spec.LegacyProblemCode()
	if publicCode == "" {
		publicCode = strconv.FormatInt(int64(spec.Code()), 10)
	}

	problem := NewProblem(
		StatusForRPC(application.RPCCode()),
		publicCode,
		application.PublicMessage(),
	)
	problem.Reason = string(application.Key())
	problem.Domain = _errors.ErrorDomain

	for _, violation := range application.Violations() {
		problem.Errors = append(problem.Errors, FieldProblem{
			Field:  violation.Field,
			Detail: violation.Description,
		})
	}
	problem.Metadata = application.Metadata()
	return problem
}

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
			problem.Instance = "urn:qhpro:request:" + requestID
		}
	}
	if operationID, ok := _request.OperationIDFromContext(ctx); ok {
		problem.OperationID = operationID
	}

	if cfg.legacyDataJSONEnvelope {
		w.Header().Set("Content-Type", LegacyJSONMediaType)
		w.WriteHeader(problem.Status)
		_ = json.NewEncoder(w).Encode(legacyDataJSONEnvelope{
			Code:    problem.Status,
			Message: problem.Detail,
		})
		return
	}

	if cfg.legacyJSONEnvelope {
		w.Header().Set("Content-Type", LegacyJSONMediaType)
		w.WriteHeader(problem.Status)
		_ = json.NewEncoder(w).Encode(legacyJSONEnvelope{
			Code:        problem.Status,
			Message:     problem.Detail,
			ErrorCode:   problem.Code,
			Status:      problem.Status,
			RequestID:   problem.RequestID,
			OperationID: problem.OperationID,
		})
		return
	}

	w.Header().Set("Content-Type", ProblemMediaType)
	w.WriteHeader(problem.Status)
	_ = json.NewEncoder(w).Encode(problem)
}
