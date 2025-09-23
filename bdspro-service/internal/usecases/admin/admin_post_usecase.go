package admin_usecases

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	admin_repo "bdspro/internal/repo/admin"
	"common/case/crud"
	_usecase "common/domain/usecase"
	"context"
)

type AdminPostUsecase struct {
	crud.BaseUsecase[domain.Post, admin_repo.AdminPostRepo]
	CodeDataUsecase *_usecase.CodeDataUsecase
}

func NewAdminPostUsecase(repo admin_repo.AdminPostRepo, codeDataUsecase *_usecase.CodeDataUsecase) *AdminPostUsecase {
	return &AdminPostUsecase{
		crud.BaseUsecase[domain.Post, admin_repo.AdminPostRepo]{
			Repo: repo,
		},
		codeDataUsecase,
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

// GetPosts - Lấy danh sách posts cho admin
func (uc *AdminPostUsecase) GetPosts(ctx context.Context, dto *dto.AdminPostSearchRequest) ([]domain.Post, int64, error) {
	result, total, err := uc.Repo.GetAdminPosts(ctx, dto)
	if err != nil {
		return nil, 0, err
	}
	for i := range result {
		result[i].Code = uc.CodeDataUsecase.GetPostCode(ctx, result[i].ID)
		// result[i].PublishedAt = _utils.FormatTimeToString(result[i].PublishedAt)
	}
	return result, total, nil
}

func (uc *AdminPostUsecase) GetPostDetail(ctx context.Context, id uint64) (*domain.Post, error) {
	post, err := uc.Repo.GetPostDetail(ctx, id)
	if err != nil {
		return nil, err
	}

	// Generate post code
	post.Code = uc.CodeDataUsecase.GetPostCode(ctx, post.ID)

	return post, nil
}
