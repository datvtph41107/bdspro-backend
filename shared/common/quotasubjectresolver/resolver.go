package quotasubjectresolver

import (
	_callercontext "common/callercontext"
	_identitycontext "common/identitycontext"
	_quotasubject "common/quotasubject"
	_transporttrust "common/transporttrust"

	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrUnverifiedTransport       = errors.New("verified transport is required")
	ErrCallerMismatch            = errors.New("caller cannot own the requested quota subject")
	ErrEvidenceRequired          = errors.New("authoritative subject evidence is required")
	ErrEvidenceMismatch          = errors.New("subject evidence does not match the caller or requirement")
	ErrSubjectConflict           = errors.New("quota subject conflicts with the existing subject")
	ErrInvalidRequirement        = errors.New("invalid quota subject requirement")
	ErrExpectedSubjectIDRequired = errors.New("expected quota subject ID is required for subject continuation")
)

/**
 * Grant records one request-scoped decision returned by an authoritative
 * domain lookup. EvidenceID is a traceable decision reference, not a secret.
 */
type Grant struct {
	SubjectKind _quotasubject.Kind
	SubjectID   string
	GranteeKind _quotasubject.Kind
	GranteeID   string
	Authority   string
	EvidenceID  string
}

/**
 * Requirement asks for one exact quota subject kind.
 *
 * ExpectedID binds the decision to the resource or account selected by the
 * operation. Authority names the domain service that must have produced a
 * supplied Grant. Neither field is inferred from arbitrary request metadata.
 */
type Requirement struct {
	Kind       _quotasubject.Kind
	ExpectedID string
	Authority  string
	Grant      *Grant
}

/**
 * Resolve selects one exact subject without fallback ordering.
 */
func Resolve(ctx context.Context, requirement Requirement) (_quotasubject.Subject, error) {
	transport, ok := _transporttrust.FromContext(ctx)
	if !ok || transport.Level != _transporttrust.LevelVerified {
		return _quotasubject.Subject{}, ErrUnverifiedTransport
	}

	requirement, err := normalizeRequirement(requirement)
	if err != nil {
		return _quotasubject.Subject{}, err
	}

	if existing, ok := _quotasubject.FromContext(ctx); ok {
		return resolveExisting(ctx, transport, existing, requirement)
	}

	return resolveFresh(ctx, transport, requirement)
}

/**
 * ResolveAndAttach stores a successfully resolved subject in context.
 */
func ResolveAndAttach(ctx context.Context, requirement Requirement) (context.Context, _quotasubject.Subject, error) {
	subject, err := Resolve(ctx, requirement)
	if err != nil {
		return ctx, _quotasubject.Subject{}, err
	}
	bound, err := _quotasubject.BindSubject(ctx, subject)
	if err != nil {
		return ctx, _quotasubject.Subject{}, err
	}
	return bound, subject, nil
}

func resolveExisting(
	ctx context.Context,
	transport _transporttrust.Caller,
	existing _quotasubject.Subject,
	requirement Requirement,
) (_quotasubject.Subject, error) {
	if existing.Kind != requirement.Kind {
		return _quotasubject.Subject{}, ErrSubjectConflict
	}
	if requirement.ExpectedID != "" && existing.ID != requirement.ExpectedID {
		return _quotasubject.Subject{}, ErrSubjectConflict
	}
	if requirement.Grant != nil {
		grantKind := requirement.Grant.SubjectKind
		grantID := strings.TrimSpace(requirement.Grant.SubjectID)
		if grantKind != existing.Kind || grantID != existing.ID {
			return _quotasubject.Subject{}, ErrSubjectConflict
		}
	}

	// A propagated resource-scoped subject is trustworthy only because it was
	// verified and signed on the incoming transport. The downstream operation
	// must still name the exact expected ID; kind-only continuation is too weak.
	if requirement.Grant == nil && requiresExpectedIDForContinuation(existing.Kind) && requirement.ExpectedID == "" {
		return _quotasubject.Subject{}, ErrExpectedSubjectIDRequired
	}

	switch existing.Kind {
	case _quotasubject.KindUser, _quotasubject.KindAPIKey:
		candidate, err := resolveFresh(ctx, transport, requirement)
		if err != nil {
			return _quotasubject.Subject{}, err
		}
		if candidate != existing {
			return _quotasubject.Subject{}, ErrSubjectConflict
		}
		return existing, nil

	case _quotasubject.KindInternalService:
		caller, ok := _callercontext.FromContext(ctx)
		if !ok || caller.Kind != _callercontext.KindInternalService {
			return _quotasubject.Subject{}, ErrCallerMismatch
		}
		return existing, nil

	case _quotasubject.KindOrganization,
		_quotasubject.KindProject,
		_quotasubject.KindSubscription,
		_quotasubject.KindAnonymous:
		caller, ok := _callercontext.FromContext(ctx)
		if !ok {
			return _quotasubject.Subject{}, ErrCallerMismatch
		}
		if existing.Kind == _quotasubject.KindAnonymous && caller.Kind != _callercontext.KindAnonymous {
			return _quotasubject.Subject{}, ErrCallerMismatch
		}
		if requirement.Grant == nil {
			return existing, nil
		}
		candidate, err := resolveFresh(ctx, transport, requirement)
		if err != nil {
			return _quotasubject.Subject{}, err
		}
		if candidate != existing {
			return _quotasubject.Subject{}, ErrSubjectConflict
		}
		return existing, nil

	default:
		return _quotasubject.Subject{}, ErrInvalidRequirement
	}
}

