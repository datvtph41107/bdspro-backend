package postgre

import (
	_utils "common/utils"
	"context"
	"fmt"
	"social/internal/domain"
	"social/internal/dto"
	"time"

	"gorm.io/gorm"
)

// @bind: social/internal/repo.CommentRepo
type PostgreComment struct {
	db *gorm.DB
}

func NewPostgreComment(db *gorm.DB) *PostgreComment {
	return &PostgreComment{db: db}
}

func (r *PostgreComment) Create(ctx context.Context, comment *domain.Comment) (*domain.Comment, error) {
	err := r.db.WithContext(ctx).Create(comment).Error
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (r *PostgreComment) Update(ctx context.Context, comment *domain.Comment) (*domain.Comment, error) {
	err := r.db.WithContext(ctx).
		Model(&domain.Comment{}).
		Where("id = ? AND deleted_at IS NULL", comment.ID).
		Updates(map[string]interface{}{
			"content": comment.Content,
		}).Error
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (r *PostgreComment) Delete(ctx context.Context, comment *domain.Comment) error {
	err := r.db.WithContext(ctx).
		Model(&domain.Comment{}).
		Where("id = ?", comment.ID).
		Update("deleted_at", time.Now()).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgreComment) SearchPublic(ctx context.Context, dto *dto.CommentRequest) ([]domain.CommentQuery, int64, error) {
	var comments []domain.CommentQuery
	var total int64

	// _query := r.db.WithContext(ctx).
	// 	Table("comment")

	query := ""
	profileCond := "from comment c"
	preQuery := "select c.* "
	parentCond := "and c.parent_id is null"
	profileId := _utils.GetProfileIdWithContext(ctx)
	if profileId != 0 {
		preQuery = "select c.*, tl.created_at as liked_at,tl.dis_like "
		profileCond = fmt.Sprintf(`from public.comment c
left join tb_like tl on tl.user_id = %d and tl.target_id = c.id and tl.target_type = 10`, profileId)
	}
	if dto.ParentID != nil {
		parentCond = fmt.Sprintf("and c.parent_id = %d", *dto.ParentID)
	}
	query = fmt.Sprintf(`select count(*) %s
where c.deleted_at is null and c.news_feed_id = ? %s
 group by c.id
order by c.created_at asc
LIMIT ? OFFSET ?;`, profileCond, parentCond)

	err := r.db.WithContext(ctx).
		Debug().
		Raw(query, dto.NewsFeedID, dto.GetLimit(), dto.GetOffset()).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	query = fmt.Sprintf(`%s %s
where c.deleted_at is null and c.news_feed_id = ? %s
order by c.created_at asc
LIMIT ? OFFSET ?;`, preQuery, profileCond, parentCond)

	err = r.db.WithContext(ctx).
		Debug().
		Raw(query, dto.NewsFeedID, dto.GetLimit(), dto.GetOffset()).
		Scan(&comments).Error
	if err != nil {
		return nil, 0, err
	}
	return comments, total, nil
}

func (r *PostgreComment) GetList(ctx context.Context, dto *dto.CommentRequest) ([]domain.CommentQuery, int64, error) {
	var comments []domain.CommentQuery
	var total int64
	query := r.db.WithContext(ctx).
		Model(&domain.CommentQuery{}).
		Where("deleted_at IS NULL and news_feed_id = ?", dto.NewsFeedID)
	if dto.ParentID != nil {
		query = query.Where("parent_id = ?", dto.ParentID)
	} else {
		query = query.Where("parent_id IS NULL")
	}

	err := query.
		Order("created_at DESC").
		Offset(dto.GetOffset()).
		Limit(dto.GetLimit()).
		Find(&comments).Error
	if err != nil {
		return nil, 0, err
	}

	query.Count(&total)
	return comments, total, nil
}

func (r *PostgreComment) GetByNewsFeedID(ctx context.Context, newsFeedID uint64) ([]*domain.Comment, error) {
	return nil, nil
}

func (r *PostgreComment) GetByParentID(ctx context.Context, parentID uint64) ([]*domain.Comment, error) {
	return nil, nil
}

func (r *PostgreComment) GetByNewsFeedIDAndUserID(ctx context.Context, newsFeedID uint64, userID uint64) ([]*domain.Comment, error) {
	return nil, nil
}

func (r *PostgreComment) ExistByID(ctx context.Context, id uint64) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).
		Model(&domain.Comment{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Count(&count).
		Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *PostgreComment) GetByID(ctx context.Context, id uint64) (*domain.Comment, error) {
	var comment domain.Comment
	if err := r.db.WithContext(ctx).
		Model(&domain.Comment{}).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&comment).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *PostgreComment) GetLastUpdated(ctx context.Context, id uint64) (*time.Time, error) {
	var updatedTime time.Time
	err := r.db.WithContext(ctx).
		Model(&domain.Comment{}).
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

func (r *PostgreComment) UpdateLikeNumber(ctx context.Context, newFeedID uint64, totalLike int, dislike bool) error {
	field := "num_like"
	if dislike {
		field = "num_dis_like"
	}
	return r.db.WithContext(ctx).
		Model(&domain.Comment{}).
		Where("id = ? and deleted_at is null", newFeedID).
		Update(field, totalLike).Error
}

func (r *PostgreComment) UpdateNumberReply(ctx context.Context, commentID uint64, totalReply int) error {
	return r.db.WithContext(ctx).
		Model(&domain.Comment{}).
		Where("id = ? and deleted_at is null", commentID).
		Update("num_reply", totalReply).Error
}

func (r *PostgreComment) CountReply(ctx context.Context, commentID uint64) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.Comment{}).
		Where("parent_id = ? and deleted_at is null", commentID).
		Count(&count).Error
	return int(count), err
}

func (r *PostgreComment) CountByNewsFeedID(ctx context.Context, newsFeedID uint64) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.Comment{}).
		Where("news_feed_id = ? and deleted_at is null", newsFeedID).
		Count(&count).Error
	return int(count), err
}

func (r *PostgreComment) GetNearestUserIds(ctx context.Context, newsFeedId uint64, parentId *uint64) ([]uint64, error) {
	var ids []uint64
	whereCond := ""
	if parentId != nil {
		whereCond = fmt.Sprintf("AND parent_id = %d", *parentId)
	}
	err := r.db.WithContext(ctx).Debug().Raw(
		fmt.Sprintf(`
		select user_id 
		from (
			SELECT  c.user_id, MAX(created_at) as latest
				FROM comment c
				WHERE news_feed_id = ? AND deleted_at IS NULL %s
				GROUP BY c.user_id
				ORDER BY latest DESC
			)
		Limit 3`, whereCond),
		newsFeedId).Scan(&ids).Error

	return ids, err
}

func (r *PostgreComment) SumUserComment(ctx context.Context, newsFeedId uint64, parentId *uint64) (int, error) {
	var count int
	whereCond := ""
	if parentId != nil {
		whereCond = fmt.Sprintf("AND parent_id = %d", *parentId)
	}
	err := r.db.WithContext(ctx).Raw(
		fmt.Sprintf(`
		select count(*) from (SELECT  c.user_id, MAX(created_at) as latest
		FROM comment c
		WHERE news_feed_id = ? AND deleted_at IS NULL %s
		group by c.user_id
		)
	`, whereCond), newsFeedId).Scan(&count).Error

	return count, err
}
