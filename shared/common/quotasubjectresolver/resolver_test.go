package quotasubjectresolver

import (
	_callercontext "common/callercontext"
	_identitycontext "common/identitycontext"
	_quotasubject "common/quotasubject"
	_transporttrust "common/transporttrust"

	"context"
	"errors"
	"strconv"
	"sync"
	"testing"
	"time"
)

func mustBindActor(ctx context.Context, actor _identitycontext.Actor) context.Context {
	bound, err := _identitycontext.BindActor(ctx, actor)
	if err != nil {
		panic(err)
	}
	return bound
}

func mustBindSubject(ctx context.Context, subject _quotasubject.Subject) context.Context {
	bound, err := _quotasubject.BindSubject(ctx, subject)
	if err != nil {
		panic(err)
	}
	return bound
}

func mustBindCaller(ctx context.Context, caller _callercontext.Caller) context.Context {
	bound, err := _callercontext.BindCaller(ctx, caller)
	if err != nil {
		panic(err)
	}
	return bound
}

func mustBindTransportCaller(ctx context.Context, caller _transporttrust.Caller) context.Context {
	bound, err := _transporttrust.BindCaller(ctx, caller)
	if err != nil {
		panic(err)
	}
	return bound
}

func verifiedContext(serviceID string) context.Context {
	return mustBindTransportCaller(context.Background(), _transporttrust.Caller{
		ServiceID: serviceID,
		Level:     _transporttrust.LevelVerified,
		IssuedAt:  time.Now(),
	})
}

func verifiedUserContext(profileID uint64) context.Context {
	ctx := verifiedContext("gateway-service")
	ctx = mustBindCaller(ctx, _callercontext.Caller{Kind: _callercontext.KindUser})
	return mustBindActor(ctx, _identitycontext.Actor{ProfileID: profileID})
}

func verifiedUserContextWithOrganization(profileID, organizationID uint64) context.Context {
	ctx := verifiedContext("gateway-service")
	ctx = mustBindCaller(ctx, _callercontext.Caller{Kind: _callercontext.KindUser})
	return mustBindActor(ctx, _identitycontext.Actor{
		ProfileID:      profileID,
		OrganizationID: &organizationID,
	})
}

func authorityGrant(
	kind _quotasubject.Kind,
	subjectID string,
	granteeKind _quotasubject.Kind,
	granteeID string,
	authority string,
) *Grant {
	return &Grant{
		SubjectKind: kind,
		SubjectID:   subjectID,
		GranteeKind: granteeKind,
		GranteeID:   granteeID,
		Authority:   authority,
		EvidenceID:  "decision-" + string(kind) + "-" + subjectID + "-" + granteeID,
	}
}

func TestResolveUserUsesProfileIDWithoutFallback(t *testing.T) {
	t.Parallel()

	ctx := verifiedUserContext(42)
	subject, err := Resolve(ctx, Requirement{Kind: _quotasubject.KindUser, ExpectedID: "42"})
	if err != nil || subject.Key() != "user:42" {
		t.Fatalf("subject = %+v, error = %v", subject, err)
	}

	_, err = Resolve(ctx, Requirement{Kind: _quotasubject.KindUser, ExpectedID: "99"})
	if !errors.Is(err, ErrEvidenceMismatch) {
		t.Fatalf("mismatched expected ID error = %v, want ErrEvidenceMismatch", err)
	}
}

func TestResolveAPIKeyForAPIKeyAndUserWithVerifiedClient(t *testing.T) {
	t.Parallel()

	for _, kind := range []_callercontext.Kind{_callercontext.KindAPIKey, _callercontext.KindUser} {
		kind := kind
		t.Run(string(kind), func(t *testing.T) {
			t.Parallel()
			ctx := verifiedContext("gateway-service")
			ctx = mustBindCaller(ctx, _callercontext.Caller{
				Kind: kind, APIKeyID: 88, APIKeyApp: "partner-app",
			})
			if kind == _callercontext.KindUser {
				ctx = mustBindActor(ctx, _identitycontext.Actor{ProfileID: 42})
			}

			subject, err := Resolve(ctx, Requirement{Kind: _quotasubject.KindAPIKey, ExpectedID: "88"})
			if err != nil || subject.Key() != "api_key:88" {
				t.Fatalf("kind %s subject = %+v, error = %v", kind, subject, err)
			}
		})
	}
}

