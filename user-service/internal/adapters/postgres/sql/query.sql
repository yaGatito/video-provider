-- name: CreateUser :one
INSERT INTO users (name, lastname, email, password_hash, password_salt, created_at, status, is_admin)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id;

-- name: GetUser :one
SELECT id, name, lastname, email, password_hash, password_salt, created_at, status, is_admin
FROM users
WHERE id = $1;
