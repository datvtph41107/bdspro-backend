package repo

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"context"
)

type PostMediaRepo interface {
	UpdateMediaItems(c context.Context, postID uint64, items []dto.PostMediaItem) error
	GetMediaList(c context.Context, dto dto.PostMediaGetDTO) ([]domain.PostMedia, error)
}
