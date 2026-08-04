package quotasubjectresolver

import (
	_callercontext "common/callercontext"
	_identitycontext "common/identitycontext"
	_quotasubject "common/quotasubject"
	_transporttrust "common/transporttrust"

	"context"
	"errors"
	"testing"
	"time"
)

func verifiedContext(serviceID string) context.Context {
	return _transporttrust.WithCaller(context.Background(), _transporttrust.Caller{
		ServiceID: serviceID,
		Level:     _transporttrust.LevelVerified,
		IssuedAt:  time.Now(),
	})
}

func TestResolveUserUsesProfileIDWithoutFallback(t *testing.T) {
	t.Parallel()

	ctx := verifiedContext("gateway-service")
	ctx = _callercontext.WithCaller(ctx, _callercontext.Caller{Kind: _callercontext.KindUser})
	ctx = _identitycontext.WithActor(ctx, _identitycontext.Actor{AuthID: 7, ProfileID: 42, OriginID: 99})

	subject, err := Resolve(ctx, Requirement{Kind: _quotasubject.KindUser})
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}
	if subject.Key() != "user:42" {
		t.Fatalf("subject = %+v, want user:42", subject)
	}
}

func TestResolveAPIKeyForAPIKeyOrUserCaller(t *testing.T) {
	t.Parallel()

	for _, kind := range []_callercontext.Kind{_callercontext.KindAPIKey, _callercontext.KindUser} {
		ctx := verifiedContext("gateway-service")
		ctx = _callercontext.WithCaller(ctx, _callercontext.Caller{
			Kind: kind, APIKeyID: 88, APIKeyApp: "partner-app",
		})
		if kind == _callercontext.KindUser {
			ctx = _identitycontext.WithActor(ctx, _identitycontext.Actor{ProfileID: 42})
		}

		subject, err := Resolve(ctx, Requirement{Kind: _quotasubject.KindAPIKey})
		if err != nil || subject.Key() != "api_key:88" {
			t.Fatalf("kind %s subject = %+v, error = %v", kind, subject, err)
		}
	}
}

func TestOrganizationRequiresMatchingAuthoritativeGrant(t *testing.T) {
	t.Parallel()

	organizationID := uint64(12)
	ctx := verifiedContext("gateway-service")
	ctx = _callercontext.WithCaller(ctx, _callercontext.Caller{Kind: _callercontext.KindUser})
	ctx = _identitycontext.WithActor(ctx, _identitycontext.Actor{ProfileID: 42, OrganizationID: &organizationID})

	_, err := Resolve(ctx, Requirement{Kind: _quotasubject.KindOrganization})
	if !errors.Is(err, ErrEvidenceRequired) {
		t.Fatalf("missing grant error = %v, want ErrEvidenceRequired", err)
	}

	subject, err := Resolve(ctx, Requirement{
		Kind: _quotasubject.KindOrganization,
		Grant: &Grant{
			SubjectKind: _quotasubject.KindOrganization,
			SubjectID:   "12",
			GranteeKind: _quotasubject.KindUser,
			GranteeID:   "42",
			Authority:   "organization-service",
		},
	})
	if err != nil || subject.Key() != "organization:12" {
		t.Fatalf("subject = %+v, error = %v", subject, err)
	}
}

func TestOrganizationAssertionAloneNeverResolves(t *testing.T) {
	t.Parallel()

	organizationID := uint64(12)
	ctx := verifiedContext("gateway-service")
	ctx = _callercontext.WithCaller(ctx, _callercontext.Caller{Kind: _callercontext.KindUser})
	ctx = _identitycontext.WithActor(ctx, _identitycontext.Actor{ProfileID: 42, OrganizationID: &organizationID})

	_, err := Resolve(ctx, Requirement{
		Kind: _quotasubject.KindOrganization,
		Grant: &Grant{
			SubjectKind: _quotasubject.KindOrganization,
			SubjectID:   "12",
			GranteeKind: _quotasubject.KindUser,
			GranteeID:   "999",
			Authority:   "organization-service",
		},
	})
	if !errors.Is(err, ErrEvidenceMismatch) {
		t.Fatalf("error = %v, want ErrEvidenceMismatch", err)
	}
}

