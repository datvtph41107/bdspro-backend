package usecases

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/repo"
	"context"
)

type UPostMediaUsecase struct {
	PostMediaRepo repo.PostMediaRepo
}

func NewUPostMediaUsecase(postMediaRepo repo.PostMediaRepo) *UPostMediaUsecase {
	return &UPostMediaUsecase{
		PostMediaRepo: postMediaRepo,
	}
}
func (s *UPostMediaUsecase) GetMediaList(c context.Context, dto dto.PostMediaGetDTO) ([]domain.PostMedia, error) {
	return s.PostMediaRepo.GetMediaList(c, dto)
}
