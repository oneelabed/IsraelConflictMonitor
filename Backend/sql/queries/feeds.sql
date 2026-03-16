-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, icon_url) 
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetFeedsForUser :many
SELECT 
    feeds.*,
    -- If there's an ID in the follows table, they follow it. 
    -- If it's NULL, they don't.
    CASE WHEN feed_follows.id IS NOT NULL THEN TRUE ELSE FALSE END as is_following
FROM feeds
LEFT JOIN feed_follows ON feeds.id = feed_follows.feed_id 
    AND feed_follows.user_id = $1
ORDER BY is_following ASC;

-- name: GetNextFeedsToFetch :many
SELECT * FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT $1;

-- name: MarkFeedAsFetched :one
UPDATE feeds
SET last_fetched_at = NOW(), updated_at = NOW() 
WHERE id = $1
RETURNING *;
