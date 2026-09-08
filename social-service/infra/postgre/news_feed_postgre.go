package postgre

import (
	_dto "common/domain/dto"
	_utils "common/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
	socialpb "pb/types/social"
	"time"

	"social/internal/domain"
	"social/internal/dto"
	"social/internal/enums"

	"gorm.io/gorm"
)

// @bind: social/internal/repo.NewsFeedRepo
type PostgreNewsFeed struct {
	db *gorm.DB
}

func NewPostgreNewsFeed(db *gorm.DB) *PostgreNewsFeed {
	return &PostgreNewsFeed{db: db}
}

func (r *PostgreNewsFeed) Create(ctx context.Context, post *domain.NewsFeed) (*domain.NewsFeed, error) {
	err := GetDB(ctx, r.db).Create(post).Error
	return post, err
}

func (r *PostgreNewsFeed) Update(ctx context.Context, post *domain.NewsFeed) (*domain.NewsFeed, error) {
	err := r.db.WithContext(ctx).
		Model(&domain.NewsFeed{}).
		Where("id = ?", post.ID).
		Update("title", post.Title).
		Update("content", post.Content).
		// Update("description", post.Description).
		Update("image", post.Image).
		Update("link", post.Link).
		Update("visibility", post.Visibility).Error
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (r *PostgreNewsFeed) Delete(ctx context.Context, post *domain.NewsFeed) error {
	err := r.db.WithContext(ctx).
		Model(&domain.NewsFeed{}).
		Where("id = ?", post.ID).
		Update("deleted_at", time.Now()).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgreNewsFeed) QuerySeach(ctx context.Context, dto *dto.NewsFeedGlobalSearch) (*gorm.DB, error) {
	query := r.db.WithContext(ctx).
		Table("news_feeds").
		Where("news_feeds.deleted_at IS NULL and news_feeds.visibility = 10")

	if dto.Text != "" {
		query = query.Where("title LIKE ?", "%"+dto.Text+"%")
	}

	return query, nil
}

type NewsFeedRawQuery struct {
	dto.NewsFeedPublic
	// LikedAt           *time.Time      `gorm:"column:liked_at" json:"likedAt"`
	// DisLike           *bool           `gorm:"column:dis_like" json:"disLike"`
	NewsFeedMediaRaws json.RawMessage `gorm:"column:news_feed_medias" json:"newsFeedMedias"` // hoặc []Media nếu dùng pgx/json decoding
	FriendTagRaws     json.RawMessage `gorm:"column:friend_tag_ids" json:"friendTagIds"`     // hoặc []uint64
}

func (r *PostgreNewsFeed) GetPublicByUserID(ctx context.Context, userId uint64, req *dto.NewsFeedGlobalSearch) ([]*dto.NewsFeedPublic, int32, error) {
	var newsFeeds []*NewsFeedRawQuery
	var total int64
	query, err := r.QuerySeach(ctx, req)
	query = query.Where("news_feeds.created_by = ? and news_feeds.removed_at is null and news_feeds.visibility = 10", userId)
	if err != nil {
		return nil, 0, err
	}
	err = query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.WithContext(ctx).
		Debug().
		Model(&domain.NewsFeed{}).
		Joins("left join tb_like tl on tl.target_id = news_feeds.id and tl.target_type = 20 and tl.user_id = ?", userId).
		Joins("left join news_feed_medias nfm on nfm.news_feed_id = news_feeds.id").
		Joins("left join friend_tag ft on ft.news_feed_id = news_feeds.id").
		Select(`news_feeds.*, tl.created_at as liked_at, tl.dis_like,
		
  COALESCE(
    JSON_AGG(
      DISTINCT JSONB_BUILD_OBJECT(
        'news_feed_id', nfm.news_feed_id,
        'url', nfm.url,
        'type', nfm.type,
        'order_number', nfm.order_number
      )
    ) FILTER (WHERE nfm.news_feed_id IS NOT NULL),
    '[]'
  ) AS news_feed_medias,

  COALESCE(
    JSON_AGG(
      DISTINCT ft.user_id
    ) FILTER (WHERE ft.user_id IS NOT NULL),
    '[]'
  ) AS friend_tag_ids
		`).
		Order("news_feeds.created_at DESC").
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Group("news_feeds.id, tl.created_at, tl.dis_like").
		Where("news_feeds.deleted_at is null and news_feeds.visibility = 10 and news_feeds.created_by = ?", userId).
		Scan(&newsFeeds).Error
	if err != nil {
		return nil, 0, err
	}

	results, err := r.buildNewsFeedResults(newsFeeds)
	if err != nil {
		return nil, 0, err
	}

	return results, int32(total), nil
}

func (r *PostgreNewsFeed) GetGlobal(ctx context.Context, req *dto.NewsFeedGlobalSearch) ([]*dto.NewsFeedPublic, int32, error) {
	var newsFeeds []*NewsFeedRawQuery
	var total int64

	profileId := _utils.GetProfileIdWithContext(ctx)

	// Query để đếm total
	countQuery := r.db.WithContext(ctx).
		Model(&domain.NewsFeed{}).
		Where("news_feeds.deleted_at IS NULL AND news_feeds.visibility = 10 AND news_feeds.removed_at IS NULL")

	if req.IsReel != nil {
		countQuery = countQuery.Where("news_feeds.is_reel = ?", *req.IsReel)
	}

	if req.IsPost != nil {
		if *req.IsPost {
			countQuery = countQuery.Where("news_feeds.post_id IS NOT NULL")
		} else {
			countQuery = countQuery.Where("news_feeds.post_id IS NULL")
		}
	}

	err := countQuery.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	// Query chi tiết với join
	detailQuery := r.db.WithContext(ctx).
		Debug().
		Model(&domain.NewsFeed{}).
		Joins("left join tb_like tl on tl.target_id = news_feeds.id and tl.target_type = 20 and tl.user_id = ?", profileId).
		Joins("left join news_feed_medias nfm on nfm.news_feed_id = news_feeds.id").
		Joins("left join friend_tag ft on ft.news_feed_id = news_feeds.id").
		Where("news_feeds.deleted_at IS NULL AND news_feeds.visibility = 10 AND news_feeds.removed_at IS NULL")

	if req.IsReel != nil {
		detailQuery = detailQuery.Where("news_feeds.is_reel = ?", *req.IsReel)
	}

	if req.IsPost != nil {
		if *req.IsPost {
			detailQuery = detailQuery.Where("news_feeds.post_id IS NOT NULL")
		} else {
			detailQuery = detailQuery.Where("news_feeds.post_id IS NULL")
		}
	}

	err = detailQuery.
		Select(`
            news_feeds.*, tl.created_at as liked_at, tl.dis_like,
            COALESCE(
              JSON_AGG(
                DISTINCT JSONB_BUILD_OBJECT(
                  'news_feed_id', nfm.news_feed_id,
                  'url', nfm.url,
                  'type', nfm.type,
                  'order_number', nfm.order_number
                )
              ) FILTER (WHERE nfm.news_feed_id IS NOT NULL),
              '[]'
            ) AS news_feed_medias,
            COALESCE(
              JSON_AGG(
                DISTINCT ft.user_id
              ) FILTER (WHERE ft.user_id IS NOT NULL),
              '[]'
            ) AS friend_tag_ids
        `).
		Order("news_feeds.created_at DESC").
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Group("news_feeds.id, tl.created_at, tl.dis_like").
		Scan(&newsFeeds).Error
	if err != nil {
		return nil, 0, err
	}

	results, err := r.buildNewsFeedResults(newsFeeds)
	if err != nil {
		return nil, 0, err
	}

	return results, int32(total), nil
}

func (r *PostgreNewsFeed) UpdateVisibility(ctx context.Context, newsFeedID uint64, visibility enums.Visibility) error {
	return r.db.WithContext(ctx).
		Model(&domain.NewsFeed{}).
		Where("id = ? and deleted_at is null", newsFeedID).
		Update("visibility", visibility).Error
}

func (r *PostgreNewsFeed) UpdateLikeNumber(ctx context.Context, newFeedID uint64, totalLike int) error {
	return r.db.WithContext(ctx).
		Model(&domain.NewsFeed{}).
		Where("id = ? and deleted_at is null", newFeedID).
		Update("num_like", totalLike).Error
}

func (r *PostgreNewsFeed) UpdateCommentNumber(ctx context.Context, newFeedID uint64, totalComment int) error {
	return r.db.WithContext(ctx).
		Model(&domain.NewsFeed{}).
		Where("id = ? and deleted_at is null", newFeedID).
		Update("num_comment", totalComment).Error
}

func (r *PostgreNewsFeed) UpdateComment(ctx context.Context, newFeedID uint64, numComment int) error {
	return nil
}

func (r *PostgreNewsFeed) UpdateShare(ctx context.Context, newFeedID uint64, numShare int) error {
	return nil
}

func (r *PostgreNewsFeed) UpdateView(ctx context.Context, newFeedID uint64, numView int) error {
	return nil
}

func (r *PostgreNewsFeed) UserAvaiableReact(ctx context.Context, profileID uint64, id uint64) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&domain.NewsFeed{}).
		Where("id = ? AND visibility = 10 AND deleted_at IS NULL", id).
		Count(&count).
		Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *PostgreNewsFeed) ExistByID(ctx context.Context, id uint64) (bool, error) {
	var exists bool
	query := `
        SELECT EXISTS (
            SELECT 1 FROM news_feeds
            WHERE id = $1 AND deleted_at IS NULL
        )
    `
	err := r.db.WithContext(ctx).
		Raw(query, id).
		Scan(&exists).
		Error
	return exists, err
}

