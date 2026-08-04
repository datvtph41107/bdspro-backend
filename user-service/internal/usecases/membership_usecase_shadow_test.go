package usecases

import (
	"context"
	"errors"
	"testing"

	_enums "common/domain/enum"
	"user/internal/interface/repo"
	"user/internal/models"
)

type membershipShadowRepo struct {
	repo.IMembershipRepo
	membership *models.Membership
	err        error
}

func (r *membershipShadowRepo) GetByProfileID(
	ctx context.Context,
	profileID uint64,
) (*models.Membership, error) {
	if profileID != 42 {
		return nil, errors.New("unexpected profile")
	}
	return r.membership, r.err
}

type membershipObserver struct {
	calls      int
	membership *models.Membership
}

func (o *membershipObserver) ObserveMembership(
	ctx context.Context,
	membership *models.Membership,
) {
	o.calls++
	o.membership = membership
}

func TestGetMyMembershipObservesAfterLegacySuccess(t *testing.T) {
	membership := &models.Membership{ProfileID: 42, PlanID: 7}
	membership.ID = 9
	observer := &membershipObserver{}
	usecase := NewMembershipUsecase(
		&membershipShadowRepo{membership: membership},
		nil,
		nil,
		observer,
	)
	ctx := context.WithValue(context.Background(), _enums.ProfileIDKey, uint64(42))

	got, errDTO := usecase.GetMyMembership(ctx)
	if errDTO != nil {
		t.Fatalf("GetMyMembership() error = %+v", errDTO)
	}
	if got != membership {
		t.Fatal("legacy membership response was replaced")
	}
	if observer.calls != 1 || observer.membership != membership {
		t.Fatalf("observer calls=%d membership=%p", observer.calls, observer.membership)
	}
}

func TestGetMyMembershipDoesNotObserveLegacyFailure(t *testing.T) {
	observer := &membershipObserver{}
	usecase := NewMembershipUsecase(
		&membershipShadowRepo{err: errors.New("not found")},
		nil,
		nil,
		observer,
	)
	ctx := context.WithValue(context.Background(), _enums.ProfileIDKey, uint64(42))

	if _, errDTO := usecase.GetMyMembership(ctx); errDTO == nil {
		t.Fatal("GetMyMembership() accepted legacy read failure")
	}
	if observer.calls != 0 {
		t.Fatalf("observer calls = %d", observer.calls)
	}
}
