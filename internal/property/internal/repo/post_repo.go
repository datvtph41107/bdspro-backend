package repo

import (
	"bdspro/internal/domain"
	"bdspro/internal/dto"
	"bdspro/internal/enums"
	_enum "common/domain/enum"
	"context"
	"time"

	"gorm.io/gorm"
)

type PostRepo interface {
	// crud2.IBaseRepo[domain.Post]
	GetByID(id uint64) (*domain.Post, error)
	GetPostItemByID(ctx context.Context, id uint64) (*domain.PostItem, error)
	GetByIDs(ctx context.Context, ids []uint64) ([]*domain.PostItem, error)
	GetByIDsWithStatus(ctx context.Context, ids []uint64) ([]*domain.PostItem, error)
	CreatePost(ctx context.Context, post *domain.Post) error
	GetPostByProductID(productID uint64) (*domain.Post, error)
	UpdatePost(c context.Context, post *domain.Post) error
	CreateOrUpdatePostStatus(ctx context.Context, productID uint64, status enums.EPostStatus) error
	CreateOrUpdatePostVisibility(ctx context.Context, productID uint64, visibility enums.EVisibility) error
	UpdatePostStatus(ctx context.Context, productID uint64, status enums.EPostStatus) error
	CountPost(ctx context.Context, profileId uint64) int64
	CountPostFromDate(ctx context.Context, from time.Time) int64
	ExistedPostWorking(ctx context.Context, productID uint64, profileId uint64, transactionType enums.TransactionType) bool
	Query(query *gorm.DB, dto *dto.PostSearchRequest) (tx *gorm.DB)
	Search(ctx context.Context, profileId uint64, dto *dto.PostSearchRequest) ([]domain.PostItem, int64, error)
	Global(ctx context.Context, dto dto.PostSearchRequest) ([]domain.Post, int64, error)
	SearchLinkToProduct(ctx context.Context, profileId uint64, dto *dto.PostLinkSearch) ([]domain.PostItem, error)
	Delete(ctx context.Context, profileId, id uint64) error
	GetPublishByProfileID(ctx context.Context, profileId uint64, dto dto.PostPublishSearch) ([]domain.Post, int64, error)
	OwnerPost(ctx context.Context, profileId uint64, postId uint64) (bool, error)

	GetPosts(ctx context.Context, dto *dto.PostSearchRequest) ([]domain.PostItem, int64, error)

	// Đếm tổng số post hiện tại (không bị xóa)
	CountCurrent(ctx context.Context) (int64, error)
	CountPostByTime(ctx context.Context, from time.Time, to time.Time) ([]domain.CountPostByTime, error)
	CountByOwner(ctx context.Context, ownerOf _enum.EOwnerOf, ownerId uint64) (uint32, error)

	// Lấy ngẫu nhiên n bản ghi post
	GetRandomPosts(ctx context.Context, limit int) ([]domain.PostItem, error)

	// Đếm số post theo productId
	CountPostByProductID(ctx context.Context, productID uint64) (uint32, error)
}