func (r *PostgreNewsFeed) ExistByIDAndProfileID(ctx context.Context, id uint64, profileID uint64) (bool, error) {
	var exists bool
	query := `
        SELECT EXISTS (
            SELECT 1 FROM news_feeds
            WHERE id = $1 AND created_by = $2 AND deleted_at IS NULL
        )
    `
	err := r.db.WithContext(ctx).
		Raw(query, id, profileID).
		Scan(&exists).
		Error
	return exists, err
}

func (r *PostgreNewsFeed) GetPublicIgnorePreloadByID(ctx context.Context, id uint64) (*domain.NewsFeed, error) {
	var newsFeed domain.NewsFeed
	err := r.db.WithContext(ctx).
		Model(&domain.NewsFeed{}).
		Where("id = ? AND visibility = 10 AND deleted_at IS NULL", id).
		First(&newsFeed).
		Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &newsFeed, nil
}

func (r *PostgreNewsFeed) GetByID(ctx context.Context, id uint64) (*domain.NewsFeed, error) {
	var newsFeed *domain.NewsFeed
	err := r.db.WithContext(ctx).
		Model(&domain.NewsFeed{}).
		Preload("NewsFeedMedias").
		Preload("FriendTags").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&newsFeed).
		Error
	if err != nil {
		return nil, err
	}

	newsFeed.FriendTagIds = make([]uint64, len(newsFeed.FriendTags))
	for i, friendTag := range newsFeed.FriendTags {
		newsFeed.FriendTagIds[i] = friendTag.UserID
	}

	return newsFeed, err
}