func TestOrganizationRequiresMatchingAuthoritativeGrant(t *testing.T) {
	t.Parallel()

	ctx := verifiedUserContextWithOrganization(42, 12)

	_, err := Resolve(ctx, Requirement{Kind: _quotasubject.KindOrganization, ExpectedID: "12", Authority: "organization-service"})
	if !errors.Is(err, ErrEvidenceRequired) {
		t.Fatalf("missing grant error = %v, want ErrEvidenceRequired", err)
	}

	subject, err := Resolve(ctx, Requirement{
		Kind:       _quotasubject.KindOrganization,
		ExpectedID: "12",
		Authority:  "organization-service",
		Grant: authorityGrant(
			_quotasubject.KindOrganization,
			"12",
			_quotasubject.KindUser,
			"42",
			"organization-service",
		),
	})
	if err != nil || subject.Key() != "organization:12" {
		t.Fatalf("subject = %+v, error = %v", subject, err)
	}
}

func TestAuthorityAndEvidenceIDArePartOfTheRequirement(t *testing.T) {
	t.Parallel()

	ctx := verifiedUserContext(42)
	grant := authorityGrant(
		_quotasubject.KindOrganization,
		"12",
		_quotasubject.KindUser,
		"42",
		"organization-service",
	)

	_, err := Resolve(ctx, Requirement{
		Kind:       _quotasubject.KindOrganization,
		ExpectedID: "12",
		Authority:  "untrusted-service",
		Grant:      grant,
	})
	if !errors.Is(err, ErrEvidenceMismatch) {
		t.Fatalf("wrong authority error = %v, want ErrEvidenceMismatch", err)
	}

	grantWithoutEvidence := *grant
	grantWithoutEvidence.EvidenceID = ""
	_, err = Resolve(ctx, Requirement{
		Kind:       _quotasubject.KindOrganization,
		ExpectedID: "12",
		Authority:  "organization-service",
		Grant:      &grantWithoutEvidence,
	})
	if !errors.Is(err, ErrEvidenceMismatch) {
		t.Fatalf("missing evidence ID error = %v, want ErrEvidenceMismatch", err)
	}
}

func TestGrantFieldsAreNormalizedBeforeUse(t *testing.T) {
	t.Parallel()

	ctx := verifiedUserContext(42)
	subject, err := Resolve(ctx, Requirement{
		Kind:       _quotasubject.KindOrganization,
		ExpectedID: " 12 ",
		Authority:  " organization-service ",
		Grant: &Grant{
			SubjectKind: _quotasubject.KindOrganization,
			SubjectID:   " 12 ",
			GranteeKind: _quotasubject.KindUser,
			GranteeID:   " 42 ",
			Authority:   " organization-service ",
			EvidenceID:  " decision-org-12-user-42 ",
		},
	})
	if err != nil || subject.Key() != "organization:12" {
		t.Fatalf("subject = %+v, error = %v", subject, err)
	}
}

func TestOrganizationAssertionAloneNeverResolves(t *testing.T) {
	t.Parallel()

	ctx := verifiedUserContextWithOrganization(42, 12)

	_, err := Resolve(ctx, Requirement{
		Kind:       _quotasubject.KindOrganization,
		ExpectedID: "12",
		Authority:  "organization-service",
		Grant: authorityGrant(
			_quotasubject.KindOrganization,
			"12",
			_quotasubject.KindUser,
			"999",
			"organization-service",
		),
	})
	if !errors.Is(err, ErrEvidenceMismatch) {
		t.Fatalf("error = %v, want ErrEvidenceMismatch", err)
	}
}

func TestAnonymousRequiresEdgeGrantFromTransportCaller(t *testing.T) {
	t.Parallel()

	ctx := verifiedContext("gateway-service")
	ctx = mustBindCaller(ctx, _callercontext.Caller{Kind: _callercontext.KindAnonymous})

	_, err := Resolve(ctx, Requirement{
		Kind:       _quotasubject.KindAnonymous,
		ExpectedID: "public_bucket_01",
		Authority:  "gateway-service",
	})
	if !errors.Is(err, ErrEvidenceRequired) {
		t.Fatalf("error = %v, want ErrEvidenceRequired", err)
	}

	subject, err := Resolve(ctx, Requirement{
		Kind:       _quotasubject.KindAnonymous,
		ExpectedID: "public_bucket_01",
		Authority:  "gateway-service",
		Grant: &Grant{
			SubjectKind: _quotasubject.KindAnonymous,
			SubjectID:   "public_bucket_01",
			Authority:   "gateway-service",
			EvidenceID:  "edge-bucket-policy-v1",
		},
	})
	if err != nil || subject.Key() != "anonymous:public_bucket_01" {
		t.Fatalf("subject = %+v, error = %v", subject, err)
	}
}

