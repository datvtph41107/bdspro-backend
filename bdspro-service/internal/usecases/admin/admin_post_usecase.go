package admin_usecases

import (
	"bdspro/internal/domain"
	admin_repo "bdspro/internal/repo/admin"
	"common/case/crud"
	"context"
)

type AdminPostUsecase struct {
	crud.BaseUsecase[domain.Post, admin_repo.AdminPostRepo]
}

func NewAdminPostUsecase(repo admin_repo.AdminPostRepo) *AdminPostUsecase {
	return &AdminPostUsecase{
		crud.BaseUsecase[domain.Post, admin_repo.AdminPostRepo]{
			Repo: repo,
		},
	}
}

func (uc *AdminPostUsecase) Approve(ctx context.Context, id uint64) error {

	return uc.Repo.Approve(ctx, id)
}

func (uc *AdminPostUsecase) Reject(ctx context.Context, id uint64) error {

	return uc.Repo.Reject(ctx, id)
}

func (uc *AdminPostUsecase) Archive(ctx context.Context, postId uint64, archived bool) error {
	return uc.Repo.Archive(ctx, postId, archived)
}

func (uc *AdminPostUsecase) Hide(ctx context.Context, id uint64) error {

	return uc.Repo.Hide(ctx, id)
}

func (uc *AdminPostUsecase) Unhide(ctx context.Context, id uint64) error {

	return uc.Repo.Unhide(ctx, id)
}
