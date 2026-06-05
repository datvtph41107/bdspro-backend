package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	_dto "common/domain/dto"

	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/interface/repo"
)

// QHLayerFamilyUpdateInput — trường nil = không đổi.
type QHLayerFamilyUpdateInput struct {
	Name       *string
	SortNumber *int32
}

type QHLayerFamilyUsecase interface {
	Create(ctx context.Context, row *qh_domain.QHLayerFamily) (*qh_domain.QHLayerFamily, error)
	Update(ctx context.Context, id uint64, in *QHLayerFamilyUpdateInput) (*qh_domain.QHLayerFamily, error)
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*qh_domain.QHLayerFamily, error)
	List(ctx context.Context, pagable *_dto.Pagable, search string) ([]qh_domain.QHLayerFamily, int64, error)
}

type qhLayerFamilyUsecase struct {
	repo repo.QHLayerFamilyRepository
}

func NewQHLayerFamilyUsecase(r repo.QHLayerFamilyRepository) QHLayerFamilyUsecase {
	return &qhLayerFamilyUsecase{repo: r}
}

func (u *qhLayerFamilyUsecase) Create(ctx context.Context, row *qh_domain.QHLayerFamily) (*qh_domain.QHLayerFamily, error) {
	if row == nil {
		return nil, errors.New("payload is required")
	}
	row.Name = strings.TrimSpace(row.Name)
	if row.Name == "" {
		return nil, errors.New("name is required")
	}
	if err := u.repo.Create(ctx, row); err != nil {
		return nil, fmt.Errorf("create layer family: %w", err)
	}
	return row, nil
}

func (u *qhLayerFamilyUsecase) Update(ctx context.Context, id uint64, in *QHLayerFamilyUpdateInput) (*qh_domain.QHLayerFamily, error) {
	if id == 0 {
		return nil, errors.New("id is required")
	}
	if in == nil {
		return nil, errors.New("payload is required")
	}
	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get by id: %w", err)
	}
	if existing == nil {
		return nil, fmt.Errorf("record %d not found", id)
	}

	if in.Name != nil {
		s := strings.TrimSpace(*in.Name)
		if s == "" {
			return nil, errors.New("name cannot be empty")
		}
		existing.Name = s
	}

	if in.SortNumber != nil {
		existing.SortNumber = *in.SortNumber
	}

	if err := u.repo.Update(ctx, existing); err != nil {
		return nil, fmt.Errorf("update layer family: %w", err)
	}
	return existing, nil
}

func (u *qhLayerFamilyUsecase) Delete(ctx context.Context, id uint64) error {
	if id == 0 {
		return errors.New("id is required")
	}
	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get by id: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("record %d not found", id)
	}
	return u.repo.Delete(ctx, id)
}

func (u *qhLayerFamilyUsecase) GetByID(ctx context.Context, id uint64) (*qh_domain.QHLayerFamily, error) {
	if id == 0 {
		return nil, errors.New("id is required")
	}
	row, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get by id: %w", err)
	}
	if row == nil {
		return nil, fmt.Errorf("record %d not found", id)
	}
	return row, nil
}

func (u *qhLayerFamilyUsecase) List(ctx context.Context, pagable *_dto.Pagable, search string) ([]qh_domain.QHLayerFamily, int64, error) {
	offset := 0
	limit := 20
	sort := ""
	if pagable != nil {
		offset = pagable.GetOffset()
		limit = pagable.GetLimit()
		sort = pagable.GetSort()
	}
	return u.repo.List(ctx, offset, limit, search, sort)
}