func TestUnverifiedTransportNeverResolves(t *testing.T) {
	t.Parallel()

	ctx := mustBindCaller(context.Background(), _callercontext.Caller{Kind: _callercontext.KindUser})
	ctx = mustBindActor(ctx, _identitycontext.Actor{ProfileID: 42})

	_, err := Resolve(ctx, Requirement{Kind: _quotasubject.KindUser})
	if !errors.Is(err, ErrUnverifiedTransport) {
		t.Fatalf("error = %v, want ErrUnverifiedTransport", err)
	}
}

func TestExistingResourceSubjectRequiresExactExpectedID(t *testing.T) {
	t.Parallel()

	ctx := verifiedContext("tqd-service")
	ctx = mustBindCaller(ctx, _callercontext.Caller{Kind: _callercontext.KindUser})
	ctx = mustBindActor(ctx, _identitycontext.Actor{ProfileID: 42})
	ctx = mustBindSubject(ctx, _quotasubject.Subject{
		Kind: _quotasubject.KindOrganization, ID: "12", Source: _quotasubject.SourceOrganizationMembership,
	})

	_, err := Resolve(ctx, Requirement{Kind: _quotasubject.KindOrganization})
	if !errors.Is(err, ErrExpectedSubjectIDRequired) {
		t.Fatalf("kind-only continuation error = %v, want ErrExpectedSubjectIDRequired", err)
	}

	subject, err := Resolve(ctx, Requirement{Kind: _quotasubject.KindOrganization, ExpectedID: "12"})
	if err != nil || subject.Key() != "organization:12" {
		t.Fatalf("exact continuation subject = %+v, error = %v", subject, err)
	}

	_, err = Resolve(ctx, Requirement{Kind: _quotasubject.KindOrganization, ExpectedID: "99"})
	if !errors.Is(err, ErrSubjectConflict) {
		t.Fatalf("different ID error = %v, want ErrSubjectConflict", err)
	}
}

func TestExistingSameKindConflictingGrantIsRejected(t *testing.T) {
	t.Parallel()

	ctx := verifiedUserContext(42)
	ctx = mustBindSubject(ctx, _quotasubject.Subject{
		Kind: _quotasubject.KindOrganization, ID: "12", Source: _quotasubject.SourceOrganizationMembership,
	})

	_, err := Resolve(ctx, Requirement{
		Kind:       _quotasubject.KindOrganization,
		ExpectedID: "12",
		Authority:  "organization-service",
		Grant: authorityGrant(
			_quotasubject.KindOrganization,
			"99",
			_quotasubject.KindUser,
			"42",
			"organization-service",
		),
	})
	if !errors.Is(err, ErrSubjectConflict) {
		t.Fatalf("conflicting same-kind grant error = %v, want ErrSubjectConflict", err)
	}
}

func TestExistingUserSubjectIsRevalidatedAgainstCurrentActor(t *testing.T) {
	t.Parallel()

	ctx := verifiedUserContext(99)
	ctx = mustBindSubject(ctx, _quotasubject.Subject{
		Kind: _quotasubject.KindUser, ID: "42", Source: _quotasubject.SourceActorProfile,
	})

	_, err := Resolve(ctx, Requirement{Kind: _quotasubject.KindUser})
	if !errors.Is(err, ErrSubjectConflict) {
		t.Fatalf("error = %v, want ErrSubjectConflict", err)
	}
}

func TestResolveUserDoesNotFallbackToAuthOrOriginID(t *testing.T) {
	t.Parallel()

	ctx := verifiedContext("gateway-service")
	ctx = mustBindCaller(ctx, _callercontext.Caller{Kind: _callercontext.KindUser})
	ctx = mustBindActor(ctx, _identitycontext.Actor{AuthID: 7, OriginID: 99})

	_, err := Resolve(ctx, Requirement{Kind: _quotasubject.KindUser})
	if !errors.Is(err, ErrCallerMismatch) {
		t.Fatalf("error = %v, want ErrCallerMismatch", err)
	}
}

