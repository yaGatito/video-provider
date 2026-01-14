-- name: CreateVideo :exec
INSERT INTO videos (
  publisherid, topic, description, createdAt
) VALUES (
  $1, $2, $3, NOW()
);
-- RETURNING *;

-- name: GetVideoByID :one
SELECT * FROM videos WHERE id = $1 LIMIT 1;

-- name: GetVideosByPublisher :many
SELECT * FROM videos WHERE publisherid = $1 ORDER BY createdAt LIMIT $2 OFFSET $3;

-- name: SearchGlobal :many
SELECT * FROM videos WHERE topic LIKE CONCAT('%', $1, '%') OR description LIKE CONCAT('%', $1, '%') ORDER BY createdAt LIMIT $2 OFFSET $3;

-- name: SearchPublisher :many
SELECT * FROM videos WHERE publisherid = $1 AND (topic LIKE CONCAT('%', $2, '%') OR description LIKE CONCAT('%', $2, '%')) ORDER BY createdAt LIMIT $3 OFFSET $4;
