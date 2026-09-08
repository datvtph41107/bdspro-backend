CREATE MATERIALIZED VIEW public.news_feed_global
TABLESPACE pg_default
AS SELECT nf.id,
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
    nf.rank_feed,
    COALESCE(json_agg(DISTINCT jsonb_build_object('news_feed_id', nfm.news_feed_id, 'url', nfm.url, 'type', nfm.type, 'order_number', nfm.order_number)) FILTER (WHERE nfm.news_feed_id IS NOT NULL), '[]'::json) AS news_feed_medias,
    COALESCE(json_agg(DISTINCT ft.user_id) FILTER (WHERE ft.user_id IS NOT NULL), '[]'::json) AS friend_tag_ids
   FROM news_feeds nf
     LEFT JOIN news_feed_medias nfm ON nf.id = nfm.news_feed_id and nfm.deleted_at is null
     LEFT JOIN friend_tag ft ON nf.id = ft.news_feed_id and ft.deleted_at is null
  GROUP BY nf.id
  ORDER BY nf.rank_feed ASC
WITH DATA;

-- View indexes:
CREATE INDEX idx_mview_news_feed_created_at ON public.news_feed_global USING btree (created_at DESC);
CREATE INDEX idx_mview_news_feed_id ON public.news_feed_global USING btree (id);