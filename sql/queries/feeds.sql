-- name: CreateFeed :one
INSERT INTO feeds (id, name, url, user_id, created_at, updated_at)
VALUES($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetAllFeeds :many
SELECT *
FROM feeds;

-- name: GetNextFeedsToFetch :many
SELECT *
FROM feeds
ORDER BY last_fetched_at asc NULLS FIRST
LIMIT $1;

-- name: MarkFeedFetched :exec
UPDATE feeds SET 
	last_fetched_at = CURRENT_TIMESTAMP,
	updated_at 	= CURRENT_TIMESTAMP
WHERE id = $1;

-- name: DeleteFeed :exec
DELETE FROM feeds
WHERE id = $1;
