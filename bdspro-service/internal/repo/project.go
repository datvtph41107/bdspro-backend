package repo

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"context"
)

type ProjectRepo interface {
	// common.IBaseRepo[domain.Project, *dto.ProjectSearchDTO]
	SearchItem(c context.Context, dto *dto.ProjectSearchDTO) ([]domain.ProjectItem, int64, error)
}
