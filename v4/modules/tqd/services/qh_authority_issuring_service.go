package services

import (
	"common/validation"
	"context"
	"errors"
	"fmt"
	"strings"

	_dto "common/domain/dto"

	qh "mdn4/modules/tqd/entity/qh"
	repo "mdn4/modules/tqd/repository"
)

// ErrAuthorityNotFound chỉ biểu thị lookup thành công nhưng không có bản ghi.
var ErrAuthorityNotFound = errors.New("authority not found")

// ErrAuthorityCodeConflict dùng chung cho precheck và vi phạm unique constraint tại DB.
var ErrAuthorityCodeConflict = errors.New("authority code already exists")

// QHAuthorityIssuringUpdateInput — trường nil = không đổi.
type QHAuthorityIssuringUpdateInput struct {
	Name        *string
	Code        *string
	Description *string
}

type QHAuthorityIssuringUsecase interface {
	Create(ctx context.Context, row *qh.QHAuthorityIssuring) (*qh.QHAuthorityIssuring, error)
	Update(ctx context.Context, id uint64, in *QHAuthorityIssuringUpdateInput) (*qh.QHAuthorityIssuring, error)
	Delete(ctx context.Context, id uint64) error
	GetByID(ctx context.Context, id uint64) (*qh.QHAuthorityIssuring, error)
	List(ctx context.Context, pagable *_dto.Pagable) ([]qh.QHAuthorityIssuring, int64, error)
}

type qhAuthorityIssuringUsecase struct {
	repo repo.QHAuthorityIssuringRepository
}

func NewQHAuthorityIssuringUsecase(r repo.QHAuthorityIssuringRepository) QHAuthorityIssuringUsecase {
	return &qhAuthorityIssuringUsecase{repo: r}
}

func (u *qhAuthorityIssuringUsecase) Create(ctx context.Context, row *qh.QHAuthorityIssuring) (*qh.QHAuthorityIssuring, error) {
	if row == nil {
		return nil, validation.New("payload", "AUTHORITY_PAYLOAD_INVALID", "payload is required")
	}
	row.Name = strings.TrimSpace(row.Name)
	row.Code = strings.TrimSpace(row.Code)
	if row.Name == "" {
		return nil, validation.New("name", "AUTHORITY_NAME_INVALID", "name is required")
	}
	if row.Code == "" {
		return nil, validation.New("code", "AUTHORITY_CODE_INVALID", "code is required")
	}
	dup, err := u.repo.GetByCode(ctx, row.Code)
	if err != nil {
		return nil, fmt.Errorf("check code: %w", err)
	}
	if dup != nil {
		return nil, fmt.Errorf("%w: code=%q", ErrAuthorityCodeConflict, row.Code)
	}
	if err := u.repo.Create(ctx, row); err != nil {
		return nil, fmt.Errorf("create authority issuring: %w", err)
	}
	return row, nil
}

func (u *qhAuthorityIssuringUsecase) Update(ctx context.Context, id uint64, in *QHAuthorityIssuringUpdateInput) (*qh.QHAuthorityIssuring, error) {
	if id == 0 {
		return nil, validation.New("id", "AUTHORITY_ID_INVALID", "id is required")
	}
	if in == nil {
		return nil, validation.New("payload", "AUTHORITY_PAYLOAD_INVALID", "payload is required")
	}
	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get by id: %w", err)
	}
	if existing == nil {
		return nil, fmt.Errorf("%w: id=%d", ErrAuthorityNotFound, id)
	}

	if in.Name != nil {
		s := strings.TrimSpace(*in.Name)
		if s == "" {
			return nil, validation.New("name", "AUTHORITY_NAME_INVALID", "name cannot be empty")
		}
		existing.Name = s
	}
	if in.Code != nil {
		s := strings.TrimSpace(*in.Code)
		if s == "" {
			return nil, validation.New("code", "AUTHORITY_CODE_INVALID", "code cannot be empty")
		}
		if s != existing.Code {
			dup, err := u.repo.GetByCode(ctx, s)
			if err != nil {
				return nil, fmt.Errorf("check code: %w", err)
			}
			if dup != nil && dup.ID != id {
				return nil, fmt.Errorf("%w: code=%q", ErrAuthorityCodeConflict, s)
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
		return validation.New("id", "AUTHORITY_ID_INVALID", "id is required")
	}
	existing, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get by id: %w", err)
	}
	if existing == nil {
		return fmt.Errorf("%w: id=%d", ErrAuthorityNotFound, id)
	}
	return u.repo.Delete(ctx, id)
}

func (u *qhAuthorityIssuringUsecase) GetByID(ctx context.Context, id uint64) (*qh.QHAuthorityIssuring, error) {
	if id == 0 {
		return nil, validation.New("id", "AUTHORITY_ID_INVALID", "id is required")
	}
	row, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get by id: %w", err)
	}
	if row == nil {
		return nil, fmt.Errorf("%w: id=%d", ErrAuthorityNotFound, id)
	}
	return row, nil
}

func (u *qhAuthorityIssuringUsecase) List(ctx context.Context, pagable *_dto.Pagable) ([]qh.QHAuthorityIssuring, int64, error) {
	offset := 0
	limit := 20
	if pagable != nil {
		offset = pagable.GetOffset()
		limit = pagable.GetLimit()
	}
	return u.repo.List(ctx, offset, limit)
}
