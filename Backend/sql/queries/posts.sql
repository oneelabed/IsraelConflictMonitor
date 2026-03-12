-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, description, published_at, url, feed_id) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetDiversePosts :many
WITH RankedPosts AS (
    SELECT 
        posts.*, 
        feeds.name AS feed_name, 
        feeds.icon_url AS feed_icon,
        ROW_NUMBER() OVER (PARTITION BY posts.feed_id ORDER BY posts.published_at DESC) as rank
    FROM posts
    JOIN feeds ON posts.feed_id = feeds.id
)
SELECT * FROM RankedPosts 
WHERE rank <= 10 
ORDER BY published_at DESC 
LIMIT $1;

-- name: GetPostsForUser :many
SELECT 
    posts.*, 
    feeds.name as feed_name, 
    feeds.icon_url as feed_icon
FROM posts
JOIN feeds ON posts.feed_id = feeds.id
JOIN feed_follows ON feed_follows.feed_id = feeds.id
WHERE feed_follows.user_id = $1
ORDER BY posts.published_at DESC;

-- name: SearchPostsForUser :many
SELECT 
    posts.*, 
    feeds.name AS feed_name, 
    feeds.icon_url AS feed_icon
FROM posts
JOIN feeds ON posts.feed_id = feeds.id
JOIN feed_follows ON feed_follows.feed_id = feeds.id
WHERE feed_follows.user_id = $1 
  AND (posts.title ILIKE '%' || $2 || '%' OR posts.description ILIKE '%' || $2 || '%')
ORDER BY posts.published_at DESC;