// AdminGetByID lấy chi tiết bài viết cho admin (không filter deleted_at, visibility, removed_at)
func (r *PostgreNewsFeed) AdminGetByID(ctx context.Context, id uint64) (*domain.NewsFeed, error) {
	var newsFeed *domain.NewsFeed
	err := r.db.WithContext(ctx).
		Model(&domain.NewsFeed{}).
		Preload("NewsFeedMedias").
		Preload("FriendTags").
		Where("id = ?", id). // Admin có thể xem cả bài đã bị soft delete
		First(&newsFeed).
		Error
	if err != nil {
		return nil, err
	}

	newsFeed.FriendTagIds = make([]uint64, len(newsFeed.FriendTags))
	for i, friendTag := range newsFeed.FriendTags {
		newsFeed.FriendTagIds[i] = friendTag.UserID
	}

	return newsFeed, err
}

func (r *PostgreNewsFeed) GetLastUpdated(ctx context.Context, id uint64) (*time.Time, error) {
	var updatedTime time.Time
	err := r.db.WithContext(ctx).
		Model(&domain.NewsFeed{}).
		Select("updated_at").
		Where("id = ? AND deleted_at IS NULL", id).
		Order("updated_at DESC").
		Limit(1).
		Pluck("updated_at", &updatedTime).Error
	if err != nil {
		return nil, err
	}

	return &updatedTime, nil
}

