package qh_usecase

import (
	_err "common/domain/err"
	"context"
	qh "userv4/modules/tqd/entity/qh"
	qh_dto "userv4/modules/tqd/dto/qh"
	repo "userv4/modules/tqd/repository"
)

type IProAIJobUsecase interface {
	GetByID(ctx context.Context, id uint64) (*qh.ProAIJob, error)
	GetList(ctx context.Context, req *qh_dto.ListProAIJobsRequest) ([]qh.ProAIJob, int64, error)
}

type proAIJobUsecase struct {
	repo repo.IProAIJobRepo
}

// NewProAIJobUsecase @bind: internal/usecase/qh.IProAIJobUsecase
func NewProAIJobUsecase(repo repo.IProAIJobRepo) IProAIJobUsecase {
	return &proAIJobUsecase{repo: repo}
}

func (u *proAIJobUsecase) GetByID(ctx context.Context, id uint64) (*qh.ProAIJob, error) {
	entity, err := u.repo.GetByID(ctx, id)
	if err != nil {
		return nil, &_err.ErrorDTO{Code: 404, Message: "Không tìm thấy job AI"}
	}
	return entity, nil
}

func (u *proAIJobUsecase) GetList(ctx context.Context, req *qh_dto.ListProAIJobsRequest) ([]qh.ProAIJob, int64, error) {
	entities, total, err := u.repo.GetList(ctx, req)
	if err != nil {
		return nil, 0, &_err.ErrorDTO{Code: 500, Message: err.Error()}
	}
	return entities, total, nil
}
