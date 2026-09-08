package usecase

import (
	"context"

	_crud "common/domain/crud"
	_dto "common/domain/dto"
	_err "common/domain/err"
	"hub/internal/domain"
	_repo "hub/internal/repo"
)

type IFAQUsecase interface {
	_crud.IBaseUsecase[domain.FAQEntity]
	GetListWithFilter(ctx context.Context, question string, groupKey string, pagable _dto.IPagable) ([]*domain.FAQEntity, int64, *_err.ErrorDTO)
	GetSimpleList(ctx context.Context, text string, groupKey string, pagable _dto.IPagable) ([]*domain.FAQEntity, int64, *_err.ErrorDTO)
}

type FAQUsecase struct {
	_crud.BaseUsecase[domain.FAQEntity, _repo.IFAQRepo]
}

func NewFAQUsecase(repo _repo.IFAQRepo) IFAQUsecase {
	return &FAQUsecase{
		BaseUsecase: _crud.BaseUsecase[domain.FAQEntity, _repo.IFAQRepo]{
			Repo: repo,
		},
	}
}

// GetListWithFilter retrieves FAQs with filters and pagination
func (u *FAQUsecase) GetListWithFilter(ctx context.Context, question string, groupKey string, pagable _dto.IPagable) ([]*domain.FAQEntity, int64, *_err.ErrorDTO) {
	data, total, err := u.Repo.GetListWithFilter(ctx, question, groupKey, pagable)
	if err != nil {
		return nil, 0, &_err.ErrorDTO{
			Code:    500,
			Message: "Lỗi khi lấy danh sách FAQ: " + err.Error(),
		}
	}

	return data, total, nil
}

// GetSimpleList retrieves FAQs for simple API (only id, question, answer) with text search
func (u *FAQUsecase) GetSimpleList(ctx context.Context, text string, groupKey string, pagable _dto.IPagable) ([]*domain.FAQEntity, int64, *_err.ErrorDTO) {
	data, total, err := u.Repo.GetSimpleListWithText(ctx, text, groupKey, pagable)
	if err != nil {
		return nil, 0, &_err.ErrorDTO{
			Code:    500,
			Message: "Lỗi khi lấy danh sách FAQ đơn giản: " + err.Error(),
		}
	}

	return data, total, nil
}