func (r *PostgreNewsFeed) IncrementNumberShare(ctx context.Context, newsFeedID uint64) error {
	return r.db.WithContext(ctx).
		Transaction(func(tx *gorm.DB) error {
			err := tx.Model(&domain.NewsFeed{}).
				Where("id = ? AND deleted_at IS NULL", newsFeedID).
				UpdateColumn("num_share", gorm.Expr("COALESCE(num_share, 0) + ?", 1)).Error
			if err != nil {
				return err
			}
			return nil
		})
}

func (r *PostgreNewsFeed) SyncData(ctx context.Context) error {
	return r.db.Exec(`REFRESH MATERIALIZED VIEW public.news_feed_global`).Error
}

func (r *PostgreNewsFeed) GetReelTx(ctx context.Context, tx *gorm.DB) *gorm.DB {
	return GetDB(ctx, tx).
		Table("news_feeds").
		Select(`news_feeds.*, nfm.url as reel_url`).
		Joins("left join news_feed_medias nfm on nfm.news_feed_id = news_feeds.id").
		Where("news_feeds.deleted_at is null and news_feeds.is_reel = true and nfm.type = 'video'").
		Order("news_feeds.rank_feed desc").
		Group("news_feeds.id, nfm.id")
}

func (r *PostgreNewsFeed) GetReel(ctx context.Context, dto _dto.Pagable) ([]*domain.ReelEntity, int64, error) {
	var results []*domain.ReelEntity
	var count int64
	err := r.GetReelTx(ctx, r.db).
		Limit(dto.GetLimit()).
		Offset(dto.GetOffset()).
		Scan(&results).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.GetReelTx(ctx, r.db).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	return results, count, nil
}

// func (r *PostgreNewsFeed) GetListPostUser(ctx context.Context, req *dto.NewsFeedGlobalSearch, userId uint64) ([]*domain.NewsFeed, int32, error) {
// 	var results []*domain.NewsFeed
// 	var total int64
// 	err := r.db.WithContext(ctx).
// 		Model(&domain.NewsFeed{}).
// 		Where("created_by = ? and deleted_at is null", userId).
// 		Order("created_at DESC").
// 		Offset(req.GetOffset()).
// 		Limit(req.GetLimit()).
// 		Scan(&results).Error

func (r *PostgreNewsFeed) GetListPostUser(ctx context.Context,
	req *dto.NewsFeedGlobalSearch,
	userId uint64,
) ([]*dto.NewsFeedPublic, int32, error) {
	return r.GetListPostWithOwner(ctx, enums.OwnerOfUser, userId, req)
}

func (r *PostgreNewsFeed) GetListPostGroup(ctx context.Context,
	req *dto.NewsFeedGlobalSearch,
	groupId uint64,
) ([]*dto.NewsFeedPublic, int32, error) {
	return r.GetListPostWithOwner(ctx, enums.OwnerOfGroup, groupId, req)
}

func (r *PostgreNewsFeed) GetListPostOrganization(ctx context.Context,
	req *dto.NewsFeedGlobalSearch,
	organizationId uint64,
) ([]*dto.NewsFeedPublic, int32, error) {
	return r.GetListPostWithOwner(ctx, enums.OwnerOfOrganization, organizationId, req)
}

