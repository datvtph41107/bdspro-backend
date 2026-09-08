package organization

import (
	"context"
	"errors"
	"testing"

	domain "user/internal/domain/organization"
)

type repositoryStub struct {
	created    domain.Organization
	updated    domain.Update
	current    domain.Organization
	forProfile []domain.Organization
}

func (r *repositoryStub) List(context.Context, domain.ListQuery) ([]domain.Organization, uint64, error) {
	return nil, 0, nil
}

func (r *repositoryStub) Get(context.Context, uint64) (domain.Organization, error) {
	if r.current.ID == 0 {
		return domain.Organization{}, ErrNotFound
	}
	return r.current, nil
}

func (r *repositoryStub) Create(_ context.Context, value domain.Organization) (domain.Organization, error) {
	r.created = value
	value.ID = 1
	return value, nil
}

func (r *repositoryStub) Update(_ context.Context, value domain.Update) (domain.Organization, error) {
	r.updated = value
	return r.current, nil
}

func (r *repositoryStub) Archive(context.Context, uint64, uint64) error { return nil }
func (r *repositoryStub) CheckMembers(_ context.Context, _ uint64, profileIDs []uint64) ([]domain.MemberCheck, error) {
	result := make([]domain.MemberCheck, 0, len(profileIDs))
	for _, profileID := range profileIDs {
		result = append(result, domain.MemberCheck{ProfileID: profileID})
	}
	return result, nil
}
func (r *repositoryStub) ListForProfile(context.Context, uint64) ([]domain.Organization, error) {
	return r.forProfile, nil
}

func TestCreateNormalizesOrganizationAndPreservesActor(t *testing.T) {
	repository := &repositoryStub{}
	service := NewService(repository)

	_, err := service.Create(context.Background(), domain.Organization{
		Code: "  QHPRO_TEAM  ", Name: "  Nhóm QHPro  ", Type: "Nhóm",
		Email: " team@example.com ", OwnerProfileID: 42, CreatedBy: 42,
	})
	if err != nil {
		t.Fatalf("Create returned error: %v", err)
	}
	if repository.created.Code != "QHPRO_TEAM" || repository.created.Name != "Nhóm QHPro" {
		t.Fatalf("organization text was not normalized: %+v", repository.created)
	}
	if repository.created.Type != "team" || repository.created.Status != "active" {
		t.Fatalf("organization defaults are invalid: %+v", repository.created)
	}
	if repository.created.OwnerProfileID != 42 || repository.created.CreatedBy != 42 {
		t.Fatalf("trusted actor identity was not preserved: %+v", repository.created)
	}
}

func TestCreateRejectsInvalidAuthorityAndEmail(t *testing.T) {
	service := NewService(&repositoryStub{})
	_, err := service.Create(context.Background(), domain.Organization{
		Code: "QHPRO", Name: "QHPro", Type: "business", Email: "not-an-email",
	})
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected invalid input for missing actor, got %v", err)
	}

	_, err = service.Create(context.Background(), domain.Organization{
		Code: "QHPRO", Name: "QHPro", Type: "business", Email: "not-an-email",
		OwnerProfileID: 1, CreatedBy: 1,
	})
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected invalid input for email, got %v", err)
	}
}

func TestUpdateValidatesMergedStateBeforePersistence(t *testing.T) {
	repository := &repositoryStub{current: domain.Organization{
		ID: 7, Code: "ORG_7", Name: "Tổ chức 7", Type: "business",
		Status: "active", VerificationStatus: "unverified", WarningLevel: "none",
	}}
	service := NewService(repository)
	invalid := "deleted"
	_, err := service.Update(context.Background(), domain.Update{ID: 7, ActorID: 9, Status: &invalid})
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected invalid merged status, got %v", err)
	}
}

func TestCheckMembersDeduplicatesAndRejectsUnboundedInput(t *testing.T) {
	service := NewService(&repositoryStub{})
	checks, err := service.CheckMembers(context.Background(), 8, []uint64{3, 0, 3, 5})
	if err != nil {
		t.Fatalf("CheckMembers returned error: %v", err)
	}
	if len(checks) != 2 || checks[0].ProfileID != 3 || checks[1].ProfileID != 5 {
		t.Fatalf("unexpected normalized checks: %+v", checks)
	}
	_, err = service.CheckMembers(context.Background(), 8, make([]uint64, 501))
	if !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected bounded input error, got %v", err)
	}
}

func TestListForProfileRequiresIdentityAndUsesCanonicalRepository(t *testing.T) {
	repository := &repositoryStub{forProfile: []domain.Organization{{ID: 9, Name: "QHPro Team"}}}
	service := NewService(repository)
	if _, err := service.ListForProfile(context.Background(), 0); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected invalid profile identity, got %v", err)
	}
	organizations, err := service.ListForProfile(context.Background(), 42)
	if err != nil || len(organizations) != 1 || organizations[0].ID != 9 {
		t.Fatalf("unexpected account organization projection: organizations=%+v err=%v", organizations, err)
	}
}
