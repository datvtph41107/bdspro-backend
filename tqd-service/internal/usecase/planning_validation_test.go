package usecase

import (
	"common/validation"
	"context"
	"errors"
	"testing"
	qh "tqd/internal/domain/qh"
	"tqd/internal/interface/repo"
)

type familyErrorRepo struct {
	repo.QHLayerFamilyRepository
	err error
}

func (r familyErrorRepo) GetByID(context.Context, uint64) (*qh.QHLayerFamily, error) {
	return nil, r.err
}

// Repository nil và lỗi truy vấn là hai kết quả khác nhau, kể cả khi message giống nhau.
func TestFamilyLookupPreservesDatabaseFailure(t *testing.T) {
	failure := errors.New("required database table not found")
	service := NewQHLayerFamilyUsecase(familyErrorRepo{err: failure}, nil, nil)
	_, err := service.GetByID(context.Background(), 1)
	if !errors.Is(err, failure) || errors.Is(err, ErrLayerFamilyNotFound) {
		t.Fatalf("lost database cause: %v", err)
	}
	err = service.BuildFamilyPMTiles(context.Background(), 1, 6, 16, "", func(BuildProgress) {})
	if !errors.Is(err, failure) || errors.Is(err, ErrLayerFamilyNotFound) {
		t.Fatalf("build lost database cause: %v", err)
	}
	service = NewQHLayerFamilyUsecase(familyErrorRepo{}, nil, nil)
	_, err = service.GetByID(context.Background(), 1)
	if !errors.Is(err, ErrLayerFamilyNotFound) {
		t.Fatalf("missing row: %v", err)
	}
}

func TestPlanningValidationRunsBeforeRepository(t *testing.T) {
	cases := []struct {
		name  string
		run   func() error
		field string
	}{
		{"family", func() error {
			_, e := NewQHLayerFamilyUsecase(nil, nil, nil).Create(context.Background(), &qh.QHLayerFamily{Name: "  "})
			return e
		}, "name"},
		{"authority", func() error {
			_, e := NewQHAuthorityIssuringUsecase(nil).Create(context.Background(), &qh.QHAuthorityIssuring{Name: "Issuer", Code: " "})
			return e
		}, "code"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			v, ok := validation.As(tc.run())
			if !ok || v.Field != tc.field {
				t.Fatalf("expected field %s, got %v", tc.field, v)
			}
		})
	}
}