// Lấy news feed của User
func (r *PostgreNewsFeed) GetListPostWithOwner(
	ctx context.Context,
	ownerType enums.OwnerOf,
	ownerID uint64,
	req *dto.NewsFeedGlobalSearch,
) ([]*dto.NewsFeedPublic, int32, error) {
	var newsFeeds []*NewsFeedRawQuery
	var total int64

	// joinCondition := ""
	whereCondition := fmt.Sprintf("owner_id = %d and owner_of = %d", ownerID, ownerType)
	if req.IsGlobal {
		whereCondition = fmt.Sprintf("owner_id = %d and owner_of = %d and visibility = 10", ownerID, ownerType)
	}
	// if ownerType == enums.OwnerOfUser {
	// 	// tableName = "news_feed_of_user"
	// 	joinCondition = "JOIN news_feed_of_user nfu ON nfu.post_id = news_feeds.id"
	// 	whereCondition = fmt.Sprintf("nfu.user_id = %d", ownerID)
	// } else if ownerType == enums.OwnerOfGroup {
	// 	// tableName = "news_feed_of_group"
	// 	joinCondition = "JOIN news_feed_of_group nfg ON nfg.post_id = news_feeds.id"
	// 	whereCondition = fmt.Sprintf("nfg.group_id = %d", ownerID)
	// } else if ownerType == enums.OwnerOfOrganization {
	// 	// tableName = "news_feed_of_organization"
	// 	joinCondition = "JOIN news_feed_of_organization nfo ON nfo.post_id = news_feeds.id"
	// 	whereCondition = fmt.Sprintf("nfo.organization_id = %d", ownerID)
	// }

	// Join với bảng mapping news_feed_of_user
	query := GetDB(ctx, r.db).
		Debug().
		Model(&domain.NewsFeed{}).
		// Joins(joinCondition).
		Where(whereCondition)

	// đếm total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	profileId := _utils.GetProfileIdWithContext(ctx)

	// query chi tiết
	err := GetDB(ctx, r.db).
		Debug().
		Model(&domain.NewsFeed{}).
		// Joins(joinCondition).
		Joins("left join tb_like tl on tl.target_id = news_feeds.id and tl.target_type = 20 and tl.user_id = ?", profileId).
		Joins("left join news_feed_medias nfm on nfm.news_feed_id = news_feeds.id").
		Joins("left join friend_tag ft on ft.news_feed_id = news_feeds.id").
		Where(whereCondition).
		Select(`
            news_feeds.*, tl.created_at as liked_at, tl.dis_like,
            COALESCE(
              JSON_AGG(
                DISTINCT JSONB_BUILD_OBJECT(
                  'news_feed_id', nfm.news_feed_id,
                  'url', nfm.url,
                  'type', nfm.type,
                  'order_number', nfm.order_number
                )
              ) FILTER (WHERE nfm.news_feed_id IS NOT NULL),
              '[]'
            ) AS news_feed_medias,
            COALESCE(
              JSON_AGG(
                DISTINCT ft.user_id
              ) FILTER (WHERE ft.user_id IS NOT NULL),
              '[]'
            ) AS friend_tag_ids
        `).
		Order("news_feeds.created_at DESC").
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Group("news_feeds.id, tl.created_at, tl.dis_like").
		Scan(&newsFeeds).Error
	if err != nil {
		return nil, 0, err
	}

	results, err := r.buildNewsFeedResults(newsFeeds)
	if err != nil {
		return nil, 0, err
	}

	return results, int32(total), nil
}

// type NewsFeedGlobalRaw struct {
// 	NewsFeedPublic
// 	RawFriendTagIds   json.RawMessage `db:"friend_tag_ids" json:"-"`
// 	RawNewsFeedMedias json.RawMessage `db:"news_feed_medias" json:"-"`
// }