func resolveFresh(
	ctx context.Context,
	transport _transporttrust.Caller,
	requirement Requirement,
) (_quotasubject.Subject, error) {
	caller, hasCaller := _callercontext.FromContext(ctx)
	actor, hasActor := _identitycontext.FromContext(ctx)

	var subject _quotasubject.Subject

	switch requirement.Kind {
	case _quotasubject.KindUser:
		if requirement.Grant != nil || requirement.Authority != "" {
			return _quotasubject.Subject{}, ErrInvalidRequirement
		}
		if !hasCaller || caller.Kind != _callercontext.KindUser || !hasActor || actor.ProfileID == 0 {
			return _quotasubject.Subject{}, ErrCallerMismatch
		}
		subject = _quotasubject.Subject{
			Kind:   _quotasubject.KindUser,
			ID:     strconv.FormatUint(actor.ProfileID, 10),
			Source: _quotasubject.SourceActorProfile,
		}

	case _quotasubject.KindAPIKey:
		if requirement.Grant != nil || requirement.Authority != "" {
			return _quotasubject.Subject{}, ErrInvalidRequirement
		}
		if !hasCaller || (caller.Kind != _callercontext.KindAPIKey && caller.Kind != _callercontext.KindUser) || caller.APIKeyID == 0 {
			return _quotasubject.Subject{}, ErrCallerMismatch
		}
		subject = _quotasubject.Subject{
			Kind:   _quotasubject.KindAPIKey,
			ID:     strconv.FormatUint(caller.APIKeyID, 10),
			Source: _quotasubject.SourceAPIKeyAuthority,
		}

	case _quotasubject.KindInternalService:
		if requirement.Grant != nil || requirement.Authority != "" {
			return _quotasubject.Subject{}, ErrInvalidRequirement
		}
		if !hasCaller || caller.Kind != _callercontext.KindInternalService {
			return _quotasubject.Subject{}, ErrCallerMismatch
		}
		subject = _quotasubject.Subject{
			Kind:   _quotasubject.KindInternalService,
			ID:     transport.ServiceID,
			Source: _quotasubject.SourceTransportService,
		}

	case _quotasubject.KindOrganization,
		_quotasubject.KindProject,
		_quotasubject.KindSubscription:
		if requirement.Grant == nil {
			return _quotasubject.Subject{}, ErrEvidenceRequired
		}
		grant, err := normalizeAndValidateGrant(*requirement.Grant, requirement)
		if err != nil {
			return _quotasubject.Subject{}, err
		}
		if !grantMatchesCaller(ctx, grant) {
			return _quotasubject.Subject{}, ErrEvidenceMismatch
		}
		if requirement.Kind == _quotasubject.KindOrganization && hasActor && actor.OrganizationID != nil {
			assertedID := strconv.FormatUint(*actor.OrganizationID, 10)
			if assertedID != grant.SubjectID {
				return _quotasubject.Subject{}, ErrEvidenceMismatch
			}
		}
		subject = subjectFromGrant(grant)

	case _quotasubject.KindAnonymous:
		if !hasCaller || caller.Kind != _callercontext.KindAnonymous {
			return _quotasubject.Subject{}, ErrCallerMismatch
		}
		if requirement.Grant == nil {
			return _quotasubject.Subject{}, ErrEvidenceRequired
		}
		grant, err := normalizeAndValidateGrant(*requirement.Grant, requirement)
		if err != nil {
			return _quotasubject.Subject{}, err
		}
		if grant.GranteeKind != "" || grant.GranteeID != "" || grant.Authority != transport.ServiceID {
			return _quotasubject.Subject{}, ErrEvidenceMismatch
		}
		subject = subjectFromGrant(grant)

	default:
		return _quotasubject.Subject{}, fmt.Errorf("%w: unsupported quota subject kind %q", ErrInvalidRequirement, requirement.Kind)
	}

	if requirement.ExpectedID != "" && subject.ID != requirement.ExpectedID {
		return _quotasubject.Subject{}, ErrEvidenceMismatch
	}
	return subject, nil
}

