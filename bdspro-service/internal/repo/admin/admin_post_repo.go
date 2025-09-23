package admin_repo

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"common/case/crud"
	"context"
)

type AdminPostRepo interface {
	crud.ICrudRepo[domain.Post]
	Approve(ctx context.Context, id uint64) error
	Reject(ctx context.Context, id uint64) error
	Archive(ctx context.Context, postId uint64, archived bool) error
	Hide(ctx context.Context, id uint64) error
	Unhide(ctx context.Context, id uint64) error
	GetAdminPosts(ctx context.Context, req *dto.AdminPostSearchRequest) ([]domain.Post, int64, error)
	GetPostDetail(ctx context.Context, id uint64) (*domain.Post, error)
}