func (r *PostgreNewsFeed) buildNewsFeedResults(newsFeeds []*NewsFeedRawQuery) ([]*dto.NewsFeedPublic, error) {
	results := make([]*dto.NewsFeedPublic, 0)
	for _, newsFeed := range newsFeeds {
		// results[len(results)-1].LikedAt = newsFeed.LikedAt
		// results[len(results)-1].DisLike = newsFeed.DisLike
		// newsFeed.LikedAt = newsFeed.LikedAt
		// newsFeed.DisLike = newsFeed.DisLike
		// newsFeed.NewsFeedPublic.FriendTagIds = make([]uint64, len(newsFeed.RawFriendTagIds))
		var medias []*dto.NewsFeedMediaDTO
		var friendIDs []uint64

		if err := json.Unmarshal(newsFeed.NewsFeedMediaRaws, &medias); err != nil {
			log.Printf("failed to unmarshal NewsFeedMedias: %v", err)
			medias = []*dto.NewsFeedMediaDTO{} // fallback
		}

		if err := json.Unmarshal(newsFeed.FriendTagRaws, &friendIDs); err != nil {
			log.Printf("failed to unmarshal FriendTagIds: %v", err)
			friendIDs = []uint64{}
		}

		newsFeed.FriendTagIds = friendIDs
		newsFeed.NewsFeedMedias = medias

		for i, friendTag := range newsFeed.FriendTags {
			newsFeed.FriendTagIds[i] = friendTag.ID
		}
		// results[len(results)-1] = &newsFeed.NewsFeedPublic
		results = append(results, &newsFeed.NewsFeedPublic)
	}

	return results, nil
}

// CountByOwner đếm số tin đăng theo ownerOf và ownerId
func (r *PostgreNewsFeed) CountByOwner(ctx context.Context, ownerOf socialpb.OwnerOf, ownerId uint64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.NewsFeed{}).
		Where("owner_of = ? AND owner_id = ? AND deleted_at IS NULL", ownerOf, ownerId).
		Count(&count).Error
	return count, err
}

// HideNewsFeed ẩn/hiện bài viết (admin only)
func (r *PostgreNewsFeed) HideNewsFeed(ctx context.Context, newsFeedId uint64, isHidden bool) error {
	var removedAt *time.Time
	if isHidden {
		now := time.Now()
		removedAt = &now
	}
	return r.db.WithContext(ctx).
		Model(&domain.NewsFeed{}).
		Where("id = ? AND deleted_at IS NULL", newsFeedId).
		Update("removed_at", removedAt).Error
}