func TestResolveProjectAndSubscriptionFromAuthorityGrant(t *testing.T) {
	t.Parallel()

	ctx := verifiedUserContext(42)
	tests := []struct {
		kind      _quotasubject.Kind
		id        string
		key       string
		authority string
	}{
		{kind: _quotasubject.KindProject, id: "project_9", key: "project:project_9", authority: "tqd-service"},
		{kind: _quotasubject.KindSubscription, id: "subscription_7", key: "subscription:subscription_7", authority: "user-service"},
	}
	for _, test := range tests {
		test := test
		t.Run(string(test.kind), func(t *testing.T) {
			t.Parallel()
			subject, err := Resolve(ctx, Requirement{
				Kind:       test.kind,
				ExpectedID: test.id,
				Authority:  test.authority,
				Grant: authorityGrant(
					test.kind,
					test.id,
					_quotasubject.KindUser,
					"42",
					test.authority,
				),
			})
			if err != nil || subject.Key() != test.key {
				t.Fatalf("subject = %+v, error = %v, want %s", subject, err, test.key)
			}
		})
	}
}

func TestPlanIDCannotStandInForSubscriptionEvidence(t *testing.T) {
	t.Parallel()

	ctx := verifiedUserContext(42)
	_, err := Resolve(ctx, Requirement{
		Kind:       _quotasubject.KindSubscription,
		ExpectedID: "plan_7",
		Authority:  "user-service",
	})
	if !errors.Is(err, ErrEvidenceRequired) {
		t.Fatalf("error = %v, want ErrEvidenceRequired", err)
	}
}

func TestResolveInternalServiceUsesVerifiedTransportServiceID(t *testing.T) {
	t.Parallel()

	ctx := verifiedContext("report-worker")
	ctx = mustBindCaller(ctx, _callercontext.Caller{Kind: _callercontext.KindInternalService})

	subject, err := Resolve(ctx, Requirement{Kind: _quotasubject.KindInternalService, ExpectedID: "report-worker"})
	if err != nil || subject.Key() != "internal_service:report-worker" {
		t.Fatalf("subject = %+v, error = %v", subject, err)
	}

	downstream := verifiedContext("tqd-service")
	downstream = mustBindCaller(downstream, _callercontext.Caller{Kind: _callercontext.KindInternalService})
	downstream = mustBindSubject(downstream, subject)
	subject, err = Resolve(downstream, Requirement{Kind: _quotasubject.KindInternalService, ExpectedID: "report-worker"})
	if err != nil || subject.Key() != "internal_service:report-worker" {
		t.Fatalf("downstream subject = %+v, error = %v", subject, err)
	}
}

func TestResolveAndAttachStoresExactlyOneSubject(t *testing.T) {
	t.Parallel()

	ctx := verifiedUserContext(42)
	attached, subject, err := ResolveAndAttach(ctx, Requirement{Kind: _quotasubject.KindUser})
	if err != nil || subject.Key() != "user:42" {
		t.Fatalf("subject = %+v, error = %v", subject, err)
	}
	stored, ok := _quotasubject.FromContext(attached)
	if !ok || stored != subject {
		t.Fatalf("stored subject = %+v, %v, want %+v", stored, ok, subject)
	}
}

func TestInvalidRequirementIsRejected(t *testing.T) {
	t.Parallel()

	ctx := verifiedUserContext(42)
	for _, requirement := range []Requirement{
		{Kind: _quotasubject.Kind("unknown")},
		{Kind: _quotasubject.KindUser, ExpectedID: "unsafe id"},
		{Kind: _quotasubject.KindUser, Authority: "organization-service"},
	} {
		_, err := Resolve(ctx, requirement)
		if err == nil {
			t.Fatalf("requirement %+v should fail", requirement)
		}
	}
}

func TestConcurrentResolutionIsDeterministic(t *testing.T) {
	t.Parallel()

	ctx := verifiedUserContext(42)
	const total = 100
	keys := make(chan string, total)
	errs := make(chan error, total)
	var wg sync.WaitGroup
	for index := 0; index < total; index++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			subject, err := Resolve(ctx, Requirement{Kind: _quotasubject.KindUser, ExpectedID: strconv.FormatUint(42, 10)})
			if err != nil {
				errs <- err
				return
			}
			keys <- subject.Key()
		}()
	}
	wg.Wait()
	close(keys)
	close(errs)

	for err := range errs {
		t.Fatalf("concurrent resolve error = %v", err)
	}
	count := 0
	for key := range keys {
		if key != "user:42" {
			t.Fatalf("concurrent key = %q, want user:42", key)
		}
		count++
	}
	if count != total {
		t.Fatalf("resolved %d subjects, want %d", count, total)
	}
}
