package usecase

import (
	"context"

	_crud "common/domain/crud"
	_dto "common/domain/dto"
	_err "common/domain/err"
	_errors "common/errors"
	"hub/internal/domain"
	_repo "hub/internal/repo"
)

type IUserGuideUsecase interface {
	_crud.IBaseUsecase[domain.UserGuideEntity]
	GetListWithFilter(ctx context.Context, title string, groupKey string, pagable _dto.IPagable) ([]*domain.UserGuideEntity, int64, *_err.ErrorDTO)
	GetGroups(ctx context.Context) ([]map[string]interface{}, *_err.ErrorDTO)
	GetSimpleList(ctx context.Context, text string, groupKey string, pagable _dto.IPagable) ([]*domain.UserGuideEntity, int64, *_err.ErrorDTO)
	GetByKey(ctx context.Context, key string) (*domain.UserGuideEntity, error)
	GetAllWithKey(ctx context.Context) ([]*domain.UserGuideEntity, error)
}

type UserGuideUsecase struct {
	_crud.BaseUsecase[domain.UserGuideEntity, _repo.IUserGuideRepo]
	stepRepo _repo.IUserGuideStepRepo
}

func NewUserGuideUsecase(repo _repo.IUserGuideRepo, stepRepo _repo.IUserGuideStepRepo) IUserGuideUsecase {
	return &UserGuideUsecase{
		BaseUsecase: _crud.BaseUsecase[domain.UserGuideEntity, _repo.IUserGuideRepo]{
			Repo: repo,
		},
		stepRepo: stepRepo,
	}
}

// GetListWithFilter retrieves user guides with filters and pagination
func (u *UserGuideUsecase) GetListWithFilter(ctx context.Context, title string, groupKey string, pagable _dto.IPagable) ([]*domain.UserGuideEntity, int64, *_err.ErrorDTO) {
	data, total, err := u.Repo.GetListWithFilter(ctx, title, groupKey, pagable)
	if err != nil {
		return nil, 0, &_err.ErrorDTO{
			Code:    500,
			Message: "Lỗi khi lấy danh sách user guide: " + err.Error(),
		}
	}

	return data, total, nil
}

// GetGroups retrieves all unique groups with count
func (u *UserGuideUsecase) GetGroups(ctx context.Context) ([]map[string]interface{}, *_err.ErrorDTO) {
	groups, err := u.Repo.GetGroups(ctx)
	if err != nil {
		return nil, &_err.ErrorDTO{
			Code:    500,
			Message: "Lỗi khi lấy danh sách nhóm: " + err.Error(),
		}
	}

	return groups, nil
}

// GetSimpleList retrieves user guides for simple API (only id, title, description) with text search
func (u *UserGuideUsecase) GetSimpleList(ctx context.Context, text string, groupKey string, pagable _dto.IPagable) ([]*domain.UserGuideEntity, int64, *_err.ErrorDTO) {
	data, total, err := u.Repo.GetSimpleListWithText(ctx, text, groupKey, pagable)
	if err != nil {
		return nil, 0, &_err.ErrorDTO{
			Code:    500,
			Message: "Lỗi khi lấy danh sách user guide đơn giản: " + err.Error(),
		}
	}

	return data, total, nil
}

// GetByKey retrieves user guide by key
func (u *UserGuideUsecase) GetByKey(ctx context.Context, key string) (*domain.UserGuideEntity, error) {
	data, err := u.Repo.GetByKey(ctx, key)
	if err != nil {
		return nil, _errors.ReturnError(404, "Không tìm thấy user guide với key: "+key)
	}
	return data, nil
}

// GetAllWithKey lấy toàn bộ user guide có key (cho API sync).
func (u *UserGuideUsecase) GetAllWithKey(ctx context.Context) ([]*domain.UserGuideEntity, error) {
	return u.Repo.GetAllWithKey(ctx)
}
