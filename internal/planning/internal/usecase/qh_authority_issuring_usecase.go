package usecase

import (
	"context"
	"fmt"
	"strings"

	_dto "common/domain/dto"

	qh_domain "tqd/internal/domain/qh"
	"tqd/internal/interface/repo"
)

// QHAuthorityIssuringUpdateInput — trường nil = không đổi.
type QHAuthorityIssuringUpdateInput struct {
	Name        *string
	Code        *string
	Description *string
}

type QHAuthorityIssuringUsecase interface {
	Create(ctx context.Context, row *qh_domain.QHAuthorityIssuring) (*qh_domain.QHAuthorityIssuring, error)
	Update(ctx context.Context, id uint64, in *QHAuthorityIssuringUpdateInput) (*qh_domain.QHAuthorityIssuring, error)
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*qh_domain.QHAuthorityIssuring, error)
	List(ctx context.Context, pagable *_dto.Pagable) ([]qh_domain.QHAuthorityIssuring, int64, error)
}

type qhAuthorityIssuringUsecase struct {
	repo repo.QHAuthorityIssuringRepository
}

func NewQHAuthorityIssuringUsecase(r repo.QHAuthorityIssuringRepository) QHAuthorityIssuringUsecase {
	return &qhAuthorityIssuringUsecase{repo: r}
}

func (u *qhAuthorityIssuringUsecase) Create(ctx context.Context, row *qh_domain.QHAuthorityIssuring) (*qh_domain.QHAuthorityIssuring, error) {
	if row == nil {
		return nil, ErrQHAuthorityIssuringPayloadRequired
	}
	row.Name = strings.TrimSpace(row.Name)
	row.Code = strings.TrimSpace(row.Code)
	if row.Name == "" {
		return nil, ErrQHAuthorityIssuringNameRequired
	}
	if row.Code == "" {
		return nil, ErrQHAuthorityIssuringCodeRequired
	}
	dup, err := u.repo.GetByCode(ctx, row.Code)
	if err != nil {
		return nil, fmt.Errorf("check code: %w", err)
	}
	if dup != nil {
		return nil, qhAuthorityIssuringCodeConflict(row.Code)
	}
	if err := u.repo.Create(ctx, row); err != nil {
		return nil, fmt.Errorf("create authority issuring: %w", err)
	}
	return row, nil
}

func (u *qhAuthorityIssuringUsecase) Update(ctx context.Context, id uint64, in *QHAuthorityIssuringUpdateInput) (*qh_domain.QHAuthorityIssuring, error) {
	if id == 0 {
		return nil, ErrQHAuthorityIssuringIDRequired
	}
	if in == nil {
		return nil, ErrQHAuthorityIssuringPayloadRequired
	}
	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get by id: %w", err)
	}
	if existing == nil {
		return nil, qhAuthorityIssuringNotFound(id)
	}

	if in.Name != nil {
		s := strings.TrimSpace(*in.Name)
		if s == "" {
			return nil, ErrQHAuthorityIssuringNameRequired
		}
		existing.Name = s
	}
	if in.Code != nil {
		s := strings.TrimSpace(*in.Code)
		if s == "" {
			return nil, ErrQHAuthorityIssuringCodeRequired
		}
		if s != existing.Code {
			dup, err := u.repo.GetByCode(ctx, s)
			if err != nil {
				return nil, fmt.Errorf("check code: %w", err)
			}
			if dup != nil && dup.ID != id {
				return nil, qhAuthorityIssuringCodeConflict(s)
			}
			existing.Code = s
		}
	}
	if in.Description != nil {
		existing.Description = *in.Description
	}

	if err := u.repo.Update(ctx, existing); err != nil {
		return nil, fmt.Errorf("update authority issuring: %w", err)
	}
	return existing, nil
}

func (u *qhAuthorityIssuringUsecase) Delete(ctx context.Context, id uint64) error {
	if id == 0 {
		return ErrQHAuthorityIssuringIDRequired
	}
	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get by id: %w", err)
	}
	if existing == nil {
		return qhAuthorityIssuringNotFound(id)
	}
	return u.repo.Delete(ctx, id)
}

func (u *qhAuthorityIssuringUsecase) GetByID(ctx context.Context, id uint64) (*qh_domain.QHAuthorityIssuring, error) {
	if id == 0 {
		return nil, ErrQHAuthorityIssuringIDRequired
	}
	row, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get by id: %w", err)
	}
	if row == nil {
		return nil, qhAuthorityIssuringNotFound(id)
	}
	return row, nil
}

func (u *qhAuthorityIssuringUsecase) List(ctx context.Context, pagable *_dto.Pagable) ([]qh_domain.QHAuthorityIssuring, int64, error) {
	offset := 0
	limit := 20
	if pagable != nil {
		offset = pagable.GetOffset()
		limit = pagable.GetLimit()
	}
	return u.repo.List(ctx, offset, limit)
}
