package usecases

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/repo"
	"context"
)

type PostMediaUsecase struct {
	PostMediaRepo repo.PostMediaRepo
}

func NewUPostMediaUsecase(postMediaRepo repo.PostMediaRepo) *PostMediaUsecase {
	return &PostMediaUsecase{
		PostMediaRepo: postMediaRepo,
	}
}
func (s *PostMediaUsecase) GetMediaList(c context.Context, postId uint64, dto dto.PostMediaGetDTO) ([]domain.PostMedia, error) {
	dto.PostID = postId
	return s.PostMediaRepo.GetMediaList(c, dto)
}
