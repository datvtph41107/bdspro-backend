package quotasubject

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"google.golang.org/grpc/metadata"
)

const (
	KindMetadataKey   = "x-qhpro-quota-subject-kind"
	IDMetadataKey     = "x-qhpro-quota-subject-id"
	SourceMetadataKey = "x-qhpro-quota-subject-source"
	MaxIDLength       = 128
)

type Kind string

const (
	KindUser            Kind = "user"
	KindOrganization    Kind = "organization"
	KindAPIKey          Kind = "api_key"
	KindProject         Kind = "project"
	KindSubscription    Kind = "subscription"
	KindAnonymous       Kind = "anonymous"
	KindInternalService Kind = "internal_service"
)

type Source string

const (
	SourceActorProfile           Source = "actor_profile"
	SourceAPIKeyAuthority        Source = "api_key_authority"
	SourceOrganizationMembership Source = "organization_membership"
	SourceProjectAccess          Source = "project_access"
	SourceSubscriptionOwnership  Source = "subscription_ownership"
	SourceAnonymousEdge          Source = "anonymous_edge"
	SourceTransportService       Source = "transport_service"
)

/**
 * Subject is the exact owner selected for one quota decision.
 */
type Subject struct {
	Kind   Kind
	ID     string
	Source Source
}

type contextKey struct{}

var (
	ErrInvalidSubject         = errors.New("quota subject is invalid")
	ErrSubjectContextConflict = errors.New("quota subject conflicts with the existing context")
)

/**
 * IsValid reports whether the subject is internally consistent.
 */
func (s Subject) IsValid() bool {
	s.ID = strings.TrimSpace(s.ID)
	if !validID(s.ID) {
		return false
	}

	switch s.Kind {
	case KindUser:
		return s.Source == SourceActorProfile
	case KindOrganization:
		return s.Source == SourceOrganizationMembership
	case KindAPIKey:
		return s.Source == SourceAPIKeyAuthority
	case KindProject:
		return s.Source == SourceProjectAccess
	case KindSubscription:
		return s.Source == SourceSubscriptionOwnership
	case KindAnonymous:
		return s.Source == SourceAnonymousEdge
	case KindInternalService:
		return s.Source == SourceTransportService
	default:
		return false
	}
}

/**
 * Key returns a stable namespace-qualified subject key.
 */
func (s Subject) Key() string {
	if !s.IsValid() {
		return ""
	}
	return fmt.Sprintf("%s:%s", s.Kind, strings.TrimSpace(s.ID))
}

/**
 * BindSubject validates and immutably binds one exact quota owner.
 */
func BindSubject(ctx context.Context, subject Subject) (context.Context, error) {
	if ctx == nil {
		return nil, errors.New("context is nil")
	}
	if !subject.IsValid() {
		return ctx, ErrInvalidSubject
	}
	subject.ID = strings.TrimSpace(subject.ID)
	if existing, exists := FromContext(ctx); exists {
		if existing == subject {
			return ctx, nil
		}
		return ctx, ErrSubjectContextConflict
	}
	return context.WithValue(ctx, contextKey{}, subject), nil
}

/**
 * FromContext returns the selected quota subject.
 */
func FromContext(ctx context.Context) (Subject, bool) {
	if ctx == nil {
		return Subject{}, false
	}
	subject, ok := ctx.Value(contextKey{}).(Subject)
	return subject, ok && subject.IsValid()
}

/**
 * AppendMetadata writes one selected subject to outgoing metadata.
 */
func AppendMetadata(md metadata.MD, subject Subject) {
	if md == nil || !subject.IsValid() {
		return
	}
	md.Set(KindMetadataKey, string(subject.Kind))
	md.Set(IDMetadataKey, strings.TrimSpace(subject.ID))
	md.Set(SourceMetadataKey, string(subject.Source))
}

/**
 * FromMetadata parses one transported quota subject.
 */
func FromMetadata(md metadata.MD) (Subject, bool, error) {
	if md == nil {
		return Subject{}, false, nil
	}

	kind, hasKind, err := singleOptional(md, KindMetadataKey)
	if err != nil {
		return Subject{}, false, err
	}
	id, hasID, err := singleOptional(md, IDMetadataKey)
	if err != nil {
		return Subject{}, false, err
	}
	source, hasSource, err := singleOptional(md, SourceMetadataKey)
	if err != nil {
		return Subject{}, false, err
	}

	if !hasKind && !hasID && !hasSource {
		return Subject{}, false, nil
	}
	if !hasKind || !hasID || !hasSource {
		return Subject{}, false, fmt.Errorf("quota subject metadata is incomplete")
	}

	subject := Subject{Kind: Kind(kind), ID: id, Source: Source(source)}
	if !subject.IsValid() {
		return Subject{}, false, fmt.Errorf("invalid quota subject metadata")
	}
	return subject, true, nil
}

func singleOptional(md metadata.MD, key string) (string, bool, error) {
	values := md.Get(key)
	if len(values) == 0 {
		return "", false, nil
	}
	if len(values) != 1 {
		return "", false, fmt.Errorf("duplicate %s metadata", key)
	}
	raw := values[0]
	value := strings.TrimSpace(raw)
	if value == "" {
		return "", false, fmt.Errorf("empty %s metadata", key)
	}
	if value != raw {
		return "", false, fmt.Errorf("invalid %s metadata", key)
	}
	return value, true, nil
}

func validID(value string) bool {
	if value == "" || len(value) > MaxIDLength {
		return false
	}
	for _, r := range value {
		if (r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') ||
			r == '-' || r == '_' || r == '.' {
			continue
		}
		return false
	}
	return true
}
