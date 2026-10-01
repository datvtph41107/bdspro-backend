DROP INDEX IF EXISTS idx_saved_news_domain;
DROP INDEX IF EXISTS idx_saved_news_user;
ALTER TABLE user_saved_planning_news DROP CONSTRAINT IF EXISTS user_saved_planning_news_user_domain_key;
