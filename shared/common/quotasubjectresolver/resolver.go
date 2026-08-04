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
	ErrUnverifiedTransport = errors.New("verified transport is required")
	ErrCallerMismatch      = errors.New("caller cannot own the requested quota subject")
	ErrEvidenceRequired    = errors.New("authoritative subject evidence is required")
	ErrEvidenceMismatch    = errors.New("subject evidence does not match the caller")
	ErrSubjectConflict     = errors.New("quota subject conflicts with the existing subject")
)

/**
 * Grant records one authoritative subject assignment.
 */
type Grant struct {
	SubjectKind _quotasubject.Kind
	SubjectID   string
	GranteeKind _quotasubject.Kind
	GranteeID   string
	Authority   string
}

/**
 * Requirement asks for one exact quota subject kind.
 */
type Requirement struct {
	Kind  _quotasubject.Kind
	Grant *Grant
}

/**
 * Resolve selects one exact subject without fallback ordering.
 */
func Resolve(ctx context.Context, requirement Requirement) (_quotasubject.Subject, error) {
	transport, ok := _transporttrust.FromContext(ctx)
	if !ok || transport.Level != _transporttrust.LevelVerified {
		return _quotasubject.Subject{}, ErrUnverifiedTransport
	}

	if existing, ok := _quotasubject.FromContext(ctx); ok {
		if existing.Kind != requirement.Kind {
			return _quotasubject.Subject{}, ErrSubjectConflict
		}
		return existing, nil
	}

	caller, hasCaller := _callercontext.FromContext(ctx)
	actor, hasActor := _identitycontext.FromContext(ctx)

	switch requirement.Kind {
	case _quotasubject.KindUser:
		if !hasCaller || caller.Kind != _callercontext.KindUser || !hasActor || actor.ProfileID == 0 {
			return _quotasubject.Subject{}, ErrCallerMismatch
		}
		return _quotasubject.Subject{
			Kind:   _quotasubject.KindUser,
			ID:     strconv.FormatUint(actor.ProfileID, 10),
			Source: _quotasubject.SourceActorProfile,
		}, nil

	case _quotasubject.KindAPIKey:
		if !hasCaller || (caller.Kind != _callercontext.KindAPIKey && caller.Kind != _callercontext.KindUser) || caller.APIKeyID == 0 {
			return _quotasubject.Subject{}, ErrCallerMismatch
		}
		return _quotasubject.Subject{
			Kind:   _quotasubject.KindAPIKey,
			ID:     strconv.FormatUint(caller.APIKeyID, 10),
			Source: _quotasubject.SourceAPIKeyAuthority,
		}, nil

	case _quotasubject.KindInternalService:
		if !hasCaller || caller.Kind != _callercontext.KindInternalService {
			return _quotasubject.Subject{}, ErrCallerMismatch
		}
		return _quotasubject.Subject{
			Kind:   _quotasubject.KindInternalService,
			ID:     transport.ServiceID,
			Source: _quotasubject.SourceTransportService,
		}, nil

	case _quotasubject.KindOrganization,
		_quotasubject.KindProject,
		_quotasubject.KindSubscription:
		if requirement.Grant == nil {
			return _quotasubject.Subject{}, ErrEvidenceRequired
		}
		if err := validateGrant(*requirement.Grant, requirement.Kind); err != nil {
			return _quotasubject.Subject{}, err
		}
		if !grantMatchesCaller(ctx, *requirement.Grant) {
			return _quotasubject.Subject{}, ErrEvidenceMismatch
		}
		if requirement.Kind == _quotasubject.KindOrganization && hasActor && actor.OrganizationID != nil {
			assertedID := strconv.FormatUint(*actor.OrganizationID, 10)
			if assertedID != requirement.Grant.SubjectID {
				return _quotasubject.Subject{}, ErrEvidenceMismatch
			}
		}
		return subjectFromGrant(*requirement.Grant), nil

	case _quotasubject.KindAnonymous:
		if !hasCaller || caller.Kind != _callercontext.KindAnonymous {
			return _quotasubject.Subject{}, ErrCallerMismatch
		}
		if requirement.Grant == nil {
			return _quotasubject.Subject{}, ErrEvidenceRequired
		}
		grant := *requirement.Grant
		if err := validateGrant(grant, _quotasubject.KindAnonymous); err != nil {
			return _quotasubject.Subject{}, err
		}
		if grant.GranteeKind != "" || grant.GranteeID != "" || grant.Authority != transport.ServiceID {
			return _quotasubject.Subject{}, ErrEvidenceMismatch
		}
		return subjectFromGrant(grant), nil

	default:
		return _quotasubject.Subject{}, fmt.Errorf("unsupported quota subject kind %q", requirement.Kind)
	}
}

/**
 * ResolveAndAttach stores a successfully resolved subject in context.
 */
func ResolveAndAttach(ctx context.Context, requirement Requirement) (context.Context, _quotasubject.Subject, error) {
	subject, err := Resolve(ctx, requirement)
	if err != nil {
		return ctx, _quotasubject.Subject{}, err
	}
	return _quotasubject.WithSubject(ctx, subject), subject, nil
}

func validateGrant(grant Grant, expected _quotasubject.Kind) error {
	grant.SubjectID = strings.TrimSpace(grant.SubjectID)
	grant.GranteeID = strings.TrimSpace(grant.GranteeID)
	grant.Authority = strings.TrimSpace(grant.Authority)

	if grant.SubjectKind != expected || grant.SubjectID == "" || grant.Authority == "" {
		return ErrEvidenceMismatch
	}
	if !validText(grant.SubjectID) || !validText(grant.Authority) {
		return ErrEvidenceMismatch
	}

	if expected == _quotasubject.KindAnonymous {
		if grant.GranteeKind != "" || grant.GranteeID != "" {
			return ErrEvidenceMismatch
		}
		return nil
	}

	switch grant.GranteeKind {
	case _quotasubject.KindUser,
		_quotasubject.KindAPIKey,
		_quotasubject.KindInternalService:
		if !validText(grant.GranteeID) {
			return ErrEvidenceMismatch
		}
	default:
		return ErrEvidenceMismatch
	}
	return nil
}

func grantMatchesCaller(ctx context.Context, grant Grant) bool {
	if subject, ok := _quotasubject.FromContext(ctx); ok &&
		subject.Kind == grant.GranteeKind && subject.ID == grant.GranteeID {
		return true
	}

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
