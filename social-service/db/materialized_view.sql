CREATE MATERIALIZED VIEW news_feed_global AS
SELECT
  nf.id,
  nf.title,
  nf.content,
  nf.image,
  nf.link,
  nf.visibility,
  nf.parent_id,
  nf.post_id,
  nf.num_like,
  nf.num_comment,
  nf.num_share,
  nf.num_view,
  nf.is_reel,
  nf.created_at,
  nf.updated_at,
  nf.created_by,

  -- Aggregate medias
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

  -- Aggregate friend tag IDs
  COALESCE(
    JSON_AGG(
      DISTINCT ft.user_id
    ) FILTER (WHERE ft.user_id IS NOT NULL),
    '[]'
  ) AS friend_tag_ids

FROM news_feeds nf
LEFT JOIN news_feed_medias nfm ON nf.id = nfm.news_feed_id
LEFT JOIN friend_tag ft ON nf.id = ft.news_feed_id

GROUP BY nf.id
ORDER BY nf.created_at DESC;

CREATE INDEX idx_mview_news_feed_id ON news_feed_global(id);
CREATE INDEX idx_mview_news_feed_created_at ON news_feed_global(created_at DESC);
REFRESH MATERIALIZED VIEW CONCURRENTLY news_feed_global;