func normalizeRequirement(requirement Requirement) (Requirement, error) {
	requirement.ExpectedID = strings.TrimSpace(requirement.ExpectedID)
	requirement.Authority = strings.TrimSpace(requirement.Authority)

	if requirement.ExpectedID != "" && !validText(requirement.ExpectedID) {
		return Requirement{}, ErrInvalidRequirement
	}
	if requirement.Authority != "" && !validText(requirement.Authority) {
		return Requirement{}, ErrInvalidRequirement
	}

	switch requirement.Kind {
	case _quotasubject.KindUser,
		_quotasubject.KindAPIKey,
		_quotasubject.KindOrganization,
		_quotasubject.KindProject,
		_quotasubject.KindSubscription,
		_quotasubject.KindAnonymous,
		_quotasubject.KindInternalService:
		return requirement, nil
	default:
		return Requirement{}, ErrInvalidRequirement
	}
}

func normalizeAndValidateGrant(grant Grant, requirement Requirement) (Grant, error) {
	grant.SubjectID = strings.TrimSpace(grant.SubjectID)
	grant.GranteeID = strings.TrimSpace(grant.GranteeID)
	grant.Authority = strings.TrimSpace(grant.Authority)
	grant.EvidenceID = strings.TrimSpace(grant.EvidenceID)

	if requirement.Authority == "" || grant.Authority != requirement.Authority {
		return Grant{}, ErrEvidenceMismatch
	}
	if grant.SubjectKind != requirement.Kind || grant.SubjectID == "" || grant.EvidenceID == "" {
		return Grant{}, ErrEvidenceMismatch
	}
	if requirement.ExpectedID != "" && grant.SubjectID != requirement.ExpectedID {
		return Grant{}, ErrEvidenceMismatch
	}
	if !validText(grant.SubjectID) || !validText(grant.Authority) || !validText(grant.EvidenceID) {
		return Grant{}, ErrEvidenceMismatch
	}

	if requirement.Kind == _quotasubject.KindAnonymous {
		if grant.GranteeKind != "" || grant.GranteeID != "" {
			return Grant{}, ErrEvidenceMismatch
		}
		return grant, nil
	}

	switch grant.GranteeKind {
	case _quotasubject.KindUser,
		_quotasubject.KindAPIKey,
		_quotasubject.KindInternalService:
		if !validText(grant.GranteeID) {
			return Grant{}, ErrEvidenceMismatch
		}
	default:
		return Grant{}, ErrEvidenceMismatch
	}
	return grant, nil
}

func grantMatchesCaller(ctx context.Context, grant Grant) bool {
	caller, hasCaller := _callercontext.FromContext(ctx)
	actor, hasActor := _identitycontext.FromContext(ctx)
	transport, hasTransport := _transporttrust.FromContext(ctx)

	switch grant.GranteeKind {
	case _quotasubject.KindUser:
		return hasCaller && caller.Kind == _callercontext.KindUser &&
			hasActor && actor.ProfileID > 0 &&
			grant.GranteeID == strconv.FormatUint(actor.ProfileID, 10)
	case _quotasubject.KindAPIKey:
		return hasCaller && caller.APIKeyID > 0 &&
			grant.GranteeID == strconv.FormatUint(caller.APIKeyID, 10)
	case _quotasubject.KindInternalService:
		return hasCaller && caller.Kind == _callercontext.KindInternalService &&
			hasTransport && transport.Level == _transporttrust.LevelVerified &&
			grant.GranteeID == transport.ServiceID
	default:
		return false
	}
}

func subjectFromGrant(grant Grant) _quotasubject.Subject {
	source := _quotasubject.SourceOrganizationMembership
	switch grant.SubjectKind {
	case _quotasubject.KindProject:
		source = _quotasubject.SourceProjectAccess
	case _quotasubject.KindSubscription:
		source = _quotasubject.SourceSubscriptionOwnership
	case _quotasubject.KindAnonymous:
		source = _quotasubject.SourceAnonymousEdge
	}
	return _quotasubject.Subject{
		Kind:   grant.SubjectKind,
		ID:     strings.TrimSpace(grant.SubjectID),
		Source: source,
	}
}

func requiresExpectedIDForContinuation(kind _quotasubject.Kind) bool {
	switch kind {
	case _quotasubject.KindOrganization,
		_quotasubject.KindProject,
		_quotasubject.KindSubscription,
		_quotasubject.KindAnonymous,
		_quotasubject.KindInternalService:
		return true
	default:
		return false
	}
}

func validText(value string) bool {
	value = strings.TrimSpace(value)
	if value == "" || len(value) > _quotasubject.MaxIDLength {
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
