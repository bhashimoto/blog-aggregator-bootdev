-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetPostsByUser :many
SELECT p.*
FROM posts p
JOIN feeds f ON p.feed_id = f.id
JOIN feed_follow ff ON f.id = ff.feed_id
JOIN users u ON ff.user_id = u.id
WHERE u.ID = $1
ORDER BY p.published_at DESC
LIMIT $2;
