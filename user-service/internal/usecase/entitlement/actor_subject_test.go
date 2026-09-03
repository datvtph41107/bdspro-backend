package evaluate

import (
	"common/identity"
	"errors"
	"testing"

	"user/internal/domain/entitlement"
)

func TestSubjectFromActorDefaultsToProfilePool(t *testing.T) {
	subject, err := SubjectFromActor(identity.Actor{ProfileID: 42})
	if err != nil {
		t.Fatalf("SubjectFromActor() error = %v", err)
	}
	if subject != (access.Subject{Type: access.SubjectProfile, ID: "42"}) {
		t.Fatalf("subject = %+v, want profile 42", subject)
	}
}

func TestSubjectFromActorUsesVerifiedOrganizationContext(t *testing.T) {
	organizationID := uint64(7)
	subject, err := SubjectFromActor(identity.Actor{
		ProfileID:      42,
		OrganizationID: &organizationID,
	})
	if err != nil {
		t.Fatalf("SubjectFromActor() error = %v", err)
	}
	if subject != (access.Subject{Type: access.SubjectOrganization, ID: "7"}) {
		t.Fatalf("subject = %+v, want organization 7", subject)
	}
}

func TestSubjectFromActorRejectsMissingProfile(t *testing.T) {
	_, err := SubjectFromActor(identity.Actor{Role: "MEMBER"})
	if !errors.Is(err, ErrActorProfileMissing) {
		t.Fatalf("SubjectFromActor() error = %v, want ErrActorProfileMissing", err)
	}
}