func TestAnonymousRequiresEdgeGrantFromTransportCaller(t *testing.T) {
	t.Parallel()

	ctx := verifiedContext("gateway-service")
	ctx = _callercontext.WithCaller(ctx, _callercontext.Caller{Kind: _callercontext.KindAnonymous})

	_, err := Resolve(ctx, Requirement{Kind: _quotasubject.KindAnonymous})
	if !errors.Is(err, ErrEvidenceRequired) {
		t.Fatalf("error = %v, want ErrEvidenceRequired", err)
	}

	subject, err := Resolve(ctx, Requirement{
		Kind: _quotasubject.KindAnonymous,
		Grant: &Grant{
			SubjectKind: _quotasubject.KindAnonymous,
			SubjectID:   "public_bucket_01",
			Authority:   "gateway-service",
		},
	})
	if err != nil || subject.Key() != "anonymous:public_bucket_01" {
		t.Fatalf("subject = %+v, error = %v", subject, err)
	}
}

func TestUnverifiedTransportNeverResolves(t *testing.T) {
	t.Parallel()

	ctx := _callercontext.WithCaller(context.Background(), _callercontext.Caller{Kind: _callercontext.KindUser})
	ctx = _identitycontext.WithActor(ctx, _identitycontext.Actor{ProfileID: 42})

	_, err := Resolve(ctx, Requirement{Kind: _quotasubject.KindUser})
	if !errors.Is(err, ErrUnverifiedTransport) {
		t.Fatalf("error = %v, want ErrUnverifiedTransport", err)
	}
}

func TestExistingSubjectIsImmutableAcrossDownstreamResolution(t *testing.T) {
	t.Parallel()

	ctx := verifiedContext("tqd-service")
	ctx = _quotasubject.WithSubject(ctx, _quotasubject.Subject{
		Kind: _quotasubject.KindOrganization, ID: "12", Source: _quotasubject.SourceOrganizationMembership,
	})

	subject, err := Resolve(ctx, Requirement{Kind: _quotasubject.KindOrganization})
	if err != nil || subject.Key() != "organization:12" {
		t.Fatalf("subject = %+v, error = %v", subject, err)
	}

	_, err = Resolve(ctx, Requirement{Kind: _quotasubject.KindUser})
	if !errors.Is(err, ErrSubjectConflict) {
		t.Fatalf("error = %v, want ErrSubjectConflict", err)
	}
}

func TestResolveUserDoesNotFallbackToAuthOrOriginID(t *testing.T) {
	t.Parallel()

	ctx := verifiedContext("gateway-service")
	ctx = _callercontext.WithCaller(ctx, _callercontext.Caller{Kind: _callercontext.KindUser})
	ctx = _identitycontext.WithActor(ctx, _identitycontext.Actor{AuthID: 7, OriginID: 99})

	_, err := Resolve(ctx, Requirement{Kind: _quotasubject.KindUser})
	if !errors.Is(err, ErrCallerMismatch) {
		t.Fatalf("error = %v, want ErrCallerMismatch", err)
	}
}

func TestResolveProjectAndSubscriptionFromAuthorityGrant(t *testing.T) {
	t.Parallel()

	ctx := verifiedContext("gateway-service")
	ctx = _callercontext.WithCaller(ctx, _callercontext.Caller{Kind: _callercontext.KindUser})
	ctx = _identitycontext.WithActor(ctx, _identitycontext.Actor{ProfileID: 42})

	tests := []struct {
		kind _quotasubject.Kind
		id   string
		key  string
	}{
		{kind: _quotasubject.KindProject, id: "project_9", key: "project:project_9"},
		{kind: _quotasubject.KindSubscription, id: "subscription_7", key: "subscription:subscription_7"},
	}
	for _, test := range tests {
		test := test
		t.Run(string(test.kind), func(t *testing.T) {
			t.Parallel()
			subject, err := Resolve(ctx, Requirement{
				Kind: test.kind,
				Grant: &Grant{
					SubjectKind: test.kind,
					SubjectID:   test.id,
					GranteeKind: _quotasubject.KindUser,
					GranteeID:   "42",
					Authority:   "subject-authority-service",
				},
			})
			if err != nil || subject.Key() != test.key {
				t.Fatalf("subject = %+v, error = %v, want %s", subject, err, test.key)
			}
		})
	}
}

func TestResolveInternalServiceUsesVerifiedTransportServiceID(t *testing.T) {
	t.Parallel()

	ctx := verifiedContext("report-worker")
	ctx = _callercontext.WithCaller(ctx, _callercontext.Caller{Kind: _callercontext.KindInternalService})

	subject, err := Resolve(ctx, Requirement{Kind: _quotasubject.KindInternalService})
	if err != nil || subject.Key() != "internal_service:report-worker" {
		t.Fatalf("subject = %+v, error = %v", subject, err)
	}
}
