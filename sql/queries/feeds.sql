-- name: CreateFeed :one
INSERT INTO feeds (id, user_id, created_at, updated_at, name, url)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: SelectAllFeeds :many
SELECT * FROM feeds;

-- name: GetFeedsToFetch :many
SELECT * 
FROM feeds
ORDER BY last_fetched_at DESC NULLS FIRST
LIMIT $1;

-- name: MarkFeedAsFetched :exec
UPDATE feeds
SET last_fetched_at = $1,
    updated_at = $2
WHERE id = $3;
