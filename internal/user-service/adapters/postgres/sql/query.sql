-- name: CreateUser :one
INSERT INTO users (name, lastname, email, password_hash, password_salt, created_at, status, is_admin)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id;

-- name: UpdateUser :exec
UPDATE users
SET name = $2, lastname = $3, email = $4, password_hash = $5, password_salt = $6, status = $7, is_admin = $8
WHERE id = $1;

-- name: GetUserById :one
SELECT id, name, lastname, email, password_hash, password_salt, created_at, status, is_admin
FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT id, name, lastname, email, password_hash, password_salt, created_at, status, is_admin
FROM users
WHERE email = $1;