// AdminGetList lấy danh sách bài viết cho admin (không filter visibility và removed_at)
func (r *PostgreNewsFeed) AdminGetList(ctx context.Context, req *dto.AdminNewsFeedSearch) ([]*dto.NewsFeedPublic, int32, error) {
	var newsFeeds []*NewsFeedRawQuery
	var total int64

	// Base query - chỉ filter deleted_at (soft delete)
	query := GetDB(ctx, r.db).
		Debug().
		Model(&domain.NewsFeed{}).
		Where("deleted_at IS NULL")

	// Filter theo text
	if req.Text != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+req.Text+"%", "%"+req.Text+"%")
	}

	// Filter theo isReel
	if req.IsReel != nil {
		query = query.Where("is_reel = ?", *req.IsReel)
	}

	// Filter theo isPost
	if req.IsPost != nil {
		if *req.IsPost {
			query = query.Where("post_id IS NOT NULL")
		} else {
			query = query.Where("post_id IS NULL")
		}
	}

	// Filter theo ownerOf và ownerId
	if req.OwnerOf != nil && req.OwnerID != nil {
		query = query.Where("owner_of = ? AND owner_id = ?", *req.OwnerOf, *req.OwnerID)
	}

	// Filter theo visibility
	if req.Visibility != nil {
		query = query.Where("visibility = ?", *req.Visibility)
	}

	// Filter theo removed_at (isHidden)
	if req.IsHidden != nil {
		if *req.IsHidden {
			query = query.Where("removed_at IS NOT NULL")
		} else {
			query = query.Where("removed_at IS NULL")
		}
	}

	// Filter theo createdBy
	if req.CreatedBy != nil {
		query = query.Where("created_by = ?", *req.CreatedBy)
	}

	// Đếm total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	profileId := _utils.GetProfileIdWithContext(ctx)

	// Query chi tiết với join
	detailQuery := GetDB(ctx, r.db).
		Debug().
		Model(&domain.NewsFeed{}).
		Joins("left join tb_like tl on tl.target_id = news_feeds.id and tl.target_type = 20 and tl.user_id = ?", profileId).
		Joins("left join news_feed_medias nfm on nfm.news_feed_id = news_feeds.id").
		Joins("left join friend_tag ft on ft.news_feed_id = news_feeds.id").
		Where("news_feeds.deleted_at IS NULL")

	// Apply filters cho query chi tiết
	if req.Text != "" {
		detailQuery = detailQuery.Where("news_feeds.title LIKE ? OR news_feeds.content LIKE ?", "%"+req.Text+"%", "%"+req.Text+"%")
	}
	if req.IsReel != nil {
		detailQuery = detailQuery.Where("news_feeds.is_reel = ?", *req.IsReel)
	}
	if req.IsPost != nil {
		if *req.IsPost {
			detailQuery = detailQuery.Where("news_feeds.post_id IS NOT NULL")
		} else {
			detailQuery = detailQuery.Where("news_feeds.post_id IS NULL")
		}
	}
	if req.OwnerOf != nil && req.OwnerID != nil {
		detailQuery = detailQuery.Where("news_feeds.owner_of = ? AND news_feeds.owner_id = ?", *req.OwnerOf, *req.OwnerID)
	}
	if req.Visibility != nil {
		detailQuery = detailQuery.Where("news_feeds.visibility = ?", *req.Visibility)
	}
	if req.IsHidden != nil {
		if *req.IsHidden {
			detailQuery = detailQuery.Where("news_feeds.removed_at IS NOT NULL")
		} else {
			detailQuery = detailQuery.Where("news_feeds.removed_at IS NULL")
		}
	}
	if req.CreatedBy != nil {
		detailQuery = detailQuery.Where("news_feeds.created_by = ?", *req.CreatedBy)
	}

	err := detailQuery.
		Select(`
            news_feeds.*, tl.created_at as liked_at, tl.dis_like,
            COALESCE(
              JSON_AGG(
                DISTINCT JSONB_BUILD_OBJECT(
                  'news_feed_id', nfm.news_feed_id,
                  'url', nfm.url,
                  'type', nfm.type,
                  'order_number', nfm.order_number
                )
              ) FILTER (WHERE nfm.news_feed_id IS NOT NULL),
              '[]'
            ) AS news_feed_medias,
            COALESCE(
              JSON_AGG(
                DISTINCT ft.user_id
              ) FILTER (WHERE ft.user_id IS NOT NULL),
              '[]'
            ) AS friend_tag_ids
        `).
		Order("news_feeds.updated_at DESC").
		Offset(req.GetOffset()).
		Limit(req.GetLimit()).
		Group("news_feeds.id, tl.created_at, tl.dis_like").
		Scan(&newsFeeds).Error
	if err != nil {
		return nil, 0, err
	}

	results, err := r.buildNewsFeedResults(newsFeeds)
	if err != nil {
		return nil, 0, err
	}

	return results, int32(total), nil
}

// UpdateCreatedBy cập nhật createdBy cho bài viết (admin only)
func (r *PostgreNewsFeed) UpdateCreatedBy(ctx context.Context, newsFeedId uint64, createdBy uint64) error {
	return r.db.WithContext(ctx).
		Model(&domain.NewsFeed{}).
		Where("id = ?", newsFeedId).
		Update("created_by", createdBy).Error
}
