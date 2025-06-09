package postgre

import (
	_utils "common/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"social/internal/domain"
	"social/internal/dto"
	"time"

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

type NewsFeedQuery struct {
	domain.NewsFeed
	LikeAt            *time.Time      `gorm:"column:like_at" json:"likeAt"`
	DisLike           *bool           `gorm:"column:dis_like" json:"disLike"`
	NewsFeedMediaRaws json.RawMessage `gorm:"column:news_feed_medias" json:"newsFeedMedias"` // hoặc []Media nếu dùng pgx/json decoding
	FriendTagRaws     json.RawMessage `gorm:"column:friend_tag_ids" json:"friendTagIds"`     // hoặc []uint64
}

func (r *PostgreNewsFeed) GetPublicByUserID(ctx context.Context, userId uint64, req *dto.NewsFeedGlobalSearch) ([]*domain.NewsFeed, int32, error) {
	var newsFeeds []*NewsFeedQuery
	var total int64
	query, err := r.QuerySeach(ctx, req)
	query = query.Where("news_feeds.created_by = ?", userId)
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
		Select(`news_feeds.*, tl.created_at as like_at, tl.dis_like,
		
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

	results := make([]*domain.NewsFeed, 0)
	for _, newsFeed := range newsFeeds {
		results = append(results, &newsFeed.NewsFeed)
		results[len(results)-1].LikeAt = newsFeed.LikeAt
		results[len(results)-1].DisLike = newsFeed.DisLike
		newsFeed.NewsFeed.LikeAt = newsFeed.LikeAt
		newsFeed.NewsFeed.DisLike = newsFeed.DisLike
		newsFeed.FriendTagIds = make([]uint64, len(newsFeed.FriendTag))
		var medias []domain.NewsFeedMedia
		var friendIDs []uint64

		if err := json.Unmarshal(newsFeed.NewsFeedMediaRaws, &medias); err != nil {
			log.Printf("failed to unmarshal NewsFeedMedias: %v", err)
			medias = []domain.NewsFeedMedia{} // fallback
		}

		if err := json.Unmarshal(newsFeed.FriendTagRaws, &friendIDs); err != nil {
			log.Printf("failed to unmarshal FriendTagIds: %v", err)
			friendIDs = []uint64{}
		}

		newsFeed.FriendTagIds = friendIDs
		newsFeed.NewsFeedMedias = medias

		for i, friendTag := range newsFeed.FriendTag {
			newsFeed.FriendTagIds[i] = friendTag.UserID
		}
	}

	return results, int32(total), nil
}

func (r *PostgreNewsFeed) GetGlobal(ctx context.Context, req *dto.NewsFeedGlobalSearch) ([]*dto.NewsFeedPublic, int32, error) {
	var newsFeeds []*dto.NewsFeedPublic
	var total int64
	_query := r.db.WithContext(ctx).
		Table("news_feed_global")

	if req.IsReel != nil {
		_query = _query.Where("is_reel = ?", req.IsReel)
	}

	err := _query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	whereCondition := ""
	if req.IsReel != nil {
		whereCondition = fmt.Sprintf("WHERE (is_reel = %t)", *req.IsReel)
	}
	query := `SELECT * FROM news_feed_global %s LIMIT ? OFFSET ?;`
	profileId := _utils.GetProfileIdWithContext(ctx)
	if profileId != 0 {
		query = fmt.Sprintf(`select g.*, tl.created_at as like_at,tl.dis_like from public.news_feed_global g
left join tb_like tl on tl.user_id = %d and tl.target_id  = g.id and tl.target_type = 20 %s LIMIT ? OFFSET ?;`, profileId, whereCondition)
	} else {
		query = fmt.Sprintf(query, whereCondition)
	}

	err = r.db.WithContext(ctx).
		Raw(query, req.GetLimit(), req.GetOffset()).
		Scan(&newsFeeds).Error
	if err != nil {
		return nil, 0, err
	}

	return newsFeeds, int32(total), nil
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

func (r *PostgreNewsFeed) GetByID(ctx context.Context, id uint64) (*domain.NewsFeed, error) {
	var newsFeed *domain.NewsFeed
	err := r.db.WithContext(ctx).
		Model(&domain.NewsFeed{}).
		Preload("NewsFeedMedias").
		Preload("FriendTag").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&newsFeed).
		Error

	if err != nil {
		return nil, err
	}

	newsFeed.FriendTagIds = make([]uint64, len(newsFeed.FriendTag))
	for i, friendTag := range newsFeed.FriendTag {
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
